---
layout: "alicloud"
page_title: "Alicloud: alicloud_pvtz_zones"
sidebar_current: "docs-alicloud-datasource-pvtz-zones"
description: |-
    Provides a list of Private Zones which owned by an Alicloud account.
---

# alicloud\_pvtz\_zones

The Private Zones data source lists a number of Private Zones resource information owned by an Alicloud account.

## Example Usage

```
data "alicloud_pvtz_zones" "keyword" {
	keyword = "${alicloud_pvtz_zone.basic.zone_name}"
}

```

## Argument Reference

The following arguments are supported:

* `keyword` - (Optional) keyword for zone name.
* `output_file` - (Optional) The name of file that can save vpcs data source after running `terraform plan`.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the Private Zone.
* `remark` - Remark of the Private Zone.
* `record_count` - Count of the Private Zone Record.
* `name` - Name of the Private Zone.
* `is_ptr` - Whether the Private Zone is ptr
* `creation_time` - Time of creation of the Private Zone.
* `update_time` - Time of update of the Private Zone.
* `bind_vpcs` - List of the VPCs is bound to the Private Zone.