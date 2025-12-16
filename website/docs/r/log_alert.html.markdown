---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_log_alert"
sidebar_current: "docs-alicloud-resource-log-alert"
description: |-
  Provides a Alicloud log alert resource.
---

# alicloud\_log\_alert

Log alert is a unit of log service, which is used to monitor and alert the user's logstore status information. 
Log Service enables you to configure alerts based on the charts in a dashboard to monitor the service status in real time.

For information about SLS Alert and how to use it, see [SLS Alert Overview](https://www.alibabacloud.com/help/en/doc-detail/209202.html)

-> **NOTE:** Available in 1.78.0

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_log_alert&exampleId=95292d19-57da-ce57-7642-9f8788f163adb3198a22&activeTab=example&spm=docs.r.log_alert.0.95292d1957&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_log_project" "example" {
  name        = "terraform-example-${random_integer.default.result}"
  description = "terraform-example"
}

resource "alicloud_log_store" "example" {
  project               = alicloud_log_project.example.name
  name                  = "example-store"
  retention_period      = 3650
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}

resource "alicloud_log_alert" "example" {
  project_name      = alicloud_log_project.example.name
  alert_name        = "example-alert"
  alert_displayname = "example-alert"
  condition         = "count> 100"
  dashboard         = "example-dashboard"
  schedule {
    type            = "FixedRate"
    interval        = "5m"
    hour            = 0
    day_of_week     = 0
    delay           = 0
    run_immediately = false
  }
  query_list {
    logstore    = alicloud_log_store.example.name
    chart_title = "chart_title"
    start       = "-60s"
    end         = "20s"
    query       = "* AND aliyun"
  }
  notification_list {
    type        = "SMS"
    mobile_list = ["12345678", "87654321"]
    content     = "alert content"
  }
  notification_list {
    type       = "Email"
    email_list = ["aliyun@alibaba-inc.com", "tf-example@123.com"]
    content    = "alert content"
  }
  notification_list {
    type        = "DingTalk"
    service_uri = "www.aliyun.com"
    content     = "alert content"
  }
}
```

Basic Usage for new alert

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_log_alert&exampleId=c21c7650-789a-fe6a-59cf-81dc055a64d950d2b31e&activeTab=example&spm=docs.r.log_alert.1.c21c765078&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_log_project" "example" {
  project_name = "terraform-example-${random_integer.default.result}"
  description  = "terraform-example"
}

resource "alicloud_log_store" "example" {
  project_name          = alicloud_log_project.example.project_name
  logstore_name         = "example-store"
  retention_period      = 3650
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}

resource "alicloud_log_alert" "example-2" {
  version           = "2.0"
  type              = "default"
  project_name      = alicloud_log_project.example.project_name
  alert_name        = "example-alert"
  alert_displayname = "example-alert"
  mute_until        = "1632486684"
  no_data_fire      = "false"
  no_data_severity  = 8
  send_resolved     = true
  auto_annotation   = true
  schedule {
    type            = "FixedRate"
    interval        = "5m"
    hour            = 0
    day_of_week     = 0
    delay           = 0
    run_immediately = false
  }
  query_list {
    store          = alicloud_log_store.example.logstore_name
    store_type     = "log"
    project        = alicloud_log_project.example.project_name
    region         = "cn-heyuan"
    chart_title    = "chart_title"
    start          = "-60s"
    end            = "20s"
    query          = "* AND aliyun | select count(1) as cnt"
    power_sql_mode = "auto"
    dashboard_id   = "example-dashboard"
  }
  query_list {
    store          = alicloud_log_store.example.logstore_name
    store_type     = "log"
    project        = alicloud_log_project.example.project_name
    region         = "cn-heyuan"
    chart_title    = "chart_title"
    start          = "-60s"
    end            = "20s"
    query          = "error | select count(1) as error_cnt"
    power_sql_mode = "enable"
    dashboard_id   = "example-dashboard"
  }
  labels {
    key   = "env"
    value = "test"
  }
  annotations {
    key   = "title"
    value = "alert title"
  }
  annotations {
    key   = "desc"
    value = "alert desc"
  }
  annotations {
    key   = "test_key"
    value = "test value"
  }
  group_configuration {
    type   = "custom"
    fields = ["cnt"]
  }
  policy_configuration {
    alert_policy_id  = "sls.bultin"
    action_policy_id = "sls_test_action"
    repeat_interval  = "4h"
  }
  severity_configurations {
    severity = 8
    eval_condition = {
      condition       = "cnt > 3"
      count_condition = "__count__ > 3"
    }
  }
  severity_configurations {
    severity = 6
    eval_condition = {
      condition       = ""
      count_condition = "__count__ > 0"
    }
  }
  severity_configurations {
    severity = 2
    eval_condition = {
      condition       = ""
      count_condition = ""
    }
  }
  join_configurations {
    type      = "cross_join"
    condition = ""
  }
}
```

