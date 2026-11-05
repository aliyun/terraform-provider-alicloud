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

-> **NOTE:** Available since v2.3.0.

-> **NOTE:** A catalog can only be dropped after all tables and user-created databases under it have been deleted, and those tables/databases must have been deleted at least 24 hours before the Delete call. Dropping a catalog that does not meet these conditions fails.

## Example Usage

Basic Usage

```terraform
terraform {
  required_providers {
    alicloud = {
      source = "aliyun/alicloud"
    }
  }
}

resource "alicloud_dlf_next_catalog" "default" {
  name = "terraform-example"
  type = "PAIMON"
  options = {
    comment = "a sample paimon catalog"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, ForceNew) The name of the catalog. The name must be unique within the region.
* `type` - (Optional, String, ForceNew) The type of the catalog. Default value: `PAIMON`. Valid values: `PAIMON`, `ICEBERG`. A `PAIMON` catalog is an OmniCatalog compatible with Paimon REST, Iceberg REST and HMS protocols and manages Paimon, Iceberg and Hive-compatible tables. `ICEBERG` is accepted because the CreateCatalog API documents it as an enum value; whether a specific account or region can create an Iceberg catalog depends on the DLF backend.
* `options` - (Optional, Map) The catalog option keys managed by Terraform, as a map of key-value pairs. `options = {k = v}` manages those keys: an update writes the new values and removes previously-managed keys that are no longer present. `options = {}` and omitting `options` both mean the managed set is expected to be empty, so previously-managed keys are removed from the remote catalog. Server-injected defaults and keys you never declared are never removed by this field; they remain visible in `remote_options`. A fresh import cannot recover which options were originally user-configured, so `options` is empty after import — re-apply the desired block to adopt managed keys.
* `is_shared` - (Optional, Bool, ForceNew) Whether the catalog is shared. The API has no documented default for `isShared`: omitting the argument does not send the field (the backend default applies), while `is_shared = false` sends an explicit `false`.
* `share_id` - (Optional, String, ForceNew) The share ID of the catalog. A share ID is issued by DLF catalog sharing and cannot be created by this resource; set it only when sharing has already been configured for the catalog.

## Attributes Reference

The following attributes are exported:

* `id` - (String) The Terraform resource id, which is the catalog name.
* `remote_options` - (Map) The full remote options returned by GetCatalog, including server-injected defaults and keys Terraform does not manage.
* `status` - (String) The status of the catalog. Valid values: `NEW`, `INITIALIZING`, `INITIALIZE_FAILED`, `RUNNING`, `TERMINATED`, `DELETING`, `DELETE_FAILED`, `DELETED`, `STORAGE_RESTRICTED`.
* `owner` - (String) The owner of the catalog.
* `created_at` - (String) The creation time of the catalog, in milliseconds since the epoch, as a string (the raw 64-bit value would overflow a 32-bit int).
* `created_by` - (String) The creator of the catalog.
* `updated_at` - (String) The update time of the catalog, in milliseconds since the epoch, as a string.
* `updated_by` - (String) The updater of the catalog.

## Import

DLF Next Catalog can be imported using the catalog name, e.g.

```shell
$ terraform import alicloud_dlf_next_catalog.example <catalog_name>
```

## RAM Permissions

The following RAM actions (operate level) are required to manage this resource:

* `dlf:CreateCatalog` - create a catalog.
* `dlf:GetCatalog` - query a catalog.
* `dlf:AlterCatalog` - modify catalog options.
* `dlf:DropCatalog` - delete a catalog.
