---
subcategory: "Elastic High Performance Computing(ehpc)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ehpc_cluster"
sidebar_current: "docs-alicloud-resource-ehpc-cluster"
description: |-
  Provides a Alicloud Ehpc Cluster resource.
---

# alicloud\_ehpc\_cluster

Provides a Ehpc Cluster resource.

For information about Ehpc Cluster and how to use it, see [What is Cluster](https://www.alibabacloud.com/help/e-hpc/latest/createcluster).

-> **NOTE:** Available in v1.173.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_zones" default {
  available_resource_creation = "VSwitch"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.zones.0.id
}
data "alicloud_instance_types" "default" {
  availability_zone = data.alicloud_zones.default.zones.0.id
}
variable "storage_type" {
  default = "Performance"
}
resource "alicloud_nas_file_system" "default" {
  storage_type  = var.storage_type
  protocol_type = "NFS"
}
resource "alicloud_nas_mount_target" "default" {
  file_system_id    = alicloud_nas_file_system.default.id
  access_group_name = "DEFAULT_VPC_GROUP_NAME"
  vswitch_id        = data.alicloud_vswitches.default.ids.0
}
data "alicloud_images" "default" {
  name_regex = "^centos_7_6_x64*"
  owners     = "system"
}
resource "alicloud_ehpc_cluster" "default" {
  cluster_name          = "example_value"
  deploy_mode           = "Simple"
  description           = "example_description"
  ha_enable             = false
  image_id              = data.alicloud_images.default.images.0.id
  image_owner_alias     = "system"
  volume_protocol       = "nfs"
  volume_id             = alicloud_nas_file_system.default.id
  volume_mountpoint     = alicloud_nas_mount_target.default.mount_target_domain
  compute_count         = 1
  compute_instance_type = data.alicloud_instance_types.default.instance_types.0.id
  login_count           = 1
  login_instance_type   = data.alicloud_instance_types.default.instance_types.0.id
  manager_count         = 1
  manager_instance_type = data.alicloud_instance_types.default.instance_types.0.id
  os_tag                = "CentOS_7.6_64"
  scheduler_type        = "pbs"
  password              = "your-password123"
  vswitch_id            = data.alicloud_vswitches.default.ids.0
  vpc_id                = data.alicloud_vpcs.default.ids.0
  zone_id               = data.alicloud_zones.default.zones.0.id
}
```

## Argument Reference

The following arguments are supported:

* `account_type` - (Optional, Computed, ForceNew) The type of the domain account service. Valid values: `nis`, `ldap`. Default value: `nis`
* `additional_volumes` - (Optional, ForceNew) The additional volumes. See the following `Block additional_volumes`.
* `application` - (Optional, Computed, ForceNew) The application. See the following `Block application`.
* `auto_renew` - (Optional) Specifies whether to enable auto-renewal for the subscription. Default value: `false`.
* `auto_renew_period` - (Optional) The auto-renewal period of the subscription compute nodes. The parameter takes effect when AutoRenew is set to true.
* `client_version` - (Optional, Computed, ForceNew) The version of the E-HPC client. By default, the parameter is set to the latest version number.
* `cluster_name` - (Required) The name of the cluster. The name must be `2` to `64` characters in length.
* `cluster_version` - (Optional) The version of the cluster. Default value: `1.0`.
* `compute_count` - (Required, ForceNew) The number of the compute nodes. Valid values: `1` to `99`.
* `compute_enable_ht` - (Optional) Specifies whether the compute nodes support hyper-threading. Default value: `true`.
* `compute_instance_type` - (Required, ForceNew) The instance type of the compute nodes.
* `compute_spot_price_limit` - (Optional) The maximum hourly price of the compute nodes. A maximum of three decimal places can be used in the value of the parameter. The parameter is valid only when the ComputeSpotStrategy parameter is set to SpotWithPriceLimit.
* `compute_spot_strategy` - (Optional) The bidding method of the compute nodes. Default value: `NoSpot`. Valid values:
  - `NoSpot`: The compute nodes are pay-as-you-go instances.
  - `SpotWithPriceLimit`: The compute nodes are preemptible instances that have a user-defined maximum hourly price.
  - `SpotAsPriceGo`: The compute nodes are preemptible instances for which the market price at the time of purchase is used as the bid price.
* `deploy_mode` - (Optional, Computed, ForceNew) The mode in which the cluster is deployed. Valid values: `Standard`, `Simple`, `Tiny`. Default value: Standard.
  - `Standard`: An account node, a scheduling node, a logon node, and multiple compute nodes are separately deployed.
  - `Simple`: A management node, a logon node, and multiple compute nodes are deployed. The management node consists of an account node and a scheduling node. The logon node and compute nodes are separately deployed.
  - `Tiny`: A management node and multiple compute nodes are deployed. The management node consists of an account node, a scheduling node, and a logon node. The compute nodes are separately deployed.
* `description` - (Optional, Computed) The description of the cluster. The description must be `2` to `256` characters in length. It cannot start with `http://` or `https://`.
* `domain` - (Optional) The domain name of the on-premises cluster. This parameter takes effect only when the AccoutType parameter is set to Idap.
* `ecs_charge_type` - (Optional) The billing method of the nodes.
* `ehpc_version` - (Optional) The version of E-HPC. By default, the parameter is set to the latest version number.
* `ha_enable` - (Optional, Computed, ForceNew) Specifies whether to enable the high availability feature. Default value: `false`.  **Note:** If high availability is enabled, a primary management node and a secondary management node are used.
* `image_id` - (Optional, Computed) The ID of the image.
* `image_owner_alias` - (Optional, Computed) The type of the image. Valid values: `others`, `self`, `system`, `marketplace`. Default value: `system`.
* `input_file_url` - (Optional) The URL of the job files that are uploaded to an Object Storage Service (OSS) bucket.
* `is_compute_ess` - (Optional) Specifies whether to enable auto scaling. Default value: `false`.
* `job_queue` - (Optional) The queue to which the compute nodes are added.
* `key_pair_name` - (Optional) The name of the AccessKey pair.
* `login_count` - (Required, ForceNew) The number of the logon nodes. Valid values: `1`.
* `login_instance_type` - (Required, ForceNew) The instance type of the logon nodes.
* `manager_instance_type` - (Required, ForceNew) The instance type of the management nodes.
* `manager_count` - (Optional, Computed, ForceNew) The number of the management nodes. Valid values: 1 and 2.
* `os_tag` - (Required, ForceNew) The image tag of the operating system.
* `password` - (Optional, ForceNew) The root password of the logon node. The password must be 8 to 30 characters in length and contain at least three of the following items: uppercase letters, lowercase letters, digits, and special characters. The password can contain the following special characters: `( ) ~ ! @ # $ % ^ & * - + = { } [ ] : ; ‘ < > , . ? /`. You must specify either `password` or `key_pair_name`. If both are specified, the Password parameter prevails.
* `period` - (Optional) The duration of the subscription. The unit of the duration is specified by the `period_unit` parameter. Default value: `1`.
  * If you set PriceUnit to Year, the valid values of the Period parameter are 1, 2, and 3.
  * If you set PriceUnit to Month, the valid values of the Period parameter are 1, 2, 3, 4, 5, 6, 7, 8, and 9.
  * If you set PriceUnit to Hour, the valid value of the Period parameter is 1.
