---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_biz_traces"
sidebar_current: "docs-alicloud-datasource-cms-biz-traces"
description: |-
  Provides a list of Cms Biz Traces to the user.
---

# alicloud\_cms\_biz\_traces

This data source provides the Cms Biz Traces of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.292.0.

## Example Usage

Basic Usage

```terraform
data "alicloud_cms_biz_traces" "ids" {
  ids = ["example_id"]
}
output "cms_biz_trace_id_1" {
  value = data.alicloud_cms_biz_traces.ids.biz_traces.0.biz_trace_id
}

data "alicloud_cms_biz_traces" "nameRegex" {
  name_regex = "^my-BizTrace"
}
output "cms_biz_trace_id_2" {
  value = data.alicloud_cms_biz_traces.nameRegex.biz_traces.0.biz_trace_id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed) A list of Biz Trace IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Biz Trace name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `workspace` - (Optional, ForceNew) The workspace to which the business traces belong.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `biz_traces` - A list of Cms Biz Traces. Each element contains the following attributes:
  * `advanced_config` - The advanced configuration, a JSON object string (the API may enrich it with server-resolved keys).
  * `biz_trace_code` - The business trace code, which is used as the value of the business trace tag.
  * `biz_trace_id` - The ID of the Biz Trace.
  * `biz_trace_name` - The name of the business trace.
  * `create_time` - The creation time of the resource.
  * `region_id` - The region ID of the resource.
  * `rule_config` - The rule configuration list, a JSON array string (the API enriches it with server-resolved keys).
  * `workspace` - The workspace to which the business trace belongs.
