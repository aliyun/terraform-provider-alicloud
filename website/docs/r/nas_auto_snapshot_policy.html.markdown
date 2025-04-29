---
subcategory: "File Storage (NAS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_nas_auto_snapshot_policy"
description: |-
  Provides a Alicloud NAS Auto Snapshot Policy resource.
---

# alicloud_nas_auto_snapshot_policy

Provides a NAS Auto Snapshot Policy resource. Automatic snapshot policy.

For information about NAS Auto Snapshot Policy and how to use it, see [What is Auto Snapshot Policy](https://www.alibabacloud.com/help/en/doc-detail/135662.html)).

-> **NOTE:** Available since v1.153.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_nas_auto_snapshot_policy&exampleId=3b2734b5-d239-3b5a-1f90-4607655ceda98588432b&activeTab=example&spm=docs.r.nas_auto_snapshot_policy.0.3b2734b5d2&intl_lang=EN_US" target="_blank">
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


resource "alicloud_nas_auto_snapshot_policy" "default" {
  time_points               = ["0", "1", "2"]
  retention_days            = "1"
  repeat_weekdays           = ["2", "3", "4"]
  auto_snapshot_policy_name = var.name
  file_system_type          = "extreme"
}
```

## Argument Reference

The following arguments are supported:
* `auto_snapshot_policy_name` - (Optional) The name of the automatic snapshot policy. Limits:
  - The name must be `2` to `128` characters in length,
  - The name must start with a letter.
  - The name can contain digits, colons (:), underscores (_), and hyphens (-). The name cannot start with `http://` or `https://`.
  - The value of this parameter is empty by default.
* `file_system_type` - (Optional, ForceNew, Computed, Available since v1.223.2) The file system type.
* `repeat_weekdays` - (Required) The day on which an auto snapshot is created.
  - A maximum of 7 time points can be selected.
  - The format is  an JSON array of ["1", "2", … "7"]  and the time points are separated by commas (,).
* `retention_days` - (Optional, Computed) The number of days for which you want to retain auto snapshots. Unit: days. Valid values:
  - `-1`: the default value. Auto snapshots are permanently retained. After the number of auto snapshots exceeds the upper limit, the earliest auto snapshot is automatically deleted.
  - `1` to `65536`: Auto snapshots are retained for the specified days. After the retention period of auto snapshots expires, the auto snapshots are automatically deleted.
* `time_points` - (Required) The point in time at which an auto snapshot is created.
  - A maximum of 24 time points can be selected.
  - The format is  an JSON array of ["0", "1", … "23"] and the time points are separated by commas (,).

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - Creation time.
* `status` - The status of the automatic snapshot policy.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Auto Snapshot Policy.
* `delete` - (Defaults to 5 mins) Used when delete the Auto Snapshot Policy.
* `update` - (Defaults to 5 mins) Used when update the Auto Snapshot Policy.

## Import

NAS Auto Snapshot Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_nas_auto_snapshot_policy.example <id>
```