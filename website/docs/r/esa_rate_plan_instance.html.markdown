---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_rate_plan_instance"
description: |-
  Provides a Alicloud ESA Rate Plan Instance resource.
---

# alicloud_esa_rate_plan_instance

Provides a ESA Rate Plan Instance resource.



For information about ESA Rate Plan Instance and how to use it, see [What is Rate Plan Instance](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.234.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_esa_rate_plan_instance&exampleId=8a610a35-0473-4250-1ee5-cec23e2ec9dad16ea3e5&activeTab=example&spm=docs.r.esa_rate_plan_instance.0.8a610a3504&intl_lang=EN_US" target="_blank">
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

### Deleting `alicloud_esa_rate_plan_instance` or removing it from your configuration

Terraform cannot destroy resource `alicloud_esa_rate_plan_instance`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `auto_pay` - (Optional) Whether to pay automatically.
* `auto_renew` - (Optional) Auto Renew:

  true: Automatic renewal.

  false: Do not renew automatically.
* `coverage` - (Optional) Acceleration area:

  domestic: Mainland China only.

  global: global.

  overseas: Global (excluding Mainland China).
* `payment_type` - (Optional, ForceNew, Computed) The payment type of the resource, Valid vales: Subscription.
* `period` - (Optional, Int) Purchase cycle (in months).
* `plan_name` - (Optional) The plan name, which is obtained from the DescribeRatePlanPrice interface.
* `type` - (Optional) Site access type:

  NS:NS access.

  CNAME:CNAME access.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The new purchase time of the package instance.
* `instance_status` - Renewing: renewing

  upgrading: upgrading

  releasePrepaidService: Prepaid overdue release

  creating: creating

  downgrading: downgrading

  ceasePrepaidService: prepaid service

  running: running
* `status` - The status of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Rate Plan Instance.
* `update` - (Defaults to 5 mins) Used when update the Rate Plan Instance.

## Import

ESA Rate Plan Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_rate_plan_instance.example <id>
```