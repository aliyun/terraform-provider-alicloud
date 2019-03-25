package alicloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccAlicloudRouteEntriesDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRouteEntriesDataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_route_entries.foo"),
					resource.TestCheckResourceAttr("data.alicloud_route_entries.foo", "entries.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_route_entries.foo", "entries.0.route_table_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_route_entries.foo", "entries.0.cidr_block"),
					resource.TestCheckResourceAttrSet("data.alicloud_route_entries.foo", "entries.0.instance_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_route_entries.foo", "entries.0.status"),
					resource.TestCheckResourceAttr("data.alicloud_route_entries.foo", "entries.0.type", "Custom"),
					resource.TestCheckResourceAttr("data.alicloud_route_entries.foo", "entries.0.next_hop_type", "Instance"),
				),
			},
			{
				Config: testAccCheckAlicloudRouteEntriesDataSourceConfig_mismatch,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_route_entries.foo"),
					resource.TestCheckResourceAttr("data.alicloud_route_entries.foo", "entries.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudRouteEntriesDataSourceRouteEntryType(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRouteEntriesDataSourceRouteEntryType,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_route_entries.foo"),
					resource.TestCheckResourceAttr("data.alicloud_route_entries.foo", "entries.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_route_entries.foo", "entries.0.route_table_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_route_entries.foo", "entries.0.cidr_block"),
					resource.TestCheckResourceAttrSet("data.alicloud_route_entries.foo", "entries.0.instance_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_route_entries.foo", "entries.0.status"),
					resource.TestCheckResourceAttr("data.alicloud_route_entries.foo", "entries.0.type", "Custom"),
					resource.TestCheckResourceAttr("data.alicloud_route_entries.foo", "entries.0.next_hop_type", "Instance"),
				),
			},
			{
				Config: testAccCheckAlicloudRouteEntriesDataSourceRouteEntryType_mismatch,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_route_entries.foo"),
					// 系统路由条目
					resource.TestCheckResourceAttr("data.alicloud_route_entries.foo", "entries.#", "2"),
					resource.TestCheckResourceAttr("data.alicloud_route_entries.foo", "entries.0.type", "System"),
					resource.TestCheckResourceAttr("data.alicloud_route_entries.foo", "entries.1.type", "System"),
				),
			},
		},
	})
}

func TestAccAlicloudRouteEntriesDataSourceCidrBlock(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRouteEntriesDataSourceCidrBlock,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_route_entries.foo"),
					resource.TestCheckResourceAttr("data.alicloud_route_entries.foo", "entries.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_route_entries.foo", "entries.0.route_table_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_route_entries.foo", "entries.0.cidr_block"),
					resource.TestCheckResourceAttrSet("data.alicloud_route_entries.foo", "entries.0.instance_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_route_entries.foo", "entries.0.status"),
					resource.TestCheckResourceAttr("data.alicloud_route_entries.foo", "entries.0.type", "Custom"),
					resource.TestCheckResourceAttr("data.alicloud_route_entries.foo", "entries.0.next_hop_type", "Instance"),
				),
			},
			{
				Config: testAccCheckAlicloudRouteEntriesDataSourceCidrBlock_mismatch,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_route_entries.foo"),
					resource.TestCheckResourceAttr("data.alicloud_route_entries.foo", "entries.#", "0"),
				),
			},
		},
	})
}

func TestAccAlicloudRouteEntriesDataSourceInstanceId(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlicloudRouteEntriesDataSourceInstanceId,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_route_entries.foo"),
					resource.TestCheckResourceAttr("data.alicloud_route_entries.foo", "entries.#", "1"),
					resource.TestCheckResourceAttrSet("data.alicloud_route_entries.foo", "entries.0.route_table_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_route_entries.foo", "entries.0.cidr_block"),
					resource.TestCheckResourceAttrSet("data.alicloud_route_entries.foo", "entries.0.instance_id"),
					resource.TestCheckResourceAttrSet("data.alicloud_route_entries.foo", "entries.0.status"),
					resource.TestCheckResourceAttr("data.alicloud_route_entries.foo", "entries.0.type", "Custom"),
					resource.TestCheckResourceAttr("data.alicloud_route_entries.foo", "entries.0.next_hop_type", "Instance"),
				),
			},
			{
				Config: testAccCheckAlicloudRouteEntriesDataSourceInstanceId_mismatch,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAlicloudDataSourceID("data.alicloud_route_entries.foo"),
					resource.TestCheckResourceAttr("data.alicloud_route_entries.foo", "entries.#", "0"),
				),
			},
		},
	})
}

