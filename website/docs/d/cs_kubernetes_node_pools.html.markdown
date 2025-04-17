---
subcategory: "Container Service for Kubernetes (ACK)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_kubernetes_node_pools"
sidebar_current: "docs-alicloud-datasource-cs-kubernetes-node-pools"
description: |-
  Provides a list of Ack Nodepool owned by an Alibaba Cloud account.
---

# alicloud_cs_kubernetes_node_pools

This data source provides Ack Nodepool available to the user.[What is Nodepool](https://next.api.alibabacloud.com/document/CS/2015-12-15/CreateClusterNodePool)

-> **NOTE:** Available since v1.246.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_enhanced_nat_available_zones" "enhanced" {
}

data "alicloud_instance_types" "cloud_efficiency" {
  availability_zone    = data.alicloud_enhanced_nat_available_zones.enhanced.zones.0.zone_id
  cpu_core_count       = 4
  memory_size          = 8
  kubernetes_node_role = "Worker"
  system_disk_category = "cloud_efficiency"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "10.4.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vswitch_name = var.name
  cidr_block   = "10.4.0.0/24"
  vpc_id       = alicloud_vpc.default.id
  zone_id      = data.alicloud_enhanced_nat_available_zones.enhanced.zones.0.zone_id
}

resource "alicloud_cs_managed_kubernetes" "default" {
  name_prefix          = var.name
  cluster_spec         = "ack.pro.small"
  vswitch_ids          = [alicloud_vswitch.default.id]
  new_nat_gateway      = true
  pod_cidr             = cidrsubnet("10.0.0.0/8", 8, 36)
  service_cidr         = cidrsubnet("172.16.0.0/16", 4, 7)
  slb_internet_enabled = true
  enable_rrsa          = true
}

resource "alicloud_key_pair" "default" {
  key_pair_name = var.name
}

resource "alicloud_cs_kubernetes_node_pool" "default" {
  node_pool_name       = "spot_auto_scaling"
  cluster_id           = alicloud_cs_managed_kubernetes.default.id
  vswitch_ids          = [alicloud_vswitch.default.id]
  instance_types       = [data.alicloud_instance_types.cloud_efficiency.instance_types.0.id]
  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  key_name             = alicloud_key_pair.default.key_pair_name

  # automatic scaling node pool configuration.
  scaling_config {
    min_size = 1
    max_size = 10
    type     = "spot"
  }
  # spot price config
  spot_strategy = "SpotWithPriceLimit"
  spot_price_limit {
    instance_type = data.alicloud_instance_types.cloud_efficiency.instance_types.0.id
    price_limit   = "0.70"
  }
}

data "alicloud_cs_kubernetes_node_pools" "default" {
  ids        = ["${alicloud_cs_kubernetes_node_pool.default.node_pool_id}"]
  cluster_id = alicloud_cs_managed_kubernetes.default.id
}

