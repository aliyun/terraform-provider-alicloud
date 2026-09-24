---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_polar_fs"
sidebar_current: "docs-alicloud-resource-polardb-polar-fs"
description: |-
  Provides a PolarDB PolarFS resource.
---

# alicloud_polardb_polar_fs

Provides a PolarDB PolarFS resource.

-> **NOTE:** Available since v1.294.0.

## Example Usage

```terraform
resource "alicloud_polardb_polar_fs" "example" {
  db_type           = "polardb_pg"
  creation_category = "high_performance"
  storage_type      = "essdpl1"
  storage_space     = 100
  pay_type          = "Postpaid"
  vpc_id            = alicloud_vpc.example.id
  vswitch_id        = alicloud_vswitch.example.id
  zone_id           = alicloud_vswitch.example.zone_id
}
```

## Argument Reference

The following arguments are supported:

* `storage_type` - (Optional, ForceNew, Available since v1.294.0) The storage type. Valid values: `local_redundancy`, `city_redundancy`, `essdpl1`, and `essdpl0`.
* `authorized_user_ids` - (Optional, ForceNew, Set, Available since v1.294.0) The IDs of the accounts authorized to use the PolarFS instance.
* `db_type` - (Optional, ForceNew, Available since v1.294.0) The database ecosystem. Valid values: `polardb_mysql`, `polardb_pg`, and `polar_agentic_db`.
* `vpc_id` - (Optional, ForceNew, Available since v1.294.0) The ID of the VPC.
* `vswitch_id` - (Optional, ForceNew, Available since v1.294.0) The ID of the vSwitch.
* `zone_id` - (Optional, ForceNew, Available since v1.294.0) The zone ID.
* `storage_space` - (Optional, ForceNew, Int, Available since v1.294.0) The storage capacity in GB. Valid values: `10` to `100000`.
* `db_cluster_id` - (Optional, ForceNew, Available since v1.294.0) The ID of the associated PolarDB cluster.
* `pay_type` - (Optional, ForceNew, Available since v1.294.0) The billing method. Valid values: `Postpaid` and `Prepaid`.
* `period` - (Optional, ForceNew, Available since v1.294.0) The subscription period unit. Valid values: `Month` and `Year`. Required when `pay_type` is `Prepaid`.
* `used_time` - (Optional, ForceNew, Available since v1.294.0) The subscription duration. Required when `pay_type` is `Prepaid`.
* `auto_renew` - (Optional, ForceNew, Bool, Available since v1.294.0) Whether to enable automatic renewal. Default: `false`.
* `accelerate_switch` - (Optional, ForceNew, Available since v1.294.0) The acceleration switch. Valid values: `ONLY` and `ON`.
* `accelerate_storage_size` - (Optional, ForceNew, Int, Available since v1.294.0) The accelerated storage capacity in GB.
* `creation_category` - (Optional, ForceNew, Available since v1.294.0) The product category. Valid values: `basic`, `cold`, and `high_performance`.
* `custom_bucket_count` - (Optional, ForceNew, Int, Available since v1.294.0) The number of custom buckets.
* `custom_bucket_path` - (Optional, ForceNew, Available since v1.294.0) The custom bucket path.
* `custom_oss_ak` - (Optional, ForceNew, Sensitive, Available since v1.294.0) The access key used to access custom OSS storage.
* `custom_oss_sk` - (Optional, ForceNew, Sensitive, Available since v1.294.0) The secret key used to access custom OSS storage.
* `auto_use_coupon` - (Optional, ForceNew, Bool, Available since v1.294.0) Whether to automatically use a coupon. Default: `true`.
* `promotion_code` - (Optional, ForceNew, Available since v1.294.0) The promotion code.
* `accelerate_type` - (Optional, ForceNew, Available since v1.294.0) The acceleration type. Valid values: `juice` and `alluxio`.
* `custom_bucket_path_list` - (Optional, ForceNew, Set, Available since v1.294.0) Custom buckets and paths. See [`custom_bucket_path_list`](#custom_bucket_path_list) below.

### `custom_bucket_path_list`

The `custom_bucket_path_list` block supports the following:

* `bucket` - (Required, Available since v1.294.0) The bucket endpoint.
* `path` - (Required, Available since v1.294.0) The path in the bucket.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - (Available since v1.294.0) The ID of the PolarFS instance.
* `region_id` - (Available since v1.294.0) The region ID.
* `polar_fs_path` - (Available since v1.294.0) The PolarFS path.
* `status` - (Available since v1.294.0) The status of the PolarFS instance.
* `polar_fs_version` - (Available since v1.294.0) The PolarFS version.
* `description` - (Available since v1.294.0) The description of the PolarFS instance.
* `security_group_id` - (Available since v1.294.0) The security group ID.
* `create_time` - (Available since v1.294.0) The creation time.
* `expire_time` - (Available since v1.294.0) The expiration time.
* `expired` - (Available since v1.294.0) Whether the instance has expired.
* `polar_fs_type` - (Available since v1.294.0) The PolarFS type.
* `storage_used` - (Available since v1.294.0) The used storage.
* `bandwidth` - (Available since v1.294.0) The bandwidth.
* `bandwidth_base_line` - (Available since v1.294.0) The baseline bandwidth.
* `lock_mode` - (Available since v1.294.0) The lock mode.
* `accelerating_enable` - (Available since v1.294.0) Whether acceleration is enabled.
* `accelerated_storage_space` - (Available since v1.294.0) The accelerated storage capacity.
* `minor_version` - (Available since v1.294.0) The minor version.
* `client_download_path` - (Available since v1.294.0) The client download path.
* `relative_pfs_cluster_id` - (Available since v1.294.0) The related backend PolarFS cluster ID.
* `bucket_id` - (Available since v1.294.0) The bucket ID.
* `file_system_id` - (Available since v1.294.0) The file system ID.
* `meta_instance_name` - (Available since v1.294.0) The metadata instance name.
* `db_endpoint_id` - (Available since v1.294.0) The database endpoint ID.
* `maxscale_endpoint_id` - (Available since v1.294.0) The MaxScale endpoint ID.
* `meta_connection_string` - (Available since v1.294.0) The metadata connection string.
* `meta_maxscale_connection_string` - (Available since v1.294.0) The metadata MaxScale connection string.
* `authorized_user_arn_ids` - (Available since v1.294.0) The authorized RAM principal ARNs.

## Import

PolarFS instances can be imported using the instance ID, e.g.

```shell
$ terraform import alicloud_polardb_polar_fs.example pfs-2ze0i74ka607wck3
```
