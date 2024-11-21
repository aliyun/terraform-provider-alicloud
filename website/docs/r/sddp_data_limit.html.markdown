---
subcategory: "Data Security Center (SDDP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sddp_data_limit"
sidebar_current: "docs-alicloud-resource-sddp-data-limit"
description: |-
  Provides a Alicloud Data Security Center Data Limit resource.
---

# alicloud_sddp_data_limit

Provides a Data Security Center Data Limit resource.

For information about Data Security Center Data Limit and how to use it, see [What is Data Limit](https://www.alibabacloud.com/help/en/doc-detail/158987.html).

-> **NOTE:** Available since v1.159.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_sddp_data_limit&exampleId=c92772ed-cd26-9642-6fe3-8952735afc0e7e647b83&activeTab=example&spm=docs.r.sddp_data_limit.0.c92772edcd&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}
variable "name" {
  default = "tf_example"
}
data "alicloud_regions" "default" {
  current = true
}
data "alicloud_db_zones" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_charge_type     = "PostPaid"
  category                 = "Basic"
  db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
  zone_id                  = data.alicloud_db_zones.default.zones.0.id
  engine                   = "MySQL"
  engine_version           = "8.0"
  category                 = "Basic"
  db_instance_storage_type = "cloud_essd"
  instance_charge_type     = "PostPaid"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_db_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_db_instance" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_type            = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
  instance_storage         = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
  instance_charge_type     = "Postpaid"
  instance_name            = var.name
  vswitch_id               = alicloud_vswitch.default.id
  monitoring_period        = "60"
  db_instance_storage_type = "cloud_essd"
  security_group_ids       = [alicloud_security_group.default.id]
}

resource "alicloud_rds_account" "default" {
  db_instance_id   = alicloud_db_instance.default.id
  account_name     = var.name
  account_password = "Example1234"
}

resource "alicloud_db_database" "default" {
  instance_id = alicloud_db_instance.default.id
  name        = var.name
}

resource "alicloud_db_account_privilege" "default" {
  instance_id  = alicloud_db_instance.default.id
  account_name = alicloud_rds_account.default.account_name
  privilege    = "ReadWrite"
  db_names     = [alicloud_db_database.default.name]
}

resource "alicloud_sddp_data_limit" "default" {
  audit_status      = 0
  engine_type       = "MySQL"
  parent_id         = join(".", [alicloud_db_account_privilege.default.instance_id, alicloud_db_database.default.name])
  resource_type     = "RDS"
  user_name         = alicloud_db_database.default.name
  password          = alicloud_rds_account.default.account_password
  port              = 3306
  service_region_id = data.alicloud_regions.default.regions.0.id
}
```

## Argument Reference

The following arguments are supported:

* `audit_status` - (Optional)  Whether to enable the log auditing feature. Valid values: `0`, `1`.
* `engine_type` - (Optional, ForceNew) The type of the database. Valid values: `MySQL`, `SQLServer`.
* `lang` - (Optional) The lang.
* `log_store_day` - (Optional, ForceNew) The retention period of raw logs after you enable the log auditing feature. Unit: day. Valid values: `180`, `30`, `365`, `90`. **NOTE:** The`log_store_day` is valid when the `audit_status` is `1`.
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

```shell
$ terraform import alicloud_sddp_data_limit.example <id>
```