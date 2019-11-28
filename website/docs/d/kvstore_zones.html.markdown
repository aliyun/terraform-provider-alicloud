---
layout: "alicloud"
page_title: "Alicloud: alicloud_kvstore_zones"
sidebar_current: "docs-alicloud-datasource-kvstore-zones"
description: |-
    Provides a list of kvstore availability zones that can be used by an Alibaba Cloud account.
---

# alicloud\_kvstore_zones

This data source provides kvstore availability zones that can be accessed by an Alibaba Cloud account within the region configured in the provider.

-> **NOTE:** Available in 1.63.0+.
-> **NOTE:** If one zone is sold out, it will not be exported.

## Example Usage

```
# Declare the data source
data "alicloud_kvstore_zones" "zones" {
  multi = true
  instance_charge_type = "PrePaid"
}

# Create an kvstore instance with the first matched zone
resource "alicloud_kvstore_instance" "instance" {
  availability_zone = "${data.alicloud_kvstore_zones.zones.ids.0}"

  # Other properties...
}
```

## Argument Reference

The following arguments are supported:

* `multi` - (Optional, type: bool) Indicate whether the zones can be used in a multi AZ configuration. Default to `false`. Multi AZ is usually used to launch kvstore instances.
* `instance_charge_type` - (Optional) Filter the results by a specific kvstore instance charge type. Valid values: `PrePaid` and `PostPaid`. Default to `PostPaid`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of zone IDs.
