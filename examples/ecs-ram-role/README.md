Alicloud ECS Instance Launched In VPC and Attach A RAM Role Examples
====================================================================

A terraform example to provide ECS instances in Alicloud. These instances will be launched in the VPC network.
After creating ECS instances, this example will attach a new RAM role to them.

- The example contains one VPC, one VSwitch, one Security Group, several Security Group Rules, one RAM role, one RAM policy, several Disks and Instances.
- If VPC, VSwitch or Security Group is not specified, the module will launch a new one using its own parameters.
- If you have no idea some parametes, such as instance type, availability zone and image id, the module will provide default values by these data source.


Module Input Variables
----------------------

The example aim to create one or more instances and disks in the VPC, and then attach RAM role to them. Its input variables contains VPC, VSwitch, Security Group, Security Group Rules, ECS Disks and ECS Instances.

#### Common Imput vairables

- `alicloud_access_key` - The Alicloud Access Key ID to launch resources
- `alicloud_secret_key` - The Alicloud Access Secret Key to launch resources
- `region` - The region to launch resources
- `zone_id` - The availability zone ID to launch VSwitch, ECS Instances and ECS Disks - default to a zone ID retrieved by zones' data source
- `number_format` - The number format used to mark multiple resources - default to "%02d"

#### VPC Input variables

- `vpc_id` - VPC ID to launch a new VSwitch and Security Group
- `vpc_name` - VPC name to mark a new VPC when `vpc_id` is not specified - default to "TF-VPC"
- `vpc_cidr` - VPC CIDR block to launch a new VPC when `vpc_id` is not specified - default to "172.16.0.0/12"

#### VSwitch Input variables

- `vswitch_id` - VSwitch ID to launch new ECS instances
- `vswitch_name` - VSwitch name to mark a new VSwitch when `vswitch_id` is not specified - default to "TF_VSwitch"
- `vswitch_cidr` - VSwitch CIDR block to launch a new VSwitch when `vswitch_id` is not specified. It has a default value '172.16.0.0/16' according `vpc_cidr's` default value.

`NOTE`: One of the `vswitch_id` and `vswitch_cidr` is required.

#### Security Group Input variables

- `sg_id` - Security Group ID to configure rules and launch new ECS instances
- `sg_name` - Security Group name to mark a new Security Group when `sg_id` is not specified - default to "TF_Security_Group"
- `ip_protocols` - List of IP protocols to configure Security Group rules - item choices: ["tcp", "udp", "icmp", "gre", "all"]
- `rule_directions` - List of directions to configure Security Group rules - item choices: ["ingress", "egress"] - default to ["ingress"]
- `policies` - List of policies to configure Security Group rules - item choices: ["accept", "drop"] - default to ["accept"]
- `port_ranges` - List of port ranges to configure Security Group rules - default to ["-1/-1"]
- `priorities` - List of priorities to configure Security Group rules - item choices: [1-100] - default to [1]
- `cidr_ips` - List of CIDR IPs to configure Security Group rules - default to ["0.0.0.0/0"]

`NOTE`:
1. The number of Security Group rules depends on the size of `ip_protocols`
2. All of the Security Group rules' network type are `intranet`

#### RAM Role Input variables

- `ram_role_name` - The name of RAM Role - default to "TF-RAM-Role-Name"
- `ram_role_ram_users` - List of RAM users used to assume RAM role document
- `ram_role_services` - List of services used to assume RAM role document - default to ["ecs.aliyuncs.com"]
- `ram_role_terminate_force` - Whether release relationship forcibly when deleting RAM role - default to `true`

`NOTE`:
1. At least one of 'ram_role_ram_users' and 'ram_role_services' must be set"

#### RAM Policy Input variables

- `ram_policy_name` - The name of RAM Policy - default to "TF-RAM-Policy-Name"
- `ram_policy_statement_effect` - The statement effect of RAM policy document - choices: ["Allow", "Deny"] - default to "Allow"
- `ram_policy_statement_action` - List of statement actions to assemble RAM policy document - like ["oss:Get*", "oss:List*"]
- `ram_policy_statement_resource` - List of statement resources to assemble RAM policy document - like ["acs:oss:*:*:*"]
- `ram_policy_terminate_force` - Whether release relationship forcibly when deleting RAM role - default to `true`


#### ECS Disk Input variables

- `number_of_disks` - The number disks you want to launch - default to 0
- `disk_name` - ECS disk name to mark data disk(s) - default to "TF_ECS_Disk"
- `disk_category` - ECS disk category to launch data disk(s) - choices to ["cloud_ssd", "cloud_efficiency"] - default to "cloud_efficiency"
- `disk_size` - ECS disk size to launch data disk(s) - default to 40
- `disk_tags` - A map for setting ECS disk tags - default to

      disk_tags = {
          created_by = "Terraform"
          created_from = "module-tf-alicloud-ecs-instance"
      }

#### ECS Instance Input variables

- `number_of_instances` - The number of instances you want to launch - default to 1
- `image_id` - The image id to use - default to an Ubuntu-64bit image ID retrieved by images' data source
- `instance_type` - The ECS instance type, e.g. ecs.n4.small, - default to a 1Core 2GB instance type retrieved by instance_types' data source
- `instance_name` - ECS instance name to mark instance(s) - default to "TF_ECS_Instance"
- `host_name` - ECS instance host name to configure instance(s) - default to "TF_ECS_Host_Name"
- `system_category` - ECS disk category to launch system disk - choices to ["cloud_ssd", "cloud_efficiency"] - default to "cloud_efficiency"
- `system_size` - ECS disk size to launch system disk - default to 40
- `allocate_public_ip` - Whether to allocate public for instance(s) - default to true
- `internet_charge_type` - The internet charge type for setting instance network - choices["PayByTraffic", "PayByBandwidth"] - default to "PayByTraffic"
- `internet_max_bandwidth_out` - The max out bandwidth for setting instance network - default to 10
- `instance_charge_type` - The instance charge type - choices to ["PrePaid", "PostPaid"] - default to "PostPaid"
- `period` - The instance charge period when instance charge type is 'PrePaid' - default to 1
- `key_name` - The instance key pair name for SSH keys
- `password` - The instance password
- `instance_tags` - A map for setting ECS Instance tags - default to

      instance_tags = {
          created_by = "Terraform"
          created_from = "module-tf-alicloud-ecs-instance"
      }


Usage
-----
You can input and specify some parameters in the variables.tf, and then execute the following commands to create and manage them:

* Planning phase

		terraform plan

* Apply phase

		terraform apply


* Destroy

		terraform destroy

Module Output Variables
-----------------------

- instance_ids - List of new instance ids
- disk_ids - List of new data disk ids
- vpc_id - A new VPC ID
- vswitch_id - A new VSwitch ID
- security_group_id - A new Security Group ID

Authors
-------
Created and maintained by He Guimin(heguimin36@163.com)

License
-------
Apache 2 Licensed. See LICENSE for full details.