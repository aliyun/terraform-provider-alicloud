package client

// Wire regressions for the legacy SDK calls migrated onto the shared RPC
// entry points (ECS regions, CMS, SLS POP) and for the NAS open-service call.
// Each migrated group keeps one legacy HMAC-SHA1 regression and one EO V4
// assertion, captured on a local HTTP server with temporary STS credentials:
// the token proves temporary-credential coverage, and the signing secret must
// never appear outside the signature.
//
// Every branch also pins the Host header to the bare endpoint. Both tea-rpc
// and the signed transport derive the Host header from the endpoint value, so
// a regression that folds the SLS POP API path into the endpoint would blank
// the Host header — the real service rejects an empty Host.

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/alibabacloud-go/tea/tea"
	sdkerrors "github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sls"
	"github.com/aliyun/credentials-go/credentials"
)

// The legacy SDK registers endpoints globally by product and region.
const legacyWireRegion = "cn-unit-test"

// newLegacyWireClient points one product at a local test server and swaps the
// credential for a temporary STS token, so the legacy path proves
// SecurityToken lands on the query and the V4 path proves
// X-Acs-Security-Token lands on the wire.
func newLegacyWireClient(t *testing.T, productCode, endpoint string) *AliyunClient {
	t.Helper()
	clearProxyEnv(t)
	credential, err := credentials.NewCredential(new(credentials.Config).
		SetType("sts").
		SetAccessKeyId("test-ak").
		SetAccessKeySecret("test-secret").
		SetSecurityToken("test-token"))
	if err != nil {
		t.Fatalf("building the credential: %s", err)
	}
	c := newLocalRpcClient(t, endpoint)
	c.config.RegionId, c.RegionId = legacyWireRegion, legacyWireRegion
	c.teaRpcOpenapiConfig.Protocol = tea.String("HTTP")
	c.config.Credential = credential
	c.config.Endpoints.Store(productCode, endpoint)
	return c
}

func isV4WireRequest(r *http.Request) bool {
	return strings.HasPrefix(r.Header.Get("Authorization"), "ACS4-HMAC-SHA256 ")
}

// captureLegacyRpcWire verifies what every migrated legacy call must keep on
// the wire: POST to the bare endpoint host, the signing secret nowhere
// outside the signature, and the old HMAC-SHA1 parameters with the temporary
// STS credentials on the legacy branch. It returns whether the request was
// legacy signed; the caller asserts the V4 branch with verifyWireV4.
func captureLegacyRpcWire(t *testing.T, r *http.Request, endpoint string) (legacy bool) {
	t.Helper()
	if r.Method != http.MethodPost {
		t.Errorf("expected POST, got %s", r.Method)
	}
	if r.Host == "" || r.Host != endpoint {
		t.Errorf("Host header %q is not the bare endpoint %q", r.Host, endpoint)
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		t.Errorf("reading the request body: %s", err)
	}
	r.Body = io.NopCloser(bytes.NewReader(body))
	if strings.Contains(r.URL.String(), "test-secret") || bytes.Contains(body, []byte("test-secret")) {
		t.Error("signing secret leaked into the request URL or body")
	}
	if legacy = !isV4WireRequest(r); legacy {
		q := r.URL.Query()
		if q.Get("SignatureMethod") != "HMAC-SHA1" || q.Get("SignatureVersion") != "1.0" || q.Get("Signature") == "" {
			t.Error("expected legacy RPC signature parameters")
		}
		if q.Get("AccessKeyId") != "test-ak" || q.Get("SecurityToken") != "test-token" {
			t.Error("expected the temporary STS credentials on the query")
		}
	}
	return legacy
}

