---
subcategory: "VPC"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_flow_log"
sidebar_current: "docs-alicloud-resource-vpc-flow-log"
description: |-
  Provides a Alicloud VPC flow log resource.
---

# alicloud\_vpc\_flow\_log

This topic provides an overview of the flow log function of Virtual Private Cloud (VPC). 
By using this function, you can capture the inbound and outbound traffic over the Elastic Network Interface (ENI) in your VPC. 
With flow logs, you can check access control rules, monitor network traffic, and troubleshoot network faults.

For information about VPC flow log and how to use it, see [Manage VPC flow_log](https://www.alibabacloud.com/help/doc-detail/127150.htm).

-> **NOTE:** Available in 1.92.0+

-> **NOTE:**  Only the following regions support create VPC flow log.
[`cn-huhehaote`,`ap-southeast-3`,`ap-southeast-5`,`eu-west-1`,`ap-south-1`]

-> **NOTE:**  If the target VPC, the VPC to which the target VSwitch belongs, or the VPC to which the target ENI belongs, contains any instance of the following instance type families, you cannot create any flow log for the target VPC, VSwitch, or ENI.
[`ecs.c1`, `ecs.c2`, `ecs.c4`, `ecs.c5`, `ecs.ce4`, `ecs.cm4`, `ecs.d1`, `ecs.e3`, `ecs.e4`, `ecs.ga1`, `ecs.gn4`, `ecs.gn5`, `ecs.i1`, `ecs.m1`, `ecs.m2`, `ecs.mn4`, `ecs.n1`, `ecs.n2`, `ecs.n4`, `ecs.s1`, `ecs.s2`, `ecs.s3`, `ecs.se1`, `ecs.sn1`, `ecs.sn2`, `ecs.t1`, `ecs.xn4`]
To create flow logs for such resources, you must upgrade the instance type. For more information, see [Instance families that support instance type changes](https://www.alibabacloud.com/help/zh/doc-detail/89743.htm).

## Example Usage

Basic Usage

```
# Create a vpc flow log resource and use it to capture the inbound and outbound traffic over the Elastic Network Interface (ENI) in your Virtual Private Cloud (VPC).

resource "alicloud_vpc" "default" {
  cidr_block = "192.168.0.0/24"
  name       = var.name
}
data "alicloud_zones" "default" {
}
resource "alicloud_vswitch" "default" {
  cidr_block        = "192.168.0.0/24"
  availability_zone = data.alicloud_zones.default.zones[0].id
  vpc_id            = alicloud_vpc.default.id
}
resource "alicloud_security_group" "default" {
  vpc_id = alicloud_vpc.default.id
}
resource "alicloud_network_interface" "default" {
  vswitch_id        = alicloud_vswitch.default.id
  security_groups   = [ alicloud_security_group.default.id ]
  private_ip        = "192.168.0.2"
  private_ips_count = 3
}
resource "alicloud_log_project" "default"{
  name        = lower(var.name)
  description = "create by terraform"
}
resource "alicloud_log_store" "default"{
  project               = alicloud_log_project.default.name
  name                  = lower(var.name)
  retention_period      = 3650
  shard_count           = 3
  auto_split            = true
  max_split_shard_count = 60
  append_meta           = true
}
resource "alicloud_vpc_flow_log" "default" {
  resource_id    = alicloud_vpc.default.id
  resource_type  ="VPC"
  traffic_type   = "All"
  log_store_name = alicloud_log_store.default.name
  project_name   = alicloud_log_project.default.name
  flow_log_name  = var.name
  description    = var.description
}

```
## Argument Reference

The following arguments are supported:

* `resource_id` - (Required, ForceNew) The ID of the resource whose traffic you want to capture.
* `resource_type` - (Required, ForceNew) The type of the resource whose traffic you want to capture. Valid values: ["NetworkInterface", "VSwitch", "VPC"].
* `traffic_type` - (Required, ForceNew) The type of the traffic to be captured. Valid values: ["All", "Allow", "Drop"].
* `project_name` - (Required, ForceNew) The name of the SLS project.
* `log_store_name` - (Required, ForceNew) The name of the log store which is in the  `project_name` SLS project.
* `flow_log_name` - (Optional) The name of flow log. It must be 2 to 128 characters in length, and must begin with a letter or Chinese character (beginning with http:// or https:// is not allowed). It can contain digits, colons (:), underscores (_), or hyphens (-). Default value: null.
* `description` - (Optional) The description of flow log. It must be 2 to 256 characters in length and must not start with http:// or https://. Default value: null.
* `status` - (Optional) The status of flow log. Valid values: ["Active", "Inactive"]. Default to "Active".

## Attributes Reference

The following attributes are exported:

* `id` - ID of the flow log.

## Import

VPC flow log can be imported using the id, e.g.

```
$ terraform import alicloud_vpc_flow_log.default flowlog-tig1xxxxxx
```

