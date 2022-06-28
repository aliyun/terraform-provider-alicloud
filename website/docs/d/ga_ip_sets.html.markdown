---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_ip_sets"
sidebar_current: "docs-alicloud-datasource-ga-ip-sets"
description: |-
  Provides a list of Global Accelerator (GA) Ip Sets to the user.
---

# alicloud\_ga\_ip\_sets

This data source provides the Global Accelerator (GA) Ip Sets of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.113.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ga_ip_sets" "example" {
  accelerator_id = "example_value"
  ids            = ["example_value"]
}

output "first_ga_ip_set_id" {
  value = data.alicloud_ga_ip_sets.example.sets.0.id
}
```

## Argument Reference

The following arguments are supported:

* `accelerator_id` - (Required, ForceNew) The ID of the Global Accelerator (GA) instance.
* `ids` - (Optional, ForceNew, Computed)  A list of Ip Set IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew)  The status of the acceleration region. Valid values: `active`, `deleting`, `init`, `updating`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `sets` - A list of Ga Ip Sets. Each element contains the following attributes:
	* `accelerate_region_id` -  The ID of an acceleration region.
	* `bandwidth` - The bandwidth allocated to the acceleration region.
	* `id` - The ID of the Ip Set.
	* `ip_address_list` - The list of accelerated IP addresses in the acceleration region.
	* `ip_set_id` -  Accelerated area ID.
	* `ip_version` - The IP protocol used by the GA instance.
	* `status` -  The status of the acceleration region.
