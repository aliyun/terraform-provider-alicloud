---
subcategory: "Application Load Balancer (ALB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_alb_acl"
sidebar_current: "docs-alicloud-resource-alb-acl"
description: |-
  Provides a Alicloud Application Load Balancer (ALB) Acl resource.
---

# alicloud\_alb\_acl

Provides a Application Load Balancer (ALB) Acl resource.

For information about ALB Acl and how to use it, see [What is Acl](https://www.alibabacloud.com/help/doc-detail/200280.html).

-> **NOTE:** Available in v1.133.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_alb_acl" "example" {
  acl_name = "example_value"
  acl_entries {
    description = "example_value"
    entry       = "10.0.0.0/24"
  }
}

```

## Argument Reference

The following arguments are supported:

* `acl_entries` - (Optional) The list of the ACL entries. You can add up to `20` entries in each call.  **NOTE:** "Field 'acl_entries' has been deprecated from provider version 1.166.0 and it will be removed in the future version. Please use the new resource 'alicloud_alb_acl_entry_attachment'.",
* `acl_name` - (Required) The name of the ACL. The name must be `2` to `128` characters in length, and can contain letters, digits, hyphens (-) and underscores (_). It must start with a letter.
* `dry_run` - (Optional) Specifies whether to precheck the API request. 
* `resource_group_id` - (Optional, Computed, ForceNew) The ID of the resource group.
* `tags` - (Optional) A mapping of tags to assign to the resource.

#### Block acl_entries

The acl_entries supports the following: 

* `description` - (Optional) The description of the ACL entry. The description must be `1` to `256` characters in length, and can contain letters, digits, hyphens (-), forward slashes (/), periods (.),and underscores (_). It can also contain Chinese characters.
* `entry` - (Optional) The IP address for the ACL entry.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Acl.
* `status` - The state of the ACL. Valid values:`Provisioning`, `Available` and `Configuring`. `Provisioning`: The ACL is being created. `Available`: The ACL is available. `Configuring`: The ACL is being configured.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 16 mins) Used when create the Acl.
* `delete` - (Defaults to 16 mins) Used when delete the Acl.
* `update` - (Defaults to 16 mins) Used when update the Acl.

## Import

ALB Acl can be imported using the id, e.g.

```
$ terraform import alicloud_alb_acl.example <id>
```
