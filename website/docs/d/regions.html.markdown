---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_regions"
sidebar_current: "docs-alicloud-datasource-regions"
description: |-
    Provides a list of regions that can be used by an Alibaba Cloud account.
---

# alicloud_regions

This data source provides Alibaba Cloud regions.

-> **NOTE:** Available since v0.1.0

## Example Usage

```terraform
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

* `ids` - A list of region IDs.
* `regions` - A list of regions.

---

### `regions`

The regions supports the following attributes:

* `id` - ID of the region.
* `region_id` - ID of the region.
* `local_name` - Name of the region in the local language.

  -> **NOTE:** The language depends on the request signature version: the legacy `v1` signing (the default) matches the old client and never sends `AcceptLanguage`, so `local_name` comes back in the local language; with `sign_version.default = "v4"` the shared EO gateway path sends `AcceptLanguage = en-US` and `local_name` comes back in English.
