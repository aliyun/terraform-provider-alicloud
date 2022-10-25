---
subcategory: "Auto Scaling(ESS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ess_eci_scaling_configuration"
sidebar_current: "docs-alicloud-resource-ess-eci-scaling-configuration"
description: |-
  Provides a ESS eci scaling configuration resource.
---

# alicloud\_ess\_eci\_scaling\_configuration

Provides a ESS eci scaling configuration resource.

-> **NOTE:** Resource `alicloud_ess_alb_server_group_attachment` is available in 1.164.0+.

## Example Usage

Basic Usage

```
variable "name" {
  default = "essscalingconfiguration"
}

resource "alicloud_vpc" "default" {
  name       = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id            = alicloud_vpc.default.id
  cidr_block        = "172.16.0.0/24"
  zone_id           = data.alicloud_zones.default.zones[0].id
  vswitch_name      = var.name
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_ess_scaling_group" "default" {
  min_size           = 0
  max_size           = 1
  scaling_group_name = var.name
  removal_policies   = ["OldestInstance", "NewestInstance"]
  vswitch_ids        = [alicloud_vswitch.default.id]
  group_type          = "ECI"
}

resource "alicloud_ess_eci_scaling_configuration" "default" {
  scaling_group_id  = alicloud_ess_scaling_group.default.id
  cpu               = 2
  memory            = 4
  security_group_id = alicloud_security_group.default.id
  force_delete      = true
  active            = true
  container_group_name = "container-group-1649839595174"
  containers {
    name = "container-1"
    image = "registry-vpc.cn-hangzhou.aliyuncs.com/eci_open/alpine:3.5"
  }
}
```

## Argument Reference

The following arguments are supported:

* `active` - (Optional) Whether active current eci scaling configuration in the specified scaling group. Note that only
  one configuration can be active. Default to `false`.
* `force_delete` - (Optional) The eci scaling configuration will be deleted forcibly with deleting its scaling group.
  Default to false.
* `scaling_group_id` - (Required, ForceNew) ID of the scaling group of a eci scaling configuration.
* `scaling_configuration_name` - (Optional) Name shown for the scheduled task. which must contain 2-64 characters (
  English or Chinese), starting with numbers, English letters or Chinese characters, and can contain number,
  underscores `_`, hypens `-`, and decimal point `.`. If this parameter value is not specified, the default value is
  EciScalingConfigurationId.
* `description` - (Optional) The description of data disk N. Valid values of N: 1 to 16. The description must be 2 to
  256 characters in length and cannot start with http:// or https://.
* `security_group_id` - (Optional) ID of the security group used to create new instance. It is conflict
  with `security_group_ids`.
* `container_group_name` - (Optional) The name of the container group.
* `restart_policy` - (Optional) The restart policy of the container group. Default to `Always`.
* `cpu` - (Optional) The amount of CPU resources allocated to the container group.
* `memory` - (Optional) The amount of memory resources allocated to the container group.
* `resource_group_id` - (Optional) ID of resource group.
* `dns_policy` - (Optional) dns policy of contain group.
* `enable_sls` - (Optional) Enable sls log service.
* `ram_role_name` - (Optional) The RAM role that the container group assumes. ECI and ECS share the same RAM role.
* `spot_strategy` - (Optional) The spot strategy for a Pay-As-You-Go instance. Valid values: `NoSpot`, `SpotAsPriceGo`
  , `SpotWithPriceLimit`.
* `spot_price_limit` - (Optional) The maximum price hourly for spot instance.
* `auto_create_eip` - (Optional) Whether create eip automatically.
* `eip_bandwidth` - (Optional) Eip bandwidth.
* `ingress_bandwidth` - (Optional) Ingress bandwidth.
* `egress_bandwidth` - (Optional) egress bandwidth.
* `host_name` - (Optional) Hostname of an ECI instance.
* `tags` - (Optional) A mapping of tags to assign to the resource. It will be applied for ECI instances finally.
    - Key: It can be up to 64 characters in length. It cannot begin with "aliyun", "http://", or "https://". It cannot
      be a null string.
    - Value: It can be up to 128 characters in length. It cannot begin with "aliyun", "http://", or "https://" It can be
      a null string.
