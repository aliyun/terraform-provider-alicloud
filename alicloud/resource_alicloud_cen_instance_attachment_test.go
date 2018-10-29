package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cbn"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudCenInstanceAttachment_basic(t *testing.T) {
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
				Config: testAccCenInstanceAttachmentBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenInstanceAttachmentExists("alicloud_cen_instance_attachment.foo", &instance),
					resource.TestCheckResourceAttr("alicloud_cen_instance_attachment.foo", "child_instance_region_id", "cn-hangzhou"),
				),
			},
		},
	})
}

func TestAccAlicloudCenInstanceAttachment_multi_same_regions(t *testing.T) {
	var instance cbn.ChildInstance

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCenInstanceAttachmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCenInstanceAttachmentMultiSameRegions,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenInstanceAttachmentExists("alicloud_cen_instance_attachment.bar1", &instance),
					resource.TestCheckResourceAttr("alicloud_cen_instance_attachment.bar1", "child_instance_region_id", "cn-beijing"),
					testAccCheckCenInstanceAttachmentExists("alicloud_cen_instance_attachment.bar2", &instance),
					resource.TestCheckResourceAttr("alicloud_cen_instance_attachment.bar2", "child_instance_region_id", "cn-beijing"),
				),
			},
		},
	})
}

func TestAccAlicloudCenInstanceAttachment_multi_different_regions(t *testing.T) {
	var instance cbn.ChildInstance

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCenInstanceAttachmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCenInstanceAttachmentMultiDifferentRegions,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCenInstanceAttachmentExists("alicloud_cen_instance_attachment.bar1", &instance),
					resource.TestCheckResourceAttr("alicloud_cen_instance_attachment.bar1", "child_instance_region_id", "eu-central-1"),
					testAccCheckCenInstanceAttachmentExists("alicloud_cen_instance_attachment.bar2", &instance),
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

		client := testAccProvider.Meta().(*connectivity.AliyunClient)
		cenService := CenService{client}

		cenId, instanceId, err := cenService.GetCenIdAndAnotherId(rs.Primary.ID)
		if err != nil {
			return err
		}

		childInstance, err := cenService.DescribeCenAttachedChildInstanceById(instanceId, cenId)
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
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	cenService := CenService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_instance_attachment" {
			continue
		}

		cenId, instanceId, err := cenService.GetCenIdAndAnotherId(rs.Primary.ID)
		if err != nil {
			return err
		}

		instance, err := cenService.DescribeCenAttachedChildInstanceById(instanceId, cenId)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		} else {
			return fmt.Errorf("CEN %s child instance %s still attach", instance.CenId, instance.ChildInstanceId)
		}
	}

	return nil
}

const testAccCenInstanceAttachmentBasic = `
provider "alicloud" {
    alias = "hz"
    region = "cn-hangzhou"
}

variable "name"{
    default = "tf-testAccCenInstanceAttachmentBasic"
}

resource "alicloud_cen_instance" "cen" {
    name = "${var.name}"
    description = "tf-testAccCenInstanceAttachmentBasicDescription"
}

resource "alicloud_vpc" "vpc" {
    provider = "alicloud.hz"
    name = "${var.name}"
    cidr_block = "192.168.0.0/16"
}

resource "alicloud_cen_instance_attachment" "foo" {
    instance_id = "${alicloud_cen_instance.cen.id}"
    child_instance_id = "${alicloud_vpc.vpc.id}"
    child_instance_region_id = "cn-hangzhou"
}
`

const testAccCenInstanceAttachmentMultiSameRegions = `
provider "alicloud" {
    alias = "bj"
    region = "cn-beijing"
}

variable "name"{
    default = "tf-testAccCenInstanceAttachmentMultiSameRegions"
}

resource "alicloud_cen_instance" "cen" {
    name = "${var.name}"
    description = "tf-testAccCenInstanceAttachmentMultiSameRegionsDescription"
}

resource "alicloud_vpc" "vpc1" {
    provider = "alicloud.bj"
    name = "${var.name}"
    cidr_block = "192.168.0.0/16"
}

resource "alicloud_vpc" "vpc2" {
    provider = "alicloud.bj"
    name = "${var.name}"
    cidr_block = "172.16.0.0/12"
}

resource "alicloud_cen_instance_attachment" "bar1" {
    instance_id = "${alicloud_cen_instance.cen.id}"
    child_instance_id = "${alicloud_vpc.vpc1.id}"
    child_instance_region_id = "cn-beijing"
}

resource "alicloud_cen_instance_attachment" "bar2" {
    instance_id = "${alicloud_cen_instance.cen.id}"
    child_instance_id = "${alicloud_vpc.vpc2.id}"
    child_instance_region_id = "cn-beijing"
}
`

const testAccCenInstanceAttachmentMultiDifferentRegions = `
variable "name"{
    default = "tf-testAccCenInstanceAttachmentMultiDifferentRegions"
}

provider "alicloud" {
    alias = "fra"
    region = "eu-central-1"
}

provider "alicloud" {
    alias = "sh"
    region = "cn-shanghai"
}

resource "alicloud_vpc" "vpc1" {
    provider = "alicloud.fra"
    name = "${var.name}"
    cidr_block = "192.168.0.0/16"
}

resource "alicloud_vpc" "vpc2" {
    provider = "alicloud.sh"
    name = "${var.name}"
    cidr_block = "172.16.0.0/12"
}

resource "alicloud_cen_instance" "cen" {
    name = "${var.name}"
    description = "tf-testAccCenInstanceAttachmentMultiDifferentRegionsDescription"
}

resource "alicloud_cen_instance_attachment" "bar1" {
    instance_id = "${alicloud_cen_instance.cen.id}"
    child_instance_id = "${alicloud_vpc.vpc1.id}"
    child_instance_region_id = "eu-central-1"
}

resource "alicloud_cen_instance_attachment" "bar2" {
    instance_id = "${alicloud_cen_instance.cen.id}"
    child_instance_id = "${alicloud_vpc.vpc2.id}"
    child_instance_region_id = "cn-shanghai"
}
`
