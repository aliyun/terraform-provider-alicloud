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
For information about Tag Policy Attachment and how to use it.

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, ForceNew) The id of the policy.
* `target_Id` - (Optional) The ID of the tag policy.
* `target_type` - (Optional) The type of the object. Valid values: `USER`, `ROOT`, `FOLDER`, `ACCOUNT`.

## Attributes Reference

* `id` - This ID of this resource. It is formatted to `<policy_Id>`:`<target_Id>`:`<target_type>`.

## Import

Tag Policy Attachment can be imported using the id, e.g.

