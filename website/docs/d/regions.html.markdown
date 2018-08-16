---
layout: "alicloud"
page_title: "Alicloud: alicloud_regions"
sidebar_current: "docs-alicloud-datasource-regions"
description: |-
    Provides a list of Availability Regions which can be used by an Alicloud account.
---

# alicloud\_regions

The Regions data source allows access to the list of Alicloud Regions.

## Example Usage

```
data "alicloud_regions" "current" {
	current = true
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The full name of the region to select.
* `current` - (Optional) Set to true to match only the region configured in the provider.
* `output_file` - (Optional) The name of file that can save regions data source after running `terraform plan`.

~> **NOTE:** You will get an error if you set `current` to true and `name` to a different value from the one you configured in the provider.
 It is better to either use `name` or `current`, but not both at the same time.

## Attributes Reference

A list of regions will be exported and its every element contains the following attributes:

* `id` - ID of the region.
* `local_name` - Name of the region in the local language.
