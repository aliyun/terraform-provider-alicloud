package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Case Dataworks网络资源管理-TF验收_北京 8963
func TestAccAliCloudDataWorksNetwork_basic8963_modify(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_data_works_network.default"
	ra := resourceAttrInit(resourceId, AlicloudDataWorksNetworkMap8963)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &DataWorksServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeDataWorksNetwork")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf_testacc_dwpt%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudDataWorksNetworkBasicDependence8963)
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
					"vpc_id":               "${data.alicloud_vpcs.default2.ids.0}",
					"vswitch_id":           "${data.alicloud_vswitches.default2.ids.0}",
					"dw_resource_group_id": "${alicloud_data_works_dw_resource_group.defaultVJvKvl.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpc_id":               CHECKSET,
						"vswitch_id":           CHECKSET,
						"dw_resource_group_id": CHECKSET,
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

var AlicloudDataWorksNetworkMap8963 = map[string]string{
	"status":      CHECKSET,
	"create_time": CHECKSET,
}

func AlicloudDataWorksNetworkBasicDependence8963(name string) string {
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

`, name)
}
