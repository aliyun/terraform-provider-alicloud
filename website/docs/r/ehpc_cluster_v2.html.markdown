---
subcategory: "Elastic High Performance Computing (Ehpc)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ehpc_cluster_v2"
description: |-
  Provides a Alicloud Ehpc Cluster V2 resource.
---

# alicloud_ehpc_cluster_v2

Provides a Ehpc Cluster V2 resource.

E-HPC Cluster Resources.

For information about Ehpc Cluster V2 and how to use it, see [What is Cluster V2](https://next.api.alibabacloud.com/document/EHPC/2024-07-30/CreateCluster).

-> **NOTE:** Available since v1.266.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ehpc_cluster_v2&exampleId=0c55a428-eb99-e56f-1d51-a5b38c350d4abddffba9&activeTab=example&spm=docs.r.ehpc_cluster_v2.0.0c55a428eb&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

resource "alicloud_vpc" "example" {
  is_default = false
  cidr_block = "10.0.0.0/24"
  vpc_name   = "example-cluster-vpc"
}

resource "alicloud_nas_access_group" "example" {
  access_group_type = "Vpc"
  description       = var.name
  access_group_name = "StandardMountTarget"
  file_system_type  = "standard"
}

resource "alicloud_nas_file_system" "example" {
  description  = "example-cluster-nas"
  storage_type = "Capacity"
  nfs_acl {
    enabled = false
  }
  zone_id          = "cn-hangzhou-k"
  encrypt_type     = "0"
  protocol_type    = "NFS"
  file_system_type = "standard"
}

resource "alicloud_vswitch" "example" {
  is_default   = false
  vpc_id       = alicloud_vpc.example.id
  zone_id      = "cn-hangzhou-k"
  cidr_block   = "10.0.0.0/24"
  vswitch_name = "example-cluster-vsw"
}

resource "alicloud_nas_access_rule" "example" {
  priority          = "1"
  access_group_name = alicloud_nas_access_group.example.access_group_name
  file_system_type  = alicloud_nas_file_system.example.file_system_type
  source_cidr_ip    = "10.0.0.0/24"
}

resource "alicloud_ecs_key_pair" "example" {
  key_pair_name = var.name
}

resource "alicloud_nas_mount_target" "example" {
  vpc_id            = alicloud_vpc.example.id
  network_type      = "Vpc"
  access_group_name = alicloud_nas_access_group.example.access_group_name
  vswitch_id        = alicloud_vswitch.example.id
  file_system_id    = alicloud_nas_file_system.example.id
}

resource "alicloud_security_group" "example" {
  vpc_id              = alicloud_vpc.example.id
  security_group_type = "normal"
}

resource "alicloud_ehpc_cluster_v2" "default" {
  cluster_credentials {
    key_pair_name = alicloud_ecs_key_pair.example.id
  }

  cluster_mode        = "Integrated"
  cluster_vpc_id      = alicloud_vpc.example.id
  deletion_protection = "true"
  shared_storages {
    mount_directory     = "/home"
    nas_directory       = "/"
    mount_target_domain = alicloud_nas_mount_target.example.mount_target_domain
    protocol_type       = "NFS"
    file_system_id      = alicloud_nas_file_system.example.id
    mount_options       = "-t nfs -o vers=3,nolock,proto=tcp,noresvport"
  }
  shared_storages {
    nas_directory       = "/"
    mount_target_domain = alicloud_nas_mount_target.example.mount_target_domain
    protocol_type       = "NFS"
    file_system_id      = alicloud_nas_file_system.example.id
    mount_options       = "-t nfs -o vers=3,nolock,proto=tcp,noresvport"
    mount_directory     = "/opt"
  }
  shared_storages {
    mount_options       = "-t nfs -o vers=3,nolock,proto=tcp,noresvport"
    mount_directory     = "/ehpcdata"
    nas_directory       = "/"
    mount_target_domain = alicloud_nas_mount_target.example.mount_target_domain
    protocol_type       = "NFS"
    file_system_id      = alicloud_nas_file_system.example.id
  }

  cluster_vswitch_id = alicloud_vswitch.example.id
  cluster_category   = "Standard"
  security_group_id  = alicloud_security_group.example.id
  cluster_name       = var.name
  manager {
    manager_node {
      spot_strategy = "NoSpot"
      system_disk {
        category = "cloud_essd"
        size     = "40"
        level    = "PL0"
      }
      enable_ht            = "true"
      instance_charge_type = "PostPaid"
      image_id             = "centos_7_6_x64_20G_alibase_20211130.vhd"
      instance_type        = "ecs.c6.xlarge"
    }
    scheduler {
      type    = "SLURM"
      version = "22.05.8"
    }
    dns {
      type    = "nis"
      version = "1.0"
    }
    directory_service {
      type    = "nis"
      version = "1.0"
    }
  }
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ehpc_cluster_v2&spm=docs.r.ehpc_cluster_v2.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `additional_packages` - (Optional, ForceNew, List, Available since v1.288.0) The list of software to be installed on the cluster. The value range of N is 0 to 10. See [`additional_packages`](#additional_packages) below.
* `addons` - (Optional, ForceNew, List) The cluster custom service component configuration. Only one component is supported. See [`addons`](#addons) below.
* `client_version` - (Optional, Computed) The version of the E-HPC client.
* `cluster_category` - (Optional, ForceNew) The cluster type. Valid values:

  - Standard
  - Serverless
* `cluster_credentials` - (Required, ForceNew, Set) Security credentials for the cluster. See [`cluster_credentials`](#cluster_credentials) below.
* `cluster_custom_configuration` - (Optional, List, Available since v1.288.0) The post-processing script configuration of the cluster. See [`cluster_custom_configuration`](#cluster_custom_configuration) below.
* `cluster_description` - (Optional, Available since v1.288.0) The description of the cluster. The description must be 2 to 128 characters in length. It can contain letters, digits, hyphens (-), and underscores (_).
* `cluster_mode` - (Optional, ForceNew) The deployment mode of the cluster. Valid values:

  - Integrated
  - Hybrid
  - Custom
* `cluster_name` - (Optional) The name of the cluster.
* `cluster_vswitch_id` - (Optional, ForceNew) The ID of the vSwitch that you want the cluster to use. The vSwitch must reside in the VPC that is specified by the `ClusterVpcId` parameter.
You can call the [DescribeVpcs](https://www.alibabacloud.com/help/en/doc-detail/448581.html) operation to query information about the created VPCs and vSwitches.
* `cluster_vpc_id` - (Optional, ForceNew) The ID of the virtual private cloud (VPC) in which the cluster resides.
* `deletion_protection` - (Optional) Specifies whether to enable deletion protection for the cluster. Valid values:

  - true
  - false
* `enable_scale_in` - (Optional, Computed, Bool, Available since v1.288.0) Specifies whether to enable auto scale-in for the cluster. Valid values:

  - true
  - false
* `enable_scale_out` - (Optional, Computed, Bool, Available since v1.288.0) Specifies whether to enable auto scale-out for the cluster. Valid values:

  - true
  - false
* `grow_interval` - (Optional, Computed, Int, Available since v1.288.0) The time interval of auto scale-out for the cluster.
* `idle_interval` - (Optional, Computed, Int, Available since v1.288.0) The idle duration of the compute nodes allowed by the cluster.
* `is_enterprise_security_group` - (Optional, ForceNew, Bool, Available since v1.288.0) Specifies whether to use an enterprise security group. Valid values:

  - true: An enterprise security group is automatically created and used.
  - false: A basic security group is automatically created and used.
* `manager` - (Optional, ForceNew, Set) The configurations of the cluster management node. See [`manager`](#manager) below.
* `max_core_count` - (Optional, Computed, Int, Available since v1.288.0) The total number of CPU cores of the compute nodes that the cluster can manage. Valid values: 0 to 100000.
* `max_count` - (Optional, Computed, Int, Available since v1.288.0) The number of compute nodes that the cluster can manage. Valid values: 0 to 5000.
* `monitor_spec` - (Optional, Computed, List, Available since v1.288.0) The monitoring configuration of the cluster. See [`monitor_spec`](#monitor_spec) below.
* `queues` - (Optional, ForceNew, List, Available since v1.288.0) The queue configurations of the cluster. The value range of N is 0 to 8. See [`queues`](#queues) below.
* `resource_group_id` - (Optional, ForceNew, Computed) The ID of the resource group to which the cluster belongs.
You can call the [ListResourceGroups](https://www.alibabacloud.com/help/en/doc-detail/158855.html) operation to obtain the IDs of the resource groups.
* `scheduler_spec` - (Optional, Computed, List, Available since v1.288.0) The scheduler configuration of the cluster. See [`scheduler_spec`](#scheduler_spec) below.
* `security_group_id` - (Optional, ForceNew) The security group ID.
* `shared_storages` - (Required, ForceNew, Set) List of cluster shared storage configurations. See [`shared_storages`](#shared_storages) below.
* `tags` - (Optional, ForceNew, Map, Available since v1.288.0) A mapping of tags to assign to the resource.

### `additional_packages`

The additional_packages supports the following:
* `name` - (Optional, ForceNew) The name of the software to be installed.
* `version` - (Optional, ForceNew) The version of the software to be installed.

### `addons`

The addons supports the following:
* `name` - (Required, ForceNew) Customize the specific configuration information of the service component.
* `resources_spec` - (Optional, ForceNew) Customize the resource configuration of the service component.
* `services_spec` - (Optional, ForceNew) Customize the service configuration of the service component.
* `version` - (Required, ForceNew) Customize the service component version.

### `cluster_credentials`

The cluster_credentials supports the following:
* `key_pair_name` - (Optional, ForceNew, Available since v1.270.0) The SSH key of root of the cluster node.
* `password` - (Optional, ForceNew) The root password of the cluster node. It is 8 to 20 characters in length and must contain three types of characters: uppercase and lowercase letters, numbers, and special symbols. Special symbols can be: () ~! @ # $ % ^ & * - = + { } [ ] : ; ',. ? /

### `cluster_custom_configuration`

The cluster_custom_configuration supports the following:
* `script` - (Optional, Available since v1.288.0) The download URL of the post-processing script.
* `args` - (Optional, Available since v1.288.0) The execution parameters of the post-processing script.

### `manager`

The manager supports the following:
* `directory_service` - (Optional, ForceNew, Set) The configurations of the domain account service. See [`directory_service`](#manager-directory_service) below.
* `dns` - (Optional, ForceNew, Set) The configurations of the domain name resolution service. See [`dns`](#manager-dns) below.
* `manager_node` - (Optional, ForceNew, Set) The hardware configurations of the management node. See [`manager_node`](#manager-manager_node) below.
* `scheduler` - (Optional, ForceNew, Set) The configurations of the scheduler service. See [`scheduler`](#manager-scheduler) below.

### `manager-directory_service`

The manager-directory_service supports the following:
* `type` - (Optional, ForceNew) The type of the domain account. Valid values:

  - NIS
* `version` - (Optional, ForceNew) The version of the domain account service.

### `manager-dns`

The manager-dns supports the following:
* `type` - (Optional, ForceNew) The domain name resolution type. Valid values:

  - NIS
* `version` - (Optional, ForceNew) The version of the domain name resolution service.

### `manager-manager_node`

The manager-manager_node supports the following:
* `auto_renew` - (Optional, ForceNew) Whether to automatically renew. This parameter takes effect only when the value of InstanceChargeType is PrePaid. Value range:
  - true: Automatic renewal.
  - false: Do not renew automatically (default).
* `auto_renew_period` - (Optional, ForceNew, Int) The renewal duration of a single automatic renewal. Value range:
  - When PeriodUnit = Week: 1, 2, 3.
  - When PeriodUnit = Month: 1, 2, 3, 6, 12, 24, 36, 48, 60.

Default value: 1.
* `duration` - (Optional, ForceNew, Computed, Int) The duration of the preemptible instance, in hours. Value:
  - : After the instance is created, Alibaba Cloud will ensure that the instance will not be automatically released after one hour of operation. After one hour, the system will compare the bid price with the market price in real time and check the resource inventory to determine the holding and recycling of the instance.
  - 0: After creation, Alibaba Cloud does not guarantee the running time of the instance. The system compares the bid price with the market price in real time and checks the resource inventory to determine the holding and recycling of the instance.

Default value: 1.
* `enable_ht` - (Optional, ForceNew) Specifies whether to enable hyper-threading on the management node.
* `image_id` - (Optional, ForceNew) The ID of the image used by the management node.
* `instance_charge_type` - (Optional, ForceNew) The instance billing method of the management node. Valid values:

  - PostPaid: pay-as-you-go
  - PrePaid: subscription
* `instance_type` - (Optional, ForceNew) The instance type of the management node.
* `period` - (Optional, ForceNew, Int) The duration of the resource purchase. The unit is specified by PeriodUnit. The parameter InstanceChargeType takes effect only when the value is PrePaid and is a required value. Once DedicatedHostId is specified, the value range cannot exceed the subscription duration of the DDH. Value range:
  - When PeriodUnit = Week, the values of Period are 1, 2, 3, and 4.
  - When PeriodUnit = Month, the values of Period are 1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36, 48, and 60.
* `period_unit` - (Optional, ForceNew) The unit of duration of the year-to-month billing method. Value range:
  - Week.
  - Month (default).
* `spot_price_limit` - (Optional, ForceNew, Float) Set the maximum price per hour for the instance. The maximum number of decimals is 3. It takes effect when the value of the SpotStrategy parameter is SpotWithPriceLimit.
* `spot_strategy` - (Optional, ForceNew) The bidding strategy for pay-as-you-go instances. This parameter takes effect when the value of the InstanceChargeType parameter is PostPaid. Value range:
  - NoSpot: normal pay-as-you-go instances (default).
  - SpotWithPriceLimit: set the upper limit price for the preemptible instance.
  - SpotAsPriceGo: The system automatically bids, following the actual price of the current market.
* `system_disk` - (Optional, ForceNew, Set) System disk configuration of the management node. See [`system_disk`](#manager-manager_node-system_disk) below.

### `manager-scheduler`

The manager-scheduler supports the following:
* `type` - (Optional, ForceNew) The scheduler type. Valid values:

  - SLURM
  - PBS
  - OPENGRIDSCHEDULER
  - LSF_PLUGIN
  - PBS_PLUGIN
* `version` - (Optional, ForceNew) The scheduler version.

### `manager-manager_node-system_disk`

The manager-manager_node-system_disk supports the following:
* `category` - (Optional, ForceNew) Manage the system disk configuration of the node. Value range:
  - cloud_efficiency: The Ultra cloud disk.
  - cloud_ssd:SSD cloud disk.
  - cloud_essd:ESSD cloud disk.
  - cloud: ordinary cloud disk.
* `level` - (Optional, ForceNew) When creating an ESSD cloud disk to use as a system disk, set the performance level of the cloud disk. Value range:
  - PL0: maximum random read/write IOPS 10000 for a single disk.
  - PL1 (default): Maximum random read/write IOPS 50000 for a single disk.
  - PL2: maximum random read/write IOPS 100000 for a single disk.
  - PL3: maximum random read/write IOPS 1 million for a single disk.
* `size` - (Optional, ForceNew, Int) The system disk size of the management node. Unit: GiB. Value range:
  - Ordinary cloud tray: 20~500.
  - ESSD cloud disk:
  - PL0:1~2048.
  - PL1:20~2048.
  - PL2:461~2048.
  - PL3:1261~2048.
  - Other cloud disk types: 20~2048.

### `monitor_spec`

The monitor_spec supports the following:
* `enable_compute_load_monitor` - (Optional, Computed, Bool, Available since v1.288.0) Specifies whether to enable the monitoring component for the compute nodes. Valid values:

  - true
  - false

### `queues`

The queues supports the following:
* `queue_name` - (Optional, ForceNew, Available since v1.288.0) The name of the queue. The name must be 1 to 15 characters in length. It can contain letters, digits, and periods (.).
* `enable_scale_out` - (Optional, ForceNew, Bool, Available since v1.288.0) Specifies whether to enable auto scale-out for the queue. Valid values:

  - true
  - false
* `enable_scale_in` - (Optional, ForceNew, Bool, Available since v1.288.0) Specifies whether to enable auto scale-in for the queue. Valid values:

  - true
  - false
* `min_count` - (Optional, ForceNew, Int, Available since v1.288.0) The minimum number of compute nodes that the queue retains.
* `max_count` - (Optional, ForceNew, Int, Available since v1.288.0) The maximum number of compute nodes that the queue can retain.
* `initial_count` - (Optional, ForceNew, Int, Available since v1.288.0) The initial number of compute nodes that the queue retains.
* `inter_connect` - (Optional, ForceNew, Available since v1.288.0) The network type between the compute nodes in the queue. Valid values:

  - vpc
  - eRDMA
* `vswitch_ids` - (Optional, ForceNew, List, Available since v1.288.0) The list of vSwitches available to the compute nodes in the queue. The value range of N is 1 to 5.
* `compute_nodes` - (Optional, ForceNew, List, Available since v1.288.0) The list of hardware configurations of the compute nodes in the queue. The value range of N is 0 to 10. See [`compute_nodes`](#queues-compute_nodes) below.
* `allocation_strategy` - (Optional, ForceNew, Available since v1.288.0) The auto scale-out strategy of the queue.
* `ram_role` - (Optional, ForceNew, Available since v1.288.0) The name of the instance role attached to the compute nodes in the queue.
* `hostname_prefix` - (Optional, ForceNew, Available since v1.288.0) The hostname prefix of the compute nodes in the queue.
* `hostname_suffix` - (Optional, ForceNew, Available since v1.288.0) The hostname suffix of the compute nodes in the queue.
* `keep_alive_nodes` - (Optional, ForceNew, List, Available since v1.288.0) The list of nodes with deletion protection enabled in the queue. The value is the hostname of the node.
* `max_count_per_cycle` - (Optional, ForceNew, Int, Available since v1.288.0) The maximum number of compute nodes that the queue can scale out in each scale-out cycle.
* `reserved_node_pool_id` - (Optional, ForceNew, Available since v1.288.0) The ID of the reserved node pool used by the queue.

### `queues-compute_nodes`

The queues-compute_nodes supports the following:
* `instance_type` - (Optional, ForceNew, Available since v1.288.0) The instance type of the compute node.
* `image_id` - (Optional, ForceNew, Available since v1.288.0) The ID of the image used by the compute node.
* `instance_charge_type` - (Optional, ForceNew, Available since v1.288.0) The instance billing method of the compute node. Valid values:

  - PostPaid: pay-as-you-go
  - PrePaid: subscription
* `period_unit` - (Optional, ForceNew, Available since v1.288.0) The unit of duration of the year-to-month billing method. Valid values:

  - Week
  - Month
* `period` - (Optional, ForceNew, Int, Available since v1.288.0) The duration of the resource purchase. The unit is specified by `period_unit`.
* `auto_renew` - (Optional, ForceNew, Bool, Available since v1.288.0) Specifies whether to enable auto-renewal. This parameter takes effect only when `instance_charge_type` is set to PrePaid.
* `auto_renew_period` - (Optional, ForceNew, Int, Available since v1.288.0) The renewal duration of a single automatic renewal.
* `spot_strategy` - (Optional, ForceNew, Available since v1.288.0) The bidding strategy for pay-as-you-go instances. Valid values:

  - NoSpot: normal pay-as-you-go instances (default).
  - SpotWithPriceLimit: set the upper limit price for the preemptible instance.
  - SpotAsPriceGo: The system automatically bids, following the actual price of the current market.
* `spot_price_limit` - (Optional, ForceNew, Float, Available since v1.288.0) The maximum price per hour for the instance. The maximum number of decimals is 3. It takes effect when `spot_strategy` is set to SpotWithPriceLimit.
* `duration` - (Optional, ForceNew, Int, Available since v1.288.0) The duration of the preemptible instance, in hours.
* `enable_ht` - (Optional, ForceNew, Bool, Available since v1.288.0) Specifies whether to enable hyper-threading on the compute node.
* `system_disk` - (Optional, ForceNew, List, Available since v1.288.0) The system disk configuration of the compute node. See [`system_disk`](#queues-compute_nodes-system_disk) below.

### `queues-compute_nodes-system_disk`

The queues-compute_nodes-system_disk supports the following:
* `category` - (Optional, ForceNew, Available since v1.288.0) The category of the system disk. Valid values:

  - cloud_efficiency: The Ultra cloud disk.
  - cloud_ssd: SSD cloud disk.
  - cloud_essd: ESSD cloud disk.
  - cloud: ordinary cloud disk.
* `size` - (Optional, ForceNew, Int, Available since v1.288.0) The system disk size of the compute node. Unit: GiB.
* `level` - (Optional, ForceNew, Available since v1.288.0) The performance level of the ESSD cloud disk used as the system disk. Valid values:

  - PL0
  - PL1
  - PL2
  - PL3

### `scheduler_spec`

The scheduler_spec supports the following:
* `enable_topology_awareness` - (Optional, Computed, Bool, Available since v1.288.0) Specifies whether to enable the topology awareness feature for the cluster. Valid values:

  - true
  - false

### `shared_storages`

The shared_storages supports the following:
* `file_system_id` - (Optional, ForceNew) The ID of the mounted file system.
* `mount_directory` - (Optional, ForceNew) The local Mount directory where the file system is mounted.
* `mount_options` - (Optional, ForceNew) Storage mount options for the mounted file system.
* `mount_target_domain` - (Optional, ForceNew) The mount point address of the mounted file system.
* `nas_directory` - (Optional, ForceNew) The remote directory to which the mounted file system needs to be mounted.
* `protocol_type` - (Optional, ForceNew) The protocol type of the mounted file system. Value range:
  - NFS
  - SMB

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `cluster_status` - (Available since v1.288.0) The status of the cluster.
* `create_time` - The time when the cluster was created.
* `ehpc_version` - (Available since v1.288.0) The version of the E-HPC cluster.
* `modify_time` - (Available since v1.288.0) The time when the cluster was modified.
* `manager` - The configurations of the cluster management node.
  * `manager_node` - The hardware configurations of the management node.
    * `expired_time` - The expiration time of the management node.
    * `instance_id` - The instance ID of the management node.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 8 mins) Used when create the Cluster V2.
* `delete` - (Defaults to 5 mins) Used when delete the Cluster V2.
* `update` - (Defaults to 5 mins) Used when update the Cluster V2.

## Import

Ehpc Cluster V2 can be imported using the id, e.g.

```shell
$ terraform import alicloud_ehpc_cluster_v2.example <cluster_id>
```
