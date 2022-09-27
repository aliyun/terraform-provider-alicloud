---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_hybrid_monitor_sls_task"
sidebar_current: "docs-alicloud-resource-cms-hybrid-monitor-sls-task"
description: |-
  Provides a Alicloud Cloud Monitor Service Hybrid Monitor Sls Task resource.
---

# alicloud\_cms\_hybrid\_monitor\_sls\_task

Provides a Cloud Monitor Service Hybrid Monitor Sls Task resource.

For information about Cloud Monitor Service Hybrid Monitor Sls Task and how to use it, see [What is Hybrid Monitor Sls Task](https://www.alibabacloud.com/help/en/cloudmonitor/latest/createhybridmonitortask).

-> **NOTE:** Available in v1.179.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_account" "this" {}

resource "alicloud_cms_sls_group" "default" {
  sls_group_config {
    sls_user_id  = data.alicloud_account.this.id
    sls_logstore = "Logstore-ECS"
    sls_project  = "aliyun-project"
    sls_region   = "cn-hangzhou"
  }
  sls_group_description = "example_value"
  sls_group_name        = "example_value"
}
resource "alicloud_cms_namespace" "default" {
  description   = var.name
  namespace     = "example-value"
  specification = "cms.s1.large"
}
resource "alicloud_cms_hybrid_monitor_sls_task" "default" {
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
  task_name           = "example_value"
  namespace           = alicloud_cms_namespace.default.id
  description         = "example_value"
  collect_interval    = 60
  collect_target_type = alicloud_cms_sls_group.default.id
  attach_labels {
    name  = "app_service"
    value = "testValue"
  }
}
```

## Argument Reference

The following arguments are supported:

* `attach_labels` - (Optional) The label of the monitoring task. See the following `Block attach_labels`.
* `collect_interval` - (Optional, Computed) The interval at which metrics are collected. Valid values: `15`, `60`(default value). Unit: seconds.
* `collect_target_type` - (Required, ForceNew) The type of the collection target, enter the name of the Logstore group.
* `description` - (Optional) The description of the metric import task.
* `namespace` - (Required, ForceNew) The name of the namespace.
* `sls_process_config` - (Required) The configurations of the logs that are imported from Log Service. See the following `Block sls_process_config`.
* `task_name` - (Required, ForceNew) The name of the metric import task, enter the name of the metric for logs imported from Log Service.

#### Block sls_process_config

The sls_process_config supports the following: 

* `express` - (Optional) The extended fields that specify the results of basic operations that are performed on aggregation results. See the following `Block express`.
* `filter` - (Optional) The conditions that are used to filter logs imported from Log Service. See the following `Block filter`.
* `group_by` - (Optional) The dimension based on which data is aggregated. This parameter is equivalent to the GROUP BY clause in SQL. See the following `Block group_by`.
* `statistics` - (Optional) The method that is used to aggregate logs imported from Log Service. See the following `Block statistics`.

#### Block statistics

The statistics supports the following: 

* `alias` - (Optional) The alias of the aggregation result.
* `function` - (Optional) The function that is used to aggregate log data within a statistical period. Valid values: `count`, `sum`, `avg`, `max`, `min`, `value`, `countps`, `sumps`, `distinct`, `distribution`, `percentile`.
* `parameter_one` - (Optional) The value of the function that is used to aggregate logs imported from Log Service.
  - If you set the `function` parameter to `distribution`, this parameter specifies the lower limit of the statistical interval. For example, if you want to calculate the number of HTTP requests whose status code is 2XX, set this parameter to 200.
  - If you set the `function` parameter to `percentile`, this parameter specifies the percentile at which the expected value is. For example, 0.5 specifies P50.
* `parameter_two` - (Optional) The value of the function that is used to aggregate logs imported from Log Service. **Note:** This parameter is required only if the `function` parameter is set to `distribution`. This parameter specifies the upper limit of the statistical interval.
* `sls_key_name` - (Optional) The name of the key that is used to aggregate logs imported from Log Service.

#### Block group_by

The group_by supports the following: 

* `alias` - (Optional) The alias of the aggregation result.
* `sls_key_name` - (Optional) The name of the key that is used to aggregate logs imported from Log Service.

#### Block filter

The filter supports the following: 

* `filters` - (Optional) The conditions that are used to filter logs imported from Log Service. See the following `Block filters`.
* `relation` - (Optional) The relationship between multiple filter conditions. Valid values: `and`(default value), `or`.

#### Block filters

The filters supports the following: 

* `operator` - (Optional) The method that is used to filter logs imported from Log Service. Valid values: `>`, `>=`, `=`, `<=`, `<`, `!=`, `contain`, `notContain`.
* `sls_key_name` - (Optional) The name of the key that is used to filter logs imported from Log Service.
* `value` - (Optional) The value of the key that is used to filter logs imported from Log Service.

#### Block express

The express supports the following: 

* `alias` - (Optional) The alias of the extended field that specifies the result of basic operations that are performed on aggregation results.
* `express` - (Optional) The extended field that specifies the result of basic operations that are performed on aggregation results.

#### Block attach_labels

The attach_labels supports the following: 

* `name` - (Optional) The tag key of the metric.
* `value` - (Optional) The tag value of the metric.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Hybrid Monitor Sls Task.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 2 mins) Used when create the Hybrid Monitor Sls Task.
* `delete` - (Defaults to 2 mins) Used when delete the Hybrid Monitor Sls Task.
* `update` - (Defaults to 2 mins) Used when update the Hybrid Monitor Sls Task.

## Import

Cloud Monitor Service Hybrid Monitor Sls Task can be imported using the id, e.g.

```
$ terraform import alicloud_cms_hybrid_monitor_sls_task.example <id>
```