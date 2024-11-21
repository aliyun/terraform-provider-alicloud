---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_rds_ddr_instance"
sidebar_current: "docs-alicloud-resource-rds-ddr-instance"
description: |-
  Provide RDS remote disaster recovery instance resources.
---

# alicloud_rds_ddr_instance

Provide RDS remote disaster recovery instance resources. 

For information about RDS remote disaster recovery instance and how to use it, see [What is ApsaraDB for RDS Remote Disaster Recovery](https://www.alibabacloud.com/help/en/rds/developer-reference/api-rds-2014-08-15-createddrinstance).

-> **NOTE:** Available since v1.198.0.

## Example Usage

Because the generation time of the disaster recovery set is uncertain, the query backup set may not have a value, so the following examples may not be executed successfully in one run.

### Create an RDS instance based on the remote disaster recovery set

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_rds_ddr_instance&exampleId=a6fc5b56-b9bd-39a0-52d8-30c6c67349f849fe0fa4&activeTab=example&spm=docs.r.rds_ddr_instance.0.a6fc5b56b9&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "tf-example"
}

data "alicloud_db_zones" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  db_instance_storage_type = "local_ssd"
  category                 = "HighAvailability"
}

data "alicloud_db_instance_classes" "default" {
  zone_id                  = data.alicloud_db_zones.default.ids.0
  engine                   = "MySQL"
  engine_version           = "8.0"
  db_instance_storage_type = "local_ssd"
  category                 = "HighAvailability"
}
data "alicloud_rds_cross_regions" "regions" {
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_db_zones.default.ids.0
  vswitch_name = var.name
}

resource "alicloud_db_instance" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_type            = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
  instance_storage         = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
  instance_charge_type     = "Postpaid"
  category                 = "HighAvailability"
  instance_name            = var.name
  vswitch_id               = alicloud_vswitch.default.id
  db_instance_storage_type = "local_ssd"
  monitoring_period        = "60"
}

resource "alicloud_rds_instance_cross_backup_policy" "default" {
  instance_id         = alicloud_db_instance.default.id
  cross_backup_region = data.alicloud_rds_cross_regions.regions.ids.0
}

data "alicloud_rds_cross_region_backups" "backups" {
  db_instance_id = alicloud_rds_instance_cross_backup_policy.default.instance_id
  start_time     = formatdate("YYYY-MM-DD'T'hh:mm'Z'", timeadd(timestamp(), "-24h"))
  end_time       = timestamp()
}

resource "alicloud_rds_ddr_instance" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  db_instance_storage_type = "local_ssd"
  instance_type            = alicloud_db_instance.default.instance_type
  instance_storage         = alicloud_db_instance.default.instance_storage
  payment_type             = "PayAsYouGo"
  vswitch_id               = alicloud_vswitch.default.id
  instance_name            = var.name
  monitoring_period        = "60"
  restore_type             = "BackupSet"
  backup_set_id            = data.alicloud_rds_cross_region_backups.backups.ids.0
}
```

### Create RDS instance according to the recovery time point of remote disaster recovery

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_rds_ddr_instance&exampleId=f2acaba7-6c8c-2823-c1d8-1c4e5658901224039095&activeTab=example&spm=docs.r.rds_ddr_instance.1.f2acaba76c&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "tf-example"
}

data "alicloud_db_zones" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  db_instance_storage_type = "local_ssd"
  category                 = "HighAvailability"
}

data "alicloud_db_instance_classes" "default" {
  zone_id                  = data.alicloud_db_zones.default.ids.0
  engine                   = "MySQL"
  engine_version           = "8.0"
  db_instance_storage_type = "local_ssd"
  category                 = "HighAvailability"
}
data "alicloud_rds_cross_regions" "regions" {
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_db_zones.default.ids.0
  vswitch_name = var.name
}

resource "alicloud_db_instance" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_type            = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
  instance_storage         = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
  instance_charge_type     = "Postpaid"
  category                 = "HighAvailability"
  instance_name            = var.name
  vswitch_id               = alicloud_vswitch.default.id
  db_instance_storage_type = "local_ssd"
  monitoring_period        = "60"
}

resource "alicloud_rds_instance_cross_backup_policy" "default" {
  instance_id         = alicloud_db_instance.default.id
  cross_backup_region = data.alicloud_rds_cross_regions.regions.ids.0
}

data "alicloud_rds_cross_region_backups" "backups" {
  db_instance_id = alicloud_rds_instance_cross_backup_policy.default.instance_id
  start_time     = formatdate("YYYY-MM-DD'T'hh:mm'Z'", timeadd(timestamp(), "-24h"))
  end_time       = timestamp()
}

resource "alicloud_rds_ddr_instance" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  db_instance_storage_type = "local_ssd"
  instance_type            = alicloud_db_instance.default.instance_type
  instance_storage         = alicloud_db_instance.default.instance_storage
  payment_type             = "PayAsYouGo"
  vswitch_id               = alicloud_vswitch.default.id
  instance_name            = var.name
  monitoring_period        = "60"
  restore_type             = "BackupTime"
  restore_time             = data.alicloud_rds_cross_region_backups.backups.backups.0.recovery_end_time
  source_region            = data.alicloud_rds_cross_region_backups.backups.backups.0.restore_regions.0
  source_db_instance_name  = data.alicloud_rds_cross_region_backups.backups.db_instance_id
}
```

