package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudSLBServerGroup_vpc(t *testing.T) {
	var v *slb.DescribeVServerGroupAttributeResponse
	resourceId := "alicloud_slb_server_group.default"
	ra := resourceAttrInit(resourceId, serverGroupMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccSlbServerGroupVpc")
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceSlbServerGroupDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_id": "${alicloud_slb_load_balancer.default.id}",
					"name":             name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_protection_validation"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"servers": []map[string]interface{}{
						{
							"server_ids": []string{"${alicloud_instance.default.0.id}", "${alicloud_instance.default.1.id}"},
							"port":       "100",
							"weight":     "10",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"servers.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": "tf-testAccSlbServerGroupVpcUpdate",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": "tf-testAccSlbServerGroupVpcUpdate",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"servers": []map[string]interface{}{
						{
							"server_ids": []string{"${alicloud_instance.default.0.id}", "${alicloud_instance.default.1.id}"},
							"port":       "100",
							"weight":     "10",
						},
						{
							"server_ids": []string{"${alicloud_network_interface.default.0.id}"},
							"port":       "70",
							"weight":     "10",
							"type":       "eni",
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
				Config: testAccConfig(map[string]interface{}{
					"servers": serversMap,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"servers.#": "14",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": "${var.name}",
					"servers": []map[string]interface{}{
						{
							"server_ids": []string{"${alicloud_instance.default.0.id}", "${alicloud_instance.default.1.id}"},
							"port":       "100",
							"weight":     "10",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":      "tf-testAccSlbServerGroupVpc",
						"servers.#": "1",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudSLBServerGroup_multi_vpc(t *testing.T) {
	var v *slb.DescribeVServerGroupAttributeResponse
	resourceId := "alicloud_slb_server_group.default.9"
	ra := resourceAttrInit(resourceId, serverGroupMultiVpcMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccSlbServerGroupVpc")
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceSlbServerGroupMultiVpcDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_id": "${alicloud_slb_load_balancer.default.id}",
					"count":            "10",
					"servers": []map[string]interface{}{
						{
							"server_ids": []string{"${alicloud_instance.default.0.id}", "${alicloud_instance.default.1.id}"},
							"port":       "100",
							"weight":     "10",
						},
						{
							"server_ids": "${alicloud_instance.default.*.id}",
							"port":       "80",
							"weight":     "100",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func TestAccAlicloudSLBServerGroup_classic(t *testing.T) {
	var v *slb.DescribeVServerGroupAttributeResponse
	resourceId := "alicloud_slb_server_group.default"
	ra := resourceAttrInit(resourceId, serverGroupMultiClassicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccSlbServerGroupClassic")
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceServerGroupClassicDependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.SlbClassicNoSupportedRegions)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_id": "${alicloud_slb_load_balancer.default.id}",
					"name":             name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": "tf-testAccSlbServerGroupClassic",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_protection_validation"},
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"servers": []map[string]interface{}{
						{
							"server_ids": []string{"${alicloud_instance.default.0.id}", "${alicloud_instance.default.1.id}"},
							"port":       "100",
							"weight":     "10",
						},
						{
							"server_ids": "${alicloud_instance.default.*.id}",
							"port":       "80",
							"weight":     "100",
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
				Config: testAccConfig(map[string]interface{}{
					"name": "tf-testAccSlbServerGroupClassicUpdate",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": "tf-testAccSlbServerGroupClassicUpdate",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"servers": []map[string]interface{}{
						{
							"server_ids": []string{"${alicloud_instance.default.0.id}", "${alicloud_instance.default.1.id}"},
							"port":       "100",
							"weight":     "10",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"servers.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": "${var.name}",
					"servers": []map[string]interface{}{
						{
							"server_ids": []string{"${alicloud_instance.default.0.id}", "${alicloud_instance.default.1.id}"},
							"port":       "100",
							"weight":     "10",
						},
						{
							"server_ids": "${alicloud_instance.default.*.id}",
							"port":       "80",
							"weight":     "100",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":      "tf-testAccSlbServerGroupClassic",
						"servers.#": "2",
					}),
				),
			},
		},
	})
}

func resourceSlbServerGroupDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
data "alicloud_instance_types" "default" {
  cpu_core_count    = 1
  memory_size       = 2
  availability_zone = data.alicloud_slb_zones.default.zones.0.id
}
data "alicloud_instance_types" "new" {
	eni_amount = 2
    availability_zone = data.alicloud_slb_zones.default.zones.0.id
}
data "alicloud_images" "default" {
  name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
  owners = "system"
}
data "alicloud_vpcs" "default"{
	name_regex = "default-NODELETING"
}
data "alicloud_slb_zones" "default" {
	available_slb_address_type = "vpc"
}

data "alicloud_vswitches" "default" {
	vpc_id  = data.alicloud_vpcs.default.ids.0
	zone_id = data.alicloud_slb_zones.default.zones.0.id
}
resource "alicloud_security_group" "default" {
  name = "${var.name}"
  vpc_id = data.alicloud_vpcs.default.ids.0
}
resource "alicloud_network_interface" "default" {
    count = 1
    name = "${var.name}"
    vswitch_id = data.alicloud_vswitches.default.vswitches.0.id
    security_groups = [ "${alicloud_security_group.default.id}" ]
}
resource "alicloud_instance" "default" {
  image_id = "${data.alicloud_images.default.images.0.id}"
  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  instance_name = "${var.name}"
  count = "21"
  security_groups = "${alicloud_security_group.default.*.id}"
  internet_charge_type = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone = "${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"
  instance_charge_type = "PostPaid"
  system_disk_category = "cloud_efficiency"
  vswitch_id = data.alicloud_vswitches.default.ids[0]
}
resource "alicloud_instance" "new" {
  image_id = "${data.alicloud_images.default.images.0.id}"
  instance_type = "${data.alicloud_instance_types.new.instance_types.0.id}"
  instance_name = "${var.name}"
  count = "1"
  security_groups = "${alicloud_security_group.default.*.id}"
  internet_charge_type = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone = "${data.alicloud_instance_types.new.instance_types.0.availability_zones.0}"
  instance_charge_type = "PostPaid"
  system_disk_category = "cloud_efficiency"
  vswitch_id = data.alicloud_vswitches.default.ids[0]
}
resource "alicloud_network_interface_attachment" "default" {
	count = 1
    instance_id = "${alicloud_instance.new.0.id}"
    network_interface_id = "${element(alicloud_network_interface.default.*.id, count.index)}"
}
resource "alicloud_slb_load_balancer" "default" {
  load_balancer_name = "${var.name}"
  vswitch_id = data.alicloud_vswitches.default.ids[0]
  load_balancer_spec  = "slb.s2.small"
}
`, name)
}

func resourceServerGroupClassicDependence(name string) string {
	return fmt.Sprintf(`

variable "name" {
  default = "%s"
}
data "alicloud_instance_types" "default" {
  cpu_core_count    = 1
  memory_size       = 2
 availability_zone = data.alicloud_slb_zones.default.zones.0.id
}
data "alicloud_images" "default" {
  name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
  owners = "system"
}
data "alicloud_vpcs" "default"{
	name_regex = "default-NODELETING"
}
data "alicloud_slb_zones" "default" {
	available_slb_address_type = "vpc"
}

data "alicloud_vswitches" "default" {
	vpc_id  = data.alicloud_vpcs.default.ids.0
	zone_id = data.alicloud_slb_zones.default.zones.0.id
}
resource "alicloud_security_group" "default" {
  name = "${var.name}"
  vpc_id = data.alicloud_vpcs.default.ids.0
}
resource "alicloud_instance" "default" {
  image_id = "${data.alicloud_images.default.images.0.id}"
  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  instance_name = "${var.name}"
  count = "2"
  security_groups = "${alicloud_security_group.default.*.id}"
  internet_charge_type = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone = "${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"
  instance_charge_type = "PostPaid"
  system_disk_category = "cloud_efficiency"
  vswitch_id = data.alicloud_vswitches.default.ids[0]
}
resource "alicloud_slb_load_balancer" "default" {
  load_balancer_name = "${var.name}"
  load_balancer_spec = "slb.s1.small"
}
`, name)

}

func resourceSlbServerGroupMultiVpcDependence(name string) string {
	return fmt.Sprintf(`

variable "name" {
  default = "tf-testAccSlbServerGroupVpc"
}
data "alicloud_instance_types" "default" {
  cpu_core_count    = 1
  memory_size       = 2
  availability_zone = data.alicloud_slb_zones.default.zones.0.id
}
data "alicloud_images" "default" {
  name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
  owners = "system"
}
data "alicloud_vpcs" "default"{
	name_regex = "default-NODELETING"
}
data "alicloud_slb_zones" "default" {
	available_slb_address_type = "vpc"
}

data "alicloud_vswitches" "default" {
	vpc_id  = data.alicloud_vpcs.default.ids.0
	zone_id = data.alicloud_slb_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  name = "${var.name}"
  vpc_id = data.alicloud_vpcs.default.ids.0
}
resource "alicloud_instance" "default" {
  image_id = "${data.alicloud_images.default.images.0.id}"
  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  instance_name = "${var.name}"
  count = "2"
  security_groups = "${alicloud_security_group.default.*.id}"
  internet_charge_type = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone = "${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"
  instance_charge_type = "PostPaid"
  system_disk_category = "cloud_efficiency"
  vswitch_id = data.alicloud_vswitches.default.ids[0]
}
resource "alicloud_slb_load_balancer" "default" {
  load_balancer_name = "${var.name}"
  vswitch_id = data.alicloud_vswitches.default.ids[0]
  load_balancer_spec = "slb.s1.small"
}
`)

}

var serverGroupMap = map[string]string{
	"servers.#":        "0",
	"load_balancer_id": CHECKSET,
}

var serverGroupMultiVpcMap = map[string]string{
	"servers.#":        "2",
	"load_balancer_id": CHECKSET,
}

var serverGroupMultiClassicMap = map[string]string{
	"servers.#":        "0",
	"load_balancer_id": CHECKSET,
}

var serversMap = []map[string]interface{}{
	{
		"server_ids": []string{"${alicloud_instance.default.0.id}"},
		"port":       "1",
		"weight":     "10",
	},
	{
		"server_ids": []string{"${alicloud_instance.default.1.id}"},
		"port":       "2",
		"weight":     "10",
	},
	{
		"server_ids": []string{"${alicloud_instance.default.2.id}"},
		"port":       "3",
		"weight":     "10",
	},
	{
		"server_ids": []string{"${alicloud_instance.default.3.id}"},
		"port":       "4",
		"weight":     "10",
	},
	{
		"server_ids": []string{"${alicloud_instance.default.4.id}"},
		"port":       "5",
		"weight":     "10",
	},
	{
		"server_ids": []string{"${alicloud_instance.default.5.id}"},
		"port":       "6",
		"weight":     "10",
	},
	{
		"server_ids": []string{"${alicloud_instance.default.6.id}"},
		"port":       "7",
		"weight":     "10",
	},
	{
		"server_ids": []string{"${alicloud_instance.default.7.id}"},
		"port":       "8",
		"weight":     "10",
	},
	{
		"server_ids": []string{"${alicloud_instance.default.8.id}"},
		"port":       "9",
		"weight":     "10",
	},
	{
		"server_ids": []string{"${alicloud_instance.default.9.id}"},
		"port":       "10",
		"weight":     "10",
	},
	{
		"server_ids": []string{"${alicloud_instance.default.10.id}"},
		"port":       "11",
		"weight":     "10",
	},
	{
		"server_ids": []string{"${alicloud_instance.default.11.id}"},
		"port":       "12",
		"weight":     "10",
	},
	{
		"server_ids": []string{"${alicloud_instance.default.12.id}"},
		"port":       "13",
		"weight":     "10",
	},
	{
		"server_ids": []string{"${alicloud_instance.default.13.id}"},
		"port":       "14",
		"weight":     "10",
	},
}
