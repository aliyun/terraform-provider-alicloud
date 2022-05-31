variable "name" {
  default = "autoscaler"
}

data "alicloud_vpcs" "default" {}

data "alicloud_images" "default" {
  owners      = "system"
  name_regex  = "^centos_7"
  most_recent = true
}

# If your account no running clusters, you need to create a new one
data "alicloud_cs_managed_kubernetes_clusters" "default" {}

data "alicloud_instance_types" "default" {
  cpu_core_count = 2
  memory_size    = 4
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = data.alicloud_vpcs.default.vpcs.0.id
}

resource "alicloud_ess_scaling_group" "default" {
  scaling_group_name = var.name

  min_size         = var.min_size
  max_size         = var.max_size
  vswitch_ids      = [data.alicloud_vpcs.default.vpcs.0.vswitch_ids.0]
  removal_policies = ["OldestInstance", "NewestInstance"]
}

resource "alicloud_ess_scaling_configuration" "default" {
  image_id             = data.alicloud_images.default.images.0.id
  security_group_id    = alicloud_security_group.default.id
  scaling_group_id     = alicloud_ess_scaling_group.default.id
  instance_type        = data.alicloud_instance_types.default.instance_types.0.id
  internet_charge_type = "PayByTraffic"
  force_delete         = true
  enable               = true
  active               = true

  # ... ignore the change about tags and user_data
  lifecycle {
    ignore_changes = [tags, user_data]
  }

}

resource "alicloud_cs_kubernetes_autoscaler" "default" {
  cluster_id = data.alicloud_cs_managed_kubernetes_clusters.default.clusters.0.id
  nodepools {
    id     = alicloud_ess_scaling_group.default.id
    labels = "a=b"
  }

  utilization             = var.utilization
  cool_down_duration      = var.cool_down_duration
  defer_scale_in_duration = var.defer_scale_in_duration

  depends_on = [alicloud_ess_scaling_group.defalut, alicloud_ess_scaling_configuration.default]
}
