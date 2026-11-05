package alb_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAlicloudALBSystemSecurityPoliciesDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	const systemPolicyIds = "tls_cipher_policy_1_0"

	idsConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudAlbSystemSecurityPolicieDataSourceName(rand, map[string]string{
			"ids": fmt.Sprintf(`["%s"]`, systemPolicyIds),
		}),
		FakeConfig: testAccCheckAlicloudAlbSystemSecurityPolicieDataSourceName(rand, map[string]string{
			"ids": fmt.Sprintf(`["%s_fake"]`, systemPolicyIds),
		}),
	}

	var existDataAlicloudAlbSystemSecurityPoliciesSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"policies.#":                    "1",
			"policies.0.tls_versions.#":     "3",
			"policies.0.ciphers.#":          "19",
			"policies.0.id":                 acctest.CHECKSET,
			"policies.0.security_policy_id": acctest.CHECKSET,
		}
	}
	var fakeDataAlicloudAlbSystemSecurityPoliciesSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"policies.#": "0",
		}
	}
	var alicloudAlbSystemSecurityPolicyCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_alb_system_security_policies.default",
		ExistMapFunc: existDataAlicloudAlbSystemSecurityPoliciesSourceNameMapFunc,
		FakeMapFunc:  fakeDataAlicloudAlbSystemSecurityPoliciesSourceNameMapFunc,
	}
	preCheck := func() {
		acctest.PreCheckWithRegions(t, true, connectivity.AlbSupportRegions)
	}
	alicloudAlbSystemSecurityPolicyCheckInfo.DataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf)
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
