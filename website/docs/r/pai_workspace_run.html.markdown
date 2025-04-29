---
subcategory: "PAI Workspace"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_workspace_run"
description: |-
  Provides a Alicloud PAI Workspace Run resource.
---

# alicloud_pai_workspace_run

Provides a PAI Workspace Run resource.



For information about PAI Workspace Run and how to use it, see [What is Run](https://next.api.alibabacloud.com/document/AIWorkSpace/2021-02-04/CreateRun).

-> **NOTE:** Available since v1.236.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_pai_workspace_run&exampleId=e7b8d398-6264-ef90-4a74-91a5f6316c6d888e541a&activeTab=example&spm=docs.r.pai_workspace_run.0.e7b8d39862&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Run.
* `delete` - (Defaults to 5 mins) Used when delete the Run.
* `update` - (Defaults to 5 mins) Used when update the Run.

## Import

PAI Workspace Run can be imported using the id, e.g.

```shell
$ terraform import alicloud_pai_workspace_run.example <id>
```