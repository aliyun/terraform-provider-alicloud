---
subcategory: "Container Service for Kubernetes (ACK)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_kubernetes_autoscaler"
sidebar_current: "docs-alicloud-resource-cs-kubernetes-autoscaler"
description: |-
  Provides a Alicloud resource to manage container kubernetes cluster-autoscaler.
---

# alicloud\_cs\_kubernetes\_autoscaler

This resource will help you to manager cluster-autoscaler in Kubernetes Cluster. 

-> **NOTE:** The scaling group must use CentOS7 or AliyunLinux2 as base image.

-> **NOTE:** The cluster-autoscaler can only use the same size of instanceTypes in one scaling group. 

-> **NOTE:** Add Policy to RAM role of the node to deploy cluster-autoscaler if you need.

-> **NOTE:** Available in 1.65.0+.

-> **NOTE:** From version v1.127.0+. Resource `alicloud_cs_kubernetes_autoscaler` is replaced by resource [alicloud_cs_autoscaling_config](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/cs_autoscaling_config). If you have used resource `alicloud_cs_kubernetes_autoscaler`, please refer to [Use Terraform to create an auto-scaling node pool](https://www.alibabacloud.com/help/doc-detail/197717.htm) to switch to `alicloud_cs_autoscaling_config`.

## Example Usage

cluster-autoscaler in Kubernetes Cluster.

```terraform
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
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) The id of kubernetes cluster.
* `nodepools` - (Required) 
* `nodepools.id` - (Required) The scaling group id of the groups configured for cluster-autoscaler.
* `nodepools.taints` - (Required) The taints for the nodes in scaling group.
* `nodepools.labels` - (Required) The labels for the nodes in scaling group.
* `utilization` - (Required) The utilization option of cluster-autoscaler.
* `cool_down_duration` (Required) The cool_down_duration option of cluster-autoscaler.  
* `defer_scale_in_duration` (Required) The defer_scale_in_duration option of cluster-autoscaler.
* `use_ecs_ram_role_token` (Optional, Available in 1.88.0+) Enable autoscaler access to alibabacloud service by ecs ramrole token. default: false

#### Ignoring Changes to tags and user_data

-> **NOTE:** You can utilize the generic Terraform resource [lifecycle configuration block](https://www.terraform.io/docs/configuration/resources.html) with `ignore_changes` to create a  a autoscaler group, then ignore any changes to that tags and user_data caused externally (e.g. Application Autoscaling).
```
  # ... ignore the change about tags and user_data in alicloud_ess_scaling_configuration
  lifecycle {
    ignore_changes = [tags,user_data]
  }
```

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 90 mins) Used when creating cluster-autoscaler in the kubernetes cluster (until it reaches the initial `running` status). 
* `update` - (Defaults to 60 mins) Used when activating the cluster-autoscaler in the kubernetes cluster when necessary during update.
* `delete` - (Defaults to 60 mins) Used when deleting cluster-autoscaler in kubernetes cluster. 