output "alicloud_cs_kubernetes_node_pool_example_id" {
  value = data.alicloud_cs_kubernetes_node_pools.default.nodepools.0.node_pool_id
}
```

## Argument Reference

The following arguments are supported:
* `cluster_id` - (Required, ForceNew) The id of kubernetes cluster.
* `node_pool_name` - (ForceNew, Optional) The name of node pool.
* `ids` - (Optional, ForceNew, Computed) A list of Nodepool IDs.
* `output_file` - (Optional, ForceNew) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Nodepool IDs.
* `nodepools` - A list of Nodepool Entries. Each element contains the following attributes:
  * `auto_renew` - Whether to enable automatic renewal for nodes in the node pool takes effect only when `instance_charge_type` is set to `PrePaid`. Default value: `false`. Valid values:- `true`: Automatic renewal. - `false`: Do not renew automatically.
  * `auto_renew_period` - The automatic renewal period of nodes in the node pool takes effect only when you select Prepaid and Automatic Renewal, and is a required value. When `PeriodUnit = Month`, the value range is {1, 2, 3, 6, 12}. Default value: 1.
  * `cis_enabled` - Whether enable worker node to support cis security reinforcement, its valid value `true` or `false`. Default to `false` and apply to AliyunLinux series. Use `security_hardening_os` instead.
  * `compensate_with_on_demand` - Specifies whether to automatically create pay-as-you-go instances to meet the required number of ECS instances if preemptible instances cannot be created due to reasons such as cost or insufficient inventory. This parameter takes effect when you set `multi_az_policy` to `COST_OPTIMIZED`. Valid values: `true`: automatically creates pay-as-you-go instances to meet the required number of ECS instances if preemptible instances cannot be created. `false`: does not create pay-as-you-go instances to meet the required number of ECS instances if preemptible instances cannot be created.
  * `cpu_policy` - Node CPU management policies. Default value: `none`. When the cluster version is 1.12.6 or later, the following two policies are supported:- `static`: allows pods with certain resource characteristics on the node to enhance its CPU affinity and exclusivity.- `none`: Enables the existing default CPU affinity scheme.
  * `data_disks` - Configure the data disk of the node in the node pool.
    * `auto_format` - Whether to automatically mount the data disk. Valid values: true and false.
    * `auto_snapshot_policy_id` - The ID of the automatic snapshot policy that you want to apply to the system disk.
    * `bursting_enabled` - Whether the data disk is enabled with Burst (performance Burst). This is configured when the disk type is cloud_auto.
    * `category` - The type of data disk. Default value: `cloud_efficiency`. Valid values:- `cloud`: basic disk.- `cloud_efficiency`: ultra disk.- `cloud_ssd`: standard SSD.- `cloud_essd`: Enterprise SSD (ESSD).- `cloud_auto`: ESSD AutoPL disk.- `cloud_essd_entry`: ESSD Entry disk.- `elastic_ephemeral_disk_premium`: premium elastic ephemeral disk.- `elastic_ephemeral_disk_standard`: standard elastic ephemeral disk.
    * `device` - The mount target of data disk N. Valid values of N: 1 to 16. If you do not specify this parameter, the system automatically assigns a mount target when Auto Scaling creates ECS instances. The name of the mount target ranges from /dev/xvdb to /dev/xvdz.
    * `encrypted` - Specifies whether to encrypt data disks. Valid values: true and false. Default to `false`.
    * `file_system` - The type of the mounted file system. Works when auto_format is true. Optional value: `ext4`, `xfs`.
    * `kms_key_id` - The kms key id used to encrypt the data disk. It takes effect when `encrypted` is true.
    * `mount_target` - The Mount path. Works when auto_format is true.
    * `name` - The length is 2~128 English or Chinese characters. It must start with an uppercase or lowr letter or a Chinese character and cannot start with http:// or https. Can contain numbers, colons (:), underscores (_), or dashes (-). It will be overwritten if auto_format is set.
    * `performance_level` - Worker node data disk performance level, when `category` values `cloud_essd`, the optional values are `PL0`, `PL1`, `PL2` or `PL3`, but the specific performance level is related to the disk capacity. For more information, see [Enhanced SSDs](https://www.alibabacloud.com/help/doc-detail/122389.htm). Default is `PL1`.
    * `provisioned_iops` - The read/write IOPS preconfigured for the data disk, which is configured when the disk type is cloud_auto.
    * `size` - The size of a data disk, Its valid value range [40~32768] in GB. Default to `40`.
    * `snapshot_id` - The ID of the snapshot that you want to use to create data disk N. Valid values of N: 1 to 16. If you specify this parameter, DataDisk.N.Size is ignored. The size of the disk is the same as the size of the specified snapshot. If you specify a snapshot that is created on or before July 15, 2013, the operation fails and InvalidSnapshot.TooOld is returned.
  * `deployment_set_id` - The deployment set of node pool. Specify the deploymentSet to ensure that the nodes in the node pool can be distributed on different physical machines.
  * `desired_size` - Number of expected nodes in the node pool.
  * `image_id` - The custom image ID. The system-provided image is used by default.
  * `image_type` - The operating system image type and the `platform` parameter can be selected from the following values:- `AliyunLinux` : Alinux2 image.- `AliyunLinux3` : Alinux3 image.- `AliyunLinux3Arm64` : Alinux3 mirror ARM version.- `AliyunLinuxUEFI` : Alinux2 Image UEFI version.- `CentOS` : CentOS image.- `Windows` : Windows image.- `WindowsCore` : WindowsCore image.- `ContainerOS` : container-optimized image.- `Ubuntu`: Ubuntu image.
  * `install_cloud_monitor` - Whether to install cloud monitoring on the ECS node. After installation, you can view the monitoring information of the created ECS instance in the cloud monitoring console and recommend enable it. Default value: `false`. Valid values:- `true` : install cloud monitoring on the ECS node.- `false` : does not install cloud monitoring on the ECS node.
  * `instance_charge_type` - Node payment type. Valid values: `PostPaid`, `PrePaid`, default is `PostPaid`. If value is `PrePaid`, the arguments `period`, `period_unit`, `auto_renew` and `auto_renew_period` are required.
  * `instance_types` - In the node instance specification list, you can select multiple instance specifications as alternatives. When each node is created, it will try to purchase from the first specification until it is created successfully. The final purchased instance specifications may vary with inventory changes.
  * `internet_charge_type` - The billing method for network usage. Valid values `PayByBandwidth` and `PayByTraffic`. Conflict with `eip_internet_charge_type`, EIP and public network IP can only choose one. 
  * `internet_max_bandwidth_out` - The maximum bandwidth of the public IP address of the node. The unit is Mbps(Mega bit per second). The value range is:\[1,100\]
  * `key_name` - The name of the key pair. When the node pool is a managed node pool, only `key_name` is supported.
  * `kubelet_configuration` - Kubelet configuration parameters for worker nodes. See [`kubelet_configuration`](#kubelet_configuration) below. More information in [Kubelet Configuration](https://kubernetes.io/docs/reference/config-api/kubelet-config.v1beta1/).
    * `allowed_unsafe_sysctls` - Allowed sysctl mode whitelist.
    * `cluster_dns` - The list of IP addresses of the cluster DNS servers.
    * `container_log_max_files` - The maximum number of log files that can exist in each container.
    * `container_log_max_size` - The maximum size that can be reached before a log file is rotated.
    * `container_log_max_workers` - Specifies the maximum number of concurrent workers required to perform log rotation operations.
    * `container_log_monitor_interval` - Specifies the duration for which container logs are monitored for log rotation.
    * `cpu_cfs_quota` - CPU CFS quota constraint switch.
    * `cpu_cfs_quota_period` - CPU CFS quota period value.
    * `cpu_manager_policy` - Same as cpuManagerPolicy. The name of the policy to use. Requires the CPUManager feature gate to be enabled. Valid value is `none` or `static`.
    * `event_burst` - Same as eventBurst. The maximum size of a burst of event creations, temporarily allows event creations to burst to this number, while still not exceeding `event_record_qps`. It is only used when `event_record_qps` is greater than 0. Valid value is `[0-100]`.
    * `event_record_qps` - Same as eventRecordQPS. The maximum event creations per second. If 0, there is no limit enforced. Valid value is `[0-50]`.
    * `eviction_hard` - Same as evictionHard. The map of signal names to quantities that defines hard eviction thresholds. For example: `{"memory.available" = "300Mi"}`.
    * `eviction_soft` - Same as evictionSoft. The map of signal names to quantities that defines soft eviction thresholds. For example: `{"memory.available" = "300Mi"}`.
    * `eviction_soft_grace_period` - Same as evictionSoftGracePeriod. The map of signal names to quantities that defines grace periods for each soft eviction signal. For example: `{"memory.available" = "30s"}`.
    * `feature_gates` - Feature switch to enable configuration of experimental features.
    * `image_gc_high_threshold_percent` - If the image usage exceeds this threshold, image garbage collection will continue.
    * `image_gc_low_threshold_percent` - Image garbage collection is not performed when the image usage is below this threshold.
    * `kube_api_burst` - Same as kubeAPIBurst. The burst to allow while talking with kubernetes api-server. Valid value is `[0-100]`.
    * `kube_api_qps` - Same as kubeAPIQPS. The QPS to use while talking with kubernetes api-server. Valid value is `[0-50]`.
    * `kube_reserved` - Same as kubeReserved. The set of ResourceName=ResourceQuantity (e.g. cpu=200m,memory=150G) pairs that describe resources reserved for kubernetes system components. Currently, cpu, memory and local storage for root file system are supported. See [compute resources](https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/) for more details.
    * `max_pods` - The maximum number of running pods.
    * `memory_manager_policy` - The policy to be used by the memory manager.
    * `pod_pids_limit` - The maximum number of PIDs that can be used in a Pod.
    * `read_only_port` - Read-only port number.
    * `registry_burst` - Same as registryBurst. The maximum size of burst pulls, temporarily allows pulls to burst to this number, while still not exceeding `registry_pull_qps`. Only used if `registry_pull_qps` is greater than 0. Valid value is `[0-100]`.
    * `registry_pull_qps` - Same as registryPullQPS. The limit of registry pulls per second. Setting it to `0` means no limit. Valid value is `[0-50]`.
    * `reserved_memory` - Reserve memory for NUMA nodes.
      * `limits` - Memory resource limit.
      * `numa_node` - The NUMA node.
    * `serialize_image_pulls` - Same as serializeImagePulls. When enabled, it tells the Kubelet to pull images one at a time. We recommend not changing the default value on nodes that run docker daemon with version < 1.9 or an Aufs storage backend. Valid value is `true` or `false`.
    * `system_reserved` - Same as systemReserved. The set of ResourceName=ResourceQuantity (e.g. cpu=200m,memory=150G) pairs that describe resources reserved for non-kubernetes components. Currently, only cpu and memory are supported. See [compute resources](https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/) for more details.
    * `topology_manager_policy` - Name of the Topology Manager policy used.
    * `tracing` - OpenTelemetry tracks the configuration information for client settings versioning.
      * `endpoint` - The endpoint of the collector.
      * `sampling_rate_per_million` - Number of samples to be collected per million span.
  * `labels` - A List of Kubernetes labels to assign to the nodes . Only labels that are applied with the ACK API are managed by this argument. Detailed below. More information in [Labels](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/).
    * `key` - The label key.
    * `value` - The label value.
  * `login_as_non_root` - Whether the ECS instance is logged on as a ecs-user user. Valid value: `true` and `false`.
  * `management` - Managed node pool configuration.
    * `auto_repair` - Whether to enable automatic repair. Valid values: `true`: Automatic repair. `false`: not automatically repaired.
    * `auto_repair_policy` - Automatic repair node policy.
      * `restart_node` - Whether to allow node restart.
    * `auto_upgrade` - Specifies whether to enable auto update. Valid values: `true`: enables auto update. `false`: disables auto update.
    * `auto_upgrade_policy` - The auto update policy.
      * `auto_upgrade_kubelet` - Specifies whether  to automatically update the kubelet. Valid values: `true`: yes; `false`: no.
    * `auto_vul_fix` - Specifies whether to automatically patch CVE vulnerabilities. Valid values: `true`, `false`.
    * `auto_vul_fix_policy` - The auto CVE patching policy.
      * `restart_node` - Specifies whether to automatically restart nodes after patching CVE vulnerabilities. Valid values: `true`, `false`.
      * `vul_level` - The severity levels of vulnerabilities that is allowed to automatically patch. Multiple severity levels are separated by commas (,).
    * `enable` - Specifies whether to enable the managed node pool feature. Valid values: `true`: enables the managed node pool feature. `false`: disables the managed node pool feature. Other parameters in this section take effect only when you specify enable=true.
    * `max_unavailable` - Maximum number of unavailable nodes. Default value: 1. Value range:\[1,1000\].
    * `surge` - Number of additional nodes. You have to specify one of surge, surge_percentage.
    * `surge_percentage` - Proportion of additional nodes. You have to specify one of surge, surge_percentage.
  * `multi_az_policy` - The scaling policy for ECS instances in a multi-zone scaling group. Valid value: `PRIORITY`, `COST_OPTIMIZED` and `BALANCE`. `PRIORITY`: scales the capacity according to the virtual switches you define (VSwitchIds.N). When an ECS instance cannot be created in the zone where the higher-priority vSwitch is located, the next-priority vSwitch is automatically used to create an ECS instance. `COST_OPTIMIZED`: try to create by vCPU unit price from low to high. When the scaling configuration is configured with multiple instances of preemptible billing, preemptible instances are created first. You can continue to use the `CompensateWithOnDemand` parameter to specify whether to automatically try to create a preemptible instance by paying for it. It takes effect only when the scaling configuration has multi-instance specifications or preemptible instances. `BALANCE`: distributes ECS instances evenly among the multi-zone specified by the scaling group. If the zones become unbalanced due to insufficient inventory, you can use the API [RebalanceInstances](~~ 71516 ~~) to balance resources.
  * `node_name_mode` - Each node name consists of a prefix, its private network IP, and a suffix, separated by commas. The input format is `customized,,ip,`.- The prefix and suffix can be composed of one or more parts separated by '.', each part can use lowercase letters, numbers and '-', and the beginning and end of the node name must be lowercase letters and numbers.- The node IP address is the complete private IP address of the node.- For example, if the string `customized,aliyun,ip,com` is passed in (where 'customized' and 'ip' are fixed strings, 'aliyun' is the prefix, and 'com' is the suffix), the name of the node is `aliyun192.168.xxx.xxxcom`.
  * `node_pool_id` - The first ID of the resource.
  * `node_pool_name` - The name of node pool.
  * `on_demand_base_capacity` - The minimum number of pay-as-you-go instances that must be kept in the scaling group. Valid values: 0 to 1000. If the number of pay-as-you-go instances is less than the value of this parameter, Auto Scaling preferably creates pay-as-you-go instances.
  * `on_demand_percentage_above_base_capacity` - The percentage of pay-as-you-go instances among the extra instances that exceed the number specified by `on_demand_base_capacity`. Valid values: 0 to 100.
  * `password` - The password of ssh login. You have to specify one of `password` and `key_name` fields. The password rule is 8 to 30 characters and contains at least three items (upper and lower case letters, numbers, and special symbols).
  * `period` - Node payment period. Its valid value is one of {1, 2, 3, 6, 12}.
  * `period_unit` - Node payment period unit, valid value: `Month`. Default is `Month`.
  * `platform` - Operating system release, using `image_type` instead.
  * `pre_user_data` - Node pre custom data, base64-encoded, the script executed before the node is initialized. 
  * `private_pool_options` - Private node pool configuration.
    * `private_pool_options_id` - The ID of the private node pool.
    * `private_pool_options_match_criteria` - The type of private node pool. This parameter specifies the type of the private pool that you want to use to create instances. A private node pool is generated when an elasticity assurance or a capacity reservation service takes effect. The system selects a private node pool to launch instances. Valid values: `Open`: specifies an open private node pool. The system selects an open private node pool to launch instances. If no matching open private node pool is available, the resources in the public node pool are used. `Target`: specifies a private node pool. The system uses the resources of the specified private node pool to launch instances. If the specified private node pool is unavailable, instances cannot be started. `None`: no private node pool is used. The resources of private node pools are not used to launch the instances.
  * `ram_role_name` - The name of the Worker RAM role.* If it is empty, the default Worker RAM role created in the cluster will be used.* If the specified RAM role is not empty, the specified RAM role must be a **Common Service role**, and its **trusted service** configuration must be **cloud server**. For more information, see [Create a common service role](https://help.aliyun.com/document_detail/116800.html). If the specified RAM role is not the default Worker RAM role created in the cluster, the role name cannot start with 'KubernetesMasterRole-'or 'KubernetesWorkerRole.-> **NOTE:**  This parameter is only supported for ACK-managed clusters of 1.22 or later versions.
  * `rds_instances` - The list of RDS instances.
  * `resource_group_id` - The ID of the resource group
  * `runtime_name` - The runtime name of containers. If not set, the cluster runtime will be used as the node pool runtime. If you select another container runtime, see [Comparison of Docker, containerd, and Sandboxed-Container](https://www.alibabacloud.com/help/doc-detail/160313.htm).
  * `runtime_version` - The runtime version of containers. If not set, the cluster runtime will be used as the node pool runtime.
  * `scaling_config` - Automatic scaling configuration.
    * `eip_bandwidth` - Peak EIP bandwidth. Its valid value range [1~500] in Mbps. It works if `is_bond_eip=true`. Default to `5`.
    * `eip_internet_charge_type` - EIP billing type. `PayByBandwidth`: Charged at fixed bandwidth. `PayByTraffic`: Billed as used traffic. Default: `PayByBandwidth`. It works if `is_bond_eip=true`, conflict with `internet_charge_type`. EIP and public network IP can only choose one.
    * `enable` - Whether to enable automatic scaling. Value:- `true`: enables the node pool auto-scaling function.- `false`: Auto scaling is not enabled. When the value is false, other `auto_scaling` configuration parameters do not take effect.
    * `is_bond_eip` - Whether to bind EIP for an instance. Default: `false`.
    * `max_size` - Max number of instances in a auto scaling group, its valid value range [0~1000]. `max_size` has to be greater than `min_size`.
    * `min_size` - Min number of instances in a auto scaling group, its valid value range [0~1000].
    * `type` - Instance classification, not required. Vaild value: `cpu`, `gpu`, `gpushare` and `spot`. Default: `cpu`. The actual instance type is determined by `instance_types`.
  * `scaling_group_id` - The ID of the scaling group.
  * `scaling_policy` - Scaling group mode, default value: `release`. Valid values:- `release`: in the standard mode, scaling is performed by creating and releasing ECS instances based on the usage of the application resource value.- `recycle`: in the speed mode, scaling is performed through creation, shutdown, and startup to increase the speed of scaling again (computing resources are not charged during shutdown, only storage fees are charged, except for local disk models).
  * `security_group_id` - The security group ID of the node pool. This field has been replaced by `security_group_ids`, please use the `security_group_ids` field instead.
  * `security_group_ids` - Multiple security groups can be configured for a node pool. If both `security_group_ids` and `security_group_id` are configured, `security_group_ids` takes effect. This field cannot be modified.
  * `security_hardening_os` - Alibaba Cloud OS security reinforcement. Default value: `false`. Value:-`true`: enable Alibaba Cloud OS security reinforcement.-`false`: does not enable Alibaba Cloud OS security reinforcement.
  * `soc_enabled` - Whether enable worker node to support soc security reinforcement, its valid value `true` or `false`. Default to `false` and apply to AliyunLinux series. See [SOC Reinforcement](https://help.aliyun.com/document_detail/196148.html).> It is forbidden to set both `security_hardening_os` and `soc_enabled` to `true` at the same time.
  * `spot_instance_pools` - The number of instance types that are available. Auto Scaling creates preemptible instances of multiple instance types that are available at the lowest cost. Valid values: 1 to 10.
  * `spot_instance_remedy` - Specifies whether to supplement preemptible instances when the number of preemptible instances drops below the specified minimum number. If you set the value to true, Auto Scaling attempts to create a new preemptible instance when the system notifies that an existing preemptible instance is about to be reclaimed. Valid values: `true`: enables the supplementation of preemptible instances. `false`: disables the supplementation of preemptible instances.
  * `spot_price_limit` - The current single preemptible instance type market price range configuration.
    * `instance_type` - The type of the preemptible instance.
    * `price_limit` - The maximum price of a single instance.
  * `spot_strategy` - The preemptible instance type. Value:- `NoSpot` : Non-preemptible instance.- `SpotWithPriceLimit` : Set the upper limit of the preemptible instance price.- `SpotAsPriceGo` : The system automatically bids, following the actual price of the current market.
  * `system_disk_bursting_enabled` - Specifies whether to enable the burst feature for system disks. Valid values:`true`: enables the burst feature. `false`: disables the burst feature. This parameter is supported only when `system_disk_category` is set to `cloud_auto`.
  * `system_disk_categories` - The multi-disk categories of the system disk. When a high-priority disk type cannot be used, Auto Scaling automatically tries to create a system disk with the next priority disk category. Valid values see `system_disk_category`.
  * `system_disk_category` - The category of the system disk for nodes. Default value: `cloud_efficiency`. Valid values:- `cloud`: basic disk.- `cloud_efficiency`: ultra disk.- `cloud_ssd`: standard SSD.- `cloud_essd`: ESSD.- `cloud_auto`: ESSD AutoPL disk.- `cloud_essd_entry`: ESSD Entry disk.
  * `system_disk_encrypt_algorithm` - The encryption algorithm used by the system disk. Value range: aes-256.
  * `system_disk_encrypted` - Whether to encrypt the system disk. Value range: `true`: encryption. `false`: Do not encrypt.
  * `system_disk_kms_key` - The ID of the KMS key used by the system disk.
  * `system_disk_performance_level` - The system disk performance of the node takes effect only for the ESSD disk.- `PL0`: maximum random read/write IOPS 10000 for a single disk.- `PL1`: maximum random read/write IOPS 50000 for a single disk.- `PL2`: highest random read/write IOPS 100000 for a single disk.- `PL3`: maximum random read/write IOPS 1 million for a single disk.
  * `system_disk_provisioned_iops` - The predefined IOPS of a system disk. Valid values: 0 to min{50,000, 1,000 × Capacity - Baseline IOPS}. Baseline IOPS = min{1,800 + 50 × Capacity, 50,000}. This parameter is supported only when `system_disk_category` is set to `cloud_auto`.
  * `system_disk_size` - The size of the system disk. Unit: GiB. The value of this parameter must be at least 1 and greater than or equal to the image size. Default value: 40 or the size of the image, whichever is larger.- Basic disk: 20 to 500.- ESSD (cloud_essd): The valid values vary based on the performance level of the ESSD. PL0 ESSD: 1 to 2048. PL1 ESSD: 20 to 2048. PL2 ESSD: 461 to 2048. PL3 ESSD: 1261 to 2048.- ESSD AutoPL disk (cloud_auto): 1 to 2048.- Other disk categories: 20 to 2048.
  * `system_disk_snapshot_policy_id` - The ID of the automatic snapshot policy used by the system disk.
  * `tags` - Add tags only for ECS instances. The maximum length of the tag key is 128 characters. The tag key and value cannot start with aliyun or acs:, or contain https:// or http://.
  * `taints` - A List of Kubernetes taints to assign to the nodes. Detailed below. More information in [Taints and Toleration](https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration/).
    * `effect` - The scheduling policy.
    * `key` - The key of a taint.
    * `value` - The value of a taint.
  * `tee_config` - The configuration about confidential computing for the cluster.
    * `tee_enable` - Specifies whether to enable confidential computing for the cluster.
  * `unschedulable` - Whether the node after expansion can be scheduled.
  * `user_data` - Node custom data, base64-encoded.
  * `vswitch_ids` - The vswitches used by node pool workers.
