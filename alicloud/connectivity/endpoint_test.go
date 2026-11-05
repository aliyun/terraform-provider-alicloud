package connectivity

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitCommonLoadEndpointFromEnv(t *testing.T) {
	t.Setenv("ALIBABA_CLOUD_ENDPOINT_ECS", "https://ecs.test.com")
	defer os.Unsetenv("ALIBABA_CLOUD_ENDPOINT_ECS")

	client := &AliyunClient{
		config: &Config{
			Endpoints: new(sync.Map),
			RegionId:  "cn-beijing",
		},
	}

	err := client.loadEndpoint("ecs")
	assert.NoError(t, err)

	val, ok := client.config.Endpoints.Load("ecs")
	assert.True(t, ok)
	assert.Equal(t, "https://ecs.test.com", val)
}

func TestUnitCommonLoadEndpointFromLocal(t *testing.T) {
	content := `<?xml version="1.0" encoding="UTF-8"?>
<Endpoints>
  <Endpoint name="cn-beijing">
    <RegionIds>
      <RegionId>cn-beijing</RegionId>
    </RegionIds>
    <Products>
      <Product>
        <ProductName>ECS</ProductName>
        <DomainName>ecs.cn-beijing.aliyuncs.com</DomainName>
      </Product>
    </Products>
  </Endpoint>
</Endpoints>`

	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "endpoints.xml")
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	assert.NoError(t, err)

	t.Setenv("TF_ENDPOINT_PATH", tmpFile)
	defer os.Unsetenv("TF_ENDPOINT_PATH")

	config := &Config{
		Endpoints: new(sync.Map),
		RegionId:  "cn-beijing",
	}

	err = config.loadEndpointFromLocal()
	assert.NoError(t, err)

	val, ok := config.Endpoints.Load("ecs")
	assert.True(t, ok)
	assert.Equal(t, "ecs.cn-beijing.aliyuncs.com", val)
}

func TestUnitCommonIrregularProductEndpoint(t *testing.T) {
	client := &AliyunClient{
		RegionId: "cn-hangzhou",
		config: &Config{
			Endpoints: new(sync.Map),
			RegionId:  "cn-hangzhou",
		},
	}

	testCases := []struct {
		productCode string
		expected    string
	}{
		{"ram", "ram.aliyuncs.com"},
		{"cloudfw", "cloudfw.aliyuncs.com"},
		{"agentloop", "agentloop.cn-hangzhou.aliyuncs.com"},
	}

	for _, tc := range testCases {
		t.Run(tc.productCode, func(t *testing.T) {
			err := client.loadEndpoint(tc.productCode)
			assert.NoError(t, err)

			val, ok := client.config.Endpoints.Load(tc.productCode)
			assert.True(t, ok)
			assert.Equal(t, tc.expected, val)
		})
	}
}

func TestUnitCommonInternationalRegionEndpoint(t *testing.T) {
	client := &AliyunClient{
		config: &Config{
			Endpoints: new(sync.Map),
			RegionId:  "ap-southeast-1",
		},
	}

	testCases := []struct {
		productCode string
		expected    string
	}{
		{"cloudfw", "cloudfw.aliyuncs.com"},
		{"sas", "tds.ap-southeast-1.aliyuncs.com"},
	}

	for _, tc := range testCases {
		t.Run(tc.productCode, func(t *testing.T) {
			err := client.loadEndpoint(tc.productCode)
			assert.NoError(t, err)

			val, ok := client.config.Endpoints.Load(tc.productCode)
			assert.True(t, ok)
			assert.Equal(t, tc.expected, val)
		})
	}
}

func TestUnitCommonEndpointErrorHandling(t *testing.T) {
	client := NewTestClient(t)

	err := client.loadEndpoint("invalid_product")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Illegal")
}

func TestUnitSlsEndpoint(t *testing.T) {
	var endpoints sync.Map
	endpoints.Store("log", "cn-shanghai.log.aliyuncs.com")

	client := NewTestClient(t)
	client.config.Region = Hangzhou
	client.config.RegionId = "cn-hangzhou"
	client.config.Endpoints = &endpoints

	testCases := []struct {
		productCode    string
		configEndpoint string
		expected       string
	}{
		{"sls", "log", "cn-shanghai.log.aliyuncs.com"},
		{"alb", "alb", "alb.cn-hangzhou.aliyuncs.com"},
	}

	for _, tc := range testCases {
		t.Run(tc.productCode, func(t *testing.T) {
			endpoint, err := client.loadApiEndpoint(tc.productCode)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, endpoint)

			val, ok := client.config.Endpoints.Load(tc.configEndpoint)
			assert.True(t, ok)
			assert.Equal(t, tc.expected, val)
		})
	}
}

