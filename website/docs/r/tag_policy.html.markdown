---
subcategory: "Tag Policy"
layout: "alicloud"
page_title: "Alicloud: alicloud_tag_policy"
sidebar_current: "docs-alicloud-resource-tag-policy"
description: |-
Provides a Alicloud Tag Policy resource.
---

# alicloud\_tag\_policy

Provides a Tag Policy resource.  
For information about Tag Policy and how to use it.

## Example Usage

Basic Usage

```terraform
resource "alicloud_tag_policy" "example" {
  policy_name     = "abc12345"
  policy_document = <<EOF
		{"tags":{"CostCenter":{"tag_value":{"@@assign":["Beijing","Shanghai"]},"tag_key":{"@@assign":"CostCenter"}}}}
    EOF
}
```
## Argument Reference

The following arguments are supported:

* `policy_name` - (Required) The name of the policy. name must be 1 to 128 characters in length and can contain letters, digits, and hyphens (-).
* `policy_content` - (Required) The content of the policy. 
* `policy_desc` - (Optional) The description of the policy. The description must be 1 to 512 characters in length.
* `user_type` - (Optional) The type of the tag policy. Valid values: `USER`, `RD`.
* `dry_run` - (Optional) Whether to execute


## Attributes Reference

* `id` - The resource ID of tag policy. 
* `user_type` - The type of the tag policy. Valid values: `USER`, `RD`.

## Import

Tag Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_tag_policy.example abc12345
```
