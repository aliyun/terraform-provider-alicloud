// Provider specific configs
provider "alicloud" {
  access_key = "${var.alicloud_access_key}"
  secret_key = "${var.alicloud_secret_key}"
  region     = "${var.region}"
}

// Images data source for image_id
data "alicloud_images" "default" {
  most_recent = true
  owners      = "system"
  name_regex  = "^ubuntu_14.*_64"
}

// Instance_types data source for instance_type
data "alicloud_instance_types" "default" {
  instance_type_family = "ecs.n4"
  cpu_core_count       = 1
  memory_size          = 2
}

// Zones data source for availability_zone
data "alicloud_zones" "default" {
  available_disk_category = "${var.disk_category}"
  available_instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
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
  vpc_id            = "${var.vpc_id == "" ? join("", alicloud_vpc.vpc.*.id) : var.vpc_id}"
}

// Security Group Resource for Module
resource "alicloud_security_group" "group" {
  count = "${var.sg_id == "" ? 1 : 0}"

  name   = "${var.sg_name}"
  vpc_id = "${var.vpc_id == "" ? join("", alicloud_vpc.vpc.*.id) : var.vpc_id}"
}

// Security Group Resource for Module
resource "alicloud_security_group_rule" "rules" {
  count = "${length(var.ip_protocols)}"

  type              = "${count.index >= length(var.rule_directions) ? "ingress" : var.rule_directions[count.index]}"
  ip_protocol       = "${var.ip_protocols[count.index]}"
  nic_type          = "intranet"
  policy            = "${count.index >= length(var.policies) ? "accept" : var.policies[count.index]}"
  port_range        = "${count.index >= length(var.ip_protocols) ? "-1/-1" : var.port_ranges[count.index]}"
  priority          = "${count.index >= length(var.priorities) ? 1 : var.priorities[count.index]}"
  security_group_id = "${var.sg_id == "" ? join("", alicloud_security_group.group.*.id) : var.sg_id}"
  cidr_ip           = "${length(var.cidr_ips) <= 0 || count.index >= length(var.cidr_ips) ? "0.0.0.0/0" : element(concat(var.cidr_ips, list("0.0.0.0/0")), count.index)}"
}

// ECS Instance Resource for Module
resource "alicloud_instance" "instances" {
  count = "${var.number_of_instances}"

  image_id        = "${var.image_id == "" ? data.alicloud_images.default.images.0.id : var.image_id }"
  instance_type   = "${var.instance_type == "" ? data.alicloud_instance_types.default.instance_types.0.id : var.instance_type}"
  security_groups = ["${var.sg_id == "" ? join("", alicloud_security_group.group.*.id) : var.sg_id}"]

  instance_name = "${var.number_of_instances < 2 ? var.instance_name : format("%s-%s", var.instance_name, format(var.number_format, count.index+1))}"
  host_name     = "${var.number_of_instances < 2 ? var.host_name : format("%s-%s", var.host_name, format(var.number_format, count.index+1))}"

  internet_charge_type       = "${var.internet_charge_type}"
  internet_max_bandwidth_out = "${var.internet_max_bandwidth_out}"

  instance_charge_type = "${var.instance_charge_type}"
  system_disk_category = "${var.system_category}"
  system_disk_size     = "${var.system_size}"

  password = "${var.password}"

  vswitch_id = "${var.vswitch_id == "" ? join("", alicloud_vswitch.vswitch.*.id) : var.vswitch_id}"

  period = "${var.period}"

  tags {
    created_by   = "${lookup(var.instance_tags, "created_by")}"
    created_from = "${lookup(var.instance_tags, "created_from")}"
  }
}

// ECS Disk Resource for Module
resource "alicloud_disk" "disks" {
  count = "${var.number_of_disks}"

  availability_zone = "${var.availability_zone == "" ? data.alicloud_zones.default.zones.0.id : var.availability_zone}"
  name              = "${var.number_of_disks < 2 ? var.disk_name : format("%s-%s", var.disk_name, format(var.number_format, count.index+1))}"
  category          = "${var.disk_category}"
  size              = "${var.disk_size}"

  tags {
    created_by   = "${lookup(var.disk_tags, "created_by")}"
    created_from = "${lookup(var.disk_tags, "created_from")}"
  }
}

// Attach ECS disks to instances for Module
resource "alicloud_disk_attachment" "disk_attach" {
  count       = "${(var.number_of_instances > 0 && var.number_of_disks > 0) ? var.number_of_disks : 0}"
  disk_id     = "${element(alicloud_disk.disks.*.id, count.index)}"
  instance_id = "${element(alicloud_instance.instances.*.id, count.index%var.number_of_instances)}"
}

// Attach key pair to instances for Module
resource "alicloud_key_pair_attchment" "default" {
  count = "${var.number_of_instances > 0 && var.key_name != "" ? 1 : 0}"

  key_name     = "${var.key_name}"
  instance_ids = ["${alicloud_instance.instances.*.id}"]
}
