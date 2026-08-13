---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_common_bandwidth_eip_attachment"
description: |-
  Provides a Alicloud ENS Common Bandwidth Eip Attachment resource.
---

# alicloud_ens_common_bandwidth_eip_attachment

Provides a ENS Common Bandwidth Eip Attachment resource.

Bindings between shared bandwidth and EIP.

For information about ENS Common Bandwidth Eip Attachment and how to use it, see [What is Common Bandwidth Eip Attachment](https://next.api.alibabacloud.com/document/Ens/2017-11-10/AddCommonBandwidthPackageIps).

-> **NOTE:** Available since v1.289.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}

variable "ens_region_id" {
  default = "cn-chenzhou-telecom_unicom_cmcc"
}

resource "alicloud_ens_eip" "default" {
  bandwidth            = "5"
  eip_name             = var.name
  ens_region_id        = var.ens_region_id
  internet_charge_type = "95BandwidthByMonth"
  payment_type         = "PayAsYouGo"
}

resource "alicloud_ens_common_bandwidth_eip_attachment" "default" {
  bandwidth_package_id = "cbwp-xxxxx"
  ip_instance_id       = alicloud_ens_eip.default.id
}
```

## Argument Reference

The following arguments are supported:
* `bandwidth_package_id` - (Required, ForceNew) The ID of the ENS common bandwidth package.
* `ip_instance_id` - (Optional, ForceNew, Computed) The ID of the EIP (AllocationId) that is attached to the common bandwidth package.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<ip_instance_id>:<bandwidth_package_id>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Common Bandwidth Eip Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Common Bandwidth Eip Attachment.

## Import

ENS Common Bandwidth Eip Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_ens_common_bandwidth_eip_attachment.example <ip_instance_id>:<bandwidth_package_id>
```