---
subcategory: "AliKafka"
layout: "alicloud"
page_title: "Alicloud: alicloud_alikafka_topic"
description: |-
  Provides a Alicloud Alikafka Topic resource.
---

# alicloud_alikafka_topic

Provides a Alikafka Topic resource.

Topic in kafka.

For information about Alikafka Topic and how to use it, see [What is Topic](https://www.alibabacloud.com/help/en/message-queue-for-apache-kafka/latest/api-alikafka-2019-09-16-createtopic).

-> **NOTE:** Available since v1.56.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_alikafka_topic&exampleId=777edb69-3f55-dbe0-541e-4240ae1486ac9df26498&activeTab=example&spm=docs.r.alikafka_topic.0.777edb693f&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_alikafka_instance" "default" {
  name            = var.name
  partition_num   = 50
  disk_type       = "1"
  disk_size       = "500"
  deploy_type     = "5"
  io_max          = "20"
  spec_type       = "professional"
  service_version = "2.2.0"
  vswitch_id      = alicloud_vswitch.default.id
  security_group  = alicloud_security_group.default.id
  config = jsonencode(
    {
      "enable.acl" : "true"
    }
  )
}

resource "alicloud_alikafka_topic" "default" {
  instance_id   = alicloud_alikafka_instance.default.id
  topic         = var.name
  remark        = var.name
  local_topic   = "true"
  compact_topic = "true"
  partition_num = "18"
  configs = jsonencode(
    {
      "message.format.version" : "2.2.0",
      "max.message.bytes" : "10485760",
      "min.insync.replicas" : "1",
      "replication-factor" : "2",
      "retention.ms" : "3600000"
    }
  )
  tags = {
    Created = "TF",
    For     = "example",
  }
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_alikafka_topic&spm=docs.r.alikafka_topic.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `compact_topic` - (Optional, ForceNew, Bool) The cleanup policy for the topic. This parameter is available only if you set the storage engine of the topic to Local storage. Valid values:
  - false: The delete cleanup policy is used.
  - true: The compact cleanup policy is used.
* `configs` - (Optional, Available since v1.262.1) The advanced configurations.
* `instance_id` - (Required, ForceNew) The ID of the instance.
* `local_topic` - (Optional, ForceNew, Bool) The storage engine of the topic. Valid values:
  - false: Cloud storage.
  - true: Local storage.
* `partition_num` - (Optional, Int) The number of partitions in the topic.
* `remark` - (Required) The description of the topic.
* `tags` - (Optional, Map, Available since v1.63.0) A mapping of tags to assign to the resource.
* `topic` - (Required, ForceNew) The topic name.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<instance_id>:<topic>`.
* `create_time` - (Available since v1.262.1) The time when the topic was created.
* `region_id` - (Available since v1.262.1) The ID of the region where the instance resides.
* `status` - (Available since v1.262.1) The status of the service.

## Timeouts

-> **NOTE:** Available since v1.119.0.

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Topic.
* `delete` - (Defaults to 16 mins) Used when delete the Topic.
* `update` - (Defaults to 5 mins) Used when update the Topic.

## Import

Alikafka Topic can be imported using the id, e.g.

```shell
$ terraform import alicloud_alikafka_topic.example <instance_id>:<topic>
```
