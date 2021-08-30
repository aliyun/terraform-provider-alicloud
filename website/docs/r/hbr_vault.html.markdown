---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_vault"
sidebar_current: "docs-alicloud-resource-hbr-vault"
description: |-
  Provides a Alicloud Hybrid Backup Recovery (HBR) Backup vault resource.
---

# alicloud\_hbr\_vault

Provides a HBR Backup vault resource.

For information about HBR Backup vault and how to use it, see [What is Backup vault](https://www.alibabacloud.com/help/doc-detail/62362.htm).

-> **NOTE:** Available in v1.129.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_hbr_vault" "example" {
  vault_name = "example_value"
}
```

## Argument Reference

The following arguments are supported:

* `vault_name` - (Required) The name of Vault.
* `description` - (Optional) The description of Vault. Defaults to an empty string.
* `vault_type` - (Optional, Computed, ForceNew) The type of Vault. Valid values: `STANDARD`. 
* `vault_storage_class` - (Optional, Computed, ForceNew) The storage class of Vault. Valid values: `STANDARD`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Vault.
* `status` - The status of Vault. Valid values: `INITIALIZING`, `CREATED`, `ERROR`, `UNKNOWN`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Vault.

## Import

HBR Vault can be imported using the id, e.g.

```
$ terraform import alicloud_hbr_vault.example <id>
```