func TestUnitEfloControllerEndpoint(t *testing.T) {
	// The location service returns InvalidRegionId for efloctrl in some
	// regions (e.g. us-southeast-1). regularProductEndpoint acts as the
	// backup endpoint so eflo_controller resolves to
	// eflo-controller.<region>.aliyuncs.com instead of failing.
	endpointFmt, ok := regularProductEndpoint["eflo_controller"]
	assert.True(t, ok)
	assert.Equal(t, "eflo-controller.%s.aliyuncs.com", endpointFmt)

	for _, region := range []string{"us-southeast-1", "cn-hangzhou", "cn-beijing"} {
		assert.Equal(t, fmt.Sprintf("eflo-controller.%s.aliyuncs.com", region), fmt.Sprintf(endpointFmt, region))
	}
}

func TestUnitEfloControllerEndpointFromEnv(t *testing.T) {
	// Explicit provider endpoint config and environment variables take
	// precedence over the regularProductEndpoint fallback.
	t.Setenv("ALIBABA_CLOUD_ENDPOINT_EFLO_CONTROLLER", "https://eflo-controller.test.com")
	defer os.Unsetenv("ALIBABA_CLOUD_ENDPOINT_EFLO_CONTROLLER")

	client := &AliyunClient{
		config: &Config{
			Endpoints: new(sync.Map),
			RegionId:  "us-southeast-1",
		},
	}

	err := client.loadEndpoint("eflo_controller")
	assert.NoError(t, err)

	val, ok := client.config.Endpoints.Load("eflo_controller")
	assert.True(t, ok)
	assert.Equal(t, "https://eflo-controller.test.com", val)
}

func TestUnitSlsWithoutEndpoint(t *testing.T) {
	testCases := []struct {
		productCode    string
		configEndpoint string
		RegionId       string
		expected       string
	}{
		{"sls", "", "cn-hangzhou", "cn-hangzhou.log.aliyuncs.com"},
		{"slb", "", "cn-hangzhou", "slb.aliyuncs.com"},
		{"sls", "", "eu-central-1", "eu-central-1.log.aliyuncs.com"},
		{"alb", "", "eu-central-1", "alb.eu-central-1.aliyuncs.com"},
	}

	client := NewTestClient(t)

	for _, tc := range testCases {
		t.Run(tc.productCode, func(t *testing.T) {
			// A fresh endpoint cache keeps the endpoint discovered for another
			// region out of this case.
			client.config.Endpoints = new(sync.Map)
			client.config.Region = Region(tc.RegionId)
			client.config.RegionId = tc.RegionId

			endpoint, err := client.loadApiEndpoint(tc.productCode)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, endpoint)

			val, ok := client.config.Endpoints.Load(tc.productCode)
			assert.True(t, ok)
			assert.Equal(t, tc.expected, val)
		})
	}
}

func TestUnitSlsEndpointCacheReuse(t *testing.T) {
	// Mirror the reported scenario: the location service is unreachable, so the
	// first call falls back to the regular endpoint form and logs a WARN. Every
	// following call must be served from the endpoint cache, without probing the
	// location service again and without logging another WARN.
	var locationProbes int32
	locationProbe := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&locationProbes, 1)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer locationProbe.Close()

	var logs bytes.Buffer
	previousLogOutput := log.Writer()
	log.SetOutput(&logs)
	defer log.SetOutput(previousLogOutput)

	t.Setenv("ALIBABA_CLOUD_ENDPOINT_SLS", "")
	t.Setenv("SLS_ENDPOINT", "")

	client := &AliyunClient{
		RegionId: "cn-hangzhou",
		config: &Config{
			Endpoints: new(sync.Map),
			RegionId:  "cn-hangzhou",
			// HTTP keeps the simulated unreachable-location failure out of the
			// SDK's HTTPS-specific retry branch, which would take minutes.
			Protocol:         "HTTP",
			LocationEndpoint: locationProbe.Listener.Addr().String(),
			AccessKey:        "test-ak",
			SecretKey:        "test-sk",
			Credential:       &mockCredential{err: fmt.Errorf("unit test uses the static access key")},
		},
	}

	first, err := client.loadApiEndpoint("sls")
	assert.NoError(t, err)
	assert.Equal(t, "cn-hangzhou.log.aliyuncs.com", first)
	assert.NotZero(t, atomic.LoadInt32(&locationProbes))

	cached, ok := client.config.Endpoints.Load("sls")
	assert.True(t, ok)
	assert.Equal(t, first, cached)
	_, logCached := client.config.Endpoints.Load("log")
	assert.False(t, logCached)
	assert.Equal(t, 1, strings.Count(logs.String(), "[WARN] loading sls endpoint got an error"))

	// Change the probe source. An un-cached implementation would return the new
	// value here; a cache hit keeps the first endpoint.
	t.Setenv("ALIBABA_CLOUD_ENDPOINT_SLS", "https://sls.second.test.com")
	probesAfterFirstCall := atomic.LoadInt32(&locationProbes)

	second, err := client.loadApiEndpoint("sls")
	assert.NoError(t, err)
	assert.Equal(t, first, second)
	assert.Equal(t, probesAfterFirstCall, atomic.LoadInt32(&locationProbes))
	assert.Equal(t, 1, strings.Count(logs.String(), "[WARN] loading sls endpoint got an error"))

	// Clear the probe source again. With the endpoint cached nothing may probe
	// the location service or log another WARN here; without the cache this call
	// re-runs discovery and breaks both counters, while still returning the same
	// regular endpoint form, so the counters are what distinguish the two paths.
	t.Setenv("ALIBABA_CLOUD_ENDPOINT_SLS", "")

	third, err := client.loadApiEndpoint("sls")
	assert.NoError(t, err)
	assert.Equal(t, first, third)
	assert.Equal(t, probesAfterFirstCall, atomic.LoadInt32(&locationProbes))
	assert.Equal(t, 1, strings.Count(logs.String(), "[WARN] loading sls endpoint got an error"))
}

