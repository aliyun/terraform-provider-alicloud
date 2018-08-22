---
layout: "alicloud"
page_title: "Alicloud: alicloud_pvtz_zone_records"
sidebar_current: "docs-alicloud-datasource-pvtz-zone-records"
description: |-
    Provides a list of Private Zone Records which owned by an Alicloud account.
---

# alicloud\_pvtz\_zone\_records

The Private Zone Records data source lists a number of Private Zone Records resource information owned by an Alicloud account.

## Example Usage

```
data "alicloud_pvtz_zone_records" "keyword" {
	zone_id = "${alicloud_pvtz_zone.basic.id}"
	keyword = "${alicloud_pvtz_zone_record.foo.value}"
}

```

## Argument Reference

The following arguments are supported:

* `keyword` - (Optional) Keyword for record rr and value.
* `zone_id` - (Required) ID of The Private Zone.
* `output_file` - (Optional) The name of file that can save vpcs data source after running `terraform plan`.

## Attributes Reference

The following attributes are exported:

* `id` - ID of the Private Zone Record.
* `resource_record` - Resource record of the Private Zone Record.
* `type` - Type of the Private Zone Record.
* `value` - Value of the Private Zone Record.
* `ttl` - Ttl of the Private Zone Record.
* `priority` - Priority of the Private Zone Record.