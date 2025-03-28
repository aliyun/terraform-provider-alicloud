---
subcategory: "Max Compute"
layout: "alicloud"
page_title: "Alicloud: alicloud_max_compute_quota"
description: |-
  Provides a Alicloud Max Compute Quota resource.
---

# alicloud_max_compute_quota

Provides a Max Compute Quota resource.



For information about Max Compute Quota and how to use it, see [What is Quota](https://next.api.alibabacloud.com/document/MaxCompute/2022-01-04/CreateQuota).

-> **NOTE:** Available since v1.247.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraformexample"
}

provider "alicloud" {
  region = "cn-chengdu"
}

variable "part_nick_name" {
  default = "TFTest17292"
}

variable "sub_quota_nickname_3" {
  default = "sub398892"
}

variable "sub_quota_nickname_1" {
  default = "sub129792"
}

variable "sub_quota_nickname_2" {
  default = "sub223192"
}

resource "alicloud_max_compute_quota" "default" {
  payment_type   = "Subscription"
  part_nick_name = var.part_nick_name
  commodity_data = "{\"CU\":80,\"ord_time\":\"1:Month\",\"autoRenew\":false} "
  commodity_code = "odpsplus"
  sub_quota_info_list {
    parameter {
      min_cu              = "10"
      max_cu              = "60"
      enable_priority     = "false"
      force_reserved_min  = "false"
      scheduler_type      = "Fifo"
      single_job_cu_limit = "10"
    }

    nick_name = "os_${var.part_nick_name}"
    type      = "FUXI_OFFLINE"
  }
  sub_quota_info_list {
    parameter {
      min_cu             = "10"
      max_cu             = "10"
      scheduler_type     = "Fair"
      enable_priority    = "false"
      force_reserved_min = "false"
    }

    nick_name = var.sub_quota_nickname_1
    type      = "FUXI_OFFLINE"
  }
  sub_quota_info_list {
    nick_name = var.sub_quota_nickname_2
    type      = "FUXI_OFFLINE"
    parameter {
      min_cu             = "60"
      max_cu             = "60"
      scheduler_type     = "Fair"
      enable_priority    = "true"
      force_reserved_min = "true"
    }

  }
  tags = {
    "tf"    = "created"
    "valid" = "true"
  }
}
```

### Deleting `alicloud_max_compute_quota` or removing it from your configuration

Terraform cannot destroy resource `alicloud_max_compute_quota`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `commodity_code` - (Required) Valid values: odps_intl/odpsplus_intl/odps/odpsplus
  - *NOTE:** --odps_intl: International Station standard post-payment -- odpsplus_intl: International Station standard pre-payment -- odps: China Station standard post-payment -- odpsplus: China Station standard pre-payment
* `commodity_data` - (Optional) Define quota rules when creating quotas, for example:{"CU":50,"ord_time":"1:Month","autoRenew":false}.
CU (minimum quota resource size is 50)
ord_time 1:Month/Year (1 means quantity: Month/Year is the unit)
autoRenew (whether to enable automatic renewal)
If PaymentType is PayAsYouGo, you do not need to fill in
* `part_nick_name` - (Optional) Quota partial nickname, supports English letters and numbers, up to 24 characters

-> **NOTE:** If PaymentType is PayAsYouGo, you do not need to fill it in

* `payment_type` - (Required, ForceNew) Payment type. Valid values: Subscription/PayAsYouGo

-> **NOTE:** -- PayAsYouGo only needs to be opened once per region

* `sub_quota_info_list` - (Optional, Set) Secondary Quota list

-> **NOTE:** -- Add: If the configuration contains a second-level Quota that does not exist, a second-level Quota is added. -- Delete: If the configuration does not contain the existing secondary Quota, it will be deleted. -- Modify: If the configuration is inconsistent with the existing secondary Quota configuration parameters, the secondary Quota configuration will be updated. -- The default secondary Quota must be configured and cannot be deleted.
 See [`sub_quota_info_list`](#sub_quota_info_list) below.

### `sub_quota_info_list`

The sub_quota_info_list supports the following:
* `nick_name` - (Required) Secondary Quota nickname.

-> **NOTE:** -- Subscription: If you enter partNickName, the first-level QuotaNickName created is os_partNickName_p. Each first-level Quota has a default second-level Quota whose QuotaNickName is os_partNickName . -- The first-level quotanicname created by PayAsYouGo is os_PayAsYouGoQuota_p  by default, the second-level quotanicname is os_PayAsYouGoQuota

* `parameter` - (Optional, List) Parameter See [`parameter`](#sub_quota_info_list-parameter) below.
* `type` - (Optional, Computed) The secondary Quota type. The default value is: FUXI_OFFLINE

### `sub_quota_info_list-parameter`

The sub_quota_info_list-parameter supports the following:
* `enable_priority` - (Optional, Computed) Enable priority. Valid values: true/false, default: false
* `force_reserved_min` - (Optional, Computed) Exclusive or not. Valid values: true/false, default: false
* `max_cu` - (Required, Int) The value of maxCU in Reserved CUs.

-> **NOTE:**  The value of maxCU must be less than or equal to the value of maxCU in the level-1 quota that you purchased.

* `min_cu` - (Required, Int) The value of minCU in Reserved CUs.

-> **NOTE:**  -- The total value of minCU in all the level-2 quotas is equal to the value of minCU in the level-1 quota.    -- The value of minCU must be less than or equal to the value of maxCU in the level-2 quota and less than or equal to the value of minCU in the level-1 quota that you purchased.

* `scheduler_type` - (Optional, Computed) Scheduling policy. Valid values: Fifo/Fair, default: Fifo
* `single_job_cu_limit` - (Optional, Int) Single job CU upper limit. Valid value: greater than or equal to 1

-> **NOTE:** -- If you want to not restrict SingleJobCuLimit, please make sure that this parameter is not included in the configuration at all. That is, do not configure SingleJobCuLimit to "null" or any other invalid value


## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 31 mins) Used when create the Quota.
* `update` - (Defaults to 5 mins) Used when update the Quota.

## Import

Max Compute Quota can be imported using the id, e.g.

```shell
$ terraform import alicloud_max_compute_quota.example <id>
```