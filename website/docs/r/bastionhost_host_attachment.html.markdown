---
subcategory: "Bastion Host"
layout: "alicloud"
page_title: "Alicloud: alicloud_bastionhost_host_attachment"
sidebar_current: "docs-alicloud-resource-bastionhost-host-attachment"
description: |-
  Provides a Alicloud Bastion Host Host Attachment resource.
---

# alicloud_bastionhost_host_attachment

Provides a Bastion Host Host Attachment resource to add host into one host group.

-> **NOTE:** Available since v1.135.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_bastionhost_host_attachment&exampleId=822d5c4a-7af2-dec5-6578-c9ab48363d139623a701&activeTab=example&spm=docs.r.bastionhost_host_attachment.0.822d5c4a7a&intl_lang=EN_US" target="_blank">
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

resource "alicloud_bastionhost_host_group" "default" {
  host_group_name = var.name
  instance_id     = alicloud_bastionhost_instance.default.id
}

resource "alicloud_bastionhost_host" "default" {
  instance_id          = alicloud_bastionhost_instance.default.id
  host_name            = var.name
  active_address_type  = "Private"
  host_private_address = "172.16.0.10"
  os_type              = "Linux"
  source               = "Local"
}

resource "alicloud_bastionhost_host_attachment" "default" {
  host_group_id = alicloud_bastionhost_host_group.default.host_group_id
  host_id       = alicloud_bastionhost_host.default.host_id
  instance_id   = alicloud_bastionhost_instance.default.id
}
```

## Argument Reference

The following arguments are supported:

* `host_group_id` - (Required, ForceNew) Specifies the added to the host group ID.
* `host_id` - (Required, ForceNew) Specified to be part of a host group of host ID.
* `instance_id` - (Required, ForceNew) The bastion host instance id.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Host Attachment. The value formats as `<instance_id>:<host_group_id>:<host_id>`.

## Import

Bastion Host Host Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_bastionhost_host_attachment.example <instance_id>:<host_group_id>:<host_id>
```
