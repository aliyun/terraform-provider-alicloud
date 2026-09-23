alicloud-network-peering
========================

A terraform Alicloud template to build a VPC network peering environment. The template mainly provides following jobs:

- Create a VPC in each region cn-beijing and cn-hangzhou.
- Create a router interface in each VPC and then connect each other.


Input Variables
---------------

#### VPC Input variables

- `init_vpc_id` - A existing VPC ID used to launch a new initiating side router interface
- `init_vpc_cidr` - VPC CIDR block to launch a new VPC when `init_vpc_id` is not specified - default to "10.0.0.0/8"
- `accpet_vpc_id` - A existing VPC ID used to launch a new accepting side router interface
- `accept_vpc_cidr` - VPC CIDR block to launch a new VPC when `accept_vpc_id` is not specified - default to "172.16.0.0/12"
- `vpc_name` - VPC name to mark a new VPC when `init_vpc_id`  or `accept_vpc_id` is not specified - default to "TF-VPC"
- `vpc_description` - VPC description used to launch a new vpc when 'vpc_id' is not specified - default to "A new VPC created by Terrafrom module tf-alicloud-vpc-cluster"

#### Router Interface Input Variables

- `region` - The initiating side region
- `opposite_region` - The accepting side region
- `interface_spec - The initiating side router interface specification


Usage
-----
You can use this in your terraform template with the following steps.


1. Plugin reinitialization

        $ terraform init

2. Planning phase

		$ terraform plan

3. Apply phase

		$ terraform apply

4. Destroy

		$ terraform destroy


Output Variables
-----------------------

- vpc_id - The initiating side VPC ID
- accpeting_vpc_id - The initiating side VPC ID
- interface_id - The initiating side route interface ID
- accepting_interface_id - The accepting side route interface ID
