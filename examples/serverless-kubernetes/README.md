Alibaba Cloud terraform example for kubernetes cluster
======================================================

A terraform example to launching a serveless kubernetes cluster in alibaba cloud.

These types of the module resource are supported:

- [VPC](https://www.terraform.io/docs/providers/alicloud/r/vpc.html)
- [Subnet](https://www.terraform.io/docs/providers/alicloud/r/vswitch.html)
- [Serverless](https://www.terraform.io/docs/providers/alicloud/r/cs_serverless_kubernetes.html)


Usage
-----
This example can specify the following arguments to create user-defined kuberntes cluster

* alicloud_access_key: The Alicloud Access Key ID
* alicloud_secret_key: The Alicloud Access Secret Key
* region: The ID of region in which launching resources
* serveless_cluster_name: The name  of serveless cluster
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

### 1. Create a new vpc, vswitch  for the cluster.

You can specify the following user-defined arguments:

* vpc_name: A new vpc name
* vpc_cidr: A new vpc cidr block
* vswitch_name: The name  of a vswitch
* vswitch_cidr: The of cidr blocks for a new vswitch

### 2. Using existing vpc and vswitch for the cluster.

You can specify the following user-defined arguments:

* vpc_id: A existing vpc ID
* vswitch_id: The of IDs for an existing vswitch




