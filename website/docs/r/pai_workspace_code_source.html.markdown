---
subcategory: "PAI Workspace"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_workspace_code_source"
description: |-
  Provides a Alicloud PAI Workspace Code Source resource.
---

# alicloud_pai_workspace_code_source

Provides a PAI Workspace Code Source resource.



For information about PAI Workspace Code Source and how to use it, see [What is Code Source](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.236.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform_example"
}

provider "alicloud" {
  region = "cn-shenzhen"
}

resource "alicloud_pai_workspace_workspace" "defaultgklBnM" {
  description    = "for-pop-example"
  display_name   = "CodeSourceTest_1732796227"
  workspace_name = var.name
  env_types      = ["prod"]
}


resource "alicloud_pai_workspace_code_source" "default" {
  mount_path             = "/mnt/code/dir_01/"
  code_repo              = "https://github.com/mattn/go-sqlite3.git"
  description            = "desc-01"
  code_repo_access_token = "token-01"
  accessibility          = "PRIVATE"
  display_name           = "codesource-example-01"
  workspace_id           = alicloud_pai_workspace_workspace.defaultgklBnM.id
  code_branch            = "master"
  code_repo_user_name    = "user-01"
}
```

## Argument Reference

The following arguments are supported:
* `accessibility` - (Required) Visibility of the code configuration, possible values:
  - PRIVATE: In this workspace, it is only visible to you and the administrator.
  - PUBLIC: In this workspace, it is visible to everyone.
* `code_branch` - (Optional) Code repository branch.
* `code_commit` - (Optional) The code CommitId.
* `code_repo` - (Required) Code repository address.
* `code_repo_access_token` - (Optional) The Token used to access the code repository.
* `code_repo_user_name` - (Optional) The user name of the code repository.
* `description` - (Optional) A detailed description of the code configuration.
* `display_name` - (Required) Code source configuration name.
* `mount_path` - (Required) The local Mount Directory of the code.
* `workspace_id` - (Required, ForceNew) The ID of the workspace.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Code Source.
* `delete` - (Defaults to 5 mins) Used when delete the Code Source.
* `update` - (Defaults to 5 mins) Used when update the Code Source.

## Import

PAI Workspace Code Source can be imported using the id, e.g.

```shell
$ terraform import alicloud_pai_workspace_code_source.example <id>
```