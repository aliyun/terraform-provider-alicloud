---
subcategory: "Private Zone"
layout: "alicloud"
page_title: "Alicloud: alicloud_pvtz_resolver_zones"
sidebar_current: "docs-alicloud-datasource-pvtz-resolver-zones"
description: |-
  Provides a list of Private Zone Resolver available zones to the user.
---

# alicloud\_pvtz\_resolver\_zones

This data source provides the available zones with the Private Zone Resolver of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.143.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_pvtz_resolver_zones" "default" {
  status = "NORMAL"
}

output "first_zones_id" {
  value = data.alicloud_pvtz_resolver_zones.default.zones.0.zone_id
}
```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the Zone. Valid values: `NORMAL`, `SOLD_OUT`.


## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `zones` - A list of Private Zone Resolver zones. Each element contains the following attributes:
	* `zone_id` - The zone ID.
	* `status` - The status of the Zone.
