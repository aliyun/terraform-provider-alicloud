---
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
* `engine` - (Optional) Database type. Options are `Redis`, `Memcache`. If no value is specified, all types are returned.
* `engine_version` - (Optional) Database version required by the user. Value options of Redis can refer to the latest docs [detail info](https://www.alibabacloud.com/help/doc-detail/60873.htm) `EngineVersion`. Value of Memcache should be empty.
* `architecture` - (Optional) The KVStore instance system architecture required by the user. Valid values: `standard`, `cluster` and `rwsplit`.
* `performance_type` - (Optional) The KVStore instance performance type required by the user. Valid values: `standard_performance_type` and `enhance_performance_type`.
* `storage_type` - (Optional) The KVStore instance storage space required by the user. Valid values: `inmemory` and `hybrid`.
* `node_type` - (Optional) The KVStore instance node type required by the user. Valid values: `double`, `single`, `readone`, `readthree` and `readfive`.
* `package_type` - (Optional) The KVStore instance package type required by the user. Valid values: `standard` and `customized`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform apply`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `instance_classes` - A list of KVStore available instance classes.
    