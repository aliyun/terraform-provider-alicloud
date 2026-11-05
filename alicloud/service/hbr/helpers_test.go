// Package hbr_test.
package hbr_test

import (
	"os"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
)

// hbrTestClientForRegion thin-delegates to the package-level sweep helper
// sweep.SharedClientForRegion. The root client helper (sharedClientForRegion)
// is not reachable from this external test package, but the migrated
// sweep.SharedClientForRegion is: it isolates the SDK endpoint cache per
// region (SharedEndpointCacheForRegion) and keeps the shared signVersion,
// AK/STS, AccountId/AccountType setup. Sharing one endpoint cache across
// regions can route a Hangzhou request to a previously cached Beijing
// endpoint, so per-region isolation must be preserved. It returns interface{}
// to mirror the root helper, so the standard
// `rawClient = rawClient.(*connectivity.AliyunClient)` type assertion at each
// call site stays valid.
func hbrTestClientForRegion(region string) (interface{}, error) {
	return sweep.SharedClientForRegion(region)
}

// hbrSweepAll mirrors the root-package sweepAll gate; it is not
// reachable from a service test package.
func hbrSweepAll() bool {
	return os.Getenv("ALICLOUD_SWEEP_ALL_RESOURCES") == "true"
}

// hbrTestSweepPreCheckWithRegions mirrors the root-package
// testSweepPreCheckWithRegions gate; it is not reachable from a service test
// package.
func hbrTestSweepPreCheckWithRegions(region string, supported bool, regions []connectivity.Region) bool {
	find := false
	for _, r := range regions {
		if region == string(r) {
			find = true
			break
		}
	}
	return (find && !supported) || (!find && supported)
}
