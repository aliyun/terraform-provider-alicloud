package cloudfirewall_test

import (
	"fmt"
	"strings"
	"testing"

	acctest "github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudCloudFirewallInstancesDataSource(t *testing.T) {
	rand := tfacctest.RandInt()
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudCloudFirewallInstanceDataSourceConfig(rand, map[string]string{}),
		FakeConfig:  "",
	}

	var existCloudFirewallInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"instances.#":              "1",
			"instances.0.id":           acctest.CHECKSET,
			"instances.0.instance_id":  acctest.CHECKSET,
			"instances.0.status":       acctest.CHECKSET,
			"instances.0.payment_type": "Subscription",
		}
	}

	var fakeCloudFirewallInstancesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"instances.#": "0",
		}
	}

	var cloudFirewallInstancesRecordsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_cloud_firewall_instances.default",
		ExistMapFunc: existCloudFirewallInstancesMapFunc,
		FakeMapFunc:  fakeCloudFirewallInstancesMapFunc,
	}

	cloudFirewallInstancesRecordsCheckInfo.DataSourceTestCheck(t, rand, allConf)

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
