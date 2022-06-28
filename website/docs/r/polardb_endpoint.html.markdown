---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_endpoint"
sidebar_current: "docs-alicloud-resource-poalrdb-endpoint"
description: |-
  Provides a PolarDB instance endpoint resource.
---

# alicloud\_polardb\_endpoint

Provides a PolarDB endpoint resource to manage endpoint of PolarDB cluster.

-> **NOTE:** After v1.80.0 and before v1.121.0, you can only use this resource to manage the custom endpoint. Since v1.121.0, you also can import the primary endpoint and the cluster endpoint, to modify their ssl status and so on. 
 
-> **NOTE:** The primary endpoint and the default cluster endpoint can not be created or deleted manually.

## Example Usage

```terraform
variable "creation" {
  default = "PolarDB"
}

variable "name" {
  default = "polardbconnectionbasic"
}

data "alicloud_zones" "default" {
  available_resource_creation = var.creation
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_zones.default.zones[0].id
  vswitch_name = var.name
}

resource "alicloud_polardb_cluster" "default" {
  db_type       = "MySQL"
  db_version    = "8.0"
  pay_type      = "PostPaid"
  db_node_class = "polar.mysql.x4.large"
  vswitch_id    = alicloud_vswitch.default.id
  description   = var.name
}

resource "alicloud_polardb_endpoint" "endpoint" {
  db_cluster_id = alicloud_polardb_cluster.default.id
  endpoint_type = "Custom"
}
```

## Argument Reference 

The following arguments are supported:

* `db_cluster_id` - (Required, ForceNew) The Id of cluster that can run database.
* `endpoint_type` - (Optional, ForceNew) Type of the endpoint. Before v1.121.0, it only can be `Custom`. since v1.121.0, `Custom`, `Cluster`, `Primary` are valid, default to `Custom`. However when creating a new endpoint, it also only can be `Custom`. 
* `read_write_mode` - (Optional) Read or write mode. Valid values are `ReadWrite`, `ReadOnly`. When creating a new custom endpoint, default to `ReadOnly`.
* `nodes` - (Optional) Node id list for endpoint configuration. At least 2 nodes if specified, or if the cluster has more than 3 nodes, read-only endpoint is allowed to mount only one node. Default is all nodes.
* `auto_add_new_nodes` - (Optional) Whether the new node automatically joins the default cluster address. Valid values are `Enable`, `Disable`. When creating a new custom endpoint, default to `Disable`.
* `endpoint_config` - (Optional) The advanced settings of the endpoint of Apsara PolarDB clusters are in JSON format. Including the settings of consistency level, transaction splitting, connection pool, and offload reads from primary node. For more details, see the [description of EndpointConfig in the Request parameters table for details](https://www.alibabacloud.com/help/doc-detail/116593.htm).
* `ssl_enabled` - (Optional, Available in v1.121.0+) Specifies how to modify the SSL encryption status. Valid values: `Disable`, `Enable`, `Update`.
* `net_type` - (Optional, Available in v1.121.0+) The network type of the endpoint address.
* `ssl_auto_rotate` - (Available in v1.132.0+) Specifies whether automatic rotation of SSL certificates is enabled. Valid values: `Enable`,`Disable`.  
* `ssl_certificate_url` - (Available in v1.132.0+) Specifies SSL certificate download link.  
    **NOTE:** For a PolarDB for MySQL cluster, this parameter is required, and only one connection string in each endpoint can enable the ssl, for other notes, see [Configure SSL encryption](https://www.alibabacloud.com/help/doc-detail/153182.htm).  
    For a PolarDB for PostgreSQL cluster or a PolarDB-O cluster, this parameter is not required, by default, SSL encryption is enabled for all endpoints.

## Attributes Reference

The following attributes are exported:

* `id` - The current instance connection resource ID. Composed of instance ID and connection string with format `<db_cluster_id>:<db_endpoint_id>`.
* `endpoint_type` - Type of endpoint.
* `ssl_connection_string` - (Available in v1.121.0+) The SSL connection string.
* `ssl_expire_time` - (Available in v1.121.0+) The time when the SSL certificate expires. The time follows the ISO 8601 standard in the yyyy-MM-ddTHH:mm:ssZ format. The time is displayed in UTC.
* `db_endpoint_id` - (Available in v1.161.0+) The ID of the cluster endpoint.

## Import

PolarDB endpoint can be imported using the id, e.g.

```
$ terraform import alicloud_polardb_endpoint.example pc-abc123456:pe-abc123456
```
