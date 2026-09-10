package alicloud

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func TestAccAlicloudCloudControlPriceDataSource(t *testing.T) {
	testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
	rand := acctest.RandIntRange(1000000, 9999999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudCloudControlPriceSourceConfig(rand, map[string]string{
			"product":       `"SLB"`,
			"resource_code": `"LoadBalancer"`,
		}),
	}

	CloudControlPriceCheckInfo.dataSourceTestCheck(t, rand, idsConf)
}

var existCloudControlPriceMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"prices.#":                     "1",
		"prices.0.original_price":      CHECKSET,
		"prices.0.discount_price":      CHECKSET,
		"prices.0.currency":            CHECKSET,
		"prices.0.module_details.#":    CHECKSET,
		"prices.0.promotion_details.#": CHECKSET,
		"prices.0.trade_price":         CHECKSET,
	}
}

var fakeCloudControlPriceMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"prices.#": "0",
	}
}

var CloudControlPriceCheckInfo = dataSourceAttr{
	resourceId:   "data.alicloud_cloud_control_prices.default",
	existMapFunc: existCloudControlPriceMapFunc,
	fakeMapFunc:  fakeCloudControlPriceMapFunc,
}

func testAccCheckAlicloudCloudControlPriceSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccCloudControlPrice%d"
}


data "alicloud_cloud_control_prices" "default" {
    desire_attributes = {
      AddressType = "internet"
      PaymentType = "PayAsYouGo"
    }
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
