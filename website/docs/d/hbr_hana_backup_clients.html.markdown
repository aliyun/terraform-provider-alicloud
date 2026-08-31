---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_hana_backup_clients"
sidebar_current: "docs-alicloud-datasource-hbr-hana-backup-clients"
description: |-
  Provides a list of Hybrid Backup Recovery (HBR) Hana Backup Clients to the user.
---

# alicloud\_hbr\_hana\_backup\_clients

This data source provides the Hybrid Backup Recovery (HBR) Hana Backup Clients of the current Alibaba Cloud user.

-> **NOTE:** Available in 1.198.0+

## Example Usage

Basic Usage

```terraform
data "alicloud_hbr_hana_backup_clients" "ids" {
  ids      = ["example_id"]
  vault_id = "your_vault_id"
}

output "hbr_hana_backup_clients_id_1" {
  value = data.alicloud_hbr_hana_backup_clients.ids.hana_backup_clients.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed) A list of Hana Backup Client IDs.
* `vault_id` - (Required, ForceNew) The ID of the backup vault.
* `client_id` - (Optional, ForceNew) The ID of the backup client.
* `cluster_id` - (Optional, ForceNew) The ID of the SAP HANA instance.
* `status` - (Optional, ForceNew) The status of the Hana Backup Client. Valid Values: `REGISTERED`, `ACTIVATED`, `DEACTIVATED`, `INSTALLING`, `INSTALL_FAILED`, `NOT_INSTALLED`, `UPGRADING`, `UPGRADE_FAILED`, `UNINSTALLING`, `UNINSTALL_FAILED`, `STOPPED`, `UNKNOWN`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `hana_backup_clients` - A list of Hana Backup Clients. Each element contains the following attributes:
  * `id` - The id of the Hana Backup Client. It formats as `<vault_id>:<client_id>`.
  * `vault_id` - The ID of the backup vault.
  * `client_id` - The ID of the backup client.
  * `client_name` - The name of the backup client.
  * `client_type` - The type of the backup client.
  * `client_version` - The version number of the backup client.
  * `max_version` - The maximum version number of the backup client.
  * `cluster_id` - The ID of the SAP HANA instance.
  * `instance_id` - The ID of the instance.
  * `instance_name` - The name of the ECS instance.
  * `alert_setting` - The alert settings.
  * `use_https` - Indicates whether data is transmitted over HTTPS.
  * `network_type` - The network type.
  * `status_message` - The status information.
  * `status` - The status of the backup client.
  