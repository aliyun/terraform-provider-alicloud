---
subcategory: "Bastion Host"
layout: "alicloud"
page_title: "Alicloud: alicloud_bastionhost_host"
sidebar_current: "docs-alicloud-resource-bastionhost-host"
description: |-
  Provides a Alicloud Bastion Host Host resource.
---

# alicloud_bastionhost_host

Provides a Bastion Host Host resource.

For information about Bastion Host Host and how to use it, see [What is Host](https://www.alibabacloud.com/help/en/doc-detail/201330.htm).

-> **NOTE:** Available since v1.135.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_bastionhost_host&exampleId=4ca63422-dd8c-93bf-845f-e113fe2731b3b8ed1b24&activeTab=example&spm=docs.r.bastionhost_host.0.4ca63422dd&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf_example"
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
  cidr_block = "10.4.0.0/16"
}

data "alicloud_vswitches" "default" {
  cidr_block = "10.4.0.0/24"
  vpc_id     = data.alicloud_vpcs.default.ids.0
  zone_id    = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_security_group" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_bastionhost_instance" "default" {
  description        = var.name
  license_code       = "bhah_ent_50_asset"
  plan_code          = "cloudbastion"
  storage            = "5"
  bandwidth          = "5"
  period             = "1"
  vswitch_id         = data.alicloud_vswitches.default.ids[0]
  security_group_ids = [alicloud_security_group.default.id]
}

resource "alicloud_bastionhost_host" "default" {
  instance_id          = alicloud_bastionhost_instance.default.id
  host_name            = var.name
  active_address_type  = "Private"
  host_private_address = "172.16.0.10"
  os_type              = "Linux"
  source               = "Local"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_bastionhost_host&spm=docs.r.bastionhost_host.example&intl_lang=EN_US)

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

```shell
$ terraform import alicloud_bastionhost_host.example <instance_id>:<host_id>
```
