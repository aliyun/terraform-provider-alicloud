package connectivity

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"sync"
	"testing"
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

func TestUnitSlsWithoutEndpoint(t *testing.T) {
	client := NewTestClient(t)

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

	for _, tc := range testCases {
		t.Run(tc.productCode, func(t *testing.T) {
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
