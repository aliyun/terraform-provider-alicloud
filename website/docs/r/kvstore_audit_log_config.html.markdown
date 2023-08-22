---
subcategory: "Redis And Memcache (KVStore)"
layout: "alicloud"
page_title: "Alicloud: alicloud_kvstore_audit_log_config"
sidebar_current: "docs-alicloud-resource-kvstore-audit-log-config"
description: |-
  Provides a Alicloud Redis And Memcache (KVStore) Audit Log Config resource.
---

# alicloud_kvstore_audit_log_config

Provides a Redis And Memcache (KVStore) Audit Log Config resource.

-> **NOTE:** Available since v1.130.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}
data "alicloud_kvstore_zones" "default" {

}
data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_kvstore_zones.default.zones.0.id
}

resource "alicloud_kvstore_instance" "default" {
  db_instance_name  = var.name
  vswitch_id        = alicloud_vswitch.default.id
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.ids.0
  zone_id           = data.alicloud_kvstore_zones.default.zones.0.id
  instance_class    = "redis.master.large.default"
  instance_type     = "Redis"
  engine_version    = "5.0"
  security_ips      = ["10.23.12.24"]
  config = {
    appendonly             = "yes"
    lazyfree-lazy-eviction = "yes"
  }
  tags = {
    Created = "TF",
    For     = "example",
  }
}

resource "alicloud_kvstore_audit_log_config" "example" {
  instance_id = alicloud_kvstore_instance.default.id
  db_audit    = true
  retention   = 1
}
```

## Argument Reference

The following arguments are supported:

* `db_audit` - (Optional) Indicates Whether to Enable the Audit Log.  Valid value: 
  * true: Default Value, Open. 
  * false: Closed. 
    
  Note: When the Instance for the Cluster Architecture Or Read/Write Split Architecture, at the Same Time to Open Or Close the Data Node and the Proxy Node of the Audit Log Doesn't Support Separate Open.
  
* `instance_id` - (Required, ForceNew) Instance ID, Call the Describeinstances Get.
* `retention` - (Optional) Audit Log Retention Period Value: 1~365. 
  
-> **NOTE**: When the Instance dbaudit Value Is Set to True, This Parameter Entry into Force. The Parameter Setting of the Current Region of All an Apsaradb for Redis Instance for a Data Entry into Force.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Audit Log Config. Its value is same as `instance_id`.
* `create_time` - Instance Creation Time.
* `status` - The status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Audit Log Config.
* `update` - (Defaults to 1 mins) Used when update the Audit Log Config.

## Import

Redis And Memcache (KVStore) Audit Log Config can be imported using the id, e.g.

```shell
$ terraform import alicloud_kvstore_audit_log_config.example <instance_id>
```
