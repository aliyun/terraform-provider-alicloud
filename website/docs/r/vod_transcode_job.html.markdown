---
subcategory: "ApsaraVideo VoD (VOD)"
layout: "alicloud"
page_title: "Alicloud: alicloud_vod_transcode_job"
sidebar_current: "docs-alicloud-resource-vod-transcode-job"
description: |-
  Provides a Alicloud VOD Transcode Job resource.
---

# alicloud_vod_transcode_job

Provides a VOD Transcode Job resource.

For information about VOD Transcode Job and how to use it, see [What is SubmitSnapshotJob](https://www.alibabacloud.com/help/en/apsaravideo-for-vod/latest/submitsnapshotjob).

-> **NOTE:** Available since v1.237.0.

## Example Usage

Basic Usage

```terraform
resource "alicloud_vod_transcode_job" "example" {
  video_id               = "your_video_id"
  reference_id           = "tf-example-ref-id"
  height                 = "480"
  width                  = "640"
  count                  = 10
  interval               = 10
  user_data              = "tf-example-userdata"
  specified_offset_times = [0, 1000, 2000]
}
```

## Argument Reference

The following arguments are supported:

* `video_id` - (Required, ForceNew) The ID of the video. This parameter specifies which video to capture snapshots from.
* `reference_id` - (Optional, ForceNew) The custom ID. Only lowercase letters, uppercase letters, digits, hyphens (-), and underscores (_) are supported. The length must be 6 to 64 characters. The custom ID must be unique among the resources that belong to the same Alibaba Cloud account.
* `height` - (Optional, ForceNew) The height of the snapshot. Unit: pixel.
* `width` - (Optional, ForceNew) The width of the snapshot. Unit: pixel.
* `specified_offset_time` - (Optional, ForceNew) The start time of the time range in which snapshots are captured. Unit: milliseconds.
* `count` - (Optional, ForceNew) The maximum number of snapshots to capture.
* `interval` - (Optional, ForceNew) The interval at which snapshots are captured. Unit: milliseconds.
* `sprite_snapshot_config` - (Optional, ForceNew) The configuration of the sprite snapshot, in JSON format.
* `snapshot_template_id` - (Optional, ForceNew) The ID of the snapshot template.
* `user_data` - (Optional, ForceNew) The custom configuration.
* `specified_offset_times` - (Optional, ForceNew) The user-defined snapshot time points. Unit: milliseconds.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Transcode Job.
* `transcode_job_id` - The ID of the Transcode Job.
* `create_time` - The creation time of the resource.
* `region_id` - The region ID of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when create the Transcode Job.
* `delete` - (Defaults to 10 mins) Used when delete the Transcode Job.

## Import

VOD Transcode Job can be imported using the id, e.g.

```shell
$ terraform import alicloud_vod_transcode_job.example <id>
```
