---
subcategory: "AliKafka"
layout: "alicloud"
page_title: "Alicloud: alicloud_alikafka_consumer_group"
description: |-
  Provides a Alicloud AliKafka Consumer Group resource.
---

# alicloud_alikafka_consumer_group

Provides a Ali Kafka Consumer Group resource.

Group in kafka.

For information about Ali Kafka Consumer Group and how to use it, see [What is Consumer Group](https://next.api.alibabacloud.com/document/alikafka/2019-09-16/CreateConsumerGroup).

-> **NOTE:** Available since v1.56.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_alikafka_consumer_group&exampleId=e0c01fec-c5b8-1e8d-7fd6-c01a9798ede2dfa2cf90&activeTab=example&spm=docs.r.alikafka_consumer_group.0.e0c01fecc5&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_alikafka_instances" "default" {
}

resource "alicloud_alikafka_consumer_group" "default" {
  instance_id = data.alicloud_alikafka_instances.default.instances.0.id
  consumer_id = var.name
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_alikafka_consumer_group&spm=docs.r.alikafka_consumer_group.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) ID of the ALIKAFKA Instance that owns the groups.
* `consumer_id` - (Required, ForceNew) ID of the consumer group.
* `tags` - (Optional, Available since v1.63.0) A mapping of tags to assign to the resource.
* `remark` - (Optional, ForceNew, Available since v1.268.0) The remark of the resource.
* `description` - (Deprecated since v1.268.0) Field `description` has been deprecated from provider version 1.268.0. New field `remark` instead.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<instance_id>:<consumer_id>`.
* `create_time` - (Available since v1.268.0) The timestamp of when the group was created.
* `region_id` - (Available since v1.268.0) The region ID.

## Timeouts

-> **NOTE:** Available since v1.268.0.

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Consumer Group.
* `delete` - (Defaults to 5 mins) Used when delete the Consumer Group.
* `update` - (Defaults to 5 mins) Used when update the Consumer Group.

## Import

AliKafka Consumer Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_alikafka_consumer_group.example <instance_id>:<consumer_id>
```
