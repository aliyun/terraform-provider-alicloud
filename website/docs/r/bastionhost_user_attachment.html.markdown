---
subcategory: "Bastion Host"
layout: "alicloud"
page_title: "Alicloud: alicloud_bastionhost_user_attachment"
sidebar_current: "docs-alicloud-resource-bastionhost-user-attachment"
description: |-
  Provides a Alicloud Bastion Host User Attachment resource.
---

# alicloud_bastionhost_user_attachment

Provides a Bastion Host User Attachment resource to add user to one user group.

-> **NOTE:** Available since v1.134.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_bastionhost_user_attachment&exampleId=fa9ede33-e08c-03dd-3e8b-005c50047fc1fef4807f&activeTab=example&spm=docs.r.bastionhost_user_attachment.0.fa9ede33e0&intl_lang=EN_US" target="_blank">
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

resource "alicloud_bastionhost_user_group" "default" {
  instance_id     = alicloud_bastionhost_instance.default.id
  user_group_name = var.name
}

resource "alicloud_bastionhost_user" "local_user" {
  instance_id         = alicloud_bastionhost_instance.default.id
  mobile_country_code = "CN"
  mobile              = "13312345678"
  password            = "YourPassword-123"
  source              = "Local"
  user_name           = "${var.name}_local_user"
}

resource "alicloud_bastionhost_user_attachment" "default" {
  instance_id   = alicloud_bastionhost_instance.default.id
  user_group_id = alicloud_bastionhost_user_group.default.user_group_id
  user_id       = alicloud_bastionhost_user.local_user.user_id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) Specifies the user group to add the user's bastion host ID of.
* `user_group_id` - (Required, ForceNew) Specifies the user group to which you want to add the user ID.
* `user_id` - (Required, ForceNew) Specify that you want to add to the policy attached to the user group ID. This includes response parameters in a Json-formatted string supports up to set up 100 USER ID.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of User Attachment. The value formats as `<instance_id>:<user_group_id>:<user_id>`.

## Import

Bastion Host User Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_bastionhost_user_attachment.example <instance_id>:<user_group_id>:<user_id>
```
