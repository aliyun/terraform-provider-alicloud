package alicloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alicloud/alicloud/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccAlicloudSLBBackendServers_vpc(t *testing.T) {
	var v *slb.DescribeLoadBalancerAttributeResponse
	resourceId := "alicloud_slb_backend_server.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccSlbBackendServersVpc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceBackendServerVpcCountConfigDependence)

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
					"backend_servers": []map[string]interface{}{
						{
							"server_id": "${alicloud_instance.instance.0.id}",
							"weight":    "80",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backend_servers.#": "1",
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
					"load_balancer_id": "${alicloud_slb_load_balancer.default.id}",
					"backend_servers": []map[string]interface{}{
						{
							"server_id": "${alicloud_instance.instance.0.id}",
							"weight":    "80",
						},
						{
							"server_id": "${alicloud_instance.instance.1.id}",
							"weight":    "80",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backend_servers.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_id": "${alicloud_slb_load_balancer.default.id}",
					"backend_servers": []map[string]interface{}{
						{
							"server_id": "${alicloud_network_interface.default.0.id}",
							"weight":    "80",
							"type":      "eni",
							"server_ip": "${alicloud_network_interface.default.0.private_ip}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backend_servers.#": "1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_id": "${alicloud_slb_load_balancer.default.id}",
					"backend_servers":  buildBackendServersMap(21),
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backend_servers.#": "21",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_id": "${alicloud_slb_load_balancer.default.id}",
					"backend_servers": []map[string]interface{}{
						{
							"server_id": "${alicloud_instance.instance.0.id}",
							"weight":    "80",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backend_servers.#": "1",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudSLBBackendServers_multi_vpc(t *testing.T) {

	var v *slb.DescribeLoadBalancerAttributeResponse
	resourceId := "alicloud_slb_backend_server.default.1"
	ra := resourceAttrInit(resourceId, slb_vpc)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccSlbBackendServersVpc_multi%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceBackendServerConfigDependence)

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
					"count":            "2",
					"load_balancer_id": "${alicloud_slb_load_balancer.default.id}",
					"backend_servers": []map[string]interface{}{
						{
							"server_id": "${alicloud_instance.instance.0.id}",
							"weight":    "80",
						},
						{
							"server_id": "${alicloud_instance.instance.1.id}",
							"weight":    "80",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backend_servers.#": "2",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudSLBBackendServers_classic(t *testing.T) {
	var v *slb.DescribeLoadBalancerAttributeResponse
	resourceId := "alicloud_slb_backend_server.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccSlbBackendServersVpc_multi%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceBackendServerConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckWithRegions(t, false, connectivity.SlbClassicNoSupportedRegions)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(), //testAccCheckSlbBackendServersDestroy,
		Steps: []resource.TestStep{
			{
				//Config: testAccSlbBackendServersClassic,
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_id": "${alicloud_slb_load_balancer.default.id}",
					"backend_servers": []map[string]interface{}{
						{
							"server_id": "${alicloud_instance.instance.0.id}",
							"weight":    "80",
						},
						{
							"server_id": "${alicloud_instance.instance.1.id}",
							"weight":    "80",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backend_servers.#": "2",
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
				//Config: testAccSlbBackendServersClassicUpdateServer,
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_id": "${alicloud_slb_load_balancer.default.id}",
					"backend_servers": []map[string]interface{}{
						{
							"server_id": "${alicloud_instance.instance.0.id}",
							"weight":    "80",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backend_servers.#": "1",
					}),
				),
			},
			{
				//Config: testAccSlbBackendServersClassic,
				Config: testAccConfig(map[string]interface{}{
					"load_balancer_id": "${alicloud_slb_load_balancer.default.id}",
					"backend_servers": []map[string]interface{}{
						{
							"server_id": "${alicloud_instance.instance.0.id}",
							"weight":    "80",
						},
						{
							"server_id": "${alicloud_instance.instance.1.id}",
							"weight":    "80",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"backend_servers.#": "2",
					}),
				),
			},
		},
	})
}

func buildBackendServersMap(count int) []map[string]interface{} {
	var result []map[string]interface{}

	str := `${alicloud_instance.instance.%d.id}`
	for i := 0; i < count; i++ {
		tmp := make(map[string]interface{}, 2)
		tmp["server_id"] = fmt.Sprintf(str, i)
		tmp["weight"] = "10"
		result = append(result, tmp)
	}
	return result
}

func resourceBackendServerVpcCountConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSlbBackendServersVpc"
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
resource "alicloud_security_group" "group" {
  name   = "${var.name}"
  vpc_id = data.alicloud_vpcs.default.ids.0
}
resource "alicloud_instance" "instance" {
  image_id                   = "${data.alicloud_images.default.images.0.id}"
  instance_type              = "${data.alicloud_instance_types.default.instance_types.0.id}"
  instance_name              = "${var.name}"
  count                      = "21"
  security_groups            = "${alicloud_security_group.group.*.id}"
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone          = "${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"
  instance_charge_type       = "PostPaid"
  system_disk_category       = "cloud_efficiency"
  vswitch_id                 = data.alicloud_vswitches.default.ids[0]
}
resource "alicloud_slb_load_balancer" "default" {
  load_balancer_name          = "${var.name}"
  vswitch_id    = data.alicloud_vswitches.default.ids[0]
  load_balancer_spec = "slb.s2.small"
}


resource "alicloud_network_interface" "default" {
    count = 1
    name = "${var.name}"
    vswitch_id = data.alicloud_vswitches.default.ids[0]
    security_groups = [ "${alicloud_security_group.group.id}" ]
}
resource "alicloud_instance" "new" {
  image_id = "${data.alicloud_images.default.images.0.id}"
  instance_type = "${data.alicloud_instance_types.new.instance_types.0.id}"
  instance_name = "${var.name}"
  count = "1"
  security_groups = "${alicloud_security_group.group.*.id}"
  internet_charge_type = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone = "${data.alicloud_instance_types.default.instance_types.0.availability_zones.0}"
  instance_charge_type = "PostPaid"
  system_disk_category = "cloud_efficiency"
  vswitch_id = data.alicloud_vswitches.default.ids[0]
}
resource "alicloud_network_interface_attachment" "default" {
	count = 1
    instance_id = "${alicloud_instance.new.0.id}"
    network_interface_id = "${element(alicloud_network_interface.default.*.id, count.index)}"
}
`)
}

func resourceBackendServerConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "tf-testAccSlbBackendServersVpc"
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
resource "alicloud_instance" "instance" {
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

var slb_vpc = map[string]string{
	"backend_servers.#": "2",
}
