---
subcategory: "SSL Certificates"
layout: "alicloud"
page_title: "Alicloud: alicloud_cas_certificates"
sidebar_current: "docs-alicloud-datasource-cas-certificates"
description: |-
  Provides a list of certs available to the user.
---

# alicloud\_cas\_certificates

-> **DEPRECATED:**  This datasource has been deprecated from version `1.129.0`. Please use new datasource [alicloud_ssl_certificates_service_certificates](https://www.terraform.io/docs/providers/alicloud/d/ssl_certificates_service_certificates).

This data source provides a list of CAS Certificates in an Alibaba Cloud account according to the specified filters.

## Example Usage

```
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

* `name_regex` - (Optional) A regex string to filter results by the certificate name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `ids` - (Optional, Available in 1.52.0+) A list of cert IDs.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of cert IDs.
* `names` - A list of cert names. 
* `certificates` - A list of apis. Each element contains the following attributes:
  * `id` - The cert's id.
  * `name` - The cert's name.
  * `common` - The cert's common name.
  * `finger_print` - The cert's finger.
  * `issuer` - The cert's .
  * `org_name` - The cert's organization.
  * `province` - The cert's province.
  * `city` - The cert's city.
  * `country` - The cert's country.
  * `start_date` - The cert's not valid before time.
  * `end_date` - The cert's not valid after time.
  * `sans` - The cert's subject alternative name.
  * `expired` - The cert is expired or not.
  * `buy_in_aliyun` - The cert is buy from aliyun or not.
