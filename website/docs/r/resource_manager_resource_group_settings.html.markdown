---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_resource_group_settings"
description: |-
  Provides a Alicloud Resource Manager Resource Group Settings resource.
---

# alicloud_resource_manager_resource_group_settings

Provides a Resource Manager Resource Group Settings resource.

Resource group product feature settings, such as automatic group transfer, resource group administrator, and default resource group transfer notification.

For information about Resource Manager Resource Group Settings and how to use it, see [What is Resource Group Settings](https://next.api.alibabacloud.com/document/ResourceManager/2020-03-31/UpdateResourceGroupAdminSetting).

-> **NOTE:** Available since v1.287.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_resource_manager_resource_group_settings" "default" {
  resource_group_admin_setting_status        = true
  resource_group_notification_setting_status = true
}
```

### Deleting `alicloud_resource_manager_resource_group_settings` or removing it from your configuration

Terraform cannot destroy resource `alicloud_resource_manager_resource_group_settings`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `resource_group_admin_setting_status` - (Required) Specifies whether to designate the resource creator as the administrator of the resource group. When enabled, the user who creates a resource is automatically set as the administrator of the resource group that the resource belongs to. This maps to the `CreatorAsAdmin` parameter of the UpdateResourceGroupAdminSetting API.
* `resource_group_notification_setting_status` - (Optional) Specifies whether to enable notifications when resources are automatically transferred to the default resource group. When enabled, a notification is sent upon automatic transfer of resources to the default resource group. This maps to the `ResourceGroupNotificationEnableStatus` field returned by GetResourceGroupNotificationSetting, and is toggled through the EnableResourceGroupNotification and DisableResourceGroupNotification APIs.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<Alibaba Cloud Account ID>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Resource Group Settings.
* `update` - (Defaults to 5 mins) Used when update the Resource Group Settings.

## Import

Resource Manager Resource Group Settings can be imported using the id, e.g.

```shell
$ terraform import alicloud_resource_manager_resource_group_settings.example <Alibaba Cloud Account ID>
```