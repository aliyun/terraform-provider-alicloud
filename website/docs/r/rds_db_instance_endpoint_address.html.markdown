---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_rds_db_instance_endpoint_address"
sidebar_current: "docs-alicloud-resource-rds-db-instance-endpoint-address"
description: |-
  Provide RDS cluster instance endpoint public connection resources.
---

# alicloud_rds_db_instance_endpoint_address

Provide RDS cluster instance endpoint public connection resources.

Information about RDS MySQL cluster endpoint address and how to use it, see [What is RDS DB Instance Endpoint Address](https://www.alibabacloud.com/help/en/apsaradb-for-rds/latest/api-rds-2014-08-15-createdbinstanceendpointaddress).

-> **NOTE:** Available since v1.204.0.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_rds_db_instance_endpoint_address&exampleId=daebd126-5021-54be-f6ba-4f296fb3049d8cf98232&activeTab=example&spm=docs.r.rds_db_instance_endpoint_address.0.daebd12650&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-beijing"
}

variable "name" {
  default = "tf-example"
}
data "alicloud_db_zones" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_charge_type     = "PostPaid"
  category                 = "cluster"
  db_instance_storage_type = "cloud_essd"
}

data "alicloud_db_instance_classes" "default" {
  zone_id                  = data.alicloud_db_zones.default.ids.0
  engine                   = "MySQL"
  engine_version           = "8.0"
  category                 = "cluster"
  db_instance_storage_type = "cloud_essd"
  instance_charge_type     = "PostPaid"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_db_zones.default.ids.0
  vswitch_name = var.name
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_db_instance" "default" {
  engine                   = "MySQL"
  engine_version           = "8.0"
  instance_type            = data.alicloud_db_instance_classes.default.instance_classes.0.instance_class
  instance_storage         = data.alicloud_db_instance_classes.default.instance_classes.0.storage_range.min
  instance_charge_type     = "Postpaid"
  instance_name            = var.name
  vswitch_id               = alicloud_vswitch.default.id
  monitoring_period        = "60"
  db_instance_storage_type = "cloud_essd"
  security_group_ids       = [alicloud_security_group.default.id]
  zone_id                  = data.alicloud_db_zones.default.ids.0
  zone_id_slave_a          = data.alicloud_db_zones.default.ids.0
}

resource "alicloud_rds_db_node" "default" {
  db_instance_id = alicloud_db_instance.default.id
  class_code     = alicloud_db_instance.default.instance_type
  zone_id        = alicloud_vswitch.default.zone_id
}

resource "alicloud_rds_db_instance_endpoint" "default" {
  db_instance_id                   = alicloud_rds_db_node.default.db_instance_id
  vpc_id                           = alicloud_vpc.default.id
  vswitch_id                       = alicloud_db_instance.default.vswitch_id
  connection_string_prefix         = "example"
  port                             = "3306"
  db_instance_endpoint_description = var.name
  node_items {
    node_id = alicloud_rds_db_node.default.node_id
    weight  = 25
  }
}

resource "alicloud_rds_db_instance_endpoint_address" "default" {
  db_instance_id           = alicloud_db_instance.default.id
  db_instance_endpoint_id  = alicloud_rds_db_instance_endpoint.default.db_instance_endpoint_id
  connection_string_prefix = "tf-example-prefix"
  port                     = "3306"
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_id` - (Required, ForceNew) The ID of the instance.
* `db_instance_endpoint_id` - (Required, ForceNew) The Endpoint ID of the instance.
* `connection_string_prefix` - (Required) The prefix of the public endpoint.
* `port` - (Required) The port number of the public endpoint.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of endpoint public connection.The value formats as `<db_instance_id>:<db_instance_endpoint_id>`.
* `ip_address` - The IP address of the endpoint.
* `connection_string` - The endpoint of the instance.
* `ip_type` - The type of the IP address.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 mins) Use when opening exclusive agent (until it reaches the initial `Running` status).
* `update` - (Defaults to 30 mins) Used when updating exclusive agent (until it reaches the initial `Running` status).
* `delete` - (Defaults to 20 mins) Use when closing exclusive agent.

## Import

RDS database endpoint public address feature can be imported using the id, e.g.

```shell
$ terraform import alicloud_rds_db_instance_endpoint_address.example <db_instance_id>:<db_instance_endpoint_id>
```
