---
subcategory: "RocketMQ"
layout: "alicloud"
page_title: "Alicloud: alicloud_rocketmq_disaster_recovery_plans"
description: |-
  This data source provides the list of RocketMQ Disaster Recovery Plans.
---

# alicloud_rocketmq_disaster_recovery_plans

This data source provides the list of RocketMQ Disaster Recovery Plans.

-> **NOTE:** Available since v1.247.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_rocketmq_disaster_recovery_plans" "default" {
  ids        = ["<plan_id>"]
  name_regex = "tf-test.*"
}

output "plan_id" {
  value = data.alicloud_rocketmq_disaster_recovery_plans.default.plans.0.plan_id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of disaster recovery plan IDs.
* `instance_id` - (Optional) The ID of the RocketMQ instance to filter by.
* `name_regex` - (Optional) A regex string to filter disaster recovery plans by plan name.
* `output_file` - (Optional) File name where to save the data source results.

## Attributes Reference

The following attributes are exported in addition to the above attributes:

* `plans` - A list of Disaster Recovery Plans. Each element contains the following attributes:
  * `plan_id` - The ID of the disaster recovery plan.
  * `plan_name` - The name of the disaster recovery plan.
  * `plan_desc` - The description of the disaster recovery plan.
  * `plan_type` - The type of the disaster recovery plan.
  * `status` - The status of the disaster recovery plan.
  * `create_time` - The creation time of the disaster recovery plan.
  * `update_time` - The update time of the disaster recovery plan.
  * `region_id` - The region ID of the disaster recovery plan.
  * `auto_sync_checkpoint` - Whether automatic consumption progress synchronization is enabled.
  * `sync_checkpoint_enabled` - Whether synchronization of consumption progress is enabled.
  * `instances` - The list of instances participating in the disaster recovery plan.
    * `instance_id` - The ID of the RocketMQ instance.
    * `instance_role` - The role of the instance.
    * `instance_type` - The type of the instance.
    * `region_id` - The region ID of the instance.
    * `endpoint_url` - The access point URL of the instance.
    * `vpc_id` - The VPC ID of the instance.
    * `vswitch_id` - The vSwitch ID of the instance.
    * `security_group_id` - The security group ID of the instance.
    * `auth_type` - The authentication method.
    * `network_type` - The network type.
    * `username` - The authentication username.
    * `password` - The authentication password.
    * `consumer_group_id` - The consumer group ID.
    * `message_property` - The message property filter.
      * `property_key` - The key of the message property.
      * `property_value` - The value of the message property.
