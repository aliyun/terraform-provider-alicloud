---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_data_works_project_role"
description: |-
  Provides a Alicloud Data Works Project Role resource.
---

# alicloud_data_works_project_role

Provides a Data Works Project Role resource.

For information about Data Works Project Role and how to use it, see [What is Project Role](https://www.alibabacloud.com/help/en/dataworks/developer-reference/api-dataworks-public-2024-05-18-createprojectrole).

-> **NOTE:** Available since v2.3.0.

## Example Usage

Basic Usage

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

variable "name" {
  default = "tf_example"
}

resource "random_integer" "randint" {
  max = 999
  min = 1
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_data_works_project" "default" {
  status            = "Available"
  description       = "tf_desc"
  project_name      = "${var.name}${random_integer.randint.id}"
  pai_task_enabled  = "false"
  display_name      = "tf_project_role_display"
  dev_role_disabled = "true"
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
}

resource "alicloud_data_works_project_role" "default" {
  project_id = alicloud_data_works_project.default.id
  role_name  = "${var.name}${random_integer.randint.id}"

  module_permissions {
    module_id       = "9"
    permission_type = "Read"
  }

  module_permissions {
    module_id       = "11"
    permission_type = "Write"
  }
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, ForceNew, String) The ID of the Data Works workspace that the role belongs to. Changing this creates a new role.
* `role_name` - (Required, ForceNew, String) The name of the project role. Up to 30 characters. Changing this creates a new role.
* `module_permissions` - (Required, Set) The module permissions granted to the role. UpdateProjectRole replaces the whole set, so every update must send the complete list. The collection is order-independent: the resource compares permissions by module, not by position. See [`module_permissions`](#module_permissions) below.

### module_permissions

The `module_permissions` block supports:

* `module_id` - (Required, String) The ID of the Data Works module that the permission applies to. Valid values: `2` (HoloStudio), `3` (StreamStudio), `4` (发布管理), `6` (数据保护伞), `7` (数据地图), `8` (数据服务), `9` (数据集成), `10` (数据建模（DataBlau DDM）), `11` (数据开发), `12` (数据质量), `13` (数据综合治理), `14` (运维中心), `15` (资源优化), `16` (迁移助手), `17` (数据分析), `18` (审批中心), `19` (安全中心), `20` (智能数据建模).
* `permission_type` - (Required, String) The permission level granted on the module. Valid values: `Write`, `Read`, `NotSet`. Note the API naming is counterintuitive: `Write` grants read-only access, `Read` grants edit access, and `NotSet` leaves the module without permission control.

## Attributes Reference

The following attributes are exported:

* `id` - (String) The ID of the resource. It is formatted to `<project_id>:<code>`.
* `code` - (String) The server-generated code that uniquely identifies the role within the workspace.
* `type` - (String) The role type. `UserCustom` indicates a custom role created by the workspace owner; `System` indicates a system role managed by Data Works. Only `UserCustom` roles can be created or deleted through this resource; `System` roles cannot.

## Import

Data Works Project Role can be imported using the id, e.g.

```shell
$ terraform import alicloud_data_works_project_role.example <project_id>:<code>
```
