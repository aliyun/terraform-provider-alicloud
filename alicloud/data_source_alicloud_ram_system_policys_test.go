package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudRamSystemPolicyDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudRamSystemPolicySourceConfig(rand, map[string]string{
			"name_regex": `"^AdministratorAccess$"`,
		}),
		fakeConfig: testAccCheckAlicloudRamSystemPolicySourceConfig(rand, map[string]string{
			"name_regex": `"AdministratorAccessInvalid"`,
		}),
	}

	RamSystemPolicyCheckInfo.dataSourceTestCheck(t, rand, allConf)
}

var existRamSystemPolicyMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"policys.#":                  "1",
		"policys.0.policy_type":      "System",
		"policys.0.update_date":      CHECKSET,
		"policys.0.description":      CHECKSET,
		"policys.0.attachment_count": CHECKSET,
		"policys.0.policy_name":      CHECKSET,
		"policys.0.create_time":      CHECKSET,
	}
}

var fakeRamSystemPolicyMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"policys.#": "0",
	}
}

var RamSystemPolicyCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_ram_system_policys.default",
	existMapFunc: existRamSystemPolicyMapFunc,
	fakeMapFunc:  fakeRamSystemPolicyMapFunc,
}

func testAccCheckAlicloudRamSystemPolicySourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccRamSystemPolicy%d"
}

data "alicloud_ram_system_policys" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
