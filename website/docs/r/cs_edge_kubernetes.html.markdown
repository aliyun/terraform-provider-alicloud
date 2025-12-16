---
subcategory: "Container Service for Kubernetes (ACK)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_edge_kubernetes"
sidebar_current: "docs-alicloud-resource-cs-edge-kubernetes"
description: |-
  Provides a Alicloud resource to manage container edge kubernetes cluster.
---

# alicloud_cs_edge_kubernetes

This resource will help you to manage a Edge Kubernetes Cluster in Alibaba Cloud Kubernetes Service, see [What is edge kubernetes](https://www.alibabacloud.com/help/en/ack/ack-managed-and-ack-dedicated/developer-reference/create-an-ack-edge-cluster).

-> **NOTE:** Kubernetes cluster only supports VPC network and it can access internet while creating kubernetes cluster.
A Nat Gateway and configuring a SNAT for it can ensure one VPC network access internet. If there is no nat gateway in the
VPC, you can set `new_nat_gateway` to "true" to create one automatically.

-> **NOTE:** Creating kubernetes cluster need to install several packages and it will cost about 15 minutes. Please be patient.

-> **NOTE:** The provider supports to download kube config, client certificate, client key and cluster ca certificate
after creating cluster successfully, and you can put them into the specified location, like '~/.kube/config'.

-> **NOTE:** The provider supports disabling internet load balancer for API Server by setting `false` to `slb_internet_enabled`.

-> **NOTE:** If you want to manage Kubernetes, you can use [Kubernetes Provider](https://www.terraform.io/docs/providers/kubernetes/index.html).

-> **NOTE:** Available since v1.103.0.

-> **NOTE:** From version 1.185.0+, support new fields `cluster_spec`, `runtime` and `load_balancer_spec`.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cs_edge_kubernetes&exampleId=6ee7e192-4503-b08d-2813-7f0849ea44ff9c8945d8&activeTab=example&spm=docs.r.cs_edge_kubernetes.0.6ee7e19245&intl_lang=EN_US" target="_blank">
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
data "alicloud_instance_types" "default" {
  availability_zone    = data.alicloud_zones.default.zones.0.id
  cpu_core_count       = 4
  memory_size          = 8
  kubernetes_node_role = "Master"
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

resource "alicloud_cs_edge_kubernetes" "default" {
  name_prefix                    = var.name
  worker_vswitch_ids             = [alicloud_vswitch.default.id]
  worker_instance_types          = [data.alicloud_instance_types.default.instance_types.0.id]
  version                        = "1.26.3-aliyun.1"
  worker_number                  = 1
  password                       = "Test12345"
  pod_cidr                       = "10.99.0.0/16"
  service_cidr                   = "172.16.0.0/16"
  worker_instance_charge_type    = "PostPaid"
  new_nat_gateway                = true
  node_cidr_mask                 = 24
  install_cloud_monitor          = true
  slb_internet_enabled           = true
  is_enterprise_security_group   = true
  skip_set_certificate_authority = true

  worker_data_disks {
    category  = "cloud_ssd"
    size      = "200"
    encrypted = "false"
  }
}
```

You could create a professional kubernetes edge cluster now.

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cs_edge_kubernetes&exampleId=edb46713-c344-72e1-9513-e5cd38937431d9719edd&activeTab=example&spm=docs.r.cs_edge_kubernetes.1.edb46713c3&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf_example"
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}
data "alicloud_instance_types" "default" {
  availability_zone    = data.alicloud_zones.default.zones.0.id
  cpu_core_count       = 4
  memory_size          = 8
  kubernetes_node_role = "Master"
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

resource "alicloud_cs_edge_kubernetes" "default" {
  name_prefix                    = var.name
  worker_vswitch_ids             = [alicloud_vswitch.default.id]
  worker_instance_types          = [data.alicloud_instance_types.default.instance_types.0.id]
  cluster_spec                   = "ack.pro.small"
  worker_number                  = 1
  password                       = "Test12345"
  pod_cidr                       = "10.99.0.0/16"
  service_cidr                   = "172.16.0.0/16"
  worker_instance_charge_type    = "PostPaid"
  new_nat_gateway                = true
  node_cidr_mask                 = 24
  load_balancer_spec             = "slb.s2.small"
  install_cloud_monitor          = true
  slb_internet_enabled           = true
  is_enterprise_security_group   = true
  skip_set_certificate_authority = true

  worker_data_disks {
    category  = "cloud_ssd"
    size      = "200"
    encrypted = "false"
  }
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cs_edge_kubernetes&spm=docs.r.cs_edge_kubernetes.example&intl_lang=EN_US)

## Argument Reference

*Global params*

* `name` - (Optional) The kubernetes cluster's name. It is unique in one Alicloud account.
* `name_prefix` - (Optional) The kubernetes cluster name's prefix. It is conflict with `name`. If it is specified, terraform will using it to build the only cluster name. Default to "Terraform-Creation".
* `version` - (Optional) Desired Kubernetes version. If you do not specify a value, the latest available version at resource creation is used and no upgrades will occur except you set a higher version number. The value must be configured and increased to upgrade the version when desired. Downgrades are not supported by ACK.
* `security_group_id` - (Optional) The ID of the security group to which the ECS instances in the cluster belong. If it is not specified, a new Security group will be built.
* `is_enterprise_security_group` - (Optional) Enable to create advanced security group. default: false. See [Advanced security group](https://www.alibabacloud.com/help/doc-detail/120621.htm).
* `addons` - (Optional) The addon you want to install in cluster. See [`addons`](#addons) below.
* `rds_instances` - (Optional, Available since v1.103.2) RDS instance list, You can choose which RDS instances whitelist to add instances to.
* `resource_group_id` - (Optional, Available since v1.103.2) The ID of the resource group,by default these cloud resources are automatically assigned to the default resource group.
* `deletion_protection` - (Optional, Available since v1.103.2)  Whether to enable cluster deletion protection.
* `tags` - (Optional, Available since v1.120.0) Default nil, A map of tags assigned to the kubernetes cluster and work node.
* `retain_resources` - (Optional, Available since v1.141.0) Resources that are automatically created during cluster creation, including NAT gateways, SNAT rules, SLB instances, and RAM Role, will be deleted. Resources that are manually created after you create the cluster, such as SLB instances for Services, will also be deleted. If you need to retain resources, please configure with `retain_resources`. There are several aspects to pay attention to when using `retain_resources` to retain resources. After configuring `retain_resources` into the terraform configuration manifest file, you first need to run `terraform apply`.Then execute `terraform destroy`.
* `cluster_spec` - (Optional, Available since v1.185.0) The cluster specifications of kubernetes cluster,which can be empty. Valid values:
  * ack.standard : Standard edge clusters.
  * ack.pro.small : Professional edge clusters.
* `runtime` - (Optional, Available since v1.185.0) The runtime of containers. If you select another container runtime, see [Comparison of Docker, containerd, and Sandboxed-Container](https://www.alibabacloud.com/help/doc-detail/160313.htm). See [`runtime`](#runtime) below.
* `availability_zone` - (Optional) The ID of availability zone.
* `skip_set_certificate_authority` - (Optional) Configure whether to save certificate authority data for your cluster to attribute `certificate_authority`. For cluster security, recommended configuration as `true`. Will be removed with attribute certificate_authority removed.

*Network params*

* `pod_cidr` - (Optional) [Flannel Specific] The CIDR block for the pod network when using Flannel.
* `new_nat_gateway` - (Optional) Whether to create a new nat gateway while creating kubernetes cluster. Default to true. Then openapi in Alibaba Cloud are not all on intranet, So turn this option on is a good choice.
* `service_cidr` - (Optional) The CIDR block for the service network. It cannot be duplicated with the VPC CIDR and CIDR used by Kubernetes cluster in VPC, cannot be modified after creation.
* `node_cidr_mask` - (Optional) The node cidr block to specific how many pods can run on single node. 24-28 is allowed. 24 means 2^(32-24)-1=255 and the node can run at most 255 pods. default: 24
* `slb_internet_enabled` - (Optional) Whether to create internet load balancer for API Server. Default to true.
* `load_balancer_spec` - (Optional, Available since v1.185.0) The cluster api server load balance instance specification. For more information on how to select a LB instance specification, see [SLB instance overview](https://help.aliyun.com/document_detail/85931.html).
->NOTE: If you want to use `Flannel` as CNI network plugin, You need to specific the `pod_cidr` field and addons with `flannel`.

*Worker params*

* `password` - (Optional, Sensitive) The password of ssh login cluster node. You have to specify one of `password`, `key_name` `kms_encrypted_password` fields.
* `key_name` - (Optional) The keypair of ssh login cluster node, you have to create it first. You have to specify one of `password` `key_name` `kms_encrypted_password` fields.
* `worker_number` - (Required) The cloud worker node number of the edge kubernetes cluster. Default to 1. It is limited up to 50 and if you want to enlarge it, please apply white list or contact with us.
* `worker_vswitch_ids` - (Required) The vswitches used by workers.
* `worker_instance_charge_type` - (Optional) Worker payment type, its valid value is `PostPaid`. Defaults to `PostPaid`. More charge details in [ACK@edge charge](https://help.aliyun.com/document_detail/178718.html).
* `worker_instance_types` - (Required) The instance types of worker node, you can set multiple types to avoid NoStock of a certain type.
* `worker_disk_category` - (Optional) The system disk category of worker node. Its valid value are `cloud_efficiency`, `cloud_ssd` and `cloud_essd` and . Default to `cloud_efficiency`.
* `worker_disk_size` - (Optional) The system disk size of worker node. Its valid value range [20~32768] in GB. Default to 40.
* `worker_data_disks` - (Optional) The data disk configurations of worker nodes, such as the disk type and disk size. See [`worker_data_disks`](#worker_data_disks) below.
* `install_cloud_monitor` - (Optional) Install cloud monitor agent on ECS. default: `true`.
* `proxy_mode` - (Optional) Proxy mode is option of kube-proxy. options: iptables|ipvs. default: ipvs.
* `user_data` - (Optional) Windows instances support batch and PowerShell scripts. If your script file is larger than 1 KB, we recommend that you upload the script to Object Storage Service (OSS) and pull it through the internal endpoint of your OSS bucket.
* `worker_disk_performance_level` - (Optional, Available since v1.120.0) Worker node system disk performance level, when `worker_disk_category` values `cloud_essd`, the optional values are `PL0`, `PL1`, `PL2` or `PL3`, but the specific performance level is related to the disk capacity. For more information, see [Enhanced SSDs](https://www.alibabacloud.com/help/doc-detail/122389.htm). Default is `PL1`.
* `worker_disk_snapshot_policy_id` - (Optional, Available since v1.120.0) Worker node system disk auto snapshot policy.

*Computed params*

You can set some file paths to save kube_config information, but this way is cumbersome. Since version 1.105.0, we've written it to tf state file. About its use，see export attribute certificate_authority. From version 1.187.0+, new DataSource `alicloud_cs_cluster_credential` is recommended to manage cluster's kube_config.

* `kube_config` - (Optional, Deprecated from v1.187.0) The path of kube config, like ~/.kube/config. Please use the attribute [output_file](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/data-sources/cs_cluster_credential#output_file) of new DataSource `alicloud_cs_cluster_credential` to replace it.
* `client_cert` - (Optional, Deprecated from v1.248.0) From version 1.248.0, new DataSource `alicloud_cs_cluster_credential` is recommended to manage cluster's kubeconfig, you can also save the [certificate_authority.client_cert](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/data-sources/cs_cluster_credential#client_cert) attribute content of new DataSource `alicloud_cs_cluster_credential` to an appropriate path(like ~/.kube/client-cert.pem) for replace it.
* `client_key` - (Optional, Deprecated from v1.248.0) From version 1.248.0, new DataSource `alicloud_cs_cluster_credential` is recommended to manage cluster's kubeconfig, you can also save the [certificate_authority.client_key](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/data-sources/cs_cluster_credential#client_key) attribute content of new DataSource `alicloud_cs_cluster_credential` to an appropriate path(like ~/.kube/client-key.pem) for replace it.
* `cluster_ca_cert` - (Optional, Deprecated from v1.248.0) From version 1.248.0, new DataSource `alicloud_cs_cluster_credential` is recommended to manage cluster's kubeconfig, you can also save the [certificate_authority.cluster_cert](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/data-sources/cs_cluster_credential#cluster_cert) attribute content of new DataSource `alicloud_cs_cluster_credential` to an appropriate path(like ~/.kube/cluster-ca-cert.pem) for replace it.

*Removed params*

* `log_config` - (Optional, Deprecated) A list of one element containing information about the associated log store. See [`log_config`](#log_config) below.
* `force_update` - (Removed) Whether to force the update of kubernetes cluster arguments. Default to false.

### `addons`

The addons supports the following:

* `name` - (Optional) Name of the ACK add-on. The name must match one of the names returned by [DescribeAddons](https://help.aliyun.com/document_detail/171524.html).
* `config` - (Optional) The ACK add-on configurations. For more config information, see [cs_kubernetes_addon_metadata](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/data-sources/cs_kubernetes_addon_metadata).
* `version` - (Optional) It specifies the version of the component.
* `disabled` - (Optional) Disables the automatic installation of a component. Default is `false`.

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
```

### `worker_data_disks`

The worker_data_disks supports the following:

* `category` - (Optional) The type of the data disks. Valid values: `cloud`, `cloud_efficiency`, `cloud_ssd` and `cloud_essd`. Default to `cloud_efficiency`.
* `size` - (Optional) The size of a data disk, at least 40. Unit: GiB.
* `encrypted` - (Optional) Specifies whether to encrypt data disks. Valid values: true and false. Default is `false`.
* `performance_level` - (Optional, Available since v1.120.0) Worker node data disk performance level, when `category` values `cloud_essd`, the optional values are `PL0`, `PL1`, `PL2` or `PL3`, but the specific performance level is related to the disk capacity. For more information, see [Enhanced SSDs](https://www.alibabacloud.com/help/doc-detail/122389.htm). Default is `PL1`.
* `auto_snapshot_policy_id` - (Optional, Available since v1.120.0) Worker node data disk auto snapshot policy.
* `snapshot_id` - (Optional) The id of snapshot.
* `kms_key_id` - (Optional) The id of the kms key.
* `name` - (Optional) The name of the data disks.
* `device` - (Optional) The device of the data disks.

### `log_config`

The log_config supports the following:

* `type` - (Required) Type of collecting logs, only `SLS` are supported currently.
* `project` - (Optional) Log Service project name, cluster logs will output to this project.

### `runtime`

* `name` - (Optional) The name of the runtime. Supported runtimes can be queried by data source [alicloud_cs_kubernetes_version](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/data-sources/cs_kubernetes_version).
* `version` - (Optional) The version of the runtime.

The following example is the definition of runtime block:

```
  runtime = {
    name = "containerd"
    version = "1.6.28"
  }
```

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the container cluster.
* `vpc_id` - The ID of VPC where the current cluster is located.
* `slb_intranet` - The ID of private load balancer where the current cluster master node is located.
* `slb_internet` - The public ip of load balancer.
* `nat_gateway_id` - The ID of nat gateway used to launch kubernetes cluster.
* `worker_nodes` - List of cluster worker nodes.
  * `id` - ID of the node.
  * `name` - Node name.
  * `private_ip` - The private IP address of node.
* `worker_ram_role_name` - The RamRole Name attached to worker node.
* `connections` - (Map) Map of kubernetes cluster connection information.
  * `api_server_internet` - API Server Internet endpoint.
  * `api_server_intranet` - API Server Intranet endpoint.
  * `master_public_ip` - Master node SSH IP address.
  * `service_domain` - Service Access Domain.
* `certificate_authority` - (Map, Deprecated from v1.248.0) Nested attribute containing certificate authority data for your cluster. Please use the attribute [certificate_authority](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/data-sources/cs_cluster_credential#certificate_authority) of new DataSource `alicloud_cs_cluster_credential` to replace it.
  * `cluster_cert` - The base64 encoded cluster certificate data required to communicate with your cluster. Add this to the certificate-authority-data section of the kubeconfig file for your cluster.
  * `client_cert` - The base64 encoded client certificate data required to communicate with your cluster. Add this to the client-certificate-data section of the kubeconfig file for your cluster.
  * `client_key` - The base64 encoded client key data required to communicate with your cluster. Add this to the client-key-data section of the kubeconfig file for your cluster.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 60 mins) Used when creating the kubernetes cluster (until it reaches the initial `running` status). 
* `update` - (Defaults to 60 mins) Used when activating the kubernetes cluster when necessary during update.
* `delete` - (Defaults to 60 mins) Used when terminating the kubernetes cluster. 

## Import

Kubernetes edge cluster can be imported using the id, e.g. Then complete the main.tf accords to the result of `terraform plan`.

```shell
$ terraform import alicloud_cs_edge_kubernetes.main cluster-id
```
