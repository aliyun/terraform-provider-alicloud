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
data "alicloud_mongodb_zones" "default" {}
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = "${data.alicloud_mongodb_zones.default.zones.0.id}"
}
resource "alicloud_mse_cluster" "default" {
  cluster_specification = "MSE_SC_1_2_200_c"
  cluster_type          = "ZooKeeper"
  cluster_version       = "ZooKeeper_3_5_5"
  instance_count        = 1
  net_type              = "privatenet"
  vswitch_id            = data.alicloud_vswitches.default.ids.0
  pub_network_flow      = "1"
  acl_entry_list        = ["127.0.0.1/32"]
  cluster_alias_name    = "example_value"
}

resource "alicloud_mse_znode" "default" {
  cluster_id = alicloud_mse_cluster.default.cluster_id
  data       = "example_value"
  path       = "example_value"
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

```
$ terraform import alicloud_mse_znode.example <cluster_id>:<path>
```