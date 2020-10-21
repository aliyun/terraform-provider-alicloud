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
  name       = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id            = alicloud_vpc.default.id
  cidr_block        = "172.16.0.0/24"
  availability_zone = data.alicloud_zones.default.zones[0].id
  name              = var.name
}

resource "alicloud_polardb_cluster" "default" {
  db_type       = "MySQL"
  db_version    = "5.6"
  db_node_class = "polar.mysql.x4.medium"
  pay_type      = "PostPaid"
  description   = var.name
  vswitch_id    = alicloud_vswitch.default.id
}
```

## Argument Reference

The following arguments are supported:

* `db_type` - (Required,ForceNew) Database type. Value options: MySQL, Oracle, PostgreSQL.
* `db_version` - (Required,ForceNew) Database version. Value options can refer to the latest docs [CreateDBCluster](https://help.aliyun.com/document_detail/98169.html) `DBVersion`.
* `db_node_class` - (Required) The db_node_class of cluster node.
* `modify_type` - (Optional, Available in 1.71.2+) Use as `db_node_class` change class, define upgrade or downgrade. Valid values are `Upgrade`, `Downgrade`, Default to `Upgrade`.
* `db_node_count` - (Optional, Available in 1.95.0+)Number of the PolarDB cluster nodes, default is 2(Each cluster must contain at least a primary node and a read-only node). Add/remove nodes by modifying this parameter, valid values: [2~16].
    **NOTE:** To avoid adding or removing multiple read-only nodes by mistake, the system allows you to add or remove one read-only node at a time.
* `zone_id` - (Optional) The Zone to launch the DB cluster. it supports multiple zone.
* `pay_type` - (Optional,ForceNew) Valid values are `PrePaid`, `PostPaid`, Default to `PostPaid`. Currently, the resource can not supports change pay type.
* `renewal_status` - (Optional) Valid values are `AutoRenewal`, `Normal`, `NotRenewal`, Default to `NotRenewal`.
* `auto_renew_period` - (Optional) Auto-renewal period of an cluster, in the unit of the month. It is valid when pay_type is `PrePaid`. Valid value:1, 2, 3, 6, 12, 24, 36, Default to 1.
* `period` - (Optional) The duration that you will buy DB cluster (in month). It is valid when pay_type is `PrePaid`. Valid values: [1~9], 12, 24, 36. Default to 1.
* `security_ips` - (Optional) List of IP addresses allowed to access all databases of an cluster. The list contains up to 1,000 IP addresses, separated by commas. Supported formats include 0.0.0.0/0, 10.23.12.24 (IP), and 10.23.12.24/24 (Classless Inter-Domain Routing (CIDR) mode. /24 represents the length of the prefix in an IP address. The range of the prefix length is [1,32]).
* `resource_group_id` (Optional, ForceNew, Computed, Available in 1.96.0+) The ID of resource group which the PolarDB cluster belongs. If not specified, then it belongs to the default resource group.
* `vswitch_id` - (Optional) The virtual switch ID to launch DB instances in one VPC.
* `maintain_time` - (Optional) Maintainable time period format of the instance: HH:MMZ-HH:MMZ (UTC time)
* `description` - (Optional) The description of cluster.
* `parameters` - (Optional) Set of parameters needs to be set after DB cluster was launched. Available parameters can refer to the latest docs [View database parameter templates](https://www.alibabacloud.com/help/doc-detail/26284.htm) .
* `tags` - (Optional, Available in v1.68.0+) A mapping of tags to assign to the resource.
    - Key: It can be up to 64 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It cannot be a null string.
    - Value: It can be up to 128 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It can be a null string.

-> **NOTE:** Because of data backup and migration, change DB cluster type and storage would cost 15~20 minutes. Please make full preparation before changing them.

### Removing alicloud_polardb_cluster from your configuration
 
The alicloud_polardb_cluster resource allows you to manage your polardb cluster, but Terraform cannot destroy it if your cluster type is pre paid(post paid type can destroy normally). Removing this resource from your configuration will remove it from your statefile and management, but will not destroy the cluster. You can resume managing the cluster via the polardb Console.
 
## Attributes Reference

The following attributes are exported:

* `id` - The PolarDB cluster ID.
* `connection_string` - (Available in 1.81.0+) PolarDB cluster connection string. When security_ips is configured, the address of cluster type endpoint will be returned, and if only "127.0.0.1" is configured, it will also be an empty string.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 mins) Used when creating the polardb cluster (until it reaches the initial `Running` status). 
* `update` - (Defaults to 30 mins) Used when updating the polardb cluster (until it reaches the initial `Running` status). 
* `delete` - (Defaults to 10 mins) Used when terminating the polardb cluster. 

## Import

PolarDB cluster can be imported using the id, e.g.

```
$ terraform import alicloud_polardb_cluster.example pc-abc12345678
```