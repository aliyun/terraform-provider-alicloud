---
subcategory: "Bastion Host"
layout: "alicloud"
page_title: "Alicloud: alicloud_bastionhost_host_share_keys"
sidebar_current: "docs-alicloud-datasource-bastionhost-host-share-keys"
description: |-
  Provides a list of Bastionhost Host Share Keys to the user.
---

# alicloud\_bastionhost\_host\_share\_keys

This data source provides the Bastionhost Host Share Keys of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.165.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_bastionhost_host_share_keys" "ids" {
  instance_id = "example_value"
  ids         = ["example_value-1", "example_value-2"]
}
output "bastionhost_host_share_key_id_1" {
  value = data.alicloud_bastionhost_host_share_keys.ids.keys.0.id
}

data "alicloud_bastionhost_host_share_keys" "nameRegex" {
  instance_id = "example_value"
  name_regex  = "^my-HostShareKey"
}
output "bastionhost_host_share_key_id_2" {
  value = data.alicloud_bastionhost_host_share_keys.nameRegex.keys.0.id
}
```

## Argument Reference

The following arguments are supported:

* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed)  A list of Host Share Key IDs.
* `instance_id` - (Required, ForceNew) The ID of the Bastion instance.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Host Share Key name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Host Share Key names.
* `keys` - A list of Bastionhost Host Share Keys. Each element contains the following attributes:
	* `host_share_key_id` - The first ID of the resource.
	* `host_share_key_name` - The name of the host shared key.
	* `id` - The ID of the Host Share Key.
	* `instance_id` - The ID of the Bastion instance.
	* `private_key_finger_print` - The fingerprint of the private key.