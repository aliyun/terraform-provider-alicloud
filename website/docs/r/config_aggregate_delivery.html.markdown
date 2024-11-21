---
subcategory: "Cloud Config (Config)"
layout: "alicloud"
page_title: "Alicloud: alicloud_config_aggregate_delivery"
description: |-
  Provides a Alicloud Config Aggregate Delivery resource.
---

# alicloud_config_aggregate_delivery

Provides a Config Aggregate Delivery resource.

Delivery channel of aggregator.

For information about Config Aggregate Delivery and how to use it, see [What is Aggregate Delivery](https://www.alibabacloud.com/help/en/cloud-config/latest/api-config-2020-09-07-createaggregateconfigdeliverychannel).

-> **NOTE:** Available since v1.172.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_config_aggregate_delivery&exampleId=0ca7d69d-f0e3-cbed-ceb6-bcac284c0108ed852513&activeTab=example&spm=docs.r.config_aggregate_delivery.0.0ca7d69df0&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf_example"
}
data "alicloud_regions" "this" {
  current = true
}
data "alicloud_account" "this" {}
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
  description     = var.name
  aggregator_type = "CUSTOM"
}

resource "random_uuid" "default" {}
resource "alicloud_log_project" "default" {
  project_name = substr("tf-example-${replace(random_uuid.default.result, "-", "")}", 0, 16)
}
resource "alicloud_log_store" "default" {
  logstore_name = var.name
  project_name  = alicloud_log_project.default.project_name
}
resource "alicloud_config_aggregate_delivery" "default" {
  aggregator_id                          = alicloud_config_aggregator.default.id
  configuration_item_change_notification = true
  non_compliant_notification             = true
  delivery_channel_name                  = var.name
  delivery_channel_target_arn            = "acs:log:${data.alicloud_regions.this.ids.0}:${data.alicloud_account.this.id}:project/${alicloud_log_project.default.project_name}/logstore/${alicloud_log_store.default.logstore_name}"
  delivery_channel_type                  = "SLS"
  description                            = var.name
}
```

## Argument Reference

The following arguments are supported:
* `aggregator_id` - (Required, ForceNew) Aggregator ID.
* `configuration_item_change_notification` - (Optional) Indicates whether the specified destination receives resource change logs. If the value of this parameter is true, Cloud Config delivers the resource change logs to OSS, Log Service, or MNS when the configurations of the resources change. Valid values:  
  - true: The specified destination receives resource change logs.  
  - false: The specified destination does not receive resource change logs.  
* `configuration_snapshot` - (Optional) Indicates whether the specified destination receives scheduled resource snapshots. Cloud Config delivers scheduled resource snapshots at 04:00Z and 16:00Z to OSS, MNS, or Log Service every day. The time is displayed in UTC. Valid values:  
  - true: The specified destination receives scheduled resource snapshots.  
  - false: The specified destination does not receive scheduled resource snapshots.  
* `delivery_channel_condition` - (Optional) The rule that is attached to the delivery channel.   

  This parameter is available when you deliver data of all types to MNS or deliver snapshots to Log Service.

  If you specify the risk level or resource types for subscription events, this is as follows:  

  The lowest risk level of the events to which you want to subscribe is in the following format: {"filterType":"RuleRiskLevel","value":"1","multiple":false}, The value field indicates the risk level of the events to which you want to subscribe. Valid values: 1, 2, and 3. The value 1 indicates the high risk level, the value 2 indicates the medium risk level, and the value 3 indicates the low risk level.  

  The setting of the resource types of the events to which you want to subscribe is in the following format: {"filterType":"ResourceType","values":["ACS::ACK::Cluster","ACS::ActionTrail::Trail","ACS::CBWP::CommonBandwidthPackage"],"multiple":true}, The values field indicates the resource types of the events to which you want to subscribe. The value of the field is a JSON array. 

  Examples:[{"filterType":"ResourceType","values":["ACS::ActionTrail::Trail","ACS::CBWP::CommonBandwidthPackage","ACS::CDN::Domain","ACS::CEN::CenBandwidthPackage","ACS::CEN::CenInstance","ACS::CEN::Flowlog","ACS::DdosCoo::Instance"],"multiple":true}].   
* `delivery_channel_name` - (Optional) The name of the delivery channel.
* `delivery_channel_target_arn` - (Required) The ARN of the delivery destination.  
  - If the value of the DeliveryChannelType parameter is OSS, the value of this parameter is the ARN of the destination OSS bucket.  
  - If the value of the DeliveryChannelType parameter is MNS, the value of this parameter is the ARN of the destination MNS topic.  
  - If the value of the DeliveryChannelType parameter is SLS, the value of this parameter is the ARN of the destination Log Service Logstore.  
* `delivery_channel_type` - (Required, ForceNew) The type of the delivery channel. Valid values:
  - OSS: Object Storage Service (OSS)
  - MNS: Message Service (MNS)
  - SLS: Log Service
* `description` - (Optional) The description of the delivery method.
* `non_compliant_notification` - (Optional) Indicates whether the specified destination receives resource non-compliance events. If the value of this parameter is true, Cloud Config delivers resource non-compliance events to Log Service or MNS when resources are evaluated as non-compliant. Valid values:  
  - true: The specified destination receives resource non-compliance events.  
  - false: The specified destination does not receive resource non-compliance events.  
* `oversized_data_oss_target_arn` - (Optional) The ARN of the OSS bucket to which the delivered data is transferred when the size of the data exceeds the specified upper limit of the delivery channel.
* `status` - (Optional, Computed) The status of the delivery method. Valid values:   
  - 0: The delivery method is disabled.   
  - 1: The delivery destination is enabled. This is the default value.  

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<aggregator_id>:<delivery_channel_id>`.
* `delivery_channel_id` - The ID of the delivery method. This parameter is required when you modify a delivery method.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Aggregate Delivery.
* `delete` - (Defaults to 5 mins) Used when delete the Aggregate Delivery.
* `update` - (Defaults to 5 mins) Used when update the Aggregate Delivery.

## Import

Config Aggregate Delivery can be imported using the id, e.g.

```shell
$ terraform import alicloud_config_aggregate_delivery.example <aggregator_id>:<delivery_channel_id>
```