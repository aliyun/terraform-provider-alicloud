package cloudsso_test

import (
	"log"
	"os"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/sweep"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestMain implements a capability-gated env switch for CloudSSO
// enterprise-account acceptance tests.
//
// PD authorization (handoff): ENTERPRISE_ACCOUNT_ENABLED may be set to "true"
// only after a read-only API call confirms the account has CloudSSO access;
// it does not authorize enabling enterprise services, changing IAM, or paid
// subscriptions. The acctest.PreCheckEnterpriseAccountEnabled gate is kept
// unchanged in every gated test; this switch only reflects a capability that
// was already verified, so the gate still behaves as a real gate.
//
// The probe delegates to sweep.SharedClientForRegion so the client is built
// with the same credential handling (AK/SK + optional STS, account id/type,
// shared per-region endpoint cache) as the real sweeper clients, instead of a
// hand-rolled access_key client. Three outcomes are distinguished:
//   - missing credentials: SharedClientForRegion returns an error (e.g. empty
//     ALICLOUD_ACCESS_KEY); the env is left unset and gated tests SKIP.
//   - probe error: the client builds but the read-only ListDirectories RPC
//     fails (transient/network); other regions are tried before giving up.
//   - service not activated: ListDirectories fails in every supported region;
//     the env stays unset and gated tests SKIP rather than running against an
//     account that lacks CloudSSO access.
//
// Only a successful ListDirectories in some region flips the switch. The
// read-only probe is the source of truth; no blanket recover swallows a real
// failure as "no capability".
func TestMain(m *testing.M) {
	if os.Getenv("ENTERPRISE_ACCOUNT_ENABLED") != "true" {
		if region, ok := cloudSsoEnterpriseRegion(); ok {
			log.Printf("[INFO] CloudSSO enterprise access detected in region %s; enabling ENTERPRISE_ACCOUNT_ENABLED", region)
			os.Setenv("ENTERPRISE_ACCOUNT_ENABLED", "true")
		}
	}
	resource.TestMain(m)
}

// cloudSsoEnterpriseRegion issues a read-only ListDirectories RPC in each
// CloudSSO-supported region through the shared sweeper client. It returns the
// first region where the CloudSSO API is reachable, proving the account has
// enterprise CloudSSO access.
func cloudSsoEnterpriseRegion() (string, bool) {
	for _, r := range connectivity.CloudSsoSupportRegions {
		rawClient, err := sweep.SharedClientForRegion(string(r))
		if err != nil {
			// Missing credentials or client setup failure: not a capability
			// verdict. Log and try the next region; the gate stays a real SKIP.
			log.Printf("[INFO] CloudSSO probe: client setup failed for region %s: %T: %v", r, err, err)
			continue
		}
		client := rawClient.(*connectivity.AliyunClient)
		_, rpcErr := client.RpcPost("cloudsso", "2021-05-15", "ListDirectories", nil, map[string]interface{}{}, true)
		if rpcErr == nil {
			return string(r), true
		}
		// RPC failed in this region: service not activated or a transient probe
		// error. Continue probing the remaining supported regions. rpcErr holds
		// the real ListDirectories failure (not the client-setup err) so the
		// probe log reflects the actual RPC cause for diagnosis.
		log.Printf("[INFO] CloudSSO probe: ListDirectories failed in region %s: %T: %v", r, rpcErr, rpcErr)
	}
	return "", false
}
