---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_resource_directory"
description: |-
  Provides a Alicloud Resource Manager Resource Directory resource.
---

# alicloud_resource_manager_resource_directory

Provides a Resource Manager Resource Directory resource.



For information about Resource Manager Resource Directory and how to use it, see [What is Resource Directory](https://www.alibabacloud.com/help/en/doc-detail/94475.htm).

-> **NOTE:** Available since v1.84.0.

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
* `member_account_display_name_sync_status` - (Optional, Computed, Available since v1.259.0) The status of the Member Display Name Synchronization feature. Valid values:
  - Enabled
  - Disabled
* `member_deletion_status` - (Optional, Computed) The status of the member deletion feature. Valid values:
  - Enabled: The feature is enabled. You can call the DeleteAccount operation to delete members of the resource account type.
  - Disabled: The feature is disabled. You cannot delete members of the resource account type.
* `status` - (Optional, Computed, Available since v1.120.0) ScpStatus

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The time when the resource directory was created
* `master_account_id` - The ID of the master account
* `master_account_name` - The name of the master account
* `root_folder_id` - The ID of the root folder

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Resource Directory.
* `delete` - (Defaults to 5 mins) Used when delete the Resource Directory.
* `update` - (Defaults to 5 mins) Used when update the Resource Directory.

## Import

Resource Manager Resource Directory can be imported using the id, e.g.

```shell
$ terraform import alicloud_resource_manager_resource_directory.example <id>
```