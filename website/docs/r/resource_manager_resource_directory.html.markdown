---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_resource_directory"
sidebar_current: "docs-alicloud-resource-resource-manager-resource-directory"
description: |-
  Provides a Alicloud Resource Manager Resource Directory resource.
---

# alicloud_resource_manager_resource_directory

Provides a Resource Manager Resource Directory resource. Resource Directory enables you to establish an organizational structure for the resources used by applications of your enterprise. You can plan, build, and manage the resources in a centralized manner by using only one resource directory.

For information about Resource Manager Resource Directory and how to use it, see [What is Resource Manager Resource Directory](https://www.alibabacloud.com/help/en/doc-detail/94475.htm).

-> **NOTE:** Available since v1.84.0.

-> **NOTE:** An account can only be used to enable a resource directory after it passes enterprise real-name verification. An account that only passed individual real-name verification cannot be used to enable a resource directory.

-> **NOTE:** Before you destroy the resource, make sure that the following requirements are met:
  - All member accounts must be removed from the resource directory. 
  - All folders except the root folder must be deleted from the resource directory.
  
## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_resource_manager_resource_directory&exampleId=0a7b7736-1528-148e-ea16-f8b3cc4cbc4474fdc65c&activeTab=example&spm=docs.r.resource_manager_resource_directory.0.0a7b773615&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_resource_manager_resource_directories" "default" {
}

resource "alicloud_resource_manager_resource_directory" "default" {
  count  = length(data.alicloud_resource_manager_resource_directories.default.directories) > 0 ? 0 : 1
  status = "Enabled"
}
```
## Argument Reference

The following arguments are supported:

* `status` - (Optional, Available since v1.120.0) The status of control policy. Valid values:`Enabled` and `Disabled`.
* `member_deletion_status` - (Optional, Available since v1.201.0) Specifies whether to enable the member deletion feature. Valid values:`Enabled` and `Disabled`.

## Attributes Reference

* `id` - The ID of the resource directory.
* `root_folder_id` - The ID of the root folder.
* `master_account_id` - The ID of the master account.
* `master_account_name` - The name of the master account.

## Timeouts

-> **NOTE:** Available since v1.120.0.

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `update` - (Defaults to 6 mins) Used when update the control policy status.

## Import

Resource Manager Resource Directory can be imported using the id, e.g.

```shell
$ terraform import alicloud_resource_manager_resource_directory.example rd-s3****
```
