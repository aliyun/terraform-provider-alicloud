package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Case DataWorks Route resource management - TF acceptance test
func TestAccAliCloudDataWorksRoute_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_data_works_route.default"
	ra := resourceAttrInit(resourceId, AlicloudDataWorksRouteMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DataWorksServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDataWorksRoute")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc_dwrt%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDataWorksRouteBasicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-shenzhen"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"network_id":       "${alicloud_data_works_network.default.id}",
					"destination_cidr": "192.168.10.0/24",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination_cidr":     "192.168.10.0/24",
						"network_id":           CHECKSET,
						"route_id":             CHECKSET,
						"create_time":          CHECKSET,
						"dw_resource_group_id": CHECKSET,
						"region_id":            CHECKSET,
						"resource_id":          CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"network_id":       "${alicloud_data_works_network.default.id}",
					"destination_cidr": "192.168.20.0/24",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"destination_cidr": "192.168.20.0/24",
						"network_id":       CHECKSET,
						"route_id":         CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

var AlicloudDataWorksRouteMap = map[string]string{
	"route_id":             CHECKSET,
	"create_time":          CHECKSET,
	"dw_resource_group_id": CHECKSET,
	"region_id":            CHECKSET,
	"resource_id":          CHECKSET,
}

func AlicloudDataWorksRouteBasicDependence(name string) string {
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

data "alicloud_vpcs" "default2" {
	name_regex = "^default-NODELETING-2$"
}

data "alicloud_vswitches" "default2" {
	vpc_id  = data.alicloud_vpcs.default2.ids.0
}


resource "alicloud_data_works_dw_resource_group" "defaultVJvKvl" {
  payment_duration_unit = "Month"
  payment_type          = "PostPaid"
  specification         = "500"
  default_vswitch_id    = data.alicloud_vswitches.default.ids.0
  remark                = "OpenAPI测试用资源组"
  resource_group_name   = var.name
  default_vpc_id        = data.alicloud_vpcs.default.ids.0
  auto_renew = false
}

resource "alicloud_data_works_network" "default" {
  vpc_id               = data.alicloud_vpcs.default2.ids.0
  vswitch_id           = data.alicloud_vswitches.default2.ids.0
  dw_resource_group_id = alicloud_data_works_dw_resource_group.defaultVJvKvl.id
}

`, name)
}
