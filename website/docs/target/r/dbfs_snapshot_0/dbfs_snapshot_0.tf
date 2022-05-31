data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

locals {
  zone_id = "cn-hangzhou-i"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = local.zone_id
}

resource "alicloud_security_group" "default" {
  name        = var.name
  description = "tf test"
  vpc_id      = data.alicloud_vpcs.default.ids.0
}

data "alicloud_images" "default" {
  owners      = "system"
  name_regex  = "^centos_8"
  most_recent = true
}

resource "alicloud_instance" "default" {
  image_id             = data.alicloud_images.default.images[0].id
  instance_name        = var.name
  instance_type        = "ecs.g7se.large"
  availability_zone    = local.zone_id
  vswitch_id           = data.alicloud_vswitches.default.ids[0]
  system_disk_category = "cloud_essd"
  security_groups = [
    alicloud_security_group.default.id
  ]
}

resource "alicloud_dbfs_instance" "default" {
  category          = "standard"
  zone_id           = alicloud_instance.default.availability_zone
  performance_level = "PL1"
  instance_name     = var.name
  size              = 100
}

resource "alicloud_dbfs_instance_attachment" "default" {
  ecs_id      = alicloud_instance.default.id
  instance_id = alicloud_dbfs_instance.default.id
}

resource "alicloud_dbfs_snapshot" "example" {
  depends_on     = [alicloud_dbfs_instance_attachment.default]
  instance_id    = data.alicloud_dbfs_instances.default.ids.0
  snapshot_name  = "example_value"
  description    = "example_value"
  retention_days = 30
}
