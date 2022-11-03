---
subcategory: "EDAS"
layout: "alicloud"
page_title: "Alicloud: alicloud_edas_cluster"
sidebar_current: "docs-alicloud-resource-edas-cluster"
description: |-
  Provides an EDAS cluster resource.
---

# alicloud\_edas\_cluster

Provides an EDAS cluster resource.

-> **NOTE:** Available in 1.82.0+

## Example Usage

Basic Usage

```terraform
resource "alicloud_edas_cluster" "default" {
  cluster_name      = var.cluster_name
  cluster_type      = var.cluster_type
  network_mode      = var.network_mode
  logical_region_id = var.logical_region_id
  vpc_id            = var.vpc_id
}
```

## Argument Reference

The following arguments are supported:

* `cluster_name` - (Required, ForceNew) The name of the cluster that you want to create.
* `cluster_type` - (Required, ForceNew) The type of the cluster that you want to create. Valid values only: 2: ECS cluster.
* `network_mode` - (Required, ForceNew) The network type of the cluster that you want to create. Valid values: 1: classic network. 2: VPC.
* `logical_region_id` - (Optional, ForceNew) The ID of the namespace where you want to create the application. You can call the ListUserDefineRegion operation to query the namespace ID.
* `vpc_id` - (Optional, ForceNew) The ID of the Virtual Private Cloud (VPC) for the cluster.

## Attributes Reference

The following attributes are exported:

* `id` - The `key` of the resource supplied above. The value is formulated as `<cluster_id>`.

## Import

EDAS cluster can be imported using the id, e.g.

```shell
$ terraform import alicloud_edas_cluster.cluster cluster_id
```
