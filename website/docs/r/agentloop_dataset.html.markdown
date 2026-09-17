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

-> **NOTE:** Available since v1.294.0.

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
* `agent_space` - (Required, ForceNew) The name of the Agent Space to which the Dataset belongs.
* `dataset_name` - (Required, ForceNew) The name of the Dataset. It must start with a lowercase letter and can contain only lowercase letters, digits, and non-consecutive underscores (`_`).
* `description` - (Optional) The description of the Dataset.
* `schema` - (Required, ForceNew, Map) The schema definition of the Dataset fields. Each map key is a field name, and the value is a JSON string describing the index key, e.g. `jsonencode({ type = "text", chn = true })`. The index key supports the following properties: `type` (the field type, such as `text`, `long`, `double`, `json`), `chn` (whether to enable Chinese word segmentation), `embedding` (the vector embedding configuration), and `jsonKeys` (the sub-field definitions when `type` is `json`). -> **NOTE:** The API automatically injects a system-defined `agentloop_annotations` field into the schema; the provider ignores it when reading, so it does not appear in the state. Formatting differences between the configured JSON text and the API representation (key order, whitespace) are suppressed.

## Attributes Reference

The following attributes are exported:
* `id` - The resource ID in terraform of Dataset. It formats as `<agent_space>:<dataset_name>`.
* `create_time` - The creation time of the Dataset.
* `region_id` - The region ID of the Dataset.
* `update_time` - The last update time of the Dataset.

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
