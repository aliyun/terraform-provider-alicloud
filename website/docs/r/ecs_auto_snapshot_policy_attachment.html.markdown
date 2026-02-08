---
subcategory: "ECS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecs_auto_snapshot_policy_attachment"
description: |-
  Provides a Alicloud ECS Auto Snapshot Policy Attachment resource.
---

# alicloud_ecs_auto_snapshot_policy_attachment

Provides a ECS Auto Snapshot Policy Attachment resource.

Automatic snapshot policy Mount relationship.

For information about ECS Auto Snapshot Policy Attachment and how to use it, see [What is Auto Snapshot Policy Attachment](https://www.alibabacloud.com/help/en/doc-detail/25531.htm).

-> **NOTE:** Available since v1.122.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ecs_auto_snapshot_policy_attachment&exampleId=7ee749ee-8e64-18a2-d1fd-de30aeb8b92d743a3145&activeTab=example&spm=docs.r.ecs_auto_snapshot_policy_attachment.0.7ee749ee8e&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_ecs_auto_snapshot_policy" "default" {
  auto_snapshot_policy_name = var.name
  repeat_weekdays           = ["1", "2", "3"]
  retention_days            = 1
  time_points               = ["1", "2", "3"]
}

resource "alicloud_ecs_disk" "default" {
  zone_id = data.alicloud_zones.default.zones.0.id
  size    = "500"
}

resource "alicloud_ecs_auto_snapshot_policy_attachment" "default" {
  auto_snapshot_policy_id = alicloud_ecs_auto_snapshot_policy.default.id
  disk_id                 = alicloud_ecs_disk.default.id
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ecs_auto_snapshot_policy_attachment&spm=docs.r.ecs_auto_snapshot_policy_attachment.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `auto_snapshot_policy_id` - (Required, ForceNew) The ID of the automatic snapshot policy that is applied to the cloud disk.
* `disk_id` - (Required, ForceNew) The ID of the disk.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<auto_snapshot_policy_id>:<disk_id>`.
* `region_id` - (Available since v1.271.0) The ID of the region where the automatic snapshot policy and the cloud disk are located.

## Timeouts

-> **NOTE:** Available since v1.271.0.

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Auto Snapshot Policy Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Auto Snapshot Policy Attachment.

## Import

ECS Auto Snapshot Policy Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_ecs_auto_snapshot_policy_attachment.example <auto_snapshot_policy_id>:<disk_id>
```
