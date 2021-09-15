package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudSddpInstancesDataSource(t *testing.T) {
	rand := acctest.RandInt()
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudSddpInstanceDataSourceConfig(rand, map[string]string{}),
		fakeConfig:  "",
	}

	var existSddpInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"instances.#":              "1",
			"instances.0.id":           CHECKSET,
			"instances.0.instance_id":  CHECKSET,
			"instances.0.status":       CHECKSET,
			"instances.0.payment_type": "Subscription",
		}
	}

	var fakeSddpInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"instances.#": "0",
		}
	}

	var sddpInstancesRecordsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_sddp_instances.default",
		existMapFunc: existSddpInstanceMapFunc,
		fakeMapFunc:  fakeSddpInstanceMapFunc,
	}

	var preCheck = func() {
		testAccPreCheckWithRegions(t, true, connectivity.SddpSupportRegions)
	}

	sddpInstancesRecordsCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, allConf)

}

func testAccCheckAlicloudSddpInstanceDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
data "alicloud_sddp_instances" "default" {
  %s
}
`, strings.Join(pairs, "\n  "))
	return config
}
