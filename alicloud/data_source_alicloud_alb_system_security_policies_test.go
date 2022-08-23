package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudALBSystemSecurityPoliciesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	const systemPolicyIds = "tls_cipher_policy_1_0"

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudAlbSystemSecurityPolicieDataSourceName(rand, map[string]string{
			"ids": fmt.Sprintf(`["%s"]`, systemPolicyIds),
		}),
		fakeConfig: testAccCheckAlicloudAlbSystemSecurityPolicieDataSourceName(rand, map[string]string{
			"ids": fmt.Sprintf(`["%s_fake"]`, systemPolicyIds),
		}),
	}

	var existDataAlicloudAlbSystemSecurityPoliciesSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"policies.#":                    "1",
			"policies.0.tls_versions.#":     "3",
			"policies.0.ciphers.#":          "19",
			"policies.0.id":                 CHECKSET,
			"policies.0.security_policy_id": CHECKSET,
		}
	}
	var fakeDataAlicloudAlbSystemSecurityPoliciesSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"policies.#": "0",
		}
	}
	var alicloudAlbSystemSecurityPolicyCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_alb_system_security_policies.default",
		existMapFunc: existDataAlicloudAlbSystemSecurityPoliciesSourceNameMapFunc,
		fakeMapFunc:  fakeDataAlicloudAlbSystemSecurityPoliciesSourceNameMapFunc,
	}
	preCheck := func() {
		testAccPreCheckWithRegions(t, true, connectivity.AlbSupportRegions)
	}
	alicloudAlbSystemSecurityPolicyCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf)
}

func testAccCheckAlicloudAlbSystemSecurityPolicieDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	return fmt.Sprintf(`

data "alicloud_alb_system_security_policies" "default" {	
	%s
}
`, strings.Join(pairs, " \n "))
}
