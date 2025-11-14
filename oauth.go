package garmin

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/dghubble/oauth1"

	"github.com/harrybrwn/go-garmin/internal/rt"
)

const (
	oAuthConsumerURL = "https://thegarth.s3.amazonaws.com/oauth_consumer.json"
	formContentType  = "application/x-www-form-urlencoded"
	authHeader       = "Authorization"
)

var (
	ErrNoCSRF              = errors.New("could not get csrf token")
	ErrNoTitle             = errors.New("could not find page title")
	ErrNoTicket            = errors.New("could not find ticket")
	ErrNotSuccessful       = errors.New("page title was not \"Success\"")
	ErrExpiredRefreshToken = errors.New("refresh token has expired")
	ErrAccessTokenExpired  = errors.New("access_token has expired")
)

type oAuthConsumer struct {
	Key    string `json:"consumer_key"`
	Secret string `json:"consumer_secret"`
}

var getOAuthConsumer = sync.OnceValues(func() (*oAuthConsumer, error) {
	res, err := http.Get(oAuthConsumerURL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var consumer oAuthConsumer
	err = json.NewDecoder(res.Body).Decode(&consumer)
	if err != nil {
		return nil, err
	}
	return &consumer, nil
})

func getOAuthConfig(c *Client, ticket string) (*oauth1.Config, error) {
	consumer, err := getOAuthConsumer()
	if err != nil {
		return nil, err
	}
	query := map[string]string{
		"login-url":          fmt.Sprintf("https://sso.%s/sso/embed", c.Domain),
		"accepts-mfa-tokens": "true",
	}
	if len(ticket) > 0 {
		query["ticket"] = ticket
	}
	return &oauth1.Config{
		ConsumerKey:    consumer.Key,
		ConsumerSecret: consumer.Secret,
		HTTPClient:     &c.http,
		Endpoint: oauth1.Endpoint{
			RequestTokenURL: new(URLBuilder).
				HTTPS().
				Host("connectapi", c.Domain).
				Path("/oauth-service/oauth/preauthorized").
				MapQuery(query).
				String(),
		},
	}, nil
}

type oauthClient struct {
	client       *Client
	signinParams url.Values
	csrf         string
	buf          bytes.Buffer
}

func login(client *Client, username, password string) (*oauth1.Token, *AccessToken, error) {
	var (
		err      error
		ssoEmbed = fmt.Sprintf("https://sso.%s/sso/embed", client.Domain)
		oc       = oauthClient{
			client: client,
			signinParams: url.Values{
				"id":                              []string{"gauth-widget"},
				"embedWidget":                     []string{"true"},
				"gauthHost":                       []string{ssoEmbed},
				"service":                         []string{ssoEmbed},
				"source":                          []string{ssoEmbed},
				"redirectAfterAccountLoginUrl":    []string{ssoEmbed},
				"redirectAfterAccountCreationUrl": []string{ssoEmbed},
			},
		}
	)
	if client.Cacher != nil {
		ot, at, err := getCachedPair(client.Cacher)
		if err == nil && !at.expired() {
			return ot, at, nil
		}
	}

	if err = oc.getCSRF(); err != nil {
		return nil, nil, err
	}
	ticket, err := oc.signin(username, password)
	if err != nil {
		return nil, nil, err
	}
	// Get tokens
	conf, err := getOAuthConfig(client, ticket)
	if err != nil {
		return nil, nil, err
	}
	// TODO This request can also yeild MFA info, I need to find a way to get it
	// without rewriting the oauth1 package.
	requestToken, requestSecret, err := conf.RequestToken()
	if err != nil {
		return nil, nil, err
	}
	token := oauth1.NewToken(requestToken, requestSecret)
	accessToken, err := exchange(client, conf, token)
	if err != nil {
		return nil, nil, err
	}
	if client.Cacher != nil {
		if err = client.Cacher.SaveOAuth1Token(token); err != nil {
			return token, accessToken, err
		}
		if err = client.Cacher.SaveAccessToken(accessToken); err != nil {
			return token, accessToken, err
		}
	}
	return token, accessToken, nil
}

func (oc *oauthClient) getCSRF() error {
	var (
		err error
		res *http.Response
	)
	// Get csrf token
	res, err = oc.client.Do(&http.Request{
		Method: "GET",
		URL:    oc.client.url("sso", "/sso/signin", oc.signinParams),
	})
	if err != nil {
		return err
	}
	_, err = io.Copy(&oc.buf, res.Body)
	if err = errors.Join(err, res.Body.Close()); err != nil {
		return err
	}
	oc.csrf, err = findCSRF(oc.buf.Bytes())
	return err
}

func (oc *oauthClient) signin(username, password string) (ticket string, err error) {
	var (
		res        *http.Response
		signinData = url.Values{
			"username": []string{username},
			"password": []string{password},
			"embed":    []string{"true"},
			"_csrf":    []string{oc.csrf},
		}
	)
	res, err = oc.client.Do(&http.Request{
		Method: "POST",
		URL:    oc.client.url("sso", "/sso/signin", oc.signinParams),
		Header: http.Header{
			"Content-Type": []string{"application/x-www-form-urlencoded"},
			"Referer":      []string{oc.client.prev.Request.URL.String()},
		},
		Body: io.NopCloser(strings.NewReader(signinData.Encode())),
	})
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	oc.buf.Reset()
	if _, err = io.Copy(&oc.buf, res.Body); err != nil {
		return "", err
	}
	title, err := findTitle(oc.buf.Bytes())
	if err != nil {
		return "", err
	}
	if strings.Contains(title, "MFA") {
		title, err = oc.handleMFA()
		if err != nil {
			return "", err
		}
	}
	if title != "Success" {
		return "", ErrNotSuccessful
	}
	return parseTicket(oc.buf.Bytes())
}

func (oc *oauthClient) handleMFA() (string, error) {
	if oc.client.MFAHandler == nil {
		return "", errors.New("no MFA handler specified, cannot get MFA code")
	}
	code, err := oc.client.MFAHandler()
	if err != nil {
		return "", err
	}
	data := url.Values{
		"mfa-code": []string{code},
		"_csrf":    []string{oc.csrf},
		"embed":    []string{"true"},
		"fromPage": []string{"setupEnterMfaCode"},
	}
	res, err := oc.client.Do(&http.Request{
		Method: "POST",
		URL: new(URLBuilder).
			HTTPS().
			Host("sso", oc.client.Domain).
			Path("/sso/verifyMFA/loginEnterMfaCode").
			Query(oc.signinParams).
			URL(),
		Body: io.NopCloser(strings.NewReader(data.Encode())),
		Header: http.Header{
			"Referer":      []string{oc.client.prev.Request.URL.String()},
			"Content-Type": []string{formContentType},
		},
	})
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	oc.buf.Reset()
	if _, err = io.Copy(&oc.buf, res.Body); err != nil {
		return "", err
	}
	return findTitle(oc.buf.Bytes())
}

type AccessToken struct {
	Scope string `json:"scope"`
	// JTI is a JWT ID.
	JTI                   string `json:"jti"`
	AccessToken           string `json:"access_token"`
	TokenType             string `json:"token_type"`
	RefreshToken          string `json:"refresh_token"`
	Expires               int64  `json:"expires"`
	ExpiresIn             int    `json:"expires_in"`
	RefreshTokenExpires   int64  `json:"refresh_token_expires"`
	RefreshTokenExpiresIn int    `json:"refresh_token_expires_in"`
}

func (at *AccessToken) expired() bool {
	return time.Now().After(at.ExpiresAt())
}

func (at *AccessToken) refreshExpired() bool {
	return time.Now().After(at.RefreshTokenExpiresAt())
}

func (at *AccessToken) setExpirations(now time.Time) {
	expires := now.Add(time.Second * time.Duration(at.ExpiresIn))
	refreshTokenExpires := now.Add(time.Second * time.Duration(at.RefreshTokenExpiresIn))
	at.Expires = expires.UnixMilli()
	at.RefreshTokenExpires = refreshTokenExpires.UnixMilli()
}

func (at *AccessToken) ExpiresAt() time.Time {
	return time.UnixMilli(at.Expires)
}

func (at *AccessToken) RefreshTokenExpiresAt() time.Time {
	return time.UnixMilli(at.RefreshTokenExpires)
}

func (at *AccessToken) marshal(r io.Reader, now time.Time) error {
	b, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, at)
	if err != nil {
		return err
	}
	at.setExpirations(now)
	return nil
}

