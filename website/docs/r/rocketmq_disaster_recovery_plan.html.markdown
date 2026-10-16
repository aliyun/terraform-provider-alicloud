---
subcategory: "RocketMQ"
layout: "alicloud"
page_title: "Alicloud: alicloud_rocketmq_disaster_recovery_plan"
description: |-
  Provides a Alicloud RocketMQ Disaster Recovery Plan resource.
---

# alicloud_rocketmq_disaster_recovery_plan

Provides a RocketMQ Disaster Recovery Plan resource.

For information about RocketMQ Disaster Recovery Plan and how to use it, see [What is RocketMQ Disaster Recovery Plan](https://www.alibabacloud.com/help/en/apsaramq-for-rocketmq/cloud-message-queue-rocketmq-5-x-series/developer-reference/api-rocketmq-2022-08-01-createdisasterrecoveryplan).

-> **NOTE:** Available since v1.295.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_rocketmq_disaster_recovery_plan" "default" {
  plan_name               = var.name
  plan_desc               = "disaster recovery plan description"
  plan_type               = "ACTIVE_PASSIVE"
  auto_sync_checkpoint    = true
  sync_checkpoint_enabled = true
  instances {
    instance_id       = "rmq-xxx"
    instance_role     = "ACTIVE"
    instance_type     = "ALIYUN_ROCKETMQ"
    region_id         = "cn-hangzhou"
    endpoint_url      = "rmq-xxx.cn-hangzhou.aliyuncs.com"
    vpc_id            = "vpc-xxx"
    vswitch_id        = "vsw-xxx"
    security_group_id = "sg-xxx"
    auth_type         = "ACL_AUTH"
    network_type      = "TCP_VPC"
    username          = "example-user"
    password          = "Example123456!"
    consumer_group_id = "cg-xxx"
    message_property {
      property_key   = "tag"
      property_value = "*"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `plan_name` - (Optional) The name of the disaster recovery plan.
* `plan_desc` - (Optional) The description of the disaster recovery plan.
* `plan_type` - (Optional, ForceNew) The type of the disaster recovery plan. Valid values: `ACTIVE_PASSIVE`, `ACTIVE_ACTIVE`.
* `auto_sync_checkpoint` - (Optional) Whether to enable automatic consumption progress synchronization. It is only effective when `sync_checkpoint_enabled` is `true`. Default to `false`.
* `sync_checkpoint_enabled` - (Optional) Whether to enable synchronization of consumption progress. Default to `false`.
* `instances` - (Optional) The list of instances participating in the disaster recovery plan. See [`instances`](#instances) below.

### `instances`

The instances block supports the following:

* `instance_id` - (Optional, ForceNew) The ID of the RocketMQ instance. It is required when `instance_type` is `ALIYUN_ROCKETMQ`.
* `instance_role` - (Optional, ForceNew) The role of the instance in the disaster recovery plan. Valid values: `ACTIVE`, `PASSIVE`.
* `instance_type` - (Optional, ForceNew) The type of the instance. Valid values: `ALIYUN_ROCKETMQ`, `EXTERNAL_ROCKETMQ`.
* `region_id` - (Optional, ForceNew) The region ID of the instance.
* `endpoint_url` - (Optional, ForceNew) The access point URL of the instance. It is required when `instance_type` is `EXTERNAL_ROCKETMQ`.
* `vpc_id` - (Optional, ForceNew) The VPC ID of the instance. It is required when `instance_type` is `EXTERNAL_ROCKETMQ` and `network_type` is `TCP_VPC`.
* `vswitch_id` - (Optional, ForceNew) The vSwitch ID of the instance. It is required when `instance_type` is `EXTERNAL_ROCKETMQ` and `network_type` is `TCP_VPC`.
* `security_group_id` - (Optional, ForceNew) The security group ID of the instance. It is required when `instance_type` is `EXTERNAL_ROCKETMQ` and `network_type` is `TCP_VPC`.
* `auth_type` - (Optional) The authentication method of the instance. Valid values: `NO_AUTH`, `ACL_AUTH`.
* `network_type` - (Optional, ForceNew) The network type of the instance. Valid values: `TCP_INTERNET`, `TCP_VPC`.
* `username` - (Optional) The authentication username. It is required when `auth_type` is `ACL_AUTH`.
* `password` - (Optional, Sensitive) The authentication password. It is required when `auth_type` is `ACL_AUTH`.
* `consumer_group_id` - (Optional, ForceNew, Computed) The consumer group ID used to pull messages from the source instance. It is automatically generated for Alibaba Cloud instances when not specified.
* `message_property` - (Optional) The message property filter. See [`message_property`](#instances-message_property) below.

### `instances-message_property`

The message_property block supports the following:

* `property_key` - (Optional) The key of the message property.
* `property_value` - (Optional) The value of the message property.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of the Disaster Recovery Plan. The value is the plan ID.
* `plan_id` - The ID of the disaster recovery plan.
* `status` - The status of the disaster recovery plan.
* `create_time` - The creation time of the disaster recovery plan.
* `update_time` - The update time of the disaster recovery plan.
* `region_id` - The region ID of the disaster recovery plan.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#custom-timeouts) for certain actions:

* `create` - (Defaults to 20 mins) Used when creating the Disaster Recovery Plan.
* `update` - (Defaults to 5 mins) Used when updating the Disaster Recovery Plan.
* `delete` - (Defaults to 5 mins) Used when deleting the Disaster Recovery Plan.

## Import

Disaster Recovery Plan can be imported using the id, e.g.

```shell
$ terraform import alicloud_rocketmq_disaster_recovery_plan.example <plan_id>
```
