---
subcategory: "PAI Feature Store"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_feature_store_llm_config"
description: |-
  Provides a Alicloud PaiFeatureStore LLM Config resource.
---

# alicloud_pai_feature_store_llm_config

Provides a PaiFeatureStore LLM Config resource. An LLM Config describes how to call a large language model (model name, base URL, API key, concurrency and token limits) and is scoped under a PaiFeatureStore instance.

For information about PaiFeatureStore LLM Config and how to use it, see [What is LLMConfig](https://next.api.alibabacloud.com/document/PaiFeatureStore/2023-06-21/CreateLLMConfig).

-> **NOTE:** Available since v1.247.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform_example"
}

resource "alicloud_pai_feature_store_llm_config" "default" {
  instance_id  = "fs-xxxxxxxx"
  workspace_id = "123"
  name         = var.name
  api_key      = "your-api-key"
  model        = "text-embedding-v1"
  base_url     = "https://dashscope.aliyuncs.com/compatible-mode/v1"
  model_type   = "TEXT_EMBEDDING"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the PaiFeatureStore instance. It can be obtained through the ListInstances API.
* `workspace_id` - (Required, ForceNew) The ID of the workspace.
* `name` - (Required) The name of the LLM config.
* `api_key` - (Required, Sensitive) The API key used to call the LLM.
* `model` - (Required) The model name, such as `text-embedding-v1`.
* `base_url` - (Required) The base URL used to call the LLM.
* `rps` - (Optional) The number of concurrent requests per second.
* `max_tokens` - (Optional) The maximum number of input tokens per row.
* `batch_size` - (Optional) The batch size.
* `embedding_dimension` - (Optional) The vector dimension. Omit to let the model use its own default dimension. Must be greater than `0`; the API treats `0` as the model default and echoes the resolved value, which would otherwise cause a perpetual diff.
* `enable_fusion` - (Optional) Whether to enable fusion.
* `model_type` - (Optional) The model type. Valid values: `TEXT_EMBEDDING`, `MULTI_MODAL_EMBEDDING`. Defaults to `TEXT_EMBEDDING` when omitted.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource, formatted as `<instance_id>:<llm_config_id>`.
* `llm_config_id` - The ID of the LLM config.
* `gmt_create_time` - The creation time.
* `gmt_modified_time` - The modification time.

## Import

PaiFeatureStore LLM Config can be imported using the composite id, e.g. `<instance_id>:<llm_config_id>`, e.g.

```shell
$ terraform import alicloud_pai_feature_store_llm_config.default fs-xxxxxxxx:llm_config1
```
