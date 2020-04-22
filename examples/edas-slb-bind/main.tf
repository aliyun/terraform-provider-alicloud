resource "alicloud_edas_slb_bind" "default" {
  app_id = "${var.app_id}"
  slb_id = "${var.slb_id}"
  slb_ip = "${var.slb_ip}"
  type = "${var.type}"
  listener_port = "${var.listener_port}"
  vserver_group_id = "${var.vserver_group_id}"
}