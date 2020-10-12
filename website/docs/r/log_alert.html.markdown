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

-> **NOTE:** Available in 1.78.0

## Example Usage

Basic Usage

```
resource "alicloud_log_project" "example" {
  name        = "test-tf"
  description = "create by terraform"
}

resource "alicloud_log_store" "example" {
  project               = alicloud_log_project.example.name
  name                  = "tf-test-logstore"
  retention_period      = 3650
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}

resource "alicloud_log_alert" "example" {
  project_name      = alicloud_log_project.example.name
  alert_name        = "tf-test-alert"
  alert_displayname = "tf-test-alert-displayname"
  condition         = "count> 100"
  dashboard         = "tf-test-dashboard"
  query_list {
    logstore    = "tf-test-logstore"
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
    email_list = ["aliyun@alibaba-inc.com", "tf-test@123.com"]
    content    = "alert content"
  }
  notification_list {
    type        = "DingTalk"
    service_uri = "www.aliyun.com"
    content     = "alert content"
  }
}
```
## Argument Reference

The following arguments are supported:

* `project_name` - (Required, ForceNew) The project name.
* `alert_name` - (Required, ForceNew) Name of logstore for configuring alarm service.
* `alert_displayname` - (Required) Alert displayname.
* `alert_description` - (Optional) Alert description.
* `condition` - (Required) Conditional expression, such as: count> 100.
* `dashboard` - (Required) The name of the dashboard associated with the alarm. The name of the instrument cluster associated with the alarm. If there is no such instrument cluster, terraform will help you create an empty instrument cluster.
* `mute_until` - (Optional) Timestamp, notifications before closing again.
* `throttling` - (Optional) Notification interval, default is no interval. Support number + unit type, for example 60s, 1h.
* `notify_threshold` - (Optional) Notification threshold, which is not notified until the number of triggers is reached. The default is 1.
* `query_list` - (Required) Multiple conditions for configured alarm query.
    * `chart_title` - (Required) chart title
    * `logstore` - (Required) Query logstore
    * `query` - (Required) query corresponding to chart. example: * AND aliyun.
    * `start` - (Required) begin time. example: -60s.
    * `end` - (Required) end time. example: 20s.
    * `time_span_type` - (Optional) default Custom. No need to configure this parameter.
* `notification_list` - (Required) Alarm information notification list.
    * `type` - (Required) Notification type. support Email, SMS, DingTalk, MessageCenter.
    * `content` - (Required) Notice content of alarm.
    * `service_uri` - (Optional) Request address.
    * `mobile_list` - (Optional) SMS sending mobile number.
    * `email_list` - (Optional) Email address list.   
* `schedule_interval` - (Optional) Execution interval. 60 seconds minimum, such as 60s, 1h.
* `schedule_type` - (Optional)  Default FixedRate. No need to configure this parameter.

## Attributes Reference

The following attributes are exported:

*  `id` - The ID of the log alert. It formats of `<project>:<alert_name>`.

## Import

Log alert can be imported using the id, e.g.

```
$ terraform import alicloud_log_alert.example tf-log:tf-log-alert
```
