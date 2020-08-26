---
subcategory: "ALIDNS"
layout: "alicloud"
page_title: "Alicloud: alicloud_alidns_instance"
sidebar_current: "docs-alicloud-resource-alidns-instance"
description: |-
  Provides a Alicloud ALIDNS Instance resource.
---

# alicloud\_alidns\_instance

Create an ALIDNS Instance resource.

-> **NOTE:** Available in v1.80.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_alidns_instance" "this" {
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

* `dns_security` - (Required, ForceNew) ALIDNS security level. Valid values: `no`, `basic`, `advanced`.
* `domain_numbers` - (Required, ForceNew) Number of domain names bound.
* `period` - (Optional, ForceNew) Creating a pre-paid instance, it must be set, the unit is month, please enter an integer multiple of 12 for annually paid products.
* `renew_period` - (Optional, ForceNew) Automatic renewal period, the unit is month. When setting RenewalStatus to AutoRenewal, it must be set.
* `renewal_status` - (Optional, ForceNew) Automatic renewal status. Valid values: `AutoRenewal`, `ManualRenewal`, default to `ManualRenewal`.
* `version_code` - (Required, ForceNew) Paid package version. Valid values: `version_personal`, `version_enterprise_basic`, `version_enterprise_advanced`.

## Attributes Reference

* `id` - ID of the ALIDNS instance.
* `version_name` - Paid package version name.

## Import

ALIDNS instance be imported using the id, e.g.

```
$ terraform import alicloud_alidns_instance.example dns-cn-v0h1ldjhfff
```