data "alicloud_simple_application_server_instances" "example" {}

data "alicloud_simple_application_server_images" "example" {}

data "alicloud_simple_application_server_plans" "example" {}

resource "alicloud_simple_application_server_instance" "example" {
  count         = length(data.alicloud_simple_application_server_instances.example.ids) > 0 ? 0 : 1
  payment_type  = "Subscription"
  plan_id       = data.alicloud_simple_application_server_plans.example.plans.0.id
  instance_name = "example_value"
  image_id      = data.alicloud_simple_application_server_images.example.images.0.id
  period        = 1
}

data "alicloud_simple_application_server_disks" "example" {
  instance_id = length(data.alicloud_simple_application_server_instances.example.ids) > 0 ? data.alicloud_simple_application_server_instances.example.ids.0 : alicloud_simple_application_server_instance.example.0.id
}

resource "alicloud_simple_application_server_snapshot" "example" {
  disk_id       = data.alicloud_simple_application_server_disks.example.ids.0
  snapshot_name = "example_value"
}

