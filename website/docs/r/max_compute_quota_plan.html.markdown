---
subcategory: "Max Compute"
layout: "alicloud"
page_title: "Alicloud: alicloud_max_compute_quota_plan"
description: |-
  Provides a Alicloud Max Compute Quota Plan resource.
---

# alicloud_max_compute_quota_plan

Provides a Max Compute Quota Plan resource.



For information about Max Compute Quota Plan and how to use it, see [What is Quota Plan](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.242.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_max_compute_quota_plan&exampleId=656ae932-d8b1-745f-2025-59b2f53acf243b1d78f0&activeTab=example&spm=docs.r.max_compute_quota_plan.0.656ae932d8&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


variable "elastic_reserved_cu" {
  default = "50"
}

resource "alicloud_max_compute_quota_plan" "default" {
  nickname = "os_terrform_p"
  quota {
    parameter {
      elastic_reserved_cu = var.elastic_reserved_cu
    }

    sub_quota_info_list {
      nick_name = "sub_quota"
      parameter {
        min_cu              = "0"
        max_cu              = "20"
        elastic_reserved_cu = var.elastic_reserved_cu
      }

    }
    sub_quota_info_list {
      nick_name = "os_terrform"
      parameter {
        min_cu              = "50"
        max_cu              = "50"
        elastic_reserved_cu = "0"
      }

    }

  }

  plan_name = "quota_plan"
}
```

## Argument Reference

The following arguments are supported:
* `is_effective` - (Optional) Whether to take effect immediately. “Valid values: true”  
.-> **NOTE:** when other quota plans in the same quota group take effect, the effective quota group will become invalid. That is, IsEffective will become false. The effective quota plan cannot be deleted.
* `nickname` - (Required, ForceNew) Quota Name
* `plan_name` - (Required, ForceNew) The Quota plan name. Start with a letter, containing letters, numbers, and underscores (_). It is no more than 64 characters long.
* `quota` - (Optional, ForceNew, List) Quota property See [`quota`](#quota) below.

### `quota`

The quota supports the following:
* `parameter` - (Optional, ForceNew, List) The parameters of level-1 quota. See [`parameter`](#quota-parameter) below.
* `sub_quota_info_list` - (Optional, ForceNew, Set) Secondary Quota list

-> **NOTE:** need to list all secondary Quota
 See [`sub_quota_info_list`](#quota-sub_quota_info_list) below.

### `quota-parameter`

The quota-parameter supports the following:
* `elastic_reserved_cu` - (Required, Int) The value of elastic Reserved CUs.

### `quota-sub_quota_info_list`

The quota-sub_quota_info_list supports the following:
* `nick_name` - (Required, ForceNew) The nickname of the level-2 quota.
* `parameter` - (Optional, List) Level 2 Quota CU configuration See [`parameter`](#quota-sub_quota_info_list-parameter) below.

### `quota-sub_quota_info_list-parameter`

The quota-sub_quota_info_list-parameter supports the following:
* `elastic_reserved_cu` - (Required, Int) The value of elastic Reserved CUs.

-> **NOTE:**  The total number of elastically reserved CUs in all the level-2 quotas is equal to the number of elastically reserved CUs in the level-1 quota..

* `max_cu` - (Required, Int) The value of maxCU in Reserved CUs.

-> **NOTE:**  The value of maxCU must be less than or equal to the value of maxCU in the level-1 quota that you purchased.

* `min_cu` - (Required, Int) The value of minCU in Reserved CUs.

-> **NOTE:**  -- The total value of minCU in all the level-2 quotas is equal to the value of minCU in the level-1 quota.    -- The value of minCU must be less than or equal to the value of maxCU in the level-2 quota and less than or equal to the value of minCU in the level-1 quota that you purchased.


## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<nickname>:<plan_name>`.
* `quota` - Quota property
  * `parameter` - The parameters of level-1 quota.
    * `max_cu` - The value of maxCU in Reserved CUs.
    * `min_cu` - The value of minCU in Reserved CUs.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Quota Plan.
* `delete` - (Defaults to 5 mins) Used when delete the Quota Plan.
* `update` - (Defaults to 5 mins) Used when update the Quota Plan.

## Import

Max Compute Quota Plan can be imported using the id, e.g.

```shell
$ terraform import alicloud_max_compute_quota_plan.example <nickname>:<plan_name>
```