---
subcategory: "AnalyticDB for PostgreSQL (GPDB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_gpdb_jdbc_data_source"
description: |-
  Provides a Alicloud GPDB Jdbc Data Source resource.
---

# alicloud_gpdb_jdbc_data_source

Provides a GPDB Jdbc Data Source resource.



For information about GPDB Jdbc Data Source and how to use it, see [What is Jdbc Data Source](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.230.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_gpdb_jdbc_data_source&exampleId=f4d9cc1e-340b-3446-56c0-b2db388ef2e68dd75211&activeTab=example&spm=docs.r.gpdb_jdbc_data_source.0.f4d9cc1e34&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-beijing"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = "cn-beijing-h"
}

resource "alicloud_gpdb_instance" "defaulttuqTmM" {
  instance_spec         = "2C8G"
  description           = var.name
  seg_node_num          = "2"
  seg_storage_type      = "cloud_essd"
  instance_network_type = "VPC"
  payment_type          = "PayAsYouGo"
  ssl_enabled           = "0"
  engine_version        = "6.0"
  zone_id               = "cn-beijing-h"
  vswitch_id            = data.alicloud_vswitches.default.ids[0]
  storage_size          = "50"
  master_cu             = "4"
  vpc_id                = data.alicloud_vpcs.default.ids.0
  db_instance_mode      = "StorageElastic"
  engine                = "gpdb"
  db_instance_category  = "Basic"
}

resource "alicloud_gpdb_account" "defaultsk1eaS" {
  account_description = "example_001"
  db_instance_id      = alicloud_gpdb_instance.defaulttuqTmM.id
  account_name        = "example_001"
  account_password    = "example_001"
}

resource "alicloud_gpdb_external_data_service" "defaultRXkfKL" {
  service_name        = var.name
  db_instance_id      = alicloud_gpdb_instance.defaulttuqTmM.id
  service_description = "myexample"
  service_spec        = "8"
}

resource "alicloud_gpdb_jdbc_data_source" "default" {
  jdbc_connection_string  = "jdbc:mysql://rm-2ze327yr44c61183c.mysql.rds.aliyuncs.com:3306/example_001"
  data_source_description = "myexample"
  db_instance_id          = alicloud_gpdb_instance.defaulttuqTmM.id
  jdbc_password           = "example_001"
  data_source_name        = alicloud_gpdb_external_data_service.defaultRXkfKL.service_name
  data_source_type        = "mysql"
  jdbc_user_name          = "example_001"
}
```

## Argument Reference

The following arguments are supported:
* `data_source_description` - (Optional) Data Source Description
* `data_source_name` - (Optional, ForceNew) Data Source Name
* `data_source_type` - (Optional) Data Source Type
* `db_instance_id` - (Required, ForceNew) The instance ID.

* `jdbc_connection_string` - (Optional) The JDBC connection string.

* `jdbc_password` - (Optional) The password of the database account.

* `jdbc_user_name` - (Required) The name of the database account.


## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<db_instance_id>:<data_source_id>`.
* `create_time` - The creation time of the resource
* `data_source_id` - The data source ID.

* `status` - Data Source Status

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Jdbc Data Source.
* `delete` - (Defaults to 5 mins) Used when delete the Jdbc Data Source.
* `update` - (Defaults to 5 mins) Used when update the Jdbc Data Source.

## Import

GPDB Jdbc Data Source can be imported using the id, e.g.

```shell
$ terraform import alicloud_gpdb_jdbc_data_source.example <db_instance_id>:<data_source_id>
```