const testAccCheckAlicloudRouteEntriesDataSourceConfig = `
data "alicloud_zones" "default" {
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
	default = "tf-testAcc-for-route-entries-datasource"
}
resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "10.1.1.0/24"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_route_entry" "foo" {
	route_table_id = "${alicloud_vpc.foo.route_table_id}"
	destination_cidrblock = "172.11.1.1/32"
	nexthop_type = "Instance"
	nexthop_id = "${alicloud_instance.foo.id}"
}

resource "alicloud_security_group" "tf_test_foo" {
	name = "${var.name}"
	description = "foo"
	vpc_id = "${alicloud_vpc.foo.id}"
}

resource "alicloud_security_group_rule" "ingress" {
	type = "ingress"
	ip_protocol = "tcp"
	nic_type = "intranet"
	policy = "accept"
	port_range = "22/22"
	priority = 1
	security_group_id = "${alicloud_security_group.tf_test_foo.id}"
	cidr_ip = "0.0.0.0/0"
}

resource "alicloud_instance" "foo" {
	# cn-beijing
	security_groups = ["${alicloud_security_group.tf_test_foo.id}"]

	vswitch_id = "${alicloud_vswitch.foo.id}"
	allocate_public_ip = true

	# series III
	instance_charge_type = "PostPaid"
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	internet_charge_type = "PayByTraffic"
	internet_max_bandwidth_out = 5

	system_disk_category = "cloud_efficiency"
	image_id = "${data.alicloud_images.default.images.0.id}"
	instance_name = "${var.name}"
}

data "alicloud_route_entries" "foo" {
  type = "Custom"
  cidr_block = "${alicloud_route_entry.foo.destination_cidrblock}"
  instance_id = "${alicloud_instance.foo.id}"
  route_table_id = "${alicloud_route_entry.foo.route_table_id}"
}
`
const testAccCheckAlicloudRouteEntriesDataSourceConfig_mismatch = `
data "alicloud_zones" "default" {
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
	default = "tf-testAcc-for-route-entries-datasource"
}
resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "10.1.1.0/24"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_route_entry" "foo" {
	route_table_id = "${alicloud_vpc.foo.route_table_id}"
	destination_cidrblock = "172.11.1.1/32"
	nexthop_type = "Instance"
	nexthop_id = "${alicloud_instance.foo.id}"
}

resource "alicloud_security_group" "tf_test_foo" {
	name = "${var.name}"
	description = "foo"
	vpc_id = "${alicloud_vpc.foo.id}"
}

resource "alicloud_security_group_rule" "ingress" {
	type = "ingress"
	ip_protocol = "tcp"
	nic_type = "intranet"
	policy = "accept"
	port_range = "22/22"
	priority = 1
	security_group_id = "${alicloud_security_group.tf_test_foo.id}"
	cidr_ip = "0.0.0.0/0"
}

resource "alicloud_instance" "foo" {
	# cn-beijing
	security_groups = ["${alicloud_security_group.tf_test_foo.id}"]

	vswitch_id = "${alicloud_vswitch.foo.id}"
	allocate_public_ip = true

	# series III
	instance_charge_type = "PostPaid"
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	internet_charge_type = "PayByTraffic"
	internet_max_bandwidth_out = 5

	system_disk_category = "cloud_efficiency"
	image_id = "${data.alicloud_images.default.images.0.id}"
	instance_name = "${var.name}"
}

data "alicloud_route_entries" "foo" {
  type = "System"
  cidr_block = "${alicloud_route_entry.foo.destination_cidrblock}-fake"
  instance_id = "${alicloud_instance.foo.id}-fake"
  route_table_id = "${alicloud_route_entry.foo.route_table_id}-fake"
}
`

const testAccCheckAlicloudRouteEntriesDataSourceRouteEntryType = `
data "alicloud_zones" "default" {
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
	default = "tf-testAcc-for-route-entries-datasource"
}
resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "10.1.1.0/24"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_route_entry" "foo" {
	route_table_id = "${alicloud_vpc.foo.route_table_id}"
	destination_cidrblock = "172.11.1.1/32"
	nexthop_type = "Instance"
	nexthop_id = "${alicloud_instance.foo.id}"
}

resource "alicloud_security_group" "tf_test_foo" {
	name = "${var.name}"
	description = "foo"
	vpc_id = "${alicloud_vpc.foo.id}"
}

resource "alicloud_security_group_rule" "ingress" {
	type = "ingress"
	ip_protocol = "tcp"
	nic_type = "intranet"
	policy = "accept"
	port_range = "22/22"
	priority = 1
	security_group_id = "${alicloud_security_group.tf_test_foo.id}"
	cidr_ip = "0.0.0.0/0"
}

resource "alicloud_instance" "foo" {
	# cn-beijing
	security_groups = ["${alicloud_security_group.tf_test_foo.id}"]

	vswitch_id = "${alicloud_vswitch.foo.id}"
	allocate_public_ip = true

	# series III
	instance_charge_type = "PostPaid"
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	internet_charge_type = "PayByTraffic"
	internet_max_bandwidth_out = 5

	system_disk_category = "cloud_efficiency"
	image_id = "${data.alicloud_images.default.images.0.id}"
	instance_name = "${var.name}"
}

data "alicloud_route_entries" "foo" {
  type = "Custom"
  route_table_id = "${alicloud_route_entry.foo.route_table_id}"
}
`

