---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_endpoints"
sidebar_current: "docs-alicloud-datasource-polardb-endpoints"
description: |-
    Provides a collection of PolarDB endpoints according to the specified filters.
---

# alicloud\_polardb\_endpoints

The `alicloud_polardb_endpoints` data source provides a collection of PolarDB endpoints available in Alibaba Cloud account.
Filters support regular expression for the cluster name, searches by clusterId, and other filters which are listed below.

-> **NOTE:** Available since v1.68.0+.

## Example Usage

```terraform

data "alicloud_polardb_node_classes" "this" {
  db_type    = "MySQL"
  db_version = "8.0"
  pay_type   = "PostPaid"
  category   = "Normal"
}

resource "alicloud_vpc" "default" {
  vpc_name   = "terraform-example"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_polardb_node_classes.this.classes[0].zone_id
  vswitch_name = "terraform-example"
}

resource "alicloud_polardb_cluster" "cluster" {
  db_type       = "MySQL"
  db_version    = "8.0"
  pay_type      = "PostPaid"
  db_node_count = "2"
  db_node_class = data.alicloud_polardb_node_classes.this.classes.0.supported_engines.0.available_resources.0.db_node_class
  vswitch_id    = alicloud_vswitch.default.id
}

data "alicloud_polardb_clusters" "polardb_clusters_ds" {
  description_regex = alicloud_polardb_cluster.cluster.description
  status            = "Running"
}

data "alicloud_polardb_endpoints" "default" {
  db_cluster_id = data.alicloud_polardb_clusters.polardb_clusters_ds.clusters.0.id
}

output "endpoint" {
  value = data.alicloud_polardb_endpoints.default.endpoints[0].db_endpoint_id
}
```

## Argument Reference

The following arguments are supported:

* `db_cluster_id` - (Required, ForceNew) PolarDB cluster ID. 
* `db_endpoint_id` - (Optional) endpoint of the cluster.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `endpoints` - A list of PolarDB cluster endpoints. Each element contains the following attributes:
  * `db_endpoint_id` - The endpoint ID.
  * `auto_add_new_nodes` - Whether the new node is automatically added to the default cluster address.Options are `Enable` and `Disable`.
  * `endpoint_config` - The Endpoint configuration. `ConsistLevel`: session consistency level, value:`0`: final consistency,`1`: session consistency;`LoadBalanceStrategy`: load balancing strategy. Based on the automatic scheduling of load, the value is: `load`.
  * `endpoint_type` - Cluster address type.`Cluster`: the default address of the Cluster.`Primary`: Primary address.`Custom`: Custom cluster addresses.
  * `nodes` - A list of nodes that connect to the address configuration.
  * `read_write_mode` - Read-write mode:`ReadWrite`: readable and writable (automatic read-write separation).`ReadOnly`: ReadOnly.
  * `address_items` - A list of endpoint addresses. Each element contains the following attributes.
      * `net_type` - IP network type:`Public` or `Private`.
      * `connection_string` - Connection instance string.
      * `port` - Intranet connection port.
      * `vpc_id` - ID of the VPC the instance belongs to.
      * `vswitch_id` - ID of the vSwitch the cluster belongs to.
      * `ip_address` - The ip address of connection string.
