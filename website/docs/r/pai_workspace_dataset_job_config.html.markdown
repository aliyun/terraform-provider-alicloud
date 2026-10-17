---
subcategory: "PAI Workspace"
layout: "alicloud"
page_title: "Alicloud: alicloud_pai_workspace_dataset_job_config"
description: |-
  Provides a PAI Workspace Dataset Job Config resource.
---

# alicloud_pai_workspace_dataset_job_config

Provides a PAI Workspace Dataset Job Config resource.

For information about PAI Workspace Dataset Job Config and how to use it, see
[What is DatasetJobConfig](https://next.api.alibabacloud.com/document/AIWorkSpace/2021-02-04/CreateDatasetJobConfig).

-> **NOTE:** Available since v1.236.0.

## Example Usage

<div class="oics-button" type="example"><a href="https://api.aliyun.com/api/aiworkspace/2021-02-04/CreateDatasetJobConfig" data-title="Open in AliCloud" data-intro="<div class=&quot;oics-map-content&quot;><p>Open in AliCloud</p></div>"><img src="/img/oics/btn-open-in-alicloud.png" alt="Open in AliCloud" class="oics-button-icon" /></a></div>

```terraform
resource "alicloud_pai_workspace_workspace" "default" {
  workspace_name = "tf-example"
  display_name   = "tf-example"
  env_types      = ["prod"]
}

resource "alicloud_pai_workspace_dataset" "default" {
  workspace_id     = alicloud_pai_workspace_workspace.default.id
  dataset_name     = "tf-example-dataset"
  data_source_type = "NAS"
  data_type        = "PIC"
  property         = "DIRECTORY"
  uri              = "nas://086b649545.cn-hangzhou/"
}

resource "alicloud_pai_workspace_dataset_job_config" "default" {
  workspace_id    = alicloud_pai_workspace_workspace.default.id
  dataset_id      = alicloud_pai_workspace_dataset.default.id
  config_type     = "MultimodalIntelligentTag"
  config          = jsonencode({ apiKey = "sk-xxxxxxxxxxxxxxxxxxxxx" })
  dataset_version = "v1"
}
```

## Argument Reference

The following arguments are supported:

* `workspace_id` - (Required, ForceNew) The workspace ID.
* `dataset_id` - (Required, ForceNew) The dataset ID.
* `config_type` - (Required, ForceNew) The config type. Valid values:
  * `MultimodalIntelligentTag`
  * `MultimodalSemanticIndex`
* `config` - (Required) The config content in JSON format. The structure depends on `config_type`:
  * `MultimodalIntelligentTag`: `{"apiKey":"sk-xxxxxxxxxxxxxxxxxxxxx"}`
  * `MultimodalSemanticIndex`: `{"defaultModelId":"xxx","defaultModelVersion":"1.0.0"}`
* `dataset_version` - (Optional, ForceNew) The dataset version name. Write-only: set at create time and not returned by Read.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in format `<workspace_id>:<dataset_id>:<dataset_job_config_id>`.
* `dataset_job_config_id` - The dataset job config ID, assigned by the server.
* `create_time` - The creation time, in ISO 8601 (UTC) format.
* `modify_time` - The modification time, in ISO 8601 (UTC) format.

## Import

Dataset Job Config can be imported using the id, e.g.

```shell
terraform import alicloud_pai_workspace_dataset_job_config.example <workspace_id>:<dataset_id>:<dataset_job_config_id>
```
