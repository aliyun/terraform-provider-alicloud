---
subcategory: "Click House"
layout: "alicloud"
page_title: "Alicloud: alicloud_click_house_db_cluster"
description: |-
  Provides a Alicloud Click House D B Cluster resource.
---

# alicloud_click_house_db_cluster

Provides a Click House D B Cluster resource. The ClickHouse instance integrates computing, storage, and network resources, and provides management functions such as database users, databases, tables, networks, and backups. You can easily customize and change the configuration of an instance.

For information about Click House D B Cluster and how to use it, see [What is D B Cluster](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.218.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_click_house_db_cluster" "default" {
  category            = "Basic"
  storage_type        = "CloudESSD_PL2"
  resource_group_id   = "rg-acfmzygvt54v5pq"
  zone_id             = "cn-hangzhou-h"
  db_node_group_count = "1"
  vswitch_id          = "vsw-bp12zcpn6zyagqr8lbrrg"
  db_cluster_version  = "21.8.10.19"
  db_node_storage     = "2000"
  db_cluster_name     = var.name

  vpc_id                  = "vpc-bp1hmr2u26m8ewyt2banw"
  db_cluster_network_type = "vpc"
  used_time               = "1"
  payment_type            = "Subscription"
  db_cluster_class        = "S8"
  encryption_type         = "CloudDisk"
  period                  = "Month"
  encryption_key          = "39875ecc-c89e-4731-aa11-2ff15c876ad1"
}
```

## Argument Reference

The following arguments are supported:
* `category` - (Required, ForceNew, Available since v1.134.0) Copy configuration, value description:  Basic: Single copy version. HighAvailability: Double Edition.
* `db_cluster_class` - (Required, Available since v1.134.0) Cluster specifications.

  - Single copy version, value:
  - **LS20**: large storage type 20 core 88GB.
  - **LS40**: Large storage 40-core 176GB.
  - **LS80**: large storage 80-core 352GB.
  - **S8**: Standard 8 core 32GB.
  - **S16**: standard 16-core 64GB.
  - **S32**: Standard 32 core 128GB.
  - **S64**: Standard 64-core 256GB.
  - **S80**: standard 80 core 384GB.
  - **S104**: Standard 104 core 384GB.
  - Double versions, value:
  - **LC20**: large storage type 20 core 88GB.
  - **LC40**: Large storage 40-core 176GB.
  - **LC80**: large storage 80-core 352GB.
  - **C8**: Standard 8 core 32GB.
  - **C16**: standard 16-core 64GB.
  - **C32**: Standard 32 core 128GB.
  - **C64**: Standard 64-core 256GB.
  - **C80**: standard 80 core 384GB.
  - **C104**: Standard 104 core 384GB.



  - Single copy version, value:
  - **S4**:4 cores 16GB.
  - **S8**:8 core 32GB.
  - **S16**:16-core 64GB.
  - **S32**:32 core 128GB.
  - **S64** :64 core 256GB.
  - **S104**:104 core 384GB.
  - Double versions, value:
  - **C4**:4 core 16GB.
  - **C8**:8 core 32GB.
  - **C16**:16-core 64GB.
  - **C32**:32 core 128GB.
  - **C64** :64-core 256GB.
  - **C104**:104 core 384GB.

.
* `db_cluster_ip_array_name` - (Optional) The name of the whitelist group to be modified.
-> **NOTE:**  If this parameter is not configured, the whitelist in the default group is modified by default.
* `db_cluster_name` - (Optional) The cluster description information.
* `db_cluster_network_type` - (Required, ForceNew, Available since v1.134.0) Network type. Currently, only VPC is supported.
* `db_cluster_version` - (Required, Available since v1.134.0) Kernel version, value:  21.8.10.19  22.8.5.29.
* `db_node_group_count` - (Required, Available since v1.134.0) Number of nodes.
  - Single version, value range: 1~48.
  - Double versions, value range: 1~24.
* `db_node_storage` - (Required, ForceNew, Available since v1.134.0) Single-node storage space. Value range: 100 GB to 32000GB.
* `duration` - (Optional) Set the instance auto-renewal duration.
* `encryption_key` - (Optional, ForceNew, Available since v1.134.0) Key management service KMS key ID.
* `encryption_type` - (Optional, ForceNew, Available since v1.134.0) Currently only supports ECS disk encryption, with a value of CloudDisk, not encrypted when empty.
* `maintain_time` - (Optional, Available since v1.134.0) Examples of the maintenance window, in the format of hh:mmZ-hh:mm Z.
* `modify_mode` - (Optional) Modification method, value description:
  - **Cover**: overwrites the original whitelist.
  - **Append**: Adds a modified whitelist to the original whitelist.
  - **Delete**: deletes the original whitelist.
-> **NOTE:**  If this parameter is not configured, the whitelist is modified by the Cover method by default.
* `payment_type` - (Required, ForceNew, Available since v1.134.0) The paymen type of the resource.
* `period` - (Optional, Available since v1.134.0) Pre-paid cluster of the pay-as-you-go cycle.
* `period_unit` - (Optional) Unit of renewal duration.
* `renewal_status` - (Optional) Set automatic renewal status.
* `resource_group_id` - (Optional, Computed) The ID of the resource group.
* `restart_time` - (Optional) Appointment of restart time. Format: yyyy-MM-ddTHH:mmZ(UTC time).
-> **NOTE:**  When the parameter is empty or filled in earlier than the current time, the cluster will restart immediately.
* `security_ips` - (Optional) The whitelist supports the following two formats:
  - IP address format, for example, 192.168.0.1 indicates that this IP address is allowed to access the cloud database ClickHouse.
  - IP segment format: for example, 192.168.0.0/24 indicates that IP addresses from 192.168.0.1 to 192.168.0.255 are allowed to access the cloud database ClickHouse.
-> **NOTE:** - Prohibit entering IP: 0.0.0.0.
  - Set to 127.0.0.1 to disable all address access.
* `storage_type` - (Required, Available since v1.134.0) Storage type CloudSSD:SSD cloud disk CloudEfficiency: Ultra cloud disk.
* `used_time` - (Optional, Available since v1.134.0) This parameter takes effect only when the payment type is Prepaid and is required.  The purchase duration of the prepaid cluster.  When Period is set to Year, the value ranges from 1 to 3 (integer). When Period is Month, the value ranges from 1 to 9 (integer).
* `vswitch_id` - (Required, ForceNew, Available since v1.134.0) Switch ID.
* `vswitch_id_bak` - (Optional) Standby VPC switch .
* `vswitch_id_bak_second` - (Optional) VSwitchIdBak2.
* `vpc_id` - (Optional, ForceNew, Available since v1.134.0) VPC ID.
* `zone_id` - (Optional, ForceNew, Available since v1.134.0) On behalf of the zone resource attribute field.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the D B Cluster.
* `delete` - (Defaults to 5 mins) Used when delete the D B Cluster.
* `update` - (Defaults to 5 mins) Used when update the D B Cluster.

## Import

Click House D B Cluster can be imported using the id, e.g.

```shell
$ terraform import alicloud_click_house_db_cluster.example <id>
```