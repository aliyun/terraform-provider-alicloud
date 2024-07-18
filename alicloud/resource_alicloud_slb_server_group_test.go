package alicloud

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudSlbServerGroup_basic0(t *testing.T) {
	var v *slb.DescribeVServerGroupAttributeResponse
	resourceId := "alicloud_slb_server_group.default"
	ra := resourceAttrInit(resourceId, AliCloudSlbServerGroupMap0)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sSlbServerGroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudGaSlbServerGroupBasicDependence0)
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
					"load_balancer_id": "${alicloud_slb_load_balancer.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "ServerGroup",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "ServerGroup",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"servers": serversMap0,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"servers.#": "16",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"servers": serversMap2,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"servers.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"servers": serversMap1,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"servers.#": "10",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_protection_validation"},
			},
		},
	})
}

func TestAccAliCloudSlbServerGroup_basic0_twin(t *testing.T) {
	var v *slb.DescribeVServerGroupAttributeResponse
	resourceId := "alicloud_slb_server_group.default"
	ra := resourceAttrInit(resourceId, AliCloudSlbServerGroupMap0)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sSlbServerGroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudGaSlbServerGroupBasicDependence0)
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
					"load_balancer_id": "${alicloud_slb_load_balancer.default.id}",
					"name":             name,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "ServerGroup",
					},
					"delete_protection_validation": "false",
					"servers":                      serversMap0,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_id": CHECKSET,
						"name":             name,
						"tags.%":           "2",
						"tags.Created":     "TF",
						"tags.For":         "ServerGroup",
						"servers.#":        "16",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_protection_validation"},
			},
		},
	})
}

func TestAccAliCloudSlbServerGroup_basic1(t *testing.T) {
	var v *slb.DescribeVServerGroupAttributeResponse
	resourceId := "alicloud_slb_server_group.default"
	ra := resourceAttrInit(resourceId, AliCloudSlbServerGroupMap0)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sSlbServerGroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudGaSlbServerGroupBasicDependence1)
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
					"load_balancer_id": "${alicloud_slb_load_balancer.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "ServerGroup",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "ServerGroup",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"servers": []map[string]interface{}{
						{
							"type":       "eni",
							"port":       "80",
							"weight":     "10",
							"server_ids": []string{"${alicloud_ecs_network_interface.default.id}"},
						},
						{
							"port":       "100",
							"weight":     "10",
							"server_ids": []string{"${alicloud_instance.default.id}"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"servers.#": "2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_protection_validation"},
			},
		},
	})
}

func TestAccAliCloudSlbServerGroup_basic1_twin(t *testing.T) {
	var v *slb.DescribeVServerGroupAttributeResponse
	resourceId := "alicloud_slb_server_group.default"
	ra := resourceAttrInit(resourceId, AliCloudSlbServerGroupMap0)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sSlbServerGroup%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudGaSlbServerGroupBasicDependence1)
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
					"load_balancer_id": "${alicloud_slb_load_balancer.default.id}",
					"name":             name,
					"tags": map[string]string{
						"Created": "TF",
						"For":     "ServerGroup",
					},
					"delete_protection_validation": "false",
					"servers": []map[string]interface{}{
						{
							"type":       "eni",
							"port":       "80",
							"weight":     "10",
							"server_ids": []string{"${alicloud_ecs_network_interface.default.id}"},
						},
						{
							"port":       "100",
							"weight":     "10",
							"server_ids": []string{"${alicloud_instance.default.id}"},
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_id": CHECKSET,
						"name":             name,
						"tags.%":           "2",
						"tags.Created":     "TF",
						"tags.For":         "ServerGroup",
						"servers.#":        "2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_protection_validation"},
			},
		},
	})
}

func AliCloudGaSlbServerGroupBasicDependence0(name string) string {
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

	resource "alicloud_instance" "default" {
  		count                      = 16
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

	resource "alicloud_slb_load_balancer" "default" {
  		load_balancer_name = var.name
  		load_balancer_spec = "slb.s1.small"
	}
`, name)

}

func AliCloudGaSlbServerGroupBasicDependence1(name string) string {
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

	resource "alicloud_ecs_network_interface" "default" {
		network_interface_name = var.name
  		vswitch_id             = alicloud_vswitch.default.id
  		security_group_ids     = [alicloud_security_group.default.id]
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

	resource "alicloud_slb_load_balancer" "default" {
  		load_balancer_name = var.name
  		vswitch_id         = alicloud_vswitch.default.id
  		load_balancer_spec = "slb.s2.small"
	}
`, name)
}

