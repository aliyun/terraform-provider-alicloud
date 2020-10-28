---
subcategory: "Container Service (CS)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_serverless_kubernetes"
sidebar_current: "docs-alicloud-resource-cs-serverless-kubernetes"
description: |-
  Provides a Alicloud resource to manage container serverless kubernetes cluster.
---

# alicloud\_cs\_serverless_kubernetes

This resource will help you to manager a Serverless Kubernetes Cluster. The cluster is same as container service created by web console.


-> **NOTE:** Serverless Kubernetes cluster only supports VPC network and it can access internet while creating kubernetes cluster.
A Nat Gateway and configuring a SNAT for it can ensure one VPC network access internet. If there is no nat gateway in the
VPC, you can set `new_nat_gateway` to "true" to create one automatically.

-> **NOTE:** Creating serverless kubernetes cluster need to install several packages and it will cost about 5 minutes. Please be patient.

-> **NOTE:** The provider supports to download kube config, client certificate, client key and cluster ca certificate
after creating cluster successfully, and you can put them into the specified location, like '~/.kube/config'.

-> **NOTE:** If you want to manage serverless Kubernetes, you can use [Kubernetes Provider](https://www.terraform.io/docs/providers/kubernetes/index.html).

-> **NOTE:** You need to activate several other products and confirm Authorization Policy used by Container Service before using this resource.
Please refer to the `Authorization management` and `Cluster management` sections in the [Document Center](https://www.alibabacloud.com/help/doc-detail/86488.htm).

-> **NOTE:** Available in 1.58.0+

## Example Usage

Basic Usage

```
variable "name" {
  default = "my-first-k8s"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alicloud_vpc" "default" {
  name       = var.name
  cidr_block = "10.1.0.0/21"
}

resource "alicloud_vswitch" "default" {
  name              = var.name
  vpc_id            = alicloud_vpc.default.id
  cidr_block        = "10.1.1.0/24"
  availability_zone = data.alicloud_zones.default.zones[0].id
}

resource "alicloud_cs_serverless_kubernetes" "serverless" {
  name_prefix                    = var.name
  vpc_id                         = alicloud_vpc.default.id
  vswitch_id                     = alicloud_vswitch.default.id
  new_nat_gateway                = true
  endpoint_public_access_enabled = true
  private_zone                   = false
  deletion_protection            = false
  tags = {
    "k-aa" = "v-aa"
    "k-bb" = "v-aa"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The kubernetes cluster's name. It is the only in one Alicloud account.
* `name_prefix` - (Optional) The kubernetes cluster name's prefix. It is conflict with `name`. If it is specified, terraform will using it to build the only cluster name. Default to "Terraform-Creation".
* `version` - (Optional, Available since 1.97.0) Desired Kubernetes version. If you do not specify a value, the latest available version at resource creation is used.
* `vpc_id` - (Required, ForceNew) The vpc where new kubernetes cluster will be located. Specify one vpc's id, if it is not specified, a new VPC  will be built.
* `vswitch_ids` - (Required, ForceNew) The vswitches where new kubernetes cluster will be located.
* `new_nat_gateway` - (Optional) Whether to create a new nat gateway while creating kubernetes cluster. Default to true.
* `endpoint_public_access_enabled` - (Optional, ForceNew) Whether to create internet  eip for API Server. Default to false.
* `private_zone` - (Optional, ForceNew) Enable Privatezone if you need to use the service discovery feature within the serverless cluster. Default to false.
* `deletion_protection` - (Optional, ForceNew) Whether enable the deletion protection or not.
    - true: Enable deletion protection.
    - false: Disable deletion protection.
* `force_update` - (Optional) Default false, when you want to change `vpc_id` and `vswitch_id`, you have to set this field to true, then the cluster will be recreated.
* `tags` - (Optional) Default nil, A map of tags assigned to the kubernetes cluster .
* `kube_config` - (Optional) The path of kube config, like `~/.kube/config`.
* `client_cert` - (Optional) The path of client certificate, like `~/.kube/client-cert.pem`.
* `client_key` - (Optional) The path of client key, like `~/.kube/client-key.pem`.
* `cluster_ca_cert` - (Optional) The path of cluster ca certificate, like `~/.kube/cluster-ca-cert.pem`
* `security_group_id` - (Optional, Available in 1.91.0+) The ID of the security group to which the ECS instances in the cluster belong. If it is not specified, a new Security group will be built.
* `resource_group_id` - (Optional, ForceNew, Available in 1.101.0+) The ID of the resource group,by default these cloud resources are automatically assigned to the default resource group.


#### Addons 
It is a new field since 1.91.0. You can specific network plugin,log component,ingress component and so on.     
 
```$xslt
  main.tf
   
  dynamic "addons" {
      for_each = var.cluster_addons
      content {
        name                    = lookup(addons.value, "name", var.cluster_addons)
        config                  = lookup(addons.value, "config", var.cluster_addons)
        disabled                = lookup(addons.value, "disabled", var.cluster_addons)
      }
  }
```
```$xslt
    varibales.tf 
    
    // Log
    variable "cluster_addons" {
        type = list(object({
            name      = string
            config    = string
        }))
    
        default = [
            {
                "name"     = "logtail-ds",
                "config"   = "{\"IngressDashboardEnabled\":\"true\",\"sls_project_name\":\"your-sls-project-name\"}",
            }
        ]
    } 
    
  
```
* `logtail-ds` - You can specific `IngressDashboardEnabled` and `sls_project_name` in config. If you switch on `IngressDashboardEnabled` and `sls_project_name`,then logtail-ds would use `sls_project_name` as default log store.

#### Removed params (Never Supported)
* `vswitch_id` - (Deprecated from version 1.91.0)(Required, ForceNew) The vswitch where new kubernetes cluster will be located. Specify one vswitch's id, if it is not specified, a new VPC and VSwicth will be built. It must be in the zone which `availability_zone` specified.


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
