package alicloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudDiskAttachment(t *testing.T) {
	var i ecs.Instance
	var v ecs.Disk

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_disk_attachment.disk-att",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDiskAttachmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDiskAttachmentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(
						"alicloud_instance.instance", &i),
					testAccCheckDiskExists(
						"alicloud_disk.disk", &v),
					testAccCheckDiskAttachmentExists(
						"alicloud_disk_attachment.disk-att", &i, &v),
				),
			},
		},
	})

}

func TestAccAlicloudDiskMultiAttachment(t *testing.T) {
	var i ecs.Instance
	var v ecs.Disk

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_disk_attachment.disks-attach.0",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDiskAttachmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccMultiDiskAttachmentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(
						"alicloud_instance.instance", &i),
					testAccCheckDiskExists(
						"alicloud_disk.disks.0", &v),
					testAccCheckDiskAttachmentExists(
						"alicloud_disk_attachment.disks-attach.0", &i, &v),
				),
			},
		},
	})
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_disk_attachment.disks-attach.1",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckDiskAttachmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccMultiDiskAttachmentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckInstanceExists(
						"alicloud_instance.instance", &i),
					testAccCheckDiskExists(
						"alicloud_disk.disks.1", &v),
					testAccCheckDiskAttachmentExists(
						"alicloud_disk_attachment.disks-attach.1", &i, &v),
				),
			},
		},
	})

}

func testAccCheckDiskAttachmentExists(n string, instance *ecs.Instance, disk *ecs.Disk) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Disk ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)

		return resource.Retry(3*time.Minute, func() *resource.RetryError {
			d, err := client.DescribeDiskById(instance.InstanceId, rs.Primary.Attributes["disk_id"])
			if err != nil {
				return resource.NonRetryableError(err)
			}
			if d.Status != string(DiskInUse) {
				return resource.RetryableError(fmt.Errorf("Disk is in attaching - trying again while it attaches"))
			}

			*disk = d
			return nil
		})
	}
}

func testAccCheckDiskAttachmentDestroy(s *terraform.State) error {

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_disk_attachment" {
			continue
		}
		// Try to find the Disk
		client := testAccProvider.Meta().(*AliyunClient)

		disk, err := client.DescribeDiskById("", rs.Primary.ID)

		if err != nil {
			if NotFoundError(err) {
				return nil
			}
			return fmt.Errorf("Describing disk %s got an error.", rs.Primary.ID)
		}
		if disk.Status != string(Available) {
			return fmt.Errorf("Error ECS Disk Attachment still exist")
		}
	}

	return nil
}

const testAccDiskAttachmentConfig = `
resource "alicloud_disk" "disk" {
  availability_zone = "cn-beijing-a"
  size = "50"

  tags {
    Name = "TerraformTest-disk"
  }
}

resource "alicloud_instance" "instance" {
  image_id = "ubuntu_140405_64_40G_cloudinit_20161115.vhd"
  instance_type = "ecs.n4.small"
  availability_zone = "cn-beijing-a"
  security_groups = ["${alicloud_security_group.group.id}"]
  instance_name = "hello"
  internet_charge_type = "PayByBandwidth"

  tags {
    Name = "TerraformTest-instance"
  }
}

resource "alicloud_disk_attachment" "disk-att" {
  disk_id = "${alicloud_disk.disk.id}"
  instance_id = "${alicloud_instance.instance.id}"
}

resource "alicloud_security_group" "group" {
  name = "terraform-test-group"
  description = "New security group"
}
`
const testAccMultiDiskAttachmentConfig = `

variable "count" {
  default = "2"
}

resource "alicloud_disk" "disks" {
  count = "${var.count}"
  availability_zone = "cn-beijing-a"
  size = "50"

  tags {
    Name = "TerraformTest-disk-${count.index}"
  }
}

resource "alicloud_instance" "instance" {
  image_id = "ubuntu_140405_64_40G_cloudinit_20161115.vhd"
  instance_type = "ecs.n4.small"
  availability_zone = "cn-beijing-a"
  security_groups = ["${alicloud_security_group.group.id}"]
  instance_name = "hello"
  internet_charge_type = "PayByBandwidth"

  tags {
    Name = "TerraformTest-instance"
  }
}

resource "alicloud_disk_attachment" "disks-attach" {
  count = "${var.count}"
  disk_id     = "${element(alicloud_disk.disks.*.id, count.index)}"
  instance_id = "${alicloud_instance.instance.id}"
}

resource "alicloud_security_group" "group" {
  name = "terraform-test-group"
  description = "New security group"
}
`
