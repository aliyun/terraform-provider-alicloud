package nas_test

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
)

// nasTestClientForRegion delegates to the public sweep.SharedClientForRegion,
// preserving the baseline per-region endpoint isolation (sweep.SharedEndpointCacheForRegion)
// instead of a lossy package-level global cache shared across regions. The sweep
// package is the public correspondent of the root-package sweeper helpers; it only
// imports connectivity and does not reverse-import the root or nas package.
func nasTestClientForRegion(region string) (interface{}, error) {
	return sweep.SharedClientForRegion(region)
}

func nasSweepAll() bool {
	return sweep.SweepAll()
}

func nasTestSweepPreCheckWithRegions(region string, supported bool, regions []connectivity.Region) bool {
	return sweep.SweepPreCheckWithRegions(region, supported, regions)
}
