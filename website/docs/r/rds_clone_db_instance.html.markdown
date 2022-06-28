---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_rds_clone_db_instance"
sidebar_current: "docs-alicloud-resource-rds-clone-db-instance"
description: |-
  Provides a Alicloud RDS Clone DB Instance resource.
---

# alicloud\_rds\_clone\_db\_instance

Provides a RDS Clone DB Instance resource.

For information about RDS Clone DB Instance and how to use it, see [What is ApsaraDB for RDS](https://www.alibabacloud.com/help/en/doc-detail/26092.htm).

-> **NOTE:** Available in v1.149.0+.

## Example Usage

### Create a RDS MySQL clone instance

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
  name       = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "example" {
  vpc_id     = alicloud_vpc.example.id
  cidr_block = "172.16.0.0/24"
  zone_id    = data.alicloud_zones.example.zones[0].id
  name       = var.name
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

resource "alicloud_rds_clone_db_instance" "example" {
  source_db_instance_id    = alicloud_db_instance.example.id
  db_instance_storage_type = "local_ssd"
  payment_type             = "PayAsYouGo"
  restore_time             = "2021-11-24T11:25:00Z"
  db_instance_storage      = "30"
}
```

## Argument Reference

The following arguments are supported:
* `source_db_instance_id` - (Required, ForceNew) The source db instance id.
* `db_instance_storage_type` - (Required) The type of storage media that is used for the new instance. Valid values:
  * **local_ssd**: local SSDs
  * **cloud_ssd**: standard SSDs
  * **cloud_essd**: enhanced SSDs (ESSDs) of performance level 1 (PL1)
  * **cloud_essd2**: ESSDs of PL2
  * **cloud_essd3**: ESSDs of PL3
* `payment_type` - (Required) The billing method of the new instance. Valid values: `PayAsYouGo` and `Subscription`.
* `db_instance_class` - (Optional, Computed) The instance type of the new instance. For information, see [Primary ApsaraDB RDS instance types](https://www.alibabacloud.com/doc-detail/26312.htm).
* `restore_time` - (Optional) The point in time to which you want to restore the data of the original instance. The point in time must fall within the specified log backup retention period. The time follows the ISO 8601 standard in the yyyy-MM-ddTHH:mm:ssZ format. The time must be in UTC.
* `backup_id` - (Optional) The ID of the data backup file you want to use. You can call the DescribeBackups operation to query the most recent data backup file list.

-> **NOTE:** You must specify at least one of the BackupId and RestoreTime parameters.
* `db_instance_storage` - (Optional, Computed) The storage capacity of the new instance. Unit: GB. The storage capacity increases in increments of 5 GB. For more information, see [Primary ApsaraDB RDS instance types](https://www.alibabacloud.com/doc-detail/26312.htm).

-> **NOTE:** The default value of this parameter is the storage capacity of the original instance.
* `restore_table` - (Optional) Specifies whether to restore only the databases and tables that you specify. The value 1 specifies to restore only the specified databases and tables. If you do not want to restore only the specified databases or tables, you can choose not to specify this parameter.
* `backup_type` - (Optional) The type of backup that is used to restore the data of the original instance. Valid values:
  * **FullBackup**: full backup
  * **IncrementalBackup**: incremental backup
* `vpc_id` - (Optional, Computed) The ID of the VPC to which the new instance belongs.

-> **NOTE:** Make sure that the VPC resides in the specified region.
* `vswitch_id` - (Optional, Computed) The ID of the vSwitch associated with the specified VPC.

-> **NOTE:** Make sure that the vSwitch belongs to the specified VPC and region.
* `private_ip_address` - (Optional, Computed) The intranet IP address of the new instance must be within the specified vSwitch IP address range. By default, the system automatically allocates by using **VPCId** and **VSwitchId**.
* `used_time` - (Optional) The subscription period of the new instance. This parameter takes effect only when you select the subscription billing method for the new instance. Valid values:
  * If you set the `Period` parameter to Year, the value of the UsedTime parameter ranges from 1 to 3.
  * If you set the `Period` parameter to Month, the value of the UsedTime parameter ranges from 1 to 9.

-> **NOTE:** If you set the payment_type parameter to Subscription, you must specify the used_time parameter.
* `period` - (Optional) The period. Valid values: `Month`, `Year`.

-> **NOTE:** If you set the payment_type parameter to Subscription, you must specify the period parameter.
* `deletion_protection` - (Optional, Available in 1.167.0+) The switch of delete protection. Valid values:
  - true: delete protect.
  - false: no delete protect.

-> **NOTE:** `deletion_protection` is valid only when attribute `payment_type` is set to `PayAsYouGo`, supported engine type: **MySQL**, **PostgresSQL**, **MariaDB**, **MSSQL**.
* `acl` - (Optional, Computed) This parameter is only supported by the RDS PostgreSQL cloud disk version. This parameter indicates the authentication method. It is allowed only when the public key of the client certificate authority is enabled. Valid values: `cert` and `perfer` and `verify-ca` and `verify-full (supported by RDS PostgreSQL above 12)`.
* `auto_upgrade_minor_version` - (Optional, Computed) How to upgrade the minor version of the instance. Valid values:
  * **Auto**: automatically upgrade the minor version.
  * **Manual**: It is not automatically upgraded. It is only mandatory when the current version is offline.
* `ca_type` - (Optional, Computed) This parameter is only supported by the RDS PostgreSQL cloud disk version. It indicates the certificate type. When the value of ssl_action is Open, the default value of this parameter is aliyun. Value range:
  * **aliyun**: using cloud certificates
  * **custom**: use a custom certificate. Valid values: `aliyun`, `custom`.
* `category` - (Optional, Computed, ForceNew) Instance series. Valid values:
  * **Basic**: Basic Edition
  * **HighAvailability**: High availability
  * **AlwaysOn**: Cluster Edition
  * **Finance**: Three-node Enterprise Edition.
* `certificate` - (Optional) The file that contains the certificate used for TDE.
* `client_ca_cert` - (Optional) This parameter is only supported by the RDS PostgreSQL cloud disk version. It indicates the public key of the client certification authority. If the value of client_ca_enabled is 1, this parameter must be configured.
* `client_ca_enabled` - (Optional) The client ca enabled.
* `client_cert_revocation_list` - (Optional) This parameter is only supported by the RDS PostgreSQL cloud disk version, which indicates that the client revokes the certificate file. If the value of client_crl_enabled is 1, this parameter must be configured.
* `client_crl_enabled` - (Optional) The client crl enabled.
* `connection_string_prefix` - (Optional) The connection string prefix.
* `db_instance_description` - (Optional) The db instance description.
* `db_name` - (Optional) The name of the database for which you want to enable TDE. Up to 50 names can be entered in a single request. If you specify multiple names, separate these names with commas (,).

-> **NOTE:** This parameter is available and must be specified only when the instance runs SQL Server 2019 SE or an Enterprise Edition of SQL Server.
* `db_names` - (Optional) The names of the databases that you want to create on the new instance.
* `dedicated_host_group_id` - (Optional) The ID of the dedicated cluster to which the new instance belongs. This parameter takes effect only when you create the new instance in a dedicated cluster.
* `direction` - (Optional) The direction. Valid values: `Auto`, `Down`, `TempUpgrade`, `Up`.
* `effective_time` - (Optional) The effective time.
* `encryption_key` - (Optional) The ID of the private key.

-> **NOTE:** This parameter is available only when the instance runs MySQL.
* `engine_version` - (Optional, Computed, ForceNew) Database version. Value:
  * MySQL:**5.5/5.6/5.7/8.0**
  * SQL Server:**2008r2/08r2_ent_ha/2012/2012_ent_ha/2012_std_ha/2012_web/2014_std_ha/2016_ent_ha/2016_std_ha/2016_web/2017_std_ha/2017_ent/2019_std_ha/2019_ent**
  * PostgreSQL:**9.4/10.0/11.0/12.0/13.0**
  * PPAS:**9.3/10.0**
  * MariaDB:**10.3**.
* `instance_network_type` - (Optional, Computed, ForceNew) The network type of the instance. Valid values:
  * **Classic**: Classic Network
  * **VPC**: VPC.
* `maintain_time` - (Optional, Computed) The maintainable time period of the instance. Format: <I> HH:mm</I> Z-<I> HH:mm</I> Z(UTC time).
* `password` - (Optional) The password of the certificate. 

-> **NOTE:** This parameter is available only when the instance runs SQL Server 2019 SE or an Enterprise Edition of SQL Server.
* `port` - (Optional, Computed) The port.
* `private_key` - (Optional) The file that contains the private key used for TDE.
* `released_keep_policy` - (Optional) The released keep policy.
* `replication_acl` - (Optional, Computed) This parameter is only supported by the RDS PostgreSQL cloud disk version, indicating the authentication method of the replication permission. It is only allowed when the public key of the client certificate authority is enabled. Valid values: `cert` and `perfer` and `verify-ca` and `verify-full (supported by RDS PostgreSQL above 12)`.
* `resource_group_id` - (Optional) The resource group id.
* `role_arn` - (Optional) The Alibaba Cloud Resource Name (ARN) of a RAM role. A RAM role is a virtual RAM identity that you can create within your Alibaba Cloud account. For more information, see [RAM role overview](https://www.alibabacloud.com/doc-detail/93689.htm).

-> **NOTE:** This parameter is available only when the instance runs MySQL.
* `security_ips` - (Optional, Computed) The IP address whitelist of the instance. Separate multiple IP addresses with commas (,) and cannot be repeated. The following two formats are supported:
  * IP address form, for example: 10.23.12.24.
  * CIDR format, for example, 10.23.12.0/24 (no Inter-Domain Routing, 24 indicates the length of the prefix in the address, ranging from 1 to 32).

-> **NOTE:** each instance can add up to 1000 IP addresses or IP segments, that is, the total number of IP addresses or IP segments in all IP whitelist groups cannot exceed 1000. When there are more IP addresses, it is recommended to merge them into IP segments, for example, 10.23.12.0/24.
* `server_cert` - (Optional, Computed) This parameter is only supported by the RDS PostgreSQL cloud disk version. It indicates the content of the server certificate. If the CAType value is custom, this parameter must be configured.
* `server_key` - (Optional, Computed) This parameter is only supported by the RDS PostgreSQL cloud disk version. It indicates the private key of the server certificate. If the value of CAType is custom, this parameter must be configured.
* `source_biz` - (Optional) The source biz.
* `ssl_enabled` - (Optional, Computed) Enable or disable SSL. Valid values: `0` and `1`.
* `switch_time` - (Optional) The time at which you want to apply the specification changes. The time follows the ISO 8601 standard in the yyyy-MM-ddTHH:mm:ssZ format. The time must be in UTC.
* `ha_mode` - (Optional, Computed) The high availability mode. Valid values:
  * **RPO**: Data persistence is preferred. The instance preferentially ensures data reliability to minimize data loss. Use this mode if you have higher requirements on data consistency.
  * **RTO**: Instance availability is preferred. The instance restores services as soon as possible to ensure availability. Use this mode if you have higher requirements on service availability.
* `sync_mode` - (Optional, Computed) [The data replication mode](https://www.alibabacloud.com/help/doc-detail/96055.htm). Valid values:
  * **Sync**: strong synchronization
  * **Semi-sync**: Semi-synchronous
  * **Async**: asynchronous

-> **NOTE:** SQL Server 2017 cluster version is currently not supported.
* `table_meta` - (Optional) The information about the databases and tables that you want to restore. Format:
  [{"type":"db","name":"The original name of Database 1","newname":"The new name of Database 1","tables":[{"type":"table","name":"The original name of Table 1 in Database 1","newname":"The new name of Table 1 in Database 1"},{"type":"table","name":"The original name of Table 2 in Database 1","newname":"The new name of Table 2 in Database 1"}]},{"type":"db","name":"The original name of Database 2","newname":"The new name of Database 2","tables":[{"type":"table","name":"The original name of Table 1 in Database 2","newname":"The new name of Table 1 in Database 2"},{"type":"table","name":"The original name of Table 2 in Database 2","newname":"The new name of Table 2 in Database 2"}]}]
* `tde_status` - (Optional) Specifies whether to enable TDE. Valid values:
  * Enabled
  * Disabled
* `zone_id` - (Optional, Computed, ForceNew) The ID of the zone to which the new instance belongs. You can call the [DescribeRegions](https://www.alibabacloud.com/doc-detail/26243.htm) operation to query the most recent region list.

-> **NOTE:** The default value of this parameter is the ID of the zone to which the original instance belongs.
* `engine` - (Optional, Computed, ForceNew) Database type. Value options: MySQL, SQLServer, PostgreSQL, and PPAS.
* `parameters` - (Optional) Set of parameters needs to be set after DB instance was launched. Available parameters can refer to the latest docs [View database parameter templates](https://www.alibabacloud.com/help/doc-detail/26284.htm).
* `force_restart` - (Optional) Set it to true to make some parameter efficient when modifying them. Default to false.
* `tcp_connection_type` - (Optional, Available in 1.171.0+) The availability check method of the instance. Valid values:
  - **SHORT**: Alibaba Cloud uses short-lived connections to check the availability of the instance.
  - **LONG**: Alibaba Cloud uses persistent connections to check the availability of the instance.
* `pg_hba_conf` - (Optional, Available in 1.155.0+) The configuration of [AD domain](https://www.alibabacloud.com/help/en/doc-detail/349288.htm) (documented below).

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

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Clone DB Instance.
* `connection_string` - The database connection address.

-> **NOTE:** The parameter **DBInstanceNetType** determines whether the address is internal or public.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 300 mins) Used when create the Clone DB Instance.
* `update` - (Defaults to 30 mins) Used when update the Clone DB Instance.
* `delete` - (Defaults to 20 mins) Used when terminating the Clone DB instance.

## Import

RDS Clone DB Instance can be imported using the id, e.g.

```
$ terraform import alicloud_rds_clone_db_instance.example <id>
```