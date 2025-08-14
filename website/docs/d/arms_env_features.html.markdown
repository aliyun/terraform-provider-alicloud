---
subcategory: "Application Real-Time Monitoring Service (ARMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_arms_env_features"
description: |-
  Provides a list of ARMS Env Features to the user.
---

# alicloud_arms_env_features

This data source provides the ARMS Env Features of the current Alibaba Cloud user.

-> **NOTE:** Available since v1.258.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

resource "random_integer" "default" {
  min = 10000
  max = 99999
}

data "alicloud_enhanced_nat_available_zones" "enhanced" {
}

resource "alicloud_vpc" "vpc" {
  description = "api-resource-test1-hz"
  cidr_block  = "192.168.0.0/16"
  vpc_name    = "${var.name}-${random_integer.default.result}"
}

resource "alicloud_vswitch" "vswitch" {
  description  = "api-resource-test1-hz"
  vpc_id       = alicloud_vpc.vpc.id
  vswitch_name = "${var.name}-${random_integer.default.result}"
  zone_id      = data.alicloud_enhanced_nat_available_zones.enhanced.zones.0.zone_id
  cidr_block   = cidrsubnet(alicloud_vpc.vpc.cidr_block, 8, 8)
}

resource "alicloud_snapshot_policy" "default" {
  name            = "${var.name}-${random_integer.default.result}"
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
  name               = "${var.name}-${random_integer.default.result}"
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
  key_name             = alicloud_key_pair.default.key_pair_name
  desired_size         = 2
}

resource "alicloud_arms_environment" "default" {
  environment_type     = "CS"
  environment_name     = "${var.name}-${random_integer.default.result}"
  bind_resource_id     = alicloud_cs_kubernetes_node_pool.default.cluster_id
  environment_sub_type = "ManagedKubernetes"
}

resource "alicloud_arms_env_feature" "default" {
  env_feature_name = "metric-agent"
  environment_id   = alicloud_arms_environment.default.id
  feature_version  = "1.1.17"
}

data "alicloud_arms_env_features" "ids" {
  environment_id = alicloud_arms_env_feature.default.environment_id
  ids            = [alicloud_arms_env_feature.default.id]
}

output "arms_env_features_id_0" {
  value = data.alicloud_arms_env_features.ids.features.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, List) A list of ARMS Env Feature IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by ARMS Env Feature name.
* `environment_id` - (Required, ForceNew) The ID of the environment instance.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of ARMS Env Feature names.
* `features` - A list of ARMS Env Features. Each element contains the following attributes:
  * `id` - The ID of the Env Feature. It formats as `<environment_id>:<env_feature_name>`.
  * `aliyun_lang` - The language.
  * `env_feature_name` - The name of the feature.
  * `environment_id` - The ID of the environment instance.
  * `feature_version` - The version of the feature.
  * `status` - The status of the feature.
