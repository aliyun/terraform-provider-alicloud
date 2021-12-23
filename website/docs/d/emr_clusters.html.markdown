---
subcategory: "E-MapReduce"
layout: "alicloud"
page_title: "Alicloud: alicloud_emr_clusters"
sidebar_current: "docs-alicloud-datasource-emr-clusters"
description: |-
  Provides a list of Emr Clusters to the user.
---

# alicloud\_emr\_clusters

This data source provides the Emr Clusters of the current Alibaba Cloud user.

-> **NOTE:** Available in v1.146.0+.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-testAccClusters"
}

data "alicloud_resource_manager_resource_groups" "default" {
  status = "OK"
}

data "alicloud_emr_main_versions" "default" {}

data "alicloud_emr_instance_types" "default" {
  destination_resource  = "InstanceType"
  cluster_type          = data.alicloud_emr_main_versions.default.main_versions.0.cluster_types.0
  support_local_storage = false
  instance_charge_type  = "PostPaid"
  support_node_type     = ["MASTER", "CORE", "TASK"]
}

data "alicloud_emr_disk_types" "data_disk" {
  destination_resource = "DataDisk"
  cluster_type         = data.alicloud_emr_main_versions.default.main_versions.0.cluster_types.0
  instance_charge_type = "PostPaid"
  instance_type        = data.alicloud_emr_instance_types.default.types.0.id
  zone_id              = data.alicloud_emr_instance_types.default.types.0.zone_id
}

data "alicloud_emr_disk_types" "system_disk" {
  destination_resource = "SystemDisk"
  cluster_type         = data.alicloud_emr_main_versions.default.main_versions.0.cluster_types.0
  instance_charge_type = "PostPaid"
  instance_type        = data.alicloud_emr_instance_types.default.types.0.id
  zone_id              = data.alicloud_emr_instance_types.default.types.0.zone_id
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = data.alicloud_vpcs.default.ids.0
}


data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}

data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_emr_instance_types.default.types.0.zone_id
}

resource "alicloud_ram_role" "default" {
  name        = var.name
  document    = <<EOF
    {
        "Statement": [
        {
            "Action": "sts:AssumeRole",
            "Effect": "Allow",
            "Principal": {
            "Service": [
                "emr.aliyuncs.com",
                "ecs.aliyuncs.com"
            ]
            }
        }
        ],
        "Version": "1"
    }
    EOF
  description = "this is a role test."
  force       = true
}

resource "alicloud_emr_cluster" "default" {
  name    = var.name
  emr_ver = data.alicloud_emr_main_versions.default.main_versions.0.emr_version

  cluster_type = data.alicloud_emr_main_versions.default.main_versions.0.cluster_types.0

  host_group {
    host_group_name   = "master_group"
    host_group_type   = "MASTER"
    node_count        = "2"
    instance_type     = data.alicloud_emr_instance_types.default.types.0.id
    disk_type         = data.alicloud_emr_disk_types.data_disk.types.0.value
    disk_capacity     = data.alicloud_emr_disk_types.data_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.data_disk.types.0.min : 160
    disk_count        = "1"
    sys_disk_type     = data.alicloud_emr_disk_types.system_disk.types.0.value
    sys_disk_capacity = data.alicloud_emr_disk_types.system_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.system_disk.types.0.min : 160
  }

  host_group {
    host_group_name   = "core_group"
    host_group_type   = "CORE"
    node_count        = "3"
    instance_type     = data.alicloud_emr_instance_types.default.types.0.id
    disk_type         = data.alicloud_emr_disk_types.data_disk.types.0.value
    disk_capacity     = data.alicloud_emr_disk_types.data_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.data_disk.types.0.min : 160
    disk_count        = "4"
    sys_disk_type     = data.alicloud_emr_disk_types.system_disk.types.0.value
    sys_disk_capacity = data.alicloud_emr_disk_types.system_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.system_disk.types.0.min : 160
  }

  host_group {
    host_group_name   = "task_group"
    host_group_type   = "TASK"
    node_count        = "2"
    instance_type     = data.alicloud_emr_instance_types.default.types.0.id
    disk_type         = data.alicloud_emr_disk_types.data_disk.types.0.value
    disk_capacity     = data.alicloud_emr_disk_types.data_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.data_disk.types.0.min : 160
    disk_count        = "4"
    sys_disk_type     = data.alicloud_emr_disk_types.system_disk.types.0.value
    sys_disk_capacity = data.alicloud_emr_disk_types.system_disk.types.0.min > 160 ? data.alicloud_emr_disk_types.system_disk.types.0.min : 160
  }

  high_availability_enable  = true
  zone_id                   = data.alicloud_emr_instance_types.default.types.0.zone_id
  security_group_id         = alicloud_security_group.default.id
  is_open_public_ip         = true
  charge_type               = "PostPaid"
  vswitch_id                = data.alicloud_vswitches.default.ids.0
  user_defined_emr_ecs_role = alicloud_ram_role.default.name
  ssh_enable                = true
  master_pwd                = "ABCtest1234!"
  tags = {
    Created = "TF"
    For     = "acceptance test"
  }
}

data "alicloud_emr_clusters" "ids" {}
output "emr_cluster_id_1" {
  value = data.alicloud_emr_clusters.ids.clusters.0.id
}

data "alicloud_emr_clusters" "nameRegex" {
  name_regex = alicloud_emr_cluster.default.name
}
output "emr_cluster_id_2" {
  value = data.alicloud_emr_clusters.nameRegex.clusters.0.id
}

