---
subcategory: "Cloud Storage Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_storage_gateway_storage_bundle"
sidebar_current: "docs-alicloud-resource-cloud-storage-gateway-storage-bundle"
description: |-
  Provides a Alicloud Cloud Storage Gateway Storage Bundle resource.
---

# alicloud\_cloud\_storage\_gateway\_storage\_bundle

Provides a Cloud Storage Gateway Storage Bundle resource.

For information about Cloud Storage Gateway Storage Bundle and how to use it, see [What is Storage Bundle](https://www.alibabacloud.com/help/en/doc-detail/53972.htm).

-> **NOTE:** Available in v1.116.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_cloud_storage_gateway_storage_bundle" "example" {
  storage_bundle_name = "example_value"
}

```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) The description of storage bundle.
* `storage_bundle_name` - (Required) The name of storage bundle.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Storage Bundle.

## Import

Cloud Storage Gateway Storage Bundle can be imported using the id, e.g.

```
$ terraform import alicloud_cloud_storage_gateway_storage_bundle.example <id>
```
