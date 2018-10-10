data "alicloud_images" main {
  most_recent = "${var.most_recent}"
  owners      = "${var.image_owners}"
  name_regex  = "${var.name_regex}"
}

data "alicloud_zones" main {
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" main {
  availability_zone = "${data.alicloud_zones.main.zones.0.id}"
  cpu_core_count    = "${var.cpu_core_count}"
  memory_size       = "${var.memory_size}"
}

resource "alicloud_vpc" "main" {
  name       = "${var.vpc_name}"
  cidr_block = "${var.vpc_cidr}"
}

resource "alicloud_vswitch" "main" {
  vpc_id            = "${alicloud_vpc.main.id}"
  cidr_block        = "${var.vswitch_cidr}"
  availability_zone = "${data.alicloud_zones.main.zones.0.id}"

  depends_on = [
    "alicloud_vpc.main",
  ]
}

resource "alicloud_container_cluster" "cs_vpc" {
  password      = "${var.password}"
  instance_type = "${data.alicloud_instance_types.main.instance_types.0.id}"
  name          = "${var.cluster_name}"
  node_number   = "${var.node_number}"
  disk_category = "${var.disk_category}"
  disk_size     = "${var.disk_size}"
  cidr_block    = "${var.cidr_block}"
  image_id      = "${data.alicloud_images.main.images.0.id}"
  vswitch_id    = "${alicloud_vswitch.main.id}"
}
