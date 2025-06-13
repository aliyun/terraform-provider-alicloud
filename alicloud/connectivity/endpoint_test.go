package connectivity

import (
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestLoadEndpointFromEnv(t *testing.T) {
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

func TestLoadEndpointFromLocal(t *testing.T) {
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

func TestIrregularProductEndpoint(t *testing.T) {
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

func TestInternationalRegionEndpoint(t *testing.T) {
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

func TestEndpointErrorHandling(t *testing.T) {
	client := NewTestClient(t)

	err := client.loadEndpoint("invalid_product")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Illegal")
}
