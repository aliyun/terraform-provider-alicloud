package scdn_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

// scdnTestClientForRegion delegates to the exported sweep.SharedClientForRegion
// so the sweeper client carries the same ALICLOUD_ACCOUNT_ID/ALICLOUD_ACCOUNT_TYPE
// environment variables and region-isolated endpoint cache
// (SharedEndpointCacheForRegion) as the flat-suite sharedClientForRegion in
// alicloud_sweeper_test.go, instead of a divergent per-package copy that omits
// the account metadata and shares endpoints across regions. The returned
// client is the concrete *connectivity.AliyunClient used by the sweeper helpers.
func scdnTestClientForRegion(region string) (*connectivity.AliyunClient, error) {
	client, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return nil, err
	}
	scdnClient, ok := client.(*connectivity.AliyunClient)
	if !ok {
		return nil, fmt.Errorf("sweep.SharedClientForRegion(%q) returned unexpected client type %T", region, client)
	}
	return scdnClient, nil
}

// scdnTestRegion mirrors the flat-suite defaultRegionToTest() call through the
// exported acctest value, so test names carry the resolved region.
func scdnTestRegion() string {
	return acctest.DefaultRegionToTest
}
