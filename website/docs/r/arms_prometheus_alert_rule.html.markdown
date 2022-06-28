---
subcategory: "Application Real-Time Monitoring Service (ARMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_arms_prometheus_alert_rule"
sidebar_current: "docs-alicloud-resource-arms-prometheus-alert-rule"
description: |-
  Provides a Alicloud Application Real-Time Monitoring Service (ARMS) Prometheus Alert Rule resource.
---

# alicloud\_arms\_prometheus\_alert\_rule

Provides a Application Real-Time Monitoring Service (ARMS) Prometheus Alert Rule resource.

For information about Application Real-Time Monitoring Service (ARMS) Prometheus Alert Rule and how to use it, see [What is Prometheus Alert Rule](https://www.alibabacloud.com/help/en/doc-detail/212056.htm).

-> **NOTE:** Available in v1.136.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_arms_prometheus_alert_rule" "example" {
  cluster_id                 = "example_value"
  duration                   = "example_value"
  expression                 = "example_value"
  message                    = "example_value"
  prometheus_alert_rule_name = "example_value"
}

```

## Argument Reference

The following arguments are supported:

* `annotations` - (Optional) The annotations of the alert rule.. See the following `Block annotations`.
* `cluster_id` - (Required, ForceNew) The ID of the cluster.
* `dispatch_rule_id` - (Optional) The ID of the notification policy. This parameter is required when the `notify_type` parameter is set to `DISPATCH_RULE`.
* `duration` - (Required, ForceNew) The duration of the alert.
* `expression` - (Required, ForceNew) The alert rule expression that follows the PromQL syntax.
* `labels` - (Optional) The labels of the resource. See the following `Block labels`.
* `message` - (Required, ForceNew) The message of the alert notification.
* `notify_type` - (Optional) The method of sending the alert notification. Valid values: `ALERT_MANAGER`, `DISPATCH_RULE`.
* `prometheus_alert_rule_name` - (Required, ForceNew) The name of the resource.
* `type` - (Optional, ForceNew) The type of the alert rule.

#### Block labels

The labels supports the following: 

* `name` - (Optional) The name of the label.
* `value` - (Optional) The value of the label.

#### Block annotations

The annotations supports the following: 

* `name` - (Optional) The name of the annotation.
* `value` - (Optional) The value of the annotation.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Prometheus Alert Rule. The value formats as `<cluster_id>:<prometheus_alert_rule_id>`.
* `prometheus_alert_rule_id` - The first ID of the resource.
* `status` -  The status of the resource. Valid values: `0`, `1`.


## Import

Application Real-Time Monitoring Service (ARMS) Prometheus Alert Rule can be imported using the id, e.g.

```
$ terraform import alicloud_arms_prometheus_alert_rule.example <cluster_id>:<prometheus_alert_rule_id>
```
