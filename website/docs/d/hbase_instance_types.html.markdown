---
subcategory: "HBase"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbase_instance_types"
sidebar_current: "docs-alicloud-datasource-hbase-instance-types"
description: |-
    Provides a list of availability instance_types for HBase that can be used by an Alibaba Cloud account.
---

# alicloud\_hbase\_instance\_types

This data source provides availability instance_types for HBase that can be accessed by an Alibaba Cloud account within the region configured in the provider.

-> **NOTE:** Available in v1.106.0+.

## Example Usage

```terraform
data "alicloud_hbase_instance_types" "default" {
  charge_type   = "Postpaid"
  region_id     = "cn-shanghai"
  zone_id       = "cn-shanghai-g"
  engine        = "hbaseue"
  version       = "2.0"
  instance_type = "hbase.sn2.large"
  disk_type     = "cloud_ssd"
}

resource "alicloud_hbase_instance" "hbase" {
  core_instance_type = data.alicloud_hbase_instance_types.default.types[0].id

  # Other properties...
}
```

## Argument Reference

The following arguments are supported:

* `instance_type` - (Optional) The hbase instance type of create hbase cluster instance.
* `charge_type` - (Optional, Available in 1.115.0+) The charge type of create hbase cluster instance, `PrePaid` or `PostPaid`.
* `region_id` - (Optional, Available in 1.115.0+) The dest region id, default client region.
* `zone_id` - (Optional, Available in 1.115.0+) The zone id, belong to regionId.
* `engine` - (Optional, Available in 1.115.0+) The engine name, `singlehbase`, `hbase`, `hbaseue`, `bds`.
* `version` - (Optional, Available in 1.115.0+) The engine version, singlehbase/hbase=1.1/2.0, bds=1.0.
* `disk_type` - (Optional, Available in 1.115.0+) The disk type, `cloud_ssd`, `cloud_essd_pl1`, `cloud_efficiency`, `local_hdd_pro`, `local_ssd_pro`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of instance types type IDs. 
* `types` - (Deprecated) A list of instance types. Each element contains the following attributes:
  * `value` - Name of the instance type.
  * `cpu_size` - Cpu size of the instance type.
  * `mem_size` - Mem size of the instance type.
* `master_instance_types` - (Available in 1.115.0+) A list of master instance types. Each element contains the following attributes:
    * `instance_type` - Name of the instance type.
    * `cpu_size` - Cpu size of the instance type.
    * `mem_size` - Mem size of the instance type.
* `core_instance_types` - (Available in 1.115.0+) A list of core instance types. Each element contains the following attributes:
    * `zone` - Name of zone id.
    * `engine` - Name of the engine.
    * `version` - The version of the engine.
    * `category` - Name of the category, single or cluster.
    * `storage_type` - Name of the storage type.
    * `instance_type` - Name of the instance type.
    * `instance_type` - Name of the instance type.
    * `cpu_size` - Cpu size of the instance type.
    * `mem_size` - Mem size of the instance type.
    * `max_core_count` - Max count of the core instance nodes.
