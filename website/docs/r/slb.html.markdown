---
layout: "alicloud"
page_title: "Alicloud: alicloud_slb"
sidebar_current: "docs-alicloud-resource-slb"
description: |-
  Provides an Application Load Banlancer resource.
---

# alicloud\_slb

Provides an Application Load Balancer resource.

-> **NOTE:** Resource `alicloud_slb` has deprecated 'listener' filed from terraform-alicloud-provider [version 1.3.0](https://github.com/alibaba/terraform-provider/releases/tag/V1.3.0) . You can create new listeners for Load Balancer by resource `alicloud_slb_listener`.
If you have had several listeners in one load balancer, you can import them via the specified listener ID. In the `alicloud_slb_listener`, listener ID is consist of load balancer ID and frontend port, and its format is `<load balancer ID>:<frontend port>`, like "lb-hr2fwnf32t:8080".

-> **NOTE:** At present, to avoid some unnecessary regulation confusion, SLB can not support alicloud international account to create "paybybandwidth" instance.

-> **NOTE:** The supported specifications vary by region. Currently not all regions support guaranteed-performance instances.
For more details about guaranteed-performance instance, see [Guaranteed-performance instances](https://www.alibabacloud.com/help/doc-detail/27657.htm).

## Example Usage

```
variable "name" {
  default = "terraformtestslbconfig"
}
data "alicloud_zones" "default" {
	available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "172.16.0.0/21"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
  name = "${var.name}"
}

resource "alicloud_slb" "default" {
  name = "${var.name}"
  specification = "slb.s2.small"
  vswitch_id = "${alicloud_vswitch.default.id}"
  tags = {
    tag_a = 1
    tag_b = 2
    tag_c = 3
    tag_d = 4
    tag_e = 5
    tag_f = 6
    tag_g = 7
    tag_h = 8
    tag_i = 9
    tag_j = 10
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the SLB. This name must be unique within your AliCloud account, can have a maximum of 80 characters,
must contain only alphanumeric characters or hyphens, such as "-","/",".","_", and must not begin or end with a hyphen. If not specified,
Terraform will autogenerate a name beginning with `tf-lb`.
* `internet` - (Optional, ForceNew) If true, the SLB addressType will be internet, false will be intranet, Default is false. If load balancer launched in VPC, this value must be "false".
* `internet_charge_type` - (Optional, ForceNew) Valid
  values are `PayByBandwidth`, `PayByTraffic`. If this value is "PayByBandwidth", then argument "internet" must be "true". Default is "PayByTraffic". If load balancer launched in VPC, this value must be "PayByTraffic".
  Before version 1.10.1, the valid values are "paybybandwidth" and "paybytraffic".
* `bandwidth` - (Optional) Valid
  value is between 1 and 1000, If argument "internet_charge_type" is "paybytraffic", then this value will be ignore.
* `listener` - (Deprecated) The field has been deprecated from terraform-alicloud-provider [version 1.3.0](https://github.com/alibaba/terraform-provider/releases/tag/V1.3.0), and use resource `alicloud_slb_listener` to replace.
* `vswitch_id` - (Required for a VPC SLB, Forces New Resource) The VSwitch ID to launch in.
* `specification` - (Optional) The specification of the Server Load Balancer instance. Default to empty string indicating it is "Shared-Performance" instance.
 Launching "[Performance-guaranteed](https://www.alibabacloud.com/help/doc-detail/27657.htm)" instance, it is must be specified and it valid values are: "slb.s1.small", "slb.s2.small", "slb.s2.medium",
 "slb.s3.small", "slb.s3.medium" and "slb.s3.large".
* `tags` - (Optional) A mapping of tags to assign to the resource. The `tags` can have a maximum of 10 tag for every load balancer instance.
* `instance_charge_type` - (Optional, Available in v1.34.0+) The billing method of the load balancer. Valid values are "PrePaid" and "PostPaid". Default to "PostPaid".
* `period` - (Optional, Available in v1.34.0+) The duration that you will buy the resource, in month. It is valid when `instance_charge_type` is `PrePaid`. Default to 1. Valid values: [1-9, 12, 24, 36].
* `master_zone_id` - (Optional, ForceNew, Available in v1.36.0+) The primary zone ID of the SLB instance. If not specified, the system will be randomly assigned. You can query the primary and standby zones in a region by calling the DescribeZone API.
* `slave_zone_id` - (Optional, ForceNew, Available in v1.36.0+) The standby zone ID of the SLB instance. If not specified, the system will be randomly assigned. You can query the primary and standby zones in a region by calling the DescribeZone API.

-> **NOTE:** A "Shared-Performance" instance can be changed to "Performance-guaranteed", but the change is irreversible.

-> **NOTE:** To change a "Shared-Performance" instance to a "Performance-guaranteed" instance, the SLB will have a short probability of business interruption (10 seconds-30 seconds). Advise to change it during the business downturn, or migrate business to other SLB Instances by using GSLB before changing.

-> **NOTE:** Currently, the alibaba cloud international account does not support creating a PrePaid SLB instance.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the load balancer.
* `name` - The name of the load balancer.
* `internet` - The internet of the load balancer.
* `internet_charge_type` - The internet_charge_type of the load balancer.
* `bandwidth` - The bandwidth of the load balancer.
* `vswitch_id` - The VSwitch ID of the load balancer. Only available on SLB launched in a VPC.
* `address` - The IP address of the load balancer.
* `specification` - The specification of the Server Load Balancer instance.

## Import

Load balancer can be imported using the id, e.g.

```
$ terraform import alicloud_slb.example lb-abc123456
```
