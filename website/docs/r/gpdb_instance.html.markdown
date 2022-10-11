---
subcategory: "AnalyticDB for PostgreSQL (GPDB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_gpdb_instance"
sidebar_current: "docs-alicloud-resource-gpdb-instance"
description: |-
  Provides a AnalyticDB for PostgreSQL instance resource.
---

# alicloud\_gpdb\_instance

Provides a AnalyticDB for PostgreSQL instance resource supports replica set instances only. the AnalyticDB for PostgreSQL provides stable, reliable, and automatic scalable database services.
You can see detail product introduction [here](https://www.alibabacloud.com/help/doc-detail/35387.htm)

-> **NOTE:**  Available in 1.47.0+

## Example Usage

Basic Usage

```terraform
data "alicloud_resource_manager_resource_groups" "default" {}
data "alicloud_gpdb_zones" "default" {}
data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_gpdb_zones.default.ids.0
}
resource "alicloud_vswitch" "vswitch" {
  count        = length(data.alicloud_vswitches.default.ids) > 0 ? 0 : 1
  vpc_id       = data.alicloud_vpcs.default.ids.0
  cidr_block   = cidrsubnet(data.alicloud_vpcs.default.vpcs[0].cidr_block, 8, 8)
  zone_id      = data.alicloud_gpdb_zones.default.ids.0
  vswitch_name = var.name
}
locals {
  vswitch_id = length(data.alicloud_vswitches.default.ids) > 0 ? data.alicloud_vswitches.default.ids[0] : concat(alicloud_vswitch.vswitch.*.id, [""])[0]
}
resource "alicloud_gpdb_instance" "example" {
  db_instance_category  = "HighAvailability"
  db_instance_class     = "gpdb.group.segsdx1"
  db_instance_mode      = "StorageElastic"
  description           = "example_value"
  engine                = "gpdb"
  engine_version        = "6.0"
  zone_id               = data.alicloud_gpdb_zones.default.ids.0
  instance_network_type = "VPC"
  instance_spec         = "2C16G"
  master_node_num       = 1
  payment_type          = "PayAsYouGo"
  private_ip_address    = "1.1.1.1"
  seg_storage_type      = "cloud_essd"
  seg_node_num          = 4
  storage_size          = 50
  vpc_id                = data.alicloud_vpcs.default.ids.0
  vswitch_id            = local.vswitch_id
  ip_whitelist {
    security_ip_list = "127.0.0.1"
  }
}

```

## Argument Reference

The following arguments are supported:

* `db_instance_category` - (Optional, ForceNew) The db instance category. Valid values: `HighAvailability`, `Basic`.
-> **NOTE:** This parameter must be passed in to create a storage reservation mode instance.

* `db_instance_class` - (Optional, ForceNew) The db instance class. see [Instance specifications](https://www.alibabacloud.com/help/doc-detail/86942.htm).
-> **NOTE:** This parameter must be passed in to create a storage reservation mode instance.

* `db_instance_mode` - (Required, ForceNew) The db instance mode. Valid values: `StorageElastic`, `Serverless`, `Classic`.
* `description` - (Optional) The description of the instance.
* `engine` - (Required, ForceNew) The database engine used by the instance. Value options can refer to the latest docs [CreateDBInstance](https://www.alibabacloud.com/help/doc-detail/86908.htm) `EngineVersion`.
* `engine_version` - (Required, ForceNew) The version of the database engine used by the instance.
* `instance_network_type` - (Optional, ForceNew) The network type of the instance.
* `instance_spec` - (Optional) The specification of segment nodes.
  * When `db_instance_category` is `HighAvailability`, Valid values: `2C16G`, `4C32G`, `16C128G`.
  * When `db_instance_category` is `Basic`, Valid values: `2C8G`, `4C16G`, `8C32G`, `16C64G`.
  * When `db_instance_category` is `Serverless`, Valid values: `4C16G`, `8C32G`.
-> **NOTE:** This parameter must be passed to create a storage elastic mode instance and a serverless version instance.

* `ip_whitelist` - (Optional) The ip whitelist.
* `security_ip_list` - (Optional) Field `security_ip_list` has been deprecated from provider version 1.187.0. New field `ip_whitelist` instead.
* `maintain_end_time` - (Optional) The end time of the maintenance window for the instance. in the format of HH:mmZ (UTC time), for example 03:00Z. start time should be later than end time.
* `maintain_start_time` - (Optional) The start time of the maintenance window for the instance. in the format of HH:mmZ (UTC time), for example 02:00Z.
* `master_node_num` - (Optional) The number of Master nodes. Valid values: 1 to 2. if it is not filled in, the default value is 1 Master node.
* `instance_group_count` - (Optional, ForceNew) The number of nodes. Valid values: `2`, `4`, `8`, `12`, `16`, `24`, `32`, `64`, `96`, `128`.
* `payment_type` - (Optional, ForceNew) The billing method of the instance. Valid values: `Subscription`, `PayAsYouGo`.
* `instance_charge_type` - (Optional, ForceNew, Deprecated) Field `instance_charge_type` has been deprecated from provider version 1.187.0. New field `payment_type` instead.
* `period` - (Optional) The duration that you will buy the resource, in month. required when `payment_type` is `Subscription`. Valid values: `Year`, `Month`.
* `private_ip_address` - (Optional) The private ip address.
* `resource_group_id` - (Optional) The ID of the enterprise resource group to which the instance belongs.
* `seg_node_num` - (Optional, Computed) Calculate the number of nodes. The value range of the high-availability version of the storage elastic mode is 4 to 512, and the value must be a multiple of 4. The value range of the basic version of the storage elastic mode is 2 to 512, and the value must be a multiple of 2. The-Serverless version has a value range of 2 to 512. The value must be a multiple of 2.
-> **NOTE:** This parameter must be passed in to create a storage elastic mode instance and a Serverless version instance. During the public beta of the Serverless version (from 0101, 2022 to 0131, 2022), a maximum of 12 compute nodes can be created.

* `seg_storage_type` - (Optional) The seg storage type. Valid values: `cloud_essd`, `cloud_efficiency`.
-> **NOTE:** This parameter must be passed in to create a storage elastic mode instance. Storage Elastic Mode Basic Edition instances only support ESSD cloud disks.

* `storage_size` - (Optional, Computed) The storage capacity. Unit: GB. Value: `50` to `4000`.
-> **NOTE:** This parameter must be passed in to create a storage reservation mode instance.

* `used_time` - (Optional) The used time. When the parameter `period` is `Year`, the `used_time` value is 1 to 3. When the parameter `period` is `Month`, the `used_time` value is 1 to 9.
* `vswitch_id` - (Required, ForceNew) The vswitch id.
* `create_sample_data` - (Optional, Computed) Whether to load the sample dataset after the instance is created. Valid values: `true`, `false`.
* `vpc_id` - (Optional, Computed, ForceNew) The vpc ID of the resource.
* `zone_id` - (Optional, ForceNew) The zone ID of the instance.
* `availability_zone` - (Optional, ForceNew, Deprecated) Field `availability_zone` has been deprecated from provider version 1.187.0. New field `zone_id` instead.
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `ssl_enabled` - (Optional, Computed, Available in v1.188.0+) Enable or disable SSL. Valid values: `0` and `1`.

#### Block ip_whitelist

The ip_whitelist supports the following:

* `ip_group_attribute` - (Optional) The value of this parameter is empty by default. The attribute of the whitelist group. The console does not display the whitelist group whose value of this parameter is hidden.
* `ip_group_name` - (Optional) IP whitelist group name
* `security_ip_list` - (Required) List of IP addresses allowed to access all databases of an instance. The list contains up to 1,000 IP addresses, separated by commas. Supported formats include 0.0.0.0/0, 10.23.12.24 (IP), and 10.23.12.24/24 (Classless Inter-Domain Routing (CIDR) mode. /24 represents the length of the prefix in an IP address. The range of the prefix length is [1,32]). System default to `["127.0.0.1"]`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of AnalyticDB for PostgreSQL.
* `status` - The status of the instance.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 mins) Used when create the DB Instance.
* `update` - (Defaults to 60 mins) Used when update the DB Instance.

## Import

AnalyticDB for PostgreSQL can be imported using the id, e.g.

```
$ terraform import alicloud_gpdb_instance.example <id>
```