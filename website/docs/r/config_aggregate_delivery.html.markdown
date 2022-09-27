---
subcategory: "Cloud Config"
layout: "alicloud"
page_title: "Alicloud: alicloud_config_aggregate_delivery"
sidebar_current: "docs-alicloud-resource-config-aggregate-delivery"
description: |-
  Provides a Alicloud Cloud Config Aggregate Delivery resource.
---

# alicloud\_config\_aggregate\_delivery

Provides a Cloud Config Aggregate Delivery resource.

For information about Cloud Config Aggregate Delivery and how to use it, see [What is Aggregate Delivery](https://www.alibabacloud.com/help/en/cloud-config/latest/delivery-destination-services-overview).

-> **NOTE:** Available in v1.172.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_resource_manager_accounts" "default" {
  status = "CreateSuccess"
}
resource "alicloud_config_aggregator" "default" {
  aggregator_accounts {
    account_id   = data.alicloud_resource_manager_accounts.default.accounts.0.account_id
    account_name = data.alicloud_resource_manager_accounts.default.accounts.0.display_name
    account_type = "ResourceDirectory"
  }
  aggregator_name = var.name
  description     = "tf-create-aggregator"
}
data "alicloud_account" "this" {}
locals {
  uid = data.alicloud_account.this.id
  sls = format("acs:log:%[2]s:%%s:project/%%s/logstore/%%s", local.uid, alicloud_log_project.this.name, alicloud_log_store.this.name)
}
resource "alicloud_log_project" "this" {
  name = "example_value"
}
resource "alicloud_log_store" "this" {
  name    = "example_value"
  project = alicloud_log_project.this.name
}
resource "alicloud_config_delivery" "example" {
  aggregator_id                          = alicloud_config_aggregator.default.id
  configuration_item_change_notification = true
  non_compliant_notification             = true
  delivery_channel_name                  = "example_value"
  delivery_channel_target_arn            = local.sls
  delivery_channel_type                  = "SLS"
  description                            = "example_value"
}
```

## Argument Reference

The following arguments are supported:

* `aggregator_id` - (Required, ForceNew) The ID of the Aggregator.
* `configuration_item_change_notification` - (Optional, Computed) Open or close delivery configuration change history.
* `configuration_snapshot` - (Optional, Computed) Open or close timed snapshot of shipping resources. **NOTE:** The attribute is valid when the attribute `delivery_channel_type` is `OSS`.
* `delivery_channel_condition` - (Optional) The rule attached to the delivery method. Please refer to api [CreateConfigDeliveryChannel](https://help.aliyun.com/document_detail/429798.html) for example format. **NOTE:** The attribute is valid when the attribute `delivery_channel_type` is `MNS`.
* `delivery_channel_name` - (Optional) The name of the delivery method.
* `delivery_channel_target_arn` - (Required) The ARN of the delivery destination. The value must be in one of the following formats:
  * `acs:oss:{RegionId}:{Aliuid}:{bucketName}`: if your delivery destination is an Object Storage Service (OSS) bucket.
  * `acs:mns:{RegionId}:{Aliuid}:/topics/{topicName}`: if your delivery destination is a Message Service (MNS) topic.
  * `acs:log:{RegionId}:{Aliuid}:project/{projectName}/logstore/{logstoreName}`: if your delivery destination is a Log Service Logstore.
* `delivery_channel_type` - (Required, ForceNew) The type of the delivery method. Valid values: `OSS`: Object Storage, `MNS`: Message Service, `SLS`: Log Service.
* `description` - (Optional) The description of the delivery method.
* `non_compliant_notification` - (Optional, Computed) Open or close non-compliance events of delivery resources. **NOTE:** The attribute is valid when the attribute `delivery_channel_type` is `SLS` or `MNS`.
* `oversized_data_oss_target_arn` - (Optional) The oss ARN of the delivery channel when the value data oversized limit.
  * The value must be in one of the following formats: `acs:oss:{RegionId}:{accountId}:{bucketName}`, if your delivery destination is an Object Storage Service (OSS) bucket.
  * Only delivery channels `SLS` and `MNS` are supported. The delivery channel limit for Log Service SLS is 1 MB, and the delivery channel limit for Message Service MNS is 64 KB.
* `status` - (Optional, Computed) The status of the delivery method. Valid values: `0`: The delivery method is disabled. `1`: The delivery destination is enabled. This is the default value.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Aggregate Delivery. The value formats as `<aggregator_id>:<delivery_channel_id>`.
* `delivery_channel_id` - The ID of the delivery method.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Aggregate Delivery.
* `delete` - (Defaults to 1 mins) Used when delete the Aggregate Delivery.
* `update` - (Defaults to 1 mins) Used when update the Aggregate Delivery.

## Import

Cloud Config Aggregate Delivery can be imported using the id, e.g.

```
$ terraform import alicloud_config_aggregate_delivery.example <aggregator_id>:<delivery_channel_id>
```