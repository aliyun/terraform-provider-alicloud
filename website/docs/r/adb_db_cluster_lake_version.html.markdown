---
subcategory: "AnalyticDB for MySQL (ADB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_adb_db_cluster_lake_version"
sidebar_current: "docs-alicloud-resource-adb-db-cluster-lake-version"
description: |-
  Provides a Alicloud AnalyticDB for MySQL (ADB) DB Cluster Lake Version resource.
---

# alicloud_adb_db_cluster_lake_version

Provides a AnalyticDB for MySQL (ADB) DB Cluster Lake Version resource.

For information about AnalyticDB for MySQL (ADB) DB Cluster Lake Version and how to use it, see [What is DB Cluster Lake Version](https://www.alibabacloud.com/help/en/analyticdb-for-mysql/developer-reference/api-adb-2021-12-01-createdbcluster).

-> **NOTE:** Available since v1.190.0.

-> **NOTE:** Destroying a `Subscription` instance fails with an error since v2.1.0. To destroy it, convert the instance to `PayAsYouGo` by setting `payment_type = "PayAsYouGo"` and re-applying, or unsubscribe it in the Alibaba Cloud console. To remove the instance from the state only, run `terraform state rm`.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_adb_db_cluster_lake_version&exampleId=44a0d0ac-5708-e0a7-dcbd-889815af140cbc29256a&activeTab=example&spm=docs.r.adb_db_cluster_lake_version.0.44a0d0ac57&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "ap-southeast-1"
}

data "alicloud_adb_zones" "default" {
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_adb_zones.default.ids.0
}

