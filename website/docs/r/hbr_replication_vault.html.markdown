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

resource "alicloud_hbr_vault" "default" {
  provider   = alicloud.source
  vault_name = "terraform-example"
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

```shell
$ terraform import alicloud_hbr_replication_vault.example <id>
```