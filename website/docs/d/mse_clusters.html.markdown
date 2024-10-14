---
subcategory: "Microservice Engine (MSE)"
layout: "alicloud"
page_title: "Alicloud: alicloud_mse_clusters"
sidebar_current: "docs-alicloud-resource-mse-clusters"
description: |-
    Provides a collection of MSE Clusters to the specified filters.
---

# alicloud_mse_clusters

This data source provides a list of MSE Clusters in an Alibaba Cloud account according to the specified filters.

-> **NOTE:** Available since v1.94.0.

## Example Usage

```terraform
# Create resource
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

# Declare the data source
data "alicloud_mse_clusters" "example" {
  enable_details = "true"
  ids            = [alicloud_mse_cluster.example.id]
  status         = "INIT_SUCCESS"
  name_regex     = alicloud_mse_cluster.example.cluster_alias_name
}

output "instance_id" {
  value = data.alicloud_mse_clusters.example.clusters.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew) A list of MSE Cluster ids. It is formatted to `<instance_id>`
* `names` - (Optional)  A list of MSE Cluster names.
* `name_regex` - (Optional) A regex string to filter the results by the cluster alias name.
* `cluster_alias_name` - (Optional, ForceNew) The alias name of MSE Cluster.
* `status` - (Optional) The status of MSE Cluster. Valid: `DESTROY_FAILED`, `DESTROY_ING`, `DESTROY_SUCCESS`, `INIT_FAILED`, `INIT_ING`, `INIT_SUCCESS`, `INIT_TIME_OUT`, `RESTART_FAILED`, `RESTART_ING`, `RESTART_SUCCESS`, `SCALE_FAILED`, `SCALE_ING`, `SCALE_SUCCESS`
* `request_pars` - (Optional) The extended request parameters. The JSON format is supported.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` -  A list of MSE Cluster ids.
* `names` -  A list of MSE Cluster names.
* `clusters` - A list of MSE Clusters. Each element contains the following attributes:
  * `id` - ID of the MSE Cluster.
  * `app_version` - The version of app.
  * `cluster_id` - ID of the MSE Cluster.
  * `cluster_name` - ID of the OOS Executions.
  * `cluster_type` - The type of MSE Cluster.
  * `instance_id` - ID of the MSE Cluster.
  * `internet_address` - The address of public network.
  * `internet_domain` - The domain of public network.
  * `intranet_address` - The address of private network.
  * `intranet_domain` - The domain of private network.
  * `instance_models` - The list of instances.
  * `status` - The status of MSE Cluster.
  * `acl_id` - The id of acl.
  * `cpu` - The num of cpu.
  * `health_status` - The health status of MSE Cluster.
  * `init_cost_time` - Time-consuming to create.
  * `instance_count` - The count of instance.
  * `internet_port` - The port of public network.
  * `intranet_port` - The port of private network.
  * `memory_capacity` - The memory size.
  * `pay_info` - The type of payment.
  * `pub_network_flow` - The public network bandwidth.
  * `instance_models` - The list of instance nodes.
    * `health_status` - The health status of the instance.
    * `vip` - (Deprecated from version 1.232.0)
    * `internet_ip` - The public IP address.
    * `single_tunnel_vip` - The single-thread IP address.
    * `pod_name` - The name of the pod.
    * `role` - The role.
    * `ip` - The IP address of the instance.
    * `instance_type` - (Deprecated from version 1.232.0)

