package connectivity

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
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
	client := NewTestClient(t)

	testCases := []struct {
		productCode    string
		configEndpoint string
		RegionId       string
		expected       string
	}{
		{"sls", "", "cn-hangzhou", "cn-hangzhou.log.aliyuncs.com"},
		//{"slb", "", "cn-hangzhou", "slb.aliyuncs.com"},
		{"sls", "", "eu-central-1", "eu-central-1.log.aliyuncs.com"},
		//{"alb", "", "eu-central-1", "alb.eu-central-1.aliyuncs.com"},
	}

	for _, tc := range testCases {
		t.Run(tc.productCode, func(t *testing.T) {
			// In production one client serves a single region for its whole
			// lifetime, and the endpoint resolved on the first API call is
			// cached and reused by every following call. Each case must
			// therefore start with its own cache: reusing the shared
			// NewTestClient cache would serve the first case's endpoint to
			// the second case (same product, different region).
			client.config.Endpoints = new(sync.Map)
			client.config.Region = Region(tc.RegionId)
			client.config.RegionId = tc.RegionId
			client.RegionId = tc.RegionId

			endpoint, err := client.loadApiEndpoint(tc.productCode)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, endpoint)

			val, ok := client.config.Endpoints.Load(tc.productCode)
			assert.True(t, ok)
			assert.Equal(t, tc.expected, val)
		})
	}
}

// TestUnitSlsEndpointCacheReuse reproduces the read/write key mismatch fixed in
// loadApiEndpoint: for sls, productCodeToConfigEndpoints maps the read key to
// "log", while loadEndpoint always stores under "sls". Before the fix the
// mapped read key never matched the write key, so every SLS API call re-ran
// endpoint discovery. The env var is the only probe source, so this case needs
// no credentials and no network and runs under ci-check-quick.
func TestUnitSlsEndpointCacheReuse(t *testing.T) {
	const envKey = "ALIBABA_CLOUD_ENDPOINT_SLS"
	oldVal, hadOld := os.LookupEnv(envKey)
	t.Cleanup(func() {
		if hadOld {
			os.Setenv(envKey, oldVal)
		} else {
			os.Unsetenv(envKey)
		}
	})

	client := &AliyunClient{
		config: &Config{
			Endpoints: new(sync.Map),
			RegionId:  "cn-hangzhou",
		},
	}

	// First call probes the env var and must cache under "sls", not "log".
	os.Setenv(envKey, "first.sls.example.com")
	ep1, err := client.loadApiEndpoint("sls")
	assert.NoError(t, err)
	assert.Equal(t, "first.sls.example.com", ep1)

	// Bug symptom: the read key is "log" but the write key is "sls", so the
	// user-config key stays empty after a probe.
	_, ok := client.config.Endpoints.Load("log")
	assert.False(t, ok, `Load("log") must be empty; the probe is cached under "sls"`)

	cached, ok := client.config.Endpoints.Load("sls")
	assert.True(t, ok)
	assert.Equal(t, "first.sls.example.com", cached)

	// Flip the probe source. Before the fix the second call re-probed and
	// returned the new value; with the Load(productCode) pre-check it is served
	// from the cache instead.
	os.Setenv(envKey, "second.sls.example.com")
	ep2, err := client.loadApiEndpoint("sls")
	assert.NoError(t, err)
	assert.Equal(t, "first.sls.example.com", ep2, "second call must reuse the cached endpoint instead of re-probing")

	cached, ok = client.config.Endpoints.Load("sls")
	assert.True(t, ok)
	assert.Equal(t, "first.sls.example.com", cached, "cache must not be overwritten by the skipped probe")
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

// TestUnitNonMappedProductEndpointCacheReuse is the control group: for
// non-mapped products the read and write keys already coincide (ecs), so
// caching must keep working unchanged after the fix.
func TestUnitNonMappedProductEndpointCacheReuse(t *testing.T) {
	const envKey = "ALIBABA_CLOUD_ENDPOINT_ECS"
	oldVal, hadOld := os.LookupEnv(envKey)
	t.Cleanup(func() {
		if hadOld {
			os.Setenv(envKey, oldVal)
		} else {
			os.Unsetenv(envKey)
		}
	})

	client := &AliyunClient{
		config: &Config{
			Endpoints: new(sync.Map),
			RegionId:  "cn-hangzhou",
		},
	}

	os.Setenv(envKey, "first.ecs.example.com")
	ep1, err := client.loadApiEndpoint("ecs")
	assert.NoError(t, err)
	assert.Equal(t, "first.ecs.example.com", ep1)

	os.Setenv(envKey, "second.ecs.example.com")
	ep2, err := client.loadApiEndpoint("ecs")
	assert.NoError(t, err)
	assert.Equal(t, "first.ecs.example.com", ep2, "second call must reuse the cached ecs endpoint instead of re-probing")

	cached, ok := client.config.Endpoints.Load("ecs")
	assert.True(t, ok)
	assert.Equal(t, "first.ecs.example.com", cached)
}
