---
subcategory: "Application Load Balancer (ALB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_alb_acl"
sidebar_current: "docs-alicloud-resource-alb-acl"
description: |-
  Provides a Alicloud Application Load Balancer (ALB) Acl resource.
---

# alicloud_alb_acl

Provides a Application Load Balancer (ALB) Acl resource.

For information about ALB Acl and how to use it, see [What is Acl](https://www.alibabacloud.com/help/en/server-load-balancer/latest/api-doc-alb-2020-06-16-api-doc-createacl).

-> **NOTE:** Available since v1.133.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_alb_acl" "default" {
  acl_name          = "tf_example"
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
}
```

## Argument Reference

The following arguments are supported:

* `acl_entries` - (Optional, Deprecated from v1.166.0+) The list of the ACL entries. You can add up to `20` entries in each call.  See [`acl_entries`](#acl_entries) below for details.
**NOTE:** "Field 'acl_entries' has been deprecated from provider version 1.166.0 and it will be removed in the future version. Please use the new resource 'alicloud_alb_acl_entry_attachment'.",
* `acl_name` - (Required) The name of the ACL. The name must be `2` to `128` characters in length, and can contain letters, digits, hyphens (-) and underscores (_). It must start with a letter.
* `dry_run` - (Optional) Specifies whether to precheck the API request. 
* `resource_group_id` - (Optional, ForceNew) The ID of the resource group.
* `tags` - (Optional) A mapping of tags to assign to the resource.

### `acl_entries`

The acl_entries supports the following: 

* `description` - (Optional) The description of the ACL entry. The description must be `1` to `256` characters in length, and can contain letters, digits, hyphens (-), forward slashes (/), periods (.),and underscores (_). It can also contain Chinese characters.
* `entry` - (Optional) The IP address for the ACL entry.
* `status` - (Optional) The status of the ACL entry. Valid values:
  - `Adding`: The ACL entry is being added.
  - `Available`: The ACL entry is added and available.
  - `Removing`: The ACL entry is being removed.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Acl.
* `status` - The state of the ACL. Valid values:`Provisioning`, `Available` and `Configuring`. `Provisioning`: The ACL is being created. `Available`: The ACL is available. `Configuring`: The ACL is being configured.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 16 mins) Used when create the Acl.
* `delete` - (Defaults to 16 mins) Used when delete the Acl.
* `update` - (Defaults to 16 mins) Used when update the Acl.

## Import

ALB Acl can be imported using the id, e.g.

```shell
$ terraform import alicloud_alb_acl.example <id>
```
