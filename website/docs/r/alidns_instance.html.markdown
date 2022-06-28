---
subcategory: "DNS"
layout: "alicloud"
page_title: "Alicloud: alicloud_alidns_instance"
sidebar_current: "docs-alicloud-resource-alidns-instance"
description: |-
  Provides a Alicloud Alidns Instance resource.
---

# alicloud\_alidns\_instance

Create an Alidns Instance resource.

-> **NOTE:** Available in v1.95.0+.

-> **NOTE:** The Alidns Instance is not support to be purchase automatically in the international site.

## Example Usage

Basic Usage

```terraform
resource "alicloud_alidns_instance" "example" {
  dns_security   = "no"
  domain_numbers = "2"
  period         = 1
  renew_period   = 1
  renewal_status = "ManualRenewal"
  version_code   = "version_personal"
}

```

## Argument Reference

The following arguments are supported:

* `dns_security` - (Required, ForceNew) Alidns security level. Valid values: `no`, `basic`, `advanced`.
* `domain_numbers` - (Required, ForceNew) Number of domain names bound.
* `period` - (Optional) Creating a pre-paid instance, it must be set, the unit is month, please enter an integer multiple of 12 for annually paid products.
* `renew_period` - (Optional, ForceNew) Automatic renewal period, the unit is month. When setting RenewalStatus to AutoRenewal, it must be set.
* `renewal_status` - (Optional, ForceNew) Automatic renewal status. Valid values: `AutoRenewal`, `ManualRenewal`, default to `ManualRenewal`.
* `version_code` - (Required, ForceNew) Paid package version. Valid values: `version_personal`, `version_enterprise_basic`, `version_enterprise_advanced`.
* `payment_type` - (Optional, ForceNew) The billing method of the Alidns instance. Valid values: `Subscription`. Default to `Subscription`.

## Attributes Reference

* `id` - ID of the Alidns instance.
* `version_name` - Paid package version name.

## Import

DNS instance be imported using the id, e.g.

```
$ terraform import alicloud_alidns_instance.example dns-cn-v0h1ldjhfff
```
