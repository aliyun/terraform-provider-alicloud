---
subcategory: "Elastic IP Address (EIP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_eip"
sidebar_current: "docs-alicloud-resource-eip"
description: |-
  Provides a ECS EIP resource.
---

# alicloud\_eip

Provides an elastic IP resource.

-> **DEPRECATED:**  This resource  has been deprecated from version `1.126.0`. Please use new resource [alicloud_eip_address](https://www.terraform.io/docs/providers/alicloud/r/eip_address).

-> **NOTE:** The resource only supports to create `PostPaid PayByTraffic`  or `PrePaid PayByBandwidth` elastic IP for international account. Otherwise, you will happened error `COMMODITY.INVALID_COMPONENT`.
Your account is international if you can use it to login in [International Web Console](https://account.alibabacloud.com/login/login.htm).

-> **NOTE:** From version 1.10.1, this resource supports creating "PrePaid" EIP. In addition, it supports setting EIP name and description.

## Example Usage

```
# Create a new EIP.
resource "alicloud_eip" "example" {
  bandwidth            = "10"
  internet_charge_type = "PayByBandwidth"
}
```

## Module Support

You can use the existing [eip module](https://registry.terraform.io/modules/terraform-alicloud-modules/eip/alicloud) 
to create several EIP instances and associate them with other resources one-click, like ECS instances, SLB, Nat Gateway and so on.

## Argument Reference

The following arguments are supported:

* `name` - (Deprecated) It has been deprecated from version 1.126.0 and using new attribute `address_name` instead.
* `address_name` - (Optional) The name of the EIP instance. This name can have a string of 2 to 128 characters, must contain only alphanumeric characters or hyphens, such as "-",".","_", and must not begin or end with a hyphen, and must not begin with http:// or https://.
* `description` - (Optional) Description of the EIP instance, This description can have a string of 2 to 256 characters, It cannot begin with http:// or https://. Default value is null.
* `bandwidth` - (Optional) Maximum bandwidth to the elastic public network, measured in Mbps (Mega bit per second). If this value is not specified, then automatically sets it to 5 Mbps.
* `internet_charge_type` - (Optional, ForceNew) Internet charge type of the EIP, Valid values are `PayByBandwidth`, `PayByTraffic`. Default to `PayByBandwidth`. **NOTE:** From version `1.7.1` to `1.125.0`, it defaults to `PayByTraffic`. It is only "PayByBandwidth" when `instance_charge_type` is PrePaid.
* `instance_charge_type` - (Depreacted, Optional, ForceNew) (It has been deprecated from version 1.126.0 and using new attribute `payment_type` instead) Elastic IP instance charge type. Valid values are "PrePaid" and "PostPaid". Default to "PostPaid".
* `payment_type` - (Optional, ForceNew) The billing method of the EIP. Valid values: `Subscription` and `PayAsYouGo`. Default value is `PayAsYouGo`.
* `period` - (Optional) The duration that you will buy the resource, in month. It is valid when `instance_charge_type` is `PrePaid`. Valid values: [1-9, 12, 24, 36]. At present, the provider does not support modify "period" and you can do that via web console.
**NOTE:** The attribute `period` is only used to create Subscription instance or modify the PayAsYouGo instance to Subscription. Once effect, it will not be modified that means running `terraform apply` will not effect the resource.
* `isp` - (Optional, ForceNew, Available in 1.47.0+) The line type of the Elastic IP instance. Default to `BGP`. Other type of the isp need to open a whitelist.
* `tags` - (Optional, Available in v1.55.3+) A mapping of tags to assign to the resource.
* `resource_group_id` - (Optional, Available in 1.58.0+, Modifiable in 1.115.0+) The Id of resource group which the eip belongs.
* `deletion_protection` - (Optional, Available in v1.124.4+) Whether enable the deletion protection or not. Default value: `false`.
  - true: Enable deletion protection.
  - false: Disable deletion protection.
  
## Attributes Reference

The following attributes are exported:

* `id` - The EIP ID.
* `status` - The EIP current status.
* `ip_address` - The elastic ip address

## Import

Elastic IP address can be imported using the id, e.g.

```
$ terraform import alicloud_eip_address.example eip-abc12345678
```
