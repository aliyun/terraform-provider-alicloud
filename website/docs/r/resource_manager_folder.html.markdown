---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_folder"
description: |-
  Provides a Alicloud Resource Manager Folder resource.
---

# alicloud_resource_manager_folder

Provides a Resource Manager Folder resource.

The management unit of the organization account in the resource directory.

For information about Resource Manager Folder and how to use it, see [What is Folder](https://www.alibabacloud.com/help/en/resource-management/resource-directory/developer-reference/api-resourcedirectorymaster-2022-04-19-createfolder).

-> **NOTE:** Available since v1.82.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_resource_manager_folder&exampleId=87cafec5-c4eb-0dd1-ca5d-a76fd768ef0e90677803&activeTab=example&spm=docs.r.resource_manager_folder.0.87cafec5c4&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
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
* `folder_name` - (Required) The name of the folder.
* `parent_folder_id` - (Optional, ForceNew) The ID of the parent folder.
* `tags` - (Optional, Map, Available since v1.259.0) The tag of the resource.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - (Available since v1.259.0) The time when the folder was created.

## Timeouts

-> **NOTE:** Available since v1.259.0.

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Folder.
* `delete` - (Defaults to 5 mins) Used when delete the Folder.
* `update` - (Defaults to 5 mins) Used when update the Folder.

## Import

Resource Manager Folder can be imported using the id, e.g.

```shell
$ terraform import alicloud_resource_manager_folder.example <id>
```
