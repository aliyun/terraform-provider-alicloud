data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
  availability_zone = data.alicloud_zones.default.zones[0].id
  eni_amount        = 2
}

data "alicloud_images" "image" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

variable "name" {
  default = "tf-testAccSlbMasterSlaveServerGroupVpc"
}

variable "number" {
  default = 2
}

resource "alicloud_vpc" "main" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "main" {
  vpc_id       = alicloud_vpc.main.id
  cidr_block   = "172.16.0.0/16"
  zone_id      = data.alicloud_zones.default.zones[0].id
  vswitch_name = var.name
}

resource "alicloud_security_group" "group" {
  name   = var.name
  vpc_id = alicloud_vpc.main.id
}

resource "alicloud_instance" "instance" {
  count                      = 2
  image_id                   = data.alicloud_images.image.images[0].id
  instance_type              = data.alicloud_instance_types.default.instance_types[0].id
  instance_name              = var.name
  security_groups            = [alicloud_security_group.group.id]
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  availability_zone          = data.alicloud_zones.default.zones[0].id
  instance_charge_type       = "PostPaid"
  system_disk_category       = "cloud_efficiency"
  vswitch_id                 = alicloud_vswitch.main.id
}

resource "alicloud_slb_load_balancer" "instance" {
  load_balancer_name              = var.name
  vswitch_id        = alicloud_vswitch.main.id
  load_balancer_spec     = "slb.s2.small"
  delete_protection = "on"
}

resource "alicloud_slb_master_slave_server_group" "this" {
  load_balancer_id = alicloud_slb_load_balancer.instance.id
  name             = var.name
  servers {
    server_id   = alicloud_instance.instance.0.id
    port        = 100
    weight      = 100
    server_type = "Master"
  }

  servers {
    server_id   = alicloud_instance.instance.1.id
    port        = 100
    weight      = 100
    server_type = "Slave"
  }
}