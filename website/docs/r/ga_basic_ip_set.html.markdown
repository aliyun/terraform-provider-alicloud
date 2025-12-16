---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_basic_ip_set"
sidebar_current: "docs-alicloud-resource-ga-basic-ip-set"
description: |-
  Provides a Alicloud Global Accelerator (GA) Basic Ip Set resource.
---

# alicloud_ga_basic_ip_set

Provides a Global Accelerator (GA) Basic Ip Set resource.

For information about Global Accelerator (GA) Basic Ip Set and how to use it, see [What is Basic Ip Set](https://www.alibabacloud.com/help/en/global-accelerator/latest/api-ga-2019-11-20-createbasicipset).

-> **NOTE:** Available since v1.194.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ga_basic_ip_set&exampleId=e914ddd6-2c21-e1e2-7757-c575c716b0323ea94d1a&activeTab=example&spm=docs.r.ga_basic_ip_set.0.e914ddd62c&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = var.region
}

variable "region" {
  default = "cn-hangzhou"
}

resource "alicloud_ga_basic_accelerator" "default" {
  duration               = 1
  pricing_cycle          = "Month"
  bandwidth_billing_type = "CDT"
  auto_pay               = true
  auto_use_coupon        = "true"
  auto_renew             = false
  auto_renew_duration    = 1
}

resource "alicloud_ga_basic_ip_set" "default" {
  accelerator_id       = alicloud_ga_basic_accelerator.default.id
  accelerate_region_id = var.region
  isp_type             = "BGP"
  bandwidth            = "5"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ga_basic_ip_set&spm=docs.r.ga_basic_ip_set.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `accelerator_id` - (Required, ForceNew) The ID of the basic GA instance.
* `accelerate_region_id` - (Required, ForceNew) The ID of the acceleration region.
* `isp_type` - (Optional, ForceNew) The line type of the elastic IP address (EIP) in the acceleration region. Default value: `BGP`. Valid values: `BGP`, `BGP_PRO`, `ChinaTelecom`, `ChinaUnicom`, `ChinaMobile`, `ChinaTelecom_L2`, `ChinaUnicom_L2`, `ChinaMobile_L2`.
* `bandwidth` - (Optional) The bandwidth of the acceleration region. Unit: Mbit/s.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Basic Ip Set.
* `status` - The status of the Basic Ip Set instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the Basic Ip Set.
* `update` - (Defaults to 3 mins) Used when update the Basic Ip Set.
* `delete` - (Defaults to 3 mins) Used when delete the Basic Ip Set.

## Import

Global Accelerator (GA) Basic Ip Set can be imported using the id, e.g.

```shell
$ terraform import alicloud_ga_basic_ip_set.example <id>
```