Basic Usage for alert template

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_log_alert&exampleId=cea67a8f-cee7-ce99-008f-0fa3e60f113cca6b547a&activeTab=example&spm=docs.r.log_alert.2.cea67a8fce&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_log_project" "example" {
  project_name = "terraform-example-${random_integer.default.result}"
  description  = "terraform-example"
}

resource "alicloud_log_store" "example" {
  project_name          = alicloud_log_project.example.project_name
  logstore_name         = "example-store"
  retention_period      = 3650
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}

resource "alicloud_log_alert" "example-3" {
  version           = "2.0"
  type              = "tpl"
  project_name      = alicloud_log_project.example.project_name
  alert_name        = "example-alert"
  alert_displayname = "example-alert"
  mute_until        = "1632486684"
  schedule {
    type            = "FixedRate"
    interval        = "5m"
    hour            = 0
    day_of_week     = 0
    delay           = 0
    run_immediately = false
  }
  template_configuration {
    id          = "sls.app.sls_ack.node.down"
    type        = "sys"
    lang        = "cn"
    annotations = {}
    tokens = {
      "interval_minute"        = "5"
      "default.action_policy"  = "sls.app.ack.builtin"
      "default.severity"       = "6"
      "sendResolved"           = "false"
      "default.project"        = "${alicloud_log_project.example.project_name}"
      "default.logstore"       = "k8s-event"
      "default.repeatInterval" = "4h"
      "trigger_threshold"      = "1"
      "default.clusterId"      = "example-cluster-id"
    }
  }
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_log_alert&spm=docs.r.log_alert.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `version` - (Optional, Available in 1.161.0+) The version of alert, new alert is 2.0.
* `type` - (Optional, Available in 1.161.0+) The type of new alert, `default` for custom alert, `tpl` for template alert.
* `project_name` - (Required, ForceNew) The project name.
* `alert_name` - (Required, ForceNew) Name of logstore for configuring alarm service.
* `alert_displayname` - (Required) Alert displayname.
* `alert_description` - (Optional) Alert description.
* `condition` - (Optional, Deprecated) Conditional expression, such as: count> 100, Deprecated from 1.161.0+.
* `dashboard` - (Optional, Deprecated) The name of the dashboard associated with the alarm. The name of the instrument cluster associated with the alarm. If there is no such instrument cluster, terraform will help you create an empty instrument cluster, Deprecated from 1.161.0+.
* `mute_until` - (Optional) Timestamp, notifications before closing again.
* `throttling` - (Optional, Deprecated) Notification interval, default is no interval. Support number + unit type, for example 60s, 1h, Deprecated from 1.161.0+.
* `notify_threshold` - (Optional, Deprecated) Notification threshold, which is not notified until the number of triggers is reached. The default is 1, Deprecated from 1.161.0+.
* `threshold` - (Optional, Available in 1.161.0+) Evaluation threshold, alert will not fire until the number of triggers is reached. The default is 1.
* `no_data_fire` - (Optional, Available in 1.161.0+) Switch for whether new alert fires when no data happens, default is false.
* `no_data_severity` - (Optional, Available in 1.161.0+) when no data happens, the severity of new alert.
* `send_resolved` - (Optional, Available in 1.161.0+) when new alert is resolved, whether to notify, default is false.
* `auto_annotation` - (Optional, Available in 1.164.0+) whether to add automatic annotation, default is false.
* `query_list` - (Optinal, Required before 1.203.0) Multiple conditions for configured alarm query.
    * `project` - (Optional, Available in 1.161.0+) Query project.
    * `region` - (Optional, Available in 1.161.0+) Query project region.
    * `role_arn` - (Optional) Query project store's ARN.
    * `dashboard_id` - (Optional, Available in 1.161.0+) Query dashboard id.
    * `chart_title` - (Optional) Chart title, optional from 1.161.0+.
    * `logstore` - (Optional, Deprecated) Query logstore, use store for new alert, Deprecated from 1.161.0+.
    * `store` - (Optional, Available in 1.161.0+) Query store for new alert.
    * `store_type` - (Optional, Available in 1.161.0+) Query store type for new alert, including log,metric,meta.
    * `query` - (Required) Query corresponding to chart. example: * AND aliyun.
    * `start` - (Required) Begin time. example: -60s.
    * `end` - (Required) End time. example: 20s.
    * `time_span_type` - (Optional) default Custom. No need to configure this parameter.
    * `power_sql_mode` - (Optional, Available in 1.164.0+) default disable, whether to use power sql. support auto, enable, disable.
* `notification_list` - (Optional, Deprecated) Alarm information notification list, Deprecated from 1.161.0+.
    * `type` - (Required) Notification type. support Email, SMS, DingTalk, MessageCenter.
    * `content` - (Required) Notice content of alarm.
    * `service_uri` - (Optional) Request address.
    * `mobile_list` - (Optional) SMS sending mobile number.
    * `email_list` - (Optional) Email address list.   
* `labels` - (Optional, Available in 1.161.0+) Labels for new alert.
    * `key` - (Required) Labels's key for new alert.
    * `value` - (Required) Labels's value for new alert.
* `annotations` - (Optional, Available in 1.161.0+) Annotations for new alert.
    * `key` - (Required) Annotations's key for new alert.
    * `value` - (Required) Annotations's value for new alert.
* `policy_configuration` - (Optional, Available in 1.161.0+) Policy configuration for new alert.
    * `alert_policy_id` - (Required) Alert Policy Id.
    * `action_policy_id` - (Optional) Action Policy Id.
    * `repeat_interval` - (Required) Repeat interval used by alert policy, 1h, 1m.e.g.
* `group_configuration` - (Optional, Available in 1.161.0+) Group configuration for new alert.
    * `type` - (Optional) Group configuration type, including no_group, labels_auto, custom.
    * `fileds` - (Optional) Group configuration's fields list when type is custom.
* `severity_configurations` - (Optional, Available in 1.161.0+) Severity configuration for new alert.
    * `severity` - (Required) Severity for new alert, including 2,4,6,8,10 for Report,Low,Medium,High,Critical.
    * `eval_condition` - (Required) Severity when this condition is met.
        * `condition` - (Optional) Condition for each row.
        * `count_condition` - (Optional) Count's condition for the rows met condition above.
* `join_configurations` - (Optional, Available in 1.161.0+) Join configuration for different queries.
    * `type` - (Required) Join type, including cross_join, inner_join, left_join, right_join, full_join, left_exclude, right_exclude, concat, no_join.
    * `condition` - (Required) Join condition.
* `template_configuration` - (Optional, Available in 1.203.0+) Template configuration for alert, when `type` is `tpl`.
    * `id` - (Required) Alert template id.
    * `type` - (Required) Alert template type including `sys`, `user`.
    * `lang` - (Optional) Alert template language including `cn`, `en`.
    * `tokens` - (Optional) Alert template tokens.
    * `annotations` - (Optional) Alert template annotations.
* `schedule_interval` - (Optional, Deprecated) Execution interval. 60 seconds minimum, such as 60s, 1h. Deprecated from 1.176.0+. use interval in schedule.
* `schedule_type` - (Optional, Deprecated)  Default FixedRate. No need to configure this parameter. Deprecated from 1.176.0+. use type in schedule.
* `schedule` - (Optional, Available in 1.176.0+) schedule for alert.
    * `type` - (Required) including FixedRate,Hourly,Daily,Weekly,Cron.
    * `interval` - (Optional) Execution interval. 60 seconds minimum, such as 60s, 1h. used when type is FixedRate.
    * `cron_expression` - (Optional) Cron expression when type is Cron.
    * `day_of_week` - (Optional) Day of week when type is Weekly, including 0,1,2,3,4,5,6, 0 for Sunday, 1 for Monday
    * `hour` - (Optional) Hour of day when type is Weekly/Daily.
    * `time_zone` - (Optional) Time zone for schedule.


## Attributes Reference

The following attributes are exported:

*  `id` - The ID of the log alert. It formats of `<project>:<alert_name>`.

## Import

Log alert can be imported using the id, e.g.

```shell
$ terraform import alicloud_log_alert.example tf-log:tf-log-alert
```
