---
layout: "alicloud"
page_title: "Alicloud: alicloud_network_interfaces"
sidebar_current: "docs-alicloud-datasource-network-interfaces"
description: |-
  Provides a data source to get a list of elastic network interfaces according to the specified filters.
---

# alicloud\_network_interfaces

Use this data source to get a list of elastic network interfaces according to the specified filters in an Alibaba Cloud account.

For information about elastic network interface and how to use it, see [Elastic Network Interface](https://www.alibabacloud.com/help/doc-detail/58496.html)

## Example Usage

```
data "alicloud_network_interfaces" "enis"  {
	ids = ["${alicloud_network_interface.eni.id}"]
	name_regex = "${alicloud_network_interface.eni.name}"
	vpc_id = "${alicloud_vpc.vpc.id}"
	vswitch_id = "${alicloud_vswitch.vswitch.id}"
	security_group_id = "${alicloud_security_group.sg.id}"
	name = "${alicloud_network_interface.eni.name}"
	tags = {
		TF-VER = "0.11.3"
	}
}

output "eni0_name" {
    value = "${data.alicloud_network_interfaces.enis.interfaces.0.name}"
}
```

##  Argument Reference

The following arguments are supported:

* `ids` - (Optional)  A list of ENI IDs.
* `name_regex` - (Optional) A regex string to filter results by ENI name.
* `vpc_id` - (Optional) The VPC ID linked to ENIs.
* `vswitch_id` - (Optional) The VSwitch ID linked to ENIs.
* `private_ip` - (Optional) The primary private IP address of the ENI.
* `security_group_id` - (Optional) The security group ID linked to ENIs.
* `name` - (Optional) The name of the ENIs.
* `type` - (Optional) The type of ENIs, Only support for "Primary" or "Secondary".
* `instance_id` - (Optional) The ECS instance ID that the ENI is attached to.
* `tags` - (Optional) A map of tags assigned to ENIs.
* `output_file` - (Optional) The name of output file that saves the filter results.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `interfaces` - A list of ENIs. Each element contains the following attributes:
    * `id` - ID of the ENI.
    * `status` - Current status of the ENI.
    * `vpc_id` - ID of the VPC that the ENI belongs to.
    * `vswitch_id` - ID of the VSwitch that the ENI is linked to.
    * `zone_id` - ID of the availability zone that the ENI belongs to.
    * `public_ip` - Public IP of the ENI.
    * `private_ip` - Primary private IP of the ENI.
    * `private_ips` - A list of secondary private IP address that is assigned to the ENI.
    * `mac` - MAC address of the ENI.
    * `security_groups` - A list of security group that the ENI belongs to.
    * `name` - Name of the ENI.
    * `description` - Description of the ENI.
    * `instance_id` - ID of the instance that the ENI is attached to.
    * `creation_time` - Creation time of the ENI.
    * `tags` - A map of tags assigned to the ENI.