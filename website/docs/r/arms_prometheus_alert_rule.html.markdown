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
<div class="oics-button" style="float: right;margin: 0 0 -40px 0;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_arms_prometheus_alert_rule&exampleId=7427e925-2c7d-d88b-edd5-e967d06cf3b78599f916&activeTab=example&spm=docs.r.arms_prometheus_alert_rule.0.7427e9252c" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; margin: 32px auto; max-width: 100%;">
  </a>
</div>

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

data "alicloud_instance_types" "default" {
  availability_zone    = alicloud_vswitch.default.zone_id
  cpu_core_count       = 2
  memory_size          = 4
  kubernetes_node_role = "Worker"
  instance_type_family = "ecs.sn1ne"
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

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_key_pair" "default" {
  key_pair_name = "${var.name}-${random_integer.default.result}"
}

resource "alicloud_cs_kubernetes_node_pool" "default" {
  name                 = "desired_size"
  cluster_id           = alicloud_cs_managed_kubernetes.default.id
  vswitch_ids          = [alicloud_vswitch.default.id]
  instance_types       = [data.alicloud_instance_types.default.instance_types.0.id]
  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  key_name             = alicloud_key_pair.default.key_pair_name
  desired_size         = 2
}

resource "alicloud_arms_prometheus" "default" {
  cluster_type        = "aliyun-cs"
  grafana_instance_id = "free"
  cluster_id          = alicloud_cs_kubernetes_node_pool.default.cluster_id
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
