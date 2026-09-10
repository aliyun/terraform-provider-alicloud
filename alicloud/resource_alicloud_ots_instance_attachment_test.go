package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ots"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudOtsInstanceAttachmentBasic(t *testing.T) {
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
					"vswitch_id":    "${local.vswitch_id}",
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

func TestAccAliCloudOtsInstanceAttachmentHighPerformance(t *testing.T) {
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
					"vswitch_id":    "${local.vswitch_id}",
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
	  vpc_name   = "${var.name}"
	}

	data "alicloud_vswitches" "default" {
		vpc_id = alicloud_vpc.default.id
		zone_id      = data.alicloud_zones.default.zones.0.id
	}
	
	resource "alicloud_vswitch" "vswitch" {
	  vpc_id            = alicloud_vpc.default.id 
	  cidr_block        = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 8)
	  zone_id           = data.alicloud_zones.default.zones.0.id
	  vswitch_name      = var.name
	}
	
	locals {
	  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
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
	  vpc_name   = "${var.name}"
	}
	data "alicloud_vswitches" "default" {
		vpc_id = alicloud_vpc.default.id
		zone_id      = data.alicloud_zones.default.zones.0.id
	}
	
	resource "alicloud_vswitch" "vswitch" {
	  vpc_id            = alicloud_vpc.default.id
	  cidr_block        = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 8)
	  zone_id           = data.alicloud_zones.default.zones.0.id
	  vswitch_name      = var.name
	}
	
	locals {
	  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
	}
	`, name, string(OtsHighPerformance))
}

var otsInstanceAttachmentBasicMap = map[string]string{
	"instance_name": CHECKSET,
	"vpc_name":      CHECKSET,
	"vswitch_id":    CHECKSET,
	"vpc_id":        CHECKSET,
}
