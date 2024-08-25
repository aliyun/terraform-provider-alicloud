---
subcategory: "Cloud Enterprise Network (CEN)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cen_flowlog"
sidebar_current: "docs-alicloud-resource-cen-flowlog"
description: |-
  Provides a Alicloud CEN manage route entried resource.
---

# alicloud_cen_flowlog

This resource used to create a flow log function in Cloud Enterprise Network (CEN). 
By using the flow log function, you can capture the traffic data of the network instances in different regions of a CEN. 
You can also use the data aggregated in flow logs to analyze cross-region traffic flows, minimize traffic costs, and troubleshoot network faults.

For information about CEN flow log and how to use it, see [Manage CEN flowlog](https://www.alibabacloud.com/help/en/cen/developer-reference/api-cbn-2017-09-12-createflowlog).

-> **NOTE:** Available since v1.73.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_cen_flowlog&exampleId=efefd3ba-1366-1815-787a-eea6b0b57f21758b6048&activeTab=example&spm=docs.r.cen_flowlog.0.efefd3ba13&intl_lang=EN_US" target="_blank">
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

* `cen_id` - (Required, ForceNew) The ID of the CEN Instance.
* `project_name` - (Required, ForceNew) The name of the SLS project.
* `log_store_name` - (Required, ForceNew) The name of the log store which is in the  `project_name` SLS project.
* `flow_log_name` - (Optional) The name of flowlog.
* `description` - (Optional) The description of flowlog.
* `status` - (Optional) The status of flowlog. Valid values: ["Active", "Inactive"]. Default to "Active".

## Attributes Reference

The following attributes are exported:

* `id` - ID of the flowlog.

## Import

CEN flowlog can be imported using the id, e.g.

```shell
$ terraform import alicloud_cen_flowlog.default flowlog-tig1xxxxxx
```
