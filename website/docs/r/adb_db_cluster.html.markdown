---
subcategory: "AnalyticDB for MySQL (ADB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_adb_db_cluster"
sidebar_current: "docs-alicloud-resource-adb-db-cluster"
description: |-
  Provides a Alicloud AnalyticDB for MySQL (ADB) DBCluster resource.
---

# alicloud_adb_db_cluster

Provides a AnalyticDB for MySQL (ADB) DBCluster resource.

For information about AnalyticDB for MySQL (ADB) DBCluster and how to use it, see [What is DBCluster](https://www.alibabacloud.com/help/en/analyticdb-for-mysql/developer-reference/api-adb-2021-12-01-createdbcluster).

-> **NOTE:** Available since v1.121.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_adb_db_cluster&exampleId=aa431eca-925a-85a6-ebd6-cd03b6e898e4704d8bbd&activeTab=example&spm=docs.r.adb_db_cluster.0.aa431eca92&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_adb_zones" "default" {}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_adb_zones.default.ids.0
}

locals {
  vswitch_id = data.alicloud_vswitches.default.ids.0
}

resource "alicloud_adb_db_cluster" "cluster" {
  db_cluster_category = "MixedStorage"
  mode                = "flexible"
  compute_resource    = "8Core32GB"
  vswitch_id          = local.vswitch_id
  description         = var.name
}
```

### Removing alicloud_adb_cluster from your configuration

The alicloud_adb_cluster resource allows you to manage your adb cluster, but Terraform cannot destroy it if your cluster type is PrePaid(PostPaid type can destroy normally). Removing this resource from your configuration will remove it from your state file and management, but will not destroy the cluster. You can resume managing the cluster via the adb Console.

## Argument Reference

The following arguments are supported:

* `auto_renew_period` - (Optional, Int) Auto-renewal period of an cluster, in the unit of the month. It is valid when `payment_type` is `Subscription`. Valid values: `1`, `2`, `3`, `6`, `12`, `24`, `36`. Default Value: `1`.
* `compute_resource` - (Optional) The specifications of computing resources in elastic mode. The increase of resources can speed up queries. AnalyticDB for MySQL automatically scales computing resources. For more information, see [ComputeResource](https://www.alibabacloud.com/help/en/analyticdb-for-mysql/developer-reference/api-adb-2019-03-15-describecomputeresource)
* `db_cluster_category` - (Required) The db cluster category. Valid values: `Basic`, `Cluster`, `MixedStorage`.
* `db_cluster_class` - (Deprecated since v1.121.2) It duplicates with attribute db_node_class and is deprecated from 1.121.2.
* `db_cluster_version` - (Optional, ForceNew) The db cluster version. Valid values: `3.0`. Default Value: `3.0`.
* `db_node_class` - (Optional) The db node class. For more information, see [DBClusterClass](https://help.aliyun.com/document_detail/190519.html)
* `db_node_count` - (Optional, Int) The db node count.
* `db_node_storage` - (Optional, Int) The db node storage.
* `description` - (Optional) The description of DBCluster.
* `elastic_io_resource` - (Optional, Int) The elastic io resource.
* `maintain_time` - (Optional) The maintenance window of the cluster. Format: hh:mmZ-hh:mmZ.
* `mode` - (Required) The mode of the cluster. Valid values: `reserver`, `flexible`.
* `modify_type` - (Optional) The modify type.
* `pay_type` - (Deprecated since v1.166.0) Field `pay_type` has been deprecated. New field `payment_type` instead.
* `payment_type` - (Optional) The payment type of the resource. Valid values: `PayAsYouGo` and `Subscription`. Default Value: `PayAsYouGo`. **Note:** The `payment_type` supports updating from v1.166.0+.
* `period` - (Optional, Int) The duration that you will buy DB cluster (in month). It is valid when `payment_type` is `Subscription`. Valid values: [1~9], 12, 24, 36.
-> **NOTE:** The attribute `period` is only used to create Subscription instance or modify the PayAsYouGo instance to Subscription. Once effect, it will not be modified that means running `terraform apply` will not affect the resource.
* `renewal_status` - (Optional) Valid values are `AutoRenewal`, `Normal`, `NotRenewal`, Default to `NotRenewal`.
* `resource_group_id` - (Optional) The ID of the resource group.
* `security_ips` - (Optional, List) List of IP addresses allowed to access all databases of an cluster. The list contains up to 1,000 IP addresses, separated by commas. Supported formats include 0.0.0.0/0, 10.23.12.24 (IP), and 10.23.12.24/24 (Classless Inter-Domain Routing (CIDR) mode. /24 represents the length of the prefix in an IP address. The range of the prefix length is [1,32]).
* `vswitch_id` - (Optional, ForceNew) The vswitch id.
* `zone_id` - (Optional, ForceNew) The zone ID of the resource.
* `vpc_id` - (Optional, ForceNew, Available since v1.178.0) The vpc ID of the resource.
* `elastic_io_resource_size` - (Optional, Available since v1.207.0) The specifications of a single elastic resource node. Default Value: `8Core64GB`. Valid values:
  - `8Core64GB`: If you set `elastic_io_resource_size` to `8Core64GB`, the specifications of an EIU are 24 cores and 192 GB memory.
  - `12Core96GB`: If you set `elastic_io_resource_size` to `12Core96GB`, the specifications of an EIU are 36 cores and 288 GB memory.
* `disk_performance_level` - (Optional, Available since v1.207.0) The ESSD performance level. Default Value: `PL1`. Valid values: `PL1`, `PL2`, `PL3`.
* `disk_encryption` - (Optional, ForceNew, Bool, Available since v1.219.0) Specifies whether to enable disk encryption. Default Value: `false`. Valid values: `true`, `false`.
* `kms_id` - (Optional, ForceNew, Available since v1.219.0) The Key Management Service (KMS) ID that is used for disk encryption. `kms_id` is valid only when `disk_encryption` is set to `true`.
* `enable_ssl` - (Optional, Bool, Available since v1.230.0) Specifies whether to enable SSL encryption. Default Value: `false`. Valid values: `true`, `false`.
* `tags` - (Optional) A mapping of tags to assign to the resource.
  - Key: It can be up to 64 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It cannot be a null string.
  - Value: It can be up to 128 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It can be a null string.

-> **NOTE:** Because of data backup and migration, change DB cluster type and storage would cost 15~30 minutes. Please make full preparation before changing them.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of DBCluster.
* `connection_string` - The connection string of the cluster.
* `port` - (Available since v1.196.0) The connection port of the ADB cluster.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 120 mins) Used when create the DBCluster.
* `update` - (Defaults to 6  hours) Used when update the DBCluster.
* `delete` - (Defaults to 3 hours) Used when delete the DBCluster.

## Import

AnalyticDB for MySQL (ADB) DBCluster can be imported using the id, e.g.

```shell
$ terraform import alicloud_adb_db_cluster.example <id>
```
