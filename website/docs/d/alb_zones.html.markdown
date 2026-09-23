---
subcategory: "Application Load Balancer (ALB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_alb_zones"
sidebar_current: "docs-alicloud-datasource-alb-zones"
description: |-
  Provides a list of Application Load Balancer (ALB) instance available zones to the user.
---

# alicloud\_alb\_zones

This data source provides the available zones with the Application Load Balancer (ALB) Instance of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.132.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_alb_zones" "example" {}

output "first_alb_zones_id" {
  value = data.alicloud_alb_zones.example.zones.0.zone_id
}
```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of alb instance zone IDs.
* `zones` - A list of alb Instance zones. Each element contains the following attributes:
	* `id` - The ID of zone.
	* `zone_id` - The zone ID.
	* `local_name` - The local name.
