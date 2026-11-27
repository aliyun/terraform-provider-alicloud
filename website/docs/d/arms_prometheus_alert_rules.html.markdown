---
subcategory: "Application Real-Time Monitoring Service (ARMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_arms_prometheus_alert_rules"
sidebar_current: "docs-alicloud-datasource-arms-prometheus-alert-rules"
description: |-
  Provides a list of Arms Prometheus Alert Rules to the user.
---

# alicloud_arms_prometheus_alert_rules

This data source provides the Arms Prometheus Alert Rules of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.136.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_arms_prometheus_alert_rules" "ids" {
  cluster_id = "example_value"
  ids        = ["example_value-1", "example_value-2"]
}
output "arms_prometheus_alert_rule_id_1" {
  value = data.alicloud_arms_prometheus_alert_rules.ids.rules.0.id
}

data "alicloud_arms_prometheus_alert_rules" "nameRegex" {
  cluster_id = "example_value"
  name_regex = "^my-PrometheusAlertRule"
}
output "arms_prometheus_alert_rule_id_2" {
  value = data.alicloud_arms_prometheus_alert_rules.nameRegex.rules.0.id
}

```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) The ID of the cluster.
* `ids` - (Optional, Computed)  A list of Prometheus Alert Rule IDs.
* `match_expressions` - (Optional) Match expressions used to filter Prometheus alert rules.
* `name_regex` - (Optional) A regex string to filter results by Prometheus Alert Rule name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional) The status of the resource. Valid values: `0`, `1`.
* `type` - (Optional) The type of the alert rule.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Prometheus Alert Rule names.
* `rules` - A list of Arms Prometheus Alert Rules.

---

### `rules`

The rules supports the following attributes:

* `annotations` - The annotations of the alert rule.
* `cluster_id` - The ID of the cluster.
* `dispatch_rule_id` - The ID of the notification policy. This parameter is required when the `notify_type` parameter is set to `DISPATCH_RULE`.
* `duration` -The duration of the alert.
* `expression` - The alert rule expression that follows the PromQL syntax..
* `id` - The ID of the Prometheus Alert Rule.
* `labels` -The labels of the resource.
* `message` - The message of the alert notification.
* `notify_type` - The method of sending the alert notification. Valid values: `ALERT_MANAGER`, `DISPATCH_RULE`.
* `prometheus_alert_rule_id` - The first ID of the resource.
* `prometheus_alert_rule_name` - The name of the resource.
* `status` - The status of the resource. Valid values: `0`, `1` (0=off,1=open).
* `type` - The type of the alert rule.

---

### `rules.annotations`

The annotations supports the following attributes:

* `name` - The name of the annotation.
* `value` - The value of the annotation.

---

### `rules.labels`

The labels supports the following attributes:

* `name` - The name of the label.
* `value` - The value of the label.
