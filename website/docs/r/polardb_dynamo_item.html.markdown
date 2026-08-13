---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_dynamo_item"
sidebar_current: "docs-alicloud-resource-polardb-dynamo-item"
description: |-
  Provides a PolarDB DynamoDB-compatible table item resource.
---

# alicloud_polardb_dynamo_item

Provides a PolarDB DynamoDB-compatible table item resource to manage a single item in a DynamoDB-compatible table of a PolarDB for PostgreSQL cluster.

-> **NOTE:** Available since v1.287.0.

-> **NOTE:** This resource is intended for managing a small amount of well-known, seed-style data. It is not recommended to manage large numbers of items with Terraform.

-> **NOTE:** All operations are performed against the DynamoDB-compatible endpoint (`http://<connection_string>:5432`) using the DynamoDB API, not the PolarDB OpenAPI.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_polardb_dynamo_item&exampleId=b2cfb953-c370-bb6a-0664-462e75287ae9847b9088&activeTab=example&spm=docs.r.polardb_dynamo_item.0.b2cfb953c3&intl_lang=EN_US" target="_blank">
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

resource "alicloud_polardb_dynamo_item" "default" {
  endpoint      = "http://${alicloud_polardb_endpoint_address.dynamo_public.connection_string}:5432"
  db_cluster_id = alicloud_polardb_cluster.default.id
  account_name  = alicloud_polardb_account.dynamo.account_name
  account_auth  = alicloud_polardb_account.dynamo.dynamodb_auth_password
  table_name    = alicloud_polardb_dynamo_table.default.table_name
  hash_key      = "pk"
  range_key     = "sk"

  item = jsonencode({
    pk    = { S = "test-item-1" }
    sk    = { S = "row1" }
    name  = { S = "Test Item" }
    count = { N = "42" }
  })
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_polardb_dynamo_item&spm=docs.r.polardb_dynamo_item.example&intl_lang=EN_US)


## Argument Reference

The following arguments are supported:

* `endpoint` - (Required) The PolarDB DynamoDB-compatible endpoint URL, in the format `http://<connection_string>:5432`.
* `account_name` - (Optional, Sensitive) The account name for PolarDB DynamoDB authentication. If not set, it is resolved from the cluster's DynamoDB-type account automatically.
* `account_auth` - (Optional, Sensitive) The authentication password for PolarDB DynamoDB. Usually references the `dynamodb_auth_password` attribute of an `alicloud_polardb_account` with `account_type = "DynamoDB"`. If not set, it is resolved from the cluster's DynamoDB-type account automatically.
* `db_cluster_id` - (Required, ForceNew) The ID of the PolarDB cluster where the DynamoDB table resides.
* `table_name` - (Required, ForceNew) The name of the DynamoDB-compatible table that contains the item.
* `hash_key` - (Required, ForceNew) The partition key (hash key) attribute name of the item. Must match the table's hash key.
* `range_key` - (Optional, ForceNew) The sort key (range key) attribute name of the item. Required if the table has a range key.
* `item` - (Required) JSON representation of the item attributes in DynamoDB attribute value format, e.g. `{"pk": {"S": "value"}, "count": {"N": "42"}}`. The item must contain the `hash_key` attribute (and the `range_key` attribute if set). Supported type descriptors: `S`, `N`, `B`, `BOOL`, `NULL`, `L`, `M`, `SS`, `NS`, `BS`.

-> **NOTE:** Changing the key attribute values inside `item` results in a new item being written; the resource ID will be recomputed accordingly.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of the item. Composed with format `<db_cluster_id>:<table_name>:<hash_key_value>` or `<db_cluster_id>:<table_name>:<hash_key_value>:<range_key_value>` when the table has a range key.

## Import

PolarDB DynamoDB-compatible table item can be imported using the id, e.g.

```shell
$ terraform import alicloud_polardb_dynamo_item.example pc-abc123456:table_name:hash_value:range_value
```

-> **NOTE:** On import, `account_name`, `account_auth`, `hash_key`, `range_key` and the endpoint address are resolved from the cluster and table schema automatically, but `endpoint`, `hash_key` and the other required arguments must still be present in the resource block.
