---
subcategory: "Elastic Desktop Service(EDS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecd_zones"
sidebar_current: "docs-alicloud-datasource-ecd-zones"
description: |-
  Provides a list of Ecd available zones to the user.
---

# alicloud\_ecd\_zones

This data source provides the available zones with the Elastic Desktop Service(EDS) of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.174.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ecd_zones" "default" {}

output "alicloud_ecd_zones" {
  value = "${data.alicloud_ecd_zones.default.zones.0.zone_id}"
}
```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `zones` - A list of availability zone information collection.
	* `zone_id` - String to filter results by zone id.
