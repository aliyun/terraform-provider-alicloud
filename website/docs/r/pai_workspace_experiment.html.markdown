---
subcategory: "PAI Workspace"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_workspace_experiment"
description: |-
  Provides a Alicloud PAI Workspace Experiment resource.
---

# alicloud_pai_workspace_experiment

Provides a PAI Workspace Experiment resource.



For information about PAI Workspace Experiment and how to use it, see [What is Experiment](https://next.api.alibabacloud.com/document/AIWorkSpace/2021-02-04/CreateExperiment).

-> **NOTE:** Available since v1.236.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_pai_workspace_experiment&exampleId=24283ef1-ac2b-a17f-8da4-6a364a2880b8daf3f17d&activeTab=example&spm=docs.r.pai_workspace_experiment.0.24283ef1ac&intl_lang=EN_US" target="_blank">
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

resource "alicloud_pai_workspace_workspace" "defaultDI9fsL" {
  description    = var.name
  display_name   = var.name
  workspace_name = var.name
  env_types      = ["prod"]
}


resource "alicloud_pai_workspace_experiment" "default" {
  accessibility   = "PRIVATE"
  artifact_uri    = "oss://yyt-409262.oss-cn-hangzhou.aliyuncs.com/example/"
  experiment_name = var.name
  workspace_id    = alicloud_pai_workspace_workspace.defaultDI9fsL.id
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_pai_workspace_experiment&spm=docs.r.pai_workspace_experiment.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `accessibility` - (Optional, Computed) Experimental Visibility
* `artifact_uri` - (Required, ForceNew) ArtifactUri is default OSS storage path of the output of trials in the experiment
* `experiment_name` - (Required) Name is the name of the experiment, unique in a namespace
* `workspace_id` - (Required, ForceNew) WorkspaceId is the workspace id which contains the experiment

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - GmtCreateTime is time when this entity is created.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Experiment.
* `delete` - (Defaults to 5 mins) Used when delete the Experiment.
* `update` - (Defaults to 5 mins) Used when update the Experiment.

## Import

PAI Workspace Experiment can be imported using the id, e.g.

```shell
$ terraform import alicloud_pai_workspace_experiment.example <id>
```