package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudConfigConfigurationRecordersDataSource(t *testing.T) {
	rand := acctest.RandIntRange(1000000, 9999999)
	resourceId := "data.alicloud_config_configuration_recorders.example"

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudConfigConfigurationRecordersSourceConfig(rand, map[string]string{}),
	}

	var existConfigConfigurationRecordersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"recorders.#":                            "1",
			"recorders.0.id":                         CHECKSET,
			"recorders.0.account_id":                 CHECKSET,
			"recorders.0.organization_enable_status": CHECKSET,
			"recorders.0.organization_master_id":     CHECKSET,
			"recorders.0.resource_types.#":           "2",
			"recorders.0.status":                     "REGISTERED",
		}
	}

	var fakeConfigConfigurationRecordersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"recorders.#": "0",
		}
	}

	var configConfigurationRecordersCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existConfigConfigurationRecordersMapFunc,
		fakeMapFunc:  fakeConfigConfigurationRecordersMapFunc,
	}

	configConfigurationRecordersCheckInfo.dataSourceTestCheck(t, rand, allConf)

}

func testAccCheckAlicloudConfigConfigurationRecordersSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccConfigConfigurationRecorders%d"
}

resource "alicloud_config_configuration_recorder" "example" {
 resource_types = ["ACS::ECS::Disk","ACS::ECS::Instance"]
}

data "alicloud_config_configuration_recorders" "example"{
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
