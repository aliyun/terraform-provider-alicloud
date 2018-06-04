// Provider specific configs
provider "alicloud" {
  access_key = "${var.alicloud_access_key}"
  secret_key = "${var.alicloud_secret_key}"
  region = "${var.region}"
}

data "alicloud_images" main {
  most_recent = "${var.most_recent}"
  owners = "${var.image_owners}"
  name_regex = "${var.image_name_regex}"
}

data "alicloud_zones" main {
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" main {
  availability_zone = "${data.alicloud_zones.main.zones.0.id}"
  cpu_core_count = "${var.cpu_core_count}"
  memory_size = "${var.memory_size}"
}

resource "alicloud_vpc" "main" {
  name = "${var.vpc_name}"
  cidr_block = "${var.vpc_cidr}"
}

resource "alicloud_vswitch" "main" {
  vpc_id = "${alicloud_vpc.main.id}"
  cidr_block = "${var.vswitch_cidr}"
  availability_zone = "${data.alicloud_zones.main.zones.0.id}"
}

resource "alicloud_cs_swarm" "cs_vpc" {
  password = "${var.password}"
  instance_type = "${data.alicloud_instance_types.main.instance_types.0.id}"
  name = "${var.cluster_name}"
  node_number = "${var.node_number}"
  disk_category = "${var.disk_category}"
  disk_size = "${var.disk_size}"
  cidr_block = "${var.cidr_block}"
  image_id = "${data.alicloud_images.main.images.0.id}"
  vswitch_id = "${alicloud_vswitch.main.id}"
}

resource "alicloud_cs_application" "wordpress" {
  cluster_name = "${alicloud_cs_swarm.cs_vpc.name}"
  name = "${var.app_name == "" ? var.resource_group_name : var.app_name}"
  version = "${var.app_version}"
  template = "${file("wordpress.yml")}"
  description = "terraform deploy consource"
  latest_image = "${var.latest_image}"
  blue_green = "${var.blue_green}"
  blue_green_confirm = "${var.confirm_blue_green}"
}

