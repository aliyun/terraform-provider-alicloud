---
subcategory: "Elastic Container Instance (ECI)"
layout: "alicloud"
page_title: "Alicloud: alicloud_eci_container_group"
sidebar_current: "docs-alicloud-resource-eci-container-group"
description: |-
  Provides a Alicloud ECI Container Group resource.
---

# alicloud_eci_container_group

Provides ECI Container Group resource.

For information about ECI Container Group and how to use it, see [What is Container Group](https://www.alibabacloud.com/help/en/elastic-container-instance/latest/api-eci-2018-08-08-createcontainergroup).

-> **NOTE:** Available since v1.111.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_eci_container_group&exampleId=f86e77f9-0e1a-75b0-7cf6-56e74d282bd54f7e8421&activeTab=example&spm=docs.r.eci_container_group.0.f86e77f90e&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-beijing"
}

variable "name" {
  default = "tf-example"
}

data "alicloud_eci_zones" "default" {
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.0.0.0/8"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.1.0.0/16"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_eci_zones.default.zones.0.zone_ids.0
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_eci_container_group" "default" {
  container_group_name = var.name
  cpu                  = 8.0
  memory               = 16.0
  restart_policy       = "OnFailure"
  security_group_id    = alicloud_security_group.default.id
  vswitch_id           = alicloud_vswitch.default.id
  auto_create_eip      = true
  tags = {
    Created = "TF",
    For     = "example",
  }
  containers {
    image             = "registry.cn-beijing.aliyuncs.com/eci_open/nginx:alpine"
    name              = "nginx"
    working_dir       = "/tmp/nginx"
    image_pull_policy = "IfNotPresent"
    commands          = ["/bin/sh", "-c", "sleep 9999"]
    volume_mounts {
      mount_path = "/tmp/example"
      read_only  = false
      name       = "empty1"
    }
    ports {
      port     = 80
      protocol = "TCP"
    }
    environment_vars {
      key   = "name"
      value = "nginx"
    }
    liveness_probe {
      period_seconds        = "5"
      initial_delay_seconds = "5"
      success_threshold     = "1"
      failure_threshold     = "3"
      timeout_seconds       = "1"
      exec {
        commands = ["cat /tmp/healthy"]
      }
    }
    readiness_probe {
      period_seconds        = "5"
      initial_delay_seconds = "5"
      success_threshold     = "1"
      failure_threshold     = "3"
      timeout_seconds       = "1"
      exec {
        commands = ["cat /tmp/healthy"]
      }
    }
  }
  init_containers {
    name              = "init-busybox"
    image             = "registry.cn-beijing.aliyuncs.com/eci_open/busybox:1.30"
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

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_eci_container_group&spm=docs.r.eci_container_group.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `container_group_name` - (Required, ForceNew) The name of the container group.
* `vswitch_id` - (Required, ForceNew) The ID of the VSwitch. Currently, container groups can only be deployed in VPC networks. The number of IP addresses in the VSwitch CIDR block determines the maximum number of container groups that can be created in the VSwitch. Before you can create an ECI instance, plan the CIDR block of the VSwitch.
**NOTE:** From version 1.208.0, You can specify up to 10 `vswitch_id`. Separate multiple vSwitch IDs with commas (,), such as vsw-***,vsw-***.  attribute `vswitch_id` updating diff will be ignored when you set multiple vSwitchIds, there is only one valid `vswitch_id` exists in the set vSwitchIds.
* `security_group_id` - (Required, ForceNew) The ID of the security group to which the container group belongs. Container groups within the same security group can access each other.
* `instance_type` - (Optional, ForceNew) The type of the ECS instance.
* `zone_id` - (Optional, ForceNew) The ID of the zone where you want to deploy the container group. If no value is specified, the system assigns a zone to the container group. By default, no value is specified.
* `cpu` - (Optional, ForceNew, Float) The amount of CPU resources allocated to the container group.
* `memory` - (Optional, ForceNew, Float) The amount of memory resources allocated to the container group.
* `ram_role_name` - (Optional, ForceNew) The RAM role that the container group assumes. ECI and ECS share the same RAM role.
* `resource_group_id` - (Optional) The ID of the resource group. **NOTE:** From version 1.208.0, `resource_group_id` can be modified.
* `restart_policy` - (Optional) The restart policy of the container group. Valid values: `Always`, `Never`, `OnFailure`.
* `auto_match_image_cache` - (Optional, ForceNew, Bool, Available since v1.166.0) Specifies whether to automatically match the image cache. Default value: `false`. Valid values: `true` and `false`.
* `plain_http_registry` - (Optional, Available since v1.170.0) The address of the self-built mirror warehouse. When creating an image cache from an image in a self-built image repository using the HTTP protocol, you need to configure this parameter so that the ECI uses the HTTP protocol to pull the image to avoid image pull failure due to different protocols.
* `insecure_registry` - (Optional, Available since v1.170.0) The address of the self-built mirror warehouse. When creating an image cache using an image in a self-built image repository with a self-signed certificate, you need to configure this parameter to skip certificate authentication to avoid image pull failure due to certificate authentication failure.
* `auto_create_eip` - (Optional, Bool, Available since v1.170.0) Specifies whether to automatically create an EIP and bind the EIP to the elastic container instance.
* `eip_bandwidth` - (Optional, Int, Available since v1.170.0) The bandwidth of the EIP. Default value: `5`.
* `eip_instance_id` - (Optional, Available since v1.170.0) The ID of the elastic IP address (EIP).
* `ephemeral_storage` - (Optional, ForceNew, Int, Available since v1.262.0) The size of the temporary storage space to add. Unit: GiB.
* `containers` - (Required, Set) The list of containers. See [`containers`](#containers) below.
* `init_containers` - (Optional, Set) The list of initContainers. See [`init_containers`](#init_containers) below.
* `dns_policy` - (Optional, ForceNew, Available since v1.232.0) The policy of DNS. Default value: `Default`. Valid values: `Default` and `None`.
* `dns_config` - (Optional, Set) The structure of dnsConfig. See [`dns_config`](#dns_config) below.
* `eci_security_context` - (Deprecated since 1.215.0, Optional, ForceNew, Set) The security context of the container group. See [`eci_security_context`](#eci_security_context) below.
* `security_context` - (Optional, ForceNew, Set, Available since v1.215.0) The security context of the container group. See [`security_context`](#security_context) below.
* `host_aliases` - (Optional, ForceNew, Set) HostAliases. See [`host_aliases`](#host_aliases) below.
* `volumes` - (Optional, Set) The list of volumes. See [`volumes`](#volumes) below.
* `image_registry_credential` - (Optional, Set, Available since v1.141.0) The image registry credential. See [`image_registry_credential`](#image_registry_credential) below.
* `acr_registry_info` - (Optional, ForceNew, Set, Available since v1.189.0) The ACR enterprise edition example properties. See [`acr_registry_info`](#acr_registry_info) below.
* `tags` - (Optional) A mapping of tags to assign to the resource.
  - Key: It can be up to 64 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It cannot be a null string.
  - Value: It can be up to 128 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It can be a null string.
* `termination_grace_period_seconds` - (Optional, ForceNew, Int, Available since v1.216.0) The buffer time during which the program handles operations before the program stops. Unit: seconds.
* `spot_strategy` - (Optional, ForceNew, Available since v1.216.0) Filter the results by ECI spot type. Valid values: `NoSpot`, `SpotWithPriceLimit` and `SpotAsPriceGo`. Default to `NoSpot`.
* `spot_price_limit` - (Optional, ForceNew, Available since v1.216.0) The maximum hourly price of the ECI spot instance.

### `acr_registry_info`

The acr_registry_info supports the following:

* `instance_id` - (Optional) The ACR enterprise edition example ID.
* `region_id` - (Optional, ForceNew) The ACR enterprise edition instance belongs to the region.
* `instance_name` - (Optional) The name of the ACR enterprise edition instance.
* `domains` - (Optional, List) The domain name of the ACR Enterprise Edition instance. Defaults to all domain names of the corresponding instance. Support specifying individual domain names, multiple separated by half comma.

### `image_registry_credential`

The image_registry_credential supports the following:

* `user_name` - (Required) The username used to log on to the image repository. It is required when `image_registry_credential` is configured.
* `password` - (Required) The password used to log on to the image repository. It is required when `image_registry_credential` is configured.
* `server` - (Required) The address of the image repository. It is required when `image_registry_credential` is configured.

### `volumes`

The volumes supports the following:

* `name` - (Optional, ForceNew) The name of the volume.
* `type` - (Optional) The type of the volume.
* `disk_volume_disk_id` - (Optional, ForceNew) The ID of DiskVolume.
* `disk_volume_fs_type` - (Optional, ForceNew) The system type of DiskVolume.
* `flex_volume_driver` - (Optional, ForceNew) The name of the FlexVolume driver.
* `flex_volume_fs_type` - (Optional, ForceNew) The type of the mounted file system. The default value is determined by the script of FlexVolume.
* `flex_volume_options` - (Optional, ForceNew) The list of FlexVolume objects. Each object is a key-value pair contained in a JSON string.
* `nfs_volume_path` - (Optional, ForceNew) The path to the NFS volume.
* `nfs_volume_server` - (Optional, ForceNew) The address of the NFS server.
* `nfs_volume_read_only` - (Optional, ForceNew, Bool) The nfs volume read only. Default value: `false`.
* `config_file_volume_config_file_to_paths` - (Optional, Set, ForceNew) The paths of the ConfigFile volume. See [`config_file_volume_config_file_to_paths`](#volumes-config_file_volume_config_file_to_paths) below.
-> **NOTE:** Every volumes mounted must have `name` and `type` attributes.

### `volumes-config_file_volume_config_file_to_paths`

The config_file_volume_config_file_to_paths supports the following:

* `content` - (Optional, ForceNew) The content of the configuration file. Maximum size: 32 KB.
* `path` - (Optional, ForceNew) The relative file path.

### `host_aliases`

The host_aliases supports the following:

* `ip` - (Optional, ForceNew) The IP address of the host.
* `hostnames` - (Optional, ForceNew, List) The information about the host.

### `eci_security_context`

The eci_security_context supports the following:

* `sysctls` - (Optional, ForceNew, Set) Sysctls hold a list of namespaced sysctls used for the pod. Pods with unsupported sysctls (by the container runtime) might fail to launch. See [`sysctls`](#eci_security_context-sysctls) below.

### `eci_security_context-sysctls`

The sysctls supports the following:

* `name` - (Optional, ForceNew) The name of the security context that the container group runs.
* `value` - (Optional, ForceNew) The variable value of the security context that the container group runs.

### `security_context`

The security_context supports the following:

* `sysctl` - (Optional, ForceNew, Set) Sysctls hold a list of namespaced sysctls used for the pod. Pods with unsupported sysctls (by the container runtime) might fail to launch. See [`sysctl`](#security_context-sysctl) below.

### `security_context-sysctl`

The sysctls supports the following:

* `name` - (Optional, ForceNew) The name of the security context that the container group runs.
* `value` - (Optional, ForceNew) The variable value of the security context that the container group runs.


### `dns_config`

The dns_config supports the following:

* `name_servers` - (Optional, List) The list of DNS server IP addresses.
* `searches` - (Optional, List) The list of DNS lookup domains.
* `options` - (Optional, Set) The structure of options. See [`options`](#dns_config-options) below.

### `dns_config-options`

The options supports the following:

* `name` - (Optional) The name of the object.
* `value` - (Optional) The value of the object.

### `init_containers`

The init_containers supports the following:

* `name` - (Optional, ForceNew) The name of the init container.
* `cpu` - (Optional, Float) The amount of CPU resources allocated to the container. Default value: `0`.
* `gpu` - (Optional, ForceNew, Int) The number GPUs. Default value: `0`.
* `memory` - (Optional, Float) The amount of memory resources allocated to the container. Default value: `0`.
* `image` - (Optional) The image of the container.
* `image_pull_policy` - (Optional) The restart policy of the image. Default value: `IfNotPresent`. Valid values: `Always`, `IfNotPresent`, `Never`.
* `working_dir` - (Optional) The working directory of the container.
* `commands` - (Optional, List) The commands run by the init container.
* `args` - (Optional, List) The arguments passed to the commands.
* `ports` - (Optional, ForceNew, Set) The structure of port. See [`ports`](#init_containers-ports) below.
* `environment_vars` - (Optional, Set) The structure of environmentVars. See [`environment_vars`](#init_containers-environment_vars) below.
* `volume_mounts` - (Optional, Set) The structure of volumeMounts. See [`volume_mounts`](#init_containers-volume_mounts) below.
* `ready` - (Optional, Available since v1.208.0) Indicates whether the container passed the readiness probe.
* `restart_count` - (Optional, Available since v1.208.0) The number of times that the container restarted.
* `security_context` - (Optional, Set, Available since v1.215.0) The security context of the container. See [`security_context`](#init_containers-security_context) below.

### `init_containers-ports`

The ports supports the following:

* `port` - (Optional, ForceNew, Int) The port number. Valid values: `1` to `65535`.
* `protocol` - (Optional, ForceNew) The type of the protocol. Valid values: `TCP` and `UDP`.

### `init_containers-environment_vars`

The environment_vars supports the following:

* `key` - (Optional) The name of the variable. The name can be 1 to 128 characters in length and can contain letters, digits, and underscores (_). It cannot start with a digit.
* `value` - (Optional) The value of the variable. The value can be 0 to 256 characters in length.
* `field_ref` - (Optional) The reference of the environment variable. See [`field_ref`](#init_containers-environment_vars-field_ref) below.

### `init_containers-environment_vars-field_ref`
* `field_path` - (Optional) The path of the reference.

### `init_containers-volume_mounts`

The volume_mounts supports the following:

* `mount_path` - (Optional) The directory of the mounted volume. Data under this directory will be overwritten by the data in the volume.
* `name` - (Optional, ForceNew) The name of the mounted volume.
* `read_only` - (Optional, Bool) Specifies whether the mount path is read-only. Default value: `false`.

### `init_containers-security_context`

The security_context supports the following:

* `capability` - (Optional, Available since v1.215.0) The permissions that you want to grant to the processes in the containers. See [`capability`](#init_containers-security_context-capability) below.
* `run_as_user` - (Optional, Long, Available since v1.215.0) The ID of the user who runs the container.

### `init_containers-security_context-capability`

The capability supports the following:
* `add` - (Optional, List, Available since v1.215.0) The permissions that you want to grant to the processes in the containers.


### `containers`

The containers supports the following:

* `name` - (Required, ForceNew) The name of the init container.
* `image` - (Required) The image of the container.
* `cpu` - (Optional, Float) The amount of CPU resources allocated to the container. Default value: `0`.
* `gpu` - (Optional, ForceNew, Int) The number GPUs. Default value: `0`.
* `memory` - (Optional, Float) The amount of memory resources allocated to the container. Default value: `0`.
* `image_pull_policy` - (Optional) The restart policy of the image. Default value: `IfNotPresent`. Valid values: `Always`, `IfNotPresent`, `Never`.
* `working_dir` - (Optional) The working directory of the container.
* `commands` - (Optional, List) The commands run by the init container.
* `args` - (Optional, List) The arguments passed to the commands.
* `ports` - (Optional, ForceNew, Set) The structure of port. See [`ports`](#containers-ports) below.
* `environment_vars` - (Optional, Set) The structure of environmentVars. See [`environment_vars`](#containers-environment_vars) below.
* `volume_mounts` - (Optional, Set) The structure of volumeMounts. See [`volume_mounts`](#containers-volume_mounts) below.
* `liveness_probe` - (Optional, Set, Available since v1.189.0) The health check of the container. See [`liveness_probe`](#containers-liveness_probe) below.
* `readiness_probe` - (Optional, Set, Available since v1.189.0) The health check of the container. See [`readiness_probe`](#containers-readiness_probe) below.
* `ready` - (Optional, Available since v1.208.0) Indicates whether the container passed the readiness probe.
* `restart_count` - (Optional, Available since v1.208.0) The number of times that the container restarted.
* `security_context` - (Optional, Set, Available since v1.215.0) The security context of the container. See [`security_context`](#containers-security_context) below.
* `lifecycle_pre_stop_handler_exec` - (Optional, List, Available since v1.216.0) The commands to be executed in containers when you use the CLI to specify the preStop callback function.

### `containers-ports`

The ports supports the following:

* `port` - (Optional, ForceNew, Int) The port number. Valid values: `1` to `65535`.
* `protocol` - (Optional, ForceNew) The type of the protocol. Valid values: `TCP` and `UDP`.

### `containers-environment_vars`

The environment_vars supports the following:

* `key` - (Optional) The name of the variable. The name can be 1 to 128 characters in length and can contain letters, digits, and underscores (_). It cannot start with a digit.
* `value` - (Optional) The value of the variable. The value can be 0 to 256 characters in length.
* `field_ref` - (Optional) The reference of the environment variable. See [`field_ref`](#containers-environment_vars-field_ref) below.

### `containers-environment_vars-field_ref`
* `field_path` - (Optional) The path of the reference.

### `containers-volume_mounts`

The volume_mounts supports the following:

* `mount_path` - (Optional) The directory of the mounted volume. Data under this directory will be overwritten by the data in the volume.
* `name` - (Optional, ForceNew) The name of the mounted volume.
* `read_only` - (Optional, Bool) Specifies whether the volume is read-only. Default value: `false`.

### `containers-liveness_probe`

The liveness_probe supports the following:

* `initial_delay_seconds` - (Optional, Int) Check the time to start execution, calculated from the completion of container startup.
* `period_seconds` - (Optional, Int) Buffer time for the program to handle operations before closing.
* `timeout_seconds` - (Optional, Int) Check the timeout, the default is 1 second, the minimum is 1 second.
* `success_threshold` - (Optional, ForceNew, Int) The check count threshold for re-identifying successful checks since the last failed check (must be consecutive successes), default is 1. Current must be 1.
* `failure_threshold` - (Optional, Int) Threshold for the number of checks that are determined to have failed since the last successful check (must be consecutive failures), default is 3.
* `exec` - (Optional, Set) Health check using command line method. See [`exec`](#containers-liveness_probe-exec) below.
* `tcp_socket` - (Optional, Set) Health check using TCP socket method. See [`tcp_socket`](#containers-liveness_probe-tcp_socket) below.
* `http_get` - (Optional, Set) Health check using HTTP request method. See [`http_get`](#containers-liveness_probe-http_get) below.

-> **NOTE:** When you configure `liveness_probe`, you can select only one of the `exec`, `tcp_socket`, `http_get`.

### `containers-liveness_probe-exec`

The exec supports the following:

* `commands` - (Optional, List) Commands to be executed inside the container when performing health checks using the command line method.

### `containers-liveness_probe-tcp_socket`

The tcp_socket supports the following:

* `port` - (Optional, Int) The port for TCP socket detection when using the TCP socket method for health check.

### `containers-liveness_probe-http_get`

The http_get supports the following:

* `scheme` - (Optional) The protocol type corresponding to the HTTP Get request when using the HTTP request method for health checks. Valid values: `HTTP`, `HTTPS`.
* `port` - (Optional, Int) When using the HTTP request method for health check, the port number for HTTP Get request detection.
* `path` - (Optional) The path of HTTP Get request detection when setting the postStart callback function using the HTTP request method.

### `containers-readiness_probe`

The readiness_probe supports the following:

* `initial_delay_seconds` - (Optional, Int) Check the time to start execution, calculated from the completion of container startup.
* `period_seconds` - (Optional, Int) Buffer time for the program to handle operations before closing.
* `timeout_seconds` - (Optional, Int) Check the timeout, the default is 1 second, the minimum is 1 second.
* `success_threshold` - (Optional, ForceNew, Int) The check count threshold for re-identifying successful checks since the last failed check (must be consecutive successes), default is 1. Current must be 1.
* `failure_threshold` - (Optional, Int) Threshold for the number of checks that are determined to have failed since the last successful check (must be consecutive failures), default is 3.
* `exec` - (Optional) Health check using command line method. See [`exec`](#containers-readiness_probe-exec) below.
* `tcp_socket` - (Optional) Health check using TCP socket method. See [`tcp_socket`](#containers-readiness_probe-tcp_socket) below.
* `http_get` - (Optional) Health check using HTTP request method. See [`http_get`](#containers-readiness_probe-http_get) below.

-> **NOTE:** When you configure `readiness_probe`, you can select only one of the `exec`, `tcp_socket`, `http_get`.

### `containers-readiness_probe-exec`

The exec supports the following:

* `commands` - (Optional, List) Commands to be executed inside the container when performing health checks using the command line method.

### `containers-readiness_probe-tcp_socket`

The tcp_socket supports the following:

* `port` - (Optional, Int) The port for TCP socket detection when using the TCP socket method for health check.

### `containers-readiness_probe-http_get`

The http_get supports the following:

* `scheme` - (Optional) The protocol type corresponding to the HTTP Get request when using the HTTP request method for health checks. Valid values: `HTTP`, `HTTPS`.
* `port` - (Optional, Int) When using the HTTP request method for health check, the port number for HTTP Get request detection.
* `path` - (Optional) The path of HTTP Get request detection when setting the postStart callback function using the HTTP request method.

### `containers-security_context`

The security_context supports the following:

* `capability` - (Optional, Available since v1.215.0) The permissions that you want to grant to the processes in the containers. See [`capability`](#containers-security_context-capability) below.
* `run_as_user` - (Optional, Long, Available since v1.215.0) The ID of the user who runs the container.
* `privileged` - (Optional, ForceNew, Bool, Available since v1.225.1) Specifies whether to give extended privileges to this container. Default value: `false`. Valid values: `true` and `false`.

### `containers-security_context-capability`

The capability supports the following:
* `add` - (Optional, List, Available since v1.215.0) The permissions that you want to grant to the processes in the containers.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Container Group. Value as `container_group_id`.
* `internet_ip` - (Available since v1.170.0) The Public IP of the container group.
* `intranet_ip` - (Available since v1.170.0) The Private IP of the container group.
* `status` - The status of container group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 20 mins) Used when create the Container Group.
* `update` - (Defaults to 20 mins) Used when update the Container Group.
* `delete` - (Defaults to 5 mins) Used when delete the Container Group.

## Import

ECI Container Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_eci_container_group.example <container_group_id>
```
