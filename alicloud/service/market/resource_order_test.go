package market_test

import (
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/market"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	marketsvc "github.com/aliyun/terraform-provider-alicloud/alicloud/service/market"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// lintignore: AT001
func TestAccAliCloudMarketOrder_basic(t *testing.T) {
	var v *market.DescribeOrderResponse
	resourceId := "alicloud_market_order.default"
	ra := acctest.ResourceAttrInit(resourceId, marketOrderMap)
	serviceFunc := func() interface{} {
		return marketsvc.NewMarketService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}
	rc := acctest.ResourceCheckInit(resourceId, &v, serviceFunc)
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, "", resourceMarketOrderDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
			acctest.PreCheckPrePaidResources(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
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
					"components":      map[string]interface{}{"package_version": "yuncode2713600001"},
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
