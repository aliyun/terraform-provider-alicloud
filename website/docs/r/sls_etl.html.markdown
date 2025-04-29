---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_sls_etl"
description: |-
  Provides a Alicloud Log Service (SLS) Etl resource.
---

# alicloud_sls_etl

Provides a Log Service (SLS) Etl resource.



For information about Log Service (SLS) Etl and how to use it, see [What is Etl](https://next.api.alibabacloud.com/document/Sls/2020-12-30/CreateETL).

-> **NOTE:** Available since v1.248.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_sls_etl&exampleId=18ea6489-61d7-eb1b-e1aa-c19d7081fa3db8bc64b8&activeTab=example&spm=docs.r.sls_etl.0.18ea648961&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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
```

## Argument Reference

The following arguments are supported:
* `configuration` - (Required, Set) The ETL configuration. See [`configuration`](#configuration) below.
* `description` - (Optional) Data Processing Task Description.
* `display_name` - (Required) Data processing task display name.
* `job_name` - (Required, ForceNew) Unique identification of data processing task.
* `project` - (Required, ForceNew) Project Name.

### `configuration`

The configuration supports the following:
* `from_time` - (Required, ForceNew, Int) The beginning of the time range for transformation.
* `lang` - (Required) Data processing syntax type.
* `logstore` - (Required, ForceNew) Source Logstore Name.
* `parameters` - (Optional, Map) Advanced parameter configuration.
* `role_arn` - (Required) The ARN role that authorizes reading of the source Logstore.
* `script` - (Required) Processing script.
* `sink` - (Required, Set) Processing result output target list See [`sink`](#configuration-sink) below.
* `to_time` - (Required, ForceNew, Int) The end of the time range for transformation.

### `configuration-sink`

The configuration-sink supports the following:
* `datasets` - (Required, List) Write Result Set.
* `endpoint` - (Required) The endpoint of the region where the target Project is located.
* `logstore` - (Required) Destination Logstore Name.
* `name` - (Required, ForceNew) Output Destination Name.
* `project` - (Required) Target Project name.
* `role_arn` - (Required) The ARN role that authorizes writing to the target Logstore.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<project>:<job_name>`.
* `create_time` - The time when the data transformation job was created.
* `status` - The status of the data transformation job.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Etl.
* `delete` - (Defaults to 5 mins) Used when delete the Etl.
* `update` - (Defaults to 5 mins) Used when update the Etl.

## Import

Log Service (SLS) Etl can be imported using the id, e.g.

```shell
$ terraform import alicloud_sls_etl.example <project>:<job_name>
```
