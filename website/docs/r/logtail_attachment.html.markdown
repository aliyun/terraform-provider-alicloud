---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_logtail_attachment"
sidebar_current: "docs-alicloud-resource-logtail-attachment"
description: |-
  Provides a Alicloud logtail attachment resource.
---

# alicloud\_logtail\_attachment

The Logtail access service is a log collection agent provided by Log Service.
You can use Logtail to collect logs from servers such as Alibaba Cloud Elastic
Compute Service (ECS) instances in real time in the Log Service console. [Refer to details](https://www.alibabacloud.com/help/doc-detail/29058.htm)

This resource amis to attach one logtail configure to a machine group.

-> **NOTE:** One logtail configure can be attached to multiple machine groups and one machine group can attach several logtail configures.

## Example Usage

Basic Usage

```terraform
resource "alicloud_log_project" "test" {
  name        = "test-tf2"
  description = "create by terraform"
}

resource "alicloud_log_store" "test" {
  project               = alicloud_log_project.test.name
  name                  = "tf-test-logstore"
  retention_period      = 3650
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}

resource "alicloud_log_machine_group" "test" {
  project       = alicloud_log_project.test.name
  name          = "tf-log-machine-group"
  topic         = "terraform"
  identify_list = ["10.0.0.1", "10.0.0.3", "10.0.0.2"]
}

resource "alicloud_logtail_config" "test" {
  project      = alicloud_log_project.test.name
  logstore     = alicloud_log_store.test.name
  input_type   = "file"
  log_sample   = "test"
  name         = "tf-log-config"
  output_type  = "LogService"
  input_detail = <<DEFINITION
  	{
		"logPath": "/logPath",
		"filePattern": "access.log",
		"logType": "json_log",
		"topicFormat": "default",
		"discardUnmatch": false,
		"enableRawLog": true,
		"fileEncoding": "gbk",
		"maxDepth": 10
	}
	
DEFINITION

}

resource "alicloud_logtail_attachment" "test" {
  project             = alicloud_log_project.test.name
  logtail_config_name = alicloud_logtail_config.test.name
  machine_group_name  = alicloud_log_machine_group.test.name
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required, ForceNew) The project name to the log store belongs.
* `logtail_config_name` - (Required, ForceNew) The Logtail configuration name, which is unique in the same project.
* `machine_group_name` - (Required, ForceNew) The machine group name, which is unique in the same project.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the logtail to machine group. It formats of `<project>:<logtail_config_name>:<machine_group_name>`.

## Import

Logtial to machine group can be imported using the id, e.g.

```shell
$ terraform import alicloud_logtail_to_machine_group.example tf-log:tf-log-config:tf-log-machine-group
```
