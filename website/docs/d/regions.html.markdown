---
layout: "alicloud"
page_title: "Alicloud: alicloud_regions"
sidebar_current: "docs-alicloud-datasource-regions"
description: |-
    Provides a list of regions that can be used by an Alibaba Cloud account.
---

# alicloud\_regions

This data source provides Alibaba Cloud regions.

## Example Usage

```
data "alicloud_regions" "current_region_ds" {
  current = true
}

output "current_region_id" {
  value = "${data.alicloud_regions.current_region_ds.regions.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the region to select, such as `eu-central-1`.
* `current` - (Optional) Set to true to match only the region configured in the provider.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

~> **NOTE:** You will get an error if you set `current` to true and `name` to a different value from the one you configured in the provider.
 It is better to either use `name` or `current`, but not both at the same time.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `regions` - A list of regions. Each element contains the following attributes:
  * `id` - ID of the region.
  * `local_name` - Name of the region in the local language.
