---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_folders"
sidebar_current: "docs-alicloud-datasource-resource-manager-folders"
description: |-
    Provides a list of resource manager folders to the user.
---

# alicloud\_resource\_manager\_folders

This data source provides the resource manager folders of the current Alibaba Cloud user.

-> **NOTE:**  Available in 1.84.0+.

-> **NOTE:**  You can view only the information of the first-level child folders of the specified folder.

## Example Usage

```terraform
data "alicloud_resource_manager_folders" "example" {
  name_regex = "tftest"
}

output "first_folder_id" {
  value = "${data.alicloud_resource_manager_folders.example.folders.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of resource manager folders IDs.
* `name_regex` - (Optional) A regex string to filter results by folder name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `parent_folder_id` (Optional) The ID of the parent folder.
* `query_keyword` (Optional, ForceNew, Available in 1.114.0+) The query keyword.
* `enable_details` -(Optional, Available in v1.114.0+) Default to `false`. Set it to true can output more details.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of folder IDs.
* `names` - A list of folder names.
* `folders` - A list of folders. Each element contains the following attributes:
    * `id` - The ID of the folder.
    * `folder_id`- The ID of the folder.
    * `folder_name`- The name of the folder.
    * `parent_folder_id`- (Available in v1.114.0+)The ID of the parent folder.
    
