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
* `master_node_num` - (Optional, ForceNew, Available in 1.134.0+) The number of master nodes. Valid values: [1~16]. Default value: `1`.
* `seg_storage_type` - (Required, ForceNew) The disk type of segment nodes. Valid values: `cloud_essd`, `cloud_efficiency`.
* `seg_node_num` - (Required, ForceNew) The number of segment nodes. Minimum is `4`, max is `256`, step is `4`.
* `storage_size` - (Required, ForceNew) The storage capacity of per segment node. Unit: GB. Minimum is `50`, max is `4000`, step is `50`. 
* `instance_spec` - (Required, ForceNew) The specification of segment nodes. Valid values: `2C16G`, `4C32G`, `16C128G`.
* `db_instance_description` - (Optional) The description of ADB PG instance. It is a string of 2 to 256 characters.
* `instance_network_type` - (Optional, ForceNew) The network type of ADB PG instance. Only `VPC` supported now.
* `payment_type` - (Optional, ForceNew) Valid values are `PayAsYouGo`, `Subscription`. Default to `PayAsYouGo`.
* `payment_duration_unit` - (Optional) The unit of the subscription period. Valid values: `Month`, `Year`. It is valid when payment_type is `Subscription`.  
  **NOTE:** Will not take effect after modifying `payment_duration_unit` for now, if you want to renew a PayAsYouGo instance, need to do in on aliyun console.
* `payment_duration` - (Optional) The subscription period. Valid values: [1~12]. It is valid when payment_type is `Subscription`.  
  **NOTE:** Will not take effect after modifying `payment_duration` for now, if you want to renew a PayAsYouGo instance, need to do in on aliyun console.
* `encryption_key` - (Optional, ForceNew, Available in 1.134.0+) If the `encryption_type` parameter is set to `CloudDisk`, you must specify this parameter to the encryption key that is in the same region with the disks that is specified by the `encryption_type` parameter. Otherwise, leave this parameter empty.
* `encryption_type` - (Optional, ForceNew, Available in 1.134.0+) The type of the encryption. Default value: `Off`. Valid values:   
    - Off: Encryption is disabled.
    - CloudDisk: Encryption is enabled on disks and the encryption key is specified by using the `encryption_type` parameter.  
  **NOTE:** Disk encryption cannot be disabled after it is enabled.
* `zone_id` - (Optional, ForceNew) The Zone to launch the ADB PG instance. If specified, must be consistent with the zone where the vswitch is located.
* `vswitch_id` - (Required, ForceNew) The virtual switch ID to launch ADB PG instances in one VPC.
* `security_ip_list` - (Optional) List of IP addresses allowed to access all databases of an instance. The list contains up to 1,000 IP addresses, separated by commas. Supported formats include 0.0.0.0/0, 10.23.12.24 (IP), and 10.23.12.24/24 (Classless Inter-Domain Routing (CIDR) mode. /24 represents the length of the prefix in an IP address. The range of the prefix length is [1,32]).
* `preferred_backup_period` - (Optional, Available in 1.134.0+) The cycle based on which you want to perform a backup. Separate multiple values with commas (,). Valid values: [Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday].
* `preferred_backup_time` - (Optional, Available in 1.134.0+) The backup window. Specify the backup window in the HH:mmZ-HH:mmZ format. The backup window must be in UTC.
* `backup_retention_period` - (Optional, Available in 1.134.0+) The number of days for which data backup files are retained. Valid values: [1~7]. Default value: `7`.
* `enable_recovery_point` - (Optional, Available in 1.134.0+) Specifies whether to enable automatic point-in-time backup. Valid values: `true`, `false`. Default value: `true`.
* `recovery_point_period` - (Optional, Available in 1.134.0+) The frequency of point-in-time backup. Valid values: `1`, `2`, `4`, `8`. Default value: `8`. Valid when the `enable_recovery_point` is `true`.
* `tags` - (Optional, Available in 1.134.0+) A mapping of tags to assign to the instance.
* `parameters` - (Optional, Available in 1.134.0+) The parameters of ADB PG instance. Available parameters can refer to the latest docs [View database parameter templates](https://www.alibabacloud.com/help/doc-detail/208310.htm).
* `force_restart_instance` - (Optional, Available in 1.134.0+) Specifies whether to forcibly restart the instance after parameters modified. Valid values: `true`, `false`.
* `ssl_enabled` - (Optional, Available in 1.134.0+) The status of SSL encryption. Valid values:  
    - 0: disables SSL encryption.
    - 1: enables SSL encryption.
    - 2: updates SSL encryption.


### Timeouts

-> **NOTE:** Available in 1.127.0+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 50 mins) Used when creating the ADB PG instance (until it reaches the initial `Running` status). 
* `update` - (Defaults to 30 mins) Used when updating the ADB PG instance (until it reaches the initial `Running` status).
* `delete` - (Defaults to 10 mins) Used when terminating the ADB PG instance.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Instance.
* `connection_string` - ADB PG instance connection string.
* `status` - Instance status.
* `ssl_expired_time` - (Available in 1.134.0+) The expiration time of the SSL certificate.

## Import

AnalyticDB for PostgreSQL can be imported using the id, e.g.

```
$ terraform import alicloud_gpdb_elastic_instance.adb_pg_instance gp-bpxxxxxxxxxxxxxx
```