resource "alicloud_adb_db_cluster_lake_version" "default" {
  db_cluster_version            = "5.0"
  vpc_id                        = data.alicloud_vpcs.default.ids.0
  vswitch_id                    = data.alicloud_vswitches.default.ids.0
  zone_id                       = data.alicloud_adb_zones.default.ids.0
  compute_resource              = "16ACU"
  storage_resource              = "0ACU"
  payment_type                  = "PayAsYouGo"
  enable_default_resource_group = false
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_adb_db_cluster_lake_version&spm=docs.r.adb_db_cluster_lake_version.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `db_cluster_version` - (Required, ForceNew) The version of the cluster. Valid values: `5.0`.
* `vpc_id` - (Required, ForceNew) The vpc ID of the resource.
* `vswitch_id` - (Required, ForceNew) The ID of the vSwitch.
* `zone_id` - (Required, ForceNew) The zone ID of the resource.
* `payment_type` - (Required) The payment type of the resource. Valid values: `PayAsYouGo`, `Subscription`. **NOTE:** From version 1.245.0, `payment_type` can be set to `Subscription`. From version 1.272.0, `payment_type` can be modified.
* `secondary_vswitch_id` - (Optional, ForceNew, Available since v1.248.0) The ID of the secondary zone.
* `secondary_zone_id` - (Optional, ForceNew, Available since v1.248.0) The ID of the secondary vSwitch.
* `product_form` - (Optional, ForceNew, Available since v1.245.0) The product form of the cluster. Valid values:
  - `IntegrationForm`: Integrated.
  - `LegacyForm`: Data Lakehouse Edition.
* `product_version` - (Optional, ForceNew, Available since v1.245.0) The edition of the cluster. Valid values:
  - `BasicVersion`: Basic Edition.
  - `EnterpriseVersion`: Enterprise Edition.

  -> **NOTE:** `product_version` must be specified only when `product_form` is set to `IntegrationForm`.
* `compute_resource` - (Optional) The computing resources of the cluster.
* `storage_resource` - (Optional) The storage resources of the cluster.
* `reserved_node_size` - (Optional, Available since v1.245.0) The specifications of reserved resource nodes.
* `reserved_node_count` - (Optional, Int, Available since v1.245.0) The number of reserved resource nodes.
* `disk_encryption` - (Optional, Bool, ForceNew, Available since v1.245.0) Specifies whether to enable disk encryption. Valid values: `true`, `false`.
* `kms_id` - (Optional, ForceNew, Available since v1.245.0) The ID of the key that is used to encrypt disk data. `kms_id` is valid only when `disk_encryption` is set to `true`.
* `security_ips` - (Optional, Available since v1.198.0) The IP addresses in an IP address whitelist of a cluster. Separate multiple IP addresses with commas (,). You can add a maximum of 500 different IP addresses to a whitelist. The entries in the IP address whitelist must be in one of the following formats:
  - IP addresses, such as 10.23.XX.XX.
  - CIDR blocks, such as 10.23.xx.xx/24. In this example, 24 indicates that the prefix of each IP address in the IP whitelist is 24 bits in length. You can replace 24 with a value within the range of 1 to 32.
* `enable_ssl` - (Optional, Bool, Available since v1.245.0) Specifies whether to enable SSL encryption. Valid values: `true`, `false`.
* `db_cluster_description` - (Optional, Available since v1.198.0) The description of the cluster.
* `resource_group_id` - (Optional, Available since v1.211.1) The ID of the resource group.
* `period` - (Optional, Int, Available since v1.245.0) The subscription period of the subscription cluster. Valid values: `1` to `9`, `12`, `24`, `36`.
* `enable_default_resource_group` - (Optional, Bool) Whether to enable default allocation of resources to user_default resource groups.
* `source_db_cluster_id` - (Optional, Available since v1.211.1) The ID of the source AnalyticDB for MySQL Data Warehouse Edition cluster.
* `backup_set_id` - (Optional, Available since v1.211.1) The ID of the backup set that you want to use to restore data.
* `restore_type` - (Optional, Available since v1.211.1) The method that you want to use to restore data. Valid values:
  - `backup`: Restores data from a backup set. **NOTE:** You must set `source_db_cluster_id` and `backup_set_id`.
  - `timepoint`: Restores data to a point in time. **NOTE:** You must set `source_db_cluster_id` and `restore_to_time`.
* `restore_to_time` - (Optional, Available since v1.211.1) The point in time to which you want to restore data from the backup set.
* `db_cluster_network_type` - (Optional, String, ForceNew, Available since v2.2.0) The network type of the cluster. Valid value: `VPC`.
* `audit_log_status` - (Optional, String, Available since v2.2.0) The status of the SQL audit log. Valid values: `on`, `off`.
* `essd_cache_size` - (Optional, Int, Available since v2.2.0) The size of the ESSD cache. Unit: GB.
* `enable_essd_cache` - (Optional, Bool, Available since v2.2.0) Specifies whether to enable the ESSD cache feature. Valid values: `true`, `false`.
* `capacity` - (Optional, Int, Available since v2.2.0) The size of the lake storage acceleration cache. Unit: GB.
* `enable_lake_cache` - (Optional, Bool, Available since v2.2.0) Specifies whether to enable lake storage acceleration. Valid values: `true`, `false`.
* `enable_compaction_service` - (Optional, Bool, Available since v2.2.0) Specifies whether to enable the remote build (compaction) service. Valid values: `true`, `false`.
* `maintain_time` - (Optional, String, Available since v2.2.0) The maintenance window of the cluster. Format: `HH:00Z-HH:00Z`. The window must start on the hour and span exactly one hour, for example `04:00Z-05:00Z`.
* `auto_renewal_status` - (Optional, String, Available since v2.2.0) The auto-renewal status of a `Subscription` cluster. Valid values: `AutoRenewal`, `Normal`, `NotRenewal`. **NOTE:** This argument is valid only when `payment_type` is set to `Subscription`.
* `auto_renewal_period` - (Optional, Int, Available since v2.2.0) The auto-renewal duration of a `Subscription` cluster. Valid values: `1` to `11` when `auto_renewal_period_unit` is `Month`, `1`, `2`, `3` or `5` when `auto_renewal_period_unit` is `Year`. **NOTE:** This argument is valid only when `payment_type` is set to `Subscription`.
* `auto_renewal_period_unit` - (Optional, String, Available since v2.2.0) The unit of the auto-renewal duration of a `Subscription` cluster. Valid values: `Year`, `Month`. **NOTE:** This argument is valid only when `payment_type` is set to `Subscription`.
* `db_cluster_ip_array_name` - (Optional, String, Available since v2.2.0) The name of the IP address whitelist group. Defaults to `default`. **NOTE:** `security_ips` is applied to this group. Changing the name does not delete the previous group; remove it manually in the console if it is no longer needed.
* `db_cluster_ip_array_attribute` - (Optional, String, Available since v2.2.0) The attribute of the IP address whitelist group. Set it to `hidden` to hide the group in the console. **NOTE:** This argument is write-only. Groups that carry the `hidden` attribute are not returned by the whitelist query API, so the attribute and the IP list of such a group cannot be refreshed and drift is not detected.
* `modify_mode` - (Optional, String, Available since v2.2.0) The modification mode of the IP address whitelist. Valid values: `Cover`, `Append`, `Delete`.
* `db_version` - (Optional, String, Available since v2.2.0) The kernel version to which you want to upgrade the cluster. **NOTE:** `db_version` must be set when you change the kernel version. Use the version returned by `DescribeKernelVersion` and make sure it is newer than the current kernel version of the cluster.
* `switch_mode` - (Optional, Int, Available since v2.2.0) The time when the kernel version is switched. Valid values: `0` (switch immediately) and `1` (switch during the maintenance window). Defaults to `0`. **NOTE:** `switch_mode` takes effect only when `db_version` is set.
* `engine_type` - (Optional, String, Available since v2.2.0) The computing engine type of the SQL audit log. Valid values: `XIHE`, `SPARK`. **NOTE:** This argument requires `audit_log_status` to be set.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of DBCluster Lake Version.
* `status` - The status of the resource.
* `commodity_code` - The name of the service.
* `engine` - The engine of the database.
* `engine_version` - The engine version of the database.
* `expire_time` - The time when the cluster expires.
* `expired` - Indicates whether the cluster has expired.
* `lock_mode` - The lock mode of the cluster.
* `lock_reason` - The reason why the cluster is locked.
* `port` - The port that is used to access the cluster.
* `connection_string` - The endpoint of the cluster.
* `create_time` - The createTime of the cluster.
* `used_time` - (Available since v2.2.0) The subscription duration in months, derived from `period`. It is `0` for `PayAsYouGo` clusters and for subscription clusters without `period`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 50 mins) Used when create the DB Cluster Lake Version.
* `update` - (Defaults to 72 mins) Used when update the DB Cluster Lake Version.
* `delete` - (Defaults to 5 mins) Used when delete the DB Cluster Lake Version.

## Import

AnalyticDB for MySQL (ADB) DB Cluster Lake Version can be imported using the id, e.g.

```shell
$ terraform import alicloud_adb_db_cluster_lake_version.example <id>
```
