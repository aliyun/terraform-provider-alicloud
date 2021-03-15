package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func testAccCheckSnatEntryDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_snat_entry" {
			continue
		}

		// Try to find the Snat entry
		_, err := vpcService.DescribeSnatEntry(rs.Primary.ID)

		//this special deal cause the DescribeSnatEntry can't find the records would be throw "cant find the snatTable error"
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}

		return WrapError(Error("Snat entry still exist"))
	}

	return nil
}

func TestAccAlicloudSnatEntryBasic(t *testing.T) {
	var v vpc.SnatTableEntry

	resourceId := "alicloud_snat_entry.default"
	ra := resourceAttrInit(resourceId, testAccCheckSnatEntryBasicMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_snat_entry.default",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSnatEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSnatEntryConfigBasic(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"snat_entry_name": fmt.Sprintf("tf-testAccSnatEntryConfig%d", rand),
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccSnatEntryConfig_snatIp(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccSnatEntryConfig_snatname(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"snat_entry_name": fmt.Sprintf("tf-testAccSnatEntryConfig%d-update", rand),
					}),
				),
			},
		},
	})

}

func TestAccAlicloudSnatEntryMulti(t *testing.T) {
	var v vpc.SnatTableEntry

	resourceId := "alicloud_snat_entry.default.9"
	ra := resourceAttrInit(resourceId, testAccCheckSnatEntryBasicMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)

	rand := acctest.RandInt()
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_snat_entry.default.9",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSnatEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSnatEntryConfigMulti(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})

}

func testAccSnatEntryConfigBasic(rand int) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "tf-testAccSnatEntryConfig%d"
}
data "alicloud_zones" "default" {
	available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
	vpc_name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
	vpc_id = "${alicloud_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_nat_gateway" "default" {
	depends_on = [alicloud_vpc.default]
	vpc_id = "${alicloud_vswitch.default.vpc_id}"
	specification = "Small"
	name = "${var.name}"
}

resource "alicloud_eip" "default" {
	name = "${var.name}"
}

resource "alicloud_eip_association" "default" {
	depends_on = [alicloud_eip.default, alicloud_nat_gateway.default]
	allocation_id = "${alicloud_eip.default.id}"
	instance_id = "${alicloud_nat_gateway.default.id}"
}

resource "alicloud_snat_entry" "default"{
	depends_on = [alicloud_eip_association.default, alicloud_nat_gateway.default]
	snat_table_id = "${alicloud_nat_gateway.default.snat_table_ids}"
	source_vswitch_id = "${alicloud_vswitch.default.id}"
	snat_ip = "${alicloud_eip.default.ip_address}"
    snat_entry_name = "${var.name}"
}

resource "alicloud_snat_entry" "ecs"{
	depends_on = [alicloud_eip_association.default, alicloud_nat_gateway.default]
	snat_table_id = "${alicloud_nat_gateway.default.snat_table_ids}"
	source_cidr = "172.16.0.10/32"
	snat_ip = "${alicloud_eip.default.ip_address}"
}
`, rand)
}

func testAccSnatEntryConfig_snatIp(rand int) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "tf-testAccSnatEntryConfig%d"
}
data "alicloud_zones" "default" {
	available_resource_creation= "VSwitch"
}

resource "alicloud_vpc" "default" {
	vpc_name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
	vpc_id = "${alicloud_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_nat_gateway" "default" {
	depends_on = [alicloud_vpc.default]
	vpc_id = "${alicloud_vswitch.default.vpc_id}"
	specification = "Small"
	name = "${var.name}"
}

resource "alicloud_eip" "default" {
	name = "${var.name}"
}

resource "alicloud_eip_association" "default" {
	depends_on = [alicloud_eip.default, alicloud_nat_gateway.default]
	allocation_id = "${alicloud_eip.default.id}"
	instance_id = "${alicloud_nat_gateway.default.id}"
}

resource "alicloud_snat_entry" "default"{
	depends_on = [alicloud_eip_association.default, alicloud_nat_gateway.default]
	snat_table_id = "${alicloud_nat_gateway.default.snat_table_ids}"
	source_vswitch_id = "${alicloud_vswitch.default.id}"
	snat_ip = "${alicloud_eip.default.ip_address}"
    snat_entry_name = "${var.name}"
}

resource "alicloud_snat_entry" "ecs"{
	depends_on = [alicloud_eip_association.default, alicloud_nat_gateway.default]
	snat_table_id = "${alicloud_nat_gateway.default.snat_table_ids}"
	source_cidr = "172.16.0.10/32"
	snat_ip = "${alicloud_eip.default.ip_address}"
}

`, rand)
}

