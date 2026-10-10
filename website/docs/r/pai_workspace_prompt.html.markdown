---
subcategory: "PAI Workspace"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_workspace_prompt"
description: |-
  Provides a Alicloud PAI Workspace Prompt resource.
---

# alicloud_pai_workspace_prompt

Provides a PAI Workspace Prompt resource.

Custom Prompt Word.

For information about PAI Workspace Prompt and how to use it, see [What is Prompt](https://next.api.alibabacloud.com/document/AIWorkSpace/2021-02-04/CreatePrompt).

-> **NOTE:** Available since v1.287.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_pai_workspace_workspace" "default7Adve7" {
  description    = "example_prompt_1790392362"
  display_name   = "用来example提示词"
  workspace_name = "prompt_1790392362"
  env_types      = ["prod"]
}


resource "alicloud_pai_workspace_prompt" "default" {
  prompt_name       = "prompt_1790392365"
  description       = "这是一个提示词example模版"
  accessibility     = "PRIVATE"
  framework_content = "{     \"PromptContext\":\"你是一个拥有十年驾龄的老司机，请你针对以下图片场景做出你的分析判断。\",     \"Tags\":{     \"侧翻的车辆\":\"车辆侧翻在地,4个车轮至少有两个离开地面\",     \"匝道\":\"只有明确看到高速路上的大弯道，一般匝道都在高速路干道的右侧，进出收费站才可判定存在。\"     } }"
  framework_type    = "ICIO"
  workspace_id      = alicloud_pai_workspace_workspace.default7Adve7.id
}
```

## Argument Reference

The following arguments are supported:
* `accessibility` - (Optional, ForceNew, Computed) Represents the visibility property of the resource in the workspace
* `description` - (Optional) Attributes that represent the context of the thread on the resource
* `framework_content` - (Optional) Cue word frame content
* `framework_type` - (Optional) Cue word frame type
* `prompt_name` - (Required, ForceNew) Resource attribute field representing the resource name
* `workspace_id` - (Required, ForceNew) Workspace id

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<workspace_id>:<prompt_id>`.
* `create_time` - Resource attribute field representing creation time.
* `modify_time` - Resource attribute representing modification time.
* `prompt_id` - Prompt word id.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Prompt.
* `delete` - (Defaults to 5 mins) Used when delete the Prompt.
* `update` - (Defaults to 5 mins) Used when update the Prompt.

## Import

PAI Workspace Prompt can be imported using the id, e.g.

```shell
$ terraform import alicloud_pai_workspace_prompt.example <workspace_id>:<prompt_id>
```