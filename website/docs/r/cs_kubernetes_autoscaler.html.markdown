---
subcategory: "Container Service (CS)"
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

## Example Usage

cluster-autoscaler in Kubernetes Cluster

```
resource "alicloud_cs_kubernetes_autoscaler" "default" {
  cluster_id              = "${var.cluster_id}"
  nodepools {
        id                = "scaling_group_id"
        taints            = "c=d:NoSchedule"
        labels            = "a=b"
  }
  utilization             = "${var.utilization}"
  cool_down_duration      = "${var.cool_down_duration}"
  defer_scale_in_duration = "${var.defer_scale_in_duration}"
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

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 90 mins) Used when creating cluster-autoscaler in the kubernetes cluster (until it reaches the initial `running` status). 
* `update` - (Defaults to 60 mins) Used when activating the cluster-autoscaler in the kubernetes cluster when necessary during update.
* `delete` - (Defaults to 60 mins) Used when deleting cluster-autoscaler in kubernetes cluster. 

