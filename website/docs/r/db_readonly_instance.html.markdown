---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_db_readonly_instance"
sidebar_current: "docs-alicloud-resource-db-readonly-instance"
description: |-
  Provides an RDS readonly instance resource.
---

# alicloud_db_readonly_instance

Provides an RDS readonly instance resource.
-> **NOTE:** Available since v1.52.1+.

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
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_zones.default.zones[0].id
  vswitch_name = var.name
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
* `parameters` - (Optional) Set of parameters needs to be set after DB instance was launched. Available parameters can refer to the latest docs [View database parameter templates](https://www.alibabacloud.com/help/doc-detail/26284.htm). See [`parameters`](#parameters) below.
* `force_restart` - (Optional, Available in 1.121.0+) Set it to true to make some parameter efficient when modifying them. Default to false.
* `zone_id` - (Optional, ForceNew) The Zone to launch the DB instance.
* `vswitch_id` - (Optional, ForceNew) The virtual switch ID to launch DB instances in one VPC.
* `deletion_protection` - (Optional, Available in 1.167.0+) The switch of delete protection. Valid values:
  - true: delete protect.
  - false: no delete protect.
* `tags` - (Optional, Available in 1.68.0+) A mapping of tags to assign to the resource.
    - Key: It can be up to 64 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It cannot be a null string.
    - Value: It can be up to 128 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It can be a null string.
* `upgrade_db_instance_kernel_version` - (Optional, Available in 1.128.0+) Whether to upgrade a minor version of the kernel. Valid values:
  - true: upgrade
  - false: not to upgrade
* `upgrade_time` - (Optional, Available in 1.128.0+) The method to update the minor engine version. Default value: Immediate. It is valid only when `upgrade_db_instance_kernel_version = true`. Valid values:
  - Immediate: The minor engine version is immediately updated.
  - MaintainTime: The minor engine version is updated during the maintenance window. For more information about how to change the maintenance window, see ModifyDBInstanceMaintainTime.
  - SpecifyTime: The minor engine version is updated at the point in time you specify.
* `switch_time` - (Optional, Available in 1.128.0+) The specific point in time when you want to perform the update. Specify the time in the ISO 8601 standard in the yyyy-MM-ddTHH:mm:ssZ format. It is valid only when `upgrade_db_instance_kernel_version = true`. The time must be in UTC.

-> **NOTE:** This parameter takes effect only when you set the UpgradeTime parameter to SpecifyTime.
* `target_minor_version` - (Optional, Available in 1.128.0+) The minor engine version to which you want to update the instance. If you do not specify this parameter, the instance is updated to the latest minor engine version. It is valid only when `upgrade_db_instance_kernel_version = true`. You must specify the minor engine version in one of the following formats:
  - PostgreSQL: rds_postgres_<Major engine version>00_<Minor engine version>. Example: rds_postgres_1200_20200830.
  - MySQL: <RDS edition>_<Minor engine version>. Examples: rds_20200229, xcluster_20200229, and xcluster80_20200229. The following RDS editions are supported:
    - rds: The instance runs RDS Basic or High-availability Edition.
    - xcluster: The instance runs MySQL 5.7 on RDS Enterprise Edition.
    - xcluster80: The instance runs MySQL 8.0 on RDS Enterprise Edition.
  - SQLServer: <Minor engine version>. Example: 15.0.4073.23.
  
-> **NOTE:** For more information about minor engine versions, see Release notes of minor AliPG versions, Release notes of minor AliSQL versions, and Release notes of minor engine versions of ApsaraDB RDS for SQL Server.
* `ssl_enabled` - (Optional, Available in v1.124.4+) Specifies whether to enable or disable SSL encryption. Valid values:
  - 1: enables SSL encryption
  - 0: disables SSL encryption
* `ca_type` - (Optional, Available in 1.124.4+) The type of the server certificate. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs. If you set the SSLEnabled parameter to 1, the default value of this parameter is aliyun. It is valid only when `ssl_enabled  = 1`. Value range:
  - aliyun: a cloud certificate
  - custom: a custom certificate
* `server_cert` - (Optional, Available in 1.124.4+) The content of the server certificate. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs. If you set the CAType parameter to custom, you must also specify this parameter. It is valid only when `ssl_enabled  = 1`.
* `server_key` - (Optional, Available in 1.124.4+) The private key of the server certificate. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs. If you set the CAType parameter to custom, you must also specify this parameter. It is valid only when `ssl_enabled  = 1`.
* `client_ca_enabled` - (Optional, Available in 1.124.4+) Specifies whether to enable the public key of the CA that issues client certificates. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs. It is valid only when `ssl_enabled  = 1`. Valid values:
  - 1: enables the public key
  - 0: disables the public key
* `client_ca_cert` - (Optional, Available in 1.124.4+) The public key of the CA that issues client certificates. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs. If you set the ClientCAEbabled parameter to 1, you must also specify this parameter. It is valid only when `ssl_enabled  = 1`.
* `client_crl_enabled` - (Optional, Available in 1.124.4+) Specifies whether to enable a certificate revocation list (CRL) that contains revoked client certificates. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs. In addition, this parameter is available only when the public key of the CA that issues client certificates is enabled. It is valid only when `ssl_enabled  = 1`. Valid values:
  - 1: enables the CRL
  - 0: disables the CRL
* `client_cert_revocation_list` - (Optional, Available in 1.124.4+) The CRL that contains revoked client certificates. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs. If you set the ClientCrlEnabled parameter to 1, you must also specify this parameter. It is valid only when `ssl_enabled  = 1`.
* `acl` - (Optional, Available in 1.124.4+) The method that is used to verify the identities of clients. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs. In addition, this parameter is available only when the public key of the CA that issues client certificates is enabled. It is valid only when `ssl_enabled  = 1`. Valid values:
  - cert
  - perfer
  - verify-ca
  - verify-full (supported only when the instance runs PostgreSQL 12 or later)
* `replication_acl` - (Optional, Available in 1.124.4+) The method that is used to verify the replication permission. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs. In addition, this parameter is available only when the public key of the CA that issues client certificates is enabled. It is valid only when `ssl_enabled  = 1`. Valid values:
  - cert
  - perfer
  - verify-ca
  - verify-full (supported only when the instance runs PostgreSQL 12 or later)
-> **NOTE:** Because of data backup and migration, change DB instance type and storage would cost 15~20 minutes. Please make full preparation before changing them.

* `security_ips` - (Optional, Available in 1.196.0+) List of IP addresses allowed to access all databases of an instance. The list contains up to 1,000 IP addresses, separated by commas. Supported formats include 0.0.0.0/0, 10.23.12.24 (IP), and 10.23.12.24/24 (Classless Inter-Domain Routing (CIDR) mode. /24 represents the length of the prefix in an IP address. The range of the prefix length is [1,32]).
* `db_instance_ip_array_name` - (Optional, Available in 1.196.0+) The name of the IP address whitelist. Default value: Default.

-> **NOTE:** A maximum of 200 IP address whitelists can be configured for each instance.
* `db_instance_ip_array_attribute` - (Optional, Available in 1.196.0+) The attribute of the IP address whitelist. By default, this parameter is empty.

-> **NOTE:** The IP address whitelists that have the hidden attribute are not displayed in the ApsaraDB RDS console. These IP address whitelists are used to access Alibaba Cloud services, such as Data Transmission Service (DTS).
* `security_ip_type` - (Optional, Available in 1.196.0+) The type of IP address in the IP address whitelist.
* `whitelist_network_type` - (Optional, Available in 1.196.0+) The network type of the IP address whitelist. Default value: MIX. Valid values:
  - Classic: classic network in enhanced whitelist mode
  - VPC: virtual private cloud (VPC) in enhanced whitelist mode
  - MIX: standard whitelist mode
    -> **NOTE:** In standard whitelist mode, IP addresses and CIDR blocks can be added only to the default IP address whitelist. In enhanced whitelist mode, IP addresses and CIDR blocks can be added to both IP address whitelists of the classic network type and those of the VPC network type.
* `modify_mode` - (Optional, Available in 1.196.0+) The method that is used to modify the IP address whitelist. Default value: Cover. Valid values:
  - Cover: Use the value of the SecurityIps parameter to overwrite the existing entries in the IP address whitelist.
  - Append: Add the IP addresses and CIDR blocks that are specified in the SecurityIps parameter to the IP address whitelist.
  - Delete: Delete IP addresses and CIDR blocks that are specified in the SecurityIps parameter from the IP address whitelist. You must retain at least one IP address or CIDR block.
* `instance_charge_type` - (Optional, Available in 1.201.0+) Valid values are `Prepaid`, `Postpaid`, Default to `Postpaid`. The interval between the two conversion operations must be greater than 15 minutes. Only when this parameter is `Postpaid`, the instance can be released. 
* `period` - (Optional, Available in 1.201.0+) The duration that you will buy DB instance (in month). It is valid when instance_charge_type is `PrePaid`. Valid values: [1~9], 12, 24, 36.
* `auto_renew` - (Optional, Available in 1.201.0+) Whether to renewal a DB instance automatically or not. It is valid when instance_charge_type is `PrePaid`. Default to `false`.
* `auto_renew_period` - (Optional, Available in 1.201.0+) Auto-renewal period of an instance, in the unit of the month. It is valid when instance_charge_type is `PrePaid`. Valid value:[1~12], Default to 1.
* `db_instance_storage_type` - (Optional, Available in 1.201.0+) The storage type of the instance. Valid values:
  - local_ssd: specifies to use local SSDs. This value is recommended.
  - cloud_ssd: specifies to use standard SSDs.
  - cloud_essd: specifies to use enhanced SSDs (ESSDs).
  - cloud_essd2: specifies to use enhanced SSDs (ESSDs).
  - cloud_essd3: specifies to use enhanced SSDs (ESSDs).
* `effective_time` - (Optional, Available in 1.207.2+) The method to change.  Default value: Immediate. Valid values:
  - Immediate: The change immediately takes effect.
  - MaintainTime: The change takes effect during the specified maintenance window. For more information, see ModifyDBInstanceMaintainTime.

### `parameters`

The parameters mapping supports the following:

* `name` - (Required) The parameter name.
* `value` - (Required) The parameter value.

## Attributes Reference

The following attributes are exported:

* `id` - The RDS instance ID.
* `engine` - Database type.
* `port` - RDS database connection port.
* `connection_string` - RDS database connection string.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 mins) Used when creating the db instance (until it reaches the initial `Running` status). 
* `update` - (Defaults to 30 mins) Used when updating the db instance (until it reaches the initial `Running` status). 
* `delete` - (Defaults to 20 mins) Used when terminating the db instance. 

## Import

RDS readonly instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_db_readonly_instance.example rm-abc12345678
```
