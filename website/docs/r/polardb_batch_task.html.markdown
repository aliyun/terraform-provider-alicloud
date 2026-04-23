---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_batch_task"
sidebar_current: "docs-alicloud-resource-polardb-batch-task"
description: |-
  Provides a PolarDB Batch Task resource.
---

# alicloud_polardb_batch_task

Provides a PolarDB Batch Task resource. This resource is used to manage batch operations, such as installing or uninstalling skills (e.g., Polar Claw), on multiple PolarDB instances simultaneously.

-> **NOTE:** Available since v1.240.0.

## Example Usage

Create a PolarDB Batch Task to install skills

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_polardb_batch_task&exampleId=example-batch-task-01&activeTab=example&spm=docs.r.polardb_batch_task.0.example&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_polardb_batch_task" "default" {
  task_name    = "terraform-batch-task-example"
  task_type    = "polarclaw_install_skills"
  region_id    = "cn-hangzhou"
  instance_ids = ["pa-xxx"]

  task_params {
    skill_name = "ontology"
    version    = "1.0.4"
  }
}
```

### Removing alicloud_polardb_batch_task from your configuration

The `alicloud_polardb_batch_task` resource allows you to manage batch tasks. Note that this resource represents an asynchronous operation. Removing this resource from your configuration will remove it from your statefile and management, but it will not undo the actions performed by the task (e.g., it will not uninstall the skills if the task was an installation). You can verify the status of the instances via the PolarDB Console.

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_polardb_batch_task&spm=docs.r.polardb_batch_task.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `task_type` - (Required, ForceNew) The type of the batch task. Valid values: `polarclaw_install_skills`, `polarclaw_uninstall_skills`.
* `instance_ids` - (Required, ForceNew) A list of PolarDB cluster IDs to which the task will be applied.
* `task_params` - (Required, ForceNew) The parameters for the task. See [`task_params`](#task_params) below.
* `task_name` - (Optional, ForceNew) The name of the batch task.
* `region_id` - (Optional, ForceNew) The region ID where the PolarDB clusters are located. If not specified, the provider's region is used.

### `task_params`

The `task_params` block supports the following:

* `skill_name` - (Required) The name of the skill to be installed or uninstalled. For example, `polar_claw`.
* `version` - (Optional) The version of the skill. If not specified, the latest version may be used depending on the task type.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Batch Task.
* `batch_id` - The ID of the batch operation returned by the API.
* `task_status` - The status of the task. Valid values may include `Running`, `Success`, `Failed`, etc.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 50 mins) Used when creating the polardb batch task (until the task reaches a final status).
* `delete` - (Defaults to 10 mins) Used when deleting the polardb batch task resource from state.

## Import

PolarDB Batch Task can be imported using the id, e.g.

```shell
$ terraform import alicloud_polardb_batch_task.example pcb-abc12345678
```