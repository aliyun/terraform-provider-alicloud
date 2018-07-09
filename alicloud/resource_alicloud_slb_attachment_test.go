package alicloud

import (
	"fmt"
	"log"
	"testing"

	"github.com/denverdino/aliyungo/slb"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAlicloudSlbAttachment_basic(t *testing.T) {
	var slb slb.LoadBalancerType

	testCheckAttr := func() resource.TestCheckFunc {
		return func(*terraform.State) error {
			log.Printf("testCheckAttr slb BackendServers is: %#v", slb.BackendServers)
			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: "alicloud_slb_attachment.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSlbDestroy,
		Steps: []resource.TestStep{
			//test internet_charge_type is paybybandwidth
			resource.TestStep{
				Config: testAccSlbAttachment,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSlbExists("alicloud_slb_attachment.foo", &slb),
					testCheckAttr(),
					testAccCheckAttachment("alicloud_instance.foo", &slb),
					resource.TestCheckResourceAttr(
						"alicloud_slb_attachment.foo",
						"weight", "90"),
				),
			},
		},
	})
}

func testAccCheckAttachment(n string, slb *slb.LoadBalancerType) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ECS ID is set")
		}

		ecsInstanceId := rs.Primary.ID

		backendServers := slb.BackendServers.BackendServer

		if len(backendServers) == 0 {
			return fmt.Errorf("no SLB backendServer: %#v", backendServers)
		}

		log.Printf("slb bacnendservers: %#v", backendServers)

		backendServersInstanceId := backendServers[0].ServerId

		if ecsInstanceId != backendServersInstanceId {
			return fmt.Errorf("SLB attachment check invalid: ECS instance %s is not equal SLB backendServer %s",
				ecsInstanceId, backendServersInstanceId)
		}
		return nil
	}
}

const testAccSlbAttachment = `
data "alicloud_zones" "default" {
	"available_disk_category"= "cloud_efficiency"
	"available_resource_creation"= "VSwitch"
}
data "alicloud_instance_types" "default" {
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}
data "alicloud_images" "default" {
        name_regex = "^ubuntu_14.*_64"
	most_recent = true
	owners = "system"
}
variable "name" {
	default = "testAccSlbAttachment"
}

resource "alicloud_vpc" "main" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "main" {
	vpc_id = "${alicloud_vpc.main.id}"
	cidr_block = "172.16.0.0/16"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	depends_on = [
	"alicloud_vpc.main"]
}

resource "alicloud_security_group" "group" {
	name = "${var.name}"
	vpc_id = "${alicloud_vpc.main.id}"
}

resource "alicloud_instance" "foo" {
	# cn-beijing
	image_id = "${data.alicloud_images.default.images.0.id}"

	# series III
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	internet_charge_type = "PayByBandwidth"
	internet_max_bandwidth_out = "5"
	system_disk_category = "cloud_efficiency"

	security_groups = ["${alicloud_security_group.group.id}"]
	instance_name = "${var.name}"
	vswitch_id = "${alicloud_vswitch.main.id}"
}

resource "alicloud_slb" "foo" {
	name = "${var.name}"
	vswitch_id = "${alicloud_vswitch.main.id}"
}

resource "alicloud_slb_attachment" "foo" {
	load_balancer_id = "${alicloud_slb.foo.id}"
	instance_ids = ["${alicloud_instance.foo.id}"]
	weight = 90
}

`
