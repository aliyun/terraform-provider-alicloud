package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudEipAssociation_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eip_association.default"
	ra := resourceAttrInit(resourceId, AliCloudEipAssociationMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EipServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEipAssociation")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sEipAssociation%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEipAssociationBasicDependence0)
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
					"allocation_id": "${alicloud_eip_address.default.id}",
					"instance_id":   "${alicloud_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allocation_id": CHECKSET,
						"instance_id":   CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force"},
			},
		},
	})
}

func TestAccAliCloudEipAssociation_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eip_association.default"
	ra := resourceAttrInit(resourceId, AliCloudEipAssociationMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EipServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEipAssociation")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sEipAssociation%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEipAssociationBasicDependence0)
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
					"allocation_id": "${alicloud_eip_address.default.id}",
					"instance_id":   "${alicloud_instance.default.id}",
					"instance_type": "EcsInstance",
					"force":         "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allocation_id": CHECKSET,
						"instance_id":   CHECKSET,
						"instance_type": "EcsInstance",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force"},
			},
		},
	})
}

func TestAccAliCloudEipAssociation_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eip_association.default"
	ra := resourceAttrInit(resourceId, AliCloudEipAssociationMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EipServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEipAssociation")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sEipAssociation%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEipAssociationBasicDependence1)
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
					"allocation_id": "${alicloud_eip_address.default.id}",
					"instance_id":   "${alicloud_ecs_network_interface.default.id}",
					"instance_type": "NetworkInterface",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allocation_id": CHECKSET,
						"instance_id":   CHECKSET,
						"instance_type": "NetworkInterface",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force"},
			},
		},
	})
}

func TestAccAliCloudEipAssociation_basic1_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eip_association.default"
	ra := resourceAttrInit(resourceId, AliCloudEipAssociationMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EipServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEipAssociation")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sEipAssociation%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEipAssociationBasicDependence1)
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
					"allocation_id":      "${alicloud_eip_address.default.id}",
					"instance_id":        "${alicloud_ecs_network_interface.default.id}",
					"instance_type":      "NetworkInterface",
					"mode":               "NAT",
					"private_ip_address": "${tolist(alicloud_ecs_network_interface.default.private_ip_addresses).0}",
					"force":              "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allocation_id":      CHECKSET,
						"instance_id":        CHECKSET,
						"instance_type":      "NetworkInterface",
						"mode":               "NAT",
						"private_ip_address": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"mode": "BINDED",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"mode": "BINDED",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force"},
			},
		},
	})
}

func TestAccAliCloudEipAssociation_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eip_association.default"
	ra := resourceAttrInit(resourceId, AliCloudEipAssociationMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EipServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEipAssociation")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sEipAssociation%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEipAssociationBasicDependence2)
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
					"allocation_id": "${alicloud_eip_address.default.id}",
					"instance_id":   "192.168.0.1",
					"instance_type": "IpAddress",
					"vpc_id":        "${alicloud_vpc_ipv4_gateway.default.vpc_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allocation_id": CHECKSET,
						"instance_id":   CHECKSET,
						"instance_type": "IpAddress",
						"vpc_id":        CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force"},
			},
		},
	})
}

func TestAccAliCloudEipAssociation_basic2_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_eip_association.default"
	ra := resourceAttrInit(resourceId, AliCloudEipAssociationMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EipServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeEipAssociation")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testAcc%sEipAssociation%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudEipAssociationBasicDependence2)
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
					"allocation_id": "${alicloud_eip_address.default.id}",
					"instance_id":   "192.168.0.1",
					"instance_type": "IpAddress",
					"vpc_id":        "${alicloud_vpc_ipv4_gateway.default.vpc_id}",
					"force":         "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"allocation_id": CHECKSET,
						"instance_id":   CHECKSET,
						"instance_type": "IpAddress",
						"vpc_id":        CHECKSET,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force"},
			},
		},
	})
}

var AliCloudEipAssociationMap0 = map[string]string{
	"instance_type": CHECKSET,
	"mode":          CHECKSET,
}

func AliCloudEipAssociationBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_zones" "default" {
  		available_disk_category     = "cloud_efficiency"
  		available_resource_creation = "VSwitch"
	}

	data "alicloud_instance_types" "default" {
  		availability_zone    = data.alicloud_zones.default.zones.0.id
  		instance_type_family = "ecs.sn1ne"
	}

	data "alicloud_images" "default" {
  		name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  		most_recent = true
		instance_type = data.alicloud_instance_types.default.instance_types.0.id
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
  		zone_id      = data.alicloud_zones.default.zones.0.id
	}

	resource "alicloud_security_group" "default" {
  		name   = var.name
  		vpc_id = alicloud_vpc.default.id
	}

	resource "alicloud_eip_address" "default" {
  		address_name = var.name
	}

	resource "alicloud_instance" "default" {
  		image_id             = data.alicloud_images.default.images.0.id
  		instance_type        = data.alicloud_instance_types.default.instance_types.0.id
  		instance_name        = var.name
  		security_groups      = alicloud_security_group.default.*.id
  		availability_zone    = data.alicloud_zones.default.zones.0.id
  		instance_charge_type = "PostPaid"
  		system_disk_category = "cloud_efficiency"
  		vswitch_id           = alicloud_vswitch.default.id
	}
`, name)
}

func AliCloudEipAssociationBasicDependence1(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_zones" "default" {
  		available_resource_creation = "VSwitch"
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
  		name   = var.name
  		vpc_id = alicloud_vpc.default.id
	}

	resource "alicloud_eip_address" "default" {
  		address_name = var.name
	}

	resource "alicloud_ecs_network_interface" "default" {
  		network_interface_name = var.name
  		vswitch_id             = alicloud_vswitch.default.id
  		security_group_ids     = [alicloud_security_group.default.id]
  		private_ip_addresses   = [cidrhost(alicloud_vswitch.default.cidr_block, 100)]
	}
`, name)
}

func AliCloudEipAssociationBasicDependence2(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	resource "alicloud_vpc" "default" {
	  cidr_block = "192.168.0.0/16"
	  name = var.name
	}

	resource "alicloud_eip_address" "default" {
  		address_name = var.name
	}

	resource "alicloud_vpc_ipv4_gateway" "default" {
  		vpc_id                   = alicloud_vpc.default.id
  		ipv4_gateway_name        = var.name
  		ipv4_gateway_description = var.name
        enabled                  = true
	}
`, name)
}
