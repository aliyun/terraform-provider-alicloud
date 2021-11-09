---
subcategory: "RabbitMQ (AMQP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_amqp_instances"
sidebar_current: "docs-alicloud-datasource-amqp-instances"
description: |-
  Provides a list of Amqp Instances to the user.
---

# alicloud\_amqp\_instances

This data source provides the Amqp Instances of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.128.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_amqp_instances" "ids" {
  ids = ["amqp-abc12345", "amqp-abc34567"]
}
output "amqp_instance_id_1" {
  value = data.alicloud_amqp_instances.ids.instances.0.id
}

data "alicloud_amqp_instances" "nameRegex" {
  name_regex = "^my-Instance"
}
output "amqp_instance_id_2" {
  value = data.alicloud_amqp_instances.nameRegex.instances.0.id
}

```

## Argument Reference

The following arguments are supported:

* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed)  A list of Instance IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Instance name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `status` - (Optional, ForceNew) The status of the resource. Valid values: "DEPLOYING", "EXPIRED", "RELEASED", "SERVING".

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Instance names.
* `instances` - A list of Amqp Instances. Each element contains the following attributes:
	* `create_time` - OrderCreateTime.
	* `expire_time` - ExpireTime.
	* `id` - The ID of the Instance.
	* `instance_id` - THe instance Id.
	* `instance_name` - THe instance name.
	* `instance_type` - The instance type.
	* `payment_type` - The Pay-as-You-Type Values Include: the Subscription of a Pre-Paid.
	* `private_end_point` - The private endPoint.
	* `public_endpoint` - The public dndpoint.
	* `renewal_duration` - Renewal duration.
	* `renewal_duration_unit` - Auto-Renewal Cycle Unit Values Include: Month: Month. Year: Years.
	* `renewal_status` - Renew status.
	* `status` - The status of the resource.
	* `support_eip` - Whether to support eip.
