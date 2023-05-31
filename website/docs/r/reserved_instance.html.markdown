---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_reserved_instance"
sidebar_current: "docs-alicloud-resource-reserved-instance"
description: |-
  Provides an ECS Reserved Instance resource.
---

# alicloud\_reserved\_instance

Provides an Reserved Instance resource.

-> **NOTE:** Available in 1.65.0+

## Example Usage

```terraform
data "alicloud_instance_types" "default" {
  instance_type_family = "ecs.g6"
}

resource "alicloud_reserved_instance" "default" {
  instance_type          = data.alicloud_instance_types.default.instance_types.0.id
  instance_amount        = "1"
  period_unit            = "Month"
  offering_type          = "All Upfront"
  reserved_instance_name = "terraform-example"
  description            = "ReservedInstance"
  zone_id                = data.alicloud_instance_types.default.instance_types.0.availability_zones.0
  scope                  = "Zone"
}
```

## Argument Reference

The following arguments are supported:

* `offering_type` - (Optional, Computed, ForceNew) Payment type of the RI. Default value: `All Upfront`. Valid values:
  - `No Upfront`: No upfront payment.
  - `Partial Upfront`: A portion of upfront payment.
  - `All Upfront`: Full upfront payment.
* `zone_id` - (Optional, ForceNew) ID of the zone to which the RI belongs. When Scope is set to Zone, this parameter is required. For information about the zone list, see [DescribeZones](https://www.alibabacloud.com/help/doc-detail/25610.html).
* `scope` - (Optional, Computed, ForceNew) Scope of the RI. Optional values: `Region`: region-level, `Zone`: zone-level. Default is `Region`.
* `instance_type` - (Required, ForceNew) Instance type of the RI. For more information, see [Instance type families](https://www.alibabacloud.com/help/doc-detail/25378.html).
* `instance_amount` - (Optional, ForceNew) Number of instances allocated to an RI (An RI is a coupon that includes one or more allocated instances.).
* `period` - (Optional, ForceNew) The validity period of the reserved instance. Default value: `1`. **NOTE:** From version 1.183.0, `period` can be set to `5`, when `period_unit` is `Year`.
  - When `period_unit` is `Year`, Valid values: `1`, `3`, `5`.
  - When `period_unit` is `Month`, Valid values: `1`.
* `period_unit` - (Optional, ForceNew) The unit of the validity period of the reserved instance. Valid value: `Month`, `Year`. Default value: `Year`. **NOTE:** From version 1.183.0, `period_unit` can be set to `Month`.
* `resource_group_id` - (Optional, ForceNew) Resource group ID.
* `description` - (Optional) Description of the RI. 2 to 256 English or Chinese characters. It cannot start with `http://` or `https://`.
* `name` - (Optional, Computed, Deprecated from v1.194.0+) Field `name` has been deprecated from provider version 1.194.0. New field `reserved_instance_name` instead.
* `platform` - (Optional, ForceNew) The operating system type of the image used by the instance. Optional values: `Windows`, `Linux`. Default is `Linux`.
* `reserved_instance_name` - (Optional, Computed, Available in v1.194.0+)  Name of the RI. The name must be a string of 2 to 128 characters in length and can contain letters, numbers, colons (:), underscores (_), and hyphens. It must start with a letter. It cannot start with http:// or https://.
* `renewal_status` - (Optional, Computed, Available in v1.194.0+) Automatic renewal status. Valid values: `AutoRenewal`,`Normal`.
* `auto_renew_period` - (Optional, Computed, Available in v1.194.0+) The auto-renewal term of the reserved instance. This parameter takes effect only when AutoRenew is set to true. Valid values: 1, 12, 36, and 60. Default value when `period_unit` is set to Month: 1 Default value when `period_unit` is set to Year: 12
* `tags` - (Optional, Available in v1.194.0+) A mapping of tags to assign to the resource.

### Removing alicloud_reserved_instance from your configuration
 
The alicloud_reserved_instance resource allows you to manage your ReservedInstance, but Terraform cannot destroy it. Removing this resource from your configuration will remove it from your statefile and management, but will not destroy the ReservedInstance.
 

## Attributes Reference

The following attributes are exported:

* `id` - ID of the ReservedInstance.
* `allocation_status` - Indicates the sharing status of the reserved instance when the AllocationType parameter is set to Shared. Valid values: `allocated`: The reserved instance is allocated to another account. `beAllocated`: The reserved instance is allocated by another account.
* `create_time` -  The time when the reserved instance was created.
* `expired_time` -  The time when the reserved instance expires.
* `operation_locks` -  Details about the lock status of the reserved instance.
  * `lock_reason` - The reason why the reserved instance was locked.
* `start_time` -  The time when the reserved instance took effect.
* `status` -  The status of the reserved instance.

#### Timeouts

-> **NOTE:** Available in 1.194.0+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the resource.
* `update` - (Defaults to 1 mins) Used when update the resource.

## Import

reservedInstance can be imported using id, e.g.

```shell
$ terraform import alicloud_reserved_instance.default ecsri-uf6df4xm0h3licit****
```

