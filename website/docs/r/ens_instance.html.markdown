---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_instance"
description: |-
  Provides a Alicloud ENS Instance resource.
---

# alicloud_ens_instance

Provides a ENS Instance resource. 

For information about ENS Instance and how to use it, see [What is Instance](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.208.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}


resource "alicloud_ens_instance" "default" {
  period = 1
  data_disk {
    size     = 20
    category = "local_ssd"
  }
  public_ip_identification   = true
  period_unit                = "Month"
  auto_renew                 = "False"
  scheduling_strategy        = "Concentrate"
  schedule_area_level        = "Region"
  image_id                   = "centos_6_08_64_20G_alibase_20171208"
  carrier                    = "cmcc"
  instance_type              = "ens.sn1.tiny"
  host_name                  = "testHost56"
  password                   = "Test123456@@"
  net_district_code          = "100102"
  amount                     = 1
  internet_charge_type       = "95BandwidthByMonth"
  instance_name              = var.name
  internet_max_bandwidth_out = 100
  ens_region_id              = "cn-hefei-cmcc-2"
  instance_charge_type       = "PrePaid"
  system_disk {
    size = 20
  }
  scheduling_price_strategy = "PriceHighPriority"
  user_data                 = "IyEvYmluL3NoCmVjaG8gIkhlbGxvIFdvcmxkLiAgVGhlIHRpbWUgaXMgbm93ICQoZGF0ZSAtUikhIiB8IHRlZSAvcm9vdC9vdXRwdXQudHh0"
  instance_charge_strategy  = "user"
}
```

## Argument Reference

The following arguments are supported:
* `amount` - (Required) Quantity.
* `auto_renew` - (Optional) Whether renew the fee automatically?it could be True,FalseDefault value is False.
* `auto_renew_period` - (Optional) The time period of auto renew. it will take effect.It could be 1, 2, 3, 6, 12. Default value is 1.
* `carrier` - (Optional) Operator.
* `data_disk` - (Optional, ForceNew) Data disk specifications. See [`data_disk`](#data_disk) below.
* `ens_region_id` - (Optional, ForceNew) Node id.
* `host_name` - (Optional, ForceNew) Host Name.
* `image_id` - (Optional, ForceNew) Image id.
* `instance_charge_strategy` - (Optional) InstanceChargeStrategy.
* `instance_charge_type` - (Required) Instance payment method, PrePaid: PrePaid, monthly PostPaid: pay-as-you-go.
* `instance_name` - (Optional, ForceNew) The instance name. It must be 2 to 128 characters in length and must start with an uppercase or lowercase letter or a Chinese character. It cannot start with http:// or https. It can contain Chinese, English, numbers, half-width colons (:), underscores (_), periods (.), or hyphens (-). The default value is the InstanceId of the instance.
* `instance_type` - (Required, ForceNew) Instance specifications.
* `internet_charge_type` - (Optional) Instance Charge type.it could be 95BandwidthByMonth, PayByBandwidth4thMonth.
* `internet_max_bandwidth_out` - (Required, ForceNew) The maximum public network bandwidth. If the value of the InternetMaxBandwidthOut parameter is greater than 0, public network IP is automatically allocated to the instance.
* `ip_type` - (Optional) ip type, It could be ipv4Andipv6,ipv4,ipv6.default value isi pv4.
* `net_district_code` - (Optional) Region code.
* `order_id` - (Optional) Order number.
* `password` - (Optional) The password of the instance.It is 8 to 30 characters in length and must contain three types of characters: uppercase and lowercase letters, numbers, and special symbols. The following special symbols can be set: '''()'~! @#$%^& *-_+ =|{}[]:;',.? /'''.
* `password_inherit` - (Optional) PasswordInherit.
* `period` - (Optional) Prepaid time period. Unit is month, it could be from 1 to 9 or 12. Default value is 1.
* `period_unit` - (Optional) Query the prices of ENS in different billing cycles. Value range:Month (default): The price unit of monthly billing.Day: Price unit of daily billing.
* `public_ip_identification` - (Optional) Whether the public IP address can be assigned to the specified instance. Value:
  - **true** (default): can be assigned.
  - **false**: cannot be assigned.
* `schedule_area_level` - (Required) Scheduling Hierarchy.
* `scheduling_price_strategy` - (Optional) Dispatch Price Policy.
* `scheduling_strategy` - (Optional) Scheduling Policy.
* `system_disk` - (Optional, ForceNew) System disk specifications. See [`system_disk`](#system_disk) below.
* `unique_suffix` - (Optional) Specifies whether to automatically append sequential suffixes to the hostnames specified by the HostName parameter and instance names specified by the InstanceName parameter when you create multiple instances at a time. The sequential suffix ranges from 001 to 999. Valid values:  true false Default value: false.
* `user_data` - (Optional) User data to pass to instance. [1, 16KB] characters.User data should not be base64 encoded. If you want to pass base64 encoded string to the property, use function Fn::Base64Decode to decode the base64 string first.

### `data_disk`

The data_disk supports the following:
* `category` - (Optional, ForceNew) Type of disk
  - High-efficiency cloud disk: cloud_efficiency
  - Full flash cloud disk: cloud_ssd
  - Local hdd disk: local_hdd
  - Local disk ssd:local_ssd.
* `size` - (Optional, ForceNew) Data disk size, unit: GB.

### `system_disk`

The system_disk supports the following:
* `size` - (Optional, ForceNew) The size of the system disk, in GB.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `payment_type` - The payment type of the resource.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Instance.
* `delete` - (Defaults to 5 mins) Used when delete the Instance.

## Import

ENS Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_ens_instance.example <id>
```