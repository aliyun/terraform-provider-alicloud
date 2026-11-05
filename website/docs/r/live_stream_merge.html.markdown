---
subcategory: "Live"
layout: "alicloud"
page_title: "Alicloud: alicloud_live_stream_merge"
description: |-
  Provides a Alicloud Live Stream Merge resource.
---

# alicloud_live_stream_merge

Provides a Live Stream Merge resource.

Live Stream Merge (主备合流) merges two ingested live streams (a primary and a backup) into one output stream for the same play domain. When the primary stream is interrupted, the merge rule switches to the backup stream automatically or manually according to the configured switch mode and merge parameters.

For information about Live Stream Merge and how to use it, see [AddLiveStreamMerge](https://www.alibabacloud.com/help/en/live/developer-reference/api-live-2016-11-01-addlivestreammerge).

-> **NOTE:** Available since v2.2.0.

## Example Usage

Basic Usage

```terraform
variable "domain_name" {
  default = "example.com"
}

provider "alicloud" {
  region = "cn-shanghai"
}

resource "alicloud_live_stream_merge" "default" {
  domain_name      = var.domain_name
  app_name         = "liveApp"
  stream_name      = "liveStream"
  in_app_name_1    = "liveApp1"
  in_stream_name_1 = "liveStream1"
  in_app_name_2    = "liveApp2"
  in_stream_name_2 = "liveStream2"
  start_time       = "2020-05-29T00:00:00Z"
  end_time         = "2020-05-29T01:00:00Z"
  protocol         = "rtmp"
  live_merger      = "off"
  switch_mode      = "manually"
}
```

## Argument Reference

The following arguments are supported:

* `domain_name` - (Required, ForceNew) The play domain name of the output stream. Changing this creates a new merge rule.
* `app_name` - (Required, ForceNew) The application name of the output stream. It must be different from `in_app_name_1` and `in_app_name_2`. Changing this creates a new merge rule.
* `stream_name` - (Required, ForceNew) The stream name of the output stream. It must be different from `in_stream_name_1` and `in_stream_name_2`. Changing this creates a new merge rule.
* `in_app_name_1` - (Required) The application name of the first input stream.
* `in_stream_name_1` - (Required) The stream name of the first input stream.
* `in_app_name_2` - (Required) The application name of the second input stream.
* `in_stream_name_2` - (Required) The stream name of the second input stream.
* `start_time` - (Required) The start time of the merge. The date format follows ISO8601 and uses UTC+0, in the format `yyyy-MM-ddTHH:mm:ssZ`.
* `end_time` - (Required) The end time of the merge. The date format follows ISO8601 and uses UTC+0, in the format `yyyy-MM-ddTHH:mm:ssZ`. The start time and end time can differ by at most 7 days.
* `protocol` - (Optional, Computed) The live stream protocol. Valid values: `rtmp` (default), `rtc`.
* `live_merger` - (Optional, Computed) Whether to use the liveswitch engine. Valid values: `on` (use the new liveswitch engine), `off` (default, use the legacy engine such as rtmpr).
* `merge_parameters` - (Optional, Computed) The configuration parameters for stream switching, in query string format such as `ali_low_frame_rate_threshold=10&ali_max_no_frame_timeout=5&block_all_jitter=0`. It controls when the primary-backup switch is triggered:
  - When a stream break (eof or net err) occurs, the switch is triggered immediately.
  - When the stall rate in the last 5 seconds exceeds 60%, the switch is triggered immediately.
  - When pulling times out and no frame data is received for 2 seconds, the switch is triggered.
  - When the average frame rate within `ali_max_no_frame_timeout` is below `ali_low_frame_rate_threshold`, the switch is triggered.
  - When `block_all_jitter` is enabled, the previous three rules do not take effect.

  Value range: `10 > timeout > 2`, `200 > threshold > 1`, `block` is `0` or `1`.
* `select_app_name` - (Optional, Computed) The application name of the input stream selected as the primary source. This is a write-only field; the configured value is preserved in state and is not returned by Describe.
* `select_stream_name` - (Optional, Computed) The stream name of the input stream selected as the primary source. This is a write-only field; the configured value is preserved in state and is not returned by Describe.
* `switch_mode` - (Optional, Computed) The switch mode of the primary-backup merge. Valid values: `manually`, `automatically`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in format `<domain_name>:<app_name>:<stream_name>`.
* `app_using` - The application name of the input stream currently in use.
* `stream_using` - The stream name of the input stream currently in use.
* `extra_in_app_streams` - The extra input application streams.

## Import

Live Stream Merge can be imported using the id, e.g.

```shell
$ terraform import alicloud_live_stream_merge.example <domain_name>:<app_name>:<stream_name>
```
