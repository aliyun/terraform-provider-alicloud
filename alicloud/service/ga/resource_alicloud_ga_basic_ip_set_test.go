package ga_test

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ga"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudGaBasicIpSet_basic0(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_basic_ip_set.default"
	ra := acctest.ResourceAttrInit(resourceId, resourceAliCloudGaBasicIpSetMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ga.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeGaBasicIpSet")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tf-testacc%s-name%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudGaBasicIpSetBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"accelerator_id":       "${alicloud_ga_basic_accelerator.default.id}",
					"accelerate_region_id": acctest.DefaultRegionToTest,
					"isp_type":             "BGP",
					"bandwidth":            "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_id":       acctest.CHECKSET,
						"accelerate_region_id": acctest.CHECKSET,
						"isp_type":             "BGP",
						"bandwidth":            "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth": "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "5",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var resourceAliCloudGaBasicIpSetMap = map[string]string{
	"isp_type":  acctest.CHECKSET,
	"bandwidth": acctest.CHECKSET,
	"status":    acctest.CHECKSET,
}

func AliCloudGaBasicIpSetBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	resource "alicloud_ga_basic_accelerator" "default" {
  		duration               = 1
  		pricing_cycle          = "Month"
  		basic_accelerator_name = var.name
  		description            = var.name
  		bandwidth_billing_type = "CDT"
  		auto_pay               = true
  		auto_use_coupon        = "true"
  		auto_renew             = false
  		auto_renew_duration    = 1
	}
`, name)
}
