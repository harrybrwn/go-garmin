package garmin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/harrybrwn/go-garmin/internal/rt"
)

var (
	BaseDomain    = "garmin.com"
	DefaultClient = NewClient()
)

const UserAgent = "com.garmin.android.apps.connectmobile"

type Client struct {
	Domain     string
	Cacher     TokenCacher
	AddReferer bool
	MFAHandler func() (string, error)
	Clock      Clock

	http http.Client
	prev *http.Response
}

func NewClient(opts ...ClientOpt) *Client {
	options := clientOpts{
		Domain:    BaseDomain,
		UserAgent: UserAgent,
		Clock:     &DefaultClock,
		Transport: http.DefaultTransport,
	}
	for _, o := range opts {
		o(&options)
	}
	uat := rt.NewUserAgent(options.UserAgent)
	cookies, _ := cookiejar.New(options.CookieOpts)
	c := http.Client{
		Transport: uat.Wrap(options.Transport),
		Jar:       cookies,
	}
	client := Client{
		Domain:     options.Domain,
		Cacher:     options.Cacher,
		MFAHandler: options.MFAHandler,
		Clock:      options.Clock,
		http:       c,
	}
	return &client
}

type clientOpts struct {
	Transport  http.RoundTripper
	CookieOpts *cookiejar.Options
	UserAgent  string
	Domain     string
	Cacher     TokenCacher
	MFAHandler func() (string, error)
	Clock      Clock
}

type ClientOpt func(*clientOpts)

func WithCacher(cacher TokenCacher) ClientOpt {
	return func(c *clientOpts) {
		if cacher != nil {
			c.Cacher = cacher
		}
	}
}

func WithDomain(domain string) ClientOpt { return func(c *clientOpts) { c.Domain = domain } }

func WithMFAHandler(fn func() (string, error)) ClientOpt {
	return func(c *clientOpts) { c.MFAHandler = fn }
}

func WithUserAgent(ua string) ClientOpt {
	return func(c *clientOpts) {
		c.UserAgent = ua
	}
}

func WithTransport(transport rt.RoundTripper) ClientOpt {
	return func(c *clientOpts) {
		if c.Transport == nil {
			c.Transport = http.DefaultTransport
		}
		c.Transport = transport.Wrap(c.Transport)
	}
}

func WithCookieOpts(opts *cookiejar.Options) ClientOpt {
	return func(co *clientOpts) { co.CookieOpts = opts }
}

func WithClock(clock Clock) ClientOpt { return func(co *clientOpts) { co.Clock = clock } }

func WithDebugging(enabled, skipBody bool) ClientOpt {
	if !enabled {
		return func(co *clientOpts) {}
	}
	return WithTransport(&rt.Debugger{SkipBody: skipBody})
}

// Login will get an access token and auto authenticate every request sent by
// the client.
func (c *Client) Login(email, password string) error {
	basic, access, err := login(c, email, password)
	if err != nil {
		return err
	}
	refresher := oauth1TokenRefresher{
		token:  basic,
		client: c,
	}
	injector := accessTokenInjector{
		AccessToken: access,
		refresher:   &refresher,
	}
	c.prependTransport(&injector)
	return nil
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	if c.AddReferer && c.prev != nil {
		req.Header.Add("Referer", c.prev.Request.URL.String())
	}
	req.Host = req.URL.Host
	res, err := c.http.Do(req)
	if err != nil {
		return res, err
	}
	c.prev = res
	return res, nil
}

func (c *Client) url(sub, path string, params url.Values) *url.URL {
	u := url.URL{
		Scheme: "https",
		Path:   path,
	}
	if params != nil {
		u.RawQuery = params.Encode()
	}
	if len(sub) > 0 {
		u.Host = fmt.Sprintf("%s.%s", sub, c.Domain)
	} else {
		u.Host = c.Domain
	}
	return &u
}