```

## Argument Reference

The following arguments are supported:

* `cluster_name` - (Optional, ForceNew) The cluster name.
* `cluster_type_list` - (Optional, ForceNew) The cluster type list.
* `create_type` - (Optional, ForceNew) How to create a cluster. Valid values: `ON-DEMAND`, `MANUAL`.
* `default_status` - (Optional, ForceNew) The default status.
* `deposit_type` - (Optional, ForceNew) The hosting type of the cluster. Valid values: `HALF_MANAGED`, `MANAGED`.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew, Computed)  A list of Cluster IDs.
* `is_desc` - (Optional, ForceNew) The is desc.
* `machine_type` - (Optional, ForceNew) The host type of the cluster. The default is ECS. Valid values: `DOCKER`, `ECS`, `PYHSICAL_MACHINE`, `ECS_FROM_ECM_HOSTPOOL`.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Cluster name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `resource_group_id` - (Optional, ForceNew) The Resource Group ID.
* `status_list` - (Optional, ForceNew) The status list. Valid values: `ABNORMAL`, `CREATE_FAILED`, `CREATING`, `IDLE`, `RELEASED`, `RELEASE_FAILED`, `RELEASING`, `RUNNING`, `WAIT_FOR_PAY`.
* `vpc_id` - (Optional, ForceNew) The VPC ID.

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Cluster names.
* `clusters` - A list of Emr Clusters. Each element contains the following attributes:
	* `access_info` - Cluster connection information.
	  * `zk_links` - Link address information list of ZooKeeper.
			* `link` - The access link address of ZooKeeper.
			* `port` - The port of ZooKeeper.
	* `auto_scaling_allowed` - Whether flexible expansion is allowed.
	* `auto_scaling_by_load_allowed` - Whether to allow expansion by load.
	* `auto_scaling_enable` - Whether to enable elastic expansion.
	* `auto_scaling_spot_with_limit_allowed` - Whether to allow the use of elastic scaling bidding instances.
	* `bootstrap_action_list` - List of boot actions.
		* `name` - The name of the boot operation.
		* `path` - Boot operation script path.
		* `arg` - Parameters of the boot operation.
	* `bootstrap_failed` - The result of the boot operation.
	* `cluster_id` - The first ID of the resource.
	* `cluster_name` - The ClusterName.
	* `create_resource` - Cluster tag, no need to pay attention.
	* `create_time` - The creation time of the resource.
	* `create_type` - How to create a cluster.
	* `deposit_type` - The hosting type of the cluster.
	* `eas_enable` - High security cluster.
	* `expired_time` - The expiration time of the cluster.
	* `extra_info` - Additional information for Stack.
	* `high_availability_enable` - High availability cluster.
	* `host_group_list` - List of cluster machine groups.
		* `host_group_name` - The name of the machine group.
		* `instance_type` - Machine Group instance.
		* `nodes` - Machine node.
			* `create_time` - Creation time.
			* `disk_infos` - Disk information.
				* `device` - The disk name.
				* `disk_id` - The ID of the disk.
				* `disk_name` - The disk name.
				* `size` - Disk capacity.
				* `type` - Disk type.
			* `expired_time` - Timeout time.
			* `inner_ip` - The Intranet IP of the EMR.
			* `emr_expired_time` - The timeout of the EMR.
			* `instance_id` - The ID of the ECS instance.
			* `pub_ip` - Public IP address.
			* `status` - Status.
			* `support_ipv6` - Whether IPV6 is supported.
			* `zone_id` - The zone ID.
		* `band_width` - Bandwidth.
		* `disk_capacity` - Data disk capacity.
		* `disk_count` - The number of data disks.
		* `disk_type` - System disk type:
		* `memory_capacity` - Memory size.
		* `node_count` - The number of machine group nodes.
		* `period` - Package year and month time (days).
		* `charge_type` - Payment Type.
		* `cpu_core` - The number of CPU cores.
		* `host_group_change_type` - The current operation type of the machine Group:
		* `host_group_id` - The ID of the machine group.
		* `host_group_type` - Role of host in cluster:
	* `host_pool_info` - Machine pool information.
		* `hp_biz_id` - Machine pool ID.
		* `hp_name` - The name of the machine pool.
	* `image_id` - The ID of the image used to create the cluster.
	* `local_meta_db` - Whether to use Hive local Metabase.
	* `machine_type` - The host type of the cluster. The default is ECS.
	* `meta_store_type` - Metadata type:
	* `net_type` - Cluster network type.
	* `payment_type` - The payment type of the resource.
	* `period` - The package year and month time of the machine group. The Valid Values : `1`, `2`, `3`, `4`, `5`, `6`, `7`, `8`, `9`, `12`, `24`, `36`.
	* `relate_cluster_info` - The information of the primary cluster associated with the Gateway.
		* `cluster_id` - The ID of the associated cluster.
		* `cluster_name` - The name of the associated cluster.
		* `cluster_type` - The cluster type of the associated cluster.
		* `status` - The status  of the associated cluster.
	* `resize_disk_enable` - Whether to allow disk expansion:
	* `running_time` - The time (in seconds) that has been running.
	* `security_group_id` - The ID of the security group.
	* `security_group_name` - The name of the security group.
	* `software_info` - Service list.
		* `cluster_type` - Cluster type:
		* `emr_ver` - E-MapReduce version number.
		* `softwares` - Service list.
			* `display_name` - The name of the service.
			* `name` - The internal name of the service.
			* `only_display` - Whether it shows.
			* `start_tpe` - Startup type.
			* `version` - Service version.
	* `start_time` - Cluster startup time.
	* `status` - The cluster status.
	* `stop_time` - Cluster stop time.
	* `tags` - A mapping of tags to assign to the resource.
	* `user_defined_emr_ecs_role` - The EMR permission name used.
	* `user_id` - The user ID.
	* `vpc_id` - The VPC ID.
	* `vswitch_id` - The vswitch id.
	* `zone_id` - The zone ID.