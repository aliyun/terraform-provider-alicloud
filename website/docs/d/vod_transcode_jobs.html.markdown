---
subcategory: "ApsaraVideo VoD (VOD)"
layout: "alicloud"
page_title: "Alicloud: alicloud_vod_transcode_jobs"
sidebar_current: "docs-alicloud-datasource-vod-transcode-jobs"
description: |-
  Provides a list of VOD Transcode Job owned by an Alibaba Cloud account.
---

# alicloud_vod_transcode_jobs

This data source provides the VOD Transcode Jobs of the current Alibaba Cloud user.

For information about VOD Transcode Job, see [What is ListSnapshots](https://www.alibabacloud.com/help/en/apsaravideo-for-vod/latest/listsnapshots).

-> **NOTE:** Available since v1.237.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_vod_transcode_jobs" "example" {
  video_id = "your_video_id"
}

output "job_ids" {
  value = data.alicloud_vod_transcode_jobs.example.jobs[*].transcode_job_id
}
```

## Argument Reference

The following arguments are supported:

* `video_id` - (Required) The ID of the video. This parameter specifies which video's snapshots to list.
* `snapshot_type` - (Optional) The type of the snapshot. Valid values: `Cover`, `Normal`.
* `ids` - (Optional) A list of transcode job IDs.
* `output_file` - (Optional) File name where to write data source results.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of transcode job IDs.
* `jobs` - A list of VOD Transcode Jobs. Each element contains the following attributes:
  * `id` - The ID of the transcode job.
  * `transcode_job_id` - The ID of the transcode job.
  * `video_id` - The ID of the video.
  * `create_time` - The creation time of the resource.
