---
subcategory: "Private Zone"
layout: "alicloud"
page_title: "Alicloud: alicloud_pvtz_zone_record"
sidebar_current: "docs-alicloud-resource-pvtz-zone-record"
description: |-
  Provides a Alicloud Private Zone Record resource.
---

# alicloud\_pvtz\_zone\_record

Provides a Private Zone Record resource.

-> **NOTE:** Terraform will auto Create a Private Zone Record while it uses `alicloud_pvtz_zone_record` to build a Private Zone Record resource.

## Example Usage

Basic Usage

```
resource "alicloud_pvtz_zone" "zone" {
  name = "foo.test.com"
}

resource "alicloud_pvtz_zone_record" "foo" {
  zone_id         = alicloud_pvtz_zone.zone.id
  resource_record = "www"
  type            = "CNAME"
  value           = "bbb.test.com"
  ttl             = 60
}
```
## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, ForceNew) The name of the Private Zone Record.
* `resource_record` - (Required, ForceNew) The resource record of the Private Zone Record.
* `type` - (Required) The type of the Private Zone Record. Valid values: A, CNAME, TXT, MX, PTR.
* `value` - (Required) The value of the Private Zone Record.
* `ttl` - (Optional) The ttl of the Private Zone Record.
* `priority` - (Optional) The priority of the Private Zone Record. At present, only can "MX" record support it. Valid values: [1-50]. Default to 1.
* `remark` - (Optional, Available in 1.103.2+) The remark of the Private Zone Record.

## Attributes Reference

The following attributes are exported:

* `id` - This ID of this resource. The value is formate as `<record_id>:<zone_id>`.
* `record_id` - The Private Zone Record ID.

## Import

Private Zone Record can be imported using the id, e.g.

```
$ terraform import alicloud_pvtz_zone_record.example abc123456
```

