package alicloud

import (
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/market"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudMarketOrder_basic(t *testing.T) {
	var v *market.DescribeOrderResponse
	resourceId := "alicloud_market_order.default"
	ra := resourceAttrInit(resourceId, marketOrderMap)
	serviceFunc := func() interface{} {
		return &MarketService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	testAccConfig := resourceTestAccConfigFunc(resourceId, "", resourceMarketOrderDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckPrePaidResources(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"product_code":    "${data.alicloud_market_products.default.ids.0}",
					"pay_type":        "PrePaid",
					"quantity":        "1",
					"duration":        "1",
					"pricing_cycle":   "Month",
					"package_version": "yuncode2713600001",
					"coupon_id":       "",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"duration", "pricing_cycle", "package_version", "coupon_id", "components"},
			},
		},
	})
}

func resourceMarketOrderDependence(name string) string {
	return `
		data "alicloud_market_products" "default" {
			name_regex = "BatchCompute"
			product_type = "MIRROR"
		}
`
}

var marketOrderMap = map[string]string{
	"pay_type":        "PrePaid",
	"quantity":        "1",
	"duration":        "1",
	"pricing_cycle":   "Month",
	"package_version": "yuncode2713600001",
	"coupon_id":       "",
}
