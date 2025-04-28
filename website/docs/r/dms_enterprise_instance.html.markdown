---
subcategory: "DMS Enterprise"
layout: "alicloud"
page_title: "Alicloud: alicloud_dms_enterprise_instance"
sidebar_current: "docs-alicloud-resource-dms-enterprise-instance"
description: |-
  Provides a DMS Enterprise Instance resource.
---

# alicloud_dms_enterprise_instance

Provides a DMS Enterprise Instance resource.

-> **NOTE:** API users must first register in DMS.

-> **NOTE:** Available since v1.81.0.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_dms_enterprise_instance&exampleId=68231718-da4c-3285-5565-c6fdfddf31a6b88a2fa2&activeTab=example&spm=docs.r.dms_enterprise_instance.0.68231718da&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}

data "alicloud_account" "current" {}
data "alicloud_regions" "default" {
  current = true
}
data "alicloud_dms_user_tenants" "default" {
  status = "ACTIVE"
}

data "alicloud_db_zones" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_charge_type     = "PostPaid"
  category                 = "HighAvailability"
  db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
  zone_id                  = data.alicloud_db_zones.default.zones.0.id
  engine                   = "MySQL"
  engine_version           = "8.0"
  category                 = "HighAvailability"
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
  db_instance_storage_type = "cloud_essd"
  instance_type            = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
  instance_storage         = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
  vswitch_id               = alicloud_vswitch.default.id
  instance_name            = var.name
  security_ips             = ["100.104.5.0/24", "192.168.0.6"]
  tags = {
    Created = "TF",
    For     = "example",
  }
}

resource "alicloud_db_account" "default" {
  db_instance_id   = alicloud_db_instance.default.id
  account_name     = "tfexamplename"
  account_password = "Example12345"
  account_type     = "Normal"
}

resource "alicloud_dms_enterprise_instance" "default" {
  tid               = data.alicloud_dms_user_tenants.default.ids.0
  instance_type     = "mysql"
  instance_source   = "RDS"
  network_type      = "VPC"
  env_type          = "dev"
  host              = alicloud_db_instance.default.connection_string
  port              = 3306
  database_user     = alicloud_db_account.default.account_name
  database_password = alicloud_db_account.default.account_password
  instance_name     = var.name
  dba_uid           = data.alicloud_account.current.id
  # The value of safe_rule can be queried through the interface: https://www.alibabacloud.com/help/en/dms/developer-reference/api-dms-enterprise-2018-11-01-liststandardgroups
  safe_rule      = "904496"
  use_dsql       = 1
  query_timeout  = 60
  export_timeout = 600
  ecs_region     = data.alicloud_regions.default.regions.0.id
}
```
## Argument Reference

The following arguments are supported:

* `tid` - (Optional) The tenant ID. 
* `instance_type` - (Required) Database type. Valid values: `MySQL`, `SQLServer`, `PostgreSQL`, `Oracle,` `DRDS`, `OceanBase`, `Mongo`, `Redis`.
* `instance_source` - (Required) The source of the database instance. Valid values: `PUBLIC_OWN`, `RDS`, `ECS_OWN`, `VPC_IDC`.
* `network_type` - (Required, ForceNew) Network type. Valid values: `CLASSIC`, `VPC`.
* `env_type` - (Required) Environment type. Valid values: `product` production environment, `dev` development environment, `pre` pre-release environment, `test` test environment, `sit` SIT environment, `uat` UAT environment, `pet` pressure test environment, `stag` STAG environment.
* `host` - (Required, ForceNew) Host address of the target database.
* `port` - (Required, ForceNew) Access port of the target database.
* `database_user` - (Required) Database access account.
* `database_password` - (Required) Database access password.
* `instance_alias` - It has been deprecated from provider version 1.100.0 and 'instance_name' instead.
* `instance_name` - (Optional, Available since v1.100.0) Instance name, to help users quickly distinguish positioning.
* `dba_uid` - (Required, ForceNew)  The DBA of the instance is passed into the Alibaba Cloud uid of the DBA.
* `safe_rule` - (Required, ForceNew) The security rule of the instance is passed into the name of the security rule in the enterprise.
* `query_timeout` - (Required) Query timeout time, unit: s (seconds).
* `export_timeout` - (Required) Export timeout, unit: s (seconds).
* `ecs_instance_id` - (Optional) ECS instance ID. The value of InstanceSource is the ECS self-built library. This value must be passed.
* `vpc_id` - (Optional) VPC ID. This value must be passed when the value of InstanceSource is VPC dedicated line IDC.
* `ecs_region` - (Optional) The region where the instance is located. This value must be passed when the value of InstanceSource is RDS, ECS self-built library, and VPC dedicated line IDC.
* `sid` - (Optional) The SID. This value must be passed when InstanceType is PostgreSQL or Oracle.
* `data_link_name` - (Optional) Cross-database query datalink name.
* `ddl_online` - (Optional) Whether to use online services, currently only supports MySQL and PolarDB. Valid values: `0` Not used, `1` Native online DDL priority, `2` DMS lock-free table structure change priority.
* `use_dsql` - (Optional) Whether to enable cross-instance query. Valid values: `0` not open, `1` open.
* `instance_id` - (Optional) The instance id of the database instance. 
* `dba_id` - (Optional) The dba id of the database instance.
* `skip_test` - (Optional) Whether the instance ignores test connectivity. Valid values: `true`, `false`.
* `safe_rule_id` - (Optional) The safe rule id of the database instance.
* `instance_alias` - (Optional, Deprecated from v1.100.0) Field `instance_alias` has been deprecated from version 1.100.0. Use `instance_name` instead.
* `state` - (Deprecated) It has been deprecated from provider version 1.100.0 and 'status' instead.

                            
## Attributes Reference

The following attributes are exported:

* `id` - The id of the DMS enterprise instance and format as `<host>:<port>`.
* `dba_nick_name` - The instance dba nickname.
* `status` - The instance status.

## Timeouts

-> **NOTE:** Available since v1.100.0.

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when creating the DMS enterprise instance. 

## Import

DMS Enterprise can be imported using host and port, e.g.

```shell
$ terraform import alicloud_dms_enterprise_instance.example rm-uf648hgs7874xxxx.mysql.rds.aliyuncs.com:3306
```
