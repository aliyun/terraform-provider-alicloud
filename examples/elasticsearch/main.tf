resource "alicloud_elasticsearch" "instance" {
  instance_charge_type = "${var.instance_charge_type}"
  data_node_amount     = "${var.data_node_amount}"
  data_node_spec       = "${var.data_node_spec}"
  data_node_disk       = "${var.data_node_disk}"
  data_node_disk_type  = "${var.data_node_disk_type}"
  vswitch_id           = "${var.vswitch_id}"
  es_admin_password    = "${var.es_admin_password}"
  es_version           = "${var.es_version}"
  description          = "${var.description}"
}

data "alicloud_elasticsearch" "instances" {
}
