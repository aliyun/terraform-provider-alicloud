---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_storage_capacity_unit"
sidebar_current: "docs-alicloud-resource-ecs-storage-capacity-unit"
description: |-
  Provides a Alicloud ECS Storage Capacity Unit resource.
---

# alicloud\_ecs\_storage\_capacity\_unit

Provides a ECS Storage Capacity Unit resource.

For information about ECS Storage Capacity Unit and how to use it, see [What is Storage Capacity Unit](https://www.alibabacloud.com/help/en/doc-detail/161157.html).

-> **NOTE:** Available in v1.155.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ecs_storage_capacity_unit" "default" {
  capacity                   = 20
  description                = "tftestdescription"
  storage_capacity_unit_name = "tftestname"
}
```

## Argument Reference

The following arguments are supported:

* `capacity` - (Required, ForceNew) The capacity of the Storage Capacity Unit. Unit: GiB. Valid values: `20`, `40`, `100`, `200`, `500`, `1024`, `2048`, `5120`, `10240`, `20480`, and `51200`.
* `description` - (Optional) The description of the Storage Capacity Unit. The description must be 2 to 256 characters in length and cannot start with `http://` or `https://`.
* `period` - (Optional, Computed) The validity period of the Storage Capacity Unit. Default value: `1`.
  * When PeriodUnit is set to Month, Valid values: `1`, `2`, `3`, `6`.
  * When PeriodUnit is set to Year, Valid values: `1`, `3`, `5`.
* `period_unit` - (Optional, Computed) The unit of the validity period of the Storage Capacity Unit. Default value: `Month`. Valid values: `Month`, `Year`.
* `start_time` - (Optional, ForceNew, Computed) The time when the Storage Capacity Unit takes effect. It cannot be earlier than or more than six months later than the time when the Storage Capacity Unit is created. Specify the time in the ISO 8601 standard in the `yyyy-MM-ddTHH:mm:ssZ` format. The time must be in UTC. **NOTE:** This parameter is empty by default. The Storage Capacity Unit immediately takes effect after it is created.
* `storage_capacity_unit_name` - (Optional, Computed) The name of the Storage Capacity Unit.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Storage Capacity Unit.
* `status` - The status of Storage Capacity Unit.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Storage Capacity Unit.

## Import

ECS Storage Capacity Unit can be imported using the id, e.g.

```
$ terraform import alicloud_ecs_storage_capacity_unit.example <id>
```