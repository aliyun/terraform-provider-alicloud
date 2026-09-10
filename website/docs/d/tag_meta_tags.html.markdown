---
subcategory: "TAG"
layout: "alicloud"
page_title: "Alicloud: alicloud_tag_meta_tags"
sidebar_current: "docs-alicloud-datasource-tag-meta-tags"
description: |-
  Provides a list of Tag Meta Tags to the user.
---

# alicloud\_tag\_meta\_tags

This data source provides the Tag Meta Tags of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.169.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_tag_meta_tags" "default" {
  key_name = "example_value"
}
output "tag_meta_tag_default_1" {
  value = data.alicloud_tag_meta_tags.default.tags.value_name
}
```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `key_name` - (Optional, ForceNew) The name of the key.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `tags` - A list of Meta Tags. Each element contains the following attributes:
  * `value_name` - The name of the value.
  * `key_name` - The name of the key.
  * `category` - The type of the resource tags.