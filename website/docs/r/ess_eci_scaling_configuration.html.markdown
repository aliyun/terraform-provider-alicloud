---
subcategory: "Auto Scaling"
layout: "alicloud"
page_title: "Alicloud: alicloud_ess_eci_scaling_configuration"
sidebar_current: "docs-alicloud-resource-ess-eci-scaling-configuration"
description: |-
  Provides a ESS eci scaling configuration resource.
---

# alicloud_ess_eci_scaling_configuration

Provides a ESS eci scaling configuration resource.

For information about ess eci scaling configuration, see [CreateEciScalingConfiguration](https://www.alibabacloud.com/help/en/auto-scaling/latest/create-eci-scaling-configuration).

-> **NOTE:** Available since v1.164.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_ess_eci_scaling_configuration&exampleId=0c41e297-00f9-b7db-8042-8d32d1564d4cb01daed9&activeTab=example&spm=docs.r.ess_eci_scaling_configuration.0.0c41e29700&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

locals {
  name = "${var.name}-${random_integer.default.result}"
}

data "alicloud_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = local.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_zones.default.zones[0].id
  vswitch_name = local.name
}

resource "alicloud_security_group" "default" {
  name   = local.name
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_ess_scaling_group" "default" {
  min_size           = 0
  max_size           = 1
  scaling_group_name = local.name
  removal_policies   = ["OldestInstance", "NewestInstance"]
  vswitch_ids        = [alicloud_vswitch.default.id]
  group_type         = "ECI"
}

