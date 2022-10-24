---
subcategory: "Elastic Container Instance (ECI)"
layout: "alicloud"
page_title: "Alicloud: alicloud_eci_container_group"
sidebar_current: "docs-alicloud-resource-eci-container-group"
description: |-
  Provides a Alicloud ECI Container Group resource.
---

# alicloud\_eci\_container\_group

Provides ECI Container Group resource.

For information about ECI Container Group and how to use it, see [What is Container Group](https://www.alibabacloud.com/help/en/doc-detail/90341.htm).

-> **NOTE:** Available in v1.111.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_eci_container_group" "example" {
  container_group_name = "tf-testacc-eci-gruop"
  cpu                  = 8.0
  memory               = 16.0
  restart_policy       = "OnFailure"
  security_group_id    = alicloud_security_group.group.id
  vswitch_id           = data.alicloud_vpcs.default.vpcs.0.vswitch_ids.0
  tags = {
    TF = "create"
  }

  containers {
    image             = "registry-vpc.cn-beijing.aliyuncs.com/eci_open/nginx:alpine"
    name              = "nginx"
    working_dir       = "/tmp/nginx"
    image_pull_policy = "IfNotPresent"
    commands          = ["/bin/sh", "-c", "sleep 9999"]
    volume_mounts {
      mount_path = "/tmp/test"
      read_only  = false
      name       = "empty1"
    }
    ports {
      port     = 80
      protocol = "TCP"
    }
    environment_vars {
      key   = "test"
      value = "nginx"
    }
  }
  containers {
    image    = "registry-vpc.cn-beijing.aliyuncs.com/eci_open/centos:7"
    name     = "centos"
    commands = ["/bin/sh", "-c", "sleep 9999"]
  }
  init_containers {
    name              = "init-busybox"
    image             = "registry-vpc.cn-beijing.aliyuncs.com/eci_open/busybox:1.30"
    image_pull_policy = "IfNotPresent"
    commands          = ["echo"]
    args              = ["hello initcontainer"]
  }
  volumes {
    name = "empty1"
    type = "EmptyDirVolume"
  }
  volumes {
    name = "empty2"
    type = "EmptyDirVolume"
  }
}
```

## Argument Reference

The following arguments are supported:

* `container_group_name` - (Required, ForceNew) The name of the container group.
* `containers` - (Required) The list of containers.
* `cpu` - (Optional, Computed) The amount of CPU resources allocated to the container group.
* `dns_config` - (Optional) The structure of dnsConfig.
* `eci_security_context` - (Optional) The security context of the container group.
* `host_aliases` - (Optional, ForceNew) HostAliases.
* `init_containers` - (Optional) The list of initContainers.
* `instance_type` - (Optional, ForceNew) The type of the ECS instance.
* `memory` - (Optional, Computed) The amount of memory resources allocated to the container group.
* `ram_role_name` - (Optional, ForceNew) The RAM role that the container group assumes. ECI and ECS share the same RAM role.
* `resource_group_id` - (Optional, Computed, ForceNew) The ID of the resource group.
* `restart_policy` - (Optional, Computed) The restart policy of the container group. Valid values: `Always`, `Never`, `OnFailure`.
* `security_group_id` - (Required, ForceNew) The ID of the security group to which the container group belongs. Container groups within the same security group can access each other.
* `volumes` - (Optional) The list of volumes.
* `vswitch_id` - (Required, ForceNew) The ID of the VSwitch. Currently, container groups can only be deployed in VPC networks. The number of IP addresses in the VSwitch CIDR block determines the maximum number of container groups that can be created in the VSwitch. Before you can create an ECI instance, plan the CIDR block of the VSwitch.
* `zone_id` - (Optional, Computed, ForceNew) The ID of the zone where you want to deploy the container group. If no value is specified, the system assigns a zone to the container group. By default, no value is specified.
* `image_registry_credential` - (Optional, Available in 1.141.0+) The image registry credential. The details see Block `image_registry_credential`.
* `auto_match_image_cache` - (Optional, Available in 1.166.0+) Specifies whether to automatically match the image cache. Default value: false.
* `insecure_registry` - (Optional, Available in 1.170.0+) The address of the self-built mirror warehouse. When creating an image cache using an image in a self-built image repository with a self-signed certificate, you need to configure this parameter to skip certificate authentication to avoid image pull failure due to certificate authentication failure.
* `plain_http_registry` - (Optional, Available in 1.170.0+) The address of the self-built mirror warehouse. When creating an image cache from an image in a self-built image repository using the HTTP protocol, you need to configure this parameter so that the ECI uses the HTTP protocol to pull the image to avoid image pull failure due to different protocols.
* `auto_create_eip` - (Optional, Available in 1.170.0+) Specifies whether to automatically create an EIP and bind the EIP to the elastic container instance.
* `eip_bandwidth` - (Optional, Available in 1.170.0+) The bandwidth of the EIP. The default value is `5`.
* `eip_instance_id` - (Optional, Available in 1.170.0+) The ID of the elastic IP address (EIP).
* `tags` - (Optional) A mapping of tags to assign to the resource.
  - Key: It can be up to 64 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It cannot be a null string.
  - Value: It can be up to 128 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It can be a null string.
  
#### Block volumes

The volumes supports the following: 
* `name` - (Optional) The name of the volume.
* `type` - (Optional) The type of the volume.
* `config_file_volume_config_file_to_paths` - (Optional) ConfigFileVolumeConfigFileToPaths.

* `disk_volume_disk_id` - (Optional, ForceNew) The ID of DiskVolume.
* `disk_volume_fs_type` - (Optional, ForceNew) The system type of DiskVolume.

* `flex_volume_driver` - (Optional, ForceNew) The name of the FlexVolume driver.
* `flex_volume_fs_type` - (Optional, ForceNew) The type of the mounted file system. The default value is determined by the script of FlexVolume.
* `flex_volume_options` - (Optional, ForceNew) The list of FlexVolume objects. Each object is a key-value pair contained in a JSON string.

* `nfs_volume_path` - (Optional) The path to the NFS volume.
* `nfs_volume_read_only` - (Optional) The nfs volume read only. Default to `false`.
* `nfs_volume_server` - (Optional) The address of the NFS server.

-> **NOTE:** Every volumes mounted must have name and type attributes.

#### Block config_file_volume_config_file_to_paths

The config_file_volume_config_file_to_paths supports the following: 

* `content` - (Optional) The content of the configuration file. Maximum size: 32 KB.
* `path` - (Optional) The relative file path.

#### Block init_containers

The init_containers supports the following: 

* `args` - (Optional) The arguments passed to the commands.
* `commands` - (Optional) The commands run by the init container.
* `cpu` - (Optional) The amount of CPU resources allocated to the container.
* `environment_vars` - (Optional) The structure of environmentVars.
* `gpu` - (Optional) The number GPUs.
* `image` - (Optional) The image of the container.
* `image_pull_policy` - (Optional) The restart policy of the image.
* `memory` - (Optional) The amount of memory resources allocated to the container.
* `name` - (Optional) The name of the init container.
* `ports` - (Optional, ForceNew) The structure of port.
* `volume_mounts` - (Optional) The structure of volumeMounts.
* `working_dir` - (Optional) The working directory of the container.

#### Block volume_mounts in init_containers

The volume_mounts supports the following: 

* `mount_path` - (Optional) The directory of the mounted volume. Data under this directory will be overwritten by the data in the volume.
* `name` - (Optional) The name of the mounted volume.
* `read_only` - (Optional) Default to `false`.

#### Block ports in init_containers

The ports supports the following: 

* `port` - (Optional, ForceNew) The port number. Valid values: 1 to 65535.
* `protocol` - (Optional, ForceNew) Valid values: TCP and UDP.

#### Block environment_vars in init_containers

The environment_vars supports the following: 

* `key` - (Optional) The name of the variable. The name can be 1 to 128 characters in length and can contain letters, digits, and underscores (_). It cannot start with a digit.
* `value` - (Optional) The value of the variable. The value can be 0 to 256 characters in length.

#### Block host_aliases

The host_aliases supports the following: 

* `hostnames` - (Optional, ForceNew) Adds a host name.
* `ip` - (Optional, ForceNew) Adds an IP address.

#### Block image_registry_credential
The image_registry_credential supports the following:
* `password` - (Optional) The password used to log on to the image repository. It is required when `image_registry_credential` is configured.
* `server` - (Optional) The address of the image repository. It is required when `image_registry_credential` is configured.
* `user_name` - (Optional) The username used to log on to the image repository. It is required when `image_registry_credential` is configured.

#### Block dns_config

The dns_config supports the following: 

* `name_servers` - (Optional) The list of DNS server IP addresses.
* `options` - (Optional) The structure of options.
* `searches` - (Optional) The list of DNS lookup domains.

#### Block options

The options supports the following: 

* `name` - (Optional) The name of the object.
* `value` - (Optional) The value of the object.

#### Block containers

The containers supports the following: 

* `args` - (Optional) The arguments passed to the commands.
* `commands` - (Optional) The commands run by the init container.
* `cpu` - (Optional) The amount of CPU resources allocated to the container.
* `environment_vars` - (Optional) The structure of environmentVars.
* `gpu` - (Optional) The number GPUs.
* `image` - (Required) The image of the container.
* `image_pull_policy` - (Optional) The restart policy of the image.
* `memory` - (Optional) The amount of memory resources allocated to the container.
* `name` - (Required) The name of the init container.
* `ports` - (Optional, ForceNew) The structure of port.
* `volume_mounts` - (Optional) The structure of volumeMounts.
* `working_dir` - (Optional) The working directory of the container.

#### Block volume_mounts

The volume_mounts supports the following: 

* `mount_path` - (Optional) The directory of the mounted volume. Data under this directory will be overwritten by the data in the volume.
* `name` - (Optional) The name of the mounted volume.
* `read_only` - (Optional) Default to `false`.

#### Block ports

The ports supports the following: 

* `port` - (Optional, ForceNew) The port number. Valid values: 1 to 65535.
* `protocol` - (Optional, ForceNew) Valid values: TCP and UDP.

#### Block environment_vars

The environment_vars supports the following: 

* `key` - (Optional) The name of the variable. The name can be 1 to 128 characters in length and can contain letters, digits, and underscores (_). It cannot start with a digit.
* `value` - (Optional) The value of the variable. The value can be 0 to 256 characters in length.

#### Block eci_security_context

The eci_security_context supports the following:
* `sysctls` - (Optional) system.

#### Block sysctls

The sysctls supports the following:

* `name` - (Optional, ForceNew) The name of the security context that the container group runs.
* `value` - (Optional, ForceNew) The variable value of the security context that the container group runs.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Container Group. Value as `container_group_id`.
* `status` - The status of container group.
* `internet_ip` - (Available in v1.170.0+) The Public IP of the container group.
* `intranet_ip` - (Available in v1.170.0+) The Private IP of the container group.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 20 mins) Used when create the Container Group.
* `update` - (Defaults to 20 mins) Used when update the Container Group.

## Import

ECI Container Group can be imported using the id, e.g.

```
$ terraform import alicloud_eci_container_group.example <container_group_id>
```
