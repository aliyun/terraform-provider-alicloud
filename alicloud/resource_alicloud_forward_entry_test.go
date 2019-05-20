package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudForwardBasic(t *testing.T) {
	var v vpc.ForwardTableEntry
	resourceId := "alicloud_forward_entry.default"

	rand := acctest.RandInt()
	testAccForwardEntryCheckMap["name"] = fmt.Sprintf("tf-testAccForwardEntryConfig%d", rand)
	ra := resourceAttrInit(resourceId, testAccForwardEntryCheckMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckForwardEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccForwardEntryConfigBasic(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccForwardEntryConfig_external_ip(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccForwardEntryConfig_external_port(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"external_port": "81",
					}),
				),
			},
			{
				Config: testAccForwardEntryConfig_ip_protocol(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_protocol": "udp",
					}),
				),
			},
			{
				Config: testAccForwardEntryConfig_internal_ip(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"internal_ip": "172.16.0.4",
					}),
				),
			},
			{
				Config: testAccForwardEntryConfig_internal_port(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"internal_port": "8081",
					}),
				),
			},
			{
				Config: testAccForwardEntryConfig_name(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": fmt.Sprintf("tf-testAccForwardEntryConfig%d_change", rand),
					}),
				),
			},
			{
				Config: testAccForwardEntryConfigBasic(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(testAccForwardEntryCheckMap),
				),
			},
		},
	})
}

func TestAccAlicloudForwardMulti(t *testing.T) {
	var v vpc.ForwardTableEntry
	resourceId := "alicloud_forward_entry.default.4"
	rand := acctest.RandInt()
	testAccForwardEntryCheckMap["name"] = fmt.Sprintf("tf-testAccForwardEntryConfig%d", rand)
	ra := resourceAttrInit(resourceId, testAccForwardEntryCheckMap)
	serviceFunc := func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckForwardEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccForwardEntryConfig_multi(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"external_port": "84",
						"internal_port": "8084",
					}),
				),
			},
		},
	})
}

func testAccCheckForwardEntryDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	vpcService := VpcService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_forward_entry" {
			continue
		}
		if _, err := vpcService.DescribeForwardEntry(rs.Primary.ID); err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}

		return WrapError(fmt.Errorf("Forward entry %s still exist", rs.Primary.ID))
	}
	return nil
}

func testAccForwardEntryConfigBasic(rand int) string {
	config := fmt.Sprintf(`
%s

resource "alicloud_forward_entry" "default"{
	name = "${var.name}"
	forward_table_id = "${alicloud_nat_gateway.default.forward_table_ids}"
	external_ip = "${alicloud_eip.default.0.ip_address}"
	external_port = "80"
	ip_protocol = "tcp"
	internal_ip = "172.16.0.3"
	internal_port = "8080"
}
`, testAccForwardEntryConfigCommon(rand))
	return config
}

func testAccForwardEntryConfig_external_ip(rand int) string {
	config := fmt.Sprintf(`
%s

resource "alicloud_forward_entry" "default"{
	name = "${var.name}"
	forward_table_id = "${alicloud_nat_gateway.default.forward_table_ids}"
	external_ip = "${alicloud_eip.default.1.ip_address}"
	external_port = "80"
	ip_protocol = "tcp"
	internal_ip = "172.16.0.3"
	internal_port = "8080"
}
`, testAccForwardEntryConfigCommon(rand))
	return config
}

func testAccForwardEntryConfig_external_port(rand int) string {
	return fmt.Sprintf(`
%s

resource "alicloud_forward_entry" "default"{
	name = "${var.name}"
	forward_table_id = "${alicloud_nat_gateway.default.forward_table_ids}"
	external_ip = "${alicloud_eip.default.1.ip_address}"
	external_port = "81"
	ip_protocol = "tcp"
	internal_ip = "172.16.0.3"
	internal_port = "8080"
}
`, testAccForwardEntryConfigCommon(rand))
}

