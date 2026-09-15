---
subcategory: "Certificate Management Service (Original SSL Certificate)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ssl_certificates_service_pca_cert_sync"
description: |-
  Provides a Alicloud SSL Certificates Pca Cert Sync action.
---

# alicloud_ssl_certificates_service_pca_cert_sync

Provides a SSL Certificates Pca Cert Sync action. The action synchronizes existing PCA (private CA) client certificates to the SSL Certificates service in one batch, the same operation as "Batch Synchronize to SSL Certificates" in the PCA certificate console.

For information about SSL Certificates Pca Cert Sync and how to use it, see [What is Pca Cert Sync](https://next.api.alibabacloud.com/document/cas/2020-06-30/UploadPcaCertToCas).

-> **NOTE:** Available since v2.0.0-beta5.

-> **NOTE:** Actions are imperative operations. Requires Terraform 1.14 or later. The synchronization is one-shot and asynchronous: the server processes it shortly after the action completes, and repeated synchronization of already-synchronized certificates is safe.

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

  lifecycle {
    action_trigger {
      events  = [after_create]
      actions = [action.alicloud_ssl_certificates_service_pca_cert_sync.default]
    }
  }
}

action "alicloud_ssl_certificates_service_pca_cert_sync" "default" {
  config {
    ids = alicloud_ssl_certificates_service_pca_cert.default[*].id
  }
}
```

## Argument Reference

The following arguments are supported:
* `ids` - (Required, List) The identifiers of the PCA client certificates to synchronize to the SSL Certificates service.
