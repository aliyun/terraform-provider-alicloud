// Provider specific configs
provider "alicloud" {
  region = "${var.region}"
}

// If there is not specifying vpc_id, the module will launch a new vpc
resource "alicloud_vpc" "init" {
  provider = "alicloud"
  count = "${var.init_vpc_id == "" ? 1 : 0}"
  name = "${var.vpc_name}"
  cidr_block = "${var.init_vpc_cidr}"
  description = "${var.vpc_description}"
}

// Provide an initiating side interface
resource "alicloud_router_interface" "init" {
  provider = "alicloud"
  opposite_region = "${var.opposite_region}"
  router_type = "VRouter"
  router_id = "${alicloud_vpc.init.router_id}"
  role = "InitiatingSide"
  specification = "${var.interface_spec}"
  name = "initiating-side"
  description = "initiating side router interface"
}

resource "alicloud_router_interface_connection" "init" {
  provider = "alicloud"
  interface_id = "${alicloud_router_interface.init.id}"
  opposite_interface_id = "${alicloud_router_interface.accept.id}"
  depends_on = ["alicloud_router_interface.accept"]
}


// Provider specific configs
provider "alicloud" {
  alias = "accept"
  region = "${var.opposite_region}"
}

// If there is not specifying vpc_id, the module will launch a new vpc
resource "alicloud_vpc" "accept" {
  provider = "alicloud.accept"
  count = "${var.accept_vpc_id == "" ? 1 : 0}"
  name = "${var.vpc_name}"
  cidr_block = "${var.accept_vpc_cidr}"
  description = "${var.vpc_description}"
}

// provide a accept side interface
resource "alicloud_router_interface" "accept" {
  provider = "alicloud.accept"
  opposite_region = "${var.region}"
  router_type = "VRouter"
  router_id = "${alicloud_vpc.accept.router_id}"
  role = "AcceptingSide"
  name = "accepting-side"
  description = "Accepting side router interface"
}

resource "alicloud_router_interface_connection" "accept" {
  provider = "alicloud.accept"
  interface_id = "${alicloud_router_interface.accept.id}"
  opposite_interface_id = "${alicloud_router_interface.init.id}"
}



