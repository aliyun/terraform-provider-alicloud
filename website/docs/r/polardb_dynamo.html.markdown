---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_dynamo"
sidebar_current: "docs-alicloud-resource-polardb-dynamo"
description: |-
  Provides a PolarDB DynamoDB-compatible table resource.
---

# alicloud_polardb_dynamo

Provides a PolarDB DynamoDB-compatible table resource to manage tables through the DynamoDB-compatible endpoint of a PolarDB for PostgreSQL cluster.

-> **NOTE:** Available since v1.286.0.

-> **NOTE:** This resource requires a PolarDB for PostgreSQL cluster with `enable_dynamodb` set to `true`, a PolarDB account of type `DynamoDB`, and a cluster endpoint of type `DynamoDB` with a reachable (e.g. public) endpoint address.

-> **NOTE:** All operations are performed against the DynamoDB-compatible endpoint (`http://<connection_string>:5432`) using the DynamoDB API, not the PolarDB OpenAPI.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_zones.default.zones.0.id
  vswitch_name = var.name
}

resource "alicloud_polardb_global_security_ip_group" "default" {
  global_ip_group_name = "tf_dynamo_whitelist"
  global_ip_list       = "0.0.0.0/0"
}

resource "alicloud_polardb_cluster" "default" {
  db_type                    = "PostgreSQL"
  db_version                 = "14"
  db_node_class              = "polar.pg.x4.medium"
  pay_type                   = "PostPaid"
  vswitch_id                 = alicloud_vswitch.default.id
  description                = var.name
  enable_dynamodb            = true
  global_security_group_list = [alicloud_polardb_global_security_ip_group.default.id]
}

resource "alicloud_polardb_account" "dynamo" {
  db_cluster_id    = alicloud_polardb_cluster.default.id
  account_name     = "tf_dynamo_acc"
  account_password = "Example1234!"
  account_type     = "DynamoDB"
}

resource "alicloud_polardb_endpoint" "dynamo" {
  db_cluster_id   = alicloud_polardb_account.dynamo.db_cluster_id
  endpoint_type   = "DynamoDB"
  read_write_mode = "ReadWrite"
}

resource "alicloud_polardb_endpoint_address" "dynamo_public" {
  db_cluster_id  = alicloud_polardb_cluster.default.id
  db_endpoint_id = alicloud_polardb_endpoint.dynamo.db_endpoint_id
  net_type       = "Public"
}

