---
subcategory: "Express Connect"
layout: "alicloud"
page_title: "Alicloud: alicloud_express_connect_traffic_qos_queue"
description: |-
  Provides a Alicloud Express Connect Traffic Qos Queue resource.
---

# alicloud_express_connect_traffic_qos_queue

Provides a Express Connect Traffic Qos Queue resource.

Express Connect Traffic QoS Queue.

For information about Express Connect Traffic Qos Queue and how to use it, see [What is Traffic Qos Queue](https://next.api.alibabacloud.com/document/Vpc/2016-04-28/CreateExpressConnectTrafficQosQueue).

-> **NOTE:** Available since v1.224.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-shanghai"
}

data "alicloud_express_connect_physical_connections" "default" {
  name_regex = "preserved-NODELETING"
}

resource "alicloud_express_connect_traffic_qos" "createQos" {
  qos_name        = var.name
  qos_description = "terraform-example"
}

resource "alicloud_express_connect_traffic_qos_association" "associateQos" {
  instance_id   = data.alicloud_express_connect_physical_connections.default.ids.1
  qos_id        = alicloud_express_connect_traffic_qos.createQos.id
  instance_type = "PHYSICALCONNECTION"
}

resource "alicloud_express_connect_traffic_qos_queue" "createQosQueue" {
  qos_id            = alicloud_express_connect_traffic_qos.createQos.id
  bandwidth_percent = "60"
  queue_description = "terraform-example"
  queue_name        = var.name
  queue_type        = "Medium"
}
```

## Argument Reference

The following arguments are supported:
* `bandwidth_percent` - (Optional, Computed) QoS queue bandwidth percentage.

  - When the QoS queue type is `Medium`, this field must be entered. Valid values: 1 to 100.
  - When the QoS queue type is `Default`, this field is "-".
* `qos_id` - (Required, ForceNew) The ID of the QoS policy.
* `queue_description` - (Optional) The description of the QoS queue.
The length is 0 to 256 characters and cannot start with 'http:// 'or 'https.
* `queue_name` - (Optional) The name of the QoS queue.
The length is 0 to 128 characters and cannot start with 'http:// 'or 'https.
* `queue_type` - (Required, ForceNew) QoS queue type, value:
  - `High`: High priority queue.
  - `Medium`: Normal priority queue.
  - `Default`: the Default priority queue.

-> **NOTE:**  Default priority queue cannot be created.


## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<qos_id>:<queue_id>`.
* `queue_id` - The ID of the QoS queue.
* `status` - The status of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 8 mins) Used when create the Traffic Qos Queue.
* `delete` - (Defaults to 8 mins) Used when delete the Traffic Qos Queue.
* `update` - (Defaults to 8 mins) Used when update the Traffic Qos Queue.

## Import

Express Connect Traffic Qos Queue can be imported using the id, e.g.

```shell
$ terraform import alicloud_express_connect_traffic_qos_queue.example <qos_id>:<queue_id>
```