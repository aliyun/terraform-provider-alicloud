---
subcategory: "Data Security Center"
layout: "alicloud"
page_title: "Alicloud: alicloud_sddp_data_limit"
sidebar_current: "docs-alicloud-resource-sddp-data-limit"
description: |-
  Provides a Alicloud Data Security Center Data Limit resource.
---

# alicloud\_sddp\_data\_limit

Provides a Data Security Center Data Limit resource.

For information about Data Security Center Data Limit and how to use it, see [What is Data Limit](https://www.alibabacloud.com/help/en/doc-detail/158987.html).

-> **NOTE:** Available in v1.159.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tfaccxinmutes"
}

variable "region" {
  default = "cn-hangzhou"
}

variable "password" {
  default = "Test12345"
}

variable "database_name" {
  default = "tftestdatabase"
}

data "alicloud_db_zones" "default" {}

data "alicloud_db_instance_classes" "default" {
  engine         = "MySQL"
  engine_version = "5.6"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids[0]
  zone_id = data.alicloud_db_zones.default.zones[0].id
}

resource "alicloud_db_instance" "default" {
  engine           = "MySQL"
  engine_version   = "5.6"
  instance_type    = data.alicloud_db_instance_classes.default.instance_classes[0].instance_class
  instance_storage = "10"
  vswitch_id       = data.alicloud_vswitches.default.ids[0]
  instance_name    = var.name
}


locals {
  parent_id = join(".", [alicloud_db_instance.default.id, var.database_name])
}

resource "alicloud_rds_account" "default" {
  db_instance_id   = alicloud_db_instance.default.id
  account_name     = var.database_name
  account_password = var.password
}

resource "alicloud_db_database" "default" {
  instance_id = alicloud_db_instance.default.id
  name        = var.database_name
}

resource "alicloud_db_account_privilege" "default" {
  instance_id  = alicloud_db_instance.default.id
  account_name = alicloud_rds_account.default.name
  privilege    = "ReadWrite"
  db_names     = [alicloud_db_database.default.name]
}

resource "alicloud_sddp_data_limit" "default" {
  audit_status      = 0
  engine_type       = "MySQL"
  parent_id         = local.parent_id
  resource_type     = "RDS"
  user_name         = var.database_name
  password          = var.password
  port              = 3306
  service_region_id = var.region
  depends_on        = [alicloud_db_account_privilege.default]
}
```

## Argument Reference

The following arguments are supported:

* `audit_status` - (Optional, Computed)  Whether to enable the log auditing feature. Valid values: `0`, `1`.
* `engine_type` - (Optional, ForceNew) The type of the database. Valid values: `MySQL`, `SQLServer`.
* `lang` - (Optional) The lang.
* `log_store_day` - (Optional) The retention period of raw logs after you enable the log auditing feature. Unit: day. Valid values: `180`, `30`, `365`, `90`. **NOTE:** The`log_store_day` is valid when the `audit_status` is `1`.
* `parent_id` - (Optional, ForceNew) The ID of the data asset.
* `password` - (Optional, ForceNew) The password that is used to connect to the database.
* `port` - (Optional, ForceNew) The port that is used to connect to the database.
* `resource_type` - (Required, ForceNew) The type of the service to which the data asset belongs. Valid values: `MaxCompute`, `OSS`, `RDS`.
* `service_region_id` - (Optional, ForceNew) The region ID of the data asset.
* `user_name` - (Optional, ForceNew) The name of the service to which the data asset belongs.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Data Limit.

## Import

Data Security Center Data Limit can be imported using the id, e.g.

```
$ terraform import alicloud_sddp_data_limit.example <id>
```