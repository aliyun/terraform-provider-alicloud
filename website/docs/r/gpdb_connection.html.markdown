---
subcategory: "AnalyticDB for PostgreSQL (GPDB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_gpdb_connection"
sidebar_current: "docs-alicloud-resource-gpdb-connection"
description: |-
  Provides an AnalyticDB for PostgreSQL instance connection resource.
---

# alicloud_gpdb_connection

Provides a connection resource to allocate an Internet connection string for instance.

-> **NOTE:** Available since v1.48.0.

-> **NOTE:** Each instance will allocate a intranet connection string automatically and its prefix is instance ID.
 To avoid unnecessary conflict, please specified a internet connection prefix before applying the resource.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_gpdb_connection&exampleId=58a46cfe-2619-aa98-09c7-4f43cf4af98ad0f79c2c&activeTab=example&spm=docs.r.gpdb_connection.0.58a46cfe26&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}
data "alicloud_resource_manager_resource_groups" "default" {}
data "alicloud_gpdb_zones" "default" {}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_gpdb_zones.default.ids.0
}

resource "alicloud_gpdb_instance" "default" {
  db_instance_category  = "HighAvailability"
  db_instance_class     = "gpdb.group.segsdx1"
  db_instance_mode      = "StorageElastic"
  description           = var.name
  engine                = "gpdb"
  engine_version        = "6.0"
  zone_id               = data.alicloud_gpdb_zones.default.ids.0
  instance_network_type = "VPC"
  instance_spec         = "2C16G"
  payment_type          = "PayAsYouGo"
  seg_storage_type      = "cloud_essd"
  seg_node_num          = 4
  storage_size          = 50
  vpc_id                = data.alicloud_vpcs.default.ids.0
  vswitch_id            = data.alicloud_vswitches.default.ids[0]
  ip_whitelist {
    security_ip_list = "127.0.0.1"
  }
}

resource "alicloud_gpdb_connection" "default" {
  instance_id       = alicloud_gpdb_instance.default.id
  connection_prefix = "exampelcon"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_gpdb_connection&spm=docs.r.gpdb_connection.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The Id of instance that can run database.
* `connection_prefix` - (Optional, ForceNew) Prefix of an Internet connection string. It must be checked for uniqueness. It may consist of lowercase letters, numbers, and underlines, and must start with a letter and have no more than 30 characters. Default to <instance_id> + '-tf'.
* `port` - (Optional) Internet connection port. Valid value: [3200-3999]. Default to 3306.
* `connection_string` - (Optional) Connection instance string.
* `ip_address` - (Optional) The ip address of connection string.

## Timeouts

-> **NOTE:** Available since v1.53.0.

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when creating the Internet connection (until DB instance reaches the initial `Running` status). 
* `update` - (Defaults to 10 mins) Used when activating the DB instance during update.
* `delete` - (Defaults to 10 mins) Used when terminating the DB instance. 

## Attributes Reference

The following attributes are exported:

* `id` - The current instance connection resource ID. Composed of instance ID and connection string with format `<instance_id>:<connection_prefix>`.

## Import

AnalyticDB for PostgreSQL's connection can be imported using the id, e.g.

```shell
$ terraform import alicloud_gpdb_connection.example abc12345678
```
