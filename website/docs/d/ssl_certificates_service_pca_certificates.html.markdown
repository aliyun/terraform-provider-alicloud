---
subcategory: "Certificate Management Service (Original SSL Certificate)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ssl_certificates_service_pca_certificates"
sidebar_current: "docs-alicloud-datasource-ssl-certificates-service-pca-certificates"
description: |-
  Provides a list of Ssl Certificates Service Pca Certificates to the user.
---

# alicloud\_ssl\_certificates\_service\_pca\_certificates

This data source provides the Ssl Certificates Service Pca Certificates of current Alibaba Cloud user.

-> **NOTE:** Available since v1.231.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_ssl_certificates_service_pca_certificates" "default" {
  cert_type   = "root"
  output_file = "pca_certificates.txt"
}
output "pca_certificates_first_id" {
  value = data.alicloud_ssl_certificates_service_pca_certificates.default.certificates.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of Pca Certificate IDs.
* `name_regex` - (Optional) A regex string to apply to the common name filter of the Pca Certificate.
* `ca_status` - (Optional) The status of the CA. Valid values: `issue`, `forbidden`, `revoke`.
* `cert_type` - (Optional) The type of the CA. Valid values: `root`, `subRoot`, `externalCa`.
* `issuer_type` - (Optional) The issuer of the CA. Valid values: `local`, `iTrusChina`, `external`.
* `valid_status` - (Optional) The validity status of the CA. Valid values: `valid`, `notValid`.
* `resource_group_id` - (Optional) The ID of the resource group to which the CA certificate belongs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - The ids list of Ssl Certificates Service Pca Certificates.
* `names` - The common names list of Ssl Certificates Service Pca Certificates.
* `certificates` - A list of Pca Certificates. Each element contains the following attributes:
  * `id` - The ID of the Pca Certificate.
  * `identifier` - The unique identifier of the CA certificate.
  * `serial_number` - The serial number of the CA certificate.
  * `x509_certificate` - The X509 certificate content of the CA certificate.
  * `certificate_type` - The type of the CA certificate. Valid values: `ROOT`, `SUB_ROOT`.
  * `algorithm` - The signature algorithm of the CA certificate.
  * `sign_algorithm` - The sign algorithm of the CA certificate.
  * `sha2` - The SHA-2 fingerprint of the CA certificate.
  * `md5` - The MD5 fingerprint of the CA certificate.
  * `locality` - The locality of the CA certificate.
  * `organization` - The organization of the CA certificate.
  * `organization_unit` - The organization unit of the CA certificate.
  * `common_name` - The common name of the CA certificate.
  * `country_code` - The country code of the CA certificate.
  * `state` - The state of the CA certificate.
  * `parent_identifier` - The unique identifier of the parent CA certificate. This parameter is returned only when `certificate_type` is `SUB_ROOT`.
  * `status` - The status of the CA certificate. Valid values: `ISSUE`, `REVOKE`.
  * `years` - The validity period of the CA certificate. Unit: years.
  * `before_date` - The time when the CA certificate takes effect.
  * `after_date` - The time when the CA certificate expires.
  * `resource_group_id` - The ID of the resource group to which the CA certificate belongs.
