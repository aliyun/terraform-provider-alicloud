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

For information about ALB Acl and how to use it, see [What is Acl](https://www.alibabacloud.com/help/en/slb/application-load-balancer/developer-reference/api-alb-2020-06-16-createacl).

-> **NOTE:** Available since v1.133.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_alb_acl&exampleId=167b3f5a-4d43-5880-f879-aafc5c887471944649c1&activeTab=example&spm=docs.r.alb_acl.0.167b3f5a4d&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_alb_acl" "default" {
  acl_name          = "tf_example"
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_alb_acl&spm=docs.r.alb_acl.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `acl_entries` - (Optional, Deprecated from v1.166.0+) The list of the ACL entries. You can add up to `20` entries in each call.  See [`acl_entries`](#acl_entries) below for details.
**NOTE:** "Field 'acl_entries' has been deprecated from provider version 1.166.0 and it will be removed in the future version. Please use the new resource 'alicloud_alb_acl_entry_attachment'.",
* `acl_name` - (Optional) The name of the ACL. The name must be `2` to `128` characters in length, and can contain letters, digits, hyphens (-) and underscores (_). It must start with a letter.
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

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 16 mins) Used when create the Acl.
* `delete` - (Defaults to 16 mins) Used when delete the Acl.
* `update` - (Defaults to 16 mins) Used when update the Acl.

## Import

ALB Acl can be imported using the id, e.g.

```shell
$ terraform import alicloud_alb_acl.example <id>
```
