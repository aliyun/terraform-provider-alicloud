data "alicloud_instance_types" "instance_type" {
  instance_type_family = "ecs.n4"
  cpu_core_count = "1"
  memory_size = "2"
}
resource "alicloud_vpc" "main" {
  name = "vpc-${var.short_name}"
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "main" {
  vpc_id = "${alicloud_vpc.main.id}"
  cidr_block = "10.1.1.0/24"
  availability_zone = "${var.availability_zones}"
  depends_on = [
    "alicloud_vpc.main"]
}

resource "alicloud_security_group" "group" {
  name = "group-${var.short_name}"
  description = "New security group"
  vpc_id = "${alicloud_vpc.main.id}"
}

resource "alicloud_security_group_rule" "allow_http_80" {
  type = "ingress"
  ip_protocol = "tcp"
  nic_type = "${var.nic_type}"
  policy = "accept"
  port_range = "80/80"
  priority = 1
  security_group_id = "${alicloud_security_group.group.id}"
  cidr_ip = "0.0.0.0/0"
}

resource "alicloud_security_group_rule" "allow_https_22" {
  type = "ingress"
  ip_protocol = "tcp"
  nic_type = "${var.nic_type}"
  policy = "accept"
  port_range = "22/22"
  priority = 1
  security_group_id = "${alicloud_security_group.group.id}"
  cidr_ip = "0.0.0.0/0"
}

resource "alicloud_instance" "instance" {
  instance_name = "${var.short_name}-${var.role}-${format(var.count_format, count.index+1)}"
  host_name = "${var.short_name}-${var.role}-${format(var.count_format, count.index+1)}"
  image_id = "${var.image_id}"
  instance_type = "${data.alicloud_instance_types.instance_type.instance_types.0.id}"
  count = "${var.count}"
  availability_zone = "${var.availability_zones}"
  security_groups = ["${alicloud_security_group.group.*.id}"]

  internet_charge_type = "${var.internet_charge_type}"
  internet_max_bandwidth_out = "${var.internet_max_bandwidth_out}"

  password = "${var.ecs_password}"

  allocate_public_ip = "${var.allocate_public_ip}"

  instance_charge_type = "PostPaid"
  system_disk_category = "${var.disk_category}"

  vswitch_id = "${alicloud_vswitch.main.id}"

  tags {
    role = "${var.role}"
    dc = "${var.datacenter}"
  }

}

resource "alicloud_disk" "disk" {
  availability_zone = "${alicloud_instance.instance.0.availability_zone}"
  category = "${var.disk_category}"
  size = "${var.disk_size}"
  count = "${var.disk_count}"
}

resource "alicloud_disk_attachment" "instance-attachment" {
  count = "${var.disk_count}"
  disk_id = "${element(alicloud_disk.disk.*.id, count.index)}"
  instance_id = "${element(alicloud_instance.instance.*.id, count.index%var.count)}"
}

resource "alicloud_key_pair" "key_pair" {
  key_name = "${var.key_name}"
  key_file = "${var.private_key_file}"
}

resource "alicloud_key_pair_attachment" "key_pair_attachment" {
  key_name = "${alicloud_key_pair.key_pair.id}"
  instance_ids = ["${alicloud_instance.instance.*.id}"]
}