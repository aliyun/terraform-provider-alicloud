---
subcategory: "Application Real-Time Monitoring Service (ARMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_arms_addon_release"
description: |-
  Provides a Alicloud ARMS Addon Release resource.
---

# alicloud_arms_addon_release

Provides a ARMS Addon Release resource. Release package of observability addon.

For information about ARMS Addon Release and how to use it, see [What is Addon Release](https://www.alibabacloud.com/help/en/).

-> **NOTE:** Available since v1.212.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "vpc" {
  description = "api-resource-test1-hz-addonrelease"
  cidr_block  = "172.16.0.0/12"
  vpc_name    = var.name

}

resource "alicloud_vswitch" "vsw" {
  description  = "api-resource-test1-hz-addonrelease"
  vpc_id       = alicloud_vpc.vpc.id
  vswitch_name = var.name

  zone_id    = data.alicloud_zones.default.zones.0.id
  cidr_block = "172.16.0.0/24"
}

resource "alicloud_ack_cluster" "ask" {
  kubernetes_version = "1.26.3-aliyun.1"
  cluster_type       = "ManagedKubernetes"
  cluster_spec       = "ack.pro.small"
  vpc_id             = alicloud_vpc.vpc.id
  service_cidr       = "192.168.0.0/24"
  cluster_name       = var.name

  container_cidr = "192.168.1.0/24"
  vswitch_id     = alicloud_vswitch.vsw.id
  profile        = "Serverless"
}

resource "alicloud_arms_environment" "env-addonrelease" {
  environment_type = "CS"
  environment_name = var.name

  bind_resource_id     = alicloud_ack_cluster.ask.id
  environment_sub_type = "ACK"
  tags {
    tag_key   = "api-cs-k1"
    tag_value = "api-cs-v1"
  }
}


resource "alicloud_arms_addon_release" "default" {
  environment_id = alicloud_arms_environment.env-addonrelease.id
  addon_version  = "0.0.1"
  values         = "{\"host\":\"mysql-service.default\",\"port\":3306,\"username\":\"root\",\"password\":\"roots\"}"
}
```

## Argument Reference

The following arguments are supported:
* `addon_name` - (Required, ForceNew) Addon Name.
* `addon_release_name` - (Optional, ForceNew, Computed) The name of the resource.
* `addon_version` - (Required) Version number of Addon. Addon information can be obtained through ListAddons.
* `aliyun_lang` - (Optional, ForceNew) The installed locale.
* `environment_id` - (Required, ForceNew) Environment id.
* `values` - (Optional) Configuration information for installing Addon. Obtain the configuration template from ListAddonSchema, for example, {"host":"mysql-service.default","port":3306,"username":"root","password":"roots"}.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<environment_id>:<addon_release_name>`.
* `create_time` - Creation time.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Addon Release.
* `delete` - (Defaults to 5 mins) Used when delete the Addon Release.
* `update` - (Defaults to 5 mins) Used when update the Addon Release.

## Import

ARMS Addon Release can be imported using the id, e.g.

```shell
$ terraform import alicloud_arms_addon_release.example <environment_id>:<addon_release_name>
```