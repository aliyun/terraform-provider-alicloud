package alicloud

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
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
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSlbServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSlbServerGroupVpc,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":      "tf-testAccSlbServerGroupVpc",
						"servers.#": "2",
					}),
				),
			},
			{
				Config: testAccSlbServerGroupVpcUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":      "tf-testAccSlbServerGroupVpcUpdate",
						"servers.#": "2",
					}),
				),
			},
			{
				Config: testAccSlbServerGroupVpcUpdateServer,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"servers.#": "1",
					}),
				),
			},
			{
				Config: testAccSlbServerGroupVpc,
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

func TestAccAlicloudSlbServerGroup_multi_vpc(t *testing.T) {
	var v *slb.DescribeVServerGroupAttributeResponse
	resourceId := "alicloud_slb_server_group.default.9"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &SlbService{testAccProvider.Meta().(*connectivity.AliyunClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSlbServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSlbServerGroupVpc_multi,
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
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSlbServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSlbServerGroupClassic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":      "tf-testAccSlbServerGroupClassic",
						"servers.#": "2",
					}),
				),
			},
			{
				Config: testAccSlbServerGroupClassicUpdateName,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": "tf-testAccSlbServerGroupClassicUpdate",
					}),
				),
			},
			{
				Config: testAccSlbServerGroupClassicUpdateServer,
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"servers.#": "1",
					}),
				),
			},
			{
				Config: testAccSlbServerGroupClassic,
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

func testAccCheckSlbServerGroupDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AliyunClient)
	slbService := SlbService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alicloud_slb_server_group" {
			continue
		}

		// Try to find the Slb server group
		if _, err := slbService.DescribeSlbServerGroup(rs.Primary.ID); err != nil {
			if NotFoundError(err) {
				continue
			}
			return err
		}
		return fmt.Errorf("SLB Server Group %s still exist.", rs.Primary.ID)
	}

	return nil
}

const testAccSlbServerGroupVpc = `
variable "name" {
	default = "tf-testAccSlbServerGroupVpc"
}
data "alicloud_zones" "default" {
	available_disk_category= "cloud_efficiency"
	available_resource_creation= "VSwitch"
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
resource "alicloud_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "172.16.0.0/16"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
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
  security_groups = ["${alicloud_security_group.default.*.id}"]
  internet_charge_type = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  instance_charge_type = "PostPaid"
  system_disk_category = "cloud_efficiency"
  vswitch_id = "${alicloud_vswitch.default.id}"
}
resource "alicloud_slb" "default" {
  name = "${var.name}"
  vswitch_id = "${alicloud_vswitch.default.id}"
}
resource "alicloud_slb_server_group" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  name = "${var.name}"
  servers = [
    {
      server_ids = ["${alicloud_instance.default.0.id}", "${alicloud_instance.default.1.id}"]
      port = 100
      weight = 10
    },
    {
      server_ids = ["${alicloud_instance.default.*.id}"]
      port = 80
      weight = 100
    }
  ]
}
`
const testAccSlbServerGroupVpcUpdate = `
variable "name" {
	default = "tf-testAccSlbServerGroupVpcUpdate"
}
data "alicloud_zones" "default" {
	available_disk_category= "cloud_efficiency"
	available_resource_creation= "VSwitch"
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

resource "alicloud_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "172.16.0.0/16"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
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
  security_groups = ["${alicloud_security_group.default.*.id}"]
  internet_charge_type = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  instance_charge_type = "PostPaid"
  system_disk_category = "cloud_efficiency"
  vswitch_id = "${alicloud_vswitch.default.id}"
}

resource "alicloud_slb" "default" {
  name = "${var.name}"
  vswitch_id = "${alicloud_vswitch.default.id}"
}

resource "alicloud_slb_server_group" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  name = "${var.name}"
  servers = [
    {
      server_ids = ["${alicloud_instance.default.0.id}", "${alicloud_instance.default.1.id}"]
      port = 100
      weight = 10
    },
    {
      server_ids = ["${alicloud_instance.default.*.id}"]
      port = 80
      weight = 100
    }
  ]
}
`
const testAccSlbServerGroupVpcUpdateServer = `
variable "name" {
	default = "tf-testAccSlbServerGroupVpcUpdate"
}
data "alicloud_zones" "default" {
	available_disk_category = "cloud_efficiency"
	available_resource_creation = "VSwitch"
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

resource "alicloud_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "172.16.0.0/16"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
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
  security_groups = ["${alicloud_security_group.default.*.id}"]
  internet_charge_type = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  instance_charge_type = "PostPaid"
  system_disk_category = "cloud_efficiency"
  vswitch_id = "${alicloud_vswitch.default.id}"
}

resource "alicloud_slb" "default" {
  name = "${var.name}"
  vswitch_id = "${alicloud_vswitch.default.id}"
}

resource "alicloud_slb_server_group" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  name = "${var.name}"
  servers = [
    {
      server_ids = ["${alicloud_instance.default.0.id}", "${alicloud_instance.default.1.id}"]
      port = 100
      weight = 10
    }  
  ]
}
`

