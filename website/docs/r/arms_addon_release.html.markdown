---
subcategory: "Application Real-Time Monitoring Service (ARMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_arms_addon_release"
description: |-
  Provides a Alicloud ARMS Addon Release resource.
---

# alicloud_arms_addon_release

Provides a ARMS Addon Release resource. Release package of observability addon.

For information about ARMS Addon Release and how to use it, see [What is Addon Release](https://www.alibabacloud.com/help/en/arms/developer-reference/api-arms-2019-08-08-installaddon).

-> **NOTE:** Available since v1.212.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_arms_addon_release&exampleId=2fd51354-9f62-72d2-fd2f-fe4dbe3fb92e86c56d77&activeTab=example&spm=docs.r.arms_addon_release.0.2fd513549f&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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
  key_pair_name = "terraform-example-${random_integer.default.result}"
}

resource "alicloud_cs_kubernetes_node_pool" "default" {
  node_pool_name       = "desired_size"
  cluster_id           = alicloud_cs_managed_kubernetes.default.id
  vswitch_ids          = [alicloud_vswitch.vswitch.id]
  instance_types       = [data.alicloud_instance_types.default.instance_types.0.id]
  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  key_name             = alicloud_key_pair.default.key_pair_name
  desired_size         = 2
}

resource "alicloud_arms_environment" "default" {
  bind_resource_id     = alicloud_cs_kubernetes_node_pool.default.cluster_id
  environment_sub_type = "ManagedKubernetes"
  environment_type     = "CS"
  environment_name     = "terraform-example-${random_integer.default.result}"
}

resource "alicloud_arms_addon_release" "default" {
  aliyun_lang    = "zh"
  addon_name     = "mysql"
  environment_id = alicloud_arms_environment.default.id
  addon_version  = "0.0.1"
  values = jsonencode(
    {
      host     = "mysql-service.default"
      password = "roots"
      port     = 3306
      username = "root"
    }
  )
}
```

## Argument Reference

The following arguments are supported:
* `addon_name` - (Required, ForceNew) Addon Name.
* `addon_release_name` - (Optional, ForceNew, Computed) The name of the resource.
* `addon_version` - (Required) Version number of Addon. Addon information can be obtained through ListAddons.
* `aliyun_lang` - (Optional, ForceNew, Computed) The installed locale.
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