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

* `name` - (Required) The name of the data quality scan.
* `project_id` - (Required, ForceNew) The ID of the DataWorks project to which the scan belongs.
* `spec` - (Required) The scan specification, a JSON document describing the datasets to scan and the rules to evaluate. See the [DataQualityScan API](https://www.alibabacloud.com/help/en/dataworks/developer-reference/api-dataworks-public-2024-05-18-createdataqualityscan) for the schema of datasets and rules.
* `description` - (Optional) A description of the scan.
* `owner` - (Optional) The owner of the scan.
* `parameters` - (Optional) A list of key/value parameters for the scan. See [`parameters`](#parameters) below.
* `trigger` - (Optional) The scan trigger configuration. See [`trigger`](#trigger) below.
* `compute_resource` - (Optional) The compute resource used to run the scan. See [`compute_resource`](#compute_resource) below.
* `runtime_resource` - (Optional) The runtime resource used to run the scan. See [`runtime_resource`](#runtime_resource) below.
* `hooks` - (Optional) A list of hooks invoked when the scan produces results. See [`hooks`](#hooks) below.

### `parameters`

The `parameters` block supports:

* `name` - (Required) The parameter name.
* `value` - (Required) The parameter value.

### `trigger`

The `trigger` block supports:

* `type` - (Required) The trigger type, e.g. `BySchedule`.
* `task_ids` - (Optional) The list of task IDs the scan is associated with.

### `compute_resource`

The `compute_resource` block supports:

* `name` - (Optional) The compute resource name.
* `env_type` - (Optional) The environment type, e.g. `Prod` or `Dev`.
* `runtime` - (Optional) The compute runtime configuration. See [`runtime`](#runtime) below.

#### `runtime`

The `runtime` block (under `compute_resource`) supports:

* `engine` - (Optional) The compute engine, e.g. `Hive` or `Spark`.
* `hive_conf` - (Optional, Map, ForceNew) Hive configuration key/value pairs.
* `spark_conf` - (Optional, Map, ForceNew) Spark configuration key/value pairs.

### `runtime_resource`

The `runtime_resource` block supports:

* `id` - (Optional) The runtime resource ID.
* `cu` - (Optional) The compute units allocated to the runtime resource.
* `image` - (Optional) The image used by the runtime resource.

### `hooks`

The `hooks` block supports:

* `condition` - (Required) The condition expression that triggers the hook.
* `type` - (Required) The hook type, e.g. `BlockTaskInstance`.

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
