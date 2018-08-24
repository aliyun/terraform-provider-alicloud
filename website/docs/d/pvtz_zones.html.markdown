---
layout: "alicloud"
page_title: "Alicloud: alicloud_pvtz_zones"
sidebar_current: "docs-alicloud-datasource-pvtz-zones"
description: |-
    Provides a list of Private Zones which owned by an Alibaba Cloud account.
---

# alicloud\_pvtz\_zones

This data source lists a number of Private Zones resource information owned by an Alibaba Cloud account.

## Example Usage

```
data "alicloud_pvtz_zones" "pvtz_zones_ds" {
	keyword = "${alicloud_pvtz_zone.basic.zone_name}"
}

output "first_zone_id" {
  value = "${data.alicloud_pvtz_zones.pvtz_zones_ds.zones.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `keyword` - (Optional) keyword for zone name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `zones` - A list of zones. Each element contains the following attributes:
  * `id` - ID of the Private Zone.
  * `remark` - Remark of the Private Zone.
  * `record_count` - Count of the Private Zone Record.
  * `name` - Name of the Private Zone.
  * `is_ptr` - Whether the Private Zone is ptr
  * `creation_time` - Time of creation of the Private Zone.
  * `update_time` - Time of update of the Private Zone.
  * `bind_vpcs` - List of the VPCs is bound to the Private Zone.
  