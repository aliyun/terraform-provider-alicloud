---
subcategory: "CDN"
layout: "alicloud"
page_title: "Alicloud: alicloud_cdn_blocked_regions"
sidebar_current: "docs-alicloud-datasource-cdn-blocked-regions"
description: |-
  Provides a list of Cdn Blocked Regions to the user.
---

# alicloud\_cdn\_blocked\_regions

This data source provides the Cdn blocked regions.

-> **NOTE:** Available in v1.173.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_cdn_blocked_regions" "example" {
  language = "zh"
}
```

## Argument Reference

The following arguments are supported:

* `language` - (Required, ForceNew) The language. Valid values: `zh`, `en`, `jp`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `regions` - A list of Cdn Blocked Regions. Each element contains the following attributes:
  * `continent` - The region to which the country belongs.
  * `countries_and_regions` - National region abbreviation.
  * `countries_and_regions_name` - The name of the country and region.