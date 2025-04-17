---
subcategory: "RocketMQ"
layout: "alicloud"
page_title: "Alicloud: alicloud_rocketmq_instance"
description: |-
  Provides a Alicloud RocketMQ Instance resource.
---

# alicloud_rocketmq_instance

Provides a RocketMQ Instance resource.


For information about RocketMQ Instance and how to use it, see [What is Instance](https://www.alibabacloud.com/help/en/apsaramq-for-rocketmq/cloud-message-queue-rocketmq-5-x-series/developer-reference/api-rocketmq-2022-08-01-createinstance).

-> **NOTE:** Available since v1.212.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_rocketmq_instance&exampleId=946de15e-604b-d7bb-90fe-117d3bfc8a676030e55e&activeTab=example&spm=docs.r.rocketmq_instance.0.946de15e60&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "createVPC" {
  description = "example"
  cidr_block  = "172.16.0.0/12"
  vpc_name    = var.name
}

resource "alicloud_vswitch" "createVSwitch" {
  description  = "example"
  vpc_id       = alicloud_vpc.createVPC.id
  cidr_block   = "172.16.0.0/24"
  vswitch_name = var.name
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_rocketmq_instance" "default" {
  product_info {
    msg_process_spec       = "rmq.u2.10xlarge"
    send_receive_ratio     = "0.3"
    message_retention_time = "70"
  }
  service_code      = "rmq"
  payment_type      = "PayAsYouGo"
  instance_name     = var.name
  sub_series_code   = "cluster_ha"
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
  remark            = "example"
  ip_whitelist      = ["192.168.0.0/16", "10.10.0.0/16", "172.168.0.0/16"]
  software {
    maintain_time = "02:00-06:00"
  }
  tags = {
    Created = "TF"
    For     = "example"
  }
  series_code = "ultimate"
  network_info {
    vpc_info {
      vpc_id = alicloud_vpc.createVPC.id
      vswitches {
        vswitch_id = alicloud_vswitch.createVSwitch.id
      }
    }
    internet_info {
      internet_spec      = "enable"
      flow_out_type      = "payByBandwidth"
      flow_out_bandwidth = "30"
    }
  }
}
```

### Deleting `alicloud_rocketmq_instance` or removing it from your configuration

The `alicloud_rocketmq_instance` resource allows you to manage  `payment_type = "Subscription"`  instance, but Terraform cannot destroy it.
Deleting the subscription resource or removing it from your configuration will remove it from your state file and management, but will not destroy the Instance.
You can resume managing the subscription instance via the AlibabaCloud Console.

## Argument Reference

The following arguments are supported:
* `acl_info` - (Optional, Set, Available since v1.245.0) The access control list for the instance. See [`acl_info`](#acl_info) below.
* `auto_renew` - (Optional, ForceNew) Whether to enable auto-renewal. This parameter is only applicable when the payment type for the instance is Subscription (prepaid).
  - true: Enable auto-renewal
  - false: Disable auto-renewal
* `auto_renew_period` - (Optional, ForceNew, Int) Auto-renewal period. This parameter is only valid when auto-renewal is enabled. Unit: months.

  The values can be as follows:
  - Monthly renewal: 1, 2, 3, 6, 12
* `auto_renew_period_unit` - (Optional) The minimum periodic unit for the duration of auto-renewal. This parameter is only valid when auto-renewal is enabled. Valid values: `Month`, `Year`.
* `commodity_code` - (Optional, ForceNew, Available since v1.231.0) Commodity code

  ons_rmqsub_public_cn: Package year and month instance

  ons_rmqpost_public_cn: Pay-As-You-Go instance

  Next: Serverless instances
* `instance_name` - (Optional) The name of instance
* `ip_whitelists` - (Optional, List, Available since v1.245.0) The ip whitelist.
* `network_info` - (Required, Set) Instance network configuration information See [`network_info`](#network_info) below.
* `payment_type` - (Required, ForceNew) The payment type for the instance. Alibaba Cloud Message Queue RocketMQ version supports two types of payment:

  The parameter values are as follows:
  - PayAsYouGo: Pay-as-you-go, a post-payment model where you pay after usage.
  - Subscription: Subscription-based, a pre-payment model where you pay before usage. 

  For more information, please refer to [Billing Methods](https://help.aliyun.com/zh/apsaramq-for-rocketmq/cloud-message-queue-rocketmq-5-x-series/product-overview/overview-2).
* `period` - (Optional, Int) Duration of purchase. This parameter is only valid when the payment type for the instance is Subscription (prepaid).

  The values can be as follows:
  - Monthly purchase: 1, 2, 3, 4, 5, 6
  - Annual purchase: 1, 2, 3
* `period_unit` - (Optional) The minimum periodic unit for the duration of purchase.

  The parameter values are as follows:
  - Month: Purchase on a monthly basis
  - Year: Purchase on an annual basis
* `product_info` - (Optional, ForceNew, List) product info See [`product_info`](#product_info) below.
* `remark` - (Optional) Custom description
* `resource_group_id` - (Optional) The ID of the resource group
* `series_code` - (Required, ForceNew) The primary series encoding for the instance. For specific differences between the primary series, please refer to [Product Selection](https://help.aliyun.com/zh/apsaramq-for-rocketmq/cloud-message-queue-rocketmq-5-x-series/product-overview/instance-selection).

  The parameter values are as follows:
  - standard: Standard Edition
  - ultimate: Platinum Edition
  - professional: Professional Edition
* `service_code` - (Required, ForceNew) The code of the service code instance. The code of the RocketMQ is rmq.
* `software` - (Optional, List) Instance software information. See [`software`](#software) below.
* `sub_series_code` - (Required, ForceNew) The sub-series encoding for the instance. For specific differences between the sub-series, please refer to [Product Selection](https://help.aliyun.com/zh/apsaramq-for-rocketmq/cloud-message-queue-rocketmq-5-x-series/product-overview/instance-selection).

  The parameter values are as follows:
  - cluster_ha: Cluster High Availability Edition
  - single_node: Single Node Testing Edition
  - serverless：Serverless instance
 **NOTE:** From version 1.245.0, `sub_series_code` can be set to `serverless`.
  When selecting the primary series as ultimate (Platinum Edition), the sub-series can only be chosen as cluster_ha (Cluster High Availability Edition).
* `tags` - (Optional, Map) The resource label.

### `acl_info`

The acl_info supports the following:
* `acl_types` - (Optional, List) The authentication type of the instance. Valid values:
  - `default`: Intelligent identity authentication.
  - `apache_acl`: Access control list (ACL) identity authentication.
* `default_vpc_auth_free` - (Optional, Bool) Indicates whether the authentication-free in VPCs feature is enabled. Indicates whether the authentication-free in VPCs feature is enabled. Valid values:
  - `true`: Enable secret-free access.
  - `false`: Turn off secret-free access.

### `network_info`

The network_info supports the following:
* `internet_info` - (Required, ForceNew, Set) instance internet info. See [`internet_info`](#network_info-internet_info) below.
* `vpc_info` - (Required, ForceNew, Set) Proprietary network information. See [`vpc_info`](#network_info-vpc_info) below.

### `network_info-internet_info`

The network_info-internet_info supports the following:
* `flow_out_bandwidth` - (Optional, ForceNew, Int) Public network bandwidth specification. Unit: Mb/s.  This field should only be filled when the public network billing type is set to payByBandwidth.  The value range is [1 - 1000].
* `flow_out_type` - (Required, ForceNew) Public network billing type.  Parameter values are as follows:
  - payByBandwidth: Fixed bandwidth billing. This parameter must be set to the value when public network access is enabled.
  - uninvolved: Not involved. This parameter must be set to the value when public network access is disabled.
* `internet_spec` - (Required, ForceNew) Whether to enable public network access.  The parameter values are as follows:
  - enable: Enable public network access
  - disable: Disable public network access   Instances by default support VPC access. If public network access is enabled, Alibaba Cloud Message Queue RocketMQ version will incur charges for public network outbound bandwidth. For specific billing information, please refer to [Public Network Access Fees](https://help.aliyun.com/zh/apsaramq-for-rocketmq/cloud-message-queue-rocketmq-5-x-series/product-overview/internet-access-fee).
* `ip_whitelist` - (Optional, List, Deprecated since v1.245.0) Field `ip_whitelist` has been deprecated from provider version 1.245.0. New field `ip_whitelists` instead.

### `network_info-vpc_info`

The network_info-vpc_info supports the following:
* `security_group_ids` - (Optional, ForceNew, Available since v1.231.0) Security group id.
* `vswitch_id` - (Optional, ForceNew, Deprecated since v1.231.0) Field `vswitch_id` has been deprecated from provider version 1.245.0. New field `vswitches` instead.
* `vswitches` - (Optional, ForceNew, Set, Available since v1.231.0) Multiple VSwitches. At least two VSwitches are required for a serverless instance. See [`vswitches`](#network_info-vpc_info-vswitches) below.
* `vpc_id` - (Required, ForceNew) Proprietary Network.

### `network_info-vpc_info-vswitches`

The network_info-vpc_info-vswitches supports the following:
* `vswitch_id` - (Optional, ForceNew) VPC switch id.

### `product_info`

The product_info supports the following:
* `auto_scaling` - (Optional) is open auto scaling.
* `message_retention_time` - (Optional, Int) Duration of message retention. Unit: hours.  For the range of values, please refer to [Usage Limits](https://help.aliyun.com/zh/apsaramq-for-rocketmq/cloud-message-queue-rocketmq-5-x-series/product-overview/usage-limits)>Resource Quotas>Limitations on Message Retention.  The message storage in AlibabaCloud RocketMQ is fully implemented in a serverless and elastic manner, with charges based on the actual storage space. You can control the storage capacity of messages by adjusting the duration of message retention. For more information, please see [Storage Fees](https://help.aliyun.com/zh/apsaramq-for-rocketmq/cloud-message-queue-rocketmq-5-x-series/product-overview/storage-fees).
* `msg_process_spec` - (Required, ForceNew) Message sending and receiving calculation specifications. For details about the upper limit for sending and receiving messages, see [Instance Specifications](https://help.aliyun.com/zh/apsaramq-for-rocketmq/cloud-message-queue-rocketmq-5-x-series/product-overview/instance-specifications).
* `send_receive_ratio` - (Optional, Float) message send receive ratio.  Value range: [0.2, 0.5].
* `storage_encryption` - (Optional, ForceNew, Bool, Available since v1.245.0) Specifies whether to enable the encryption at rest feature. Valid values: `true`, `false`.
* `storage_secret_key` - (Optional, ForceNew, Available since v1.245.0) The key for encryption at rest.
* `trace_on` - (Optional, Bool, Available since v1.245.0) Whether to enable the message trace function. Valid values: `true`, `false`.

### `software`

The software supports the following:
* `maintain_time` - (Optional) Upgrade time period.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time of the resource
* `network_info` - Instance network configuration information
  * `endpoints` - Access point list.
    * `endpoint_type` - Access point type.
    * `endpoint_url` - Access point address.
    * `ip_white_list` - White list of access addresses.
* `product_info` - product info
  * `support_auto_scaling` - is support auto scaling.
* `region_id` - (Available since v1.245.0) The ID of the region in which the instance resides.
* `software` - Instance software information.
  * `software_version` - Software version.
  * `upgrade_method` - Upgrade method.
* `status` - The status of the instance


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Instance.
* `delete` - (Defaults to 5 mins) Used when delete the Instance.
* `update` - (Defaults to 5 mins) Used when update the Instance.

## Import

RocketMQ Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_rocketmq_instance.example <id>
```