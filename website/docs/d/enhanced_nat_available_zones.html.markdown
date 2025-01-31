---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_enhanced_nat_available_zones"
sidebar_current: "docs-alicloud-datasource-enhanced-nat-available-zones"
description: |-
    Provides a list of available zones by the enhanced Nat Gateway.
---

# alicloud_enhanced_nat_available_zones

This data source provides a list of available zones by the enhanced Nat Gateway.

-> **NOTE:** Available since 1.102.0+.

## Example Usage

```terraform
data "alicloud_enhanced_nat_available_zones" "default" {
}

output "zones" {
  value = data.alicloud_enhanced_nat_available_zones.default.zones.0.zone_id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of alarm IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - (Optional) A list of available zones IDs by the enhanced NAT gateway.
* `zones` - A list of available zones. Each element contains the following attributes:
  * `zone_id` - The ID of the available zone.
  * `local_name` - Name of the available zone.

