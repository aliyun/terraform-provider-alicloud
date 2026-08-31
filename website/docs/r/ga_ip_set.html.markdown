---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_ip_set"
sidebar_current: "docs-alicloud-resource-ga-ip-set"
description: |-
  Provides a Alicloud Global Accelerator (GA) Ip Set resource.
---

# alicloud_ga_ip_set

Provides a Global Accelerator (GA) Ip Set resource.

For information about Global Accelerator (GA) Ip Set and how to use it, see [What is Ip Set](https://www.alibabacloud.com/help/en/global-accelerator/latest/api-ga-2019-11-20-createipsets).

-> **NOTE:** Available since v1.113.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ga_ip_set&exampleId=6cf17f10-dac7-8038-d7b3-074b7ff5ababfc85d97b&activeTab=example&spm=docs.r.ga_ip_set.0.6cf17f10da&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "region" {
  default = "cn-hangzhou"
}

provider "alicloud" {
  region  = var.region
  profile = "default"
}

resource "alicloud_ga_accelerator" "default" {
  duration        = 1
  auto_use_coupon = true
  spec            = "1"
}

resource "alicloud_ga_bandwidth_package" "default" {
  bandwidth      = 100
  type           = "Basic"
  bandwidth_type = "Basic"
  payment_type   = "PayAsYouGo"
  billing_type   = "PayBy95"
  ratio          = 30
}

resource "alicloud_ga_bandwidth_package_attachment" "default" {
  accelerator_id       = alicloud_ga_accelerator.default.id
  bandwidth_package_id = alicloud_ga_bandwidth_package.default.id
}

resource "alicloud_ga_ip_set" "example" {
  accelerate_region_id = var.region
  bandwidth            = "5"
  accelerator_id       = alicloud_ga_bandwidth_package_attachment.default.accelerator_id
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ga_ip_set&spm=docs.r.ga_ip_set.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `accelerator_id` - (Required, ForceNew) The ID of the Global Accelerator (GA) instance.
* `accelerate_region_id` - (Required, ForceNew) The ID of an acceleration region.
* `bandwidth` - (Optional, Int) The bandwidth allocated to the acceleration region.
-> **NOTE:** The minimum bandwidth of each accelerated region is 2Mbps. The total bandwidth of the acceleration region should be less than or equal to the bandwidth of the basic bandwidth package you purchased.
* `ip_version` - (Optional, ForceNew) The IP protocol used by the GA instance. Default value: `IPv4`. Valid values: `IPv4`, `IPv6`, `DUAL_STACK`. **NOTE:** From version 1.220.0, `ip_version` can be set to `DUAL_STACK`.
* `isp_type` - (Optional, ForceNew, Available since v1.207.0) The line type of the elastic IP address (EIP) in the acceleration region. Valid values: `BGP`, `BGP_PRO`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Ip Set.
* `ip_address_list` - The list of accelerated IP addresses in the acceleration region.
* `status` -  The status of the acceleration region.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 15 mins) Used when create the Ip Set.
* `update` - (Defaults to 3 mins) Used when update the Ip Set.
* `delete` - (Defaults to 10 mins) Used when delete the Ip Set.

## Import

Ga Ip Set can be imported using the id, e.g.

```shell
$ terraform import alicloud_ga_ip_set.example <id>
```
