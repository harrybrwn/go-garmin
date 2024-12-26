package garmin

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/harrybrwn/go-garmin/internal/rt"
)

func TestParsers(t *testing.T) {
	t.Run("FindCSRF", func(t *testing.T) {
		v, err := findCSRF([]byte(`
            <input type="hidden" name="embed" value="true"/>
            <input type="hidden" name="_csrf" value="__csrf-token-value__" />
            <button type="submit" id="login-btn-signin" class="btn1" accesskey="l">Sign In</button>
        `))
		if err != nil {
			t.Fatal(err)
		}
		exp := "__csrf-token-value__"
		if v != exp {
			t.Errorf("expected %q to be %q", v, exp)
		}
	})
	t.Run("FindTitle", func(t *testing.T) {
		v, err := findTitle([]byte(`<head>
            <title>Success</title>
            <meta charset="utf-8">
            <meta http-equiv="X-UA-Compatible" content="IE=edge;" />
            <meta name="description" content="">
            <meta name="viewport" content="width=device-width, initial-scale=1">
            <meta http-equiv="cleartype" content="on">`))
		if err != nil {
			t.Fatal(err)
		}
		exp := "Success"
		if v != exp {
			t.Errorf("expected %q to be %q", v, exp)
		}
	})
	t.Run("ParseTicket", func(t *testing.T) {
		v, err := parseTicket([]byte(`
            var consumeServiceTicket                  = "true";
            var service_url                           = "https:\/\/sso.garmin.com\/sso\/embed";
            var parent_url                            = "https:\/\/sso.garmin.com\/sso\/embed";
            var response_url                          = "https:\/\/sso.garmin.com\/sso\/embed?ticket=ST-ticket-cas";
            var logintoken                            = "";
            var socialLogin                           = "";
        `))
		if err != nil {
			t.Fatal(err)
		}
		exp := "ST-ticket-cas"
		if v != exp {
			t.Errorf("expected %q to be %q", v, exp)
		}
	})
}

func TestAccessToken(t *testing.T) {
	at := AccessToken{
		Scope:     "COMMUNITY_COURSE_READ GOLF_API_READ GHS_HID ATP_READ GHS_SAMD INSIGHTS_READ DIVE_API_READ DIVE_API_IMAGE_PREVIEW COMMUNITY_COURSE_WRITE CONNECT_WRITE DIVE_SHARED_READ GHS_REGISTRATION DT_CLIENT_ANALYTICS_WRITE GOLF_API_WRITE INSIGHTS_WRITE PRODUCT_SEARCH_READ GOLF_SHARED_READ OMT_CAMPAIGN_READ CONNECT_NON_SOCIAL_SHARED_READ CONNECT_READ ATP_WRITE",
		JTI:       "8fb56000-0b37-0000-0000-000000000000",
		TokenType: "Bearer",
		AccessToken: "" +
			b64strip(jwtHeader{
				Alg: "RS256",
				Typ: "JWT",
				Kid: "di-oauth-signer-prod-2024-q1",
			}, nil) +
			"." +
			b64strip(garminClaims{
				Scope: []string{
					"ATP_READ",
					"ATP_WRITE",
					"COMMUNITY_COURSE_READ",
					"COMMUNITY_COURSE_WRITE",
					"CONNECT_NON_SOCIAL_SHARED_READ",
					"CONNECT_READ",
					"CONNECT_WRITE",
					"DIVE_API_IMAGE_PREVIEW",
					"DIVE_API_READ",
					"DIVE_SHARED_READ",
					"DT_CLIENT_ANALYTICS_WRITE",
					"GHS_HID",
					"GHS_REGISTRATION",
					"GHS_SAMD",
					"GOLF_API_READ",
					"GOLF_API_WRITE",
					"GOLF_SHARED_READ",
					"INSIGHTS_READ",
					"INSIGHTS_WRITE",
					"OMT_CAMPAIGN_READ",
					"PRODUCT_SEARCH_READ",
				},
				Iss:                   "https://diauth.garmin.com",
				RevocationEligibility: []string{"GLOBAL_SIGNOUT"},
				ClientType:            "UNDEFINED",
				Exp:                   1723676069,
				Iat:                   1723675769,
				GarminGUID:            "44444444-5555-6666-cccc-dddddddddddd",
				Jti:                   "11111111-2222-3333-aaaa-bbbbbbbbbbbb",
				ClientID:              "CONNECT_WEB",
				FGP:                   "00000000000000000000000000000000000000000000000000000000000",
			}, nil) +
			"." +
			"D3-TvI2XvPUJQn9xA3Or8K4gYesaWfkDNccP-0000000_00000000000000000000000-p8IzKsM9_Eoby5mLfQOasD3WY-Xp",
		RefreshToken: b64(&refreshToken{
			RefreshTokenValue: "00000000-ae12-45e1-ab12-000000000000",
			GarminGuid:        "00000000-ab12-45e1-9f18-000000000000",
		}, base64.StdEncoding),
		ExpiresIn:             299,
		Expires:               1723676008817,
		RefreshTokenExpiresIn: 7199,
		RefreshTokenExpires:   1723682908817,
	}
	if !at.expired() {
		t.Error("access token should be marked as expired")
	}
	_ = at.refreshExpired() // TODO check this
}

