---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_flow_log"
sidebar_current: "docs-alicloud-resource-vpc-flow-log"
description: |-
  Provides a Alicloud VPC Flow Log resource.
---

# alicloud\_vpc\_flow\_log

Provides a VPC Flow Log resource.

For information about VPC Flow log and how to use it, see [Flow log overview](https://www.alibabacloud.com/help/doc-detail/127150.htm).

-> **NOTE:** Available in v1.117.0+

-> **NOTE:** While it uses `alicloud_vpc_flow_log` to build a vpc flow log resource, it will be active by default.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terratest_vpc_flow_log"
}

variable "log_store_name" {
  default = "vpc-flow-log-for-vpc"
}

variable "project_name" {
  default = "vpc-flow-log-for-vpc"
}

resource "alicloud_vpc" "default" {
  cidr_block = "192.168.0.0/24"
  name       = var.name
}

resource "alicloud_vpc_flow_log" "default" {
  depends_on     = ["alicloud_vpc.default"]
  resource_id    = alicloud_vpc.default.id
  resource_type  = "VPC"
  traffic_type   = "All"
  log_store_name = var.log_store_name
  project_name   = var.project_name
  flow_log_name  = var.name
  status         = "Active"
}

```
## Argument Reference

The following arguments are supported:

* `flow_log_name` - (Optional) The Name of the VPC Flow Log.
* `description` - (Optional) The Description of the VPC Flow Log.
* `resource_type` - (Required, ForceNew) The type of the resource to capture traffic. Valid values `NetworkInterface`, `VPC`, and `VSwitch`.
* `resource_id` - (Required, ForceNew) The ID of the resource.
* `traffic_type` - (Required, ForceNew) The type of traffic collected. Valid values `All`, `Drop` and `Allow`.
* `project_name` - (Required, ForceNew) The name of the project.
* `log_store_name` - (Required, ForceNew) The name of the logstore.
* `status` - (Optional, Computed) The status of the VPC Flow Log. Valid values `Active` and `Inactive`.

## Attributes Reference

The following attributes are exported:

* `id` - The Id of the VPC Flow Log.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when creating the VPC Flow Log (until it reaches the initial `Active` status). 
* `update` - (Defaults to 10 mins) Used when updating the VPC Flow Log (until it reaches the `Active` or `Inactive` status when you set `status`). 
* `delete` - (Defaults to 10 mins) Used when terminating the VPC Flow Log. 

## Import

VPC Flow Log can be imported using the id, e.g.

```shell
$ terraform import alicloud_vpc_flow_log.example fl-abc123456
```
