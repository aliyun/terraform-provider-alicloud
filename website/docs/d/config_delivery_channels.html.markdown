---
subcategory: "Cloud Config"
layout: "alicloud"
page_title: "Alicloud: alicloud_config_delivery_channels"
sidebar_current: "docs-alicloud-datasource-config-delivery-channels"
description: |-
    Provides a list of Config Delivery Channels to the user.
---

# alicloud\_config\_delivery\_channels

-> **DEPRECATED:**  This resource is based on Config's old version OpenAPI, and it has been deprecated from version `1.173.0`.
Please use new datasource [alicloud_config_deliveries](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/data-sources/config_deliveries) instead.

This data source provides the Config Delivery Channels of the current Alibaba Cloud user.

-> **NOTE:**  Available in 1.99.0+.

-> **NOTE:** The Cloud Config region only support `cn-shanghai` and `ap-southeast-1`.

## Example Usage

```terraform
data "alicloud_config_delivery_channels" "example" {
  ids        = ["cdc-49a2ad756057********"]
  name_regex = "tftest"
}

output "first_config_delivery_channel_id" {
  value = data.alicloud_config_delivery_channels.example.channels.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew) A list of Config Delivery Channel IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by delivery channel name.
* `status` - (Optional, ForceNew) The status of the config delivery channel. Valid values `0`: Disable delivery channel, `1`: Enable delivery channel.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Config Delivery Channel IDs.
* `names` - A list of Config Delivery Channel names.
* `channels` - A list of Config Delivery Channels. Each element contains the following attributes:
    * `id` - The ID of the Config Delivery Channel.
    * `delivery_channel_assume_role_arn` - The Alibaba Cloud Resource Name (ARN) of the role assumed by delivery method.
    * `delivery_channel_condition` - The rule attached to the delivery method. This parameter is applicable only to delivery methods of the Message Service (MNS) type.
    * `delivery_channel_id` - The ID of the delivery channel.
    * `delivery_channel_name` - The name of the delivery channel.
    * `delivery_channel_target_arn` - The ARN of the delivery destination.
    * `delivery_channel_type` - The type of the delivery channel.
    * `description` - The description of the delivery channel.
    * `status` - The status of the delivery channel.
