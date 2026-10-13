---
subcategory: "RAM"
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_mfa_devices"
description: |-
  Provides a list of RAM MFA Devices to the user.
---

# alicloud_ram_mfa_devices

This data source provides a list of RAM MFA Devices available to the user.

-> **NOTE:** Available since v1.287.0.

## Example Usage

```terraform
data "alicloud_ram_mfa_devices" "default" {
  name_regex = "my-mfa-device"
}

output "first_mfa_serial" {
  value = data.alicloud_ram_mfa_devices.default.mfa_devices[0].serial_number
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to filter results by the MFA device name. The device name is extracted from the serial number.
* `ids` - (Optional) A list of MFA device serial numbers. If specified, only devices with matching serial numbers are returned.
* `output_file` - (Optional) File name where to write the data source results in JSON format.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of MFA device serial numbers.
* `mfa_devices` - A list of MFA Devices. Each element contains the following attributes:
  * `serial_number` - The serial number of the MFA device.
  * `virtual_mfa_device_name` - The name of the MFA device, extracted from the serial number.
  * `activate_date` - The activation time of the MFA device. Empty if the device has not been activated.