resource "alicloud_ess_eci_scaling_configuration" "default" {
  scaling_group_id     = alicloud_ess_scaling_group.default.id
  cpu                  = 2
  memory               = 4
  security_group_id    = alicloud_security_group.default.id
  force_delete         = true
  active               = true
  container_group_name = "container-group-1649839595174"
  containers {
    name  = "container-1"
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
* `container_group_name` - (Optional) The name series of the elastic container instances created from the scaling configuration. If you want to use an ordered instance name, specify the value for this parameter in the following format: name_prefix(AUTO_INCREMENT)[begin_number,bits]name_suffix.
  name_prefix: the prefix of the hostname.
  (AUTO_INCREMENT): the sort method. This is a static field.
  begin_number: the start value of the sequential values. Valid values: 0 to 999999.
  bits: the number of digits in sequential values. Valid values: 1 to 6. If the number of digits in the specified begin_number value is greater than the value of the bits field, the bits field is automatically set to 6.
  name_suffix: the suffix of the hostname. This field is optional.
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
* `image_registry_credentials` - (Optional)  The image registry credential.   See [`image_registry_credentials`](#image_registry_credentials) below for
  details.
* `security_context_sysctls` - (Optional, Available since v1.232.0)  The system information about the security context in which the elastic container instance is run.   See [`security_context_sysctls`](#security_context_sysctls) below for
  details.
* `dns_config_options` - (Optional, Available since v1.232.0)  The options. Each option is a name-value pair. The value in the name-value pair is optional.   See [`dns_config_options`](#dns_config_options) below for
  details.
* `containers` - (Optional) The list of containers. See [`containers`](#containers) below for details.
* `init_containers` - (Optional) The list of initContainers. See [`init_containers`](#init_containers) below for details.
* `volumes` - (Optional) The list of volumes. See [`volumes`](#volumes) below for details.
* `host_aliases` - (Optional) HostAliases. See [`host_aliases`](#host_aliases) below.
* `acr_registry_infos` - (Optional, Available in 1.193.1+) Information about the Container Registry Enterprise Edition instance. See [`acr_registry_infos`](#acr_registry_infos) below for details.
* `image_snapshot_id` - (Optional) The ID of image cache.
* `termination_grace_period_seconds` - (Optional) The program's buffering time before closing.
* `auto_match_image_cache` - (Optional) Whether to automatically match the image cache.
* `ipv6_address_count` - (Optional) Number of IPv6 addresses.
* `cpu_options_core` - (Optional, Available since v1.227.0) The number of physical CPU cores. You can specify this parameter for only specific instance types.
* `cpu_options_threads_per_core` - (Optional, Available since v1.227.0) The number of threads per core. You can specify this parameter for only specific instance types. If you set this parameter to 1, Hyper-Threading is disabled.
* `active_deadline_seconds` - (Optional) The duration in seconds relative to the startTime that the job may be active before the system tries to terminate it.
* `ephemeral_storage` - (Optional) The size of ephemeral storage.
* `load_balancer_weight` - (Optional) The weight of an ECI instance attached to the Server Group.
* `instance_types` - (Optional, Available since v1.223.0) The specified ECS instance types. You can specify up to five ECS instance types.
* `cost_optimization` - (Optional, Available since v1.232.0) Indicates whether the Cost Optimization feature is enabled. Valid values: true,false.
* `instance_family_level` - (Optional, Available since v1.232.0) The level of the instance family, which is used to filter instance types that meet the specified criteria. This parameter takes effect only if you set CostOptimization to true. Valid values: EntryLevel, EnterpriseLevel, CreditEntryLevel.

### `volumes`

The volume supports the following:

* `name` - (Optional) The name of the volume.
* `host_path_volume_type` - (Optional, Available since v1.232.0) The type of the host path. Examples: File, Directory, and Socket.
* `host_path_volume_path` - (Optional, Available since v1.232.0) The absolute path on the host.
* `config_file_volume_default_mode` - (Optional, Available since v1.232.0) The default permissions on the ConfigFileVolume.
* `empty_dir_volume_medium` - (Optional, Available since v1.232.0) The storage medium of the EmptyDirVolume. If you leave this parameter empty, the file system of the node is used as the storage medium. If you set this parameter to memory, the memory is used as the storage medium.
* `empty_dir_volume_size_limit` - (Optional, Available since v1.232.0) The storage size of the EmptyDirVolume. Unit: GiB or MiB.
* `type` - (Optional) The type of the volume.
* `config_file_volume_config_file_to_paths` - (Optional) ConfigFileVolumeConfigFileToPaths.
  See [`config_file_volume_config_file_to_paths`](#volumes-config_file_volume_config_file_to_paths) below for details.
* `disk_volume_disk_id` - (Optional) The ID of DiskVolume.
* `disk_volume_fs_type` - (Optional) The system type of DiskVolume.
* `disk_volume_disk_size` - (Optional) The disk size of DiskVolume.
* `flex_volume_driver` - (Optional) The name of the FlexVolume driver.
* `flex_volume_fs_type` - (Optional) The type of the mounted file system. The default value is determined by the script
  of FlexVolume.
* `flex_volume_options` - (Optional) The list of FlexVolume objects. Each object is a key-value pair contained in a JSON
  string.
* `nfs_volume_path` - (Optional) The path to the NFS volume.
* `nfs_volume_read_only` - (Optional) The nfs volume read only. Default to `false`.
* `nfs_volume_server` - (Optional) The address of the NFS server.

-> **NOTE:** Every volume mounted must have a name and type attributes.

### `volumes-config_file_volume_config_file_to_paths`

The config_file_volume_config_file_to_path supports the following:

* `content` - (Optional) The content of the configuration file. Maximum size: 32 KB.
* `path` - (Optional) The relative file path.
* `mode` - (Optional, Available since v1.232.0) The permissions on the ConfigFileVolume directory.

### `init_containers`

The init_container supports the following:

* `args` - (Optional) The arguments passed to the commands.
* `commands` - (Optional) The commands run by the init container.
* `cpu` - (Optional) The amount of CPU resources allocated to the container.
* `environment_vars` - (Optional) The structure of environmentVars. 
  See [`environment_vars`](#init_containers-environment_vars) below for details.
* `gpu` - (Optional) The number GPUs.
* `image` - (Optional) The image of the container.
* `image_pull_policy` - (Optional) The restart policy of the image.
* `memory` - (Optional) The amount of memory resources allocated to the container.
* `name` - (Optional) The name of the init container.
* `ports` - (Optional) The structure of port. See [`ports`](#init_containers-ports) below for details.
* `volume_mounts` - (Optional) The structure of volumeMounts. See [`volume_mounts`](#init_containers-volume_mounts) below for details.
* `working_dir` - (Optional) The working directory of the container.
* `security_context_capability_adds` - (Optional, Available since 1.215.0) Grant certain permissions to processes within container. Optional values:
  - NET_ADMIN: Allow network management tasks to be performed.
  - NET_RAW: Allow raw sockets.
* `security_context_read_only_root_file_system` - (Optional, Available since 1.215.0) Mounts the container's root filesystem as read-only.
* `security_context_run_as_user` - (Optional, Available since 1.215.0) Specifies user ID  under which all processes run.

### `init_containers-environment_vars`

The environment_var supports the following:

* `key` - (Optional) The name of the variable. The name can be 1 to 128 characters in length and can contain letters,
  digits, and underscores (_). It cannot start with a digit.
* `value` - (Optional) The value of the variable. The value can be 0 to 256 characters in length.
* `field_ref_field_path` - (Optional, Available since 1.215.0) Environment variable value reference. Optional values:
  - status.podIP: IP of pod.

### `init_containers-ports`

The port supports the following:

* `port` - (Optional, ForceNew) The port number. Valid values: 1 to 65535.
* `protocol` - (Optional, ForceNew) Valid values: TCP and UDP.

### `init_containers-volume_mounts`

The volume_mount supports the following:

* `mount_path` - (Optional) The directory of the mounted volume. Data under this directory will be overwritten by the
  data in the volume.
* `name` - (Optional) The name of the mounted volume.
* `read_only` - (Optional) Default to `false`.
* `sub_path` - (Optional, Available since 1.232.0) The subdirectory of volume N.
* `mount_propagation` - (Optional, Available since 1.232.0) The mount propagation settings of volume N. Mount propagation enables volumes mounted on one container to be shared among other containers within the same pod or across distinct pods residing on the same node. Valid values: None, HostToCotainer, Bidirectional.

### `image_registry_credentials`

The image_registry_credential supports the following:

* `password` - (Optional) The password used to log on to the image repository. It is required
  when `image_registry_credential` is configured.
* `server` - (Optional) The address of the image repository. It is required when `image_registry_credential` is
  configured.
* `username` - (Optional) The username used to log on to the image repository. It is required
  when `image_registry_credential` is configured.

### `security_context_sysctls`

The security_context_sysctl supports the following:

* `name` - (Optional, Available since v1.232.0) The system name of the security context in which the elastic container instance is run.
* `value` - (Optional, Available since v1.232.0) The system value of the security context in which the elastic container instance is run.

### `dns_config_options`

The dns_config_option supports the following:

* `name` - (Optional, Available since v1.232.0) The option name.
* `value` - (Optional, Available since v1.232.0) The option value.

### `host_aliases`

The host_aliases supports the following:

* `hostnames` - (Optional) Adds a host name.
* `ip` - (Optional) Adds an IP address.

### `containers`

The container supports the following:

* `args` - (Optional) The arguments passed to the commands.
* `commands` - (Optional) The commands run by the init container.
* `cpu` - (Optional) The amount of CPU resources allocated to the container.
* `environment_vars` - (Optional) The structure of environmentVars.
  See [`environment_vars`](#containers-environment_vars) below for details.
* `gpu` - (Optional) The number GPUs.
* `image` - (Optional) The image of the container.
* `image_pull_policy` - (Optional) The restart policy of the image.
* `memory` - (Optional) The amount of memory resources allocated to the container.
* `name` - (Optional) The name of the init container.
* `ports` - (Optional) The structure of port. See [`ports`](#containers-ports) below for details.
* `volume_mounts` - (Optional) The structure of volumeMounts. 
   See [`volume_mounts`](#containers-volume_mounts) below for details.
* `working_dir` - (Optional) The working directory of the container.
* `liveness_probe_exec_commands` - (Optional, Available in 1.193.1+) Commands that you want to run in containers when you use the CLI to perform liveness probes.
* `liveness_probe_period_seconds` - (Optional, Available in 1.193.1+) The interval at which the liveness probe is performed. Unit: seconds. Default value: 10. Minimum value: 1.
* `liveness_probe_http_get_path` - (Optional, Available in 1.193.1+) The path to which HTTP GET requests are sent when you use HTTP requests to perform liveness probes.
* `liveness_probe_failure_threshold` - (Optional, Available in 1.193.1+) The minimum number of consecutive failures for the liveness probe to be considered failed after having been successful. Default value: 3.
* `liveness_probe_initial_delay_seconds` - (Optional, Available in 1.193.1+) The number of seconds after container has started before liveness probes are initiated.
* `liveness_probe_http_get_port` - (Optional, Available in 1.193.1+) The port to which HTTP GET requests are sent when you use HTTP requests to perform liveness probes.
* `liveness_probe_http_get_scheme` - (Optional, Available in 1.193.1+) The protocol type of HTTP GET requests when you use HTTP requests for liveness probes.Valid values:HTTP and HTTPS.
* `liveness_probe_tcp_socket_port` - (Optional, Available in 1.193.1+) The port detected by TCP sockets when you use TCP sockets to perform liveness probes.
* `liveness_probe_success_threshold` - (Optional, Available in 1.193.1+) The minimum number of consecutive successes for the liveness probe to be considered successful after having failed. Default value: 1. Set the value to 1.
* `liveness_probe_timeout_seconds` - (Optional, Available in 1.193.1+) The timeout period for the liveness probe. Unit: seconds. Default value: 1. Minimum value: 1.
* `readiness_probe_exec_commands` - (Optional, Available in 1.193.1+) Commands that you want to run in containers when you use the CLI to perform readiness probes.
* `readiness_probe_period_seconds` - (Optional, Available in 1.193.1+) The interval at which the readiness probe is performed. Unit: seconds. Default value: 10. Minimum value: 1.
* `readiness_probe_http_get_path` - (Optional, Available in 1.193.1+) The path to which HTTP GET requests are sent when you use HTTP requests to perform readiness probes.
* `readiness_probe_failure_threshold` - (Optional, Available in 1.193.1+) The minimum number of consecutive failures for the readiness probe to be considered failed after having been successful. Default value: 3.
* `readiness_probe_initial_delay_seconds` - (Optional, Available in 1.193.1+) The number of seconds after container N has started before readiness probes are initiated.
* `readiness_probe_http_get_port` - (Optional, Available in 1.193.1+) The port to which HTTP GET requests are sent when you use HTTP requests to perform readiness probes.
* `readiness_probe_http_get_scheme` - (Optional, Available in 1.193.1+) The protocol type of HTTP GET requests when you use HTTP requests for readiness probes. Valid values: HTTP and HTTPS.
* `readiness_probe_tcp_socket_port` - (Optional, Available in 1.193.1+) The port detected by Transmission Control Protocol (TCP) sockets when you use TCP sockets to perform readiness probes.
* `readiness_probe_success_threshold` - (Optional, Available in 1.193.1+) The minimum number of consecutive successes for the readiness probe to be considered successful after having failed. Default value: 1. Set the value to 1.
* `readiness_probe_timeout_seconds` - (Optional, Available in 1.193.1+) The timeout period for the readiness probe. Unit: seconds. Default value: 1. Minimum value: 1.
* `security_context_capability_adds` - (Optional, Available since 1.215.0) Grant certain permissions to processes within container. Optional values:
  - NET_ADMIN: Allow network management tasks to be performed.
  - NET_RAW: Allow raw sockets.
* `lifecycle_pre_stop_handler_execs` - (Optional, Available since 1.216.0) The commands to be executed in containers when you use the CLI to specify the preStop callback function.
* `security_context_read_only_root_file_system` - (Optional, Available since 1.215.0) Mounts the container's root filesystem as read-only.
* `security_context_run_as_user` - (Optional, Available since 1.215.0) Specifies user ID  under which all processes run.
* `tty` - (Optional, Available since v1.232.0) Specifies whether to enable the Interaction feature. Valid values: true, false.
* `stdin` - (Optional, Available since v1.232.0) Specifies whether container N allocates buffer resources to standard input streams during its active runtime. If you do not specify this parameter, an end-of-file (EOF) error occurs.

### `containers-environment_vars`

The environment_var supports the following:

* `key` - (Optional) The name of the variable. The name can be 1 to 128 characters in length and can contain letters,
  digits, and underscores (_). It cannot start with a digit.
* `value` - (Optional) The value of the variable. The value can be 0 to 256 characters in length.
* `field_ref_field_path` - (Optional, Available since 1.215.0) Environment variable value reference. Optional values: 
  - status.podIP: IP of pod.

### `containers-ports`

The port supports the following:

* `port` - (Optional, ForceNew) The port number. Valid values: 1 to 65535.
* `protocol` - (Optional, ForceNew) Valid values: TCP and UDP.

### `containers-volume_mounts`

The volume_mount supports the following:

* `mount_path` - (Optional) The directory of the mounted volume. Data under this directory will be overwritten by the
  data in the volume.
* `name` - (Optional) The name of the mounted volume.
* `read_only` - (Optional) Default to `false`.
* `sub_path` - (Optional, Available since v1.232.0) The subdirectory of volume N.
* `mount_propagation` - (Optional, Available since v1.232.0) The mount propagation settings of volume N. Mount propagation enables volumes mounted on one container to be shared among other containers within the same pod or across distinct pods residing on the same node. Valid values: None, HostToCotainer, Bidirectional. 

### `acr_registry_infos`

The acr_registry_info supports the following:

* `domains` - (Optional) Endpoint of Container Registry Enterprise Edition instance. By default, all endpoints of the Container Registry Enterprise Edition instance are displayed. It is required
  when `acr_registry_info` is configured.
* `instance_name` - (Optional) The name of Container Registry Enterprise Edition instance. It is required when `acr_registry_info` is
  configured.
* `instance_id` - (Optional) The ID of Container Registry Enterprise Edition instance. It is required
  when `acr_registry_info` is configured.
* `region_id` - (Optional) The region ID of Container Registry Enterprise Edition instance. It is required
  when `acr_registry_info` is configured.

## Attributes Reference

The following attributes are exported:

* `id` - The eci scaling configuration ID.

## Import

ESS eci scaling configuration can be imported using the id, e.g.

```shell
$ terraform import alicloud_ess_eci_scaling_configuration.example asc-abc123456
```