func exchange(client *Client, conf *oauth1.Config, token *oauth1.Token) (*AccessToken, error) {
	body := url.Values{} // TODO add mfa info
	c := conf.Client(context.Background(), token)
	c.Transport.(*oauth1.Transport).Base = client.http.Transport
	req := http.Request{
		Method: "POST",
		URL: new(URLBuilder).
			HTTPS().
			Host("connectapi", client.Domain).
			Path("/oauth-service/oauth/exchange/user/2.0").
			URL(),
		Header: http.Header{"Content-Type": []string{formContentType}},
		Body:   io.NopCloser(strings.NewReader(body.Encode())),
	}
	res, err := c.Do(&req)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}
	var at AccessToken
	err = at.marshal(res.Body, now)
	if err != nil {
		return nil, err
	}
	return &at, nil
}

type Refresher interface {
	Refresh(*AccessToken) (*AccessToken, error)
}

type accessTokenInjector struct {
	AccessToken *AccessToken
	refresher   Refresher
	base        http.RoundTripper
}

func (ati *accessTokenInjector) Wrap(rt http.RoundTripper) rt.RoundTripper {
	ati.base = rt
	return ati
}

func (ati *accessTokenInjector) Unwrap() http.RoundTripper { return ati.base }

func (ati *accessTokenInjector) RoundTrip(req *http.Request) (*http.Response, error) {
	if ati.AccessToken.expired() {
		if ati.AccessToken.refreshExpired() {
			return nil, ErrExpiredRefreshToken
		}
		at, err := ati.refresher.Refresh(ati.AccessToken)
		if err != nil {
			return nil, err
		}
		ati.AccessToken = at
	}
	if req.Header == nil {
		req.Header = make(http.Header)
	}
	req.Header.Set(authHeader, fmt.Sprintf("%s %s", ati.AccessToken.TokenType, ati.AccessToken.AccessToken))
	return ati.base.RoundTrip(req)
}

