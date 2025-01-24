---
subcategory: "Max Compute"
layout: "alicloud"
page_title: "Alicloud: alicloud_max_compute_quota_schedule"
description: |-
  Provides a Alicloud Max Compute Quota Schedule resource.
---

# alicloud_max_compute_quota_schedule

Provides a Max Compute Quota Schedule resource.



For information about Max Compute Quota Schedule and how to use it, see [What is Quota Schedule](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.242.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_max_compute_quota_schedule&exampleId=8409f2ba-063c-2e64-7b88-6b78d13eeaa4f93b2301&activeTab=example&spm=docs.r.max_compute_quota_schedule.0.8409f2ba06&intl_lang=EN_US" target="_blank">
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
  default = "0"
}

variable "quota_nick_name" {
  default = "os_terrform_p"
}

resource "alicloud_max_compute_quota_plan" "default" {
  quota {
    parameter {
      elastic_reserved_cu = 50
    }
    sub_quota_info_list {
      nick_name = "sub_quota"
      parameter {
        min_cu              = "0"
        max_cu              = "20"
        elastic_reserved_cu = "30"
      }
    }
    sub_quota_info_list {
      nick_name = "os_terrform"
      parameter {
        min_cu              = "50"
        max_cu              = "50"
        elastic_reserved_cu = "20"
      }

    }
  }

  plan_name = "quota_plan1"
  nickname  = "os_terrform_p"
}

resource "alicloud_max_compute_quota_plan" "default2" {
  quota {
    parameter {
      elastic_reserved_cu = 50
    }
    sub_quota_info_list {
      nick_name = "sub_quota"
      parameter {
        min_cu              = "0"
        max_cu              = "20"
        elastic_reserved_cu = "20"
      }
    }
    sub_quota_info_list {
      nick_name = "os_terrform"
      parameter {
        min_cu              = "50"
        max_cu              = "50"
        elastic_reserved_cu = "30"
      }

    }
  }

  plan_name = "quota_plan2"
  nickname  = "os_terrform_p"
}

resource "alicloud_max_compute_quota_plan" "default3" {
  quota {
    parameter {
      elastic_reserved_cu = 50
    }
    sub_quota_info_list {
      nick_name = "sub_quota"
      parameter {
        min_cu              = "40"
        max_cu              = "40"
        elastic_reserved_cu = "40"
      }
    }
    sub_quota_info_list {
      nick_name = "os_terrform"
      parameter {
        min_cu              = "10"
        max_cu              = "10"
        elastic_reserved_cu = "10"
      }

    }
  }

  plan_name = "quota_plan3"
  nickname  = "os_terrform_p"
}

resource "alicloud_max_compute_quota_schedule" "default" {
  timezone = "UTC+8"
  nickname = var.quota_nick_name
  schedule_list {
    plan = "Default"
    condition {
      at = "00:00"
    }

    type = "daily"
  }

  # schedule_list {
  #   plan = "${alicloud_max_compute_quota_plan.default.plan_name}"
  #     condition {
  #     at = "00:00"
  #   }

  #   type = "daily"
  # }
  # schedule_list {
  #   type = "daily"
  #   plan = "${alicloud_max_compute_quota_plan.default2.plan_name}"
  #     condition {
  #     at = "01:00"
  #   }

  # }
  # schedule_list {
  #   plan = "${alicloud_max_compute_quota_plan.default3.plan_name}"
  #     condition {
  #     at = "02:00"
  #   }

  #   type = "daily"
  # }

}
```

### Deleting `alicloud_max_compute_quota_schedule` or removing it from your configuration

Terraform cannot destroy resource `alicloud_max_compute_quota_schedule`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `nickname` - (Required, ForceNew) The nickname of level-1 compute quota.
* `schedule_list` - (Optional, Set) schedule list See [`schedule_list`](#schedule_list) below.
* `timezone` - (Required, ForceNew) Time zone, reference value: UTC +8

### `schedule_list`

The schedule_list supports the following:
* `condition` - (Optional, List) The value of effective condition. See [`condition`](#schedule_list-condition) below.
* `plan` - (Required) The name of the quota plan.
* `type` - (Required, ForceNew) The type of the quota plan. Valid values: daily 

-> **NOTE:** Currently, only daily is supported.


### `schedule_list-condition`

The schedule_list-condition supports the following:
* `at` - (Required) Effective time. The format is HH:mm, sample value: 00:00

-> **NOTE:** The configuration must start from the effective time of 00:00. The input time must be either a whole hour or a half hour, and the minimum interval between each schedule is 30 minutes.


## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<nickname>:<timezone>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Quota Schedule.
* `update` - (Defaults to 5 mins) Used when update the Quota Schedule.

## Import

Max Compute Quota Schedule can be imported using the id, e.g.

```shell
$ terraform import alicloud_max_compute_quota_schedule.example <nickname>:<timezone>
```