resource "alicloud_vpc" "main" {
  name = "${var.long_name}"
  cidr_block = "${var.vpc_cidr}"
}

resource "alicloud_vswitch" "main" {
  vpc_id = "${alicloud_vpc.main.id}"
  count = "${length(split(",", var.availability_zones))}"
  cidr_block = "${lookup(var.cidr_blocks, "az${count.index}")}"
  availability_zone = "${element(split(",", var.availability_zones), count.index)}"
  depends_on = [
    "alicloud_vpc.main"]
}

resource "alicloud_slb" "instance" {
  name = "${var.name}"
  vswitch_id = "${alicloud_vswitch.main.id}"
  internet_charge_type = "${var.internet_charge_type}"
}

resource "alicloud_slb_listener" "listener" {
  load_balancer_id = "${alicloud_slb.instance.id}"
  backend_port = "2111"
  frontend_port = "21"
  protocol = "tcp"
  bandwidth = "5"
}

resource "alicloud_router_interface" "interface" {
  opposite_region = "cn-beijing"
  router_type = "VRouter"
  router_id = "${alicloud_vpc.main.router_id}"
  role = "InitiatingSide"
  specification = "Large.2"
  name = "test1"
  description = "test1"
}