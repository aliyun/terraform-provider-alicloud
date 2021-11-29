---
subcategory: "ApsaraDB for MyBase"
layout: "alicloud"
page_title: "Alicloud: alicloud_cddc_zones"
sidebar_current: "docs-alicloud-datasource-cddc-zones"
description: |-
  Provides a list of Cddc Zones to the user.
---

# alicloud\_cddc\_zones

This data source provides the Cddc Zones of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.147.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_cddc_zones" "example" {}
output "cddc_zones_id_1" {
  value = data.alicloud_cddc_zones.example.zones.0.id
}

```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `zones` - A list of Cddc Zones. Each element contains the following attributes:
	* `region_id` - The ID of the region.
	* `zone_id` - The ID of the zone.
	* `id` - The ID of the zone.
	