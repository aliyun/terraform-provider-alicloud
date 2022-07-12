---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_db_instance"
sidebar_current: "docs-alicloud-resource-db-instance"
description: |-
  Provides an RDS instance resource.
---

# alicloud\_db\_instance

Provides an RDS instance resource. A DB instance is an isolated database environment in the cloud. A DB instance can contain multiple user-created databases.

For information about RDS and how to use it, see [What is ApsaraDB for RDS](https://www.alibabacloud.com/help/en/doc-detail/26092.htm).

-> **NOTE:** This resource has a fatal bug in the version v1.155.0. If you want to use new feature, please upgrade it to v1.156.0.

## Example Usage

### Create a RDS MySQL instance

```terraform
variable "name" {
  default = "tf-testaccdbinstance"
}

variable "creation" {
  default = "Rds"
}

data "alicloud_zones" "example" {
  available_resource_creation = var.creation
}

resource "alicloud_vpc" "example" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "example" {
  vpc_id       = alicloud_vpc.example.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_zones.example.zones[0].id
  vswitch_name = var.name
}

resource "alicloud_db_instance" "example" {
  engine               = "MySQL"
  engine_version       = "5.6"
  instance_type        = "rds.mysql.s2.large"
  instance_storage     = "30"
  instance_charge_type = "Postpaid"
  instance_name        = var.name
  vswitch_id           = alicloud_vswitch.example.id
  monitoring_period    = "60"
}
```

### Create a RDS MySQL instance with specific parameters

```terraform
resource "alicloud_vpc" "example" {
  vpc_name   = "vpc-123456"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "example" {
  vpc_id       = alicloud_vpc.example.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_zones.example.zones.0.id
  vswitch_name = "vpc-123456"
}

resource "alicloud_db_instance" "default" {
  engine              = "MySQL"
  engine_version      = "5.6"
  db_instance_class   = "rds.mysql.t1.small"
  db_instance_storage = "10"
  vswitch_id          = alicloud_vswitch.example.id
}

resource "alicloud_db_instance" "example" {
  engine              = "MySQL"
  engine_version      = "5.6"
  db_instance_class   = "rds.mysql.t1.small"
  db_instance_storage = "10"
  parameters {
    name  = "innodb_large_prefix"
    value = "ON"
  }
  parameters {
    name  = "connect_timeout"
    value = "50"
  }
}
```
### Create a High Availability RDS MySQL Instance

```terraform
variable "name" {
  default = "tf-testaccdbinstance"
}

data "alicloud_zones" "example" {
  available_resource_creation = "Rds"
}

resource "alicloud_vpc" "example" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "example" {
  count        = 2
  vpc_id       = alicloud_vpc.example.id
  cidr_block   = format("172.16.%d.0/24", count.index + 1)
  zone_id      = data.alicloud_zones.example.zones[count.index].id
  vswitch_name = format("vswich_%d", var.name, count.index)
}

resource "alicloud_db_instance" "example" {
  engine               = "MySQL"
  engine_version       = "5.6"
  instance_storage     = "30"
  instance_type        = "rds.mysql.t1.small"
  instance_charge_type = "Postpaid"
  instance_name        = var.name
  zone_id              = data.alicloud_zones.example.zones.0.id
  zone_id_slave_a      = data.alicloud_zones.example.zones.1.id
  vswitch_id           = join(",", alicloud_vswitch.example.*.id)
  monitoring_period    = "60"
}

```

### Create a High Availability RDS MySQL Instance with multi zones

```terraform
variable "name" {
  default = "tf-testaccdbinstance"
}

data "alicloud_zones" "example" {
  available_resource_creation = "Rds"
  multi                       = true
}

resource "alicloud_vpc" "example" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "example" {
  count        = length(data.alicloud_zones.example.zones.0.multi_zone_ids)
  vpc_id       = alicloud_vpc.example.id
  cidr_block   = format("172.16.%d.0/24", count.index + 1)
  zone_id      = data.alicloud_zones.example.zones.0.multi_zone_ids[count.index]
  vswitch_name = format("vswitch_%d", count.index)
}

resource "alicloud_db_instance" "this" {
  engine               = "MySQL"
  engine_version       = "5.6"
  instance_storage     = "30"
  instance_type        = "rds.mysql.t1.small"
  instance_charge_type = "Postpaid"
  instance_name        = var.name
  zone_id              = data.alicloud_zones.example.zones.0.id
  vswitch_id           = join(",", alicloud_vswitch.example.*.id)
  monitoring_period    = "60"
}
```

### Create a Enterprise Edition RDS MySQL Instance 
```terraform
variable "name" {
  default = "tf-testaccdbinstance"
}

data "alicloud_zones" "example" {
  available_resource_creation = "Rds"
}

resource "alicloud_vpc" "example" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "example" {
  count        = 3
  vpc_id       = alicloud_vpc.example.id
  cidr_block   = format("172.16.%d.0/24", count.index + 1)
  zone_id      = data.alicloud_zones.example.zones[count.index].id
  vswitch_name = format("vswich_%d", var.name, count.index)
}

resource "alicloud_db_instance" "example" {
  engine               = "MySQL"
  engine_version       = "8.0"
  instance_storage     = "30"
  instance_type        = "mysql.n2.small.25"
  instance_charge_type = "Postpaid"
  instance_name        = var.name
  zone_id              = data.alicloud_zones.example.zones.0.id
  zone_id_slave_a      = data.alicloud_zones.example.zones.1.id
  zone_id_slave_b      = data.alicloud_zones.example.zones.2.id
  vswitch_id           = join(",", alicloud_vswitch.example.*.id)
  monitoring_period    = "60"
}
```

### Deleting `alicloud_db_instance` or removing it from your configuration

The `alicloud_db_instance` resource allows you to manage `instance_charge_type = "Prepaid"` db instance, but Terraform cannot destroy it.
Deleting the subscription resource or removing it from your configuration will remove it from your statefile and management, but will not destroy the DB Instance.
You can resume managing the subscription db instance via the AlibabaCloud Console.

## Argument Reference

The following arguments are supported:

* `engine` - (Required,ForceNew) Database type. Value options: MySQL, SQLServer, PostgreSQL, and PPAS.
* `engine_version` - (Required,ForceNew) Database version. Value options can refer to the latest docs [CreateDBInstance](https://www.alibabacloud.com/help/doc-detail/26228.htm) `EngineVersion`.
* `instance_type` - (Required) DB Instance type. For details, see [Instance type table](https://www.alibabacloud.com/help/doc-detail/26312.htm).
* `instance_storage` - (Required) User-defined DB instance storage space. Value range:
    - [5, 2000] for MySQL/PostgreSQL/PPAS HA dual node edition;
    - [20,1000] for MySQL 5.7 basic single node edition;
    - [10, 2000] for SQL Server 2008R2;
    - [20,2000] for SQL Server 2012 basic single node edition
    Increase progressively at a rate of 5 GB. For details, see [Instance type table](https://www.alibabacloud.com/help/doc-detail/26312.htm).
    Note: There is extra 5 GB storage for SQL Server Instance and it is not in specified `instance_storage`.

* `db_instance_storage_type` - (Optional, Available in 1.68.0+) The storage type of the instance. Valid values:
    - local_ssd: specifies to use local SSDs. This value is recommended.
    - cloud_ssd: specifies to use standard SSDs.
    - cloud_essd: specifies to use enhanced SSDs (ESSDs).
    - cloud_essd2: specifies to use enhanced SSDs (ESSDs).
    - cloud_essd3: specifies to use enhanced SSDs (ESSDs).

* `db_time_zone` - (Optional, Available in 1.136.0+) The time zone of the instance. This parameter takes effect only when you set the `Engine` parameter to MySQL or PostgreSQL.
  - If you set the `Engine` parameter to MySQL.
    - This time zone of the instance is in UTC. Valid values: -12:59 to +13:00.
    - You can specify this parameter when the instance is equipped with local SSDs. For example, you can specify the time zone to Asia/Hong_Kong. For more information about time zones, see [Time zones](https://www.alibabacloud.com/help/doc-detail/297356.htm).
  - If you set the `Engine` parameter to PostgreSQL.
    - This time zone of the instance is not in UTC. For more information about time zones, see [Time zones](https://www.alibabacloud.com/help/doc-detail/297356.htm).
    - You can specify this parameter only when the instance is equipped with standard SSDs or ESSDs.

-> **NOTE:**
- You can specify the time zone when you create a primary instance. You cannot specify the time zone when you create a read-only instance. Read-only instances inherit the time zone of their primary instance.
- If you do not specify this parameter, the system assigns the default time zone of the region where the instance resides.
* `sql_collector_status` - (Optional, Available in 1.70.0+) The sql collector status of the instance. Valid values are `Enabled`, `Disabled`, Default to `Disabled`.
* `sql_collector_config_value` - (Optional, Available in 1.70.0+) The sql collector keep time of the instance. Valid values are `30`, `180`, `365`, `1095`, `1825`, Default to `30`.
    
* `instance_name` - (Optional) The name of DB instance. It a string of 2 to 256 characters.
* `connection_string_prefix` - (Optional, Available in 1.126.0+) The private connection string prefix. If you want to update public connection string prefix, please use resource alicloud_db_connection [connection_prefix](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/db_connection#connection_prefix). 
-> **NOTE:** The prefix must be 8 to 64 characters in length and can contain letters, digits, and hyphens (-). It cannot contain Chinese characters and special characters ~!#%^&*=+\|{};:'",<>/?
* `port` - (Optional, Available in 1.126.0+) The private port of the database service. If you want to update public port, please use resource alicloud_db_connection [port](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/db_connection#port).
* `instance_charge_type` - (Optional) Valid values are `Prepaid`, `Postpaid`, Default to `Postpaid`. Currently, the resource only supports PostPaid to PrePaid.
* `resource_group_id` (Optional, Computed, Available in 1.86.0+, Modifiable in 1.115.0+) The ID of resource group which the DB instance belongs.
* `period` - (Optional) The duration that you will buy DB instance (in month). It is valid when instance_charge_type is `PrePaid`. Valid values: [1~9], 12, 24, 36.
-> **NOTE:** The attribute `period` is only used to create Subscription instance or modify the PayAsYouGo instance to Subscription. Once effect, it will not be modified that means running `terraform apply` will not effect the resource.
* `monitoring_period` - (Optional) The monitoring frequency in seconds. Valid values are 5, 60, 300. Defaults to 300. 
* `auto_renew` - (Optional, Available in 1.34.0+) Whether to renewal a DB instance automatically or not. It is valid when instance_charge_type is `PrePaid`. Default to `false`.
* `auto_renew_period` - (Optional, Available in 1.34.0+) Auto-renewal period of an instance, in the unit of the month. It is valid when instance_charge_type is `PrePaid`. Valid value:[1~12], Default to 1.
* `zone_id` - (ForceNew) The Zone to launch the DB instance. From version 1.8.1, it supports multiple zone.
If it is a multi-zone and `vswitch_id` is specified, the vswitch must in the one of them.
The multiple zone ID can be retrieved by setting `multi` to "true" in the data source `alicloud_zones`.
* `vswitch_id` - (ForceNew) The virtual switch ID to launch DB instances in one VPC. If there are multiple vswitches, separate them with commas.
* `private_ip_address` - (Optional, Available in v1.125.0+) The private IP address of the instance. The private IP address must be within the Classless Inter-Domain Routing (CIDR) block of the vSwitch that is specified by the VSwitchId parameter.
* `security_ips` - (Optional) List of IP addresses allowed to access all databases of an instance. The list contains up to 1,000 IP addresses, separated by commas. Supported formats include 0.0.0.0/0, 10.23.12.24 (IP), and 10.23.12.24/24 (Classless Inter-Domain Routing (CIDR) mode. /24 represents the length of the prefix in an IP address. The range of the prefix length is [1,32]).
* `db_instance_ip_array_name` - (Optional, Available in 1.125.0+) The name of the IP address whitelist. Default value: Default.

-> **NOTE:** A maximum of 200 IP address whitelists can be configured for each instance.
* `db_instance_ip_array_attribute` - (Optional, Available in 1.125.0+) The attribute of the IP address whitelist. By default, this parameter is empty.

-> **NOTE:** The IP address whitelists that have the hidden attribute are not displayed in the ApsaraDB RDS console. These IP address whitelists are used to access Alibaba Cloud services, such as Data Transmission Service (DTS).
* `security_ip_type` - (Optional, Available in 1.125.0+) The type of IP address in the IP address whitelist.
* `db_is_ignore_case` - (Optional, Available in 1.168.0+) Specifies whether table names on the instance are case-sensitive. Valid values: `true`, `false`.
  * `true` - Table names are not case-sensitive. This is the default value.
  * `false` - Table names are case-sensitive.
* `whitelist_network_type` - (Optional, Available in 1.125.0+) The network type of the IP address whitelist. Default value: MIX. Valid values:
    - Classic: classic network in enhanced whitelist mode
    - VPC: virtual private cloud (VPC) in enhanced whitelist mode
    - MIX: standard whitelist mode

-> **NOTE:** In standard whitelist mode, IP addresses and CIDR blocks can be added only to the default IP address whitelist. In enhanced whitelist mode, IP addresses and CIDR blocks can be added to both IP address whitelists of the classic network type and those of the VPC network type.
* `modify_mode` - (Optional, Available in 1.125.0+) The method that is used to modify the IP address whitelist. Default value: Cover. Valid values:
    - Cover: Use the value of the SecurityIps parameter to overwrite the existing entries in the IP address whitelist.
    - Append: Add the IP addresses and CIDR blocks that are specified in the SecurityIps parameter to the IP address whitelist.
    - Delete: Delete IP addresses and CIDR blocks that are specified in the SecurityIps parameter from the IP address whitelist. You must retain at least one IP address or CIDR block.
* `security_ip_mode` - (Optional, Available in 1.62.1+)  Valid values are `normal`, `safety`, Default to `normal`. support `safety` switch to high security access mode.
* `fresh_white_list_readins` - (Optional, Available in 1.148.0+) The read-only instances to which you want to synchronize the IP address whitelist.
  * If the instance is attached with a read-only instance, you can use this parameter to synchronize the IP address whitelist to the read-only instance. If the instance is attached with multiple read-only instances, the read-only instances must be separated by commas (,).
  * If the instance is not attached with a read-only instance, this parameter is empty.
* `parameters` - (Optional) Set of parameters needs to be set after DB instance was launched. Available parameters can refer to the latest docs [View database parameter templates](https://www.alibabacloud.com/help/doc-detail/26284.htm) .
* `force_restart` - (Optional, Available in 1.75.0+) Set it to true to make some parameter efficient when modifying them. Default to false.
* `tags` - (Optional) A mapping of tags to assign to the resource.
    - Key: It can be up to 64 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It cannot be a null string.
    - Value: It can be up to 128 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It can be a null string.

   Note: From 1.63.0, the tag key and value are case sensitive. Before that, they are not case sensitive.

* `security_group_id` - (Deprecated) It has been deprecated from 1.69.0 and use `security_group_ids` instead.
* `security_group_ids` - (Optional, List(string), Available in 1.69.0+) The list IDs to join ECS Security Group. At most supports three security groups.
* `maintain_time` - (Optional, Available in 1.56.0+) Maintainable time period format of the instance: HH:MMZ-HH:MMZ (UTC time)
* `auto_upgrade_minor_version` - (Optional, Available in 1.62.1+) The upgrade method to use. Valid values:
   - Auto: Instances are automatically upgraded to a higher minor version.
   - Manual: Instances are forcibly upgraded to a higher minor version when the current version is unpublished.
   
   Default to "Manual". See more [details and limitation](https://www.alibabacloud.com/help/doc-detail/123605.htm).
* `upgrade_db_instance_kernel_version` - (Optional, Available in 1.126.0+) Whether to upgrade a minor version of the kernel. Valid values:
    - true: upgrade
    - false: not to upgrade
* `upgrade_time` - (Optional, Available in 1.126.0+) The method to update the minor engine version. Default value: Immediate. It is valid only when `upgrade_db_instance_kernel_version = true`. Valid values:
    - Immediate: The minor engine version is immediately updated.
    - MaintainTime: The minor engine version is updated during the maintenance window. For more information about how to change the maintenance window, see ModifyDBInstanceMaintainTime.
    - SpecifyTime: The minor engine version is updated at the point in time you specify.
* `switch_time` - (Optional, Available in 1.126.0+) The specific point in time when you want to perform the update. Specify the time in the ISO 8601 standard in the yyyy-MM-ddTHH:mm:ssZ format. It is valid only when `upgrade_db_instance_kernel_version = true`. The time must be in UTC.

-> **NOTE:** This parameter takes effect only when you set the UpgradeTime parameter to SpecifyTime.
* `target_minor_version` - (Optional, Available in 1.126.0+) The minor engine version to which you want to update the instance. If you do not specify this parameter, the instance is updated to the latest minor engine version. It is valid only when `upgrade_db_instance_kernel_version = true`. You must specify the minor engine version in one of the following formats:
    - PostgreSQL: rds_postgres_<Major engine version>00_<Minor engine version>. Example: rds_postgres_1200_20200830.
    - MySQL: <RDS edition>_<Minor engine version>. Examples: rds_20200229, xcluster_20200229, and xcluster80_20200229. The following RDS editions are supported:
      - rds: The instance runs RDS Basic or High-availability Edition.
      - xcluster: The instance runs MySQL 5.7 on RDS Enterprise Edition.
      - xcluster80: The instance runs MySQL 8.0 on RDS Enterprise Edition.
    - SQLServer: <Minor engine version>. Example: 15.0.4073.23.

-> **NOTE:** For more information about minor engine versions, see Release notes of minor AliPG versions, Release notes of minor AliSQL versions, and Release notes of minor engine versions of ApsaraDB RDS for SQL Server.
* `zone_id_slave_a` - (Optional, ForceNew, Available in 1.101.0+) The region ID of the secondary instance if you create a secondary instance. If you set this parameter to the same value as the ZoneId parameter, the instance is deployed in a single zone. Otherwise, the instance is deployed in multiple zones.
* `zone_id_slave_b`- (Optional, ForceNew, Available in 1.101.0+) The region ID of the log instance if you create a log instance. If you set this parameter to the same value as the ZoneId parameter, the instance is deployed in a single zone. Otherwise, the instance is deployed in multiple zones.
* `ssl_action` - (Optional, Available in v1.90.0+) Actions performed on SSL functions, Valid values: `Open`: turn on SSL encryption; `Close`: turn off SSL encryption; `Update`: update SSL certificate. See more [engine and engineVersion limitation](https://www.alibabacloud.com/help/zh/doc-detail/26254.htm).
* `tde_status` - (Optional, ForceNew, Available in 1.90.0+) The TDE(Transparent Data Encryption) status. See more [engine and engineVersion limitation](https://www.alibabacloud.com/help/zh/doc-detail/26256.htm).
* `encryption_key` - (Optional, Available in 1.109.0+) The key id of the KMS. Used for encrypting a disk if not null. Only for PostgreSQL, MySQL and SQLServer.
* `ca_type` - (Optional, Available in 1.124.1+) The type of the server certificate. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs. If you set the SSLEnabled parameter to 1, the default value of this parameter is aliyun. Value range:
    - aliyun: a cloud certificate
    - custom: a custom certificate
* `server_cert` - (Optional, Available in 1.124.1+) The content of the server certificate. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs. If you set the CAType parameter to custom, you must also specify this parameter.
* `server_key` - (Optional, Available in 1.124.1+) The private key of the server certificate. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs. If you set the CAType parameter to custom, you must also specify this parameter.
* `client_ca_enabled` - (Optional, Available in 1.124.1+) Specifies whether to enable the public key of the CA that issues client certificates. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs. Valid values:
    - 1: enables the public key
    - 0: disables the public key
* `client_ca_cert` - (Optional, Available in 1.124.1+) The public key of the CA that issues client certificates. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs. If you set the ClientCAEbabled parameter to 1, you must also specify this parameter.
* `client_crl_enabled` - (Optional, Available in 1.124.1+) Specifies whether to enable a certificate revocation list (CRL) that contains revoked client certificates. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs. In addition, this parameter is available only when the public key of the CA that issues client certificates is enabled. Valid values:
    - 1: enables the CRL
    - 0: disables the CRL
* `client_cert_revocation_list` - (Optional, Available in 1.124.1+) The CRL that contains revoked client certificates. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs. If you set the ClientCrlEnabled parameter to 1, you must also specify this parameter.
* `acl` - (Optional, Available in 1.124.1+) The method that is used to verify the identities of clients. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs. In addition, this parameter is available only when the public key of the CA that issues client certificates is enabled. Valid values:
    - cert
    - perfer
    - verify-ca
    - verify-full (supported only when the instance runs PostgreSQL 12 or later)
* `replication_acl` - (Optional, Available in 1.124.1+) The method that is used to verify the replication permission. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs. In addition, this parameter is available only when the public key of the CA that issues client certificates is enabled. Valid values:
    - cert
    - perfer
    - verify-ca
    - verify-full (supported only when the instance runs PostgreSQL 12 or later)
* `ha_config` - (Optional, Available in 1.128.0+) The primary/secondary switchover mode of the instance. Default value: Auto. Valid values:
    - Auto: The system automatically switches over services from the primary to secondary instances in the event of a fault.
    - Manual: You must manually switch over services from the primary to secondary instances in the event of a fault.

-> **NOTE:** If you set this parameter to Manual, you must specify the ManualHATime parameter.
* `manual_ha_time` - (Optional, Available in 1.128.0+) The time after when you want to enable automatic primary/secondary switchover. At most, you can set this parameter to 23:59:59 seven days later. Specify the time in the ISO 8601 standard in the yyyy-MM-ddTHH:mm:ssZ format. The time must be in UTC.

-> **NOTE:** This parameter only takes effect when the HAConfig parameter is set to Manual.
* `released_keep_policy` - (Optional, Available in 1.136.0+) The policy based on which ApsaraDB RDS retains archived backup files after the instance is released. Valid values:
  - None: No archived backup files are retained.
  - Lastest: Only the last archived backup file is retained.
  - All: All the archived backup files are retained.

-> **NOTE:** This parameter is supported only when the instance runs the MySQL database engine.
* `storage_auto_scale` - (Optional, Available in 1.129.0+)Automatic storage space expansion switch. Valid values:
  - Enable
  - Disable

-> **NOTE:** This parameter only takes effect when the StorageAutoScale parameter is set to Enable.
* `storage_threshold` - (Optional, Available in 1.129.0+)The trigger threshold (percentage) for automatic storage space expansion. Valid values:
  - 10
  - 20
  - 30
  - 40
  - 50

-> **NOTE:** This parameter only takes effect when the StorageAutoScale parameter is set to Enable. The value must be greater than or equal to the total size of the current storage space of the instance.
* `storage_upper_bound` - (Optional, Available in 1.129.0+) The upper limit of the total storage space for automatic expansion of the storage space, that is, automatic expansion will not cause the total storage space of the instance to exceed this value. Unit: GB. The value must be ≥0.

-> **NOTE:** Because of data backup and migration, change DB instance type and storage would cost 15~20 minutes. Please make full preparation before changing them.
* `deletion_protection` - (Optional, Available in 1.165.0+) The switch of delete protection. Valid values: 
  - true: delete protect.
  - false: no delete protect.

-> **NOTE:** `deletion_protection` is valid only when attribute `instance_charge_type` is set to `Postpaid`, supported engine type: **MySQL**, **PostgresSQL**, **MariaDB**, **MSSQL**.
* `tcp_connection_type` - (Optional, Available in 1.171.0+) The availability check method of the instance. Valid values:
  - **SHORT**: Alibaba Cloud uses short-lived connections to check the availability of the instance.
  - **LONG**: Alibaba Cloud uses persistent connections to check the availability of the instance.

-> **NOTE:** `zone_id_slave_a` and `zone_id_slave_b` can specify slave zone ids when creating the high-availability or enterprise edition instances. Meanwhile, `vswitch_id` needs to pass in the corresponding vswitch id to the slave zone by order (If the `vswitch_id` is not specified, the classic network version will be created). For example, `zone_id` = "zone-a" and `zone_id_slave_a` = "zone-c", `zone_id_slave_b` = "zone-b", then the `vswitch_id` must be "vsw-zone-a,vsw-zone-c,vsw-zone-b". Of course, you can also choose automatic allocation , for example, `zone_id` = "zone-a" and `zone_id_slave_a` = "Auto",`zone_id_slave_b` = "Auto", then the `vswitch_id` must be "vsw-zone-a,Auto,Auto". The list contains up to 2 slave zone ids , separated by commas.
* `pg_hba_conf` - (Optional, Available in 1.155.0+) The configuration of [AD domain](https://www.alibabacloud.com/help/en/doc-detail/349288.htm) (documented below).
* `babelfish_port` - (Optional, Available in 1.176.0+) The TDS port of the instance for which Babelfish is enabled.
  
-> **NOTE:** This parameter applies only to ApsaraDB RDS for PostgreSQL instances. For more information about Babelfish for ApsaraDB RDS for PostgreSQL, see [Introduction to Babelfish](https://www.alibabacloud.com/help/en/apsaradb-for-rds/latest/babelfish-for-pg).
* `babelfish_config` - (Optional, Available in 1.176.0+) The configuration of an ApsaraDB RDS for PostgreSQL instance for which Babelfish is enabled. (documented below).

-> **NOTE:** This parameter takes effect only when you create an ApsaraDB RDS for PostgreSQL instance. For more information, see [Introduction to Babelfish](https://www.alibabacloud.com/help/en/apsaradb-for-rds/latest/babelfish-for-pg).

#### Block pg_hba_conf

The pg_hba_conf mapping supports the following:

* `type` - (Required) The type of connection to the instance. Valid values:
  * **host**: specifies to verify TCP/IP connections, including SSL connections and non-SSL connections.
  * **hostssl**: specifies to verify only TCP/IP connections that are established over SSL connections.
  * **hostnossl**: specifies to verify only TCP/IP connections that are established over non-SSL connections.

-> **NOTE:** You can set this parameter to hostssl only when SSL encryption is enabled for the instance. For more information, see [Configure SSL encryption for an ApsaraDB RDS for PostgreSQL instance](https://www.alibabacloud.com/help/en/doc-detail/229518.htm).
* `mask` - (Optional) The mask of the instance. If the value of the `Address` parameter is an IP address, you can use this parameter to specify the mask of the IP address.
* `database` - (Required) The name of the database that the specified users are allowed to access. If you set this parameter to all, the specified users are allowed to access all databases in the instance. If you specify multiple databases, separate the database names with commas (,).
* `priority_id` - (Required) The priority of an AD domain. If you set this parameter to 0, the AD domain has the highest priority. Valid values: 0 to 10000. This parameter is used to identify each AD domain. When you add an AD domain, the value of the PriorityId parameter of the new AD domain cannot be the same as the value of the PriorityId parameter for any existing AD domain. When you modify or delete an AD domain, you must also modify or delete the value of the PriorityId parameter for this AD domain.
* `address` - (Required) The IP addresses from which the specified users can access the specified databases. If you set this parameter to 0.0.0.0/0, the specified users are allowed to access the specified databases from all IP addresses.
* `user` - (Required) The user that is allowed to access the instance. If you specify multiple users, separate the usernames with commas (,).
* `method` - (Required) The authentication method of Lightweight Directory Access Protocol (LDAP). Valid values: `trust`, `reject`, `scram-sha-256`, `md5`, `password`, `gss`, `sspi`, `ldap`, `radius`, `cert`, `pam`.
* `option` - (Optional) Optional. The value of this parameter is based on the value of the HbaItem.N.Method parameter. In this topic, LDAP is used as an example. You must configure this parameter. For more information, see [Authentication Methods](https://www.postgresql.org/docs/11/auth-methods.html).

#### Block babelfish_config

The babelfish_config mapping supports the following:

* `babelfish_enabled` - (Required) specifies whether to enable the Babelfish for the instance. If you set this parameter to **true**, you enable Babelfish for the instance. If you leave this parameter empty, you disable Babelfish for the instance.
* `migration_mode` - (Required) The migration mode of the instance. Valid values: **single-db** and **multi-db**.
* `master_username` - (Required) The name of the administrator account. The name can contain lowercase letters, digits, and underscores (_). It must start with a letter and end with a letter or digit. It can be up to 63 characters in length and cannot start with pg.
* `master_user_password` - (Required) The password of the administrator account. The password must contain at least three of the following character types: uppercase letters, lowercase letters, digits, and special characters. It must be 8 to 32 characters in length. The password can contain any of the following characters:! @ # $ % ^ & * () _ + - =

## Attributes Reference

The following attributes are exported:

* `id` - The RDS instance ID.
* `connection_string` - RDS database connection string.
* `ssl_status` - Status of the SSL feature. `Yes`: SSL is turned on; `No`: SSL is turned off.

### Timeouts

-> **NOTE:** Available in 1.52.1+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 20 mins) Used when creating the db instance (until it reaches the initial `Running` status). 
* `update` - (Defaults to 30 mins) Used when updating the db instance (until it reaches the initial `Running` status). 
* `delete` - (Defaults to 20 mins) Used when terminating the db instance. 

## Import

RDS instance can be imported using the id, e.g.

```
$ terraform import alicloud_db_instance.example rm-abc12345678
```
