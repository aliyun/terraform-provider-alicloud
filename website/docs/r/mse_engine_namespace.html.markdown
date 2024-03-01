---
subcategory: "Microservice Engine (MSE)"
layout: "alicloud"
page_title: "Alicloud: alicloud_mse_engine_namespace"
sidebar_current: "docs-alicloud-resource-mse-engine-namespace"
description: |-
  Provides a Alicloud Microservice Engine (MSE) Engine Namespace resource.
---

# alicloud\_mse\_engine\_namespace

Provides a Microservice Engine (MSE) Engine Namespace resource.

For information about Microservice Engine (MSE) Engine Namespace and how to use it, see [What is Engine Namespace](https://www.alibabacloud.com/help/en/mse/developer-reference/api-mse-2019-05-31-createenginenamespace).

-> **NOTE:** Available in v1.166.0+.

## Example Usage

Basic Usage

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "tf-example"
}

data "alicloud_zones" "example" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "example" {
  vpc_name   = "terraform-example"
  cidr_block = "172.17.3.0/24"
}

resource "alicloud_vswitch" "example" {
  vswitch_name = "terraform-example"
  cidr_block   = "172.17.3.0/24"
  vpc_id       = alicloud_vpc.example.id
  zone_id      = data.alicloud_zones.example.zones.0.id
}

resource "alicloud_mse_cluster" "default" {
  connection_type       = "slb"
  net_type              = "privatenet"
  vswitch_id            = alicloud_vswitch.example.id
  cluster_specification = "MSE_SC_1_2_60_c"
  cluster_version       = "NACOS_2_0_0"
  instance_count        = "1"
  pub_network_flow      = "1"
  cluster_alias_name    = var.name
  mse_version           = "mse_dev"
  cluster_type          = "Nacos-Ans"
}

resource "alicloud_mse_engine_namespace" "example" {
  cluster_id          = alicloud_mse_cluster.default.id
  namespace_show_name = var.name
  namespace_id        = var.name
}
```

## Argument Reference

The following arguments are supported:

* `accept_language` - (Optional) The language type of the returned information. Valid values: `zh`, `en`.
* `cluster_id` - (Required, ForceNew) The id of the cluster.
* `namespace_id` - (Required, ForceNew) The id of Namespace.
* `namespace_show_name` - (Required) The name of the Engine Namespace.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Engine Namespace. It is formatted to `<cluster_id>:<namespace_id>`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when creating the Engine Namespace.
* `update` - (Defaults to 1 mins) Used when updating the Engine Namespace.
* `delete` - (Defaults to 1 mins) Used when deleting adb Engine Namespace.

## Import

Microservice Engine (MSE) Engine Namespace can be imported using the id, e.g.

```shell
$ terraform import alicloud_mse_engine_namespace.example <cluster_id>:<namespace_id>
```