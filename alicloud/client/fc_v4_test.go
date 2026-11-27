package client

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/credentials-go/credentials"
	fc "github.com/aliyun/fc-go-sdk"
	"gopkg.in/resty.v1"
)

// newFcCaptureClient points the FC endpoint at a local HTTP capture server
// with an STS credential, covering the temporary-token configuration the EO
// environment uses.
func newFcCaptureClient(t *testing.T, endpoint, region string) *AliyunClient {
	t.Helper()
	c := newLocalRpcClient(t, endpoint)
	// Satisfies the fc account-id endpoint formatting without a caller
	// identity lookup.
	c.accountId = "1234567890"
	c.config.Endpoints.Store("fc", endpoint)
	c.config.RegionId, c.RegionId = region, region
	cred, err := credentials.NewCredential(new(credentials.Config).SetType("sts").
		SetAccessKeyId("test-ak").SetAccessKeySecret("test-secret").SetSecurityToken("test-token"))
	if err != nil {
		t.Fatal(err)
	}
	c.config.Credential = cred
	return c
}

// runFcCaptureCall drives one real fc-go-sdk call through WithFcClient and
// captures the resulting HTTP exchange on the local server.
func runFcCaptureCall(t *testing.T, c *AliyunClient, call string, tlsConfig ...*tls.Config) {
	t.Helper()
	_, err := c.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
		if len(tlsConfig) > 0 && tlsConfig[0] != nil {
			c.fcHTTP.GetClient().Transport.(*http.Transport).TLSClientConfig = tlsConfig[0].Clone()
		}
		switch call {
		case "account settings":
			return fcClient.GetAccountSettings(fc.NewGetAccountSettingsInput())
		case "list services":
			return fcClient.ListServices(fc.NewListServicesInput().WithNextToken("token-2").WithLimit(10))
		default:
			name := "wire-test"
			return fcClient.CreateService(&fc.CreateServiceInput{
				ServiceName: tea.String(name), Description: tea.String("captured")})
		}
	})
	if err != nil {
		t.Fatalf("unexpected FC API error: %v", err)
	}
}

func fcCaptureServer(t *testing.T, v4 bool, region string) *httptest.Server {
	t.Helper()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Error(err)
		}
		// verifyWireV4 re-reads the request body to recompute the payload
		// digest, so replay the captured bytes for it.
		r.Body = io.NopCloser(bytes.NewReader(body))
		switch {
		case v4:
			// Recomputes the ACS4 signature from the captured request, which
			// also asserts the physical signing region and the FC product id.
			verifyWireV4(t, r, "fc", region)
			if r.Header.Get("X-Acs-Version") != "2016-08-15" {
				t.Error("expected the FC 1.0 API version header")
			}
			if r.Header.Get("X-Acs-Action") != "" {
				t.Error("the FC 1.0 REST surface must not carry an action header")
			}
			if r.Header.Get("Date") != "" || r.Header.Get("X-Fc-Security-Token") != "" {
				t.Error("legacy FC credential headers must be dropped on V4 requests")
			}
			if !strings.HasPrefix(r.Header.Get("Authorization"), "ACS4-HMAC-SHA256 ") {
				t.Errorf("expected an ACS4-HMAC-SHA256 authorization, got %q", r.Header.Get("Authorization"))
			}
			for name, values := range r.Header {
				for _, value := range values {
					if strings.Contains(value, "test-secret") {
						t.Errorf("header %s leaked the access key secret", name)
					}
				}
			}
			if len(body) > 0 && r.Header.Get("Content-Type") != "application/json" {
				t.Errorf("expected the SDK JSON content type on the wire, got %q", r.Header.Get("Content-Type"))
			}
		default:
			if !strings.HasPrefix(r.Header.Get("Authorization"), "FC test-ak:") {
				t.Errorf("expected the legacy FC signature, got %q", r.Header.Get("Authorization"))
			}
			if r.Header.Get("Date") == "" || r.Header.Get("X-Fc-Security-Token") != "test-token" {
				t.Error("expected the legacy FC date and STS token headers")
			}
			for name := range r.Header {
				if strings.HasPrefix(name, "X-Acs-") {
					t.Errorf("unexpected V4 header %s on a legacy FC request", name)
				}
			}
		}
		if !strings.HasPrefix(r.URL.Path, "/2016-08-15/") {
			t.Errorf("expected the FC 1.0 API path prefix, got %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"RequestId":"request-id"}`)
	})
	if !v4 {
		return httptest.NewTLSServer(handler)
	}
	return httptest.NewServer(handler)
}

