---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_cache_reserve_instances"
description: |-
  Provides a list of Esa Cache Reserve Instances to the user.
---

# alicloud_esa_cache_reserve_instances

This data source provides the Esa Cache Reserve Instances of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.282.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_esa_cache_reserve_instance" "default" {
  quota_gb     = "10240"
  cr_region    = "CN-beijing"
  auto_renew   = true
  period       = "1"
  payment_type = "Subscription"
  auto_pay     = true
}

data "alicloud_esa_cache_reserve_instances" "ids" {
  ids = [alicloud_esa_cache_reserve_instance.default.id]
}

output "esa_cache_reserve_instances_id_0" {
  value = data.alicloud_esa_cache_reserve_instances.ids.instances.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, List) A list of Cache Reserve Instance IDs.
* `cache_reserve_instance_id` - (Optional) The ID of the Cache Reserve Instance.
* `status` - (Optional) The status of the cache reserve instance. Valid values: `online`, `offline`, `disable`, `overdue`.
* `sort_by` - (Optional) The field to sort the results by. Valid values: `CreateTime`, `ExpireTime`.
* `sort_order` - (Optional) The sort order. Valid values: `asc`, `desc`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `instances` - A list of Cache Reserve Instances. Each element contains the following attributes:
  * `id` - The ID of the Cache Reserve Instance.
  * `cache_reserve_instance_id` - The ID of the Cache Reserve Instance.
  * `quota_gb` - The cache reserve capacity, in GB.
  * `cr_region` - The region where the cache reserve instance is used.
  * `payment_type` - The payment type of the resource.
  * `period` - The purchase duration of the instance, in months.
  * `status` - The status of the instance.
  * `create_time` - The time when the instance was created.
  * `expire_time` - The expiration time of the instance.
