---
subcategory: "Microservice Engine (MSE)"
layout: "alicloud"
page_title: "Alicloud: alicloud_mse_engine_namespaces"
sidebar_current: "docs-alicloud-datasource-mse-engine-namespaces"
description: |-
  Provides a list of Mse Engine Namespaces to the user.
---

# alicloud_mse_engine_configs

This data source provides the Mse Engine Configs of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.230.2.

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

data "alicloud_mse_engine_configs" "example" {
  instance_id    = data.alicloud_mse_clusters.example.clusters.0.instance_id
  enable_details = "true"
  napespace_id   = "example"
}

output "mse_engine_configs_example" {
  value = data.alicloud_mse_engine_configs.example
}
```

## Argument Reference

The following arguments are supported:
* `ids` - (Optional, ForceNew, Computed) A list of MSE Engine Configs ids. It is formatted to `<instance_id>:<namespace_id>:<data_id>:<group>`.
* `accept_language` - (Optional) The language type of the returned information. Valid values: `zh`, `en`.
* `instance_id` - (Required, ForceNew) The ID of the instance.
* `data_id` - (Optional) The ID of the data.
* `namespace_id` - (Optional, ForceNew) The id of Namespace.
* `group` - (Optional) The ID of the group.
* `app_name` - (Optional) The name of the application.
* `tags` - (Optional) The tags of the configuration.
* `request_pars` - (Optional) The extended request parameters. The JSON format is supported.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
The following attributes are exported in addition to the arguments listed above:

* `configs` - A list of Mse Engine Namespaces. Each element contains the following attributes:
  * `id` -  The ID of the Engine Config. It is formatted to `<instance_id>:<namespace_id>:<data_id>:<group>`.
  * `data_id` -  The ID of the data.
  * `namespace_id` -  The id of Namespace.
  * `group` -  The ID of the group.
  * `app_name` -  The name of the application.
  * `tags` -  The tags of the configuration.
  * `content` -  The content of the configuration.
  * `md5` - The message digest of the configuration.
  * `beta_ips` - The list of IP addresses where the beta release of the configuration is performed.
  * `desc` - The description of the configuration.
  * `encrypted_data_key` - The encryption key.

