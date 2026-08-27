package connectivity

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"sync"
	"testing"
	"time"

	ossgateway "github.com/alibabacloud-go/alibabacloud-gateway-oss/client"
	spi "github.com/alibabacloud-go/alibabacloud-gateway-spi/client"
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