func TestUnitSlsConfiguredEndpointPriority(t *testing.T) {
	t.Setenv("ALIBABA_CLOUD_ENDPOINT_SLS", "https://sls.env.test.com")

	var endpoints sync.Map
	endpoints.Store("log", "cn-shanghai.log.aliyuncs.com")

	client := &AliyunClient{
		config: &Config{
			Endpoints: &endpoints,
			RegionId:  "cn-hangzhou",
		},
	}

	endpoint, err := client.loadApiEndpoint("sls")
	assert.NoError(t, err)
	assert.Equal(t, "cn-shanghai.log.aliyuncs.com", endpoint)

	cached, ok := client.config.Endpoints.Load("log")
	assert.True(t, ok)
	assert.Equal(t, "cn-shanghai.log.aliyuncs.com", cached)
}

func TestUnitProductEndpointCacheReuse(t *testing.T) {
	testCases := []struct {
		productCode string
		envName     string
	}{
		{"ecs", "ALIBABA_CLOUD_ENDPOINT_ECS"},
		{"alb", "ALIBABA_CLOUD_ENDPOINT_ALB"},
	}

	for _, tc := range testCases {
		t.Run(tc.productCode, func(t *testing.T) {
			firstName := fmt.Sprintf("https://%s.first.test.com", tc.productCode)
			secondName := fmt.Sprintf("https://%s.second.test.com", tc.productCode)
			t.Setenv(tc.envName, firstName)

			client := &AliyunClient{
				config: &Config{
					Endpoints: new(sync.Map),
					RegionId:  "cn-hangzhou",
				},
			}

			first, err := client.loadApiEndpoint(tc.productCode)
			assert.NoError(t, err)
			assert.Equal(t, firstName, first)

			cached, ok := client.config.Endpoints.Load(tc.productCode)
			assert.True(t, ok)
			assert.Equal(t, firstName, cached)

			t.Setenv(tc.envName, secondName)

			second, err := client.loadApiEndpoint(tc.productCode)
			assert.NoError(t, err)
			assert.Equal(t, first, second)

			cached, ok = client.config.Endpoints.Load(tc.productCode)
			assert.True(t, ok)
			assert.Equal(t, first, cached)
		})
	}
}

// TestUnitSlsEndpointUserConfigTakesPrecedenceOverCache verifies the fix keeps
// the v1.253.0 user-config priority: endpoints.log for sls wins over a
// previously probed endpoint cached under "sls".
func TestUnitSlsEndpointUserConfigTakesPrecedenceOverCache(t *testing.T) {
	var endpoints sync.Map
	endpoints.Store("log", "user.log.example.com")
	endpoints.Store("sls", "cached.sls.example.com")

	client := &AliyunClient{
		config: &Config{
			Endpoints: &endpoints,
			RegionId:  "cn-hangzhou",
		},
	}

	ep, err := client.loadApiEndpoint("sls")
	assert.NoError(t, err)
	assert.Equal(t, "user.log.example.com", ep, "user-configured log endpoint must take precedence over the cached sls endpoint")

	val, ok := client.config.Endpoints.Load("log")
	assert.True(t, ok)
	assert.Equal(t, "user.log.example.com", val)

	cached, ok := client.config.Endpoints.Load("sls")
	assert.True(t, ok)
	assert.Equal(t, "cached.sls.example.com", cached, "the cached sls entry must be left untouched")
}
