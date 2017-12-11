alicloud-network-peering
========================

A terraform Alicloud template to build a VPC network peering environment. The template mainly provides following jobs:

- Create a VPC in each region cn-beijing and cn-hangzhou.
- Create serveral subnets in each VPC.
- Create a router interface in each VPC and support to connect two VPCs each other.


Input Variables
---------------

#### Common Input Vairables

- `alicloud_access_key` - The Alicloud Access Key ID to launch resources
- `alicloud_secret_key` - The Alicloud Access Secret Key to launch resources

#### VPC Input variables

- `vpc_cidr` - VPC CIDR block to launch a new VPC when `vpc_id` is not specified - default to "172.16.0.0/12"
- `vpc_name` - VPC name to mark a new VPC when `vpc_id` is not specified - default to "TF-VPC"
- `vpc_description` - VPC description used to launch a new vpc when 'vpc_id' is not specified - default to "A new VPC created by Terrafrom module tf-alicloud-vpc-cluster"

#### VSwitch Input Variables

- `vswitch_name` - VSwitch name prefix to mark a new VSwitch - default to "TF_VSwitch"
- `vswitch_cidr` - VSwitch CIDR block to launch several new VSwitches
- `vswitch_description` - VSwitch description used to describe new vswitch - default to "New VSwitch created by Terrafrom module tf-alicloud-vpc-cluster."

#### Router Interface Input Variables

- `opposite_region` - The opposite region to launch resources
- `interface_role` - The router interface role. Choices are 'InitiatingSide' and 'AcceptingSide'."
- `interface_spec - The router interface specification


Usage
-----
You can use this in your terraform template with the following steps.

1. Get module

        $ cd ./initiate
        $ terraform get

2. Plugin reinitialization

        $ terraform init

3. Planning phase

		$ terraform plan

4. Apply phase

		$ terraform apply

5. Destroy

		$ terraform destroy


Output Variables
-----------------------

- vpc_id - A new VPC ID
- vswitch_ids - A list of new VSwitch IDs
- availability_zones - A list of availability zones
- router_id - The virtual router ID in which new route entries are launched
- route_table_id - The route table ID in which new route entries are launched
