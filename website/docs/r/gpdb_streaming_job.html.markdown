---
subcategory: "AnalyticDB for PostgreSQL (GPDB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_gpdb_streaming_job"
description: |-
  Provides a Alicloud GPDB Streaming Job resource.
---

# alicloud_gpdb_streaming_job

Provides a GPDB Streaming Job resource.

Real-time data tasks.

For information about GPDB Streaming Job and how to use it, see [What is Streaming Job](https://www.alibabacloud.com/help/en/analyticdb/analyticdb-for-postgresql/developer-reference/api-gpdb-2016-05-03-createstreamingjob).

-> **NOTE:** Available since v1.231.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_gpdb_streaming_job&exampleId=e48965b0-2aeb-f14a-4f57-371f9c9a080fb9a5cfe0&activeTab=example&spm=docs.r.gpdb_streaming_job.0.e48965b02a&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-beijing"
}

resource "alicloud_vpc" "defaultTXqb15" {
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "defaultaSWhbT" {
  vpc_id     = alicloud_vpc.defaultTXqb15.id
  zone_id    = "cn-beijing-h"
  cidr_block = "192.168.1.0/24"
}

resource "alicloud_gpdb_instance" "defaulth2ghc1" {
  instance_spec         = "2C8G"
  description           = var.name
  seg_node_num          = "2"
  seg_storage_type      = "cloud_essd"
  instance_network_type = "VPC"
  db_instance_category  = "Basic"
  payment_type          = "PayAsYouGo"
  ssl_enabled           = "0"
  engine_version        = "6.0"
  zone_id               = "cn-beijing-h"
  vswitch_id            = alicloud_vswitch.defaultaSWhbT.id
  storage_size          = "50"
  master_cu             = "4"
  vpc_id                = alicloud_vpc.defaultTXqb15.id
  db_instance_mode      = "StorageElastic"
  engine                = "gpdb"
}

resource "alicloud_gpdb_streaming_data_service" "default2dUszY" {
  service_name        = "example"
  db_instance_id      = alicloud_gpdb_instance.defaulth2ghc1.id
  service_description = "example"
  service_spec        = "8"
}

resource "alicloud_gpdb_streaming_data_source" "defaultcDQItu" {
  db_instance_id          = alicloud_gpdb_instance.defaulth2ghc1.id
  data_source_name        = "example"
  data_source_config      = jsonencode({ "brokers" : "alikafka-post-cn-g4t3t4eod004-1-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-g4t3t4eod004-2-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-g4t3t4eod004-3-vpc.alikafka.aliyuncs.com:9092", "delimiter" : "|", "format" : "delimited", "topic" : "ziyuan_example" })
  data_source_type        = "kafka"
  data_source_description = "example"
  service_id              = alicloud_gpdb_streaming_data_service.default2dUszY.service_id
}


resource "alicloud_gpdb_streaming_job" "default" {
  account         = "example_001"
  dest_schema     = "public"
  mode            = "professional"
  job_name        = "example-kafka"
  job_description = "example-kafka"
  dest_database   = "adb_sampledata_tpch"
  db_instance_id  = alicloud_gpdb_instance.defaulth2ghc1.id
  dest_table      = "customer"
  data_source_id  = alicloud_gpdb_streaming_data_source.defaultcDQItu.data_source_id
  password        = "example_001"
  job_config      = <<EOF
ATABASE: adb_sampledata_tpch
USER: example_001
PASSWORD: example_001
HOST: gp-2zean69451zsjj139-master.gpdb.rds.aliyuncs.com
PORT: 5432
KAFKA:
  INPUT:
    SOURCE:
      BROKERS: alikafka-post-cn-3mp3t4ekq004-1-vpc.alikafka.aliyuncs.com:9092
      TOPIC: ziyuan_example
      FALLBACK_OFFSET: LATEST
    KEY:
      COLUMNS:
      - NAME: c_custkey
        TYPE: int
      FORMAT: delimited
      DELIMITED_OPTION:
        DELIMITER: \'|\'
    VALUE:
      COLUMNS:
      - NAME: c_comment
        TYPE: varchar
      FORMAT: delimited
      DELIMITED_OPTION:
        DELIMITER: \'|\'
    ERROR_LIMIT: 10
  OUTPUT:
    SCHEMA: public
    TABLE: customer
    MODE: MERGE
    MATCH_COLUMNS:
    - c_custkey
    ORDER_COLUMNS:
    - c_custkey
    UPDATE_COLUMNS:
    - c_custkey
    MAPPING:
    - NAME: c_custkey
      EXPRESSION: c_custkey
  COMMIT:
    MAX_ROW: 1000
    MINIMAL_INTERVAL: 1000
    CONSISTENCY: ATLEAST
  POLL:
    BATCHSIZE: 1000
    TIMEOUT: 1000
  PROPERTIES:
    group.id: ziyuan_example_01
EOF
}
```

## Argument Reference

The following arguments are supported:
* `account` - (Optional) The name of the database account.
* `consistency` - (Optional) The delivery guarantee setting.

  Valid values:

  - ATLEAST
  - EXACTLY
* `db_instance_id` - (Required, ForceNew) The instance ID.
* `data_source_id` - (Required, ForceNew) The data source ID.
* `dest_columns` - (Optional, List) Target Field
* `dest_database` - (Optional) The name of the destination database.
* `dest_schema` - (Optional) Target Schema
* `dest_table` - (Optional) The name of the destination table.
* `error_limit_count` - (Optional, Int) The number of allowed error rows. Write failures occur when Kafka data does not match the destination table in AnalyticDB for PostgreSQL. If the specified value is exceeded, the job fails.
* `fallback_offset` - (Optional) Automatic offset reset
* `group_name` - (Optional) Group Name
* `job_config` - (Optional) The YAML configuration file of the job. This parameter must be specified when Mode is set to professional.
* `job_description` - (Optional) The description of the job.
* `job_name` - (Required, ForceNew) The name of the job.
* `match_columns` - (Optional, List) Match Field
* `mode` - (Optional, ForceNew) The configuration mode. Valid values:

  1.  basic: In basic mode, you must configure the configuration parameters.

  2.  professional: In professional mode, you can submit a YAML configuration file.
* `password` - (Optional) The password of the database account.
* `src_columns` - (Optional, List) Source Field
* `try_run` - (Optional) Specifies whether to test the real-time job. Valid values:

  - true
  - false

  Default value: false.
* `update_columns` - (Optional, List) Update Field
* `write_mode` - (Optional) The write mode.

  Valid values:

  - insert
  - update
  - merge

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<db_instance_id>:<job_id>`.
* `create_time` - The creation time of the resource
* `job_id` - The job ID.
* `status` - Service status, value:

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Streaming Job.
* `delete` - (Defaults to 5 mins) Used when delete the Streaming Job.
* `update` - (Defaults to 5 mins) Used when update the Streaming Job.

## Import

GPDB Streaming Job can be imported using the id, e.g.

```shell
$ terraform import alicloud_gpdb_streaming_job.example <db_instance_id>:<job_id>
```