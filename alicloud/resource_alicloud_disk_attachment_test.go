package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudDiskAttachment(t *testing.T) {
	var i ecs.Instance
	var v ecs.Disk
	var attachment ecs.Disk
	serverFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	diskRc := resourceCheckInit("alicloud_disk.default", &v, serverFunc)

	instanceRc := resourceCheckInit("alicloud_instance.default", &i, serverFunc)

	attachmentRc := resourceCheckInit("alicloud_disk_attachment.default", &attachment, serverFunc)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_disk_attachment.default",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDiskAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDiskAttachmentConfig(EcsInstanceCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					diskRc.checkResourceExists(),
					instanceRc.checkResourceExists(),
					attachmentRc.checkResourceExists(),
				),
			},
		},
	})

}

func TestAccAlicloudDiskMultiAttachment(t *testing.T) {
	var i ecs.Instance
	var v ecs.Disk
	var attachment ecs.Disk
	serverFunc := func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	diskRc := resourceCheckInit("alicloud_disk.default.1", &v, serverFunc)

	instanceRc := resourceCheckInit("alicloud_instance.default", &i, serverFunc)

	attachmentRc := resourceCheckInit("alicloud_disk_attachment.default.1", &attachment, serverFunc)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_disk_attachment.default.1",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDiskAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMultiDiskAttachmentConfig(EcsInstanceCommonTestCase),
				Check: resource.ComposeTestCheckFunc(
					diskRc.checkResourceExists(),
					instanceRc.checkResourceExists(),
					attachmentRc.checkResourceExists(),
				),
			},
		},
	})

}

func testAccCheckDiskAttachmentDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_disk_attachment" {
			continue
		}
		// Try to find the Disk
		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		ecsService := EcsService{client}
		_, err := ecsService.DescribeDiskAttachment(rs.Primary.ID)

		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}
	}

	return nil
}

func testAccDiskAttachmentConfig(common string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEcsDiskAttachmentConfig"
	}

	resource "alicloud_disk" "default" {
	  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	  size = "50"
	  name = "${var.name}"

	  tags {
	    Name = "TerraformTest-disk"
	  }
	}

	resource "alicloud_instance" "default" {
		image_id = "${data.alicloud_images.default.images.0.id}"
		availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		system_disk_category = "cloud_ssd"
		system_disk_size = 40
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		security_groups = ["${alicloud_security_group.default.id}"]
		instance_name = "${var.name}"
		vswitch_id = "${alicloud_vswitch.default.id}"
	}

	resource "alicloud_disk_attachment" "default" {
	  disk_id = "${alicloud_disk.default.id}"
	  instance_id = "${alicloud_instance.default.id}"
	}
	`, common)
}
func testAccMultiDiskAttachmentConfig(common string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEcsDiskAttachmentConfig"
	}

	variable "count" {
		default = "2"
	}

	resource "alicloud_disk" "default" {
		name = "${var.name}-${count.index}"
		count = "${var.count}"
		availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		size = "50"

		tags {
			Name = "TerraformTest-disk-${count.index}"
		}
	}

	resource "alicloud_instance" "default" {
		image_id = "${data.alicloud_images.default.images.0.id}"
		availability_zone = "${data.alicloud_zones.default.zones.0.id}"
		system_disk_category = "cloud_ssd"
		system_disk_size = 40
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		security_groups = ["${alicloud_security_group.default.id}"]
		instance_name = "${var.name}"
		vswitch_id = "${alicloud_vswitch.default.id}"
	}

	resource "alicloud_disk_attachment" "default" {
		count = "${var.count}"
		disk_id     = "${element(alicloud_disk.default.*.id, count.index)}"
		instance_id = "${alicloud_instance.default.id}"
	}
	`, common)
}
