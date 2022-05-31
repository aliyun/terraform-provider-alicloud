data "alicloud_resource_manager_resource_groups" "default" {
  name_regex = "default"
}
data "alicloud_zones" default {
  available_resource_creation = "Instance"
  available_disk_category     = "cloud_essd"
}

data "alicloud_instance_types" "default" {
  availability_zone    = data.alicloud_zones.default.zones.0.id
  cpu_core_count       = 2
  memory_size          = 4
  system_disk_category = "cloud_essd"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  name        = "${var.name}"
  description = "New security group"
  vpc_id      = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_disk" "default" {
  count     = 2
  disk_name = "${var.name}"
  zone_id   = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  category  = "cloud_essd"
  size      = "20"
}

data "alicloud_images" "default" {
  owners = "system"
}

resource "alicloud_instance" "default" {
  availability_zone = data.alicloud_zones.default.zones.0.id
  instance_name     = "${var.name}"
  host_name         = "tf-testAcc"
  image_id          = data.alicloud_images.default.images.0.id
  instance_type     = data.alicloud_instance_types.default.instance_types.0.id
  security_groups   = [alicloud_security_group.default.id]
  vswitch_id        = data.alicloud_vswitches.default.ids.0
}

resource "alicloud_disk_attachment" "default" {
  count       = 2
  disk_id     = element(alicloud_disk.default.*.id, count.index)
  instance_id = alicloud_instance.default.id
}

resource "alicloud_ecs_snapshot_group" "example" {
  description                   = "example_value"
  disk_id                       = [alicloud_disk_attachment.default.0.disk_id, alicloud_disk_attachment.default.1.disk_id]
  snapshot_group_name           = "example_value"
  resource_group_id             = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  instance_id                   = alicloud_disk_attachment.default.0.instance_id
  instant_access                = true
  instant_access_retention_days = 1
  tags = {
    Created = "TF"
    For     = "Acceptance-test"
  }
}
