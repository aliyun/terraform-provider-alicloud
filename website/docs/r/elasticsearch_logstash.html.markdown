---
subcategory: "Elasticsearch"
layout: "alicloud"
page_title: "Alicloud: alicloud_elasticsearch_logstash"
description: |-
  Provides a Alicloud Elasticsearch Logstash resource.
---

# alicloud_elasticsearch_logstash

Provides a Elasticsearch Logstash resource. 

For information about Elasticsearch Logstash and how to use it, see [What is Logstash](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.208.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "defaultZFtcRh" {
  vpc_name   = var.name
  cidr_block = "10.0.0.0/8"
}

resource "alicloud_vswitch" "defaultMiMSn6" {
  vpc_id       = alicloud_vpc.defaultZFtcRh.id
  cidr_block   = "10.0.10.0/24"
  vswitch_name = "${var.name}1"
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_resource_manager_resource_group" "defaultWUcLGe" {
  display_name        = "rdktest05"
  resource_group_name = "${var.name}2"
}


resource "alicloud_elasticsearch_logstash" "default" {
  description       = "tf-acc-create-test-1"
  resource_group_id = alicloud_resource_manager_resource_group.defaultWUcLGe.id
  version           = "7.4_with_X-Pack"
  node_spec {
    disk_type = "cloud_efficiency"
    spec      = "elasticsearch.sn1ne.large"
    disk      = 20
  }
  network_config {
    type       = "vpc"
    vpc_id     = alicloud_vpc.defaultZFtcRh.id
    vswitch_id = alicloud_vswitch.defaultMiMSn6.id
    vs_area    = "cn-hangzhou-i"
  }
  payment_type = "PayAsYouGo"
  node_amount  = 1
  payment_info {
  }
}
```

## Argument Reference

The following arguments are supported:
* `description` - (Optional) Description.
* `network_config` - (Required, ForceNew) VPC configuration. See [`network_config`](#network_config) below.
* `node_amount` - (Required) Number of nodes.
* `node_spec` - (Required) Elasticsearch node disk. See [`node_spec`](#node_spec) below.
* `payment_info` - (Optional) Monthly payment renewal information. See [`payment_info`](#payment_info) below.
* `payment_type` - (Optional, ForceNew, Computed) The payment type of the resource.
* `resource_group_id` - (Optional, Computed) The ID of the resource group.
* `tags` - (Optional, Map) The tag of the resource.
* `version` - (Required, ForceNew) Version number.

### `network_config`

The network_config supports the following:
* `type` - (Required, ForceNew) Network type.
* `vswitch_id` - (Required, ForceNew) Unique identification of switch network.
* `vpc_id` - (Required, ForceNew) VPC unique identifier.
* `vs_area` - (Required, ForceNew) Unique ID of zone.

### `node_spec`

The node_spec supports the following:
* `disk` - (Optional) Disk size.
* `disk_type` - (Optional) Disk type.
* `spec` - (Required) Disk specification.

### `payment_info`

The payment_info supports the following:
* `auto_renew` - (Optional) Whether to enable automatic renewal.
* `auto_renew_duration` - (Optional) Automatic renewal cycle.
* `duration` - (Optional) Package year monthly purchase time.
* `pricing_cycle` - (Optional) Package year monthly unit.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - Instance creation time.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Logstash.
* `delete` - (Defaults to 5 mins) Used when delete the Logstash.
* `update` - (Defaults to 5 mins) Used when update the Logstash.

## Import

Elasticsearch Logstash can be imported using the id, e.g.

```shell
$ terraform import alicloud_elasticsearch_logstash.example <id>
```