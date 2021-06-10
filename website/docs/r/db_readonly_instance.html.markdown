---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_db_readonly_instance"
sidebar_current: "docs-alicloud-resource-db-readonly-instance"
description: |-
  Provides an RDS readonly instance resource.
---

# alicloud\_db\_readonly\_instance

Provides an RDS readonly instance resource. 

## Example Usage

```terraform
variable "creation" {
  default = "Rds"
}

variable "name" {
  default = "dbInstancevpc"
}

data "alicloud_zones" "default" {
  available_resource_creation = var.creation
}

resource "alicloud_vpc" "default" {
  name       = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id            = alicloud_vpc.default.id
  cidr_block        = "172.16.0.0/24"
  zone_id           = data.alicloud_zones.default.zones[0].id
  vswitch_name      = var.name
}

resource "alicloud_db_instance" "default" {
  engine               = "MySQL"
  engine_version       = "5.6"
  instance_type        = "rds.mysql.t1.small"
  instance_storage     = "20"
  instance_charge_type = "Postpaid"
  instance_name        = var.name
  vswitch_id           = alicloud_vswitch.default.id
  security_ips         = ["10.168.1.12", "100.69.7.112"]
}

resource "alicloud_db_readonly_instance" "default" {
  master_db_instance_id = alicloud_db_instance.default.id
  zone_id               = alicloud_db_instance.default.zone_id
  engine_version        = alicloud_db_instance.default.engine_version
  instance_type         = alicloud_db_instance.default.instance_type
  instance_storage      = "30"
  instance_name         = "${var.name}ro"
  vswitch_id            = alicloud_vswitch.default.id
}
```

## Argument Reference

The following arguments are supported:

* `engine_version` - (Required, ForceNew) Database version. Value options can refer to the latest docs [CreateDBInstance](https://www.alibabacloud.com/help/doc-detail/26228.htm) `EngineVersion`.
* `master_db_instance_id` - (Required) ID of the master instance.
* `instance_type` - (Required) DB Instance type. For details, see [Instance type table](https://www.alibabacloud.com/help/doc-detail/26312.htm).
* `instance_storage` - (Required) User-defined DB instance storage space. Value range: [5, 2000] for MySQL/SQL Server HA dual node edition. Increase progressively at a rate of 5 GB. For details, see [Instance type table](https://www.alibabacloud.com/help/doc-detail/26312.htm).
* `instance_name` - (Optional) The name of DB instance. It a string of 2 to 256 characters.
* `resource_group_id` (Optional, Computed, resource_group_id, Modifiable in 1.115.0+) The ID of resource group which the DB read-only instance belongs.
* `parameters` - (Optional) Set of parameters needs to be set after DB instance was launched. Available parameters can refer to the latest docs [View database parameter templates](https://www.alibabacloud.com/help/doc-detail/26284.htm).
* `force_restart` - (Optional, Available in 1.121.0+) Set it to true to make some parameter efficient when modifying them. Default to false.
* `zone_id` - (Optional, ForceNew) The Zone to launch the DB instance.
* `vswitch_id` - (Optional, ForceNew) The virtual switch ID to launch DB instances in one VPC.
* `tags` - (Optional, Available in 1.68.0+) A mapping of tags to assign to the resource.
    - Key: It can be up to 64 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It cannot be a null string.
    - Value: It can be up to 128 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It can be a null string.
* `ssl_enabled` - (Optional, Available in v1.125.0+) Specifies whether to enable or disable SSL encryption. Valid values:
  - 1: enables SSL encryption
  - 0: disables SSL encryption
* `ca_type` - (Optional, Available in 1.125.0+) The type of the server certificate. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs. If you set the SSLEnabled parameter to 1, the default value of this parameter is aliyun. It is valid only when `ssl_enabled  = 1`. Value range:
  - aliyun: a cloud certificate
  - custom: a custom certificate
* `server_cert` - (Optional, Available in 1.125.0+) The content of the server certificate. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs. If you set the CAType parameter to custom, you must also specify this parameter. It is valid only when `ssl_enabled  = 1`.
* `server_key` - (Optional, Available in 1.125.0+) The private key of the server certificate. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs. If you set the CAType parameter to custom, you must also specify this parameter. It is valid only when `ssl_enabled  = 1`.
* `client_ca_enabled` - (Optional, Available in 1.125.0+) Specifies whether to enable the public key of the CA that issues client certificates. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs. It is valid only when `ssl_enabled  = 1`. Valid values:
  - 1: enables the public key
  - 0: disables the public key
* `client_ca_cert` - (Optional, Available in 1.125.0+) The public key of the CA that issues client certificates. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs. If you set the ClientCAEbabled parameter to 1, you must also specify this parameter. It is valid only when `ssl_enabled  = 1`.
* `client_crl_enabled` - (Optional, Available in 1.125.0+) Specifies whether to enable a certificate revocation list (CRL) that contains revoked client certificates. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs. In addition, this parameter is available only when the public key of the CA that issues client certificates is enabled. It is valid only when `ssl_enabled  = 1`. Valid values:
  - 1: enables the CRL
  - 0: disables the CRL
* `client_cert_revocation_list` - (Optional, Available in 1.125.0+) The CRL that contains revoked client certificates. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs. If you set the ClientCrlEnabled parameter to 1, you must also specify this parameter. It is valid only when `ssl_enabled  = 1`.
* `acl` - (Optional, Available in 1.125.0+) The method that is used to verify the identities of clients. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs. In addition, this parameter is available only when the public key of the CA that issues client certificates is enabled. It is valid only when `ssl_enabled  = 1`. Valid values:
  - cert
  - perfer
  - verify-ca
  - verify-full (supported only when the instance runs PostgreSQL 12 or later)
* `replication_acl` - (Optional, Available in 1.125.0+) The method that is used to verify the replication permission. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs. In addition, this parameter is available only when the public key of the CA that issues client certificates is enabled. It is valid only when `ssl_enabled  = 1`. Valid values:
  - cert
  - perfer
  - verify-ca
  - verify-full (supported only when the instance runs PostgreSQL 12 or later)
-> **NOTE:** Because of data backup and migration, change DB instance type and storage would cost 15~20 minutes. Please make full preparation before changing them.

## Attributes Reference

The following attributes are exported:

* `id` - The RDS instance ID.
* `engine` - Database type.
* `port` - RDS database connection port.
* `connection_string` - RDS database connection string.

### Timeouts

-> **NOTE:** Available in 1.52.1+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 mins) Used when creating the db instance (until it reaches the initial `Running` status). 
* `update` - (Defaults to 30 mins) Used when updating the db instance (until it reaches the initial `Running` status). 
* `delete` - (Defaults to 20 mins) Used when terminating the db instance. 

## Import

RDS readonly instance can be imported using the id, e.g.

```
$ terraform import alicloud_db_readonly_instance.example rm-abc12345678
```
