---
subcategory: "Direct Mail"
layout: "alicloud"
page_title: "Alicloud: alicloud_direct_mail_tag"
sidebar_current: "docs-alicloud-resource-direct-mail-tag"
description: |-
  Provides a Alicloud Direct Mail Tag resource.
---

# alicloud\_direct\_mail\_tag

Provides a Direct Mail Tag resource.

For information about Direct Mail Tag and how to use it, see [What is Tag](https://www.alibabacloud.com/help/zh/doc-detail/119007.htm).

-> **NOTE:** Available in v1.144.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_direct_mail_tag" "example" {
  tag_name = "example_value"
}

```

## Argument Reference

The following arguments are supported:

* `tag_name` - (Required) The name of the tag. The name must be `1` to `50` characters in length, and can contain letters and digits.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Tag.

## Import

Direct Mail Tag can be imported using the id, e.g.

```
$ terraform import alicloud_direct_mail_tag.example <id>
```