---
subcategory: "Aligreen"
layout: "alicloud"
page_title: "Alicloud: alicloud_aligreen_oss_stock_task"
description: |-
  Provides a Alicloud Aligreen Oss Stock Task resource.
---

# alicloud_aligreen_oss_stock_task

Provides a Aligreen Oss Stock Task resource.

OSS stock file scanning task.

For information about Aligreen Oss Stock Task and how to use it, see [What is Oss Stock Task](https://next.api.alibabacloud.com/document/Green/2017-08-23/CreateOssStockTask).

-> **NOTE:** Available since v1.228.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform"
}

provider "alicloud" {
  region = "cn-shanghai"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_oss_bucket" "defaultPyhXOV" {
  storage_class = "Standard"
  bucket        = "${var.name}-${random_integer.default.result}"
}

resource "alicloud_aligreen_callback" "defaultJnW8Na" {
  callback_url         = "https://www.aliyun.com/"
  crypt_type           = "0"
  callback_name        = "${var.name}${random_integer.default.result}"
  callback_types       = ["machineScan"]
  callback_suggestions = ["block"]
}


resource "alicloud_aligreen_oss_stock_task" "default" {
  image_opened                       = true
  auto_freeze_type                   = "acl"
  audio_max_size                     = "200"
  image_scan_limit                   = "1"
  video_frame_interval               = "1"
  video_scan_limit                   = "1000"
  audio_scan_limit                   = "1000"
  video_max_frames                   = "200"
  video_max_size                     = "500"
  start_date                         = "2024-08-01 00:00:00 +0800"
  end_date                           = "2024-12-31 09:06:42 +0800"
  buckets                            = jsonencode([{ "Bucket" : "${alicloud_oss_bucket.defaultPyhXOV.bucket}", "Selected" : true, "Prefixes" : [] }])
  image_scenes                       = ["porn"]
  audio_antispam_freeze_config       = jsonencode({ "Type" : "suggestion", "Value" : "block" })
  image_live_freeze_config           = jsonencode({ "Type" : "suggestion", "Value" : "block" })
  video_terrorism_freeze_config      = jsonencode({ "Type" : "suggestion", "Value" : "block" })
  image_terrorism_freeze_config      = jsonencode({ "Type" : "suggestion", "Value" : "block" })
  callback_id                        = alicloud_aligreen_callback.defaultJnW8Na.id
  image_ad_freeze_config             = jsonencode({ "Type" : "suggestion", "Value" : "block" })
  biz_type                           = "recommend_massmedia_template_01"
  audio_scenes                       = jsonencode(["antispam"])
  image_porn_freeze_config           = jsonencode({ "Type" : "suggestion", "Value" : "block" })
  video_live_freeze_config           = jsonencode({ "Type" : "suggestion", "Value" : "block" })
  video_porn_freeze_config           = jsonencode({ "Type" : "suggestion", "Value" : "block" })
  video_voice_antispam_freeze_config = jsonencode({ "Type" : "suggestion", "Value" : "block" })
  video_scenes                       = jsonencode(["ad", "terrorism", "live", "porn", "antispam"])
  video_ad_freeze_config             = jsonencode({ "Type" : "suggestion", "Value" : "block" })
}
```

### Deleting `alicloud_aligreen_oss_stock_task` or removing it from your configuration

Terraform cannot destroy resource `alicloud_aligreen_oss_stock_task`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `audio_antispam_freeze_config` - (Optional, ForceNew) Voice auto freeze configuration. Example:{"type":"suggestion","value":"block,review"}. The results are frozen according to the suggestion in the speech detection results.
* `audio_auto_freeze_opened` - (Optional, ForceNew) Audio detection auto freeze switch. Value: true: automatically freeze, false: not automatically freeze.
* `audio_max_size` - (Optional, ForceNew, Int) Resource property field representing the maximum size of a single audio. 1~2048MB, the default is 200MB, more than not detected.
* `audio_opened` - (Optional, ForceNew) oss stock scan task detect audio. true: scan audio, false: do not scan audio
* `audio_scan_limit` - (Optional, ForceNew, Int) The upper limit of voice scan in the oss stock scan task. The default value is 1000/Bucket.
* `audio_scenes` - (Optional, ForceNew) The audio detection scenarios included in the oss stock scan task. Set the value to antispam.
* `auto_freeze_type` - (Optional, ForceNew) Automatic freeze type. Value: acl: modify permissions, copy: Move files
* `biz_type` - (Optional, ForceNew) Business scenarios used by the oss stock scan task
* `buckets` - (Optional, ForceNew) The bucket configuration list of the oss stock scan task. Example:[{"Bucket":"bucket_01","Selected":true,"Prefixes":["img/test_"],"Type":"exclude"}]
* `callback_id` - (Optional, ForceNew, Int) The ID of the primary key of the notification message bound to the oss stock scan task.
* `end_date` - (Optional, ForceNew) The end time of the file upload time range indicates the scanning of files uploaded before this time point.
* `image_ad_freeze_config` - (Optional, ForceNew) Picture automatically freezes the configuration of ad scenes. Example: {"type": "suggestion", "value": "block,review"}. The result will be frozen according to the suggestion in the picture detection result.
* `image_auto_freeze_opened` - (Optional, ForceNew) Picture detection auto freeze switch. Value: true: auto freeze, false: not auto freeze.
* `image_live_freeze_config` - (Optional, ForceNew) Picture automatic freezing live scene configuration. Example:{"type":"suggestion","value":"block,review"}. The result will be frozen according to the suggestion in the picture detection result.
* `image_opened` - (Optional, ForceNew) oss stock scan task detect images. true: scan images, false: do not scan images
* `image_porn_freeze_config` - (Optional, ForceNew) Picture automatic freezing porn scene configuration. Example: {"type": "suggestion", "value": "block,review"}. The result will be frozen according to the suggestion in the picture detection result.
* `image_scan_limit` - (Optional, ForceNew, Int) The upper limit for scanning images in the oss stock scan task. The default value is 10000 images per Bucket.
* `image_scenes` - (Optional, ForceNew, List) The image moderation scenario included in the oss stock scan task.Valid values:
porn: pornography detection
terrorism: terrorist content detection
ad: ad violation detection
live: undesirable scene detection
* `image_terrorism_freeze_config` - (Optional, ForceNew) The picture automatically freezes the configuration of terrorism scenes. Example: {"type": "suggestion", "value": "block,review"}. The result will be frozen according to the suggestion in the picture detection result.
* `scan_image_no_file_type` - (Optional, ForceNew) Whether the oss stock scan task detects images with file names without suffixes. true: Detect pictures with file names without suffixes, false: Do not detect pictures with file names without suffixes
* `start_date` - (Optional, ForceNew) The start time of the file upload time range represents the files uploaded after scanning this time point.
* `video_ad_freeze_config` - (Optional, ForceNew) The video automatically freezes the configuration of ad scenarios. Example:{"type":"suggestion","value":"block,review"}. The results will be frozen according to the suggestion in the video detection results.
* `video_auto_freeze_opened` - (Optional, ForceNew) Video detection auto freeze switch. Value: true: automatically freeze, false: not automatically freeze.
* `video_frame_interval` - (Optional, ForceNew, Int) Resource attribute field representing the framing frequency. 1~60 seconds/frame, the default is 1 second/frame
* `video_live_freeze_config` - (Optional, ForceNew) Video automatic freeze live scene configuration. Example:{"type":"suggestion","value":"block,review"}. The results will be frozen according to the suggestion in the video detection results.
* `video_max_frames` - (Optional, ForceNew, Int) A resource attribute field that represents the upper limit of a single video frame cut. 5 to 20000 frames, the default is 200 frames
* `video_max_size` - (Optional, ForceNew, Int) Resource property field representing the maximum size of a single video. 1~2048MB, the default is 500MB, more than not detected.
* `video_opened` - (Optional, ForceNew) oss stock scan task detect video. true: scan video, false: do not scan video
* `video_porn_freeze_config` - (Optional, ForceNew) Video automatic freezing porn scene configuration. Example: {"type": "suggestion", "value": "block,review"}. The result will be frozen according to the suggestion in the video detection result.
* `video_scan_limit` - (Optional, ForceNew, Int) The upper limit of video scanning in the oss stock scan task. The default value is 1000/Bucket.
* `video_scenes` - (Optional, ForceNew) The video detection scenarios included in the oss stock scan task.
porn: pornography detection
terrorism: terrorist content detection
ad: ad violation detection
live: undesirable scene detection
antispam: Video voice antispam
* `video_terrorism_freeze_config` - (Optional, ForceNew) The video automatically freezes the configuration of terrorism scenes. Example:{"type":"suggestion","value":"block,review"}. The results will be frozen according to the suggestion in the video detection results.
* `video_voice_antispam_freeze_config` - (Optional, ForceNew) Voice auto freeze configuration in video. Example:{"type":"suggestion","value":"block,review"}. The results will be frozen according to the suggestion in the video detection results.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Oss Stock Task.

## Import

Aligreen Oss Stock Task can be imported using the id, e.g.

```shell
$ terraform import alicloud_aligreen_oss_stock_task.example <id>
```