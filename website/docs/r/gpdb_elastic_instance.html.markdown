---
subcategory: "AnalyticDB for PostgreSQL (GPDB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_gpdb_elastic_instance"
sidebar_current: "docs-alicloud-resource-gpdb-elastic-instance"
description: |-
  Provides a flexible storage mode AnalyticDB for PostgreSQL instance resource.
---

# alicloud\_gpdb\_elastic\_instance

Provides a AnalyticDB for PostgreSQL instance resource which storage type is flexible. Compared to the reserved storage ADB PG instance, you can scale up each disk and smoothly scale out nodes online.  
For more detail product introduction, see [here](https://www.alibabacloud.com/help/doc-detail/141368.htm).

-> **NOTE:**  Available in 1.127.0+



## Example Usage

### Create a AnalyticDB for PostgreSQL instance

```
data "alicloud_zones" "default" {
  available_resource_creation = "Gpdb"
}

resource "alicloud_vpc" "default" {
  name       = "vpc-123456"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  zone_id           = data.alicloud_zones.default.zones[0].id
  vpc_id            = alicloud_vpc.default.id
  cidr_block        = "172.16.0.0/24"
  vswitch_name      = "vpc-123456"
}

resource "alicloud_gpdb_elastic_instance" "adb_pg_instance" {
  engine                   = "gpdb"
  engine_version           = "6.0"
  seg_storage_type         = "cloud_essd"
  seg_node_num             = 4
  storage_size             = 50
  instance_spec            = "2C16G"
  db_instance_description  = "Created by terraform"
  instance_network_type    = "VPC"
  payment_type             = "PayAsYouGo"
  vswitch_id               = alicloud_vswitch.default.id
}

```

## Argument Reference

The following arguments are supported:

* `engine` (Required, ForceNew) Database engine: `gpdb`.
* `engine_version` - (Required, ForceNew) Database version. Valid value is `6.0`.
* `seg_storage_type` - (Required, ForceNew) The disk type of segment nodes. Valid values: `cloud_essd`, `cloud_efficiency`.
* `seg_node_num` - (Required, ForceNew) The number of segment nodes. Minimum is `4`, max is `256`, step is `4`.
* `storage_size` - (Required, ForceNew) The storage capacity of per segment node. Unit: GB. Minimum is `50`, max is `4000`, step is `50`. 
* `instance_spec` - (Required, ForceNew) The specification of segment nodes. 
   * When `db_instance_category` is `HighAvailability`, Valid values: `2C16G`, `4C32G`, `16C128G`.
   * When `db_instance_category` is `Basic`, Valid values: `2C8G`, `4C16G`, `8C32G`, `16C64G`.
* `db_instance_category` - (Optional, ForceNew, Available in v1.158.0+) The edition of the instance. Valid values: `Basic`, `HighAvailability`. Default value: `HighAvailability`.
* `db_instance_description` - (Optional) The description of ADB PG instance. It is a string of 2 to 256 characters.
* `instance_network_type` - (Optional, ForceNew) The network type of ADB PG instance. Only `VPC` supported now.
* `payment_type` - (Optional, ForceNew) Valid values are `PayAsYouGo`, `Subscription`. Default to `PayAsYouGo`.
* `payment_duration_unit` - (Optional) The unit of the subscription period. Valid values: `Month`, `Year`. It is valid when payment_type is `Subscription`.  
  **NOTE:** Will not take effect after modifying `payment_duration_unit` for now, if you want to renew a PayAsYouGo instance, need to do in on aliyun console.
* `payment_duration` - (Optional) The subscription period. Valid values: [1~12]. It is valid when payment_type is `Subscription`.  
  **NOTE:** Will not take effect after modifying `payment_duration` for now, if you want to renew a PayAsYouGo instance, need to do in on aliyun console.
* `zone_id` - (Optional, ForceNew) The Zone to launch the ADB PG instance. If specified, must be consistent with the zone where the vswitch is located.
* `vswitch_id` - (Required, ForceNew) The virtual switch ID to launch ADB PG instances in one VPC.
* `security_ip_list` - (Optional) List of IP addresses allowed to access all databases of an instance. The list contains up to 1,000 IP addresses, separated by commas. Supported formats include 0.0.0.0/0, 10.23.12.24 (IP), and 10.23.12.24/24 (Classless Inter-Domain Routing (CIDR) mode. /24 represents the length of the prefix in an IP address. The range of the prefix length is [1,32]).
* `tags` - (Optional, Available in v1.158.0+) A mapping of tags to assign to the resource.
* `encryption_key` - (Optional, ForceNew, Available in v1.158.0+) The ID of the encryption key. **Note:** If the `encryption_type` parameter is set to `CloudDisk`, you must specify this parameter to the encryption key that is in the same region as the disk that is specified by the EncryptionType parameter. Otherwise, leave this parameter empty.
* `encryption_type` - (Optional, ForceNew, Available in v1.158.0+)  The type of the encryption. Valid values: `CloudDisk`. **Note:** Disk encryption cannot be disabled after it is enabled.


### Timeouts

-> **NOTE:** Available in 1.127.0+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 mins) Used when creating the ADB PG instance (until it reaches the initial `Running` status). 
* `update` - (Defaults to 30 mins) Used when updating the ADB PG instance (until it reaches the initial `Running` status).
* `delete` - (Defaults to 10 mins) Used when terminating the ADB PG instance.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Instance.
* `connection_string` - ADB PG instance connection string.
* `status` - Instance status.

## Import

AnalyticDB for PostgreSQL can be imported using the id, e.g.

```
$ terraform import alicloud_gpdb_elastic_instance.adb_pg_instance gp-bpxxxxxxxxxxxxxx
```
