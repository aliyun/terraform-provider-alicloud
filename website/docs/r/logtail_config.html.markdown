---
layout: "alicloud"
page_title: "Alicloud: alicloud_logtail_config"
sidebar_current: "docs-alicloud-resource-logtail-config"
description: |-
  Provides a Alicloud logtail config resource.
---

# alicloud\_logtail\_config

The Logtail access service is a log collection agent provided by Log Service. 
You can use Logtail to collect logs from servers such as Alibaba Cloud Elastic
Compute Service (ECS) instances in real time in the Log Service console. [Refer to details](https://www.alibabacloud.com/help/doc-detail/29058.htm
)

## Example Usage

Basic Usage

```
resource "alicloud_log_project" "example"{
	name = "test-tf"
	description = "create by terraform"
}
resource "alicloud_log_store" "example"{
  	project = "${alicloud_log_project.example.name}"
  	name = "tf-test-logstore"
  	retention_period = 3650
  	shard_count = 3
  	auto_split = true
  	max_split_shard_count = 60
  	append_meta = true
}
resource "alicloud_logtail_config" "example"{
	project = "${alicloud_log_project.example.name}"
  	logstore = "${alicloud_log_store.example.name}"
  	input_type = "file"
  	log_sample = "test"
  	config_name = "tf-log-config"
	output_type = "LogService"
  	input_detail = "${file("config.json")}"
}
```
## Argument Reference

The following arguments are supported:

* `project` - (Required, ForceNew) The project name to the log store belongs.
* `logstore` - (Required, ForceNew) The log store name to the query index belongs.
* `input_type` - (Required) The input type. Currently only two types of files and plugin are supported.
* `log_sample` - （Optional）The log sample of the Logtail configuration. The log size cannot exceed 1,000 bytes.
* `config_name` - (Required, ForceNew) The Logtail configuration name, which is unique in the same project.
* `output_type` - (Required) The output type. Currently, only LogService is supported.
* `input_detail` - (Required) The logtail configure the required JSON files.([Refer to details](https://www.alibabacloud.com/help/doc-detail/29058.htm))

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the log store index. It formats of `<project>:<logstore>:<config_name>`.

## Import

Logtial config can be imported using the id, e.g.

```
$ terraform import alicloud_logtail_config.example tf-log:tf-log-store:tf-log-config
```