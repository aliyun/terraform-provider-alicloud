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
  disk_type   = "System"
  instance_id = length(data.alicloud_simple_application_server_instances.example.ids) > 0 ? data.alicloud_simple_application_server_instances.example.ids.0 : alicloud_simple_application_server_instance.example.0.id
}

resource "alicloud_simple_application_server_snapshot" "example" {
  disk_id       = data.alicloud_simple_application_server_disks.example.ids.0
  snapshot_name = "example_value"
}

resource "alicloud_simple_application_server_custom_image" "example" {
  custom_image_name  = "example_value"
  instance_id        = data.alicloud_simple_application_server_disks.example.disks.0.instance_id
  system_snapshot_id = alicloud_simple_application_server_snapshot.example.id
  status             = "Share"
  description        = "example_value"
}

