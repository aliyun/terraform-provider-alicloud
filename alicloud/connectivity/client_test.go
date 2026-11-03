package connectivity

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	ossgateway "github.com/alibabacloud-go/alibabacloud-gateway-oss/client"
	spi "github.com/alibabacloud-go/alibabacloud-gateway-spi/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeOssOpenAPIResponse_HttpsConfigurationRegression(t *testing.T) {
	gateway, err := ossgateway.NewClient()
	if err != nil {
		t.Fatalf("create OSS gateway client: %v", err)
	}
	context := &spi.InterceptorContext{
		Request: &spi.InterceptorContextRequest{
			Action:   tea.String("GetBucketHttpsConfig"),
			BodyType: tea.String("xml"),
		},
		Response: &spi.InterceptorContextResponse{
			StatusCode: tea.Int(200),
			Headers:    map[string]*string{},
			Body: strings.NewReader(`<HttpsConfiguration>
  <TLS><Enable>true</Enable><TLSVersion>TLSv1.2</TLSVersion></TLS>
</HttpsConfiguration>`),
		},
	}
	if err := gateway.ModifyResponse(context, &spi.AttributeMap{Key: map[string]*string{}}); err != nil {
		t.Fatalf("parse response through OSS gateway: %v", err)
	}
	sdkBody, ok := context.Response.DeserializedBody.(map[string]interface{})
	if !ok {
		t.Fatalf("OSS gateway body should be map[string]interface{}, got %T", context.Response.DeserializedBody)
	}

	response := map[string]interface{}{
		"statusCode": 200,
		"body":       sdkBody,
	}
	getHttpsConfiguration := func(response map[string]interface{}) map[string]interface{} {
		body := response["body"].(map[string]interface{})
		return body["HttpsConfiguration"].(map[string]interface{})
	}

	assert.Panics(t, func() {
		getHttpsConfiguration(response)
	}, "the unnormalized **HttpsConfiguration fixture must reproduce the provider crash")

	normalized, err := normalizeOssOpenAPIResponse(response)
	if err != nil {
		t.Fatalf("normalize OSS response: %v", err)
	}
	var config map[string]interface{}
	assert.NotPanics(t, func() {
		config = getHttpsConfiguration(normalized)
	}, "the normalized response must be safe for the provider's map assertion")
	tls := config["TLS"].(map[string]interface{})
	if tls["Enable"] != true {
		t.Fatalf("TLS.Enable = %#v, want true", tls["Enable"])
	}
	if status, ok := normalized["statusCode"].(json.Number); !ok || status.String() != "200" {
		t.Fatalf("statusCode should preserve legacy json.Number behavior, got %#v (%T)", normalized["statusCode"], normalized["statusCode"])
	}
}

var endpointMap sync.Map
var signVersion sync.Map

func NewTestClient(t *testing.T) *AliyunClient {
	accessKey := os.Getenv("ALICLOUD_ACCESS_KEY")
	secretKey := os.Getenv("ALICLOUD_SECRET_KEY")

	if accessKey == "" || secretKey == "" {
		t.Skipf("Skipping the test case as some necessary params are empty")
		t.Skipped()
	}

	config := &Config{
		Region:      Beijing,
		RegionId:    "cn-beijing",
		AccessKey:   accessKey,
		SecretKey:   secretKey,
		Protocol:    "HTTPS",
		Endpoints:   &endpointMap,
		SignVersion: &signVersion,
	}

	credentialConfig := new(credentials.Config).
		SetType("access_key").
		SetAccessKeyId(accessKey).
		SetAccessKeySecret(secretKey)

	credential, err := credentials.NewCredential(credentialConfig)
	if err != nil {
		t.Fatalf("create credential failed: %v", err)
	}
	config.Credential = credential

	client, err := config.Client()
	if err != nil {
		t.Fatalf("initial client failed: %v", err)
	}
	return client
}

func TestGetTableStoreConfigProxy(t *testing.T) {
	assert.True(t, getTableStoreConfig().ProxyFromEnvironment)
}

