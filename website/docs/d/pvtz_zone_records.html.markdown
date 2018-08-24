---
layout: "alicloud"
page_title: "Alicloud: alicloud_pvtz_zone_records"
sidebar_current: "docs-alicloud-datasource-pvtz-zone-records"
description: |-
    Provides a list of Private Zone Records which owned by an Alibaba Cloud account.
---

# alicloud\_pvtz\_zone\_records

This data source provides Private Zone Records resource information owned by an Alibaba Cloud account.

## Example Usage

```
data "alicloud_pvtz_zone_records" "records_ds" {
	zone_id = "${alicloud_pvtz_zone.basic.id}"
	keyword = "${alicloud_pvtz_zone_record.foo.value}"
}

output "first_record_id" {
  value = "${data.alicloud_pvtz_zone_records.records_ds.records.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `keyword` - (Optional) Keyword for record rr and value.
* `zone_id` - (Required) ID of the Private Zone.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `records` - A list of zone records. Each element contains the following attributes:
  * `id` - ID of the Private Zone Record.
  * `resource_record` - Resource record of the Private Zone Record.
  * `type` - Type of the Private Zone Record.
  * `value` - Value of the Private Zone Record.
  * `ttl` - Ttl of the Private Zone Record.
  * `priority` - Priority of the Private Zone Record.
 