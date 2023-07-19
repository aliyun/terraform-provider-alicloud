---
subcategory: "Cloud Config (Config)"
layout: "alicloud"
page_title: "Alicloud: alicloud_config_delivery"
sidebar_current: "docs-alicloud-resource-config-delivery"
description: |-
  Provides a Alicloud Cloud Config Delivery resource.
---

# alicloud_config_delivery

Provides a Cloud Config Delivery resource.

For information about Cloud Config Delivery and how to use it, see [What is Delivery](https://www.alibabacloud.com/help/en/cloud-config/latest/api-doc-config-2020-09-07-api-doc-createconfigdeliverychannel).

-> **NOTE:** Available since v1.171.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example-sls"
}
data "alicloud_account" "this" {}
data "alicloud_regions" "this" {
  current = true
}

resource "alicloud_log_project" "default" {
  name = var.name
}
resource "alicloud_log_store" "default" {
  name    = var.name
  project = alicloud_log_project.default.name
}

resource "alicloud_config_delivery" "default" {
  configuration_item_change_notification = true
  non_compliant_notification             = true
  delivery_channel_name                  = var.name
  delivery_channel_target_arn            = "acs:log:${data.alicloud_regions.this.ids.0}:${data.alicloud_account.this.id}:project/${alicloud_log_project.default.name}/logstore/${alicloud_log_store.default.name}"
  delivery_channel_type                  = "SLS"
  description                            = var.name
}
```

## Argument Reference

The following arguments are supported:

* `configuration_item_change_notification` - (Optional) Open or close delivery configuration change history. true: open, false: close.
* `configuration_snapshot` - (Optional) Open or close timed snapshot of shipping resources. **NOTE:** The attribute is valid when the attribute `delivery_channel_type` is `OSS`.
* `delivery_channel_condition` - (Optional) The rule attached to the delivery method. Please refer to api [CreateConfigDeliveryChannel](https://help.aliyun.com/document_detail/429798.html) for example format. **NOTE:** The attribute is valid when the attribute `delivery_channel_type` is `MNS`.
* `delivery_channel_name` - (Optional) The name of the delivery method.
* `delivery_channel_target_arn` - (Required) The ARN of the delivery destination. The value must be in one of the following formats:
  * `acs:oss:{RegionId}:{Aliuid}:{bucketName}`: if your delivery destination is an Object Storage Service (OSS) bucket. 
  * `acs:mns:{RegionId}:{Aliuid}:/topics/{topicName}`: if your delivery destination is a Message Service (MNS) topic. 
  * `acs:log:{RegionId}:{Aliuid}:project/{projectName}/logstore/{logstoreName}`: if your delivery destination is a Log Service Logstore.
* `delivery_channel_type` - (Required, ForceNew) The type of the delivery method. Valid values: `OSS`: Object Storage, `MNS`: Message Service, `SLS`: Log Service.
* `description` - (Optional) The description of the delivery method.
* `non_compliant_notification` - (Optional) Open or close non-compliance events of delivery resources. **NOTE:** The attribute is valid when the attribute `delivery_channel_type` is `SLS` or `MNS`.
* `oversized_data_oss_target_arn` - (Optional) The oss ARN of the delivery channel when the value data oversized limit. 
  * The value must be in one of the following formats: `acs:oss:{RegionId}:{accountId}:{bucketName}`, if your delivery destination is an Object Storage Service (OSS) bucket.
  * Only delivery channels `SLS` and `MNS` are supported. The delivery channel limit for Log Service SLS is 1 MB, and the delivery channel limit for Message Service MNS is 64 KB.
* `status` - (Optional) The status of the delivery method. Valid values: `0`: The delivery method is disabled. `1`: The delivery destination is enabled. This is the default value.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Delivery.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Config Delivery Channel.
* `update` - (Defaults to 1 mins) Used when update the Config Delivery Channel.
* `delete` - (Defaults to 1 mins) Used when update the Config Delivery Channel.

## Import

Cloud Config Delivery can be imported using the id, e.g.

```shell
$ terraform import alicloud_config_delivery.example <id>
```