func TestTableStoreRequestUsesEnvironmentProxy(t *testing.T) {
	const helperEnv = "TABLESTORE_PROXY_TEST_HELPER"
	if os.Getenv(helperEnv) == "1" {
		credential, err := credentials.NewCredential(new(credentials.Config).
			SetType("access_key").
			SetAccessKeyId("access-key").
			SetAccessKeySecret("secret-key"))
		if err != nil {
			t.Fatalf("create credential: %v", err)
		}
		client := &AliyunClient{
			config: &Config{
				OtsEndpoint: "http://tablestore.invalid",
				Credential:  credential,
			},
			tablestoreconnByInstanceName: make(map[string]*tablestore.TableStoreClient),
		}
		if _, err := client.WithTableStoreClient("instance", func(client *tablestore.TableStoreClient) (interface{}, error) {
			return client.ListTable()
		}); err != nil {
			t.Fatalf("list tables through proxy: %v", err)
		}
		return
	}

	type proxyRequest struct {
		method string
		host   string
		path   string
	}
	requests := make(chan proxyRequest, 1)
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case requests <- proxyRequest{method: r.Method, host: r.URL.Host, path: r.URL.Path}:
		default:
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer proxy.Close()

	t.Setenv(helperEnv, "1")
	t.Setenv("HTTP_PROXY", proxy.URL)
	t.Setenv("http_proxy", "")
	t.Setenv("NO_PROXY", "")
	t.Setenv("no_proxy", "")
	t.Setenv("REQUEST_METHOD", "")
	cmd := exec.Command(os.Args[0], "-test.run=^TestTableStoreRequestUsesEnvironmentProxy$")
	cmd.Env = os.Environ()
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("proxy test subprocess failed: %v\n%s", err, output)
	}

	select {
	case request := <-requests:
		assert.Equal(t, http.MethodPost, request.method)
		assert.Equal(t, "tablestore.invalid", request.host)
		assert.Equal(t, "/ListTable", request.path)
	default:
		t.Fatal("proxy did not receive the TableStore request")
	}
}

func TestUnitCommonWithEcsClient_UsingHttpMock(t *testing.T) {
	client := NewTestClient(t)

	res, _ := client.WithEcsClient(func(c *ecs.Client) (interface{}, error) {
		req := ecs.CreateDescribeInstancesRequest()
		return c.DescribeInstances(req)
	})

	assert.NotNil(t, res)
}

func TestUnitCommonWithEcsClient_Proxy(t *testing.T) {
	client := NewTestClient(t)

	testCases := []struct {
		name          string
		proxyURL      string
		skipProxy     bool
		expectedHTTP  string
		expectedHTTPS string
	}{
		{
			name:         "HTTP协议代理设置",
			proxyURL:     "http://proxy.example.com:8080",
			expectedHTTP: "http://proxy.example.com:8080",
		},
		{
			name:          "HTTPS协议代理设置",
			proxyURL:      "https://proxy.example.com:8443",
			expectedHTTPS: "https://proxy.example.com:8443",
		},
		{
			name: "无代理设置",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ecsconn, err := ecs.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(time.Duration(60)*time.Second), client.config.getAuthCredential(true))
			assert.NoError(t, err)

			if tc.expectedHTTP != "" {
				client.config.Protocol = "HTTP"
				t.Setenv("HTTP_PROXY", tc.proxyURL)
				proxy, err := client.getHttpProxy()
				ecsconn.SetHttpProxy(proxy.String())
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedHTTP, ecsconn.GetHttpProxy())
				assert.Empty(t, ecsconn.GetHttpsProxy())
			} else if tc.expectedHTTPS != "" {
				client.config.Protocol = "HTTPS"
				t.Setenv("HTTPS_PROXY", tc.proxyURL)
				proxy, err := client.getHttpProxy()
				ecsconn.SetHttpsProxy(proxy.String())
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedHTTPS, ecsconn.GetHttpsProxy())
				assert.Empty(t, ecsconn.GetHttpProxy())
			} else {
				assert.Empty(t, ecsconn.GetHttpProxy())
				assert.Empty(t, ecsconn.GetHttpsProxy())
			}
		})
	}
}

