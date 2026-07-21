---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_data_works_folders"
sidebar_current: "docs-alicloud-datasource-data-works-folders"
description: |-
  Provides a list of Data Works Folders to the user.
---

# alicloud\_data\_works\_folders

This data source provides the Data Works Folders of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.131.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_data_works_folder" "default" {
  project_id  = "xxxx"
  folder_path = "Business Flow/tfTestAcc/folderDi"
}

data "alicloud_data_works_folders" "ids" {
  ids                = [alicloud_data_works_folder.default.folder_id]
  project_id         = alicloud_data_works_folder.default.project_id
  parent_folder_path = "Business Flow/tfTestAcc/folderDi"
}

output "data_works_folder_id_1" {
  value = data.alicloud_data_works_folders.ids.folders.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Folder IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `parent_folder_path` - (Required, ForceNew) The parent folder path.
* `project_id` - (Optional, ForceNew, Available in v1.131.0+) The ID of the project.
* `project_identifier` - (Optional, ForceNew) The project identifier.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `folders` - A list of Data Works Folders. Each element contains the following attributes:
	* `id` - The Folder ID.
	* `folder_path` - Folder Path.
	* `project_id` - The ID of the project.
