---
subcategory: "Express Connect"
layout: "alicloud"
page_title: "Alicloud: alicloud_express_connect_traffic_qos_association"
description: |-
  Provides a Alicloud Express Connect Traffic Qos Association resource.
---

# alicloud_express_connect_traffic_qos_association

Provides a Express Connect Traffic Qos Association resource. Express Connect QoS associated resources.

For information about Express Connect Traffic Qos Association and how to use it, see [What is Traffic Qos Association](https://www.alibabacloud.com/help/en/).

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
```

## Argument Reference

The following arguments are supported:
* `instance_id` - (Optional, ForceNew, Computed) The ID of the associated instance.
* `instance_type` - (Optional, ForceNew, Computed) The type of the associated instance. Value: **physical connection** physical connection.
* `qos_id` - (Required, ForceNew) The QoS policy ID.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<qos_id>:<instance_id>:<instance_type>`.
* `status` - The status of the associated instance. Value:

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Traffic Qos Association.
* `delete` - (Defaults to 5 mins) Used when delete the Traffic Qos Association.

## Import

Express Connect Traffic Qos Association can be imported using the id, e.g.

```shell
$ terraform import alicloud_express_connect_traffic_qos_association.example <qos_id>:<instance_id>:<instance_type>
```