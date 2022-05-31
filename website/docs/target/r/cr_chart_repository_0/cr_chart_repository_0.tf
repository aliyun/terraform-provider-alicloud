resource "alicloud_cr_ee_instance" "default" {
  payment_type  = "Subscription"
  period        = 1
  instance_type = "Advanced"
  instance_name = "name"
}

resource "alicloud_cr_chart_namespace" "default" {
  instance_id    = alicloud_cr_ee_instance.default.id
  namespace_name = "name"
}

resource "alicloud_cr_chart_repository" "default" {
  repo_namespace_name = alicloud_cr_chart_namespace.default.namespace_name
  instance_id         = local.instance
  repo_name           = "repo_name"
}

