---
subcategory: "AnalyticDB for PostgreSQL (GPDB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_gpdb_instance"
description: |-
  Provides a Alicloud GPDB D B Instance resource.
---

# alicloud_gpdb_instance

Provides a GPDB D B Instance resource.

AnalyticDB for PostgreSQL database instances.

For information about GPDB D B Instance and how to use it, see [What is D B Instance](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.226.0.

## Example Usage

Basic Usage

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

### Deleting `alicloud_gpdb_instance` or removing it from your configuration

The `alicloud_gpdb_instance` resource allows you to manage  `payment_type = "Subscription"`  instance, but Terraform cannot destroy it.
Deleting the subscription resource or removing it from your configuration will remove it from your state file and management, but will not destroy the Instance.
You can resume managing the subscription instance via the AlibabaCloud Console.

## Argument Reference

The following arguments are supported:
* `connection_string` - (Optional, Available since v1.231.0) The endpoint of the instance.
* `db_instance_category` - (Optional, ForceNew) Instance series. Description of value:
  - `HighAvailability`: High availability.
  - `Basic`: Basic Edition.

-> **NOTE:**  This parameter must be passed in to create a storage elastic mode instance.

* `db_instance_ip_array_attribute` - (Optional, Available since v1.231.0) The default is empty. To distinguish between different attribute values, the console does not display groups with the 'hidden' attribute.
* `db_instance_ip_array_name` - (Optional, Available since v1.231.0) The name of the IP whitelist group.
* `data_share_status` - (Optional, Available since v1.231.0) The status of data sharing. Value Description:
  - `opening`: opening.
  - `opened`: opened.
  - `closing`: closing.
  - `closed`: closed.
* `db_instance_mode` - (Required, ForceNew) Instance resource type. Valid values:
  - `StorageElastic`: storage elastic mode.
  - `Serverless`:Serverless version.
  - `Classic`: Storage reservation mode.

-> **NOTE:**  This parameter is required.

* `description` - (Optional) The description of the instance.
* `encryption_key` - (Optional, ForceNew) The key ID.

-> **NOTE:**  If the value of the `EncryptionType` parameter is `CloudDisk`, you must use this parameter to specify the ID of the encryption key in the same region. Otherwise, it is null.

* `encryption_type` - (Optional, ForceNew) Encryption type, value description:
  - `Off`: does not enable encryption (default).
  - `CloudDisk`: Enable cloud disk encryption and specify the key by using the `EncryptionKey` parameter.

-> **NOTE:**  The current cloud disk encryption cannot be turned off after it is turned on.

* `engine_version` - (Required, ForceNew) Engine version. The value is as follows:

  6.0:6.0 edition.

  7.0:7.0 edition.
* `idle_time` - (Optional, ForceNew, Int, Available since v1.231.0) 空闲释放等待时长 ->
        即当无业务流量的时长达到指定时长后，实例转为空闲状态。单位为秒，最小值为 60，默认值为 600。
* `instance_network_type` - (Optional, ForceNew) The network type of the instance. The value is VPC.

  Description

  Public cloud only supports VPC networks.

  If this parameter is not specified, the default value is the VPC type.
* `instance_spec` - (Optional) Compute node specifications.

  The storage elastic mode high availability edition value is as follows:

  2C16G

  4C32G

  16C128G

  The value of the storage elastic mode base version is as follows:

  2C8G

  4C16G:

  8C32G:

  16C64G

  The value of the Serverless mode is as follows:

  4C16G

  8C32G

  Description

  This parameter must be set when creating a storage elastic mode instance and a Serverless mode instance.
* `maintain_end_time` - (Optional) The end time of the maintenance window for the instance.
* `maintain_start_time` - (Optional) The start time of the maintenance window for the instance.
* `master_cu` - (Optional, Int) Used to describe master node specifications
* `modify_mode` - (Optional, Available since v1.231.0) The IP address whitelist modification mode. The value is as follows:
  - `0` (default): overwrites the original IP address in the target IP whitelist group.
  - `1`: adds an IP address to the destination IP whitelist group.
  - `2`: deletes an IP address from the destination IP whitelist group.
* `payment_type` - (Optional, Computed) The billing method of the instance.
* `period` - (Optional) The unit of time for purchasing resources. The values are as follows:
  - `Month`: Month
  - `Year`: Year

-> **NOTE:**  This parameter must be passed in when creating an instance of the subscription billing type.

* `prod_type` - (Optional, ForceNew, Available since v1.231.0) Product type, can be divided into standard version and Economic version
* `resource_group_id` - (Optional, Computed) The ID of the enterprise resource group to which the instance belongs.
* `resource_management_mode` - (Optional) Resource management mode
* `sample_data_status` - (Optional, Available since v1.231.0) The loading status of the sample dataset. Value Description:
  - `loaded`: loaded.
  - `loading`: loading.
  - `unload`: Not loaded.
* `security_ip_list` - (Optional) The IP address whitelist contains a maximum of 1000 IP addresses separated by commas in the following three formats:
  - 0.0.0.0/0
  - 10.23.12.24(IP)
  - 10.23.12.24/24(CIDR mode, Classless Inter-Domain Routing, '/24' indicates the length of the prefix in the address, and the range is '[1,32]')
* `seg_disk_performance_level` - (Optional, Available since v1.231.0) The performance level of the ESSD cloud disk. The value is as follows:
  - `pl0`:PL0 level.
  - `pl1`:PL1 level.
  - `pl2`:PL2 level.

-> **NOTE:** - This parameter takes effect only when the disk storage type is ESSD cloud disk.
  - If not filled, the default is PL1 level.
* `seg_node_num` - (Optional) Calculate the number of nodes. Valid values:
  - The value range of the high-availability version of the storage elastic mode is 4 to 512, and the value must be a multiple of 4.
  - The value range of the basic version of the storage elastic mode is 2 to 512, and the value must be a multiple of 2.

  The-Serverless version has a value range of 2 to 512. The value must be a multiple of 2.

-> **NOTE:** - this parameter must be passed in to create a storage elastic mode instance and a Serverless version instance.
  - During the public beta of the Serverless version (from 0101, 2022 to 0131, 2022), a maximum of 12 compute nodes can be created.
* `seg_storage_type` - (Optional) The disk storage type. Currently, only ESSD disks are supported. Set the value to **cloud_essd * *.

-> **NOTE:**  This parameter must be set when creating a storage elastic mode instance.

* `serverless_mode` - (Optional, ForceNew, Available since v1.231.0) The mode of the Serverless instance. The values are as follows:
  - `Manual`: Manual scheduling, the default value.
  - `Auto`: automatic scheduling.

-> **NOTE:**  This parameter is required for only Serverless mode instances.

* `serverless_resource` - (Optional, ForceNew, Int, Available since v1.231.0) Calculate resource thresholds. The value range is 8 to 32, the step size is 8, the unit is ACU. The default value is 32.

-> **NOTE:**  This parameter is required for only Serverless automatic scheduling mode instances.

* `ssl_enabled` - (Optional, Int) ssl status
* `status` - (Optional, Computed, Available since v1.231.0) The status of the resource
* `storage_size` - (Optional, Computed, Int) 

  The storage space size, in GB, in the range of 50 to 4000.

  Description

  This parameter must be set when creating a storage elastic mode instance.
* `tags` - (Optional, Map) The tags of the instance.
* `upgrade_type` - (Optional, Int, Available since v1.231.0) The type of instance type change. Valid values:
  - `0` (default): changes the number of Segment nodes.
  - `1`: Change the size of the Segment node and the size of the storage space.
  - `2`: Change the number of Master nodes.

-> **NOTE:** - Different instance resource types have different support levels for computing node configuration changes. For more information, see [Note](~~ 50956 ~~).
  - After the corresponding change type is selected, only the corresponding parameters take effect, and other parameters do not take effect. For example, if the `UpgradeType` parameter is 0, the parameter that changes the number of Segment nodes and the number of Master nodes is passed in at the same time, only the parameter that changes the number of Segment nodes takes effect.
  - Only the China site supports changing the number of Master nodes.
* `used_time` - (Optional) Length of purchase of resources. The values are as follows:
  - When `Period` is `Month`, the value is 1 to 9.
  - When `Period` is `Year`, the value is 1 to 3.

-> **NOTE:**  This parameter must be passed in when creating an instance of the subscription billing type.

* `vswitch_id` - (Optional, ForceNew) The ID of the VSwitch.
* `vector_configuration_status` - (Optional) Whether to enable vector engine optimization. Value Description:
  - `enabled`: enables vector engine optimization.
  - `disabled` (default): does not enable vector engine optimization.

-> **NOTE:** - For mainstream analysis scenarios, data warehouse scenarios, and real-time data warehouse scenarios, we recommend that you **do not enable** vector engine optimization.
  - For users who use the vector analysis engine for scenarios such as AIGC and vector retrieval, we recommend that you `enable` vector engine optimization.
* `vpc_id` - (Optional, ForceNew) The ID of the VPC。
* `zone_id` - (Required, ForceNew) The zone ID of the instance.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The time when the instance was created. The time is in the YYYY-MM-DDThh:mm:ssZ format, such as 2011-05-30T12:11:4Z.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the D B Instance.
* `delete` - (Defaults to 5 mins) Used when delete the D B Instance.
* `update` - (Defaults to 5 mins) Used when update the D B Instance.

## Import

GPDB D B Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_gpdb_instance.example <id>
```