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

func TestAccAliCloudGaCustomRoutingEndpointGroup_basic0(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_custom_routing_endpoint_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudGaCustomRoutingEndpointGroupMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ga.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeGaCustomRoutingEndpointGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sGaCustomRoutingEndpointGroup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudGaCustomRoutingEndpointGroupBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"accelerator_id":        "${alicloud_ga_listener.default.accelerator_id}",
					"listener_id":           "${alicloud_ga_listener.default.id}",
					"endpoint_group_region": acctest.DefaultRegionToTest,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_id":        acctest.CHECKSET,
						"listener_id":           acctest.CHECKSET,
						"endpoint_group_region": acctest.CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"custom_routing_endpoint_group_name": name + "-update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"custom_routing_endpoint_group_name": name + "-update",
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

func TestAccAliCloudGaCustomRoutingEndpointGroup_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	acctest.CheckoutSupportedRegions(t, true, connectivity.GaSupportRegions)
	resourceId := "alicloud_ga_custom_routing_endpoint_group.default"
	ra := acctest.ResourceAttrInit(resourceId, AlicloudGaCustomRoutingEndpointGroupMap)
	rc := acctest.ResourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ga.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeGaCustomRoutingEndpointGroup")
	rac := acctest.ResourceAttrCheckInit(rc, ra)
	testAccCheck := rac.ResourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sGaCustomRoutingEndpointGroup%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := acctest.ResourceTestAccConfig(resourceId, name, AlicloudGaCustomRoutingEndpointGroupBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProviderFactories: acctest.SDKProviderFactories,
		CheckDestroy:      rac.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"accelerator_id":                     "${alicloud_ga_listener.default.accelerator_id}",
					"listener_id":                        "${alicloud_ga_listener.default.id}",
					"endpoint_group_region":              acctest.DefaultRegionToTest,
					"custom_routing_endpoint_group_name": name,
					"description":                        name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"accelerator_id":                     acctest.CHECKSET,
						"listener_id":                        acctest.CHECKSET,
						"endpoint_group_region":              acctest.CHECKSET,
						"custom_routing_endpoint_group_name": name,
						"description":                        name,
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

var AlicloudGaCustomRoutingEndpointGroupMap = map[string]string{
	"status": acctest.CHECKSET,
}

func AlicloudGaCustomRoutingEndpointGroupBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_ga_accelerators" "default" {
  		status = "active"
		bandwidth_billing_type = "BandwidthPackage"
	}
	
	resource "alicloud_ga_bandwidth_package" "default" {
		bandwidth      = 100
  		type           = "Basic"
  		bandwidth_type = "Basic"
  		payment_type   = "PayAsYouGo"
  		billing_type   = "PayBy95"
  		ratio          = 30
	}
	
	resource "alicloud_ga_bandwidth_package_attachment" "default" {
  		accelerator_id       = data.alicloud_ga_accelerators.default.accelerators.0.id
  		bandwidth_package_id = alicloud_ga_bandwidth_package.default.id
	}
	
	resource "alicloud_ga_listener" "default" {
  		accelerator_id = alicloud_ga_bandwidth_package_attachment.default.accelerator_id
		listener_type  = "CustomRouting"
  		port_ranges {
    		from_port = 10000
    		to_port   = 16000
  		}
	}
`, name)
}
