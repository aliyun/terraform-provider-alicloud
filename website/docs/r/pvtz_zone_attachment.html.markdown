---
layout: "alicloud"
page_title: "Alicloud: alicloud_pvtz_zone_attachment"
sidebar_current: "docs-alicloud-resource-pvtz-zone-attachment"
description: |-
  Provides vpcs bound to Alicloud Private Zone resource.
---

# alicloud\_pvtz\_zone\_attachment

Provides vpcs bound to Alicloud Private Zone resource.

~> **NOTE:** Terraform will auto bind vpc to a Private Zone while it uses `alicloud_pvtz_zone_attachment` to build a Private Zone and VPC binding resource.

## Example Usage

Basic Usage

```
resource "alicloud_pvtz_zone" "zone" {
	name = "foo.test.com"
}

resource "alicloud_vpc" "vpc" {
	name = "tf_test_foo"
	cidr_block = "172.16.0.0/12"
}

resource "alicloud_pvtz_zone_attachment" "zone-attachment" {
	zone_id = "${alicloud_pvtz_zone.zone.id}"
	vpc_ids = ["${alicloud_vpc.vpc.id}"]
}
```
## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, Forces new resource) The name of the Private Zone Record.
* `vpc_ids` - (Required) The id List of the VPC, for example:["vpc-1","vpc-2"].

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Private Zone VPC Attachment.
