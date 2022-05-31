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
  engine               = "MySQL"
  engine_version       = "5.6"
  instance_type        = "rds.mysql.s2.large"
  instance_storage     = "30"
  instance_charge_type = "Postpaid"
  instance_name        = var.name
  vswitch_id           = alicloud_vswitch.example.id
  monitoring_period    = "60"
}

resource "alicloud_rds_clone_db_instance" "example" {
  source_db_instance_id    = alicloud_db_instance.example.id
  db_instance_storage_type = "local_ssd"
  payment_type             = "PayAsYouGo"
  restore_time             = "2021-11-24T11:25:00Z"
  db_instance_storage      = "30"
}
