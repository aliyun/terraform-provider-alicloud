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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_tag_meta_tag&exampleId=a08885db-2c17-c64b-d934-398212575530a5fa66e1&activeTab=example&spm=docs.r.tag_meta_tag.0.a08885db2c&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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