---
subcategory: "SSL Certificates"
layout: "alicloud"
page_title: "Alicloud: alicloud_cas_certificate"
sidebar_current: "docs-alicloud-resource-cas-certificate"
description: |-
  Provides a CAS Certificate resource.
---

# alicloud\_cas\_certificate

-> **DEPRECATED:**  This datasource has been deprecated from version `1.129.0`. Please use new datasource [alicloud_ssl_certificates_service_certificate](https://www.terraform.io/docs/providers/alicloud/r/ssl_certificates_service_certificate).

Provides a CAS Certificate resource.

-> **NOTE:** The Certificate name which you want to add must be already registered and had not added by another account. Every Certificate name can only exist in a unique group.

-> **NOTE:** The Cas Certificate region only support cn-hangzhou, ap-south-1, me-east-1, eu-central-1, ap-northeast-1, ap-southeast-2.

-> **NOTE:** Available in 1.35.0+ .

## Example Usage

```
# Add a new Certificate.
resource "alicloud_cas_certificate" "cert" {
  name = "test"
  cert = file("${path.module}/test.crt")
  key  = file("${path.module}/test.key")
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required, ForcesNew) Name of the Certificate. This name without suffix can have a string of 1 to 63 characters, must contain only alphanumeric characters or "-", and must not begin or end with "-", and "-" must not in the 3th and 4th character positions at the same time. Suffix `.sh` and `.tel` are not supported.
* `cert` - (Required, ForcesNew) Cert of the Certificate in which the Certificate will add.
* `key` - (Required, ForcesNew) Key of the Certificate in which the Certificate will add.


## Attributes Reference

The following attributes are exported:

* `id` - The cert id.
