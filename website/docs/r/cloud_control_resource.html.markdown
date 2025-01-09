---
subcategory: "Cloud Control"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_control_resource"
description: |-
  Provides a Alicloud Cloud Control Resource resource.
---

# alicloud_cloud_control_resource

Provides a Cloud Control Resource resource.



For information about Cloud Control Resource and how to use it, see [What is Resource](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.241.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_cloud_control_resource" "mq_instance" {
  desire_attributes = jsonencode({ "InstanceName" : "terraform-example-ons-instance" })
  product           = "Ons"
  resource_code     = "Instance"
}

resource "alicloud_cloud_control_resource" "default" {
  product           = "Ons"
  resource_code     = "Instance::Topic"
  resource_id       = alicloud_cloud_control_resource.mq_instance.resource_id
  desire_attributes = jsonencode({ "InstanceId" : "${alicloud_cloud_control_resource.mq_instance.resource_id}", "TopicName" : "terraform-example-ons-topic", "MessageType" : "1" })
}
```

## Argument Reference

The following arguments are supported:
* `desire_attributes` - (Optional, JsonString) Resource attributes specified when users create and update resources
* `product` - (Required, ForceNew) Products
* `resource_code` - (Required, ForceNew) Resource Code, if there is a parent resource, split with::, such as VPC::VSwitch.
* `resource_id` - (Optional, ForceNew) If there is a parent resource, you need to enter the id of the parent resource, for example, in the VPC::VSwtich resource, you need to enter the id of the VPC: vpc-dexadfe3r4ad. If there are more than one level of parent resources, you need to split them.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<product>:<resource_code>:<resource_id>`.
* `resource_attributes` - Resource Attribute Collection

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Resource.
* `delete` - (Defaults to 5 mins) Used when delete the Resource.
* `update` - (Defaults to 5 mins) Used when update the Resource.

## Import

Cloud Control Resource can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_control_resource.example <product>:<resource_code>:<resource_id>
```