---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sls_etls"
sidebar_current: "docs-alicloud-datasource-sls-etls"
description: |-
  Provides a list of Sls Etl owned by an Alibaba Cloud account.
---

# alicloud_sls_etls

This data source provides Sls Etl available to the user.[What is Etl](https://next.api.alibabacloud.com/document/Sls/2020-12-30/CreateETL)

-> **NOTE:** Available since v1.258.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_log_project" "defaulthhAPo6" {
  description  = "terraform-etl-example-813"
  project_name = "terraform-etl-example-330"
}

resource "alicloud_log_store" "defaultzWKLkp" {
  hot_ttl          = "8"
  retention_period = "30"
  shard_count      = "2"
  project_name     = alicloud_log_project.defaulthhAPo6.id
  logstore_name    = "example"
}

resource "alicloud_sls_etl" "default" {
  project     = alicloud_log_project.defaulthhAPo6.id
  description = "etl-1740472705-185721"
  configuration {
    script   = "* | extend a=1"
    lang     = "SPL"
    role_arn = var.name
    sink {
      name     = "11111"
      endpoint = "cn-hangzhou-intranet.log.aliyuncs.com"
      project  = "gy-hangzhou-huolang-1"
      logstore = "gy-rm2"
      datasets = ["__UNNAMED__"]
      role_arn = var.name
    }
    logstore  = alicloud_log_store.defaultzWKLkp.logstore_name
    from_time = "1706771697"
    to_time   = "1738394097"
  }
  job_name     = "etl-1740472705-185721"
  display_name = "etl-1740472705-185721"
}

data "alicloud_sls_etls" "default" {
  logstore = alicloud_log_store.defaultzWKLkp.name
  project  = alicloud_log_project.defaulthhAPo6.id
}

output "alicloud_sls_etl_example_id" {
  value = data.alicloud_sls_etls.default.etls.0.id
}
```

## Argument Reference

The following arguments are supported:
* `project` - (ForceNew, Required) Project Name
* `logstore` - (ForceNew, Required) Source Logstore Name.
* `offset` - (ForceNew, Optional) Query start row. The default value is 0.
* `size` - (ForceNew, Optional) Specify the number of data processing tasks returned by the query
* `ids` - (Optional, ForceNew, Computed) A list of Etl IDs. The value is formulated as `<project>:<job_name>`.
* `output_file` - (Optional, ForceNew) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Etl IDs.
* `etls` - A list of Etl Entries. Each element contains the following attributes:
  * `configuration` - Detailed configuration of data processing tasks
    * `from_time` - Processing time start timestamp (accurate to the second). Enter 0 when the first log received from the source Logstore is consumed.
    * `lang` - Data processing syntax type.
    * `logstore` - Source Logstore Name.
    * `parameters` - Advanced parameter configuration.
    * `role_arn` - The ARN role that authorizes reading of the source Logstore.
    * `script` - Processing script.
    * `sink` - Processing result output target list.
      * `datasets` - Write Result Set.
      * `endpoint` - The endpoint of the region where the target Project is located.
      * `logstore` - Destination Logstore Name.
      * `name` - Output Destination Name.
      * `project` - Target Project name.
      * `role_arn` - The ARN role that authorizes writing to the target Logstore.
    * `to_time` - Processing time end timestamp (accurate to seconds). When continuous consumption is stopped manually, fill in 0.
  * `create_time` - Task creation time. Example value: 1718787534
  * `description` - Data Processing Task Description
  * `display_name` - Data processing task display name
  * `job_name` - Unique identification of data processing task
  * `last_modified_time` - The time when the task was last modified. Example value: 1718787681
  * `schedule_id` - The task ID. Example values:
  * `status` - Task Status
  * `id` - The ID of the resource supplied above.
