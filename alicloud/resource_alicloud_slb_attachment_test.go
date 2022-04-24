package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudSLBAttachment_basic(t *testing.T) {

	var v *slb.DescribeLoadBalancerAttributeResponse
	resourceId := "alicloud_slb_attachment.default"
	ra := resourceAttrInit(resourceId, attachmentMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlb")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccSlbAttachment")
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceSlbAttachmentBasicdependence)
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
					"instance_ids":     []string{"${alicloud_instance.default.0.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
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
					"weight": "70",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"weight": "70",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_ids": []string{"${alicloud_instance.default.0.id}", "${alicloud_instance.default.1.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_ids.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"depends_on":   []string{"alicloud_network_interface_attachment.default"},
					"instance_ids": []string{"${alicloud_network_interface.default.0.id}"},
					"server_type":  "eni",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_ids.#": "1",
						"server_type":    "eni",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_id": "${alicloud_slb_load_balancer.default.id}",
					"instance_ids":     []string{"${alicloud_instance.default.0.id}"},
					"weight":           "90",
					"server_type":      "ecs",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_id": CHECKSET,
						"weight":           "90",
						"instance_ids.#":   "1",
						"backend_servers":  CHECKSET,
						"server_type":      "ecs",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudSLBAttachment_multi(t *testing.T) {

	var v *slb.DescribeLoadBalancerAttributeResponse
	resourceId := "alicloud_slb_attachment.default.9"
	ra := resourceAttrInit(resourceId, attachmentMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlb")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccSlbAttachmentMulti")
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceSlbAttachmentMultidependence)

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
					"instance_ids":     []string{"${alicloud_instance.default.id}"},
					"count":            "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func TestAccAlicloudSLBAttachment_classic_basic(t *testing.T) {

	var v *slb.DescribeLoadBalancerAttributeResponse
	resourceId := "alicloud_slb_attachment.default"
	ra := resourceAttrInit(resourceId, attachmentMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	}, "DescribeSlb")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	name := fmt.Sprintf("tf-testAccSlbAttachment")
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceSlbAttachmentClassBasicdependence)
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
					"instance_ids":     []string{"${alicloud_instance.default.0.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"weight": "70",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"weight": "70",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"instance_ids": []string{"${alicloud_instance.default.0.id}", "${alicloud_instance.default.1.id}"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_ids.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"depends_on":   []string{"alicloud_network_interface_attachment.default"},
					"instance_ids": []string{"${alicloud_network_interface.default.0.id}"},
					"server_type":  "eni",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"instance_ids.#": "1",
						"server_type":    "eni",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_id": "${alicloud_slb_load_balancer.default.id}",
					"instance_ids":     []string{"${alicloud_instance.default.0.id}"},
					"weight":           "90",
					"server_type":      "ecs",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"load_balancer_id": CHECKSET,
						"weight":           "90",
						"instance_ids.#":   "1",
						"backend_servers":  CHECKSET,
						"server_type":      "ecs",
					}),
				),
			},
		},
	})
}

func resourceSlbAttachmentBasicdependence(name string) string {
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
  name_regex = "^ubuntu_[0-9]+_[0-9]+_x64*"
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
    vswitch_id = data.alicloud_vswitches.default.ids[0]
    security_groups = [ "${alicloud_security_group.default.id}" ]
}
resource "alicloud_instance" "default" {
  # cn-beijing
  image_id = "${data.alicloud_images.default.images.0.id}"

  # series III
  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  internet_charge_type = "PayByTraffic"
  internet_max_bandwidth_out = "5"
  system_disk_category = "cloud_efficiency"
  count = 2
  security_groups = ["${alicloud_security_group.default.id}"]
  instance_name = "${var.name}"
  vswitch_id = data.alicloud_vswitches.default.ids[0]
  availability_zone = "${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"
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

func resourceSlbAttachmentMultidependence(name string) string {
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
  name_regex = "^ubuntu_[0-9]+_[0-9]+_x64*"
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
  # cn-beijing
  image_id = "${data.alicloud_images.default.images.0.id}"

  # series III
  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  internet_charge_type = "PayByTraffic"
  internet_max_bandwidth_out = "5"
  system_disk_category = "cloud_efficiency"

  security_groups = ["${alicloud_security_group.default.id}"]
  instance_name = "${var.name}"
  vswitch_id = data.alicloud_vswitches.default.ids[0]
  availability_zone = "${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"
}
resource "alicloud_slb_load_balancer" "default" {
  load_balancer_name = "${var.name}"
  vswitch_id = data.alicloud_vswitches.default.ids[0]
  load_balancer_spec = "slb.s1.small"
}
`, name)
}

func resourceSlbAttachmentClassBasicdependence(name string) string {
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
  name_regex = "^ubuntu_[0-9]+_[0-9]+_x64*"
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
    vswitch_id = data.alicloud_vswitches.default.ids[0]
    security_groups = [ "${alicloud_security_group.default.id}" ]
}
resource "alicloud_instance" "default" {
  # cn-beijing
  image_id = "${data.alicloud_images.default.images.0.id}"

  # series III
  instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
  internet_charge_type = "PayByTraffic"
  internet_max_bandwidth_out = "5"
  system_disk_category = "cloud_efficiency"
  count = 2
  security_groups = ["${alicloud_security_group.default.id}"]
  instance_name = "${var.name}"
  vswitch_id = data.alicloud_vswitches.default.ids[0]
  availability_zone = "${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"
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
  load_balancer_spec  = "slb.s2.small"
}
`, name)
}

var attachmentMap = map[string]string{
	"load_balancer_id": CHECKSET,
	"instance_ids.#":   "1",
}
