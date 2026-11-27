package client

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	pop "github.com/alibabacloud-go/alibabacloud-gateway-pop/client"
	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/credentials-go/credentials"
)

// Recompute the authorization from what the HTTP server received, including
// the actual payload and Host. An algorithm name alone does not prove validity.
func verifyWireV4(t *testing.T, r *http.Request, product, region string) {
	t.Helper()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		t.Error(err)
		return
	}
	payload := fmt.Sprintf("%x", sha256.Sum256(body))
	date := r.Header.Get("X-Acs-Date")
	if len(date) < 10 || r.Header.Get("X-Acs-Content-Sha256") != payload || r.Header.Get("X-Acs-Security-Token") != "test-token" {
		t.Error("missing date, token, or incorrect wire payload digest")
		return
	}
	headers := map[string]*string{"host": &r.Host}
	for name := range r.Header {
		headers[strings.ToLower(name)] = tea.String(r.Header.Get(name))
	}
	query := make(map[string]*string)
	for name := range r.URL.Query() {
		query[name] = tea.String(r.URL.Query().Get(name))
	}
	gateway, err := pop.NewClient()
	if err != nil {
		t.Error(err)
		return
	}
	algorithm := "ACS4-HMAC-SHA256"
	day := strings.ReplaceAll(date[:10], "-", "")
	key := gateway.GetSigningkey(&algorithm, tea.String("test-secret"), &product, &region, &day)
	expected, err := gateway.GetAuthorization(tea.String(r.URL.EscapedPath()), &r.Method, query, headers, &algorithm, &payload, tea.String("test-ak"), key, &product, &region, &day)
	if err != nil || tea.StringValue(expected) != r.Header.Get("Authorization") {
		t.Error("wire signature does not match the expected region, product, headers and body")
	}
}

func TestDefaultSignVersionOnWire(t *testing.T) {
	clearProxyEnv(t)
	for _, v4 := range []bool{false, true} {
		for _, path := range []string{"RPC", "context RPC", "ROA", "OpenAPI RPC", "legacy Tea RPC", "legacy Tea ROA"} {
			t.Run(fmt.Sprintf("%s/v4=%t", path, v4), func(t *testing.T) {
				const region = "ap-jakarta-idpt-2"
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					switch {
					case v4 && !strings.HasPrefix(path, "legacy Tea "):
						verifyWireV4(t, r, "ecs", region)
					case path == "ROA" || path == "legacy Tea ROA":
						if r.Header.Get("x-acs-signature-method") != "HMAC-SHA1" || !strings.HasPrefix(r.Header.Get("Authorization"), "acs ") {
							t.Error("expected legacy ROA signature")
						}
					case path == "OpenAPI RPC":
						if !strings.HasPrefix(r.Header.Get("Authorization"), "ACS3-HMAC-SHA256 ") {
							t.Error("expected existing OpenAPI default signature")
						}
					default:
						q := r.URL.Query()
						if q.Get("SignatureMethod") != "HMAC-SHA1" || q.Get("SignatureVersion") != "1.0" || q.Get("Signature") == "" {
							t.Error("expected legacy RPC signature")
						}
					}
					w.Header().Set("Content-Type", "application/json")
					io.WriteString(w, `{"RequestId":"request-id"}`)
				}))
				defer server.Close()
				endpoint := strings.TrimPrefix(server.URL, "http://")
				c := newLocalRpcClient(t, endpoint)
				c.config.Endpoints.Store("ecs", endpoint)
				c.config.RegionId, c.RegionId = region, region
				c.teaRoaSdkConfig.Protocol = tea.String("HTTP")
				c.teaRpcOpenapiConfig.Protocol = tea.String("HTTP")
				cred, err := credentials.NewCredential(new(credentials.Config).SetType("sts").SetAccessKeyId("test-ak").SetAccessKeySecret("test-secret").SetSecurityToken("test-token"))
				if err != nil {
					t.Fatal(err)
				}
				c.config.Credential = cred
				if v4 {
					c.config.SignVersion.Store("default", "v4")
				}
				input := map[string]interface{}{"Name": "test value"}
				var response map[string]interface{}
				switch path {
				case "RPC":
					response, err = c.RpcPostWithEndpoint("Ecs", "2014-05-26", "DescribeInstances", input, input, false, endpoint)
				case "context RPC":
					response, err = c.RpcPostWithContext(context.Background(), "Ecs", "2014-05-26", "DescribeInstances", input, input)
				case "ROA":
					response, err = c.RoaPostWithApiName("Ecs", "2014-05-26", "Example", "/test", nil, nil, input, false)
				case "legacy Tea RPC":
					conn, initErr := c.NewTeaCommonClient(endpoint)
					if initErr != nil {
						t.Fatal(initErr)
					}
					response, err = conn.DoRequest(tea.String("DescribeInstances"), nil, tea.String("POST"), tea.String("2014-05-26"), tea.String("AK"), nil, nil, &util.RuntimeOptions{})
				case "legacy Tea ROA":
					conn, initErr := c.NewTeaRoaCommonClient(endpoint)
					if initErr != nil {
						t.Fatal(initErr)
					}
					response, err = conn.DoRequestWithAction(tea.String("Example"), tea.String("2014-05-26"), nil, tea.String("POST"), tea.String("AK"), tea.String("/test"), nil, nil, nil, &util.RuntimeOptions{})
					if err == nil {
						response, _ = response["body"].(map[string]interface{})
					}
				case "OpenAPI RPC":
					params := new(openapi.Params).SetAction("DescribeInstances").SetVersion("2014-05-26").SetMethod("POST").SetStyle("RPC").SetProtocol("HTTP").SetAuthType("AK").SetReqBodyType("formData").SetBodyType("json").SetPathname("/")
					response, err = c.Do("Ecs", params, nil, input, nil, nil, false)
				}
				if err != nil || response["RequestId"] != "request-id" {
					t.Fatalf("unexpected API response: %v, %v", response, err)
				}
			})
		}
	}
}

