package connectivity

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUnitCms2EndpointTableEntry verifies that the CloudMonitor 2.0 (CMS2, API
// version 2024-03-30) product code is registered in the regularProductEndpoint
// table with the dedicated CMS 2.0 regional domain `cms.%s.aliyuncs.com` rather
// than the legacy CMS 1.0 `metrics.%s.aliyuncs.com` domain.
func TestUnitCms2EndpointTableEntry(t *testing.T) {
	endpointFmt, ok := regularProductEndpoint["cms2"]
	assert.True(t, ok, "cms2 product code must be registered in regularProductEndpoint")
	assert.Equal(t, "cms.%s.aliyuncs.com", endpointFmt)

	assert.Equal(t, "cms.cn-hangzhou.aliyuncs.com", fmt.Sprintf(endpointFmt, "cn-hangzhou"))
	assert.Equal(t, "cms.us-east-1.aliyuncs.com", fmt.Sprintf(endpointFmt, "us-east-1"))
}

// TestUnitCms2EndpointDomainDistinctFromCms1 guards against the regression where
// CMS 2.0 requests route to the CMS 1.0 domain: cms and cms2 must resolve to
// different regional domains so the dedicated 2.0 endpoint is reachable.
func TestUnitCms2EndpointDomainDistinctFromCms1(t *testing.T) {
	region := "cn-hangzhou"
	cms1 := fmt.Sprintf(regularProductEndpoint["cms"], region)
	cms2 := fmt.Sprintf(regularProductEndpoint["cms2"], region)
	assert.NotEqual(t, cms1, cms2)
	assert.Equal(t, "metrics.cn-hangzhou.aliyuncs.com", cms1)
	assert.Equal(t, "cms.cn-hangzhou.aliyuncs.com", cms2)
}

// TestUnitCms2EndpointOverrideFromConfig verifies that the provider `cms2` endpoint
// override (stored in config.Endpoints by the provider endpoints loop) takes
// precedence over the regional default, so users can point CMS2 resources at a
// custom CloudMonitor 2.0 endpoint.
func TestUnitCms2EndpointOverrideFromConfig(t *testing.T) {
	client := &AliyunClient{
		config: &Config{
			Endpoints: new(sync.Map),
			RegionId:  "cn-hangzhou",
		},
	}
	client.config.Endpoints.Store("cms2", "https://cms2.custom.example.com")

	endpoint, err := client.loadApiEndpoint("cms2")
	assert.NoError(t, err)
	assert.Equal(t, "https://cms2.custom.example.com", endpoint)
}
