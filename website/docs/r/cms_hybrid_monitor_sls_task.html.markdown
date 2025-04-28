---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_hybrid_monitor_sls_task"
sidebar_current: "docs-alicloud-resource-cms-hybrid-monitor-sls-task"
description: |-
  Provides a Alicloud Cloud Monitor Service Hybrid Monitor Sls Task resource.
---

# alicloud_cms_hybrid_monitor_sls_task

Provides a Cloud Monitor Service Hybrid Monitor Sls Task resource.

For information about Cloud Monitor Service Hybrid Monitor Sls Task and how to use it, see [What is Hybrid Monitor Sls Task](https://www.alibabacloud.com/help/en/cloudmonitor/latest/createhybridmonitortask).

-> **NOTE:** Available since v1.179.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cms_hybrid_monitor_sls_task&exampleId=1889b555-04b8-1245-2111-5e323084845c1c708dc9&activeTab=example&spm=docs.r.cms_hybrid_monitor_sls_task.0.1889b55504&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf_example"
}
data "alicloud_account" "default" {}
data "alicloud_regions" "default" {
  current = true
}
resource "random_uuid" "default" {
}
resource "alicloud_log_project" "default" {
  project_name = substr("tf-example-${replace(random_uuid.default.result, "-", "")}", 0, 16)
}

resource "alicloud_log_store" "default" {
  project_name          = alicloud_log_project.default.project_name
  logstore_name         = var.name
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}

resource "alicloud_cms_sls_group" "default" {
  sls_group_config {
    sls_user_id  = data.alicloud_account.default.id
    sls_logstore = alicloud_log_store.default.logstore_name
    sls_project  = alicloud_log_project.default.project_name
    sls_region   = data.alicloud_regions.default.regions.0.id
  }
  sls_group_description = var.name
  sls_group_name        = var.name
}

resource "alicloud_cms_namespace" "default" {
  namespace     = substr("tf-example-${replace(random_uuid.default.result, "-", "")}", 0, 16)
  specification = "cms.s1.large"
}

resource "alicloud_cms_hybrid_monitor_sls_task" "default" {
  task_name           = var.name
  namespace           = alicloud_cms_namespace.default.id
  description         = var.name
  collect_interval    = 60
  collect_target_type = alicloud_cms_sls_group.default.id

  sls_process_config {
    filter {
      relation = "and"
      filters {
        operator     = "="
        value        = "200"
        sls_key_name = "code"
      }
    }
    statistics {
      function      = "count"
      alias         = "level_count"
      sls_key_name  = "name"
      parameter_one = "200"
      parameter_two = "299"
    }
    group_by {
      alias        = "code"
      sls_key_name = "ApiResult"
    }
    express {
      express = "success_count"
      alias   = "SuccRate"
    }
  }

  attach_labels {
    name  = "app_service"
    value = "example_Value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `attach_labels` - (Optional) The label of the monitoring task. See [`attach_labels`](#attach_labels) below. 
* `collect_interval` - (Optional) The interval at which metrics are collected. Valid values: `15`, `60`(default value). Unit: seconds.
* `collect_target_type` - (Required, ForceNew) The type of the collection target, enter the name of the Logstore group.
* `description` - (Optional) The description of the metric import task.
* `namespace` - (Required, ForceNew) The name of the namespace.
* `sls_process_config` - (Required) The configurations of the logs that are imported from Log Service. See [`sls_process_config`](#sls_process_config) below. 
* `task_name` - (Required, ForceNew) The name of the metric import task, enter the name of the metric for logs imported from Log Service.

### `sls_process_config`

The sls_process_config supports the following: 

* `express` - (Optional) The extended fields that specify the results of basic operations that are performed on aggregation results. See [`express`](#sls_process_config-express) below. 
* `filter` - (Optional) The conditions that are used to filter logs imported from Log Service. See [`filter`](#sls_process_config-filter) below. 
* `group_by` - (Optional) The dimension based on which data is aggregated. This parameter is equivalent to the GROUP BY clause in SQL. See [`group_by`](#sls_process_config-group_by) below. 
* `statistics` - (Optional) The method that is used to aggregate logs imported from Log Service. See [`statistics`](#sls_process_config-statistics) below. 

### `sls_process_config-statistics`

The statistics supports the following: 

* `alias` - (Optional) The alias of the aggregation result.
* `function` - (Optional) The function that is used to aggregate log data within a statistical period. Valid values: `count`, `sum`, `avg`, `max`, `min`, `value`, `countps`, `sumps`, `distinct`, `distribution`, `percentile`.
* `parameter_one` - (Optional) The value of the function that is used to aggregate logs imported from Log Service.
  - If you set the `function` parameter to `distribution`, this parameter specifies the lower limit of the statistical interval. For example, if you want to calculate the number of HTTP requests whose status code is 2XX, set this parameter to 200.
  - If you set the `function` parameter to `percentile`, this parameter specifies the percentile at which the expected value is. For example, 0.5 specifies P50.
* `parameter_two` - (Optional) The value of the function that is used to aggregate logs imported from Log Service. **Note:** This parameter is required only if the `function` parameter is set to `distribution`. This parameter specifies the upper limit of the statistical interval.
* `sls_key_name` - (Optional) The name of the key that is used to aggregate logs imported from Log Service.

### `sls_process_config-group_by`

The group_by supports the following: 

* `alias` - (Optional) The alias of the aggregation result.
* `sls_key_name` - (Optional) The name of the key that is used to aggregate logs imported from Log Service.

### `sls_process_config-filter`

The filter supports the following: 

* `filters` - (Optional) The conditions that are used to filter logs imported from Log Service. See [`filters`](#sls_process_config-filter-filters) below. 
* `relation` - (Optional) The relationship between multiple filter conditions. Valid values: `and`(default value), `or`.

### `sls_process_config-filter-filters`

The filters supports the following: 

* `operator` - (Optional) The method that is used to filter logs imported from Log Service. Valid values: `>`, `>=`, `=`, `<=`, `<`, `!=`, `contain`, `notContain`.
* `sls_key_name` - (Optional) The name of the key that is used to filter logs imported from Log Service.
* `value` - (Optional) The value of the key that is used to filter logs imported from Log Service.

### `sls_process_config-express`

The express supports the following: 

* `alias` - (Optional) The alias of the extended field that specifies the result of basic operations that are performed on aggregation results.
* `express` - (Optional) The extended field that specifies the result of basic operations that are performed on aggregation results.

### `attach_labels`

The attach_labels supports the following: 

* `name` - (Optional) The tag key of the metric.
* `value` - (Optional) The tag value of the metric.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Hybrid Monitor Sls Task.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 2 mins) Used when create the Hybrid Monitor Sls Task.
* `delete` - (Defaults to 2 mins) Used when delete the Hybrid Monitor Sls Task.
* `update` - (Defaults to 2 mins) Used when update the Hybrid Monitor Sls Task.

## Import

Cloud Monitor Service Hybrid Monitor Sls Task can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_hybrid_monitor_sls_task.example <id>
```