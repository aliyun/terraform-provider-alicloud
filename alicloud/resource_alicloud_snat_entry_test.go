package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudVpcSnatEntry_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_snat_entry.default"
	ra := resourceAttrInit(resourceId, AlicloudSnatEntryMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSnatEntry")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%ssnatentry%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSnatEntryBasicDependence0)
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
					"snat_ip":           "${alicloud_eip_address.default.ip_address}",
					"snat_table_id":     "${alicloud_nat_gateway.default.snat_table_ids}",
					"source_vswitch_id": "${local.vswitch_id}",
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
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"snat_entry_name": name + "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"snat_entry_name": name + "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"snat_entry_name": "${var.name}",
					"snat_ip":         "${alicloud_eip_address.default.ip_address}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"snat_entry_name": name,
						"snat_ip":         CHECKSET,
					}),
				),
			},
		},
	})
}

func TestAccAlicloudVpcSnatEntry_multi(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_snat_entry.default.1"
	ra := resourceAttrInit(resourceId, AlicloudSnatEntryMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSnatEntry")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%ssnatentry%d", defaultRegionToTest, rand)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: AlicloudSnatEntryMultiDependence0(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"snat_ip":           CHECKSET,
						"snat_table_id":     CHECKSET,
						"source_vswitch_id": CHECKSET,
					}),
				),
			},
		},
	})
}

func TestAccAlicloudVpcSnatEntry_basic1(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_snat_entry.default"
	ra := resourceAttrInit(resourceId, AlicloudSnatEntryMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSnatEntry")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%ssnatentry%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudSnatEntryBasicDependence0)
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
					"snat_ip":         "${alicloud_eip_address.default.ip_address}",
					"snat_table_id":   "${alicloud_nat_gateway.default.snat_table_ids}",
					"source_cidr":     "${cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)}",
					"snat_entry_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"snat_ip":         CHECKSET,
						"snat_table_id":   CHECKSET,
						"source_cidr":     CHECKSET,
						"snat_entry_name": name,
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

var AlicloudSnatEntryMap0 = map[string]string{}

func AlicloudSnatEntryBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
	vpc_id = data.alicloud_vpcs.default.ids.0
	zone_id      = data.alicloud_zones.default.zones.0.id
}
resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_zones.default.zones.0.id
  vswitch_name      = var.name
}
locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_nat_gateway" "default" {
  vpc_id        = data.alicloud_vpcs.default.ids.0
  specification = "Middle"
  network_type  = "internet"
  nat_gateway_name = "${var.name}"
  vswitch_id    = local.vswitch_id
  nat_type      = "Enhanced"
  period        = "5"
  instance_charge_type = "PrePaid"
  internet_charge_type = "PayBySpec"
}

resource "alicloud_eip_address" "default" {
  address_name = "${var.name}"
}

resource "alicloud_eip_association" "default" {
  depends_on    = [alicloud_eip_address.default, alicloud_nat_gateway.default]
  allocation_id = "${alicloud_eip_address.default.id}"
  instance_id   = "${alicloud_nat_gateway.default.id}"
}

`, name)
}

func AlicloudSnatEntryMultiDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_vpcs" "default" {
	name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
	vpc_id = data.alicloud_vpcs.default.ids.0
	zone_id      = data.alicloud_zones.default.zones.0.id
}
resource "alicloud_vswitch" "vswitch" {
  count             = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id            = data.alicloud_vpcs.default.ids.0
  cidr_block        = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id           = data.alicloud_zones.default.zones.0.id
  vswitch_name      = var.name
}
locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}

resource "alicloud_nat_gateway" "default" {
  vpc_id        = data.alicloud_vpcs.default.ids.0
  specification = "Middle"
  network_type  = "internet"
  nat_gateway_name = "${var.name}"
  vswitch_id    = local.vswitch_id
  nat_type      = "Enhanced"
  period        = "5"
  instance_charge_type = "PrePaid"
  internet_charge_type = "PayBySpec"
}

resource "alicloud_eip_address" "default" {
  address_name = "${var.name}"
}

resource "alicloud_eip_association" "default" {
  depends_on    = [alicloud_eip_address.default, alicloud_nat_gateway.default]
  allocation_id = "${alicloud_eip_address.default.id}"
  instance_id   = "${alicloud_nat_gateway.default.id}"
}

resource "alicloud_snat_entry" "default" {
	count = 3
	snat_ip = "${alicloud_eip_address.default.ip_address}"
    snat_table_id = "${alicloud_nat_gateway.default.snat_table_ids}"
    source_vswitch_id = local.vswitch_id
	depends_on = [alicloud_eip_association.default]
}

`, name)
}
