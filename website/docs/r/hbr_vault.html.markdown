---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_vault"
sidebar_current: "docs-alicloud-resource-hbr-vault"
description: |-
  Provides a Alicloud Hybrid Backup Recovery (HBR) Backup vault resource.
---

# alicloud_hbr_vault

Provides a HBR Backup vault resource.

For information about HBR Backup vault and how to use it, see [What is Backup vault](https://www.alibabacloud.com/help/en/hybrid-backup-recovery/latest/api-hbr-2017-09-08-createvault).

-> **NOTE:** Available since v1.129.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_hbr_vault&exampleId=3769646c-17aa-6388-ab06-042ff80e0d1e4978e73c&activeTab=example&spm=docs.r.hbr_vault.0.3769646c17&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "random_integer" "default" {
  min = 10000
  max = 99999
}
resource "alicloud_hbr_vault" "example" {
  vault_name = "example_value_${random_integer.default.result}"
}
```

## Argument Reference

The following arguments are supported:

* `vault_name` - (Required) The name of Vault.
* `vault_type` - (Optional, ForceNew) The type of Vault. Valid values: `STANDARD`, `OTS_BACKUP`.
* `vault_storage_class` - (Optional, ForceNew) The storage class of Vault. Valid values: `STANDARD`.
* `encrypt_type` - (Optional, ForceNew, Available since v1.173.0) Source Encryption Type，It is valid only when vault_type is `STANDARD` or `OTS_BACKUP`. Default value: `HBR_PRIVATE`. Valid values:
  - `HBR_PRIVATE`: HBR is fully hosted, uses the backup service's own encryption method.
  - `KMS`: Use Alibaba Cloud Kms to encryption.
* `kms_key_id` - (Optional, ForceNew, Available since v1.173.0) The key id or alias name of Alibaba Cloud Kms. It is required and valid only when encrypt_type is `KMS`.
* `description` - (Optional) The description of Vault. Defaults to an empty string.
* `redundancy_type` - (Removed since v1.209.1) The redundancy type of the vault. **NOTE:** Field `redundancy_type` has been removed from provider version 1.209.1.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Vault.
* `status` - The status of the Vault.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 min) Used when create the Vault.
* `update` - (Defaults to 1 min) Used when update the Vault.
* `delete` - (Defaults to 1 min) Used when delete the Vault.

## Import

HBR Vault can be imported using the id, e.g.

```shell
$ terraform import alicloud_hbr_vault.example <id>
```
