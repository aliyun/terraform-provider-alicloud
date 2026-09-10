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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_eflo_experiment_plan&exampleId=24bb02bc-489e-e03d-6b04-dc8d18d0c9ea061b843c&activeTab=example&spm=docs.r.eflo_experiment_plan.0.24bb02bc48&intl_lang=EN_US" target="_blank">
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

resource "random_integer" "default" {
  min = 10000
  max = 99999
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

resource "alicloud_eflo_resource" "default" {
  user_access_param {
    access_id    = "your_access_id"
    access_key   = "your_access_key"
    workspace_id = "your_workspace_id"
    endpoint     = "your_endpoint"
  }
  cluster_id = "terraform-${random_integer.default.result}"
  machine_types {
    memory_info  = "32x 64GB DDR4 4800 Memory"
    type         = "Private"
    bond_num     = "5"
    node_count   = "1"
    cpu_info     = "2x Intel Saphhire Rapid 8469C 48C CPU"
    network_info = "1x 200Gbps Dual Port BF3 DPU for VPC 4x 200Gbps Dual Port EIC"
    gpu_info     = "8x OAM 810 GPU"
    disk_info    = "2x 480GB SATA SSD 4x 3.84TB NVMe SSD"
    network_mode = "net"
    name         = "lingjun"
  }
  cluster_name = var.name
  cluster_desc = var.name
}

resource "alicloud_eflo_experiment_plan" "default" {
  resource_id = alicloud_eflo_resource.default.resource_id
  plan_name   = var.name
  template_id = alicloud_eflo_experiment_plan_template.defaultpSZN7t.id
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_eflo_experiment_plan&spm=docs.r.eflo_experiment_plan.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `external_params` - (Optional, Map) Additional operating parameters. You can include information about the specified node.
* `plan_name` - (Optional) Indicates the name of the experiment plan, which is used to distinguish different experiment plans.
* `resource_group_id` - (Optional) The ID of the resource group.
* `resource_id` - (Required, ForceNew, Int) The ID of the resource.
* `tags` - (Optional, Map) The tag of the resource.
* `template_id` - (Required, ForceNew, Int) The ID of the template.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Experiment Plan.
* `delete` - (Defaults to 5 mins) Used when delete the Experiment Plan.
* `update` - (Defaults to 5 mins) Used when update the Experiment Plan.

## Import

Eflo Experiment Plan can be imported using the id, e.g.

```shell
$ terraform import alicloud_eflo_experiment_plan.example <id>
```
