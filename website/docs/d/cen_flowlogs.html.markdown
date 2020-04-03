---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_flowlogs"
sidebar_current: "docs-alicloud-datasource-cen-flowlogs"
description: |-
    Provides a list of CEN flow logs owned by an Alibaba Cloud account.
---

# alicloud\_cen\_flowlogs

This data source provides CEN flow logs available to the user.

-> **NOTE:** Available in 1.78.0+

## Example Usage

Basic Usage

```
data "alicloud_cen_flowlogs" "default" {
  ids        = ["flowlog-tig1xxxxx"]
  name_regex = "^foo"
}

output "first_cen_flowlog_id" {
  value = "${data.alicloud_cen_instances.default.flowlogs.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of CEN flow log IDs.
* `name_regex` - (Optional) A regex string to filter CEN flow logs by name.
* `cen_id` - (Optional) The ID of the CEN Instance.
* `project_name` - (Optional) The name of the SLS project.
* `log_store_name` - (Optional) The name of the log store which is in the  `project_name` SLS project.
* `flow_log_name` - (Optional) The name of flowlog.
* `description` - (Optional) The description of flowlog.
* `status` - (Optional) The status of flowlog. Valid values: ["Active", "Inactive"]. Default to "Active".
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of CEN flow log IDs.
* `names` - A list of CEN flow log names. 
* `instances` - A list of CEN flow logs. Each element contains the following attributes:
  * `id` - ID of the CEN flow log.
  * `flow_log_id` - ID of the CEN flow log.
  * `cen_id` -  The ID of the CEN Instance.
  * `project_name` - The name of the SLS project.
  * `log_store_name` - The name of the log store which is in the  `project_name` SLS project.
  * `flow_log_name` - The name of flowlog.
  * `description` -  The description of flowlog.
  * `status` - The status of flowlog.
