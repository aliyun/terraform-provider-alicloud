---
subcategory: "AnalyticDB for PostgreSQL (GPDB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_gpdb_streaming_data_source"
description: |-
  Provides a Alicloud GPDB Streaming Data Source resource.
---

# alicloud_gpdb_streaming_data_source

Provides a GPDB Streaming Data Source resource.

Real-time data source.

For information about GPDB Streaming Data Source and how to use it, see [What is Streaming Data Source](https://www.alibabacloud.com/help/en/analyticdb/analyticdb-for-postgresql/developer-reference/api-gpdb-2016-05-03-createstreamingdatasource).

-> **NOTE:** Available since v1.227.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_gpdb_streaming_data_source&exampleId=4237e505-f0f6-3076-c73f-0946959b47f2010c8412&activeTab=example&spm=docs.r.gpdb_streaming_data_source.0.4237e505f0&intl_lang=EN_US" target="_blank">
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

variable "kafka-config-modify" {
  default = <<EOF
{
    "brokers": "alikafka-post-cn-g4t3t4eod004-1-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-g4t3t4eod004-2-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-g4t3t4eod004-3-vpc.alikafka.aliyuncs.com:9092",
    "delimiter": "#",
    "format": "delimited",
    "topic": "ziyuan_example"
}
EOF
}

variable "kafka-config" {
  default = <<EOF
{
    "brokers": "alikafka-post-cn-g4t3t4eod004-1-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-g4t3t4eod004-2-vpc.alikafka.aliyuncs.com:9092,alikafka-post-cn-g4t3t4eod004-3-vpc.alikafka.aliyuncs.com:9092",
    "delimiter": "|",
    "format": "delimited",
    "topic": "ziyuan_example"
}
EOF
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "defaultDfkYOR" {
  cidr_block = "192.168.0.0/16"
}

resource "alicloud_vswitch" "default59ZqyD" {
  vpc_id     = alicloud_vpc.defaultDfkYOR.id
  zone_id    = "cn-beijing-h"
  cidr_block = "192.168.1.0/24"
}

resource "alicloud_gpdb_instance" "default7mX6ld" {
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
  vswitch_id            = alicloud_vswitch.default59ZqyD.id
  storage_size          = "50"
  master_cu             = "4"
  vpc_id                = alicloud_vpc.defaultDfkYOR.id
  db_instance_mode      = "StorageElastic"
  engine                = "gpdb"
}

resource "alicloud_gpdb_streaming_data_service" "defaultwruvdv" {
  service_name        = "example"
  db_instance_id      = alicloud_gpdb_instance.default7mX6ld.id
  service_description = "example"
  service_spec        = "8"
}


resource "alicloud_gpdb_streaming_data_source" "default" {
  db_instance_id          = alicloud_gpdb_instance.default7mX6ld.id
  data_source_name        = "example-kafka3"
  data_source_config      = var.kafka-config
  data_source_type        = "kafka"
  data_source_description = "example-kafka"
  service_id              = alicloud_gpdb_streaming_data_service.defaultwruvdv.service_id
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_gpdb_streaming_data_source&spm=docs.r.gpdb_streaming_data_source.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `db_instance_id` - (Required, ForceNew) The instance ID.

-> **NOTE:**   You can call the [DescribeDBInstances](https://www.alibabacloud.com/help/en/doc-detail/196830.html) operation to query the information about all AnalyticDB for PostgreSQL instances within a region, including instance IDs.

* `data_source_config` - (Required) The configurations of the data source. 
* `data_source_description` - (Optional) The description of the data source. 
* `data_source_name` - (Required, ForceNew) Data Source Name
* `data_source_type` - (Required, ForceNew) Data Source Type
* `service_id` - (Required, ForceNew) The real-time data service ID. 

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<db_instance_id>:<data_source_id>`.
* `create_time` - Creation time
* `data_source_id` - The data source ID. 
* `status` - Service Status:

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Streaming Data Source.
* `delete` - (Defaults to 5 mins) Used when delete the Streaming Data Source.
* `update` - (Defaults to 5 mins) Used when update the Streaming Data Source.

## Import

GPDB Streaming Data Source can be imported using the id, e.g.

```shell
$ terraform import alicloud_gpdb_streaming_data_source.example <db_instance_id>:<data_source_id>
```