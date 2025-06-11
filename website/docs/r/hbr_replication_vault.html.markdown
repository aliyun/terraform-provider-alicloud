---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_replication_vault"
description: |-
  Provides a Alicloud Hybrid Backup Recovery (HBR) Replication Vault resource.
---

# alicloud_hbr_replication_vault

Provides a Hybrid Backup Recovery (HBR) Replication Vault resource.

The replication vault used by the cross-region backup function of Cloud Backup.

For information about Hybrid Backup Recovery (HBR) Replication Vault and how to use it, see [What is Replication Vault](https://www.alibabacloud.com/help/en/doc-detail/345603.html).

-> **NOTE:** Available since v1.252.0.

## Example Usage

Basic Usage

```terraform
variable "source_region" {
  default = "cn-hangzhou"
}

provider "alicloud" {
  alias  = "source"
  region = var.source_region
}

data "alicloud_hbr_replication_vault_regions" "default" {}

provider "alicloud" {
  alias  = "replication"
  region = data.alicloud_hbr_replication_vault_regions.default.regions.0.replication_region_id
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_hbr_vault" "default" {
  provider   = alicloud.source
  vault_name = "terraform-example-${random_integer.default.result}"
}

resource "alicloud_hbr_replication_vault" "default" {
  provider                     = alicloud.replication
  replication_source_region_id = var.source_region
  replication_source_vault_id  = alicloud_hbr_vault.default.id
  vault_name                   = "terraform-example"
  vault_storage_class          = "STANDARD"
  description                  = "terraform-example"
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Optional) The description of the backup vault.
* `encrypt_type` - (Optional, ForceNew, Computed, Available since v1.252.0) The encryption type of the backup vault.
* `kms_key_id` - (Optional, ForceNew, Available since v1.252.0) Alibaba Cloud KMS custom Key or Alias. This parameter is required only when EncryptType = KMS.
* `replication_source_region_id` - (Required, ForceNew) The region ID of the source backup vault.
* `replication_source_vault_id` - (Required, ForceNew) The vault ID of the source backup vault.
* `vault_name` - (Required) The name of the backup vault.
* `vault_storage_class` - (Optional, ForceNew, Computed) Backup Vault Storage Class

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `region_id` - RegionId
* `status` - The status of the mirror backup vault.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Replication Vault.
* `delete` - (Defaults to 15 mins) Used when delete the Replication Vault.
* `update` - (Defaults to 5 mins) Used when update the Replication Vault.

## Import

Hybrid Backup Recovery (HBR) Replication Vault can be imported using the id, e.g.

```shell
$ terraform import alicloud_hbr_replication_vault.example <id>
```