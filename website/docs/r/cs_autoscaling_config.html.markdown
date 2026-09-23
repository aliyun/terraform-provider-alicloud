---
subcategory: "Container Service for Kubernetes (ACK)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_autoscaling_config"
sidebar_current: "docs-alicloud-cs-autoscaling-config"
description: |-
  Provides a Alicloud resource to configure auto scaling for for ACK cluster.
---

# alicloud_cs_autoscaling_config

This resource will help you configure auto scaling for the kubernetes cluster, see [What is autoscaling config](https://www.alibabacloud.com/help/en/ack/ack-managed-and-ack-dedicated/developer-reference/api-configure-auto-scaling).

-> **NOTE:** Available since v1.127.0.

-> **NOTE:** From version 1.164.0, support for specifying whether to allow the scale-in of nodes by parameter `scale_down_enabled`.

-> **NOTE:** From version 1.164.0, support for selecting the policy for selecting which node pool to scale by parameter `expander`.

-> **NOTE:** From version 1.237.0, support for selecting the type of autoscaler by parameter `scaler_type`.

-> **NOTE:** From version 1.256.0, support for setting the priority of scaling groups by parameter `priorities`.

## Example Usage
If you do not have an existing cluster, you need to create an ACK cluster through [alicloud_cs_managed_kubernetes](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/cs_managed_kubernetes) first, and then configure automatic scaling.

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cs_autoscaling_config&exampleId=580a36e2-202d-1eae-0cfa-f1b6c560ee5863927387&activeTab=example&spm=docs.r.cs_autoscaling_config.0.580a36e220&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
data "alicloud_instance_types" "default" {
  availability_zone    = "${data.alicloud_zones.default.zones.0.id}"
  cpu_core_count       = 4
  memory_size          = 8
  kubernetes_node_role = "Worker"
  system_disk_category = "cloud_essd"
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

resource "alicloud_cs_kubernetes_node_pool" "default" {
  count                = 3
  node_pool_name       = format("%s-%d", var.name, count.index)
  cluster_id           = alicloud_cs_managed_kubernetes.default.id
  vswitch_ids          = [alicloud_vswitch.default.id]
  instance_types       = [data.alicloud_instance_types.default.instance_types.0.id]
  system_disk_category = "cloud_essd"
  system_disk_size     = 40
  image_type           = "AliyunLinux3"
  scaling_config {
    enable   = true
    min_size = 0
    max_size = 10
  }
}

resource "alicloud_cs_autoscaling_config" "default" {
  cluster_id = alicloud_cs_managed_kubernetes.default.id
  // configure auto scaling
  cool_down_duration            = "10m"
  unneeded_duration             = "10m"
  utilization_threshold         = "0.5"
  gpu_utilization_threshold     = "0.5"
  scan_interval                 = "30s"
  scale_down_enabled            = true
  expander                      = "priority"
  skip_nodes_with_system_pods   = true
  skip_nodes_with_local_storage = false
  daemonset_eviction_for_nodes  = false
  max_graceful_termination_sec  = 14400
  min_replica_count             = 0
  recycle_node_deletion_enabled = false
  scale_up_from_zero            = true
  scaler_type                   = "cluster-autoscaler"
  priorities = {
    10 = join(",", [
      alicloud_cs_kubernetes_node_pool.default[0].scaling_group_id,
      alicloud_cs_kubernetes_node_pool.default[1].scaling_group_id,
    ])
    20 = alicloud_cs_kubernetes_node_pool.default[2].scaling_group_id
  }
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cs_autoscaling_config&spm=docs.r.cs_autoscaling_config.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported.

* `cluster_id` - (Optional, ForceNew) The id of kubernetes cluster.
* `cool_down_duration` - (Optional) Specify the time interval between detecting a scale-in requirement (when the threshold is reached) and actually executing the scale-in operation (reducing the number of Pods). Default is `10m`. If the delay (cooldown) value is set too long, there could be complaints that the Horizontal Pod Autoscaler is not responsive to workload changes. However, if the delay value is set too short, the scale of the replicas set may keep thrashing as usual.
* `unneeded_duration` - (Optional) Specify the time interval during which autoscaler does not perform scale-in operations after the most recent scale-out completion. Nodes added through scale-out can only be considered for scale-in after the period has elapsed. Default is `10m`.
* `utilization_threshold` - (Optional) The scale-in a threshold. Default is `0.5`. 
* `gpu_utilization_threshold` - (Optional) The scale-in threshold for GPU instance. Default is `0.5`. 
* `scan_interval` - (Optional) The interval at which the cluster is reevaluated for scaling. Default is `30s`.
* `scale_down_enabled` - (Optional) Specify whether to allow the scale-in of nodes. Default is `true`.
* `expander` - (Optional) The policy for selecting which node pool to scale. Valid values: `least-waste`, `random`, `priority`. For scaler type `goatscaler`, only the `least-waste` expander is currently supported. For more information on these policies, see [Configure auto scaling](https://www.alibabacloud.com/help/en/container-service-for-kubernetes/latest/auto-scaling-of-nodes#section-3bg-2ko-inl)
* `skip_nodes_with_system_pods` - (Optional, Available since v1.209.0) If true cluster autoscaler will never delete nodes with pods from kube-system (except for DaemonSet or mirror pods). Default is `true`.
* `skip_nodes_with_local_storage` - (Optional, Available since v1.209.0) If true cluster autoscaler will never delete nodes with pods with local storage, e.g. EmptyDir or HostPath. Default is `false`.
* `daemonset_eviction_for_nodes` - (Optional, Available since v1.209.0) If true DaemonSet pods will be  terminated from nodes. Default is `false`. 
* `max_graceful_termination_sec` - (Optional, Available since v1.209.0) Maximum number of seconds CA waits for pod termination when trying to scale down a node. Default is `14400`.
* `min_replica_count` - (Optional, Available since v1.209.0) Minimum number of replicas that a replica set or replication controller should have to allow their pods deletion in scale down. Default is `0`.
* `recycle_node_deletion_enabled` - (Optional, Available since v1.209.0) Should CA delete the K8s node object when recycle node has scaled down successfully. Default is `false`.
* `scale_up_from_zero` - (Optional, Available since v1.209.0) Should CA scale up when there 0 ready nodes. Default is `true`.
* `scaler_type` - (Optional, Available since v1.237.0) The type of autoscaler. Valid values: `cluster-autoscaler`, `goatscaler`. For cluster version 1.22 and below, we only support `cluster-autoscaler`. When switching from `cluster-autoscaler` to `goatscaler`, all configuration parameters will be automatically migrated.
* `priorities` - (Optional, Available since v1.256.0) Priority settings for autoscaling node pool scaling groups. This parameter only takes effect when `expander` is set to `priority`. Only supports scaler type `cluster-autoscaler`. Uses key-value pairs where the key is the priority value, and the value is a comma-separated list of scaling group IDs. High numerical values indicate higher priority.

## Attributes Reference

The following attributes are exported:
* `id` - Resource id.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 90 mins) Used when creating the kubernetes cluster (until it reaches the initial `running` status).
* `update` - (Defaults to 60 mins) Used when activating the kubernetes cluster when necessary during update.
* `delete` - (Defaults to 60 mins) Used when terminating the kubernetes cluster.

