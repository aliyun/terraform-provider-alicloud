---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_instance"
description: |-
  Provides a Alicloud ENS Instance resource.
---

# alicloud_ens_instance

Provides a ENS Instance resource. 

For information about ENS Instance and how to use it, see [What is Instance](https://www.alibabacloud.com/help/en/ens/latest/create-instances).

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
    category = "cloud_efficiency"
  }
  data_disk {
    size     = 30
    category = "cloud_efficiency"
  }
  data_disk {
    size     = 40
    category = "cloud_efficiency"
  }
  public_ip_identification   = true
  period_unit                = "Month"
  scheduling_strategy        = "Concentrate"
  schedule_area_level        = "Region"
  image_id                   = "centos_6_08_64_20G_alibase_20171208"
  carrier                    = "cmcc"
  instance_type              = "ens.sn1.tiny"
  host_name                  = "exampleHost80"
  password                   = "Example123456@@"
  net_district_code          = "100102"
  internet_charge_type       = "95BandwidthByMonth"
  instance_name              = var.name
  internet_max_bandwidth_out = 100
  ens_region_id              = "cn-wuxi-telecom_unicom_cmcc-2"
  system_disk {
    size = 20
  }
  scheduling_price_strategy = "PriceHighPriority"
  user_data                 = "IyEvYmluL3NoCmVjaG8gIkhlbGxvIFdvcmxkLiAgVGhlIHRpbWUgaXMgbm93ICQoZGF0ZSAtUikhIiB8IHRlZSAvcm9vdC9vdXRwdXQudHh0"
  instance_charge_strategy  = "user"
  payment_type              = "Subscription"
}
```

### Deleting `alicloud_ens_instance` or removing it from your configuration

The `alicloud_ens_instance` resource allows you to manage  `payment_type = "Subscription"`  instance, but Terraform cannot destroy it.
Deleting the subscription resource or removing it from your configuration will remove it from your state file and management, but will not destroy the Instance.
You can resume managing the subscription instance via the AlibabaCloud Console.

## Argument Reference

The following arguments are supported:
* `auto_renew` - (Optional) Whether to automatically renew, default to False, this parameter is invalid when paying by volume.
* `carrier` - (Optional) Operator, required for regional level scheduling, invalid for node level scheduling.
* `data_disk` - (Optional, ForceNew) Data disk specifications. See [`data_disk`](#data_disk) below.
* `ens_region_id` - (Optional, ForceNew) Node id. When ScheduleAreaLevel is Region, EnsRegionId is required. When ScheduleAreaLevel is Big, Middle, Small, EnsRegionId is not required.
* `host_name` - (Optional, ForceNew, Computed) Host Name.
* `image_id` - (Optional, ForceNew) The Image Id field. If InstanceType is arm_bmi, the image Id is a non-required parameter. If instanceType is another specification value, the image Id is a required parameter.
* `instance_charge_strategy` - (Optional) Instance billing strategy, instance: instance granularity (prepaid method currently does not support instance), user: by user dimension (not transferred or prepaid method supports user).
* `instance_name` - (Optional, ForceNew) The instance name. It must be 2 to 128 characters in length and must start with an uppercase or lowercase letter or a Chinese character. It cannot start with http:// or https. It can contain Chinese, English, numbers, half-width colons (:), underscores (_), periods (.), or hyphens (-). The default value is the InstanceId of the instance.
* `instance_type` - (Required, ForceNew) Instance specifications type.
* `internet_charge_type` - (Optional) Instance Charge type.it could be BandwidthByDay, 95BandwidthByMonth, PayByBandwidth4thMonth.
* `internet_max_bandwidth_out` - (Required, ForceNew) The maximum public network bandwidth.
* `net_district_code` - (Optional) Region code, required for regional level scheduling, invalid for node level scheduling.
* `password` - (Optional) The password of the instance。It is 8 to 30 characters in length and must contain three types of characters: uppercase and lowercase letters, numbers, and special symbols. The following special symbols can be set: '''()'~! @#$%^& *-_+ =|{}[]:;',.? /'''.
* `password_inherit` - (Optional) Whether to use image preset password prompt: Password and KeyPairNamePasswordInherit must be passed.
* `payment_type` - (Required, ForceNew) Instance payment method, Subscription: prepaid, monthly package; PayAsYouGo: Pay as you go.
* `period` - (Optional) The duration of purchasing resources. If PeriodUnit is not specified, it defaults to purchasing on a monthly basis. Currently, only days and months are supported. If PeriodUnit=Day, Period can only be 3. If PeriodUnit=Monthc, then Period can be 1-9,12.
* `period_unit` - (Optional) The unit of time for purchasing resources. If PeriodUnit is not specified, it defaults to purchasing by Month. Currently, only days and months are supported. If PeriodUnit=Day, Period can only be 3. If PeriodUnit=Month, then Period can be 1-9,12.
* `public_ip_identification` - (Optional) Whether to allocate public IP. Value：true (default): can be assigned，false: cannot be assigned.
* `quantity` - (Optional) Number of instances.
* `schedule_area_level` - (Required) Scheduling level, which is used to perform node level or regional scheduling.
* `scheduling_price_strategy` - (Optional) Dispatch price strategy. If left blank, it defaults to prioritizing low prices. Values: PriceLowPriority (priority high price), PriceLowPriority (priority low price).
* `scheduling_strategy` - (Optional) When scheduling at the node level, it is Concentrate. When scheduling at the regional level, it is selected according to customer needs. Concentrate: Centralized; Disperse: Disperse.
* `system_disk` - (Optional, ForceNew) The field representing the system disk specification. SystemDisk is a non-required parameter when InstanceType is x86_pm,x86_bmi,x86_bm,pc_bmi, or arm_bmi. SystemDisk is a required parameter when instanceType is other specification families. See [`system_disk`](#system_disk) below.
* `unique_suffix` - (Optional) Specifies whether to automatically append sequential suffixes to the hostnames specified by the HostName parameter and instance names specified by the InstanceName parameter when you create multiple instances at a time. The sequential suffix ranges from 001 to 999. Valid values:  true false Default value: false.
* `user_data` - (Optional) User defined data, with a maximum support of 16KB. You can input UserData information. UserData encoded in Base64 format.

### `data_disk`

The data_disk supports the following:
* `category` - (Optional, ForceNew) Type of dataDisk
  - cloud_efficiency：High-efficiency cloud disk
  - cloud_ssd：Full flash cloud disk
  - local_hdd：Local hdd disk
  - local_ssd：Local disk ssd.
* `size` - (Optional, ForceNew) Data disk size, cloud_efficiency is 20-32000,cloud_ssd/local_hdd/local_ssd is 20-25000, unit: GB.

### `system_disk`

The system_disk supports the following:
* `size` - (Optional, ForceNew) System disk size, cloud_efficiency is 20-32000,cloud_ssd/local_hdd/local_ssd is 20-25000, unit: GB.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `status` - the status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 10 mins) Used when create the Instance.
* `delete` - (Defaults to 10 mins) Used when delete the Instance.

## Import

ENS Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_ens_instance.example <id>
```