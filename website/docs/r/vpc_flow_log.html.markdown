---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_flow_log"
sidebar_current: "docs-alicloud-resource-vpc-flow-log"
description: |-
  Provides a Alicloud Vpc Flow Log resource.
---

# alicloud_vpc_flow_log

Provides a Vpc Flow Log resource. 

For information about Vpc Flow Log and how to use it, see [What is Flow Log](https://www.alibabacloud.com/help/en/virtual-private-cloud/latest/flow-logs-overview).

-> **NOTE:** Available in v1.117.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "alicloud_resource_manager_resource_group" "qWOSqC" {
  display_name        = "test01"
  resource_group_name = var.name
}

resource "alicloud_vpc" "f9wsFd" {
  vpc_name          = "${var.name}1"
  cidr_block        = "10.0.0.0/8"
  resource_group_id = alicloud_resource_manager_resource_group.qWOSqC.id
}

resource "alicloud_resource_manager_resource_group" "ModifyRG" {
  display_name        = "test02"
  resource_group_name = "${var.name}2"
}


resource "alicloud_vpc_flow_log" "default" {
  flow_log_name        = var.name
  log_store_name       = "rdktest"
  description          = "test"
  traffic_path         = ["all"]
  project_name         = "rdktest"
  resource_type        = "VPC"
  resource_group_id    = alicloud_resource_manager_resource_group.qWOSqC.id
  resource_id          = alicloud_vpc.f9wsFd.id
  aggregation_interval = "1"
  traffic_type         = "All"
}
```


## Argument Reference

The following arguments are supported:
* `aggregation_interval` - (Optional, Computed, Available in v1.207.0+) Data aggregation interval.
* `description` - (Optional) The Description of flow log.
* `flow_log_name` - (Optional) The flow log name.
* `log_store_name` - (Required, ForceNew) The log store name.
* `project_name` - (Required, ForceNew) The project name.
* `resource_group_id` - (Optional, Computed, Available in v1.207.0+) The ID of the resource group to which the VPC belongs.
* `resource_id` - (Required, ForceNew) The resource id.
* `resource_type` - (Required) The resource type of the traffic captured by the flow log:
  - **NetworkInterface**: ENI.
  - **VSwitch**: All ENIs in the VSwitch.
  - **VPC**: All ENIs in the VPC.
* `status` - (Optional, Computed) The status of  flow log.
* `tags` - (Optional, Map, Available in v1.207.0+) The tags of PrefixList.
* `traffic_path` - (Optional, ForceNew, Computed, Available in v1.207.0+) 采集的流量路径。取值：    all（默认值）：表示全量采集。     internetGateway：表示公网流量采集。.
* `traffic_type` - (Required, ForceNew) The traffic type.



## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - the time of creation.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Flow Log.
* `delete` - (Defaults to 5 mins) Used when delete the Flow Log.
* `update` - (Defaults to 5 mins) Used when update the Flow Log.

## Import

Vpc Flow Log can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_flow_log.example <id>
```