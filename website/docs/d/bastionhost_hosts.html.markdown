---
subcategory: "Bastion Host"
layout: "alicloud"
page_title: "Alicloud: alicloud_bastionhost_hosts"
sidebar_current: "docs-alicloud-datasource-bastionhost-hosts"
description: |-
  Provides a list of Bastionhost Hosts to the user.
---

# alicloud\_bastionhost\_hosts

This data source provides the Bastionhost Hosts of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.135.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_bastionhost_hosts" "ids" {
  instance_id = "example_value"
  ids         = ["1", "2"]
}
output "bastionhost_host_id_1" {
  value = data.alicloud_bastionhost_hosts.ids.hosts.0.id
}

data "alicloud_bastionhost_hosts" "nameRegex" {
  instance_id = "example_value"
  name_regex  = "^my-Host"
}
output "bastionhost_host_id_2" {
  value = data.alicloud_bastionhost_hosts.nameRegex.hosts.0.id
}

```

## Argument Reference

The following arguments are supported:

* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `host_address` - (Optional, ForceNew) The host address.
* `host_name` - (Optional, ForceNew) Specify the new create a host name of the supports up to 128 characters.
* `ids` - (Optional, ForceNew, Computed)  A list of Host IDs.
* `instance_id` - (Required, ForceNew) Specify the new create a host where the Bastion host ID of.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Host name.
* `os_type` - (Optional, ForceNew) Specify the new create the host's operating system. Valid values: Linux Windows.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `source` - (Optional, ForceNew) Specify the new create a host of source. Valid values: Local: localhost Ecs:ECS instance Rds:RDS exclusive cluster host.
* `source_instance_id` - (Optional, ForceNew) Specify the newly created ECS instance ID or dedicated cluster host ID.
* `source_instance_state` - (Optional, ForceNew) The source instance state.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Host names.
* `hosts` - A list of Bastionhost Hosts. Each element contains the following attributes:
	* `active_address_type` - Specify the new create a host of address types. Valid values: Public: the IP address of a Public network Private: Private network address.
	* `comment` - Specify a host of notes, supports up to 500 characters.
	* `host_id` - The host ID.
	* `host_name` - Specify the new create a host name of the supports up to 128 characters.
	* `host_private_address` - Specify the new create a host of the private network address, it is possible to use the domain name or IP ADDRESS.
	* `host_public_address` - Specify the new create a host of the IP address of a public network, it is possible to use the domain name or IP ADDRESS.
	* `id` - The ID of the Host.
	* `instance_id` - Specify the new create a host where the Bastion host ID of.
	* `os_type` - Specify the new create the host's operating system. Valid values: Linux Windows.
	* `protocols` - The host of the protocol information.
		* `host_finger_print` - Host fingerprint information, it is possible to uniquely identify a host.
		* `port` - Host the service port of the RDS.
		* `protocol_name` - The host uses the protocol name. 
	* `source` - Specify the new create a host of source. Valid values: Local: localhost Ecs:ECS instance Rds:RDS exclusive cluster host.
	* `source_instance_id` - Specify the newly created ECS instance ID or dedicated cluster host ID.