func TestDefaultV4Precedence(t *testing.T) {
	c := newLocalRpcClient(t, "unused.example")
	c.config.SignVersion.Store("default", "v4")
	// The SDKv2 schema materializes omitted optional string fields as "".
	c.config.SignVersion.Store("sls", "")
	if got := effectiveSignVersion(c.config.SignVersion, "sls"); got != "v4" {
		t.Fatalf("inherited SLS signature = %q", got)
	}
	c.config.SignVersion.Store("sls", "v1")
	if c.useV4ForProduct("Sls") || effectiveSignVersion(c.config.SignVersion, "sls") != "v1" {
		t.Fatal("explicit SLS signature did not override the default")
	}
}

func TestDefaultV4LegacySDKServiceDecision(t *testing.T) {
	clearProxyEnv(t)
	for _, status := range []int{http.StatusOK, http.StatusForbidden} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				q := r.URL.Query()
				if q.Get("SignatureMethod") != "HMAC-SHA1" || q.Get("Signature") == "" || q.Get("Action") != "DescribeRegions" {
					t.Error("expected a signed legacy DescribeRegions request")
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(status)
				if status == http.StatusForbidden {
					io.WriteString(w, `{"Code":"SignatureVersionNotSupported","Message":"endpoint requires V4","RequestId":"denied"}`)
					return
				}
				io.WriteString(w, `{"RequestId":"accepted"}`)
			}))
			defer server.Close()
			endpoint := strings.TrimPrefix(server.URL, "http://")
			c := newLocalRpcClient(t, endpoint)
			c.RegionId, c.config.RegionId = "test-legacy-signing", "test-legacy-signing"
			c.config.Endpoints.Store("ecs", endpoint)
			c.config.SignVersion.Store("default", "v4")
			result, err := c.WithEcsClient(func(conn *ecs.Client) (interface{}, error) {
				return conn.DescribeRegions(ecs.CreateDescribeRegionsRequest())
			})
			if status == http.StatusForbidden {
				if err == nil || !strings.Contains(err.Error(), "SignatureVersionNotSupported") {
					t.Fatalf("expected service rejection, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if result.(*ecs.DescribeRegionsResponse).RequestId != "accepted" {
				t.Fatal("missing service response")
			}
		})
	}
}
