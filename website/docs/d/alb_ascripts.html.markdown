---
subcategory: "Application Load Balancer (ALB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_alb_ascripts"
sidebar_current: "docs-alicloud-datasource-alb-ascripts"
description: |-
  Provides a list of Alb Ascript owned by an Alibaba Cloud account.
---

# alicloud_alb_ascripts

This data source provides Alb Ascript available to the user.[What is AScript](https://www.alibabacloud.com/help/en/)

-> **NOTE:** Available in 1.195.0+

## Example Usage

```
data "alicloud_alb_ascripts" "default" {
  ids          = ["${alicloud_alb_ascript.default.id}"]
  name_regex   = alicloud_alb_ascript.default.name
  ascript_name = "test"
  listener_id  = var.listenerId
}

output "alicloud_alb_ascript_example_id" {
  value = data.alicloud_alb_ascripts.default.ascripts.0.id
}
```

## Argument Reference

The following arguments are supported:
* `ids` - (Optional, ForceNew, Computed) A list of AScript IDs.
* `ascript_name` - (ForceNew,Optional) Script name.
* `listener_id` - (ForceNew,Optional) Listener ID of script attribution
* `names` - (Optional, ForceNew) The name of the AScript. You can specify at most 10 names.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Group Metric Rule name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of AScript IDs.
* `names` - A list of name of AScripts.
* `ascripts` - A list of AScript Entries. Each element contains the following attributes:
  * `ascript_id` - Script identification.
  * `ascript_name` - Script name.
  * `enabled` - Whether scripts are enabled.
  * `ext_attribute_enabled` - Whether extension parameters are enabled.
  * `ext_attributes` - Extended attribute list.
    * `attribute_key` - The key of the extended attribute.
    * `attribute_value` - The value of the extended attribute.
  * `listener_id` - Listener ID of script attribution.
  * `position` - Script execution location.
  * `script_content` - Script content.
  * `status` - Script status.
