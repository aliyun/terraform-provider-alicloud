---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_ecs_backup_clients"
sidebar_current: "docs-alicloud-datasource-hbr-ecs-backup-clients"
description: |-
  Provides a list of Hybrid Backup Recovery (HBR) Ecs Backup Clients to the user.
---

# alicloud\_hbr\_ecs\_backup\_clients

This data source provides the Hbr Ecs File Backup Clients of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.132.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_instances" "default" {
  name_regex = "ecs_instance_name"
  status     = "Running"
}

data "alicloud_hbr_ecs_backup_clients" "ids" {
  ids          = [alicloud_hbr_ecs_backup_client.default.id]
  instance_ids = [alicloud_hbr_ecs_backup_client.default.instance_id]
}

output "hbr_ecs_backup_client_id_1" {
  value = data.alicloud_hbr_ecs_backup_clients.ids.clients.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Ecs Backup Client IDs.
* `instance_ids` - (Optional, ForceNew) A list of ECS Instance IDs.
* `status` - (Optional, ForceNew) The status of the resource. Valid values: `ACTIVATED`, `DEACTIVATED`, `INSTALLING`, `INSTALL_FAILED`, `NOT_INSTALLED`, `REGISTERED`, `STOPPED`, `UNINSTALLING`, `UNINSTALL_FAILED`, `UNKNOWN`, `UPGRADE_FAILED`, `UPGRADING`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `clients` - A list of Hbr Ecs Backup Clients. Each element contains the following attributes:
    * `id` - The ID of the Ecs Backup Client.
    * `instance_id` - The ID of ECS instance. When the client type is ECS file backup client, it indicates the ID of ECS instance. When the client type is a local file backup client, it is a hardware fingerprint generated based on system information.
    * `instance_name` - The name of ECS instance.
    * `arch_type` - The system architecture of client, only the ECS File Backup Client is available. Valid values: `AMD64` , `386`.
    * `backup_status` - Client protected status. Valid values: `UNPROTECTED`, `PROTECTED`.
    * `client_type` - The type of client. Valid values: `ECS_CLIENT` (ECS File Backup Client).
    * `client_version` - The version of client.
    * `create_time` - The creation time of client. Unix time in seconds.
    * `data_network_type` - The data plane access point type. Valid Values: `PUBLIC`, `VPC`, `CLASSIC`.
    * `data_proxy_setting` - The data plane proxy settings. Valid Values: `DISABLE`, `USE_CONTROL_PROXY`, `CUSTOM`.
        * `USE_CONTROL_PROXY` (Default, the same with control plane)
        * `CUSTOM` (Custom configuration items for the HTTP protocol).
    * `ecs_backup_client_id` - The first ID of the resource.
    * `hostname` - The hostname of ECS instance.
    * `last_heart_beat_time` - The last heartbeat time of client. Unix Time Seconds.
    * `max_client_version` - The latest version of client.
    * `max_cpu_core` - The number of CPU cores used by a single backup task, 0 means no restrictions.
    * `max_worker` - The number of concurrent jobs for a single backup task, 0 means no restrictions.
    * `os_type` - The operating system type of client, only the ECS File Backup Client is available. Valid values: `windows`, `linux`.
    * `private_ipv4` - Intranet IP address of the instance, only available for ECS file backup client.
    * `proxy_host` - Custom data plane proxy server host address.
    * `proxy_port` - Custom data plane proxy server host port.
    * `proxy_user` - The username of custom data plane proxy server.
    * `proxy_password` -  The password of custom data plane proxy server.
    * `status` - The status of the resource.
    * `updated_time` - The update time of client. Unix Time Seconds.
    * `use_https` - Indicates whether to use the HTTPS protocol. Valid values: `true`, `false`.
    * `zone_id` - The ID of Zone.
