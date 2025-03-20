package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudSLBServerGroupServerAttachment_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_slb_server_group_server_attachment.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlbServerGroupServerAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccSlbServerGroupServerAttachment%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudSlbServerGroupServerAttachmentBasicDependence0)
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
					"server_group_id": "${alicloud_slb_server_group.default.id}",
					"server_id":       "${alicloud_instance.default.id}",
					"port":            "8080",
					"type":            "ecs",
					"weight":          "0",
					"description":     name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_group_id": CHECKSET,
						"server_id":       CHECKSET,
						"port":            "8080",
						"type":            "ecs",
						"weight":          "0",
						"description":     name,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func AliCloudSlbServerGroupServerAttachmentBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_slb_zones" "default" {
  		available_slb_address_type = "vpc"
	}

	data "alicloud_instance_types" "default" {
  		availability_zone    = data.alicloud_slb_zones.default.zones.0.id
  		instance_type_family = "ecs.sn1ne"
	}

	data "alicloud_images" "default" {
  		name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  		most_recent = true
  		owners      = "system"
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "192.168.0.0/16"
	}

	resource "alicloud_vswitch" "default" {
  		vswitch_name = var.name
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "192.168.192.0/24"
  		zone_id      = data.alicloud_slb_zones.default.zones.0.id
	}

	resource "alicloud_security_group" "default" {
  		name   = var.name
  		vpc_id = alicloud_vpc.default.id
	}

	resource "alicloud_slb_load_balancer" "default" {
  		load_balancer_name = var.name
  		vswitch_id         = alicloud_vswitch.default.id
  		load_balancer_spec = "slb.s2.small"
  		address_type       = "intranet"
	}

	resource "alicloud_slb_server_group" "default" {
  		load_balancer_id = alicloud_slb_load_balancer.default.id
  		name             = var.name
	}

	resource "alicloud_instance" "default" {
  		image_id                   = data.alicloud_images.default.images[0].id
  		instance_type              = data.alicloud_instance_types.default.instance_types[0].id
  		instance_name              = var.name
  		security_groups            = alicloud_security_group.default.*.id
  		internet_charge_type       = "PayByTraffic"
  		internet_max_bandwidth_out = "10"
  		availability_zone          = data.alicloud_slb_zones.default.zones.0.id
  		instance_charge_type       = "PostPaid"
  		system_disk_category       = "cloud_efficiency"
  		vswitch_id                 = alicloud_vswitch.default.id
	}
`, name)
}

func TestAccAliCloudSLBServerGroupServerAttachment_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_slb_server_group_server_attachment.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlbServerGroupServerAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAccSlbServerGroupServerAttachment%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudSlbServerGroupServerAttachmentBasicDependence1)
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
					"server_group_id": "${alicloud_slb_server_group.default.id}",
					"server_id":       "${alicloud_ecs_network_interface.default.id}",
					"port":            "8080",
					"type":            "eni",
					"weight":          "10",
					"description":     name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_group_id": CHECKSET,
						"server_id":       CHECKSET,
						"port":            "8080",
						"type":            "eni",
						"weight":          "10",
						"description":     name,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func AliCloudSlbServerGroupServerAttachmentBasicDependence1(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_slb_zones" "default" {
  		available_slb_address_type = "vpc"
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "192.168.0.0/16"
	}

	resource "alicloud_vswitch" "default" {
  		vswitch_name = var.name
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "192.168.192.0/24"
  		zone_id      = data.alicloud_slb_zones.default.zones.0.id
	}

	resource "alicloud_security_group" "default" {
  		name   = var.name
  		vpc_id = alicloud_vpc.default.id
	}

	resource "alicloud_slb_load_balancer" "default" {
  		load_balancer_name = var.name
  		vswitch_id         = alicloud_vswitch.default.id
  		load_balancer_spec = "slb.s2.small"
  		address_type       = "intranet"
	}

	resource "alicloud_slb_server_group" "default" {
  		load_balancer_id = alicloud_slb_load_balancer.default.id
  		name             = var.name
	}

	resource "alicloud_ecs_network_interface" "default" {
		network_interface_name = var.name
  		vswitch_id             = alicloud_vswitch.default.id
  		security_group_ids     = [alicloud_security_group.default.id]
	}
`, name)
}

func TestAccAliCloudSLBServerGroupServerAttachment_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_slb_server_group_server_attachment.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlbServerGroupServerAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testccslb%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudSlbServerGroupServerAttachmentBasicDependence2)
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
					"server_group_id": "${alicloud_slb_server_group.default.id}",
					"server_id":       "${alicloud_eci_container_group.default.id}",
					"port":            "80",
					"type":            "eci",
					"weight":          "100",
					"description":     name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"server_group_id": CHECKSET,
						"server_id":       CHECKSET,
						"port":            "80",
						"type":            "eci",
						"weight":          "100",
						"description":     name,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func AliCloudSlbServerGroupServerAttachmentBasicDependence2(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_slb_zones" "default" {
  		available_slb_address_type = "vpc"
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "192.168.0.0/16"
	}

	resource "alicloud_vswitch" "default" {
  		vswitch_name = var.name
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "192.168.192.0/24"
  		zone_id      = data.alicloud_slb_zones.default.zones.0.id
	}

	resource "alicloud_security_group" "default" {
  		name   = var.name
  		vpc_id = alicloud_vpc.default.id
	}

	resource "alicloud_slb_load_balancer" "default" {
  		load_balancer_name = var.name
  		vswitch_id         = alicloud_vswitch.default.id
  		load_balancer_spec = "slb.s2.small"
  		address_type       = "intranet"
	}

	resource "alicloud_slb_server_group" "default" {
  		load_balancer_id = alicloud_slb_load_balancer.default.id
  		name             = var.name
	}

	resource "alicloud_eci_container_group" "default" {
  		container_group_name = var.name
  		vswitch_id           = alicloud_vswitch.default.id
		security_group_id    = alicloud_security_group.default.id
  		containers {
    		image             = "registry-vpc.%s.aliyuncs.com/eci_open/nginx:alpine"
    		name              = "nginx"
    		working_dir       = "/tmp/nginx"
    		image_pull_policy = "IfNotPresent"
    		commands          = ["/bin/sh", "-c", "sleep 9999"]
			volume_mounts {
      			mount_path = "/tmp/test"
      			read_only  = false
      			name       = "empty1"
			}
    		ports {
      			port     = 80
      			protocol = "TCP"
    		}
    		environment_vars {
      			key   = "test"
      			value = "nginx"
    		}
  		}
  		volumes {
    		name = "empty1"
    		type = "EmptyDirVolume"
  		}
	}
`, name, defaultRegionToTest)
}
