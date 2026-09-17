---
subcategory: "AgentLoop"
layout: "alicloud"
page_title: "Alicloud: alicloud_agentloop_dataset"
description: |-
  Provides a Alicloud AgentLoop Dataset resource.
---

# alicloud_agentloop_dataset

Provides a AgentLoop Dataset resource.

Dataset is used to define and manage the data sets in the Agent Space, which can be used as the data source or data sink of Pipelines and Evaluation Tasks.

For information about AgentLoop Dataset and how to use it, see [What is Dataset](https://next.api.alibabacloud.com/document/AgentLoop/2026-05-20/CreateDataset).

-> **NOTE:** Available since v1.295.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_agentloop_agent_space" "default" {
  agent_space = "${var.name}-as"
}

resource "alicloud_agentloop_dataset" "default" {
  agent_space  = alicloud_agentloop_agent_space.default.agent_space
  dataset_name = "terraform_example"
  description  = "terraform-example"
  schema = {
    input  = jsonencode({ type = "text", chn = true })
    output = jsonencode({ type = "text", chn = false })
  }
}
```

## Argument Reference

The following arguments are supported:
* `agent_space` - (Required, ForceNew) The name of the Agent Space to which the Dataset belongs. It must be 2 to 64 characters in length.
* `dataset_name` - (Required, ForceNew) The name of the Dataset. It must be 2 to 64 characters in length, must start with a lowercase letter, can contain only lowercase letters, digits, and non-consecutive underscores (`_`), and cannot end with an underscore.
* `description` - (Optional) The description of the Dataset. It can be modified in place.
* `schema` - (Required, ForceNew, Map) The schema definition of the Dataset fields. Each map key is a field name, and the value is a JSON string describing the index key, e.g. `jsonencode({ type = "text", chn = true })`. The index key supports the following properties:
  - `type` - The field type, such as `text`, `long`, `double`, `json`.
  - `chn` - Whether to enable Chinese word segmentation. It is required when `type` is `text`; set it explicitly even when the value is `false`.
  - `description` - An optional string description of the field.
  - `embedding` - The name of the embedding model used to build the vector index (a string identifier, e.g. `agentloop-embedding-v4`), not a configuration object.
  - `jsonKeys` - The sub-field definitions when `type` is `json`. Each sub-field only supports `type` and `chn`, and `chn` is required when the sub-field `type` is `text`.

-> **NOTE:** The API automatically injects a system-defined `agentloop_annotations` field into the schema; the provider ignores it when reading, so it does not appear in the state. The JSON values are compared semantically: differences in key order or whitespace between the configured JSON text and the API representation do not produce a diff and do not recreate the Dataset.

## Attributes Reference

The following attributes are exported:
* `id` - The resource ID in terraform of Dataset. It formats as `<agent_space>:<dataset_name>`.
* `create_time` - The creation time of the Dataset, in ISO 8601 UTC format (`yyyy-MM-ddTHH:mm:ssZ`).
* `region_id` - The region ID of the Dataset.
* `update_time` - The last update time of the Dataset, in ISO 8601 UTC format (`yyyy-MM-ddTHH:mm:ssZ`).

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Dataset.
* `delete` - (Defaults to 5 mins) Used when delete the Dataset.
* `update` - (Defaults to 5 mins) Used when update the Dataset.

## Import

AgentLoop Dataset can be imported using the id, e.g.

```shell
$ terraform import alicloud_agentloop_dataset.example <agent_space>:<dataset_name>
```
