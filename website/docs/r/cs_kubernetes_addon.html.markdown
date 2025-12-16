---
subcategory: "Container Service for Kubernetes (ACK)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_kubernetes_addon"
sidebar_current: "docs-alicloud-resource-cs-kubernetes-addon"
description: |-
  Provides a Alicloud resource to manage container kubernetes addon.
---

# alicloud_cs_kubernetes_addon

This resource will help you to manage addon in Kubernetes Cluster, see [What is kubernetes addon](https://www.alibabacloud.com/help/en/ack/ack-managed-and-ack-dedicated/developer-reference/api-install-a-component-in-an-ack-cluster). For more usage information, see [Use Terraform to manage addons](https://www.alibabacloud.com/help/en/ack/serverless-kubernetes/developer-reference/use-terraform-to-manage-components).

-> **NOTE:** Available since v1.150.0.

-> **NOTE:** From version 1.166.0, support specifying addon customizable configuration.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cs_kubernetes_addon&exampleId=10c0d82e-ea4b-da2e-f004-27e6e9ae89e1a5236d82&activeTab=example&spm=docs.r.cs_kubernetes_addon.0.10c0d82eea&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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
  new_nat_gateway      = false
  pod_cidr             = cidrsubnet("10.0.0.0/8", 8, 36)
  service_cidr         = cidrsubnet("172.16.0.0/16", 4, 7)
  slb_internet_enabled = true
  # By defining the addons attribute in cluster resource, it indicates that the addon will be installed when creating a cluster
  addons {
    name = "logtail-ds"
    # Specify the addon config. Or it can be written as
    # config = "{\"IngressDashboardEnabled\":\"true\"}
    config = jsonencode(
      {
        IngressDashboardEnabled = "true"
      }
    )
    # The default value of this parameter is false.Some addon will be installed by default to facilitate users to manage the cluster. If you do not need to install these addons when creating a cluster, you can set disabled=true.
    disabled = false
  }
}
# data source provides the information of available addons
data "alicloud_cs_kubernetes_addons" "default" {
  cluster_id = alicloud_cs_managed_kubernetes.default.id
  name_regex = "logtail-ds"
}
# Manage addon resource
resource "alicloud_cs_kubernetes_addon" "logtail-ds" {
  cluster_id = alicloud_cs_managed_kubernetes.default.id
  name       = "logtail-ds"
  # Manage addon version
  version = "v1.6.0.0-aliyun"
  # Manage addon config
  config = jsonencode(
    {}
  )
}
```
**Installing of addon**
When a cluster is created, some system addons and those specified at the time of cluster creation will be installed, so when an addon resource is applied:
* If the addon already exists in the cluster and its version is the same as the specified version, it will be skipped and will not be reinstalled.
* If the addon already exists in the cluster and its version is different from the specified version, the addon will be upgraded.
* If the addon does not exist in the cluster, it will be installed.

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cs_kubernetes_addon&spm=docs.r.cs_kubernetes_addon.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, ForceNew) The id of kubernetes cluster.
* `name` - (Required, ForceNew) The name of addon.
* `version` - (Optional, Computed) The current version of addon.
* `config` - (Optional, Available since v1.166.0) The customized configuration of addon. Your customized configuration will be merged to existed configuration stored in server. If you want to clean one configuration, you must set the configuration to empty value, removing from code cannot make effect. You can checkout the customized configuration of the addon through datasource `alicloud_cs_kubernetes_addon_metadata`, the returned format is the standard json schema. If return empty, it means that the addon does not support custom configuration yet. You can also checkout the current custom configuration through the data source `alicloud_cs_kubernetes_addons`.
* `cleanup_cloud_resources` - (Optional) Whether to clean up cloud resources when deleting. Currently only works for addon `ack-virtual-node` and you must specify it when uninstall addon `ack-virtual-node`. Valid values: `true`: clean up, `false`: do not clean up.

## Attributes Reference

The following attributes are exported:
* `id` - The id of addon, which consists of the cluster id and the addon name, with the structure <cluster_ud>:<addon_name>.
* `next_version` - The version which addon can be upgraded to.
* `can_upgrade` - Is the addon ready for upgrade.
* `required` - Is it a mandatory addon to be installed.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when installing addon in the kubernetes cluster. 
* `update` - (Defaults to 10 mins) Used when upgrading addon in the kubernetes cluster.
* `delete` - (Defaults to 10 mins) Used when deleting addon in kubernetes cluster. 

## Import

Cluster addon can be imported by cluster id and addon name. Then write the addon.tf file according to the result of `terraform plan`.

```shell
$ terraform import alicloud_cs_kubernetes_addon.my_addon <cluster_id>:<addon_name>
```
