resource "alicloud_eip" "foo" {
  count                = "${length(var.internet_charge_type)}"
  period               = "${var.period}"
  bandwidth            = "${var.bandwidth}"
  instance_charge_type = "${var.instance_charge_type}"
  internet_charge_type = "${element(var.internet_charge_type, count.index)}"

}

resource "alicloud_eip_association" "foo" {
  allocation_id = "${alicloud_eip.foo.0.id}"
  instance_id   = "${alicloud_instance.default.id}"
  instance_type = "${var.associate_instance_type}"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

data "alicloud_images" "default" {
  name_regex  = "^ubuntu_14\\w{1,5}[64]{1}.*"
  most_recent = true
  owners      = "system"
}

resource "alicloud_vpc" "default" {
  name       = "alicloud"
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "default" {
  vpc_id            = "${alicloud_vpc.default.id}"
  cidr_block        = "10.1.1.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_security_group" "default" {
  description = "default"
  vpc_id      = "${alicloud_vpc.default.id}"
}

resource "alicloud_instance" "default" {
  security_groups = ["${alicloud_security_group.default.id}"]

  vswitch_id = "${alicloud_vswitch.default.id}"

  instance_charge_type = "PostPaid"
  instance_type        = "${data.alicloud_instance_types.default.instance_types.0.id}"
  internet_charge_type = "PayByTraffic"
  image_id             = "${data.alicloud_images.default.images.0.id}"
}

