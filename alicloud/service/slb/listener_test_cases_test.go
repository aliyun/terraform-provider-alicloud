package slb_test

// Shared listener test-config constants for the split slb test package.
// The flat-layout originals live in the root alicloud package
// (alicloud/service_alicloud_common_test.go) and cannot be referenced
// unqualified from a service test package. They are pure HCL fixtures with no
// Go symbol references, so the values are copied verbatim to keep behaviour
// unchanged; sibling service packages keep the same local-copy convention.

const SlbListenerCommonTestCase = `
variable "ip_version" {
  default = "ipv4"
}
resource "alicloud_slb_load_balancer" "default" {
  load_balancer_name = "${var.name}"
  internet_charge_type = "PayByTraffic"
  address_type = "internet"
  load_balancer_spec = "slb.s1.small"
}
resource "alicloud_slb_acl" "default" {
  name = "${var.name}"
  ip_version = "${var.ip_version}"
  entry_list {
      entry="10.10.10.0/24"
      comment="first"
  }
  entry_list {
      entry="168.10.10.0/24"
      comment="second"
  }
}
`
const SlbListenerVserverCommonTestCase = `
data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_[0-9]+_[0-9]+_x64*"
  most_recent = true
  owners      = "system"
}

data "alicloud_instance_types" "default" {
  availability_zone    = data.alicloud_zones.default.zones.0.id
  image_id             = data.alicloud_images.default.images.0.id
  instance_type_family = "ecs.g6"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "192.168.192.0/24"
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_instance" "default" {
  count                      = "2"
  image_id                   = data.alicloud_images.default.images.0.id
  instance_type              = data.alicloud_instance_types.default.instance_types.0.id
  security_groups            = alicloud_security_group.default.*.id
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone          = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  instance_charge_type       = "PostPaid"
  system_disk_category       = "cloud_efficiency"
  vswitch_id                 = alicloud_vswitch.default.id
  instance_name              = var.name
}

resource "alicloud_slb_load_balancer" "default" {
  load_balancer_name = var.name
  vswitch_id         = alicloud_vswitch.default.id
  load_balancer_spec = "slb.s1.small"
}

resource "alicloud_slb_server_group" "default" {
  load_balancer_id = alicloud_slb_load_balancer.default.id
  name             = var.name
}

resource "alicloud_slb_master_slave_server_group" "default" {
  load_balancer_id = alicloud_slb_load_balancer.default.id
  name             = var.name
  servers {
    server_id   = alicloud_instance.default.0.id
    port        = 80
    weight      = 100
    server_type = "Master"
  }
  servers {
    server_id   = alicloud_instance.default.1.id
    port        = 80
    weight      = 100
    server_type = "Slave"
  }
}
`
