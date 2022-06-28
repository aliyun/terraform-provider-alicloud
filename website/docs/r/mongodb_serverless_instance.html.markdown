---
subcategory: "MongoDB"
layout: "alicloud"
page_title: "Alicloud: alicloud_mongodb_serverless_instance"
sidebar_current: "docs-alicloud-resource-mongodb-serverless-instance"
description: |-
  Provides a Alicloud MongoDB Serverless Instance resource.
---

# alicloud\_mongodb\_serverless\_instance

Provides a MongoDB Serverless Instance resource.

For information about MongoDB Serverless Instance and how to use it, see [What is Serverless Instance](https://www.alibabacloud.com/help/doc-detail/26558.html).

-> **NOTE:** Available in v1.148.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_mongodb_zones" "default" {}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_mongodb_zones.default.zones.0.id
}

data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_mongodb_serverless_instance" "example" {
  account_password        = "Abc12345"
  db_instance_description = "example_value"
  db_instance_storage     = 5
  storage_engine          = "WiredTiger"
  capacity_unit           = 100
  engine                  = "MongoDB"
  resource_group_id       = data.alicloud_resource_manager_resource_groups.default.groups.0.id
  engine_version          = "4.2"
  period                  = 1
  period_price_type       = "Month"
  vpc_id                  = data.alicloud_vpcs.default.ids.0
  zone_id                 = data.alicloud_mongodb_zones.default.zones.0.id
  vswitch_id              = data.alicloud_vswitches.default.ids.0
  tags = {
    Created = "MongodbServerlessInstance"
    For     = "TF"
  }
  security_ip_groups {
    security_ip_group_attribute = "example_value"
    security_ip_group_name      = "example_value"
    security_ip_list            = "192.168.0.1"
  }
}

```

## Argument Reference

The following arguments are supported:

* `account_password` - (Required, Sensitive) The password of the database logon account.
    * The password length is `8` to `32` bits.
    * The password consists of at least any three of uppercase letters, lowercase letters, numbers, and special characters. The special character is `!#$%^&*()_+-=`. The MongoDB Serverless instance provides a default database login account. This account cannot be modified. You can only set or modify the password for this account.
* `auto_renew` - (Optional) Set whether the instance is automatically renewed.
* `capacity_unit` - (Required) The I/O throughput consumed by the instance. Valid values: `100` to `8000`.
* `db_instance_description` - (Optional) The db instance description.
* `db_instance_storage` - (Required) The db instance storage. Valid values: `1` to `100`.
* `engine` - (Optional) The database engine of the instance. Valid values: `MongoDB`.
* `engine_version` - (Required) The database version number. Valid values: `4.2`.
* `maintain_end_time` - (Optional) The end time of the maintenance window. Specify the time in the `HH:mmZ` format. The time must be in UTC. **NOTE:** The difference between the start time and end time must be one hour. For example, if `maintain_start_time` is `01:00Z`, `maintain_end_time` must be `02:00Z`.
* `maintain_start_time` - (Optional) The start time of the maintenance window. Specify the time in the `HH:mmZ` format. The time must be in UTC.
* `period` - (Optional) The purchase duration of the instance, in months. Valid values: `1` to `9`, `12`, `24`, `36`, `60`.
* `period_price_type` - (Optional) The period price type. Valid values: `Day`, `Month`.
* `resource_group_id` - (Optional, ForceNew) The ID of the resource group.
* `security_ip_groups` - (Optional) An array that consists of the information of IP whitelists.
* `storage_engine` - (Optional) The storage engine used by the instance. Valid values: `WiredTiger`.
* `vswitch_id` - (Required, ForceNew) The of the vswitch.
* `zone_id` - (Required, ForceNew) The ID of the zone. Use this parameter to specify the zone created by the instance.
* `vpc_id` - (Required, ForceNew) The ID of the VPC network.
* `tags` - (Optional) A mapping of tags to assign to the resource.

#### Block security_ip_groups

The security_ip_groups supports the following:

* `security_ip_group_attribute` - (Optional) The attribute of the IP whitelist. This parameter is empty by default.
* `security_ip_group_name` - (Optional) The name of the IP whitelist.
* `security_ip_list` - (Optional) The IP addresses in the whitelist.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Serverless Instance.
* `status` - The instance status. For more information, see the instance Status Table.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when create the Serverless Instance.
* `update` - (Defaults to 10 mins) Used when update the Serverless Instance.

## Import

MongoDB Serverless Instance can be imported using the id, e.g.

```
$ terraform import alicloud_mongodb_serverless_instance.example <id>
```