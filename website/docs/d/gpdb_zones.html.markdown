---
layout: "alicloud"
page_title: "Alicloud: alicloud_gpdb_zones"
sidebar_current: "docs-alicloud-datasource-gpdb-zones"
description: |-
    Provides a list of gpdb availability zones that can be used by an Alibaba Cloud account.
---

# alicloud\_gpdb_zones

This data source provides gpdb availability zones that can be accessed by an Alibaba Cloud account within the region configured in the provider.

-> **NOTE:** Available in 1.63.0+.

## Example Usage

```
# Declare the data source
data "alicloud_gpdb_zones" "zones" {
  multi = true
  instance_charge_type = "PrePaid"
}
```

## Argument Reference

The following arguments are supported:

* `multi` - (Optional, type: bool) Indicate whether the zones can be used in a multi AZ configuration. Default to `false`. Multi AZ is usually used to launch gpdb instances.
* `instance_charge_type` - (Optional) Filter the results by a specific gpdb instance charge type. Valid values: `PrePaid` and `PostPaid`. Default to `PostPaid`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of zone IDs.
