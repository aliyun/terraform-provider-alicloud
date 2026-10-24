---
subcategory: "Elasticsearch"
layout: "alicloud"
page_title: "Alicloud: alicloud_elasticsearch_snapshot_policy"
sidebar_current: "docs-alicloud-resource-elasticsearch-snapshot-policy"
description: |-
  Provides a Elasticsearch Snapshot Policy resource.
---

# alicloud_elasticsearch_snapshot_policy

Provides a Elasticsearch Snapshot Policy resource.

For information about Elasticsearch Snapshot Policy and how to use it, see [What is Snapshot Policy](https://next.api.alibabacloud.com/document/elasticsearch/2017-06-13/UpdateSnapshotSetting).

-> **NOTE:** Available since v2.1.1.

-> **NOTE:** The snapshot policy is bound to an Elasticsearch instance. Destroying this resource only sets `enable` to `false` and does not delete the instance or existing snapshot data.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-beijing"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.0.0.0/8"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.1.0.0/16"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = "cn-beijing-l"
}

resource "alicloud_elasticsearch_instance" "default" {
  description       = var.name
  vswitch_id        = alicloud_vswitch.default.id
  password          = "Examplw1234"
  version           = "7.10.0_with_X-Pack"
  payment_type      = "PayAsYouGo"
  instance_category = "x-pack"
  zone_count        = 1

  data_node_configuration {
    amount            = 2
    spec              = "elasticsearch.sn1ne.large.new"
    disk              = 20
    disk_type         = "cloud_essd"
    performance_level = "PL1"
  }

  kibana_configuration {
    spec   = "elasticsearch.n4.small"
    amount = 1
    disk   = 0
  }
}

resource "alicloud_elasticsearch_snapshot_policy" "default" {
  instance_id  = alicloud_elasticsearch_instance.default.id
  enable       = true
  quartz_regex = "0 0 01 ? * * *"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_elasticsearch_snapshot_policy&spm=docs.r.elasticsearch_snapshot_policy.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew, String) The ID of the Elasticsearch instance.
* `enable` - (Required, Bool) Whether to enable the automatic snapshot creation feature.
* `quartz_regex` - (Optional, String) The cron expression for the automatic snapshot creation. For example, `"0 0 01 ? * * *"` indicates that the snapshot is created at 01:00 every day. This parameter is required when `enable` is set to `true`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Snapshot Policy. It is the ID of the Elasticsearch instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 mins) Used when create the Snapshot Policy.
* `update` - (Defaults to 30 mins) Used when update the Snapshot Policy.

## Import

Elasticsearch Snapshot Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_elasticsearch_snapshot_policy.example <id>
```
