---
subcategory: "Cms"
layout: "alicloud"
page_title: "Alicloud: alicloud_cms_biz_trace"
description: |-
  Provides a Alicloud Cms Biz Trace resource.
---

# alicloud_cms_biz_trace

Provides a Cms Biz Trace resource.

For information about Cms Biz Trace and how to use it, see [What is Biz Trace](https://next.api.alibabacloud.com/document/Cms/2024-03-30/CreateBizTrace).

-> **NOTE:** Available since v1.292.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_cms_biz_trace" "default" {
  biz_trace_code  = "example-biz-trace-code"
  biz_trace_name  = var.name
  rule_config     = jsonencode([{ "entrancePid" : "xxx@d9w3jd9j3", "rpcMatcher" : { "matchType" : "EQUALS", "pattern" : "/" } }])
  advanced_config = jsonencode({ "sample" : { "strategy" : "BY_APP" } })
  workspace       = "example-workspace"
}
```

## Argument Reference

The following arguments are supported:

* `biz_trace_code` - (Required, ForceNew) The business trace code, which is used as the value of the business trace tag.
* `biz_trace_name` - (Optional) The name of the business trace.
* `rule_config` - (Optional) The rule configuration list of the business trace. It must be a valid JSON array string, e.g. `[{"entrancePid":"xxx@d9w3jd9j3","rpcMatcher":{"matchType":"EQUALS","pattern":"/"},"characteristics":{"operation":"OR","rules":[{"target":"CUSTOM_EXTRACT","id":"oi0b3bb7","key":"biz.test","matcher":{"matchType":"CONTAINS","pattern":["1"]}}]}}]`. **NOTE:** The API enriches this value with server-resolved keys when reading it back; the provider suppresses the resulting diff as long as the configured JSON is a subset of the stored JSON.
* `advanced_config` - (Optional) The advanced configuration of the business trace. It must be a valid JSON object string, e.g. `{"sample":{"strategy":"BY_APP"}}`. **NOTE:** The API may enrich this value with server-resolved keys when reading it back; the provider suppresses the resulting diff as long as the configured JSON is a subset of the stored JSON.
* `workspace` - (Optional, ForceNew) The workspace to which the business trace belongs. **NOTE:** The `UpdateBizTrace` API silently ignores workspace changes, so modifying this field recreates the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource supplied above. Its value is the Biz Trace ID.
* `biz_trace_id` - The ID of the Biz Trace.
* `create_time` - The creation time of the resource.
* `region_id` - The region ID of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Biz Trace.
* `delete` - (Defaults to 5 mins) Used when delete the Biz Trace.
* `update` - (Defaults to 5 mins) Used when update the Biz Trace.

## Import

Cms Biz Trace can be imported using the id, e.g.

```shell
$ terraform import alicloud_cms_biz_trace.example <biz_trace_id>
```
