---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_resource_directory_sharing"
description: |-
  Provides a Alicloud Resource Manager Resource Directory Sharing resource.
---

# alicloud_resource_manager_resource_directory_sharing

Provides a Resource Manager Resource Directory Sharing resource.

Resource directory sharing, which enables sharing with the resource directory.

For information about Resource Manager Resource Directory Sharing and how to use it, see [What is Resource Directory Sharing](https://next.api.alibabacloud.com/document/ResourceSharing/2020-01-10/EnableSharingWithResourceDirectory).

-> **NOTE:** Available since v1.283.0.

-> **NOTE:** Sharing with the resource directory is an account-level capability. Once enabled, the underlying service API does not provide a way to disable it, so this resource cannot be torn down.

## Example Usage

Basic Usage

```terraform
resource "alicloud_resource_manager_resource_directory_sharing" "default" {
}
```

### Deleting `alicloud_resource_manager_resource_directory_sharing` or removing it from your configuration

Sharing with the resource directory is enabled at the Alibaba Cloud account level, and once enabled it cannot be disabled through the underlying API. Terraform cannot destroy resource `alicloud_resource_manager_resource_directory_sharing`: removing it from your configuration only removes it from the Terraform state file, and the feature will remain enabled on the account.

## Argument Reference

This resource has no configurable arguments.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<Alibaba Cloud Account ID>`.
* `enable_sharing_with_rd` - Indicates whether sharing with the resource directory is enabled.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Resource Directory Sharing.

## Import

Resource Manager Resource Directory Sharing can be imported using the id, e.g.

```shell
$ terraform import alicloud_resource_manager_resource_directory_sharing.example <Alibaba Cloud Account ID>
```