const testAccCheckAlicloudRouteEntriesDataSourceRouteEntryType_mismatch = `
data "alicloud_zones" "default" {
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
	default = "tf-testAcc-for-route-entries-datasource"
}
resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "10.1.1.0/24"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_route_entry" "foo" {
	route_table_id = "${alicloud_vpc.foo.route_table_id}"
	destination_cidrblock = "172.11.1.1/32"
	nexthop_type = "Instance"
	nexthop_id = "${alicloud_instance.foo.id}"
}

resource "alicloud_security_group" "tf_test_foo" {
	name = "${var.name}"
	description = "foo"
	vpc_id = "${alicloud_vpc.foo.id}"
}

resource "alicloud_security_group_rule" "ingress" {
	type = "ingress"
	ip_protocol = "tcp"
	nic_type = "intranet"
	policy = "accept"
	port_range = "22/22"
	priority = 1
	security_group_id = "${alicloud_security_group.tf_test_foo.id}"
	cidr_ip = "0.0.0.0/0"
}

resource "alicloud_instance" "foo" {
	# cn-beijing
	security_groups = ["${alicloud_security_group.tf_test_foo.id}"]

	vswitch_id = "${alicloud_vswitch.foo.id}"
	allocate_public_ip = true

	# series III
	instance_charge_type = "PostPaid"
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	internet_charge_type = "PayByTraffic"
	internet_max_bandwidth_out = 5

	system_disk_category = "cloud_efficiency"
	image_id = "${data.alicloud_images.default.images.0.id}"
	instance_name = "${var.name}"
}

data "alicloud_route_entries" "foo" {
  type = "System"
  route_table_id = "${alicloud_route_entry.foo.route_table_id}"
}
`

const testAccCheckAlicloudRouteEntriesDataSourceCidrBlock = `
data "alicloud_zones" "default" {
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
	default = "tf-testAcc-for-route-entries-datasource"
}
resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "10.1.1.0/24"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_route_entry" "foo" {
	route_table_id = "${alicloud_vpc.foo.route_table_id}"
	destination_cidrblock = "172.11.1.1/32"
	nexthop_type = "Instance"
	nexthop_id = "${alicloud_instance.foo.id}"
}

resource "alicloud_security_group" "tf_test_foo" {
	name = "${var.name}"
	description = "foo"
	vpc_id = "${alicloud_vpc.foo.id}"
}

resource "alicloud_security_group_rule" "ingress" {
	type = "ingress"
	ip_protocol = "tcp"
	nic_type = "intranet"
	policy = "accept"
	port_range = "22/22"
	priority = 1
	security_group_id = "${alicloud_security_group.tf_test_foo.id}"
	cidr_ip = "0.0.0.0/0"
}

resource "alicloud_instance" "foo" {
	# cn-beijing
	security_groups = ["${alicloud_security_group.tf_test_foo.id}"]

	vswitch_id = "${alicloud_vswitch.foo.id}"
	allocate_public_ip = true

	# series III
	instance_charge_type = "PostPaid"
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	internet_charge_type = "PayByTraffic"
	internet_max_bandwidth_out = 5

	system_disk_category = "cloud_efficiency"
	image_id = "${data.alicloud_images.default.images.0.id}"
	instance_name = "${var.name}"
}

data "alicloud_route_entries" "foo" {
  cidr_block = "${alicloud_route_entry.foo.destination_cidrblock}"
  route_table_id = "${alicloud_route_entry.foo.route_table_id}"
}
`

