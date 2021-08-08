package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudSnatEntry_basic(t *testing.T) {
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

func TestAccAlicloudSnatEntry_multi(t *testing.T) {
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

var AlicloudSnatEntryMap0 = map[string]string{}

func AlicloudSnatEntryBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = "${alicloud_vpc.default.id}"
  cidr_block   = "172.16.0.0/21"
  zone_id      = "${data.alicloud_zones.default.zones.0.id}"
  vswitch_name = "${var.name}"
}

resource "alicloud_nat_gateway" "default" {
  depends_on    = [alicloud_vpc.default]
  vpc_id        = alicloud_vswitch.default.vpc_id
  specification = "Small"
  nat_gateway_name          = "${var.name}"
  vswitch_id    = alicloud_vswitch.default.id
  nat_type      = "Enhanced"
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

resource "alicloud_vpc" "default" {
  vpc_name   = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  count = 3
  vpc_id       = "${alicloud_vpc.default.id}"
  cidr_block   = "172.16.${count.index}.0/24"
  zone_id      = "${data.alicloud_zones.default.zones.0.id}"
  vswitch_name = "${var.name}"
}

resource "alicloud_nat_gateway" "default" {
  depends_on    = [alicloud_vpc.default]
  vpc_id        = alicloud_vswitch.default[0].vpc_id
  specification = "Small"
  nat_gateway_name          = "${var.name}"
  vswitch_id    = alicloud_vswitch.default[0].id
  nat_type      = "Enhanced"
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
    source_vswitch_id = "${alicloud_vswitch.default[count.index].id}"
	depends_on = [alicloud_eip_association.default]
}

`, name)
}
