---
subcategory: "ENS"
layout: "alicloud"
page_title: "Alicloud: alicloud_ens_network_interfaces"
sidebar_current: "docs-alicloud-datasource-ens-network-interfaces"
description: |-
  Provides a list of ENS Network Interface owned by an Alibaba Cloud account.
---

# alicloud_ens_network_interfaces

This data source provides ENS Network Interface available to the user.[What is Network Interface](https://next.api.alibabacloud.com/document/Ens/2017-11-10/CreateNetworkInterface)

-> **NOTE:** Available since v1.287.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = ""
}

variable "ens_region_id" {
  default = "cn-hangzhou-58"
}

resource "alicloud_ens_security_group" "defaultvdnzzs" {
  security_group_name = "弹性网卡测试用例使用"
}

resource "alicloud_ens_network" "defaultsLhpIw" {
  network_name  = "弹性网卡测试用例"
  cidr_block    = "10.0.0.0/8"
  ens_region_id = var.ens_region_id
}

resource "alicloud_ens_vswitch" "defaultwsMJ1N" {
  cidr_block    = "10.0.5.0/24"
  vswitch_name  = "弹性网卡测试用例"
  ens_region_id = var.ens_region_id
  network_id    = alicloud_ens_network.defaultsLhpIw.id
}


resource "alicloud_ens_network_interface" "default" {
  description            = "desc"
  network_interface_name = "弹性网卡测试用例"
  security_group_ids     = ["${alicloud_ens_security_group.defaultvdnzzs.id}"]
  vswitch_id             = alicloud_ens_vswitch.defaultwsMJ1N.id
}

data "alicloud_ens_network_interfaces" "default" {
  ids                    = ["${alicloud_ens_network_interface.default.id}"]
  name_regex             = alicloud_ens_network_interface.default.network_interface_name
  network_interface_name = "弹性网卡测试用例"
  vswitch_id             = alicloud_ens_vswitch.defaultwsMJ1N.id
}

output "alicloud_ens_network_interface_example_id" {
  value = data.alicloud_ens_network_interfaces.default.interfaces.0.id
}
```

## Argument Reference

The following arguments are supported:
* `ens_region_id` - (ForceNew, Optional) The node ID of ENS.
* `instance_id` - (ForceNew, Optional) The ID of the instance bound to the Eni.
* `network_id` - (ForceNew, Optional) The network ID.
* `network_interface_id` - (ForceNew, Optional) The ID of the Eni.
* `network_interface_name` - (ForceNew, Optional) The name of the ENI.
* `status` - (ForceNew, Optional) Status of network card, value:
  - Available: Available.
  - Attaching: Attaching.
  - InUse: Attached.
  - Detaching: in the process of Detaching.
  - Deleting: Deleting.
* `vswitch_id` - (ForceNew, Optional) The vSwitch ID.
* `ids` - (Optional, Computed) A list of Network Interface IDs.
* `name_regex` - (Optional) A regex string to filter results by Group Metric Rule name.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Network Interface IDs.
* `names` - A list of name of Network Interfaces.
* `interfaces` - A list of Network Interface Entries. Each element contains the following attributes:
  * `create_time` - Creation time.
  * `description` - The description of the ENI.
  * `ens_region_id` - The node ID of ENS.
  * `instance_id` - The ID of the instance bound to the Eni.
  * `mac_address` - The MAC address of the Eni.
  * `network_id` - The network ID.
  * `network_interface_id` - The ID of the Eni.
  * `network_interface_name` - The name of the ENI.
  * `primary_ip` - The primary private network IP address.
  * `primary_ip_type` - The primary IP address type.
  * `security_group_ids` - **NOTE:** This field is only available when `enable_details` is `true`. The ID of the security group.
  * `status` - Status of network card, value:.
  * `vswitch_id` - The vSwitch ID.
  * `vmnc_learn` - Whether to enable NIC route learning.
  * `id` - The ID of the resource supplied above.
