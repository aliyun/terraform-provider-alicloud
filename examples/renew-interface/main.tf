resource "alicloud_renew_interface" "foo" {
  instance_id = var.instance_id
  instance_type = var.instance_type
  pricing_cycle = var.pricing_cycle
  duration = var.duration
  renewal_status = var.renewal_status
}

data "alicloud_auto_renew_instances" "foo" {
  instance_type = var.instance_type
}