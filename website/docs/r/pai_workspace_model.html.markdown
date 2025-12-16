---
subcategory: "PAI Workspace"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_workspace_model"
description: |-
  Provides a Alicloud PAI Workspace Model resource.
---

# alicloud_pai_workspace_model

Provides a PAI Workspace Model resource.



For information about PAI Workspace Model and how to use it, see [What is Model](https://www.alibabacloud.com/help/en/pai/developer-reference/api-aiworkspace-2021-02-04-createmodel).

-> **NOTE:** Available since v1.249.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_pai_workspace_model&exampleId=04812c08-8b7a-24b8-1279-05f15fc0136d5e9463ae&activeTab=example&spm=docs.r.pai_workspace_model.0.04812c088b&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform_example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_pai_workspace_workspace" "defaultENuC6u" {
  description    = "156"
  display_name   = var.name
  workspace_name = "${var.name}_${random_integer.default.result}"
  env_types      = ["prod"]
}

resource "alicloud_pai_workspace_model" "default" {
  origin        = "Civitai"
  task          = "text-to-image-synthesis"
  model_name    = var.name
  accessibility = "PRIVATE"
  workspace_id  = alicloud_pai_workspace_workspace.defaultENuC6u.id
  model_type    = "Checkpoint"
  labels {
    key   = "base_model"
    value = "SD 1.5"
  }
  order_number = "1"
  extra_info = {
    test = "15"
  }
  model_description = "ModelDescription."
  model_doc         = "https://eas-***.oss-cn-hangzhou.aliyuncs.com/s**.safetensors"
  domain            = "aigc"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_pai_workspace_model&spm=docs.r.pai_workspace_model.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `accessibility` - (Optional) The visibility of the model in the workspace. Default value: `PRIVATE`. Valid values:
  - `PRIVATE`: In this workspace, it is only visible to you and the administrator.
  - `PUBLIC`: In this workspace, it is visible to everyone.
* `domain` - (Optional) The domain of the model. Describe the domain in which the model solves the problem. For example: nlp (natural language processing), cv (computer vision), etc.
* `extra_info` - (Optional, Map) Other information about the model.
* `labels` - (Optional, List) A list of tags. See [`labels`](#labels) below.
* `model_description` - (Optional) The model description, used to distinguish different models.
* `model_doc` - (Optional) The documentation of the model.
* `model_name` - (Required) The name of the model. The name must be 1 to 127 characters in length.
* `model_type` - (Optional) The model type. Example: Checkpoint or LoRA.
* `order_number` - (Optional, Int) The sequence number of the model. Can be used for custom sorting.
* `origin` - (Optional) The source of the model. The community or organization to which the source model belongs, such as ModelScope or HuggingFace.
* `task` - (Optional) The task of the model. Describes the specific problem that the model solves. Example: text-classification.
* `workspace_id` - (Optional, ForceNew) The ID of the workspace.

### `labels`

The labels supports the following:
* `key` - (Optional) label key
* `value` - (Optional) label value

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Model.
* `delete` - (Defaults to 5 mins) Used when delete the Model.
* `update` - (Defaults to 5 mins) Used when update the Model.

## Import

PAI Workspace Model can be imported using the id, e.g.

```shell
$ terraform import alicloud_pai_workspace_model.example <id>
```
