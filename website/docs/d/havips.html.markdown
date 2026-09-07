---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_havips"
sidebar_current: "docs-alicloud-datasource-havips"
description: |-
  Provides a list of Havips to the user.
---

# alicloud\_havips

This data source provides the Havips of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.120.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_havips" "example" {
  ids        = ["example_value"]
  name_regex = "the_resource_name"
}

output "first_havip_id" {
  value = data.alicloud_havips.example.havips.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Ha Vip IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Ha Vip name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of HaVip instance. Valid value: `Available`, `InUse` and `Pending`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Ha Vip names.
* `havips` - A list of Havips. Each element contains the following attributes:
	* `associated_eip_addresses` - EIP bound to HaVip.
	* `associated_instances` - An ECS instance that is bound to HaVip.
	* `description` - Dependence of a HaVip instance.
	* `havip_id` - The  ID of the resource.
	* `havip_name` - The name of the HaVip instance.
	* `id` - The ID of the Ha Vip.
	* `ip_address` - IP address of private network.
	* `master_instance_id` - The primary instance ID bound to HaVip.
	* `status` - The status.
	* `vpc_id` - The VPC ID to which the HaVip instance belongs.
	* `vswitch_id` - The vswitch id.
