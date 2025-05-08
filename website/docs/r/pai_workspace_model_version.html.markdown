---
subcategory: "PAI Workspace"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_workspace_model_version"
description: |-
  Provides a Alicloud PAI Workspace Model Version resource.
---

# alicloud_pai_workspace_model_version

Provides a PAI Workspace Model Version resource.



For information about PAI Workspace Model Version and how to use it, see [What is Model Version](https://next.api.alibabacloud.com/document/AIWorkSpace/2021-02-04/CreateModelVersion).

-> **NOTE:** Available since v1.248.0.

## Example Usage

Basic Usage

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

resource "alicloud_pai_workspace_workspace" "defaultDI9fsL" {
  description    = "802"
  display_name   = var.name
  workspace_name = "${var.name}_${random_integer.default.result}"
  env_types      = ["prod"]
}

resource "alicloud_pai_workspace_model" "defaultsHptEL" {
  model_name        = var.name
  workspace_id      = alicloud_pai_workspace_workspace.defaultDI9fsL.id
  origin            = "Civitai"
  task              = "text-to-image-synthesis"
  accessibility     = "PRIVATE"
  model_type        = "Checkpoint"
  order_number      = "1"
  model_description = "ModelDescription."
  model_doc         = "https://eas-***.oss-cn-hangzhou.aliyuncs.com/s**.safetensors"
  domain            = "aigc"
  labels {
    key   = "base_model"
    value = "SD 1.5"
  }
  extra_info = {
    test = "15"
  }
}

resource "alicloud_pai_workspace_model_version" "default" {
  version_description = "VersionDescription."
  source_type         = "TrainingService"
  source_id           = "region=$${region},workspaceId=$${workspace_id},kind=TrainingJob,id=job-id"
  extra_info = {
    test = "15"
  }
  training_spec = {
    test = "TrainingSpec"
  }
  uri = "oss://hz-example-0701.oss-cn-hangzhou-internal.aliyuncs.com/checkpoints/"
  inference_spec = {
    test = "InferenceSpec"
  }
  model_id        = alicloud_pai_workspace_model.defaultsHptEL.id
  format_type     = "SavedModel"
  approval_status = "Pending"
  framework_type  = "PyTorch"
  version_name    = "1.0.0"
  metrics = {
  }
  labels {
    key   = "k1"
    value = "vs1"
  }
}
```

## Argument Reference

The following arguments are supported:
* `approval_status` - (Optional) The approval status. Valid values:
  - Pending: To be determined.
  - Approved: Allow to go online.
  - Rejected: Online is not allowed.
* `extra_info` - (Optional, Map) Other information.
* `format_type` - (Optional, ForceNew) The format of the model. Valid values:
  - OfflineModel
  - SavedModel
  - Keras H5
  - Frozen Pb
  - Caffe Prototxt
  - TorchScript
  - XGBoost
  - PMML
  - AlinkModel
  - ONNX
* `framework_type` - (Optional, ForceNew) The framework of the model. Valid values:
  - Pytorch
  - XGBoost
  - Keras
  - Caffe
  - Alink
  - Xflow
  - TensorFlow
* `inference_spec` - (Optional, Map) Describes how to apply to downstream inference services.
* `labels` - (Optional, List) List of model version labels. See [`labels`](#labels) below.
* `metrics` - (Optional, Map) The metrics for the model. The serialized length is limited to 8192.
* `model_id` - (Required, ForceNew) The model ID.
* `options` - (Optional) The extended field. This is a JSON string.
* `source_id` - (Optional) The source ID.
* `source_type` - (Optional) The type of the model source. Valid values:
  - Custom: Custom.
  - PAIFlow:PAI workflow.
  - TrainingService:PAI training service.
* `training_spec` - (Optional, Map) The training configurations. Used for fine-tuning and incremental training.
* `uri` - (Required, ForceNew) The URI of the model version.
* `version_description` - (Optional) The version descriptions.
* `version_name` - (Optional, ForceNew, Computed) The Model version.

### `labels`

The labels supports the following:
* `key` - (Optional) label key.
* `value` - (Optional) label value.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<model_id>:<version_name>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Model Version.
* `delete` - (Defaults to 5 mins) Used when delete the Model Version.
* `update` - (Defaults to 5 mins) Used when update the Model Version.

## Import

PAI Workspace Model Version can be imported using the id, e.g.

```shell
$ terraform import alicloud_pai_workspace_model_version.example <model_id>:<version_name>
```