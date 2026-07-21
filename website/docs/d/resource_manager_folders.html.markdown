---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_folders"
sidebar_current: "docs-alicloud-datasource-resource-manager-folders"
description: |-
    Provides a list of Resource Manager Folders to the user.
---

# alicloud_resource_manager_folders

This data source provides the Resource Manager Folders of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.84.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_resource_manager_folder" "default" {
  folder_name = var.name
}

data "alicloud_resource_manager_folders" "ids" {
  ids = [alicloud_resource_manager_folder.default.id]
}

output "resource_manager_folder_id_0" {
  value = data.alicloud_resource_manager_folders.ids.folders.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, List) A list of Folders IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Folder name.
* `parent_folder_id` (Optional, ForceNew) The ID of the parent folder. **NOTE:** If `parent_folder_id` is not set, the information of the first-level subfolders of the Root folder is queried.
* `query_keyword` (Optional, ForceNew, Available since v1.114.0) The keyword used for the query, such as a folder name. Fuzzy match is supported.
* `enable_details` -(Optional, Bool, Available since v1.114.0) Whether to query the detailed list of resource attributes. Default value: `false`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Folder names.
* `folders` - A list of Folder. Each element contains the following attributes:
  * `id` - The ID of the Resource Manager Folder.
  * `folder_id`- The ID of the Folder.
  * `folder_name`- The Name of the Folder.
  * `parent_folder_id`- (Available since v1.114.0) The ID of the parent folder. **Note:** `parent_folder_id` takes effect only if `enable_details` is set to `true`.
