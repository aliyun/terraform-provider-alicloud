---
subcategory: "PolarDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_polardb_polar_fs_endpoint"
sidebar_current: "docs-alicloud-resource-polardb-polar-fs-endpoint"
description: |-
  Provides a PolarDB PolarFS endpoint resource.
---

# alicloud_polardb_polar_fs_endpoint

Provides a PolarDB PolarFS endpoint resource.

-> **NOTE:** Available since v1.294.0.

## Example Usage

```terraform
resource "alicloud_polardb_polar_fs_endpoint" "example" {
  db_cluster_id        = alicloud_polardb_cluster.example.id
  polar_fs_instance_id = alicloud_polardb_polar_fs.example.id
  endpoint_type        = "S3Gateway"
  vpc_id               = alicloud_vpc.example.id
  vswitch_id           = alicloud_vswitch.example.id
}
```

## Argument Reference

* `db_cluster_id` - (Required, ForceNew, Available since v1.294.0) The ID of the PolarDB cluster associated with the PolarFS instance.
* `polar_fs_instance_id` - (Required, ForceNew, Available since v1.294.0) The PolarFS instance ID.
* `endpoint_type` - (Required, ForceNew, Available since v1.294.0) The endpoint type. Valid values: `Nas`, `S3Gateway`, and `S3Pvtz`.
* `vpc_id` - (Optional, ForceNew, Available since v1.294.0) The VPC ID.
* `vswitch_id` - (Optional, ForceNew, Available since v1.294.0) The vSwitch ID.
* `db_endpoint_description` - (Optional, ForceNew, Available since v1.294.0) The endpoint description.

## Attributes Reference

* `id` - (Available since v1.294.0) The resource ID in the format `<db_cluster_id>:<db_endpoint_id>`.
* `db_endpoint_id` - (Available since v1.294.0) The endpoint ID.
* `address_items` - (Available since v1.294.0) Endpoint addresses.
  * `connection_string` - (Available since v1.294.0) The connection string.
  * `private_zone_connection_string` - (Available since v1.294.0) The PrivateZone connection string.
  * `ip_address` - (Available since v1.294.0) The IP address.
  * `port` - (Available since v1.294.0) The port.
  * `vpc_id` - (Available since v1.294.0) The VPC ID.
  * `vswitch_id` - (Available since v1.294.0) The vSwitch ID.
  * `net_type` - (Available since v1.294.0) The network type.

## Import

Use `<db_cluster_id>:<polar_fs_instance_id>:<db_endpoint_id>` to import an endpoint, e.g.

```shell
$ terraform import alicloud_polardb_polar_fs_endpoint.example pc-abc:pfs-abc:pe-abc
```