* `image_registry_credentials` - (Optional)  The image registry credential. The details see
  Block `image_registry_credential`.See [Block image_registry_credential](#block-image_registry_credential) below for
  details.
* `containers` - (Optional) The list of containers.See [Block container](#block-container) below for details.
* `init_containers` - (Optional) The list of initContainers.See [Block init_container](#block-init_container) below for
  details.
* `volumes` - (Optional) The list of volumes.See [Block volume](#block-volume) below for details.
* `host_aliases` - (Optional) HostAliases.See [Block host_alias](#block-host_alias) below for details.

#### Block volume

The volume supports the following:

* `name` - (Optional) The name of the volume.
* `type` - (Optional) The type of the volume.
* `config_file_volume_config_file_to_paths` - (Optional) ConfigFileVolumeConfigFileToPaths.
  See [Block_config_file_volume_config_file_to_path](#block-config_file_volume_config_file_to_path) below for details.
* `disk_volume_disk_id` - (Optional) The ID of DiskVolume.
* `disk_volume_fs_type` - (Optional) The system type of DiskVolume.
* `flex_volume_driver` - (Optional) The name of the FlexVolume driver.
* `flex_volume_fs_type` - (Optional) The type of the mounted file system. The default value is determined by the script
  of FlexVolume.
* `flex_volume_options` - (Optional) The list of FlexVolume objects. Each object is a key-value pair contained in a JSON
  string.
* `nfs_volume_path` - (Optional) The path to the NFS volume.
* `nfs_volume_read_only` - (Optional) The nfs volume read only. Default to `false`.
* `nfs_volume_server` - (Optional) The address of the NFS server.

-> **NOTE:** Every volume mounted must have a name and type attributes.

#### Block config_file_volume_config_file_to_path

The config_file_volume_config_file_to_path supports the following:

* `content` - (Optional) The content of the configuration file. Maximum size: 32 KB.
* `path` - (Optional) The relative file path.

#### Block init_container

The init_container supports the following:

* `args` - (Optional) The arguments passed to the commands.
* `commands` - (Optional) The commands run by the init container.
* `cpu` - (Optional) The amount of CPU resources allocated to the container.
* `environment_vars` - (Optional) The structure of environmentVars.
  See [Block_environment_var_in_init_container](#block-environment_var_in_init_containers) below for details.
* `gpu` - (Optional) The number GPUs.
* `image` - (Optional) The image of the container.
* `image_pull_policy` - (Optional) The restart policy of the image.
* `memory` - (Optional) The amount of memory resources allocated to the container.
* `name` - (Optional) The name of the init container.
* `ports` - (Optional) The structure of port. See [Block_port_in_init_container](#block-port_in_init_container) below
  for details.
* `volume_mounts` - (Optional) The structure of volumeMounts.
  See [Block_volume_mount_in_init_container](#block-volume_mount_in_init_container) below for details.
* `working_dir` - (Optional) The working directory of the container.

#### Block environment_var_in_init_container

The environment_var supports the following:

* `key` - (Optional) The name of the variable. The name can be 1 to 128 characters in length and can contain letters,
  digits, and underscores (_). It cannot start with a digit.
* `value` - (Optional) The value of the variable. The value can be 0 to 256 characters in length.

#### Block port_in_init_container

The port supports the following:

* `port` - (Optional, ForceNew) The port number. Valid values: 1 to 65535.
* `protocol` - (Optional, ForceNew) Valid values: TCP and UDP.

#### Block volume_mount_in_init_container

The volume_mount supports the following:

* `mount_path` - (Optional) The directory of the mounted volume. Data under this directory will be overwritten by the
  data in the volume.
* `name` - (Optional) The name of the mounted volume.
* `read_only` - (Optional) Default to `false`.

#### Block host_alias

The host_alias supports the following:

* `hostnames` - (Optional) Adds a host name.
* `ip` - (Optional) Adds an IP address.

#### Block image_registry_credential

The image_registry_credential supports the following:

* `password` - (Optional) The password used to log on to the image repository. It is required
  when `image_registry_credential` is configured.
* `server` - (Optional) The address of the image repository. It is required when `image_registry_credential` is
  configured.
* `user_name` - (Optional) The username used to log on to the image repository. It is required
  when `image_registry_credential` is configured.

#### Block container

The container supports the following:

* `args` - (Optional) The arguments passed to the commands.
* `commands` - (Optional) The commands run by the init container.
* `cpu` - (Optional) The amount of CPU resources allocated to the container.
* `environment_vars` - (Optional) The structure of environmentVars.
  See [Block_environment_var_in_container](#block-environment_var_in_container) below for details.
* `gpu` - (Optional) The number GPUs.
* `image` - (Optional) The image of the container.
* `image_pull_policy` - (Optional) The restart policy of the image.
* `memory` - (Optional) The amount of memory resources allocated to the container.
* `name` - (Optional) The name of the init container.
* `ports` - (Optional) The structure of port. See [Block_port_in_container](#block-port_in_container) below for details.
* `volume_mounts` - (Optional) The structure of volumeMounts.
  See [Block_volume_mount_in_container](#block-volume_mount_in_container) below for details.
* `working_dir` - (Optional) The working directory of the container.

#### Block environment_var_in_container

The environment_var supports the following:

* `key` - (Optional) The name of the variable. The name can be 1 to 128 characters in length and can contain letters,
  digits, and underscores (_). It cannot start with a digit.
* `value` - (Optional) The value of the variable. The value can be 0 to 256 characters in length.

#### Block port_in_container

The port supports the following:

* `port` - (Optional, ForceNew) The port number. Valid values: 1 to 65535.
* `protocol` - (Optional, ForceNew) Valid values: TCP and UDP.

#### Block volume_mount_in_container

The volume_mount supports the following:

* `mount_path` - (Optional) The directory of the mounted volume. Data under this directory will be overwritten by the
  data in the volume.
* `name` - (Optional) The name of the mounted volume.
* `read_only` - (Optional) Default to `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The eci scaling configuration ID.

## Import

ESS eci scaling configuration can be imported using the id, e.g.

```
$ terraform import alicloud_ess_eci_scaling_configuration.example asc-abc123456
```

