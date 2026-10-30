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

func TestAccAliCloudGaBasicEndpointGroup_basic0(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_basic_endpoint_group.default"
	ra := acctest.ResourceAttrInit(resourceId, resourceAliCloudGaBasicEndpointGroupMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ga.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeGaBasicEndpointGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(100, 999)
	name := fmt.Sprintf("tf-testacc%s-name%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AliCloudGaBasicEndpointGroupBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"accelerator_id":            "${alicloud_ga_basic_accelerator.default.id}",
					"endpoint_group_region":     "cn-beijing",
					"endpoint_type":             "SLB",
					"endpoint_address":          "${alicloud_slb_load_balancer.default.id}",
					"endpoint_sub_address":      "192.168.0.1",
					"basic_endpoint_group_name": name,
					"description":               name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_id":            acctest.CHECKSET,
						"endpoint_group_region":     "cn-beijing",
						"endpoint_type":             "SLB",
						"endpoint_address":          acctest.CHECKSET,
						"endpoint_sub_address":      "192.168.0.1",
						"basic_endpoint_group_name": name,
						"description":               name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"basic_endpoint_group_name": name + "-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"basic_endpoint_group_name": name + "-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "-update",
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

var resourceAliCloudGaBasicEndpointGroupMap = map[string]string{
	"endpoint_type":        acctest.CHECKSET,
	"endpoint_address":     acctest.CHECKSET,
	"endpoint_sub_address": acctest.CHECKSET,
	"status":               acctest.CHECKSET,
}

func AliCloudGaBasicEndpointGroupBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_vpcs" "default" {
		name_regex = "^default-NODELETING$"
	}
	
	data "alicloud_vswitches" "default" {
  		vpc_id  = data.alicloud_vpcs.default.ids.0
	}

	resource "alicloud_slb_load_balancer" "default" {
  		load_balancer_spec = "slb.s2.small"
  		vswitch_id = data.alicloud_vswitches.default.ids.0
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
