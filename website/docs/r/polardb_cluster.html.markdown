---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_cluster"
sidebar_current: "docs-alicloud-resource-polardb-cluster"
description: |-
  Provides a PolarDB cluster resource.
---

# alicloud_polardb_cluster

Provides an PolarDB cluster resource. An PolarDB cluster is an isolated database
environment in the cloud. An PolarDB cluster can contain multiple user-created
databases.

-> **NOTE:** Available since v1.66.0.

## Example Usage

Create a PolarDB MySQL cluster

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_polardb_cluster&exampleId=14aedf61-7507-af1c-c239-50c54127fb66a22cd321&activeTab=example&spm=docs.r.polardb_cluster.0.14aedf6175&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_polardb_node_classes" "default" {
  db_type    = "MySQL"
  db_version = "8.0"
  category   = "Normal"
  pay_type   = "PostPaid"
}

resource "alicloud_vpc" "default" {
  vpc_name   = "terraform-example"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_polardb_node_classes.default.classes[0].zone_id
  vswitch_name = "terraform-example"
}

resource "alicloud_polardb_cluster" "default" {
  db_type       = "MySQL"
  db_version    = "8.0"
  db_node_class = data.alicloud_polardb_node_classes.default.classes.0.supported_engines.0.available_resources.0.db_node_class
  pay_type      = "PostPaid"
  vswitch_id    = alicloud_vswitch.default.id
  description   = "terraform-example"

  db_cluster_ip_array {
    db_cluster_ip_array_name = "default"
    security_ips             = ["1.2.3.4", "1.2.3.5"]
  }
  db_cluster_ip_array {
    db_cluster_ip_array_name = "default2"
    security_ips             = ["1.2.3.6"]
  }
}

```

When enabling TDE encryption, it is necessary to ensure that there is an AliyunRDSInstanceEncryptionDefaultRole role, and it is authorized under the account. If not, the following code can be used to create it.
Note: If there is only the role AliyunRDSSInceEncryptionDefaultRole under the account, this example may not be applicable.

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_polardb_cluster&exampleId=ad8f542a-4a74-b8d8-7a2c-228ccb357ce050ece07e&activeTab=example&spm=docs.r.polardb_cluster.1.ad8f542a4a&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_account" "current" {
}

data "alicloud_ram_roles" "roles" {
  name_regex = "AliyunRDSInstanceEncryptionDefaultRole"
}

resource "alicloud_ram_role" "default" {
  count       = length(data.alicloud_ram_roles.roles.roles) > 0 ? 0 : 1
  name        = "AliyunRDSInstanceEncryptionDefaultRole"
  document    = <<DEFINITION
    {
        "Statement": [
            {
               "Action": "sts:AssumeRole",
                "Effect": "Allow",
                "Principal": {
                    "Service": [
                        "rds.aliyuncs.com"
                    ]
                }
            }
        ],
        "Version": "1"
    }
	DEFINITION
  description = "RDS使用此角色来访问您在其他云产品中的资源"
}

resource "alicloud_resource_manager_policy_attachment" "default" {
  count             = length(data.alicloud_ram_roles.roles.roles) > 0 ? 0 : 1
  policy_name       = "AliyunRDSInstanceEncryptionRolePolicy"
  policy_type       = "System"
  principal_name    = length(data.alicloud_ram_roles.roles.roles) > 0 ? "${data.alicloud_ram_roles.roles.roles.0.name}@role.${data.alicloud_account.current.id}.onaliyunservice.com" : "${alicloud_ram_role.default[0].name}@role.${data.alicloud_account.current.id}.onaliyunservice.com"
  principal_type    = "ServiceRole"
  resource_group_id = "${data.alicloud_account.current.id}"
}

```

### Removing alicloud_polardb_cluster from your configuration

The alicloud_polardb_cluster resource allows you to manage your polardb cluster, but Terraform cannot destroy it if your cluster type is pre paid(post paid type can destroy normally). Removing this resource from your configuration will remove it from your statefile and management, but will not destroy the cluster. You can resume managing the cluster via the polardb Console.