func TestUnitCommonWithVpcClient_Proxy(t *testing.T) {
	client := NewTestClient(t)

	testCases := []struct {
		name          string
		proxyURL      string
		skipProxy     bool
		expectedHTTP  string
		expectedHTTPS string
	}{
		{
			name:         "HTTP协议代理设置",
			proxyURL:     "http://proxy.example.com:8080",
			expectedHTTP: "http://proxy.example.com:8080",
		},
		{
			name:          "HTTPS协议代理设置",
			proxyURL:      "https://proxy.example.com:8443",
			expectedHTTPS: "https://proxy.example.com:8443",
		},
		{
			name: "无代理设置",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			vpcconn, err := vpc.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(time.Duration(60)*time.Second), client.config.getAuthCredential(true))
			assert.NoError(t, err)

			if tc.expectedHTTP != "" {
				client.config.Protocol = "HTTP"
				t.Setenv("HTTP_PROXY", tc.proxyURL)
				proxy, err := client.getHttpProxy()
				vpcconn.SetHttpProxy(proxy.String())
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedHTTP, vpcconn.GetHttpProxy())
				assert.Empty(t, vpcconn.GetHttpsProxy())
			} else if tc.expectedHTTPS != "" {
				client.config.Protocol = "HTTPS"
				t.Setenv("HTTPS_PROXY", tc.proxyURL)
				proxy, err := client.getHttpProxy()
				vpcconn.SetHttpsProxy(proxy.String())
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedHTTPS, vpcconn.GetHttpsProxy())
				assert.Empty(t, vpcconn.GetHttpProxy())
			} else {
				assert.Empty(t, vpcconn.GetHttpProxy())
				assert.Empty(t, vpcconn.GetHttpsProxy())
			}
		})
	}
}

func TestUnitCommonWithRamClient_Proxy(t *testing.T) {
	client := NewTestClient(t)

	testCases := []struct {
		name          string
		proxyURL      string
		skipProxy     bool
		expectedHTTP  string
		expectedHTTPS string
	}{
		{
			name:         "HTTP协议代理设置",
			proxyURL:     "http://proxy.example.com:8080",
			expectedHTTP: "http://proxy.example.com:8080",
		},
		{
			name:          "HTTPS协议代理设置",
			proxyURL:      "https://proxy.example.com:8443",
			expectedHTTPS: "https://proxy.example.com:8443",
		},
		{
			name: "无代理设置",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ramconn, err := ram.NewClientWithOptions(client.config.RegionId, client.getSdkConfig(0), client.config.getAuthCredential(true))
			assert.NoError(t, err)

			if tc.expectedHTTP != "" {
				client.config.Protocol = "HTTP"
				t.Setenv("HTTP_PROXY", tc.proxyURL)
				proxy, err := client.getHttpProxy()
				ramconn.SetHttpProxy(proxy.String())
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedHTTP, ramconn.GetHttpProxy())
				assert.Empty(t, ramconn.GetHttpsProxy())
			} else if tc.expectedHTTPS != "" {
				client.config.Protocol = "HTTPS"
				t.Setenv("HTTPS_PROXY", tc.proxyURL)
				proxy, err := client.getHttpProxy()
				ramconn.SetHttpsProxy(proxy.String())
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedHTTPS, ramconn.GetHttpsProxy())
				assert.Empty(t, ramconn.GetHttpProxy())
			} else {
				assert.Empty(t, ramconn.GetHttpProxy())
				assert.Empty(t, ramconn.GetHttpsProxy())
			}
		})
	}
}

// TestUnitTeaNoProxyHost verifies the host[:port] normalization that hands tea
// the exact request host it compares against its NO_PROXY list. tea-roa sets
// the host header to client.EndpointHost and never sets request.Port, so an
// endpoint such as "host:443" compares as "host:443"; the bypass value must
// preserve the port to match.
func TestUnitTeaNoProxyHost(t *testing.T) {
	cases := []struct {
		name     string
		endpoint string
		want     string
	}{
		{"bare hostname", "tablestore.cn-shanghai.aliyuncs.com", "tablestore.cn-shanghai.aliyuncs.com"},
		{"https scheme", "https://tablestore.cn-shanghai.aliyuncs.com", "tablestore.cn-shanghai.aliyuncs.com"},
		{"host with port", "tablestore.cn-shanghai.aliyuncs.com:443", "tablestore.cn-shanghai.aliyuncs.com:443"},
		{"scheme and port", "http://foo.com:8080", "foo.com:8080"},
		{"empty endpoint", "", ""},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, teaNoProxyHost(tc.endpoint))
		})
	}
}

