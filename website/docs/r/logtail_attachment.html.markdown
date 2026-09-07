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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_logtail_attachment&exampleId=efc3eb85-d6ae-b13c-578a-d54b0589f151a4c74fa3&activeTab=example&spm=docs.r.logtail_attachment.0.efc3eb85d6&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_log_project" "example" {
  project_name = "terraform-example-${random_integer.default.result}"
  description  = "terraform-example"
}

resource "alicloud_log_store" "example" {
  project_name          = alicloud_log_project.example.project_name
  logstore_name         = "example-store"
  retention_period      = 3650
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}

resource "alicloud_logtail_config" "example" {
  project      = alicloud_log_project.example.project_name
  logstore     = alicloud_log_store.example.logstore_name
  input_type   = "file"
  name         = "terraform-example"
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

resource "alicloud_log_machine_group" "example" {
  project       = alicloud_log_project.example.project_name
  name          = "terraform-example"
  identify_type = "ip"
  topic         = "terraform"
  identify_list = ["10.0.0.1", "10.0.0.2"]
}

resource "alicloud_logtail_attachment" "example" {
  project             = alicloud_log_project.example.project_name
  logtail_config_name = alicloud_logtail_config.example.name
  machine_group_name  = alicloud_log_machine_group.example.name
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_logtail_attachment&spm=docs.r.logtail_attachment.example&intl_lang=EN_US)

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
