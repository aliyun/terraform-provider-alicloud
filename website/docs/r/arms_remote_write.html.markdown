---
subcategory: "Application Real-Time Monitoring Service (ARMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_arms_remote_write"
description: |-
  Provides a Alicloud Application Real-Time Monitoring Service (ARMS) Remote Write resource.
---

# alicloud_arms_remote_write

Provides a Application Real-Time Monitoring Service (ARMS) Remote Write resource.

For information about Application Real-Time Monitoring Service (ARMS) Remote Write and how to use it, see [What is Remote Write](https://www.alibabacloud.com/help/en/arms/developer-reference/api-arms-2019-08-08-addprometheusremotewrite).

-> **NOTE:** Available since v1.204.0.

-> **DEPRECATED:** This resource has been deprecated since v1.228.0.

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

resource "alicloud_arms_remote_write" "default" {
  cluster_id        = alicloud_arms_prometheus.default.id
  remote_write_yaml = "remote_write:\n- name: ArmsRemoteWrite\n  url: http://47.96.227.137:8080/prometheus/xxx/yyy/cn-hangzhou/api/v3/write\n  basic_auth: {username: 666, password: '******'}\n  write_relabel_configs:\n  - source_labels: [instance_id]\n    separator: ;\n    regex: si-6e2ca86444db4e55a7c1\n    replacement: $1\n    action: keep\n"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, ForceNew) The ID of the Prometheus instance.
* `remote_write_yaml` - (Required) The details of the Remote Write configuration item. Specify the value in the YAML format.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Remote Write. It formats as `<cluster_id>:<remote_write_name>`.
* `remote_write_name` - The name of the Remote Write configuration item.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Remote Write.
* `delete` - (Defaults to 5 mins) Used when delete the Remote Write.
* `update` - (Defaults to 5 mins) Used when update the Remote Write.

## Import

Application Real-Time Monitoring Service (ARMS) Remote Write can be imported using the id, e.g.

```shell
$ terraform import alicloud_arms_remote_write.example <cluster_id>:<remote_write_name>
```