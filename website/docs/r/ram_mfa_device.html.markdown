---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_mfa_device"
description: |-
  Provides an Alibaba Cloud RAM MFA Device resource.
---

# alicloud_ram_mfa_device

Provides a RAM MFA Device resource.

For information about RAM MFA Device and how to use it, see [CreateVirtualMFADevice](https://next.api.alibabacloud.com/api/Ims/2019-08-15/CreateVirtualMFADevice).

-> **NOTE:** Available since v1.287.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_ram_mfa_device" "default" {
  virtual_mfa_device_name = "terraform-example"
}
```

## Argument Reference

The following arguments are supported:

* `virtual_mfa_device_name` - (Required, ForceNew) The name of the MFA device. The name must be 1 to 64 characters in length and can contain letters, digits, periods (.), and hyphens (-).

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource, which is the serial number of the MFA device.
* `activate_date` - The activation time of the MFA device. This value is empty if the device has not been activated.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the MFA Device.
* `delete` - (Defaults to 5 mins) Used when delete the MFA Device.

## Import

RAM MFA Device can be imported using the id (serial number), e.g.

```shell
$ terraform import alicloud_ram_mfa_device.example <serial_number>
```
