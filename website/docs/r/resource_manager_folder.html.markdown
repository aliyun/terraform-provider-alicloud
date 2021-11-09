---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_folder"
sidebar_current: "docs-alicloud-resource-resource-manager-folder"
description: |-
  Provides a Alicloud Resource Manager Folder resource.
---

# alicloud\_resource\_manager\_folder

Provides a Resource Manager Folder resource. A folder is an organizational unit in a resource directory. You can use folders to build an organizational structure for resources.
For information about Resource Manager Foler and how to use it, see [What is Resource Manager Folder](https://www.alibabacloud.com/help/en/doc-detail/111221.htm).

-> **NOTE:** Available in v1.82.0+.

-> **NOTE:** A maximum of five levels of folders can be created under the root folder.

## Example Usage

Basic Usage

```terraform
resource "alicloud_resource_manager_folder" "example" {
  folder_name = "test"
}
```
## Argument Reference

The following arguments are supported:

* `folder_name` - (Required) The name of the folder. The name must be 1 to 24 characters in length and can contain letters, digits, underscores (_), periods (.), and hyphens (-).
* `parent_folder_id` (Optional, ForceNew) The ID of the parent folder. If not set, the system default value will be used.
                                         
## Attributes Reference

* `id` - The ID of the folder.

## Import

Resource Manager Folder can be imported using the id, e.g.

```
$ terraform import alicloud_resource_manager_folder.example fd-u8B321****	
```
