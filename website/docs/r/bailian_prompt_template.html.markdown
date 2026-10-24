---
subcategory: "Bailian (SFM)"
layout: "alicloud"
page_title: "Alicloud: alicloud_bailian_prompt_template"
sidebar_current: "docs-alicloud-resource-bailian-prompt-template"
description: |-
  Provides a Bailian Prompt Template resource.
---

# alicloud_bailian_prompt_template

Provides a Bailian Prompt Template resource.

For information about Bailian Prompt Template and how to use it, see [What is Prompt Template](https://www.alibabacloud.com/help/en/model-studio/developer-api/api-reference-26).

-> **NOTE:** Available since v1.243.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_bailian_prompt_template" "default" {
  workspace_id = "llm-xxx"
  name         = "tf-example-prompt-template"
  content      = "You are a helpful assistant."
}
```

## Argument Reference

The following arguments are supported:

* `workspace_id` - (Required, ForceNew) The workspace ID of Bailian.
* `name` - (Required) The name of the prompt template.
* `content` - (Required) The content of the prompt template.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource, format as `<workspace_id>:<prompt_template_id>`.
* `prompt_template_id` - The ID of the prompt template.

## Import

Bailian Prompt Template can be imported using the id, e.g.

```shell
$ terraform import alicloud_bailian_prompt_template.example <workspace_id>:<prompt_template_id>
```
