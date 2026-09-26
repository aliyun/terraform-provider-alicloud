---
subcategory: "PAI"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_llm_trace_eval"
description: |-
  Provides a Alicloud PAI LLM Trace Eval resource
---

# alicloud_pai_llm_trace_eval

Provides a PAI LLM Trace Eval resource.

For information about PAI LLM Trace Eval and how to use it, see [What is PAI LLM Trace Eval](https://next.api.alibabacloud.com/document/paillmtrace/2024-03-11/CreateEval).

-> **NOTE:** Available since v1.249.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_pai_llm_trace_eval" "default" {
  eval_name       = var.name
  description     = "terraform-example-description"
  app_name        = "terraform-example-app"
  data_source     = "terraform-example-data-source"
  evaluation_data = "{\"key\":\"value\"}"
  metadata        = "{\"env\":\"test\"}"
}
```

## Argument Reference

The following arguments are supported:

* `app_name` - (Optional) The name of the application associated with the evaluation.
* `data_source` - (Optional) The data source for the evaluation.
* `description` - (Required) The description of the evaluation.
* `eval_name` - (Required) The name of the evaluation.
* `evaluation_data` - (Optional) The evaluation data in JSON format.
* `metadata` - (Optional) Additional metadata in JSON format.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the evaluation. It is the same as `eval_id`.
* `eval_id` - The ID of the evaluation.
* `gmt_create_time` - The creation time of the evaluation.
* `record_count` - The record count of the evaluation.
* `region_id` - The region ID of the evaluation.

## Import

PAI LLM Trace Eval can be imported using the id, e.g.

```shell
terraform import alicloud_pai_llm_trace_eval.example <eval_id>
```