func (c *Client) apiGet(out any, path string, params url.Values) error {
	host := fmt.Sprintf("connectapi.%s", c.Domain)
	req := http.Request{
		Method: "GET",
		Host:   host,
		URL: &url.URL{
			Scheme: "https",
			Host:   host,
			Path:   path,
		},
		Header: http.Header{
			"Accept": []string{"application/json"},
			"Nk":     []string{"NT"},
		},
	}
	if len(params) > 0 {
		req.URL.RawQuery = params.Encode()
	}
	res, err := c.Do(&req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("recived bad status code %q", res.Status)
	}
	return json.NewDecoder(res.Body).Decode(out)
}

func (c *Client) api(out any, method, path string, params url.Values, payload any) (int, error) {
	host := fmt.Sprintf("connectapi.%s", c.Domain)
	req := http.Request{
		Method: method,
		Host:   host,
		URL: &url.URL{
			Scheme: "https",
			Host:   host,
			Path:   path,
		},
		Header: http.Header{
			"Accept": []string{"application/json"},
		},
	}
	if len(params) > 0 {
		req.URL.RawQuery = params.Encode()
	}
	var (
		err  error
		body bytes.Buffer
	)
	if payload != nil {
		req.Header.Set("Content-Type", "application/json")
		err = json.NewEncoder(&body).Encode(payload)
		if err != nil {
			return 0, err
		}
		req.Body = io.NopCloser(&body)
	}
	res, err := c.Do(&req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()
	if out != nil {
		err = json.NewDecoder(res.Body).Decode(out)
	}
	return res.StatusCode, err
}

func (c *Client) prependTransport(rt rt.RoundTripper) {
	c.http.Transport = rt.Wrap(c.http.Transport)
}

type URLBuilder struct {
	u url.URL
}

func (ub *URLBuilder) URL() *url.URL {
	if len(ub.u.Scheme) == 0 {
		ub.Scheme("https")
	}
	return &ub.u
}

func (ub *URLBuilder) String() string { return ub.URL().String() }

func (ub *URLBuilder) Scheme(s string) *URLBuilder {
	ub.u.Scheme = s
	return ub
}

func (ub *URLBuilder) HTTPS() *URLBuilder {
	ub.Scheme("https")
	return ub
}

func (ub *URLBuilder) Host(h string, domain ...string) *URLBuilder {
	if len(domain) > 0 {
		ub.u.Host = strings.Join(append([]string{h}, domain...), ".")
	} else {
		ub.u.Host = h
	}
	return ub
}

func (ub *URLBuilder) Subdomain(sub, domain string) *URLBuilder {
	ub.u.Host = fmt.Sprintf("%s.%s", sub, domain)
	return ub
}

func (ub *URLBuilder) Path(p string) *URLBuilder {
	ub.u.Path = p
	return ub
}

func (ub *URLBuilder) Pathf(p string, v ...any) *URLBuilder {
	ub.u.Path = fmt.Sprintf(p, v...)
	return ub
}

func (ub *URLBuilder) Query(q url.Values) *URLBuilder {
	ub.u.RawQuery = q.Encode()
	return ub
}

func (ub *URLBuilder) MapQuery(q map[string]string) *URLBuilder {
	vals := make(url.Values, len(q))
	for k, v := range q {
		vals[k] = []string{v}
	}
	ub.u.RawQuery = vals.Encode()
	return ub
}

func (ub *URLBuilder) Fragment(f string) *URLBuilder {
	ub.u.Fragment = f
	return ub
}

func (ub *URLBuilder) Clear() *URLBuilder {
	ub.u.Scheme = ""
	ub.u.Opaque = ""
	ub.u.User = nil
	ub.u.Host = ""
	ub.u.Path = ""
	ub.u.RawPath = ""
	ub.u.RawQuery = ""
	ub.u.Fragment = ""
	ub.u.RawFragment = ""
	return ub
}
