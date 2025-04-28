---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_data_works_di_job"
description: |-
  Provides a Alicloud Data Works Di Job resource.
---

# alicloud_data_works_di_job

Provides a Data Works Di Job resource.

Data Integration Tasks.

For information about Data Works Di Job and how to use it, see [What is Di Job](https://www.alibabacloud.com/help/en/dataworks/developer-reference/api-dataworks-public-2024-05-18-createdijob).

-> **NOTE:** Available since v1.241.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_data_works_di_job&exampleId=308b052b-c8d0-78bc-f2b2-013cd620d55b0d91ebae&activeTab=example&spm=docs.r.data_works_di_job.0.308b052bc8&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-chengdu"
}

resource "alicloud_data_works_project" "defaultMMHL8U" {
  description  = var.name
  project_name = var.name
  display_name = var.name
}


resource "alicloud_data_works_di_job" "default" {
  description    = var.name
  project_id     = alicloud_data_works_project.defaultMMHL8U.id
  job_name       = "zhenyuan_example_case"
  migration_type = "api_FullAndRealtimeIncremental"
  source_data_source_settings {
    data_source_name = "dw_mysql"
    data_source_properties {
      encoding = "utf-8"
      timezone = "Asia/Shanghai"
    }
  }
  destination_data_source_type = "Hologres"
  table_mappings {
    source_object_selection_rules {
      action          = "Include"
      expression      = "dw_mysql"
      expression_type = "Exact"
      object_type     = "Datasource"
    }
    source_object_selection_rules {
      action          = "Include"
      expression      = "example_db1"
      expression_type = "Exact"
      object_type     = "Database"
    }
    source_object_selection_rules {
      action          = "Include"
      expression      = "lsc_example01"
      expression_type = "Exact"
      object_type     = "Table"
    }
    transformation_rules {
      rule_name        = "my_table_rename_rule"
      rule_action_type = "Rename"
      rule_target_type = "Table"
    }
  }
  source_data_source_type = "MySQL"
  resource_settings {
    offline_resource_settings {
      requested_cu              = 2
      resource_group_identifier = "S_res_group_524257424564736_1716799673667"
    }
    realtime_resource_settings {
      requested_cu              = 2
      resource_group_identifier = "S_res_group_524257424564736_1716799673667"
    }
    schedule_resource_settings {
      requested_cu              = 2
      resource_group_identifier = "S_res_group_524257424564736_1716799673667"
    }
  }
  transformation_rules {
    rule_action_type = "Rename"
    rule_expression  = "{\"expression\":\"table2\"}"
    rule_name        = "my_table_rename_rule"
    rule_target_type = "Table"
  }
  destination_data_source_settings {
    data_source_name = "dw_example_holo"
  }
  job_settings {
    column_data_type_settings {
      destination_data_type = "bigint"
      source_data_type      = "longtext"
    }
    ddl_handling_settings {
      action = "Ignore"
      type   = "CreateTable"
    }
    runtime_settings {
      name  = "runtime.realtime.concurrent"
      value = "1"
    }
    channel_settings = "1"
    cycle_schedule_settings {
      cycle_migration_type = "2"
      schedule_parameters  = "3"
    }
  }
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Optional) Description of the integration task
* `destination_data_source_settings` - (Required, ForceNew, List) Destination data source See [`destination_data_source_settings`](#destination_data_source_settings) below.
* `destination_data_source_type` - (Required, ForceNew) The type of the target data source. Enumerated values: Hologres and Hive.
* `job_name` - (Required, ForceNew) Task Name.
* `job_settings` - (Optional, ForceNew, List) The dimension settings of the synchronization task, including the DDL processing policy, the source and destination column data type mapping policy, and the task runtime parameters. See [`job_settings`](#job_settings) below.
* `migration_type` - (Required, ForceNew) Synchronization type, optional enumeration values are:

  Fulllandrealtimeincremental (full and real-time incremental)

  RealtimeIncremental

  Full

  Offflineincremental

  FullAndOfflineIncremental (full amount + offline increment)
* `project_id` - (Optional, ForceNew, Computed, Int) Project Id
* `resource_settings` - (Required, ForceNew, List) Resource Group Properties See [`resource_settings`](#resource_settings) below.
* `source_data_source_settings` - (Required, ForceNew, List) Source data source setting List See [`source_data_source_settings`](#source_data_source_settings) below.
* `source_data_source_type` - (Required, ForceNew) The type of the source data source. The enumerated value is MySQL.
* `table_mappings` - (Required, List) Synchronize object transformation mapping list See [`table_mappings`](#table_mappings) below.
* `transformation_rules` - (Optional, List) Definition list of synchronization object conversion rules See [`transformation_rules`](#transformation_rules) below.

### `destination_data_source_settings`

The destination_data_source_settings supports the following:
* `data_source_name` - (Optional, ForceNew) Destination data source name

### `job_settings`

The job_settings supports the following:
* `channel_settings` - (Optional) Channel-related task settings, in the form of a Json String.

  For example, 
{"structInfo":"MANAGED","storageType":"TEXTFILE","writeMode":"APPEND","partitionColumns":[{"columnName":"pt","columnType":"STRING","comment":""}],"fieldDelimiter":""}
* `column_data_type_settings` - (Optional, List) Column type mapping of the synchronization task See [`column_data_type_settings`](#job_settings-column_data_type_settings) below.
* `cycle_schedule_settings` - (Optional, List) Periodic scheduling settings See [`cycle_schedule_settings`](#job_settings-cycle_schedule_settings) below.
* `ddl_handling_settings` - (Optional, List) List of DDL processing settings for synchronization tasks See [`ddl_handling_settings`](#job_settings-ddl_handling_settings) below.
* `runtime_settings` - (Optional, List) Run-time setting parameter list See [`runtime_settings`](#job_settings-runtime_settings) below.

### `job_settings-column_data_type_settings`

The job_settings-column_data_type_settings supports the following:
* `destination_data_type` - (Optional) The destination type of the mapping relationship
* `source_data_type` - (Optional) The source type of the mapping type

### `job_settings-cycle_schedule_settings`

The job_settings-cycle_schedule_settings supports the following:
* `cycle_migration_type` - (Optional, ForceNew) The type of synchronization that requires periodic scheduling. Value range:

  Full: Full

  OfflineIncremental: offline increment
* `schedule_parameters` - (Optional) Scheduling Parameters

### `job_settings-ddl_handling_settings`

The job_settings-ddl_handling_settings supports the following:
* `action` - (Optional) Processing action, optional enumeration value:

  Ignore (Ignore)

  Critical (error)

  Normal (Normal processing)
* `type` - (Optional) DDL type, optional enumeration value:

  RenameColumn (rename column)

  ModifyColumn (rename column)

  CreateTable (Rename Column)

  TruncateTable (empty table)

  DropTable (delete table)

### `job_settings-runtime_settings`

The job_settings-runtime_settings supports the following:
* `name` - (Optional) Set name, optional ENUM value:

  runtime.offline.speed.limit.mb (valid when runtime.offline.speed.limit.enable = true)

  runtime.offline.speed.limit.enable

  dst.offline.connection.max (the maximum number of write connections for offline batch tasks)

  runtime.offline.concurrent (offline batch synchronization task concurrency)

  dst.realtime.connection.max (maximum number of write connections for real-time tasks)

  runtime.enable.auto.create.schema (whether to automatically create a schema on the target side)

  src.offline.datasource.max.connection (maximum number of source connections for offline batch tasks)

  runtime.realtime.concurrent (real-time task concurrency)
* `value` - (Optional) Runtime setting value

### `resource_settings`

The resource_settings supports the following:
* `offline_resource_settings` - (Optional, List) Offline Resource Group configuration See [`offline_resource_settings`](#resource_settings-offline_resource_settings) below.
* `realtime_resource_settings` - (Optional, List) Real-time Resource Group See [`realtime_resource_settings`](#resource_settings-realtime_resource_settings) below.
* `schedule_resource_settings` - (Optional, List) Scheduling Resource Groups See [`schedule_resource_settings`](#resource_settings-schedule_resource_settings) below.

### `resource_settings-offline_resource_settings`

The resource_settings-offline_resource_settings supports the following:
* `requested_cu` - (Optional, Float) Offline resource group cu
* `resource_group_identifier` - (Optional) Offline resource group name

### `resource_settings-realtime_resource_settings`

The resource_settings-realtime_resource_settings supports the following:
* `requested_cu` - (Optional, Float) Real-time resource group cu
* `resource_group_identifier` - (Optional) Real-time resource group name

### `resource_settings-schedule_resource_settings`

The resource_settings-schedule_resource_settings supports the following:
* `requested_cu` - (Optional, Float) Scheduling resource group cu
* `resource_group_identifier` - (Optional) Scheduling resource group name

### `source_data_source_settings`

The source_data_source_settings supports the following:
* `data_source_name` - (Optional, ForceNew) Data source name of a single source
* `data_source_properties` - (Optional, ForceNew, List) Single Source Data Source Properties See [`data_source_properties`](#source_data_source_settings-data_source_properties) below.

### `source_data_source_settings-data_source_properties`

The source_data_source_settings-data_source_properties supports the following:
* `encoding` - (Optional, ForceNew) Data Source Encoding
* `timezone` - (Optional, ForceNew) Data Source Time Zone

### `table_mappings`

The table_mappings supports the following:
* `source_object_selection_rules` - (Optional, List) Each rule can select different types of source objects to be synchronized, such as source database and source data table. See [`source_object_selection_rules`](#table_mappings-source_object_selection_rules) below.
* `transformation_rules` - (Optional, List) A list of conversion rule definitions for a synchronization object. Each element in the list defines a conversion rule. See [`transformation_rules`](#table_mappings-transformation_rules) below.

### `table_mappings-source_object_selection_rules`

The table_mappings-source_object_selection_rules supports the following:
* `action` - (Optional) Select an action. Value range: Include/Exclude
* `expression` - (Optional) Expression, such as mysql_table_1
* `expression_type` - (Optional) Expression type, value range: Exact/Regex
* `object_type` - (Optional) Object type, optional enumeration value:

  Table (Table)

  Database

### `table_mappings-transformation_rules`

The table_mappings-transformation_rules supports the following:
* `rule_action_type` - (Optional) Action type, optional enumeration value:

  DefinePrimaryKey (defines the primary key)

  Rename

  AddColumn (increase column)

  HandleDml(DML handling)

  DefineIncrementalCondition

  DefineCycleScheduleSettings (defines periodic scheduling settings)

  DefineRuntimeSettings (defines advanced configuration parameters)

  DefinePartitionKey (defines partition column)
* `rule_name` - (Optional) The rule name, which is unique under an action type + the target type of the action action.
* `rule_target_type` - (Optional) Target type of action, optional enumeration value:

  Table (Table)

  Schema(schema)

### `transformation_rules`

The transformation_rules supports the following:
* `rule_action_type` - (Optional) Action type, optional enumeration value:

  DefinePrimaryKey (defines the primary key)

  Rename

  AddColumn (increase column)

  HandleDml(DML handling)

  DefineIncrementalCondition
* `rule_expression` - (Optional) Regular expression, in json string format.

  Example renaming rule (Rename): {"expression":"${srcDatasourceName}_${srcDatabaseName}_0922","variables":[{"variableName":"srcDatabaseName","variableRules":[{"from":"fromdb","to":"todb"}]}]}
* `rule_name` - (Optional) Rule Name
* `rule_target_type` - (Optional) Target type of action, optional enumeration value:

  Table (Table)

  Schema(schema)

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<project_id>:<di_job_id>`.
* `di_job_id` - Integration Task Id

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Di Job.
* `delete` - (Defaults to 5 mins) Used when delete the Di Job.
* `update` - (Defaults to 5 mins) Used when update the Di Job.

## Import

Data Works Di Job can be imported using the id, e.g.

```shell
$ terraform import alicloud_data_works_di_job.example <project_id>:<di_job_id>
```