func testAccSnatEntryConfig_snatname(rand int) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "tf-testAccSnatEntryConfig%d"
}
data "alicloud_zones" "default" {
	available_resource_creation= "VSwitch"
}

resource "alicloud_vpc" "default" {
	vpc_name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
	vpc_id = "${alicloud_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_nat_gateway" "default" {
	depends_on = [alicloud_vpc.default]
	vpc_id = "${alicloud_vswitch.default.vpc_id}"
	specification = "Small"
	name = "${var.name}"
}

resource "alicloud_eip" "default" {
	name = "${var.name}"
}

resource "alicloud_eip_association" "default" {
	depends_on = [alicloud_eip.default, alicloud_nat_gateway.default]
	allocation_id = "${alicloud_eip.default.id}"
	instance_id = "${alicloud_nat_gateway.default.id}"
}

resource "alicloud_snat_entry" "default"{
	depends_on = [alicloud_eip_association.default, alicloud_nat_gateway.default]
	snat_table_id = "${alicloud_nat_gateway.default.snat_table_ids}"
	source_vswitch_id = "${alicloud_vswitch.default.id}"
	snat_ip = "${alicloud_eip.default.ip_address}"
    snat_entry_name = "${var.name}-update"
}

resource "alicloud_snat_entry" "ecs"{
	depends_on = [alicloud_eip_association.default, alicloud_nat_gateway.default]
	snat_table_id = "${alicloud_nat_gateway.default.snat_table_ids}"
	source_cidr = "172.16.0.10/32"
	snat_ip = "${alicloud_eip.default.ip_address}"
}

`, rand)
}

func testAccSnatEntryConfigMulti(rand int) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "tf-testAccSnatEntryMulti%d"
}

data "alicloud_zones" "default" {
	available_resource_creation= "VSwitch"
}

resource "alicloud_vpc" "default" {
	vpc_name = "${var.name}"
	cidr_block = "10.1.0.0/16"
}

resource "alicloud_vswitch" "default" {
    count = 10
    vpc_id            = "${alicloud_vpc.default.id}"
    cidr_block        = "10.1.${count.index + 1}.0/24"
    availability_zone = "${data.alicloud_zones.default.zones.0.id}"
    name = "${var.name}"
}

resource "alicloud_nat_gateway" "default" {
	depends_on = [alicloud_vpc.default]
	vpc_id = "${alicloud_vpc.default.id}"
	specification = "Small"
	name = "${var.name}"
}

resource "alicloud_eip" "default" {
	name = "${var.name}"
}

resource "alicloud_eip_association" "default" {
	depends_on = [alicloud_eip.default, alicloud_nat_gateway.default]
	allocation_id = "${alicloud_eip.default.id}"
	instance_id = "${alicloud_nat_gateway.default.id}"
}

resource "alicloud_snat_entry" "default"{
	depends_on = [alicloud_eip_association.default, alicloud_nat_gateway.default]
	count = "10"
	snat_table_id = "${alicloud_nat_gateway.default.snat_table_ids}"
	source_vswitch_id = "${element(alicloud_vswitch.default.*.id, count.index)}"
	snat_ip = "${alicloud_eip.default.ip_address}"
    snat_entry_name = "${var.name}"
}

resource "alicloud_snat_entry" "ecs"{
	depends_on = [alicloud_eip_association.default, alicloud_nat_gateway.default]
	snat_table_id = "${alicloud_nat_gateway.default.snat_table_ids}"
	source_cidr = "10.1.2.10/32"
	snat_ip = "${alicloud_eip.default.ip_address}"
}
`, rand)
}

var testAccCheckSnatEntryBasicMap = map[string]string{
	"snat_table_id":     CHECKSET,
	"source_vswitch_id": CHECKSET,
	"snat_ip":           CHECKSET,
	"snat_entry_id":     CHECKSET,
}
