Alibaba Cloud terraform example for kubernetes cluster
======================================================

A terraform example to launching a kubernetes cluster in alibaba cloud.

These types of the module resource are supported:

- [VPC](https://www.terraform.io/docs/providers/alicloud/r/vpc.html)
- [Subnet](https://www.terraform.io/docs/providers/alicloud/r/vswitch.html)
- [ECS Instance](https://www.terraform.io/docs/providers/alicloud/r/instance.html)
- [Security Group](https://www.terraform.io/docs/providers/alicloud/r/security_group.html)
- [Nat Gateway](https://www.terraform.io/docs/providers/alicloud/r/nat_gateway.html)
- [Kubernetes](https://www.terraform.io/docs/providers/alicloud/r/cs_kubernetes.html)


Usage
-----
This example can specify the following arguments to create user-defined kuberntes cluster
* availability_zone: The availability zones of vswitches of kubernetes clusters.
* vpc_id: Existing vpc id used to create kubernetes vswitches and other resources.Leave it to empty then terraform would create a new one.
* vpc_cidr: Conflict with `vpc_id`, terraform using `vpc_cidr` to create a new vpc.
* vswitch_ids: The vswitch_ids of masters and workers. You can also use `master_vswitch_ids` and `worker_vswitch_ids` instead.
* vswitch_cidrs: Conflict with `vswitch_ids`, terraform using vswitch_cidrs to create vswitches.
* new_nat_gateway: Create a SNAT gateway for kubernetes cluster.Because of the openapi in Alibaba Cloud is not all on intranet.
* master_instance_types: The ecs instance types used to launch master nodes. 3 or 5 instance types are allowed. Be careful of the matching relation between instanceType and availability_zone. Not all instance types are available in different zones.
* worker_instance_types: The ecs instance types used to launch worker nodes.Configure more than 1 instance types is a better choice.
* node_cidr_mask: The node cidr block to specific how many pods can run on single node. 24-28 is allowed. 24 means 2^(32-24)-1=255 and the node can run at most 255 pods. default: 24
* enable_ssh: Enable login to the node through SSH. default: false 
* install_cloud_monitor: Install cloud monitor agent on ECS. default: true 
* cpu_policy: kubelet cpu policy. options: static|none. default: none.
* proxy_mode: Proxy mode is option of kube-proxy. options: iptables|ipvs. default: ipvs.
* password: The password of ECS instance.
* worker_number: The number of worker nodes in kubernetes cluster.
* pod_cidr: The kubernetes pod cidr block. It cannot be equals to vpc's or vswitch's and cannot be in them. [Flannel]
* service_cidr: The kubernetes service cidr block. It cannot be equals to vpc's or vswitch's or pod's and cannot be in them.
* cluster_addons: Addon components in kubernetes cluster
                                                                                                                             
Planning phase

    terraform plan

Apply phase

	terraform apply


Destroy

    terraform destroy

Terraform version
-----------------
Terraform version 0.11.0 or newer and Provider version 1.9.0 or newer are required for this example to work.

Authors
-------
Created and maintained by He Guimin(@xiaozhu36, heguimin36@163.com)

License
-------
Mozilla Public License 2.0. See LICENSE for full details.

Reference
---------
* [Terraform-Provider-Alicloud Github](https://github.com/aliyun/terraform-provider-alicloud)
* [Terraform-Provider-Alicloud Release](https://releases.hashicorp.com/terraform-provider-alicloud/)
* [Terraform-Provider-Alicloud Docs](https://www.terraform.io/docs/providers/alicloud/)


