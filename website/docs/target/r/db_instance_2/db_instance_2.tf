variable "name" {
  default = "tf-testaccdbinstance"
}

data "alicloud_zones" "example" {
  available_resource_creation = "Rds"
}

resource "alicloud_vpc" "example" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "example" {
  count        = 2
  vpc_id       = alicloud_vpc.example.id
  cidr_block   = format("172.16.%d.0/24", count.index + 1)
  zone_id      = data.alicloud_zones.example.zones[count.index].id
  vswitch_name = format("vswich_%d", var.name, count.index)
}

resource "alicloud_db_instance" "example" {
  engine               = "MySQL"
  engine_version       = "5.6"
  instance_storage     = "30"
  instance_type        = "rds.mysql.t1.small"
  instance_charge_type = "Postpaid"
  instance_name        = var.name
  zone_id              = data.alicloud_zones.example.zones.0.id
  zone_id_slave_a      = data.alicloud_zones.example.zones.1.id
  vswitch_id           = join(",", alicloud_vswitch.example.*.id)
  monitoring_period    = "60"
}

