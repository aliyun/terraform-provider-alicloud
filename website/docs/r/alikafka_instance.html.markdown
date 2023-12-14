---
subcategory: "AliKafka"
layout: "alicloud"
page_title: "Alicloud: alicloud_alikafka_instance"
sidebar_current: "docs-alicloud-resource-alikafka-instance"
description: |-
  Provides a Alicloud AliKafka Instance resource.
---

# alicloud_alikafka_instance

Provides an AliKafka instance resource.

For information about Kafka instance and how to use it, see [What is alikafka instance](https://www.alibabacloud.com/help/en/message-queue-for-apache-kafka/latest/api-alikafka-2019-09-16-startinstance).

-> **NOTE:** Available since v1.59.0.

-> **NOTE:** Creation or modification may took about 10-40 minutes.

-> **NOTE:** Only the following regions support create alikafka pre paid instance.
[`cn-hangzhou`,`cn-beijing`,`cn-shenzhen`,`cn-shanghai`,`cn-qingdao`,`cn-hongkong`,`cn-huhehaote`,`cn-zhangjiakou`,`cn-chengdu`,`cn-heyuan`,`ap-southeast-1`,`ap-southeast-3`,`ap-southeast-5`,`ap-south-1`,`ap-northeast-1`,`eu-central-1`,`eu-west-1`,`us-west-1`,`us-east-1`]

-> **NOTE:** Only the following regions support create alikafka post paid instance. 
[`cn-hangzhou`,`cn-beijing`,`cn-shenzhen`,`cn-shanghai`,`cn-qingdao`,`cn-hongkong`,`cn-huhehaote`,`cn-zhangjiakou`,`cn-chengdu`,`cn-heyuan`,`ap-southeast-1`,`ap-southeast-3`,`ap-southeast-5`,`ap-south-1`,`ap-northeast-1`,`eu-central-1`,`eu-west-1`,`us-west-1`,`us-east-1`]

## Example Usage

Basic Usage

```terraform
variable "instance_name" {
  default = "alikafkaInstanceName"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id     = alicloud_vpc.default.id
  cidr_block = "172.16.0.0/24"
  zone_id    = data.alicloud_zones.default.zones[0].id
}

resource "alicloud_security_group" "default" {
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_alikafka_instance" "default" {
  name           = var.instance_name
  partition_num  = 50
  disk_type      = 1
  disk_size      = 500
  deploy_type    = 5
  io_max         = 20
  vswitch_id     = alicloud_vswitch.default.id
  security_group = alicloud_security_group.default.id
}
```

### Removing alicloud_alikafka_instance from your configuration

The alicloud_alikafka_instance resource allows you to manage your alikafka instance, but Terraform cannot destroy it if your instance type is pre paid(post paid type can destroy normally). Removing this resource from your configuration will remove it from your statefile and management, but will not destroy the instance. You can resume managing the instance via the alikafka Console.

## Argument Reference

The following arguments are supported:

* `name` - (Optional) Name of your Kafka instance. The length should between 3 and 64 characters. If not set, will use instance id as instance name.
* `partition_num` - (Optional, Available since v1.194.0) The number of partitions.
* `topic_quota` - (Deprecated since v1.194.0) The max num of topic can be creation of the instance.
  It has been deprecated since version 1.194.0 and using `partition_num` instead.
  Currently, its value only can be set to 50 when creating it, and finally depends on `partition_num` value: <`topic_quota`> = 1000 + <`partition_num`>.
  Therefore, you can update it by updating the `partition_num`, and it is the only updating path.
* `disk_type` - (Required, ForceNew) The disk type of the instance. 0: efficient cloud disk , 1: SSD.
* `disk_size` - (Required) The disk size of the instance. When modify this value, it only supports adjust to a greater value.
* `deploy_type` - (Required) The deployment type of the instance. **NOTE:** From version 1.161.0, this attribute supports to be updated. Valid values:
  - 4: eip/vpc instance
  - 5: vpc instance.
* `io_max` - (Optional) The max value of io of the instance. When modify this value, it only support adjust to a greater value.
* `io_max_spec` - (Optional, Available since v1.201.0) The traffic specification of the instance. We recommend that you configure this parameter.
  - You should specify one of the `io_max` and `io_max_spec` parameters, and `io_max_spec` is recommended.
  - For more information about the valid values, see [Billing](https://www.alibabacloud.com/help/en/message-queue-for-apache-kafka/latest/billing-overview).
* `eip_max` - (Optional) The max bandwidth of the instance. It will be ignored when `deploy_type = 5`. When modify this value, it only supports adjust to a greater value.
* `paid_type` - (Optional) The paid type of the instance. Support two type, "PrePaid": pre paid type instance, "PostPaid": post paid type instance. Default is PostPaid. When modify this value, it only support adjust from post pay to pre pay. 
* `spec_type` - (Optional) The spec type of the instance. Support two type, "normal": normal version instance, "professional": professional version instance. Default is normal. When modify this value, it only support adjust from normal to professional. Note only pre paid type instance support professional specific type.
* `vswitch_id` - (Required, ForceNew) The ID of attaching vswitch to instance.
* `security_group` - (Optional, ForceNew, Available since v1.93.0) The ID of security group for this instance. If the security group is empty, system will create a default one.
* `service_version` - (Optional, Available since v1.112.0) The kafka openSource version for this instance. Only 0.10.2 or 2.2.0 is allowed, default is 0.10.2.
* `config` - (Optional, Available since v1.112.0) The basic config for this instance. The input should be json type, only the following key allowed: enable.acl, enable.vpc_sasl_ssl, kafka.log.retention.hours, kafka.message.max.bytes.
* `kms_key_id` - (Optional, ForceNew, Available since v1.180.0) The ID of the key that is used to encrypt data on standard SSDs in the region of the instance.
* `tags` - (Optional, Available since v1.63.0) A mapping of tags to assign to the resource.
* `vpc_id` - (Optional, ForceNew, Available since v1.185.0) The VPC ID of the instance.
* `zone_id` - (Optional, ForceNew, Available since v1.185.0) The zone ID of the instance.
* `selected_zones` - (Optional, Available since v1.195.0) The zones among which you want to deploy the instance.

-> **NOTE:** Arguments io_max, disk_size, topic_quota, eip_max should follow the following constraints.

| io_max | disk_size(min-max:lag) | topic_quota(min-max:lag) | eip_max(min-max:lag) | 
|------|-------------|:----:|:-----:|
|20          |  500-6100:100   |   50-450:1  |    1-160:1  |
|30          |  800-6100:100   |   50-450:1  |    1-240:1  |
|60          |  1400-6100:100  |   80-450:1  |    1-500:1  |
|90          |  2100-6100:100  |   100-450:1 |    1-500:1  |
|120         |  2700-6100:100  |   150-450:1 |    1-500:1  |

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Instance.
* `end_point` - The EndPoint to access the kafka instance.
* `topic_num_of_buy` - (Available since v1.213.1) The number of purchased topics.
* `topic_used` - (Available since v1.213.1) The number of used topics.
* `topic_left` - (Available since v1.213.1) The number of available topics.
* `partition_used` - (Available since v1.213.1) The number of used partitions.
* `partition_left` - (Available since v1.213.1) The number of available partitions.
* `group_used` - (Available since v1.213.1) The number of used groups.
* `group_left` - (Available since v1.213.1) The number of available groups.
* `is_partition_buy` - (Available since v1.213.1) The method that you use to purchase partitions.
* `status` - The status of the instance.

## Timeouts

-> **NOTE:** Available since v1.180.0.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 mins) Used when create the resource.
* `update` - (Defaults to 120 mins) Used when update the resource.
* `delete` - (Defaults to 30 mins) Used when delete the resource.

## Import

AliKafka instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_alikafka_instance.instance <id>
```
