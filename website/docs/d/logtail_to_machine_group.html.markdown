---
layout: "alicloud"
page_title: "Alicloud: alicloud_logtail_to_machine_group"
sidebar_current: "docs-alicloud-logtail-to-machine-group"
description: |-
    Provides a list of log config and log machine groups.
---

# alicloud\_logtail\_to\_machine\_group

This data source provides the apis of the configuration and machine group under this project.
                                          
                                          
## Example Usage

```
 data "alicloud_logtail_to_machine_group" "example" {
    project = "tf-project"
    output_file = "~/newdata/map.json"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required, ForceNew) The project name to the log store belongs
* `output_file` - (Optional) .
* `offset` - (Optional) .
* `size` - (Optional) .

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - The unique ID hash value of this query resource。
* `machine_group` - A list of log machine group. 
* `logtail_config` - A list of logtail config.
