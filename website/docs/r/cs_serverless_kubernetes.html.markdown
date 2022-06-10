---
subcategory: "Container Service for Kubernetes (ACK)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_serverless_kubernetes"
sidebar_current: "docs-alicloud-resource-cs-serverless-kubernetes"
description: |-
  Provides a Alicloud resource to manage container serverless kubernetes cluster.
---

# alicloud\_cs\_serverless_kubernetes

This resource will help you to manager a Serverless Kubernetes Cluster. The cluster is same as container service created by web console.

-> **NOTE:** Available in 1.58.0+

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

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "ask-example"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "10.1.1.0/24"
  zone_id      = data.alicloud_zones.default.zones[0].id
}

resource "alicloud_cs_serverless_kubernetes" "serverless" {
  name_prefix                    = var.name
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
  # Select an existing sls project
  # sls_project_name             = ""

  # tags 
  tags = {
    "k-aa" = "v-aa"
    "k-bb" = "v-aa"
  }

  # addons 
  addons {
    # SLB Ingress
    name = "alb-ingress-controller"
  }
  addons {
    name = "metrics-server"
  }
  addons {
    name = "knative"
  }

}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The kubernetes cluster's name. It is the only in one Alicloud account.
* `name_prefix` - (Optional) The kubernetes cluster name's prefix. It is conflict with `name`. If it is specified, terraform will using it to build the only cluster name. Default to "Terraform-Creation".
* `version` - (Optional, Available since 1.97.0) Desired Kubernetes version. If you do not specify a value, the latest available version at resource creation is used.
* `vpc_id` - (Required, ForceNew) The vpc where new kubernetes cluster will be located. Specify one vpc's id, if it is not specified, a new VPC  will be built.
* `vswitch_ids` - (Required) The vswitches where new kubernetes cluster will be located.
* `new_nat_gateway` - (Optional) Whether to create a new nat gateway while creating kubernetes cluster. SNAT must be configured when a new VPC is automatically created. Default is `true`.
* `endpoint_public_access_enabled` - (Optional, ForceNew) Whether to create internet  eip for API Server. Default to false.
* `service_discovery_types` - (ForceNew, Available in 1.123.1+) Service discovery type. If the value is empty, it means that service discovery is not enabled. Valid values are `CoreDNS` and `PrivateZone`. 
* `deletion_protection` - (Optional, ForceNew) Whether enable the deletion protection or not.
    - true: Enable deletion protection.
    - false: Disable deletion protection.
