---
subcategory: "Direct Mail"
layout: "alicloud"
page_title: "Alicloud: alicloud_direct_mail_receiverses"
sidebar_current: "docs-alicloud-datasource-direct-mail-receiverses"
description: |-
  Provides a list of Direct Mail Receiverses to the user.
---

# alicloud\_direct\_mail\_receiverses

This data source provides the Direct Mail Receiverses of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.125.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_direct_mail_receiverses" "example" {
  ids        = ["ca73b1e4fb0df7c935a5097a****"]
  name_regex = "the_resource_name"
}

output "first_direct_mail_receivers_id" {
  value = data.alicloud_direct_mail_receiverses.example.receiverses.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Receivers IDs.
* `key_word` - (Optional, ForceNew) The key word.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Receivers name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the resource. Valid Values: `0` means uploading, `1` means upload completed. 

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Receivers names.
* `receiverses` - A list of Direct Mail Receiverses. Each element contains the following attributes:
	* `create_time` - The creation time of the resource.
	* `description` - The description.
	* `id` - The ID of the Receivers.
	* `receivers_alias` -The Receivers Alias.
	* `receivers_id` - The first ID of the resource.
	* `receivers_name` - The name of the resource.
	* `status` - The status of the resource.
