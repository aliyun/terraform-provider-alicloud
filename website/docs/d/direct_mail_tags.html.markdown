---
subcategory: "Direct Mail"
layout: "alicloud"
page_title: "Alicloud: alicloud_direct_mail_tags"
sidebar_current: "docs-alicloud-datasource-direct-mail-tags"
description: |-
  Provides a list of Direct Mail Tags to the user.
---

# alicloud\_direct\_mail\_tags

This data source provides the Direct Mail Tags of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.144.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_direct_mail_tags" "ids" {
  ids = ["example_id"]
}
output "direct_mail_tag_id_1" {
  value = data.alicloud_direct_mail_tags.ids.tags.0.id
}

data "alicloud_direct_mail_tags" "nameRegex" {
  name_regex = "^my-Tag"
}
output "direct_mail_tag_id_2" {
  value = data.alicloud_direct_mail_tags.nameRegex.tags.0.id
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Tag IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Tag name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Tag names.
* `tags` -  A list of Direct Mail Tags. Each element contains the following attributes:
    * `id` - The ID of the tag.
    * `tag_id` - The ID of the tag.
    * `tag_name` - The name of the tag.