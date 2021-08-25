package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudVpcForwardEntry_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alicloud_forward_entry.default"
	ra := resourceAttrInit(resourceId, AlicloudForwardEntryMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeForwardEntry")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sforwardentry%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlicloudForwardEntryBasicDependence0)
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
					"forward_table_id": "${alicloud_nat_gateway.default.forward_table_ids}",
					"external_ip":      "${alicloud_eip_address.default.0.ip_address}",
					"external_port":    `80`,
					"internal_ip":      "172.16.0.3",
					"internal_port":    `8080`,
					"ip_protocol":      "tcp",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"forward_table_id": CHECKSET,
						"external_ip":      CHECKSET,
						"external_port":    "80",
						"internal_ip":      "172.16.0.3",
						"internal_port":    "8080",
						"ip_protocol":      "tcp",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"port_break"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"external_ip": "${alicloud_eip_address.default.1.ip_address}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"external_ip": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"external_port": `90`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"external_port": "90",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"forward_entry_name": name + "1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"forward_entry_name": name + "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"internal_ip": "172.16.0.4",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"internal_ip": "172.16.0.4",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"internal_port": `9090`,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"internal_port": "9090",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ip_protocol": "udp",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ip_protocol": "udp",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"external_ip":        "${alicloud_eip_address.default.0.ip_address}",
					"external_port":      `80`,
					"forward_entry_name": "${var.name}",
					"internal_ip":        "172.16.0.3",
					"internal_port":      `8080`,
					"ip_protocol":        "tcp",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"external_ip":        CHECKSET,
						"external_port":      "80",
						"forward_entry_name": name,
						"internal_ip":        "172.16.0.3",
						"internal_port":      "8080",
						"ip_protocol":        "tcp",
					}),
				),
			},
		},
	})
}

var AlicloudForwardEntryMap0 = map[string]string{
	"port_break": "false",
}

func AlicloudForwardEntryBasicDependence0(name string) string {
	return fmt.Sprintf(`
variable "name" {
			default = "%s"
		}

variable "number" {
	default = "2"
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
	zone_id = "${data.alicloud_zones.default.zones.0.id}"
	vswitch_name = "${var.name}"
}

resource "alicloud_nat_gateway" "default" {
	vpc_id = "${alicloud_vswitch.default.vpc_id}"
	specification = "Small"
	nat_gateway_name = "${var.name}"
	nat_type = "Enhanced"
	vswitch_id = alicloud_vswitch.default.id
}

resource "alicloud_eip_address" "default" {
	count = "${var.number}"
	address_name = "${var.name}"
}

resource "alicloud_eip_association" "default" {
	count = "${var.number}"
	allocation_id = "${element(alicloud_eip_address.default.*.id,count.index)}"
	instance_id = "${alicloud_nat_gateway.default.id}"
}
`, name)
}
