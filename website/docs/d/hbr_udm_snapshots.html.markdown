---
subcategory: "Hybrid Backup Recovery (HBR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hbr_udm_snapshots"
sidebar_current: "docs-alicloud-datasource-hbr-udm-snapshots"
description: |-
  Provides a list of Hbr Udm Snapshot owned by an Alibaba Cloud account.
---

# alicloud_hbr_udm_snapshots

This data source provides Hbr Udm Snapshot available to the user.[What is Udm Snapshot](https://next.api.alibabacloud.com/document/hbr/2017-09-08/DescribeUdmSnapshots)

-> **NOTE:** Available since v1.253.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-beijing"
}

data "alicloud_hbr_udm_snapshots" "default" {
  source_type = "UDM_ECS"
  start_time  = "1642057551"
  end_time    = "1750927687"
  instance_id = "i-08qv5q4c4j****"
}

output "alicloud_hbr_udm_snapshot_example_id" {
  value = data.alicloud_hbr_udm_snapshots.default.snapshots.0.id
}
```

## Argument Reference

The following arguments are supported:
* `disk_id` - (ForceNew, Optional) Cloud disk ID. This field is valid only when SourceType = UDM_ECS_DISK.
* `end_time` - (Required, ForceNew) End Time
* `instance_id` - (Required, ForceNew) ECS instance ID
* `job_id` - (ForceNew, Optional) The ID of the backup job that creates the snapshot.
* `source_type` - (Required, ForceNew) Data source type. Only UDM_ECS and UDM_ECS_DISK are supported.
* `start_time` - (Required, ForceNew) Start Time
* `ids` - (Optional, ForceNew, Computed) A list of Udm Snapshot IDs. 
* `output_file` - (Optional, ForceNew) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Udm Snapshot IDs.
* `snapshots` - A list of Udm Snapshot Entries. Each element contains the following attributes:
  * `create_time` - The creation time of the resource
  * `disk_id` - Cloud disk ID. This field is valid only when SourceType = UDM_ECS_DISK.
  * `instance_id` - ECS instance ID
  * `job_id` - The ID of the backup job that creates the snapshot.
  * `source_type` - Data source type. Only UDM_ECS and UDM_ECS_DISK are supported.
  * `start_time` - Start Time
  * `udm_snapshot_id` - The first ID of the resource
  * `id` - The ID of the resource supplied above.
