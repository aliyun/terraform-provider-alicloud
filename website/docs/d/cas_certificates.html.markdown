---
subcategory: "Certificate Management Service (Original SSL Certificate)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cas_certificates"
sidebar_current: "docs-alicloud-datasource-cas-certificates"
description: |-
  Provides a list of certs available to the user.
---

# alicloud_cas_certificates

-> **DEPRECATED:**  This datasource has been deprecated from version `1.129.0`. Please use new datasource [alicloud_ssl_certificates_service_certificates](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/data-sources/ssl_certificates_service_certificates).

This data source provides a list of CAS Certificates in an Alibaba Cloud account according to the specified filters.

## Example Usage

```terraform
data "alicloud_cas_certificates" "certs" {
  name_regex  = "^cas"
  output_file = "${path.module}/cas_certificates.json"
}

output "cert" {
  value = "${data.alicloud_cas_certificates.certs.certificates.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, List) A list of Ssl Certificates IDs.
* `name_regex` - (Optional) A regex string to filter results by Ssl Certificates name.
* `keyword` - (Optional, Available since v1.232.0) The domain name that is bound or the ID of the resource. Fuzzy match is supported.
* `enable_details` - (Optional, Bool) Whether to query the detailed list of resource attributes. Default value: `false`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `lang` - (Optional, ForceNew, Deprecated since v1.232.0) The language of the content within the request and response.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Ssl Certificates names.
* `certificates` - A list of Ssl Certificates Service Certificates.

---

### `certificates`

The certificates supports the following attributes:

* `id` - The ID of the Ssl Certificates Service Certificate.
* `cert_id` - The ID of the certificate.
* `certificate_name` - The name of the certificate.
* `fingerprint` - The fingerprint of the certificate.
* `common` - The parent domain name of the certificate.
* `sans` - All domain names that are bound to the certificate.
* `org_name` - The name of the organization that is associated with the certificate.
* `issuer` - The issuer of the certificate.
* `country` - The code of the country in which the organization is located.
* `province` - The province or autonomous region in which the organization is located.
* `city` - The city in which the organization is located.
* `expired` - Whether the certificate has expired.
* `start_date` - The time at which the certificate starts to take effect.
* `end_date` - The time at which the certificate expires.
* `cert` - The Cert of the certificate.

  -> **NOTE:** `cert` takes effect only if `enable_details` is set to `true`.
* `key` - The Key of the certificate.

  -> **NOTE:** `key` takes effect only if `enable_details` is set to `true`.
* `buy_in_aliyun` - Whether the certificate was purchased from Aliyun.

  -> **NOTE:** From version 1.232.0, `buy_in_aliyun` takes effect only if `enable_details` is set to `true`.
* `name` - (Deprecated since v1.129.0) The name of the certificate. Use `certificate_name` instead.
