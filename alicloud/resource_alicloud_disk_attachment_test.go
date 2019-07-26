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
				Config: testAccDiskAttachmentConfig(),
				Check: resource.ComposeTestCheckFunc(
					diskRc.checkResourceExists(),
					instanceRc.checkResourceExists(),
					attachmentRc.checkResourceExists(),
					resource.TestCheckResourceAttrSet(
						"alicloud_disk_attachment.default", "device_name"),
				),
			},
			{
				Config: testAccDiskAttachmentConfigResize(),
				Check: resource.ComposeTestCheckFunc(
					diskRc.checkResourceExists(),
					instanceRc.checkResourceExists(),
					attachmentRc.checkResourceExists(),
					resource.TestCheckResourceAttrSet(
						"alicloud_disk_attachment.default", "device_name"),
					resource.TestCheckResourceAttr(
						"alicloud_disk.default", "size", "70"),
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
				Config: testAccMultiDiskAttachmentConfig(EcsInstanceCommonNoZonesTestCase),
				Check: resource.ComposeTestCheckFunc(
					diskRc.checkResourceExists(),
					instanceRc.checkResourceExists(),
					attachmentRc.checkResourceExists(),
					resource.TestCheckResourceAttrSet(
						"alicloud_disk_attachment.default.1", "device_name"),
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

func testAccDiskAttachmentConfig() string {
	return fmt.Sprintf(`
    data "alicloud_instance_types" "default" {
      cpu_core_count    = 2
      memory_size       = 4
    }
    data "alicloud_images" "default" {
	  # test for windows service
      name_regex  = "^win*"

      most_recent = true
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
    resource "alicloud_security_group_rule" "default" {
      	type = "ingress"
      	ip_protocol = "tcp"
      	nic_type = "intranet"
      	policy = "accept"
      	port_range = "22/22"
      	priority = 1
      	security_group_id = "${alicloud_security_group.default.id}"
      	cidr_ip = "172.16.0.0/24"
    }	

	variable "name" {
		default = "tf-testAccEcsDiskAttachmentConfig"
	}

	resource "alicloud_disk" "default" {
	  availability_zone = "${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"
	  size = "50"
	  name = "${var.name}"
	  category = "cloud_efficiency"

	  tags = {
	    Name = "TerraformTest-disk"
	  }
	}

	resource "alicloud_instance" "default" {
		image_id = "${data.alicloud_images.default.images.0.id}"
		availability_zone = "${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"
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
	`)
}
func testAccDiskAttachmentConfigResize() string {
	return fmt.Sprintf(`
    data "alicloud_instance_types" "default" {
      cpu_core_count    = 2
      memory_size       = 4
    }
    data "alicloud_images" "default" {
	  # test for windows service
      name_regex  = "^win*"

      most_recent = true
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
    resource "alicloud_security_group_rule" "default" {
      	type = "ingress"
      	ip_protocol = "tcp"
      	nic_type = "intranet"
      	policy = "accept"
      	port_range = "22/22"
      	priority = 1
      	security_group_id = "${alicloud_security_group.default.id}"
      	cidr_ip = "172.16.0.0/24"
    }	

	variable "name" {
		default = "tf-testAccEcsDiskAttachmentConfig"
	}

	resource "alicloud_disk" "default" {
	  availability_zone = "${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"
	  size = "70"
	  name = "${var.name}"
	  category = "cloud_efficiency"

	  tags = {
	    Name = "TerraformTest-disk"
	  }
	}

	resource "alicloud_instance" "default" {
		image_id = "${data.alicloud_images.default.images.0.id}"
		availability_zone = "${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"
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
	`)
}
func testAccMultiDiskAttachmentConfig(common string) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEcsDiskAttachmentConfig"
	}

	variable "number" {
		default = "2"
	}

	resource "alicloud_disk" "default" {
		name = "${var.name}-${count.index}"
		count = "${var.number}"
		availability_zone = "${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"
		size = "50"

		tags = {
			Name = "TerraformTest-disk-${count.index}"
		}
	}

	resource "alicloud_instance" "default" {
		image_id = "${data.alicloud_images.default.images.0.id}"
		availability_zone = "${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"
		system_disk_category = "cloud_ssd"
		system_disk_size = 40
		instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
		security_groups = ["${alicloud_security_group.default.id}"]
		instance_name = "${var.name}"
		vswitch_id = "${alicloud_vswitch.default.id}"
	}

	resource "alicloud_disk_attachment" "default" {
		count = "${var.number}"
		disk_id     = "${element(alicloud_disk.default.*.id, count.index)}"
		instance_id = "${alicloud_instance.default.id}"
	}
	`, common)
}