* `period_unit` - (Optional) The unit of the subscription duration. Valid values: `Year`, `Month`, `Hour`. Default value: `Month`.
* `plugin` - (Optional) The mode configurations of the plug-in. This parameter takes effect only when the SchedulerType parameter is set to custom. The value must be a JSON string. The parameter contains the following parameters: pluginMod, pluginLocalPath, and pluginOssPath.
  - pluginMod: the mode of the plug-in. The following modes are supported:
    - oss: The plug-in is downloaded and decompressed from OSS to a local path. The local path is specified by the pluginLocalPath parameter.
    - image: By default, the plug-in is stored in a pre-defined local path. The local path is specified by the pluginLocalPath parameter.
  - pluginLocalPath: the local path where the plug-in is stored. We recommend that you select a shared directory in oss mode and a non-shared directory in image mode.
  - pluginOssPath: the remote path where the plug-in is stored in OSS. This parameter takes effect only when the pluginMod parameter is set to oss.
* `post_install_script` - (Optional, ForceNew) The post install script. See the following `Block post_install_script`.
* `ram_node_types` - (Optional) The node of the RAM role.
* `ram_role_name` - (Optional) The name of the Resource Access Management (RAM) role.
* `release_instance` - (Optional) The release instance. Valid values: `true`.
* `remote_directory` - (Optional, Computed, ForceNew) The remote directory to which the file system is mounted.
* `remote_vis_enable` - (Optional) Specifies whether to enable Virtual Network Computing (VNC). Default value: `false`.
* `resource_group_id` - (Optional) The ID of the resource group.
* `scc_cluster_id` - (Optional, Computed, ForceNew) The ID of the Super Computing Cluster (SCC) instance. If you specify the parameter, the SCC instance is moved to a new SCC cluster.
* `scheduler_type` - (Optional, Computed, ForceNew) The type of the scheduler. Valid values: `pbs`, `slurm`, `opengridscheduler` and `deadline`. Default value: `pbs`.
* `security_group_id` - (Optional, Computed, ForceNew) The ID of the security group to which the cluster belongs.
* `security_group_name` - (Optional) If you do not use an existing security group, set the parameter to the name of a new security group. A default policy is applied to the new security group.
* `system_disk_level` - (Optional) The performance level of the ESSD that is used as the system disk. Default value: `PL1` For more information, see [ESSDs](https://www.alibabacloud.com/help/en/elastic-compute-service/latest/essds). Valid values:
  * `PL0`: A single ESSD can deliver up to 10,000 random read/write IOPS.
  * `PL1`: A single ESSD can deliver up to 50,000 random read/write IOPS.
  * `PL2`: A single ESSD can deliver up to 100,000 random read/write IOPS.
  * `PL3`: A single ESSD can deliver up to 1,000,000 random read/write IOPS.
* `system_disk_size` - (Optional) The size of the system disk. Unit: `GB`. Valid values: `40` to `500`. Default value: `40`.
* `system_disk_type` - (Optional) The type of the system disk. Valid values: `cloud_efficiency`, `cloud_ssd`, `cloud_essd` or `cloud`. Default value: `cloud_ssd`.
* `volume_id` - (Optional, Computed, ForceNew) The ID of the file system. If you leave the parameter empty, a Performance NAS file system is created by default.
* `volume_mount_option` - (Optional) The mount options of the file system.
* `volume_mountpoint` - (Optional, Computed, ForceNew) The mount target of the file system. Take note of the following information:
  - If you do not specify the VolumeId parameter, you can leave the VolumeMountpoint parameter empty. A mount target is created by default.
  - If you specify the VolumeId parameter, the VolumeMountpoint parameter is required.
* `volume_protocol` - (Optional, Computed, ForceNew) The type of the protocol that is used by the file system. Valid values: `NFS`, `SMB`. Default value: `NFS`.
* `volume_type` - (Optional, Computed, ForceNew) The type of the shared storage. Only Apsara File Storage NAS file systems are supported.
* `vswitch_id` - (Optional, Computed, ForceNew) The ID of the vSwitch. E-HPC supports only VPC networks.
* `without_agent` - (Optional) Specifies whether not to install the agent. Default value: `false`.
* `without_elastic_ip` - (Optional) Specifies whether the logon node uses an elastic IP address (EIP). Default value: `false`.
* `zone_id` - (Optional, Computed, ForceNew) The ID of the zone.
* `vpc_id` - (Optional, Computed, ForceNew) The ID of the virtual private cloud (VPC) to which the cluster belongs.

#### Block application

The application supports the following:

* `tag` - (Optional) The tag of the software.

#### Block post_install_script

The post_install_script supports the following: 

* `args` - (Optional) The parameter that is used to run the script after the cluster is created.
* `url` - (Optional) The URL that is used to download the script after the cluster is created.

#### Block additional_volumes

The additional_volumes supports the following: 

* `job_queue` - (Optional) The queue of the nodes to which the additional file system is attached.
* `local_directory` - (Optional) The local directory on which the additional file system is mounted.
* `location` - (Optional) The type of the cluster. Valid value: `PublicCloud`.
* `remote_directory` - (Optional) The remote directory to which the additional file system is mounted.
* `roles` - (Optional) The roles. See the following `Block roles`.
* `volume_id` - (Optional) The ID of the additional file system.
* `volume_mount_option` - (Optional) The mount options of the file system.
* `volume_mountpoint` - (Optional) The mount target of the additional file system.
* `volume_protocol` - (Optional) The type of the protocol that is used by the additional file system. Valid values: `NFS`, `SMB`. Default value: `NFS`
* `volume_type` - (Optional) The type of the additional shared storage. Only NAS file systems are supported.

#### Block roles

The roles supports the following:

* `name` - (Optional) The type of the nodes to which the additional file system is attached.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Cluster.
* `status` - The status of the resource.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 20 mins) Used when create the Cluster.
* `delete` - (Defaults to 5 mins) Used when delete the Cluster.
* `update` - (Defaults to 5 mins) Used when update the Cluster.

## Import

Ehpc Cluster can be imported using the id, e.g.

```
$ terraform import alicloud_ehpc_cluster.example <id>
```