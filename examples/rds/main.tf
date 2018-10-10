// Zones data source for availability_zone
data "alicloud_zones" "default" {
  available_resource_creation = "Rds"
}

// VPC Resource for Module
resource "alicloud_vpc" "vpc" {
  count = "${var.vpc_id == "" ? 1 : 0}"

  name       = "${var.vpc_name}"
  cidr_block = "${var.vpc_cidr}"
}

// VSwitch Resource for Module
resource "alicloud_vswitch" "vswitch" {
  count = "${var.vswitch_id == "" ? 1 : 0}"

  availability_zone = "${var.availability_zone == "" ? data.alicloud_zones.default.zones.0.id : var.availability_zone}"
  name              = "${var.vswitch_name}"
  cidr_block        = "${var.vswitch_cidr}"
  vpc_id            = "${var.vpc_id == "" ? alicloud_vpc.vpc.id : var.vpc_id}"
}

resource "alicloud_db_instance" "instance" {
  engine           = "${var.engine}"
  engine_version   = "${var.engine_version}"
  instance_type    = "${var.instance_class}"
  instance_storage = "${var.storage}"
  vswitch_id       = "${var.vswitch_id == "" ? alicloud_vswitch.vswitch.id : var.vswitch_id}"
}

resource "alicloud_db_account" "account" {
  count       = 2
  instance_id = "${alicloud_db_instance.instance.id}"
  name        = "tf_account_${count.index}"
  password    = "${var.password}"
}

resource "alicloud_db_backup_policy" "backup" {
  instance_id   = "${alicloud_db_instance.instance.id}"
  backup_period = ["Tuesday", "Wednesday"]
  backup_time   = "10:00Z-11:00Z"
}

resource "alicloud_db_connection" "connection" {
  instance_id       = "${alicloud_db_instance.instance.id}"
  connection_prefix = "tf-example"
}

resource "alicloud_db_database" "db" {
  count       = 2
  instance_id = "${alicloud_db_instance.instance.id}"
  name        = "${var.database_name}_${count.index}"
}

resource "alicloud_db_account_privilege" "privilege" {
  count        = 2
  instance_id  = "${alicloud_db_instance.instance.id}"
  account_name = "${element(alicloud_db_account.account.*.name, count.index)}"
  db_names     = ["${alicloud_db_database.db.*.name}"]
}
