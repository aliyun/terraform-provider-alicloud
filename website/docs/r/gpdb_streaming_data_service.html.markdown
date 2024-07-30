---
subcategory: "AnalyticDB for PostgreSQL (GPDB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_gpdb_streaming_data_service"
description: |-
  Provides a Alicloud GPDB Streaming Data Service resource.
---

# alicloud_gpdb_streaming_data_service

Provides a GPDB Streaming Data Service resource.



For information about GPDB Streaming Data Service and how to use it, see [What is Streaming Data Service](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.227.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-beijing"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "defaultTXZPBL" {
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "defaultrJ5mmz" {
  vpc_id     = alicloud_vpc.defaultTXZPBL.id
  zone_id    = "cn-beijing-h"
  cidr_block = "192.168.1.0/24"
}

resource "alicloud_gpdb_instance" "default1oSPzX" {
  instance_spec         = "2C8G"
  description           = var.name
  seg_node_num          = "2"
  seg_storage_type      = "cloud_essd"
  instance_network_type = "VPC"
  db_instance_category  = "Basic"
  payment_type          = "PayAsYouGo"
  ssl_enabled           = "0"
  engine_version        = "6.0"
  zone_id               = "cn-beijing-h"
  vswitch_id            = alicloud_vswitch.defaultrJ5mmz.id
  storage_size          = "50"
  master_cu             = "4"
  vpc_id                = alicloud_vpc.defaultTXZPBL.id
  db_instance_mode      = "StorageElastic"
  engine                = "gpdb"
}


resource "alicloud_gpdb_streaming_data_service" "default" {
  service_name        = "example"
  db_instance_id      = alicloud_gpdb_instance.default1oSPzX.id
  service_description = "example"
  service_spec        = "8"
}
```

## Argument Reference

The following arguments are supported:
* `db_instance_id` - (Required, ForceNew) The ID of the associated instance.
* `service_description` - (Optional) The creation time of the resource
* `service_name` - (Required, ForceNew) Service Name
* `service_spec` - (Required) Resource Specifications

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<db_instance_id>:<service_id>`.
* `create_time` - Create time
* `service_id` - Service ID
* `status` - The status of the resource

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Streaming Data Service.
* `delete` - (Defaults to 5 mins) Used when delete the Streaming Data Service.
* `update` - (Defaults to 5 mins) Used when update the Streaming Data Service.

## Import

GPDB Streaming Data Service can be imported using the id, e.g.

```shell
$ terraform import alicloud_gpdb_streaming_data_service.example <db_instance_id>:<service_id>
```