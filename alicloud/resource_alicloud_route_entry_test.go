package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAliCloudRouteEntry_basic0(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_route_entry.default"
	ra := resourceAttrInit(resourceId, AliCloudRouteEntryMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRouteEntry")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAcc-RouteEntry-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRouteEntryBasicDependence0)
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
					"route_table_id":        "${alicloud_vpc.default.route_table_id}",
					"destination_cidrblock": "172.11.1.1/32",
					"nexthop_type":          "Instance",
					"nexthop_id":            "${alicloud_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"route_table_id":        CHECKSET,
						"destination_cidrblock": "172.11.1.1/32",
						"nexthop_type":          "Instance",
						"nexthop_id":            CHECKSET,
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

func TestAccAliCloudRouteEntry_basic0_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_route_entry.default"
	ra := resourceAttrInit(resourceId, AliCloudRouteEntryMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRouteEntry")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAcc-RouteEntry-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRouteEntryBasicDependence0)
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
					"route_table_id":        "${alicloud_vpc.default.route_table_id}",
					"destination_cidrblock": "172.11.1.1/32",
					"nexthop_type":          "Instance",
					"nexthop_id":            "${alicloud_instance.default.id}",
					"name":                  name,
					"description":           name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"route_table_id":        CHECKSET,
						"destination_cidrblock": "172.11.1.1/32",
						"nexthop_type":          "Instance",
						"nexthop_id":            CHECKSET,
						"name":                  name,
						"description":           name,
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

func TestAccAliCloudRouteEntry_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_route_entry.default"
	ra := resourceAttrInit(resourceId, AliCloudRouteEntryMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRouteEntry")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAcc-RouteEntry-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRouteEntryBasicDependence1)
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
					"route_table_id":        "${alicloud_vpc.default.route_table_id}",
					"destination_cidrblock": "2001:ffff:ffff:ffff::/64",
					"nexthop_type":          "Instance",
					"nexthop_id":            "${alicloud_instance.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"route_table_id":        CHECKSET,
						"destination_cidrblock": "2001:ffff:ffff:ffff::/64",
						"nexthop_type":          "Instance",
						"nexthop_id":            CHECKSET,
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

func TestAccAliCloudRouteEntry_basic1_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_route_entry.default"
	ra := resourceAttrInit(resourceId, AliCloudRouteEntryMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRouteEntry")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAcc-RouteEntry-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRouteEntryBasicDependence1)
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
					"route_table_id":        "${alicloud_vpc.default.route_table_id}",
					"destination_cidrblock": "2001:ffff:ffff:ffff::/64",
					"nexthop_type":          "Instance",
					"nexthop_id":            "${alicloud_instance.default.id}",
					"name":                  name,
					"description":           name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"route_table_id":        CHECKSET,
						"destination_cidrblock": "2001:ffff:ffff:ffff::/64",
						"nexthop_type":          "Instance",
						"nexthop_id":            CHECKSET,
						"name":                  name,
						"description":           name,
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

func TestAccAliCloudRouteEntry_basic2(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_route_entry.default"
	ra := resourceAttrInit(resourceId, AliCloudRouteEntryMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRouteEntry")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAcc-RouteEntry-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRouteEntryBasicDependence2)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithRegions(t, true, []connectivity.Region{"cn-hangzhou"})
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"route_table_id":        "${alicloud_vpc.defaultVpc.route_table_id}",
					"destination_cidrblock": "1.1.1.1/32",
					"nexthop_type":          "Ipv4Gateway",
					"nexthop_id":            "${alicloud_vpc_ipv4_gateway.defaultIpv4Gateway.id}",
					"depends_on":            []string{"alicloud_vpc_ipv4_gateway.defaultIpv4Gateway"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"nexthop_type":          "Ipv4Gateway",
						"route_table_id":        CHECKSET,
						"nexthop_id":            CHECKSET,
						"destination_cidrblock": "1.1.1.1/32",
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

func TestAccAliCloudRouteEntry_basic3(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_route_entry.default"
	ra := resourceAttrInit(resourceId, AliCloudRouteEntryMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRouteEntry")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAcc-RouteEntry-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRouteEntryBasicDependence3)
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
					"route_table_id":        "${alicloud_vpc.default.route_table_id}",
					"destination_cidrblock": "172.11.1.1/32",
					"nexthop_type":          "NatGateway",
					"nexthop_id":            "${alicloud_nat_gateway.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"route_table_id":        CHECKSET,
						"destination_cidrblock": "172.11.1.1/32",
						"nexthop_type":          "NatGateway",
						"nexthop_id":            CHECKSET,
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

func TestAccAliCloudRouteEntry_basic3_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_route_entry.default"
	ra := resourceAttrInit(resourceId, AliCloudRouteEntryMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRouteEntry")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAcc-RouteEntry-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRouteEntryBasicDependence3)
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
					"route_table_id":        "${alicloud_vpc.default.route_table_id}",
					"destination_cidrblock": "172.11.1.1/32",
					"nexthop_type":          "NatGateway",
					"nexthop_id":            "${alicloud_nat_gateway.default.id}",
					"name":                  name,
					"description":           name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"route_table_id":        CHECKSET,
						"destination_cidrblock": "172.11.1.1/32",
						"nexthop_type":          "NatGateway",
						"nexthop_id":            CHECKSET,
						"name":                  name,
						"description":           name,
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

func TestAccAliCloudRouteEntry_Multi(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_route_entry.default.5"
	ra := resourceAttrInit(resourceId, AliCloudRouteEntryMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeRouteEntry")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(1000, 9999)
	name := fmt.Sprintf("tf-testAcc-RouteEntry-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudRouteEntryBasicDependence4)
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
					"count":                 "6",
					"route_table_id":        "${alicloud_vpc.default.route_table_id}",
					"destination_cidrblock": "172.16.${count.index}.0/24",
					"nexthop_type":          "NetworkInterface",
					"nexthop_id":            "${alicloud_ecs_network_interface.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"route_table_id":        CHECKSET,
						"destination_cidrblock": CHECKSET,
						"nexthop_type":          "NetworkInterface",
						"nexthop_id":            CHECKSET,
					}),
				),
			},
		},
	})
}

var AliCloudRouteEntryMap0 = map[string]string{
	"router_id": CHECKSET,
}

func AliCloudRouteEntryBasicDependence0(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_zones" "default" {
  		available_disk_category     = "cloud_efficiency"
  		available_resource_creation = "VSwitch"
	}

	data "alicloud_images" "default" {
  		name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  		most_recent = true
  		owners      = "system"
	}

	data "alicloud_instance_types" "default" {
  		availability_zone    = data.alicloud_zones.default.zones.0.id
  		image_id             = data.alicloud_images.default.images.0.id
  		instance_type_family = "ecs.g6"
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

func AliCloudRouteEntryBasicDependence1(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_zones" "default" {
  		available_disk_category     = "cloud_efficiency"
  		available_resource_creation = "VSwitch"
	}

	data "alicloud_images" "default" {
  		name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  		most_recent = true
  		owners      = "system"
	}

	data "alicloud_instance_types" "default" {
  		availability_zone                 = data.alicloud_zones.default.zones.0.id
  		image_id                          = data.alicloud_images.default.images.0.id
  		minimum_eni_ipv6_address_quantity = 1
	}

	resource "alicloud_vpc" "default" {
		vpc_name    = var.name
		cidr_block  = "192.168.0.0/16"
		enable_ipv6 = "true"
	}

	resource "alicloud_vswitch" "default" {
		vswitch_name         = var.name
		vpc_id               = alicloud_vpc.default.id
  		cidr_block           = "192.168.192.0/24"
  		zone_id              = data.alicloud_zones.default.zones.0.id
  		ipv6_cidr_block_mask = 64
	}

	resource "alicloud_security_group" "default" {
  		name   = var.name
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
		ipv6_address_count         = 1
	}
`, name)
}

func AliCloudRouteEntryBasicDependence2(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	resource "alicloud_vpc" "defaultVpc" {
	  vpc_name   = "TFRouteEntry1"
	  cidr_block = "192.168.0.0/16"
	}
	
	resource "alicloud_vswitch" "defaultVswitch" {
	  vpc_id       = alicloud_vpc.defaultVpc.id
	  zone_id      = "cn-hangzhou-i"
	  cidr_block   = "192.168.0.0/24"
	  vswitch_name = "TFRouteEntry1"
	}
	
	resource "alicloud_vpc_ipv4_gateway" "defaultIpv4Gateway" {
	  ipv4_gateway_name = "TFRouteEntry"
	  vpc_id            = alicloud_vpc.defaultVpc.id
	  enabled           = true
	}
	
	resource "alicloud_havip" "defaultHavip" {
	  vswitch_id  = alicloud_vswitch.defaultVswitch.id
	  ha_vip_name = "TFRouteEntry1"
	}
`, name)
}

func AliCloudRouteEntryBasicDependence3(name string) string {
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

	resource "alicloud_nat_gateway" "default" {
  		vpc_id           = alicloud_vpc.default.id
  		nat_type         = "Enhanced"
  		vswitch_id       = alicloud_vswitch.default.id
		nat_gateway_name = var.name
	}
`, name)
}

func AliCloudRouteEntryBasicDependence4(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_zones" "default" {
  		available_disk_category     = "cloud_efficiency"
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

	resource "alicloud_ecs_network_interface" "default" {
		network_interface_name = var.name
		vswitch_id             = alicloud_vswitch.default.id
  		security_group_ids     = alicloud_security_group.default.*.id
	}
`, name)
}
