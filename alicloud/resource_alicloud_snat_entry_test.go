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
					"source_vswitch_id": "${data.alicloud_vswitches.default.ids[0]}",
					"depends_on":        []string{"alicloud_eip_association.default"},
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

data "alicloud_enhanced_nat_available_zones" "default" {}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_enhanced_nat_available_zones.default.zones.0.zone_id
}

resource "alicloud_nat_gateway" "default" {
  vpc_id           = data.alicloud_vpcs.default.ids.0
  nat_gateway_name = var.name
  payment_type     = "PayAsYouGo"
  vswitch_id       = data.alicloud_vswitches.default.ids[0]
  nat_type         = "Enhanced"
}

resource "alicloud_eip_address" "default" {
  address_name = var.name
}

resource "alicloud_eip_association" "default" {
  allocation_id = alicloud_eip_address.default.id
  instance_id   = alicloud_nat_gateway.default.id
}
`, name)
}

func AlicloudSnatEntryMultiDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_enhanced_nat_available_zones" "default" {}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_enhanced_nat_available_zones.default.zones.0.zone_id
}

resource "alicloud_nat_gateway" "default" {
  vpc_id           = data.alicloud_vpcs.default.ids.0
  nat_gateway_name = var.name
  payment_type     = "PayAsYouGo"
  vswitch_id       = data.alicloud_vswitches.default.ids[0]
  nat_type         = "Enhanced"
}

resource "alicloud_eip_address" "default" {
  address_name = var.name
}

resource "alicloud_eip_association" "default" {
  allocation_id = alicloud_eip_address.default.id
  instance_id   = alicloud_nat_gateway.default.id
}

resource "alicloud_snat_entry" "default" {
  count             = 3
  snat_ip           = alicloud_eip_address.default.ip_address
  snat_table_id     = alicloud_nat_gateway.default.snat_table_ids
  source_vswitch_id = data.alicloud_vswitches.default.ids[count.index]
  depends_on        = [alicloud_eip_association.default]
}
`, name)
}
