---
subcategory: "Bastion Host"
layout: "alicloud"
page_title: "Alicloud: alicloud_bastionhost_host"
sidebar_current: "docs-alicloud-resource-bastionhost-host"
description: |-
  Provides a Alicloud Bastion Host Host resource.
---

# alicloud\_bastionhost\_host

Provides a Bastion Host Host resource.

For information about Bastion Host Host and how to use it, see [What is Host](https://www.alibabacloud.com/help/en/doc-detail/201330.htm).

-> **NOTE:** Available in v1.135.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_bastionhost_host" "example" {
  active_address_type  = "Private"
  host_private_address = "172.16.0.10"
  host_name            = "example_value"
  instance_id          = "bastionhost-cn-tl3xxxxxxx"
  os_type              = "Linux"
  source               = "Local"
}

```

## Argument Reference

The following arguments are supported:

* `active_address_type` - (Required) Specify the new create a host of address types. Valid values: `Public`: the IP address of a Public network. `Private`: Private network address.
* `comment` - (Optional) Specify a host of notes, supports up to 500 characters.
* `host_name` - (Required) Specify the new create a host name of the supports up to 128 characters.
* `host_private_address` - (Optional) Specify the new create a host of the private network address, it is possible to use the domain name or IP ADDRESS. **NOTE:**  This parameter is required if the `active_address_type` parameter is set to `Private`.
* `host_public_address` - (Optional) Specify the new create a host of the IP address of a public network, it is possible to use the domain name or IP ADDRESS.
* `instance_id` - (Required, ForceNew) Specify the new create a host where the Bastion host ID of.
* `instance_region_id` - (Optional) The instance region id.
* `os_type` - (Required) Specify the new create the host's operating system. Valid values: `Linux`,`Windows`.
* `source` - (Required, ForceNew) Specify the new create a host of source. Valid values: 
  * `Local`: localhost 
  * `Ecs`:ECS instance 
  * `Rds`:RDS exclusive cluster host.
* `source_instance_id` - (Optional, ForceNew) Specify the newly created ECS instance ID or dedicated cluster host ID. **NOTE:** This parameter is required if the `source` parameter is set to `Ecs` or `Rds`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Host. The value formats as `<instance_id>:<host_id>`.
* `host_id` - The host ID.

## Import

Bastion Host Host can be imported using the id, e.g.

```
$ terraform import alicloud_bastionhost_host.example <instance_id>:<host_id>
```
