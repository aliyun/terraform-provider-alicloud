---
layout: "alicloud"
page_title: "Alicloud: alicloud_slb_zones"
sidebar_current: "docs-alicloud-datasource-slb-zones"
description: |-
    Provides a list of slb availability zones that can be used by an Alibaba Cloud account.
---

# alicloud\_slb_zones

This data source provides slb availability zones that can be accessed by an Alibaba Cloud account within the region configured in the provider.

-> **NOTE:** Available in 1.63.0+.
-> **NOTE:** If one zone is sold out, it will not be exported.

## Example Usage

```
# Declare the data source
data "alicloud_slb_zones" "zones" {
  multi = true
  instance_charge_type = "PrePaid"
}
```

## Argument Reference

The following arguments are supported:

* `instance_charge_type` - (Optional) Filter the results by a specific slb instance charge type. Valid values: `PrePaid` and `PostPaid`. Default to `PostPaid`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of zone IDs.
* `zones` - A list of availability zones. Each element contains the following attributes:
  * `zone_id` - ID of the zone.
  * `slave_zone_ids` - A list of slave zone ids.