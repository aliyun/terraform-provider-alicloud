---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sls_alert"
description: |-
  Provides a Alicloud SLS Alert resource.
---

# alicloud_sls_alert

Provides a SLS Alert resource. 

For information about SLS Alert and how to use it, see [What is Alert](https://www.alibabacloud.com/help/en/doc-detail/209202.html).

-> **NOTE:** Available since v1.223.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_sls_alert&exampleId=6ab20446-3b87-b005-9f15-e6b5ded6ed461c1915db&activeTab=example&spm=docs.r.sls_alert.0.6ab204463b&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

variable "alert_name" {
  default = "openapi-terraform-alert"
}

variable "project_name" {
  default = "terraform-alert-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_log_project" "defaultINsMgl" {
  description  = "${var.project_name}-${random_integer.default.result}"
  project_name = "${var.project_name}-${random_integer.default.result}"
}

resource "alicloud_sls_alert" "default" {
  schedule {
    type           = "FixedRate"
    run_immdiately = "true"
    interval       = "1m"
    time_zone      = "+0800"
    delay          = "10"
  }

  display_name = "openapi-terraform"
  description  = "create alert"
  status       = "ENABLED"
  configuration {
    group_configuration {
      fields = [
        "a",
        "b"
      ]
      type = "no_group"
    }

    no_data_fire = "false"
    version      = "2"
    severity_configurations {
      severity = "6"
      eval_condition {
        count_condition = "cnt > 0"
        condition       = "__count__ > 1"
      }

    }

    labels {
      key   = "a"
      value = "b"
    }

    auto_annotation = "true"
    template_configuration {
      lang = "cn"
      tokens = {
        "a" = "b"
      }
      annotations = {
        "x" = "y"
      }
      template_id = "sls.app.ack.autoscaler.cluster_unhealthy"
      type        = "sys"
      version     = "1.0"
    }

    mute_until = "0"
    annotations {
      key   = "x"
      value = "y"
    }

    send_resolved = "false"
    threshold     = "1"
    sink_cms {
      enabled = "false"
    }

    condition_configuration {
      condition       = "cnt > 3"
      count_condition = "__count__ < 3"
    }

    policy_configuration {
      alert_policy_id  = "sls.builtin.dynamic"
      action_policy_id = "wkb-action"
      repeat_interval  = "1m"
    }

    dashboard = "internal-alert"
    type      = "tpl"
    query_list {
      ui             = "{}"
      role_arn       = "acs:ram::1654218965343050:role/aliyunslsalertmonitorrole"
      query          = "* | select *"
      time_span_type = "Relative"
      project        = alicloud_log_project.defaultINsMgl.project_name
      power_sql_mode = "disable"
      dashboard_id   = "wkb-dashboard"
      chart_title    = "wkb-chart"
      start          = "-15m"
      end            = "now"
      store_type     = "log"
      store          = "alert"
      region         = "cn-shanghai"
    }
    query_list {
      store_type = "meta"
      store      = "user.rds_ip_whitelist"
    }
    query_list {
      store_type = "meta"
      store      = "myexample1"
    }

    join_configurations {
      type      = "no_join"
      condition = "aa"
    }
    join_configurations {
      type      = "cross_join"
      condition = "qqq"
    }
    join_configurations {
      type      = "inner_join"
      condition = "fefefe"
    }

    sink_event_store {
      enabled     = "true"
      endpoint    = "cn-shanghai-intranet.log.aliyuncs.com"
      project     = "wkb-wangren"
      event_store = "alert"
      role_arn    = "acs:ram::1654218965343050:role/aliyunlogetlrole"
    }

    sink_alerthub {
      enabled = "false"
    }

    no_data_severity = "6"
    tags = [
      "wkb",
      "wangren",
      "sls"
    ]
  }

  alert_name   = var.alert_name
  project_name = alicloud_log_project.defaultINsMgl.project_name
}
```

## Argument Reference

The following arguments are supported:
* `alert_name` - (Required, ForceNew) Alert rule ID, unique under Project.
* `configuration` - (Required) Detailed configuration of alarm monitoring rules. See [`configuration`](#configuration) below.
* `description` - (Optional) Compatible fields, set to empty strings.
* `display_name` - (Required) Display name of the alarm rule.
* `project_name` - (Required, ForceNew) Project Name.
* `schedule` - (Required) Check the frequency-dependent configuration. See [`schedule`](#schedule) below.
* `status` - (Optional, Computed) Resource attribute field representing alarm status.

### `configuration`

The configuration supports the following:
* `annotations` - (Optional) Label. See [`annotations`](#configuration-annotations) below.
* `auto_annotation` - (Optional, Computed) Whether to turn on automatic labeling. true (default): The automatic annotation function is enabled, and the system automatically adds information such as__county__to the alarm. For more information, see Automatic Labeling. false: Turn off the automatic annotation function.
* `condition_configuration` - (Optional) Alarm trigger condition. See [`condition_configuration`](#configuration-condition_configuration) below.
* `dashboard` - (Optional) The instrument cluster associated with the alarm. It is recommended to set to internal-alert-analysis.
* `group_configuration` - (Optional) Group evaluation configuration. See [`group_configuration`](#configuration-group_configuration) below.
* `join_configurations` - (Optional) Set operation configuration. See [`join_configurations`](#configuration-join_configurations) below.
* `labels` - (Optional) Label. See [`labels`](#configuration-labels) below.
* `mute_until` - (Optional) Second-level timestamp representing the temporary shutdown deadline.
* `no_data_fire` - (Optional, Computed) Whether no data triggers an alarm. true: If the number of times the query and analysis results (if there are multiple results, the result after the collection operation) is no data exceeds the continuous trigger threshold, an alarm is generated. false (default): Turn off the no data alarm function.
* `no_data_severity` - (Optional) Alarm severity when no data triggers an alarm.
* `policy_configuration` - (Optional) Alert policy configuration. See [`policy_configuration`](#configuration-policy_configuration) below.
* `query_list` - (Optional) Query the statistical list. See [`query_list`](#configuration-query_list) below.
* `send_resolved` - (Optional, Computed) Whether to send a recovery notification. true: A recovery alarm is triggered when the alarm is restored. false (default): Turn off the alarm recovery notification function.
* `severity_configurations` - (Optional) Trigger condition, set at least one trigger condition. See [`severity_configurations`](#configuration-severity_configurations) below.
* `sink_alerthub` - (Optional) Configuration of Alerts Sent to Alerthub. See [`sink_alerthub`](#configuration-sink_alerthub) below.
* `sink_cms` - (Optional) Configure alerts sent to CloudMonitor. See [`sink_cms`](#configuration-sink_cms) below.
* `sink_event_store` - (Optional) Configuration of sending alarms to EventStore. See [`sink_event_store`](#configuration-sink_event_store) below.
* `tags` - (Optional) Customize the category of alarm monitoring rules.
* `template_configuration` - (Optional) Alarm rule template configuration. See [`template_configuration`](#configuration-template_configuration) below.
* `threshold` - (Optional) Set the continuous trigger threshold. When the cumulative number of triggers reaches this value, an alarm is generated. The statistics are not counted when the trigger condition is not met.
* `type` - (Optional, Computed) Fixed as default.
* `version` - (Optional) Fixed as 2.0.

### `configuration-annotations`

The configuration-annotations supports the following:
* `key` - (Optional) Field name.
* `value` - (Optional) Field value.

### `configuration-condition_configuration`

The configuration-condition_configuration supports the following:
* `condition` - (Optional) Data matching expression.
* `count_condition` - (Optional) The number of data items to determine the expression, indicating the number of data items that meet the conditions. Indicates that if there is data that is satisfied, it can be set to an empty string. In other cases, you need to set it to an expression, for example,__count __> 3.

### `configuration-group_configuration`

The configuration-group_configuration supports the following:
* `fields` - (Optional) The field used for grouping evaluation. When type is set to custom, fields must be set.
* `type` - (Optional) Group evaluation type. The values are as follows: no_group: not grouped. custom: custom. labels_auto: label automatic. Only applicable to time series data.

### `configuration-join_configurations`

The configuration-join_configurations supports the following:
* `condition` - (Optional) When type is set to inner_join, left_join, right_join, full_join, left_exclude, or right_exclude, condition must be set, for example, set to $0.host = = $1.ip.
* `type` - (Optional) Collection operation type. The values are as follows: cross_join: Cartesian product. inner_join: inline. left_join: left joint. right_join: right link. full_join: full union. left_exclude: left repulsion. right_exclude: right repulsion. concat: stitching, traversing each dataset sequentially. no_join: not merged, only the first dataset is taken.

### `configuration-labels`

The configuration-labels supports the following:
* `key` - (Optional) Field name.
* `value` - (Optional) Field value.

### `configuration-policy_configuration`

The configuration-policy_configuration supports the following:
* `action_policy_id` - (Optional) The ID of the action policy used. If the alert policy is in advanced mode and the selected alert policy is not configured with a dynamic action policy, set the actionPolicyId to an empty string.
* `alert_policy_id` - (Optional) Alarm policy ID. If it is in simple mode or normal mode, set it to sls.builtin.dynamic (dynamic alarm policy built into the system). For advanced mode, set it to the specified alarm policy ID.
* `repeat_interval` - (Optional) Repeat the waiting time. For example, 5m, 1H, etc.

### `configuration-query_list`

The configuration-query_list supports the following:
* `chart_title` - (Optional) Chart Name.
* `dashboard_id` - (Optional) Dashboard ID.
* `end` - (Optional) End time. When storeType is set to log or metric, end must be set.
* `power_sql_mode` - (Optional) Whether to use exclusive SQL. The value is as follows: auto: automatic switching. enable: Starts. disable: disable.
* `project` - (Optional) Query the Project associated with the statistics.
* `query` - (Optional) Query and analysis statements. When storeType is set to log or metric, query is set to the query analysis statement. When storeType is set to meta, set query to an empty string.
* `region` - (Optional) Region of the target Project.
* `role_arn` - (Optional) The ARN of the RAM role required to access the data.
* `start` - (Optional) Start time. When storeType is set to log or metric, start must be set.
* `store` - (Optional) Query the Logstore, Metricstore, or resource data associated with the statistics. When storeType is set to log, store is set to the target Logstore. When storeType is set to metric, store is set to the target Metricstore. When storeType is set to meta, store is set to the target resource data name.
* `store_type` - (Optional) Query the data source type. The value is as follows: log: Logstore. metric: Time series Library. meta: resource data.
* `time_span_type` - (Optional) Time Type.
* `ui` - (Optional) Use of specific scene alarm front end.

### `configuration-severity_configurations`

The configuration-severity_configurations supports the following:
* `eval_condition` - (Optional) Trigger condition. See [`eval_condition`](#configuration-severity_configurations-eval_condition) below.
* `severity` - (Optional) Alarm severity.

### `configuration-sink_alerthub`

The configuration-sink_alerthub supports the following:
* `enabled` - (Optional) Open.

### `configuration-sink_cms`

The configuration-sink_cms supports the following:
* `enabled` - (Optional) Open.

### `configuration-sink_event_store`

The configuration-sink_event_store supports the following:
* `enabled` - (Optional) Open.
* `endpoint` - (Optional) SLS service endpoint.
* `event_store` - (Optional) Event Library Name.
* `project` - (Optional) Project Name.
* `role_arn` - (Optional) Roles used to write alarm data to the event Library.

### `configuration-template_configuration`

The configuration-template_configuration supports the following:
* `annotations` - (Optional, Map) Template Annotations.
* `lang` - (Optional) Template Language.
* `template_id` - (Optional) Template ID.
* `tokens` - (Optional, Map) Template Variables.
* `type` - (Optional) Template Type.
* `version` - (Optional) Template Version.

### `configuration-severity_configurations-eval_condition`

The configuration-severity_configurations-eval_condition supports the following:
* `condition` - (Optional) Data matching expression. When the data content does not need to be determined, set it to an empty string. In other cases, it needs to be set as an expression, for example, errCnt> 10.
* `count_condition` - (Optional) The number of pieces of data to determine the number of pieces of data to indicate how many pieces of data meet the conditions. If data exists, it is satisfied. Set it to an empty string. In other cases, it needs to be set as an expression, such as__count__> 3.

### `schedule`

The schedule supports the following:
* `cron_expression` - (Optional) Cron expression, the minimum accuracy is minutes, 24 hours. For example, 0 0/1 * * * means that the check is conducted every 1 hour from 00:00. When type is set to Cron, cronExpression must be set.
* `delay` - (Optional) Timed task execution delay (unit: s).
* `interval` - (Optional) Fixed interval for scheduling.
* `run_immdiately` - (Optional) Dispatch immediately.
* `time_zone` - (Optional) The time zone where the Cron expression is located. The default value is null, indicating the eighth zone in the east.
* `type` - (Optional) Check the frequency type. Log Service checks the query and analysis results according to the frequency you configured. The values are as follows: Fixedate: checks query and analysis results at regular intervals. Cron: specifies the time interval by using the Cron expression, and checks the query and analysis results at the specified time interval.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<project_name>:<alert_name>`.
* `create_time` - Alarm rule creation time.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Alert.
* `delete` - (Defaults to 5 mins) Used when delete the Alert.
* `update` - (Defaults to 5 mins) Used when update the Alert.

## Import

SLS Alert can be imported using the id, e.g.

```shell
$ terraform import alicloud_sls_alert.example <project_name>:<alert_name>
```