---
subcategory: "Eflo"
layout: "alicloud"
page_title: "Alicloud: alicloud_eflo_experiment_plan_template"
description: |-
  Provides a Alicloud Eflo Experiment Plan Template resource.
---

# alicloud_eflo_experiment_plan_template

Provides a Eflo Experiment Plan Template resource.



For information about Eflo Experiment Plan Template and how to use it, see [What is Experiment Plan Template](https://www.alibabacloud.com/help/en/pai/developer-reference/api-eflo-cnp-2023-08-28-createexperimentplantemplate).

-> **NOTE:** Available since v1.248.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_eflo_experiment_plan_template&exampleId=9e3bfa4d-7a5a-06f1-c820-ada9938540c5d6ab1ec8&activeTab=example&spm=docs.r.eflo_experiment_plan_template.0.9e3bfa4d7a&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-wulanchabu"
}

resource "alicloud_eflo_experiment_plan_template" "default" {
  template_pipeline {
    workload_id   = "2"
    workload_name = "MatMul"
    env_params {
      cpu_per_worker     = "90"
      gpu_per_worker     = "8"
      memory_per_worker  = "500"
      share_memory       = "500"
      worker_num         = "1"
      py_torch_version   = "1"
      gpu_driver_version = "1"
      cuda_version       = "1"
      nccl_version       = "1"
    }
    pipeline_order = "1"
    scene          = "baseline"
  }
  privacy_level        = "private"
  template_name        = var.name
  template_description = var.name
}
```

## Argument Reference

The following arguments are supported:
* `privacy_level` - (Required, ForceNew) Used to indicate the privacy level of the content or information. It can have the following optional parameters:
  - private: Indicates that the content is private and restricted to specific users or permission groups. Private content is usually not publicly displayed, and only authorized users can view or edit it.
  - public: Indicates that the content is public and can be accessed by anyone. Public content is usually viewable by all users and is suitable for sharing information or resources
* `template_description` - (Optional, ForceNew) Describe the purpose of this template.
* `template_name` - (Required, ForceNew) Help users identify and select specific templates.
* `template_pipeline` - (Required, Set) Representative Template Pipeline. See [`template_pipeline`](#template_pipeline) below.

### `template_pipeline`

The template_pipeline supports the following:
* `env_params` - (Required, Set) Contains a series of parameters related to the environment. See [`env_params`](#template_pipeline-env_params) below.
* `pipeline_order` - (Required, Int) Indicates the sequence number of the pipeline node.
* `scene` - (Required) The use of the template scenario. It can have the following optional parameters:
  - baseline: benchmark evaluation
* `setting_params` - (Optional, Map) Represents additional parameters for the run.
* `workload_id` - (Required, Int) Used to uniquely identify a specific payload.
* `workload_name` - (Required) The name used to represent a specific payload.

### `template_pipeline-env_params`

The template_pipeline-env_params supports the following:
* `cpu_per_worker` - (Required, Int) Number of central processing units (CPUs) allocated. This parameter affects the processing power of the computation, especially in tasks that require a large amount of parallel processing.
* `cuda_version` - (Optional) The version of CUDA(Compute Unified Device Architecture) used. CUDA is a parallel computing platform and programming model provided by NVIDIA. A specific version may affect the available GPU functions and performance optimization.
* `gpu_driver_version` - (Optional) The version of the GPU driver used. Driver version may affect GPU performance and compatibility, so it is important to ensure that the correct version is used
* `gpu_per_worker` - (Required, Int) Number of graphics processing units (GPUs). GPUs are a key component in deep learning and large-scale data processing, so this parameter is very important for tasks that require graphics-accelerated computing.
* `memory_per_worker` - (Required, Int) The amount of memory available. Memory size has an important impact on the performance and stability of the program, especially when dealing with large data sets or high-dimensional data.
* `nccl_version` - (Optional) The NVIDIA Collective Communications Library(NCCL) version used. NCCL is a library for multi-GPU and multi-node communication. This parameter is particularly important for optimizing data transmission in distributed computing.
* `py_torch_version` - (Optional) The version of the PyTorch framework used. PyTorch is a widely used deep learning library, and differences between versions may affect the performance and functional support of model training and inference.
* `share_memory` - (Required, Int) Shared memory GB allocation
* `worker_num` - (Required, Int) The total number of nodes. This parameter directly affects the parallelism and computing speed of the task, and a higher number of working nodes usually accelerates the completion of the task.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource.
* `template_id` - The ID of the template.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Experiment Plan Template.
* `delete` - (Defaults to 5 mins) Used when delete the Experiment Plan Template.
* `update` - (Defaults to 5 mins) Used when update the Experiment Plan Template.

## Import

Eflo Experiment Plan Template can be imported using the id, e.g.

```shell
$ terraform import alicloud_eflo_experiment_plan_template.example <id>
```
