---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_client_ca_certificate"
description: |-
  Provides a Alicloud ESA Client Ca Certificate resource.
---

# alicloud_esa_client_ca_certificate

Provides a ESA Client Ca Certificate resource.



For information about ESA Client Ca Certificate and how to use it, see [What is Client Ca Certificate](https://next.api.alibabacloud.com/document/ESA/2024-09-10/UploadClientCaCertificate).

-> **NOTE:** Available since v1.244.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
  site_name           = "gositecdn.cn"
}

resource "alicloud_esa_client_ca_certificate" "default" {
  certificate         = "-----BEGIN CERTIFICATE-----\n****-----END CERTIFICATE-----"
  client_ca_cert_name = "example"
  site_id             = data.alicloud_esa_sites.default.sites.0.id
}
```

## Argument Reference

The following arguments are supported:
* `certificate` - (Required, ForceNew) Certificate content.
* `client_ca_cert_name` - (Optional, ForceNew) The certificate name.
* `site_id` - (Required, ForceNew) Site Id

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<site_id>:<client_ca_cert_id>`.
* `client_ca_cert_id` - ClientCaCertificate Id
* `create_time` - Creation time.
* `status` - Certificate status.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Client Ca Certificate.
* `delete` - (Defaults to 5 mins) Used when delete the Client Ca Certificate.

## Import

ESA Client Ca Certificate can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_client_ca_certificate.example <site_id>:<client_ca_cert_id>
```