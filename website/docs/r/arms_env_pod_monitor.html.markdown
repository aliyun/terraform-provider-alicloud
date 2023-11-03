---
subcategory: "Application Real-Time Monitoring Service (ARMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_arms_env_pod_monitor"
description: |-
  Provides a Alicloud ARMS Env Pod Monitor resource.
---

# alicloud_arms_env_pod_monitor

Provides a ARMS Env Pod Monitor resource. PodMonitor for the arms environment.

For information about ARMS Env Pod Monitor and how to use it, see [What is Env Pod Monitor](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.212.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "vpc" {
  description = "api-resource-sub-test1-hz-pm"
  cidr_block  = "172.16.0.0/12"
  vpc_name    = var.name

}

resource "alicloud_vswitch" "vsw" {
  description  = "api-resource-test1-hz"
  vpc_id       = alicloud_vpc.vpc.id
  vswitch_name = var.name

  zone_id    = data.alicloud_zones.default.zones.0.id
  cidr_block = "172.16.0.0/24"
}

resource "alicloud_ack_cluster" "ask" {
  kubernetes_version = "1.26.3-aliyun.1"
  cluster_type       = "ManagedKubernetes"
  cluster_spec       = "ack.pro.small"
  vpc_id             = alicloud_vpc.vpc.id
  service_cidr       = "192.168.0.0/24"
  cluster_name       = var.name

  container_cidr = "192.168.1.0/24"
  vswitch_id     = alicloud_vswitch.vsw.id
  profile        = "Serverless"
}

resource "alicloud_arms_environment" "environment-cs" {
  environment_type = "CS"
  environment_name = var.name

  bind_resource_id     = alicloud_ack_cluster.ask.id
  environment_sub_type = "ACK"
}


resource "alicloud_arms_env_pod_monitor" "default" {
  environment_id = alicloud_arms_environment.environment-cs.id
  config_yaml    = <<EOF
apiVersion: monitoring.coreos.com/v1
kind: PodMonitor
metadata:
  name: arms-admin-pm1
  namespace: arms-prom
  annotations:
    arms.prometheus.io/discovery: 'true'
    o11y.aliyun.com/addon-name: mysql
    o11y.aliyun.com/addon-version: 1.0.2
    o11y.aliyun.com/release-name: mysql2
spec:
  selector:
    matchLabels:
      app: arms-prometheus-ack-arms-prometheus
      release: arms-prometheus
  namespaceSelector:
    any: true    
  podMetricsEndpoints:
  - interval: 30s
    targetPort: 9335
    path: /metrics
  - interval: 10s
    targetPort: 9335
    path: /metrics1
EOF
  aliyun_lang    = "zh"
}
```

## Argument Reference

The following arguments are supported:
* `aliyun_lang` - (Optional) Language environment, default is Chinese zh | en.
* `config_yaml` - (Required) Yaml configuration string.
* `environment_id` - (Required, ForceNew) Environment id.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<environment_id>:<namespace>:<env_pod_monitor_name>`.
* `env_pod_monitor_name` - The name of the resource.
* `namespace` - Resource namespace.
* `status` - PodMonitor status.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Env Pod Monitor.
* `delete` - (Defaults to 5 mins) Used when delete the Env Pod Monitor.
* `update` - (Defaults to 5 mins) Used when update the Env Pod Monitor.

## Import

ARMS Env Pod Monitor can be imported using the id, e.g.

```shell
$ terraform import alicloud_arms_env_pod_monitor.example <environment_id>:<namespace>:<env_pod_monitor_name>
```