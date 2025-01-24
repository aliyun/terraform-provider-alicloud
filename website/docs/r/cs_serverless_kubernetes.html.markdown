---
subcategory: "Container Service for Kubernetes (ACK)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_serverless_kubernetes"
sidebar_current: "docs-alicloud-resource-cs-serverless-kubernetes"
description: |-
  Provides a Alicloud resource to manage container serverless kubernetes cluster.
---

# alicloud_cs_serverless_kubernetes

This resource will help you to manager a Serverless Kubernetes Cluster, see [What is serverless kubernetes](https://www.alibabacloud.com/help/en/ack/ack-managed-and-ack-dedicated/developer-reference/create-a-dedicated-kubernetes-cluster-that-supports-sandboxed-containers). The cluster is same as container service created by web console.

-> **NOTE:** Available since v1.58.0.

-> **NOTE:** Serverless Kubernetes cluster only supports VPC network and it can access internet while creating kubernetes cluster.
A Nat Gateway and configuring a SNAT for it can ensure one VPC network access internet. If there is no nat gateway in the
VPC, you can set `new_nat_gateway` to "true" to create one automatically.

-> **NOTE:** Creating serverless kubernetes cluster need to install several packages and it will cost about 5 minutes. Please be patient.

-> **NOTE:** The provider supports to download kube config, client certificate, client key and cluster ca certificate
after creating cluster successfully, and you can put them into the specified location, like '~/.kube/config'.

-> **NOTE:** If you want to manage serverless Kubernetes, you can use [Kubernetes Provider](https://www.terraform.io/docs/providers/kubernetes/index.html).

-> **NOTE:** You need to activate several other products and confirm Authorization Policy used by Container Service before using this resource.
Please refer to the `Authorization management` and `Cluster management` sections in the [Document Center](https://www.alibabacloud.com/help/doc-detail/86488.htm).

-> **NOTE:** From version 1.162.0, support for creating professional serverless cluster.

-> **NOTE:** From version 1.229.1, support to migrate basic serverless cluster to professional serverless cluster.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cs_serverless_kubernetes&exampleId=d818bebd-0275-45d3-0929-c5c965ef152289c52fa8&activeTab=example&spm=docs.r.cs_serverless_kubernetes.0.d818bebd02&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "ask-example-pro"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.2.0.0/21"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "10.2.1.0/24"
  zone_id      = data.alicloud_zones.default.zones[0].id
}

resource "alicloud_cs_serverless_kubernetes" "serverless" {
  name_prefix                    = var.name
  cluster_spec                   = "ack.pro.small"
  vpc_id                         = alicloud_vpc.default.id
  vswitch_ids                    = [alicloud_vswitch.default.id]
  new_nat_gateway                = true
  endpoint_public_access_enabled = true
  deletion_protection            = false

  load_balancer_spec      = "slb.s2.small"
  time_zone               = "Asia/Shanghai"
  service_cidr            = "172.21.0.0/20"
  service_discovery_types = ["PrivateZone"]
  # Enable log service, A project named k8s-log-{ClusterID} will be automatically created
  logging_type = "SLS"

  # tags
  tags = {
    "k-aa" = "v-aa"
    "k-bb" = "v-aa"
  }

  # addons
  addons {
    # ALB Ingress
    name = "alb-ingress-controller"
  }
  addons {
    name = "metrics-server"
  }
  addons {
    name = "knative"
  }
  addons {
    name = "arms-prometheus"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The kubernetes cluster's name. It is the only in one Alicloud account.
* `name_prefix` - (Optional, ForceNew) The kubernetes cluster name's prefix. It is conflict with `name`. If it is specified, terraform will using it to build the only cluster name. Default to "Terraform-Creation".
* `version` - (Optional, Available since v1.97.0) Desired Kubernetes version. If you do not specify a value, the latest available version at resource creation is used.  Do not specify if cluster auto upgrade is enabled, see [cluster_auto_upgrade](#operation_policy-cluster_auto_upgrade) for more information.  
* `vpc_id` - (Optional, ForceNew) The vpc where new kubernetes cluster will be located. Specify one vpc's id, if it is not specified, a new VPC will be built.
* `vswitch_ids` - (Optional, ForceNew) The vswitches where new kubernetes cluster will be located.
* `new_nat_gateway` - (Optional) Whether to create a new nat gateway while creating kubernetes cluster. SNAT must be configured when a new VPC is automatically created. Default is `true`.
* `endpoint_public_access_enabled` - (Optional, ForceNew) Whether to create internet eip for API Server. Default to false. Only works for **Create** Operation. 
* `service_discovery_types` - (Optional, Available since v1.123.1) Service discovery type. Only works for **Create** Operation. If the value is empty, it means that service discovery is not enabled. Valid values are `CoreDNS` and `PrivateZone`.
* `deletion_protection` - (Optional, ForceNew) Whether enable the deletion protection or not.
    - true: Enable deletion protection.
    - false: Disable deletion protection.
* `enable_rrsa` - (Optional, Available since v1.171.0) Whether to enable cluster to support RRSA for version 1.22.3+. Default to `false`. Once the RRSA function is turned on, it is not allowed to turn off. If your cluster has enabled this function, please manually modify your tf file and add the rrsa configuration to the file, learn more [RAM Roles for Service Accounts](https://www.alibabacloud.com/help/zh/container-service-for-kubernetes/latest/use-rrsa-to-enforce-access-control).
* `tags` - (Optional) Default nil, A map of tags assigned to the kubernetes cluster and work nodes.
* `kube_config` - (Optional, Deprecated from v1.187.0) The path of kube config, like `~/.kube/config`.
* `client_cert` - (Optional) The path of client certificate, like `~/.kube/client-cert.pem`.
* `client_key` - (Optional) The path of client key, like `~/.kube/client-key.pem`.
* `cluster_ca_cert` - (Optional) The path of cluster ca certificate, like `~/.kube/cluster-ca-cert.pem`
* `security_group_id` - (Optional, ForceNew, Available since v1.91.0) The ID of the security group to which the ECS instances in the cluster belong. If it is not specified, a new Security group will be built.
* `resource_group_id` - (Optional, Available since v1.101.0) The ID of the resource group,by default these cloud resources are automatically assigned to the default resource group.
* `load_balancer_spec` - (Optional, Deprecated since v1.229.1) The cluster api server load balance instance specification, default `slb.s2.small`. For more information on how to select a LB instance specification, see [SLB instance overview](https://help.aliyun.com/document_detail/85931.html). Only works for **Create** Operation. 
* `addons` - (Optional, Available since v1.91.0) You can specific network plugin, log component, ingress component and so on. See [`addons`](#addons) below. Only works for **Create** Operation, use [resource cs_kubernetes_addon](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/resources/cs_kubernetes_addon) to manage addons if cluster is created.
* `time_zone` - (Optional, Available since v1.123.1) The time zone of the cluster.
* `zone_id` - (Optional, Available since v1.123.1) When creating a cluster using automatic VPC creation, you need to specify the zone where the VPC is located. Only works for **Create** Operation.  
* `service_cidr` - (Optional, ForceNew, Available since v1.123.1) CIDR block of the service network. The specified CIDR block cannot overlap with that of the VPC or those of the ACK clusters that are deployed in the VPC. The CIDR block cannot be modified after the cluster is created.
* `logging_type` - (Optional, ForceNew, Deprecated since v1.229.1) Enable log service, Valid value `SLS`. Only works for **Create** Operation. 
* `sls_project_name` - (Optional, ForceNew, Deprecated since v1.229.1) If you use an existing SLS project, you must specify `sls_project_name`. Only works for **Create** Operation. 
* `retain_resources` - (Optional, Available since v1.141.0) Resources that are automatically created during cluster creation, including NAT gateways, SNAT rules, SLB instances, and RAM Role, will be deleted. Resources that are manually created after you create the cluster, such as SLB instances for Services, will also be deleted. If you need to retain resources, please configure with `retain_resources`. There are several aspects to pay attention to when using `retain_resources` to retain resources. After configuring `retain_resources` into the terraform configuration manifest file, you first need to run `terraform apply`.Then execute `terraform destroy`.
* `delete_options` - (Optional, Available since v1.229.1) Delete options, only work for deleting resource. Make sure you have run `terraform apply` to make the configuration applied. See [`delete_options`](#delete_options) below.
* `cluster_spec` - (Optional, ForceNew, Available since v1.162.0) The cluster specifications of serverless kubernetes cluster, which can be empty. Valid values:
    - ack.standard: Standard serverless clusters.
    - ack.pro.small: Professional serverless clusters.
* `custom_san` - (Optional, Available since v1.229.1) Customize the certificate SAN, multiple IP or domain names are separated by English commas (,).
-> **NOTE:** Make sure you have specified all certificate SANs before updating. Updating this field will lead APIServer to restart.
* `maintenance_window` - (Optional, Available since v1.232.0) The cluster maintenance window，effective only in the professional managed cluster. Managed node pool will use it. See [`maintenance_window`](#maintenance_window) below.
* `operation_policy` - (Optional, Available since v1.232.0) The cluster automatic operation policy. See [`operation_policy`](#operation_policy) below.

*Removed params*

* `vswitch_id` - (Removed since v1.229.1) The vswitch where new kubernetes cluster will be located. Specify one vswitch's id, if it is not specified, a new VPC and VSwicth will be built. It must be in the zone which `availability_zone` specified.
* `private_zone` - (Deprecated since v1.123.1) Has been deprecated from provider version 1.123.1. `PrivateZone` is used as the enumeration value of `service_discovery_types`.
* `create_v2_cluster` - (Removed since v1.229.1) whether to create a v2 version cluster.
* `force_update` - (Removed since v1.229.1) Default false, when you want to change `vpc_id` and `vswitch_id`, you have to set this field to true, then the cluster will be recreated.

### `maintenance_window`

The following arguments are supported in the `maintenance_window` configuration block:

* `enable` - (Optional) Whether to open the maintenance window. The following parameters take effect only `enable = true`.
* `maintenance_time` - (Optional) Initial maintenance time, RFC3339 format. For example: "2024-10-15T12:31:00.000+08:00".
* `duration` - (Optional) The maintenance time, values range from 1 to 24,unit is hour. For example: "3h".
* `weekly_period` - (Optional) Maintenance cycle, you can set the values from Monday to Sunday, separated by commas when the values are multiple. The default is Thursday.

for example:
```
  maintenance_window {
    enable            = true
    maintenance_time  = "2024-10-15T12:31:00.000+08:00"
    duration          = "3h"
    weekly_period     = "Monday,Friday"
  }
```

### `operation_policy`
* `cluster_auto_upgrade` - (Optional) Automatic cluster upgrade policy. See [`cluster_auto_upgrade`](#operation_policy-cluster_auto_upgrade) below.

### `operation_policy-cluster_auto_upgrade`
Automatic cluster upgrade policy. If `enabled` is set to `true`, ACK will automatically upgrade cluster depends on the `channel` value. The `version` field may show diffs if set in config, please remove the field or ignore it.  

* `enabled` - (Optional) Whether to enable automatic cluster upgrade.
* `channel` - (Optional) The automatic cluster upgrade channel. Valid values: `patch`, `stable`, `rapid`.

for example:
```
  operation_policy {
    cluster_auto_upgrade {
      enabled = true
      channel = "stable"
    }
  }
```


### `addons`

The addons supports the following:

* `name` - (Optional) Name of the ACK add-on. The name must match one of the names returned by [DescribeAddons](https://help.aliyun.com/document_detail/171524.html).
* `config` - (Optional) The ACK add-on configurations. For more config information, see [cs_kubernetes_addon_metadata](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/data-sources/cs_kubernetes_addon_metadata).
* `version` - (Optional, Available since v1.229.1) It specifies the version of the component.
* `disabled` - (Optional) Disables the automatic installation of a component. Default is `false`.

The following example is the definition of addons block, The type of this field is list:

```
# install nginx ingress, conflict with ALB ingress
addons {
  name = "nginx-ingress-controller"
  # use internet
  config = "{\"IngressSlbNetworkType\":\"internet",\"IngressSlbSpec\":\"slb.s2.small\"}"
  # if use intranet, detail below.
  # config = "{\"IngressSlbNetworkType\":\"intranet",\"IngressSlbSpec\":\"slb.s2.small\"}"
}
# install ALB ingress, conflict with nginx ingress
addons {
  name = "alb-ingress-controller"
}
# install metric server
addons {
  name = "metrics-server"
}
# install knative
addons {
  name = "knative"
}
# install prometheus
addons {
  name = "arms-prometheus"
}
```
### `delete_options`

The following arguments are supported in the `delete_options` configuration block:
* `delete_mode` - (Optional) The deletion mode of the cluster. Different resources may have different default behavior, see `resource_type` for details. Valid values:
  - `delete`: delete resources created by the cluster.
  - `retain`: retain resources created by the cluster.
* `resource_type` - (Optional) The type of resources that are created by cluster. Valid values:
  - `SLB`: SLB resources created by the Nginx Ingress Service, default behavior is to delete, option to retain is available. 
  - `ALB`: ALB resources created by the ALB Ingress Controller, default behavior is to retain, option to delete is available. 
  - `SLS_Data`: SLS Project used by the cluster logging feature, default behavior is to retain, option to delete is available. 
  - `SLS_ControlPlane`: SLS Project used for the managed cluster control plane logs, default behavior is to retain, option to delete is available.
  - `PrivateZone`: PrivateZone resources created by the cluster, default behavior is to retain, option to delete is available.
```
  ...
  // Specify delete_options as below when deleting cluster
  // delete SLB resources created by the Nginx Ingress Service
  delete_options {
    delete_mode = "delete"
    resource_type = "SLB"
  }
  // delete ALB resources created by the ALB Ingress Controller
  delete_options {
    delete_mode = "delete"
    resource_type = "ALB"
  }
  // delete SLS Project used by the cluster logging feature
  delete_options {
    delete_mode = "delete"
    resource_type = "SLS_Data"
  }
  // delete SLS Project used for the managed cluster control plane logs
  delete_options {
    delete_mode = "delete"
    resource_type = "SLS_ControlPlane"
  }
  // delete PrivateZone resources created by the cluster
  delete_options {
    delete_mode = "delete"
    resource_type = "PrivateZone"
  }
```

## Attributes Reference

The following attributes are exported:

* `id` - The Cluster ID of the serverless cluster.
* `rrsa_metadata` - Nested attribute containing RRSA related data for your cluster.
  * `enabled` - Whether the RRSA feature has been enabled.
  * `rrsa_oidc_issuer_url` - The issuer URL of RRSA OIDC Token.
  * `ram_oidc_provider_name` - The name of OIDC Provider that was registered in RAM.
  * `ram_oidc_provider_arn` -  The arn of OIDC provider that was registered in RAM.

## Timeouts

-> **NOTE:** Available since v1.58.0.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 mins) Used when creating the kubernetes cluster (until it reaches the initial `running` status). 
* `delete` - (Defaults to 30 mins) Used when terminating the kubernetes cluster. 


## Import

Serverless Kubernetes cluster can be imported using the id, e.g. Then complete the main.tf accords to the result of `terraform plan`.

```shell
$ terraform import alicloud_cs_serverless_kubernetes.main ce4273f9156874b46bb
```
