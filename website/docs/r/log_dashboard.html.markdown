---
subcategory: "Log Service (SLS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_log_dashboard"
sidebar_current: "docs-alicloud-resource-log-dashboard"
description: |-
  Provides a Alicloud Log Dashboard resource.
---

# alicloud\_log\_dashboard
The dashboard is a real-time data analysis platform provided by the log service. You can display frequently used query and analysis statements in the form of charts and save statistical charts to the dashboard.
[Refer to details](https://www.alibabacloud.com/help/doc-detail/102530.htm).

-> **NOTE:** Available in 1.86.0, parameter "action" in char_list is supported since 1.164.0+. 

## Example Usage

Basic Usage

```terraform
resource "alicloud_log_project" "default" {
  name        = "tf-project"
  description = "tf unit test"
}
resource "alicloud_log_store" "default" {
  project          = "tf-project"
  name             = "tf-logstore"
  retention_period = "3000"
  shard_count      = 1
}
resource "alicloud_log_dashboard" "example" {
  project_name   = "tf-project"
  dashboard_name = "tf-dashboard"
  char_list      = <<EOF
  [
    {
      "action": {},
      "title":"new_title",
      "type":"map",
      "search":{
        "logstore":"tf-logstore",
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
        "displayName":"xixihaha911"
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
* `display_name` - (Optional) Dashboard alias.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Log Dashboard. It sames as its name.

## Import

Log Dashboard can be imported using the id or name, e.g.

```
$ terraform import alicloud_log_dashboard.example tf-project:tf-logstore:tf-dashboard
```
