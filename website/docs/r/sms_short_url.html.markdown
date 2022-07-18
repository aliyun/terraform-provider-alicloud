---
subcategory: "Short Message Service (SMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sms_short_url"
sidebar_current: "docs-alicloud-resource-sms-short-url"
description: |-
  Provides a Alicloud SMS Short Url resource.
---

# alicloud\_sms\_short\_url

Provides a SMS Short Url resource.

For information about SMS Short Url and how to use it, see [What is Short Url](https://help.aliyun.com/document_detail/419291.html).

-> **NOTE:** Available in v1.178.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_sms_short_url" "example" {
  effective_days = 30
  short_url_name = "example_value"
  source_url     = "example_value"
}
```

## Argument Reference

The following arguments are supported:

* `effective_days` - (Required, ForceNew) Short chain service use validity period. Valid values: `30`, `60`, `90`. The unit is days, and the maximum validity period is 90 days.
* `short_url_name` - (Required, ForceNew) The name of the resource.
* `source_url` - (Required, ForceNew) The original link address.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Short Url.
* `status` - Short chain status.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 2 mins) Used when create the Short Url.
* `delete` - (Defaults to 1 mins) Used when delete the Short Url.

## Import

SMS Short Url can be imported using the id, e.g.

```
$ terraform import alicloud_sms_short_url.example <id>
```