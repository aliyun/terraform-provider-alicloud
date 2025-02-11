---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_flow_log"
description: |-
  Provides a Alicloud VPC Flow Log resource.
---

# alicloud_vpc_flow_log

Provides a VPC Flow Log resource.

While it uses alicloud_vpc_flow_log to build a vpc flow log resource, it will be active by default.

For information about VPC Flow Log and how to use it, see [What is Flow Log](https://www.alibabacloud.com/help/en/virtual-private-cloud/latest/flow-logs-overview).

-> **NOTE:** Available since v1.117.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_vpc_flow_log&exampleId=736bfeac-3fca-11b4-b45a-d933bca5c287ac089a4f&activeTab=example&spm=docs.r.vpc_flow_log.0.736bfeac3f&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}
data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}
resource "alicloud_vpc" "example" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}
resource "random_uuid" "example" {
}
resource "alicloud_log_project" "example" {
  project_name = substr("tf-example-${replace(random_uuid.example.result, "-", "")}", 0, 16)
  description  = var.name
}

resource "alicloud_log_store" "example" {
  project_name          = alicloud_log_project.example.project_name
  logstore_name         = var.name
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}

resource "alicloud_vpc_flow_log" "example" {
  flow_log_name        = var.name
  log_store_name       = alicloud_log_store.example.logstore_name
  description          = var.name
  traffic_path         = ["all"]
  project_name         = alicloud_log_project.example.project_name
  resource_type        = "VPC"
  resource_group_id    = data.alicloud_resource_manager_resource_groups.default.ids.0
  resource_id          = alicloud_vpc.example.id
  aggregation_interval = "1"
  traffic_type         = "All"
}
```

## Argument Reference

The following arguments are supported:
* `aggregation_interval` - (Optional, Available since v1.205.0) Data aggregation interval
* `description` - (Optional) The Description of the VPC Flow Log.
* `flow_log_name` - (Optional) The Name of the VPC Flow Log.
* `ip_version` - (Optional, Available since v1.243.0) The IP address type of the collected traffic.
* `log_store_name` - (Required, ForceNew) The name of the logstore.
* `project_name` - (Required, ForceNew) The name of the project.
* `resource_group_id` - (Optional, Available since v1.205.0) The ID of the resource group.
* `resource_id` - (Required, ForceNew) The ID of the resource.
* `resource_type` - (Required, ForceNew) The resource type of the traffic captured by the flow log:
  - `NetworkInterface`: ENI.
  - `VSwitch`: All ENIs in the VSwitch.
  - `VPC`: All ENIs in the VPC.
* `status` - (Optional) The status of the VPC Flow Log. Valid values: `Active` and `Inactive`.
* `tags` - (Optional, Map, Available since v1.205.0) The tag of the current instance resource.
* `traffic_path` - (Optional, ForceNew, List, Available since v1.205.0) The collected flow path. Value:
  - *all**: indicates full acquisition.
  - *internetGateway**: indicates public network traffic collection.
* `traffic_type` - (Required, ForceNew) The type of traffic collected. Valid values:
  - *All**: All traffic.
  - *Allow**: Access control allowedtraffic.
  - *Drop**: Access control denied traffic.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `business_status` - Business status
* `create_time` - Creation time
* `flow_log_id` - The flow log ID.
* `region_id` - (Available since v1.243.0) The region ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Flow Log.
* `delete` - (Defaults to 5 mins) Used when delete the Flow Log.
* `update` - (Defaults to 5 mins) Used when update the Flow Log.

## Import

VPC Flow Log can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_flow_log.example <id>
```
