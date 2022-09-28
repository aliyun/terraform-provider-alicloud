---
subcategory: "Click House"
layout: "alicloud"
page_title: "Alicloud: alicloud_click_house_regions"
sidebar_current: "docs-alicloud-datasource-click-house-regions"
description: |-
  Provides a list of Click House Accounts to the user.
---

# alicloud\_click\_house\_regions

This data source provides the Click House Accounts of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.138.0+.

## Example Usage

Basic Usage

```terraform
data alicloud_click_house_regions "default1" {
  current = true
}
data alicloud_click_house_regions "default2" {
  region_id = "cn-hangzhou"
}

```

## Argument Reference

The following arguments are supported:

* `region_id` - (Option) You can use specified region_id to find the region and available zones information that supports ClickHouse.
* `current` - (Optional) Set to true to match only the region configured in the provider. Default value: `true`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `regions` - A list of Click House Regions. Each element contains the following attributes:
    * `region_id` - The Region ID.
    * `zone_ids` -  A list of available zone ids in the region_id.
      * `vpc_enabled` - Whether to support vpc network.
      * `zone_id` - The zone ID.
