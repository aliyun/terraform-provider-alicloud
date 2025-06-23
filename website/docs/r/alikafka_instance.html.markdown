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
[`cn-hangzhou`,`cn-beijing`,`cn-shenzhen`,`cn-shanghai`,`cn-qingdao`,`cn-hongkong`,`cn-huhehaote`,`cn-zhangjiakou`,`cn-chengdu`,`cn-heyuan`,`ap-southeast-1`,`ap-southeast-3`,`ap-southeast-5`,`ap-northeast-1`,`eu-central-1`,`eu-west-1`,`us-west-1`,`us-east-1`]

-> **NOTE:** Only the following regions support create alikafka post paid instance. 
[`cn-hangzhou`,`cn-beijing`,`cn-shenzhen`,`cn-shanghai`,`cn-qingdao`,`cn-hongkong`,`cn-huhehaote`,`cn-zhangjiakou`,`cn-chengdu`,`cn-heyuan`,`ap-southeast-1`,`ap-southeast-3`,`ap-southeast-5`,`ap-northeast-1`,`eu-central-1`,`eu-west-1`,`us-west-1`,`us-east-1`]

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_alikafka_instance&exampleId=5e353142-9caa-4f03-436d-fad8a0e042b719460d1f&activeTab=example&spm=docs.r.alikafka_instance.0.5e3531429c&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "instance_name" {
  default = "terraform-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
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
  name           = "${var.instance_name}-${random_integer.default.result}"
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

* `deploy_type` - (Required) The deployment type of the instance. **NOTE:** From version 1.161.0, this attribute supports to be updated. Valid values:
  - `4`: eip/vpc instance
  - `5`: vpc instance.
* `instance_type` - (Optional, ForceNew, Available since v1.253.0) The type of the Instance. Default value: `alikafka`. Valid values:
  - `alikafka`: Kafka Instance.
  - `alikafka_serverless`: Serverless ApsaraMQ for Kafka Instance.
  - `alikafka_confluent`: ApsaraMQ for Confluent Instance.
* `partition_num` - (Optional, Available since v1.194.0) The number of partitions.
* `topic_quota` - (Deprecated since v1.194.0) The max num of topic can be creation of the instance.
  It has been deprecated since version 1.194.0 and using `partition_num` instead.
  Currently, its value only can be set to 50 when creating it, and finally depends on `partition_num` value: <`topic_quota`> = 1000 + <`partition_num`>.
  Therefore, you can update it by updating the `partition_num`, and it is the only updating path.
* `disk_type` - (Optional, ForceNew) The disk type of the instance. Valid values:
  - `0`: efficient cloud disk.
  - `1`: SSD.
-> **NOTE:** If `instance_type` is set to `alikafka`, `disk_type` is required.
* `disk_size` - (Optional, Int) The disk size of the instance. When modify this value, it only supports adjust to a greater value.
-> **NOTE:** If `instance_type` is set to `alikafka`, `disk_size` is required.
* `io_max` - (Optional) The max value of io of the instance. When modify this value, it only support adjust to a greater value.
* `io_max_spec` - (Optional, Available since v1.201.0) The traffic specification of the instance. We recommend that you configure this parameter.
  - You should specify one of the `io_max` and `io_max_spec` parameters, and `io_max_spec` is recommended.
  - For more information about the valid values, see [Billing](https://www.alibabacloud.com/help/en/message-queue-for-apache-kafka/latest/billing-overview).
* `eip_max` - (Optional) The max bandwidth of the instance. It will be ignored when `deploy_type = 5`. When modify this value, it only supports adjust to a greater value.
* `resource_group_id` - (Optional, Available since v1.224.0) The ID of the resource group. **Note:** Once you set a value of this property, you cannot set it to an empty string anymore.
* `name` - (Optional) Name of your Kafka instance. The length should between 3 and 64 characters. If not set, will use instance id as instance name.
* `paid_type` - (Optional) The billing method of the instance. Default value: `PostPaid`. Valid values: `PostPaid`, `PrePaid`. When modify this value, it only support adjust from `PostPaid` to `PrePaid`.
* `spec_type` - (Optional) The instance edition. Default value: `normal`. Valid values:
  - If `instance_type` is set to `alikafka`. Valid values: `normal`, `professional`, `professionalForHighRead`.
  - If `instance_type` is set to `alikafka_serverless`. Valid values: `normal`.
  - If `instance_type` is set to `alikafka_confluent`. Valid values: `professional`, `enterprise`.
* `vswitch_id` - (Optional, ForceNew) The ID of attaching vswitch to instance.
* `security_group` - (Optional, ForceNew, Available since v1.93.0) The ID of security group for this instance. If the security group is empty, system will create a default one.
* `service_version` - (Optional, Available since v1.112.0) The version of the Instance. Valid values:
  - If `instance_type` is set to `alikafka`. Default value: `2.2.0`. Valid values: `2.2.0`, `2.6.2`.
  - If `instance_type` is set to `alikafka_serverless`. Default value: `3.3.1`. Valid values: `3.3.1`.
  - If `instance_type` is set to `alikafka_confluent`. Default value: `7.4.0`. Valid values: `7.4.0`.
* `config` - (Optional, Available since v1.112.0) The initial configurations of the ApsaraMQ for Kafka instance. The values must be valid JSON strings. The `config` supports the following parameters:
  * `enable.vpc_sasl_ssl`: Specifies whether to enable VPC transmission encryption. Default value: `false`. Valid values:
    - `true`: Enables VPC transmission encryption. If you enable VPC transmission encryption, you must also enable access control list (ACL).
    - `false`: Disables VPC transmission encryption. This is the default value.
  * `enable.acl`: Specifies whether to enable ACL. Default value: `false`. Valid values:
    - `true`: Enables the ACL feature.
    - `false`: Disables the ACL feature.
  * `kafka.log.retention.hours`: The maximum message retention period when the disk capacity is sufficient. Unit: hours. Default value: `72`. Valid values: `24` to `480`.
  * `kafka.message.max.bytes`: The maximum size of a message that can be sent and received by ApsaraMQ for Kafka. Unit: bytes. Default value: `1048576`. Valid values: `1048576` to `10485760`.
* `kms_key_id` - (Optional, ForceNew, Available since v1.180.0) The ID of the key that is used to encrypt data on standard SSDs in the region of the instance.
* `tags` - (Optional, Available since v1.63.0) A mapping of tags to assign to the resource.
* `vpc_id` - (Optional, ForceNew, Available since v1.185.0) The VPC ID of the instance.
* `zone_id` - (Optional, ForceNew, Available since v1.185.0) The zone ID of the instance. The value can be in zone x or region id-x format. **NOTE**: When the available zone is insufficient, another availability zone may be deployed.
* `enable_auto_group` - (Optional, Bool, Available since v1.241.0) Specify whether to enable the flexible group creation feature. Default value: `false`. Valid values:
  - `true`: Enables the flexible group creation feature.
  - `false`: Disabled the flexible group creation feature.
* `enable_auto_topic` - (Optional, Available since v1.241.0) Specify whether to enable the automatic topic creation feature. Default value: `disable`. Valid values:
  - `enable`: Enables the automatic topic creation feature.
  - `disable`: Disabled the automatic topic creation feature.
* `default_topic_partition_num` - (Optional, Int, Available since v1.241.0) The number of partitions in a topic that is automatically created.
* `password` - (Optional, Available since v1.253.0) The instance password. **NOTE:** If `instance_type` is set to `alikafka_confluent`, `password` is required.
* `vswitch_ids` - (Optional, ForceNew, List, Available since v1.241.0) The IDs of the vSwitches with which the instance is associated.
* `selected_zones` - (Optional, List, Available since v1.195.0) The zones among which you want to deploy the instance.
* `serverless_config` - (Optional, Set, Available since v1.253.0) The parameters configured for the serverless ApsaraMQ for Kafka instance. See [`serverless_config`](#serverless_config) below.
-> **NOTE:** If `instance_type` is set to `alikafka_serverless`, `serverless_config` is required.
* `confluent_config` - (Optional, Set, Available since v1.253.0) The configurations of Confluent. See [`confluent_config`](#confluent_config) below.
-> **NOTE:** If `instance_type` is set to `alikafka_confluent`, `confluent_config` is required.

-> **NOTE:** Field `io_max`, `disk_size`, `topic_quota`, `eip_max` should follow the following constraints.

| io_max | disk_size(min-max:lag) | topic_quota(min-max:lag) | eip_max(min-max:lag) | 
|------|-------------|:----:|:-----:|
|20          |  500-6100:100   |   50-450:1  |    1-160:1  |
|30          |  800-6100:100   |   50-450:1  |    1-240:1  |
|60          |  1400-6100:100  |   80-450:1  |    1-500:1  |
|90          |  2100-6100:100  |   100-450:1 |    1-500:1  |
|120         |  2700-6100:100  |   150-450:1 |    1-500:1  |

### `serverless_config`

The serverless_config supports the following:

* `reserved_publish_capacity` - (Optional, Int) The reserved capacity for publishing messages.
* `reserved_subscribe_capacity` - (Optional, Int) The reserved capacity for subscribing to message.

### `confluent_config`

The confluent_config supports the following:

* `kafka_cu` - (Optional, Int) The number of CPU cores of the Kafka broker.
* `kafka_storage` - (Optional, Int) The disk capacity of the Kafka broker.
* `kafka_replica` - (Optional, Int) The number of replicas of the Kafka broker.
* `kafka_rest_proxy_cu` - (Optional, Int) The number of CPU cores of Kafka Rest Proxy.
* `kafka_rest_proxy_replica` - (Optional, Int) The number of replicas of Kafka Rest Proxy.
* `zookeeper_cu` - (Optional, Int) The number of CPU cores of ZooKeeper.
* `zookeeper_storage` - (Optional, Int) The disk capacity of ZooKeeper.
* `zookeeper_replica` - (Optional, ForceNew, Int) The number of replicas of ZooKeeper.
* `control_center_cu` - (Optional, Int) The number of CPU cores of Control Center.
* `control_center_storage` - (Optional, Int) The disk capacity of Control Center.
* `control_center_replica` - (Optional, ForceNew, Int) The number of replicas of Control Center.
* `schema_registry_cu` - (Optional, Int) The number of CPU cores of Schema Registry.
* `schema_registry_replica` - (Optional, Int) The number of replicas of Schema Registry.
* `connect_cu` - (Optional, Int) The number of CPU cores of Connect.
* `connect_replica` - (Optional, Int) The number of replicas of Connect.
* `ksql_cu` - (Optional, Int) The number of CPU cores of ksqlDB.
* `ksql_storage` - (Optional, Int) The disk capacity of ksqlDB.
* `ksql_replica` - (Optional, Int) The number of replicas of ksqlDB.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Instance.
* `end_point` - The EndPoint to access the kafka instance.
* `ssl_endpoint` - (Available since v1.234.0) The Secure Sockets Layer (SSL) endpoint of the instance in IP address mode.
* `domain_endpoint` - (Available since v1.234.0) The default endpoint of the instance in domain name mode.
* `ssl_domain_endpoint` - (Available since v1.234.0) The SSL endpoint of the instance in domain name mode.
* `sasl_domain_endpoint` - (Available since v1.234.0) The Simple Authentication and Security Layer (SASL) endpoint of the instance in domain name mode.
* `topic_num_of_buy` - (Available since v1.214.1) The number of purchased topics.
* `topic_used` - (Available since v1.214.1) The number of used topics.
* `topic_left` - (Available since v1.214.1) The number of available topics.
* `partition_used` - (Available since v1.214.1) The number of used partitions.
* `partition_left` - (Available since v1.214.1) The number of available partitions.
* `group_used` - (Available since v1.214.1) The number of used groups.
* `group_left` - (Available since v1.214.1) The number of available groups.
* `is_partition_buy` - (Available since v1.214.1) The method that you use to purchase partitions.
* `status` - The status of the instance.

## Timeouts

-> **NOTE:** Available since v1.180.0.

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 60 mins) Used when create the resource.
* `update` - (Defaults to 120 mins) Used when update the resource.
* `delete` - (Defaults to 30 mins) Used when delete the resource.

## Import

AliKafka instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_alikafka_instance.example <id>
```
