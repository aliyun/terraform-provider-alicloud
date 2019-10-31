data "alicloud_zones" "default" {
  available_disk_category = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alicloud_instance_types" "default" {
  availability_zone = data.alicloud_zones.default.zones[0].id
  eni_amount = 2
}

data "alicloud_images" "image" {
  name_regex = "^ubuntu_18.*_64"
  most_recent = true
  owners = "system"
}

variable "name" {
  default = "tf-testAccCloudConnectNetworkGrant"
}

data "alicloud_account" "default"{
}

resource "alicloud_cen_instance" "default" {
  name = "${var.name}"
}

resource "alicloud_cloud_connect_network" "default" {
  name = "${var.name}"
  is_default = "true"
}

resource "alicloud_cloud_connect_network_grant" "default" {
  ccn_id = "${alicloud_cloud_connect_network.default.id}"
  cen_id = "${alicloud_cen_instance.default.id}"
  cen_uid = "${data.alicloud_account.default.id}"
}