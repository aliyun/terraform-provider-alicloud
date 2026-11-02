// Package smartag_test.
package smartag_test

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
)

// smartagTestClientForRegion thin-delegates to the shared sweep.SharedClientForRegion,
// which isolates the endpoint cache per region via SharedEndpointCacheForRegion and
// retains signVersion, AK/STS and AccountId/AccountType. The shared helper is reachable
// from this external test package, so the client is no longer rebuilt inline against a
// package-level global endpoint cache (which would be shared across regions and risk
// cross-region routing during multi-region sweepers). Return type, error handling and
// the *connectivity.AliyunClient assertion are preserved.
func smartagTestClientForRegion(region string) (*connectivity.AliyunClient, error) {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return nil, err
	}
	return rawClient.(*connectivity.AliyunClient), nil
}

// smartagTestSweepPreCheckWithRegions mirrors the root-package
// testSweepPreCheckWithRegions gate; it is not reachable from a service test
// package.
func smartagTestSweepPreCheckWithRegions(region string, supported bool, regions []connectivity.Region) bool {
	find := false
	for _, r := range regions {
		if region == string(r) {
			find = true
			break
		}
	}
	return (find && !supported) || (!find && supported)
}