// TestUnitResolveTeaProxy verifies the tea-roa proxy decision reuses the
// provider's regexp-based skipProxy (suffix/wildcard matching) instead of
// tea's own exact host comparison. It covers every NO_PROXY form the source
// task requires: suffix (.aliyuncs.com), wildcard (*.aliyuncs.com), exact
// host, star (*), host:443 (port-bearing miss), a suffix entry bypassing a
// port-bearing endpoint (bypass value preserves the port), empty, and no
// HTTPS_PROXY.
func TestUnitResolveTeaProxy(t *testing.T) {
	const proxyURL = "http://proxy.example.com:8080"
	const host = "tablestore.cn-shanghai.aliyuncs.com"

	cases := []struct {
		name            string
		httpsProxy      string
		noProxy         string
		protocol        string
		endpoint        string
		wantUseProxy    bool
		wantProxyURL    string
		wantNoProxyHost string
	}{
		{
			name:       "no proxy environment configured",
			httpsProxy: "", noProxy: "", protocol: "HTTPS", endpoint: host,
			wantUseProxy: false, wantProxyURL: "", wantNoProxyHost: "",
		},
		{
			name:       "suffix no_proxy entry bypasses proxy",
			httpsProxy: proxyURL, noProxy: ".aliyuncs.com", protocol: "HTTPS", endpoint: host,
			wantUseProxy: false, wantProxyURL: "", wantNoProxyHost: host,
		},
		{
			name:       "wildcard no_proxy entry bypasses proxy",
			httpsProxy: proxyURL, noProxy: "*.aliyuncs.com", protocol: "HTTPS", endpoint: host,
			wantUseProxy: false, wantProxyURL: "", wantNoProxyHost: host,
		},
		{
			name:       "exact host no_proxy entry bypasses proxy",
			httpsProxy: proxyURL, noProxy: host, protocol: "HTTPS", endpoint: host,
			wantUseProxy: false, wantProxyURL: "", wantNoProxyHost: host,
		},
		{
			name:       "star no_proxy entry bypasses proxy",
			httpsProxy: proxyURL, noProxy: "*", protocol: "HTTPS", endpoint: host,
			wantUseProxy: false, wantProxyURL: "", wantNoProxyHost: host,
		},
		{
			name:       "suffix no_proxy bypasses port-bearing endpoint",
			httpsProxy: proxyURL, noProxy: ".aliyuncs.com", protocol: "HTTPS", endpoint: host + ":443",
			wantUseProxy: false, wantProxyURL: "", wantNoProxyHost: host + ":443",
		},
		{
			name:       "port-bearing no_proxy entry does not match",
			httpsProxy: proxyURL, noProxy: host + ":443", protocol: "HTTPS", endpoint: host,
			wantUseProxy: true, wantProxyURL: proxyURL, wantNoProxyHost: "",
		},
		{
			name:       "unmatched no_proxy routes through proxy",
			httpsProxy: proxyURL, noProxy: "other.example.com", protocol: "HTTPS", endpoint: host,
			wantUseProxy: true, wantProxyURL: proxyURL, wantNoProxyHost: "",
		},
		{
			name:       "HTTP protocol reads HTTP_PROXY",
			httpsProxy: "http://p.example.com:3128", noProxy: "", protocol: "HTTP", endpoint: host,
			wantUseProxy: true, wantProxyURL: "http://p.example.com:3128", wantNoProxyHost: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// env is process-global; subtests run sequentially to avoid races.
			t.Setenv("HTTPS_PROXY", tc.httpsProxy)
			t.Setenv("https_proxy", "")
			t.Setenv("HTTP_PROXY", tc.httpsProxy)
			t.Setenv("http_proxy", "")
			t.Setenv("NO_PROXY", tc.noProxy)
			t.Setenv("no_proxy", "")

			client := &AliyunClient{config: &Config{Protocol: tc.protocol}}
			useProxy, proxyURL, noProxyHost, err := client.resolveTeaProxy(tc.endpoint)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantUseProxy, useProxy)
			assert.Equal(t, tc.wantProxyURL, proxyURL)
			assert.Equal(t, tc.wantNoProxyHost, noProxyHost)
		})
	}
}

