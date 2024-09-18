---
subcategory: "Microservice Engine (MSE)"
layout: "alicloud"
page_title: "Alicloud: alicloud_mse_engine_config"
sidebar_current: "docs-alicloud-resource-mse-engine-config"
description: |-
  Provides a Alicloud Microservice Engine (MSE) Engine Config resource.
---

# alicloud_mse_engine_config

Provides a Microservice Engine (MSE) Engine Config resource.

For information about Microservice Engine (MSE) Engine Config and how to use it, see [What is Nacos configuration](https://www.alibabacloud.com/help/en/mse/developer-reference/api-mse-2019-05-31-createnacosconfig)
-> **NOTE:** Available since v1.230.1.

## Example Usage

Basic Usage

```terraform
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "example" {
  vpc_name   = "example"
  cidr_block = "172.17.3.0/24"
}

resource "alicloud_vswitch" "example" {
  vswitch_name = "example"
  cidr_block   = "172.17.3.0/24"
  vpc_id       = alicloud_vpc.example.id
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_mse_cluster" "example" {
  cluster_specification = "MSE_SC_1_2_60_c"
  cluster_type          = "Nacos-Ans"
  cluster_version       = "NACOS_2_0_0"
  instance_count        = 1
  net_type              = "privatenet"
  vswitch_id            = alicloud_vswitch.example.id
  connection_type       = "slb"
  pub_network_flow      = "1"
  mse_version           = "mse_dev"
  vpc_id                = alicloud_vpc.example.id
  cluster_alias_name    = "example"
}

data "alicloud_mse_clusters" "example" {
  ids    = alicloud_mse_cluster.example.clusters.0.id
  status = "INIT_SUCCESS"
}

resource "alicloud_mse_engine_namespace" "example" {
  instance_id         = data.alicloud_mse_clusters.example.clusters.0.instance_id
  namespace_show_name = "example"
  namespace_id        = "example"
}

resource "alicloud_mse_engine_config" "example" {
  instance_id  = data.alicloud_mse_clusters.tf.clusters.0.instance_id
  data_id      = "example"
  group        = "example"
  content      = "test"
  app_name     = "test"
  desc         = "test"
  type         = "text"
  napespace_id = "example"
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