var AliCloudSlbServerGroupMap0 = map[string]string{}

var serversMap0 = []map[string]interface{}{
	{
		"port":       "1",
		"weight":     "10",
		"server_ids": []string{"${alicloud_instance.default.0.id}"},
	},
	{
		"port":       "2",
		"weight":     "10",
		"server_ids": []string{"${alicloud_instance.default.1.id}"},
	},
	{
		"port":       "3",
		"weight":     "10",
		"server_ids": []string{"${alicloud_instance.default.2.id}"},
	},
	{
		"port":       "4",
		"weight":     "10",
		"server_ids": []string{"${alicloud_instance.default.3.id}"},
	},
	{
		"port":       "5",
		"weight":     "10",
		"server_ids": []string{"${alicloud_instance.default.4.id}"},
	},
	{
		"port":       "6",
		"weight":     "10",
		"server_ids": []string{"${alicloud_instance.default.5.id}"},
	},
	{
		"port":       "7",
		"weight":     "10",
		"server_ids": []string{"${alicloud_instance.default.6.id}"},
	},
	{
		"port":       "8",
		"weight":     "10",
		"server_ids": []string{"${alicloud_instance.default.7.id}"},
	},
	{
		"port":       "9",
		"weight":     "10",
		"server_ids": []string{"${alicloud_instance.default.8.id}"},
	},
	{
		"port":       "10",
		"weight":     "10",
		"server_ids": []string{"${alicloud_instance.default.9.id}"},
	},
	{
		"port":       "11",
		"weight":     "10",
		"server_ids": []string{"${alicloud_instance.default.10.id}"},
	},
	{
		"port":       "12",
		"weight":     "10",
		"server_ids": []string{"${alicloud_instance.default.11.id}"},
	},
	{
		"port":       "13",
		"weight":     "10",
		"server_ids": []string{"${alicloud_instance.default.12.id}"},
	},
	{
		"port":       "14",
		"weight":     "10",
		"server_ids": []string{"${alicloud_instance.default.13.id}"},
	},
	{
		"port":       "15",
		"weight":     "10",
		"server_ids": []string{"${alicloud_instance.default.14.id}"},
	},
	{
		"port":       "16",
		"weight":     "10",
		"server_ids": []string{"${alicloud_instance.default.15.id}"},
	},
}

var serversMap1 = []map[string]interface{}{
	{
		"port":       "1",
		"weight":     "10",
		"server_ids": []string{"${alicloud_instance.default.0.id}"},
	},
	{
		"port":       "2",
		"weight":     "10",
		"server_ids": []string{"${alicloud_instance.default.1.id}"},
	},
	{
		"port":       "3",
		"weight":     "10",
		"server_ids": []string{"${alicloud_instance.default.2.id}"},
	},
	{
		"port":       "4",
		"weight":     "10",
		"server_ids": []string{"${alicloud_instance.default.3.id}"},
	},
	{
		"port":       "5",
		"weight":     "10",
		"server_ids": []string{"${alicloud_instance.default.4.id}"},
	},
	{
		"port":       "6",
		"weight":     "10",
		"server_ids": []string{"${alicloud_instance.default.5.id}"},
	},
	{
		"port":       "7",
		"weight":     "10",
		"server_ids": []string{"${alicloud_instance.default.6.id}"},
	},
	{
		"port":       "8",
		"weight":     "10",
		"server_ids": []string{"${alicloud_instance.default.7.id}"},
	},
	{
		"port":       "9",
		"weight":     "10",
		"server_ids": []string{"${alicloud_instance.default.8.id}"},
	},
	{
		"port":       "10",
		"weight":     "10",
		"server_ids": []string{"${alicloud_instance.default.9.id}"},
	},
}

var serversMap2 = []map[string]interface{}{
	{
		"port":       "15",
		"weight":     "10",
		"server_ids": []string{"${alicloud_instance.default.14.id}"},
	},
}
