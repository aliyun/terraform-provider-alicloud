---
subcategory: "Application Real-Time Monitoring Service (ARMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_arms_env_feature"
description: |-
  Provides a Alicloud ARMS Env Feature resource.
---

# alicloud_arms_env_feature

Provides a ARMS Env Feature resource. Feature of the arms environment.

For information about ARMS Env Feature and how to use it, see [What is Env Feature](https://www.alibabacloud.com/help/en/).

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
  description = "api-resource-test1-hz-feature"
  cidr_block  = "172.16.0.0/12"
  vpc_name    = var.name

}

resource "alicloud_vswitch" "vsw" {
  description  = "api-resource-test1-hz"
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

resource "alicloud_arms_environment" "env-feature" {
  environment_type = "CS"
  environment_name = var.name

  bind_resource_id     = alicloud_ack_cluster.ask.id
  environment_sub_type = "ACK"
  tags {
    tag_key   = "api-cs-k1"
    tag_value = "api-cs-v1"
  }
}


resource "alicloud_arms_env_feature" "default" {
  env_feature_name = var.name

  environment_id  = alicloud_arms_environment.env-feature.id
  feature_version = "1.1.17"
}
```

## Argument Reference

The following arguments are supported:
* `aliyun_lang` - (Optional, ForceNew) The locale. The default is Chinese zh | en.
* `config` - (Optional, ForceNew) The configuration information of the Feature.
* `env_feature_name` - (Required, ForceNew) The name of the resource.
* `environment_id` - (Required, ForceNew) The first ID of the resource.
* `feature_version` - (Required) Version information of the Feature. You can query Feature information by using ListEnvironmentFeatures.
* `region` - (Optional) Feature Region to be installed.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<environment_id>:<env_feature_name>`.
* `namespace` - Namespace.
* `status` - Status.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Env Feature.
* `delete` - (Defaults to 5 mins) Used when delete the Env Feature.
* `update` - (Defaults to 5 mins) Used when update the Env Feature.

## Import

ARMS Env Feature can be imported using the id, e.g.

```shell
$ terraform import alicloud_arms_env_feature.example <environment_id>:<env_feature_name>
```