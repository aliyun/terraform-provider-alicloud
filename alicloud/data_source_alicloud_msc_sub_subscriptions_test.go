package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudMscSubSubscriptionsDataSource(t *testing.T) {
	checkoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	rand := acctest.RandInt()
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudMscSubSubscriptionDataSourceConfig(rand, map[string]string{}),
		fakeConfig:  "",
	}

	var existMscSubSubscriptionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"subscriptions.#":           CHECKSET,
			"subscriptions.0.id":        CHECKSET,
			"subscriptions.0.item_name": CHECKSET,
		}
	}

	var fakeMscSubSubscriptionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"subscriptions.#": "0",
		}
	}

	var sddpInstancesRecordsCheckInfo = dataSourceAttr{
		resourceId:   "data.alicloud_msc_sub_subscriptions.default",
		existMapFunc: existMscSubSubscriptionsMapFunc,
		fakeMapFunc:  fakeMscSubSubscriptionsMapFunc,
	}

	sddpInstancesRecordsCheckInfo.dataSourceTestCheck(t, rand, allConf)

}

func testAccCheckAlicloudMscSubSubscriptionDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
data "alicloud_msc_sub_subscriptions" "default" {
  %s
}
`, strings.Join(pairs, "\n  "))
	return config
}
