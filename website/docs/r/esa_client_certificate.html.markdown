---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_client_certificate"
description: |-
  Provides a Alicloud ESA Client Certificate resource.
---

# alicloud_esa_client_certificate

Provides a ESA Client Certificate resource.



For information about ESA Client Certificate and how to use it, see [What is Client Certificate](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateClientCertificate).

-> **NOTE:** Available since v1.244.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
  site_name           = "gositecdn.cn"
}

resource "alicloud_esa_client_certificate" "default" {
  site_id       = data.alicloud_esa_sites.default.sites.0.id
  pkey_type     = "RSA"
  validity_days = "365"
}
```

## Argument Reference

The following arguments are supported:
* `csr` - (Optional) Certificate signing request content.
* `pkey_type` - (Optional) The private key algorithm type.
* `site_id` - (Required, ForceNew, Int) Site Id
* `validity_days` - (Required) Certificate validity period.
* `status` - (Optional) The certificate status. Valid values: `revoked`, `active`.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<site_id>:<client_cert_id>`.
* `client_cert_id` - ClientCertificate Id
* `create_time` - The time when the certificate was created.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Client Certificate.
* `delete` - (Defaults to 5 mins) Used when delete the Client Certificate.
* `update` - (Defaults to 5 mins) Used when update the Client Certificate.

## Import

ESA Client Certificate can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_client_certificate.example <site_id>:<client_cert_id>
```