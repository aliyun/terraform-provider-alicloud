---
subcategory: "Microservice Engine (MSE)"
layout: "alicloud"
page_title: "Alicloud: alicloud_mse_engine_config"
sidebar_current: "docs-alicloud-resource-mse-engine-config"
description: |-
  Provides a Alicloud Microservice Engine (MSE) Engine Config resource.
---

# alicloud\_mse\_engine\_config

Provides a Microservice Engine (MSE) Engine Config resource.

For information about Microservice Engine (MSE) Engine Config and how to use it, see [What is Nacos configuration](https://www.alibabacloud.com/help/en/mse/developer-reference/api-mse-2019-05-31-createnacosconfig)
-> **NOTE:** Available since v1.230.0.

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
  instance_id          = alicloud_mse_cluster.default.id
  namespace_show_name  = var.name
  namespace_id         = var.name
}

resource "alicloud_mse_engine_coonfig" "example" {
  instance_id          = alicloud_mse_cluster.default.id
  data_id              = var.name
  group                = var.name
  namespace_id         = var.name
}
```

## Argument Reference

The following arguments are supported:

* `accept_language` - (Optional) The language type of the returned information. Valid values: `zh`, `en`.
* `instance_id` - (Required, ForceNew) The ID of the instance.
* `data_id` - (Required) The ID of the data.
* `namespace_id` - (Optional, ForceNew) The id of Namespace. If you want to create a config under the `public` namespace, this parameter can be set to an empty string  *`""`* or just not set this parameter.
* `group` - (Required) The ID of the group.
* `app_name` - (Optional) The name of the application.
* `tags` - (Optional) The tags of the configuration.
* `Desc` - (Optional) The description of the configuration.
* `type` - (Optional) The format of the configuration. Supported formats include TEXT, JSON, and XML.
* `content` - (Optional) The content of the configuration.
* `beta_ips` - (Optional) The list of IP addresses where the beta release of the configuration is performed.
* `desc` - (Optional) The description of the configuration.


## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Engine Config. It is formatted to `<instance_id>:<namespace_id>:<data_id>:<group>`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when creating the Engine Config.
* `update` - (Defaults to 1 mins) Used when updating the Engine Config.
* `delete` - (Defaults to 1 mins) Used when deleting adb Engine Config.

## Import

Microservice Engine (MSE) Engine Config can be imported using the id, e.g.

```shell
$ terraform import alicloud_mse_engine_config.example <instance_id>:<namespace_id>:<data_id>:<group>
```