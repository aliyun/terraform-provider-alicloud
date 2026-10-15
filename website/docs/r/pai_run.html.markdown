---
subcategory: "Pai Flow"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_run"
description: |-
  Provides a Alicloud Pai Flow Pipeline Run resource.
---

# alicloud_pai_run

Provides a PAI Flow Pipeline Run resource. A pipeline run is a single execution of a pipeline; creating the resource submits and starts the run.

For information about PAI Flow Pipeline Run and how to use it, see [CreatePipelineRun](https://www.alibabacloud.com/help/en/pai/developer-reference/api-paiflow-2021-02-02-createpipelinerun).

-> **NOTE:** Available since v1.286.0.

## Example Usage

Basic Usage - run an existing pipeline by ID

```terraform
variable "name" {
  default = "terraform_example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

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

resource "alicloud_pai_run" "default" {
  workspace_id = alicloud_pai_workspace_workspace.default.id
  pipeline_id  = alicloud_pai_flow_pipeline.default.id
}
```

## Argument Reference

The following arguments are supported:

* `workspace_id` - (Required, ForceNew) The ID of the workspace.
* `pipeline_id` - (Optional, ForceNew) The ID of the pipeline to run. Exactly one of `pipeline_id` and `pipeline_manifest` must be set.
* `pipeline_manifest` - (Optional, ForceNew) The manifest of the pipeline to run, in YAML. Exactly one of `pipeline_id` and `pipeline_manifest` must be set.
* `name` - (Optional, ForceNew) The name of the pipeline run. If omitted, PAI Flow generates one.
* `arguments` - (Optional, ForceNew) The run arguments in YAML.
* `options` - (Optional, ForceNew) The run options in JSON.
* `source_type` - (Optional, ForceNew) The source type of the pipeline run. Valid values: `UNKNOWN`, `SDK`, `DESIGNER`, `M6`. Default: `UNKNOWN`.
* `source_id` - (Optional, ForceNew) The source ID of the pipeline run.
* `accessibility` - (Optional, ForceNew) The accessibility of the pipeline run. Valid values: `PUBLIC`, `PRIVATE`. Default: `PUBLIC`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource. Same as the pipeline run ID.
* `status` - The status of the pipeline run, e.g. `Running`, `Succeeded`, `Failed`, `Stopped`.
* `manifest` - The manifest of the pipeline run.
* `pipeline_run_uri` - The URI of the pipeline run.
* `started_at` - The time when the pipeline run started.
* `finished_at` - The time when the pipeline run finished.
* `duration` - The duration of the pipeline run, in milliseconds.
* `node_id` - The node ID of the pipeline run.
* `user_id` - The user ID that started the pipeline run.
* `parent_user_id` - The parent user ID of the pipeline run.
* `message` - The message of the pipeline run, populated on failure.
* `gmt_create_time` - The creation time of the pipeline run.
* `gmt_modified_time` - The modification time of the pipeline run.

## Import

Pai Flow Pipeline Run can be imported using the id, e.g.

```shell
$ terraform import alicloud_pai_run.example <pipeline_run_id>
```
