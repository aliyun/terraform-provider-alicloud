---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_flowlog"
description: |-
  Provides a Alicloud CEN Flow Log resource.
---

# alicloud_cen_flowlog

Provides a CEN Flow Log resource.



For information about CEN Flow Log and how to use it, see [What is Flow Log](https://www.alibabacloud.com/help/en/cen/developer-reference/api-cbn-2017-09-12-createflowlog).

-> **NOTE:** Available since v1.73.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cen_flowlog&exampleId=efefd3ba-1366-1815-787a-eea6b0b57f21758b6048&activeTab=example&spm=docs.r.cen_flowlog.0.efefd3ba13&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
# Create a cen flowlog resource and use it to publish a route entry pointing to an ECS.
resource "alicloud_cen_instance" "default" {
  name = "my-cen"
}
resource "alicloud_log_project" "default" {
  name        = "sls-for-flowlog"
  description = "create by terraform"
}
resource "alicloud_log_store" "default" {
  project               = alicloud_log_project.default.name
  name                  = "sls-for-flowlog"
  retention_period      = 3650
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}

resource "alicloud_cen_flowlog" "default" {
  flow_log_name  = "my-flowlog"
  cen_id         = alicloud_cen_instance.default.id
  project_name   = alicloud_log_project.default.name
  log_store_name = alicloud_log_store.default.name
}
```

## Argument Reference

The following arguments are supported:
* `cen_id` - (Required, ForceNew) cen id
* `description` - (Optional) The description of the flowlog.
* `flow_log_name` - (Optional) The name of the flowlog.
* `interval` - (Optional, Int, Available since v1.235.0) The duration of the capture window for the flow log to capture traffic. Unit: seconds. Valid values: `60` or **600 * *. Default value: **600 * *.
* `log_format_string` - (Optional, ForceNew, Available since v1.235.0) Log Format
* `log_store_name` - (Required, ForceNew) The LogStore that stores the flowlog.
* `project_name` - (Required, ForceNew) The Project that stores the flowlog.
* `status` - (Optional, Computed) The status of the flow log. Valid values:
  - `Active`: started.
  - `InActive`: not started.
* `tags` - (Optional, Map, Available since v1.235.0) The tag of the resource
* `transit_router_attachment_id` - (Optional, ForceNew, Available since v1.235.0) Cross-region Connection ID or VBR connection ID.

-> **NOTE:**  This parameter is required.

* `transit_router_id` - (Optional, ForceNew, Available since v1.235.0) Transit Router ID

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - CreateTime
* `region_id` - region id

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Flow Log.
* `delete` - (Defaults to 5 mins) Used when delete the Flow Log.
* `update` - (Defaults to 5 mins) Used when update the Flow Log.

## Import

CEN Flow Log can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_flowlog.example <id>
```