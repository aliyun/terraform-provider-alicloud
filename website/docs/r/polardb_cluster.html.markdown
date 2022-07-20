---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_cluster"
sidebar_current: "docs-alicloud-resource-polardb-cluster"
description: |-
  Provides a PolarDB cluster resource.
---

# alicloud\_polardb\_cluster

Provides a PolarDB cluster resource. A PolarDB cluster is an isolated database
environment in the cloud. A PolarDB cluster can contain multiple user-created
databases.

-> **NOTE:** Available in v1.66.0+.

## Example Usage

### Create a PolarDB MySQL cluster

```terraform
variable "name" {
  default = "polardbClusterconfig"
}

variable "creation" {
  default = "PolarDB"
}

data "alicloud_zones" "default" {
  available_resource_creation = var.creation
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_zones.default.zones[0].id
  vswitch_name = var.name
}

resource "alicloud_polardb_cluster" "default" {
  db_type       = "MySQL"
  db_version    = "5.6"
  db_node_class = "polar.mysql.x4.medium"
  pay_type      = "PostPaid"
  description   = var.name
  vswitch_id    = alicloud_vswitch.default.id
  db_cluster_ip_array {
    db_cluster_ip_array_name = "default"
    security_ips             = ["1.2.3.4", "1.2.3.5"]
  }
  db_cluster_ip_array {
    db_cluster_ip_array_name = "test_ips1"
    security_ips             = ["1.2.3.6"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `db_type` - (Required,ForceNew) Database type. Value options: MySQL, Oracle, PostgreSQL.
* `db_version` - (Required,ForceNew) Database version. Value options can refer to the latest docs [CreateDBCluster](https://help.aliyun.com/document_detail/98169.html) `DBVersion`.
* `db_node_class` - (Required) The db_node_class of cluster node.
-> **NOTE:** Node specifications are divided into cluster version, single node version and History Library version. They can't change each other, but the general specification and exclusive specification of cluster version can be changed.
* `modify_type` - (Optional, Available in 1.71.2+) Use as `db_node_class` change class, define upgrade or downgrade. Valid values are `Upgrade`, `Downgrade`, Default to `Upgrade`.
* `db_node_count` - (Optional, Available in 1.95.0+)Number of the PolarDB cluster nodes, default is 2(Each cluster must contain at least a primary node and a read-only node). Add/remove nodes by modifying this parameter, valid values: [2~16].  
-> **NOTE:** To avoid adding or removing multiple read-only nodes by mistake, the system allows you to add or remove one read-only node at a time.
* `zone_id` - (Optional) The Zone to launch the DB cluster. it supports multiple zone.
* `pay_type` - (Optional) Valid values are `PrePaid`, `PostPaid`, Default to `PostPaid`.
* `renewal_status` - (Optional) Valid values are `AutoRenewal`, `Normal`, `NotRenewal`, Default to `NotRenewal`.
* `auto_renew_period` - (Optional) Auto-renewal period of an cluster, in the unit of the month. It is valid when pay_type is `PrePaid`. Valid value:1, 2, 3, 6, 12, 24, 36, Default to 1.
* `period` - (Optional) The duration that you will buy DB cluster (in month). It is valid when pay_type is `PrePaid`. Valid values: [1~9], 12, 24, 36.
-> **NOTE:** The attribute `period` is only used to create Subscription instance or modify the PayAsYouGo instance to Subscription. Once effect, it will not be modified that means running `terraform apply` will not effect the resource.
* `db_cluster_ip_array` - (Optional, Type: list, Available in 1.130.0+) db_cluster_ip_array defines how users can send requests to your API.
* `security_ips` - (Optional, Deprecated) This attribute has been deprecated from v1.130.0 and using `db_cluster_ip_array` sub-element `security_ips` instead.
  Its value is same as `db_cluster_ip_array` sub-element `security_ips` value and its db_cluster_ip_array_name is "default".
* `resource_group_id` (Optional, ForceNew, Computed, Available in 1.96.0+) The ID of resource group which the PolarDB cluster belongs. If not specified, then it belongs to the default resource group.
* `vswitch_id` - (Optional) The virtual switch ID to launch DB instances in one VPC.
-> **NOTE:** If vswitch_id is not specified, system will get a vswitch belongs to the user automatically.
* `maintain_time` - (Optional) Maintainable time period format of the instance: HH:MMZ-HH:MMZ (UTC time)
* `description` - (Optional) The description of cluster.
* `collector_status` - (Optional, Available in 1.114.0+) Specifies whether to enable or disable SQL data collector. Valid values are `Enable`, `Disabled`.
* `parameters` - (Optional) Set of parameters needs to be set after DB cluster was launched. Available parameters can refer to the latest docs [View database parameter templates](https://www.alibabacloud.com/help/doc-detail/98122.htm) .
* `tags` - (Optional, Available in 1.68.0+) A mapping of tags to assign to the resource.
  - Key: It can be up to 64 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It cannot be a null string.
  - Value: It can be up to 128 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It can be a null string.
* `tde_status` - (Optional, Available in 1.121.3+) turn on TDE encryption. Valid values are `Enabled`, `Disabled`. Default to `Disabled`. TDE cannot be closed after it is turned on. 
-> **NOTE:** `tde_status` Cannot modify after created when `db_type` is `PostgreSQL` or `Oracle`.`tde_status` only support modification from `Disabled` to `Enabled` when `db_type` is `MySQL`.
* `encrypt_new_tables` - (Optional, Available in 1.124.1+) turn on table auto encryption. Valid values are `ON`, `OFF`. Only MySQL 8.0 supports. 
-> **NOTE:** `encrypt_new_tables` Polardb MySQL 8.0 cluster, after TDE and Automatic Encryption are enabled, all newly created tables are automatically encrypted in the cluster.
* `security_group_ids` - (Optional, Available in 1.128.0+) The ID of the security group. Separate multiple security groups with commas (,). You can add a maximum of three security groups to a cluster.
-> **NOTE:** Because of data backup and migration, change DB cluster type and storage would cost 15~20 minutes. Please make full preparation before changing them.
* `deletion_lock` - (Optional, Available in 1.169.0+) turn on table deletion_lock. Valid values are 0, 1. 1 means to open the cluster protection lock, 0 means to close the cluster protection lock
-> **NOTE:**  Cannot modify after created when `pay_type` is `Prepaid` .`deletion_lock` the cluster protection lock can be turned on or off when `pay_type` is `Postpaid`.
* `backup_retention_policy_on_cluster_deletion` - (Optional, Available in 1.170.0+) The retention policy for the backup sets when you delete the cluster.  Valid values are `ALL`, `LATEST`, `NONE`. Value options can refer to the latest docs [DeleteDBCluster](https://help.aliyun.com/document_detail/98170.html)
* `imci_switch` - (Optional, Available in 1.173.0+) Specifies whether to enable the In-Memory Column Index (IMCI) feature. Valid values are `ON`, `OFF`.
-> **NOTE:**  Only polardb MySQL Cluster version is available. The cluster with minor version number of 8.0.1 supports the column index feature, and the specific kernel version must be 8.0.1.1.22 or above.
-> **NOTE:**  The single node, the single node version of the history library, and the cluster version of the history library do not support column save indexes.
* `sub_category` - (Optional, Available in 1.177.0+)  The category of the cluster. Valid values are `Exclusive`, `General`. Only MySQL supports.
### Block db_cluster_ip_array

The db_cluster_ip_array mapping supports the following:

* `security_ips` - (Optional) List of IP addresses allowed to access all databases of a cluster. The list contains up to 1,000 IP addresses, separated by commas. Supported formats include 0.0.0.0/0, 10.23.12.24 (IP), and 10.23.12.24/24 (Classless Inter-Domain Routing (CIDR) mode. /24 represents the length of the prefix in an IP address. The range of the prefix length is [1,32]).
* `db_cluster_ip_array_name` - (Optional) The name of the IP whitelist group. The group name must be 2 to 120 characters in length and consists of lowercase letters and digits. It must start with a letter, and end with a letter or a digit.
  **NOTE:** If the specified whitelist group name does not exist, the whitelist group is created. If the specified whitelist group name exists, the whitelist group is modified. If you do not specify this parameter, the default group is modified. You can create a maximum of 50 IP whitelist groups for a cluster.
* `modify_mode` - (Optional) The method for modifying the IP whitelist. Valid values are `Cover`, `Append`, `Delete`.
  **NOTE:** There does not recommend setting modify_mode to `Append` or `Delete` and it will bring a potential diff error.  
### Removing alicloud_polardb_cluster from your configuration

The alicloud_polardb_cluster resource allows you to manage your polardb cluster, but Terraform cannot destroy it if your cluster type is pre paid(post paid type can destroy normally). Removing this resource from your configuration will remove it from your statefile and management, but will not destroy the cluster. You can resume managing the cluster via the polardb Console.

## Attributes Reference

The following attributes are exported:

* `id` - The PolarDB cluster ID.
* `connection_string` - (Available in 1.81.0+) PolarDB cluster connection string. When security_ips is configured, the address of cluster type endpoint will be returned, and if only "127.0.0.1" is configured, it will also be an empty string.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 50 mins) Used when creating the polardb cluster (until it reaches the initial `Running` status).
* `update` - (Defaults to 50 mins) Used when updating the polardb cluster (until it reaches the initial `Running` status).
* `delete` - (Defaults to 10 mins) Used when terminating the polardb cluster.

## Import

PolarDB cluster can be imported using the id, e.g.

```
$ terraform import alicloud_polardb_cluster.example pc-abc12345678
```
