---
subcategory: "Certificate Management Service (Original SSL Certificate)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ssl_certificates_service_instance_certificates"
sidebar_current: "docs-alicloud-datasource-ssl-certificates-service-instance-certificates"
description: |-
  Provides a list of Ssl Certificates Service Instance Certificate owned by an Alibaba Cloud account.
---

# alicloud_ssl_certificates_service_instance_certificates

This data source provides Ssl Certificates Service Instance Certificate available to the user.[What is Instance Certificate](https://next.api.alibabacloud.com/document/cas/2020-04-07/ListCertificates)

-> **NOTE:** Available since v1.287.0.

## Example Usage

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_ssl_certificates_service_instance_certificates" "default" {
  certificate_status = "issued"
}

output "alicloud_ssl_certificates_service_instance_certificate_example_id" {
  value = data.alicloud_ssl_certificates_service_instance_certificates.default.certificates.0.id
}
```

## Argument Reference

The following arguments are supported:
* `instance_id` - (ForceNew, Optional) The ID of the certificate instance that issued the certificate. Use it to list the certificates a single instance has issued.
* `certificate_status` - (ForceNew, Optional) The status of the certificate. Valid values: `issued`, `revoked`, `willExpire` and `expired`.
* `certificate_source` - (ForceNew, Optional) The source of the certificate. Valid values: `BUY`, `UPLOAD` and `TEST`.
* `ids` - (Optional, Computed) A list of Instance Certificate IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Instance Certificate IDs.
* `certificates` - A list of Instance Certificate Entries. Each element contains the following attributes:
  * `certificate_id` - The ID of the certificate.
  * `instance_id` - The ID of the certificate instance that issued this certificate.
  * `cert_identifier` - The global certificate identifier, formatted as the certificate ID plus `-` plus the site region ID. Alibaba Cloud services such as ALB, CDN and WAF reference a certificate by this value.
  * `certificate_name` - The name of the certificate.
  * `certificate_status` - The status of the certificate.
  * `certificate_source` - The source of the certificate.
  * `common_name` - The common name of the certificate.
  * `domain` - All domain names covered by the certificate, separated by commas.
  * `algorithm` - The encryption algorithm of the certificate.
  * `key_size` - The key length of the certificate algorithm.
  * `not_before` - The start of the certificate validity period.
  * `not_after` - The end of the certificate validity period.
  * `issuer` - The certificate authority that issued the certificate.
  * `finger_print` - The public key fingerprint of the certificate.
  * `exist_private_key` - Indicates whether the private key of the certificate is held on the server side.
  * `id` - The ID of the resource supplied above.
