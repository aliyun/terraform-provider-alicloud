package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/terraform-providers/terraform-provider-alicloud/alicloud/connectivity"
)

func TestAccAlicloudSlbServerGroup_vpc(t *testing.T) {
	var v *slb.DescribeVServerGroupAttributeResponse
	resourceId := "alicloud_slb_server_group.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccSlbServerGroupVpc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceSlbServerGroupDependence)
	testAccVpcUpdateServerConfig := resourceTestAccConfigFunc(resourceId, name, resourceSlbServerGroupVpcUpdateServerDependence)

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
					"load_balancer_id": "${alicloud_slb.default.id}",
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
						"name":      "tf-server-group",
						"servers.#": "1",
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
					"name": "${var.update_name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": "tf-testAccSlbServerGroupVpcUpdate",
					}),
				),
			},
			{
				Config: testAccVpcUpdateServerConfig(map[string]interface{}{
					"depends_on":       []string{"alicloud_network_interface_attachment.default"},
					"load_balancer_id": "${alicloud_slb.default.id}",
					"name":             "${var.name}",
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
				Config: testAccSlbServerGroupVpcUpdateServerSize(21),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"servers.#": "21",
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

func TestAccAlicloudSlbServerGroup_multi_vpc(t *testing.T) {
	var v *slb.DescribeVServerGroupAttributeResponse
	resourceId := "alicloud_slb_server_group.default.9"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccSlbServerGroupMultiVpc%d", rand)
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
					"load_balancer_id": "${alicloud_slb.default.id}",
					"name":             "${var.name}",
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
					testAccCheck(map[string]string{
						"name":      "tf-testAccSlbServerGroupVpc",
						"servers.#": "2",
					}),
				),
			},
		},
	})
}

func TestAccAlicloudSlbServerGroup_classic(t *testing.T) {
	var v *slb.DescribeVServerGroupAttributeResponse
	resourceId := "alicloud_slb_server_group.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testAccSlbServerGroupClassic%d", rand)
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
					"load_balancer_id": "${alicloud_slb.default.id}",
					//"name":             "${var.name}",
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
						//"name":      "tf-testAccSlbServerGroupClassic",
						"servers.#": "2",
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
					//"load_balancer_id": "${alicloud_slb.default.id}",
					"name": "${var.update_name}",
					//"servers": []map[string]interface{}{
					//	{
					//		"server_ids": []string{"${alicloud_instance.default.0.id}", "${alicloud_instance.default.1.id}"},
					//		"port":       "100",
					//		"weight":     "10",
					//	},
					//	{
					//		"server_ids": "${alicloud_instance.default.*.id}",
					//		"port":       "80",
					//		"weight":     "100",
					//	},
					//},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": "tf-testAccSlbServerGroupClassicUpdate",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					//"load_balancer_id": "${alicloud_slb.default.id}",
					//"name":             "${var.update_name}",
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
					//"load_balancer_id": "${alicloud_slb.default.id}",
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

func buildservers(count int) string {
	var result string
	temp := `
  servers {
      server_ids = ["${alicloud_instance.default.%d.id}"]
    port       = %d
    weight     = 10
  }
  `
	for i := 0; i < count; i++ {
		result += fmt.Sprintf(temp, i, i+1)
	}
	return result
}

func resourceSlbServerGroupDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSlbServerGroupVpc"
}

variable "update_name" {
  default = "tf-testAccSlbServerGroupVpcUpdate"
}

data "alicloud_instance_types" "default" {
  cpu_core_count    = 1
  memory_size       = 2
}
data "alicloud_images" "default" {
  name_regex = "^ubuntu_18.*_64"
  most_recent = true
  owners = "system"
}
resource "alicloud_vpc" "default" {
  name = "${var.name}"
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
  vswitch_id = "${alicloud_vswitch.default.id}"
}
resource "alicloud_slb" "default" {
  name = "${var.name}"
  vswitch_id = "${alicloud_vswitch.default.id}"
  specification  = "slb.s2.small"
}

  `)
}

func resourceSlbServerGroupVpcUpdateServerDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSlbServerGroupVpcUpdate"
}

data "alicloud_instance_types" "default" {
  cpu_core_count    = 1
  memory_size       = 2
}
data "alicloud_instance_types" "new" {
  eni_amount = 2
}
data "alicloud_images" "default" {
        name_regex = "^ubuntu_18.*_64"
  most_recent = true
  owners = "system"
}

resource "alicloud_vpc" "default" {
  name = "${var.name}"
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
  vswitch_id = "${alicloud_vswitch.default.id}"
}
resource "alicloud_network_interface_attachment" "default" {
  count = 1
    instance_id = "${alicloud_instance.new.0.id}"
    network_interface_id = "${element(alicloud_network_interface.default.*.id, count.index)}"
}

resource "alicloud_slb" "default" {
  name = "${var.name}"
  vswitch_id = "${alicloud_vswitch.default.id}"
  specification  = "slb.s2.small"
}
`)
}

func testAccSlbServerGroupVpcUpdateServerSize(count int) string {
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccSlbServerGroupVpcUpdate"
}

data "alicloud_instance_types" "default" {
  cpu_core_count    = 1
  memory_size       = 2
}
data "alicloud_instance_types" "new" {
  eni_amount = 2
}
data "alicloud_images" "default" {
        name_regex = "^ubuntu_18.*_64"
  most_recent = true
  owners = "system"
}

resource "alicloud_vpc" "default" {
  name = "${var.name}"
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
  vswitch_id = "${alicloud_vswitch.default.id}"
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
  vswitch_id = "${alicloud_vswitch.default.id}"
}
resource "alicloud_network_interface_attachment" "default" {
  count = 1
    instance_id = "${alicloud_instance.new.0.id}"
    network_interface_id = "${element(alicloud_network_interface.default.*.id, count.index)}"
}

resource "alicloud_slb" "default" {
  name = "${var.name}"
  vswitch_id = "${alicloud_vswitch.default.id}"
  specification  = "slb.s2.small"
}

resource "alicloud_slb_server_group" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  name = "${var.name}"
  %s
}
`, buildservers(count))
}

func resourceServerGroupClassicDependence(name string) string {
	return fmt.Sprintf(`

variable "name" {
  default = "tf-testAccSlbServerGroupClassic"
}
variable "update_name" {
  default = "tf-testAccSlbServerGroupClassicUpdate"
}

data "alicloud_instance_types" "default" {
  cpu_core_count    = 1
  memory_size       = 2
}
data "alicloud_images" "default" {
  name_regex = "^ubuntu_18.*_64"
  most_recent = true
  owners = "system"
}

resource "alicloud_vpc" "default" {
  name = "${var.name}"
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

resource "alicloud_slb" "default" {
  name = "${var.name}"
}

`)

}

func resourceSlbServerGroupMultiVpcDependence(name string) string {
	return fmt.Sprintf(`

variable "name" {
  default = "tf-testAccSlbServerGroupVpc"
}

data "alicloud_instance_types" "default" {
  cpu_core_count    = 1
  memory_size       = 2
}
data "alicloud_images" "default" {
        name_regex = "^ubuntu_18.*_64"
  most_recent = true
  owners = "system"
}
resource "alicloud_vpc" "default" {
  name = "${var.name}"
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
resource "alicloud_slb" "default" {
  name = "${var.name}"
  vswitch_id = "${alicloud_vswitch.default.id}"
}
`)

}
