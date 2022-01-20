---
subcategory: "Container Service for Kubernetes (ACK)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_kubernetes_addon"
sidebar_current: "docs-alicloud-resource-cs-kubernetes-addon"
description: |-
  Provides a Alicloud resource to manage container kubernetes addon.
---

# alicloud\_cs\_kubernetes\_addon

This resource will help you to manage addon in Kubernetes Cluster. 

-> **NOTE:** Available in 1.150.0+.

## Example Usage

**Create a managed cluster**

```terraform
variable "name" {
  default = "tf-test"
}
data "alicloud_zones" default {
  available_resource_creation = "VSwitch"
}
data "alicloud_instance_types" "default" {
  availability_zone    = data.alicloud_zones.default.zones.0.id
  cpu_core_count       = 2
  memory_size          = 4
  kubernetes_node_role = "Worker"
}
resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.1.0.0/21"
}
resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "10.1.1.0/24"
  zone_id      = data.alicloud_zones.default.zones.0.id
}
resource "alicloud_key_pair" "default" {
  key_pair_name = var.name
}
resource "alicloud_cs_managed_kubernetes" "default" {
  name                         = var.name
  count                        = 1
  cluster_spec                 = "ack.pro.small"
  is_enterprise_security_group = true
  worker_number                = 2
  password                     = "Hello1234"
  pod_cidr                     = "172.20.0.0/16"
  service_cidr                 = "172.21.0.0/20"
  worker_vswitch_ids           = [alicloud_vswitch.default.id]
  worker_instance_types        = [data.alicloud_instance_types.default.instance_types.0.id]
}
```
**Installing of addon**
When a cluster is created, some system addons and those specified at the time of cluster creation will be installed, so when an addon resource is applied:
* If the addon already exists in the cluster and its version is the same as the specified version, it will be skipped and will not be reinstalled.
* If the addon already exists in the cluster and its version is different from the specified version, the addon will be upgraded.
* If the addon does not exist in the cluster, it will be installed.

```terraform
resource "alicloud_cs_kubernetes_addon" "ack-node-problem-detector" {
  cluster_id = alicloud_cs_managed_kubernetes.default.0.id
  name       = "ack-node-problem-detector"
  version    = "1.2.7"
}
```
**Upgrading of addon**
First, check the `next_version` field of the addon that can be upgraded to through the `.tfstate file`, then overwrite the `version` field with the value of `next_version` and apply.
```terraform
resource "alicloud_cs_kubernetes_addon" "ack-node-problem-detector" {
  cluster_id = alicloud_cs_managed_kubernetes.default.0.id
  name       = "ack-node-problem-detector"
  version    = "1.2.8" # upgrade from 1.2.7 to 1.2.8
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, ForceNew) The id of kubernetes cluster.
* `name` - (Required, ForceNew) The name of addon.
* `version` - (Required) The current version of addon.

## Attributes Reference

The following attributes are exported:
* `id` - The id of addon, which consists of the cluster id and the addon name, with the structure <cluster_ud>:<addon_name>.
* `next_version` - The version which addon can be upgraded to.
* `can_upgrade` - Is the addon ready for upgrade.
* `required` - Is it a mandatory addon to be installed.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when installing addon in the kubernetes cluster. 
* `update` - (Defaults to 10 mins) Used when upgrading addon in the kubernetes cluster.
* `delete` - (Defaults to 10 mins) Used when deleting addon in kubernetes cluster. 

## Import

Cluster addon can be imported by cluster id and addon name. Then write the addon.tf file according to the result of `terraform plan`.

```
  $ terraform import alicloud_cs_kubernetes_addon.my_addon <cluster_id>:<addon_name>
```
