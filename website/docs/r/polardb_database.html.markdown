---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_database"
description: |-
  Provides a Alicloud Polar Db Database resource.
---

# alicloud_polardb_database

Provides a Polar Db Database resource.

Manage linked databases.

For information about Polar Db Database and how to use it, see [What is Database](https://next.api.alibabacloud.com/document/polardb/2017-08-01/CreateDatabase).

-> **NOTE:** Available since v1.66.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_polardb_database&exampleId=150536ad-ef26-da12-83db-3ddff47dbf191e1a5755&activeTab=example&spm=docs.r.polardb_database.0.150536adef&intl_lang=EN_US" target="_blank">
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

resource "alicloud_polardb_database" "default" {
  db_cluster_id = alicloud_polardb_cluster.default.id
  db_name       = "terraform-example"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_polardb_database&spm=docs.r.polardb_database.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `account_name` - (Optional) The name of the account that is authorized to access the database. **NOTE:** From version 1.265.0, `account_name` can be modified. However, only PolarDB for PostgreSQL (Compatible with Oracle) and PolarDB for PostgreSQL cluster can be modified.
* `character_set_name` - (Optional, ForceNew) The character set that is used by the cluster. For more information, see [Character set tables](https://www.alibabacloud.com/help/en/doc-detail/99716.html).
* `collate` - (Optional, Available since v1.265.0) The language that defines the collation rules in the database.
-> **NOTE:** The locale must be compatible with the character set set set by `character_set_name`. This parameter is required for a PolarDB for PostgreSQL (Compatible with Oracle) or PolarDB for PostgreSQL cluster. This parameter is optional for a PolarDB for MySQL cluster.
* `ctype` - (Optional, Available since v1.265.0) The language that indicates the character type of the database.
-> **NOTE:** The language must be compatible with the character set that is specified by `character_set_name`. The value that you specify must be the same as the value of `collate`. This parameter is required for PolarDB for PostgreSQL (Compatible with Oracle) clusters or PolarDB for PostgreSQL clusters. This parameter is optional for PolarDB for MySQL clusters.This parameter is required for a PolarDB for PostgreSQL (Compatible with Oracle) or PolarDB for PostgreSQL cluster. This parameter is optional for a PolarDB for MySQL cluster.
* `db_cluster_id` - (Required, ForceNew) The ID of cluster.
* `db_name` - (Required, ForceNew) The name of the database. It may consist of lower case letters, numbers, and underlines, and must start with a letterand have no more than 64 characters.
* `db_description` - (Optional) The description of the database. The description must meet the following requirements:
  - It cannot start with `http://` or `https://`.
  - It must be 2 to 256 characters in length.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<db_cluster_id>:<db_name>`.
* `status` - (Available since v1.265.0) The state of the database.

## Timeouts

-> **NOTE:** Available since v1.265.0.

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Database.
* `delete` - (Defaults to 30 mins) Used when delete the Database.
* `update` - (Defaults to 5 mins) Used when update the Database.

## Import

Polar Db Database can be imported using the id, e.g.

```shell
$ terraform import alicloud_polardb_database.example <db_cluster_id>:<db_name>
```
