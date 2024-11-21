---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_cluster_endpoint"
sidebar_current: "docs-alicloud-resource-poalrdb-cluster-endpoint"
description: |-
  Provides a PolarDB instance cluster endpoint resource.
---

# alicloud_polardb_cluster_endpoint

Provides a PolarDB endpoint resource to manage cluster endpoint of PolarDB cluster.

-> **NOTE:** Available since v1.217.0

-> **NOTE:** The default cluster endpoint can not be created or deleted manually.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_polardb_cluster_endpoint&exampleId=e13b5fcc-5a07-520c-4f55-bfcd06741fb0906fd342&activeTab=example&spm=docs.r.polardb_cluster_endpoint.0.e13b5fcc5a&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_polardb_node_classes" "default" {
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
  zone_id      = data.alicloud_polardb_node_classes.default.classes[0].zone_id
  vswitch_name = "terraform-example"
}

resource "alicloud_polardb_cluster" "default" {
  db_type       = "MySQL"
  db_version    = "8.0"
  db_node_class = data.alicloud_polardb_node_classes.default.classes.0.supported_engines.0.available_resources.0.db_node_class
  pay_type      = "PostPaid"
  vswitch_id    = alicloud_vswitch.default.id
  description   = "terraform-example"
}

resource "alicloud_polardb_cluster_endpoint" "default" {
  db_cluster_id = alicloud_polardb_cluster.default.id
}
```

## Argument Reference

The following arguments are supported:

* `db_cluster_id` - (Required, ForceNew) The Id of cluster that can run database.
* `read_write_mode` - (Optional) Read or write mode. Valid values are `ReadWrite`, `ReadOnly`. When creating a new custom endpoint, default to `ReadOnly`.
* `nodes` - (Optional) Node id list for endpoint configuration. At least 2 nodes if specified, or if the cluster has more than 3 nodes, read-only endpoint is allowed to mount only one node. Default is all nodes.
* `auto_add_new_nodes` - (Optional) Whether the new node automatically joins the default cluster address. Valid values are `Enable`, `Disable`. When creating a new custom endpoint, default to `Disable`.
* `endpoint_config` - (Optional) The advanced settings of the endpoint of Apsara PolarDB clusters are in JSON format. Including the settings of consistency level, transaction splitting, connection pool, and offload reads from primary node. For more details, see the [description of EndpointConfig in the Request parameters table for details](https://www.alibabacloud.com/help/doc-detail/116593.htm).
* `ssl_enabled` - (Optional) Specifies how to modify the SSL encryption status. Valid values: `Disable`, `Enable`, `Update`.
* `net_type` - (Optional) The network type of the endpoint address.
* `ssl_auto_rotate` - (Optional) Specifies whether automatic rotation of SSL certificates is enabled. Valid values: `Enable`,`Disable`.
  **NOTE:** For a PolarDB for MySQL cluster, this parameter is required, and only one connection string in each endpoint can enable the ssl, for other notes, see [Configure SSL encryption](https://www.alibabacloud.com/help/doc-detail/153182.htm).  
  For a PolarDB for PostgreSQL cluster or a PolarDB-O cluster, this parameter is not required, by default, SSL encryption is enabled for all endpoints.
* `db_endpoint_description` - (Optional) The name of the endpoint.
* `connection_prefix` - (Optional) Prefix of the specified endpoint. The prefix must be 6 to 30 characters in length, and can contain lowercase letters, digits, and hyphens (-), must start with a letter and end with a digit or letter.
* `port` - (Optional) Port of the specified endpoint. Valid values: 3000 to 5999.

## Attributes Reference

The following attributes are exported:

* `id` - The current instance connection resource ID. Composed of instance ID and connection string with format `<db_cluster_id>:<db_endpoint_id>`.
* `endpoint_type` - Type of endpoint.
* `ssl_connection_string` - The SSL connection string.
* `ssl_expire_time` - The time when the SSL certificate expires. The time follows the ISO 8601 standard in the yyyy-MM-ddTHH:mm:ssZ format. The time is displayed in UTC.
* `db_endpoint_id` - The ID of the cluster endpoint.
* `ssl_certificate_url` - The specifies SSL certificate download link.

## Import

PolarDB endpoint can be imported using the id, e.g.

```shell
$ terraform import alicloud_polardb_cluster_endpoint.example pc-abc123456:pe-abc123456
```
