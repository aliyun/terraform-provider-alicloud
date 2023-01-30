---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_rds_collation_time_zones"
sidebar_current: "docs-alicloud-rds-collation-time-zones"
description: |-
  Operation to query the character set collations and time zones available for use in ApsaraDB RDS.
---

# alicloud\_rds\_collation\_time\_zones

Operation to query the character set collations and time zones available for use in ApsaraDB RDS.

-> **NOTE:** Available in v1.198.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_rds_collation_time_zones" "zones" {
  output_file = "./classes.txt"
}

output "first_rds_collation_time_zones" {
  value = data.alicloud_rds_collation_time_zones.zones.collation_time_zones.0
}
```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional) File name where to save data source results (after running `terraform apply`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of time zones.
* `collation_time_zones` - An array that consists of the character set collations and time zones that are available for
  use in ApsaraDB RDS.
    * `description` - The code of the instance type.
    * `standard_time_offset` - The offset of the UTC time. The offset is in the following format: (UTC+<i>HH:mm</i>).
    * `time_zone` - The time zone that is available for use in ApsaraDB RDS.