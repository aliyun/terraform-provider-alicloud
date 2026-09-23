---
subcategory: "Realtime Compute"
layout: "alicloud"
page_title: "Alicloud: alicloud_realtime_compute_sql_file"
description: |-
  Provides a Alicloud Realtime Compute Sql File resource.
---

# alicloud_realtime_compute_sql_file

Provides a Realtime Compute Sql File resource.

SQL query script file of Realtime Compute for Apache Flink.

For information about Realtime Compute Sql File and how to use it, see [What is Sql File](https://next.api.alibabacloud.com/document/ververica/2022-07-18/CreateSqlFile).

-> **NOTE:** Available since v1.292.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_realtime_compute_sql_file&exampleId=68a83acd-9ce2-142d-adbe-ff3aeda729394505d03d&activeTab=example&spm=docs.r.realtime_compute_sql_file.0.68a83acd9c&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tfexample-sql-file"
}

resource "alicloud_vpc" "default" {
  is_default = false
  cidr_block = "172.16.0.0/16"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "default" {
  is_default   = false
  vpc_id       = alicloud_vpc.default.id
  zone_id      = "cn-beijing-g"
  cidr_block   = "172.16.0.0/24"
  vswitch_name = var.name
}

resource "alicloud_oss_bucket" "default" {
}

resource "alicloud_realtime_compute_vvp_instance" "default" {
  vvp_instance_name = var.name
  storage {
    oss {
      bucket = alicloud_oss_bucket.default.id
    }
  }
  vpc_id      = alicloud_vpc.default.id
  vswitch_ids = [alicloud_vswitch.default.id]
  resource_spec {
    cpu       = "4"
    memory_gb = "16"
  }
  payment_type = "PayAsYouGo"
  zone_id      = alicloud_vswitch.default.zone_id
}

resource "alicloud_realtime_compute_sql_file" "default" {
  workspace            = alicloud_realtime_compute_vvp_instance.default.resource_id
  namespace            = "${alicloud_realtime_compute_vvp_instance.default.vvp_instance_name}-default"
  name                 = var.name
  sql_script           = "SELECT * FROM `vvp`.`default`.example_table;"
  session_cluster_name = "example-session-cluster"
  description          = "example sql file"
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_realtime_compute_sql_file&spm=docs.r.realtime_compute_sql_file.example&intl_lang=EN_US)


## Argument Reference

The following arguments are supported:
* `batch_mode` - (Optional, Computed) Whether the SQL query script runs in batch mode.
* `description` - (Optional) The description of the SQL file.
* `name` - (Required) The name of the SQL file.
* `namespace` - (Required, ForceNew) The name of the namespace.
* `parent_id` - (Optional, Computed) The ID of the parent folder of the SQL file. If not specified, the SQL file is created under the root folder of SQL scripts.
* `session_cluster_name` - (Required) The name of the session cluster that runs the SQL query script.
* `sql_script` - (Required) The SQL script content.
* `workspace` - (Optional, ForceNew, Computed) The ID of the workspace.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<workspace>:<namespace>:<sql_file_id>`.
* `sql_file_id` - The ID of the SQL file.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Sql File.
* `delete` - (Defaults to 5 mins) Used when delete the Sql File.
* `update` - (Defaults to 5 mins) Used when update the Sql File.

## Import

Realtime Compute Sql File can be imported using the id, e.g.

```shell
$ terraform import alicloud_realtime_compute_sql_file.example <workspace>:<namespace>:<sql_file_id>
```