func testAccForwardEntryConfig_ip_protocol(rand int) string {
	return fmt.Sprintf(`
%s

resource "alicloud_forward_entry" "default"{
	name = "${var.name}"
	forward_table_id = "${alicloud_nat_gateway.default.forward_table_ids}"
	external_ip = "${alicloud_eip.default.1.ip_address}"
	external_port = "81"
	ip_protocol = "udp"
	internal_ip = "172.16.0.3"
	internal_port = "8080"
}
`, testAccForwardEntryConfigCommon(rand))
}

func testAccForwardEntryConfig_internal_ip(rand int) string {
	return fmt.Sprintf(`
%s

resource "alicloud_forward_entry" "default"{
	name = "${var.name}"
	forward_table_id = "${alicloud_nat_gateway.default.forward_table_ids}"
	external_ip = "${alicloud_eip.default.1.ip_address}"
	external_port = "81"
	ip_protocol = "udp"
	internal_ip = "172.16.0.4"
	internal_port = "8080"
}
`, testAccForwardEntryConfigCommon(rand))
}

func testAccForwardEntryConfig_internal_port(rand int) string {
	return fmt.Sprintf(`
%s

resource "alicloud_forward_entry" "default"{
	name = "${var.name}"
	forward_table_id = "${alicloud_nat_gateway.default.forward_table_ids}"
	external_ip = "${alicloud_eip.default.1.ip_address}"
	external_port = "81"
	ip_protocol = "udp"
	internal_ip = "172.16.0.4"
	internal_port = "8081"
}
`, testAccForwardEntryConfigCommon(rand))
}

func testAccForwardEntryConfig_name(rand int) string {
	return fmt.Sprintf(`
%s

resource "alicloud_forward_entry" "default"{
	name = "${var.name}_change"
	forward_table_id = "${alicloud_nat_gateway.default.forward_table_ids}"
	external_ip = "${alicloud_eip.default.1.ip_address}"
	external_port = "81"
	ip_protocol = "udp"
	internal_ip = "172.16.0.4"
	internal_port = "8081"
}
`, testAccForwardEntryConfigCommon(rand))
}

func testAccForwardEntryConfig_multi(rand int) string {
	config := fmt.Sprintf(`
%s

resource "alicloud_forward_entry" "default"{
	count = 5
	forward_table_id = "${alicloud_nat_gateway.default.forward_table_ids}"
	external_ip = "${alicloud_eip.default.0.ip_address}"
	external_port = "${80 + count.index}"
	ip_protocol = "tcp"
	internal_ip = "172.16.0.3"
	internal_port = "${8080 + count.index}"
}
`, testAccForwardEntryConfigCommon(rand))
	return config
}

func testAccForwardEntryConfigCommon(rand int) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "tf-testAccForwardEntryConfig%d"
}

variable "count" {
	default = "2"
}

data "alicloud_zones" "default" {
	"available_resource_creation"= "VSwitch"
}

resource "alicloud_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
	vpc_id = "${alicloud_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_nat_gateway" "default" {
	vpc_id = "${alicloud_vpc.default.id}"
	specification = "Small"
	name = "${var.name}"
}

resource "alicloud_eip" "default" {
	count = "${var.count}"
	name = "${var.name}"
}

resource "alicloud_eip_association" "default" {
	count = "${var.count}"
	allocation_id = "${element(alicloud_eip.default.*.id,count.index)}"
	instance_id = "${alicloud_nat_gateway.default.id}"
}
`, rand)
}

var testAccForwardEntryCheckMap = map[string]string{
	"forward_table_id": CHECKSET,
	"external_ip":      CHECKSET,
	"external_port":    "80",
	"ip_protocol":      "tcp",
	"internal_ip":      "172.16.0.3",
	"internal_port":    "8080",
	"forward_entry_id": CHECKSET,
}
