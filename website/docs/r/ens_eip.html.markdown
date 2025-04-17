---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_eip"
description: |-
  Provides a Alicloud ENS Eip resource.
---

# alicloud_ens_eip

Provides a ENS Eip resource.

Edge elastic public network IP. When you use it for the first time, please contact the product classmates to add a resource whitelist.

For information about ENS Eip and how to use it, see [What is Eip](https://www.alibabacloud.com/help/en/ens/developer-reference/api-createeipinstance).

-> **NOTE:** Available since v1.213.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ens_eip&exampleId=4516c73b-69f9-683b-790b-d52399bd7124a8581863&activeTab=example&spm=docs.r.ens_eip.0.4516c73b69&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_ens_eip" "default" {
  description          = "EipDescription_autotest"
  bandwidth            = "5"
  isp                  = "cmcc"
  payment_type         = "PayAsYouGo"
  ens_region_id        = "cn-chenzhou-telecom_unicom_cmcc"
  eip_name             = var.name
  internet_charge_type = "95BandwidthByMonth"
}
```

## Argument Reference

The following arguments are supported:
* `bandwidth` - (Optional, Computed) The maximum bandwidth of the EIP. Default value: `5`. Valid values: `5` to `10000`. Unit: Mbit/s.
* `description` - (Optional) The description of the EIP.
* `eip_name` - (Optional) The name of the EIP.
* `ens_region_id` - (Required, ForceNew) Ens node ID.
* `internet_charge_type` - (Required, ForceNew) The metering method of the EIP. Valid value: `95BandwidthByMonth`.
* `isp` - (Optional, ForceNew) The Internet service provider. Valid value: `cmcc`, `unicom`, `telecom`.
* `payment_type` - (Required, ForceNew) The billing method of the EIP. Valid value: `PayAsYouGo`.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the EIP instance.
* `status` - The status of the EIP.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Eip.
* `delete` - (Defaults to 5 mins) Used when delete the Eip.
* `update` - (Defaults to 5 mins) Used when update the Eip.

## Import

ENS Eip can be imported using the id, e.g.

```shell
$ terraform import alicloud_ens_eip.example <id>
```
