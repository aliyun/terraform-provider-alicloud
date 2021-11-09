---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_hpc_cluster"
sidebar_current: "docs-alicloud-resource-ecs-hpc-cluster"
description: |-
  Provides a Alicloud ECS Hpc Cluster resource.
---

# alicloud\_ecs\_hpc\_cluster

Provides a ECS Hpc Cluster resource.

For information about ECS Hpc Cluster and how to use it, see [What is Hpc Cluster](https://www.alibabacloud.com/help/en/doc-detail/109138.htm).

-> **NOTE:** Available in v1.116.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ecs_hpc_cluster" "example" {
  name        = "tf-testAcc"
  description = "For Terraform Test"
}

```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) The description of ECS Hpc Cluster.
* `name` - (Required) The name of ECS Hpc Cluster.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Hpc Cluster.

## Import

ECS Hpc Cluster can be imported using the id, e.g.

```
$ terraform import alicloud_ecs_hpc_cluster.example <id>
```
