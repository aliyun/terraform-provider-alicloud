---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_site_monitors"
sidebar_current: "docs-alicloud-datasource-cms-site-monitors"
description: |-
  Provides a list of Cloud Monitor Service Site Monitor owned by an Alibaba Cloud account.
---

# alicloud_cms_site_monitors

This data source provides Cloud Monitor Service Site Monitor available to the user.[What is Site Monitor](https://www.alibabacloud.com/help/en/cms/developer-reference/api-cms-2019-01-01-createsitemonitor)

-> **NOTE:** Available since v1.224.0.

## Example Usage

```terraform
variable "name" {
  default = "tf_example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_cms_site_monitor" "default" {
  address   = "http://www.alibabacloud.com"
  task_name = "terraform-example-${random_integer.default.result}"
  task_type = "HTTP"
  interval  = 5
  isp_cities {
    city = "546"
    isp  = "465"
  }
  options_json = <<EOT
{
    "http_method": "get",
    "waitTime_after_completion": null,
    "ipv6_task": false,
    "diagnosis_ping": false,
    "diagnosis_mtr": false,
    "assertions": [
        {
            "operator": "lessThan",
            "type": "response_time",
            "target": 1000
        }
    ],
    "time_out": 30000
}
EOT
}

data "alicloud_cms_site_monitors" "default" {
  ids       = ["${alicloud_cms_site_monitor.default.id}"]
  task_type = "HTTP"
}

output "alicloud_cms_site_monitor_example_id" {
  value = data.alicloud_cms_site_monitors.default.monitors.0.task_id
}
```

## Argument Reference

The following arguments are supported:
* `task_id` - (ForceNew, Optional) Task ID.
* `task_type` - (ForceNew, Optional) Task Type.
* `ids` - (Optional, ForceNew, Computed) A list of Site Monitor IDs.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Site Monitor IDs.
* `monitors` - A list of Site Monitor Entries. Each element contains the following attributes:
  * `address` - Address.
  * `create_time` - CreateTime.
  * `interval` - Monitoring frequency.
  * `task_id` - Task Id.
  * `task_name` - Task Name.
  * `task_type` - Task Type.
