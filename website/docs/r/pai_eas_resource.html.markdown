---
subcategory: "EAS"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_eas_resource"
sidebar_current: "docs-alicloud-resource-pai-eas-resource"
description: |-
  Provides a PAI EAS Resource Group resource.
---

# alicloud_pai_eas_resource

Provides a PAI EAS Resource Group resource.

An EAS resource group is a container that groups EAS service instances together
with the backing compute resources. For information about EAS Resource Group and
how to use it, see
[CreateResource](https://www.alibabacloud.com/help/en/pai/user-guide/eas/eas-resource-group-management).

-> **NOTE:** Available since v1.246.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_pai_eas_resource" "example" {
  resource_name = "terraform-example"
  charge_type   = "PostPaid"
}
```

## Argument Reference

The following arguments are supported:

* `resource_name` - (Required) The name of the resource group. This is the only
  attribute that can be updated after creation.
* `ecs_instance_count` - (Optional, ForceNew) The number of instances. Required
  when `resource_type` is `Dedicated`.
* `ecs_instance_type` - (Optional, ForceNew) The ECS instance type.
* `system_disk_size` - (Optional, ForceNew) The system disk size in GB. Valid
  values: 200 to 2000. Defaults to 200.
* `charge_type` - (Optional, ForceNew) The charge type. Valid values: `PrePaid`
  and `PostPaid`.
* `auto_renewal` - (Optional, ForceNew) Whether to auto renew. Valid only when
  `charge_type` is `PrePaid`.
* `zone` - (Optional, ForceNew) The zone ID.
* `resource_type` - (Optional, ForceNew) The resource type.
* `usage_mode` - (Optional, ForceNew) The usage mode.
* `labels` - (Optional, ForceNew) The labels of the resource group. Map of
  string key/value pairs.
* `external_cluster_id` - (Optional, ForceNew) The external cluster ID. Part of
  the self-managed resource options.
* `role_name` - (Optional, ForceNew) The role name. Part of the self-managed
  resource options.
* `node_match_labels` - (Optional, ForceNew) The node match labels. Map of
  string key/value pairs. Part of the self-managed resource options.

  -> **NOTE** The self-managed resource options
    (`external_cluster_id`, `role_name` and `node_match_labels`) are create-only:
    the DescribeResource API does not return them, so they cannot be read back or
    updated in place.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource group.
* `region_id` - The region (cluster) ID of the resource group.
* `status` - The status of the resource group.
* `create_time` - The time when the resource group was created.

## Import

PAI EAS Resource Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_pai_eas_resource.example <resource_id>
```
