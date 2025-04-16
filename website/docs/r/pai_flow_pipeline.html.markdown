---
subcategory: "Pai Flow"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_flow_pipeline"
description: |-
  Provides a Alicloud Pai Flow Pipeline resource.
---

# alicloud_pai_flow_pipeline

Provides a Pai Flow Pipeline resource.



For information about Pai Flow Pipeline and how to use it, see [What is Pipeline](https://www.alibabacloud.com/help/en/pai/developer-reference/api-paiflow-2021-02-02-createpipeline).

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

data "alicloud_account" "default" {
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_pai_workspace_workspace" "defaultWjQD1e" {
  description    = var.name
  display_name   = var.name
  workspace_name = var.name
  env_types      = ["dev"]
}

resource "alicloud_pai_flow_pipeline" "default" {
  workspace_id = alicloud_pai_workspace_workspace.defaultWjQD1e.id
  manifest     = <<EOF
apiVersion: "core/v1"
metadata:
  provider: "${data.alicloud_account.default.id}"
  version: "v1"
  identifier: "my_pipeline"
  name: "source-transform"
  uuid: "terraformuuid${random_integer.default.result}"
  annotations: {}
  labels: {}
spec:
  inputs:
    artifacts: []
    parameters:
    - name: "execution_maxcompute"
      type: "Map"
      value:
        spec:
          endpoint: "http://service.cn.maxcompute.aliyun-inc.com/api"
          odpsProject: "test_i****"
  outputs:
    artifacts: []
    parameters: []
  arguments:
    artifacts: []
    parameters: []
  dependencies: []
  initContainers: []
  sideCarContainers: []
  pipelines:
  - apiVersion: "core/v1"
    metadata:
      provider: "pai"
      version: "v1"
      identifier: "data_source"
      name: "data-source"
      uuid: "2ftahdnzcod2rt6u9q"
      displayName: "读数据表-1"
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
        parameters:
        - name: "inputTableName"
          value: "pai_online_project.wumai_data"
        - name: "execution"
          from: "{{inputs.parameters.execution_maxcompute}}"
      dependencies: []
      initContainers: []
      sideCarContainers: []
      pipelines: []
      volumes: []
  - apiVersion: "core/v1"
    metadata:
      provider: "pai"
      version: "v1"
      identifier: "type_transform"
      name: "type-transform"
      uuid: "gacnnnl4ksvbabfh6l"
      displayName: "类型转换-1"
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
        artifacts:
        - name: "inputTable"
          from: "{{pipelines.data_source.outputs.artifacts.outputTable}}"
        parameters:
        - name: "cols_to_double"
          value: "time,hour,pm2,pm10,so2,co,no2"
        - name: "execution"
          from: "{{inputs.parameters.execution_maxcompute}}"
      dependencies:
      - "data_source"
      initContainers: []
      sideCarContainers: []
      pipelines: []
      volumes: []
  volumes: []
EOF
}
```

## Argument Reference

The following arguments are supported:
* `manifest` - (Required) The pipeline definition. For more information, see the [sample pipeline definition]((https://www.alibabacloud.com/help/en/pai/developer-reference/api-paiflow-2021-02-02-createpipeline)).
* `workspace_id` - (Required, ForceNew) The ID of the workspace.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Pipeline.
* `delete` - (Defaults to 5 mins) Used when delete the Pipeline.
* `update` - (Defaults to 5 mins) Used when update the Pipeline.

## Import

Pai Flow Pipeline can be imported using the id, e.g.

```shell
$ terraform import alicloud_pai_flow_pipeline.example <id>
```
