// Package privatezone_test hosts the shared sweeper/precheck helpers that
// previously lived in the flat alicloud test suite. The root-package
// sharedClientForRegion is not reachable from an external service test
// package, so the sweeper client is rebuilt here as a thin delegate over the
// exported alicloud/sweep.SharedClientForRegion. That shared helper keeps the
// per-region endpoint cache (SharedEndpointCacheForRegion), the shared
// signVersion, the AK/STS credentials and the ALICLOUD_/ALIBABA_CLOUD_
// account-id and account-type fallbacks, matching the merge-base root helper
// and the shared sweep helper; a per-package global cache must not be used,
// because sharing one endpoint cache across regions routes a Hangzhou request
// to a previously cached Beijing endpoint. The per-test region gate uses
// acctest.PreCheckWithRegions with the exported connectivity.TestPvtzRegions
// directly, so no local region-gate mirror is needed.
package privatezone_test

import (
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
)

// privatezoneTestClientForRegion returns the package-level sweeper client by
// delegating to alicloud/sweep.SharedClientForRegion, which isolates the
// endpoint cache per region and keeps the account/credential fallbacks. It
// preserves the original typed (*connectivity.AliyunClient, error) contract:
// on error it returns (nil, err), otherwise it type-asserts the shared client
// to *connectivity.AliyunClient.
func privatezoneTestClientForRegion(region string) (*connectivity.AliyunClient, error) {
	client, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return nil, err
	}
	return client.(*connectivity.AliyunClient), nil
}

// privatezoneSweepAll mirrors the root-package sweepAll gate; it is not
// reachable from a service test package.
func privatezoneSweepAll() bool {
	return os.Getenv("ALICLOUD_SWEEP_ALL_RESOURCES") == "true"
}

// privatezoneTestPreCheckEnterpriseAccountEnabled mirrors the root-package
// testAccPreCheckEnterpriseAccountEnabled gate; it is not reachable from a
// service test package.
func privatezoneTestPreCheckEnterpriseAccountEnabled(t *testing.T) {
	enable := os.Getenv("ENTERPRISE_ACCOUNT_ENABLED")
	if enable != "true" {
		t.Skipf("Skipping the test case because the enterprise account is not enabled.")
		t.Skipped()
	}
}
