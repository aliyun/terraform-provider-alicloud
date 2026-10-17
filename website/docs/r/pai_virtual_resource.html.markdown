---
subcategory: "PAI (EAS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_virtual_resource"
sidebar_current: "docs-alicloud-resource-pai-virtual-resource"
description: |-
  Provides a PAI EAS Virtual Resource resource.
---

# alicloud_pai_virtual_resource

Provides a PAI EAS Virtual Resource resource. A virtual resource group organizes
the compute resources that an EAS service can be scheduled onto. Each entry in the
`resources` list identifies a backing resource group by exactly one of `instance_type`
(public resource group instance type), `resource_id` (EAS dedicated resource group)
or `quota_id` (Lingjun智算 resource quota).

For information about Virtual Resource and how to use it, see
[CreateVirtualResource](https://www.alibabacloud.com/help/en/pai/developer-reference/api-eas-2021-07-01-createvirtualresource).

-> **NOTE:** Available since v1.237.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_pai_virtual_resource" "default" {
  virtual_resource_name          = "tf-example"
  disable_spot_protection_period = true
  resources {
    region           = "cn-hangzhou"
    instance_type    = "ecs.s6-c1m2.xlarge"
    priority         = 6
    spot_price_limit = 10.5
  }
}
```

## Argument Reference

The following arguments are supported:

* `virtual_resource_name` - (Required) The name of the virtual resource group.
* `disable_spot_protection_period` - (Optional) Whether to disable the spot protection period for spot instances in the group. Defaults to `false`.
* `resources` - (Required) The list of resources contained in the virtual resource group. Each entry must specify exactly one of `instance_type`, `resource_id` and `quota_id` to identify the backing resource group.
  * `region` - (Optional) The region of the resource entry.
  * `priority` - (Optional) The scheduling priority of the resource. A larger number means a higher priority.
  * `instance_type` - (Optional) The instance type of the public resource group, e.g. `ecs.s6-c1m2.xlarge`. Exactly one of `instance_type`, `resource_id` and `quota_id` must be set.
  * `quota_id` - (Optional) The Lingjun resource quota ID, for Lingjun智算 resource groups. Exactly one of `instance_type`, `resource_id` and `quota_id` must be set.
  * `spot_price_limit` - (Optional) The spot price upper limit for public resource group spot instances. Leave unset to not use spot instances.
  * `resource_id` - (Optional) The EAS dedicated resource group ID. Exactly one of `instance_type`, `resource_id` and `quota_id` must be set.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource. Same as `virtual_resource_id`.
* `virtual_resource_id` - The first-level ID of the virtual resource group.
* `region_id` - The region ID of the resource, mapped from the EAS ClusterId.
* `create_time` - The creation time of the virtual resource group.

## Import

PAI Virtual Resource can be imported using the id, e.g.

```shell
$ terraform import alicloud_pai_virtual_resource.example <virtual_resource_id>
```
