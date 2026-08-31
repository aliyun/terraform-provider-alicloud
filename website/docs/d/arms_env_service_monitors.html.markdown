---
subcategory: "Application Real-Time Monitoring Service (ARMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_arms_env_service_monitors"
description: |-
  Provides a list of ARMS Env Service Monitors to the user.
---

# alicloud_arms_env_service_monitors

This data source provides the ARMS Env Service Monitors of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.258.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

data "alicloud_enhanced_nat_available_zones" "enhanced" {
}

resource "alicloud_vpc" "vpc" {
  description = "api-resource-test1-hz"
  cidr_block  = "192.168.0.0/16"
  vpc_name    = "${var.name}-${random_integer.default.result}"
}

resource "alicloud_vswitch" "vswitch" {
  description  = "api-resource-test1-hz"
  vpc_id       = alicloud_vpc.vpc.id
  vswitch_name = "${var.name}-${random_integer.default.result}"
  zone_id      = data.alicloud_enhanced_nat_available_zones.enhanced.zones.0.zone_id
  cidr_block   = cidrsubnet(alicloud_vpc.vpc.cidr_block, 8, 8)
}


resource "alicloud_snapshot_policy" "default" {
  name            = "${var.name}-${random_integer.default.result}"
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
  name               = "${var.name}-${random_integer.default.result}"
  cluster_spec       = "ack.pro.small"
  version            = "1.24.6-aliyun.1"
  new_nat_gateway    = true
  node_cidr_mask     = 26
  proxy_mode         = "ipvs"
  service_cidr       = "172.23.0.0/16"
  pod_cidr           = "10.95.0.0/16"
  worker_vswitch_ids = [alicloud_vswitch.vswitch.id]
}

resource "alicloud_key_pair" "default" {
  key_pair_name = "${var.name}-${random_integer.default.result}"
}

resource "alicloud_cs_kubernetes_node_pool" "default" {
  name                 = "desired_size"
  cluster_id           = alicloud_cs_managed_kubernetes.default.id
  vswitch_ids          = [alicloud_vswitch.vswitch.id]
  instance_types       = [data.alicloud_instance_types.default.instance_types.0.id]
  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  key_name             = alicloud_key_pair.default.key_pair_name
  desired_size         = 2
}

resource "alicloud_arms_environment" "default" {
  environment_type     = "CS"
  environment_name     = "${var.name}-${random_integer.default.result}"
  bind_resource_id     = alicloud_cs_kubernetes_node_pool.default.cluster_id
  environment_sub_type = "ManagedKubernetes"
}

resource "alicloud_arms_env_service_monitor" "default" {
  aliyun_lang    = "en"
  environment_id = alicloud_arms_environment.default.id
  config_yaml    = <<EOF
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: arms-admin1
  namespace: arms-prom
  annotations:
    arms.prometheus.io/discovery: 'true'
    o11y.aliyun.com/addon-name: mysql
    o11y.aliyun.com/addon-version: 1.0.1
    o11y.aliyun.com/release-name: mysql1
spec:
  endpoints:
  - interval: 30s
    port: operator
    path: /metrics
  - interval: 10s
    port: operator1
    path: /metrics
  namespaceSelector:
    any: true
  selector:
    matchLabels:
     app: arms-prometheus-ack-arms-prometheus
EOF
}

data "alicloud_arms_env_service_monitors" "ids" {
  environment_id = alicloud_arms_env_service_monitor.default.environment_id
  ids            = [alicloud_arms_env_service_monitor.default.id]
}

output "arms_env_service_monitors_id_0" {
  value = data.alicloud_arms_env_service_monitors.ids.monitors.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, List) A list of ARMS Env Service Monitor IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by ARMS Env Service Monitor name.
* `environment_id` - (Required, ForceNew) The environment ID.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of ARMS Env Service Monitor names.
* `monitors` - A list of ARMS Env Service Monitors. Each element contains the following attributes:
  * `id` - The ID of the ServiceMonitor. It formats as `<environment_id>:<namespace>:<env_service_monitor_name>`.
  * `config_yaml` - The YAML configuration string.
  * `env_service_monitor_name` - The name of the ServiceMonitor.
  * `environment_id` - The environment ID.
  * `namespace` - The namespace.
  * `region_id` - The region ID.
  * `status` - The status of the ServiceMonitor.
