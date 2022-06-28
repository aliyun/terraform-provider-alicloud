---
subcategory: "Elastic Container Instance (ECI)"
layout: "alicloud"
page_title: "Alicloud: alicloud_eci_zones"
sidebar_current: "docs-alicloud-datasource-eci-zones"
description: |-
  Provides a list of ECI available zones to the user.
---

# alicloud\_eci\_zones

This data source provides the available zones with the Application Load Balancer (ALB) Instance of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.145.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_eci_zones" "default" {}

output "first_eci_zones_id" {
  value = data.alicloud_eci_zones.default.zones.0.zone_ids.0
}
```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `zones` - A list of eci Instance zones. Each element contains the following attributes:
	* `zone_ids` - The list of available zone ids.
	* `region_endpoint` - The endpoint of the region.
