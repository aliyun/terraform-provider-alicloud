---
subcategory: "Classic Load Balancer (CLB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_slb"
sidebar_current: "docs-alicloud-resource-slb"
description: |-
  Provides an Application Load Banlancer resource.
---

# alicloud\_slb

-> **DEPRECATED:** This resource has been renamed to [alicloud_slb_load_balancer](https://www.terraform.io/docs/providers/alicloud/r/slb_load_balancer) from version 1.123.1.

Provides an Application Load Balancer resource.

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
  name       = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "default" {
  vpc_id            = alicloud_vpc.default.id
  cidr_block        = "172.16.0.0/21"
  zone_id           = data.alicloud_zones.default.zones[0].id
  vswitch_name      = var.name
}

resource "alicloud_slb" "default" {
  name          = var.name
  specification = "slb.s2.small"
  vswitch_id    = alicloud_vswitch.default.id
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
* `internet` - (Deprecated) Field 'internet' has been deprecated from provider version 1.55.3. Use 'address_type' replaces it.
* `address_type` - (Optional, ForceNew, Available in 1.55.3+) The network type of the SLB instance. Valid values: ["internet", "intranet"]. If load balancer launched in VPC, this value must be "intranet".
    - internet: After an Internet SLB instance is created, the system allocates a public IP address so that the instance can forward requests from the Internet.
    - intranet: After an intranet SLB instance is created, the system allocates an intranet IP address so that the instance can only forward intranet requests.
* `internet_charge_type` - (Optional, ForceNew) Valid
  values are `PayByBandwidth`, `PayByTraffic`. If this value is "PayByBandwidth", then argument "internet" must be "true". Default is "PayByTraffic". If load balancer launched in VPC, this value must be "PayByTraffic".
  Before version 1.10.1, the valid values are "paybybandwidth" and "paybytraffic".
* `bandwidth` - (Optional) Valid
  value is between 1 and 1000, If argument "internet_charge_type" is "paybytraffic", then this value will be ignore.
* `vswitch_id` - (Required for a VPC SLB, Forces New Resource) The VSwitch ID to launch in. If `address_type` is internet, it will be ignore.
* `specification` - (Optional) The specification of the Server Load Balancer instance. Default to empty string indicating it is "Shared-Performance" instance.
 Launching "[Performance-guaranteed](https://www.alibabacloud.com/help/doc-detail/27657.htm)" instance, it is must be specified and it valid values are: "slb.s1.small", "slb.s2.small", "slb.s2.medium",
 "slb.s3.small", "slb.s3.medium", "slb.s3.large" and "slb.s4.large".
* `tags` - (Optional) A mapping of tags to assign to the resource. The `tags` can have a maximum of 10 tag for every load balancer instance.
* `instance_charge_type` - (Optional, Available in v1.34.0+) The billing method of the load balancer. Valid values are "PrePaid" and "PostPaid". Default to "PostPaid".
* `period` - (Optional, Available in v1.34.0+) The duration that you will buy the resource, in month. It is valid when `instance_charge_type` is `PrePaid`. Valid values: [1-9, 12, 24, 36].
-> **NOTE:** The attribute `period` is only used to create Subscription instance or modify the PayAsYouGo instance to Subscription. Once effect, it will not be modified that means running `terraform apply` will not effect the resource.
* `master_zone_id` - (Optional, ForceNew, Available in v1.36.0+) The primary zone ID of the SLB instance. If not specified, the system will be randomly assigned. You can query the primary and standby zones in a region by calling the DescribeZone API.
* `slave_zone_id` - (Optional, ForceNew, Available in v1.36.0+) The standby zone ID of the SLB instance. If not specified, the system will be randomly assigned. You can query the primary and standby zones in a region by calling the DescribeZone API.
* `delete_protection` - (Optional, Available in v1.51.0+) Whether enable the deletion protection or not. on: Enable deletion protection. off: Disable deletion protection. Default to off. Only postpaid instance support this function.   
* `address_ip_version` - (Optional, Available in v1.55.2+) The IP version of the SLB instance to be created, which can be set to ipv4 or ipv6 . Default to "ipv4". Now, only internet instance support ipv6 address.
* `address` - (Optional, Available in v1.55.2+) Specify the IP address of the private network for the SLB instance, which must be in the destination CIDR block of the correspond ing switch.
* `resource_group_id` - (Optional, ForceNew, Available in v1.55.3+) The Id of resource group which the SLB belongs.

-> **NOTE:** A "Shared-Performance" instance can be changed to "Performance-guaranteed", but the change is irreversible.

-> **NOTE:** To change a "Shared-Performance" instance to a "Performance-guaranteed" instance, the SLB will have a short probability of business interruption (10 seconds-30 seconds). Advise to change it during the business downturn, or migrate business to other SLB Instances by using GSLB before changing.

-> **NOTE:** Currently, the alibaba cloud international account does not support creating a PrePaid SLB instance.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the load balancer.
* `address` - The IP address of the load balancer.

## Import

Load balancer can be imported using the id, e.g.

```
$ terraform import alicloud_slb_load_balancer.example lb-abc123456
```
