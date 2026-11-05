---
subcategory: "Container Service for Kubernetes (ACK)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_kubernetes_nodes"
sidebar_current: "docs-alicloud-datasource-cs-kubernetes-nodes"
description: |-
  Provides a list of nodes of an ACK cluster.
---

# alicloud_cs_kubernetes_nodes

This data source provides the nodes of an ACK cluster, optionally limited to a single node pool. See [What is a node](https://next.api.alibabacloud.com/document/CS/2015-12-15/DescribeClusterNodes)

-> **NOTE:** Available since v2.2.0.

-> **NOTE:** This data source only reads existing nodes. It does not create, modify or delete any node.

## Example Usage

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_instance_types" "default" {
  availability_zone    = "cn-hangzhou-i"
  cpu_core_count       = 4
  memory_size          = 8
  kubernetes_node_role = "Worker"
  system_disk_category = "cloud_essd"
  eni_amount           = 2
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 8)
  zone_id      = "cn-hangzhou-i"
  vswitch_name = var.name
}

resource "alicloud_cs_managed_kubernetes" "default" {
  name         = var.name
  cluster_spec = "ack.pro.small"
  vswitch_ids  = [alicloud_vswitch.default.id]
  pod_cidr     = "172.20.0.0/16"
  service_cidr = "172.21.0.0/20"
}

resource "alicloud_cs_kubernetes_node_pool" "default" {
  node_pool_name       = "${var.name}-pool"
  cluster_id           = alicloud_cs_managed_kubernetes.default.id
  vswitch_ids          = [alicloud_vswitch.default.id]
  instance_types       = [data.alicloud_instance_types.default.instance_types.0.id]
  desired_size         = "2"
  system_disk_category = "cloud_essd"
  system_disk_size     = "40"
  image_type           = "AliyunLinux3ContainerOptimized"
}

data "alicloud_cs_kubernetes_nodes" "default" {
  cluster_id    = alicloud_cs_managed_kubernetes.default.id
  nodepool_name = alicloud_cs_kubernetes_node_pool.default.node_pool_name
}

output "node_names" {
  value = data.alicloud_cs_kubernetes_nodes.default.node_names
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String, Available since v2.2.0) The ID of the ACK cluster.
* `nodepool_id` - (Optional, Computed, String, Available since v2.2.0) The ID of the node pool whose nodes are queried. A configured value takes precedence over `nodepool_name`. When only `nodepool_name` is configured, the resolved node pool ID is exported here.
* `nodepool_name` - (Optional, String, Available since v2.2.0) The name of the node pool whose nodes are queried. It is resolved to the node pool ID on every read, so a node pool created in the same configuration can be referenced without creating a dependency cycle. Because the resolution is repeated on each read, changing the name queries the newly named node pool, and a node pool that was replaced by a new one of the same name is picked up as well. If no node pool of the cluster has this name, the data source returns empty lists. When neither `nodepool_id` nor `nodepool_name` is set, all nodes of the cluster are returned.

## Attributes Reference

The following attributes are exported:

* `id` - (String, Available since v2.2.0) The ID of the data source. It is the cluster ID, followed by the node pool ID separated by a colon when a node pool is selected.
* `ids` - (List, Available since v2.2.0) A list of instance IDs of the nodes.
* `node_names` - (List, Available since v2.2.0) A list of Kubernetes node names, sorted in ascending order.
* `nodes` - (List, Available since v2.2.0) A list of nodes. The nodes are sorted by node name. Each element contains:
  * `id` - (String, Available since v2.2.0) The ID of the instance.
  * `name` - (String, Available since v2.2.0) The name of the node in the cluster.
  * `nodepool_id` - (String, Available since v2.2.0) The ID of the node pool to which the node belongs.
  * `host_name` - (String, Available since v2.2.0) The hostname of the node.
  * `instance_role` - (String, Available since v2.2.0) The role of the node. Valid values: `Master`, `Worker`.
  * `state` - (String, Available since v2.2.0) The running status of the node.
