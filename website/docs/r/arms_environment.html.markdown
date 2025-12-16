---
subcategory: "Application Real-Time Monitoring Service (ARMS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_arms_environment"
description: |-
  Provides a Alicloud ARMS Environment resource.
---

# alicloud_arms_environment

Provides a ARMS Environment resource. The arms environment.

For information about ARMS Environment and how to use it, see [What is Environment](https://www.alibabacloud.com/help/en/arms/developer-reference/api-arms-2019-08-08-createenvironment).

-> **NOTE:** Available since v1.212.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_arms_environment&exampleId=cf792b6f-bfa7-3c78-3e01-14bce9a4bf5fc3886448&activeTab=example&spm=docs.r.arms_environment.0.cf792b6fbf&intl_lang=EN_US" target="_blank">
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
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_arms_environment&spm=docs.r.arms_environment.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `aliyun_lang` - (Optional) The locale. The default is Chinese zh | en.
* `bind_resource_id` - (Optional, ForceNew) The id or vpcId of the bound container instance.
* `drop_metrics` - (Optional) List of abandoned indicators.
* `environment_name` - (Optional) The name of the resource.
* `environment_sub_type` - (Required, ForceNew) Subtype of environment:
  - Type of CS: ACK is currently supported. ManagedKubernetes, Kubernetes, ExternalKubernetes, and One are also supported.
  - Type of ECS: currently supports ECS.
  - Type of Cloud: currently supports Cloud.
* `environment_type` - (Required, ForceNew) Type of environment.
* `managed_type` - (Optional, ForceNew) Hosting type:
  - none: unmanaged. The default value of the ACK cluster.
  - agent: Managed agent (including ksm). Default values of ASK, ACS, and Acone clusters.
  - agent-exproter: Managed agent and exporter. The default value of the cloud service type.
* `resource_group_id` - (Optional, Computed) The ID of the resource group.
* `tags` - (Optional, Map) The tag of the resource.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `environment_id` - The first ID of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Environment.
* `delete` - (Defaults to 5 mins) Used when delete the Environment.
* `update` - (Defaults to 5 mins) Used when update the Environment.

## Import

ARMS Environment can be imported using the id, e.g.

```shell
$ terraform import alicloud_arms_environment.example <id>
```