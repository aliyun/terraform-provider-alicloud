---
subcategory: "Video Surveillance System"
layout: "alicloud"
page_title: "Alicloud: alicloud_video_surveillance_system_group"
sidebar_current: "docs-alicloud-resource-video-surveillance-system-group"
description: |-
  Provides a Alicloud Video Surveillance System Group resource.
---

# alicloud\_video\_surveillance\_system\_group

Provides a Video Surveillance System Group resource.

For information about Video Surveillance System Group and how to use it, see [What is Group](https://help.aliyun.com/product/108765.html).

-> **NOTE:** Available in v1.135.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_video_surveillance_system_group" "default" {
  group_name   = "your_group_name"
  in_protocol  = "rtmp"
  out_protocol = "flv"
  play_domain  = "your_plan_domain"
  push_domain  = "your_push_domain"
}
```

## Argument Reference

The following arguments are supported:
* `group_name` - (Required) The Group Name.
* `in_protocol` - (Required) The use of the access protocol support gb28181, Real Time Messaging Protocol (rtmp). Valid values: `gb28181`, `rtmp`.
* `out_protocol` - (Required) The playback protocol used by the space, multiple values are separated by commas (,). Valid values: `flv`,`hls`, `rtmp`.
* `play_domain` - (Required,ForceNew) The domain name of plan streaming used by the group.
* `push_domain` - (Required,ForceNew) The domain name of push streaming used by the group.
* `callback` - (Optional) The space within the device status update of the callback, need to start with http:// or https:// at the beginning.
* `enabled` - (Optional) Whether to open Group.
* `description` - (Optional) The description of Group.


## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Group.
* `capture_image` - The capture image.
* `capture_interval` - The capture interval.
* `capture_oss_bucket` - The capture oss bucket.
* `capture_oss_path` - The capture oss path.
* `capture_video` - The capture video.
* `lazy_pull` - Whether to enable on-demand streaming. Default value:`false`.
* `status` - Whether to open Group. Valid values: `on`,`off`.

## Import

Video Surveillance System Group can be imported using the id, e.g.

```
$ terraform import alicloud_video_surveillance_system_group.example <id>
```
