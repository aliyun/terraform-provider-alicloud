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

For information about Microservice Engine (MSE) Engine Namespace and how to use it, see [What is Engine Namespace](https://www.alibabacloud.com/help/zh/microservices-engine/latest/api-doc-mse-2019-05-31-api-doc-createenginenamespace).

-> **NOTE:** Available in v1.166.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}
resource "alicloud_mse_cluster" "default" {
  cluster_specification = "MSE_SC_1_2_200_c"
  cluster_type          = "Nacos-Ans"
  cluster_version       = "NACOS_ANS_1_2_1"
  instance_count        = 1
  net_type              = "privatenet"
  vswitch_id            = data.alicloud_vswitches.default.ids.0
  pub_network_flow      = "1"
  acl_entry_list        = ["127.0.0.1/32"]
  cluster_alias_name    = "example_value"
}
resource "alicloud_mse_engine_namespace" "example" {
  cluster_id          = alicloud_mse_cluster.default.cluster_id
  namespace_show_name = "example_value"
  namespace_id        = "example_value"
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

```
$ terraform import alicloud_mse_engine_namespace.example <cluster_id>:<namespace_id>
```