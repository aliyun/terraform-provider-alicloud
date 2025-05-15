---
subcategory: "Data Works"
layout: "alicloud"
page_title: "Alicloud: alicloud_data_works_di_alarm_rule"
description: |-
  Provides a Alicloud Data Works Di Alarm Rule resource.
---

# alicloud_data_works_di_alarm_rule

Provides a Data Works Di Alarm Rule resource.

Data Integration alarm rules.

For information about Data Works Di Alarm Rule and how to use it, see [What is Di Alarm Rule](https://www.alibabacloud.com/help/en/dataworks/developer-reference/api-dataworks-public-2024-05-18-createdialarmrule).

-> **NOTE:** Available since v1.241.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_data_works_di_alarm_rule&exampleId=8c7d98e8-66d1-2257-1072-87c4cf9f6d462719ee7f&activeTab=example&spm=docs.r.data_works_di_alarm_rule.0.8c7d98e866&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform_example"
}

provider "alicloud" {
  region = "cn-chengdu"
}

resource "alicloud_data_works_project" "defaulteNv8bu" {
  project_name     = var.name
  display_name     = var.name
  description      = var.name
  pai_task_enabled = true
}

resource "alicloud_data_works_di_job" "defaultUW8inp" {
  description             = "xxxx"
  project_id              = alicloud_data_works_project.defaulteNv8bu.id
  job_name                = "xxx"
  migration_type          = "api_xxx"
  source_data_source_type = "xxx"
  resource_settings {
    offline_resource_settings {
      requested_cu              = 2
      resource_group_identifier = "xx"
    }
    realtime_resource_settings {
      requested_cu              = 2
      resource_group_identifier = "xx"
    }
    schedule_resource_settings {
      requested_cu              = 2
      resource_group_identifier = "xx"
    }
  }
  job_settings {
    channel_settings = "xxxx"
    column_data_type_settings {
      destination_data_type = "xxxx"
      source_data_type      = "xxxx"
    }
    cycle_schedule_settings {
      cycle_migration_type = "xxxx"
      schedule_parameters  = "xxxx"
    }
  }
  source_data_source_settings {
    data_source_name = "xxxx"
    data_source_properties {
      encoding = "xxxx"
      timezone = "xxxx"
    }
  }
  destination_data_source_type = "xxxx"
  table_mappings {
    source_object_selection_rules {
      action          = "Include"
      expression      = "xxxx"
      expression_type = "Exact"
      object_type     = "xxxx"
    }
    source_object_selection_rules {
      action          = "Include"
      expression      = "xxxx"
      expression_type = "Exact"
      object_type     = "xxxx"
    }
    transformation_rules {
      rule_name        = "xxxx"
      rule_action_type = "xxxx"
      rule_target_type = "xxxx"
    }
  }
  transformation_rules {
    rule_action_type = "xxxx"
    rule_expression  = "xxxx"
    rule_name        = "xxxx"
    rule_target_type = "xxxx"
  }
  destination_data_source_settings {
    data_source_name = "xxx"
  }
}


resource "alicloud_data_works_di_alarm_rule" "default" {
  description = "Description"
  trigger_conditions {
    ddl_report_tags = ["ALTERADDCOLUMN"]
    threshold       = "20"
    duration        = "10"
    severity        = "Warning"
  }
  metric_type = "DdlReport"
  notification_settings {
    notification_channels {
      severity = "Warning"
      channels = ["Ding"]
    }
    notification_receivers {
      receiver_type   = "DingToken"
      receiver_values = ["1107550004253538"]
    }
    inhibition_interval = "10"
  }
  di_job_id          = alicloud_data_works_di_job.defaultUW8inp.di_job_id
  di_alarm_rule_name = var.name
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Optional) 
* `di_alarm_rule_name` - (Required) Data Integration alarm rule name
* `di_job_id` - (Required, ForceNew, Int) Task ID: the ID of the task associated with the alert rule.
* `enabled` - (Optional) 
* `metric_type` - (Required) Alarm indicator type. Optional enumerated values:
  - Heartbeat (task status alarm)
  - FailoverCount(failover times alarm)
  - Delay (task Delay alarm)
* `notification_settings` - (Required, List) Alarm notification settings See [`notification_settings`](#notification_settings) below.
* `trigger_conditions` - (Required, List) Alarm trigger condition list, supporting multiple conditions See [`trigger_conditions`](#trigger_conditions) below.

### `notification_settings`

The notification_settings supports the following:
* `inhibition_interval` - (Optional, Int) Alarm suppression interval, in minutes
* `notification_channels` - (Optional, List) Alarm notification Channel See [`notification_channels`](#notification_settings-notification_channels) below.
* `notification_receivers` - (Optional, List) List of alert notification recipients See [`notification_receivers`](#notification_settings-notification_receivers) below.

### `notification_settings-notification_channels`

The notification_settings-notification_channels supports the following:
* `channels` - (Optional, List) Channel, optional enumeration value:

  Mail (Mail)

  Phone (Phone)

  Sms (Sms)

  Ding (DingTalk)
* `severity` - (Optional) Severity, optional enumeration value:

  Warning

  Critical

### `notification_settings-notification_receivers`

The notification_settings-notification_receivers supports the following:
* `receiver_type` - (Optional) The type of the receiver. Valid values: AliyunUid/DingToken/FeishuToken/WebHookUrl.
* `receiver_values` - (Optional, List) Receiver Value List

### `trigger_conditions`

The trigger_conditions supports the following:
* `ddl_report_tags` - (Optional, List) It takes effect only when the DDL notification is issued. The list of effective DDLs is required.
* `duration` - (Optional, Int) Alarm calculation time interval, unit minute
* `severity` - (Optional) Severity, optional enumeration value:

  Warning

  Critical
* `threshold` - (Optional, Int) Alarm threshold.

  Task status alarm: no need to fill in the threshold.

  failover alarm: The threshold is the number of failover alarms.

  Task Delay Alarm: The threshold is the delay duration, in seconds.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<di_job_id>:<di_alarm_rule_id>`.
* `di_alarm_rule_id` - Resource attribute field representing resource level ID

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Di Alarm Rule.
* `delete` - (Defaults to 5 mins) Used when delete the Di Alarm Rule.
* `update` - (Defaults to 5 mins) Used when update the Di Alarm Rule.

## Import

Data Works Di Alarm Rule can be imported using the id, e.g.

```shell
$ terraform import alicloud_data_works_di_alarm_rule.example <di_job_id>:<di_alarm_rule_id>
```