package connectivity

import (
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	otsTunnel "github.com/aliyun/aliyun-tablestore-go-sdk/tunnel"
	"github.com/aliyun/credentials-go/credentials"
	"github.com/stretchr/testify/assert"
)

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

// otsProxyTestEndpoint mimics a cn-hangzhou TableStore data-plane endpoint.
const otsProxyTestEndpoint = "https://my-inst.cn-hangzhou.ots.aliyuncs.com"

// TestUnitGetOtsProxy verifies the OTS data-plane proxy resolution: honor HTTP(S)_PROXY,
// skip when NO_PROXY matches, and stay direct when no proxy is configured.
func TestUnitGetOtsProxy(t *testing.T) {
	client := &AliyunClient{config: &Config{Protocol: "HTTPS"}}
	cases := []struct {
		name       string
		httpsProxy string
		noProxy    string
		wantProxy  string // "" means direct (nil)
	}{
		{"no proxy -> direct", "", "", ""},
		{"proxy set -> via proxy", "https://de.coia.siemens.net:9400", "", "https://de.coia.siemens.net:9400"},
		{"no_proxy matches -> direct", "https://de.coia.siemens.net:9400", "ots.aliyuncs.com", ""},
		{"no_proxy non-match -> via proxy", "https://de.coia.siemens.net:9400", "example.com", "https://de.coia.siemens.net:9400"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("HTTPS_PROXY", tc.httpsProxy)
			t.Setenv("NO_PROXY", tc.noProxy)
			proxy, err := client.getOtsProxy(otsProxyTestEndpoint)
			assert.NoError(t, err)
			if tc.wantProxy == "" {
				assert.Nil(t, proxy)
			} else {
				assert.NotNil(t, proxy)
				assert.Equal(t, tc.wantProxy, proxy.String())
			}
		})
	}
}

// TestUnitWithTableStoreClient_Proxy verifies the table-store data-plane transport routes
// through the proxy when set and stays direct otherwise.
func TestUnitWithTableStoreClient_Proxy(t *testing.T) {
	client := &AliyunClient{config: &Config{Protocol: "HTTPS"}}
	req, _ := http.NewRequest("POST", otsProxyTestEndpoint, nil)

	t.Run("with proxy", func(t *testing.T) {
		t.Setenv("HTTPS_PROXY", "https://de.coia.siemens.net:9400")
		t.Setenv("NO_PROXY", "")
		cfg := tablestore.NewDefaultTableStoreConfig()
		proxy, err := client.getOtsProxy(otsProxyTestEndpoint)
		assert.NoError(t, err)
		assert.NotNil(t, proxy)
		cfg.Transport = &http.Transport{Proxy: http.ProxyURL(proxy)}
		used, err := cfg.Transport.(*http.Transport).Proxy(req)
		assert.NoError(t, err)
		assert.Equal(t, "https://de.coia.siemens.net:9400", used.String())
	})

	t.Run("without proxy -> direct", func(t *testing.T) {
		t.Setenv("HTTPS_PROXY", "")
		t.Setenv("NO_PROXY", "")
		proxy, err := client.getOtsProxy(otsProxyTestEndpoint)
		assert.NoError(t, err)
		assert.Nil(t, proxy) // transport untouched -> SDK default direct connection
	})
}

// TestUnitWithTableStoreTunnelClient_Proxy verifies the tunnel data-plane transport routes
// through the proxy when set and stays on the SDK default otherwise.
func TestUnitWithTableStoreTunnelClient_Proxy(t *testing.T) {
	client := &AliyunClient{config: &Config{Protocol: "HTTPS"}}
	req, _ := http.NewRequest("POST", otsProxyTestEndpoint, nil)

	t.Run("with proxy", func(t *testing.T) {
		t.Setenv("HTTPS_PROXY", "https://de.coia.siemens.net:9400")
		t.Setenv("NO_PROXY", "")
		proxy, err := client.getOtsProxy(otsProxyTestEndpoint)
		assert.NoError(t, err)
		cfg := *otsTunnel.DefaultTunnelConfig
		cfg.Transport = &http.Transport{Proxy: http.ProxyURL(proxy)}
		used, err := cfg.Transport.(*http.Transport).Proxy(req)
		assert.NoError(t, err)
		assert.Equal(t, "https://de.coia.siemens.net:9400", used.String())
	})

	t.Run("without proxy -> direct", func(t *testing.T) {
		t.Setenv("HTTPS_PROXY", "")
		t.Setenv("NO_PROXY", "")
		proxy, err := client.getOtsProxy(otsProxyTestEndpoint)
		assert.NoError(t, err)
		assert.Nil(t, proxy) // tunnel keeps DefaultTunnelConfig (no forced proxy)
	})
}