resource "alicloud_polardb_dynamo" "default" {
  endpoint      = "http://${alicloud_polardb_endpoint_address.dynamo_public.connection_string}:5432"
  db_cluster_id = alicloud_polardb_cluster.default.id
  account_name  = alicloud_polardb_account.dynamo.account_name
  account_auth  = alicloud_polardb_account.dynamo.dynamodb_auth_password
  table_name    = var.name
  hash_key      = "pk"
  range_key     = "sk"
  billing_mode  = "PAY_PER_REQUEST"

  attribute {
    name = "pk"
    type = "S"
  }
  attribute {
    name = "sk"
    type = "S"
  }
}
```

## Argument Reference

The following arguments are supported:

* `endpoint` - (Required) The PolarDB DynamoDB-compatible endpoint URL, in the format `http://<connection_string>:5432`.
* `account_name` - (Optional, Sensitive) The account name for PolarDB DynamoDB authentication. If not set, the provider's access key will be used.
* `account_auth` - (Optional, Sensitive) The authentication password for PolarDB DynamoDB. Usually references the `dynamodb_auth_password` attribute of an `alicloud_polardb_account` with `account_type = "DynamoDB"`. If not set, the provider's secret key will be used.
* `db_cluster_id` - (Required, ForceNew) The ID of the PolarDB cluster where DynamoDB is enabled.
* `table_name` - (Required, ForceNew) The name of the DynamoDB-compatible table.
* `attribute` - (Optional) List of attribute definitions for the table key schema and indexes. See [`attribute`](#attribute) below.
* `hash_key` - (Optional, ForceNew) The attribute name used as the partition key (hash key) of the table.
* `range_key` - (Optional, ForceNew) The attribute name used as the sort key (range key) of the table.
* `billing_mode` - (Optional) The billing mode of the table. Valid values: `PROVISIONED`, `PAY_PER_REQUEST`. Default to `PROVISIONED`.
* `read_capacity` - (Optional) The number of read capacity units. Required when `billing_mode` is `PROVISIONED`.
* `write_capacity` - (Optional) The number of write capacity units. Required when `billing_mode` is `PROVISIONED`.
* `global_secondary_index` - (Optional) Describe a GSI for the table. See [`global_secondary_index`](#global_secondary_index) below.
* `local_secondary_index` - (Optional, ForceNew) Describe an LSI on the table. See [`local_secondary_index`](#local_secondary_index) below.
* `stream_enabled` - (Optional) Whether Streams are enabled on the table.
* `stream_view_type` - (Optional) When a stream is enabled, how the data in the stream is written. Valid values: `NEW_IMAGE`, `OLD_IMAGE`, `NEW_AND_OLD_IMAGES`, `KEYS_ONLY`.
* `ttl` - (Optional) Configuration block for TTL. See [`ttl`](#ttl) below.
* `point_in_time_recovery` - (Optional) Configuration block for point-in-time recovery. See [`point_in_time_recovery`](#point_in_time_recovery) below.
* `server_side_encryption` - (Optional) Configuration block for server-side encryption. See [`server_side_encryption`](#server_side_encryption) below.
* `tags` - (Optional) A mapping of tags to assign to the table.

### `attribute`

The attribute supports the following:

* `name` - (Required) The name of the attribute.
* `type` - (Required) The attribute data type. Valid values: `S` (string), `N` (number), `B` (binary).

### `global_secondary_index`

The global_secondary_index supports the following:

* `name` - (Required) The name of the index.
* `hash_key` - (Optional) The attribute name used as the partition key of the index.
* `range_key` - (Optional) The attribute name used as the sort key of the index.
* `projection_type` - (Required) The set of attributes projected into the index. Valid values: `ALL`, `KEYS_ONLY`, `INCLUDE`.
* `non_key_attributes` - (Optional) A set of non-key attribute names projected into the index. Only valid when `projection_type` is `INCLUDE`.
* `read_capacity` - (Optional) The number of read capacity units for the index. Only valid when `billing_mode` is `PROVISIONED`.
* `write_capacity` - (Optional) The number of write capacity units for the index. Only valid when `billing_mode` is `PROVISIONED`.

### `local_secondary_index`

The local_secondary_index supports the following:

* `name` - (Required, ForceNew) The name of the index.
* `range_key` - (Required, ForceNew) The attribute name used as the sort key of the index.
* `projection_type` - (Required, ForceNew) The set of attributes projected into the index. Valid values: `ALL`, `KEYS_ONLY`, `INCLUDE`.
* `non_key_attributes` - (Optional, ForceNew) A list of non-key attribute names projected into the index. Only valid when `projection_type` is `INCLUDE`.

### `ttl`

The ttl supports the following:

* `enabled` - (Optional) Whether TTL is enabled. Default to `false`.
* `attribute_name` - (Optional) The name of the attribute that stores the TTL timestamp.

### `point_in_time_recovery`

The point_in_time_recovery supports the following:

* `enabled` - (Required) Whether point-in-time recovery is enabled.

### `server_side_encryption`

The server_side_encryption supports the following:

* `enabled` - (Required) Whether server-side encryption is enabled.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of the table. Composed of the cluster ID and the table name with format `<db_cluster_id>:<table_name>`.
* `arn` - The ARN of the table, if returned by the endpoint.
* `stream_arn` - The ARN of the table stream, if streams are enabled.
* `stream_label` - The timestamp-based label of the table stream, if streams are enabled.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 mins) Used when creating the table.
* `update` - (Defaults to 30 mins) Used when updating the table.
* `delete` - (Defaults to 10 mins) Used when deleting the table.

## Import

PolarDB DynamoDB-compatible table can be imported using the id, e.g.

```shell
$ terraform import alicloud_polardb_dynamo.example pc-abc123456:table_name
```

-> **NOTE:** Since `endpoint`, `account_name` and `account_auth` cannot be resolved from the resource ID, they must be configured in the resource block before running `terraform plan` after import.
