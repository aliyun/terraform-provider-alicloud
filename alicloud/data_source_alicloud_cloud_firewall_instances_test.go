package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCloudFirewallInstancesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudFirewallInstanceDataSourceConfig(rand, map[string]string{}),
		fakeConfig:  "",
	}

	var existCloudFirewallInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"instances.#":              "1",
			"instances.0.id":           CHECKSET,
			"instances.0.instance_id":  CHECKSET,
			"instances.0.status":       CHECKSET,
			"instances.0.payment_type": "Subscription",
		}
	}

	var fakeCloudFirewallInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"instances.#": "0",
		}
	}

	var cloudFirewallInstancesRecordsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_cloud_firewall_instances.default",
		existMapFunc: existCloudFirewallInstancesMapFunc,
		fakeMapFunc:  fakeCloudFirewallInstancesMapFunc,
	}

	cloudFirewallInstancesRecordsCheckInfo.dataSourceTestCheck(t, rand, allConf)

}

func testAccCheckAlicloudCloudFirewallInstanceDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
data "alicloud_cloud_firewall_instances" "default" {
  %s
}
`, strings.Join(pairs, "\n  "))
	return config
}
