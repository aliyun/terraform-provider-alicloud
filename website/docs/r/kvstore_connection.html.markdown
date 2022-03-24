---
subcategory: "Redis And Memcache (KVStore)"
layout: "alicloud"
page_title: "Alicloud: alicloud_kvstore_connection"
sidebar_current: "docs-alicloud-resource-kvstore-connection"
description: |-
  Operate the public network ip of the specified resource.
---

# alicloud\_kvstore\_connection

Operate the public network ip of the specified resource. How to use it, see [What is Resource Alicloud KVStore Connection](https://www.alibabacloud.com/help/doc-detail/125795.htm).

-> **NOTE:** Available in v1.101.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_kvstore_connection" "default" {
  connection_string_prefix = "allocatetestupdate"
  instance_id              = "r-abc123456"
  port                     = "6370"
}
```

## Argument Reference

The following arguments are supported:
* `connection_string_prefix` - (Required) The prefix of the public endpoint. The prefix can be 8 to 64 characters in length, and can contain lowercase letters and digits. It must start with a lowercase letter.
* `instance_id`- (Required) The ID of the instance.
* `port` - (Required) The service port number of the instance.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of KVStore DBInstance.
* `connection_string` - The public connection string of KVStore DBInstance.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when creating the KVStore connection (until it reaches the initial `Normal` status). 
* `update` - (Defaults to 5 mins) Used when updating the KVStore connection (until it reaches the initial `Normal` status). 
* `delete` - (Defaults to 5 mins) Used when deleting the KVStore connection (until it reaches the initial `Normal` status). 

## Import

KVStore connection can be imported using the id, e.g.

```
$ terraform import alicloud_kvstore_connection.example r-abc12345678
```

