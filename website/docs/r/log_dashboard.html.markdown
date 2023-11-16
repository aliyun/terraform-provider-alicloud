---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_log_dashboard"
sidebar_current: "docs-alicloud-resource-log-dashboard"
description: |-
  Provides a Alicloud Log Dashboard resource.
---

# alicloud_log_dashboard

The dashboard is a real-time data analysis platform provided by the log service. You can display frequently used query and analysis statements in the form of charts and save statistical charts to the dashboard.
[Refer to details](https://www.alibabacloud.com/help/doc-detail/102530.htm).

-> **NOTE:** Available since v1.86.0.

## Example Usage

Basic Usage

```terraform
resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_log_project" "example" {
  name        = "terraform-example-${random_integer.default.result}"
  description = "terraform-example"
}

resource "alicloud_log_store" "example" {
  project               = alicloud_log_project.example.name
  name                  = "example-store"
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}

resource "alicloud_log_dashboard" "example" {
  project_name   = alicloud_log_project.example.name
  dashboard_name = "terraform-example"
  display_name   = "terraform-example"
  attribute      = <<EOF
  {
    "type":"grid"
  }
EOF
  char_list      = <<EOF
  [
    {
      "action": {},
      "title":"new_title",
      "type":"map",
      "search":{
        "logstore":"example-store",
        "topic":"new_topic",
        "query":"* | SELECT COUNT(name) as ct_name, COUNT(product) as ct_product, name,product GROUP BY name,product",
        "start":"-86400s",
        "end":"now"
      },
      "display":{
        "xAxis":[
          "ct_name"
        ],
        "yAxis":[
          "ct_product"
        ],
        "xPos":0,
        "yPos":0,
        "width":10,
        "height":12,
        "displayName":"terraform-example"
      }
    }
  ]
EOF
}
```

## Argument Reference

The following arguments are supported:

* `project_name` - (Required, ForceNew) The name of the log project. It is the only in one Alicloud account.
* `dashboard_name` - (Required, ForceNew) The name of the Log Dashboard.
* `char_list` - (Required) Configuration of charts in the dashboard.
  **Note:** From version 1.164.0, `char_list` can set parameter "action".
* `display_name` - (Optional) Dashboard alias.
* `attribute` - (Optional, Available since v1.183.0) Dashboard attribute.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Log Dashboard. It formats as `<project_name>:<dashboard_name>`.

## Import

Log Dashboard can be imported using the id, e.g.

```shell
$ terraform import alicloud_log_dashboard.example <project_name>:<dashboard_name>
```
