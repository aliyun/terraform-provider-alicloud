---
layout: "alicloud"
page_title: "Alicloud: alicloud_network_interface"
sidebar_current: "docs-alicloud-resource-network-interface"
description: |-
	Provides an ECS Elastic Network Interface resource.
---

# alicloud\_network\_interface

Provides an ECS Elastic Network Interface resource.

For information about Elastic Network Interface and how to use it, see [Elastic Network Interface](https://www.alibabacloud.com/help/doc-detail/58496.html).

~> **NOTE** Only one of private_ips or private_ips_count can be specified when assign private IPs. 

## Example Usage

```
resource "alicloud_networt_interface" "eni0" {
    name = "terraform-test-eni0"
    vswitch_id = "${alicloud_vswith.vswith.id}"
    security_groups = [ "${alicloud_security_group.sg.id}" ]
    private_ips = [ "192.168.*.2", "192.168.*.3", "192.168.*.4" ]
}

resource "alicloud_networt_interface" "eni1" {
    name = "terraform-test-eni1"
    vswitch_id = "{alicloud_vswith.vswith.id}"
    primary_ip_address = "192.168.*.8"
    security_groups = [ "${alicloud_security_group.sg.id}" ]
    private_ips = [ "192.168.*.5", "192.168.*.6", "192.168.*.7" ]
}

resource "alicloud_networt_interface" "eni2" {
    name = "terraform-test-eni2"
    vswitch_id = "{alicloud_vswith.vswith.id}"
    security_groups = [ "${alicloud_security_group.sg.id}" ]
    private_ips_count = 10
}
```

## Argument Reference

The following arguments are supported:

* `vswitch_id` - (Required, ForceNew) The VSwitch to create the ENI in.
* `security_groups` - (Require) A list of security group ids to associate with.
* `private_ip` - (Optional) The primary private IP of the ENI.
* `name` - (Optional) Name of the ENI. This name can have a string of 2 to 128 characters, must contain only alphanumeric characters or hyphens, such as "-", ".", "_", and must not begin or end with a hyphen, and must not begin with http:// or https://. Default value is null.
* `description` - (Optional) Description of the ENI. This description can have a string of 2 to 256 characters, It cannot begin with http:// or https://. Default value is null.
* `private_ips`  - (Optional) List of secondary private IPs to assign to the ENI. Don't use both private_ips and private_ips_count in the same ENI resource block.
* `private_ips_count` - (Optional) Number of secondary private IPs to assign to the ENI. Don't use both private_ips and private_ips_count in the same ENI resource block.
* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported:

* `id` - The ENI ID.

## Import

ENI can be imported using the id, e.g.

```
$ terraform import alicloud_network_interface.eni eni-abc1234567890000
```
