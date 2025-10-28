---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_multi_account_delivery_channel"
description: |-
  Provides a Alicloud Resource Manager Multi Account Delivery Channel resource.
---

# alicloud_resource_manager_multi_account_delivery_channel

Provides a Resource Manager Multi Account Delivery Channel resource.

Multi-account Resource Delivery Channel.

For information about Resource Manager Multi Account Delivery Channel and how to use it, see [What is Multi Account Delivery Channel](https://next.api.alibabacloud.com/document/ResourceCenter/2022-12-01/CreateMultiAccountDeliveryChannel).

-> **NOTE:** Available since v1.262.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_resource_manager_folder" "defaultuHQ8Cu" {
  folder_name = "folder-aone-example-1"
}

resource "alicloud_resource_manager_folder" "defaultioI16p" {
  folder_name = "folder-aone-example-2"
}

resource "alicloud_resource_manager_folder" "default55Uum4" {
  folder_name = "folder-aone-example-3"
}

resource "alicloud_resource_manager_folder" "defaultiEjEbe" {
  folder_name = "folder-aone-example-4"
}

resource "alicloud_resource_manager_folder" "defaultdNL2TN" {
  folder_name = "folder-aone-example-5"
}


resource "alicloud_resource_manager_multi_account_delivery_channel" "default" {
  resource_change_delivery {
    sls_properties {
      oversized_data_oss_target_arn = "acs:oss:cn-hangzhou:1511928242963727:resourcecenter-aone-example-delivery-oss"
    }
    target_arn = "acs:log:cn-hangzhou:1511928242963727:project/delivery-aone-example/logstore/resourcecenter-delivery-aone-example-sls"
  }
  delivery_channel_description        = "multi_delivery_channel_resource_spec_mq_example"
  multi_account_delivery_channel_name = "multi_delivery_channel_resource_spec_mq_example"
  delivery_channel_filter {
    account_scopes = ["${alicloud_resource_manager_folder.defaultuHQ8Cu.id}", "${alicloud_resource_manager_folder.defaultioI16p.id}", "${alicloud_resource_manager_folder.default55Uum4.id}"]
    resource_types = ["ACS::ACK::Cluster", "ACS::ActionTrail::Trail", "ACS::BPStudio::Application"]
  }
  resource_snapshot_delivery {
    delivery_time     = "16:00Z"
    target_arn        = "acs:log:cn-hangzhou:1511928242963727:project/delivery-aone-example/logstore/resourcecenter-delivery-aone-example-sls"
    target_type       = "SLS"
    custom_expression = "select * from resources limit 10;"
    sls_properties {
      oversized_data_oss_target_arn = "acs:oss:cn-hangzhou:1511928242963727:resourcecenter-aone-example-delivery-oss"
    }
  }
}
```

## Argument Reference

The following arguments are supported:
* `delivery_channel_description` - (Required) The description of the delivery channel.
* `delivery_channel_filter` - (Required, List) The effective scope of the delivery channel. See [`delivery_channel_filter`](#delivery_channel_filter) below.
* `multi_account_delivery_channel_name` - (Required) The name of the delivery channel.
* `resource_change_delivery` - (Optional, List) The configurations for delivery of resource configuration change events. See [`resource_change_delivery`](#resource_change_delivery) below.
* `resource_snapshot_delivery` - (Optional, List) The configurations for delivery of scheduled resource snapshots. See [`resource_snapshot_delivery`](#resource_snapshot_delivery) below.

### `delivery_channel_filter`

The delivery_channel_filter supports the following:
* `account_scopes` - (Required, List) The account scopes of the delivery channel.
* `resource_types` - (Optional, List) An array of effective resource types for the delivery channel.
  - Example: ["ACS::VPC::VPC", "ACS::ECS::Instance"].
  - If you want to deliver items of all resource types supported by Resource Center, set this parameter to ["ALL"].

### `resource_change_delivery`

The resource_change_delivery supports the following:
* `enabled` - (Optional, Computed) Specifies whether to enable delivery of resource configuration change events. Valid values:
  - true
  - false
* `sls_properties` - (Optional, List) The Simple Log Service configurations. See [`sls_properties`](#resource_change_delivery-sls_properties) below.
* `target_arn` - (Optional) The ARN of the delivery destination.
  - If you set TargetType to`OSS`, you must set TargetArn to the ARN of a bucket whose name is prefixed with `resourcecenter-`.
  - If you set TargetType to`SLS`, you must set TargetArn to the ARN of a Logstore whose name is prefixed with `resourcecenter-`.
* `target_type` - (Optional, ForceNew, Computed) The type of the delivery destination.

Valid values:
  - SLS

### `resource_change_delivery-sls_properties`

The resource_change_delivery-sls_properties supports the following:
* `oversized_data_oss_target_arn` - (Optional) The ARN of the destination to which large files are delivered.
  - If the size of a resource configuration change event exceeds 1 MB, the event is delivered as an OSS object.
  - You need to set this parameter to the ARN of a bucket whose name is prefixed with resourcecenter-.

### `resource_snapshot_delivery`

The resource_snapshot_delivery supports the following:
* `custom_expression` - (Optional) The custom expression.
* `delivery_time` - (Optional) The delivery time.
* `enabled` - (Optional, Computed) Specifies whether to enable delivery of scheduled resource snapshots. Valid values:
  - true
  - false
* `sls_properties` - (Optional, List) The Simple Log Service configurations. See [`sls_properties`](#resource_snapshot_delivery-sls_properties) below.
* `target_arn` - (Optional) - The Alibaba Cloud Resource Name (ARN) of the delivery destination.
  - If you set TargetType to`OSS`, you must set TargetArn to the ARN of a bucket whose name is prefixed with `resourcecenter-`.
  - If you set TargetType to `SLS`, you must set TargetArn to the ARN of a Logstore whose name is prefixed with `resourcecenter-`.
* `target_type` - (Optional) The type of the delivery destination.

Valid values:
  - `OSS` for standard delivery
  - `OSS` or `SLS` for custom delivery

### `resource_snapshot_delivery-sls_properties`

The resource_snapshot_delivery-sls_properties supports the following:
* `oversized_data_oss_target_arn` - (Optional) exceeds 1 MB, the event is delivered as an OSS object.
  - You need to set this parameter to the ARN of a bucket whose name is prefixed with resourcecenter-.
  - This parameter takes effect only if you use custom delivery for scheduled resource snapshots. You do not need to configure this parameter if you use standard delivery for scheduled resource snapshots.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Multi Account Delivery Channel.
* `delete` - (Defaults to 5 mins) Used when delete the Multi Account Delivery Channel.
* `update` - (Defaults to 5 mins) Used when update the Multi Account Delivery Channel.

## Import

Resource Manager Multi Account Delivery Channel can be imported using the id, e.g.

```shell
$ terraform import alicloud_resource_manager_multi_account_delivery_channel.example <id>
```