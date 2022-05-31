variable "name" {
  default = "tf-testaccdbinstance"
}

variable "creation" {
  default = "Rds"
}

data "alicloud_zones" "example" {
  available_resource_creation = var.creation
}

resource "alicloud_vpc" "example" {
  name       = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "example" {
  vpc_id     = alicloud_vpc.example.id
  cidr_block = "172.16.0.0/24"
  zone_id    = data.alicloud_zones.example.zones[0].id
  name       = var.name
}

resource "alicloud_db_instance" "example" {
  engine               = "PostgreSQL"
  engine_version       = "12.0"
  instance_type        = "pg.n2.small.2c"
  instance_storage     = "20"
  instance_charge_type = "Postpaid"
  instance_name        = var.name
  vswitch_id           = alicloud_vswitch.example.id
}

resource "alicloud_rds_upgrade_db_instance" "example" {
  source_db_instance_id    = alicloud_db_instance.example.id
  target_major_version     = "13.0"
  db_instance_class        = "pg.n2.small.2c"
  db_instance_storage      = "20"
  instance_network_type    = "VPC"
  db_instance_storage_type = "cloud_ssd"
  collect_stat_mode        = "After"
  switch_over              = "false"
  payment_type             = "PayAsYouGo"
  db_instance_description  = var.name
  vswitch_id               = alicloud_vswitch.example.id
}
