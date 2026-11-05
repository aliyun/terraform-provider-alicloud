---
subcategory: "Certificate Management Service (Original SSL Certificate)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ssl_certificates_service_warehouse"
description: |-
  Provides a Alicloud SSL Certificates Service Warehouse resource.
---

# alicloud_ssl_certificates_service_warehouse

Provides a SSL Certificates Service Warehouse resource. A warehouse is the certificate storage container created through the SSL Certificates Service (CAS).

For information about SSL Certificates Service Warehouse and how to use it, see [What is Warehouse](https://next.api.alibabacloud.com/document/cas/2020-04-07/CreateWarehouse).

-> **NOTE:** Available since v1.295.0.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ssl_certificates_service_warehouse&exampleId=warehouse-basic&activeTab=example&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_ssl_certificates_service_warehouse" "default" {
  name = "terraform-example-warehouse"
  type = "ssl"
  biz  = "test"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?resource=alicloud_ssl_certificates_service_warehouse&activeTab=sample&source=OfficialSample:alicloud_ssl_certificates_service_warehouse&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) The name of the warehouse. The name must be 1 to 64 characters in length.
* `type` - (Required, ForceNew) The type of the warehouse. Valid values: `pcaAuthorityCertificate`, `pcaCertificate`, `personal`, `ssl`, `uploadedAuthorityCertificate`, `uploadedCertificate`.
* `biz` - (Required, ForceNew) The business scenario of the warehouse.

-> **NOTE:** There is no UpdateWarehouse API, so every argument is replacement-only (`ForceNew`). Changing any argument recreates the warehouse instead of updating it in place.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of the Warehouse, which is the same as `warehouse_instance_id`.
* `warehouse_instance_id` - The ID of the warehouse instance allocated by CreateWarehouse.

## Import

SSL Certificates Service Warehouse can be imported using the `id`, e.g.

```shell
$ terraform import alicloud_ssl_certificates_service_warehouse.example <warehouse_instance_id>
```

-> **NOTE:** `type` and `biz` are write-only or use a different enum value domain on read-back (CreateWarehouse `Type` vs ListWarehouse `WarehouseType`), so they are retained from configuration during import and are not refreshed from the API.