// TestUnitTeaRoaProxyRouting sends requests through the production constructors
// and the pinned SDK transport. The origin and proxy independently record where
// requests arrive; no proxy-selection algorithm or SDK transport is mocked.
func TestUnitTeaRoaProxyRouting(t *testing.T) {
	type observedRequest struct {
		method string
		host   string
		path   string
	}
	entries := []struct {
		name      string
		protocols []string
		request   func(*AliyunClient, string, *util.RuntimeOptions) (map[string]interface{}, error)
	}{
		{
			name: "OTS", protocols: []string{"HTTP", "HTTPS"},
			request: func(client *AliyunClient, endpoint string, runtime *util.RuntimeOptions) (map[string]interface{}, error) {
				conn, err := client.NewOtsRoaClient("ots")
				if err != nil {
					return nil, err
				}
				return conn.DoRequest(tea.String("2020-01-01"), nil, tea.String(http.MethodGet), tea.String("AK"), tea.String("/proxy-regression"), nil, nil, nil, runtime)
			},
		},
		{
			name: "common", protocols: []string{"HTTP", "HTTPS"},
			request: func(client *AliyunClient, endpoint string, runtime *util.RuntimeOptions) (map[string]interface{}, error) {
				conn, err := client.NewTeaRoaCommonClient(endpoint)
				if err != nil {
					return nil, err
				}
				return conn.DoRequestWithAction(tea.String("ProxyRegression"), tea.String("2020-01-01"), nil, tea.String(http.MethodGet), tea.String("AK"), tea.String("/proxy-regression"), nil, nil, nil, runtime)
			},
		},
		{
			// This wrapper owns its RuntimeOptions, so use HTTP here rather
			// than changing production TLS verification for a test certificate.
			name: "roaRequest", protocols: []string{"HTTP"},
			request: func(client *AliyunClient, endpoint string, runtime *util.RuntimeOptions) (map[string]interface{}, error) {
				return client.RoaDeleteWithApiName("cms", "2020-01-01", "ProxyRegression", "/proxy-regression", nil, nil, nil, false)
			},
		},
	}
	cases := []struct {
		name    string
		noProxy string
		proxy   bool
		bypass  bool
	}{
		{name: "suffix", noProxy: ".0.0.1", proxy: true, bypass: true},
		{name: "wildcard", noProxy: "*.0.0.1", proxy: true, bypass: true},
		{name: "exact hostname", noProxy: "127.0.0.1", proxy: true, bypass: true},
		{name: "star", noProxy: "*", proxy: true, bypass: true},
		{name: "nonmatching", noProxy: "unrelated.invalid", proxy: true},
		{name: "no proxy configured", bypass: true},
	}
	for _, entry := range entries {
		t.Run(entry.name, func(t *testing.T) {
			for _, protocol := range entry.protocols {
				t.Run(protocol, func(t *testing.T) {
					for _, tc := range cases {
						t.Run(tc.name, func(t *testing.T) {
							// Environment and SDK pools are process-wide. Keep cases
							// sequential and give each case a fresh origin/proxy port.
							var originCount, proxyCount, wrongProxyCount atomic.Int32
							originRequests := make(chan observedRequest, 1)
							proxyRequests := make(chan observedRequest, 1)
							handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
								originCount.Add(1)
								select {
								case originRequests <- observedRequest{r.Method, r.Host, r.URL.Path}:
								default:
								}
								w.Header().Set("Content-Type", "application/json")
								if err := json.NewEncoder(w).Encode(map[string]string{"route": "origin"}); err != nil {
									t.Error(err)
								}
							})
							var origin *httptest.Server
							if protocol == "HTTPS" {
								origin = httptest.NewTLSServer(handler)
							} else {
								origin = httptest.NewServer(handler)
							}
							defer origin.Close()
							endpoint := strings.TrimPrefix(strings.TrimPrefix(origin.URL, "https://"), "http://")

							proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
								proxyCount.Add(1)
								select {
								case proxyRequests <- observedRequest{r.Method, r.Host, r.URL.Path}:
								default:
								}
								if r.Method == http.MethodConnect {
									// Reproduce a corporate proxy refusing the tunnel.
									http.Error(w, "proxy rejected CONNECT", http.StatusForbidden)
									return
								}
								w.Header().Set("Content-Type", "application/json")
								if err := json.NewEncoder(w).Encode(map[string]string{"route": "proxy"}); err != nil {
									t.Error(err)
								}
							}))
							defer proxy.Close()
							wrongProxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
								wrongProxyCount.Add(1)
								http.Error(w, "wrong protocol proxy", http.StatusBadGateway)
							}))
							defer wrongProxy.Close()
							for _, key := range []string{"HTTP_PROXY", "http_proxy", "HTTPS_PROXY", "https_proxy", "NO_PROXY", "no_proxy"} {
								t.Setenv(key, "")
							}
							if tc.proxy {
								t.Setenv("HTTP_PROXY", wrongProxy.URL)
								t.Setenv("HTTPS_PROXY", wrongProxy.URL)
								t.Setenv(protocol+"_PROXY", proxy.URL)
							}
							t.Setenv("NO_PROXY", tc.noProxy)

							credential, err := credentials.NewCredential(new(credentials.Config).
								SetType("access_key").SetAccessKeyId("test-access-key").SetAccessKeySecret("test-secret-key"))
							if err != nil {
								t.Fatal(err)
							}
							endpoints := &sync.Map{}
							endpoints.Store("ots", endpoint)
							endpoints.Store("cms", endpoint)
							config := &Config{Protocol: protocol, RegionId: "cn-hangzhou", Credential: credential, Endpoints: endpoints,
								AccessKey: "test-access-key", SecretKey: "test-secret-key",
								ClientReadTimeout: 2000, ClientConnectTimeout: 2000}
							sdkConfig, err := config.getTeaRoaDslSdkConfig(true)
							if err != nil {
								t.Fatal(err)
							}
							client := &AliyunClient{config: config, teaRoaSdkConfig: sdkConfig}
							// Only the isolated httptest TLS certificate is untrusted.
							// Proxy and NoProxy runtime fields remain unset.
							runtime := &util.RuntimeOptions{IgnoreSSL: tea.Bool(protocol == "HTTPS"), Autoretry: tea.Bool(false)}
							response, err := entry.request(client, endpoint, runtime)
							if protocol == "HTTPS" && !tc.bypass {
								if err == nil || !strings.Contains(err.Error(), "Forbidden") {
									t.Fatalf("request error = %v, want proxy CONNECT Forbidden", err)
								}
							} else {
								if err != nil {
									t.Fatal(err)
								}
								body := response
								if value, ok := response["body"].(map[string]interface{}); ok {
									body = value
								}
								wantRoute := "proxy"
								if tc.bypass {
									wantRoute = "origin"
								}
								assert.Equal(t, wantRoute, body["route"])
							}
							assert.Zero(t, wrongProxyCount.Load(), "must use the proxy for the configured protocol")
							wantMethod := http.MethodGet
							if entry.name == "roaRequest" {
								wantMethod = http.MethodDelete
							}
							if tc.bypass {
								assert.EqualValues(t, 1, originCount.Load())
								assert.Zero(t, proxyCount.Load(), "bypassed request must not reach proxy")
								select {
								case request := <-originRequests:
									assert.Equal(t, observedRequest{wantMethod, endpoint, "/proxy-regression"}, request)
								default:
									t.Error("origin did not receive the SDK request")
								}
								select {
								case request := <-proxyRequests:
									t.Errorf("bypassed request reached proxy: %+v", request)
								default:
								}
							} else {
								assert.Zero(t, originCount.Load(), "proxied request must not reach origin")
								assert.EqualValues(t, 1, proxyCount.Load())
								select {
								case request := <-proxyRequests:
									assert.Equal(t, endpoint, request.host)
									if protocol == "HTTPS" {
										assert.Equal(t, http.MethodConnect, request.method)
									} else {
										assert.Equal(t, wantMethod, request.method)
										assert.Equal(t, "/proxy-regression", request.path)
									}
								default:
									t.Error("proxy did not receive the SDK request")
								}
							}
						})
					}
				})
			}
		})
	}
}
