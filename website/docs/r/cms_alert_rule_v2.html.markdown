---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_alert_rule_v2"
description: |-
  Provides a Alicloud Cms Alert Rule V2 resource.
---

# alicloud_cms_alert_rule_v2

Provides a Cms Alert Rule V2 resource.

CloudMonitor 2.0 alert rules (Unified Action V2 OpenAPI, based on ManageAlertRules and QueryAlertRules).

For information about Cms Alert Rule V2 and how to use it, see [What is Alert Rule V2](https://next.api.alibabacloud.com/document/Cms/2024-03-30/ManageAlertRules).

-> **NOTE:** Available since v1.285.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cms_alert_rule_v2&exampleId=6e622797-a599-2d94-9e32-e4edab9be5ab25c643a2&activeTab=example&spm=docs.r.cms_alert_rule_v2.0.6e622797a5&intl_lang=EN_US" target="_blank">
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


resource "alicloud_cms_alert_rule_v2" "default" {
  content_template = "umodel example alert on $${metric}"
  schedule_config {
    type          = "FIXED"
    interval_secs = "60"
  }
  datasource_config {
    type = "UMODEL"
  }
  action_integration_config {
    enabled = false
  }
  arms_integration_config {
    enabled = false
  }
  query_config {
    entity_type   = "instance"
    type          = "UMODEL_METRICSET_QUERY"
    entity_domain = "ecs"
    metric        = "CPUUtilization"
    label_filters {
      operator = "="
      value    = "web-server"
      name     = "app"
    }
    label_filters {
      operator = "="
      value    = "production"
      name     = "env"
    }
    metric_set = "acs_ecs_dashboard"
  }
  display_name = "regression-umodel-10"
  enabled      = true
  notify_config {
    type = "DIRECT_NOTIFY"
    channels {
      type        = "GROUP"
      identifiers = ["regression-example"]
    }
  }
  condition_config {
    operator      = "GT"
    type          = "UMODEL_METRICSET_CONDITION"
    severity      = "CRITICAL"
    duration_secs = "60"
    threshold     = 90
  }
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cms_alert_rule_v2&spm=docs.r.cms_alert_rule_v2.example&intl_lang=EN_US)


## Argument Reference

The following arguments are supported:
* `action_integration_config` - (Optional, Set) Action integration configuration. See [`action_integration_config`](#action_integration_config) below.
* `alert_rule_v2_id` - (Computed) The unique identifier of the alert rule, mapped to a UUID (system-generated).
* `annotations` - (Optional, Map) Annotations.
* `arms_integration_config` - (Optional, Set) ARMS integration configuration. See [`arms_integration_config`](#arms_integration_config) below.
* `condition_config` - (Optional, Set) Unified alert condition configuration. See [`condition_config`](#condition_config) below.
* `content_template` - (Optional) The alert content template.
* `datasource_config` - (Optional, ForceNew, Set) Unified data source configuration. See [`datasource_config`](#datasource_config) below.
* `display_name` - (Optional) The display name of the alert rule.
* `enabled` - (Optional) Specifies whether the alert rule is enabled.
* `labels` - (Optional, Map) Labels.
* `notify_config` - (Optional, Set) Unified notification configuration. See [`notify_config`](#notify_config) below.
* `query_config` - (Optional, Set) Unified query configuration. See [`query_config`](#query_config) below.
* `schedule_config` - (Optional, Set) Unified scheduling configuration. See [`schedule_config`](#schedule_config) below.
* `workspace` - (Optional, ForceNew) Workspace.

### `action_integration_config`

The action_integration_config supports the following:
* `actions` - (Optional, ForceNew, List) List of actions
* `enabled` - (Optional, ForceNew) Indicates whether action integration is enabled

### `arms_integration_config`

The arms_integration_config supports the following:
* `enabled` - (Optional, ForceNew) Specifies whether to enable ARMS integration.

### `condition_config`

The condition_config supports the following:
* `aggregate` - (Optional, ForceNew) Aggregation function (used when APM_SIMPLE is specified). Valid values include AVG, MAX, MIN, and SUM.
* `compare_list` - (Optional, ForceNew, List) Multiple comparison list (used when type=APM_COMPOSITE) See [`compare_list`](#condition_config-compare_list) below.
* `composite_escalation` - (Computed) The multi-metric composite trigger configuration. This parameter is required when type is set to CLOUD_MONITORING and escalationType is set to composite. See [`composite_escalation`](#condition_config-composite_escalation) below.
* `duration_secs` - (Optional, ForceNew, Int) Duration (seconds). Used when type=PROMETHEUS_SIMPLE or UMODEL_METRICSET.
* `escalation_type` - (Computed) The escalation policy type (type=CLOUD_MONITORING). Valid values: SIMPLE, COMPOSITE, EXPRESS, and PROMETHEUS.
* `express_escalation` - (Computed) Expression trigger configuration. This parameter is required when type=CLOUD_MONITORING and escalationType=express. See [`express_escalation`](#condition_config-express_escalation) below.
* `legacy_raw` - (Computed) The original V1 condition JSON string returned when type=UNKNOWN and parsing fails. The frontend displays this field as read-only.
* `legacy_type` - (Computed) Returned when type=UNKNOWN, indicating that this rule cannot be edited by using the new API.
* `no_data_policy` - (Computed) No-data processing policy (type=CLOUD_MONITORING)
* `operator` - (Optional, ForceNew) The comparison operator. Valid values include GT, GE, LT, LE, EQ, and NE.
* `prometheus` - (Computed) The PromQL trigger configuration. This field is not empty when type=CLOUD_MONITORING and escalationType=prometheus. See [`prometheus`](#condition_config-prometheus) below.
* `relation` - (Optional, ForceNew) The logical relationship between conditions (type=APM_COMPOSITE). Valid values: AND and OR.
* `severity` - (Optional, ForceNew) The severity level. Valid values: CRITICAL, ERROR, WARNING, and INFO.
* `simple_escalation` - (Computed) Single-metric trigger configuration (Required when type=CLOUD_MONITORING and escalationType=simple) See [`simple_escalation`](#condition_config-simple_escalation) below.
* `threshold` - (Optional, ForceNew, Float) Threshold (used when UMODEL_METRICSET is specified)
* `threshold_list` - (Optional, ForceNew, List) Multi-Threshold List (Used for APM_SIMPLE or UMODEL_METRICSET_MULTI_SEVERITY) See [`threshold_list`](#condition_config-threshold_list) below.
* `type` - (Required, ForceNew) The detection condition type. Valid values: PROMETHEUS_SIMPLE, UMODEL_METRICSET, APM_SIMPLE, APM_COMPOSITE, CLOUD_MONITORING, and UNKNOWN.
* `yoy_time_unit` - (Optional, ForceNew) Year-over-year time unit. This parameter takes effect only when type is set to APM_SIMPLE and operator is set to YOY_UP or YOY_DOWN. Valid values: minute, hour, day, week, and month
* `yoy_time_value` - (Optional, ForceNew, Int) Year-over-Year Time Value (Effective only when type=APM_SIMPLE and operator=YOY_UP or YOY_DOWN)

### `condition_config-compare_list`

The condition_config-compare_list supports the following:
* `aggregate` - (Optional, ForceNew) Aggregation Function
* `operator` - (Optional, ForceNew) Comparison Operator
* `threshold` - (Optional, ForceNew, Float) Threshold
* `yoy_time_unit` - (Optional, ForceNew) Year-over-year time unit. This parameter takes effect only when operator=YOY_UP or YOY_DOWN.
* `yoy_time_value` - (Optional, ForceNew, Int) Year-over-year time value. This parameter takes effect only when operator=YOY_UP or YOY_DOWN.

### `condition_config-composite_escalation`

The condition_config-composite_escalation supports the following:
* `escalations` - (Computed) List of trigger conditions See [`escalations`](#condition_config-composite_escalation-escalations) below.
* `relation` - (Computed) The logical relationship between multiple metrics
* `severity` - (Computed) Severity Level
* `times` - (Computed) Consecutive Trigger Count

### `condition_config-express_escalation`

The condition_config-express_escalation supports the following:
* `raw_expression` - (Computed) Raw Expression
* `severity` - (Computed) Severity Level
* `times` - (Computed) Consecutive Trigger Count

### `condition_config-prometheus`

The condition_config-prometheus supports the following:
* `prom_ql` - (Computed) Prometheus Query Expression
* `severity` - (Computed) Severity Level
* `times` - (Computed) Consecutive Trigger Count

### `condition_config-simple_escalation`

The condition_config-simple_escalation supports the following:
* `escalations` - (Computed) Trigger Condition List See [`escalations`](#condition_config-simple_escalation-escalations) below.
* `metric_name` - (Computed) Metric Name
* `period` - (Computed) Collection Period (Seconds)

### `condition_config-threshold_list`

The condition_config-threshold_list supports the following:
* `severity` - (Optional, ForceNew) Severity Level
* `threshold` - (Optional, ForceNew, Float) Threshold

### `condition_config-simple_escalation-escalations`

The condition_config-simple_escalation-escalations supports the following:
* `comparison_operator` - (Computed) Comparison Operator
* `pre_condition` - (Computed) Precondition
* `severity` - (Computed) Severity level. Valid values: CRITICAL, ERROR, WARNING, and INFO
* `statistics` - (Computed) Statistical Method
* `threshold` - (Computed) Threshold
* `times` - (Computed) Consecutive Trigger Count

### `condition_config-composite_escalation-escalations`

The condition_config-composite_escalation-escalations supports the following:
* `comparison_operator` - (Computed) Comparison operator
* `metric_name` - (Computed) Metric name
* `period` - (Computed) Collection period (s)
* `pre_condition` - (Computed) Precondition
* `statistics` - (Computed) Statistical Method
* `threshold` - (Computed) Threshold

### `datasource_config`

The datasource_config supports the following:
* `instance_id` - (Optional, ForceNew) The Prometheus instance ID. This parameter is used when type=PROMETHEUS.
* `legacy_raw` - (Computed) The original V1 datasource JSON string returned when type=UNKNOWN and parsing fails. The frontend displays this string in read-only mode.
* `legacy_type` - (Computed) Returned when type=UNKNOWN, indicating that this rule cannot be edited by using the new API.
* `product_category` - (Computed) The cloud service category. This parameter is used when type=CLOUD_MONITORING. If this parameter is not specified, unknown is returned.
* `region_id` - (Optional, ForceNew) The region ID. This parameter is available for all types. By default, the value is the same as the region where the rule resides.
* `type` - (Required, ForceNew) The data source type. Valid values: PROMETHEUS, UMODEL, APM, CLOUD_MONITORING, and UNKNOWN.

### `notify_config`

The notify_config supports the following:
* `active_days` - (Optional, ForceNew, List) The days of the week on which notifications are sent, 1-7 (type=DIRECT_NOTIFY). Default: [1,2,3,4,5,6,7]
* `active_end_time` - (Optional, ForceNew) The daily end time of the effective notification period (HH:mm, type=DIRECT_NOTIFY). Default: 23:59
* `active_start_time` - (Optional, ForceNew) The daily start time of the effective notification period (HH:mm, type=DIRECT_NOTIFY). Default: 00:00
* `channels` - (Optional, ForceNew, List) The list of notification channels (type=DIRECT_NOTIFY) See [`channels`](#notify_config-channels) below.
* `notify_strategies` - (Optional, ForceNew, List) List of notification policy IDs (type=NOTIFY_POLICY, up to 1 for the current business)
* `silence_time_secs` - (Optional, ForceNew, Int) The channel silence period in seconds (type=DIRECT_NOTIFY). Default: 86400
* `type` - (Required, ForceNew) Notification configuration type. Valid values: DIRECT_NOTIFY and NOTIFY_POLICY
* `utc_offset` - (Optional, ForceNew) UTC time zone offset (type=DIRECT_NOTIFY). Default: +08:00

### `notify_config-channels`

The notify_config-channels supports the following:
* `identifiers` - (Optional, ForceNew, List) List of channel identifiers
* `type` - (Optional, ForceNew) Notification Channel Type

### `query_config`

The query_config supports the following:
* `dimensions` - (Computed) The dimension list (type=CLOUD_MONITORING_QUERY). Each dimension is a key-value string mapping. See [`dimensions`](#query_config-dimensions) below.
* `enable_data_complete_check` - (Optional, ForceNew) Whether to Enable Data Completeness Check (type=PROMETHEUS_SINGLE_QUERY)
* `entity_domain` - (Optional, ForceNew) Domain to which the entity belongs (type=UMODEL_METRICSET_QUERY)
* `entity_fields` - (Optional, ForceNew, List) List of Entity Fields to Return (type=UMODEL_METRICSET_QUERY) See [`entity_fields`](#query_config-entity_fields) below.
* `entity_filters` - (Optional, ForceNew, List) The list of entity filter conditions (type=UMODEL_METRICSET_QUERY). See [`entity_filters`](#query_config-entity_filters) below.
* `entity_type` - (Optional, ForceNew) Entity type (type=UMODEL_METRICSET_QUERY)
* `expr` - (Optional, ForceNew) Prometheus query statement (type=PROMETHEUS_SINGLE_QUERY, recommended field)
* `filter_list` - (Optional, ForceNew, List) The APM filter condition list (type=APM_MULTI_QUERY). See [`filter_list`](#query_config-filter_list) below.
* `group_id` - (Computed) Resource group ID (used when type=CLOUD_MONITORING_QUERY and relationType=GROUP)
* `label_filters` - (Optional, ForceNew, List) List of label filter conditions (type=UMODEL_METRICSET_QUERY) See [`label_filters`](#query_config-label_filters) below.
* `legacy_raw` - (Computed) The raw V1 query JSON string returned when type=UNKNOWN_QUERY and parsing fails. The frontend displays this field as read-only only.
* `legacy_type` - (Computed) Returned when type=UNKNOWN_QUERY, indicating that this rule cannot be edited through the new API.
* `measure_list` - (Optional, ForceNew, List) APM measure configuration list (type=APM_MULTI_QUERY) See [`measure_list`](#query_config-measure_list) below.
* `metric` - (Optional, ForceNew) Metric name (type=UMODEL_METRICSET_QUERY)
* `metric_set` - (Optional, ForceNew) Metric set name (type=UMODEL_METRICSET_QUERY)
* `namespace` - (Computed) CloudMonitor namespace (cloud service name, type=CLOUD_MONITORING_QUERY)
* `prom_ql` - (Computed) [Deprecated] Legacy PromQL field. Use Expr instead. The backend automatically normalizes this field to Expr.
* `relation_type` - (Computed) Resource association type (type=CLOUD_MONITORING_QUERY). Valid values: INSTANCE, GROUP, and USER.
* `service_id_list` - (Optional, ForceNew, List) Application Service ID List (type=APM_MULTI_QUERY)
* `type` - (Required, ForceNew) Query type. Valid values: PROMETHEUS_SINGLE_QUERY, UMODEL_METRICSET_QUERY, APM_MULTI_QUERY, CLOUD_MONITORING_QUERY, and UNKNOWN_QUERY.

### `query_config-dimensions`

The query_config-dimensions supports the following:

### `query_config-entity_fields`

The query_config-entity_fields supports the following:
* `field` - (Optional, ForceNew) Entity Field Name
* `value` - (Optional, ForceNew) Entity Field Value

### `query_config-entity_filters`

The query_config-entity_filters supports the following:
* `field` - (Optional, ForceNew) The entity filter field name.
* `operator` - (Optional, ForceNew) The entity filter operator. Valid values: = and !=.
* `value` - (Optional, ForceNew) The entity filter value.

### `query_config-filter_list`

The query_config-filter_list supports the following:
* `key` - (Optional, ForceNew) APM filter dimension key
* `type` - (Optional, ForceNew) The APM filter type. Valid values: ALL, EQ, NE, and DISABLED.
* `value` - (Optional, ForceNew) APM filter value (can be empty when type=ALL/DISABLED)

### `query_config-label_filters`

The query_config-label_filters supports the following:
* `name` - (Optional, ForceNew) Label name
* `operator` - (Optional, ForceNew) Label filter operator. Valid values: = / !=
* `value` - (Optional, ForceNew) Label value

### `query_config-measure_list`

The query_config-measure_list supports the following:
* `group_by` - (Optional, ForceNew, List) Grouping dimension list
* `measure_code` - (Optional, ForceNew) APM metric code
* `window_secs` - (Optional, ForceNew, Int) Query Time Window (Seconds)

### `schedule_config`

The schedule_config supports the following:
* `interval_secs` - (Optional, ForceNew, Int) The scheduling interval in seconds. This parameter is used when the type is set to FIXED.
* `type` - (Required, ForceNew) The scheduling type. Valid values: FIXED and CRON.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. 
* `created_at` - Creation time (read-only), in ISO 8601 format.
* `datasource_type` - Data source type (read-only, derived).
* `notify_strategy_id` - Notification policy ID (read-only, derived).
* `observe_resource_global_scope` - Indicates whether the rule applies to all resources of this resource type (read-only, derived).
* `observe_resource_type` - Observable resource type (read-only, derived).
* `partition_key` - The partition key.
* `severity_levels` - The severity levels covered by this rule, separated by commas (read-only derived).
* `status` - Alert status (read-only).
* `updated_at` - The update time (read-only), in ISO 8601 format.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Alert Rule V2.
* `delete` - (Defaults to 5 mins) Used when delete the Alert Rule V2.
* `update` - (Defaults to 5 mins) Used when update the Alert Rule V2.

## Import

Cms Alert Rule V2 can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_alert_rule_v2.example <alert_rule_v2_id>
```