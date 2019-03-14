---
layout: "alicloud"
page_title: "Alicloud: alicloud_cas_certificate"
sidebar_current: "docs-alicloud-resource-cas-certificate"
description: |-
  Provides a CAS Certificate resource.
---

# alicloud\_cas\_certificate

Provides a CAS Certificate resource.

~> **NOTE:** The domain name which you want to add must be already registered and had not added by another account. Every domain name can only exist in a unique group.

## Example Usage

```
# Add a new Domain.
resource "alicloud_cas_certificate" "cert" {
   name = "test"
   cert = "./test.crt"
   key = "./test.key"
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the domain. This name without suffix can have a string of 1 to 63 characters, must contain only alphanumeric characters or "-", and must not begin or end with "-", and "-" must not in the 3th and 4th character positions at the same time. Suffix `.sh` and `.tel` are not supported.
* `cert` - (Required) Cert of the Certificate in which the Certificate will add.
* `key` - (Required) Key of the Certificate in which the Certificate will add. If not supplied.


## Attributes Reference

The following attributes are exported:

* `id` - The cert id.
