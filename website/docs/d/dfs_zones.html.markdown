---
subcategory: "Apsara File Storage for HDFS"
layout: "alicloud"
page_title: "Alicloud: alicloud_dfs_zones"
sidebar_current: "docs-alicloud-datasource-dfs-zones"
description: |-
    Provides a list of DFS Zones And Configurations to the user.
---

# alicloud\_dfs\_zones

This data source provides the DFS Zones And Configurations of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.140.0+.

## Example Usage

Basic Usage

```terraform

data alicloud_dfs_zones "default" {}
```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `zones` - A list of HDFS Zones And Configurations. Each element contains the following attributes:
    * `zone_id` - The zone ID.
    * `options` -  A list of available configurations of the Zone.
      * `storage_type` - The storage specifications of the File system. Valid values: `PERFORMANCE`, `STANDARD`.
      * `protocol_type` - The protocol type. Valid values: `HDFS`.
