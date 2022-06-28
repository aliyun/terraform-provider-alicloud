package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudECSImageExport(t *testing.T) {
	var v ecs.Image
	resourceId := "alicloud_image_export.default"
	ra := resourceAttrInit(resourceId, testAccExportImageCheckMap)
	serviceFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}

	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, serviceFunc, "DescribeImageById")
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandIntRange(1000, 9999)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testaccecsimageexportconfigbasic%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceImageExportBasicConfigDependence)
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
					"image_id":   "${alicloud_image.default.id}",
					"oss_bucket": "${alicloud_oss_bucket.default.bucket}",
					"oss_prefix": "ecsExport",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"oss_prefix": "ecsExport",
					}),
				),
			},
		},
	})
}

var testAccExportImageCheckMap = map[string]string{
	"image_id":   CHECKSET,
	"oss_bucket": CHECKSET,
}

func resourceImageExportBasicConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

data "alicloud_instance_types" "default" {
	cpu_core_count    = 1
	memory_size       = 2
}
data "alicloud_images" "default" {
 name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
 owners      = "system"
}

data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
	vpc_id = data.alicloud_vpcs.default.ids.0
	zone_id      = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
}
resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  vswitch_name      = var.name
}

locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}
resource "alicloud_security_group" "default" {
 name   = "${var.name}"
 vpc_id = data.alicloud_vpcs.default.ids.0
}
resource "alicloud_instance" "default" {
 image_id = "${data.alicloud_images.default.ids[0]}"
 instance_type = "${data.alicloud_instance_types.default.ids[0]}"
 security_groups = "${[alicloud_security_group.default.id]}"
 vswitch_id = local.vswitch_id
 instance_name = "${var.name}"
}
resource "alicloud_image" "default" {
 instance_id = "${alicloud_instance.default.id}"
 image_name        = "${var.name}"
}
resource "alicloud_oss_bucket" "default" {
  bucket = "${var.name}"
}
`, name)
}
