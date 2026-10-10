---
subcategory: "PAI Workspace"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_workspace_prompts"
description: |-
  Provides a list of PAI Workspace Prompts to the user.
---

# alicloud_pai_workspace_prompts

This data source provides a list of PAI Workspace Prompts available to the user.

-> **NOTE:** Available since v1.287.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_pai_workspace_prompts" "default" {
  workspace_id   = alicloud_pai_workspace_workspace.default.id
  framework_type = "ICIO"
}

output "first_prompt_id" {
  value = data.alicloud_pai_workspace_prompts.default.prompts.0.id
}
```

## Argument Reference

The following arguments are supported:

* `workspace_id` - (Required, ForceNew) Workspace id.
* `framework_type` - (Optional) Cue word frame type.
* `ids` - (Optional) A list of Prompt IDs. The value is formulated as `<workspace_id>:<prompt_id>`.
* `name_regex` - (Optional) A regex string to filter results by Prompt name.
* `output_file` - (Optional) File path where data source results are saved.

## Attributes Reference

The following attributes are exported in addition to the `Argument Reference` above:

* `ids` - A list of Prompt IDs.
* `names` - A list of Prompt names.
* `prompts` - A list of Prompts. Each element contains the following attributes:
  * `id` - The ID of the Prompt. The value is formulated as `<workspace_id>:<prompt_id>`.
  * `prompt_id` - Prompt word id.
  * `prompt_name` - Resource attribute field representing the resource name.
  * `accessibility` - Represents the visibility property of the resource in the workspace.
  * `create_time` - Resource attribute field representing creation time.
  * `description` - Attributes that represent the context of the thread on the resource.
  * `framework_content` - Cue word frame content.
  * `framework_type` - Cue word frame type.
  * `modify_time` - Resource attribute representing modification time.
