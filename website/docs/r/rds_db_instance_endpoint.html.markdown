---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_rds_db_instance_endpoint"
sidebar_current: "docs-alicloud-resource-rds-db-instance-endpoint"
description: |-
  Provide RDS cluster instance endpoint connection resources.
---

# alicloud_rds_db_instance_endpoint

Provide RDS cluster instance endpoint connection resources, see [What is RDS DB Instance Endpoint](https://www.alibabacloud.com/help/en/apsaradb-for-rds/latest/api-rds-2014-08-15-createdbinstanceendpoint).

-> **NOTE:** Available since v1.203.0.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_rds_db_instance_endpoint&exampleId=cea73954-bc97-69b0-b51b-3e07e8ffb6367c055ff2&activeTab=example&spm=docs.r.rds_db_instance_endpoint.0.cea73954bc&intl_lang=EN_US" target="_blank">
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
```

## Argument Reference

The following arguments are supported:

* `db_instance_id` - (Required, ForceNew) The ID of the instance.
* `vpc_id` - (Required) The virtual private cloud (VPC) ID of the internal endpoint.
* `vswitch_id` - (Required) The vSwitch ID of the internal endpoint.
* `connection_string_prefix` - (Required) The IP address of the internal endpoint.
* `port` - (Required) The port number of the internal endpoint. You can specify the port number for the internal endpoint.Valid values: 3000 to 5999.
* `db_instance_endpoint_description` - (Optional) The user-defined description of the endpoint.
* `node_items` - (Required) The information about the node that is configured for the endpoint.  It contains two sub-fields(node_id and weight). See [`node_items`](#node_items) below.

### `node_items`

The node_items mapping supports the following:

* `node_id` - (Required) The ID of the node.
* `weight` - (Required) The weight of the node. Read requests are distributed based on the weight.Valid values: 0 to 100.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of endpoint.The value formats as `<db_instance_id>:<db_instance_endpoint_id>`.
* `private_ip_address` - The IP address of the internal endpoint.
* `connection_string` - The internal endpoint.
* `db_instance_endpoint_type` - The type of the endpoint.
* `ip_type` - The type of the IP address.
* `db_instance_endpoint_id` - The Endpoint ID of the instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 mins) Use when opening exclusive agent (until it reaches the initial `Running` status).
* `update` - (Defaults to 30 mins) Used when updating exclusive agent (until it reaches the initial `Running` status).
* `delete` - (Defaults to 20 mins) Use when closing exclusive agent.

## Import

RDS database endpoint feature can be imported using the id, e.g.

```shell
$ terraform import alicloud_rds_db_instance_endpoint.example <db_instance_id>:<db_instance_endpoint_id>
```
