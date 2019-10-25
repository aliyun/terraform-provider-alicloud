---
layout: "alicloud"
page_title: "Alicloud: alicloud_alikafka_instance"
sidebar_current: "docs-alicloud-resource-alikafka-instance"
description: |-
  Provides a Alicloud ALIKAFKA Instance resource.
---

# alicloud\_alikafka\_instance

Provides an ALIKAFKA instance resource.

-> **NOTE:** Available in 1.59.0+

-> **NOTE:** ALIKAFKA instance resource only support create post pay instance. Creation or modification may took about 10-40 minutes.

-> **NOTE:** Only the following regions support create alikafka instance.
[`cn-hangzhou`,`cn-beijing`,`cn-shenzhen`,`cn-shanghai`,`cn-qingdao`,`cn-hongkong`,`cn-huhehaote`,`cn-zhangjiakou`,`ap-southeast-1`,`ap-south-1`,`ap-southeast-5`]

## Example Usage

Basic Usage

```
variable "instance_name" {
 default = "alikafkaInstanceName"
}

data "alicloud_zones" "default" {
    available_resource_creation= "VSwitch"
}
resource "alicloud_vpc" "default" {
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "172.16.0.0/24"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_alikafka_instance" "default" {
  name = "${var.instance_name}"
  topic_quota = "50"
  disk_type = "1"
  disk_size = "500"
  deploy_type = "4"
  io_max = "20"
  vswitch_id = "${alicloud_vswitch.default.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) Name of your Kafka instance. The length should between 3 and 64 characters.
* `topic_quota` - (Required) The max num of topic can be create of the instance. When modify this value, it only adjust to a greater value.
* `disk_type` - (Required, ForceNew) The disk type of the instance. 0: efficient cloud disk , 1: SSD.
* `disk_size` - (Required) The disk size of the instance. When modify this value, it only adjust to a greater value.
* `deploy_type` - (Required, ForceNew) The deploy type of the instance. Currently only support two deploy type, 4: eip/vpc instance, 5: vpc instance.
* `io_max` - (Required) The max value of io of the instance. When modify this value, it only adjust to a greater value.
* `eip_max` - (Optional) The max bandwidth of the instance. When modify this value, it only adjust to a greater value.
* `vswitch_id` - (Required, ForceNew) The ID of attaching vswitch to instance.

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

* `id` - The `key` of the resource supplied above, also call instance id.
* `vpc_id` - The ID of attaching VPC to instance.
* `zone_id` - The Zone to launch the kafka instance.

## Import

ALIKAFKA TOPIC can be imported using the id, e.g.

```
$ terraform import alicloud_alikafka_instance.instance KAFKA_INST_1234567890_Baso1234567
```