const testAccSlbServerGroupClassic = `
variable "name" {
	default = "tf-testAccSlbServerGroupClassic"
}

data "alicloud_zones" "default" {
	available_disk_category= "cloud_efficiency"
	available_resource_creation= "VSwitch"
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

resource "alicloud_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "172.16.0.0/16"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
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
  security_groups = ["${alicloud_security_group.default.*.id}"]
  internet_charge_type = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  instance_charge_type = "PostPaid"
  system_disk_category = "cloud_efficiency"
  vswitch_id = "${alicloud_vswitch.default.id}"
}

resource "alicloud_slb" "default" {
  name = "${var.name}"
}

resource "alicloud_slb_server_group" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  name = "${var.name}"
  servers = [
    {
      server_ids = ["${alicloud_instance.default.0.id}", "${alicloud_instance.default.1.id}"]
      port = 100
      weight = 10
    },
	{
      server_ids = ["${alicloud_instance.default.*.id}"]
      port = 80
      weight = 100
    }
  ]
}
`

const testAccSlbServerGroupClassicUpdateName = `
variable "name" {
	default = "tf-testAccSlbServerGroupClassicUpdate"
}

data "alicloud_zones" "default" {
	available_disk_category= "cloud_efficiency"
	available_resource_creation= "VSwitch"
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

resource "alicloud_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "172.16.0.0/16"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
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
  security_groups = ["${alicloud_security_group.default.*.id}"]
  internet_charge_type = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  instance_charge_type = "PostPaid"
  system_disk_category = "cloud_efficiency"
  vswitch_id = "${alicloud_vswitch.default.id}"
}

resource "alicloud_slb" "default" {
  name = "${var.name}"
}

resource "alicloud_slb_server_group" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  name = "${var.name}"
  servers = [
    {
      server_ids = ["${alicloud_instance.default.0.id}", "${alicloud_instance.default.1.id}"]
      port = 100
      weight = 10
    },
	{
      server_ids = ["${alicloud_instance.default.*.id}"]
      port = 80
      weight = 100
    }
  ]
}
`

const testAccSlbServerGroupClassicUpdateServer = `
variable "name" {
	default = "tf-testAccSlbServerGroupClassicUpdate"
}

data "alicloud_zones" "default" {
	available_disk_category= "cloud_efficiency"
	available_resource_creation= "VSwitch"
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

resource "alicloud_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "172.16.0.0/16"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
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
  security_groups = ["${alicloud_security_group.default.*.id}"]
  internet_charge_type = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  instance_charge_type = "PostPaid"
  system_disk_category = "cloud_efficiency"
  vswitch_id = "${alicloud_vswitch.default.id}"
}

resource "alicloud_slb" "default" {
  name = "${var.name}"
}

resource "alicloud_slb_server_group" "default" {
  load_balancer_id = "${alicloud_slb.default.id}"
  name = "${var.name}"
  servers = [
    {
      server_ids = ["${alicloud_instance.default.0.id}", "${alicloud_instance.default.1.id}"]
      port = 100
      weight = 10
    }
  ]
}
`

const testAccSlbServerGroupVpc_multi = `
variable "name" {
	default = "tf-testAccSlbServerGroupVpc"
}
data "alicloud_zones" "default" {
	available_disk_category= "cloud_efficiency"
	available_resource_creation= "VSwitch"
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
resource "alicloud_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "172.16.0.0/16"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
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
  security_groups = ["${alicloud_security_group.default.*.id}"]
  internet_charge_type = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  instance_charge_type = "PostPaid"
  system_disk_category = "cloud_efficiency"
  vswitch_id = "${alicloud_vswitch.default.id}"
}
resource "alicloud_slb" "default" {
  name = "${var.name}"
  vswitch_id = "${alicloud_vswitch.default.id}"
}
resource "alicloud_slb_server_group" "default" {
	count = "10"
  load_balancer_id = "${alicloud_slb.default.id}"
  name = "${var.name}"
  servers = [
    {
      server_ids = ["${alicloud_instance.default.0.id}", "${alicloud_instance.default.1.id}"]
      port = 100
      weight = 10
    },
    {
      server_ids = ["${alicloud_instance.default.*.id}"]
      port = 80
      weight = 100
    }
  ]
}
`
