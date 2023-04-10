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
see [What is Attach Policy](https://www.alibabacloud.com/help/en/resource-management/latest/attach-policy).

-> **NOTE:** Available in v1.203.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_tag_policy_attachment" "example" {
  policy_id   = "p-de62a0bf400e4b69****"
  target_id   = "151266687691****"
  target_type = "USER"
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, ForceNew) The id of the policy.
* `target_id` - (Optional) The ID of the tag policy.
* `target_type` - (Optional) The type of the object. Valid values: `USER`, `ROOT`, `FOLDER`, `ACCOUNT`.

## Attributes Reference

* `id` - This ID of this resource. It is formatted to `<policy_Id>`:`<target_Id>`:`<target_type>`.

## Import

Tag Policy Attachment can be imported using the id, e.g.
