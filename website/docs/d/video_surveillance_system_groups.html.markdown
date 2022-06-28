---
subcategory: "Video Surveillance System"
layout: "alicloud"
page_title: "Alicloud: alicloud_video_surveillance_system_groups"
sidebar_current: "docs-alicloud-datasource-video-surveillance-system-groups"
description: |-
  Provides a list of Video Surveillance System Groups to the user.
---

# alicloud\_video\_surveillance\_system\_groups

This data source provides the Video Surveillance System Groups of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.135.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_video_surveillance_system_group" "default" {
  group_name   = "groupname"
  in_protocol  = "rtmp"
  out_protocol = "flv"
  play_domain  = "your_plan_domain"
  push_domain  = "your_push_domain"
}
data "alicloud_video_surveillance_system_groups" "default" {
  ids = [alicloud_video_surveillance_system_group.default.id]
}
output "vs_group" {
  value = data.alicloud_video_surveillance_system_groups.default.ids.0
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Group IDs.
* `in_protocol` - (Optional, ForceNew) The use of the access protocol support gb28181, Real Time Messaging Protocol (rtmp). Valid values: `gb28181`, `rtmp`.
* `name` - (Optional, ForceNew) The name.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Group name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status. Valid values: `on`,`off`.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Group names.
* `groups` - A list of Video Surveillance System Groups. Each element contains the following attributes:
	* `app` - The App Name of Group.
	* `callback` - The space within the device status update of the callback, need to start with http:// or https:// at the beginning.
	* `create_time` - The creation time of the Group.
	* `description` - The description of the Group.
	* `enabled` - Whether to open Group.
	* `gb_id` - Space of national standard ID. **NOTE:** Available only in the national standard access space.
	* `gb_ip` - Space of national standard signaling server address. **NOTE:** Available only in the national standard access space.
	* `group_id` - The ID of Group.
	* `group_name` - The name of Group.
	* `id` - The ID of the Group.
	* `in_protocol` - The use of the access protocol support `gb28181`,`rtmp`(Real Time Messaging Protocol). 
	* `out_protocol` - The use of space play Protocol multi-valued separate them with commas (,). Valid values: `flv`,`hls`, `rtmp`(Real Time Messaging Protocol).
	* `play_domain` -The domain name of plan streaming used by the group.
	* `push_domain` - The domain name of push streaming used by the group.
	* `stats` - The Device statistics of Group.
		* `device_num` - The total number of devices in the group.
		* `ied_num` - The total number of smart devices in the group.
		* `ipc_num` - The total number of cameras in the group.
		* `platform_num` - The total number of platforms in the group.
