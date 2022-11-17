---
subcategory: "E-MapReduce (EMR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_emr_disk_types"
sidebar_current: "docs-alicloud-datasource-emr-disk-types"
description: |-
    Provides a collection of data disk and system disk types when create emr cluster according to the specified filters.
---

# alicloud\_emr\_disk\_types

The `alicloud_emr_disk_types` data source provides a collection of data disk and 
system disk types available in Alibaba Cloud account when create a emr cluster.

-> **NOTE:** Available in 1.60.0+

## Example Usage

```
data "alicloud_emr_disk_types" "default" {
  destination_resource = "DataDisk"
  instance_charge_type = "PostPaid"
  cluster_type         = "HADOOP"
  instance_type        = "ecs.g5.xlarge"
  zone_id              = "cn-huhehaote-a"
}

output "data_disk_type" {
  value = "${data.alicloud_emr_disk_types.default.types.0.value}"
}
```

## Argument Reference

The following arguments are supported:

* `destination_resource` - (Required) The destination resource of emr cluster instance
* `instance_charge_type` - (Required) Filter the results by charge type. Valid values: `PrePaid` and `PostPaid`. Default to `PostPaid`.
* `cluster_type` - (Required) The cluster type of the emr cluster instance. Possible values: `HADOOP`, `KAFKA`, `ZOOKEEPER`, `DRUID`.
* `instance_type` - (Required) The ecs instance type of create emr cluster instance.
* `zone_id` - (Optional) The Zone to create emr cluster instance.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of data disk and system disk type IDs. 
* `types` - A list of emr instance types. Each element contains the following attributes:
  * `value` - The value of the data disk or system disk
  * `min` - The mininum value of the data disk to supported the specific instance type
  * `max` - The maximum value of the data disk to supported the specific instance type
