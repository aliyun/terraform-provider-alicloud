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

```
# Declare the data source
data "alicloud_hbase_instance_types" "default" {}

# Create an HBase instance with the first matched type
resource "alicloud_hbase_instance" "hbase" {
    core_instance_type = data.alicloud_hbase_instance_types.default.types[0].id

  # Other properties...
}
```

## Argument Reference

The following arguments are supported:

* `instance_type` - (Optional) The hbase instance type of create hbase cluster instance.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of instance types type IDs. 
* `types` - A list of instance types. Each element contains the following attributes:
  * `value` - Name of the instance type.
  * `cpu_size` - Cpu size of the instance type.
  * `mem_size` - Mem size of the instance type.
  