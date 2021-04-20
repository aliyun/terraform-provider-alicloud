package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func testAccCheckDiskAttachmentDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_ecs_disk_attachment" {
			continue
		}
		// Try to find the Disk
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		ecsService := EcsService{client}
		_, err := ecsService.DescribeEcsDiskAttachment(rs.Primary.ID)

		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
	}

	return nil
}

func TestAccAlicloudEcsDiskAttachmentBasic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_disk_attachment.default"
	ra := resourceAttrInit(resourceId, AlicloudEcsDiskAttachmentMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccAlicloudEcsDiskAttachment%d", rand)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDiskAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: AlicloudEcsDiskAttachmentBasicDependence(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func TestAccAlicloudEcsDiskAttachmentMulti(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_disk_attachment.default.1"
	ra := resourceAttrInit(resourceId, AlicloudEcsDiskAttachmentMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccAlicloudEcsDiskAttachment%d", rand)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDiskAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: AlicloudEcsDiskAttachmentBasicDependenceMulti(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

var AlicloudEcsDiskAttachmentMap = map[string]string{
	"disk_id":     CHECKSET,
	"instance_id": CHECKSET,
}

func AlicloudEcsDiskAttachmentBasicDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
			default = "%s"
		}

data "alicloud_zones" default {
  available_resource_creation = "Instance"
}

data "alicloud_instance_types" "default" {
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  	cpu_core_count    = 1
	memory_size       = 2
}

data "alicloud_vpcs" "default" {
	is_default = true
}

data "alicloud_vswitches" "default" {
 vpc_id = data.alicloud_vpcs.default.ids.0
 zone_id = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  name = "tf-test"
  description = "New security group"
  vpc_id = data.alicloud_vpcs.default.ids.0
}


data "alicloud_images" "default" {
  owners = "system"
}

resource "alicloud_instance" "default" {
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  instance_name   = "${var.name}"
  host_name       = "tf-testAcc"
  image_id        = data.alicloud_images.default.images.0.id
  instance_type   = data.alicloud_instance_types.default.instance_types.0.id
  security_groups = [alicloud_security_group.default.id]
  vswitch_id      = data.alicloud_vswitches.default.ids.0
}

data "alicloud_zones" "disk" {
	available_resource_creation= "VSwitch"
}
resource "alicloud_ecs_disk" "default" {
	availability_zone = "${data.alicloud_zones.disk.zones.0.id}"
	category = "cloud_efficiency"
	delete_auto_snapshot = "true"
	description = "Test For Terraform"
	disk_name = var.name
	enable_auto_snapshot = "true"
	encrypted = "true"
	size = "500"
  	tags = {
    	Created     = "TF"
    	Environment = "Acceptance-test"
  	}
}

resource "alicloud_ecs_disk_attachment" "default" {
  disk_id = alicloud_ecs_disk.default.id
  instance_id = alicloud_instance.default.id
}
`, name)
}
func AlicloudEcsDiskAttachmentBasicDependenceMulti(name string) string {
	return fmt.Sprintf(`
variable "name" {
			default = "%s"
		}

variable "number" {
	default = "2"
}

data "alicloud_zones" default {
  available_resource_creation = "Instance"
}

data "alicloud_instance_types" "default" {
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  	cpu_core_count    = 1
	memory_size       = 2
}

data "alicloud_vpcs" "default" {
	is_default = true
}

data "alicloud_vswitches" "default" {
 vpc_id = data.alicloud_vpcs.default.ids.0
 zone_id = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  name = "tf-test"
  description = "New security group"
  vpc_id = data.alicloud_vpcs.default.ids.0
}


data "alicloud_images" "default" {
  owners = "system"
}

resource "alicloud_instance" "default" {
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  instance_name   = "${var.name}"
  host_name       = "tf-testAcc"
  image_id        = data.alicloud_images.default.images.0.id
  instance_type   = data.alicloud_instance_types.default.instance_types.0.id
  security_groups = [alicloud_security_group.default.id]
  vswitch_id      = data.alicloud_vswitches.default.ids.0
}

data "alicloud_zones" "disk" {
	available_resource_creation= "VSwitch"
}
resource "alicloud_ecs_disk" "default" {
	count = "${var.number}"
	availability_zone = "${data.alicloud_zones.disk.zones.0.id}"
	category = "cloud_efficiency"
	delete_auto_snapshot = "true"
	description = "Test For Terraform"
	disk_name = var.name
	enable_auto_snapshot = "true"
	encrypted = "true"
	size = "500"
  	tags = {
    	Created     = "TF"
    	Environment = "Acceptance-test"
  	}
}

resource "alicloud_ecs_disk_attachment" "default" {
  count = "${var.number}"
  disk_id = "${element(alicloud_ecs_disk.default.*.id, count.index)}"
  instance_id = alicloud_instance.default.id
}
`, name)
}
