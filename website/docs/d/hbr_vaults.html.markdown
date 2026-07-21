---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_vaults"
sidebar_current: "docs-alicloud-datasource-hbr-vaults"
description: |-
  Provides a list of Hybrid Backup Recovery (HBR) Backup vaults to the user.
---

# alicloud\_hbr\_vaults

This data source provides the Hbr Vaults of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.129.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_hbr_vaults" "ids" {
  name_regex = "^my-Vault"
}

output "hbr_vault_id_1" {
  value = data.alicloud_hbr_vaults.ids.vaults.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Vault IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Vault name.
* `vault_type` - (Optional, ForceNew) VaultType. Valid values: `STANDARD`,`OTS_BACKUP`.
  - `STANDARD` - used in OSS, NAS and ECS File backup.
  - `OTS_BACKUP` -  used in OTS backup.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew, Computed) The status of Vault. Valid values: `CREATED`, `ERROR`, `UNKNOWN`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of vault names.
* `vaults` - A list of Hbr vaults. Each element contains the following attributes:
    * `id` - The ID of vault.
    * `vault_id` - The ID of vault, same as `id`.
    * `vault_name` - The name of vault.
    * `description` - The description of the vault.
    * `vault_storage_class` - The storage class of vault. Valid values: `STANDARD`.
    * `vault_type` - The type of Vault. Valid values: `STANDARD`,`OTS_BACKUP`.
    * `vault_status_message` - Error status information of Vault. Only valid for remote backup warehouses. Only the remote backup warehouse is valid.
    * `storage_size` - Backup vault storage usage. The unit is Byte.
    * `bucket_name` - The name of the OSS bucket of the Vault.
    * `bytes_done` - The amount of backup data. The unit is Byte.
    * `created_time` - The creation time of the Vault. UNIX time in seconds.
    * `updated_time` - The update time of the Vault. UNIX time in seconds.
    * `latest_replication_time` - The time of the last remote backup synchronization.
    * `status` - The status of Vault. Valid values: `CREATED`, `ERROR`, `UNKNOWN`. 
    * `payment_type` - Billing model, possible values:
        * `FREE` is not billed
        * `V1` common vault billing model, including back-end storage capacity, client licenses and other billing items
        * `V2` new version of metering mode
        * `AEGIS` Billing method for cloud security use
        * `UNI_BACKUP` the backup of deduplication database
        * `ARCHIVE` archive library.
    * `replication` - Whether it is a remote backup warehouse. It's a boolean value.
    * `replication_source_region_id` - The region ID to which the remote backup Vault belongs.
    * `replication_source_vault_id` - The source vault ID of the remote backup Vault.
    * `dedup` - (Internal use) Whether to enable the deduplication function for the database backup Vault.
    * `retention` - (Not yet open) Warehouse-level data retention days, only valid for archive libraries.
    * `search_enabled` - (Not yet open) Whether to enable the backup search function.
    * `index_available` - (Not yet open) Index available.
    * `index_level` - (Not yet open) Index level.
    * `index_update_time` - (Not yet open) Index update time.


