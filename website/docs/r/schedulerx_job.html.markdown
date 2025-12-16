---
subcategory: "Schedulerx"
layout: "alicloud"
page_title: "Alicloud: alicloud_schedulerx_job"
description: |-
  Provides a Alicloud Schedulerx Job resource.
---

# alicloud_schedulerx_job

Provides a Schedulerx Job resource.



For information about Schedulerx Job and how to use it, see [What is Job](https://www.alibabacloud.com/help/en/schedulerx/schedulerx-serverless/developer-reference/api-schedulerx2-2019-04-30-createjob).

-> **NOTE:** Available since v1.240.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_schedulerx_job&exampleId=322a964d-9f42-7242-e4d8-63cd499b02da5f00da4a&activeTab=example&spm=docs.r.schedulerx_job.0.322a964d9f&intl_lang=EN_US" target="_blank">
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

resource "alicloud_schedulerx_namespace" "CreateNameSpace" {
  namespace_name = var.name
  description    = var.name
}

resource "alicloud_schedulerx_app_group" "default" {
  max_jobs              = "100"
  monitor_contacts_json = jsonencode([{ "userName" : "name1", "userPhone" : "89756******" }, { "userName" : "name2", "ding" : "http://www.example.com" }])
  delete_jobs           = "false"
  app_type              = "1"
  namespace_source      = "schedulerx"
  group_id              = "example-appgroup-pop-autoexample"
  namespace_name        = "default"
  description           = var.name
  monitor_config_json   = jsonencode({ "sendChannel" : "sms,ding" })
  app_version           = "1"
  app_name              = "example-appgroup-pop-autoexample"
  namespace             = alicloud_schedulerx_namespace.CreateNameSpace.namespace_uid
  enable_log            = "false"
  schedule_busy_workers = "false"
}

resource "alicloud_schedulerx_job" "default" {
  timezone        = "GTM+7"
  status          = "Enable"
  max_attempt     = "0"
  description     = var.name
  parameters      = "hello word"
  job_name        = var.name
  max_concurrency = "1"
  time_config {
    data_offset     = "1"
    time_expression = "100000"
    time_type       = "3"
    calendar        = "workday"
  }
  map_task_xattrs {
    task_max_attempt      = "1"
    task_attempt_interval = "1"
    consumer_size         = "5"
    queue_size            = "10000"
    dispatcher_size       = "5"
    page_size             = "100"
  }
  namespace = alicloud_schedulerx_namespace.CreateNameSpace.namespace_uid
  group_id  = alicloud_schedulerx_app_group.default.group_id
  job_type  = "java"
  job_monitor_info {
    contact_info {
      user_phone = "12345678910"
      user_name  = "tangtao-1"
      ding       = "https://alidocs.dingtalk.com"
      user_mail  = "12345678@xx.com"
    }
    contact_info {
      user_phone = "12345678910"
      user_name  = "tangtao-2"
      ding       = "https://alidocs.dingtalk.com1"
      user_mail  = "123456789@xx.com"
    }
    monitor_config {
      timeout             = "7200"
      send_channel        = "sms"
      timeout_kill_enable = true
      timeout_enable      = true
      fail_enable         = true
      miss_worker_enable  = true
    }
  }
  class_name       = "com.aliyun.schedulerx.example.processor.SimpleJob"
  namespace_source = "schedulerx"
  attempt_interval = "30"
  fail_times       = "1"
  execute_mode     = "batch"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_schedulerx_job&spm=docs.r.schedulerx_job.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `attempt_interval` - (Optional, Int) Error retry interval, unit s, default value 30.
* `class_name` - (Optional) Full path of the task interface class.

  This field is available only when your task is of the Java task type.
* `content` - (Optional) The script code of the python, shell, and go task types.
* `description` - (Optional) Task description.
* `execute_mode` - (Required) Task execution mode, the possible results are as follows:
  - `standalone`: stand-alone operation
  - `broadcast`: broadcast running
  - `parallel`: parallel computing
  - `grid`: Memory grid
  - `batch`: Grid Computing
  - `shard`: shard running
* `fail_times` - (Optional, Int) Number of consecutive failed alarms
* `group_id` - (Required, ForceNew) The application ID, which is obtained from the **application management** page of the console.
* `job_monitor_info` - (Optional, List) Task monitoring information See [`job_monitor_info`](#job_monitor_info) below.
* `job_name` - (Required) JobName
* `job_type` - (Required, ForceNew) Job Type
* `map_task_xattrs` - (Optional, List) Advanced configuration, limited to parallel computing, memory grid, and grid computing. See [`map_task_xattrs`](#map_task_xattrs) below.
* `max_attempt` - (Optional, Int) The maximum number of error retries, which is set based on business requirements. The default value is 0.
* `max_concurrency` - (Optional) The maximum number of instances running at the same time. The default value is 1, that is, the last trigger is not completed, and the next trigger will not be performed even at the running time.
* `namespace` - (Required, ForceNew) Namespace, which is obtained on the `Namespace` page of the console.
* `namespace_source` - (Optional) Special third parties are required.
* `parameters` - (Optional) User-defined parameters, which can be obtained at runtime.
* `status` - (Optional, Computed) Task status. The values are as follows:
  - `1`: Enabled and can be triggered normally.
  - `0`: Disabled and will not be triggered.
* `success_notice_enable` - (Optional) Success Notification Switch
* `task_dispatch_mode` - (Optional) Advanced configuration of parallel grid tasks, push model or pull model
* `template` - (Optional) K8s task type custom task template
* `time_config` - (Required, List) Time configuration information See [`time_config`](#time_config) below.
* `timezone` - (Optional) Time Zone
* `x_attrs` - (Optional, Computed) Task Extension Field

### `job_monitor_info`

The job_monitor_info supports the following:
* `contact_info` - (Optional, List) Contact information. See [`contact_info`](#job_monitor_info-contact_info) below.
* `monitor_config` - (Optional, Computed, List) Alarm switch and threshold configuration. See [`monitor_config`](#job_monitor_info-monitor_config) below.

### `job_monitor_info-contact_info`

The job_monitor_info-contact_info supports the following:
* `ding` - (Optional) DingTalk swarm robot webhook address
* `user_mail` - (Optional) User Email Address
* `user_name` - (Optional) The user name
* `user_phone` - (Optional) The user's mobile phone number

### `job_monitor_info-monitor_config`

The job_monitor_info-monitor_config supports the following:
* `fail_enable` - (Optional, Computed) Enable failure alarm
* `miss_worker_enable` - (Optional, Computed) Whether no available Machine alarm is on
* `send_channel` - (Optional, Computed) Alarm sending form
  - sms: sms alarm
  - phone: phone alarm
  - mail: mail alarm
  - webhook:webhook alarm
* `timeout` - (Optional, Computed, Int) Timeout threshold, unit s, default 7200.
* `timeout_enable` - (Optional, Computed) Time-out alarm switch. The values are as follows:
  - `true`: On
  - `false`: closed
* `timeout_kill_enable` - (Optional, Computed) The trigger switch is terminated by timeout and is turned off by default.
  - `true`: On
  - `false`: closed

### `map_task_xattrs`

The map_task_xattrs supports the following:
* `consumer_size` - (Optional, Int) The number of threads to execute a single trigger. The default value is 5.
* `dispatcher_size` - (Optional, Int) The number of subtask distribution threads. The default value is 5.
* `page_size` - (Optional, Int) The number of sub-tasks pulled by a parallel task at a time. The default value is 100.
* `queue_size` - (Optional, Int) The upper limit of the sub-task queue cache. The default value is 10000.
* `task_attempt_interval` - (Optional, Int) Subtask failure retry interval.
* `task_max_attempt` - (Optional, Int) The number of failed sub-task retries.

### `time_config`

The time_config supports the following:
* `calendar` - (Optional) The cron type can optionally fill in a custom calendar.
* `data_offset` - (Optional, Int) Cron type can choose time offset, unit s.
* `time_expression` - (Optional, Computed) Time expressions. Currently, the following types of time expressions are supported:
  - `api`: No time expression.
  - `fix_rate`: the specific fixed frequency value. For example, 30 indicates that the frequency is triggered every 30 seconds.
  - `cron`: a standard cron expression.
  - `second_delay`: The number of seconds to be delayed (1s to 60s).
* `time_type` - (Required, Int) Time configuration type. Currently, the following time types are supported:
  - `1`:cron
  - `3`:fix_rate
  - `4`:second_delay
  - `100`:api

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<namespace>:<group_id>:<job_id>`.
* `job_id` - JobId

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Job.
* `delete` - (Defaults to 5 mins) Used when delete the Job.
* `update` - (Defaults to 5 mins) Used when update the Job.

## Import

Schedulerx Job can be imported using the id, e.g.

```shell
$ terraform import alicloud_schedulerx_job.example <namespace>:<group_id>:<job_id>
```