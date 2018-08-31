package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudCen_Instance_Attachment_basic(t *testing.T) {
	var instance cbn.ChildInstance

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_cen_instance_attachment.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckCenInstanceAttachmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCenInstanceAttachmentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenInstanceAttachmentExists("alicloud_cen_instance_attachment.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_cen_instance_attachment.foo", "child_instance_type", "VPC"),
					resource.TestCheckResourceAttr("alicloud_cen_instance_attachment.foo", "child_instance_region_id", "cn-beijing"),
				),
			},
		},
	})
}

func TestAccAlicloudCen_Instance_Attachment_multi(t *testing.T) {
	var instance cbn.ChildInstance

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCenInstanceAttachmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCenInstanceAttachmentMulti,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenInstanceAttachmentExists("alicloud_cen_instance_attachment.bar1", &instance),
					resource.TestCheckResourceAttr("alicloud_cen_instance_attachment.bar1", "child_instance_type", "VPC"),
					resource.TestCheckResourceAttr("alicloud_cen_instance_attachment.bar1", "child_instance_region_id", "cn-beijing"),
					testAccCheckCenInstanceAttachmentExists("alicloud_cen_instance_attachment.bar2", &instance),
					resource.TestCheckResourceAttr("alicloud_cen_instance_attachment.bar2", "child_instance_type", "VPC"),
					resource.TestCheckResourceAttr("alicloud_cen_instance_attachment.bar2", "child_instance_region_id", "cn-shanghai"),
				),
			},
		},
	})
}

func testAccCheckCenInstanceAttachmentExists(n string, instance *cbn.ChildInstance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Cen Child Instance ID is set")
		}

		client := testAccProvider.Meta().(*AliyunClient)

		instanceId, cenId, err := getCenIdAndAnotherId(rs.Primary.ID)
		if err != nil {
			return err
		}

		childInstance, err := client.DescribeCenAttachedChildInstanceById(instanceId, cenId)
		if err != nil {
			return err
		}

		if childInstance.Status != "Attached" {
			return fmt.Errorf("CEN id %s instance id %s status error", cenId, instanceId)
		}

		*instance = childInstance
		return nil
	}
}

func testAccCheckCenInstanceAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*AliyunClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_cen_instance_attachment" {
			continue
		}

		instanceId, cenId, err := getCenIdAndAnotherId(rs.Primary.ID)
		if err != nil {
			return err
		}

		instance, err := client.DescribeCenAttachedChildInstanceById(instanceId, cenId)

		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}

		if instance.ChildInstanceId != "" {
			return fmt.Errorf("CEN %s child instance %s still attach", instance.CenId, instance.ChildInstanceId)
		}
	}

	return nil
}

const testAccCenInstanceAttachmentConfig = `
resource "alicloud_cen" "cen" {
	name = "terraform-01"
	description = "terraform01"
}

resource "alicloud_vpc" "vpc" {
	name = "terraform-01"
	cidr_block = "192.168.0.0/16"
}

resource "alicloud_cen_instance_attachment" "foo" {
    cen_id = "${alicloud_cen.cen.id}"
    child_instance_id = "${alicloud_vpc.vpc.id}"
    child_instance_type = "VPC"
    child_instance_region_id = "cn-beijing"
}
`

const testAccCenInstanceAttachmentMulti = `
provider "alicloud" {
	alias = "bj"
	region = "cn-beijing"
}

provider "alicloud" {
	alias = "sh"
	region = "cn-shanghai"
}

resource "alicloud_vpc" "vpc1" {
	provider = "alicloud.bj"
  	name = "terraform-01"
  	cidr_block = "192.168.0.0/16"
}

resource "alicloud_vpc" "vpc2" {
  	provider = "alicloud.sh"
  	name = "terraform-02"
  	cidr_block = "172.16.0.0/12"
}

resource "alicloud_cen" "cen" {
	name = "terraform-yl-01"
	description = "terraform01"
}

resource "alicloud_cen_instance_attachment" "bar1" {
    cen_id = "${alicloud_cen.cen.id}"
    child_instance_id = "${alicloud_vpc.vpc1.id}"
    child_instance_type = "VPC"
    child_instance_region_id = "cn-beijing"
}

resource "alicloud_cen_instance_attachment" "bar2" {
    cen_id = "${alicloud_cen.cen.id}"
    child_instance_id = "${alicloud_vpc.vpc2.id}"
    child_instance_type = "VPC"
    child_instance_region_id = "cn-shanghai"
}
`