const testAccCheckAlicloudRouteEntriesDataSourceCidrBlock_mismatch = `
data "alicloud_zones" "default" {
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
	default = "tf-testAcc-for-route-entries-datasource"
}
resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "10.1.1.0/24"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_route_entry" "foo" {
	route_table_id = "${alicloud_vpc.foo.route_table_id}"
	destination_cidrblock = "172.11.1.1/32"
	nexthop_type = "Instance"
	nexthop_id = "${alicloud_instance.foo.id}"
}

resource "alicloud_security_group" "tf_test_foo" {
	name = "${var.name}"
	description = "foo"
	vpc_id = "${alicloud_vpc.foo.id}"
}

resource "alicloud_security_group_rule" "ingress" {
	type = "ingress"
	ip_protocol = "tcp"
	nic_type = "intranet"
	policy = "accept"
	port_range = "22/22"
	priority = 1
	security_group_id = "${alicloud_security_group.tf_test_foo.id}"
	cidr_ip = "0.0.0.0/0"
}

resource "alicloud_instance" "foo" {
	# cn-beijing
	security_groups = ["${alicloud_security_group.tf_test_foo.id}"]

	vswitch_id = "${alicloud_vswitch.foo.id}"
	allocate_public_ip = true

	# series III
	instance_charge_type = "PostPaid"
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	internet_charge_type = "PayByTraffic"
	internet_max_bandwidth_out = 5

	system_disk_category = "cloud_efficiency"
	image_id = "${data.alicloud_images.default.images.0.id}"
	instance_name = "${var.name}"
}

data "alicloud_route_entries" "foo" {
  cidr_block = "${alicloud_route_entry.foo.destination_cidrblock}-fake"
  route_table_id = "${alicloud_route_entry.foo.route_table_id}"
}
`

const testAccCheckAlicloudRouteEntriesDataSourceInstanceId = `
data "alicloud_zones" "default" {
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
	default = "tf-testAcc-for-route-entries-datasource"
}
resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "10.1.1.0/24"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_route_entry" "foo" {
	route_table_id = "${alicloud_vpc.foo.route_table_id}"
	destination_cidrblock = "172.11.1.1/32"
	nexthop_type = "Instance"
	nexthop_id = "${alicloud_instance.foo.id}"
}

resource "alicloud_security_group" "tf_test_foo" {
	name = "${var.name}"
	description = "foo"
	vpc_id = "${alicloud_vpc.foo.id}"
}

resource "alicloud_security_group_rule" "ingress" {
	type = "ingress"
	ip_protocol = "tcp"
	nic_type = "intranet"
	policy = "accept"
	port_range = "22/22"
	priority = 1
	security_group_id = "${alicloud_security_group.tf_test_foo.id}"
	cidr_ip = "0.0.0.0/0"
}

resource "alicloud_instance" "foo" {
	# cn-beijing
	security_groups = ["${alicloud_security_group.tf_test_foo.id}"]

	vswitch_id = "${alicloud_vswitch.foo.id}"
	allocate_public_ip = true

	# series III
	instance_charge_type = "PostPaid"
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	internet_charge_type = "PayByTraffic"
	internet_max_bandwidth_out = 5

	system_disk_category = "cloud_efficiency"
	image_id = "${data.alicloud_images.default.images.0.id}"
	instance_name = "${var.name}"
}

data "alicloud_route_entries" "foo" {
  instance_id = "${alicloud_instance.foo.id}"
  route_table_id = "${alicloud_route_entry.foo.route_table_id}"
}
`

const testAccCheckAlicloudRouteEntriesDataSourceInstanceId_mismatch = `
data "alicloud_zones" "default" {
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
	default = "tf-testAcc-for-route-entries-datasource"
}
resource "alicloud_vpc" "foo" {
	name = "${var.name}"
	cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "foo" {
	vpc_id = "${alicloud_vpc.foo.id}"
	cidr_block = "10.1.1.0/24"
	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alicloud_route_entry" "foo" {
	route_table_id = "${alicloud_vpc.foo.route_table_id}"
	destination_cidrblock = "172.11.1.1/32"
	nexthop_type = "Instance"
	nexthop_id = "${alicloud_instance.foo.id}"
}

resource "alicloud_security_group" "tf_test_foo" {
	name = "${var.name}"
	description = "foo"
	vpc_id = "${alicloud_vpc.foo.id}"
}

resource "alicloud_security_group_rule" "ingress" {
	type = "ingress"
	ip_protocol = "tcp"
	nic_type = "intranet"
	policy = "accept"
	port_range = "22/22"
	priority = 1
	security_group_id = "${alicloud_security_group.tf_test_foo.id}"
	cidr_ip = "0.0.0.0/0"
}

resource "alicloud_instance" "foo" {
	# cn-beijing
	security_groups = ["${alicloud_security_group.tf_test_foo.id}"]

	vswitch_id = "${alicloud_vswitch.foo.id}"
	allocate_public_ip = true

	# series III
	instance_charge_type = "PostPaid"
	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	internet_charge_type = "PayByTraffic"
	internet_max_bandwidth_out = 5

	system_disk_category = "cloud_efficiency"
	image_id = "${data.alicloud_images.default.images.0.id}"
	instance_name = "${var.name}"
}

data "alicloud_route_entries" "foo" {
  instance_id = "${alicloud_instance.foo.id}-fake"
  route_table_id = "${alicloud_route_entry.foo.route_table_id}"
}
`
