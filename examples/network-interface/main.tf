resource "alicloud_vpc" "vpc" {
    name = "tf-testAcc-vpc"
    cidr_block = "${var.vpc_cidr}"
}

data "alicloud_zones" "default" {
    "available_resource_creation"= "VSwitch"
}

resource "alicloud_vswitch" "vswitch" {
    name = "tf-testAcc-vswitch"
    cidr_block = "${var.vswitch_cidr}"
    availability_zone = "${data.alicloud_zones.default.zones.0.id}"
    vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_security_group" "sg" {
    name = "tf-testAcc-sg"
    vpc_id = "${alicloud_vpc.vpc.id}"
}

resource "alicloud_network_interface" "eni" {
    name = "tf-testAcc-eni"
    vswitch_id = "${alicloud_vswitch.vswitch.id}"
    security_groups = [ "${alicloud_security_group.sg.id}" ]
}
