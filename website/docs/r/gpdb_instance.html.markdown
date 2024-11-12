---
subcategory: "AnalyticDB for PostgreSQL (GPDB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_gpdb_instance"
sidebar_current: "docs-alicloud-resource-gpdb-instance"
description: |-
  Provides a AnalyticDB for PostgreSQL instance resource.
---

# alicloud_gpdb_instance

Provides a AnalyticDB for PostgreSQL instance resource supports replica set instances only. the AnalyticDB for PostgreSQL provides stable, reliable, and automatic scalable database services.
You can see detail product introduction [here](https://www.alibabacloud.com/help/en/analyticdb-for-postgresql/latest/api-gpdb-2016-05-03-createdbinstance)

-> **NOTE:** Available since v1.47.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_gpdb_instance&exampleId=beac6fb2-7bb3-c67b-00d3-4e76327bc40e0120436d&activeTab=example&spm=docs.r.gpdb_instance.0.beac6fb27b&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}

data "alicloud_resource_manager_resource_groups" "default" {
}

data "alicloud_gpdb_zones" "default" {
}

data "alicloud_vpcs" "default" {
  # You need to modify name_regex to an existing VPC under your account
  name_regex = "^default-NODELETING$"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_gpdb_zones.default.ids.0
}

resource "alicloud_gpdb_instance" "default" {
  db_instance_category  = "HighAvailability"
  db_instance_class     = "gpdb.group.segsdx1"
  db_instance_mode      = "StorageElastic"
  description           = var.name
  engine                = "gpdb"
  engine_version        = "6.0"
  zone_id               = data.alicloud_gpdb_zones.default.ids.0
  instance_network_type = "VPC"
  instance_spec         = "2C16G"
  payment_type          = "PayAsYouGo"
  seg_storage_type      = "cloud_essd"
  seg_node_num          = 4
  storage_size          = 50
  vpc_id                = data.alicloud_vpcs.default.ids.0
  vswitch_id            = data.alicloud_vswitches.default.ids.0
  ip_whitelist {
    security_ip_list = "127.0.0.1"
  }
}
```

## Argument Reference

The following arguments are supported:

* `engine` - (Required, ForceNew) The database engine used by the instance. Value options can refer to the latest docs [CreateDBInstance](https://www.alibabacloud.com/help/en/analyticdb-for-postgresql/latest/api-gpdb-2016-05-03-createdbinstance) `EngineVersion`.
* `engine_version` - (Required, ForceNew) The version of the database engine used by the instance.
* `vswitch_id` - (Required, ForceNew) The vswitch id.
* `db_instance_class` - (Optional, ForceNew) The db instance class. see [Instance specifications](https://www.alibabacloud.com/help/en/analyticdb-for-postgresql/latest/instance-types).
-> **NOTE:** This parameter must be passed in to create a storage reservation mode instance.
* `db_instance_category` - (Optional, ForceNew) The db instance category. Valid values: `Basic`, `HighAvailability`.
-> **NOTE:** This parameter must be passed in to create a storage reservation mode instance.
* `db_instance_mode` - (Required, ForceNew) The db instance mode. Valid values: `StorageElastic`, `Serverless`, `Classic`.
* `instance_spec` - (Optional) The specification of segment nodes.
  * When `db_instance_category` is `HighAvailability`, Valid values: `2C16G`, `4C32G`, `16C128G`.
  * When `db_instance_category` is `Basic`, Valid values: `2C8G`, `4C16G`, `8C32G`, `16C64G`.
  * When `db_instance_category` is `Serverless`, Valid values: `4C16G`, `8C32G`.
-> **NOTE:** This parameter must be passed to create a storage elastic mode instance and a serverless version instance.
* `storage_size` - (Optional, Int) The storage capacity. Unit: GB. Valid values: `50` to `4000`.
-> **NOTE:** This parameter must be passed in to create a storage reservation mode instance.
* `instance_network_type` - (Optional, ForceNew) The network type of the instance. Valid values: `VPC`.
* `vpc_id` - (Optional, ForceNew) The vpc ID of the resource.
* `zone_id` - (Optional, ForceNew) The zone ID of the instance.
* `instance_group_count` - (Optional, ForceNew, Int) The number of nodes. Valid values: `2`, `4`, `8`, `12`, `16`, `24`, `32`, `64`, `96`, `128`.
* `payment_type` - (Optional, ForceNew) The billing method of the instance. Valid values: `Subscription`, `PayAsYouGo`.
* `period` - (Optional) The duration that you will buy the resource, in month. required when `payment_type` is `Subscription`. Valid values: `Year`, `Month`.
* `resource_group_id` - (Optional) The ID of the enterprise resource group to which the instance belongs.
* `master_cu` - (Optional, Int, Available since v1.213.0) The amount of coordinator node resources. Valid values: `2`, `4`, `8`, `16`, `32`.
* `seg_node_num` - (Optional, Int) Calculate the number of nodes. Valid values: `2` to `512`. The value range of the high-availability version of the storage elastic mode is `4` to `512`, and the value must be a multiple of `4`. The value range of the basic version of the storage elastic mode is `2` to `512`, and the value must be a multiple of `2`. The-Serverless version has a value range of `2` to `512`. The value must be a multiple of `2`.
-> **NOTE:** This parameter must be passed in to create a storage elastic mode instance and a Serverless version instance. During the public beta of the Serverless version (from 0101, 2022 to 0131, 2022), a maximum of 12 compute nodes can be created.
* `seg_storage_type` - (Optional, ForceNew) The seg storage type. Valid values: `cloud_essd`. **NOTE:** If `db_instance_mode` is set to `StorageElastic`, `seg_storage_type` is required. From version 1.233.1, `seg_storage_type` cannot be modified, or set to `cloud_efficiency`. `seg_storage_type` can only be set to `cloud_essd`.
* `seg_disk_performance_level` - (Optional, Available since v1.233.1) The ESSD cloud disk performance level. Valid values: `pl0`, `pl1`, `pl2`.
* `create_sample_data` - (Optional, Bool) Whether to load the sample dataset after the instance is created. Valid values: `true`, `false`.
* `ssl_enabled` - (Optional, Int, Available since v1.188.0) Enable or disable SSL. Valid values: `0` and `1`.
* `encryption_type` - (Optional, ForceNew, Available since v1.207.2) The encryption type. Valid values: `CloudDisk`.
-> **NOTE:** Disk encryption cannot be disabled after it is enabled.
* `encryption_key` - (Optional, ForceNew, Available since v1.207.2) The ID of the encryption key.
-> **NOTE:** If `encryption_type` is set to `CloudDisk`, you must specify an encryption key that resides in the same region as the cloud disk that is specified by EncryptionType. Otherwise, leave this parameter empty.
* `vector_configuration_status` - (Optional, Available since v1.207.2) Specifies whether to enable vector engine optimization. Default value: `disabled`. Valid values: `enabled` and `disabled`.
* `maintain_start_time` - (Optional) The start time of the maintenance window for the instance. in the format of HH:mmZ (UTC time), for example 02:00Z.
* `maintain_end_time` - (Optional) The end time of the maintenance window for the instance. in the format of HH:mmZ (UTC time), for example 03:00Z. start time should be later than end time.
* `resource_management_mode` - (Optional, Available since v1.225.0) Resource management mode. Valid values: `resourceGroup`, `resourceQueue`.
* `serverless_mode` - (Optional, ForceNew, Available since v1.233.1) The mode of the Serverless instance. Valid values: `Manual`, `Auto`. **NOTE:** `serverless_mode` is valid only when `db_instance_mode` is set to `Serverless`.
* `prod_type` - (Optional, ForceNew, Available since v1.233.1) The type of the product. Default value: `standard`. Valid values: `standard`, `cost-effective`.
* `data_share_status` - (Optional, Available since v1.233.1) Specifies whether to enable or disable data sharing. Default value: `closed`. Valid values:
  - `opened`: Enables data sharing.
  - `closed`: Disables data sharing.
-> **NOTE:** `data_share_status` is valid only when `db_instance_mode` is set to `Serverless`.
* `used_time` - (Optional) The used time. When the parameter `period` is `Year`, the `used_time` value is `1` to `3`. When the parameter `period` is `Month`, the `used_time` value is `1` to `9`.
* `description` - (Optional) The description of the instance.
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `ip_whitelist` - (Optional, Set, Available since v1.187.0) The ip whitelist. See [`ip_whitelist`](#ip_whitelist) below.
  Default to creating a whitelist group with the group name "default" and security_ip_list "127.0.0.1".
* `parameters` - (Optional, Set, Available since v1.231.0) The parameters. See [`parameters`](#parameters) below.
* `security_ip_list` - (Optional, List, Deprecated since v1.187.0) Field `security_ip_list` has been deprecated from provider version 1.187.0. New field `ip_whitelist` instead.
* `instance_charge_type` - (Optional, ForceNew, Deprecated since v1.187.0) Field `instance_charge_type` has been deprecated from provider version 1.187.0. New field `payment_type` instead.
* `availability_zone` - (Optional, ForceNew, Deprecated since v1.187.0) Field `availability_zone` has been deprecated from provider version 1.187.0. New field `zone_id` instead.
* `master_node_num` - (Optional, Int, Deprecated since v1.213.0) The number of Master nodes. **NOTE:** Field `master_node_num` has been deprecated from provider version 1.213.0.
* `private_ip_address` - (Optional, Deprecated since v1.213.0) The private ip address. **NOTE:** Field `private_ip_address` has been deprecated from provider version 1.213.0.


### `ip_whitelist`

The ip_whitelist supports the following:

* `ip_group_attribute` - (Optional) The value of this parameter is empty by default. The attribute of the whitelist group. 
  If the value contains `hidden`, this white list item will not output.
* `ip_group_name` - (Optional) IP whitelist group name.
* `security_ip_list` - (Optional) List of IP addresses allowed to access all databases of an instance. The list contains up to 1,000 IP addresses, separated by commas. Supported formats include 0.0.0.0/0, 10.23.12.24 (IP), and 10.23.12.24/24 (Classless Inter-Domain Routing (CIDR) mode. /24 represents the length of the prefix in an IP address. The range of the prefix length is [1,32]). System default to `["127.0.0.1"]`.

### `parameters`

The parameters supports the following:

* `name` - (Required, Available since v1.231.0) The name of the parameter.
* `value` - (Required, Available since v1.231.0) The value of the parameter.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of AnalyticDB for PostgreSQL.
* `status` - The status of the instance.
* `connection_string` - (Available since v1.196.0) The connection string of the instance.
* `port` - (Available since v1.196.0) The connection port of the instance.
* `parameters` - (Available since v1.231.0) A list of parameters. Each element contains the following attributes:
  * `default_value` - (Available since v1.231.0) The default value of the parameter.
  * `force_restart_instance` - (Available since v1.231.0) Whether to force restart the instance to config the parameter.
  * `parameter_description` - (Available since v1.231.0) The description of the parameter.
  * `optional_range` - (Available since v1.231.0) The optional range of the parameter.
  * `is_changeable_config` - (Available since v1.231.0) Whether the parameter is changeable.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 mins) Used when create the DB Instance.
* `update` - (Defaults to 60 mins) Used when update the DB Instance.
* `delete` - (Defaults to 10 mins) Used when update the DB Instance.

## Import

AnalyticDB for PostgreSQL can be imported using the id, e.g.

```shell
$ terraform import alicloud_gpdb_instance.example <id>
```
