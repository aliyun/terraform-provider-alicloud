---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_replication_vault"
sidebar_current: "docs-alicloud-resource-hbr-replication-vault"
description: |-
  Provides a Alicloud Hybrid Backup Recovery (HBR) Replication Vault resource.
---

# alicloud\_hbr\_replication\_vault

Provides a Hybrid Backup Recovery (HBR) Replication Vault resource.

For information about Hybrid Backup Recovery (HBR) Replication Vault and how to use it, see [What is Replication Vault](https://www.alibabacloud.com/help/en/doc-detail/345603.html).

-> **NOTE:** Available in v1.152.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-testAccReplicationVault"
}

variable "region_source" {
  default = "you Replication value source region"
}

provider "alicloud" {
  alias  = "source"
  region = var.region_source
}

resource "alicloud_hbr_vault" "default" {
  vault_name = var.name
  provider   = alicloud.source
}

data "alicloud_hbr_replication_vault_regions" "default" {}

locals {
  region_replication = data.alicloud_hbr_replication_vault_regions.default.regions.0.replication_region_id
}

provider "alicloud" {
  alias  = "replication"
  region = local.region_replication
}

resource "alicloud_hbr_replication_vault" "default" {
  replication_source_region_id = local.region_replication
  replication_source_vault_id  = alicloud_hbr_vault.default.id
  vault_name                   = var.name
  vault_storage_class          = "STANDARD"
  description                  = var.name
  provider                     = alicloud.replication
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) The description of the backup vault. The description must be 0 to 255 characters in length.
* `replication_source_region_id` - (Required, ForceNew) The ID of the region where the source vault resides.
* `replication_source_vault_id` - (Required, ForceNew) The ID of the source vault.
* `vault_name` - (Required) The name of the backup vault. The name must be 1 to 64 characters in length.
* `vault_storage_class` - (Optional, Computed, ForceNew) The storage type of the backup vault. Valid values: `STANDARD`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Replication Vault.
* `status` - The status of the resource.

## Import

Hybrid Backup Recovery (HBR) Replication Vault can be imported using the id, e.g.

```
$ terraform import alicloud_hbr_replication_vault.example <id>
```