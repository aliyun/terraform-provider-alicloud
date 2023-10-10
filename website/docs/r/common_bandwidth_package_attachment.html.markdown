---
subcategory: "EIP Bandwidth Plan (CBWP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_common_bandwidth_package_attachment"
description: |-
  Provides a Alicloud cbwp Common Bandwidth Package Attachment resource.
---

# alicloud_common_bandwidth_package_attachment

Provides a cbwp Common Bandwidth Package Attachment resource. -> **NOTE:** Terraform will auto build common bandwidth package attachment while it uses `alicloud_common_bandwidth_package_attachment` to build a common bandwidth package attachment resource.

For information about common bandwidth package and how to use it, see [What is Common Bandwidth Package](https://www.alibabacloud.com/help/product/55092.htm).

-> **NOTE:** From version 1.194.0, the resource can set the maximum bandwidth of an EIP that is associated with an EIP bandwidth plan by `bandwidth_package_bandwidth`. see [how to use it](https://www.alibabacloud.com/help/en/eip-bandwidth-plan/latest/120327).

For information about cbwp Common Bandwidth Package Attachment and how to use it, see [What is Common Bandwidth Package Attachment](https://www.alibabacloud.com/help/product/55092.htm).

-> **NOTE:** Available since v1.94.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_common_bandwidth_package" "default" {
  bandwidth            = 3
  internet_charge_type = "PayByBandwidth"
}

resource "alicloud_eip_address" "default" {
  bandwidth            = "3"
  internet_charge_type = "PayByTraffic"
}

resource "alicloud_common_bandwidth_package_attachment" "default" {
  bandwidth_package_id        = alicloud_common_bandwidth_package.default.id
  instance_id                 = alicloud_eip_address.default.id
  bandwidth_package_bandwidth = "2"
  ip_type                     = "EIP"
}
```

## Argument Reference

The following arguments are supported:
* `bandwidth_package_bandwidth` - (Optional, Computed) The maximum bandwidth for the EIP. This value cannot be larger than the maximum bandwidth of the EIP bandwidth plan. Unit: Mbit/s.
* `bandwidth_package_id` - (Required, ForceNew) The bandwidth_package_id of the common bandwidth package attachment, the field can't be changed.
* `cancel_common_bandwidth_package_ip_bandwidth` - (Optional) Whether to cancel the maximum bandwidth configuration for the EIP. Default: false.
* `instance_id` - (Required, ForceNew) The instance_id of the common bandwidth package attachment, the field can't be changed.
* `ip_type` - (Optional, Available since v1.211.0) IP type. Set the value to **EIP**, which indicates that the EIP is added to the Internet shared bandwidth.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<bandwidth_package_id>:<instance_id>`.
* `status` - The status of the Internet Shared Bandwidth instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Common Bandwidth Package Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Common Bandwidth Package Attachment.
* `update` - (Defaults to 5 mins) Used when update the Common Bandwidth Package Attachment.

## Import

cbwp Common Bandwidth Package Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_common_bandwidth_package_attachment.example <bandwidth_package_id>:<instance_id>
```