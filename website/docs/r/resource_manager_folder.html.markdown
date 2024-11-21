---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_folder"
sidebar_current: "docs-alicloud-resource-resource-manager-folder"
description: |-
  Provides a Alicloud Resource Manager Folder resource.
---

# alicloud_resource_manager_folder

Provides a Resource Manager Folder resource. A folder is an organizational unit in a resource directory. You can use folders to build an organizational structure for resources.
For information about Resource Manager Foler and how to use it, see [What is Resource Manager Folder](https://www.alibabacloud.com/help/en/doc-detail/111221.htm).

-> **NOTE:** Available since v1.82.0.

-> **NOTE:** A maximum of five levels of folders can be created under the root folder.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_resource_manager_folder&exampleId=87cafec5-c4eb-0dd1-ca5d-a76fd768ef0e90677803&activeTab=example&spm=docs.r.resource_manager_folder.0.87cafec5c4&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_resource_manager_folder" "example" {
  folder_name = "${var.name}-${random_integer.default.result}"
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

```shell
$ terraform import alicloud_resource_manager_folder.example fd-u8B321****	
```
