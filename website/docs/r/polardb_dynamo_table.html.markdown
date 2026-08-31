---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_dynamo_table"
sidebar_current: "docs-alicloud-resource-polardb-dynamo-table"
description: |-
  Provides a PolarDB DynamoDB-compatible table resource.
---

# alicloud_polardb_dynamo_table

Provides a PolarDB DynamoDB-compatible table resource to manage tables through the DynamoDB-compatible endpoint of a PolarDB for PostgreSQL cluster.

-> **NOTE:** Available since v1.287.0.

-> **NOTE:** This resource requires a PolarDB for PostgreSQL cluster with `enable_dynamodb` set to `true`, a PolarDB account of type `DynamoDB`, and a cluster endpoint of type `DynamoDB` with a reachable (e.g. public) endpoint address.

-> **NOTE:** All operations are performed against the DynamoDB-compatible endpoint (`http://<connection_string>:5432`) using the DynamoDB API, not the PolarDB OpenAPI.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_polardb_dynamo_table&exampleId=49a57a77-1e9a-fae0-4062-1d715f6a8b73e6ba75a6&activeTab=example&spm=docs.r.polardb_dynamo_table.0.49a57a771e&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_polardb_node_classes" "default" {
  db_type       = "PostgreSQL"
  db_version    = "16"
  pay_type      = "PostPaid"
  db_node_class = "polar.pg.x4.medium"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_polardb_node_classes.default.classes.0.zone_id
  vswitch_name = var.name
}

resource "alicloud_polardb_global_security_ip_group" "default" {
  global_ip_group_name = "tf_dynamo_whitelist"
  global_ip_list       = "0.0.0.0/0"
}

resource "alicloud_polardb_cluster" "default" {
  db_type                    = "PostgreSQL"
  db_version                 = "16"
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

resource "alicloud_polardb_dynamo_table" "default" {
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


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_polardb_dynamo_table&spm=docs.r.polardb_dynamo_table.example&intl_lang=EN_US)


## Argument Reference

The following arguments are supported:

* `endpoint` - (Required) The PolarDB DynamoDB-compatible endpoint URL, in the format `http://<connection_string>:5432`.
* `account_name` - (Optional, Sensitive) The account name for PolarDB DynamoDB authentication. If not set, it is resolved from the cluster's DynamoDB-type account automatically.
* `account_auth` - (Optional, Sensitive) The authentication password for PolarDB DynamoDB. Usually references the `dynamodb_auth_password` attribute of an `alicloud_polardb_account` with `account_type = "DynamoDB"`. If not set, it is resolved from the cluster's DynamoDB-type account automatically.
* `db_cluster_id` - (Required, ForceNew) The ID of the PolarDB cluster where DynamoDB is enabled.
* `table_name` - (Required, ForceNew) The name of the DynamoDB-compatible table.
* `attribute` - (Optional) List of attribute definitions for the table key schema and indexes. See [`attribute`](#attribute) below.
* `hash_key` - (Optional, ForceNew) The attribute name used as the partition key (hash key) of the table.
* `range_key` - (Optional, ForceNew) The attribute name used as the sort key (range key) of the table.
* `billing_mode` - (Optional) The billing mode of the table. Valid values: `PROVISIONED`, `PAY_PER_REQUEST`. Default to `PROVISIONED`.
* `read_capacity` - (Optional) The number of read capacity units. Required when `billing_mode` is `PROVISIONED`.
* `write_capacity` - (Optional) The number of write capacity units. Required when `billing_mode` is `PROVISIONED`.
* `global_secondary_index` - (Optional) Describe a GSI for the table. See [`global_secondary_index`](#global_secondary_index) below. Changing the key schema or projection of an existing index recreates that index.
* `local_secondary_index` - (Optional, ForceNew) Describe an LSI on the table. See [`local_secondary_index`](#local_secondary_index) below.
* `ttl` - (Optional) Configuration block for TTL. See [`ttl`](#ttl) below.

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

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of the table. Composed of the cluster ID and the table name with format `<db_cluster_id>:<table_name>`.
* `arn` - The ARN of the table, if returned by the endpoint.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 mins) Used when creating the table.
* `update` - (Defaults to 30 mins) Used when updating the table.
* `delete` - (Defaults to 10 mins) Used when deleting the table.

## Import

PolarDB DynamoDB-compatible table can be imported using the id, e.g.

```shell
$ terraform import alicloud_polardb_dynamo_table.example pc-abc123456:table_name
```

-> **NOTE:** On import, `account_name`, `account_auth` and the endpoint address are resolved from the cluster automatically, but `endpoint` is a required argument and must still be present in the resource block. In addition, the DynamoDB-compatible endpoint does not return billing and capacity information, so `billing_mode`, `read_capacity` and `write_capacity` are not populated on import and the first plan may show a diff for them.
