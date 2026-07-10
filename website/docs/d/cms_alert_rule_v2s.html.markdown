---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_alert_rule_v2s"
sidebar_current: "docs-alicloud-datasource-cms-alert-rule-v2s"
description: |-
  Provides a list of Cms Alert Rule V2 owned by an Alibaba Cloud account.
---

# alicloud_cms_alert_rule_v2s

This data source provides Cms Alert Rule V2 available to the user.[What is Alert Rule V2](https://next.api.alibabacloud.com/document/Cms/2024-03-30/ManageAlertRules)

-> **NOTE:** Available since v1.285.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_cms_alert_rule_v2" "default" {
  content_template = "umodel test alert on $${metric}"
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
      identifiers = ["regression-test"]
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

data "alicloud_cms_alert_rule_v2s" "default" {
  ids = ["${alicloud_cms_alert_rule_v2.default.id}"]
}

output "alicloud_cms_alert_rule_v2_example_id" {
  value = data.alicloud_cms_alert_rule_v2s.default.v2s.0.id
}
```

## Argument Reference

The following arguments are supported:
* `filter_datasource_type_eq` - (Optional) Filter by exact match of data source type
* `filter_display_name_contains` - (Optional) Filter by display name using contains matching
* `filter_display_name_not_contains` - (Optional) Filter by display name that does not contain the specified value
* `filter_enabled_eq` - (Optional) Filter by exact match on enabled status.
* `filter_labels_all_of_key` - (Optional) Filter by full match of label keys (all specified keys must exist)
* `filter_labels_all_of_value` - (Optional) Filter by label values using full matching (all specified key-value pairs must match)
* `filter_labels_any_of_key` - (Optional) Filter by label keys using any matching (any of the specified keys exist)
* `filter_labels_any_of_value` - (Optional) Filter by matching any label value (matches if any specified key-value pair matches)
* `filter_notify_strategy_id_eq` - (Optional) Filter by Exact Match on Notification Policy ID
* `filter_observe_resource_global_scope_eq` - (Optional) Filters by exact match based on whether the rule applies globally to the resource type.
* `filter_observe_resource_instance_id` - (Optional) Filter by observable resource instance ID
* `filter_observe_resource_list_contains` - (Optional) Filter by matching based on the inclusion of observable resources in the list
* `filter_observe_resource_type_eq` - (Optional) Filters by exact match of the observable resource type.
* `filter_partition_key_eq` - (Optional) Filter by partition key using exact matching
* `filter_severity_levels_contains` - (Optional) Filter by severity levels that contain the specified value
* `filter_status_eq` - (Optional) Filter by alert status using exact matching
* `filter_uuid_eq` - (Optional) Filter by exact UUID match
* `filter_uuid_in` - (Optional) Filter by UUID list inclusion match. Separate multiple UUIDs with commas (,).
* `workspace` - (Optional) Workspace.
* `ids` - (Optional, Computed) A list of Alert Rule V2 IDs. 
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Alert Rule V2 IDs.
* `v2s` - A list of Alert Rule V2 Entries. Each element contains the following attributes:
    * `action_integration_config` - Action integration configuration.
        * `actions` - List of actions.
        * `enabled` - Indicates whether action integration is enabled.
    * `alert_rule_v2_id` - The unique identifier of the alert rule, mapped to a UUID (system-generated).
    * `annotations` - Annotations.
    * `arms_integration_config` - ARMS integration configuration.
        * `enabled` - Specifies whether to enable ARMS integration.
    * `condition_config` - Unified alert condition configuration.
        * `aggregate` - Aggregation function (used when APM_SIMPLE is specified).
        * `compare_list` - Comparison condition list (type=APM_COMPOSITE).
            * `aggregate` - Aggregation Function.
            * `operator` - Comparison Operator.
            * `threshold` - Threshold.
            * `yoy_time_unit` - Year-over-year time unit.
            * `yoy_time_value` - Year-over-year time value.
        * `composite_escalation` - The multi-metric composite trigger configuration.
            * `escalations` - List of trigger conditions.
                * `comparison_operator` - Comparison operator.
                * `metric_name` - Metric name.
                * `period` - Collection period (s).
                * `pre_condition` - Precondition.
                * `statistics` - Statistical Method.
                * `threshold` - Threshold.
            * `relation` - The logical relationship between multiple metrics.
            * `severity` - Severity Level.
            * `times` - Consecutive Trigger Count.
        * `duration_secs` - Duration (seconds).
        * `escalation_type` - The escalation policy type (type=CLOUD_MONITORING).
        * `express_escalation` - Expression trigger configuration.
            * `raw_expression` - Raw Expression.
            * `severity` - Severity Level.
            * `times` - Consecutive Trigger Count.
        * `legacy_raw` - The original V1 condition JSON string returned when type=UNKNOWN and parsing fails.
        * `legacy_type` - Returned when type=UNKNOWN, indicating that this rule cannot be edited by using the new API.
        * `no_data_policy` - No-data processing policy (type=CLOUD_MONITORING).
        * `operator` - The comparison operator.
        * `prometheus` - The PromQL trigger configuration.
            * `prom_ql` - Prometheus Query Expression.
            * `severity` - Severity Level.
            * `times` - Consecutive Trigger Count.
        * `relation` - The logical relationship between conditions (type=APM_COMPOSITE).
        * `severity` - The severity level.
        * `simple_escalation` - Single-metric trigger configuration (Required when type=CLOUD_MONITORING and escalationType=simple).
            * `escalations` - Trigger Condition List.
                * `comparison_operator` - Comparison Operator.
                * `pre_condition` - Precondition.
                * `severity` - Severity level.
                * `statistics` - Statistical Method.
                * `threshold` - Threshold.
                * `times` - Consecutive Trigger Count.
            * `metric_name` - Metric Name.
            * `period` - Collection Period (Seconds).
        * `threshold` - Threshold (used when UMODEL_METRICSET is specified).
        * `threshold_list` - Multi-severity threshold list (used when UMODEL_METRICSET is specified).
            * `severity` - Severity Level.
            * `threshold` - Threshold.
        * `type` - The detection condition type.
        * `yoy_time_unit` - Year-over-year time unit.
        * `yoy_time_value` - Year-over-Year Time Value (Effective only when type=APM_SIMPLE and operator=YOY_UP or YOY_DOWN).
    * `content_template` - The alert content template.
    * `created_at` - Creation time (read-only), in ISO 8601 format.
    * `datasource_config` - Unified data source configuration.
        * `instance_id` - The Prometheus instance ID.
        * `legacy_raw` - The original V1 datasource JSON string returned when type=UNKNOWN and parsing fails.
        * `legacy_type` - Returned when type=UNKNOWN, indicating that this rule cannot be edited by using the new API.
        * `product_category` - The cloud service category.
        * `region_id` - The region ID.
        * `type` - The data source type.
    * `datasource_type` - Data source type (read-only, derived).
    * `display_name` - The display name of the alert rule.
    * `enabled` - Specifies whether the alert rule is enabled.
    * `labels` - Labels.
    * `notify_config` - Unified notification configuration.
        * `active_days` - The days of the week on which notifications are sent, 1-7 (type=DIRECT_NOTIFY).
        * `active_end_time` - The daily end time of the effective notification period (HH:mm, type=DIRECT_NOTIFY).
        * `active_start_time` - The daily start time of the effective notification period (HH:mm, type=DIRECT_NOTIFY).
        * `channels` - List of notification channels (type=DIRECT_NOTIFY).
            * `identifiers` - List of channel identifiers.
            * `type` - Notification Channel Type.
        * `notify_strategies` - List of notification policy IDs (type=NOTIFY_POLICY, up to 1 for the current business).
        * `silence_time_secs` - The channel silence period in seconds (type=DIRECT_NOTIFY).
        * `type` - Notification configuration type.
        * `utc_offset` - UTC time zone offset (type=DIRECT_NOTIFY).
    * `notify_strategy_id` - Notification policy ID (read-only, derived).
    * `observe_resource_global_scope` - Indicates whether the rule applies to all resources of this resource type (read-only, derived).
    * `observe_resource_type` - Observable resource type (read-only, derived).
    * `partition_key` - The partition key.
    * `query_config` - Unified query configuration.
        * `dimensions` - The dimension list (type=CLOUD_MONITORING_QUERY).
        * `enable_data_complete_check` - Whether to Enable Data Completeness Check (type=PROMETHEUS_SINGLE_QUERY).
        * `entity_domain` - Domain to which the entity belongs (type=UMODEL_METRICSET_QUERY).
        * `entity_fields` - List of Entity Fields to Return (type=UMODEL_METRICSET_QUERY).
            * `field` - Entity Field Name.
            * `value` - Entity Field Value.
        * `entity_filters` - The list of entity filter conditions (type=UMODEL_METRICSET_QUERY).
            * `field` - The entity filter field name.
            * `operator` - The entity filter operator.
            * `value` - The entity filter value.
        * `entity_type` - Entity type (type=UMODEL_METRICSET_QUERY).
        * `expr` - Prometheus query statement (type=PROMETHEUS_SINGLE_QUERY, recommended field).
        * `filter_list` - The APM filter condition list (type=APM_MULTI_QUERY).
            * `key` - APM filter dimension key.
            * `type` - The APM filter type.
            * `value` - APM filter value (can be empty when type=ALL/DISABLED).
        * `group_id` - Resource group ID (used when type=CLOUD_MONITORING_QUERY and relationType=GROUP).
        * `label_filters` - List of label filter conditions (type=UMODEL_METRICSET_QUERY).
            * `name` - Label name.
            * `operator` - Label filter operator.
            * `value` - Label value.
        * `legacy_raw` - The raw V1 query JSON string returned when type=UNKNOWN_QUERY and parsing fails.
        * `legacy_type` - Returned when type=UNKNOWN_QUERY, indicating that this rule cannot be edited through the new API.
        * `measure_list` - APM measure configuration list (type=APM_MULTI_QUERY).
            * `group_by` - Grouping dimension list.
            * `measure_code` - APM metric code.
            * `window_secs` - Query Time Window (Seconds).
        * `metric` - Metric name (type=UMODEL_METRICSET_QUERY).
        * `metric_set` - Metric set name (type=UMODEL_METRICSET_QUERY).
        * `namespace` - CloudMonitor namespace (cloud service name, type=CLOUD_MONITORING_QUERY).
        * `prom_ql` - [Deprecated] Legacy PromQL field.
        * `relation_type` - Resource association type (type=CLOUD_MONITORING_QUERY).
        * `service_id_list` - Application Service ID List (type=APM_MULTI_QUERY).
        * `type` - Query type.
    * `schedule_config` - Unified scheduling configuration.
        * `interval_secs` - The scheduling interval in seconds.
        * `type` - The scheduling type.
    * `severity_levels` - The severity levels covered by this rule, separated by commas (read-only derived).
    * `status` - Alert status (read-only).
    * `updated_at` - The update time (read-only), in ISO 8601 format.
    * `workspace` - Workspace.
    * `id` - The ID of the resource supplied above.