// The legacy FC client must keep its original wire signature and endpoint
// resolution when no V4 signing is configured, and must carry the shared V4
// signature (physical region, EO endpoint passthrough, STS token, no key
// material) once sign_version.default=v4 is set. The assertions run against
// the captured HTTP exchange rather than SDK-internal state.
func TestUnitFcClientSignVersionOnWire(t *testing.T) {
	clearProxyEnv(t)
	const region = "cn-edge-xy-01"
	for _, v4 := range []bool{false, true} {
		for _, call := range []string{"account settings", "list services", "create service"} {
			t.Run(fmt.Sprintf("%s/v4=%t", call, v4), func(t *testing.T) {
				server := fcCaptureServer(t, v4, region)
				defer server.Close()
				endpoint := server.Listener.Addr().String()
				c := newFcCaptureClient(t, endpoint, region)
				if v4 {
					c.config.SignVersion.Store("default", "v4")
				}
				if v4 {
					runFcCaptureCall(t, c, call)
				} else {
					runFcCaptureCall(t, c, call, server.Client().Transport.(*http.Transport).TLSClientConfig)
				}
			})
		}
	}
}

// An explicit product-level sign_version must override the provider default,
// and V4 without a physical region fails fast instead of signing.
func TestUnitFcClientV4PrecedenceAndRegion(t *testing.T) {
	clearProxyEnv(t)
	const region = "cn-edge-xy-01"
	server := fcCaptureServer(t, false, region)
	defer server.Close()
	endpoint := server.Listener.Addr().String()

	for _, tc := range []struct {
		name        string
		signVersion func(c *AliyunClient)
		expectError string
	}{
		{name: "fc v1 overrides the v4 default", signVersion: func(c *AliyunClient) {
			c.config.SignVersion.Store("default", "v4")
			c.config.SignVersion.Store("fc", "v1")
		}},
		{name: "fc v4 without the v4 default", signVersion: func(c *AliyunClient) {
			c.config.SignVersion.Store("fc", "v4")
		}},
		{name: "v4 default without a region", signVersion: func(c *AliyunClient) {
			c.config.SignVersion.Store("default", "v4")
			c.config.RegionId, c.RegionId = "", ""
		}, expectError: "requires a signing region"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := newFcCaptureClient(t, endpoint, region)
			tc.signVersion(c)
			if tc.expectError == "" {
				runFcCaptureCall(t, c, "account settings", server.Client().Transport.(*http.Transport).TLSClientConfig)
				return
			}
			_, err := c.WithFcClient(func(fcClient *fc.Client) (interface{}, error) {
				return fcClient.GetAccountSettings(fc.NewGetAccountSettingsInput())
			})
			if err == nil || !strings.Contains(err.Error(), tc.expectError) {
				t.Fatalf("expected a %q diagnostic, got %v", tc.expectError, err)
			}
		})
	}
}

func TestUnitFcCachedAliasesKeepCredentials(t *testing.T) {
	clearProxyEnv(t)
	var seen []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = append(seen, r.Header.Get("Authorization"))
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, "{}")
	}))
	defer server.Close()
	original := resty.DefaultClient
	var clients []*AliyunClient
	for _, ak := range []string{"alias-one", "alias-two"} {
		c := newFcCaptureClient(t, server.Listener.Addr().String(), "cn-hangzhou")
		c.config.SignVersion.Store("default", "v4")
		cred, err := credentials.NewCredential(new(credentials.Config).SetType("access_key").SetAccessKeyId(ak).SetAccessKeySecret("unit-secret"))
		if err != nil {
			t.Fatal(err)
		}
		c.config.Credential = cred
		clients = append(clients, c)
	}
	for _, i := range []int{0, 1, 0} {
		runFcCaptureCall(t, clients[i], "account settings")
		if resty.DefaultClient != original {
			t.Fatal("FC left its HTTP client bound after callback")
		}
	}
	if len(seen) != 3 {
		t.Fatalf("got %d requests", len(seen))
	}
	for i, ak := range []string{"alias-one", "alias-two", "alias-one"} {
		if !strings.Contains(seen[i], "Credential="+ak+"/") {
			t.Errorf("request %d used the wrong signer", i)
		}
	}
}
