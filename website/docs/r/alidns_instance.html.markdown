---
subcategory: "Alidns"
layout: "alicloud"
page_title: "Alicloud: alicloud_alidns_instance"
sidebar_current: "docs-alicloud-resource-alidns-instance"
description: |-
  Provides a Alicloud Alidns Instance resource.
---

# alicloud_alidns_instance

Create an Alidns Instance resource.

-> **NOTE:** Available since v1.95.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_alidns_instance&exampleId=25ff8178-a15e-63da-877e-1f5b5704d10ab81f993a&activeTab=example&spm=docs.r.alidns_instance.0.25ff8178a1&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_alidns_instance&spm=docs.r.alidns_instance.example&intl_lang=EN_US)

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

```shell
$ terraform import alicloud_alidns_instance.example dns-cn-v0h1ldjhfff
```
