---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_rds_db_proxy_public"
description: |-
  Provide RDS instance proxy external network connection resources.
---

# alicloud_rds_db_proxy_public

Provides a RDS database proxy public network address resource.



For information about Resource AlicloudRdsDBProxyPublic and how to use it, see [What is proxy](https://www.alibabacloud.com/help/en/rds/developer-reference/api-rds-2014-08-15-createdbproxyendpointaddress).

-> **NOTE:** Available since v1.250.0.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_rds_db_proxy_public&exampleId=dd8413b6-00fe-a3e4-66f6-5c0ac2abc03335013f07&activeTab=example&spm=docs.r.rds_db_proxy_public.0.dd8413b600&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}

data "alicloud_db_zones" "default" {
  engine         = "MySQL"
  engine_version = "5.6"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_db_zones.default.zones.0.id
  vswitch_name = var.name
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_db_instance" "default" {
  engine                   = "MySQL"
  engine_version           = "5.7"
  instance_type            = "rds.mysql.c1.large"
  instance_storage         = "20"
  instance_charge_type     = "Postpaid"
  instance_name            = var.name
  vswitch_id               = alicloud_vswitch.default.id
  db_instance_storage_type = "local_ssd"
}

resource "alicloud_db_readonly_instance" "default" {
  zone_id               = alicloud_db_instance.default.zone_id
  master_db_instance_id = alicloud_db_instance.default.id
  engine_version        = alicloud_db_instance.default.engine_version
  instance_storage      = alicloud_db_instance.default.instance_storage
  instance_type         = alicloud_db_instance.default.instance_type
  instance_name         = "${var.name}readonly"
  vswitch_id            = alicloud_vswitch.default.id
}

resource "alicloud_rds_db_proxy" "default" {
  instance_id                          = alicloud_db_instance.default.id
  instance_network_type                = "VPC"
  vpc_id                               = alicloud_db_instance.default.vpc_id
  vswitch_id                           = alicloud_db_instance.default.vswitch_id
  db_proxy_instance_num                = 2
  db_proxy_connection_prefix           = "example"
  db_proxy_connect_string_port         = 3306
  db_proxy_endpoint_read_write_mode    = "ReadWrite"
  read_only_instance_max_delay_time    = 90
  db_proxy_features                    = "TransactionReadSqlRouteOptimizeStatus:1;ConnectionPersist:1;ReadWriteSpliting:1"
  read_only_instance_distribution_type = "Custom"

  read_only_instance_weight {
    instance_id = alicloud_db_instance.default.id
    weight      = "100"
  }

  read_only_instance_weight {
    instance_id = alicloud_db_readonly_instance.default.id
    weight      = "500"
  }
}

resource "alicloud_rds_db_proxy_public" "default" {
  db_instance_id                      = alicloud_db_instance.default.id
  db_proxy_endpoint_id                = alicloud_rds_db_proxy.default.db_proxy_endpoint_id
  connection_string_prefix            = "exampleabc"
  db_proxy_connection_string_net_type = "Public"
  db_proxy_new_connect_string_port    = "3306"
}

```

-> **NOTE:** Resource `alicloud_rds_db_proxy_public` should be created after `alicloud_rds_db_proxy`, so the `depends_on` statement is necessary.

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_rds_db_proxy_public&spm=docs.r.rds_db_proxy_public.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `db_instance_id` - (Required, ForceNew) Instance ID.
* `db_proxy_endpoint_id` - (Required, ForceNew)Database proxy connection address ID.
* `connection_string_prefix` - (Required)The prefix for the new database proxy connection address can be customized.
* `db_proxy_connection_string_net_type` - (Required, ForceNew)The network type of the new database proxy connection address,This resource defaults to Public.
* `db_proxy_new_connect_string_port` - (Not required, Optional)The port for the new database proxy connection address is 3306 by default for MySQL and 5432 by default for PostgreSQL, which can be customized.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 60 mins) Used when create the Proxy Public Network Address.
* `delete` - (Defaults to 20 mins) Used when delete the Proxy Public Network Address.
* `update` - (Defaults to 30 mins) Used when update the Proxy Public Network Address.

## Import

RDS Database Proxy Public Network Address can be imported using the id, e.g.

```shell
$ terraform import alicloud_rds_db_proxy_public.example <id>
```