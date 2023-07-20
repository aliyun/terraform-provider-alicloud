---
subcategory: "Container Service for Kubernetes (ACK)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_kubernetes_addon"
sidebar_current: "docs-alicloud-resource-cs-kubernetes-addon"
description: |-
  Provides a Alicloud resource to manage container kubernetes addon.
---

# alicloud_cs_kubernetes_addon

This resource will help you to manage addon in Kubernetes Cluster, see [What is kubernetes addon](https://www.alibabacloud.com/help/en/ack/ack-managed-and-ack-dedicated/developer-reference/api-install-a-component-in-an-ack-cluster).

-> **NOTE:** Available since v1.150.0.

-> **NOTE:** From version 1.166.0, support specifying addon customizable configuration.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_cs_managed_kubernetes" "default" {
  name_prefix          = var.name
  cluster_spec         = "ack.pro.small"
  worker_vswitch_ids   = [alicloud_vswitch.default.id]
  new_nat_gateway      = true
  pod_cidr             = cidrsubnet("10.0.0.0/8", 8, 36)
  service_cidr         = cidrsubnet("172.16.0.0/16", 4, 7)
  slb_internet_enabled = true
}

resource "alicloud_cs_kubernetes_addon" "default" {
  cluster_id = alicloud_cs_managed_kubernetes.default.id
  name       = "ack-node-problem-detector"
  version    = "1.2.7"
}

# First, check the `next_version` field of the addon that can be upgraded to through the `.tfstate file`, then overwrite the `version` field with the value of `next_version` and apply.
# upgrade from 1.2.7 to 1.2.8
resource "alicloud_cs_kubernetes_addon" "ack-addon-upgrade" {
  cluster_id = alicloud_cs_kubernetes_addon.default.cluster_id
  name       = "ack-node-problem-detector"
  version    = "1.2.8"
}

resource "alicloud_cs_kubernetes_addon" "nginx_ingress_controller" {
  cluster_id = alicloud_cs_kubernetes_addon.ack-addon-upgrade.cluster_id
  name       = "nginx-ingress-controller"
  version    = "v1.1.2-aliyun.2"
  // Specify custom configuration for addon. You can checkout the customizable configuration of the addon through data source alicloud_cs_kubernetes_addon_metadata.
  config = jsonencode(
    {
      CpuLimit              = ""
      CpuRequest            = "100m"
      EnableWebhook         = true
      HostNetwork           = false
      IngressSlbNetworkType = "internet"
      IngressSlbSpec        = "slb.s2.small"
      MemoryLimit           = ""
      MemoryRequest         = "200Mi"
      NodeSelector          = []
    }
  )
}
```
**Installing of addon**
When a cluster is created, some system addons and those specified at the time of cluster creation will be installed, so when an addon resource is applied:
* If the addon already exists in the cluster and its version is the same as the specified version, it will be skipped and will not be reinstalled.
* If the addon already exists in the cluster and its version is different from the specified version, the addon will be upgraded.
* If the addon does not exist in the cluster, it will be installed.

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, ForceNew) The id of kubernetes cluster.
* `name` - (Required, ForceNew) The name of addon.
* `version` - (Required) The current version of addon.
* `config` - (Optional, Available since v1.166.0) The custom configuration of addon. You can checkout the customizable configuration of the addon through datasource `alicloud_cs_kubernetes_addon_metadata`, the returned format is the standard json schema. If return empty, it means that the addon does not support custom configuration yet. You can also checkout the current custom configuration through the data source `alicloud_cs_kubernetes_addons`.

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

```shell
$ terraform import alicloud_cs_kubernetes_addon.my_addon <cluster_id>:<addon_name>
```
