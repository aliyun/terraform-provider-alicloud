---
layout: "alicloud"
page_title: "Alicloud: alicloud_pvtz_zone_record"
sidebar_current: "docs-alicloud-resource-pvtz-zone-record"
description: |-
  Provides a Alicloud Private Zone Record resource.
---

# alicloud\_pvtz\_zone\_record

Provides a Private Zone Record resource.

~> **NOTE:** Terraform will auto Create a Private Zone Record while it uses `alicloud_pvtz_zone_record` to build a Private Zone Record resource.

## Example Usage

Basic Usage

```
resource "alicloud_pvtz_zone" "zone" {
	name = "foo.test.com"
}

resource "alicloud_pvtz_zone_record" "foo" {
	zone_id = "${alicloud_pvtz_zone.zone.id}"
	resource_record = "www"
	type = "CNAME"
	value = "bbb.test.com"
	ttl="60
	priority = "6"
}
```
## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, Forces new resource) The name of the Private Zone Record.
* `resource_record` - (Required) The resource record of the Private Zone Record.
* `type` - (Required) The type of the Private Zone Record.
* `value` - (Required) The value of the Private Zone Record.
* `status` - (Optional) The status of the Private Zone Record.
* `ttl` - (Optional) The ttl of the Private Zone Record.
* `priority` - (Optional) The priority of the Private Zone Record.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Private Zone Record.

## Import

Private Zone Record can be imported using the id, e.g.

```
$ terraform import alicloud_pvtz_zone_record.example abc123456
```

