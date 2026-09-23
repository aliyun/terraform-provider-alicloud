---
subcategory: "Certificate Management Service (Original SSL Certificate)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ssl_certificates_service_pca_cert_sync"
description: |-
  Provides a Alicloud SSL Certificates Pca Cert Sync resource.
---

# alicloud_ssl_certificates_service_pca_cert_sync

Provides a SSL Certificates Pca Cert Sync resource. The resource synchronizes existing PCA (private CA) client certificates to the SSL Certificates service in one batch, the same operation as "Batch Synchronize to SSL Certificates" in the PCA certificate console.

For information about SSL Certificates Pca Cert Sync and how to use it, see [What is Pca Cert Sync](https://next.api.alibabacloud.com/document/cas/2020-06-30/UploadPcaCertToCas).

-> **NOTE:** Available since v1.294.0.

-> **NOTE:** The synchronization is a one-shot operation. The API returns no server-side object, and the synchronized certificates cannot be located by PCA identifier on the SSL Certificates side, so the resource does not refresh anything after creation. Changing the `ids` list recreates the resource, which re-runs the synchronization. Destroying the resource only removes it from the Terraform state; the synchronized certificates remain in the SSL Certificates service. The server processes the synchronization asynchronously: shortly after the apply, attempts to delete the synchronized certificates from the PCA side may be rejected until the processing finishes.

## Example Usage

Basic Usage

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_ssl_certificates_service_pca_certificate" "root" {
  organization      = "a"
  years             = "1"
  locality          = "a"
  organization_unit = "a"
  state             = "a"
  common_name       = "cbc.certqa.cn"
}

resource "alicloud_ssl_certificates_service_pca_certificate" "sub" {
  parent_identifier = alicloud_ssl_certificates_service_pca_certificate.root.id
  organization      = "a"
  years             = "1"
  locality          = "a"
  organization_unit = "a"
  state             = "a"
  common_name       = "cbc.certqa.cn"
  algorithm         = "RSA_2048"
  certificate_type  = "SUB_ROOT"
  enable_crl        = true
}

resource "alicloud_ssl_certificates_service_pca_cert" "default" {
  count             = 2
  days              = "1"
  parent_identifier = alicloud_ssl_certificates_service_pca_certificate.sub.id
  algorithm         = "RSA_2048"
  common_name       = "terraform-${count.index}"
  organization      = "terraform"
  state             = "Beijing"
  country_code      = "cn"
}

resource "alicloud_ssl_certificates_service_pca_cert_sync" "default" {
  ids = alicloud_ssl_certificates_service_pca_cert.default[*].id
}
```

## Argument Reference

The following arguments are supported:
* `ids` - (Required, ForceNew, List) The identifiers of the PCA client certificates to synchronize to the SSL Certificates service. Multiple identifiers are joined into a single comma-separated request.

## Attributes Reference

The following attributes are exported:
* `id` - The resource ID. The value is the comma-separated list of the synchronized PCA certificate identifiers.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the PCA Cert Sync.
* `delete` - (Defaults to 5 mins) Used when delete the PCA Cert Sync.

## Import

SSL Certificates Pca Cert Sync can be imported using the id, which is the comma-separated list of the synchronized PCA certificate identifiers, e.g.

```shell
$ terraform import alicloud_ssl_certificates_service_pca_cert_sync.example 59425,59426
```
