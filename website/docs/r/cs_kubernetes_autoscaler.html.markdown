---
subcategory: "Container Service for Kubernetes (ACK)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_kubernetes_autoscaler"
sidebar_current: "docs-alicloud-resource-cs-kubernetes-autoscaler"
description: |-
  Provides a Alicloud resource to manage container kubernetes cluster-autoscaler.
---

# alicloud_cs_kubernetes_autoscaler

This resource will help you to manager cluster-autoscaler in Kubernetes Cluster. 

-> **NOTE:** The scaling group must use CentOS7 or AliyunLinux2 as base image.

-> **NOTE:** The cluster-autoscaler can only use the same size of instanceTypes in one scaling group. 

-> **NOTE:** Add Policy to RAM role of the node to deploy cluster-autoscaler if you need.

-> **NOTE:** Available since v1.65.0.

-> **DEPRECATED:**  This resource has been deprecated from version `1.127.0`. Please use new resource [alicloud_cs_autoscaling_config](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/cs_autoscaling_config). If you have used resource `alicloud_cs_kubernetes_autoscaler`, please refer to [Use Terraform to create an auto-scaling node pool](https://www.alibabacloud.com/help/doc-detail/197717.htm) to switch to `alicloud_cs_autoscaling_config`.

## Example Usage

cluster-autoscaler in Kubernetes Cluster.

```terraform
variable "name" {
  default = "tf-example"
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
data "alicloud_images" "default" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}
data "alicloud_instance_types" "default" {
  availability_zone    = data.alicloud_zones.default.zones.0.id
  cpu_core_count       = 4
  memory_size          = 8
  kubernetes_node_role = "Worker"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_cs_managed_kubernetes" "default" {
  name_prefix          = var.name
  cluster_spec         = "ack.pro.small"
  worker_vswitch_ids   = [alicloud_vswitch.default.id]
  new_nat_gateway      = true
  pod_cidr             = cidrsubnet("10.0.0.0/8", 8, 36)
  service_cidr         = cidrsubnet("172.16.0.0/16", 4, 7)
  slb_internet_enabled = true
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_ess_scaling_group" "default" {
  scaling_group_name = var.name
  min_size           = 1
  max_size           = 1
  vswitch_ids        = [alicloud_vswitch.default.id]
  removal_policies   = ["OldestInstance", "NewestInstance"]
}

resource "alicloud_ess_scaling_configuration" "default" {
  scaling_group_id  = alicloud_ess_scaling_group.default.id
  image_id          = data.alicloud_images.default.images[0].id
  instance_type     = data.alicloud_instance_types.default.instance_types[0].id
  security_group_id = alicloud_security_group.default.id
  force_delete      = true
  active            = true
}

resource "alicloud_cs_kubernetes_autoscaler" "default" {
  cluster_id              = alicloud_cs_managed_kubernetes.default.id
  utilization             = "0.5"
  cool_down_duration      = "10m"
  defer_scale_in_duration = "10m"
  nodepools {
    id     = alicloud_ess_scaling_configuration.default.scaling_group_id
    labels = "a=b"
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) The id of kubernetes cluster.
* `nodepools` - (Optional) The list of the node pools. See [`nodepools`](#nodepools) below.
* `utilization` - (Required) The utilization option of cluster-autoscaler.
* `cool_down_duration` (Required) The cool_down_duration option of cluster-autoscaler.  
* `defer_scale_in_duration` (Required) The defer_scale_in_duration option of cluster-autoscaler.
* `use_ecs_ram_role_token` (Optional, Available since v1.88.0) Enable autoscaler access to alibabacloud service by ecs ramrole token. default: false

### `nodepools`

The nodepools mapping supports the following:

* `id` - (Optional) The scaling group id of the groups configured for cluster-autoscaler.
* `taints` - (Optional) The taints for the nodes in scaling group.
* `labels` - (Optional) The labels for the nodes in scaling group.

## Ignoring Changes to tags and user_data

-> **NOTE:** You can utilize the generic Terraform resource [lifecycle configuration block](https://www.terraform.io/docs/configuration/resources.html) with `ignore_changes` to create a  a autoscaler group, then ignore any changes to that tags and user_data caused externally (e.g. Application Autoscaling).
```
  # ... ignore the change about tags and user_data in alicloud_ess_scaling_configuration
  lifecycle {
    ignore_changes = [tags,user_data]
  }
```

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 90 mins) Used when creating cluster-autoscaler in the kubernetes cluster (until it reaches the initial `running` status). 
* `update` - (Defaults to 60 mins) Used when activating the cluster-autoscaler in the kubernetes cluster when necessary during update.
* `delete` - (Defaults to 60 mins) Used when deleting cluster-autoscaler in kubernetes cluster. 

