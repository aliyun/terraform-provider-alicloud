---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_eip"
description: |-
  Provides a Alicloud ENS Eip resource.
---

# alicloud_ens_eip

Provides a ENS Eip resource. Edge elastic public network IP. When you use it for the first time, please contact the product classmates to add a resource whitelist.

For information about ENS Eip and how to use it, see [What is Eip](https://www.alibabacloud.com/help/en/ens/developer-reference/api-createeipinstance).

-> **NOTE:** Available since v1.213.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_ens_eip" "default" {
  description   = "EipDescription_autotest"
  bandwidth     = "5"
  isp           = "cmcc"
  payment_type  = "PayAsYouGo"
  ens_region_id = "cn-chenzhou-telecom_unicom_cmcc"
  eip_name      = var.name

  internet_charge_type = "95BandwidthByMonth"
}
```

## Argument Reference

The following arguments are supported:
* `bandwidth` - (Optional, Computed) The peak bandwidth of the EIP to be specified.Rules:Default value: 5, value range: 5~10000, unit: Mbps. Example value: 5.
* `description` - (Optional) The description of the EIP.
* `eip_name` - (Optional) Name of the EIP instance.
* `ens_region_id` - (Required, ForceNew) Ens node ID.
* `internet_charge_type` - (Required, ForceNew) Billing type of the EIP instance. Valid value: 95bandwidthbymonth.
* `isp` - (Optional, ForceNew) Internet service provider, if not filled in, it will be assigned according to the default rules.
* `payment_type` - (Required, ForceNew) The billing type of the EIP instanceValue: PayAsYouGo.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the EIP instance.
* `status` - The status of the EIP instance.Rules:Associating: bindingUnassociating: UnbindingInUse: AssignedAvailable: AvailableCreating: CreatingReleasing: Releasing.

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