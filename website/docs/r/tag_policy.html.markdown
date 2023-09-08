---
subcategory: "TAG"
layout: "alicloud"
page_title: "Alicloud: alicloud_tag_policy"
sidebar_current: "docs-alicloud-resource-tag-policy"
description: |-
  Provides a Alicloud Tag Policy resource.
---

# alicloud_tag_policy

Provides a Tag Policy resource.

For information about Tag Policy and how to use it,
see [What is Policy](https://www.alibabacloud.com/help/en/resource-management/latest/create-policy).

-> **NOTE:** Available since v1.203.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}
provider "alicloud" {
  region = "cn-shanghai"
}
resource "alicloud_tag_policy" "example" {
  policy_name    = var.name
  policy_desc    = var.name
  user_type      = "USER"
  policy_content = <<EOF
		{"tags":{"CostCenter":{"tag_value":{"@@assign":["Beijing","Shanghai"]},"tag_key":{"@@assign":"CostCenter"}}}}
    EOF
}
```

## Argument Reference

The following arguments are supported:

* `policy_name` - (Required) The name of the policy. name must be 1 to 128 characters in length and can contain letters,
  digits, and hyphens (-).
* `policy_content` - (Required) The content of the policy.
* `policy_desc` - (Optional) The description of the policy. The description must be 1 to 512 characters in length.
* `user_type` - (Optional, ForceNew) The type of the tag policy. Valid values: `USER`, `RD`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of tag policy.

## Import

Tag Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_tag_policy.example <id>
```
