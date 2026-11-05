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

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

variable "name" {
  default = "terraform-example-cms-alert-rule-v2"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_log_project" "default" {
  project_name = var.name
}

resource "alicloud_cms_workspace" "default" {
  workspace_name = var.name
  sls_project    = alicloud_log_project.default.project_name
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
    window_secs   = 60
    aggregate     = "AVG"
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

  workspace = alicloud_cms_workspace.default.id
}
```

## Argument Reference

The following arguments are supported:
* `action_integration_config` - (Optional, Computed, List) Action integration configuration. See [`action_integration_config`](#action_integration_config) below.
* `alert_rule_v2_id` - (Computed) The unique identifier of the alert rule, mapped to a UUID (system-generated).
* `annotations` - (Optional, Computed, Map) Annotations.
* `arms_integration_config` - (Optional, Computed, List) ARMS integration configuration. See [`arms_integration_config`](#arms_integration_config) below.
* `condition_config` - (Optional, Computed, List) Unified alert condition configuration. See [`condition_config`](#condition_config) below.
* `content_template` - (Optional, Computed, String) The alert content template.
* `datasource_config` - (Optional, Computed, List) Unified data source configuration. See [`datasource_config`](#datasource_config) below.
* `display_name` - (Optional, Computed, String) The display name of the alert rule.
* `enabled` - (Optional, Computed, Bool) Specifies whether the alert rule is enabled.
* `labels` - (Optional, Computed, Map) Labels.
* `notify_config` - (Optional, Computed, List) Unified notification configuration. See [`notify_config`](#notify_config) below.
* `observe_resource_config` - (Optional, List, Available since v2.2.0) Observe resource configuration See [`observe_resource_config`](#observe_resource_config) below.
* `query_config` - (Optional, Computed, List) Unified query configuration. See [`query_config`](#query_config) below.
* `rca_config` - (Optional, List, Available since v2.2.0) RCA (Root Cause Analysis) configuration See [`rca_config`](#rca_config) below.
* `schedule_config` - (Optional, Computed, List) Unified scheduling configuration. See [`schedule_config`](#schedule_config) below.
* `workspace` - (Optional, Computed, ForceNew, String) Workspace.

### `action_integration_config`

The action_integration_config supports the following:
* `actions` - (Optional, ForceNew, Computed, List) List of actions
* `enabled` - (Optional, ForceNew, Computed) Indicates whether action integration is enabled

### `arms_integration_config`

The arms_integration_config supports the following:
* `enabled` - (Optional, ForceNew, Computed) Specifies whether to enable ARMS integration.

### `condition_config`

The condition_config supports the following:
* `aggregate` - (Optional, ForceNew, Computed, String) Aggregation function (used when type=APM_SIMPLE_CONDITION). Valid values: `AVG`, `MAX`, `MIN`, `SUM`.
* `alert_count` - (Optional, ForceNew, Computed, Int, Available since v2.2.0) Consecutive trigger count threshold (type=SLS_MULTI_CONDITION). The alert fires after the condition is met N times. Defaults to `1`.
* `compare_list` - (Optional, ForceNew, Computed, List) Multiple comparison list (used when type=APM_COMPOSITE_CONDITION) See [`compare_list`](#condition_config-compare_list) below.
* `composite_escalation` - (Optional, Computed, List) The multi-metric composite trigger configuration. This parameter is required when type=CLOUD_MONITORING_CONDITION and escalation_type=COMPOSITE. See [`composite_escalation`](#condition_config-composite_escalation) below.
* `count_operator` - (Optional, ForceNew, String, Available since v2.2.0) Count comparison operator (used when type=UMODEL_LOGSET_CONDITION). Valid values: `GTE`, `GT`, `EQ`, `LTE`, `LT`.
* `count_threshold` - (Optional, ForceNew, Int, Available since v2.2.0) Count threshold (used when type=UMODEL_LOGSET_CONDITION)
* `duration_secs` - (Optional, ForceNew, Computed, Int) Duration (seconds). Used when type=PROMETHEUS_SIMPLE_CONDITION, UMODEL_METRICSET_CONDITION, UMODEL_LOGSET_CONDITION or UMODEL_METRICSET_MULTI_CONDITION.
* `escalation_type` - (Optional, Computed, String) The escalation policy type (type=CLOUD_MONITORING_CONDITION). Valid values: `SIMPLE`, `COMPOSITE`, `EXPRESS`, `PROMETHEUS`.
* `express_escalation` - (Computed, List) Expression trigger configuration. This field is only populated when type=CLOUD_MONITORING_CONDITION and escalation_type=EXPRESS; it cannot be configured. See [`express_escalation`](#condition_config-express_escalation) below.
* `legacy_raw` - (Computed, String) The original V1 condition JSON string returned when type=UNKNOWN_CONDITION and parsing fails. The frontend displays this field as read-only.
* `legacy_type` - (Computed, String) Returned when type=UNKNOWN_CONDITION, indicating that this rule cannot be edited by using the new API.
* `match_field` - (Optional, ForceNew, String, Available since v2.2.0) Log field name (used when type=UMODEL_LOGSET_CONDITION and match_operator=CONTAINS/EQUALS/REGEX)
* `match_operator` - (Optional, ForceNew, String, Available since v2.2.0) Log match operator (used when type=UMODEL_LOGSET_CONDITION). Valid values: `PRESENT`, `NOT_PRESENT`, `CONTAINS`, `EQUALS`, `REGEX`.
* `match_value` - (Optional, ForceNew, String, Available since v2.2.0) Log match value (used when type=UMODEL_LOGSET_CONDITION and match_operator=CONTAINS/EQUALS/REGEX)
* `no_data_alert_level` - (Optional, ForceNew, Computed, String, Available since v2.2.0) No data alert level (type=SLS_MULTI_CONDITION). Valid values: `INFO`, `WARN`, `WARNING`, `ERROR`, `CRITICAL`.
* `no_data_alert_severity` - (Optional, ForceNew, String, Available since v2.2.0) No data alert severity (type=PROMETHEUS_SIMPLE_CONDITION / PROMETHEUS_MULTI_CONDITION / UMODEL_METRICSET_CONDITION / UMODEL_METRICSET_MULTI_CONDITION, effective when no_data_policy=NO_DATA_TO_ALERT and forbidden for other policies). Valid values: `INFO`, `WARN`, `WARNING`, `ERROR`, `CRITICAL`.
* `no_data_duration_secs` - (Optional, ForceNew, Int, Available since v2.2.0) No data tolerance window in seconds (type=UMODEL_METRICSET_CONDITION / UMODEL_METRICSET_MULTI_CONDITION, effective when no_data_policy is set). The policy applies only after data has been absent for longer than this window. Defaults to max(collection interval x 2, 300) when omitted
* `no_data_policy` - (Optional, Computed, String) No-data processing policy (type=CLOUD_MONITORING_CONDITION / PROMETHEUS_SIMPLE_CONDITION / APM_SIMPLE_CONDITION / APM_COMPOSITE_CONDITION / UMODEL_METRICSET_CONDITION / UMODEL_METRICSET_MULTI_CONDITION). Valid values: `NO_DATA_TO_OK`, `NO_DATA_TO_ALERT`, `KEEP_LAST_STATE`, `APPEND_VALUE`.
* `operator` - (Optional, ForceNew, Computed, String) The comparison operator. Valid values: `GT`, `GE`, `LT`, `LE`, `EQ`, `NE`, `IN_RANGE`, `OUT_OF_RANGE`, `PRESENT`, `NOT_PRESENT`, `YOY_UP`, `YOY_DOWN`, `ABOVE_UPPER`, `BELOW_LOWER`, `OUT_OF_BAND`.
* `prometheus` - (Computed, List) The PromQL trigger configuration. This field is only populated when type=CLOUD_MONITORING_CONDITION and escalation_type=PROMETHEUS; it cannot be configured. See [`prometheus`](#condition_config-prometheus) below.
* `relation` - (Optional, Computed, String) The logical relationship between conditions (type=APM_COMPOSITE_CONDITION). Valid values: `AND`, `OR`.
* `severity` - (Optional, ForceNew, Computed, String) The severity level. Valid values: `CRITICAL`, `ERROR`, `WARN`, `WARNING`, `INFO`.
* `simple_escalation` - (Optional, Computed, List) Single-metric trigger configuration (required when type=CLOUD_MONITORING_CONDITION and escalation_type=SIMPLE) See [`simple_escalation`](#condition_config-simple_escalation) below.
* `threshold` - (Optional, ForceNew, Computed, Float) Threshold (used when UMODEL_METRICSET_CONDITION is specified)
* `threshold_list` - (Optional, ForceNew, Computed, List) Multi-Threshold List (Used for APM_SIMPLE_CONDITION or UMODEL_METRICSET_MULTI_SEVERITY) See [`threshold_list`](#condition_config-threshold_list) below.
* `triggers` - (Optional, ForceNew, Computed, List, Available since v2.2.0) Trigger list (at least one when type=SLS_MULTI_CONDITION; the triggers are sorted by severity priority and the first match triggers). See [`triggers`](#condition_config-triggers) below.
* `type` - (Required, ForceNew, String) The detection condition type. Valid values: `PROMETHEUS_SIMPLE_CONDITION`, `APM_SIMPLE_CONDITION`, `APM_COMPOSITE_CONDITION`, `UMODEL_METRICSET_CONDITION`, `UMODEL_LOGSET_CONDITION`, `SLS_MULTI_CONDITION`, `CLOUD_MONITORING_CONDITION`.
* `yoy_time_unit` - (Optional, ForceNew, Computed, String) Year-over-year time unit. This parameter takes effect only when type=APM_SIMPLE_CONDITION and operator is set to YOY_UP or YOY_DOWN. Valid values: `minute`, `hour`, `day`, `week`, `month`.
* `yoy_time_value` - (Optional, ForceNew, Computed, Int) Year-over-Year Time Value (Effective only when type=APM_SIMPLE_CONDITION and operator=YOY_UP or YOY_DOWN)

### `condition_config-compare_list`

The condition_config-compare_list supports the following:
* `aggregate` - (Optional, ForceNew, Computed, String) Aggregation Function
* `operator` - (Optional, ForceNew, Computed, String) Comparison Operator
* `threshold` - (Optional, ForceNew, Computed, Float) Threshold
* `yoy_time_unit` - (Optional, ForceNew, Computed, String) Year-over-year time unit. This parameter takes effect only when operator=YOY_UP or YOY_DOWN.
* `yoy_time_value` - (Optional, ForceNew, Computed, Int) Year-over-year time value. This parameter takes effect only when operator=YOY_UP or YOY_DOWN.

### `condition_config-composite_escalation`

The condition_config-composite_escalation supports the following:
* `escalations` - (Optional, Computed, List) List of trigger conditions See [`escalations`](#condition_config-composite_escalation-escalations) below.
* `relation` - (Optional, Computed, String) The logical relationship between multiple metrics. Valid values: `and`, `or`.
* `severity` - (Optional, Computed, String) Severity Level
* `times` - (Optional, Computed, Int) Consecutive Trigger Count

### `condition_config-express_escalation`

The condition_config-express_escalation supports the following:
* `raw_expression` - (Computed, String) Raw Expression
* `severity` - (Computed, String) Severity Level
* `times` - (Computed, Int) Consecutive Trigger Count

### `condition_config-prometheus`

The condition_config-prometheus supports the following:
* `prom_ql` - (Computed, String) Prometheus Query Expression
* `severity` - (Computed, String) Severity Level
* `times` - (Computed, Int) Consecutive Trigger Count

### `condition_config-simple_escalation`

The condition_config-simple_escalation supports the following:
* `escalations` - (Optional, Computed, List) Trigger Condition List See [`escalations`](#condition_config-simple_escalation-escalations) below.
* `metric_name` - (Optional, Computed, String) Metric Name
* `period` - (Optional, Computed, Int) Collection Period (Seconds)

### `condition_config-threshold_list`

The condition_config-threshold_list supports the following:
* `severity` - (Optional, ForceNew, Computed, String) Severity Level
* `threshold` - (Optional, ForceNew, Computed, Float) Threshold

### `condition_config-triggers`

The condition_config-triggers supports the following:
* `condition` - (Optional, ForceNew, Computed, String, Available since v2.2.0) Match expression (used when type=SLS_MULTI_CONDITION). The V1 condition expression is preserved verbatim and not parsed.
* `count_condition` - (Optional, ForceNew, Computed, String, Available since v2.2.0) Count match expression (used when type=SLS_MULTI_CONDITION). The V1 countCondition expression is preserved verbatim and not parsed.
* `operator` - (Optional, ForceNew, Computed, String, Available since v2.2.0) Operation of the trigger (used when type=SLS_MULTI_CONDITION). Valid values: `HAS_DATA`, `HAS_DATA_COUNT`, `HAS_DATA_MATCH`, `HAS_DATA_MATCH_COUNT`. This parameter is required for the SLS multi-condition rule.
* `severity` - (Optional, ForceNew, Computed, String, Available since v2.2.0) Severity level (used when type=SLS_MULTI_CONDITION). Valid values: `CRITICAL`, `ERROR`, `WARN`, `WARNING`, `INFO`.

### `condition_config-simple_escalation-escalations`

The condition_config-simple_escalation-escalations supports the following:
* `comparison_operator` - (Optional, Computed, String) Comparison Operator
* `pre_condition` - (Optional, Computed, String) Precondition
* `severity` - (Optional, Computed, String) Severity level. Valid values: `CRITICAL`, `ERROR`, `WARNING`, `INFO`.
* `statistics` - (Optional, Computed, String) Statistical Method
* `threshold` - (Optional, Computed, String) Threshold
* `times` - (Optional, Computed, Int) Consecutive Trigger Count

### `condition_config-composite_escalation-escalations`

The condition_config-composite_escalation-escalations supports the following:
* `comparison_operator` - (Optional, Computed, String) Comparison operator
* `metric_name` - (Optional, Computed, String) Metric name
* `period` - (Optional, Computed, Int) Collection period (s)
* `pre_condition` - (Optional, Computed, String) Precondition
* `statistics` - (Optional, Computed, String) Statistical Method
* `threshold` - (Optional, Computed, String) Threshold

### `datasource_config`

The datasource_config supports the following:
* `instance_id` - (Optional, ForceNew, Computed, String) The Prometheus instance ID. This parameter is used when type=PROMETHEUS.
* `legacy_raw` - (Computed, String) The original V1 datasource JSON string returned when type=UNKNOWN and parsing fails. The frontend displays this string in read-only mode.
* `legacy_type` - (Computed, String) Returned when type=UNKNOWN, indicating that this rule cannot be edited by using the new API.
* `product_category` - (Optional, Computed, String) The cloud service category. This parameter is used when type=CLOUD_MONITORING. If this parameter is not specified, unknown is returned.
* `project` - (Optional, ForceNew, String, Available since v2.2.0) SLS project name (used when type=SLS; all stores share the same project)
* `region_id` - (Optional, ForceNew, Computed, String) The region ID. This parameter is available for all types. By default, the value is the same as the region where the rule resides.
* `stores` - (Optional, ForceNew, List, Available since v2.2.0) SLS store list (used when type=SLS, at least 1). Each store contains `store` and `store_type` See [`stores`](#datasource_config-stores) below.
* `type` - (Required, ForceNew, String) The data source type. Valid values: `PROMETHEUS`, `SLS`, `UMODEL`, `APM`, `CLOUD_MONITORING`, `UNKNOWN`.

### `datasource_config-stores`

The datasource_config-stores supports the following:
* `store` - (Optional, ForceNew, String, Available since v2.2.0) SLS store name
* `store_type` - (Optional, ForceNew, String, Available since v2.2.0) Store type. Valid values: `LOG`, `METRIC`.

### `notify_config`

The notify_config supports the following:
* `active_days` - (Optional, ForceNew, Computed, List) The days of the week on which notifications are sent, 1-7 (type=DIRECT_NOTIFY). Defaults to `[1,2,3,4,5,6,7]`.
* `active_end_time` - (Optional, ForceNew, Computed, String) The daily end time of the effective notification period (HH:mm, type=DIRECT_NOTIFY). Defaults to `23:59`.
* `active_start_time` - (Optional, ForceNew, Computed, String) The daily start time of the effective notification period (HH:mm, type=DIRECT_NOTIFY). Defaults to `00:00`.
* `channels` - (Optional, ForceNew, Computed, List) The list of notification channels (type=DIRECT_NOTIFY) See [`channels`](#notify_config-channels) below.
* `notify_strategies` - (Optional, ForceNew, Computed, List) List of notification policy IDs (type=NOTIFY_POLICY, up to 1 for the current business)
* `silence_time_secs` - (Optional, ForceNew, Computed, Int) The channel silence period in seconds (type=DIRECT_NOTIFY). Defaults to `86400`.
* `type` - (Required, ForceNew, String) Notification configuration type. Valid values: `DIRECT_NOTIFY`, `NOTIFY_POLICY`.
* `utc_offset` - (Optional, ForceNew, Computed, String) UTC time zone offset (type=DIRECT_NOTIFY). Defaults to `+08:00`.

### `notify_config-channels`

The notify_config-channels supports the following:
* `identifiers` - (Optional, ForceNew, Computed, List) List of channel identifiers
* `type` - (Optional, ForceNew, Computed, String) Notification Channel Type

### `observe_resource_config`

The observe_resource_config supports the following:
* `entity_domain` - (Optional, ForceNew, String, Available since v2.2.0) UModel entity domain
* `entity_type` - (Optional, ForceNew, String, Available since v2.2.0) UModel entity type
* `namespace` - (Optional, String, Available since v2.2.0) Cloud monitoring namespace
* `product_category` - (Optional, String, Available since v2.2.0) Cloud monitoring product category
* `relation_type` - (Optional, String, Available since v2.2.0) Relation type. Valid values: `ALL`, `UMODEL_ENTITY`, `CLOUD_INSTANCE`, `GROUP_V1`, `GROUP_V2`, `TAG`.
* `resources` - (Optional, ForceNew, List, Available since v2.2.0) Resource list; can be empty when relation_type=ALL

### `query_config`

The query_config supports the following:
* `aggregate` - (Optional, ForceNew, String, Available since v2.2.0) Aggregate function (used when type=UMODEL_METRICSET_QUERY / UMODEL_LOGSET_QUERY). Valid values: `AVG`, `MAX`, `MIN`, `SUM`, `LAST`.
* `dimensions` - (Optional, Computed, List) The dimension list (type=CLOUD_MONITORING_QUERY). Each dimension is a key-value string mapping. See [`dimensions`](#query_config-dimensions) below.
* `enable_data_complete_check` - (Optional, ForceNew, Computed, Bool) Whether to Enable Data Completeness Check (type=PROMETHEUS_SINGLE_QUERY)
* `entity_domain` - (Optional, ForceNew, Computed, String) Domain to which the entity belongs (type=UMODEL_METRICSET_QUERY)
* `entity_fields` - (Optional, ForceNew, Computed, List) List of Entity Fields to Return (type=UMODEL_METRICSET_QUERY) See [`entity_fields`](#query_config-entity_fields) below.
* `entity_filters` - (Optional, ForceNew, Computed, List) The list of entity filter conditions (type=UMODEL_METRICSET_QUERY). See [`entity_filters`](#query_config-entity_filters) below.
* `entity_type` - (Optional, ForceNew, Computed, String) Entity type (type=UMODEL_METRICSET_QUERY)
* `expr` - (Optional, ForceNew, Computed, String) Prometheus query statement (type=PROMETHEUS_SINGLE_QUERY, recommended field)
* `filter_list` - (Optional, ForceNew, Computed, List) The APM filter condition list (type=APM_MULTI_QUERY). See [`filter_list`](#query_config-filter_list) below.
* `group_id` - (Optional, Computed, String) Resource group ID (used when type=CLOUD_MONITORING_QUERY and relation_type=GROUP)
* `group_type` - (Optional, ForceNew, Computed, String, Available since v2.2.0) Group policy (used when type=SLS_MULTI_QUERY). Valid values: `none`, `label`, `custom`.
* `joinings` - (Optional, ForceNew, Computed, List, Available since v2.2.0) Join list (used when type=SLS_MULTI_QUERY, at most 2: joinings[0] joins query 0 and query 1, joinings[1] joins query 1 and query 2). See [`joinings`](#query_config-joinings) below.
* `label_filters` - (Optional, ForceNew, Computed, List) List of label filter conditions (type=UMODEL_METRICSET_QUERY) See [`label_filters`](#query_config-label_filters) below.
* `legacy_raw` - (Computed, String) The raw V1 query JSON string returned when type=UNKNOWN_QUERY and parsing fails. The frontend displays this field as read-only only.
* `legacy_type` - (Computed, String) Returned when type=UNKNOWN_QUERY, indicating that this rule cannot be edited through the new API.
* `log_set` - (Optional, ForceNew, String, Available since v2.2.0) Logset name (used when type=UMODEL_LOGSET_QUERY; expr must provide the SPL query in that case)
* `measure_list` - (Optional, ForceNew, Computed, List) APM measure configuration list (type=APM_MULTI_QUERY) See [`measure_list`](#query_config-measure_list) below.
* `metric` - (Optional, ForceNew, Computed, String) Metric name (type=UMODEL_METRICSET_QUERY)
* `metric_set` - (Optional, ForceNew, Computed, String) Metric set name (type=UMODEL_METRICSET_QUERY)
* `namespace` - (Optional, Computed, String) CloudMonitor namespace (cloud service name, type=CLOUD_MONITORING_QUERY)
* `offset_secs` - (Optional, ForceNew, Int, Available since v2.2.0) Query time offset in seconds (used when type=UMODEL_METRICSET_QUERY / UMODEL_LOGSET_QUERY). Works with `window_secs` to query the shifted window [T - window_secs - offset_secs, T - offset_secs]. Valid values: `0-86400`.
* `prom_ql` - (Computed, String) [Deprecated] Legacy PromQL field. Use `expr` instead. The backend automatically normalizes this field to `expr`; it cannot be configured.
* `queries` - (Optional, ForceNew, Computed, List, Available since v2.2.0) Sub-query list (for type=SLS_MULTI_QUERY each entry describes one SPL sub-query with `expr`/`time_unit`/`start`/`end`/`window`). See [`queries`](#query_config-queries) below.
* `relation_type` - (Optional, Computed, String) Resource association type (type=CLOUD_MONITORING_QUERY). Valid values: `INSTANCE`, `GROUP`, `USER`.
* `service_id_list` - (Optional, ForceNew, Computed, List) Application Service ID List (type=APM_MULTI_QUERY)
* `type` - (Required, ForceNew, String) Query type. Valid values: `PROMETHEUS_SINGLE_QUERY`, `PROMETHEUS_MULTI_QUERY`, `SLS_MULTI_QUERY`, `APM_MULTI_QUERY`, `UMODEL_METRICSET_QUERY`, `UMODEL_LOGSET_QUERY`, `CLOUD_MONITORING_QUERY`, `UNKNOWN_QUERY`.
* `window_secs` - (Optional, ForceNew, Int, Available since v2.2.0) Aggregation window in seconds (used when type=UMODEL_METRICSET_QUERY / UMODEL_LOGSET_QUERY). Valid values: `60-86400`.

### `query_config-dimensions`

The query_config-dimensions supports the following:

Each entry is a map whose keys are dimension names and whose values are the corresponding dimension values (string to string).

### `query_config-entity_fields`

The query_config-entity_fields supports the following:
* `field` - (Optional, ForceNew, Computed, String) Entity Field Name
* `value` - (Optional, ForceNew, Computed, String) Entity Field Value

### `query_config-entity_filters`

The query_config-entity_filters supports the following:
* `field` - (Optional, ForceNew, Computed, String) The entity filter field name.
* `operator` - (Optional, ForceNew, Computed, String) The entity filter operator. Valid values: `=`, `!=`.
* `value` - (Optional, ForceNew, Computed, String) The entity filter value.

### `query_config-filter_list`

The query_config-filter_list supports the following:
* `key` - (Optional, ForceNew, Computed, String) APM filter dimension key
* `type` - (Optional, ForceNew, Computed, String) The APM filter type. Valid values: `ALL`, `EQ`, `NE`, `DISABLED`.
* `value` - (Optional, ForceNew, Computed, String) APM filter value (can be empty when type=ALL/DISABLED)

### `query_config-joinings`

The query_config-joinings supports the following:
* `conditions` - (Optional, ForceNew, Computed, List) Join condition list. See [`conditions`](#query_config-joinings-conditions) below.
* `type` - (Optional, ForceNew, Computed, String) Set operation type. Valid values: `InnerJoin`, `LeftJoin`, `RightJoin`, `FullJoin`, `LeftExclude`, `RightExclude`, `CrossJoin`, `NoJoin`, `Concat`.

### `query_config-joinings-conditions`

The query_config-joinings-conditions supports the following:
* `lhs_field` - (Optional, ForceNew, Computed, String) Left field in the format $<query_idx>.<field>, for example $0.hostIp
* `operator` - (Optional, ForceNew, Computed, String) Comparison operator. Valid values: `==`, `!=`, `<`, `>`, `<=`, `>=`.
* `rhs_field` - (Optional, ForceNew, Computed, String) Right field in the format $<query_idx>.<field>, for example $1.hostIp

### `query_config-label_filters`

The query_config-label_filters supports the following:
* `name` - (Optional, ForceNew, Computed, String) Label name
* `operator` - (Optional, ForceNew, Computed, String) Label filter operator. Valid values: `=`, `!=`.
* `value` - (Optional, ForceNew, Computed, String) Label value

### `query_config-measure_list`

The query_config-measure_list supports the following:
* `group_by` - (Optional, ForceNew, Computed, List) Grouping dimension list
* `measure_code` - (Optional, ForceNew, Computed, String) APM metric code
* `window_secs` - (Optional, ForceNew, Computed, Int) Query Time Window (Seconds)

### `query_config-queries`

The query_config-queries supports the following:
* `end` - (Optional, ForceNew, Computed, Int, Available since v2.2.0) End offset of the time range (mutually exclusive with `window`)
* `expr` - (Optional, ForceNew, Computed, String, Available since v2.2.0) SPL query expression (type=SLS_MULTI_QUERY)
* `start` - (Optional, ForceNew, Computed, Int, Available since v2.2.0) Start offset of the time range (mutually exclusive with `window`)
* `time_unit` - (Optional, ForceNew, Computed, String, Available since v2.2.0) Time unit (type=SLS_MULTI_QUERY). Valid values: `day`, `hour`, `minute`, `second`.
* `window` - (Optional, ForceNew, Computed, Int, Available since v2.2.0) Alignment window size (mutually exclusive with `start` and `end`)

### `rca_config`

The rca_config supports the following:
* `digital_employee_name` - (Optional, String, Available since v2.2.0) Digital employee name
* `enable_rca` - (Optional, ForceNew, Bool, Available since v2.2.0) Whether RCA is enabled

### `schedule_config`

The schedule_config supports the following:
* `interval_secs` - (Optional, ForceNew, Computed, Int) The scheduling interval in seconds. This parameter is used when the type is set to FIXED.
* `type` - (Required, ForceNew, String) The scheduling type. Valid values: `FIXED`.

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