* `enable_rrsa` - (Optional, Available in 1.171.0+) Whether to enable cluster to support rrsa for version 1.22.3+. Default to `false`. Once the rrsa function is turned on, it is not allowed to turn off. If your cluster has enabled this function, please manually modify your tf file and add the rrsa configuration to the file, learn more [RAM Roles for Service Accounts](https://www.alibabacloud.com/help/zh/container-service-for-kubernetes/latest/use-rrsa-to-enforce-access-control).
* `force_update` - (Optional) Default false, when you want to change `vpc_id` and `vswitch_id`, you have to set this field to true, then the cluster will be recreated.
* `tags` - (Optional) Default nil, A map of tags assigned to the kubernetes cluster and work nodes.
* `kube_config` - (Optional) The path of kube config, like `~/.kube/config`.
* `client_cert` - (Optional) The path of client certificate, like `~/.kube/client-cert.pem`.
* `client_key` - (Optional) The path of client key, like `~/.kube/client-key.pem`.
* `cluster_ca_cert` - (Optional) The path of cluster ca certificate, like `~/.kube/cluster-ca-cert.pem`
* `security_group_id` - (Optional, Available in 1.91.0+) The ID of the security group to which the ECS instances in the cluster belong. If it is not specified, a new Security group will be built.
* `resource_group_id` - (Optional, ForceNew, Available in 1.101.0+) The ID of the resource group,by default these cloud resources are automatically assigned to the default resource group.
* `load_balancer_spec` - (Optional, Available in 1.117.0+) The cluster api server load balance instance specification, default `slb.s2.small`. For more information on how to select a LB instance specification, see [SLB instance overview](https://help.aliyun.com/document_detail/85931.html).
* `addons` - (Available in 1.91.0+)) You can specific network plugin,log component,ingress component and so on.Detailed below.
* `time_zone` - (Optional, ForceNew, Available in 1.123.1+) The time zone of the cluster.
* `zone_id` - (Optional, ForceNew, Available in 1.123.1+) When creating a cluster using automatic VPC creation, you need to specify the zone where the VPC is located. 
* `service_cidr` - (Optional, ForceNew, Available in 1.123.1+) CIDR block of the service network. The specified CIDR block cannot overlap with that of the VPC or those of the ACK clusters that are deployed in the VPC. The CIDR block cannot be modified after the cluster is created.
* `logging_type` - (ForceNew, Available in 1.123.1+) Enable log service, Valid value `SLS`. 
* `sls_project_name` - (ForceNew, Available in 1.123.1+) If you use an existing SLS project, you must specify `sls_project_name`.
* `retain_resources` - (Optional, Available in 1.141.0+) Resources that are automatically created during cluster creation, including NAT gateways, SNAT rules, SLB instances, and RAM Role, will be deleted. Resources that are manually created after you create the cluster, such as SLB instances for Services, will also be deleted. If you need to retain resources, please configure with `retain_resources`. There are several aspects to pay attention to when using `retain_resources` to retain resources. After configuring `retain_resources` into the terraform configuration manifest file, you first need to run `terraform apply`.Then execute `terraform destroy`.
* `cluster_spec` - (Optional, ForceNew, Available in 1.162.0+) The cluster specifications of serverless kubernetes cluster, which can be empty. Valid values:
    - ack.standard: Standard serverless clusters.
    - ack.pro.small: Professional serverless clusters.

#### addons 
It is a new field since 1.91.0. You can specific network plugin,log component,ingress component and so on. 
The following arguments are optional:
* `name` - Name of the ACK add-on. The name must match one of the names returned by [DescribeAddons](https://help.aliyun.com/document_detail/171524.html).
* `config` - The ACK add-on configurations.
* `disabled` -  Disables the automatic installation of a component. Default is `false`.

The following example is the definition of addons block, The type of this field is list:

```
# install nginx ingress, conflict with SLB ingress
addons {
  name = "nginx-ingress-controller"
  # use internet
  config = "{\"IngressSlbNetworkType\":\"internet",\"IngressSlbSpec\":\"slb.s2.small\"}"
  # if use intranet, detail below.
  # config = "{\"IngressSlbNetworkType\":\"intranet",\"IngressSlbSpec\":\"slb.s2.small\"}"
}
# install SLB ingress, conflict with nginx ingress
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
```

#### Removed params (Never Supported)
* `vswitch_id` - (Deprecated from version 1.91.0) (Required, ForceNew) The vswitch where new kubernetes cluster will be located. Specify one vswitch's id, if it is not specified, a new VPC and VSwicth will be built. It must be in the zone which `availability_zone` specified.
* `private_zone` - (Deprecated from version 1.123.1) (Optional, ForceNew) Has been deprecated from provider version 1.123.1. `PrivateZone` is used as the enumeration value of `service_discovery_types`.

### Timeouts

-> **NOTE:** Available in 1.58.0+

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 mins) Used when creating the kubernetes cluster (until it reaches the initial `running` status). 
* `delete` - (Defaults to 30 mins) Used when terminating the kubernetes cluster. 

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the container cluster.
* `name` - The name of the container cluster.
* `vpc_id` - The ID of VPC where the current cluster is located.
* `vswitch_id` - The ID of VSwicth where the current cluster is located.
* `security_group_id` - The ID of security group where the current cluster worker node is located.
* `deletion_protection` - Whether enable the deletion protection or not.

## Import

Serverless Kubernetes cluster can be imported using the id, e.g.

```
$ terraform import alicloud_cs_serverless_kubernetes.main ce4273f9156874b46bb
```
