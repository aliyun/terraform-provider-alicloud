---
layout: "alicloud"
page_title: "Alicloud: alicloud_slb"
sidebar_current: "docs-alicloud-resource-slb"
description: |-
  Provides an Application Load Banlancer resource.
---

# alicloud\_slb

Provides an Application Load Balancer resource.

~> **NOTE:** Resource `alicloud_slb` has deprecated 'listener' filed from terraform-alicloud-provider [version 1.3.0](https://github.com/alibaba/terraform-provider/releases/tag/V1.3.0) . You can create new listeners for Load Balancer by resource `alicloud_slb_listener`.
If you have had several listeners in one load balancer, you can import them via the specified listener ID. In the `alicloud_slb_listener`, listener ID is consist of load balancer ID and frontend port, and its format is "<load balancer ID>:<frontend port>", like "lb-hr2fwnf32t:8080".

~> **NOTE:** At present, to avoid some unnecessary regulation confusion, SLB can not support alicloud international account to create "paybybandwidth" instance.

## Example Usage

```
# Create a new load balancer for classic
resource "alicloud_slb" "classic" {
  name                 = "test-slb-tf"
  internet             = true
  internet_charge_type = "paybybandwidth"
  bandwidth            = 5
}

# Create a new load balancer for VPC
resource "alicloud_vpc" "default" {
  # Other parameters...
}

resource "alicloud_vswitch" "default" {
  # Other parameters...
}

resource "alicloud_slb" "vpc" {
  name       = "test-slb-tf"
  vswitch_id = "${alicloud_vswitch.default.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the SLB. This name must be unique within your AliCloud account, can have a maximum of 80 characters,
must contain only alphanumeric characters or hyphens, such as "-","/",".","_", and must not begin or end with a hyphen. If not specified,
Terraform will autogenerate a name beginning with `tf-lb`.
* `internet` - (Optional, Forces New Resource) If true, the SLB addressType will be internet, false will be intranet, Default is false. If load balancer launched in VPC, this value must be "false".
* `internet_charge_type` - (Optional, Forces New Resource) Valid
  values are `paybybandwidth`, `paybytraffic`. If this value is "paybybandwidth", then argument "internet" must be "true". Default is "paybytraffic". If load balancer launched in VPC, this value must be "paybytraffic".
* `bandwidth` - (Optional) Valid
  value is between 1 and 1000, If argument "internet_charge_type" is "paybytraffic", then this value will be ignore.
* `listener` - (Deprecated) The field has been deprecated from terraform-alicloud-provider [version 1.3.0](https://github.com/alibaba/terraform-provider/releases/tag/V1.3.0), and use resource `alicloud_slb_listener` to replace.
* `vswitch_id` - (Required for a VPC SLB, Forces New Resource) The VSwitch ID to launch in.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the load balancer.
* `name` - The name of the load balancer.
* `internet` - The internet of the load balancer.
* `internet_charge_type` - The internet_charge_type of the load balancer.
* `bandwidth` - The bandwidth of the load balancer.
* `vswitch_id` - The VSwitch ID of the load balancer. Only available on SLB launched in a VPC.
* `address` - The IP address of the load balancer.