---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_resource_group"
sidebar_current: "docs-alicloud-resource-resource-manager-resource-group"
description: |-
  Provides a Alicloud Resource Manager Resource Group resource.
---

# alicloud_resource_manager_resource_group

Provides a Resource Manager Resource Group resource. If you need to group cloud resources according to business departments, projects, and other dimensions, you can create resource groups.

For information about Resource Manager Resource Group and how to use it, see [What is Resource Group](https://www.alibabacloud.com/help/en/resource-management/developer-reference/api-createresourcegroup).

-> **NOTE:** Available since v1.82.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_resource_manager_resource_group&exampleId=52f254ea-d09b-f7dc-e42f-98de82011ce2bec797e3&activeTab=example&spm=docs.r.resource_manager_resource_group.0.52f254ead0&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tfexample"
}

resource "alicloud_resource_manager_resource_group" "example" {
  resource_group_name = var.name
  display_name        = var.name
}
```

## Argument Reference

The following arguments are supported:

* `display_name` - (Required) The display name of the resource group. The name must be 1 to 50 characters in length.
* `resource_group_name` - (Optional, ForceNew, Available since v1.114.0) The unique identifier of the resource group. The identifier must be 3 to 50 characters in length and can contain letters, digits, and hyphens (-). The identifier must start with a letter.
* `tags` - (Optional, Available since v1.220.0) A mapping of tags to assign to the resource.
* `name` - (Optional, ForceNew, Deprecated since v1.114.0) Field `name` has been deprecated from provider version 1.114.0. New field `resource_group_name` instead.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Resource Group.
* `account_id` - The ID of the Alibaba Cloud account to which the resource group belongs.
* `status` - The status of the resource group.
* `region_statuses` - The status of the resource group in all regions.
  * `region_id` - The status of the region.
  * `status` - The status of the resource group.
* `create_date` - (Removed since v1.114.0) The time when the resource group was created. **NOTE:** Field `create_date` has been removed from provider version 1.114.0.

## Timeouts

-> **NOTE:** Available since v1.220.0.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 11 mins) Used when create the Resource Group.

## Import

Resource Manager Resource Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_resource_manager_resource_group.example <id>
```
