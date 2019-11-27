package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ots"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudOtsInstanceAttachment_basic(t *testing.T) {
	var v ots.VpcInfo

	resourceId := "alicloud_ots_instance_attachment.default"
	ra := resourceAttrInit(resourceId, otsInstanceAttachmentBasicMap)

	serviceFunc := func() interface{} {
		return &OtsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceOtsInstanceAttachmentConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, false, connectivity.OtsCapacityNoSupportedRegions)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": "${alicloud_ots_instance.default.name}",
					"vpc_name":      "test",
					"vswitch_id":    "${alicloud_vswitch.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name,
						"vpc_name":      "test",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudOtsInstanceAttachment_highPerformance(t *testing.T) {
	var v ots.VpcInfo

	resourceId := "alicloud_ots_instance_attachment.default"
	ra := resourceAttrInit(resourceId, otsInstanceAttachmentBasicMap)

	serviceFunc := func() interface{} {
		return &OtsService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceOtsInstanceAttachmentConfigDependenceHighperformance)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, false, connectivity.OtsHighPerformanceNoSupportedRegions)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_name": "${alicloud_ots_instance.default.name}",
					"vpc_name":      "test",
					"vswitch_id":    "${alicloud_vswitch.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_name": name,
						"vpc_name":      "test",
					}),
				),
			},
		},
	})
}

func resourceOtsInstanceAttachmentConfigDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "%s"
	}

	resource "alicloud_ots_instance" "default" {
	  name = "${var.name}"
	  description = "${var.name}"
	  accessed_by = "Vpc"
	  instance_type = "%s"
	}

	data "alicloud_zones" "default" {
	  available_resource_creation = "VSwitch"
	}
	resource "alicloud_vpc" "default" {
	  cidr_block = "172.16.0.0/16"
	  name = "${var.name}"
	}

	resource "alicloud_vswitch" "default" {
	  vpc_id = "${alicloud_vpc.default.id}"
	  name = "${var.name}"
	  cidr_block = "172.16.1.0/24"
	  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	}
	`, name, string(OtsCapacity))
}

func resourceOtsInstanceAttachmentConfigDependenceHighperformance(name string) string {
	return fmt.Sprintf(`
	variable "name" {
	  default = "%s"
	}

	resource "alicloud_ots_instance" "default" {
	  name = "${var.name}"
	  description = "${var.name}"
	  accessed_by = "Vpc"
	  instance_type = "%s"
	}

	data "alicloud_zones" "default" {
	  available_resource_creation = "VSwitch"
	}
	resource "alicloud_vpc" "default" {
	  cidr_block = "172.16.0.0/16"
	  name = "${var.name}"
	}

	resource "alicloud_vswitch" "default" {
	  vpc_id = "${alicloud_vpc.default.id}"
	  name = "${var.name}"
	  cidr_block = "172.16.1.0/24"
	  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	}
	`, name, string(OtsHighPerformance))
}

var otsInstanceAttachmentBasicMap = map[string]string{
	"instance_name": CHECKSET,
	"vpc_name":      CHECKSET,
	"vswitch_id":    CHECKSET,
	"vpc_id":        CHECKSET,
}
