---
subcategory: "Application Real-Time Monitoring Service (ARMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_arms_prometheus_alert_rule"
sidebar_current: "docs-alicloud-resource-arms-prometheus-alert-rule"
description: |-
  Provides a Alicloud Application Real-Time Monitoring Service (ARMS) Prometheus Alert Rule resource.
---

# alicloud_arms_prometheus_alert_rule

Provides a Application Real-Time Monitoring Service (ARMS) Prometheus Alert Rule resource.

For information about Application Real-Time Monitoring Service (ARMS) Prometheus Alert Rule and how to use it, see [What is Prometheus Alert Rule](https://www.alibabacloud.com/help/en/arms/prometheus-monitoring/api-arms-2019-08-08-createprometheusalertrule).

-> **NOTE:** Available since v1.136.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_arms_prometheus_alert_rule&exampleId=c81ce7f0-4cf3-e06c-ff95-ffcf4c31108a3d4a6f83&activeTab=example&spm=docs.r.arms_prometheus_alert_rule.0.c81ce7f04c&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_arms_prometheus" "default" {
  cluster_type        = "remote-write"
  cluster_name        = "${var.name}-${random_integer.default.result}"
  grafana_instance_id = "free"
}

resource "alicloud_arms_prometheus_alert_rule" "example" {
  cluster_id                 = alicloud_arms_prometheus.default.cluster_id
  duration                   = 1
  expression                 = "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes * 100 < 10"
  message                    = "node available memory is less than 10%"
  prometheus_alert_rule_name = var.name
  notify_type                = "ALERT_MANAGER"
}
```

## Argument Reference

The following arguments are supported:

* `annotations` - (Optional) The annotations of the alert rule. See [`annotations`](#annotations) below.
* `cluster_id` - (Required, ForceNew) The ID of the cluster.
* `dispatch_rule_id` - (Optional) The ID of the notification policy. This parameter is required when the `notify_type` parameter is set to `DISPATCH_RULE`.
* `duration` - (Required, ForceNew) The duration of the alert.
* `expression` - (Required, ForceNew) The alert rule expression that follows the PromQL syntax.
* `labels` - (Optional) The labels of the resource. See [`labels`](#labels) below.
* `message` - (Required, ForceNew) The message of the alert notification.
* `notify_type` - (Optional) The method of sending the alert notification. Valid values: `ALERT_MANAGER`, `DISPATCH_RULE`.
* `prometheus_alert_rule_name` - (Required, ForceNew) The name of the resource.
* `type` - (Optional, ForceNew) The type of the alert rule.

### `labels`

The labels supports the following: 

* `name` - (Optional) The name of the label.
* `value` - (Optional) The value of the label.

### `annotations`

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

```shell
$ terraform import alicloud_arms_prometheus_alert_rule.example <cluster_id>:<prometheus_alert_rule_id>
```
