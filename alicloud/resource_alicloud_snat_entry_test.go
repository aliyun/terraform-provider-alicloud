package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

// Test NATGateway SnatEntry. >>> Resource test cases, automatically generated.
// Case 全生命周期_SnatEntry source_cidr 8016
func TestAccAliCloudNATGatewaySnatEntry_basic8016(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_snat_entry.default"
	ra := resourceAttrInit(resourceId, AliCloudNATGatewaySnatEntryMap8016)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NATGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNATGatewaySnatEntry")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snatgatewaysnatentry%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudNATGatewaySnatEntryBasicDependence8016)
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
					"snat_ip":       "${alicloud_vpc_nat_ip.default.nat_ip}" + "," + "${alicloud_vpc_nat_ip.update.nat_ip}",
					"snat_table_id": "${alicloud_nat_gateway.default.snat_table_ids}",
					"source_cidr":   "${cidrsubnet(alicloud_vpc.default.cidr_block, 8, 8)}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"snat_ip":       CHECKSET,
						"snat_table_id": CHECKSET,
						"source_cidr":   CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"eip_affinity": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"eip_affinity": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"eip_affinity": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"eip_affinity": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"snat_entry_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"snat_entry_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"snat_ip": "${alicloud_vpc_nat_ip.default.nat_ip}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"snat_ip": CHECKSET,
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

func TestAccAliCloudNATGatewaySnatEntry_basic8016_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_snat_entry.default"
	ra := resourceAttrInit(resourceId, AliCloudNATGatewaySnatEntryMap8016)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NATGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNATGatewaySnatEntry")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snatgatewaysnatentry%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudNATGatewaySnatEntryBasicDependence8016)
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
					"snat_ip":         "${alicloud_vpc_nat_ip.default.nat_ip}" + "," + "${alicloud_vpc_nat_ip.update.nat_ip}",
					"snat_table_id":   "${alicloud_nat_gateway.default.snat_table_ids}",
					"source_cidr":     "${cidrsubnet(alicloud_vpc.default.cidr_block, 8, 8)}",
					"snat_entry_name": name,
					"eip_affinity":    "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"snat_ip":         CHECKSET,
						"snat_table_id":   CHECKSET,
						"source_cidr":     CHECKSET,
						"snat_entry_name": name,
						"eip_affinity":    "1",
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

// Case 全生命周期_SnatEntry source_vswitch_id 8018
func TestAccAliCloudNATGatewaySnatEntry_basic8018(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_snat_entry.default"
	ra := resourceAttrInit(resourceId, AliCloudNATGatewaySnatEntryMap8016)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NATGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNATGatewaySnatEntry")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snatgatewaysnatentry%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudNATGatewaySnatEntryBasicDependence8016)
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
					"snat_ip":           "${alicloud_vpc_nat_ip.default.nat_ip}" + "," + "${alicloud_vpc_nat_ip.update.nat_ip}",
					"snat_table_id":     "${alicloud_nat_gateway.default.snat_table_ids}",
					"source_vswitch_id": "${alicloud_vswitch.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"snat_ip":           CHECKSET,
						"snat_table_id":     CHECKSET,
						"source_vswitch_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"eip_affinity": "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"eip_affinity": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"eip_affinity": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"eip_affinity": "0",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"snat_entry_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"snat_entry_name": name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"snat_ip": "${alicloud_vpc_nat_ip.default.nat_ip}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"snat_ip": CHECKSET,
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

func TestAccAliCloudNATGatewaySnatEntry_basic8018_twin(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_snat_entry.default"
	ra := resourceAttrInit(resourceId, AliCloudNATGatewaySnatEntryMap8016)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NATGatewayServiceV2{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeNATGatewaySnatEntry")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snatgatewaysnatentry%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AliCloudNATGatewaySnatEntryBasicDependence8016)
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
					"snat_ip":           "${alicloud_vpc_nat_ip.default.nat_ip}" + "," + "${alicloud_vpc_nat_ip.update.nat_ip}",
					"snat_table_id":     "${alicloud_nat_gateway.default.snat_table_ids}",
					"source_vswitch_id": "${alicloud_vswitch.default.id}",
					"snat_entry_name":   name,
					"eip_affinity":      "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"snat_ip":           CHECKSET,
						"snat_table_id":     CHECKSET,
						"source_vswitch_id": CHECKSET,
						"snat_entry_name":   name,
						"eip_affinity":      "1",
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

var AliCloudNATGatewaySnatEntryMap8016 = map[string]string{
	"snat_entry_id": CHECKSET,
	"status":        CHECKSET,
}

func AliCloudNATGatewaySnatEntryBasicDependence8016(name string) string {
	return fmt.Sprintf(`
	variable "name" {
  		default = "%s"
	}

	data "alicloud_zones" "default" {
	}

	resource "alicloud_vpc" "default" {
  		vpc_name   = var.name
  		cidr_block = "172.16.0.0/12"
	}

	resource "alicloud_vswitch" "default" {
  		vpc_id       = alicloud_vpc.default.id
  		cidr_block   = "172.16.0.0/24"
  		zone_id      = data.alicloud_zones.default.zones.0.id
  		vswitch_name = var.name
	}

	resource "alicloud_nat_gateway" "default" {
  		vpc_id               = alicloud_vpc.default.id
  		network_type         = "intranet"
  		nat_gateway_name     = var.name
  		vswitch_id           = alicloud_vswitch.default.id
  		nat_type             = "Enhanced"
  		internet_charge_type = "PayByLcu"
	}

	resource "alicloud_vpc_nat_ip" "default" {
  		nat_ip         = "172.16.0.66"
  		nat_ip_cidr    = alicloud_vswitch.default.cidr_block
  		nat_gateway_id = alicloud_nat_gateway.default.id
	}

	resource "alicloud_vpc_nat_ip" "update" {
  		nat_ip         = "172.16.0.88"
  		nat_ip_cidr    = alicloud_vswitch.default.cidr_block
  		nat_gateway_id = alicloud_nat_gateway.default.id
	}
`, name)
}

// Test NATGateway SnatEntry. <<< Resource test cases, automatically generated.
