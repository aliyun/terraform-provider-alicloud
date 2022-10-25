---
subcategory: "Cloud Config"
layout: "alicloud"
page_title: "Alicloud: alicloud_config_delivery_channel"
sidebar_current: "docs-alicloud-resource-config-delivery-channel"
description: |-
  Provides a Alicloud Config Delivery Channel resource.
---

# alicloud\_config\_delivery\_channel

-> **DEPRECATED:**  This resource is based on Config's old version OpenAPI, and it has been deprecated from version `1.171.0`.
Please use new resource [alicloud_config_delivery](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/config_delivery) instead.

Provides an Alicloud Config Delivery Channel resource. You can receive configuration audit event changes by configuring OSS, MNS and SLS services provided by Alibaba Cloud.
For information about Alicloud Config Delivery Channel and how to use it, see [What is Delivery Channel](https://www.alibabacloud.com/help/en/doc-detail/307022.html).

-> **NOTE:** Available in v1.99.0+.

-> **NOTE:** The Cloud Config region only support `cn-shanghai` and `ap-southeast-1`.

-> **NOTE:** Once each type of delivery channel is created, it does not support destroyed by terraform. Only support through the `status` attribute control enable and disable.

## Example Usage

```terraform
provider "alicloud" {
  region = "cn-shanghai"
}

resource "alicloud_mns_topic" "example" {
  name = "test-topic"
}

# Example for create a MNS delivery channel
resource "alicloud_config_delivery_channel" "example" {
  description                      = "channel_description"
  delivery_channel_name            = "channel_name"
  delivery_channel_assume_role_arn = "acs:ram::11827252********:role/aliyunserviceroleforconfig"
  delivery_channel_type            = "MNS"
  delivery_channel_target_arn      = format("acs:oss:cn-shanghai:11827252********:/topics/%s", alicloud_mns_topic.example.name)
  delivery_channel_condition       = <<EOF
  [
      {
          "filterType":"ResourceType",
          "values":[
              "ACS::CEN::CenInstance",
              "ACS::CEN::Flowlog",
          ],
          "multiple":true
      }
  ]
    EOF
}
```
## Argument Reference

The following arguments are supported:

* `delivery_channel_name` - (Optional, Computed) The name of the delivery channel.
* `description` - (Optional, Computed) The description of the delivery method.
* `status` - (Optional, Computed) The status of the delivery method. Valid values: `0`: The delivery method is disabled., `1`: The delivery destination is enabled. This is the default value. 
* `delivery_channel_assume_role_arn` - (Required) The Alibaba Cloud Resource Name (ARN) of the role to be assumed by the delivery method.
* `delivery_channel_type` - (Required, ForceNew) - The type of the delivery method. This parameter is required when you create a delivery method. Valid values: `OSS`: Object Storage, `MNS`: Message Service, `SLS`: Log Service.
* `delivery_channel_target_arn` - (Required) - The ARN of the delivery destination. This parameter is required when you create a delivery method. The value must be in one of the following formats:
    - `acs:oss:{RegionId}:{Aliuid}:{bucketName}`: if your delivery destination is an Object Storage Service (OSS) bucket. 
    - `acs:mns:{RegionId}:{Aliuid}:/topics/{topicName}`: if your delivery destination is a Message Service (MNS) topic.
    - `acs:log:{RegionId}:{Aliuid}:project/{projectName}/logstore/{logstoreName}`: if your delivery destination is a Log Service Logstore.
* `delivery_channel_condition` - (Optional, Computed) The rule attached to the delivery method. This parameter is applicable only to delivery methods of the MNS type. Please refer to api [PutDeliveryChannel](https://www.alibabacloud.com/help/en/doc-detail/174253.htm) for example format. 

## Attributes Reference

The following attributes are exported:

* `id` - This ID of Config Delivery Channel.  

### Timeouts

-> **NOTE:** Available in 1.104.0+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 3 mins) Used when create the Config Delivery Channel.
* `update` - (Defaults to 3 mins) Used when update the Config Delivery Channel.

## Import

Alicloud Config Delivery Channel can be imported using the id, e.g.

```
$ terraform import alicloud_config_delivery_channel.example cdc-49a2ad756057********
```
