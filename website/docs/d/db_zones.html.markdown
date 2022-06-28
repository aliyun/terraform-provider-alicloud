---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_db_zones"
sidebar_current: "docs-alicloud-datasource-db-zones"
description: |-
    Provides a list of availability zones for RDS that can be used by an Alibaba Cloud account.
---

# alicloud\_db\_zones

This data source provides availability zones for RDS that can be accessed by an Alibaba Cloud account within the region configured in the provider.

-> **NOTE:** Available in v1.73.0+.

## Example Usage

```
# Declare the data source
data "alicloud_db_zones" "zones_ids" {}

# Create an RDS instance with the first matched zone
resource "alicloud_db_instance" "db" {
    zone_id = ${data.alicloud_db_zones.zones_ids.zones.0.id}

  # Other properties...
}
```

## Argument Reference

The following arguments are supported:

* `multi` - (Deprecated) It has been deprecated from version 1.137.0 and using `multi_zone` instead.
* `multi_zone` - (Optional, Available in 1.137.0+) Indicate whether the zones can be used in a multi AZ configuration. Default to `false`. Multi AZ is usually used to launch RDS instances.
* `instance_charge_type` - (Optional) Filter the results by a specific instance charge type. Valid values: `PrePaid` and `PostPaid`. Default to `PostPaid`.
* `engine` - (Optional, Available in 1.134.0+) Database type. Valid values: "MySQL", "SQLServer", "PostgreSQL", "PPAS", "MariaDB". If not set, it will match all of engines.
* `engine_version` - (Optional, Available in 1.134.0+) Database version required by the user. Value options can refer to the latest docs [detail info](https://www.alibabacloud.com/help/doc-detail/26228.htm) `EngineVersion`.
* `category` - (Optional, Available in 1.134.0+) DB Instance category. the value like [`Basic`, `HighAvailability`, `Finance`, `AlwaysOn`], [detail info](https://www.alibabacloud.com/help/doc-detail/69795.htm).
* `db_instance_storage_type` - (Optional, Available in 1.134.0+) The DB instance storage space required by the user. Valid values: "cloud_ssd", "local_ssd", "cloud_essd", "cloud_essd2", "cloud_essd3".
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of zone IDs.
* `zones` - A list of availability zones. Each element contains the following attributes:
  * `id` - ID of the zone.
  * `multi_zone_ids` - A list of zone ids in which the multi zone.
