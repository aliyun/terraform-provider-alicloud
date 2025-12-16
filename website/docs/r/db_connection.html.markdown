---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_db_connection"
sidebar_current: "docs-alicloud-resource-db-connection"
description: |-
  Provides an RDS instance connection resource.
---

# alicloud_db_connection

Provides an RDS connection resource to allocate an Internet connection string for RDS instance, see [What is DB Connection](https://www.alibabacloud.com/help/en/apsaradb-for-rds/latest/api-rds-2014-08-15-allocateinstancepublicconnection).

-> **NOTE:** Each RDS instance will allocate a intranet connnection string automatically and its prifix is RDS instance ID.
 To avoid unnecessary conflict, please specified a internet connection prefix before applying the resource.
 
-> **NOTE:** Available since v1.5.0.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_db_connection&exampleId=63fee0ed-49ad-ba93-36a8-8dd9f909e2df12b23b59&activeTab=example&spm=docs.r.db_connection.0.63fee0ed49&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf_example"
}

data "alicloud_db_zones" "default" {
  engine         = "MySQL"
  engine_version = "5.6"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_db_zones.default.zones.0.id
  vswitch_name = var.name
}

resource "alicloud_db_instance" "default" {
  engine           = "MySQL"
  engine_version   = "5.6"
  instance_type    = "rds.mysql.t1.small"
  instance_storage = "10"
  vswitch_id       = alicloud_vswitch.default.id
  instance_name    = var.name
}

resource "alicloud_db_connection" "default" {
  instance_id       = alicloud_db_instance.default.id
  connection_prefix = "testabc"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_db_connection&spm=docs.r.db_connection.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The Id of instance that can run database.
* `connection_prefix` - (Optional, ForceNew) Prefix of an Internet connection string. It must be checked for uniqueness. It may consist of lowercase letters, numbers, and underlines, and must start with a letter and have no more than 40 characters. Default to <instance_id> + 'tf'.
* `port` - (Optional) Internet connection port. Valid value: [1000-5999]. Default to 3306.
* `babelfish_port` - (Optional, Available since v1.176.0) The Tabular Data Stream (TDS) port of the instance for which Babelfish is enabled.

-> **NOTE:** This parameter applies only to ApsaraDB RDS for PostgreSQL instances. For more information about Babelfish for ApsaraDB RDS for PostgreSQL, see [Introduction to Babelfish](https://www.alibabacloud.com/help/en/apsaradb-for-rds/latest/babelfish-for-pg).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The current instance connection resource ID. Composed of instance ID and connection string with format `<instance_id>:<connection_prefix>`.
* `connection_string` - Connection instance string.
* `ip_address` - The ip address of connection string.

## Import

RDS connection can be imported using the id, e.g.

```shell
$ terraform import alicloud_db_connection.example abc12345678
```
