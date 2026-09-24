---
subcategory: "Elasticsearch"
layout: "alicloud"
page_title: "Alicloud: alicloud_elasticsearch_logstash"
description: |-
  Provides a Alicloud Elasticsearch Logstash resource.
---

# alicloud_elasticsearch_logstash

Provides a Elasticsearch Logstash resource.

For information about Elasticsearch Logstash and how to use it, see [What is Logstash](https://next.api.alibabacloud.com/document/elasticsearch/2017-06-13/CreateLogstash).

-> **NOTE:** Available since v1.287.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_elasticsearch_zones" "default" {}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.0.0.0/8"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.1.0.0/16"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_elasticsearch_zones.default.zones.0.id
}

resource "alicloud_elasticsearch_logstash" "default" {
  description = var.name
  version     = "7.4_with_X-Pack"
  node_amount = 2

  node_spec {
    disk_type = "cloud_efficiency"
    spec      = "elasticsearch.sn1ne.large"
    disk      = 20
  }

  network_config {
    type       = "vpc"
    vpc_id     = alicloud_vpc.default.id
    vswitch_id = alicloud_vswitch.default.id
    vs_area    = data.alicloud_elasticsearch_zones.default.zones.0.id
  }

  payment_type = "PayAsYouGo"
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) The description of the Logstash instance.
* `network_config` - (Required, ForceNew) Network configuration of the Logstash instance. See [`network_config`](#network_config) below.
* `node_amount` - (Required) The number of data nodes.
* `node_spec` - (Required, ForceNew) The node specification. See [`node_spec`](#node_spec) below.
* `payment_info` - (Optional, ForceNew) Payment configuration for Subscription instances. See [`payment_info`](#payment_info) below.
* `payment_type` - (Optional, Computed, ForceNew) The payment type. Valid values: `PayAsYouGo`, `Subscription`.
* `resource_group_id` - (Optional, Computed, ForceNew) The ID of the resource group.
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `version` - (Required, ForceNew) The version of Logstash. Valid values: `7.4_with_X-Pack`, `6.7_with_X-Pack`.

### `network_config`

The network_config supports the following:

* `type` - (Required, ForceNew) The network type. Valid value: `vpc`.
* `vpc_id` - (Required, ForceNew) The ID of the VPC.
* `vswitch_id` - (Required, ForceNew) The ID of the VSwitch.
* `vs_area` - (Required, ForceNew) The zone ID of the VSwitch.

### `node_spec`

The node_spec supports the following:

* `disk_type` - (Optional, ForceNew) The disk type of the node. Valid values: `cloud_ssd`, `cloud_efficiency`.
* `spec` - (Required, ForceNew) The specification of the node.
* `disk` - (Optional, ForceNew) The disk size of the node.

### `payment_info`

The payment_info supports the following:

* `auto_renew` - (Optional, ForceNew) Whether to enable auto-renew.
* `auto_renew_duration` - (Optional, ForceNew) The duration of auto-renewal.
* `duration` - (Optional, ForceNew) The subscription duration.
* `pricing_cycle` - (Optional, ForceNew) The pricing cycle. Valid values: `Year`, `Month`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `create_time` - The creation time of the Logstash instance.
* `id` - The ID of the Logstash instance.
* `status` - The status of the Logstash instance.
* `updated_at` - The update time of the Logstash instance.

## Import

Elasticsearch Logstash can be imported using the instance id, e.g.

```shell
terraform import alicloud_elasticsearch_logstash.example <instance_id>
```
