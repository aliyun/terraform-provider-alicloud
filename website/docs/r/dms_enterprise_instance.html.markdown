---
subcategory: "DMS Enterprise"
layout: "alicloud"
page_title: "Alicloud: alicloud_dms_enterprise_instance"
sidebar_current: "docs-alicloud-resource-dms-enterprise-instance"
description: |-
  Provides a DMS Enterprise Instance resource.
---

# alicloud\_dms\_enterprise\_instance

Provides a DMS Enterprise Instance resource.

-> **NOTE:** API users must first register in DMS.
-> **NOTE:** Available in 1.81.0+.

## Example Usage

```terraform
resource "alicloud_dms_enterprise_instance" "default" {
  tid               = "12345"
  instance_type     = "MySQL"
  instance_source   = "RDS"
  network_type      = "VPC"
  env_type          = "test"
  host              = "rm-uf648hgsxxxxxx.mysql.rds.aliyuncs.com"
  port              = 3306
  database_user     = "your_user_name"
  database_password = "Yourpassword123"
  instance_name     = "your_alias_name"
  dba_uid           = "1182725234xxxxxxx"
  safe_rule         = "自由操作"
  query_timeout     = 60
  export_timeout    = 600
  ecs_region        = "cn-shanghai"
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
* `instance_name` - (Optional, Computed, Available in v1.100.0+) Instance name, to help users quickly distinguish positioning.
* `dba_uid` - (Required, ForceNew)  The DBA of the instance is passed into the Alibaba Cloud uid of the DBA.
* `safe_rule` - (Required, ForceNew) The security rule of the instance is passed into the name of the security rule in the enterprise.
* `query_timeout` - (Required) Query timeout time, unit: s (seconds).
* `export_timeout` - (Required) Export timeout, unit: s (seconds).
* `ecs_instance_id` - (Optional, Computed) ECS instance ID. The value of InstanceSource is the ECS self-built library. This value must be passed.
* `vpc_id` - (Optional) VPC ID. This value must be passed when the value of InstanceSource is VPC dedicated line IDC.
* `ecs_region` - (Optional) The region where the instance is located. This value must be passed when the value of InstanceSource is RDS, ECS self-built library, and VPC dedicated line IDC.
* `sid` - (Optional) The SID. This value must be passed when InstanceType is PostgreSQL or Oracle.
* `data_link_name` - (Optional) Cross-database query datalink name.
* `ddl_online` - (Optional) Whether to use online services, currently only supports MySQL and PolarDB. Valid values: `0` Not used, `1` Native online DDL priority, `2` DMS lock-free table structure change priority.
* `use_dsql` - (Optional) Whether to enable cross-instance query. Valid values: `0` not open, `1` open.
* `instance_id` - (Optional, Computed) The instance id of the database instance. 
* `dba_id` - (Optional, Computed) The dba id of the database instance.
* `skip_test` - (Optional) Whether the instance ignores test connectivity. Valid values: `true`, `false`.
* `safe_rule_id` - (Optional, Computed) The safe rule id of the database instance.
* `instance_alias` - (Optional, Computed, Deprecated from v1.100.0) Field `instance_alias` has been deprecated from version 1.100.0. Use `instance_name` instead.
                            
## Attributes Reference

The following attributes are exported:

* `id` - The id of the DMS enterprise instance and format as `<host>:<port>`.
* `dba_nick_name` - The instance dba nickname.
* `state` - It has been deprecated from provider version 1.100.0 and 'status' instead.
* `status` - The instance status.

### Timeouts

-> **NOTE:** Available in 1.100.0+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when creating the DMS enterprise instance. 

## Import

DMS Enterprise can be imported using host and port, e.g.

```
$ terraform import alicloud_dms_enterprise_instance.example rm-uf648hgs7874xxxx.mysql.rds.aliyuncs.com:3306
```
