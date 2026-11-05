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

-> **NOTE:** Available since v2.1.0.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_polardb_polar_fs&exampleId=2a7cdad9-8635-e304-1ec0-59569242e04cd43635aa&activeTab=example&spm=docs.r.polardb_polar_fs.0.2a7cdad986&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_polardb_polar_fs" "example" {
  db_type           = "PostgreSQL"
  creation_category = "high_performance"
  storage_type      = "essdpl1"
  storage_space     = 100
  pay_type          = "Postpaid"
  vpc_id            = alicloud_vpc.example.id
  vswitch_id        = alicloud_vswitch.example.id
  zone_id           = alicloud_vswitch.example.zone_id
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_polardb_polar_fs&spm=docs.r.polardb_polar_fs.example&intl_lang=EN_US)


## Argument Reference

The following arguments are supported:

* `storage_type` - (Optional, ForceNew, Available since v2.1.0) The storage type. Valid values: `local_redundancy`, `city_redundancy`, `essdpl1`, and `essdpl0`.
* `db_type` - (Optional, ForceNew, Available since v2.1.0) The database ecosystem. Valid values: `MySQL` and `PostgreSQL`.
* `vpc_id` - (Optional, ForceNew, Available since v2.1.0) The ID of the VPC.
* `vswitch_id` - (Optional, ForceNew, Available since v2.1.0) The ID of the vSwitch.
* `zone_id` - (Optional, ForceNew, Available since v2.1.0) The zone ID.
* `storage_space` - (Optional, ForceNew, Int, Available since v2.1.0) The storage capacity in GB. Valid values: `10` to `100000`.
* `db_cluster_id` - (Optional, ForceNew, Available since v2.1.0) The ID of the associated PolarDB cluster.
* `pay_type` - (Optional, ForceNew, Available since v2.1.0) The billing method. Valid value: `Postpaid`.
* `accelerate_switch` - (Optional, ForceNew, Available since v2.1.0) The acceleration switch. Valid values: `ONLY` and `ON`.
* `accelerate_storage_size` - (Optional, ForceNew, Int, Available since v2.1.0) The accelerated storage capacity in GB.
* `creation_category` - (Optional, ForceNew, Available since v2.1.0) The product category. Valid values: `basic`, `cold`, and `high_performance`.
* `accelerate_type` - (Optional, ForceNew, Available since v2.1.0) The acceleration type. Valid values: `juice` and `alluxio`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - (Available since v2.1.0) The ID of the PolarFS instance.
* `region_id` - (Available since v2.1.0) The region ID.
* `polar_fs_path` - (Available since v2.1.0) The PolarFS path.
* `status` - (Available since v2.1.0) The status of the PolarFS instance.
* `polar_fs_version` - (Available since v2.1.0) The PolarFS version.
* `description` - (Available since v2.1.0) The description of the PolarFS instance.
* `security_group_id` - (Available since v2.1.0) The security group ID.
* `create_time` - (Available since v2.1.0) The creation time.
* `expire_time` - (Available since v2.1.0) The expiration time.
* `expired` - (Available since v2.1.0) Whether the instance has expired.
* `polar_fs_type` - (Available since v2.1.0) The PolarFS type.
* `storage_used` - (Available since v2.1.0) The used storage.
* `bandwidth` - (Available since v2.1.0) The bandwidth.
* `bandwidth_base_line` - (Available since v2.1.0) The baseline bandwidth.
* `lock_mode` - (Available since v2.1.0) The lock mode.
* `accelerating_enable` - (Available since v2.1.0) Whether acceleration is enabled.
* `accelerated_storage_space` - (Available since v2.1.0) The accelerated storage capacity.
* `minor_version` - (Available since v2.1.0) The minor version.
* `client_download_path` - (Available since v2.1.0) The client download path.
* `relative_pfs_cluster_id` - (Available since v2.1.0) The related backend PolarFS cluster ID.
* `bucket_id` - (Available since v2.1.0) The bucket ID.
* `file_system_id` - (Available since v2.1.0) The file system ID.
* `meta_instance_name` - (Available since v2.1.0) The metadata instance name.
* `db_endpoint_id` - (Available since v2.1.0) The database endpoint ID.
* `maxscale_endpoint_id` - (Available since v2.1.0) The MaxScale endpoint ID.
* `meta_connection_string` - (Available since v2.1.0) The metadata connection string.
* `meta_maxscale_connection_string` - (Available since v2.1.0) The metadata MaxScale connection string.
* `authorized_user_arn_ids` - (Available since v2.1.0) The authorized RAM principal ARNs.

## Import

PolarFS instances can be imported using the instance ID, e.g.

```shell
$ terraform import alicloud_polardb_polar_fs.example pfs-2ze0i74ka607wck3
```
