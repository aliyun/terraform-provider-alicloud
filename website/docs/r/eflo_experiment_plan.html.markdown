---
subcategory: "Eflo"
layout: "alicloud"
page_title: "Alicloud: alicloud_eflo_experiment_plan"
description: |-
  Provides a Alicloud Eflo Experiment Plan resource.
---

# alicloud_eflo_experiment_plan

Provides a Eflo Experiment Plan resource.



For information about Eflo Experiment Plan and how to use it, see [What is Experiment Plan](https://www.alibabacloud.com/help/en/pai/developer-reference/api-eflo-cnp-2023-08-28-createexperimentplan).

-> **NOTE:** Available since v1.248.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-wulanchabu"
}

resource "alicloud_eflo_experiment_plan_template" "defaultpSZN7t" {
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

resource "alicloud_eflo_experiment_plan" "default" {
  resource_id = 36
  plan_name   = var.name
  template_id = alicloud_eflo_experiment_plan_template.defaultpSZN7t.id
}
```

## Argument Reference

The following arguments are supported:
* `external_params` - (Optional, Map) Additional operating parameters. You can include information about the specified node.
* `plan_name` - (Optional) Indicates the name of the experiment plan, which is used to distinguish different experiment plans.
* `resource_group_id` - (Optional) The ID of the resource group.
* `resource_id` - (Required, ForceNew, Int) The ID of the resource.
* `tags` - (Optional, Map) The tag of the resource
* `template_id` - (Required, ForceNew, Int) The ID of the template.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Experiment Plan.
* `delete` - (Defaults to 5 mins) Used when delete the Experiment Plan.
* `update` - (Defaults to 5 mins) Used when update the Experiment Plan.

## Import

Eflo Experiment Plan can be imported using the id, e.g.

```shell
$ terraform import alicloud_eflo_experiment_plan.example <id>
```
