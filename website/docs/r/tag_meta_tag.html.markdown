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

-> **NOTE:** Available since v1.209.0.

-> **NOTE:** Meta Tag Only Support `cn-hangzhou` Region

## Example Usage

Basic Usage

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_tag_meta_tag" "example" {
  key    = "Name1"
  values = ["Desc2"]
}
```

## Argument Reference

The following arguments are supported:

* `key` - (Required, ForceNew) The key of the tag meta tag. key must be 1 to 128 characters in length.
* `values` - (Required, ForceNew) The values of the tag meta tag. 
## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Meta Tag. It formats as `<regionId>`:`<key>`.

## Import

Tag Meta Tag can be imported using the id, e.g.

```shell
$ terraform import alicloud_tag_meta_tag.example <regionId>:<key>
```