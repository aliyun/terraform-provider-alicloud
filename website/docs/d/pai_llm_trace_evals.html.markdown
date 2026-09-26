---
subcategory: "PAI"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_llm_trace_evals"
description: |-
  Provides a list of PAI LLM Trace Evals to the user
---

# alicloud_pai_llm_trace_evals

This data source provides the PAI LLM Trace Evals of the current Alibaba Cloud user.

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
  eval_name   = var.name
  description = "terraform-example-description"
  app_name    = "terraform-example-app"
  data_source = "terraform-example-data-source"
}

data "alicloud_pai_llm_trace_evals" "ids" {
  ids = ["${alicloud_pai_llm_trace_eval.default.id}"]
}

output "eval_id" {
  value = data.alicloud_pai_llm_trace_evals.ids.evals.0.eval_id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of eval IDs to filter the results.
* `name_regex` - (Optional) A regex string to filter results by eval name.
* `output_file` - (Optional) File path to save the results.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of eval IDs.
* `names` - A list of eval names.
* `evals` - A list of PAI LLM Trace Evals. Each element contains the following attributes:
  * `app_name` - The application name.
  * `data_source` - The data source.
  * `description` - The description.
  * `eval_id` - The eval ID.
  * `eval_name` - The eval name.
  * `gmt_create_time` - The creation time.
  * `id` - The eval ID, same as `eval_id`.
  * `metadata` - The metadata.
  * `record_count` - The record count.
  * `region_id` - The region ID.
