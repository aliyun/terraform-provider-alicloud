---
subcategory: "EDAS"
layout: "alicloud"
page_title: "Alicloud: alicloud_edas_k8s_cluster"
sidebar_current: "docs-alicloud-resource-edas-k8s-cluster"
description: |-
  Provides an EDAS K8s cluster resource.
---

# alicloud\_edas\_k8s\_cluster

Provides an EDAS K8s cluster resource.

-> **NOTE:** Available in 1.92.0+

## Example Usage

Basic Usage

```
resource "alicloud_edas_k8s_cluster" "default" {
  cs_cluster_id = var.cs_cluster_id
}

```

## Argument Reference

The following arguments are supported:

* `cs_cluster_id` - (Required, ForceNew) The ID of the alicloud container service kubernetes cluster that you want to import.
* `namespace_id` - (Optional, ForceNew) The ID of the namespace where you want to import. You can call the ListUserDefineRegion operation to query the namespace ID.


## Attributes Reference

The following attributes are exported:

* `cluster_name` - The name of the cluster that you want to create.
* `cluster_type` - The type of the cluster that you want to create. Valid values only: 2: ECS cluster.
* `network_mode` - The network type of the cluster that you want to create. Valid values: 1: classic network. 2: VPC.
* `region_id` - (Optional, ForceNew) The ID of the region.
* `vpc_id` - (Optional, ForceNew) The ID of the Virtual Private Cloud (VPC) for the cluster.
* `cluster_import_status` - The import status of cluster, 1 for success, 2 for failed, 3 for importing, 4 for deleted.

## Import

EDAS cluster can be imported using the id, e.g.

```
$ terraform import alicloud_edas_k8s_cluster.cluster cluster_id
```