// TestUnitLegacySdkRpcRegionsWire covers the alicloud_regions migration off
// ecsClient: unconfigured sign_version keeps the old signature, an EO
// deployment with sign_version.default = v4 signs DescribeRegions through the
// V4 pop gateway.
func TestUnitLegacySdkRpcRegionsWire(t *testing.T) {
	var endpoint string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		legacy := captureLegacyRpcWire(t, r, endpoint)
		if r.URL.Path != "/" {
			t.Errorf("expected the RPC root path, got %q", r.URL.Path)
		}
		q := r.URL.Query()
		if legacy {
			if q.Get("Action") != "DescribeRegions" || q.Get("Version") != "2014-05-26" || q.Get("RegionId") != legacyWireRegion {
				t.Error("expected the DescribeRegions action, version and region on the query")
			}
		} else {
			verifyWireV4(t, r, "ecs", legacyWireRegion)
			if r.Header.Get("x-acs-action") != "DescribeRegions" || r.Header.Get("x-acs-version") != "2014-05-26" {
				t.Error("expected the DescribeRegions action and version headers")
			}
			if q.Get("RegionId") != legacyWireRegion {
				t.Error("expected the region on the V4 query")
			}
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"RequestId":"request-id"}`)
	}))
	defer server.Close()
	endpoint = strings.TrimPrefix(server.URL, "http://")

	for _, tc := range []struct {
		name        string
		signVersion string
	}{
		{name: "unconfigured-keeps-legacy-signature"},
		{name: "default-v4-signs-eo-v4", signVersion: "v4"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := newLegacyWireClient(t, "ecs", endpoint)
			if tc.signVersion != "" {
				c.config.SignVersion.Store("default", tc.signVersion)
			}
			request := ecs.CreateDescribeRegionsRequest()
			request.RegionId = legacyWireRegion
			response, err := c.LegacySdkRpc(request, "/", ecs.CreateDescribeRegionsResponse())
			if err != nil || response["RequestId"] != "request-id" {
				t.Fatalf("unexpected DescribeRegions response: %v, %v", response, err)
			}
		})
	}
}

// TestUnitLegacySdkRpcCmsContactGroupWire covers the CMS migration off
// WithCmsClient: both signature paths carry the typed SDK serialization of
// PutContactGroup — repeated parameters (ContactNames.1, ContactNames.2),
// scalars and the Boolean wrapper — and a per-product v1 override keeps the
// old signature even when the default is v4.
func TestUnitLegacySdkRpcCmsContactGroupWire(t *testing.T) {
	var endpoint string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		legacy := captureLegacyRpcWire(t, r, endpoint)
		if r.URL.Path != "/" {
			t.Errorf("expected the RPC root path, got %q", r.URL.Path)
		}
		q := r.URL.Query()
		if legacy {
			if q.Get("Action") != "PutContactGroup" || q.Get("Version") != "2019-01-01" {
				t.Error("expected the PutContactGroup action and version on the query")
			}
		} else {
			verifyWireV4(t, r, "cms", legacyWireRegion)
			if r.Header.Get("x-acs-action") != "PutContactGroup" || r.Header.Get("x-acs-version") != "2019-01-01" {
				t.Error("expected the PutContactGroup action and version headers")
			}
		}
		for _, check := range []struct{ name, value string }{
			{"ContactGroupName", "tf-test-group"},
			{"Describe", "unit test"},
			{"EnableSubscribed", "true"},
			{"ContactNames.1", "alice"},
			{"ContactNames.2", "bob"},
		} {
			if got := q.Get(check.name); got != check.value {
				t.Errorf("expected %s=%q on the query, got %q", check.name, check.value, got)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"Code":200,"RequestId":"request-id"}`)
	}))
	defer server.Close()
	endpoint = strings.TrimPrefix(server.URL, "http://")

	for _, tc := range []struct {
		name        string
		defaultSign string
		productSign string
	}{
		{name: "unconfigured-keeps-legacy-signature"},
		{name: "default-v4-signs-eo-v4", defaultSign: "v4"},
		{name: "cms-v1-overrides-default-v4", defaultSign: "v4", productSign: "v1"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := newLegacyWireClient(t, "cms", endpoint)
			if tc.defaultSign != "" {
				c.config.SignVersion.Store("default", tc.defaultSign)
			}
			if tc.productSign != "" {
				c.config.SignVersion.Store("cms", tc.productSign)
			}
			request := cms.CreatePutContactGroupRequest()
			request.ContactGroupName = "tf-test-group"
			contactNames := []string{"alice", "bob"}
			request.ContactNames = &contactNames
			request.Describe = "unit test"
			request.EnableSubscribed = requests.NewBoolean(true)
			response, err := c.LegacySdkRpc(request, "/", cms.CreatePutContactGroupResponse())
			if tc.defaultSign == "" || tc.productSign == "v1" {
				if response["Code"] != "200" {
					t.Fatalf("legacy SDK should decode a numeric Code into a string, got %#v", response["Code"])
				}
			}
			if err != nil || response["RequestId"] != "request-id" {
				t.Fatalf("unexpected PutContactGroup response: %v, %v", response, err)
			}
		})
	}
}

