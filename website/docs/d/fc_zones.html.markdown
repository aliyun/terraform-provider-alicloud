---
layout: "alicloud"
page_title: "Alicloud: alicloud_fc_zones"
sidebar_current: "docs-alicloud-datasource-fc-zones"
description: |-
    Provides a list of fc availability zones that can be used by an Alibaba Cloud account.
---

# alicloud\_fc_zones

This data source provides fc availability zones that can be accessed by an Alibaba Cloud account within the region configured in the provider.

-> **NOTE:** Available in 1.63.0+.

## Example Usage

```
# Declare the data source
data "alicloud_fc_zones" "zones" {
  multi = true
  instance_charge_type = "PrePaid"
}
```

## Argument Reference

The following arguments are supported:

* `multi` - (Optional, type: bool) Indicate whether the zones can be used in a multi AZ configuration. Default to `false`. Multi AZ is usually used to launch FC instances.
* `instance_charge_type` - (Optional) Filter the results by a specific fc instance charge type. Valid values: `PrePaid` and `PostPaid`. Default to `PostPaid`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of zone IDs.
* `zones` - A list of availability zones. Each element contains the following attributes:
  * `id` - ID of the zone.
  * `multi_zone_ids` - A list of zone ids in which the multi zone.