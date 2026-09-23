---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_polar_fs_instances"
sidebar_current: "docs-alicloud-datasource-polardb-polar-fs-instances"
description: |-
  Provides a list of PolarDB PolarFS instances.
---

# alicloud_polardb_polar_fs_instances

Provides a list of PolarDB PolarFS instances.

-> **NOTE:** Available since v1.294.0.

## Example Usage

```terraform
data "alicloud_polardb_polar_fs_instances" "example" {
  ids = [alicloud_polardb_polar_fs.example.id]
}
```

## Argument Reference

* `ids` - (Optional, List, Available since v1.294.0) A list of PolarFS instance IDs.
* `description` - (Optional, Available since v1.294.0) The instance description.
* `relative_db_cluster_id` - (Optional, Available since v1.294.0) The associated PolarDB cluster ID.
* `polar_fs_type` - (Optional, Available since v1.294.0) The PolarFS type.
* `db_cluster_id` - (Optional, Available since v1.294.0) The cluster ID used by the API filter.
* `tags` - (Optional, Map, Available since v1.294.0) A mapping of tags to assign to the query.

## Attributes Reference

* `polar_fs_instances` - (Available since v1.294.0) The matched instances. It contains the following attributes:
  * `id` - (Available since v1.294.0) The PolarFS instance ID.
  * `path` - (Available since v1.294.0) The PolarFS path.
  * `status` - (Available since v1.294.0) The instance status.
  * `description` - (Available since v1.294.0) The instance description.
  * `create_time` - (Available since v1.294.0) The creation time.
  * `expire_time` - (Available since v1.294.0) The expiration time.
  * `pay_type` - (Available since v1.294.0) The billing method.
  * `region_id` - (Available since v1.294.0) The region ID.
  * `zone_id` - (Available since v1.294.0) The zone ID.
  * `storage_space` - (Available since v1.294.0) The storage capacity.
  * `storage_type` - (Available since v1.294.0) The storage type.
  * `expired` - (Available since v1.294.0) Whether the instance has expired.
  * `bandwidth` - (Available since v1.294.0) The bandwidth.
  * `polar_fs_type` - (Available since v1.294.0) The PolarFS type.
  * `vpc_id` - (Available since v1.294.0) The VPC ID.
  * `vswitch_id` - (Available since v1.294.0) The vSwitch ID.
  * `security_group_id` - (Available since v1.294.0) The security group ID.
  * `relative_db_cluster_id` - (Available since v1.294.0) The associated database cluster ID.
  * `category` - (Available since v1.294.0) The product category.
  * `accelerating_enable` - (Available since v1.294.0) Whether acceleration is enabled.
  * `accelerated_storage_space` - (Available since v1.294.0) The accelerated storage capacity.
  * `accelerate_type` - (Available since v1.294.0) The acceleration type.
  * `tags` - (Available since v1.294.0) The instance tags.
