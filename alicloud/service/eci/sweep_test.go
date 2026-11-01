// Package eci_test sweep helpers. The ECI resources were migrated out of the
// flat alicloud test suite, so the package-private sweeper client and region
// gates (sharedClientForRegion / testSweepPreCheckWithRegions) are no longer
// reachable. These helpers delegate to the exported sweep package, mirroring
// the actiontrail service-test pattern.
package eci_test

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
)

// eciTestClientForRegion delegates to the shared sweeper client. The migrated
// service test package cannot reach the package-private sharedClientForRegion;
// sweep.SharedClientForRegion is its per-region, account-aware equivalent: it
// isolates the endpoint cache per region via SharedEndpointCacheForRegion and
// reads ALICLOUD_ACCOUNT_ID / ALICLOUD_ACCOUNT_TYPE, matching the current root
// helper. Thin delegation avoids reintroducing a single global endpoint cache
// shared across all regions.
func eciTestClientForRegion(region string) (interface{}, error) {
	return sweep.SharedClientForRegion(region)
}

// eciTestSweepPreCheckWithRegions mirrors the root-package
// testSweepPreCheckWithRegions gate; it is not reachable from a service test
// package.
func eciTestSweepPreCheckWithRegions(region string, supported bool, regions []connectivity.Region) bool {
	find := false
	for _, r := range regions {
		if region == string(r) {
			find = true
			break
		}
	}
	return (find && !supported) || (!find && supported)
}
