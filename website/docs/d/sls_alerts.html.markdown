---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sls_alerts"
sidebar_current: "docs-alicloud-datasource-sls-alerts"
description: |-
  Provides a list of Sls Alert owned by an Alibaba Cloud account.
---

# alicloud_sls_alerts

This data source provides Sls Alert available to the user.[What is Alert](https://next.api.alibabacloud.com/document/Sls/2020-12-30/CreateAlert)

-> **NOTE:** Available since v1.250.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-beijing"
}

variable "alert_name" {
  default = "openapi-terraform-alert"
}

variable "project_name" {
  default = "terraform-alert-example"
}

resource "alicloud_log_project" "defaultINsMgl" {
  description = "terraform-alert-example"
  name        = "terraform-alert-example"
}


resource "alicloud_sls_alert" "default" {
  configuration {
    type    = "tpl"
    version = "2"
    query_list {
      query          = "* | select *"
      time_span_type = "Relative"
      start          = "-15m"
      end            = "now"
      store_type     = "log"
      project        = alicloud_log_project.defaultINsMgl.id
      store          = "alert"
      region         = "cn-beijing"
      power_sql_mode = "disable"
      chart_title    = "wkb-chart"
      dashboard_id   = "wkb-dashboard"
      ui             = "{}"
      role_arn       = "acs:ram::1654218965343050:role/aliyunslsalertmonitorrole"
    }
    query_list {
      store_type = "meta"
      store      = "user.rds_ip_whitelist"
    }
    query_list {
      store_type = "meta"
      store      = "myexample1"
    }
    group_configuration {
      type   = "no_group"
      fields = ["a", "b"]
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
    severity_configurations {
      severity = "6"
      eval_condition {
        condition       = "__count__ > 1"
        count_condition = "cnt > 0"
      }
    }
    labels {
      key   = "a"
      value = "b"
    }
    annotations {
      key   = "x"
      value = "y"
    }
    auto_annotation = true
    send_resolved   = false
    threshold       = "1"
    no_data_fire    = false
    sink_event_store {
      enabled     = true
      endpoint    = "cn-shanghai-intranet.log.aliyuncs.com"
      project     = "wkb-wangren"
      event_store = "alert"
      role_arn    = "acs:ram::1654218965343050:role/aliyunlogetlrole"
    }
    sink_cms {
      enabled = false
    }
    sink_alerthub {
      enabled = false
    }
    template_configuration {
      template_id = "sls.app.ack.autoscaler.cluster_unhealthy"
      type        = "sys"
      version     = "1.0"
      lang        = "cn"
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
    dashboard        = "internal-alert"
    mute_until       = "0"
    no_data_severity = "6"
    tags             = ["wkb", "wangren", "sls"]
  }
  alert_name   = var.alert_name
  project_name = alicloud_log_project.defaultINsMgl.id
  schedule {
    type            = "Cron"
    run_immdiately  = true
    time_zone       = "+0800"
    delay           = "10"
    cron_expression = "0/5 * * * *"
  }
  display_name = "openapi-terraform"
  description  = "create alert"
}

data "alicloud_sls_alerts" "default" {
  ids          = ["${alicloud_sls_alert.default.id}"]
  name_regex   = alicloud_sls_alert.default.alert_name
  project_name = alicloud_log_project.defaultINsMgl.id
}

output "alicloud_sls_alert_example_id" {
  value = data.alicloud_sls_alerts.default.alerts.0.id
}
```

## Argument Reference

The following arguments are supported:
* `project_name` - (Required, ForceNew) Project Name
* `ids` - (Optional, ForceNew, Computed) A list of Alert IDs. The value is formulated as `<project_name>:<alert_name>`.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Group Metric Rule name.
* `output_file` - (Optional, ForceNew) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Alert IDs.
* `names` - A list of name of Alerts.
* `alerts` - A list of Alert Entries. Each element contains the following attributes:
  * `alert_name` - Alert rule ID, unique under Project
  * `configuration` - Detailed configuration of alarm monitoring rules
    * `annotations` - Label.
      * `key` - Field name.
      * `value` - Field value.
    * `auto_annotation` - Whether to turn on automatic labeling.true (default): The automatic annotation function is enabled, and the system automatically adds information such as__county__to the alarm. For more information, see Automatic Labeling.false: Turn off the automatic annotation function.
    * `condition_configuration` - Alarm trigger condition.
      * `condition` - Data matching expression.
      * `count_condition` - The number of data items to determine the expression, indicating the number of data items that meet the conditions. Indicates that if there is data that is satisfied, it can be set to an empty string. In other cases, you need to set it to an expression, for example,__count __> 3.
    * `dashboard` - The instrument cluster associated with the alarm. It is recommended to set to internal-alert-analysis.
    * `group_configuration` - Group evaluation configuration.
      * `fields` - The field used for grouping evaluation.When type is set to custom, fields must be set.
      * `type` - Group evaluation type. The values are as follows:no_group: not grouped.custom: custom.labels_auto: label automatic.Only applicable to time series data.
    * `join_configurations` - Set operation configuration.
      * `condition` - When type is set to inner_join, left_join, right_join, full_join, left_exclude, or right_exclude, condition must be set, for example, set to $0.host = = $1.ip.
      * `type` - Collection operation type. The values are as follows:cross_join: Cartesian product.inner_join: inline.left_join: left joint.right_join: right link.full_join: full union.left_exclude: left repulsion.right_exclude: right repulsion.concat: stitching, traversing each dataset sequentially.no_join: not merged, only the first dataset is taken.
    * `labels` - Label.
      * `key` - Field name.
      * `value` - Field value.
    * `mute_until` - Second-level timestamp representing the temporary shutdown deadline.
    * `no_data_fire` - Whether no data triggers an alarm.true: If the number of times the query and analysis results (if there are multiple results, the result after the collection operation) is no data exceeds the continuous trigger threshold, an alarm is generated.false (default): Turn off the no data alarm function.
    * `no_data_severity` - Alarm severity when no data triggers an alarm.
    * `policy_configuration` - Alert policy configuration.
      * `action_policy_id` - The ID of the action policy used.If the alert policy is in advanced mode and the selected alert policy is not configured with a dynamic action policy, set the actionPolicyId to an empty string.
      * `alert_policy_id` - Alarm policy ID.If it is in simple mode or normal mode, set it to sls.builtin.dynamic (dynamic alarm policy built into the system).For advanced mode, set it to the specified alarm policy ID.
      * `repeat_interval` - Repeat the waiting time. For example, 5m, 1H, etc.
    * `query_list` - Query the statistical list.
      * `chart_title` - Chart Name.
      * `dashboard_id` - Dashboard ID.
      * `end` - End time. When storeType is set to log or metric, end must be set.
      * `power_sql_mode` - Whether to use exclusive SQL. The value is as follows: auto: automatic switching. enable: Starts. disable: disable.
      * `project` - Query the Project associated with the statistics.
      * `query` - Query and analysis statements. When storeType is set to log or metric, query is set to the query analysis statement. When storeType is set to meta, set query to an empty string.
      * `region` - Region of the target Project.
      * `role_arn` - The ARN of the RAM role required to access the data.
      * `start` - Start time. When storeType is set to log or metric, start must be set.
      * `store` - Query the Logstore, Metricstore, or resource data associated with the statistics. When storeType is set to log, store is set to the target Logstore. When storeType is set to metric, store is set to the target Metricstore. When storeType is set to meta, store is set to the target resource data name.
      * `store_type` - Query the data source type. The value is as follows: log: Logstore. metric: Time series Library. meta: resource data.
      * `time_span_type` - Time Type.
      * `ui` - Use of specific scene alarm front end.
    * `send_resolved` - Whether to send a recovery notification.true: A recovery alarm is triggered when the alarm is restored.false (default): Turn off the alarm recovery notification function.
    * `severity_configurations` - Trigger condition, set at least one trigger condition.
      * `eval_condition` - Trigger condition.
        * `condition` - Data matching expression.When the data content does not need to be determined, set it to an empty string.In other cases, it needs to be set as an expression, for example, errCnt> 10.
        * `count_condition` - The number of pieces of data to determine the number of pieces of data to indicate how many pieces of data meet the conditions.If data exists, it is satisfied. Set it to an empty string.In other cases, it needs to be set as an expression, such as__count__> 3.
      * `severity` - Alarm severity.
    * `sink_alerthub` - Configuration of Alerts Sent to Alerthub.
      * `enabled` - Open.
    * `sink_cms` - Configure alerts sent to CloudMonitor.
      * `enabled` - Open.
    * `sink_event_store` - Configuration of sending alarms to EventStore.
      * `enabled` - Open.
      * `endpoint` - SLS service endpoint.
      * `event_store` - Event Library Name.
      * `project` - Project Name.
      * `role_arn` - Roles used to write alarm data to the event Library.
    * `tags` - Customize the category of alarm monitoring rules.
    * `template_configuration` - Alarm rule template configuration.
      * `annotations` - Template Annotations.
      * `lang` - Template Language.
      * `template_id` - Template ID.
      * `tokens` - Template Variables.
      * `type` - Template Type.
      * `version` - Template Version.
    * `threshold` - Set the continuous trigger threshold. When the cumulative number of triggers reaches this value, an alarm is generated. The statistics are not counted when the trigger condition is not met.
    * `type` - Fixed as default.
    * `version` - Fixed as 2.0.
  * `description` - Compatible fields, set to empty strings.
  * `display_name` - Display name of the alarm rule
  * `schedule` - Check the frequency-dependent configuration
    * `cron_expression` - Cron expression, the minimum accuracy is minutes, 24 hours. For example, 0 0/1 * * * means that the check is conducted every 1 hour from 00:00.When type is set to Cron, cronExpression must be set.
    * `delay` - Timed task execution delay (unit: s).
    * `interval` - Fixed interval for scheduling.
    * `run_immdiately` - Dispatch immediately.
    * `time_zone` - The time zone where the Cron expression is located. The default value is null, indicating the eighth zone in the east.
    * `type` - Check the frequency type. Log Service checks the query and analysis results according to the frequency you configured. The values are as follows:Fixedate: checks query and analysis results at regular intervals.Cron: specifies the time interval by using the Cron expression, and checks the query and analysis results at the specified time interval.
  * `id` - The ID of the resource supplied above.
