---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_flowlogs"
sidebar_current: "docs-alicloud-datasource-cen-flowlogs"
description: |-
  Provides a list of Cen Flow Log owned by an Alibaba Cloud account.
---

# alicloud_cen_flowlogs

This data source provides CEN flow logs available to the user.

-> **NOTE:** Available since v1.78.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

resource "alicloud_cen_instance" "defaultc5kxyC" {
  cen_instance_name = var.name
}

resource "alicloud_cen_transit_router" "defaultVw2U9u" {
  cen_id = alicloud_cen_instance.defaultc5kxyC.id
}

resource "alicloud_log_project" "default" {
  project_name = "${var.name}-${random_integer.default.result}"
  description  = "terraform-example"
}

resource "alicloud_log_store" "default" {
  project_name          = alicloud_log_project.default.project_name
  logstore_name         = "${var.name}-${random_integer.default.result}"
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}

resource "alicloud_cen_flowlog" "default" {
  project_name      = alicloud_log_store.default.project_name
  flow_log_name     = "${var.name}-${random_integer.default.result}"
  log_format_string = "$${srcaddr}$${dstaddr}$${bytes}"
  cen_id            = alicloud_cen_instance.defaultc5kxyC.id
  log_store_name    = alicloud_log_store.default.logstore_name
  interval          = "600"
  status            = "Active"
  transit_router_id = alicloud_cen_transit_router.defaultVw2U9u.transit_router_id
  description       = "flowlog-resource-example-1"
}

data "alicloud_cen_flowlogs" "default" {
  ids = ["${alicloud_cen_flowlog.default.id}"]
}

output "first_cen_flowlog_id" {
  value = "${data.alicloud_cen_flowlogs.default.flowlogs.0.id}"
}
```

## Argument Reference

The following arguments are supported:
* `cen_id` - (ForceNew, Optional) The ID of Cen instance.
* `description` - (ForceNew, Optional) The description of the flowlog.
* `flow_log_id` - (ForceNew, Optional) The ID of FlowLog.
* `flow_log_name` - (ForceNew, Optional) The name of the flowlog.
* `flow_log_version` - (ForceNew, Optional) Flowlog Version.
* `interval` - (ForceNew, Optional) The duration of the capture window for the flow log to capture traffic. Unit: seconds. Valid values: **60** or **600 * *. Default value: **600 * *.
* `log_store_name` - (ForceNew, Optional) The LogStore that stores the flowlog.
* `page_number` - (ForceNew, Optional) Current page number.
* `page_size` - (ForceNew, Optional) Number of records per page.
* `project_name` - (ForceNew, Optional) The Project that stores the flowlog.
* `region_id` - (ForceNew, Optional) Region id
* `status` - (ForceNew, Optional) The status of the flow log. Valid values:-**Active**: started.-**InActive**: not started.
* `transit_router_id` - (ForceNew, Optional) Transit Router ID
* `ids` - (Optional, ForceNew, Computed) A list of Flow Log IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Group Metric Rule name.
* `output_file` - (Optional, ForceNew) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Flow Log IDs.
* `names` - A list of name of Flow Logs.
* `flowlogs` - A list of Flow Log Entries. Each element contains the following attributes:
  * `id` - The ID of FlowLog.
  * `cen_id` - The ID of Cen instance.
  * `create_time` - The createTime of flowlog.
  * `description` - The description of the flowlog.
  * `flow_log_id` - The ID of FlowLog.
  * `flow_log_name` - The name of the flowlog.
  * `flow_log_version` - (Available since v1.236.0) Flowlog Version.
  * `tags` - The tag of the resource.
  * `interval` - (Available since v1.236.0) The duration of the capture window for the flow log to capture traffic. Unit: seconds. Valid values: **60** or **600 * *. Default value: **600 * *.
  * `log_format_string` - (Available since v1.236.0) Log Format.
  * `log_store_name` - The LogStore that stores the flowlog.
  * `project_name` - The Project that stores the flowlog.
  * `record_total` - (Available since v1.236.0) Total number of records.
  * `region_id` - (Available since v1.236.0) Region Id.
  * `status` - The status of the flow log. Valid values:-**Active**: started.-**InActive**: not started.
  * `transit_router_attachment_id` - (Available since v1.236.0) Cross-region Connection ID or VBR connection ID.> This parameter is required.
  * `transit_router_id` - (Available since v1.236.0) Transit Router ID.
