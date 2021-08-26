---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_ecs_backup_clients"
sidebar_current: "docs-alicloud-datasource-hbr-ecs-backup-clients"
description: |-
  Provides a list of Hbr Ecs Backup Clients to the user.
---

# alicloud\_hbr\_ecs\_backup\_clients

This data source provides the Hbr Ecs Backup Clients of the current Alibaba Cloud user.

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
* `status` - (Optional, ForceNew) The status of the resource. Valid values: `ACTIVATED`, `STOPPED`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `clients` - A list of Hbr Ecs Backup Clients. Each element contains the following attributes:
	* `id` - The ID of the Ecs Backup Client.
	* `instance_id` - The ID of ECS Instance.
	* `instance_name` - ECS Instance Names.
	* `arch_type` - The Client System Architecture (Only the ECS File Backup Client Is Available. Valid Values: `AMD64` , `386`.
	* `backup_status` - Client protected status.
	* `client_type` - The Client Type. Valid Values: `ECS_CLIENT` (ECS File Backup Client).
	* `client_version` - Client Version.
	* `create_time` - The Client Creates a Time. Unix Time Seconds.
	* `data_network_type` - The Data Plane Data Access Point Type. Valid Values: `PUBLIC`, `VPC`, `CLASSIC`.
	* `data_proxy_setting` - The Data Plane Proxy Settings. Valid Values: `DISABLE`, `USE_CONTROL_PROXY`, `CUSTOM`. **Note**: `USE_CONTROL_PROXY` (Default, the same with Control Plane), `CUSTOM` (Custom Configuration Items for the HTTP Protocol).
	* `ecs_backup_client_id` - The first ID of the resource.
	* `hostname` - The ECS Host Name.
	* `last_heart_beat_time` - Client Last Heartbeat Time. Unix Time Seconds.
	* `max_client_version` - The Latest Client Version.
	* `max_cpu_core` - A Single Backup Task Uses for Example, Instances Can Be Grouped According to CPU Core Count, 0 Means No Restrictions.
	* `max_worker` - A Single Backup Task Parallel Work, the Number of 0 Means No Restrictions.
	* `os_type` - The Client System Type (Only the ECS File Backup Client Is Available. Possible Values: * windows * linux.
	* `private_ipv4` - Instance Must Not Use the Intranet IP Address.
	* `proxy_host` - Custom Data Plane Proxy Server Host Address.
	* `proxy_password` - Custom Data Plane Proxy Password.
	* `proxy_port` - Custom Data Plane Proxy Server Host Port.
	* `proxy_user` - Custom Data Plane Proxy Server User Name.
	* `status` - The status of the resource.
	* `updated_time` - Client Update Time. Unix Time Seconds.
	* `use_https` - Indicates Whether to Use the Https Transport Data Plane Data.
	* `zone_id` - The Zone ID.