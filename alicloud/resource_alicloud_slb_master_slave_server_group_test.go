package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudSLBMasterSlaveServerGroup_vpc(t *testing.T) {
	var v *slb.DescribeMasterSlaveServerGroupAttributeResponse
	resourceId := "alicloud_slb_master_slave_server_group.default"
	ra := resourceAttrInit(resourceId, testAccSlbMasterSlaveServerGroupCheckMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccSlbMasterSlaveServerGroupVpc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceMasterSlaveServerGroupConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		//module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				//Config: testAccSlbMasterSlaveServerGroupVpc,
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_id": "${alicloud_slb_load_balancer.default.id}",
					"name":             "${var.name}",
					"servers": []map[string]interface{}{
						{
							"server_id":   "${alicloud_instance.default.0.id}",
							"port":        "100",
							"weight":      "100",
							"server_type": "Master",
						},
						{
							"server_id":   "${alicloud_instance.default.1.id}",
							"port":        "100",
							"weight":      "100",
							"server_type": "Slave",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":      name,
						"servers.#": "2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_protection_validation"},
			},
		},
	})
}

func TestAccAlicloudSLBMasterSlaveServerGroup_multi_vpc(t *testing.T) {
	var v *slb.DescribeMasterSlaveServerGroupAttributeResponse
	resourceId := "alicloud_slb_master_slave_server_group.default.1"
	ra := resourceAttrInit(resourceId, testAccSlbMasterSlaveServerGroupCheckMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccSlbMasterSlaveServerGroupVpc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceMasterSlaveServerGroupConfigDependence)

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
					"count":            "10",
					"load_balancer_id": "${alicloud_slb_load_balancer.default.id}",
					"name":             "${var.name}",
					"servers": []map[string]interface{}{
						{
							"server_id":   "${alicloud_instance.default.0.id}",
							"port":        "100",
							"weight":      "100",
							"server_type": "Master",
						},
						{
							"server_id":   "${alicloud_instance.default.1.id}",
							"port":        "100",
							"weight":      "100",
							"server_type": "Slave",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":      name,
						"servers.#": "2",
					}),
				),
			},
		},
	})
}

func resourceMasterSlaveServerGroupConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}
data "alicloud_zones" "default" {
    available_disk_category = "cloud_efficiency"
    available_resource_creation = "VSwitch"
}
data "alicloud_instance_types" "default" {
	eni_amount        = 2
}
data "alicloud_images" "default" {
    name_regex = "^ubuntu"
    most_recent = true
    owners = "system"
}
resource "alicloud_vpc" "default" {
    vpc_name = "${var.name}"
    cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
    vpc_id = "${alicloud_vpc.default.id}"
    cidr_block = "172.16.0.0/16"
    availability_zone = "${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"
    name = "${var.name}"
}
resource "alicloud_security_group" "default" {
    name = "${var.name}"
    vpc_id = "${alicloud_vpc.default.id}"
}
resource "alicloud_network_interface" "default" {
    count = 1
    name = "${var.name}"
    vswitch_id = "${alicloud_vswitch.default.id}"
    security_groups = [ "${alicloud_security_group.default.id}" ]
}

resource "alicloud_network_interface_attachment" "default" {
    count = 1
    instance_id = "${alicloud_instance.default.0.id}"
    network_interface_id = "${element(alicloud_network_interface.default.*.id, count.index)}"
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
    vswitch_id = "${alicloud_vswitch.default.id}"
}
resource "alicloud_slb_load_balancer" "default" {
    load_balancer_name = "${var.name}"
    vswitch_id = "${alicloud_vswitch.default.id}"
    load_balancer_spec  = "slb.s2.small"
}
`, name)
}

var testAccSlbMasterSlaveServerGroupCheckMap = map[string]string{
	"name":      CHECKSET,
	"servers.#": "2",
}
