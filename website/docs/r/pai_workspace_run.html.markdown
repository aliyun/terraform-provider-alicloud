---
subcategory: "PAI Workspace"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_workspace_run"
description: |-
  Provides a Alicloud PAI Workspace Run resource.
---

# alicloud_pai_workspace_run

Provides a PAI Workspace Run resource.



For information about PAI Workspace Run and how to use it, see [What is Run](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.236.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform_example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_pai_workspace_workspace" "defaultCAFUa9" {
  description    = var.name
  display_name   = var.name
  workspace_name = var.name
  env_types      = ["prod"]
}

resource "alicloud_pai_workspace_experiment" "defaultQRwWbv" {
  accessibility   = "PRIVATE"
  artifact_uri    = "oss://example.oss-cn-hangzhou.aliyuncs.com/example/"
  experiment_name = format("%s1", var.name)
  workspace_id    = alicloud_pai_workspace_workspace.defaultCAFUa9.id
}


resource "alicloud_pai_workspace_run" "default" {
  source_type   = "TrainingService"
  source_id     = "759"
  run_name      = var.name
  experiment_id = alicloud_pai_workspace_experiment.defaultQRwWbv.id
}
```

## Argument Reference

The following arguments are supported:
* `experiment_id` - (Required, ForceNew) Resource attribute field of the experiment ID to which Run belongs
* `run_name` - (Optional) The name of the resource
* `source_id` - (Optional, ForceNew) Attribute Resource field representing the source task ID
* `source_type` - (Optional, ForceNew) Resource attribute fields representing the source type

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Run.
* `delete` - (Defaults to 5 mins) Used when delete the Run.
* `update` - (Defaults to 5 mins) Used when update the Run.

## Import

PAI Workspace Run can be imported using the id, e.g.

```shell
$ terraform import alicloud_pai_workspace_run.example <id>
```