provider "alicloud" {
  alias = "bj"
  region = "cn-beijing"
}

provider "alicloud" {
  alias = "sh"
  region = "cn-shanghai"
}

data "alicloud_zones" "default" {
    provider = "alicloud.bj"
	available_disk_category = "cloud_efficiency"
	available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
    provider = "alicloud.bj"
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
	cpu_core_count = 1
	memory_size = 2
}

data "alicloud_images" "default" {
    provider = "alicloud.bj"
    name_regex = "^ubuntu_14.*_64"
	most_recent = true
	owners = "system"
}

resource "alicloud_vpc" "vpc_1" {
    provider = "alicloud.bj"
  	name = "${var.name}"
  	cidr_block = "172.16.0.0/12"
}

resource "alicloud_vpc" "vpc_2" {
    provider = "alicloud.sh"
  	name = "${var.name}"
  	cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
    provider = "alicloud.bj"
 	vpc_id = "${alicloud_vpc.vpc_1.id}"
 	cidr_block = "172.16.0.0/21"
 	availability_zone = "${data.alicloud_zones.default.zones.0.id}"
 	name = "${var.name}"
}

resource "alicloud_security_group" "default" {
    provider = "alicloud.bj"
	name = "${var.name}"
	description = "${var.description}"
	vpc_id = "${alicloud_vpc.vpc_1.id}"
}

resource "alicloud_instance" "default" {
    provider = "alicloud.bj"
	vswitch_id = "${alicloud_vswitch.default.id}"
	image_id = "${data.alicloud_images.default.images.0.id}"

	instance_type = "${data.alicloud_instance_types.default.instance_types.0.id}"
	system_disk_category = "cloud_efficiency"

	internet_charge_type = "PayByTraffic"
	internet_max_bandwidth_out = 5
	security_groups = ["${alicloud_security_group.default.id}"]
	instance_name = "${var.name}"
}

resource "alicloud_cen" "cen" {
	name = "${var.name}"
	description = "${var.description}"
}

resource "alicloud_cen_instance_attachment" "vpc_attach_1" {
    cen_id = "${alicloud_cen.cen.id}"
    child_instance_id = "${alicloud_vpc.vpc_1.id}"
    child_instance_type = "VPC"
    child_instance_region_id = "cn-beijing"
}

resource "alicloud_cen_instance_attachment" "vpc_attach_2" {
    cen_id = "${alicloud_cen.cen.id}"
    child_instance_id = "${alicloud_vpc.vpc_2.id}"
    child_instance_type = "VPC"
    child_instance_region_id = "cn-shanghai"
}

resource "alicloud_cen_bandwidthpackage" "bwp" {
    bandwidth = 20
    geographic_region_id = [
		"China",
		"China"]
}

resource "alicloud_cen_bandwidthpackage_attachment" "bwp_attach" {
    cen_id = "${alicloud_cen.cen.id}"
    cen_bandwidthpackage_id = "${alicloud_cen_bandwidthpackage.bwp.id}"
}

resource "alicloud_cen_bandwidthlimit" "bwp_limit" {
     cen_id = "${alicloud_cen.cen.id}"
     regions_id = [
        "cn-beijing",
        "cn-shanghai"]
     bandwidth_limit = 15
     depends_on = [
        "alicloud_cen_bandwidthpackage_attachment.bwp_attach",
        "alicloud_cen_instance_attachment.vpc_attach_1",
        "alicloud_cen_instance_attachment.vpc_attach_2"]
}