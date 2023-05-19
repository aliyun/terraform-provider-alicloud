---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_flow_log"
sidebar_current: "docs-alicloud-resource-vpc-flow-log"
description: |-
  Provides a Alicloud Vpc Flow Log resource.
---

# alicloud_vpc_flow_log

Provides a Vpc Flow Log resource. While it uses alicloud_vpc_flow_log to build a vpc flow log resource, it will be active by default.

For information about Vpc Flow Log and how to use it, see [What is Flow Log](https://www.alibabacloud.com/help/en/virtual-private-cloud/latest/flow-logs-overview).

-> **NOTE:** Available in v1.117.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-testacc-example"
}

resource "alicloud_resource_manager_resource_group" "defaultRg" {
  resource_group_name = var.name
  display_name        = "tf-testAcc-rg78"
}

resource "alicloud_vpc" "defaultVpc" {
  vpc_name   = "${var.name}1"
  cidr_block = "10.0.0.0/8"
}

resource "alicloud_resource_manager_resource_group" "ModifyRG" {
  display_name        = "tf-testAcc-rg405"
  resource_group_name = "${var.name}2"
}

resource "alicloud_log_project" "default" {
  name = "${var.name}3"
}

resource "alicloud_log_store" "default" {
  project = alicloud_log_project.default.name
  name    = "${var.name}4"
}


resource "alicloud_vpc_flow_log" "default" {
  flow_log_name        = var.name
  log_store_name       = alicloud_log_store.default.name
  description          = "tf-testAcc-flowlog"
  traffic_path         = ["all"]
  project_name         = alicloud_log_project.default.name
  resource_type        = "VPC"
  resource_group_id    = alicloud_resource_manager_resource_group.defaultRg.id
  resource_id          = alicloud_vpc.defaultVpc.id
  aggregation_interval = "1"
  traffic_type         = "All"
}
```

## Argument Reference

The following arguments are supported:
* `aggregation_interval` - (Optional, Computed, Available in v1.205.0+) Data aggregation interval.
* `description` - (Optional) The Description of the VPC Flow Log.
* `flow_log_name` - (Optional) The Name of the VPC Flow Log.
* `log_store_name` - (Required, ForceNew) The name of the logstore.
* `project_name` - (Required, ForceNew) The name of the project.
* `resource_group_id` - (Optional, Computed, Available in v1.205.0+) The ID of the resource group.
* `resource_id` - (Required, ForceNew) The ID of the resource.
* `resource_type` - (Required, ForceNew) The resource type of the traffic captured by the flow log:-**NetworkInterface**: ENI.-**VSwitch**: All ENIs in the VSwitch.-**VPC**: All ENIs in the VPC.
* `status` - (Optional, Computed) The status of the VPC Flow Log. Valid values: **Active** and **Inactive**.
* `tags` - (Optional, Map, Available in v1.205.0+) The tag of the current instance resource.
* `traffic_path` - (Optional, ForceNew, Computed, Available in v1.205.0+) The collected flow path. Value:**all**: indicates full acquisition.**internetGateway**: indicates public network traffic collection.
* `traffic_type` - (Required, ForceNew) The type of traffic collected. Valid values:**All**: All traffic.**Allow**: Access control allowedtraffic.**Drop**: Access control denied traffic.



## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `business_status` - Business status.
* `create_time` - Creation time.
* `flow_log_id` - The flow log ID.

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