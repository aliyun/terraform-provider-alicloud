---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_rate_plan_instance"
description: |-
  Provides a Alicloud ESA Rate Plan Instance resource.
---

# alicloud_esa_rate_plan_instance

Provides a ESA Rate Plan Instance resource.



For information about ESA Rate Plan Instance and how to use it, see [What is Rate Plan Instance](https://www.alibabacloud.com/help/en/edge-security-acceleration/esa/product-overview/query-package-information).

-> **NOTE:** Available since v1.234.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_esa_rate_plan_instance&exampleId=8a610a35-0473-4250-1ee5-cec23e2ec9dad16ea3e5&activeTab=example&spm=docs.r.esa_rate_plan_instance.0.8a610a3504&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_esa_rate_plan_instance" "default" {
  type         = "NS"
  auto_renew   = true
  period       = "1"
  payment_type = "Subscription"
  coverage     = "overseas"
  plan_name    = "basic"
  auto_pay     = true
}
```

## Argument Reference

The following arguments are supported:
* `payment_type` - (Optional, ForceNew) The payment type of the resource. Valid values:
  - `Subscription`: subscription.
* `auto_pay` - (Optional) Specifies whether to enable auto payment.
* `auto_renew` - (Optional) Auto-renewal:
  - `true`: Enable auto-renewal.
  - `false`: Disable auto-renewal.
* `coverage` - (Optional) The service locations for the websites that can be associated with the plan. Multiple values are separated by commas (,). Valid values:
  - `domestic`: the Chinese mainland.
  - `overseas`: outside the Chinese mainland.
  - `global`: global.
* `period` - (Optional, Int) Subscription period (in months).
* `plan_name` - (Optional) Package name.Value range:
  - `entranceplan_intl`：Entrance Plan
  -  `basicplan_intl`：Pro Plan
  - `vipplan_intl`：Premium Plan
* `type` - (Optional) The DNS setup option for the website. Valid values:
  - `NS`
  - `CNAME`

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The time when the plan was purchased.
* `status` - The status of the resource.
* `instance_status` - The plan instance status.
  - `Renewing`: renewing
  - `upgrading`: upgrading
  - `releasePrepaidService`: Prepaid overdue release
  - `creating`: creating
  - `downgrading`: downgrading
  - `ceasePrepaidService`: prepaid service
  - `running`: running

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Rate Plan Instance.
* `delete` - (Defaults to 5 mins) Used when delete the Rate Plan Instance.
* `update` - (Defaults to 5 mins) Used when update the Rate Plan Instance.

## Import

ESA Rate Plan Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_rate_plan_instance.example <id>
```