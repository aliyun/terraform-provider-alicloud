package cloudstoragegateway_test

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
)

// csgTestClientForRegion delegates to the shared sweep client so the
// package-level sweepers reuse the same connectivity client as the rest of
// the provider test suite.
func csgTestClientForRegion(region string) (interface{}, error) {
	return sweep.SharedClientForRegion(region)
}
