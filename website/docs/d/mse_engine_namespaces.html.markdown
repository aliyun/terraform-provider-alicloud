---
subcategory: "Microservice Engine (MSE)"
layout: "alicloud"
page_title: "Alicloud: alicloud_mse_engine_namespaces"
sidebar_current: "docs-alicloud-datasource-mse-engine-namespaces"
description: |-
  Provides a list of Mse Engine Namespaces to the user.
---

# alicloud_mse_engine_namespaces

This data source provides the Mse Engine Namespaces of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.166.0.

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

# Declare the data source
data "alicloud_mse_engine_namespaces" "example" {
  instance_id = alicloud_mse_engine_namespace.example.instance_id
}

output "mse_engine_namespace_id_public" {
  value = data.alicloud_mse_engine_namespaces.example.namespaces.0.id
}

output "mse_engine_namespace_id_example" {
  value = data.alicloud_mse_engine_namespaces.example.namespaces.1.id
}
```

## Argument Reference

The following arguments are supported:

* `accept_language` - (Optional) The language type of the returned information. Valid values: `zh`, `en`.
* `ids` - (Optional, ForceNew, Computed)  A list of Engine Namespace IDs. It is formatted to `<instance_id>:<namespace_id>`.
* `cluster_id` - (Optional, ForceNew)  The ID of the cluster.
* `instance_id` - (Optional, ForceNew) The ID of the MSE Cluster Instance.It is formatted to `mse-cn-xxxxxxxxxxx`.Available since v1.232.0
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

**NOTE:** You must set `cluster_id` or `instance_id` or both.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `namespaces` - A list of Mse Engine Namespaces. Each element contains the following attributes:
  * `config_count` - The Number of Configuration of the Namespace.
  * `id` - The ID of the Engine Namespace. It is formatted to `<instance_id>:<namespace_id>`.
  * `namespace_id` - The id of Namespace.
  * `namespace_desc` - The description of the Namespace.
  * `namespace_show_name` - The name of the Namespace.
  * `quota` - The Quota of the Namespace.
  * `service_count` - The number of active services.
  * `type` - The type of the Namespace, the value is as follows:
    - '0': Global Configuration.
    - '1': default namespace.
    - '2': Custom Namespace.