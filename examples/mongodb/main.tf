data "alicloud_mongodb_instances" "mongo" {
  output_file = "out.dat"
}

// VPC Resource for Module
resource "alicloud_vpc" "vpc" {
  count      = "${var.vpc_id == "" ? 1 : 0}"
  name       = "${var.vpc_name}"
  cidr_block = "${var.vpc_cidr}"
}

// VSwitch Resource for Module
resource "alicloud_vswitch" "vswitch" {
  count             = "${var.vswitch_id == "" ? 1 : 0}"
  availability_zone = "eu-central-1a"
  name              = "${var.vswitch_name}"
  cidr_block        = "${var.vswitch_cidr}"
  vpc_id            = "${var.vpc_id == "" ? alicloud_vpc.vpc.id : var.vpc_id}"
}

resource "alicloud_mongodb_instance" "mymongo" {
  instance_class   = "dds.mongo.mid"
  instance_storage = "10"
  engine_version   = "3.4"
  description      = "my-description"
  security_ips     = ["127.0.0.1", "2.2.2.2"]
  vswitch_id       = "${var.vswitch_id == "" ? alicloud_vswitch.vswitch.id : var.vswitch_id}"
}

resource "alicloud_mongodb_backup_policy" "mongodb_backup" {
  instance_id             = "${alicloud_mongodb_instance.mymongo.id}"
  preferred_backup_time   = "03:00Z-04:00Z"
  preferred_backup_period = ["Monday", "Wednesday", "Friday"]
}