func TestAccessTokenInjector(t *testing.T) {
	type TT struct {
		at  AccessToken
		exp string
		err error
	}
	for _, tt := range []TT{
		{
			at: AccessToken{
				Expires:             time.Now().Add(time.Hour * 24).UnixMilli(),
				RefreshTokenExpires: time.Now().Add(time.Hour * 48).UnixMilli(),
			},
			exp: "Bearer access-token0",
			err: nil,
		},
		{
			at: AccessToken{
				Expires:             time.Now().Add(-1 * time.Hour).UnixMilli(),
				RefreshTokenExpires: time.Now().Add(time.Hour * 48).UnixMilli(),
			},
			exp: "Bearer refreshed-access-token1",
			err: nil,
		},
		{
			at: AccessToken{
				RefreshTokenExpires: time.Now().Add(-1 * time.Hour * 48).UnixMilli(),
			},
			exp: "",
			err: ErrExpiredRefreshToken,
		},
	} {
		tt.at.AccessToken = "access-token0"
		tt.at.TokenType = "Bearer"
		ij := accessTokenInjector{
			AccessToken: &tt.at,
			refresher: refresherFunc(func(at *AccessToken) (*AccessToken, error) {
				var refed = *at
				refed.AccessToken = "refreshed-access-token1"
				return &refed, nil
			}),
			base: rt.RoundTripperFunc(func(r *http.Request) (*http.Response, error) {
				return &http.Response{Request: r}, nil
			}),
		}
		req := http.Request{Header: make(http.Header)}
		_, err := ij.RoundTrip(&req)
		if !errors.Is(err, tt.err) {
			t.Fatalf("want error %v, got error %v", tt.err, err)
		}
		if a := req.Header.Get(authHeader); a != "" && a != tt.exp {
			t.Errorf("authorization header: got %q, want %q", a, tt.exp)
		}
	}
}

func TestDates(t *testing.T) {
	t.Skip()
	const dateFormat = "2006-01-02T15:04:05.99"
	now := time.Now()
	fmt.Println(now.Format(dateFormat))
	fmt.Println(now.In(time.UTC).Format(dateFormat))
}

func loadEnv(name string) {
	f, err := os.Open(name)
	if os.IsNotExist(err) {
		return
	} else if err != nil {
		panic(err)
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	for s.Scan() {
		line := s.Text()
		parts := strings.SplitN(line, "=", 2)
		err = os.Setenv(parts[0], parts[1])
		if err != nil {
			panic(err)
		}
	}
}

// fields are in a fixed order
type jwtHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
	Kid string `json:"kid"`
}

// fields are in a fixed order
type garminClaims struct {
	Scope                 []string `json:"scope"`
	Iss                   string   `json:"iss"`
	RevocationEligibility []string `json:"revocation_eligibility"`
	ClientType            string   `json:"client_type"`
	Exp                   int      `json:"exp"`
	Iat                   int      `json:"iat"`
	GarminGUID            string   `json:"garmin_guid"`
	Jti                   string   `json:"jti"`
	ClientID              string   `json:"client_id"`
	FGP                   string   `json:"fgp"`
}

// fields are in a fixed order
type refreshToken struct {
	RefreshTokenValue string `json:"refreshTokenValue"`
	GarminGuid        string `json:"garminGuid"`
}

func b64(m any, enc *base64.Encoding) string {
	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	if enc == nil {
		enc = base64.StdEncoding
	}
	res := make([]byte, enc.EncodedLen(len(b)))
	base64.StdEncoding.Encode(res, b)
	return string(res)
}

func b64strip(m any, enc *base64.Encoding) string {
	return strings.TrimRight(b64(m, enc), "=")
}

type refresherFunc func(*AccessToken) (*AccessToken, error)

func (fn refresherFunc) Refresh(at *AccessToken) (*AccessToken, error) { return fn(at) }
