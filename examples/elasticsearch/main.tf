resource "alicloud_elasticsearch_instance" "instance" {
  instance_charge_type = "${var.instance_charge_type}"
  data_node_amount     = "${var.data_node_amount}"
  data_node_spec       = "${var.data_node_spec}"
  data_node_disk_size  = "${var.data_node_disk_size}"
  data_node_disk_type  = "${var.data_node_disk_type}"
  vswitch_id           = "${var.vswitch_id}"
  password             = "${var.password}"
  version              = "${var.version}"
  description          = "${var.description}"
}

data "alicloud_elasticsearch_instances" "myeslist" {
  description_regex = "tf-testAccES.*"
}
