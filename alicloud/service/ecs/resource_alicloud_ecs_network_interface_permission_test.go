package ecs_test

import (
	"github.com/aliyun/terraform-provider-alicloud/alicloud/acctest"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/service/ecs"

	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	tfacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAliCloudECSNetworkInterfacePermission_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_network_interface_permission.default"
	ra := resourceAttrInit(resourceId, AlicloudECSNetworkInterfacePermissionMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return ecs.NewService(acctest.Provider.Meta().(*connectivity.AliyunClient))
	}, "DescribeEcsNetworkInterfacePermission")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := tfacctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsnetworkinterfacepermission%d", acctest.DefaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudECSNetworkInterfacePermissionBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(t)
		},
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(),
		CheckDestroy:             rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"account_id":           "${data.alicloud_account.default.id}",
					"permission":           "InstanceAttach",
					"network_interface_id": "${alicloud_ecs_network_interface.default.id}",
					"force":                "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_id":           CHECKSET,
						"permission":           "InstanceAttach",
						"network_interface_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force"},
			},
		},
	})
}

var AlicloudECSNetworkInterfacePermissionMap0 = map[string]string{
	"status": CHECKSET,
}

func AlicloudECSNetworkInterfacePermissionBasicDependence0(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}
data "alicloud_zones" "default" {
    available_resource_creation= "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  zone_id = "${data.alicloud_zones.default.zones.0.id}"
  vswitch_name              = "${var.name}"
}

resource "alicloud_security_group" "default" {
    name = var.name
    vpc_id = alicloud_vpc.default.id
}
data "alicloud_resource_manager_resource_groups" "default"{}

resource "alicloud_ecs_network_interface" "default" {
    network_interface_name = var.name
    vswitch_id = alicloud_vswitch.default.id
    security_group_ids = [alicloud_security_group.default.id]
	description = "Basic test"
	primary_ip_address = cidrhost(alicloud_vswitch.default.cidr_block, 16)
	tags = {
		Created = "TF",
		For =    "Test",
	}
	resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
}
data "alicloud_account" "default" {}
`, name)
}
