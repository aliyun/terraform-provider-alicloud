---
layout: "alicloud"
page_title: "Alicloud: alicloud_logtail_to_machine_group"
sidebar_current: "docs-alicloud-resource-logtail-to-machine-group"
description: |-
  Provides a Alicloud logtail to machine  resource.
---

# alicloud\_logtail\_to\_machine_group

The Logtail access service is a log collection agent provided by Log Service. 
You can use Logtail to collect logs from servers such as Alibaba Cloud Elastic
Compute Service (ECS) instances in real time in the Log Service console. [Refer to details](https://www.alibabacloud.com/help/doc-detail/29058.htm
)

## Example Usage

Basic Usage

```
 data "alicloud_logtail_to_machine_group" "example" {
    project = "tf-project"
}
resource "alicloud_logtail_to_machine_group" "test" {
	project = "tf-project"
	logtail_config_name = "${data.alicloud_logtail_to_machine_group.example.logtail_config.0}"
	machine_group_name = "${data.alicloud_logtail_to_machine_group.example.machine_group.0}"
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
* `project` - The project name.
* `logtail_config_name` - The Logtail configuration name..
* `machine_group_name` - The machine group name.

## Import

Logtial to machine group can be imported using the id, e.g.

```
$ terraform import alicloud_logtail_to_machine_group.example tf-log:$(logtail_config_name):$(machine_group_name)
```