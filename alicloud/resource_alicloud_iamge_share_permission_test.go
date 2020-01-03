package alicloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAliCloudImageSharePermission(t *testing.T) {
	var v *ecs.DescribeImageSharePermissionResponse
	resourceId := "alicloud_image_share_permission.default"
	ra := resourceAttrInit(resourceId, testAccImageSharePermissionCheckMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeImageShareByImageId")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAccEcsImageShareConfigBasic%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceImageSharePermissionConfigDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithMultipleAccount(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"image_id":   "${alicloud_image.default.id}",
					"account_id": os.Getenv("ALICLOUD_ACCESS_KEY_2"),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"account_id": CHECKSET,
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

var testAccImageSharePermissionCheckMap = map[string]string{
	"image_id": CHECKSET,
}

func resourceImageSharePermissionConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
data "alicloud_instance_types" "default" {
 	cpu_core_count    = 1
	memory_size       = 2
}
data "alicloud_images" "default" {
  name_regex  = "^ubuntu*"
  owners      = "system"
}
resource "alicloud_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"
  name              = "${var.name}"
}
resource "alicloud_security_group" "default" {
  name   = "${var.name}"
  vpc_id = "${alicloud_vpc.default.id}"
}
resource "alicloud_instance" "default" {
  image_id = "${data.alicloud_images.default.ids[0]}"
  instance_type = "${data.alicloud_instance_types.default.ids[0]}"
  security_groups = "${[alicloud_security_group.default.id]}"
  vswitch_id = "${alicloud_vswitch.default.id}"
  instance_name = "${var.name}"
}
resource "alicloud_image" "default" {
  instance_id = "${alicloud_instance.default.id}"
  name        = "${var.name}"
}
`, name)
}
