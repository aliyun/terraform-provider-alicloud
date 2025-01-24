---
subcategory: "Max Compute"
layout: "alicloud"
page_title: "Alicloud: alicloud_max_compute_tunnel_quota_timer"
description: |-
  Provides a Alicloud Max Compute Tunnel Quota Timer resource.
---

# alicloud_max_compute_tunnel_quota_timer

Provides a Max Compute Tunnel Quota Timer resource.



For information about Max Compute Tunnel Quota Timer and how to use it, see [What is Tunnel Quota Timer](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.242.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_max_compute_tunnel_quota_timer&exampleId=7184f4c9-2fad-08f1-e2c6-4b540adb251118b1222d&activeTab=example&spm=docs.r.max_compute_tunnel_quota_timer.0.7184f4c92f&intl_lang=EN_US" target="_blank">
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


resource "alicloud_max_compute_tunnel_quota_timer" "default" {
  quota_timer {
    begin_time = "00:00"
    end_time   = "01:00"
    tunnel_quota_parameter {
      slot_num                  = "50"
      elastic_reserved_slot_num = "50"
    }
  }
  quota_timer {
    begin_time = "01:00"
    end_time   = "02:00"
    tunnel_quota_parameter {
      slot_num                  = "50"
      elastic_reserved_slot_num = "50"
    }
  }
  quota_timer {
    begin_time = "02:00"
    end_time   = "24:00"
    tunnel_quota_parameter {
      slot_num                  = "50"
      elastic_reserved_slot_num = "50"
    }
  }
  nickname  = "ot_terraform_p"
  time_zone = "Asia/Shanghai"
}
```

### Deleting `alicloud_max_compute_tunnel_quota_timer` or removing it from your configuration

Terraform cannot destroy resource `alicloud_max_compute_tunnel_quota_timer`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `nickname` - (Required, ForceNew) The nickname of the exclusive Resource Group (Tunnel Quota) for the level - 1 data transmission service.
* `quota_timer` - (Optional, Set) Time-Sharing configuration

-> **NOTE:** -- The same reserved Quota resource group supports up to 48 time intervals. The minimum duration of a time interval is 30 minutes. -- After the current data transmission service is configured for time-sharing, if you need to perform a downgrade operation on the data transmission service (package year and month), please reduce the time-sharing concurrency first. -- The effective time of the time-sharing configuration is 0 to 5 minutes, and the billing will be calculated according to the actual effective time. -- Please make sure to set the time range completely from 00:00 to 24:00
 See [`quota_timer`](#quota_timer) below.
* `time_zone` - (Optional) Time zone, reference: Asia/Shanghai
In general, the system will automatically generate the time zone according to the region without configuration.

### `quota_timer`

The quota_timer supports the following:
* `begin_time` - (Required) The time-sharing configuration start time. Reference value: 00:00
* `end_time` - (Required) The end time of the timesharing configuration. Reference value: 24:00
* `tunnel_quota_parameter` - (Optional, List) Time-sharing configuration parameters. See [`tunnel_quota_parameter`](#quota_timer-tunnel_quota_parameter) below.

### `quota_timer-tunnel_quota_parameter`

The quota_timer-tunnel_quota_parameter supports the following:
* `elastic_reserved_slot_num` - (Required, Int) The number of elastic reserved concurrency (Slot).
* `slot_num` - (Required, Int) The number of reserved concurrency (Slot).

-> **NOTE:** The reserved concurrency (Slot) cannot be modified. The number of concurrency slots must be the same as that of the purchased tunnel quota.


## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Tunnel Quota Timer.
* `update` - (Defaults to 5 mins) Used when update the Tunnel Quota Timer.

## Import

Max Compute Tunnel Quota Timer can be imported using the id, e.g.

```shell
$ terraform import alicloud_max_compute_tunnel_quota_timer.example <id>
```