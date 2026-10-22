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

-> **NOTE:** Available since v2.1.0.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_polardb_polar_fs_endpoint&exampleId=b9edaed9-8743-442a-f655-7daa9621bc4e8c3a848a&activeTab=example&spm=docs.r.polardb_polar_fs_endpoint.0.b9edaed987&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_polardb_polar_fs_endpoint" "example" {
  db_cluster_id        = alicloud_polardb_cluster.example.id
  polar_fs_instance_id = alicloud_polardb_polar_fs.example.id
  endpoint_type        = "S3Gateway"
  vpc_id               = alicloud_vpc.example.id
  vswitch_id           = alicloud_vswitch.example.id
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_polardb_polar_fs_endpoint&spm=docs.r.polardb_polar_fs_endpoint.example&intl_lang=EN_US)


## Argument Reference

* `db_cluster_id` - (Required, ForceNew, Available since v2.1.0) The ID of the PolarDB cluster associated with the PolarFS instance.
* `polar_fs_instance_id` - (Required, ForceNew, Available since v2.1.0) The PolarFS instance ID.
* `endpoint_type` - (Required, ForceNew, Available since v2.1.0) The endpoint type. Valid values: `Nas` and `S3Gateway`.
* `vpc_id` - (Optional, ForceNew, Available since v2.1.0) The VPC ID.
* `vswitch_id` - (Optional, ForceNew, Available since v2.1.0) The vSwitch ID.
* `db_endpoint_description` - (Optional, ForceNew, Available since v2.1.0) The endpoint description.

## Attributes Reference

* `id` - (Available since v2.1.0) The resource ID in the format `<db_cluster_id>:<db_endpoint_id>`.
* `db_endpoint_id` - (Available since v2.1.0) The endpoint ID.
* `address_items` - (Available since v2.1.0) Endpoint addresses.
  * `connection_string` - (Available since v2.1.0) The connection string.
  * `private_zone_connection_string` - (Available since v2.1.0) The PrivateZone connection string.
  * `ip_address` - (Available since v2.1.0) The IP address.
  * `port` - (Available since v2.1.0) The port.
  * `vpc_id` - (Available since v2.1.0) The VPC ID.
  * `vswitch_id` - (Available since v2.1.0) The vSwitch ID.
  * `net_type` - (Available since v2.1.0) The network type.

## Import

Use `<db_cluster_id>:<polar_fs_instance_id>:<db_endpoint_id>` to import an endpoint, e.g.

```shell
$ terraform import alicloud_polardb_polar_fs_endpoint.example pc-abc:pfs-abc:pe-abc
```
