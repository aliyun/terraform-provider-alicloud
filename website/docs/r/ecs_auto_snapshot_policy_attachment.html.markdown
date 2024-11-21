---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_auto_snapshot_policy_attachment"
sidebar_current: "docs-alicloud-resource-ecs-auto-snapshot-policy-attachment"
description: |-
  Provides a Alicloud ECS Auto Snapshot Policy Attachment resource.
---

# alicloud\_ecs\_auto\_snapshot\_policy\_attachment

Provides a ECS Auto Snapshot Policy Attachment resource.

For information about ECS Auto Snapshot Policy Attachment and how to use it, see [What is Auto Snapshot Policy Attachment](https://www.alibabacloud.com/help/en/doc-detail/25531.htm).

-> **NOTE:** Available in v1.122.0+.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ecs_auto_snapshot_policy_attachment&exampleId=7ee749ee-8e64-18a2-d1fd-de30aeb8b92d743a3145&activeTab=example&spm=docs.r.ecs_auto_snapshot_policy_attachment.0.7ee749ee8e&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
data "alicloud_zones" "example" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_kms_key" "example" {
  description            = "terraform-example"
  pending_window_in_days = "7"
  status                 = "Enabled"
}

resource "alicloud_ecs_auto_snapshot_policy" "example" {
  name            = "terraform-example"
  repeat_weekdays = ["1", "2", "3"]
  retention_days  = -1
  time_points     = ["1", "22", "23"]
}

resource "alicloud_ecs_disk" "example" {
  zone_id     = data.alicloud_zones.example.zones.0.id
  disk_name   = "terraform-example"
  description = "Hello ecs disk."
  category    = "cloud_efficiency"
  size        = "30"
  encrypted   = true
  kms_key_id  = alicloud_kms_key.example.id
  tags = {
    Name = "terraform-example"
  }
}

resource "alicloud_ecs_auto_snapshot_policy_attachment" "example" {
  auto_snapshot_policy_id = alicloud_ecs_auto_snapshot_policy.example.id
  disk_id                 = alicloud_ecs_disk.example.id
}
```

## Argument Reference

The following arguments are supported:

* `auto_snapshot_policy_id` - (Required, ForceNew) The auto snapshot policy id.
* `disk_id` - (Required, ForceNew) The disk id.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Auto Snapshot Policy Attachment. The value is formatted `<auto_snapshot_policy_id>:<disk_id>`.

## Import

ECS Auto Snapshot Policy Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecs_auto_snapshot_policy_attachment.example s-abcd12345:d-abcd12345
```
