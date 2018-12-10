---
layout: "alicloud"
page_title: "Alicloud: alicloud_kvstore_parameter"
sidebar_current: "docs-alicloud-resource-kvstore-parameter"
description: |-
  Provides an ApsaraDB Redis / Memcache parameter resource.
---

# alicloud\_kvstore\_parameter

Provides an ApsaraDB Redis / Memcache parameter resource. Documentation of the available parameters for ApsaraDB Redis / Memcache can be found at: [Instance configurations table](https://www.alibabacloud.com/help/doc-detail/61209.htm)

## Example Usage

```
resource "alicloud_kvstore_parameter" "compat" {
	instance_id = "rm-2eps..."
	name = "list-max-ziplist-entries"
	value = "256"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) The Id of instance in which parameter belongs.
* `name` - (Required) The name of the DB parameter.
* `value` - (Optional) The value of the DB parameter.

## Attributes Reference

The following attributes are exported:

* `instance_id` - The Id of instance in which parameter belongs.
* `name` - The name of the DB parameter.
* `value` - The value of the DB parameter.

## Import

KVStore instance can be imported using the id, e.g.

```
$ terraform import alicloud_kvstore_parameter.example r-abc123456:list-max-ziplist-entries
```
