---
subcategory: "Tair (Redis OSS-Compatible) And Memcache (KVStore)"
layout: "alicloud"
page_title: "Alicloud: alicloud_kvstore_permission"
sidebar_current: "docs-alicloud-datasource-kvstore-permission"
description: |-
  Provides data source to assigns a RAM role to Tair (Redis OSS-Compatible) And Memcache (KVStore).
---

# alicloud_kvstore_permission

Assigns a RAM role to Tair (Redis OSS-Compatible) And Memcache (KVStore).

The log management feature of Tair (Redis OSS-Compatible) And Memcache (KVStore) requires the resources of [Log Service](https://www.alibabacloud.com/help/doc-detail/48869.htm). 
To use the log management feature of Tair (Redis OSS-Compatible) And Memcache (KVStore), you can call this operation to associate the RAM role named AliyunServiceRoleForKvstore with the Tair (Redis OSS-Compatible) And Memcache (KVStore) instance. 
For more information, see [Associated RAM roles of Tair (Redis OSS-Compatible) And Memcache (KVStore)](https://www.alibabacloud.com/help/doc-detail/184337.htm)

-> **NOTE:** Available since v1.128.0

## Example Usage

```terraform
data "alicloud_kvstore_permission" "open" {
  enable = "On"
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Optional) Setting the value to `On` to assigns a RAM role to Tair (Redis OSS-Compatible) And Memcache (KVStore). If has been enabled, return the result. Valid values: "On" or "Off". Default to "Off".

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `status` - The current service enable status. 
