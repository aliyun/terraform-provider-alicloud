---
subcategory: "TAG"
layout: "alicloud"
page_title: "Alicloud: alicloud_tag_meta_tag"
sidebar_current: "docs-alicloud-resource-tag-meta-tag"
description: |-
  Provides a Alicloud Tag Meta Tag resource.
---

# alicloud\_tag\_meta\_tag

Provides a Tag Meta Tag resource.

For information about Tag Meta Tag and how to use it,
see [What is Meta Tag](https://www.alibabacloud.com/help/en/resource-management/latest/createtags).

-> **NOTE:** Available since v1.208.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_tag_meta_tag" "example" {
  key_name   = "Name"
  value_name = "Desc"
}
```

## Argument Reference

The following arguments are supported:

* `key_name` - (Required, ForceNew) The key of the tag meta tag. key_name must be 1 to 128 characters in length.
* `value_name` - (Required, ForceNew) The value of the tag meta tag. value_name must be 1 to 128 characters in length.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Meta Tag. It formats as `<key_name>`:`<value_name>`.

## Import

Tag Meta Tag can be imported using the id, e.g.

```shell
$ terraform import alicloud_tag_meta_tag.example <key_name>:<value_name>
```
