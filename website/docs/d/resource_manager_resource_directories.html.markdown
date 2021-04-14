---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_resource_directories"
sidebar_current: "docs-alicloud-datasource-resource-manager-resource-directories"
description: |-
    Provides a list of Resource Manager Resource Directories to the user.
---

# alicloud\_resource\_manager\_resource\_directories

This data source provides the Resource Manager Resource Directories of the current Alibaba Cloud user.

-> **NOTE:**  Available in 1.86.0+.

## Example Usage

```terraform
data "alicloud_resource_manager_resource_directories" "default" {}

output "resource_directory_id" {
  value = "${data.alicloud_resource_manager_resource_directories.default.directories.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `directories` - A list of resource directories. Each element contains the following attributes:
    * `id` - The ID of resource directory.
    * `master_account_id`- The ID of the master account.
    * `master_account_name`- The name of the master account.
    * `resource_directory_id` - The ID of the resource directory.
    * `root_folder_id` - The ID of the root folder.
    * `status` - (Available in 1.120.0+.) The status of the control policy.
    
