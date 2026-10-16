---
subcategory: "Pai Flow"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_flow_pipeline_run"
description: |-
  Provides a Alicloud Pai Flow Pipeline Run resource.
---

# alicloud_pai_flow_pipeline_run

Provides a Pai Flow Pipeline Run resource. A pipeline run is a one-shot
execution of a pipeline: Create triggers the run and waits for it to settle
into a terminal status (Succeeded, Failed, Terminated or Skipped). The managed
state is the run's existence, so the outcome is reflected in the Computed
`status` attribute rather than in apply success.

For information about Pai Flow Pipeline Run and how to use it, see [What is PipelineRun](https://www.alibabacloud.com/help/en/pai/developer-reference/api-paiflow-2021-02-02-createpipelinerun).

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

data "alicloud_account" "default" {}

resource "alicloud_pai_workspace_workspace" "default" {
  description    = var.name
  display_name   = var.name
  workspace_name = var.name
  env_types      = ["dev"]
}

resource "alicloud_pai_flow_pipeline" "default" {
  workspace_id = alicloud_pai_workspace_workspace.default.id
  manifest     = <<EOF
apiVersion: "core/v1"
metadata:
  provider: "${data.alicloud_account.default.id}"
  version: "v1"
  identifier: "my_pipeline"
  name: "source-transform"
  uuid: "terraformuuid"
  annotations: {}
  labels: {}
spec:
  inputs:
    artifacts: []
    parameters: []
  outputs:
    artifacts: []
    parameters: []
  arguments:
    artifacts: []
    parameters: []
  dependencies: []
  initContainers: []
  sideCarContainers: []
  pipelines: []
  volumes: []
EOF
}

resource "alicloud_pai_flow_pipeline_run" "default" {
  workspace_id        = alicloud_pai_workspace_workspace.default.id
  pipeline_id         = alicloud_pai_flow_pipeline.default.id
  name                = var.name
  no_confirm_required = true
}
```

## Argument Reference

The following arguments are supported:
* `workspace_id` - (Required, ForceNew) The ID of the workspace.
* `pipeline_id` - (Optional, ForceNew) The ID of the pipeline to run. Required unless `manifest` is supplied as an inline pipeline definition.
* `manifest` - (Optional, ForceNew) The manifest of the pipeline run, in YAML. An inline pipeline definition used when `pipeline_id` is not supplied.
* `arguments` - (Optional, ForceNew) The run arguments, as a JSON string.
* `name` - (Optional) The name of the pipeline run. Updatable.
* `no_confirm_required` - (Optional, ForceNew) Whether to skip the run confirmation. Defaults to false. Create-only; not read back by GetPipelineRun.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource. Same as `pipeline_run_id`.
* `pipeline_run_id` - The ID of the pipeline run. It is the unique key of the resource.
* `status` - The status of the pipeline run. Valid values: `Initialized`, `Running`, `Succeeded`, `Failed`, `Suspended`, `Terminated`, `Unknown`, `Skipped`, `Terminating`.
* `gmt_create_time` - The time when the pipeline run was created.
* `gmt_modified_time` - The time when the pipeline run was last modified.

## Import

Pai Flow Pipeline Run can be imported using the id, e.g.

```shell
$ terraform import alicloud_pai_flow_pipeline_run.example <id>
```
