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

* alicloud_access_key: The Alicloud Access Key ID
* alicloud_secret_key: The Alicloud Access Secret Key
* region: The ID of region in which launching resources
* k8s_name_prefix: The name prefix of kubernetes cluster
* k8s_number: The number of kubernetes cluster
* k8s_worker_number: The number of worker nodes in each kubernetes cluster
* k8s_pod_cidr: The kubernetes pod cidr block. It cannot be equals to vpc's or vswitch's and cannot be in them. If vpc's cidr block is `172.16.XX.XX/XX`,
it had better to `192.168.XX.XX/XX` or `10.XX.XX.XX/XX`
* k8s_service_cidr: The kubernetes service cidr block. Its setting rule is same as `k8s_pod_cidr`
* Other kubernetes cluster arguments

**Note:** In order to avoid some needless error, you had better to set `new_nat_gateway` to `true`.
Otherwise, you must you must ensure you specified vswitches can access internet before running the example.

Planning phase

    terraform plan

Apply phase

	terraform apply


Destroy

    terraform destroy


Conditional creation
--------------------
This example can support the following creating kubernetes cluster scenario by setting different arguments.

### 1. Create a new vpc, vswitches and nat gateway for the cluster.

You can specify the following user-defined arguments:

* vpc_name: A new vpc name
* vpc_cidr: A new vpc cidr block
* vswitch_name_prefix: The name prefix of several vswitches
* vswitch_cidrs: List of cidr blocks for several new vswitches

### 2. Using existing vpc and vswitches for the cluster.

You can specify the following user-defined arguments:

* vpc_id: A existing vpc ID
* vswitch_ids: List of IDs for several existing vswitches

### 3. Using existing vpc, vswitches and nat gateway for the cluster.

You can specify the following user-defined arguments:

* vpc_id: A existing vpc ID
* vswitch_ids: List of IDs for several existing vswitches
* new_nat_gateway: Set it to false. But you must ensure you specified vswitches can access internet.
In other words, you must set snat entry for each vswitch before running the example.


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
* [Terraform-Provider-Alicloud Github](https://github.com/terraform-providers/terraform-provider-alicloud)
* [Terraform-Provider-Alicloud Release](https://releases.hashicorp.com/terraform-provider-alicloud/)
* [Terraform-Provider-Alicloud Docs](https://www.terraform.io/docs/providers/alicloud/)


