---
subcategory: "Realtime Compute for Apache Flink"
layout: "alicloud"
page_title: "Alicloud: alicloud_realtime_compute_sql_file"
sidebar_current: "docs-alicloud-resource-realtime-compute-sql-file"
description: |-
  Provides a Realtime Compute for Apache Flink SQL File resource.
---

# alicloud_realtime_compute_sql_file

Provides a Realtime Compute for Apache Flink SQL File resource.

For information about Realtime Compute for Apache Flink and how to use it, refer to [Realtime Compute for Apache Flink](https://www.alibabacloud.com/help/en/flink/developer-reference/api-ververica-2022-07-18).

-> **NOTE:** Available since v1.292.0.

## Example Usage

```terraform
variable "name" {
  default = "tf-example"
}

resource "alicloud_realtime_compute_vvp_instance" "default" {
  vvp_instance_name = var.name
  storage {
    oss {
      bucket = "your-oss-bucket"
    }
  }
  vpc_id      = "your-vpc-id"
  vswitch_ids = ["your-vswitch-id"]
  resource_spec {
    cpu       = "4"
    memory_gb = "16"
  }
  payment_type = "PayAsYouGo"
  zone_id      = "cn-beijing-g"
}

resource "alicloud_realtime_compute_sql_file" "default" {
  workspace   = alicloud_realtime_compute_vvp_instance.default.resource_id
  namespace   = "${alicloud_realtime_compute_vvp_instance.default.vvp_instance_name}-default"
  name        = var.name
  sql_script  = "CREATE TABLE datagen (id VARCHAR, name VARCHAR) WITH ('connector' = 'datagen'); INSERT INTO blackhole SELECT * FROM datagen;"
  parent_id   = "0"
  description = "This is a test sql file."
}
```

## Argument Reference

The following arguments are supported:

* `workspace` - (Optional, Computed, ForceNew) The ID of the workspace. The workspace is derived from the VVP instance if not specified.
* `namespace` - (Required, ForceNew) The name of the namespace. The namespace is part of the project space within the workspace.
* `name` - (Required) The name of the SQL file.
* `sql_script` - (Required) The SQL script content of the SQL file.
* `batch_mode` - (Optional) The batch mode of the SQL file.
* `session_cluster_name` - (Optional) The name of the session cluster associated with the SQL file.
* `parent_id` - (Required) The parent folder ID of the SQL file. Set to `0` for root-level files.
* `description` - (Optional) The description of the SQL file.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of the SQL file. It is formatted to `<workspace>:<namespace>:<sql_file_id>`.
* `sql_file_id` - The ID of the SQL file.
* `region_id` - The region ID of the resource.

## Import

Realtime Compute SQL File can be imported using the id, e.g.

```shell
terraform import alicloud_realtime_compute_sql_file.example <workspace>:<namespace>:<sql_file_id>
```
