data "alicloud_simple_application_server_images" "default" {}
data "alicloud_simple_application_server_plans" "default" {}

resource "alicloud_simple_application_server" "default" {
  payment_type   = "Subscription"
  plan_id        = data.alicloud_simple_application_server_plans.default.plans.0.id
  instance_name  = var.name
  image_id       = data.alicloud_simple_application_server_images.default.images.0.id
  period         = 1
  data_disk_size = 100
}

