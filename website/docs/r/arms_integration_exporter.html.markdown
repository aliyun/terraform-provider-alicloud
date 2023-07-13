---
subcategory: "Application Real-Time Monitoring Service (ARMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_arms_integration_exporter"
sidebar_current: "docs-alicloud-resource-arms-integration-exporter"
description: |-
  Provides a Alicloud Application Real-Time Monitoring Service (ARMS) Integration Exporter resource.
---

# alicloud_arms_integration_exporter

Provides a Application Real-Time Monitoring Service (ARMS) Integration Exporter resource.

For information about Application Real-Time Monitoring Service (ARMS) Integration Exporter and how to use it, see [What is Integration Exporter](https://www.alibabacloud.com/help/en/application-real-time-monitoring-service/latest/api-doc-arms-2019-08-08-api-doc-addprometheusintegration).

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

resource "alicloud_arms_integration_exporter" "default" {
  cluster_id       = alicloud_arms_prometheus.default.id
  integration_type = "kafka"
  param            = "{\"tls_insecure-skip-tls-verify\":\"none=tls.insecure-skip-tls-verify\",\"tls_enabled\":\"none=tls.enabled\",\"sasl_mechanism\":\"\",\"name\":\"kafka1\",\"sasl_enabled\":\"none=sasl.enabled\",\"ip_ports\":\"abc:888\",\"scrape_interval\":30,\"version\":\"0.10.1.0\"}"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, ForceNew) The ID of the Prometheus instance.
* `integration_type` - (Required, ForceNew) The type of prometheus integration.
* `param` - (Required) Exporter configuration parameter json string.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Integration Exporter. It formats as `<cluster_id>:<integration_type>:<instance_id>`.
* `instance_id` - The ID of the Integration Exporter instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the Integration Exporter.
* `update` - (Defaults to 3 mins) Used when update the Integration Exporter.
* `delete` - (Defaults to 3 mins) Used when delete the Integration Exporter.

## Import

Application Real-Time Monitoring Service (ARMS) Integration Exporter can be imported using the id, e.g.

```shell
$ terraform import alicloud_arms_integration_exporter.example <cluster_id>:<integration_type>:<instance_id>
```
