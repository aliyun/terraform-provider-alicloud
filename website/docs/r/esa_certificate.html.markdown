---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_certificate"
description: |-
  Provides a Alicloud ESA Certificate resource.
---

# alicloud_esa_certificate

Provides a ESA Certificate resource.



For information about ESA Certificate and how to use it, see [What is Certificate](https://next.api.alibabacloud.com/document/ESA/2024-09-10/SetCertificate).

-> **NOTE:** Available since v1.244.0.

## Example Usage

Basic Usage

```terraform

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
  site_name           = "gositecdn.cn"
}

resource "alicloud_esa_certificate" "default" {
  created_type = "free"
  domains      = "101.gositecdn.cn"
  site_id      = data.alicloud_esa_sites.default.sites.0.id
  type         = "lets_encrypt"
}
```

## Argument Reference

The following arguments are supported:
* `cas_id` - (Optional) The certificate ID on Certificate Management Service.
* `cert_id` - (Optional, ForceNew, Computed) The certificate ID on ESA.
* `cert_name` - (Optional) The certificate name.
* `certificate` - (Optional) The certificate content.
* `created_type` - (Required) The certificate ID on Certificate Management Service. Valid values:
  - free: a free certificate.
  - cas: a certificate purchased by using Certificate Management Service.
  - upload: a custom certificate that you upload.
* `domains` - (Optional, ForceNew) The Subject Alternative Name (SAN) of the certificate.
* `private_key` - (Optional) The certificate content.
* `region` - (Optional) The private key of the certificate.
* `site_id` - (Required, ForceNew, Int) Site ID, which can be obtained by calling the [ListSites](https://www.alibabacloud.com/help/en/doc-detail/2850189.html) interface.
* `type` - (Optional, ForceNew) The certificate type. Valid values:
  - cas: a certificate purchased by using Certificate Management Service.
  - upload: a custom certificate that you upload.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<site_id>:<cert_id>`.
* `create_time` - Creation time.
* `status` - Certificate status.(within 30 days).- issued.- applying.- application failed.- canceled.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Certificate.
* `delete` - (Defaults to 5 mins) Used when delete the Certificate.
* `update` - (Defaults to 5 mins) Used when update the Certificate.

## Import

ESA Certificate can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_certificate.example <site_id>:<cert_id>
```