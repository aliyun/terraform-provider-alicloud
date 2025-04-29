---
subcategory: "Microservice Engine (MSE)"
layout: "alicloud"
page_title: "Alicloud: alicloud_mse_nacos_config"
sidebar_current: "docs-alicloud-resource-mse-nacos-config"
description: |-
  Provides a Alicloud Microservice Engine (MSE) Nacos Config resource.
---

# alicloud_mse_nacos_config

Provides a Microservice Engine (MSE) Nacos Config resource.

For information about Microservice Engine (MSE) Nacos Config and how to use it, see [What is Nacos configuration](https://www.alibabacloud.com/help/en/mse/developer-reference/api-mse-2019-05-31-createnacosconfig)

-> **NOTE:** Available since v1.233.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_mse_nacos_config&exampleId=1aa53e35-384f-a407-db1e-7a19213cd273abcfc5b9&activeTab=example&spm=docs.r.mse_nacos_config.0.1aa53e3538&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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
```

## Argument Reference

The following arguments are supported:

* `accept_language` - (Optional) The language type of the returned information. Valid values: `zh`, `en`.
* `instance_id` - (Required, ForceNew) The ID of the instance.
* `data_id` - (Required, ForceNew) The ID of the data.
* `namespace_id` - (Optional, ForceNew) The id of Namespace. If you want to create a config under the `public` namespace, this parameter can be set to an empty string  *`""`* or just not set this parameter.
* `group` - (Required, ForceNew) The ID of the group.
* `app_name` - (Optional) The name of the application.
* `tags` - (Optional) The tags of the configuration.
* `type` - (Optional) The format of the configuration. Supported formats include TEXT, JSON, and XML.
* `content` - (Required) The content of the configuration.
* `desc` - (Optional) The description of the configuration.
* `beta_ips`-(Optional, Computed) The list of IP addresses where the beta release of the configuration is performed.


## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Nacos Config. It is formatted to `<instance_id>:<namespace_id>:<data_id>:<group>`.
* `encrypted_data_key` - The encryption key.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when creating the Nacos Config.
* `update` - (Defaults to 1 mins) Used when updating the Nacos Config.
* `delete` - (Defaults to 1 mins) Used when deleting adb Nacos Config.

## Import

Microservice Engine (MSE) Nacos Config can be imported using the id, e.g.

**Note**: If instance_id, namespace_id, data_id, and group contain ":", please replace it with "\\\\:", available since v1.243.0
```shell
$ terraform import alicloud_mse_nacos_config.example <instance_id>:<namespace_id>:<data_id>:<group>
```