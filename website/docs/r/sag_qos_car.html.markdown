---
subcategory: "Smart Access Gateway (Smartag)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sag_qos_car"
sidebar_current: "docs-alicloud-resource-sag-qos-car"
description: |-
  Provides a Alicloud Sag Qos Car resource.
---

# alicloud_sag_qos_car

Provides a Sag Qos Car resource.

For information about Sag Qos Car and how to use it, see [What is Qos Car](https://www.alibabacloud.com/help/en/smart-access-gateway/latest/createqoscar).

-> **NOTE:** Available since v1.60.0.

-> **NOTE:** Only the following regions support. [`cn-shanghai`, `cn-shanghai-finance-1`, `cn-hongkong`, `ap-southeast-1`, `ap-southeast-2`, `ap-southeast-3`, `ap-southeast-5`, `ap-northeast-1`, `eu-central-1`]

## Example Usage

Basic Usage

```terraform
provider "alicloud" {
  region = "cn-shanghai"
}

variable "name" {
  default = "tf_example"
}

resource "alicloud_sag_qos" "default" {
  name = var.name
}

resource "alicloud_sag_qos_car" "default" {
  qos_id              = alicloud_sag_qos.default.id
  name                = var.name
  description         = var.name
  priority            = "1"
  limit_type          = "Absolute"
  min_bandwidth_abs   = "10"
  max_bandwidth_abs   = "20"
  percent_source_type = "InternetUpBandwidth"
}
```

## Argument Reference

The following arguments are supported:

* `qos_id` - (Required) The instance ID of the QoS.
* `name` - (Optional) The name of the QoS speed limiting rule..
* `description` - (Optional) The description of the QoS speed limiting rule.
* `priority` - (Required) The priority of the specified stream.
* `limit_type` - (Required) The speed limiting method. Valid values: `Absolute`, `Percent`.
* `min_bandwidth_abs` - (Optional) The minimum bandwidth allowed for the stream specified in the quintuple rule. This parameter is required when the value of the LimitType parameter is Absolute.
* `max_bandwidth_abs` - (Optional) The maximum bandwidth allowed for the stream specified in the quintuple rule. This parameter is required when the value of the LimitType is Absolute.
* `min_bandwidth_percent` - (Optional) The minimum bandwidth percentage allowed for the stream specified in the quintuple rule. It is based on the maximum upstream bandwidth you set for the associated SAG instance.This parameter is required when the value of the LimitType parameter is Percent.
* `max_bandwidth_percent` - (Optional) The maximum bandwidth percentage allowed for the stream specified in the quintuple rule. It is based on the maximum upstream bandwidth you set for the associated Smart Access Gateway (SAG) instance.This parameter is required when the value of the LimitType parameter is Percent.
* `percent_source_type` - (Optional) The bandwidth type when the speed is limited based on percentage. Valid values: CcnBandwidth, InternetUpBandwidth.The default value is InternetUpBandwidth.


## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Qos Car. The value formats as `<qos_id>:<qos_car_id>`.

## Import

The Sag Qos Car can be imported using the id, e.g.

```shell
$ terraform import alicloud_sag_qos_car.example <qos_id>:<qos_car_id>
```
