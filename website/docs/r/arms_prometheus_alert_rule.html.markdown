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

For information about Application Real-Time Monitoring Service (ARMS) Prometheus Alert Rule and how to use it, see [What is Prometheus Alert Rule](https://www.alibabacloud.com/help/en/doc-detail/212056.htm).

-> **NOTE:** Available since v1.136.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_cs_managed_kubernetes" "default" {
  name_prefix          = var.name
  cluster_spec         = "ack.pro.small"
  worker_vswitch_ids   = [alicloud_vswitch.default.id]
  new_nat_gateway      = true
  pod_cidr             = cidrsubnet("10.0.0.0/8", 8, 36)
  service_cidr         = cidrsubnet("172.16.0.0/16", 4, 7)
  slb_internet_enabled = true
}

resource "alicloud_arms_prometheus_alert_rule" "example" {
  cluster_id                 = alicloud_cs_managed_kubernetes.default.id
  duration                   = 1
  expression                 = "node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes * 100 < 10"
  message                    = "node available memory is less than 10%"
  prometheus_alert_rule_name = var.name
  notify_type                = "DISPATCH_RULE"
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
