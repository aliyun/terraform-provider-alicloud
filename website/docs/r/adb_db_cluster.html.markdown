---
subcategory: "AnalyticDB for MySQL (ADB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_adb_db_cluster"
sidebar_current: "docs-alicloud-resource-adb-db-cluster"
description: |-
  Provides a Alicloud AnalyticDB for MySQL (ADB) DBCluster resource.
---

# alicloud\_adb\_db\_cluster

Provides a AnalyticDB for MySQL (ADB) DBCluster resource.

For information about AnalyticDB for MySQL (ADB) DBCluster and how to use it, see [What is DBCluster](https://www.alibabacloud.com/help/en/doc-detail/190519.htm).

-> **NOTE:** Available in v1.121.0+.

## Example Usage

Basic Usage

```terraform
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
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_zones.default.zones[0].id
  vswitch_name = var.name
}

resource "alicloud_adb_db_cluster" "this" {
  db_cluster_category = "Cluster"
  db_cluster_class    = "C8"
  db_node_count       = "4"
  db_node_storage     = "400"
  mode                = "reserver"
  db_cluster_version  = "3.0"
  payment_type        = "PayAsYouGo"
  vswitch_id          = alicloud_vswitch.default.id
  description         = "Test new adb again."
  maintain_time       = "23:00Z-00:00Z"
  tags = {
    Created = "TF-update"
    For     = "acceptance-test-update"
  }
  resource_group_id = "rg-aek2s7ylxx6****"
  security_ips      = ["10.168.1.12", "10.168.1.11"]
}
```

## Argument Reference

The following arguments are supported:

* `auto_renew_period` - (Optional) Auto-renewal period of an cluster, in the unit of the month. It is valid when `payment_type` is `Subscription`. Valid values: `1`, `2`, `3`, `6`, `12`, `24`, `36`. Default to `1`.
* `compute_resource` - (Optional) The specifications of computing resources in elastic mode. The increase of resources can speed up queries. AnalyticDB for MySQL automatically scales computing resources. For more information, see [ComputeResource](https://www.alibabacloud.com/help/en/doc-detail/144851.htm)
* `db_cluster_category` - (Required) The db cluster category. Valid values: `Basic`, `Cluster`, `MixedStorage`.
* `db_cluster_class` - (Deprecated) It duplicates with attribute db_node_class and is deprecated from 1.121.2.
* `db_cluster_version` - (Optional, ForceNew) The db cluster version. Value options: `3.0`, Default to `3.0`.
* `db_node_class` - (Optional, Computed) The db node class. For more information, see [DBClusterClass](https://help.aliyun.com/document_detail/190519.html)
* `db_node_count` - (Optional) The db node count.
* `db_node_storage` - (Optional) The db node storage.
* `description` - (Optional, Computed) The description of DBCluster.
* `elastic_io_resource` - (Optional) The elastic io resource.
* `maintain_time` - (Optional, Computed) The maintenance window of the cluster. Format: hh:mmZ-hh:mmZ.
* `mode` - (Required) The mode of the cluster. Valid values: `reserver`, `flexible`.
* `modify_type` - (Optional) The modify type.
* `pay_type` - (Optional, Computed) Field `pay_type` has been deprecated. New field `payment_type` instead.
* `payment_type` - (Optional, Computed) The payment type of the resource. Valid values are `PayAsYouGo` and `Subscription`. Default to `PayAsYouGo`. **Note:** The `payment_type` supports updating from v1.166.0+. 
* `period` - (Optional) The duration that you will buy DB cluster (in month). It is valid when `payment_type` is `Subscription`. Valid values: [1~9], 12, 24, 36.
-> **NOTE:** The attribute `period` is only used to create Subscription instance or modify the PayAsYouGo instance to Subscription. Once effect, it will not be modified that means running `terraform apply` will not affect the resource.
* `renewal_status` - (Optional) Valid values are `AutoRenewal`, `Normal`, `NotRenewal`, Default to `NotRenewal`.
* `resource_group_id` - (Optional, Computed) The ID of the resource group.
* `security_ips` - (Optional, Computed) List of IP addresses allowed to access all databases of an cluster. The list contains up to 1,000 IP addresses, separated by commas. Supported formats include 0.0.0.0/0, 10.23.12.24 (IP), and 10.23.12.24/24 (Classless Inter-Domain Routing (CIDR) mode. /24 represents the length of the prefix in an IP address. The range of the prefix length is [1,32]).
* `vswitch_id` - (Optional, ForceNew) The vswitch id.
* `zone_id` - (Optional, Computed, ForceNew) The zone ID of the resource.
* `tags` - (Optional) A mapping of tags to assign to the resource.
    - Key: It can be up to 64 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It cannot be a null string.
    - Value: It can be up to 128 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It can be a null string.

-> **NOTE:** Because of data backup and migration, change DB cluster type and storage would cost 15~30 minutes. Please make full preparation before changing them.

### Removing alicloud_adb_cluster from your configuration
 
The alicloud_adb_cluster resource allows you to manage your adb cluster, but Terraform cannot destroy it if your cluster type is pre paid(post paid type can destroy normally). Removing this resource from your configuration will remove it from your statefile and management, but will not destroy the cluster. You can resume managing the cluster via the adb Console.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of DBCluster.
* `connection_string` - The endpoint of the cluster.
* `status` - The status of the resource.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 50 mins) Used when create the DBCluster.
* `delete` - (Defaults to 50 mins) Used when delete the DBCluster.
* `update` - (Defaults to 72 mins) Used when update the DBCluster.

## Import

AnalyticDB for MySQL (ADB) DBCluster can be imported using the id, e.g.

```
$ terraform import alicloud_adb_db_cluster.example <id>
```
