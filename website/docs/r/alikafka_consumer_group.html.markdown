---
layout: "alicloud"
page_title: "Alicloud: alicloud_alikafka_consumer_group"
sidebar_current: "docs-alicloud-resource-alikafka-consumer-group"
description: |-
  Provides a Alicloud Alikafka Consumer Group resource.
---

# alicloud\_alikafka\_consumer\_group

Provides an ALIKAFKA consumer group resource.

-> **NOTE:** Available in 1.56.0+

-> **NOTE:**  Only the following regions support create alikafka consumer group.
[`cn-hangzhou`,`cn-beijing`,`cn-shenzhen`,`cn-shanghai`,`cn-qingdao`,`cn-hongkong`,`cn-huhehaote`,`cn-zhangjiakou`,`ap-southeast-1`,`ap-south-1`,`ap-southeast-5`]

## Example Usage

Basic Usage

```
variable "consumer_id" {
  default = "CID-alikafkaGroupDatasourceName"
}

variable "instance_id" {
  default = "xxx"
}

resource "alicloud_alikafka_consumer_group" "default" {
  consumer_id = "${var.consumer_id}"
  instance_id = "${instance_id}"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) ID of the ALIKAFKA Instance that owns the groups.
* `consumer_id` - (Required, ForceNew) ID of the consumer group. The length cannot exceed 64 characters..

## Attributes Reference

The following attributes are exported:

* `id` - The `key` of the resource supplied above. The value is formulated as `<instance_id>:<consumer_id>`.

## Import

ALIKAFKA GROUP can be imported using the id, e.g.

```
$ terraform import alicloud_alikafka_consumer_group.group ALIKAFKA_INST_1234567890_Baso1234567:CID-alikafkaGroupDemo
```