// TestUnitLegacySdkRpcSlsPopPathWire covers the SLS migration off
// WithLogPopClient: both signature paths must reach the fixed POP path with a
// clean Host header — the old client appended the path to its Domain, and
// folding it into the tea-rpc endpoint instead would blank the Host header.
func TestUnitLegacySdkRpcSlsPopPathWire(t *testing.T) {
	// Fixed SLS POP API path — the old client appended "/open-api" to its
	// Domain and the old SDK's query builder turned that into the "/open-api/"
	// request path; it lives in the sls service package as slsPopApiPath and is
	// repeated here because these tests cannot import the sls service package.
	const slsPopPath = "/open-api/"
	var endpoint string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		legacy := captureLegacyRpcWire(t, r, endpoint)
		if r.URL.Path != slsPopPath {
			t.Errorf("expected the SLS POP path %q, got %q", slsPopPath, r.URL.Path)
		}
		q := r.URL.Query()
		if legacy {
			if q.Get("Action") != "DescribeApp" || q.Get("Version") != "2019-10-23" || q.Get("RegionId") != legacyWireRegion {
				t.Error("expected the DescribeApp action, version and region on the query")
			}
		} else {
			verifyWireV4(t, r, "sls", legacyWireRegion)
			if r.Header.Get("x-acs-action") != "DescribeApp" || r.Header.Get("x-acs-version") != "2019-10-23" {
				t.Error("expected the DescribeApp action and version headers")
			}
		}
		if q.Get("AppName") != "audit" {
			t.Errorf("expected AppName=audit on the query, got %q", q.Get("AppName"))
		}
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"RequestId":"request-id"}`)
	}))
	defer server.Close()
	endpoint = strings.TrimPrefix(server.URL, "http://")

	for _, tc := range []struct {
		name        string
		signVersion string
	}{
		{name: "unconfigured-keeps-legacy-signature"},
		{name: "default-v4-signs-eo-v4", signVersion: "v4"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := newLegacyWireClient(t, "sls", endpoint)
			if tc.signVersion != "" {
				c.config.SignVersion.Store("default", tc.signVersion)
			}
			request := sls.CreateDescribeAppRequest()
			request.AppName = "audit"
			response, err := c.LegacySdkRpc(request, slsPopPath, sls.CreateDescribeAppResponse())
			if err != nil || response["RequestId"] != "request-id" {
				t.Fatalf("unexpected DescribeApp response: %v, %v", response, err)
			}
		})
	}
}

func TestUnitNasOpenServiceWire(t *testing.T) {
	for _, tc := range []struct {
		name           string
		defaultVersion string
		productVersion string
		wantV4         bool
	}{
		{"unconfigured", "", "", false},
		{"nas-v1-overrides-default", "v4", "v1", false},
		{"default-v4", "v4", "", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			const activationHost = "nas.cn-hangzhou.aliyuncs.com"
			var endpoint string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				wantHost := activationHost
				if tc.wantV4 {
					wantHost = endpoint
				}
				legacy := captureLegacyRpcWire(t, r, wantHost)
				if legacy == tc.wantV4 || r.URL.Path != "/" {
					t.Errorf("unexpected signature or path: legacy=%t path=%q", legacy, r.URL.Path)
				}
				if tc.wantV4 {
					verifyWireV4(t, r, "nas", legacyWireRegion)
				}
				action, version := r.URL.Query().Get("Action"), r.URL.Query().Get("Version")
				if tc.wantV4 {
					action, version = r.Header.Get("x-acs-action"), r.Header.Get("x-acs-version")
				}
				if action != "OpenNASService" || version != "2017-06-26" {
					t.Errorf("unexpected action/version %s/%s", action, version)
				}
				w.Header().Set("Content-Type", "application/json")
				io.WriteString(w, `{"RequestId":"request-id","OrderId":"order-123"}`)
			}))
			defer server.Close()
			endpoint = strings.TrimPrefix(server.URL, "http://")
			c := newLegacyWireClient(t, "nas", endpoint)
			c.teaSdkConfig.SetHttpProxy(server.URL)
			c.config.SecurityToken = "test-token"
			c.config.NasEndpoint = "HTTPS://" + endpoint
			if tc.defaultVersion != "" {
				c.config.SignVersion.Store("default", tc.defaultVersion)
			}
			if tc.productVersion != "" {
				c.config.SignVersion.Store("nas", tc.productVersion)
			}
			response, err := c.OpenNasService()
			if err != nil || response["OrderId"] != "order-123" {
				t.Fatalf("OpenNASService: %v, %v", response, err)
			}
		})
	}
}

func TestUnitLegacySdkRpcRetry(t *testing.T) {
	for _, tc := range []struct {
		product, path string
		request       func() requests.AcsRequest
		response      func() responses.AcsResponse
	}{
		{"ecs", "/", func() requests.AcsRequest { return ecs.CreateDescribeRegionsRequest() }, func() responses.AcsResponse { return ecs.CreateDescribeRegionsResponse() }},
		{"cms", "/", func() requests.AcsRequest { return cms.CreateDescribeContactGroupListRequest() }, func() responses.AcsResponse { return cms.CreateDescribeContactGroupListResponse() }},
		{"sls", "/open-api/", func() requests.AcsRequest { return sls.CreateDescribeAppRequest() }, func() responses.AcsResponse { return sls.CreateDescribeAppResponse() }},
	} {
		for _, status := range []int{http.StatusServiceUnavailable, http.StatusBadRequest} {
			t.Run(tc.product+"/"+http.StatusText(status), func(t *testing.T) {
				var attempts atomic.Int32
				var endpoint string
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					if !captureLegacyRpcWire(t, r, endpoint) || r.URL.Path != tc.path {
						t.Error("expected original signature and path")
					}
					w.Header().Set("Content-Type", "application/json")
					if attempts.Add(1) == 1 {
						w.WriteHeader(status)
						io.WriteString(w, `{"Code":"TestError","Message":"failed","RequestId":"id"}`)
						return
					}
					io.WriteString(w, `{"RequestId":"ok"}`)
				}))
				defer server.Close()
				endpoint = strings.TrimPrefix(server.URL, "http://")
				c := newLegacyWireClient(t, tc.product, endpoint)
				result, err := c.LegacySdkRpc(tc.request(), tc.path, tc.response())
				if status == http.StatusServiceUnavailable {
					if err != nil || attempts.Load() != 2 || result["RequestId"] != "ok" {
						t.Fatalf("expected retry success: attempts=%d result=%v err=%v", attempts.Load(), result, err)
					}
				} else {
					sdkErr, ok := err.(*sdkerrors.ServerError)
					if !ok || sdkErr.HttpStatus() != status || sdkErr.ErrorCode() != "TestError" || attempts.Load() != 1 {
						t.Fatalf("expected one SDK server error: attempts=%d err=%T %v", attempts.Load(), err, err)
					}
				}
			})
		}
	}
}

// TestUnitLegacySdkRpcAllowsSdkTrackingHeader pins the guard regression: a
// plain len() check on the headers map rejected every real call, because no
// migrated caller sets a business header — the SDK's own factory tracking
// header (x-sdk-invoke-type) is already in the map. Nothing listens on the
// endpoint, so a call that passes the guard reaches the network and fails on
// the connection; the guard's own rejection message is the failure this test
// must not see.
func TestUnitLegacySdkRpcAllowsSdkTrackingHeader(t *testing.T) {
	c := newLegacyWireClient(t, "ecs", "127.0.0.1:1")
	c.config.SignVersion.Store("default", "v4")
	_, err := c.LegacySdkRpc(ecs.CreateDescribeRegionsRequest(), "/", ecs.CreateDescribeRegionsResponse())
	if err != nil && strings.Contains(err.Error(), "does not support") {
		t.Fatalf("the factory tracking header must stay allowed: %v", err)
	}
}

// legacyHeaderParamRequest mirrors the generated request shape but carries one
// position:"Header" business field — the SDK's own capitalization, as in the
// generated dyplsapi/ververica requests — the only thing the guard exists to
// catch: InitParams serializes such a field into request.Headers via
// addHeaderParam while the shared entry points read only QueryParams/FormParams,
// so the value would otherwise be silently dropped.
type legacyHeaderParamRequest struct {
	*requests.RpcRequest
	ChannelToken string `name:"X-Channel-Token" position:"Header"`
}

// TestUnitLegacySdkRpcRejectsHeaderParams is the positive counterpart of the
// tracking-header test: the baseline snapshot is taken before InitParams, so a
// business header field the serialization adds must fail loudly instead of
// being dropped. The guard fires before any network access.
func TestUnitLegacySdkRpcRejectsHeaderParams(t *testing.T) {
	c := newLegacyWireClient(t, "cms", "127.0.0.1:1")
	c.config.SignVersion.Store("default", "v4")
	request := &legacyHeaderParamRequest{RpcRequest: &requests.RpcRequest{}}
	request.InitWithApiInfo("Cms", "2019-01-01", "PutContact", "cms", "openAPI")
	request.ChannelToken = "channel-secret"
	_, err := c.LegacySdkRpc(request, "/", cms.CreatePutContactResponse())
	if err == nil || !strings.Contains(err.Error(), "does not support requests with header parameters") {
		t.Fatalf("expected the header parameter guard to reject the request, got: %v", err)
	}
}

// TestUnitLegacyRpcResponse covers the response-side counterpart of
// LegacySdkRpc: a shared-path RPC map rehydrates into the typed SDK response
// with the envelope and nested list fields populated, so a caller migrated
// off its per-product SDK client keeps its downstream field access.
func TestUnitLegacyRpcResponse(t *testing.T) {
	raw := map[string]interface{}{
		"Code": "200",
		"Contacts": map[string]interface{}{
			"Contact": []interface{}{
				map[string]interface{}{"Name": "alarm-contact-1"},
			},
		},
	}
	resp := &cms.DescribeContactListResponse{}
	if err := LegacyRpcResponse(raw, resp); err != nil {
		t.Fatalf("decoding: %s", err)
	}
	if resp.Code != "200" || len(resp.Contacts.Contact) != 1 || resp.Contacts.Contact[0].Name != "alarm-contact-1" {
		t.Fatalf("typed fields lost after rehydration: %+v", resp)
	}
}
