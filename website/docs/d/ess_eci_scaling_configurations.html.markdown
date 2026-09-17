---
subcategory: "Auto Scaling"
layout: "alicloud"
page_title: "Alicloud: alicloud_ess_eci_scaling_configurations"
sidebar_current: "docs-alicloud_ess_eci_scaling_configurations"
description: |-
    Provides a list of ECI scaling configurations available to the user.
---

# alicloud_ess_eci_scaling_configurations

This data source provides available ECI scaling configuration resources.

-> **NOTE:** Available since v1.294.0

## Example Usage

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
  security_group_name = local.name
  vpc_id              = alicloud_vpc.default.id
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
  container_group_name = local.name
  containers {
    name  = "container-1"
    image = "registry-vpc.cn-hangzhou.aliyuncs.com/eci_open/alpine:3.5"
  }
}

data "alicloud_ess_eci_scaling_configurations" "default" {
  scaling_group_id = alicloud_ess_scaling_group.default.id
  ids              = [alicloud_ess_eci_scaling_configuration.default.id]
  name_regex       = var.name
}

output "first_scaling_configuration" {
  value = data.alicloud_ess_eci_scaling_configurations.default.configurations[0].id
}
```

## Argument Reference

The following arguments are supported:

* `scaling_group_id` - (Optional) ID of the scaling group to which the scaling configurations belong.
* `name_regex` - (Optional) A regex string to filter scaling configurations by name.
* `ids` - (Optional) A list of scaling configuration IDs.
* `output_file` - (Optional) File name where to save data source results after running `terraform plan`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of scaling configuration IDs.
* `names` - A list of scaling configuration names.
* `configurations` - A list of ECI scaling configurations. Each element contains the following attributes:
  * `id` - ID of the scaling configuration.
  * `scaling_group_id` - ID of the scaling group.
  * `scaling_configuration_name` - Name of the scaling configuration.
  * `description` - Description of the scaling configuration.
  * `security_group_id` - ID of the security group.
  * `container_group_name` - Name of the container group.
  * `restart_policy` - Restart policy for containers.
  * `cpu` - CPU size of the container group.
  * `memory` - Memory size of the container group.
  * `resource_group_id` - ID of the resource group.
  * `dns_policy` - DNS policy of the container group.
  * `cost_optimization` - Whether cost optimization is enabled.
  * `enable_sls` - Whether Simple Log Service is enabled.
  * `instance_family_level` - Instance family level.
  * `image_snapshot_id` - ID of the image cache.
  * `ram_role_name` - Name of the RAM role.
  * `termination_grace_period_seconds` - Grace period before container termination.
  * `auto_match_image_cache` - Whether to automatically match image caches.
  * `ipv6_address_count` - Number of IPv6 addresses.
  * `cpu_options_core` - Number of CPU cores.
  * `cpu_options_threads_per_core` - Number of threads per CPU core.
  * `active_deadline_seconds` - Maximum running time of the container group.
  * `spot_strategy` - Spot strategy.
  * `spot_price_limit` - Maximum hourly spot price.
  * `auto_create_eip` - Whether to automatically create an EIP.
  * `eip_bandwidth` - EIP bandwidth.
  * `host_name` - Hostname of the container group.
  * `ingress_bandwidth` - Ingress bandwidth.
  * `egress_bandwidth` - Egress bandwidth.
  * `ephemeral_storage` - Ephemeral storage size.
  * `load_balancer_weight` - Weight of the container group in a load balancer.
  * `tags` - Tags of the scaling configuration.
  * `instance_types` - Instance types.
  * `creation_time` - Creation time of the scaling configuration.
  * `lifecycle_state` - Lifecycle state of the scaling configuration.
  * `acr_registry_infos` - ACR registry information.
    * `domains` - Domains of the ACR registry.
    * `instance_name` - Name of the ACR instance.
    * `instance_id` - ID of the ACR instance.
    * `region_id` - Region ID of the ACR instance.
  * `image_registry_credentials` - Image registry credentials.
    * `password` - Password of the image registry.
    * `server` - Server address of the image registry.
    * `username` - Username of the image registry.
  * `dns_config_options` - DNS configuration options.
    * `name` - Name of the option.
    * `value` - Value of the option.
  * `security_context_sysctls` - Security context sysctls.
    * `name` - Name of the sysctl.
    * `value` - Value of the sysctl.
  * `containers` - Containers in the container group.
    * `security_context_capability_adds` - Linux capabilities to add.
    * `lifecycle_pre_stop_handler_execs` - Commands run before the container stops.
    * `security_context_read_only_root_file_system` - Whether the root file system is read-only.
    * `tty` - Whether to allocate a TTY.
    * `stdin` - Whether to allocate stdin.
    * `security_context_run_as_user` - UID to run the container as.
    * `ports` - Ports exposed by the container.
      * `port` - Port number.
      * `protocol` - Protocol of the port.
    * `environment_vars` - Environment variables.
      * `key` - Name of the environment variable.
      * `value` - Value of the environment variable.
      * `field_ref_field_path` - Field reference path.
    * `working_dir` - Working directory.
    * `args` - Arguments passed to the container.
    * `cpu` - CPU size.
    * `gpu` - Number of GPUs.
    * `memory` - Memory size.
    * `name` - Name of the container.
    * `image` - Container image.
    * `image_pull_policy` - Image pull policy.
    * `volume_mounts` - Volume mounts.
      * `mount_path` - Path where the volume is mounted.
      * `mount_propagation` - Mount propagation mode.
      * `sub_path` - Subpath within the volume.
      * `name` - Name of the volume.
      * `read_only` - Whether the mount is read-only.
    * `commands` - Commands run by the container.
    * `liveness_probe_exec_commands` - Commands in the liveness probe.
    * `liveness_probe_period_seconds` - Liveness probe period.
    * `liveness_probe_http_get_path` - HTTP path in the liveness probe.
    * `liveness_probe_failure_threshold` - Liveness probe failure threshold.
    * `liveness_probe_initial_delay_seconds` - Liveness probe initial delay.
    * `liveness_probe_http_get_port` - HTTP port in the liveness probe.
    * `liveness_probe_http_get_scheme` - HTTP scheme in the liveness probe.
    * `liveness_probe_tcp_socket_port` - TCP socket port in the liveness probe.
    * `liveness_probe_success_threshold` - Liveness probe success threshold.
    * `liveness_probe_timeout_seconds` - Liveness probe timeout.
    * `readiness_probe_exec_commands` - Commands in the readiness probe.
    * `readiness_probe_period_seconds` - Readiness probe period.
    * `readiness_probe_http_get_path` - HTTP path in the readiness probe.
    * `readiness_probe_failure_threshold` - Readiness probe failure threshold.
    * `readiness_probe_initial_delay_seconds` - Readiness probe initial delay.
    * `readiness_probe_http_get_port` - HTTP port in the readiness probe.
    * `readiness_probe_http_get_scheme` - HTTP scheme in the readiness probe.
    * `readiness_probe_tcp_socket_port` - TCP socket port in the readiness probe.
    * `readiness_probe_success_threshold` - Readiness probe success threshold.
    * `readiness_probe_timeout_seconds` - Readiness probe timeout.
  * `init_containers` - Init containers in the container group.
    * `security_context_capability_adds` - Linux capabilities to add.
    * `security_context_read_only_root_file_system` - Whether the root file system is read-only.
    * `security_context_run_as_user` - UID to run the container as.
    * `ports` - Ports exposed by the init container.
    * `environment_vars` - Environment variables.
    * `working_dir` - Working directory.
    * `args` - Arguments passed to the init container.
    * `cpu` - CPU size.
    * `gpu` - Number of GPUs.
    * `memory` - Memory size.
    * `name` - Name of the init container.
    * `image` - Init container image.
    * `image_pull_policy` - Image pull policy.
    * `volume_mounts` - Volume mounts.
    * `commands` - Commands run by the init container.
  * `volumes` - Volumes in the container group.
    * `config_file_volume_config_file_to_paths` - Files mounted from a config file volume.
      * `content` - File content.
      * `path` - Destination path.
      * `mode` - File mode.
    * `disk_volume_disk_id` - ID of the disk volume.
    * `host_path_volume_type` - Type of the host path volume.
    * `host_path_volume_path` - Path of the host path volume.
    * `config_file_volume_default_mode` - Default mode of config file volume files.
    * `empty_dir_volume_medium` - Storage medium for the empty directory volume.
    * `empty_dir_volume_size_limit` - Size limit of the empty directory volume.
    * `disk_volume_fs_type` - File system type of the disk volume.
    * `disk_volume_disk_size` - Size of the disk volume.
    * `flex_volume_driver` - Driver of the FlexVolume.
    * `flex_volume_fs_type` - File system type of the FlexVolume.
    * `flex_volume_options` - Options of the FlexVolume.
    * `nfs_volume_path` - Path of the NFS volume.
    * `nfs_volume_read_only` - Whether the NFS volume is read-only.
    * `nfs_volume_server` - Server of the NFS volume.
    * `name` - Name of the volume.
    * `type` - Type of the volume.
  * `host_aliases` - Host aliases.
    * `hostnames` - Hostnames for the alias.
    * `ip` - IP address for the alias.
