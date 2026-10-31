package ram_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudRamSystemPolicyDataSource(t *testing.T) {
	acctest.PreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := tfacctest.RandIntRange(1000000, 9999999)

	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudRamSystemPolicySourceConfig(rand, map[string]string{
			"name_regex": `"^AdministratorAccess$"`,
		}),
		FakeConfig: testAccCheckAlicloudRamSystemPolicySourceConfig(rand, map[string]string{
			"name_regex": `"AdministratorAccessInvalid"`,
		}),
	}

	RamSystemPolicyCheckInfo.DataSourceTestCheck(t, rand, allConf)
}

var existRamSystemPolicyMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"policys.#":                  "1",
		"policys.0.policy_type":      "System",
		"policys.0.update_date":      acctest.CHECKSET,
		"policys.0.description":      acctest.CHECKSET,
		"policys.0.attachment_count": acctest.CHECKSET,
		"policys.0.policy_name":      acctest.CHECKSET,
		"policys.0.create_time":      acctest.CHECKSET,
	}
}

var fakeRamSystemPolicyMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"policys.#": "0",
	}
}

var RamSystemPolicyCheckInfo = acctest.DataSourceAttr{
	ResourceId:   "data.alicloud_ram_system_policys.default",
	ExistMapFunc: existRamSystemPolicyMapFunc,
	FakeMapFunc:  fakeRamSystemPolicyMapFunc,
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
