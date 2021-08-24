package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudNasMountTarget_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_nas_mount_target.default"
	ra := resourceAttrInit(resourceId, AlicloudNasMountTarget0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NasService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNasMountTarget")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sAlicloudNasMountTarget%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudNasMountTargetBasicDependence0)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"access_group_name": "${alicloud_nas_access_group.example.access_group_name}",
					"file_system_id":    "${alicloud_nas_file_system.example.id}",
					"vswitch_id":        "${data.alicloud_vpcs.example.vpcs.0.vswitch_ids.0}",
					"security_group_id": "${alicloud_security_group.example.id}",
					"status":            "Active",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_group_name": name,
						"file_system_id":    CHECKSET,
						"vswitch_id":        CHECKSET,
						"security_group_id": CHECKSET,
						"status":            "Active",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"security_group_id"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"access_group_name": "${alicloud_nas_access_group.example1.access_group_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_group_name": name + "change",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Inactive",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Inactive",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"access_group_name": "${alicloud_nas_access_group.example.access_group_name}",
					"file_system_id":    "${alicloud_nas_file_system.example.id}",
					"vswitch_id":        "${data.alicloud_vpcs.example.vpcs.0.vswitch_ids.0}",
					"security_group_id": "${alicloud_security_group.example.id}",
					"status":            "Active",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_group_name": name,
						"file_system_id":    CHECKSET,
						"vswitch_id":        CHECKSET,
						"security_group_id": CHECKSET,
						"status":            "Active",
					}),
				),
			},
		},
	})
}

var AlicloudNasMountTarget0 = map[string]string{}

func AlicloudNasMountTargetBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

variable "name1" {
	default = "%schange"
}

data "alicloud_nas_protocols" "example" {
	type = "Performance"
}

data "alicloud_vpcs" "example" {
	name_regex = "default-NODELETING"
}

resource "alicloud_security_group" "example" {
	name = var.name
	vpc_id = "${data.alicloud_vpcs.example.vpcs.0.id}"
}

resource "alicloud_nas_file_system" "example" {
	protocol_type = "${data.alicloud_nas_protocols.example.protocols.0}"
	storage_type = "Performance"
}

resource "alicloud_nas_access_group" "example" {
	access_group_name = "${var.name}"
	access_group_type = "Vpc"
}

resource "alicloud_nas_access_group" "example1" {
	access_group_name = "${var.name1}"
	access_group_type = "Vpc"
}

resource "alicloud_nas_mount_target" "example" {
	file_system_id = "${alicloud_nas_file_system.example.id}"
	access_group_name = "${alicloud_nas_access_group.example.access_group_name}"
	vswitch_id = "${data.alicloud_vpcs.example.vpcs.0.vswitch_ids.0}"
	security_group_id = "${alicloud_security_group.example.id}"
}
`, name, name)
}
