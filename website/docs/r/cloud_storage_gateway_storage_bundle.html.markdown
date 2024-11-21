---
subcategory: "Cloud Storage Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_storage_gateway_storage_bundle"
sidebar_current: "docs-alicloud-resource-cloud-storage-gateway-storage-bundle"
description: |-
  Provides a Alicloud Cloud Storage Gateway Storage Bundle resource.
---

# alicloud_cloud_storage_gateway_storage_bundle

Provides a Cloud Storage Gateway Storage Bundle resource.

For information about Cloud Storage Gateway Storage Bundle and how to use it, see [What is Storage Bundle](https://www.alibabacloud.com/help/en/cloud-storage-gateway/latest/createstoragebundle).

-> **NOTE:** Available since v1.116.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cloud_storage_gateway_storage_bundle&exampleId=94a65ca5-3d39-301b-ed19-6f16d51339ee9cf0492c&activeTab=example&spm=docs.r.cloud_storage_gateway_storage_bundle.0.94a65ca53d&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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

```shell
$ terraform import alicloud_cloud_storage_gateway_storage_bundle.example <id>
```
