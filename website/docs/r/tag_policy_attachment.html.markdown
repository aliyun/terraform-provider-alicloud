---
subcategory: "Tag"
layout: "alicloud"
page_title: "Alicloud: alicloud_tag_policy_attachment"
sidebar_current: "docs-alicloud-resource-tag-policy-attachment"
description: |-
  Provides a Alicloud Tag Policy Attachment resource.
---

# alicloud\_tag\_policy\_attachment

Provides a Tag Policy Attachment resource to attaches a policy to an object. After you attach a tag policy to an object.
For information about Tag Policy Attachment and how to use it,
see [What is Policy Attachment](https://www.alibabacloud.com/help/en/resource-management/latest/attach-policy).

-> **NOTE:** Available in v1.204.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_tag_policy" "example" {
  policy_name     = "testName"
  policy_desc     = "testDesc"
  user_type       = "USER"
  policy_document = <<EOF
		{"tags":{"CostCenter":{"tag_value":{"@@assign":["Beijing","Shanghai"]},"tag_key":{"@@assign":"CostCenter"}}}}
    EOF
}

resource "alicloud_tag_policy_attachment" "example" {
  policy_id   = alicloud_tag_policy.example.id
  target_id   = "151266687691****"
  target_type = "USER"
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, ForceNew) The ID of the tag policy.
* `target_id` - (Required, ForceNew) The ID of the object.
* `target_type` - (Required, ForceNew) The type of the object. Valid values: `USER`, `ROOT`, `FOLDER`, `ACCOUNT`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Policy Attachment. It formats as `<policy_id>`:`<target_id>`:`<target_type>`.

## Import

Tag Policy Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_tag_policy_attachment.example <policy_id>:<target_id>:<target_type>
```