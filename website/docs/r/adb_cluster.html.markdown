---
subcategory: "AnalyticDB for MySQL (ADB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_adb_cluster"
sidebar_current: "docs-alicloud-resource-adb-cluster"
description: |-
  Provides a ADB cluster resource.
---

# alicloud\_adb\_cluster

Provides a ADB cluster resource. An ADB cluster is an isolated database
environment in the cloud. An ADB cluster can contain multiple user-created
databases.

-> **DEPRECATED:**  This resource  has been deprecated from version `1.121.0`. Please use new resource [alicloud_adb_db_cluster](https://www.terraform.io/docs/providers/alicloud/r/adb_db_cluster).

-> **NOTE:** Available in v1.71.0+.

## Example Usage

### Create a ADB MySQL cluster

```
variable "name" {
  default = "adbClusterconfig"
}

variable "creation" {
  default = "ADB"
}

data "alicloud_zones" "default" {
  available_resource_creation = var.creation
}

resource "alicloud_vpc" "default" {
  vpc_name       = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id            = alicloud_vpc.default.id
  cidr_block        = "172.16.0.0/24"
  zone_id           = data.alicloud_zones.default.zones[0].id
  vswitch_name      = var.name
}

resource "alicloud_adb_cluster" "default" {
  db_cluster_version  = "3.0"
  db_cluster_category = "Cluster"
  db_node_class       = "C8"
  db_node_count       = 2
  db_node_storage     = 200
  pay_type            = "PostPaid"
  description         = var.name
  vswitch_id          = alicloud_vswitch.default.id
}
```

## Argument Reference

The following arguments are supported:

* `db_cluster_version` - (Optional, ForceNew) Cluster version. Value options: `3.0`, Default to `3.0`.
* `db_cluster_category` - (Required, ForceNew) Cluster category. Value options: `Basic`, `Cluster`.
* `db_node_class` - (Required) The db_node_class of cluster node.
* `db_node_count` - (Required) The db_node_count of cluster node.
* `db_node_storage` - (Required) The db_node_storage of cluster node.
* `zone_id` - (Optional) The Zone to launch the DB cluster.
* `pay_type` - (Optional) Field `pay_type` has been deprecated. New field `payment_type` instead.
* `payment_type` - (Optional) The payment type of the resource. Valid values are `PayAsYouGo` and `Subscription`. Default to `PayAsYouGo`. **Note:** The `payment_type` supports updating from v1.166.0+.
* `renewal_status` - (Optional) Valid values are `AutoRenewal`, `Normal`, `NotRenewal`, Default to `NotRenewal`.
* `auto_renew_period` - (Optional) Auto-renewal period of an cluster, in the unit of the month. It is valid when pay_type is `PrePaid`. Valid value:1, 2, 3, 6, 12, 24, 36, Default to 1.
* `period` - (Optional) The duration that you will buy DB cluster (in month). It is valid when pay_type is `PrePaid`. Valid values: [1~9], 12, 24, 36. Default to 1.
* `security_ips` - (Optional) List of IP addresses allowed to access all databases of an cluster. The list contains up to 1,000 IP addresses, separated by commas. Supported formats include 0.0.0.0/0, 10.23.12.24 (IP), and 10.23.12.24/24 (Classless Inter-Domain Routing (CIDR) mode. /24 represents the length of the prefix in an IP address. The range of the prefix length is [1,32]).
* `vswitch_id` - (Required, ForceNew) The virtual switch ID to launch DB instances in one VPC.
* `maintain_time` - (Optional) Maintainable time period format of the instance: HH:MMZ-HH:MMZ (UTC time)
* `description` - (Optional) The description of cluster.
* `tags` - (Optional) A mapping of tags to assign to the resource.
    - Key: It can be up to 64 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It cannot be a null string.
    - Value: It can be up to 128 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It can be a null string.

-> **NOTE:** Because of data backup and migration, change DB cluster type and storage would cost 15~30 minutes. Please make full preparation before changing them.

### Removing alicloud_adb_cluster from your configuration
 
The alicloud_adb_cluster resource allows you to manage your adb cluster, but Terraform cannot destroy it if your cluster type is pre paid(post paid type can destroy normally). Removing this resource from your configuration will remove it from your statefile and management, but will not destroy the cluster. You can resume managing the cluster via the adb Console.
 
## Attributes Reference

The following attributes are exported:

* `id` - The ADB cluster ID.
* `connection_string` - (Available in 1.93.0+) The connection string of the ADB cluster.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 50 mins) Used when creating the adb cluster (until it reaches the initial `Running` status). 
* `update` - (Defaults to 72 mins) Used when updating the adb cluster (until it reaches the initial `Running` status). 
* `delete` - (Defaults to 50 mins) Used when terminating the adb cluster. 

## Import

ADB cluster can be imported using the id, e.g.

```
$ terraform import alicloud_adb_cluster.example am-abc12345678
```
