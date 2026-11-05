---
subcategory: "Hologres (Hologram)"
layout: "alicloud"
page_title: "Alicloud: alicloud_hologram_holo_web_login_setting"
description: |-
  Provides a Alicloud Hologres (Hologram) Holo Web Login Setting resource.
---

# alicloud_hologram_holo_web_login_setting

Provides a Hologres (Hologram) Holo Web Login Setting resource.

The Hologres web login setting controls whether external accounts are allowed to log in to a Hologres instance through HoloWeb. It is a single boolean toggle bound to one instance. Because the Hologres API only exposes `GetHoloWebLoginSetting` and `UpdateHoloWebLoginSetting`, this resource manages the toggle in place: creating or updating the resource sets the configured value, and destroying it reverts the setting to `false` (the service default). The Hologres instance itself is not created or destroyed by this resource.

For information about Hologres (Hologram) Holo Web Login Setting and how to use it, see [GetHoloWebLoginSetting](https://www.alibabacloud.com/help/zh/hologres/developer-reference/api-hologram-2022-06-01-getholowebloginsetting) and [UpdateHoloWebLoginSetting](https://www.alibabacloud.com/help/zh/hologres/developer-reference/api-hologram-2022-06-01-updateholowebloginsetting).

-> **NOTE:** Available since v1.243.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_vpc" "defaultVpc" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = var.name
}

resource "alicloud_vswitch" "defaultVSwitch" {
  vpc_id       = alicloud_vpc.defaultVpc.id
  zone_id      = "cn-hangzhou-j"
  cidr_block   = "172.16.53.0/24"
  vswitch_name = var.name
}

resource "alicloud_hologram_instance" "default" {
  instance_type = "Standard"
  pricing_cycle = "Hour"
  cpu           = "32"
  zone_id       = alicloud_vswitch.defaultVSwitch.zone_id
  endpoints {
    type = "Intranet"
  }
  endpoints {
    type       = "VPCSingleTunnel"
    vswitch_id = alicloud_vswitch.defaultVSwitch.id
    vpc_id     = alicloud_vswitch.defaultVSwitch.vpc_id
  }
  payment_type = "PayAsYouGo"
}

resource "alicloud_hologram_holo_web_login_setting" "default" {
  instance_id                 = alicloud_hologram_instance.default.id
  allow_external_accounts_login = true
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the Hologres instance. The web login setting is unique per instance; changing the instance creates a new setting resource.
* `allow_external_accounts_login` - (Required) Whether to allow external accounts to log in to the instance through HoloWeb. Valid values: `true`, `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform state, equal to the `instance_id` of the Hologres instance.

## Import

Hologram Holo Web Login Setting can be imported using the id, e.g.

```shell
terraform import alicloud_hologram_holo_web_login_setting.example <instance_id>
```
