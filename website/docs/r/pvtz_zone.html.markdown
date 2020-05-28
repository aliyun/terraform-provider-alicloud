---
subcategory: "Private Zone"
layout: "alicloud"
page_title: "Alicloud: alicloud_pvtz_zone"
sidebar_current: "docs-alicloud-resource-pvtz-zone"
description: |-
  Provides a Alicloud Private Zone Zone resource.
---

# alicloud\_pvtz\_zone

Provides a Private Zone Zone resource.
For information about Private Zone Zone and how to use it, see [What is Private Zone Zone](https://www.alibabacloud.com/help/en/doc-detail/66240.htm).

-> **NOTE:** Available in v1.85.0+.

-> **NOTE:** Terraform will auto Create a Private Zone while it uses `alicloud_pvtz_zone` to build a Private Zone resource.

## Example Usage

Basic Usage

```
resource "alicloud_pvtz_zone" "example" {
  zone_name="zone.com"
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Optional, ForceNew) Attribute `name` has been deprecated from version 1.85.0. Use resource alicloud_private_zone_zone's `zone_name` instead.
* `zone_name` - (Optional, ForceNew) The zone_name of the Private Zone.
* `remark` - (Optional) The remark of the Private Zone.
* `proxy_pattern` - (Optional, Available in 1.69.0+) The recursive DNS proxy. Valid values:
    - ZONE: indicates that the recursive DNS proxy is disabled.
    - RECORD: indicates that the recursive DNS proxy is enabled. Default to "ZONE".
* `user_client_ip` - (Optional, Available in 1.69.0+) The IP address of the client.
* `lang` - (Optional, Available in 1.69.0+) The language. Valid values: "zh", "en", "jp".

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Private Zone.

## Import

Private Zone can be imported using the id, e.g.

```
$ terraform import alicloud_pvtz_zone.example abc123456
```

