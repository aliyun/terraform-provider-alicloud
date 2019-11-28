---
layout: "alicloud"
page_title: "Alicloud: alicloud_rds_zones"
sidebar_current: "docs-alicloud-datasource-rds-zones"
description: |-
    Provides a list of rds availability zones that can be used by an Alibaba Cloud account.
---

# alicloud\_rds_zones

This data source provides rds availability zones that can be accessed by an Alibaba Cloud account within the region configured in the provider.

-> **NOTE:** Available in 1.63.0+.
-> **NOTE:** If one zone is sold out, it will not be exported.

## Example Usage

```
# Declare the data source
data "alicloud_rds_zones" "zones" {
  multi = true
  instance_charge_type = "PrePaid"
}

# Create an DB instance with the first matched zone
resource "alicloud_db_instance" "instance" {
  zone_id = "${data.alicloud_rds_zones.zones.ids.0}"

  # Other properties...
}
```

## Argument Reference

The following arguments are supported:

* `multi` - (Optional, type: bool) Indicate whether the zones can be used in a multi AZ configuration. Default to `false`. Multi AZ is usually used to launch RDS instances.
* `instance_charge_type` - (Optional) Filter the results by a specific rds instance charge type. Valid values: `PrePaid` and `PostPaid`. Default to `PostPaid`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of zone IDs.
