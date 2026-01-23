---
subcategory: "Container Service for Kubernetes (ACK)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_clusters"
sidebar_current: "docs-alicloud-datasource-cs-clusters"
description: |-
  Provides a list of Ack Cluster owned by an Alibaba Cloud account.
---

# alicloud_cs_clusters

This data source provides Ack Cluster available to the user.[What is Cluster](https://next.api.alibabacloud.com/document/CS/2015-12-15/CreateCluster)

-> **NOTE:** Available since v1.269.0.

## Example Usage

```terraform
provider "alicloud" {
  region = "cn-hangzhou"
}

variable "name" {
  default = "terraform-example"
}

variable "zone_1" {
  default = "cn-hangzhou-k"
}

variable "zone_2" {
  default = "cn-hangzhou-g"
}

variable "vsw1_cidr" {
  default = "10.1.0.0/24"
}

variable "vsw2_cidr" {
  default = "10.1.1.0/24"
}

variable "container_cidr" {
  default = "172.17.3.0/24"
}

variable "service_cidr" {
  default = "172.17.2.0/24"
}

resource "alicloud_vpc" "default" {
  cidr_block = "10.0.0.0/8"
}

resource "alicloud_security_group" "default" {
  vpc_id              = alicloud_vpc.default.id
  security_group_name = "tf-example-security-group"
  security_group_type = "normal"
}


resource "alicloud_vswitch" "default0" {
  vpc_id     = alicloud_vpc.default.id
  cidr_block = var.vsw1_cidr
  zone_id    = var.zone_1
}

resource "alicloud_vswitch" "default1" {
  vpc_id     = alicloud_vpc.default.id
  zone_id    = var.zone_2
  cidr_block = var.vsw2_cidr
}

resource "alicloud_cs_managed_kubernetes" "default" {
  pod_cidr          = var.container_cidr
  vswitch_ids       = ["${alicloud_vswitch.default0.id}", "${alicloud_vswitch.default1.id}"]
  service_cidr      = var.service_cidr
  security_group_id = alicloud_security_group.default.id
  cluster_spec      = "ack.pro.small"
}

data "alicloud_cs_clusters" "default" {
  ids        = ["${alicloud_cs_managed_kubernetes.default.id}"]
  name_regex = alicloud_cs_managed_kubernetes.default.name
}

output "alicloud_cs_managed_kubernetes_example_id" {
  value = data.alicloud_cs_clusters.default.clusters.0.id
}
```

## Argument Reference

The following arguments are supported:
* `cluster_id` - (ForceNew, Optional) The cluster ID.
* `cluster_name` - (ForceNew, Optional) Custom cluster name.
* `cluster_spec` - (ForceNew, Optional) The specification of the clusters to query. Valid values:
  - `ack.pro.small`: ACK Pro clusters.
  - `ack.standard`: ACK Basic clusters.
* `cluster_type` - (ForceNew, Optional) The type of the clusters to query. Valid values:
  - `Kubernetes`: ACK dedicated clusters.
  - `ManagedKubernetes`: ACK managed clusters. ACK managed clusters include ACK Basic clusters, ACK Pro clusters, ACK Serverless Basic clusters, ACK Serverless Pro clusters, ACK Edge Basic clusters, ACK Edge Pro clusters, and ACK Lingjun Pro clusters.
  - `ExternalKubernetes`: registered clusters.
* `profile` - (ForceNew, Optional) The subtype of the clusters to query. Valid values:
  - `Default`: ACK managed clusters. ACK managed clusters include ACK Basic clusters and ACK Pro clusters.
  - `Edge`: ACK Edge clusters. ACK Edge clusters include ACK Edge Basic clusters and ACK Edge Pro clusters.
  - `Serverless`: ACK Serverless clusters. ACK Serverless clusters include ACK Serverless Basic clusters and ACK Serverless Pro clusters.
  - `Lingjun`: ACK Lingjun Pro clusters.
* `ids` - (Optional, ForceNew, Computed) A list of Cluster IDs. 
* `name_regex` - (Optional, ForceNew) A regex string to filter results by cluster name.
* `enable_details` - (Optional, ForceNew) Default to `false`. Set it to `true` can output more details about resource attributes.
* `output_file` - (Optional, ForceNew) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Cluster IDs.
* `names` - A list of name of Clusters.
* `clusters` - A list of Cluster Entries. Each element contains the following attributes:
    * `auto_mode` - **NOTE:** This field is only available when `enable_details` is `true`. Intelligent managed mode configuration.
        * `enabled` - Specifies whether to enable the intelligent managed mode.
    * `cluster_domain` - The local domain name of the cluster.
    * `cluster_id` - The cluster ID.
    * `cluster_name` - Custom cluster name.
    * `cluster_spec` - After you set `cluster_type` to `ManagedKubernetes` and configure `profile`, you can further specify the cluster specification.
    * `cluster_type` - The cluster type.
    * `current_version` - The current version of the cluster.
    * `deletion_protection` - Cluster deletion protection prevents accidental deletion of the cluster through the console or API.
    * `ip_stack` - The IP protocol stack of the cluster.
    * `maintenance_window` - **NOTE:** This field is only available when `enable_details` is `true`. Cluster maintenance window.
        * `duration` - The duration of the maintenance window.
        * `enable` - Indicates whether to enable the maintenance window.
        * `maintenance_time` - Maintenance start time.
        * `recurrence` - The recurrence rule for the maintenance window, defined using RFC5545 Recurrence Rule syntax.
        * `weekly_period` - The maintenance cycle.
    * `node_cidr_mask` - **NOTE:** This field is only available when `enable_details` is `true`. The number of IP addresses per node, determined by specifying the CIDR block of the network.
    * `operation_policy` - **NOTE:** This field is only available when `enable_details` is `true`. The automatic operations and maintenance policy for the cluster.
        * `cluster_auto_upgrade` - Cluster automatic upgrade.
            * `channel` - Cluster automatic upgrade frequency.
            * `enabled` - Whether to enable cluster automatic upgrade.
    * `pod_cidr` - The CIDR block for the pod network.
    * `profile` - ACK managed cluster profile.
    * `proxy_mode` - kube-proxy proxy mode.
    * `region_id` - The region ID where the cluster is deployed.
    * `resource_group_id` - The resource group ID of the cluster.
    * `security_group_id` - The security group ID for the control plane.
    * `service_cidr` - The Service CIDR block.
    * `state` - Cluster operational status.
    * `tags` - Cluster resource tags.
    * `timezone` - Cluster time zone.
    * `vpc_id` - The Virtual Private Cloud (VPC) used by the cluster.
    * `vswitch_ids` - Virtual switches for the cluster control plane.
    * `id` - The ID of the resource supplied above.
