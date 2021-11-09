---
subcategory: "Smart Access Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_sag_qos_car"
sidebar_current: "docs-alicloud-resource-sag-qos-car"
description: |-
  Provides a Sag Qos Car resource.
---

# alicloud\_sag\_qos\_car

Provides a Sag qos car resource. 
You need to create a QoS car to set priorities, rate limits, and quintuple rules for different messages.

For information about Sag Qos Car and how to use it, see [What is Qos Car](https://www.alibabacloud.com/help/doc-detail/140065.htm).

-> **NOTE:** Available in 1.60.0+

-> **NOTE:** Only the following regions support. [`cn-shanghai`, `cn-shanghai-finance-1`, `cn-hongkong`, `ap-southeast-1`, `ap-southeast-2`, `ap-southeast-3`, `ap-southeast-5`, `ap-northeast-1`, `eu-central-1`]

## Example Usage

Basic Usage

```
resource "alicloud_sag_qos" "default" {
  name = "tf-testAccSagQosName"
}
resource "alicloud_sag_qos_car" "default" {
  qos_id                = alicloud_sag_qos.default.id
  name                  = "tf-testSagQosCarName"
  description           = "tf-testSagQosCarDescription"
  priority              = "1"
  limit_type            = "Absolute"
  min_bandwidth_abs     = "10"
  max_bandwidth_abs     = "20"
  min_bandwidth_percent = "10"
  max_bandwidth_percent = "20"
  percent_source_type   = "InternetUpBandwidth"
}
```
## Argument Reference

The following arguments are supported:

* `qos_id` - (Required) The instance ID of the QoS.
* `name` - (Optional) The name of the QoS speed limiting rule..
* `description` - (Optional) The description of the QoS speed limiting rule.
* `priority` - (Required) The priority of the specified stream.
* `limit_type` - (Required) The speed limiting method. Valid values: Absolute, Percent.
* `min_bandwidth_abs` - (Optional) The minimum bandwidth allowed for the stream specified in the quintuple rule. This parameter is required when the value of the LimitType parameter is Absolute.
* `max_bandwidth_abs` - (Optional) The maximum bandwidth allowed for the stream specified in the quintuple rule. This parameter is required when the value of the LimitType is Absolute.
* `min_bandwidth_percent` - (Optional) The minimum bandwidth percentage allowed for the stream specified in the quintuple rule. It is based on the maximum upstream bandwidth you set for the associated SAG instance.This parameter is required when the value of the LimitType parameter is Percent.
* `max_bandwidth_percent` - (Optional) The maximum bandwidth percentage allowed for the stream specified in the quintuple rule. It is based on the maximum upstream bandwidth you set for the associated Smart Access Gateway (SAG) instance.This parameter is required when the value of the LimitType parameter is Percent.
* `percent_source_type` - (Optional) The bandwidth type when the speed is limited based on percentage. Valid values: CcnBandwidth, InternetUpBandwidth.The default value is InternetUpBandwidth.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Qos Car id and formates as `<qos_id>:<qos_car_id>`.

## Import

The Sag Qos Car can be imported using the id, e.g.

```
$ terraform import alicloud_sag_qos_car.example qos-abc123456:qoscar-abc123456
```

