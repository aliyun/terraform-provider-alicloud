---
subcategory: "Microservice Engine (MSE)"
layout: "alicloud"
page_title: "Alicloud: alicloud_mse_engine_namespace"
sidebar_current: "docs-alicloud-resource-mse-engine-namespace"
description: |-
  Provides a Alicloud Microservice Engine (MSE) Engine Namespace resource.
---

# alicloud_mse_engine_namespace

Provides a Microservice Engine (MSE) Engine Namespace resource.

For information about Microservice Engine (MSE) Engine Namespace and how to use it, see [What is Engine Namespace](https://www.alibabacloud.com/help/en/mse/developer-reference/api-mse-2019-05-31-createenginenamespace).

-> **NOTE:** Available since v1.166.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_mse_engine_namespace&exampleId=d901b9e6-1487-89f2-734a-abc449c9ee88c0765c97&activeTab=example&spm=docs.r.mse_engine_namespace.0.d901b9e614&intl_lang=EN_US" target="_blank">
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
  cluster_specification = "MSE_SC_1_2_60_c"
  cluster_type          = "Nacos-Ans"
  cluster_version       = "NACOS_2_0_0"
  instance_count        = 3
  net_type              = "privatenet"
  pub_network_flow      = "1"
  connection_type       = "slb"
  cluster_alias_name    = "terraform-example"
  mse_version           = "mse_pro"
  vswitch_id            = alicloud_vswitch.example.id
  vpc_id                = alicloud_vpc.example.id
}

resource "alicloud_mse_engine_namespace" "example" {
  instance_id         = alicloud_mse_cluster.example.id
  namespace_show_name = "terraform-example"
  namespace_id        = "terraform-example"
  namespace_desc      = "description"
}
```

## Argument Reference

The following arguments are supported:

* `accept_language` - (Optional) The language type of the returned information. Valid values: `zh`, `en`.
* `instance_id` - (Optional, ForceNew, Computed) The instance id of the cluster. It is formatted to `mse-cn-xxxxxxxxxxx`.Available since v1.232.0.
* `cluster_id` - (Optional since v1.232.0, ForceNew, Computed) The id of the cluster.It is formatted to `mse-xxxxxxxx`.
* `namespace_id` - (Optional since v1.232.0, ForceNew, Computed) The id of Namespace. 
* `namespace_show_name` - (Required) The name of the Engine Namespace.
* `namespace_desc` - (Optional, Computed, Available since v1.232.0)The description of the namespace.

**NOTE:** You must set `cluster_id` or `instance_id` or both.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Engine Namespace. It is formatted to `<instance_id>:<namespace_id>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when creating the Engine Namespace.
* `update` - (Defaults to 1 mins) Used when updating the Engine Namespace.
* `delete` - (Defaults to 1 mins) Used when deleting the Engine Namespace.

## Import

Microservice Engine (MSE) Engine Namespace can be imported using the id, e.g.

```shell
$ terraform import alicloud_mse_engine_namespace.example <instance_id>:<namespace_id>
```