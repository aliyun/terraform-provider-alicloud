---
layout: "alicloud"
page_title: "Alicloud: alicloud_pvtz_zone"
sidebar_current: "docs-alicloud-resource-pvtz-zone"
description: |-
  Provides a Alicloud Private Zone resource.
---

# alicloud\_pvtz\_zone

Provides a Private Zone resource.

~> **NOTE:** Terraform will auto Create a Private Zone while it uses `alicloud_pvtz_zone` to build a Private Zone resource.

## Example Usage

Basic Usage

```
resource "alicloud_pvtz_zone" "foo" {
	name = "foo.test.com"
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required, Forces new resource) The name of the Private Zone.
* `remark` - (Optional) The remark of the Private Zone.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Private Zone.
* `record_count` - The count of the Private Zone Record.

## Import

Private Zone can be imported using the id, e.g.

```
$ terraform import alicloud_pvtz_zone.example abc123456
```

