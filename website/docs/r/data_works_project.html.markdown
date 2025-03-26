---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_data_works_project"
description: |-
  Provides a Alicloud Data Works Project resource.
---

# alicloud_data_works_project

Provides a Data Works Project resource.



For information about Data Works Project and how to use it, see [What is Project](https://www.alibabacloud.com/help/en/dataworks/developer-reference/api-dataworks-public-2024-05-18-createproject).

-> **NOTE:** Available since v1.229.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_data_works_project&exampleId=090310a5-16b0-0567-1731-6cc93058c3b703a3c8bb&activeTab=example&spm=docs.r.data_works_project.0.090310a516&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf_example"
}

resource "random_integer" "randint" {
  max = 999
  min = 1
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_data_works_project" "default" {
  status                  = "Available"
  description             = "tf_desc"
  project_name            = "${var.name}${random_integer.randint.id}"
  pai_task_enabled        = "false"
  display_name            = "tf_new_api_display"
  dev_role_disabled       = "true"
  dev_environment_enabled = "false"
  resource_group_id       = data.alicloud_resource_manager_resource_groups.default.ids.0
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Optional) Workspace Description
* `dev_environment_enabled` - (Optional, Computed, Available since v1.237.0) Is Development Environment Enabled
* `dev_role_disabled` - (Optional, Computed, Available since v1.237.0) Is Development Role Disabled
* `display_name` - (Required) Workspace Display Name
* `pai_task_enabled` - (Required, Available since v1.237.0) Create PAI Workspace Together
* `project_name` - (Required, ForceNew) Workspace Name
* `resource_group_id` - (Optional, Computed, Available since v1.237.0) Aliyun Resource Group Id
* `status` - (Optional, Computed) Workspace Status
* `tags` - (Optional, Map, Available since v1.237.0) Aliyun Resource Tag

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 15 mins) Used when create the Project.
* `delete` - (Defaults to 5 mins) Used when delete the Project.
* `update` - (Defaults to 5 mins) Used when update the Project.

## Import

Data Works Project can be imported using the id, e.g.

```shell
$ terraform import alicloud_data_works_project.example <id>
```