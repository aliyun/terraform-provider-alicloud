// Package alicloud. This file is generated automatically. Please do not modify it manually, thank you!
package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test Ecs RamRoleAttachment. >>> Resource test cases, automatically generated.
// Case RamRoleAttachment 10656
func TestAccAliCloudEcsRamRoleAttachment_basic10656(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_ram_role_attachment.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsRamRoleAttachmentMap10656)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsRamRoleAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccecs%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsRamRoleAttachmentBasicDependence10656)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"ram_role_name": "${alicloud_ram_role.default.id}",
					"instance_id":   "${alicloud_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ram_role_name": CHECKSET,
						"instance_id":   CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"policy"},
			},
		},
	})
}

func TestAccAliCloudEcsRamRoleAttachment_basic10656_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_ecs_ram_role_attachment.default"
	ra := resourceAttrInit(resourceId, AliCloudEcsRamRoleAttachmentMap10656)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEcsRamRoleAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tfaccecs%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEcsRamRoleAttachmentBasicDependence10656)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"policy":        `{\"Statement\": [{\"Action\": [\"*\"],\"Effect\": \"Allow\",\"Resource\": [\"*\"]}],\"Version\":\"1\"}`,
					"ram_role_name": "${alicloud_ram_role.default.id}",
					"instance_id":   "${alicloud_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ram_role_name": CHECKSET,
						"instance_id":   CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"policy"},
			},
		},
	})
}

var AliCloudEcsRamRoleAttachmentMap10656 = map[string]string{}

func AliCloudEcsRamRoleAttachmentBasicDependence10656(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_zones" "default" {
  		available_disk_category     = "cloud_efficiency"
  		available_resource_creation = "VSwitch"
	}

	data "alicloud_images" "default" {
  		most_recent = true
  		owners      = "system"
	}

	data "alicloud_instance_types" "default" {
  		availability_zone = data.alicloud_zones.default.zones.0.id
  		image_id          = data.alicloud_images.default.images.0.id
	}

	resource "alicloud_ram_role" "default" {
  		name     = var.name
  		document = <<EOF
		{
			"Statement": [
				{
					"Action": "sts:AssumeRole",
					"Effect": "Allow",
					"Principal": {
						"Service": [
							"ecs.aliyuncs.com"
						]
					}
				}
		  	],
			"Version": "1"
		}
	  	EOF
		force    = true
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "192.168.0.0/16"
	}

	resource "alicloud_vswitch" "default" {
  		vswitch_name = var.name
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "192.168.192.0/24"
  		zone_id      = data.alicloud_zones.default.zones.0.id
	}

	resource "alicloud_security_group" "default" {
  		vpc_id = alicloud_vpc.default.id
	}

	resource "alicloud_instance" "default" {
  		image_id                   = data.alicloud_images.default.images.0.id
  		instance_type              = data.alicloud_instance_types.default.instance_types.0.id
  		security_groups            = alicloud_security_group.default.*.id
  		internet_charge_type       = "PayByTraffic"
  		internet_max_bandwidth_out = "10"
  		availability_zone          = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  		instance_charge_type       = "PostPaid"
  		system_disk_category       = "cloud_efficiency"
  		vswitch_id                 = alicloud_vswitch.default.id
  		instance_name              = var.name
	}
`, name)
}

// Test Ecs RamRoleAttachment. <<< Resource test cases, automatically generated.
