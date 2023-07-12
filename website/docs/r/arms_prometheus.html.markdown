---
subcategory: "Application Real-Time Monitoring Service (ARMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_arms_prometheus"
sidebar_current: "docs-alicloud-resource-arms-prometheus"
description: |-
  Provides a Alicloud Application Real-Time Monitoring Service (ARMS) Prometheus resource.
---

# alicloud_arms_prometheus

Provides a Application Real-Time Monitoring Service (ARMS) Prometheus resource.

For information about Application Real-Time Monitoring Service (ARMS) Prometheus and how to use it, see [What is Prometheus](https://www.alibabacloud.com/help/en/application-real-time-monitoring-service/latest/api-doc-arms-2019-08-08-api-doc-createprometheusinstance).

-> **NOTE:** Available since v1.203.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf_example"
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
  zone_id      = data.alicloud_zones.default.zones[length(data.alicloud_zones.default.zones) - 1].id
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}

data "alicloud_resource_manager_resource_groups" "default" {
}

resource "alicloud_arms_prometheus" "default" {
  cluster_type        = "ecs"
  grafana_instance_id = "free"
  vpc_id              = alicloud_vpc.default.id
  vswitch_id          = alicloud_vswitch.default.id
  security_group_id   = alicloud_security_group.default.id
  cluster_name        = "${var.name}-${alicloud_vpc.default.id}"
  resource_group_id   = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  tags = {
    Created = "TF"
    For     = "Prometheus"
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_type` - (Required, ForceNew) The type of the Prometheus instance. Valid values: `remote-write`, `ecs`, `global-view`, `aliyun-cs`.
* `grafana_instance_id` - (Required) The ID of the Grafana dedicated instance. When using the shared version of Grafana, you can set `grafana_instance_id` to `free`.
* `vpc_id` - (Optional, ForceNew) The ID of the VPC. This parameter is required, if you set `cluster_type` to `ecs` or `aliyun-cs`(ASK instance).
* `vswitch_id` - (Optional, ForceNew) The ID of the VSwitch. This parameter is required, if you set `cluster_type` to `ecs` or `aliyun-cs`(ASK instance).
* `security_group_id` - (Optional, ForceNew) The ID of the security group. This parameter is required, if you set `cluster_type` to `ecs` or `aliyun-cs`(ASK instance).
* `cluster_id` - (Optional, ForceNew) The ID of the Kubernetes cluster. This parameter is required, if you set `cluster_type` to `aliyun-cs`.
* `cluster_name` - (Optional, ForceNew) The name of the created cluster. This parameter is required, if you set `cluster_type` to `remote-write`, `ecs` or `global-view`.
* `sub_clusters_json` - (Optional) The child instance json string of the globalView instance.
* `resource_group_id` - (Optional) The ID of the resource group.
* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Prometheus.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Prometheus.
* `update` - (Defaults to 5 mins) Used when update the Prometheus.
* `delete` - (Defaults to 5 mins) Used when delete the Prometheus.

## Import

Application Real-Time Monitoring Service (ARMS) Prometheus can be imported using the id, e.g.

```shell
$ terraform import alicloud_arms_prometheus.example <id>
```
