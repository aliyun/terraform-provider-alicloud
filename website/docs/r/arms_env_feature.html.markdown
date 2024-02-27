---
subcategory: "Application Real-Time Monitoring Service (ARMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_arms_env_feature"
description: |-
  Provides a Alicloud ARMS Env Feature resource.
---

# alicloud_arms_env_feature

Provides a ARMS Env Feature resource. Feature of the arms environment.

For information about ARMS Env Feature and how to use it, see [What is Env Feature](https://www.alibabacloud.com/help/en/arms/developer-reference/api-arms-2019-08-08-installenvironmentfeature).

-> **NOTE:** Available since v1.212.0.

## Example Usage

Basic Usage

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

resource "random_integer" "default" {
  max = 99999
  min = 10000
}

variable "name" {
  default = "terraform-example"
}

data "alicloud_enhanced_nat_available_zones" "enhanced" {
}

resource "alicloud_vpc" "vpc" {
  description = var.name
  cidr_block  = "192.168.0.0/16"
  vpc_name    = var.name
}

resource "alicloud_vswitch" "vswitch" {
  description  = var.name
  vpc_id       = alicloud_vpc.vpc.id
  vswitch_name = var.name
  zone_id      = data.alicloud_enhanced_nat_available_zones.enhanced.zones.0.zone_id
  cidr_block   = cidrsubnet(alicloud_vpc.vpc.cidr_block, 8, 8)
}

resource "alicloud_snapshot_policy" "default" {
  name            = var.name
  repeat_weekdays = ["1", "2", "3"]
  retention_days  = -1
  time_points     = ["1", "22", "23"]
}

data "alicloud_instance_types" "default" {
  availability_zone    = alicloud_vswitch.vswitch.zone_id
  cpu_core_count       = 2
  memory_size          = 4
  kubernetes_node_role = "Worker"
  instance_type_family = "ecs.sn1ne"
}

resource "alicloud_cs_managed_kubernetes" "default" {
  name               = "terraform-example-${random_integer.default.result}"
  cluster_spec       = "ack.pro.small"
  version            = "1.24.6-aliyun.1"
  new_nat_gateway    = true
  node_cidr_mask     = 26
  proxy_mode         = "ipvs"
  service_cidr       = "172.23.0.0/16"
  pod_cidr           = "10.95.0.0/16"
  worker_vswitch_ids = [alicloud_vswitch.vswitch.id]
}

resource "alicloud_key_pair" "default" {
  key_pair_name = "${var.name}-${random_integer.default.result}"
}

resource "alicloud_cs_kubernetes_node_pool" "default" {
  name                 = "desired_size"
  cluster_id           = alicloud_cs_managed_kubernetes.default.id
  vswitch_ids          = [alicloud_vswitch.vswitch.id]
  instance_types       = [data.alicloud_instance_types.default.instance_types.0.id]
  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  key_name             = alicloud_key_pair.default.key_name
  desired_size         = 2
}

resource "alicloud_arms_environment" "default" {
  bind_resource_id     = alicloud_cs_kubernetes_node_pool.default.cluster_id
  environment_sub_type = "ACK"
  environment_type     = "CS"
  environment_name     = "terraform-example-${random_integer.default.result}"
}

resource "alicloud_arms_env_feature" "default" {
  env_feature_name = "metric-agent"

  environment_id  = alicloud_arms_environment.default.id
  feature_version = "1.1.17"
}
```

## Argument Reference

The following arguments are supported:
* `env_feature_name` - (Required, ForceNew) The name of the resource.
* `environment_id` - (Required, ForceNew) The first ID of the resource.
* `feature_version` - (Required) Version information of the Feature. You can query Feature information by using ListEnvironmentFeatures.

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