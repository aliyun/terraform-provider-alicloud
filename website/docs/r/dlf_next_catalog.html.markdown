---
subcategory: "dlf_next"
layout: "alicloud"
page_title: "Alicloud: alicloud_dlf_next_catalog"
sidebar_current: "docs-alicloud-resource-dlf-next-catalog"
description: |-
  Provides a DLF Next Catalog resource.
---

# alicloud_dlf_next_catalog

Provides a DLF Next Catalog resource.

For information about DLF Next Catalog and how to use it, see [What is DLF Next Catalog](https://www.alibabacloud.com/help/en/dlf/).

-> **NOTE:** Available since v1.293.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_dlf_next_catalog" "default" {
  name = "tf-demo-catalog"
  type = "PAIMON"
  options = {
    comment = "a sample paimon catalog"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) The name of the catalog. The name must be unique within the region.
* `type` - (Optional, ForceNew) The type of the catalog. Defaults to `PAIMON`. Valid values: `PAIMON`. DLF 3.0 uses OmniCatalog as the unified catalog model, and `PAIMON` is the catalog type enum DLFNext CreateCatalog accepts for an OmniCatalog. A `PAIMON` catalog does not manage only Paimon tables: an OmniCatalog is compatible with Paimon REST, Iceberg REST and HMS protocols and manages Paimon, Iceberg and Hive-compatible tables. `ICEBERG` is not a separate CreateCatalog type — Iceberg is a table format / access protocol surfaced through an OmniCatalog, so it is not accepted as a `type` value.
* `options` - (Optional) The catalog configuration options managed by Terraform, as a map of key-value pairs. This field holds ONLY the option keys you declare in the configuration. The contract is **two-state**: `options = {k: v}` means Terraform manages those keys (updates are derived from the new managed set, removals from the previously-managed keys that are now absent); `options = {}` OR omitting `options` both mean the managed set is expected to be empty and any previously-managed keys are removed via AlterCatalog `removals`. This does not rely on the legacy SDK to distinguish an omitted block from an explicit empty map. On refresh, each managed key is reconciled against the value returned by GetCatalog, and a managed key the API no longer returns is dropped from state so a server-side deletion or an ignored update surfaces as drift in the next plan. Server-injected defaults and keys you never declared are NOT promoted into `options`; they remain visible in `remote_options`. A fresh `terraform import` cannot recover which remote options were originally user-configured, so `options` is empty after import — re-apply with the desired `options` block to adopt managed keys. AlterCatalog removals are derived only from the previously-managed keys that are now absent from `options`, so a clear or a partial removal can never delete a server default or an unmanaged key.
* `is_shared` - (Optional, ForceNew) Whether the catalog is shared.
* `share_id` - (Optional, ForceNew) The share ID of the catalog.

## Attributes Reference

The following attributes are exported:

* `id` - The Terraform resource id, which is the catalog name. The catalog also has a separate server-generated id (of the form `clg-xxx`); that server id is not the Terraform resource id and is surfaced only through the `alicloud_dlf_next_catalogs` data source.
* `remote_options` - The full read-only mirror of the options returned by GetCatalog, including server-injected defaults and keys Terraform does not manage. This field is computed and never drives an update; it exists so a fresh import (where `options` is empty) does not lose the remote picture, and so the complete remote options — including defaults you did not configure — remain inspectable. Because `options` only ever contains managed keys, AlterCatalog removals are computed from the managed set and can never reach into `remote_options` to delete a server default or an unmanaged key.
* `region_id` - The region ID of the catalog.
* `status` - The status of the catalog. Valid values: `NEW`, `INITIALIZING`, `INITIALIZE_FAILED`, `RUNNING`, `TERMINATED`, `DELETING`, `DELETE_FAILED`, `DELETED`, `STORAGE_RESTRICTED`.
* `owner` - The owner of the catalog.
* `created_at` - The creation time of the catalog.
* `created_by` - The creator of the catalog.
* `updated_at` - The update time of the catalog.
* `updated_by` - The updater of the catalog.

## Import

DLF Next Catalog can be imported using the catalog name. The Terraform resource id is the catalog name, e.g.

```shell
$ terraform import alicloud_dlf_next_catalog.example <catalog_name>
```

The catalog also has a separate server-generated id (of the form `clg-xxx`); that server id is not the Terraform resource id and is only surfaced through the `alicloud_dlf_next_catalogs` data source.

## RAM Permissions

The following RAM actions (operate level) are required to manage this resource:

* `dlf:CreateCatalog` - create a catalog.
* `dlf:GetCatalog` - query a catalog.
* `dlf:AlterCatalog` - modify catalog options.
* `dlf:DropCatalog` - delete a catalog.
