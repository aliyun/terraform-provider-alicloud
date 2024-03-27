---
subcategory: "Auto Scaling"
layout: "alicloud"
page_title: "Alicloud: alicloud_ess_eci_scaling_configurations"
sidebar_current: "docs-alicloud_ess_eci_scaling_configurations"
description: |-
    Provides a list of eci scaling configurations available to the user.
---

# alicloud_ess_eci_scaling_configurations

This data source provides available eci scaling configuration resources.

-> **NOTE:** Available in 1.212.0+

## Example Usage

Basic Usage

```terraform
data "alicloud_ess_eci_scaling_configurations" "scalingconfigurations_ds" {
  scaling_group_id = "scaling_group_id"
  ids              = ["scaling_configuration_id1", "scaling_configuration_id2"]
  name_regex       = "scaling_configuration_name"
}

output "first_scaling_rule" {
  value = "${data.alicloud_ess_eci_scaling_configurations.scalingconfigurations_ds.configurations.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `scaling_group_id` - (Optional) Scaling group id the eci scaling configurations belong to.
* `name_regex` - (Optional) A regex string to filter resulting eci scaling configurations by name.
* `ids` - (Optional) A list of eci scaling configuration IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of eci scaling configuration ids.
* `names` - A list of eci scaling configuration names.
* `configurations` - A list of eci scaling configurations. Each element contains the following attributes:
    * `id` - ID of the eci scaling configuration.
    * `scaling_group_id` - ID of the scaling group.
    * `name` - Name of the eci scaling configuration.
    * `security_group_id` - Security group ID of the eci scaling configuration.
    * `description` - Description of the eci scaling configuration.
    * `tags` - The tags of the elastic container instance.
    * `restart_policy` - The restart policy of the container group.
    * `container_group_name` - The name of the container group.
    * `resource_group_id` - ID of resource group.
    * `dns_policy` - Dns policy of contain group.
    * `spot_price_limit` - The maximum price hourly for spot instance.
    * `egress_bandwidth` - Egress bandwidth.
    * `auto_create_eip` - Whether create eip automatically.
    * `memory` - The amount of memory resources allocated to the container group.
    * `lifecycle_state` - The state of the eci scaling configuration in the scaling group.
    * `creation_time` - The time at which the eci scaling configuration was created.
    * `eip_bandwidth` - Eip bandwidth.
    * `ram_role_name` - The RAM role that the container group assumes. ECI and ECS share the same RAM role.
    * `ingress_bandwidth` - Ingress bandwidth.
    * `host_name` - Hostname of an ECI instance.
    * `spot_strategy` - The spot strategy for a Pay-As-You-Go instance.
    * `cpu` - The amount of CPU resources allocated to the container group.
    * `containers` - Containers of the scaling configuration.
        * `ports` - Ports of container.
            * `port` - The port number.
            * `protocol` - The type of the protocol.
        * `environment_vars` - The structure of environmentVars.
            * `key` - The key of the environment variable.
            * `value` - The value of the environment variable.
        * `working_dir` - The working directory of the container.
        * `args` - The arguments passed to the commands.
        * `cpu` - The amount of CPU resources allocated to the container.
        * `gpu` - (Optional) The number GPUs.
        * `memory` - The amount of memory resources allocated to the container.
        * `name` - The name of the init container.
        * `image` - The image of the container.
        * `image_pull_policy` - The restart policy of the image.
        * `volume_mounts` - The structure of volumeMounts.
            * `mount_path` - The directory of the mounted volume. Data under this directory will be overwritten by the
              data in the volume.
            * `name` - The name of the mounted volume.
            * `read_only` - Default to `false`.
        * `commands` - The commands run by the init container.
        * `liveness_probe_exec_commands` - Commands that you want to run in containers when you use the CLI to perform
          liveness probes.
        * `liveness_probe_period_seconds` - The interval at which the liveness probe is performed.
        * `liveness_probe_http_get_path` - The path to which HTTP GET requests are sent when you use HTTP requests to
          perform liveness probes.
        * `liveness_probe_failure_threshold` - The minimum number of consecutive failures for the liveness probe to be
          considered failed after having been successful.
        * `liveness_probe_initial_delay_seconds` - The number of seconds after container has started before liveness
          probes are initiated.
        * `liveness_probe_http_get_port` - The port to which HTTP GET requests are sent when you use HTTP requests to
          perform liveness probes.
        * `liveness_probe_http_get_scheme` - The protocol type of HTTP GET requests when you use HTTP requests for
          liveness probes.
        * `liveness_probe_tcp_socket_port` - The port detected by TCP sockets when you use TCP sockets to perform
          liveness probes.
        * `liveness_probe_success_threshold` - The minimum number of consecutive successes for the liveness probe to be
          considered successful after having failed.
        * `liveness_probe_timeout_seconds` - The timeout period for the liveness probe.
        * `readiness_probe_exec_commands` - Commands that you want to run in containers when you use the CLI to perform
          readiness probes.
        * `readiness_probe_period_seconds` - The interval at which the readiness probe is performed.
        * `readiness_probe_http_get_path` - The path to which HTTP GET requests are sent when you use HTTP requests to
          perform readiness probes.
        * `readiness_probe_failure_threshold` - The minimum number of consecutive failures for the readiness probe to be
          considered failed after having been successful.
        * `readiness_probe_initial_delay_seconds` - The number of seconds after container N has started before readiness
          probes are initiated.
        * `readiness_probe_http_get_port` - The port to which HTTP GET requests are sent when you use HTTP requests to
          perform readiness probes.
        * `readiness_probe_http_get_scheme` - The protocol type of HTTP GET requests when you use HTTP requests for
          readiness probes.
        * `readiness_probe_tcp_socket_port` - The port detected by Transmission Control Protocol (TCP) sockets when you
          use TCP sockets to perform readiness probes.
        * `readiness_probe_success_threshold` - The minimum number of consecutive successes for the readiness probe to
          be considered successful after having failed.
        * `readiness_probe_timeout_seconds` - The timeout period for the readiness probe.
    * `volumes` - The list of volumes.
        * `config_file_volume_config_file_to_paths` - The paths to configuration files.
            * `content` - The content of the configuration file.
            * `path` - The relative file path.
        * `disk_volume_disk_id` - The ID of DiskVolume.
        * `disk_volume_fs_type` - The system type of DiskVolume.
        * `disk_volume_disk_size` - The disk size of DiskVolume.
        * `flex_volume_driver` - The name of the FlexVolume driver.
        * `flex_volume_fs_type` - The type of the mounted file system.
        * `flex_volume_options` - The list of FlexVolume objects.
        * `nfs_volume_path` - The path to the NFS volume.
        * `nfs_volume_read_only` - The nfs volume read only.
        * `nfs_volume_server` - The address of the NFS server.
        * `name` - The name of the volume.
        * `type` - The type of the volume.
    * `host_aliases` - The hostnames and IP addresses of a container that are added to the hosts file of the elastic container instance.
        * `hostnames` - Adds a host name.
        * `ip` - Adds an IP address.
    * `image_registry_credentials` - The image registry credential.
        * `password` - The password used to log on to the image repository.
        * `server` - The address of the image repository.
        * `username` - The username used to log on to the image repository.
    * `acr_registry_infos` - Information about the Container Registry Enterprise Edition instance.
        * `domains` - Endpoint of Container Registry Enterprise Edition instance.
        * `instance_name` - The name of Container Registry Enterprise Edition instance.
        * `instance_id` - The ID of Container Registry Enterprise Edition instance.
        * `region_id` - The region ID of Container Registry Enterprise Edition instance.
    * `init_containers` - The list of initContainers.
        * `ports` - The structure of port.
            * `port` - The port number.
            * `protocol` - The type of the protocol.
        * `environment_vars` - The environment variables.
            * `key` - The name of the variable.
            * `value` - The value of the variable.
        * `working_dir` - The working directory of the container.
        * `args` - The arguments passed to the commands.
        * `cpu` - The amount of CPU resources allocated to the container.
        * `gpu` - The number GPUs.
        * `memory` - The amount of memory resources allocated to the container.
        * `name` - The name of the init container.
        * `image` - The image of the container.
        * `image_pull_policy` - The restart policy of the image.
        * `volume_mounts` - The structure of volumeMounts.
            * `mount_path` - The directory of the mounted volume.
            * `name` - The name of the mounted volume.
            * `read_only` - Indicates whether the volume is read-only.
        * `commands` - The commands run by the init container.
