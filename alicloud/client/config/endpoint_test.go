package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
