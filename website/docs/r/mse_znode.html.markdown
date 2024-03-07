---
subcategory: "Microservice Engine (MSE)"
layout: "alicloud"
page_title: "Alicloud: alicloud_mse_znode"
sidebar_current: "docs-alicloud-resource-mse-znode"
description: |-
  Provides a Alicloud Microservice Engine (MSE) Znode resource.
---

# alicloud\_mse\_znode

Provides a Microservice Engine (MSE) Znode resource.

For information about Microservice Engine (MSE) Znode and how to use it, see [What is Znode](https://help.aliyun.com/document_detail/393622.html).

-> **NOTE:** Available in v1.162.0+.

## Example Usage

Basic Usage

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
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

resource "alicloud_mse_cluster" "example" {
  cluster_specification = "MSE_SC_1_2_60_c"
  cluster_type          = "ZooKeeper"
  cluster_version       = "ZooKeeper_3_8_0"
  instance_count        = 1
  net_type              = "privatenet"
  pub_network_flow      = "1"
  acl_entry_list        = ["127.0.0.1/32"]
  cluster_alias_name    = "terraform-example"
  mse_version           = "mse_dev"
  vswitch_id            = alicloud_vswitch.example.id
  vpc_id                = alicloud_vpc.example.id
}

resource "alicloud_mse_znode" "example" {
  cluster_id = alicloud_mse_cluster.example.cluster_id
  data       = "terraform-example"
  path       = "/example"
}
```

## Argument Reference

The following arguments are supported:

* `accept_language` - (Optional) The language type of the returned information. Valid values: `zh` or `en`.
* `data` - (Optional) The Node data.
* `cluster_id` - (Required, ForceNew) The ID of the Cluster.
* `path` - (Required, ForceNew) The Node path. The value must start with a forward slash (/).

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Znode. The value formats as `<cluster_id>:<path>`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Znode.
* `delete` - (Defaults to 1 mins) Used when delete the Znode.
* `update` - (Defaults to 1 mins) Used when update the Znode.

## Import

Microservice Engine (MSE) Znode can be imported using the id, e.g.

```shell
$ terraform import alicloud_mse_znode.example <cluster_id>:<path>
```