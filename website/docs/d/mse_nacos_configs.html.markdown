---
subcategory: "Microservice Engine (MSE)"
layout: "alicloud"
page_title: "Alicloud: alicloud_mse_nacos_configs"
sidebar_current: "docs-alicloud-datasource-mse-nacos-configs"
description: |-
  Provides a list of Mse Nacos Configs to the user.
---

# alicloud\_mse\_nacos\_configs

This data source provides the Mse Nacos Configs of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.233.0.

## Example Usage

Basic Usage

```terraform
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

resource "alicloud_mse_cluster" "example" {
  connection_type       = "slb"
  net_type              = "privatenet"
  vswitch_id            = alicloud_vswitch.example.id
  cluster_specification = "MSE_SC_1_2_60_c"
  cluster_version       = "NACOS_2_0_0"
  instance_count        = "3"
  pub_network_flow      = "1"
  cluster_alias_name    = "example"
  mse_version           = "mse_pro"
  cluster_type          = "Nacos-Ans"
}

resource "alicloud_mse_engine_namespace" "example" {
  instance_id         = alicloud_mse_cluster.example.id
  namespace_show_name = "example"
  namespace_id        = "example"
}

resource "alicloud_mse_nacos_config" "example" {
  instance_id  = alicloud_mse_cluster.example.id
  data_id      = "example"
  group        = "example"
  namespace_id = alicloud_mse_engine_namespace.example.namespace_id
  content      = "example"
  type         = "text"
  tags         = "example"
  app_name     = "example"
  desc         = "example"
}

data "alicloud_mse_nacos_configs" "example" {
  instance_id    = alicloud_mse_cluster.example.id
  enable_details = "true"
  namespace_id   = alicloud_mse_engine_namespace.example.namespace_id
}
```

## Argument Reference

The following arguments are supported:
* `ids` - (Optional, ForceNew, Computed) A list of MSE Engine Configs ids. It is formatted to `<instance_id>:<namespace_id>:<data_id>:<group>`.
* `accept_language` - (Optional) The language type of the returned information. Valid values: `zh`, `en`.
* `instance_id` - (Required, ForceNew) The ID of the instance.
* `data_id` - (Optional, ForceNew) The ID of the data.
* `namespace_id` - (Optional, ForceNew) The id of Namespace.
* `group` - (Optional, ForceNew) The ID of the group.
* `app_name` - (Optional) The name of the application.
* `tags` - (Optional) The tags of the configuration.
* `request_pars` - (Optional) The extended request parameters. The JSON format is supported.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* 
## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `configs` - A list of Mse Nacos Configs. Each element contains the following attributes:
  * `id` -  The ID of the Nacos Config. It is formatted to `<instance_id>:<namespace_id>:<data_id>:<group>`.
  * `data_id` -  The ID of the data.
  * `group` -  The ID of the group.
  * `app_name` -  The name of the application.
  * `tags` -  The tags of the configuration.
  * `content` -  The content of the configuration.
  * `md5` - The message digest of the configuration.
  * `beta_ips` - The list of IP addresses where the beta release of the configuration is performed.
  * `desc` - The description of the configuration.
  * `encrypted_data_key` - The encryption key.
  * `type` - The format of the configuration. Supported formats include TEXT, JSON, and XML.

