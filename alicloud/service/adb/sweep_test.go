package adb_test

// sweep_test.go wires the migrated adb sweeper tests to the shared sweep
// helpers. The flat-suite sharedClientForRegion / sweepAll are unexported, but
// the sweep package now re-exports faithful equivalents: SharedClientForRegion
// builds the sweeper client (AccountId/AccountType/Credential and a
// per-region endpoint cache via SharedEndpointCacheForRegion), and SweepAll is
// the shared ALICLOUD_SWEEP_ALL_RESOURCES gate. Delegating to them keeps the
// migrated package behavior identical to the root suite instead of carrying a
// divergent local copy with a single global endpoint cache.
import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
)

// adbTestClientForRegion returns the shared sweeper client for a region.
// SharedClientForRegion returns the same *connectivity.AliyunClient the root
// suite builds, so the typed return is preserved for the migrated tests.
func adbTestClientForRegion(region string) (*connectivity.AliyunClient, error) {
	rawClient, err := sweep.SharedClientForRegion(region)
	if err != nil {
		return nil, err
	}
	return rawClient.(*connectivity.AliyunClient), nil
}

// adbSweepAll delegates to the shared sweep gate.
func adbSweepAll() bool {
	return sweep.SweepAll()
}
