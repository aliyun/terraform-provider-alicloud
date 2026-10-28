package mscopensubscription_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

func TestAccAliCloudMscSubSubscriptionsDataSource(t *testing.T) {
	acctest.CheckoutSupportedRegions(t, true, connectivity.TestSalveRegions)
	rand := tfacctest.RandInt()
	allConf := acctest.DataSourceTestAccConfig{
		ExistConfig: testAccCheckAlicloudMscSubSubscriptionDataSourceConfig(rand, map[string]string{}),
		FakeConfig:  "",
	}

	var existMscSubSubscriptionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"subscriptions.#":           acctest.CHECKSET,
			"subscriptions.0.id":        acctest.CHECKSET,
			"subscriptions.0.item_name": acctest.CHECKSET,
		}
	}

	var fakeMscSubSubscriptionsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"subscriptions.#": "0",
		}
	}

	var sddpInstancesRecordsCheckInfo = acctest.DataSourceAttr{
		ResourceId:   "data.alicloud_msc_sub_subscriptions.default",
		ExistMapFunc: existMscSubSubscriptionsMapFunc,
		FakeMapFunc:  fakeMscSubSubscriptionsMapFunc,
	}

	sddpInstancesRecordsCheckInfo.DataSourceTestCheck(t, rand, allConf)

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
