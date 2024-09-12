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
# Declare the data source

data "alicloud_mse_clusters" "example" {
  ids = ["mse-cn-xxxxxxxxxxx"]
  status = "INIT_SUCCESS"
}

output "instance_id" {
  value = "${data.alicloud_mse_clusters.example.clusters.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of MSE Cluster ids. It is formatted to `<instance_id>`
* `names` - (Optional)  A list of MSE Cluster names.
* `name_regex` - (Optional) A regex string to filter the results by the cluster alias name.
* `cluster_alias_name` - (Optional) The alias name of MSE Cluster.
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
