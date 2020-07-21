---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_flow_logs"
sidebar_current: "docs-alicloud-datasource-vpc-flow-logs"
description: |-
    Provides a list of VPC flow logs owned by an Alibaba Cloud account.
---

# alicloud\_vpc\_flow\_logs

This data source provides VPC flow logs available to the user.

-> **NOTE:** Available in 1.92.0+

## Example Usage

Basic Usage

```
data "alicloud_vpc_flow_logs" "default" {
  ids        = ["flowlog-tig1xxxxx"]
  name_regex = "^foo"
}

output "first_vpc_flow_log_id" {
  value = "${data.alicloud_vpc_flow_logs.default.flow_logs.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of VPC flow log IDs.
* `name_regex` - (Optional) A regex string to filter VPC flow logs by name.
* `resource_id` - (Required, ForceNew) The ID of the resource whose traffic you want to capture.
* `resource_type` - (Required, ForceNew) The type of the resource whose traffic you want to capture. Valid values: ["NetworkInterface", "VSwitch", "VPC"].
* `traffic_type` - (Required, ForceNew) The type of the traffic to be captured. Valid values: ["All", "Allow", "Drop"].
* `project_name` - (Optional) The name of the SLS project.
* `log_store_name` - (Optional) The name of the log store which is in the  `project_name` SLS project.
* `description` - (Optional) The description of flow log.
* `status` - (Optional) The status of flow log. Valid values: ["Active", "Inactive"]. Default to "Active".
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of VPC flow log IDs.
* `names` - A list of VPC flow log names. 
* `flow_logs` - A list of VPC flow logs. Each element contains the following attributes:
  * `id` - ID of the VPC flow log.
  * `project_name` - The name of the SLS project.
  * `log_store_name` - The name of the log store which is in the  `project_name` SLS project.
  * `flow_log_name` - The name of flow log.
  * `description` -  The description of flow log.
  * `creation_time` - Time of creation.
  * `region_id` - The region to which the flow log belongs.
