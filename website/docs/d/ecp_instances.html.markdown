---
subcategory: "Elastic Cloud Phone (ECP)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ecp_instances"
sidebar_current: "docs-alicloud-datasource-ecp-instances"
description: |-
  Provides a list of Ecp Instances to the user.
---

# alicloud\_ecp\_instances

This data source provides the Ecp Instances of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.158.0+.

## Example Usage

Basic Usage

```terraform
data "alicloud_ecp_zones" "default" {
}

data "alicloud_ecp_instance_types" "default" {
}

locals {
  count_size               = length(data.alicloud_ecp_zones.default.zones)
  zone_id                  = data.alicloud_ecp_zones.default.zones[local.count_size - 1].zone_id
  instance_type_count_size = length(data.alicloud_ecp_instance_types.default.instance_types)
  instance_type            = data.alicloud_ecp_instance_types.default.instance_types[local.instance_type_count_size - 1].instance_type
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = local.zone_id
}

resource "alicloud_security_group" "group" {
  name   = var.name
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_ecp_key_pair" "default" {
  key_pair_name   = var.name
  public_key_body = "ssh-rsa AAAAB3Nza12345678qwertyuudsfsg"
}

resource "alicloud_ecp_instance" "default" {
  instance_name = var.name
  description   = var.name
  force         = "true"
  key_pair_name = "${alicloud_ecp_key_pair.default.key_pair_name}"
  vswitch_id    = "${data.alicloud_vswitches.default.ids.0}"
  image_id      = "android_9_0_0_release_2851157_20211201.vhd"
  instance_type = "${data.alicloud_ecp_instance_types.default.instance_types[local.instance_type_count_size - 1].instance_type}"
  payment_type  = "PayAsYouGo"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Ecp Instances IDs.
* `names` - (Optional, ForceNew, Computed)  A list of Ecp Instances ames.
* `image_id` - (Required, ForceNew) The ID Of The Image.
* `payment_type` - (Optional) The payment type.Valid values: `PayAsYouGo`,`Subscription`
* `instance_name` - (Optional) The name of the instance. It must be 2 to 128 characters in length and must start with an
  uppercase letter or Chinese. It cannot start with http:// or https. It can contain Chinese, English, numbers,
  half-width colons (:), underscores (_), half-width periods (.), or dashes (-). The default value is the InstanceId of
  the instance.
* `instance_type` - (Required, ForceNew) Instance Type.
* `key_pair_name` - (Optional) The name of the key pair of the mobile phone instance.
* `resolution` - (Optional, ForceNew) The selected resolution for the cloud mobile phone instance.
* `status` - (Optional, Computed) Instance status. Valid values: `Pending`, `Running`, `Starting`, `Stopped`, `Stopping`
  .
* `name_regex` - (Optional, ForceNew) A regex string to filter results by mobile phone name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Instance names.
* `instances` - A list of Ecp Instances. Each element contains the following attributes:
    * `description` - Instance Description.
    * `id` - The ID of the Instance.
    * `image_id` - The ID Of The Image.
    * `instance_id` - InstanceId.
    * `instance_name` - Instance Name.
    * `instance_type` - Instance Type.
    * `key_pair_name` - The Key Name.
    * `resolution` - Resolution.
    * `security_group_id` - Security Group ID.
    * `status` - Instance Status.
    * `vnc_url` - VNC login address.
    * `vswitch_id` - The vswitch id.