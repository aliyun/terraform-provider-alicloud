---
subcategory: "EDAS"
layout: "alicloud"
page_title: "Alicloud: alicloud_edas_k8s_cluster"
sidebar_current: "docs-alicloud-resource-edas-k8s-cluster"
description: |-
  Provides an EDAS K8s cluster resource.
---

# alicloud_edas_k8s_cluster

Provides an EDAS K8s cluster resource. For information about EDAS K8s Cluster and how to use it, see[What is EDAS K8s Cluster](https://www.alibabacloud.com/help/en/doc-detail/85108.htm).

-> **NOTE:** Available since v1.93.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_edas_k8s_cluster&exampleId=c057f961-f5ce-43e7-748d-5ad1f3687623e6099d92&activeTab=example&spm=docs.r.edas_k8s_cluster.0.c057f961f5&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
data "alicloud_images" "default" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}
data "alicloud_instance_types" "default" {
  availability_zone    = data.alicloud_zones.default.zones.0.id
  cpu_core_count       = 4
  memory_size          = 8
  kubernetes_node_role = "Worker"
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

resource "alicloud_cs_kubernetes_node_pool" "default" {
  name                 = var.name
  cluster_id           = alicloud_cs_managed_kubernetes.default.id
  vswitch_ids          = [alicloud_vswitch.default.id]
  instance_types       = [data.alicloud_instance_types.default.instance_types.0.id]
  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  desired_size         = 2
}
resource "alicloud_edas_k8s_cluster" "default" {
  cs_cluster_id = alicloud_cs_kubernetes_node_pool.default.cluster_id
}
```

## Argument Reference

The following arguments are supported:
https://www.alibabacloud.com/help/zh/edas/developer-reference/api-edas-2017-08-01-importk8scluster

* `cs_cluster_id` - (Required, ForceNew) The ID of the alicloud container service kubernetes cluster that you want to import.
* `namespace_id` - (Optional, ForceNew) The ID of the namespace where you want to import. You can call the [ListUserDefineRegion](https://www.alibabacloud.com/help/en/doc-detail/149377.htm?spm=a2c63.p38356.879954.34.331054faK2yNvC#doc-api-Edas-ListUserDefineRegion) operation to query the namespace ID.


## Attributes Reference

The following attributes are exported:

* `cluster_name` - The name of the cluster that you want to create.
* `cluster_type` - The type of the cluster that you want to create. Valid values only: 5: K8s cluster. 
* `network_mode` - The network type of the cluster that you want to create. Valid values: 1: classic network. 2: VPC.
* `vpc_id` - The ID of the Virtual Private Cloud (VPC) for the cluster.
* `cluster_import_status` - The import status of cluster: 
    `1`: success.
    `2`: failed.
    `3`: importing. 
    `4`: deleted.

## Import

EDAS cluster can be imported using the id, e.g.

```shell
$ terraform import alicloud_edas_k8s_cluster.cluster cluster_id
```
