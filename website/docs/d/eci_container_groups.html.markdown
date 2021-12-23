---
subcategory: "Elastic Container Instance (ECI)"
layout: "alicloud"
page_title: "Alicloud: alicloud_eci_container_groups"
sidebar_current: "docs-alicloud-datasource-eci-container-groups"
description: |-
  Provides a list of Eci Container Groups to the user.
---

# alicloud\_eci\_container\_groups

This data source provides the Eci Container Groups of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.111.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_eci_container_groups" "example" {
  ids = ["example_value"]
}

output "first_eci_container_group_id" {
  value = data.alicloud_eci_container_groups.example.groups.0.id
}
```

## Argument Reference

The following arguments are supported:

* `container_group_name` - (Optional, ForceNew) The name of ContainerGroup.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed) A list of Container Group IDs.
* `limit` - (Optional, ForceNew) The maximum number of resources returned in the response. Default value is `20`. Maximum value: `20`. The number of returned results is no greater than the specified number.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Container Group name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `resource_group_id` - (Optional, ForceNew) The ID of the resource group to which the container group belongs. If you have not specified a resource group for the container group, it is added to the default resource group.
* `status` - (Optional, ForceNew) The status list. For more information, see the description of ContainerGroup arrays.
* `vswitch_id` - (Optional, ForceNew) The ID of the VSwitch. Currently, container groups can only be deployed in VPC networks.
* `zone_id` - (Optional, ForceNew) The ID of the zone where you want to deploy the container group. If no value is specified, the system assigns a zone to the container group. By default, no value is specified.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Container Group names.
* `groups` - A list of Eci Container Groups. Each element contains the following attributes:
	* `container_group_id` - The id if ContainerGroup.
	* `container_group_name` - The name of ContainerGroup.
	* `containers` - A list of containers. Each element contains the following attributes:
		* `ready` - Indicates whether the container is ready.
		* `commands` - The commands run by the container. You can define a maximum of 20 commands. Minimum length per string: 256 characters.
		* `cpu` - The amount of CPU resources allocated to the container.
		* `ports` - The list of exposed ports and protocols. Maximum: 100.
			* `port` - The port number. Valid values: 1 to 65535.
			* `protocol` - Valid values: `TCP` and `UDP`.
		* `volume_mounts` - The list of volumes mounted to the container.
			* `mount_path` - The directory of the mounted volume. Data under this directory will be overwritten by the data in the volume.
			* `name` - The name of the volume. The name is the same as the volume you selected when you purchased the container.
			* `read_only` - Default value: `false`.
		* `args` - The arguments passed to the commands. Maximum: `10`.
		* `image_pull_policy` - The policy for pulling an image.
		* `working_dir` - The working directory of the container.
		* `image` - The image of the container.
		* `memory` - The amount of memory resources allocated to the container.
		* `name` - The name of the container.
		* `restart_count` - The number of times that the container has restarted.
		* `environment_vars` - The environment variables.
			* `key` - The name of the variable.
			* `value` - The value of the variable.
		* `gpu` - The amount of GPU resources allocated to the container.
	* `cpu` - The amount of CPU resources allocated to the container group.
	* `dns_config` - The DNS settings.
		* `name_servers` - The list of DNS server IP addresses.
		* `options` - The list of objects. Each object is a name-value pair. The value is optional.
			* `name` - The name of the object variable.
			* `value` - The value of the object variable.
		* `searches` - The list of DNS lookup domains.
	* `eci_security_context` - The security context of the container group.
	    * `sysctls` - The system information.
	        * `name` - The name of the variable.
	        * `value` - The value of the variable.
	* `eni_instance_id` - The ID of the ENI instance.
	* `events` - The events of the container group. Maximum: `50`.
		* `count` - The number of events.
		* `first_timestamp` - The time when the event started.
		* `last_timestamp` - The time when the event ended.
		* `message` - The content of the event.
		* `name` - The name of the object to which the event belongs.
		* `reason` - The name of the event.
		* `type` - The type of the event. Valid values: Normal and Warning.
	* `expired_time` - The time when the container group failed to run due to overdue payments. The timestamp follows the UTC and RFC3339 formats.
	* `failed_time` - The time when the container failed to run tasks. The timestamp follows the UTC and RFC3339 formats.
	* `host_aliases` - The mapping between host names and IP addresses for a container in the container group.
		* `hostnames` - The name of the host.
		* `ip` - The IP address of the container.
	* `id` - The ID of the Container Group.
	* `init_containers` -  A list of init containers. Each element contains the following attributes:
		* `image_pull_policy` - The policy for pulling an image.
		* `ports` - The exposed ports and protocols. Maximum: `100`.
			* `port` - The port number. Valid values: 1 to 65535.
			* `protocol` - Valid values: `TCP` and `UDP`.
		* `volume_mounts` - The list of volumes mounted to the container.
        	* `mount_path` - The directory of the mounted volume. Data under this directory will be overwritten by the data in the volume.
       		* `name` - The name of the volume. The name is the same as the volume you selected when you purchased the container.
    		* `read_only` - Default value: `false`.
		* `working_dir` - The working directory of the container.
		* `commands` - The commands run by the container.
		* `cpu` - The amount of CPU resources allocated to the container.
		* `environment_vars` - The environment variables.
			* `value` - The value of the variable.
			* `key` - The name of the variable.
		* `gpu` - The amount of GPU resources allocated to the container.
		* `memory` - The amount of memory resources allocated to the container.
		* `args` - The arguments passed to the commands.
		* `image` - The image of the container.
		* `restart_count` - The number of times that the container has restarted.
		* `name` - The name of the init container.
		* `ready` - Indicates whether the container is ready.
	* `instance_type` - The type of the ECS instance.
	* `internet_ip` - The public IP address of the container group.
	* `intranet_ip` - The internal IP address of the container group.
	* `ipv6_address` - The IPv6 address.
	* `memory` - The amount of memory resources allocated to the container group.
	* `ram_role_name` - The RAM role that the container group assumes. ECI and ECS share the same RAM role.
	* `resource_group_id` - The ID of the resource group to which the container group belongs. If you have not specified a resource group for the container group, it is added to the default resource group.
	* `restart_policy` - The restart policy of the container group.
	* `security_group_id` - The ID of the security group.
	* `status` - The status of container.
	* `succeeded_time` - The time when all containers in the container group completed running the specified tasks. The timestamp follows the UTC and RFC 3339 formats. For example, 2018-08-02T15:00:00Z.
	* `tags` - The tags attached to the container group. Each tag is a key-value pair. You can attach up to 20 tags to a container group.
		* `tag_key` - The key of the tag.
		* `tag_value` - The value of the tag.
	* `volumes` - The information about the mounted volume. You can mount up to 20 volumes.
		* `disk_volume_disk_id` - The ID of DiskVolume.
		* `disk_volume_fs_type` - The type of DiskVolume.
		* `flex_volume_driver` - The name of the FlexVolume driver.
		* `flex_volume_options` - The list of FlexVolume objects.
		* `nfs_volume_path` - The path to the NFS volume.
		* `nfs_volume_read_only` - Default value: `false`.
		* `nfs_volume_server` - The address of the NFS server.
		* `config_file_volume_config_file_to_paths` - The list of configuration file paths.
			* `content` - The content of the configuration file. Maximum size: 32 KB.
			* `path` - The relative file path.
		* `flex_volume_fs_type` - The type of the mounted file system. The default value is determined by the script of FlexVolume.
		* `name` - The name of the volume.
		* `type` - The type of the volume. Currently, the following types of volumes are supported: EmptyDirVolume, NFSVolume, ConfigFileVolume, and FlexVolume.
	* `vpc_id` - The if of vpc.
	* `vswitch_id` - The vswitch id.
	* `zone_id` - The IDs of the zones where the container groups are deployed. If this parameter is not set, the system automatically selects the zones. By default, no value is specified.
