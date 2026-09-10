---
subcategory: "Network Load Balancer (NLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nlb_zones"
sidebar_current: "docs-alicloud-datasource-nlb-zones"
description: |-
  Provides a list of Network Load Balancer (NLB) instance available zones to the user.
---

# alicloud\_nlb\_zones

This data source provides the available zones with the Network Load Balancer (NLB) Instance of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.191.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_nlb_zones" "example" {}

output "first_nlb_zones_id" {
  value = data.alicloud_nlb_zones.example.zones.0.zone_id
}
```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of nlb instance zone IDs.
* `zones` - A list of nlb Instance zones. Each element contains the following attributes:
	* `id` - The ID of zone.
	* `zone_id` - The zone ID.
	* `local_name` - The local name.
