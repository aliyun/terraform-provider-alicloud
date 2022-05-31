variable "name" {
  default = "example_name"
}
data "alicloud_cr_ee_instances" "default" {}
resource "alicloud_cr_ee_instance" "default" {
  count          = length(data.alicloud_cr_ee_instances.default.ids) > 0 ? 0 : 1
  payment_type   = "Subscription"
  period         = 1
  renewal_status = "ManualRenewal"
  instance_type  = "Advanced"
  instance_name  = var.name
}
data "alicloud_cr_endpoint_acl_service" "default" {
  endpoint_type = "internet"
  enable        = true
  instance_id   = length(data.alicloud_cr_ee_instances.default.ids) > 0 ? data.alicloud_cr_ee_instances.default.ids[0] : concat(alicloud_cr_ee_instance.default.*.id, [""])[0]
  module_name   = "Registry"
}
resource "alicloud_cr_endpoint_acl_policy" "default" {
  instance_id   = length(data.alicloud_cr_ee_instances.default.ids) > 0 ? data.alicloud_cr_ee_instances.default.ids[0] : concat(alicloud_cr_ee_instance.default.*.id, [""])[0]
  entry         = "192.168.1.0/24"
  description   = var.name
  module_name   = "Registry"
  endpoint_type = "internet"
  depends_on    = [data.alicloud_cr_endpoint_acl_service.default]
}
