---
layout: "alicloud"
page_title: "Alicloud: alicloud_alikafka_consumer_group"
sidebar_current: "docs-alicloud-resource-alikafka-consumer-group"
description: |-
  Provides a Alicloud Alikafka Group resource.
---

# alicloud\_alikafka\_consumer\_group

Provides an ALIKAFKA group resource.

-> **NOTE:** Available in 1.56.0+

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
* `consumer_id` - (Required, ForceNew) Id of the consumer group.

## Attributes Reference

The following attributes are exported:

* `id` - The `key` of the resource supplied above. The value is formulated as `<instance_id>:<consumer_id>`.

## Import

ALIKAFKA GROUP can be imported using the id, e.g.

```
$ terraform import alicloud_alikafka_consumer_group.group ALIKAFKA_INST_1234567890_Baso1234567:CID-alikafkaGroupDemo
```