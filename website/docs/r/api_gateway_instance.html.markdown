---
subcategory: "Api Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_api_gateway_instance"
description: |-
  Provides a Alicloud Api Gateway Instance resource.
---

# alicloud_api_gateway_instance

Provides a Api Gateway Instance resource. 

For information about Api Gateway Instance and how to use it, see [What is Instance](https://www.alibabacloud.com/help/en/api-gateway/product-overview/dedicated-instances).

-> **NOTE:** Available since v1.218.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_api_gateway_instance" "default" {
  instance_name = var.name

  instance_spec = "api.s1.small"
  https_policy  = "HTTPS2_TLS1_0"
  zone_id       = "cn-hangzhou-MAZ6"
  payment_type  = "PayAsYouGo"
  user_vpc_id   = "1709116870"
  instance_type = "normal"
}
```

## Argument Reference

The following arguments are supported:
* `duration` - (Optional) The time of the instance package. Valid values:
  - PricingCycle is **Month**, indicating monthly payment. The value range is **1** to **9**.
  - PricingCycle is **Year**, indicating annual payment. The value range is **1** to **3**.

When the value of> ChargeType is **PrePaid**, this parameter is available and must be passed in.
* `egress_ipv6_enable` - (Optional) Does IPV6 Capability Support.
* `https_policy` - (Required) Https policy.
* `instance_name` - (Required) Instance name.
* `instance_spec` - (Required, ForceNew) Instance type.
* `instance_type` - (Optional, ForceNew, Computed) Instance type-normal: traditional exclusive instance.
* `payment_type` - (Required, ForceNew) The payment type of the resource.
* `pricing_cycle` - (Optional) The subscription instance is of the subscription year or month type. The value range is as follows:
  - **year**: year
  - **month**: month
-> **NOTE:**  If the Payment type is PrePaid, this parameter is required.
* `support_ipv6` - (Optional) Does ipv6 support.
* `user_vpc_id` - (Optional) User's VpcID.
* `vpc_slb_intranet_enable` - (Optional) Whether the slb of the Vpc supports.
* `zone_id` - (Optional, ForceNew) The zone where the instance is deployed.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - Creation time.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Instance.
* `delete` - (Defaults to 5 mins) Used when delete the Instance.
* `update` - (Defaults to 5 mins) Used when update the Instance.

## Import

Api Gateway Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_api_gateway_instance.example <id>
```