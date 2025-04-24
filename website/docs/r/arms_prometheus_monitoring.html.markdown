---
subcategory: "Application Real-Time Monitoring Service (ARMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_arms_prometheus_monitoring"
description: |-
  Provides a Alicloud ARMS Prometheus Monitoring resource.
---

# alicloud_arms_prometheus_monitoring

Provides a ARMS Prometheus Monitoring resource. Including serviceMonitor, podMonitor, customJob, probe and other four types of monitoring.

For information about ARMS Prometheus Monitoring and how to use it, see [What is Prometheus Monitoring](https://www.alibabacloud.com/help/en/arms/prometheus-monitoring/api-arms-2019-08-08-createprometheusmonitoring).

-> **NOTE:** Available since v1.209.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_arms_prometheus_monitoring&exampleId=ffdf6a46-95be-2df8-a19a-0c3ebaf8c3b1c83ce318&activeTab=example&spm=docs.r.arms_prometheus_monitoring.0.ffdf6a4695&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_vpc" "default" {
  cidr_block = "192.168.0.0/16"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "vswitch" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 8)
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = var.name
}

resource "alicloud_snapshot_policy" "default" {
  name            = var.name
  repeat_weekdays = ["1", "2", "3"]
  retention_days  = -1
  time_points     = ["1", "22", "23"]
}

data "alicloud_instance_types" "default" {
  availability_zone    = alicloud_vswitch.vswitch.zone_id
  cpu_core_count       = 2
  memory_size          = 4
  kubernetes_node_role = "Worker"
  instance_type_family = "ecs.sn1ne"
}

resource "alicloud_cs_managed_kubernetes" "default" {
  name               = var.name
  cluster_spec       = "ack.pro.small"
  version            = "1.24.6-aliyun.1"
  new_nat_gateway    = true
  node_cidr_mask     = 26
  proxy_mode         = "ipvs"
  service_cidr       = "172.23.0.0/16"
  pod_cidr           = "10.95.0.0/16"
  worker_vswitch_ids = [alicloud_vswitch.vswitch.id]
}

resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_key_pair" "default" {
  key_pair_name = "${var.name}-${random_integer.default.result}"
}

resource "alicloud_cs_kubernetes_node_pool" "default" {
  node_pool_name       = "desired_size"
  cluster_id           = alicloud_cs_managed_kubernetes.default.id
  vswitch_ids          = [alicloud_vswitch.vswitch.id]
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

resource "alicloud_arms_prometheus_monitoring" "default" {
  status      = "run"
  type        = "serviceMonitor"
  cluster_id  = alicloud_arms_prometheus.default.cluster_id
  config_yaml = <<EOF
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: tomcat-demo
  namespace: default
spec:
  endpoints:
  - bearerTokenSecret:
      key: ''
    interval: 30s
    path: /metrics
    port: tomcat-monitor
  namespaceSelector:
    any: true
  selector:
    matchLabels:
      app: tomcat
EOF
}
```

## Argument Reference

The following arguments are supported:
* `cluster_id` - (Required, ForceNew) The ID of the prometheus instance.
* `config_yaml` - (Required) Yaml configuration for monitoring.
* `status` - (Optional, Computed) Valid values: `stop`, `run`.
* `type` - (Required, ForceNew) Monitoring type: `serviceMonitor`, `podMonitor`, `customJob`, `probe`.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<cluster_id>:<monitoring_name>:<type>`.
* `monitoring_name` - The name of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Prometheus Monitoring.
* `delete` - (Defaults to 5 mins) Used when delete the Prometheus Monitoring.
* `update` - (Defaults to 5 mins) Used when update the Prometheus Monitoring.

## Import

ARMS Prometheus Monitoring can be imported using the id, e.g.

```shell
$ terraform import alicloud_arms_prometheus_monitoring.example <cluster_id>:<monitoring_name>:<type>
```