type oauth1TokenRefresher struct {
	token  *oauth1.Token
	client *Client
}

func (otr *oauth1TokenRefresher) Refresh(*AccessToken) (*AccessToken, error) {
	conf, err := getOAuthConfig(otr.client, "")
	if err != nil {
		return nil, err
	}
	accessToken, err := exchange(otr.client, conf, otr.token)
	if err != nil {
		return nil, err
	}
	if otr.client.Cacher != nil {
		if err = otr.client.Cacher.SaveOAuth1Token(otr.token); err != nil {
			return accessToken, err
		}
		if err = otr.client.Cacher.SaveAccessToken(accessToken); err != nil {
			return accessToken, err
		}
	}
	return accessToken, nil
}

var (
	csrfRe   = regexp.MustCompile(`name="_csrf"\s+value="(.+?)"`)
	titleRe  = regexp.MustCompile(`<title>(.+?)</title>`)
	ticketRe = regexp.MustCompile(`embed\?ticket=([^"]+)`)
)

func findCSRF(b []byte) (string, error) {
	idx := csrfRe.FindSubmatchIndex(b)
	if len(idx) < 4 {
		return "", ErrNoCSRF
	}
	// index the first match group
	return string(b[idx[2]:idx[3]]), nil
}

func findTitle(b []byte) (string, error) {
	idx := titleRe.FindSubmatchIndex(b)
	if len(idx) < 4 {
		return "", ErrNoTitle
	}
	// index the first match group
	return string(b[idx[2]:idx[3]]), nil
}

func parseTicket(b []byte) (string, error) {
	idx := ticketRe.FindSubmatchIndex(b)
	if len(idx) < 4 {
		return "", ErrNoTitle
	}
	// index the first match group
	return string(b[idx[2]:idx[3]]), nil
}
