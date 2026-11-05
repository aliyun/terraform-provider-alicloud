package ram_test

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
)

// ramTestClientForRegion delegates to the shared alicloud/sweep.SharedClientForRegion
// so the migrated ram service test package reuses the same account metadata
// (ALICLOUD_ACCOUNT_ID/ALIBABA_CLOUD_ACCOUNT_ID, ALICLOUD_ACCOUNT_TYPE/
// ALIBABA_CLOUD_ACCOUNT_TYPE) and per-region product endpoint cache
// (SharedEndpointCacheForRegion) as the root package, instead of a lossy local
// copy that omitted the account fields and shared a single global endpoint map
// across regions. The (interface{}, error) signature is preserved, so callers
// keep their rawClient.(*connectivity.AliyunClient) type assertion and their own
// error/Skip and cleanup logic unchanged.
func ramTestClientForRegion(region string) (interface{}, error) {
	return sweep.SharedClientForRegion(region)
}
