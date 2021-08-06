---
subcategory: "Redis And Memcache (KVStore)"
layout: "alicloud"
page_title: "Alicloud: alicloud_kvstore_instance_classes"
sidebar_current: "docs-alicloud-datasource-kvstore-instance-classes"
description: |-
    Provides a list of KVStore instacne classes info.
---

# alicloud\_kvstore\_instances\_classes

This data source provides the KVStore instance classes resource available info of Alibaba Cloud.

-> **NOTE:** Available in v1.49.0+

## Example Usage

```tf
data "alicloud_zones" "resources" {
  available_resource_creation = "KVStore"
}

data "alicloud_kvstore_instance_classes" "resources" {
  zone_id              = "${data.alicloud_zones.resources.zones.0.id}"
  instance_charge_type = "PrePaid"
  engine               = "Redis"
  engine_version       = "5.0"
  output_file          = "./classes.txt"
}

output "first_kvstore_instance_class" {
  value = "${data.alicloud_kvstore_instance_classes.resources.instance_classes}"
}
```

## Argument Reference
    
The following arguments are supported:

* `zone_id` - (Required) The Zone to launch the KVStore instance.
* `instance_charge_type` - (Optional) Filter the results by charge type. Valid values: `PrePaid` and `PostPaid`. Default to `PrePaid`.
* `engine` - (Optional) Database type. Options are `Redis`, `Memcache`. Default to `Redis`.
* `engine_version` - (Optional) Database version required by the user. Value options of Redis can refer to the latest docs [detail info](https://www.alibabacloud.com/help/doc-detail/60873.htm) `EngineVersion`. Value of Memcache should be empty.
* `architecture` - (Optional) The KVStore instance system architecture required by the user. Valid values: `standard`, `cluster` and `rwsplit`.
* `performance_type` - (Optional, Deprecated) It has been deprecated from 1.68.0. 
* `storage_type` - (Optional, Deprecated) It has been deprecated from 1.68.0. 
* `node_type` - (Optional) The KVStore instance node type required by the user. Valid values: `double`, `single`, `readone`, `readthree` and `readfive`.
* `package_type` - (Optional, Deprecated) It has been deprecated from 1.68.0.
* `output_file` - (Optional) File name where to save data source results (after running `terraform apply`).
* `edition_type` - (Optional, Available in 1.68.0+) The KVStore instance edition type required by the user. Valid values: `Community` and `Enterprise`.
* `series_type` - (Optional, Available in 1.68.0+) The KVStore instance series type required by the user. Valid values: `enhanced_performance_type` and `hybrid_storage`.
* `shard_number` - (Optional, Available in 1.68.0+) The number of shard.Valid values: `1`, `2`, `4`, `8`, `16`, `32`, `64`, `128`, `256`.
* product_type - (Optional, Available in v1.130.0+) The type of the service. Valid values:
    * Local: an ApsaraDB for Redis instance with a local disk.
    * OnECS: an ApsaraDB for Redis instance with a standard disk. This type is available only on the Alibaba Cloud China site.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `instance_classes` - A list of KVStore available instance classes.
* `classes` - A list of KVStore available instance classes when the `sorted_by` is "Price". include:
  * `instance_class` - KVStore available instance class.
    
