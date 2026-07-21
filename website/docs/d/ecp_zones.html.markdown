---
subcategory: "Elastic Cloud Phone (ECP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecp_zones"
sidebar_current: "docs-alicloud-datasource-ecp-zones"
description: |-
  Provides a list of Ecp available zones to the user.
---

# alicloud\_ecp\_zones

This data source provides the available zones with the Cloud Phone (ECP) Instance of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.158.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ecp_zones" "default" {}

output "alicloud_nas_zones_id" {
  value = "${data.alicloud_ecp_zones.default.zones.0.zone_id}"
}
```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `zones` - A list of availability zone information collection.
	* `zone_id` - String to filter results by zone id.
