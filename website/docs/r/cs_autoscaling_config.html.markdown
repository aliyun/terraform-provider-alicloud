---
subcategory: "Container Service for Kubernetes (ACK)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_autoscaling_config"
sidebar_current: "docs-alicloud-cs-autoscaling-config"
description: |-
  Provides a Alicloud resource to configure auto scaling for for ACK cluster.
---

# alicloud\_cs\_autoscaling\_config

This resource will help you configure auto scaling for the kubernetes cluster. 

-> **NOTE:** Available in v1.127.0+.
-> **NOTE:** From version 1.164.0, support for specifying whether to allow the scale-in of nodes by parameter `scale_down_enabled`.
-> **NOTE:** From version 1.164.0, support for selecting the policy for selecting which node pool to scale by parameter `expander`.

## Example Usage
If you do not have an existing cluster, you need to create an ACK cluster through [alicloud_cs_managed_kubernetes](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/cs_managed_kubernetes) first, and then configure automatic scaling.

```terraform
resource "alicloud_cs_autoscaling_config" "default" {
  cluster_id = alicloud_cs_managed_kubernetes.default.0.id
  // configure auto scaling
  cool_down_duration        = "10m"
  unneeded_duration         = "10m"
  utilization_threshold     = "0.5"
  gpu_utilization_threshold = "0.5"
  scan_interval             = "30s"
  scale_down_enabled        = true
  expander                  = "least-waste"
}
```

## Argument Reference

The following arguments are supported.

* `cluster_id` - (Optional, ForceNew) The id of kubernetes cluster.
* `cool_down_duration` - (Optional) The cool down duration. Default is `10m`. If the delay (cooldown) value is set too long, there could be complaints that the Horizontal Pod Autoscaler is not responsive to workload changes. However, if the delay value is set too short, the scale of the replicas set may keep thrashing as usual.
* `unneeded_duration` - (Optional) The unneeded duration. Default is `10m`.
* `utilization_threshold` - (Optional) The scale-in threshold. Default is `0.5`. 
* `gpu_utilization_threshold` - (Optional)  The scale-in threshold for GPU instance. Default is `0.5`. 
* `scan_interval` - (Optional) The interval at which the cluster is reevaluated for scaling. Default is `30s`.
* `scale_down_enabled` - (Optional) Specify whether to allow the scale-in of nodes. Default is `true`.
* `expander` - (Optional) The policy for selecting which node pool to scale. Valid values: `least-waste`, `random`, `priority`. For more information on these policies, see [Configure auto scaling](https://www.alibabacloud.com/help/en/container-service-for-kubernetes/latest/auto-scaling-of-nodes#section-3bg-2ko-inl)

## Attributes Reference

The following attributes are exported:
* `id` - Resource id.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 90 mins) Used when creating the kubernetes cluster (until it reaches the initial `running` status).
* `update` - (Defaults to 60 mins) Used when activating the kubernetes cluster when necessary during update.
* `delete` - (Defaults to 60 mins) Used when terminating the kubernetes cluster.

