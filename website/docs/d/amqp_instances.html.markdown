---
subcategory: "RabbitMQ (AMQP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_amqp_instances"
sidebar_current: "docs-alicloud-datasource-amqp-instances"
description: |-
  Provides a list of Amqp Instances to the user.
---

# alicloud_amqp_instances

This data source provides the Amqp Instances of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.128.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_amqp_instance" "default" {
  instance_name         = var.name
  instance_type         = "enterprise"
  max_tps               = 3000
  max_connections       = 2000
  queue_capacity        = 200
  payment_type          = "Subscription"
  renewal_status        = "AutoRenewal"
  renewal_duration      = 1
  renewal_duration_unit = "Year"
  support_eip           = true
}

data "alicloud_amqp_instances" "ids" {
  ids = [alicloud_amqp_instance.default.id]
}

output "amqp_instance_id_0" {
  value = data.alicloud_amqp_instances.ids.instances.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, List) A list of Instance IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Instance name.
* `status` - (Optional, ForceNew) The status of the resource. Valid values: `DEPLOYING`, `SERVING`, `EXPIRED`, `RELEASED`.
* `enable_details` - (Optional, Bool) Whether to query the detailed list of resource attributes. Default value: `false`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Instance names.
* `instances` - A list of Amqp Instances. Each element contains the following attributes:
  * `id` - The ID of the Instance.
  * `instance_id` - THe instance Id.
  * `instance_type` - The instance type.
  * `instance_name` - THe instance name.
  * `public_endpoint` - The public endpoint of the instance.
  * `private_end_point` - The virtual private cloud (VPC) endpoint of the instance.
  * `support_eip` - Indicates whether the instance supports elastic IP addresses (EIPs).
  * `payment_type` - The billing method of the instance. **Note:** `payment_type` takes effect only if `enable_details` is set to `true`.
  * `renewal_status` - Whether to renew an instance automatically or not. **Note:** `renewal_status` takes effect only if `enable_details` is set to `true`.
  * `renewal_duration` - Auto renewal period of an instance. **Note:** `renewal_duration` takes effect only if `enable_details` is set to `true`.
  * `renewal_duration_unit` - Automatic renewal period unit. **Note:** `renewal_duration_unit` takes effect only if `enable_details` is set to `true`.
  * `status` - The status of the instance.
  * `create_time` - The timestamp that indicates when the order was created.
  * `expire_time` - The timestamp that indicates when the instance expires.