### Deleting `alicloud_rds_ddr_instance` or removing it from your configuration

The `alicloud_rds_ddr_instance` resource allows you to manage `payment_type = "Subscription"` db instance, but Terraform cannot destroy it.
Deleting the subscription resource or removing it from your configuration will remove it from your statefile and management, but will not destroy the DB Instance.
You can resume managing the subscription db instance via the AlibabaCloud Console.

## Argument Reference

The following arguments are supported:

* `engine` - (Required, ForceNew) Database type. Value options: MySQL, SQLServer.

-> **NOTE:** When the 'EngineVersion' changes, it can be used as the target database version for the large version upgrade of RDS for MySQL instance.
* `engine_version` - (Required) Database version. Value options can refer to the latest docs [CreateDdrInstance](https://www.alibabacloud.com/help/en/apsaradb-for-rds/latest/restore-data-to-a-new-instance-across-regions) `EngineVersion`.
* `instance_type` - (Required) DB Instance type.

-> **NOTE:** When `storage_auto_scale="Enable"`, do not perform `instance_storage` check. when `storage_auto_scale="Disable"`, if the instance itself `instance_storage`has changed. You need to manually revise the `instance_storage` in the template value.
* `instance_storage` - (Required) The storage capacity of the destination instance. Valid values: 5 to 2000. Unit: GB.

This value must be a multiple of 5 GB. For more information, see Primary ApsaraDB RDS instance types.

* `db_instance_storage_type` - (Optional) The storage type of the instance. Valid values:
    - local_ssd: specifies to use local SSDs. This value is recommended.
    - cloud_ssd: specifies to use standard SSDs.
    - cloud_essd: specifies to use enhanced SSDs (ESSDs).
    - cloud_essd2: specifies to use enhanced SSDs (ESSDs).
    - cloud_essd3: specifies to use enhanced SSDs (ESSDs).

-> **NOTE:** You can specify the time zone when you create a primary instance. You cannot specify the time zone when you create a read-only instance. Read-only instances inherit the time zone of their primary instance. If you do not specify this parameter, the system assigns the default time zone of the region where the instance resides.
* `sql_collector_status` - (Optional) The sql collector status of the instance. Valid values are `Enabled`, `Disabled`, Default to `Disabled`.
* `sql_collector_config_value` - (Optional) The sql collector keep time of the instance. Valid values are `30`, `180`, `365`, `1095`, `1825`, Default to `30`.
    
* `instance_name` - (Optional) The name of DB instance. It a string of 2 to 256 characters.
* `connection_string_prefix` - (Optional) The private connection string prefix. If you want to update public connection string prefix, please use resource alicloud_db_connection [connection_prefix](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/db_connection#connection_prefix). 
-> **NOTE:** The prefix must be 8 to 64 characters in length and can contain letters, digits, and hyphens (-). It cannot contain Chinese characters and special characters ~!#%^&*=+\|{};:'",<>/?
* `port` - (Optional) The private port of the database service. If you want to update public port, please use resource alicloud_db_connection [port](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/db_connection#port).
* `payment_type` - (Required) Valid values are `Subscription`, `PayAsYouGo`, Default to `PayAsYouGo`. 
* `resource_group_id` (Optional) The ID of resource group which the DB instance belongs.
* `period` - (Optional) The duration that you will buy DB instance (in month). It is valid when payment_type is `Subscription`. Valid values: [1~9], 12, 24, 36.
-> **NOTE:** The attribute `period` is only used to create Subscription instance or modify the PayAsYouGo instance to Subscription. Once effect, it will not be modified that means running `terraform apply` will not effect the resource.
* `monitoring_period` - (Optional) The monitoring frequency in seconds. Valid values are 5, 60, 300. Defaults to 300. 
* `auto_renew` - (Optional) Whether to renewal a DB instance automatically or not. It is valid when payment_type is `Subscription`. Default to `false`.
* `auto_renew_period` - (Optional) Auto-renewal period of an instance, in the unit of the month. It is valid when payment_type is `Subscription`. Valid value:[1~12], Default to 1.
* `zone_id` - (ForceNew, Optional) The Zone to launch the DB instance. It supports multiple zone.
If it is a multi-zone and `vswitch_id` is specified, the vswitch must in the one of them.
The multiple zone ID can be retrieved by setting `multi` to "true" in the data source `alicloud_zones`.
* `vswitch_id` - (ForceNew, Optional) The virtual switch ID to launch DB instances in one VPC. If there are multiple vswitches, separate them with commas.
* `private_ip_address` - (Optional) The private IP address of the instance. The private IP address must be within the Classless Inter-Domain Routing (CIDR) block of the vSwitch that is specified by the VSwitchId parameter.
* `security_ips` - (Optional) List of IP addresses allowed to access all databases of an instance. The list contains up to 1,000 IP addresses, separated by commas. Supported formats include 0.0.0.0/0, 10.23.12.24 (IP), and 10.23.12.24/24 (Classless Inter-Domain Routing (CIDR) mode. /24 represents the length of the prefix in an IP address. The range of the prefix length is [1,32]).
* `db_instance_ip_array_name` - (Optional) The name of the IP address whitelist. Default value: Default.

-> **NOTE:** A maximum of 200 IP address whitelists can be configured for each instance.
* `db_instance_ip_array_attribute` - (Optional) The attribute of the IP address whitelist. By default, this parameter is empty.

-> **NOTE:** The IP address whitelists that have the hidden attribute are not displayed in the ApsaraDB RDS console. These IP address whitelists are used to access Alibaba Cloud services, such as Data Transmission Service (DTS).
* `security_ip_type` - (Optional) The type of IP address in the IP address whitelist.
* `whitelist_network_type` - (Optional) The network type of the IP address whitelist. Default value: MIX. Valid values:
    - Classic: classic network in enhanced whitelist mode
    - VPC: virtual private cloud (VPC) in enhanced whitelist mode
    - MIX: standard whitelist mode

-> **NOTE:** In standard whitelist mode, IP addresses and CIDR blocks can be added only to the default IP address whitelist. In enhanced whitelist mode, IP addresses and CIDR blocks can be added to both IP address whitelists of the classic network type and those of the VPC network type.
* `modify_mode` - (Optional) The method that is used to modify the IP address whitelist. Default value: Cover. Valid values:
    - Cover: Use the value of the SecurityIps parameter to overwrite the existing entries in the IP address whitelist.
    - Append: Add the IP addresses and CIDR blocks that are specified in the SecurityIps parameter to the IP address whitelist.
    - Delete: Delete IP addresses and CIDR blocks that are specified in the SecurityIps parameter from the IP address whitelist. You must retain at least one IP address or CIDR block.
* `security_ip_mode` - (Optional)  Valid values are `normal`, `safety`, Default to `normal`. support `safety` switch to high security access mode.
* `fresh_white_list_readins` - (Optional) The read-only instances to which you want to synchronize the IP address whitelist.
  * If the instance is attached with a read-only instance, you can use this parameter to synchronize the IP address whitelist to the read-only instance. If the instance is attached with multiple read-only instances, the read-only instances must be separated by commas (,).
  * If the instance is not attached with a read-only instance, this parameter is empty.
* `parameters` - (Optional) Set of parameters needs to be set after DB instance was launched. Available parameters can refer to the latest docs [View database parameter templates](https://www.alibabacloud.com/help/doc-detail/26284.htm) . See [`parameters`](#parameters) below.
* `force_restart` - (Optional) Set it to true to make some parameter efficient when modifying them. Default to false.
* `tags` - (Optional) A mapping of tags to assign to the resource.
    - Key: It can be up to 64 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It cannot be a null string.
    - Value: It can be up to 128 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It can be a null string.
* `security_group_ids` - (Optional, List(string)) The list IDs to join ECS Security Group. At most supports three security groups.
* `maintain_time` - (Optional) Maintainable time period format of the instance: HH:MMZ-HH:MMZ (UTC time)
* `auto_upgrade_minor_version` - (Optional) The upgrade method to use. Valid values:
   - Auto: Instances are automatically upgraded to a higher minor version.
   - Manual: Instances are forcibly upgraded to a higher minor version when the current version is unpublished.
 
   See more [details and limitation](https://www.alibabacloud.com/help/doc-detail/123605.htm).
* `upgrade_db_instance_kernel_version` - (Optional) Whether to upgrade a minor version of the kernel. Valid values:
    - true: upgrade
    - false: not to upgrade
* `upgrade_time` - (Optional) The method to update the minor engine version. Default value: Immediate. It is valid only when `upgrade_db_instance_kernel_version = true`. Valid values:
    - Immediate: The minor engine version is immediately updated.
    - MaintainTime: The minor engine version is updated during the maintenance window. For more information about how to change the maintenance window, see ModifyDBInstanceMaintainTime.
    - SpecifyTime: The minor engine version is updated at the point in time you specify.
* `switch_time` - (Optional) The specific point in time when you want to perform the update. Specify the time in the ISO 8601 standard in the yyyy-MM-ddTHH:mm:ssZ format. It is valid only when `upgrade_db_instance_kernel_version = true`. The time must be in UTC.

-> **NOTE:** This parameter takes effect only when you set the UpgradeTime parameter to SpecifyTime.
* `target_minor_version` - (Optional) The minor engine version to which you want to update the instance. If you do not specify this parameter, the instance is updated to the latest minor engine version. It is valid only when `upgrade_db_instance_kernel_version = true`. You must specify the minor engine version in one of the following formats:
    - PostgreSQL: rds_postgres_<Major engine version>00_<Minor engine version>. Example: rds_postgres_1200_20200830.
    - MySQL: <RDS edition>_<Minor engine version>. Examples: rds_20200229, xcluster_20200229, and xcluster80_20200229. The following RDS editions are supported:
      - rds: The instance runs RDS Basic or High-availability Edition.
      - xcluster: The instance runs MySQL 5.7 on RDS Enterprise Edition.
      - xcluster80: The instance runs MySQL 8.0 on RDS Enterprise Edition.
    - SQLServer: <Minor engine version>. Example: 15.0.4073.23.

-> **NOTE:** For more information about minor engine versions, see Release notes of minor AliPG versions, Release notes of minor AliSQL versions, and Release notes of minor engine versions of ApsaraDB RDS for SQL Server.
* `ssl_action` - (Optional) Actions performed on SSL functions, Valid values: `Open`: turn on SSL encryption; `Close`: turn off SSL encryption; `Update`: update SSL certificate. See more [engine and engineVersion limitation](https://www.alibabacloud.com/help/zh/doc-detail/26254.htm).
* `tde_status` - (Optional, ForceNew) The TDE(Transparent Data Encryption) status. See more [engine and engineVersion limitation](https://www.alibabacloud.com/help/zh/doc-detail/26256.htm).
* `encryption_key` - (Optional) The key id of the KMS. Used for encrypting a disk if not null. Only for PostgreSQL, MySQL and SQLServer.
* `ca_type` - (Optional) The type of the server certificate. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs. If you set the SSLEnabled parameter to 1, the default value of this parameter is aliyun. Value range:
    - aliyun: a cloud certificate
    - custom: a custom certificate
* `server_cert` - (Optional) The content of the server certificate. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs. If you set the CAType parameter to custom, you must also specify this parameter.
* `server_key` - (Optional) The private key of the server certificate. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs. If you set the CAType parameter to custom, you must also specify this parameter.
* `client_ca_enabled` - (Optional) Specifies whether to enable the public key of the CA that issues client certificates. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs. Valid values:
    - 1: enables the public key
    - 0: disables the public key
* `client_ca_cert` - (Optional) The public key of the CA that issues client certificates. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs. If you set the ClientCAEbabled parameter to 1, you must also specify this parameter.
* `client_crl_enabled` - (Optional) Specifies whether to enable a certificate revocation list (CRL) that contains revoked client certificates. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs. In addition, this parameter is available only when the public key of the CA that issues client certificates is enabled. Valid values:
    - 1: enables the CRL
    - 0: disables the CRL
* `client_cert_revocation_list` - (Optional) The CRL that contains revoked client certificates. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs. If you set the ClientCrlEnabled parameter to 1, you must also specify this parameter.
* `acl` - (Optional) The method that is used to verify the identities of clients. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs. In addition, this parameter is available only when the public key of the CA that issues client certificates is enabled. Valid values:
    - cert
    - perfer
    - verify-ca
    - verify-full (supported only when the instance runs PostgreSQL 12 or later)
* `replication_acl` - (Optional) The method that is used to verify the replication permission. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs. In addition, this parameter is available only when the public key of the CA that issues client certificates is enabled. Valid values:
    - cert
    - perfer
    - verify-ca
    - verify-full (supported only when the instance runs PostgreSQL 12 or later)
* `ha_config` - (Optional) The primary/secondary switchover mode of the instance. Default value: Auto. Valid values:
    - Auto: The system automatically switches over services from the primary to secondary instances in the event of a fault.
    - Manual: You must manually switch over services from the primary to secondary instances in the event of a fault.

-> **NOTE:** If you set this parameter to Manual, you must specify the ManualHATime parameter.
* `manual_ha_time` - (Optional) The time after when you want to enable automatic primary/secondary switchover. At most, you can set this parameter to 23:59:59 seven days later. Specify the time in the ISO 8601 standard in the yyyy-MM-ddTHH:mm:ssZ format. The time must be in UTC.

-> **NOTE:** This parameter only takes effect when the HAConfig parameter is set to Manual.
* `released_keep_policy` - (Optional) The policy based on which ApsaraDB RDS retains archived backup files after the instance is released. Valid values:
  - None: No archived backup files are retained.
  - Lastest: Only the last archived backup file is retained.
  - All: All the archived backup files are retained.

-> **NOTE:** This parameter is supported only when the instance runs the MySQL database engine.
* `storage_auto_scale` - (Optional)Automatic storage space expansion switch. Valid values:
  - Enable
  - Disable

-> **NOTE:** This parameter only takes effect when the StorageAutoScale parameter is set to Enable.
* `storage_threshold` - (Optional)The trigger threshold (percentage) for automatic storage space expansion. Valid values:
  - 10
  - 20
  - 30
  - 40
  - 50

-> **NOTE:** This parameter only takes effect when the StorageAutoScale parameter is set to Enable. The value must be greater than or equal to the total size of the current storage space of the instance.
* `storage_upper_bound` - (Optional) The upper limit of the total storage space for automatic expansion of the storage space, that is, automatic expansion will not cause the total storage space of the instance to exceed this value. Unit: GB. The value must be ≥0.

-> **NOTE:** Because of data backup and migration, change DB instance type and storage would cost 15~20 minutes. Please make full preparation before changing them.
* `deletion_protection` - (Optional) The switch of delete protection. Valid values: 
  - true: delete protect.
  - false: no delete protect.

-> **NOTE:** `deletion_protection` is valid only when attribute `payment_type` is set to `PayAsYouGo`, supported engine type: **MySQL**, **PostgreSQL**, **MariaDB**, **MSSQL**.
* `tcp_connection_type` - (Optional) The availability check method of the instance. Valid values:
  - **SHORT**: Alibaba Cloud uses short-lived connections to check the availability of the instance.
  - **LONG**: Alibaba Cloud uses persistent connections to check the availability of the instance.
  
* `pg_hba_conf` - (Optional) The configuration of [AD domain](https://www.alibabacloud.com/help/en/doc-detail/349288.htm) . See [`pg_hba_conf`](#pg_hba_conf) below.
* `vpc_id` - (Optional) The VPC ID of the instance.

-> **NOTE:** This parameter applies only to ApsaraDB RDS for MySQL instances. For more information about Upgrade the major engine version of an ApsaraDB RDS for MySQL instance, see [Upgrade the major engine version of an RDS instance in the ApsaraDB RDS console](https://www.alibabacloud.com/help/en/apsaradb-for-rds/latest/upgrade-the-major-engine-version-of-an-apsaradb-rds-for-mysql-instance-1).
* `effective_time` - (Optional) The method to update the engine version and change.  Default value: Immediate. Valid values:
  - Immediate: The change immediately takes effect.
  - MaintainTime: The change takes effect during the specified maintenance window. For more information, see ModifyDBInstanceMaintainTime.

* `restore_type` - (Required, ForceNew) The method that is used to restore data. Valid values:
  - BackupSet: Data is restored from a backup set. If you use this value, you must also specify the BackupSetID parameter.
  - BackupTime: restores data to a point in time. You must also specify the RestoreTime, SourceRegion, and SourceDBInstanceName parameters.
* `backup_set_id` - (Optional, ForceNew) The ID of the backup set that is used for the restoration. You can call the DescribeCrossRegionBackups operation to query the ID of the backup set.
* `restore_time` - (Optional, ForceNew) The point in time to which you want to restore data. The point in time that you specify must be earlier than the current time. Specify the time in the ISO 8601 standard in the yyyy-MM-ddTHH:mm:ssZ format. The time must be in UTC.
* `source_region` - (Optional, ForceNew) The region ID of the source instance if you want to restore data to a point in time.
* `source_db_instance_name` - (Optional, ForceNew) The ID of the source instance if you want to restore data to a point in time.

### `parameters`

The parameters mapping supports the following:

* `name` - (Required) The parameter name.
* `value` - (Required) The parameter value.

### `pg_hba_conf`

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
## Attributes Reference

The following attributes are exported:

* `id` - The RDS instance ID.
* `connection_string` - RDS database connection string.
* `ssl_status` - Status of the SSL feature. `Yes`: SSL is turned on; `No`: SSL is turned off.
* `zone_id_slave_a` - The region ID of the secondary instance if you create a secondary instance. 
* `zone_id_slave_b`- The region ID of the log instance if you create a log instance. 
* `category` - The RDS edition of the instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 90 mins) Used when creating the db instance (until it reaches the initial `Running` status). 
* `update` - (Defaults to 30 mins) Used when updating the db instance (until it reaches the initial `Running` status). 
* `delete` - (Defaults to 20 mins) Used when terminating the db instance. 

## Import

RDS ddr instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_rds_ddr_instance.example rm-abc12345678
```
