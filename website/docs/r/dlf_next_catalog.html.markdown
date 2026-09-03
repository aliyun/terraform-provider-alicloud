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

-> **NOTE:** Available since v1.232.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_dlf_next_catalog" "default" {
  name = "tf-demo-catalog"
  type = "PAIMON"
  options = {
    key1 = "value1"
    key2 = "value2"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) The name of the catalog. The name must be unique within the region.
* `type` - (Required, ForceNew) The type of the catalog. Valid values: `PAIMON`, `ICEBERG`.
* `options` - (Optional) The configuration options of the catalog. It is a map of key-value pairs.
* `is_shared` - (Optional, ForceNew) Whether the catalog is shared.
* `share_id` - (Optional, ForceNew) The share ID of the catalog.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the catalog.
* `region_id` - The region ID of the catalog.
* `status` - The status of the catalog. Valid values: `Creating`, `Active`, `Deleting`, `Failed`.
* `owner` - The owner of the catalog.
* `created_at` - The creation time of the catalog.
* `created_by` - The creator of the catalog.
* `updated_at` - The update time of the catalog.
* `updated_by` - The updater of the catalog.

## Import

DLF Next Catalog can be imported using the name, e.g.

```shell
$ terraform import alicloud_dlf_next_catalog.example <catalog_name>
```
