---
subcategory: "Data Transmission Service (DTS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_dts_consumer_channels"
sidebar_current: "docs-alicloud-datasource-dts-consumer-channels"
description: |-
  Provides a list of Dts Consumer Channels to the user.
---

# alicloud\_dts\_consumer\_channels

This data source provides the Dts Consumer Channels of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.146.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_dts_consumer_channels" "ids" {}
output "dts_consumer_channel_id_1" {
  value = data.alicloud_dts_consumer_channels.ids.channels.0.id
}
```

## Argument Reference

The following arguments are supported:

* `dts_instance_id` - (Required, ForceNew) Subscription instance ID.
* `ids` - (Optional, ForceNew, Computed)  A list of Consumer Channel IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `channels` - A list of Dts Consumer Channels. Each element contains the following attributes:
	* `consumer_group_id` - The ID of the consumer group.
	* `consumer_group_name` - The name of the consumer group.
	* `consumer_group_user_name` - The username of the consumer group.
	* `consumption_checkpoint` - The time point when the client consumed the last message in the subscription channel.
	* `id` - The ID of the Consumer Channel.
	* `message_delay` - The message delay time, for the current time data subscription channel in the earliest time of unconsumed messages of the difference, in Unix timestamp format, which is measured in seconds.
	* `unconsumed_data` - The total number of unconsumed messages.