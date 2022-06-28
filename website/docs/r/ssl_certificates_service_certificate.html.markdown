---
subcategory: "SSL Certificates"
layout: "alicloud"
page_title: "Alicloud: alicloud_ssl_certificates_service_certificate"
sidebar_current: "docs-alicloud-resource-ssl-certificates-service-certificate"
description: |-
  Provides a Alicloud SSL Certificates Certificate resource.
---

# alicloud\_ssl\_certificates\_service\_certificate

Provides a SSL Certificates Certificate resource.

For information about SSL Certificates Certificate and how to use it, see [What is Certificate](https://www.alibabacloud.com/help/product/28533.html).

-> **NOTE:** Available in v1.129.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_ssl_certificates_service_certificate" "example" {
  certificate_name = "test"
  cert             = file("${path.module}/test.crt")
  key              = file("${path.module}/test.key")
}

```

## Argument Reference

The following arguments are supported:

* `cert` - (Required, ForceNew) Cert of the Certificate in which the Certificate will add.
* `certificate_name` - (Optional, ForceNew) Name of the Certificate. 
  This name without suffix can have a string of 1 to 63 characters, must contain only alphanumeric characters or "-", 
  and must not begin or end with "-", and "-" must not in the 3th and 4th character positions at the same time. 
  Suffix .sh and .tel are not supported.
  **NOTE:** One of `certificate_name` and `name` must be specified.
* `name` - (Deprecated, Optional, ForceNew) It has been deprecated from version 1.129.0 and using `certificate_name` instead.
* `key` - (Required, ForceNew) Key of the Certificate in which the Certificate will add.
* `lang` - (Optional) The lang.

## Attributes Reference

The following attributes are exported:

* `id` - The cert id.

## Import

SSL Certificates Certificate can be imported using the id, e.g.

```
$ terraform import alicloud_ssl_certificates_service_certificate.example <id>
```
