---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_dedicated_host_cluster"
sidebar_current: "docs-alicloud-resource-ecs-dedicated-host-cluster"
description: |-
  Provides a Alicloud ECS Dedicated Host Cluster resource.
---

# alicloud\_ecs\_dedicated\_host\_cluster

Provides a ECS Dedicated Host Cluster resource.

For information about ECS Dedicated Host Cluster and how to use it, see [What is Dedicated Host Cluster](https://www.alibabacloud.com/help/en/doc-detail/184667.html).

-> **NOTE:** Available in v1.146.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_zones" example {}

resource "alicloud_ecs_dedicated_host_cluster" "example" {
  dedicated_host_cluster_name = "example_value"
  description                 = "example_value"
  zone_id                     = data.alicloud_zones.example.zones.0.id
  tags = {
    Create = "TF"
    For    = "DDH_Cluster_Test",
  }
}

```

## Argument Reference

The following arguments are supported:

* `dedicated_host_cluster_name` - (Optional) The name of the dedicated host cluster. The name must be `2` to `128` characters in length and can contain letters, digits, periods (.), underscores (_), and hyphens (-). It must start with a letter. It cannot contain `http://` or `https://`.
* `description` - (Optional) The description of the dedicated host cluster. The description must be `2` to `256` characters in length. It cannot start with `http://` or `https://`.
* `dry_run` - (Optional) The dry run.
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `zone_id` - (Required, ForceNew) The ID of the zone in which to create the dedicated host cluster.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Dedicated Host Cluster.

## Import

ECS Dedicated Host Cluster can be imported using the id, e.g.

```
$ terraform import alicloud_ecs_dedicated_host_cluster.example <id>
```