---
subcategory: "Bastion Host"
layout: "alicloud"
page_title: "Alicloud: alicloud_bastionhost_host_group"
sidebar_current: "docs-alicloud-resource-bastionhost-host-group"
description: |-
  Provides a Alicloud Bastion Host Host Group resource.
---

# alicloud_bastionhost_host_group

Provides a Bastion Host Host Group resource.

For information about Bastion Host Host Group and how to use it, see [What is Host Group](https://www.alibabacloud.com/help/en/doc-detail/204307.htm).

-> **NOTE:** Available since v1.134.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_bastionhost_host_group&exampleId=8a03897d-8bbb-e152-8d2b-4023f0dd9350a83b2c02&activeTab=example&spm=docs.r.bastionhost_host_group.0.8a03897d8b&intl_lang=EN_US" target="_blank">
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
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_bastionhost_host_group&spm=docs.r.bastionhost_host_group.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `comment` - (Optional) Specify the New Host Group of Notes, Supports up to 500 Characters.
* `host_group_name` - (Required) Specify the New Host Group Name, Supports up to 128 Characters.
* `instance_id` - (Required, ForceNew) Specify the New Host Group Where the Bastion Host ID of.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Host Group. The value formats as `<instance_id>:<host_group_id>`.
* `host_group_id` - Host Group ID.

## Import

Bastion Host Host Group can be imported using the id, e.g.

```shell
$ terraform import alicloud_bastionhost_host_group.example <instance_id>:<host_group_id>
```
