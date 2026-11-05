---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_data_works_data_quality_scan"
description: |-
  Provides a Alicloud Data Works Data Quality Scan resource.
---

# alicloud_data_works_data_quality_scan

Provides a Data Works Data Quality Scan resource.

For information about Data Works Data Quality Scan and how to use it, see [What is Data Quality Scan](https://www.alibabacloud.com/help/en/dataworks/developer-reference/api-dataworks-public-2024-05-18-createdataqualityscan).

-> **NOTE:** Available since v1.246.0.

## Example Usage

Basic Usage

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

variable "name" {
  default = "terraform_example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_data_works_project" "default" {
  description      = "default_pj_dqs"
  project_name     = var.name
  display_name     = "default_pj_dqs"
  pai_task_enabled = false
}

resource "alicloud_data_works_data_quality_scan" "default" {
  name        = var.name
  project_id  = alicloud_data_works_project.default.id
  spec        = <<JSON
{
  "datasets": [
    {
      "type": "Table",
      "dataSource": {"name": "odps_first", "envType": "Prod"},
      "tables": ["ods_d_user_info"],
      "filter": "pt = $[yyyymmdd-1]"
    }
  ],
  "rules": [
    {"assertion": "row_count > 0"}
  ]
}
JSON
  description = "Daily quality scan of ods tables."
  parameters {
    name  = "e2e_autolabel"
    value = "10"
  }
  trigger {
    type     = "BySchedule"
    task_ids = [1023777390]
  }
  hooks {
    condition = "results.any { r -> r.status == 'Fail' && r.rule.severity == 'High' }"
    type      = "BlockTaskInstance"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) The name of the data quality scan.
* `project_id` - (Required, ForceNew, String) The ID of the DataWorks project to which the scan belongs.
* `spec` - (Required, String) The scan specification, a JSON document describing the datasets to scan and the rules to evaluate. See the [DataQualityScan API](https://www.alibabacloud.com/help/en/dataworks/developer-reference/api-dataworks-public-2024-05-18-createdataqualityscan) for the schema of datasets and rules.
* `description` - (Optional, String) A description of the scan.
* `owner` - (Optional, String) The owner of the scan.
* `parameters` - (Optional, Set) A set of key/value parameters for the scan. The order of `parameters` blocks is not significant. See [`parameters`](#parameters) below.
* `trigger` - (Optional, Object) The scan trigger configuration. See [`trigger`](#trigger) below.
* `compute_resource` - (Optional, Object) The compute resource used to run the scan. See [`compute_resource`](#compute_resource) below.
* `runtime_resource` - (Optional, Object) The runtime resource used to run the scan. See [`runtime_resource`](#runtime_resource) below.
* `hooks` - (Optional, Set) A set of hooks invoked when the scan produces results. The order of `hooks` blocks is not significant. See [`hooks`](#hooks) below.

### `parameters`

The `parameters` block supports:

* `name` - (Required, String) The parameter name.
* `value` - (Required, String) The parameter value.

### `trigger`

The `trigger` block supports:

* `type` - (Required, String) The trigger type. Valid values: `ByManual` (default, manual trigger), `BySchedule` (triggered by a scheduled task instance). When set to `BySchedule`, `task_ids` must also be configured.
* `task_ids` - (Optional, List) The list of scheduling task IDs the scan is associated with. Required when `type` is `BySchedule`.

### `compute_resource`

The `compute_resource` block supports:

* `name` - (Optional, String) The compute resource name.
* `env_type` - (Optional, String) The environment type. Valid values: `Prod` (production environment), `Dev` (development environment).
* `runtime` - (Optional, Object) The compute runtime configuration. See [`runtime`](#runtime) below.

#### `runtime`

The `runtime` block (under `compute_resource`) supports:

* `engine` - (Optional, String) The compute engine. Only EMR compute engines support setting this. Valid values: `Hive` (Hive SQL), `Spark` (Spark SQL), `Kyuubi`.
* `hive_conf` - (Optional, Map) Hive configuration key/value pairs. Currently only `mapreduce.job.queuename` is supported to set the queue.
* `spark_conf` - (Optional, Map) Spark configuration key/value pairs. Currently only `spark.yarn.queue` is supported to set the queue.

### `runtime_resource`

The `runtime_resource` block supports:

* `id` - (Optional, String) The runtime resource ID.
* `cu` - (Optional, Float) The compute units allocated to the runtime resource.
* `image` - (Optional, String) The image used by the runtime resource.

### `hooks`

The `hooks` block supports:

* `condition` - (Required, String) The condition expression that triggers the hook.
* `type` - (Required, String) The hook type, e.g. `BlockTaskInstance`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the data quality scan.
* `create_time` - The time the scan was created.
* `modify_time` - The time the scan was last modified.
* `create_user` - The user who created the scan.
* `modify_user` - The user who last modified the scan.

## Import

Data Works Data Quality Scan can be imported using the id, e.g.

```shell
$ terraform import alicloud_data_works_data_quality_scan.example <id>
```
