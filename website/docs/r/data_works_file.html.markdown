---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_data_works_file"
sidebar_current: "docs-alicloud-resource-data-works-file"
description: |-
  Provides a Alicloud Data Works File resource.
---

# alicloud\_data\_works\_file

Provides a Data Works File resource. A DataWorks file is a development artifact
(for example a SQL, Shell or Data Integration node) under a DataWorks
workspace. It maps to the `File` entity of the DataWorks OpenAPI.

For information about DataWorks File and how to use it, see
[CreateFile](https://www.alibabacloud.com/help/en/dataworks/developer-reference/api-dataworks-public-2020-05-18-createfile).

-> **NOTE:** Available since v1.235.0.

## Example Usage

Basic Usage

```hcl
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

resource "alicloud_data_works_file" "example" {
  project_id = "320687"
  file_name  = "tf_test_file"
  file_type  = "10"
  content    = "SHOW TABLES;"
}
```

## Argument Reference

The following arguments are supported:

* `file_name` - (Required) The name of the file.
* `file_type` - (Required, ForceNew) The code type of the file. Different file types use different code. You can call [ListFileType](https://www.alibabacloud.com/help/en/dataworks/developer-reference/api-dataworks-public-2020-05-18-listfiletype) to query the supported code types. For example, `10` indicates an ODPS SQL node.
* `project_id` - (Optional, ForceNew) The ID of the DataWorks workspace. You must set one of `project_id` and `project_identifier` to identify the workspace.
* `project_identifier` - (Optional, ForceNew) The identifier of the DataWorks workspace. You must set one of `project_id` and `project_identifier`.
* `content` - (Optional) The code content of the file. The format varies by file type.
* `owner` - (Optional) The Alibaba Cloud user ID of the file owner. If empty, the caller user ID is used.
* `file_description` - (Optional) The description of the file.
* `file_folder_path` - (Optional) The path of the file.
* `connection_name` - (Optional) The data source connected when the file is published as a task.
* `cycle_type` - (Optional) The scheduling cycle type. Valid values: `DAY`, `NOT_DAY`. Default: `DAY`.
* `dependent_type` - (Optional) The way to depend on the previous cycle. Valid values: `SELF`, `CHILD`, `USER_DEFINE`, `NONE`, `USER_DEFINE_AND_SELF`, `CHILD_AND_SELF`. Default: `NONE`.
* `scheduler_type` - (Optional) The scheduling type. Valid values: `NORMAL`, `MANUAL`, `PAUSE`, `SKIP`.
* `resource_group_identifier` - (Optional) The identifier of the resource group used to run the task.
* `auto_rerun_times` - (Optional) The number of automatic reruns after a task fails.
* `auto_rerun_interval_millis` - (Optional) The interval between automatic reruns, in milliseconds.
* `rerun_mode` - (Optional) The rerun mode.
* `stop` - (Optional) Whether to stop the task.
* `para_value` - (Optional) The parameter value used to run the task.
* `start_effect_date` - (Optional) The start time of the effective period of the scheduling.
* `end_effect_date` - (Optional) The end time of the effective period of the scheduling.
* `cron_express` - (Optional) The CRON expression for the scheduling.
* `dependent_node_id_list` - (Optional) The list of dependent node IDs.
* `input_list` - (Optional) The input list of the file.
* `output_list` - (Optional) The output list of the file.
* `auto_parsing` - (Optional) Whether to automatically parse the code.
* `advanced_settings` - (Optional) The advanced settings of the scheduling.
* `start_immediately` - (Optional) Whether to start the task immediately after it is published.
* `input_parameters` - (Optional) The input parameters of the file.
* `output_parameters` - (Optional) The output parameters of the file.
* `ignore_parent_skip_running_property` - (Optional) Whether to ignore the parent skip running property.
* `apply_schedule_immediately` - (Optional) Whether to apply the scheduling configuration immediately.
* `timeout` - (Optional) The timeout period for the task, in minutes.
* `image_id` - (Optional) The custom image ID.
* `deploy` - (Optional) Whether to call DeployFile to deploy the file after create or update. Default: `false`.
* `submit` - (Optional) Whether to call SubmitFile to submit the file after create or update. Default: `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of the file. The value formats as `<file_id>:<project_id>`.
* `file_id` - The ID of the file.
* `node_id` - The ID of the node generated after the file is published.

## Import

Data Works File can be imported using the id, e.g.

```shell
$ terraform import alicloud_data_works_file.example <file_id>:<project_id>
```