## Argument Reference

The following arguments are supported:

* `db_type` - (Required, ForceNew) Database type. Value options: MySQL, Oracle, PostgreSQL.
* `db_version` - (Required, ForceNew) Database version. Value options can refer to the latest docs [CreateDBCluster](https://www.alibabacloud.com/help/en/polardb/latest/createdbcluster-1) `DBVersion`.
* `db_node_class` - (Required) The db_node_class of cluster node.
-> **NOTE:** Node specifications are divided into cluster version, single node version and History Library version. They can't change each other, but the general specification and exclusive specification of cluster version can be changed. 
  From version 1.204.0, If you need to create a Serverless cluster with MySQL , `db_node_class` can be set to `polar.mysql.sl.small`.
  From version 1.229.1, If you need to create a Serverless cluster with PostgreSQL 14 using the SENormal edition, `db_node_class` can be set to `polar.pg.sl.small.c`(x86 Architecture). Region can refer to the latest docs(https://help.aliyun.com/zh/polardb/polardb-for-postgresql/the-public-preview-of-polardb-for-postgresql-serverless-ends?spm=a2c4g.11186623.0.0.2e9f6cf0B4rIfC).
* `modify_type` - (Optional, Available since 1.71.2) Use as `db_node_class` change class, define upgrade or downgrade. Valid values are `Upgrade`, `Downgrade`, Default to `Upgrade`.
* `db_node_count` - (Optional, Available since 1.95.0)Number of the PolarDB cluster nodes, default is 2(Each cluster must contain at least a primary node and a read-only node). Add/remove nodes by modifying this parameter, valid values: [2~16].  
-> **NOTE:** To avoid adding or removing multiple read-only nodes by mistake, the system allows you to add or remove one read-only node at a time.
* `zone_id` - (Optional, ForceNew) The Zone to launch the DB cluster. it supports multiple zone.
* `pay_type` - (Optional) Valid values are `PrePaid`, `PostPaid`, Default to `PostPaid`.
* `renewal_status` - (Optional) Valid values are `AutoRenewal`, `Normal`, `NotRenewal`, Default to `NotRenewal`.
* `auto_renew_period` - (Optional) Auto-renewal period of an cluster, in the unit of the month. It is valid when pay_type is `PrePaid`. Valid value:1, 2, 3, 6, 12, 24, 36, Default to 1.
* `period` - (Optional) The duration that you will buy DB cluster (in month). It is valid when pay_type is `PrePaid`. Valid values: [1~9], 12, 24, 36.
-> **NOTE:** The attribute `period` is only used to create Subscription instance or modify the PayAsYouGo instance to Subscription. Once effect, it will not be modified that means running `terraform apply` will not effect the resource.
* `db_cluster_ip_array` - (Optional, Type: list, Available since 1.130.0) db_cluster_ip_array defines how users can send requests to your API. See [`db_cluster_ip_array`](#db_cluster_ip_array) below.
* `security_ips` - (Optional, Deprecated) This attribute has been deprecated from v1.130.0 and using `db_cluster_ip_array` sub-element `security_ips` instead.
  Its value is same as `db_cluster_ip_array` sub-element `security_ips` value and its db_cluster_ip_array_name is "default".
* `resource_group_id` (Optional, ForceNew, Computed, Available since 1.96.0) The ID of resource group which the PolarDB cluster belongs. If not specified, then it belongs to the default resource group.
* `vswitch_id` - (Optional, ForceNew) The virtual switch ID to launch DB instances in one VPC.
-> **NOTE:** If vswitch_id is not specified, system will get a vswitch belongs to the user automatically.
* `maintain_time` - (Optional) Maintainable time period format of the instance: HH:MMZ-HH:MMZ (UTC time)
* `description` - (Optional, Computed) The description of cluster.
* `collector_status` - (Optional, Available since 1.114.0) Specifies whether to enable or disable SQL data collector. Valid values are `Enable`, `Disabled`.
* `parameters` - (Optional) Set of parameters needs to be set after DB cluster was launched. Available parameters can refer to the latest docs [View database parameter templates](https://www.alibabacloud.com/help/en/polardb/latest/modifydbclusterparameters) .See [`parameters`](#parameters) below.
* `tags` - (Optional, Available since 1.68.0) A mapping of tags to assign to the resource.
  - Key: It can be up to 64 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It cannot be a null string.
  - Value: It can be up to 128 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It can be a null string.
* `tde_status` - (Optional, Available since 1.121.3) turn on TDE encryption. Valid values are `Enabled`, `Disabled`. Default to `Disabled`. TDE cannot be closed after it is turned on. 
-> **NOTE:** `tde_status` Cannot modify after created when `db_type` is `PostgreSQL` or `Oracle`.`tde_status` only support modification from `Disabled` to `Enabled` when `db_type` is `MySQL`.
* `encrypt_new_tables` - (Optional, Available since 1.124.1) turn on table auto encryption. Valid values are `ON`, `OFF`. Only MySQL 8.0 supports. 
-> **NOTE:** `encrypt_new_tables` Polardb MySQL 8.0 cluster, after TDE and Automatic Encryption are enabled, all newly created tables are automatically encrypted in the cluster.
* `encryption_key` - (Optional, Available since 1.200.0) The ID of the custom key. `encryption_key` cannot be modified after TDE is opened.
* `role_arn` - (Optional, Available since 1.200.0) The Alibaba Cloud Resource Name (ARN) of the RAM role. A RAM role is a virtual identity that you can create within your Alibaba Cloud account. For more information see [RAM role overview](https://www.alibabacloud.com/help/en/resource-access-management/latest/ram-role-overview).
* `security_group_ids` - (Optional, Available since 1.128.0) The ID of the security group. Separate multiple security groups with commas (,). You can add a maximum of three security groups to a cluster.
-> **NOTE:** Because of data backup and migration, change DB cluster type and storage would cost 15~20 minutes. Please make full preparation before changing them.
* `deletion_lock` - (Optional, Available since 1.169.0) turn on table deletion_lock. Valid values are 0, 1. 1 means to open the cluster protection lock, 0 means to close the cluster protection lock
-> **NOTE:**  Cannot modify after created when `pay_type` is `Prepaid` .`deletion_lock` the cluster protection lock can be turned on or off when `pay_type` is `Postpaid`.
* `backup_retention_policy_on_cluster_deletion` - (Optional, Available since 1.170.0) The retention policy for the backup sets when you delete the cluster.  Valid values are `ALL`, `LATEST`, `NONE`. Value options can refer to the latest docs [DeleteDBCluster](https://www.alibabacloud.com/help/en/polardb/latest/deletedbcluster-1)
* `imci_switch` - (Optional, Available since 1.173.0) Specifies whether to enable the In-Memory Column Index (IMCI) feature. Valid values are `ON`, `OFF`.
-> **NOTE:**  Only polardb MySQL Cluster version is available. The cluster with minor version number of 8.0.1 supports the column index feature, and the specific kernel version must be 8.0.1.1.22 or above.
-> **NOTE:**  The single node, the single node version of the history library, and the cluster version of the history library do not support column save indexes.
* `sub_category` - (Optional, Available since 1.177.0)  The category of the cluster. Valid values are `Exclusive`, `General`. Only MySQL supports.
* `creation_option` - (Optional, Available since 1.179.0) The method that is used to create a cluster. Valid values are `Normal`,`CloneFromPolarDB`,`CloneFromRDS`,`MigrationFromRDS`,`CreateGdnStandby`,`RecoverFromRecyclebin`. **NOTE:** From version 1.233.0, `creation_option` can be set to `RecoverFromRecyclebin`. Value options can refer to the latest docs [CreateDBCluster](https://www.alibabacloud.com/help/en/polardb/latest/createdbcluster-1) `CreationOption`.
* -> **NOTE:** The default value is Normal. If DBType is set to MySQL and DBVersion is set to 5.6 or 5.7, this parameter can be set to CloneFromRDS or MigrationFromRDS. If DBType is set to MySQL and DBVersion is set to 8.0, this parameter can be set to CreateGdnStandby. If `creation_option` is RecoverFromRecyclebin, you need to pass in the released source PolarDB cluster ID for this parameter. The DBType of the cluster recovered from the recycle bin and the source cluster must be consistent. For example, if the source cluster is MySQL 8.0, the cluster recovered from the recycle bin also needs to have its DBType set to MySQL and DBVersion set to 8.0.
* `creation_category` - (Optional, ForceNew, Computed, Available since 1.179.0) The edition of the PolarDB service. Valid values are `Normal`,`Basic`,`ArchiveNormal`,`NormalMultimaster`,`SENormal`.Value options can refer to the latest docs [CreateDBCluster](https://www.alibabacloud.com/help/en/polardb/latest/createdbcluster-1) `CreationCategory`.
-> **NOTE:** You can set this parameter to Basic only when DBType is set to MySQL and DBVersion is set to 5.6, 5.7, or 8.0. You can set this parameter to Archive only when DBType is set to MySQL and DBVersion is set to 8.0. From version 1.188.0, `creation_category` can be set to `NormalMultimaster`. From version 1.203.0, `creation_category` can be set to `SENormal`.
* `source_resource_id` - (Optional, Available since 1.179.0) The ID of the source RDS instance or the ID of the source PolarDB cluster. This parameter is required only when CreationOption is set to MigrationFromRDS, CloneFromRDS, or CloneFromPolarDB.Value options can refer to the latest docs [CreateDBCluster](https://www.alibabacloud.com/help/en/polardb/latest/createdbcluster-1) `SourceResourceId`.
* `gdn_id` - (Optional, Available since 1.179.0) The ID of the global database network (GDN).
-> **NOTE:** This parameter is required if CreationOption is set to CreateGdnStandby.
* `clone_data_point` - (Optional, Available since 1.179.0) The time point of data to be cloned. Valid values are `LATEST`,`BackupID`,`Timestamp`.Value options can refer to the latest docs [CreateDBCluster](https://www.alibabacloud.com/help/en/polardb/latest/createdbcluster-1) `CloneDataPoint`.
-> **NOTE:** If CreationOption is set to CloneFromRDS, the value of this parameter must be LATEST.
* `vpc_id` - (Optional, ForceNew, Computed, Available since v1.185.0) The id of the VPC.
* `storage_type` - (Optional, ForceNew, Available since v1.203.0) The storage type of the cluster. Enterprise storage type values are `PSL5`, `PSL4`. The standard version storage type values are `ESSDPL1`, `ESSDPL2`, `ESSDPL3`, `ESSDPL0`, `ESSDAUTOPL`. The standard version only supports MySQL and PostgreSQL.
* `provisioned_iops` - (Optional, ForceNew, Available since v1.229.1) The provisioned read/write IOPS of the ESSD AutoPL disk. Valid values: 0 to min{50,000, 1,000 × Capacity - Baseline IOPS}. Baseline IOPS = min{1,800 + 50 × Capacity, 50,000}.
-> **NOTE:** This parameter is available only if the StorageType parameter is set to ESSDAUTOPL.
* `storage_space` - (Optional, Computed, Available since v1.203.0) Storage space charged by space (monthly package). Unit: GB.
-> **NOTE:**  Valid values for PolarDB for MySQL Standard Edition: 20 to 32000. It is valid when pay_type are `PrePaid` ,`PostPaid`.
-> **NOTE:**  Valid values for PolarDB for MySQL Enterprise Edition: 50 to 100000.It is valid when pay_type is `PrePaid`.
* `storage_pay_type` - (Optional, ForceNew, Computed, Available since v1.210.0) The billing method of the storage. Valid values `Postpaid`, `Prepaid`.
* `hot_standby_cluster` - (Optional, Computed, ForceNew, Available since v1.203.0) Whether to enable the hot standby cluster. Valid values are `ON`, `OFF`.
* `strict_consistency` - (Optional, Computed, ForceNew, Available since v1.239.0) Whether the cluster has enabled strong data consistency across multiple zones. Valid values are `ON`, `OFF`. Available parameters can refer to the latest docs [CreateDBCluster](https://www.alibabacloud.com/help/en/polardb/latest/createdbcluster-1)
* `serverless_type` - (Optional, Available since v1.204.0) The type of the serverless cluster. Valid values `AgileServerless`, `SteadyServerless`. This parameter is valid only for serverless clusters.
* `serverless_steady_switch` - (Optional, Available since v1.211.1) Serverless steady-state switch. Valid values are `ON`, `OFF`. This parameter is valid only for serverless clusters.
  -> **NOTE:** When serverless_steady_switch is `ON` and serverless_type is `SteadyServerless`, parameters `scale_min`, `scale_max`, `scale_ro_num_min` and `scale_ro_num_max` are all required.
* `scale_min` - (Optional, Available since v1.204.0) The minimum number of PCUs per node for scaling. Valid values: 1 PCU to 31 PCUs. It is valid when serverless_type is `AgileServerless`. Valid values: 1 PCU to 8 PCUs.It is valid when serverless_type is `SteadyServerless`.· This parameter is valid only for serverless clusters.
* `scale_max` - (Optional, Available since v1.204.0) The maximum number of PCUs per node for scaling. Valid values: 1 PCU to 32 PCUs. It is valid when serverless_type is `AgileServerless`. Valid values: 1 PCU to 8 PCUs.It is valid when serverless_type is `SteadyServerless`. This parameter is valid only for serverless clusters.
* `scale_ro_num_min` - (Optional, Available since v1.204.0) The minimum number of read-only nodes for scaling. Valid values: 0 to 15 . It is valid when serverless_type is `AgileServerless`. Valid values: 0 to 7 .It is valid when serverless_type is `SteadyServerless`. This parameter is valid only for serverless clusters.
* `scale_ro_num_max` - (Optional, Available since v1.204.0) The maximum number of read-only nodes for scaling. Valid values: 0 to 15. It is valid when serverless_type is `AgileServerless`. Valid values: 0 to 7. It is valid when serverless_type is `SteadyServerless`. This parameter is valid only for serverless clusters.
* `allow_shut_down` - (Optional, Available since v1.204.0) Specifies whether to enable the no-activity suspension feature. Default value: false. Valid values are `true`, `false`. This parameter is valid only for serverless clusters.
* `seconds_until_auto_pause` - (Optional, Computed, Available since v1.204.0) The detection period for No-activity Suspension. Valid values: 300 to 86,4005. Unit: seconds. The detection duration must be a multiple of 300 seconds. This parameter is valid only for serverless clusters.
* `scale_ap_ro_num_min` - (Optional, Available since v1.211.1) Number of Read-only Columnar Nodes. Valid values: 0 to 7. This parameter is valid only for serverless clusters. This parameter is required when there are column nodes that support steady-state serverless.
* `scale_ap_ro_num_max` - (Optional, Available since v1.211.1) Number of Read-only Columnar Nodes. Valid values: 0 to 7. This parameter is valid only for serverless clusters. This parameter is required when there are column nodes that support steady-state serverless.
* `upgrade_type` - (Optional, Available since v1.208.1) Version upgrade type. Valid values are PROXY, DB, ALL. PROXY means upgrading the proxy version, DB means upgrading the db version, ALL means upgrading both db and proxy versions simultaneously.
* `from_time_service` - (Optional, Available since v1.208.1) Immediate or scheduled kernel version upgrade. Valid values are `true`, `false`. True means immediate execution, False means scheduled execution.
* `planned_start_time` - (Optional, Available since v1.208.1) The earliest time to start executing a scheduled (i.e. within the target time period) kernel version upgrade task. The format is YYYY-MM-DDThh: mm: ssZ (UTC).
-> **NOTE:** The starting time range is any time point within the next 24 hours. For example, the current time is 2021-01-14T09:00:00Z, and the allowed start time range for filling in here is 2021-01-14T09:00:00Z~2021-01-15T09:00:00Z. If this parameter is left blank, the kernel version upgrade task will be executed immediately by default.
* `planned_end_time` - (Optional, Available since v1.208.1) The latest time to start executing the target scheduled task. The format is YYYY-MM-DDThh: mm: ssZ (UTC).
-> **NOTE:** The latest time must be 30 minutes or more later than the start time. If PlannedStartTime is set but this parameter is not specified, the latest time to execute the target task defaults to the start time+30 minutes. For example, when the PlannedStartTime is set to 2021-01-14T09:00:00Z and this parameter is left blank, the target task will start executing at the latest on 2021-01-14T09:30:00Z.
* `proxy_type` - (Optional, Available since 1.210.0) The type of PolarProxy. Valid values are `EXCLUSIVE` `GENERAL`.
  -> **NOTE:** This parameter is valid for both standard and enterprise clusters.
* `proxy_class` - (Optional, Available since 1.210.0) The specifications of the Standard Edition PolarProxy. Available parameters can refer to the latest docs [CreateDBCluster](https://www.alibabacloud.com/help/en/polardb/latest/createdbcluster-1)
  -> **NOTE:** This parameter is valid only for standard edition clusters.
* `loose_polar_log_bin` - (Optional, Computed, Available since 1.210.0) Enable the Binlog function. Default value: `OFF`. Valid values are `OFF`, `ON`.
  -> **NOTE:** This parameter is valid only MySQL Engine supports.
* `loose_xengine` - (Optional, Available since v1.232.0) Specifies whether to enable X-Engine. Valid values are `ON`, `OFF`.
  -> **NOTE:** This parameter takes effect only if you do not set `creation_option` to CreateGdnStandby and you set `db_type` to MySQL and `db_version` to 8.0. To enable X-Engine on a node, make sure that the memory of the node is greater than or equal to 8 GB in size.
* `loose_xengine_use_memory_pct` - (Optional, Available since v1.232.0) Set the ratio to enable the X-Engine storage engine. Valid values: 10 to 90.
  -> **NOTE:** When the parameter `loose_xengine` is ON, `loose_xengine_use_memory_pct` takes effect.
* `db_node_num` - (Optional, Available since 1.210.0) The number of Standard Edition nodes. Default value: `1`. Valid values are `1`, `2`. From version 1.235.0, Valid values for PolarDB for MySQL Standard Edition: `1` to `8`. Valid values for PolarDB for MySQL Enterprise Edition: `1` to `16`.
* `parameter_group_id` - (Optional, Available since 1.210.0) The ID of the parameter template
  -> **NOTE:** You can call the [DescribeParameterGroups](https://www.alibabacloud.com/help/en/polardb/latest/describeparametergroups) operation to query the details of all parameter templates of a specified region, such as the ID of a parameter template.
* `lower_case_table_names`  - (Optional, ForceNew, Computed, Available since 1.210.0)  Specifies whether the table names are case-sensitive. Default value: `1`.  Valid values are `1`, `0`.
  -> **NOTE:** This parameter is valid only when the DBType parameter is set to MySQL.
* `default_time_zone` - (Optional, Computed, Available since 1.210.0) The time zone of the cluster. You can set the parameter to a value that is on the hour from -12:00 to +13:00 based on UTC. Example: 00:00. Default value: SYSTEM. This value indicates that the time zone of the cluster is the same as the time zone of the region.
  -> **NOTE:** This parameter is valid only when the DBType parameter is set to MySQL.
* `db_node_id` - (Optional, Available since v1.211.2) The ID of the node or node subscript. Node subscript values: 1 to 15.
* `hot_replica_mode` - (Optional, Available since v1.211.2) Indicates whether the hot standby feature is enabled. Valid values are `ON`, `OFF`. Only MySQL supports.
* `target_db_revision_version_code` - (Optional, Available since v1.216.0) The Version Code of the target version, whose parameter values can be obtained from the [DescribeDBClusterVersion](https://www.alibabacloud.com/help/en/polardb/latest/describedbclusterversion) interface.
* `compress_storage` - (Optional, Available since v1.232.0) Enable storage compression function. The value of this parameter is `ON`. Only MySQL supports.
  -> **NOTE:** When the value of db_type is not MySQL, the value of creation_option is neither empty nor Normal, and the value of storage_type is not PSL4, this field will be ignored.
### `db_cluster_ip_array`

The db_cluster_ip_array supports the following:

* `security_ips` - (Optional) List of IP addresses allowed to access all databases of a cluster. The list contains up to 1,000 IP addresses, separated by commas. Supported formats include 0.0.0.0/0, 10.23.12.24 (IP), and 10.23.12.24/24 (Classless Inter-Domain Routing (CIDR) mode. /24 represents the length of the prefix in an IP address. The range of the prefix length is [1,32]).
* `db_cluster_ip_array_name` - (Optional) The name of the IP whitelist group. The group name must be 2 to 120 characters in length and consists of lowercase letters and digits. It must start with a letter, and end with a letter or a digit.
  **NOTE:** If the specified whitelist group name does not exist, the whitelist group is created. If the specified whitelist group name exists, the whitelist group is modified. If you do not specify this parameter, the default group is modified. You can create a maximum of 50 IP whitelist groups for a cluster.
* `modify_mode` - (Optional) The method for modifying the IP whitelist. Valid values are `Cover`, `Append`, `Delete`.
  **NOTE:** There does not recommend setting modify_mode to `Append` or `Delete` and it will bring a potential diff error.  

### `parameters`

The parameters supports the following:

* `name` - (Required) Kernel parameter name.
* `value` - (Required) Kernel parameter value.

## Attributes Reference

The following attributes are exported:

* `id` - The PolarDB cluster ID.
* `connection_string` - (Available since 1.81.0) PolarDB cluster connection string. 
* `port` - (Available since 1.196.0) PolarDB cluster connection port. 
* `status` - (Available since 1.204.1) PolarDB cluster status.
* `create_time` - (Available since 1.204.1) PolarDB cluster creation time.
* `tde_region` - (Available since 1.200.0) The region where the TDE key resides.
-> **NOTE:** TDE can be enabled on clusters that have joined a global database network (GDN). After TDE is enabled on the primary cluster in a GDN, TDE is enabled on the secondary clusters in the GDN by default. The key used by the secondary clusters and the region for the key resides must be the same as the primary cluster. The region of the key cannot be modified.
-> **NOTE:** You cannot enable TDE for the secondary clusters in a GDN. Used to view user KMS activation status.

* `db_revision_version_list` - (Available since v1.216.0) The db_revision_version_list supports the following:
  * `release_type` - (Available since v1.216.0) Database version release status. Valid values are `Stable`, `Old`, `HighRisk`.
  * `revision_version_code` - (Available since v1.216.0) The revised version Code of the database engine is used to specify the upgrade to the target version.
  * `revision_version_name` - (Available since v1.216.0) The revision version number of the database engine.
  * `release_note` - (Available since v1.216.0) The revised version Code of the database engine is used to specify the upgrade to the target version.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 50 mins) Used when creating the polardb cluster (until it reaches the initial `Running` status).
* `update` - (Defaults to 50 mins) Used when updating the polardb cluster (until it reaches the initial `Running` status).
* `delete` - (Defaults to 10 mins) Used when terminating the polardb cluster.

## Import

PolarDB cluster can be imported using the id, e.g.

```shell
$ terraform import alicloud_polardb_cluster.example pc-abc12345678
```
