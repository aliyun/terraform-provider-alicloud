---
subcategory: "Container Service for Kubernetes (ACK)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cs_kubernetes_node_pool"
description: |-
  Provides a Alicloud Container Service for Kubernetes (ACK) Nodepool resource.
---

# alicloud_cs_kubernetes_node_pool

Provides a Container Service for Kubernetes (ACK) Nodepool resource.

This resource will help you to manage node pool in Kubernetes Cluster, see [What is kubernetes node pool](https://www.alibabacloud.com/help/en/ack/ack-managed-and-ack-dedicated/developer-reference/api-create-node-pools). 

-> **NOTE:** Available since v1.97.0.

-> **NOTE:** From version 1.109.1, support managed node pools, but only for the professional managed clusters.

-> **NOTE:** From version 1.109.1, support remove node pool nodes.

-> **NOTE:** From version 1.111.0, support auto scaling node pool. For more information on how to use auto scaling node pools, see [Use Terraform to create an elastic node pool](https://www.alibabacloud.com/help/en/ack/ack-managed-and-ack-dedicated/developer-reference/api-create-node-pools). With auto-scaling is enabled, the nodes in the node pool will be labeled with `k8s.aliyun.com=true` to prevent system pods such as coredns, metrics-servers from being scheduled to elastic nodes, and to prevent node shrinkage from causing business abnormalities.

-> **NOTE:** ACK adds a new RamRole (AliyunCSManagedAutoScalerRole) for the permission control of the node pool with auto-scaling enabled. If you are using a node pool with auto scaling, please click [AliyunCSManagedAutoScalerRole](https://ram.console.aliyun.com/role/authorization?request=%7B%22Services%22%3A%5B%7B%22Service%22%3A%22CS%22%2C%22Roles%22%3A%5B%7B%22RoleName%22%3A%22AliyunCSManagedAutoScalerRole%22%2C%22TemplateId%22%3A%22AliyunCSManagedAutoScalerRole%22%7D%5D%7D%5D%2C%22ReturnUrl%22%3A%22https%3A%2F%2Fcs.console.aliyun.com%2F%22%7D) to complete the authorization. 

-> **NOTE:** ACK adds a new RamRole（AliyunCSManagedNlcRole） for the permission control of the management node pool. If you use the management node pool, please click [AliyunCSManagedNlcRole](https://ram.console.aliyun.com/role/authorization?spm=5176.2020520152.0.0.387f16ddEOZxMv&request=%7B%22Services%22%3A%5B%7B%22Service%22%3A%22CS%22%2C%22Roles%22%3A%5B%7B%22RoleName%22%3A%22AliyunCSManagedNlcRole%22%2C%22TemplateId%22%3A%22AliyunCSManagedNlcRole%22%7D%5D%7D%5D%2C%22ReturnUrl%22%3A%22https%3A%2F%2Fcs.console.aliyun.com%2F%22%7D) to complete the authorization.

-> **NOTE:** From version 1.123.1, supports the creation of a node pool of spot instance.

-> **NOTE:** It is recommended to create a cluster with zero worker nodes, and then use a node pool to manage the cluster nodes. 

-> **NOTE:** From version 1.127.0, support for adding existing nodes to the node pool. In order to distinguish automatically created nodes, it is recommended that existing nodes be placed separately in a node pool for management. 

-> **NOTE:** From version 1.149.0, support for specifying deploymentSet for node pools. 

-> **NOTE:** From version 1.158.0, Support for specifying the desired size of nodes for the node pool, for more information, visit [Modify the expected number of nodes in a node pool](https://www.alibabacloud.com/help/en/doc-detail/160490.html#title-mpp-3jj-oo3)

-> **NOTE:** From version 1.166.0, Support configuring system disk encryption.

-> **NOTE:** From version 1.177.0+, Support `kms_encryption_context`, `rds_instances`, `system_disk_snapshot_policy_id` and `cpu_policy`, add spot strategy `SpotAsPriceGo` and `NoSpot`.

-> **NOTE:** From version 1.180.0+, Support worker nodes customized kubelet parameters by field `kubelet_configuration` and `rollout_policy`.

-> **NOTE:** From version 1.185.0+, Field `rollout_policy` will be deprecated and please use field `rolling_policy` instead.

For information about Container Service for Kubernetes (ACK) Nodepool and how to use it, see [What is Nodepool](https://www.alibabacloud.com/help/en/ack/ack-managed-and-ack-dedicated/developer-reference/api-create-node-pools).

-> **NOTE:** Available since v1.97.0.

## Example Usage

Basic Usage

```terraform
resource "random_integer" "default" {
  max = 99999
  min = 10000
}

variable "name" {
  default = "terraform-example"
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
  name_prefix          = "terraform-example-${random_integer.default.result}"
  cluster_spec         = "ack.pro.small"
  worker_vswitch_ids   = [alicloud_vswitch.default.id]
  new_nat_gateway      = true
  pod_cidr             = cidrsubnet("10.0.0.0/8", 8, 36)
  service_cidr         = cidrsubnet("172.16.0.0/16", 4, 7)
  slb_internet_enabled = true
  enable_rrsa          = true
}

resource "alicloud_key_pair" "default" {
  key_pair_name = "terraform-example-${random_integer.default.result}"
}

resource "alicloud_cs_kubernetes_node_pool" "default" {
  node_pool_name       = var.name
  cluster_id           = alicloud_cs_managed_kubernetes.default.id
  vswitch_ids          = [alicloud_vswitch.default.id]
  instance_types       = [data.alicloud_instance_types.cloud_efficiency.instance_types.0.id]
  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  key_name             = alicloud_key_pair.default.key_pair_name
  // define with multi-labels by defining with labels blocks
  labels {
    key   = "test1"
    value = "nodepool"
  }
  labels {
    key   = "test2"
    value = "nodepool"
  }
  // define with multi-taints by defining with taints blocks
  taints {
    key    = "tf"
    effect = "NoSchedule"
    value  = "example"
  }
  taints {
    key    = "tf2"
    effect = "NoSchedule"
    value  = "example2"
  }
}

#The parameter `node_count` is deprecated from version 1.158.0. Please use the new parameter `desired_size` instead, you can update it as follows.
resource "alicloud_cs_kubernetes_node_pool" "desired_size" {
  node_pool_name       = "desired_size"
  cluster_id           = alicloud_cs_managed_kubernetes.default.id
  vswitch_ids          = [alicloud_vswitch.default.id]
  instance_types       = [data.alicloud_instance_types.cloud_efficiency.instance_types.0.id]
  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  key_name             = alicloud_key_pair.default.key_pair_name
  desired_size         = 0
}

# Create a managed node pool. If you need to enable maintenance window, you need to set the maintenance window in `alicloud_cs_managed_kubernetes`.
resource "alicloud_cs_kubernetes_node_pool" "maintenance" {
  node_pool_name       = "maintenance"
  cluster_id           = alicloud_cs_managed_kubernetes.default.id
  vswitch_ids          = [alicloud_vswitch.default.id]
  instance_types       = [data.alicloud_instance_types.cloud_efficiency.instance_types.0.id]
  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40

  # only key_name is supported in the management node pool
  key_name = alicloud_key_pair.default.key_pair_name

  # you need to specify the number of nodes in the node pool, which can be zero
  desired_size = 1

  # management node pool configuration.
  management {
    enable      = true
    auto_repair = true
    auto_repair_policy {
      restart_node = true
    }
    auto_upgrade = true
    auto_upgrade_policy {
      auto_upgrade_kubelet = true
    }
    auto_vul_fix = true
    auto_vul_fix_policy {
      vul_level    = "asap"
      restart_node = true
    }
    max_unavailable = 1
  }

  # Enable with automatic scaling node pool configuration.
  # With auto-scaling is enabled, the nodes in the node pool will be labeled with `k8s.aliyun.com=true` to prevent system pods such as coredns, metrics-servers from being scheduled to elastic nodes, and to prevent node shrinkage from causing business abnormalities.
  #  scaling_config {
  #    min_size = 1
  #    max_size = 10
  #    type     = "cpu"
  #  }
}

#Create a node pool with spot instance.
resource "alicloud_cs_kubernetes_node_pool" "spot_instance" {
  node_pool_name       = "spot_instance"
  cluster_id           = alicloud_cs_managed_kubernetes.default.id
  vswitch_ids          = [alicloud_vswitch.default.id]
  instance_types       = [data.alicloud_instance_types.cloud_efficiency.instance_types.0.id, data.alicloud_instance_types.cloud_efficiency.instance_types.1.id]
  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  key_name             = alicloud_key_pair.default.key_pair_name

  # you need to specify the number of nodes in the node pool, which can be 0
  desired_size = 1

  # spot config
  spot_strategy = "SpotWithPriceLimit"
  spot_price_limit {
    instance_type = data.alicloud_instance_types.cloud_efficiency.instance_types.0.id
    # Different instance types have different price caps
    price_limit = "0.70"
  }
  // define with multi-spot_price_limit by defining with spot_price_limit blocks
  spot_price_limit {
    instance_type = data.alicloud_instance_types.cloud_efficiency.instance_types.1.id
    price_limit   = "0.72"
  }
}


#Use Spot instances to create a node pool with auto-scaling enabled
resource "alicloud_cs_kubernetes_node_pool" "spot_auto_scaling" {
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

#Create a `PrePaid` node pool.
resource "alicloud_cs_kubernetes_node_pool" "prepaid_node" {
  node_pool_name       = "prepaid_node"
  cluster_id           = alicloud_cs_managed_kubernetes.default.id
  vswitch_ids          = [alicloud_vswitch.default.id]
  instance_types       = [data.alicloud_instance_types.cloud_efficiency.instance_types.0.id]
  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  key_name             = alicloud_key_pair.default.key_pair_name
  # use PrePaid
  instance_charge_type = "PrePaid"
  period               = 1
  period_unit          = "Month"
  auto_renew           = true
  auto_renew_period    = 1

  # open cloud monitor
  install_cloud_monitor = true
}

##Create a node pool with customized kubelet parameters
resource "alicloud_cs_kubernetes_node_pool" "customized_kubelet" {
  node_pool_name       = "customized_kubelet"
  cluster_id           = alicloud_cs_managed_kubernetes.default.id
  vswitch_ids          = [alicloud_vswitch.default.id]
  instance_types       = [data.alicloud_instance_types.cloud_efficiency.instance_types.0.id]
  system_disk_category = "cloud_efficiency"
  system_disk_size     = 40
  instance_charge_type = "PostPaid"
  desired_size         = 0

  # kubelet configuration parameters
  kubelet_configuration {
    registry_pull_qps     = 10
    registry_burst        = 5
    event_record_qps      = 10
    event_burst           = 5
    serialize_image_pulls = true
    eviction_hard = {
      "memory.available"  = "1024Mi"
      "nodefs.available"  = "10%"
      "nodefs.inodesFree" = "5%"
      "imagefs.available" = "10%"
    }
    system_reserved = {
      "cpu"               = "1"
      "memory"            = "1Gi"
      "ephemeral-storage" = "10Gi"
    }
    kube_reserved = {
      "cpu"    = "500m"
      "memory" = "1Gi"
    }
    container_log_max_size  = "200Mi"
    container_log_max_files = 3
    max_pods                = 100
    read_only_port          = 0
    allowed_unsafe_sysctls  = ["net.ipv4.route.min_pmtu"]
  }

  # rolling policy: works when updating
  rolling_policy {
    max_parallelism = 1
  }
}
```

ACK Auto Mode NodePool:

ACK nodepool with Auto Mode
Nodepools enable Auto Mode can only be created in Auto Mode clusters. An Auto Mode cluster automatically creates a Auto Mode node pool when creating cluster, which you can import using terraform import. You can also create a new Auto Mode node pool using the following code.

```terraform
provider "alicloud" {
  region = var.region_id
}

variable "region_id" {
  type    = string
  default = "cn-hangzhou"
}

variable "cluster_spec" {
  type        = string
  description = "The cluster specifications of kubernetes cluster,which can be empty. Valid values:ack.standard : Standard managed clusters; ack.pro.small : Professional managed clusters."
  default     = "ack.pro.small"
}

variable "availability_zone" {
  description = "The availability zones of vswitches."
  default     = ["cn-hangzhou-i", "cn-hangzhou-j", "cn-hangzhou-k"]
}

variable "node_vswitch_ids" {
  description = "List of existing node vswitch ids for terway."
  type        = list(string)
  default     = []
}

variable "node_vswitch_cidrs" {
  description = "List of cidr blocks used to create several new vswitches when 'node_vswitch_ids' is not specified."
  type        = list(string)
  default     = ["172.16.0.0/23", "172.16.2.0/23", "172.16.4.0/23"]
}

variable "terway_vswitch_ids" {
  description = "List of existing pod vswitch ids for terway."
  type        = list(string)
  default     = []
}

variable "terway_vswitch_cidrs" {
  description = "List of cidr blocks used to create several new vswitches when 'terway_vswitch_ids' is not specified."
  type        = list(string)
  default     = ["172.16.208.0/20", "172.16.224.0/20", "172.16.240.0/20"]
}

variable "cluster_addons" {
  type = list(object({
    name     = string
    config   = optional(string)
    disabled = optional(bool, false)
  }))
  default = [
    { name = "metrics-server" },
    { name = "managed-coredns" },
    { name = "managed-security-inspector" },
    { name = "ack-cost-exporter" },
    {
      name   = "terway-controlplane"
      config = "{\"ENITrunking\":\"true\"}"
    },
    {
      name   = "terway-eniip"
      config = "{\"NetworkPolicy\":\"false\",\"ENITrunking\":\"true\",\"IPVlan\":\"false\"}"
    },
    { name = "csi-plugin" },
    { name = "managed-csiprovisioner" },
    {
      name   = "storage-operator"
      config = "{\"CnfsOssEnable\":\"false\",\"CnfsNasEnable\":\"false\"}"
    },
    {
      name   = "loongcollector"
      config = "{\"IngressDashboardEnabled\":\"true\"}"
    },
    {
      name   = "ack-node-problem-detector"
      config = "{\"sls_project_name\":\"\"}"
    },
    {
      name     = "nginx-ingress-controller"
      disabled = true
    },
    {
      name   = "alb-ingress-controller"
      config = "{\"albIngress\":{\"CreateDefaultALBConfig\":false}}"
    },
    {
      name   = "arms-prometheus"
      config = "{\"prometheusMode\":\"default\"}"
    },
    { name = "alicloud-monitor-controller" },
    { name = "managed-aliyun-acr-credential-helper" },
  ]
}

variable "k8s_name_prefix" {
  description = "The name prefix used to create managed kubernetes cluster."
  default     = "tf-ack-hangzhou"
}

locals {
  k8s_name_terway = substr(join("-", [var.k8s_name_prefix, "terway"]), 0, 63)
  new_vpc_name    = "tf-vpc-172-16"
  nodepool_name   = "default-nodepool"
}

resource "alicloud_vpc" "default" {
  vpc_name   = local.new_vpc_name
  cidr_block = "172.16.0.0/12"
}

resource "alicloud_vswitch" "vswitches" {
  count      = length(var.node_vswitch_ids) > 0 ? 0 : length(var.node_vswitch_cidrs)
  vpc_id     = alicloud_vpc.default.id
  cidr_block = element(var.node_vswitch_cidrs, count.index)
  zone_id    = element(var.availability_zone, count.index)
}

resource "alicloud_vswitch" "terway_vswitches" {
  count      = length(var.terway_vswitch_ids) > 0 ? 0 : length(var.terway_vswitch_cidrs)
  vpc_id     = alicloud_vpc.default.id
  cidr_block = element(var.terway_vswitch_cidrs, count.index)
  zone_id    = element(var.availability_zone, count.index)
}

resource "alicloud_cs_managed_kubernetes" "default" {
  name                         = local.k8s_name_terway
  cluster_spec                 = var.cluster_spec
  vswitch_ids                  = split(",", join(",", alicloud_vswitch.vswitches.*.id))
  pod_vswitch_ids              = split(",", join(",", alicloud_vswitch.terway_vswitches.*.id))
  new_nat_gateway              = true
  service_cidr                 = "10.11.0.0/16"
  slb_internet_enabled         = true
  enable_rrsa                  = true
  control_plane_log_components = ["apiserver", "kcm", "scheduler", "ccm"]
  dynamic "addons" {
    for_each = var.cluster_addons
    content {
      name     = lookup(addons.value, "name", var.cluster_addons)
      config   = lookup(addons.value, "config", var.cluster_addons)
      disabled = lookup(addons.value, "disabled", var.cluster_addons)
    }
  }

  auto_mode {
    enabled = true
  }

  maintenance_window {
    duration         = "3h"
    weekly_period    = "Monday"
    enable           = true
    maintenance_time = "2025-07-07T00:00:00.000+08:00"
  }

  operation_policy {
    cluster_auto_upgrade {
      channel = "stable"
      enabled = true
    }
  }
}

resource "alicloud_cs_kubernetes_node_pool" "auto_mode_example" {
  node_pool_name = local.nodepool_name
  cluster_id     = alicloud_cs_managed_kubernetes.default.id
  vswitch_ids    = split(",", join(",", alicloud_vswitch.vswitches.*.id))

  auto_mode {
    enabled = true
  }

  # Configure modifiable parameters
  scaling_config {
    max_size = 50
    min_size = 0
  }

  # instance_types and instance_patterns are mutually exclusive - use only one of them
  instance_patterns {
    min_cpu_cores           = 4
    max_cpu_cores           = 16
    min_memory_size         = 8
    max_memory_size         = 32
    instance_family_level   = "EnterpriseLevel"
    instance_type_families  = ["ecs.u1", "ecs.g6", "ecs.c6", "ecs.r6", "ecs.g7", "ecs.c7", "ecs.r7", "ecs.g8i", "ecs.c8i", "ecs.r8i"]
    excluded_instance_types = ["ecs.c6t.*", "ecs.g6t.*", "ecs.t5.*", "ecs.t6.*", "ecs.vgn.*", "ecs.sgn.*"] # ACK not support instance families
    instance_categories     = ["General-purpose"]
    cpu_architectures       = ["X86"]
  }

  data_disks {
    size      = 120
    encrypted = "false"
    category  = "cloud_essd"
  }

  labels {
    key   = "example1"
    value = "nodepool"
  }
  labels {
    key   = "example2"
    value = "nodepool"
  }

  taints {
    key    = "tf"
    effect = "NoSchedule"
    value  = "example"
  }
  taints {
    key    = "tf2"
    effect = "NoSchedule"
    value  = "example2"
  }


  # Alternative: use instance_types instead of instance_patterns
  # instance_types = ["ecs.c6.large"]

  # Ignore service-side default values to prevent configuration drift.
  # In Auto Mode nodepools, the parameters in ignore_changes below are not supported and must not be specified. Please do not remove them from the ignore_changes list.
  lifecycle {
    ignore_changes = [
      management,
      install_cloud_monitor,
      cpu_policy,
      node_name_mode,
      runtime_name,
      runtime_version,
      unschedulable,
      user_data,
      pre_user_data,
      auto_renew,
      auto_renew_period,
      cis_enabled,
      compensate_with_on_demand,
      deployment_set_id,
      image_id,
      image_type,
      instance_charge_type,
      instance_metadata_options,
      internet_charge_type,
      internet_max_bandwidth_out,
      key_name,
      login_as_non_root,
      password,
      multi_az_policy,
      on_demand_base_capacity,
      on_demand_percentage_above_base_capacity,
      period,
      period_unit,
      platform,
      private_pool_options,
      ram_role_name,
      rds_instances,
      scaling_policy,
      security_group_id,
      security_group_ids,
      security_hardening_os,
      soc_enabled,
      spot_instance_pools,
      spot_instance_remedy,
      spot_price_limit,
      spot_strategy,
      system_disk_category,
      system_disk_categories,
      system_disk_size,
      system_disk_bursting_enabled,
      system_disk_performance_level,
      system_disk_encrypted,
      system_disk_kms_key,
      system_disk_snapshot_policy_id,
      system_disk_encrypt_algorithm,
      system_disk_provisioned_iops,
      tee_config,
    ]
  }
}
```

## Argument Reference

The following arguments are supported:
* `auto_mode` - (Optional, ForceNew, Computed, List, Available since v1.266.0) Whether to enable auto mode. When enabled, the system will automatically manage the node pool with optimized default configurations. **Note:** When `auto_mode` is enabled, many parameters will be automatically set to default values and cannot be modified. See `auto_mode.enable` below for details. See [`auto_mode`](#auto_mode) below.
* `auto_renew` - (Optional) Whether to enable automatic renewal for nodes in the node pool takes effect only when `instance_charge_type` is set to `PrePaid`. Default value: `false`. Valid values:
  - `true`: Automatic renewal. 
  - `false`: Do not renew automatically.
* `auto_renew_period` - (Optional, Int) The automatic renewal period of nodes in the node pool takes effect only when you select Prepaid and Automatic Renewal, and is a required value. When `PeriodUnit = Month`, the value range is {1, 2, 3, 6, 12}. Default value: 1.
* `cis_enabled` - (Optional, ForceNew, Deprecated since v1.223.1) Whether enable worker node to support cis security reinforcement, its valid value `true` or `false`. Default to `false` and apply to AliyunLinux series. Use `security_hardening_os` instead.
* `cluster_id` - (Required, ForceNew) The id of kubernetes cluster.
* `compensate_with_on_demand` - (Optional) Specifies whether to automatically create pay-as-you-go instances to meet the required number of ECS instances if preemptible instances cannot be created due to reasons such as cost or insufficient inventory. This parameter takes effect when you set `multi_az_policy` to `COST_OPTIMIZED`. Valid values: `true`: automatically creates pay-as-you-go instances to meet the required number of ECS instances if preemptible instances cannot be created. `false`: does not create pay-as-you-go instances to meet the required number of ECS instances if preemptible instances cannot be created.
* `cpu_policy` - (Optional, Computed) Node CPU management policies. Default value: `none`. When the cluster version is 1.12.6 or later, the following two policies are supported:
  - `static`: allows pods with certain resource characteristics on the node to enhance its CPU affinity and exclusivity.
  - `none`: Enables the existing default CPU affinity scheme.
* `data_disks` - (Optional, List) Configure the data disk of the node in the node pool. See [`data_disks`](#data_disks) below.
* `deployment_set_id` - (Optional, ForceNew) The deployment set of node pool. Specify the deploymentSet to ensure that the nodes in the node pool can be distributed on different physical machines.
* `desired_size` - (Optional) Number of expected nodes in the node pool.
* `eflo_node_group` - (Optional, List, Available since v1.252.0) Lingjun node pool configuration. See [`eflo_node_group`](#eflo_node_group) below.

-> **NOTE:** The parameter is immutable after resource creation. It only applies during resource creation and has no effect when modified post-creation.

* `force_delete` - (Optional) Whether to force deletion.

-> **NOTE:** This parameter only takes effect when deletion is triggered.

* `image_id` - (Optional, Computed) The custom image ID. The system-provided image is used by default.
* `image_type` - (Optional, Computed, Available since v1.236.0) The operating system image type and the `platform` parameter can be selected from the following values:
  - `AliyunLinux` : Alinux2 image.
  - `AliyunLinux3` : Alinux3 image.
  - `AliyunLinux3Arm64` : Alinux3 mirror ARM version.
  - `AliyunLinuxUEFI` : Alinux2 Image UEFI version.
  - `CentOS` : CentOS image.
  - `Windows` : Windows image.
  - `WindowsCore` : WindowsCore image.
  - `ContainerOS` : container-optimized image.
  - `Ubuntu`: Ubuntu image.
  - `AliyunLinux3ContainerOptimized`: Alinux3 container-optimized image.
  - `Custom`：Custom image.
  - `AliyunLinux4ContainerOptimized`：Alinux4 container-optimized image.
* `install_cloud_monitor` - (Optional) Whether to install cloud monitoring on the ECS node. After installation, you can view the monitoring information of the created ECS instance in the cloud monitoring console and recommend enable it. Default value: `false`. Valid values:
  - `true` : install cloud monitoring on the ECS node.
  - `false` : does not install cloud monitoring on the ECS node.
* `instance_charge_type` - (Optional, Computed) Node payment type. Valid values: `PostPaid`, `PrePaid`, default is `PostPaid`. If value is `PrePaid`, the arguments `period`, `period_unit`, `auto_renew` and `auto_renew_period` are required.
* `instance_metadata_options` - (Optional, ForceNew, Computed, List, Available since v1.266.0) ECS instance metadata access configuration. See [`instance_metadata_options`](#instance_metadata_options) below.
* `instance_patterns` - (Optional, List, Available since v1.266.0) Instance property configuration. See [`instance_patterns`](#instance_patterns) below.
* `instance_types` - (Optional, List) In the node instance specification list, you can select multiple instance specifications as alternatives. When each node is created, it will try to purchase from the first specification until it is created successfully. The final purchased instance specifications may vary with inventory changes.
* `internet_charge_type` - (Optional) The billing method for network usage. Valid values `PayByBandwidth` and `PayByTraffic`. Conflict with `eip_internet_charge_type`, EIP and public network IP can only choose one. 
* `internet_max_bandwidth_out` - (Optional, Int) The maximum bandwidth of the public IP address of the node. The unit is Mbps(Mega bit per second). The value range is:\[1,100\]
* `key_name` - (Optional) The name of the key pair. When the node pool is a managed node pool, only `key_name` is supported.
* `kubelet_configuration` - (Optional, List) Kubelet configuration parameters for worker nodes. See [`kubelet_configuration`](#kubelet_configuration) below. More information in [Kubelet Configuration](https://kubernetes.io/docs/reference/config-api/kubelet-config.v1beta1/). See [`kubelet_configuration`](#kubelet_configuration) below.
* `labels` - (Optional, List) A List of Kubernetes labels to assign to the nodes . Only labels that are applied with the ACK API are managed by this argument. Detailed below. More information in [Labels](https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/). See [`labels`](#labels) below.
* `login_as_non_root` - (Optional, ForceNew) Whether the ECS instance is logged on as a ecs-user user. Valid value: `true` and `false`.
* `management` - (Optional, Computed, List) Managed node pool configuration. See [`management`](#management) below.
* `multi_az_policy` - (Optional, Computed) The scaling policy for ECS instances in a multi-zone scaling group. Valid value: `PRIORITY`, `COST_OPTIMIZED` and `BALANCE`. `PRIORITY`: scales the capacity according to the virtual switches you define (VSwitchIds.N). When an ECS instance cannot be created in the zone where the higher-priority vSwitch is located, the next-priority vSwitch is automatically used to create an ECS instance. `COST_OPTIMIZED`: try to create by vCPU unit price from low to high. When the scaling configuration is configured with multiple instances of preemptible billing, preemptible instances are created first. You can continue to use the `CompensateWithOnDemand` parameter to specify whether to automatically try to create a preemptible instance by paying for it. It takes effect only when the scaling configuration has multi-instance specifications or preemptible instances. `BALANCE`: distributes ECS instances evenly among the multi-zone specified by the scaling group. If the zones become unbalanced due to insufficient inventory, you can use the API [RebalanceInstances](~~ 71516 ~~) to balance resources.
* `node_name_mode` - (Optional, ForceNew, Computed) Each node name consists of a prefix, its private network IP, and a suffix, separated by commas. The input format is `customized,,ip,`.
  - The prefix and suffix can be composed of one or more parts separated by '.', each part can use lowercase letters, numbers and '-', and the beginning and end of the node name must be lowercase letters and numbers.
  - The node IP address is the complete private IP address of the node.
  - For example, if the string `customized,aliyun,ip,com` is passed in (where 'customized' and 'ip' are fixed strings, 'aliyun' is the prefix, and 'com' is the suffix), the name of the node is `aliyun192.168.xxx.xxxcom`.
* `node_pool_name` - (Optional) The name of node pool.
* `on_demand_base_capacity` - (Optional) The minimum number of pay-as-you-go instances that must be kept in the scaling group. Valid values: 0 to 1000. If the number of pay-as-you-go instances is less than the value of this parameter, Auto Scaling preferably creates pay-as-you-go instances.
* `on_demand_percentage_above_base_capacity` - (Optional) The percentage of pay-as-you-go instances among the extra instances that exceed the number specified by `on_demand_base_capacity`. Valid values: 0 to 100.
* `password` - (Optional) The password of ssh login. You have to specify one of `password` and `key_name` fields. The password rule is 8 to 30 characters and contains at least three items (upper and lower case letters, numbers, and special symbols).
* `period` - (Optional, Int) Node payment period. Its valid value is one of {1, 2, 3, 6, 12}.
* `period_unit` - (Optional) Node payment period unit, valid value: `Month`. Default is `Month`.
* `platform` - (Optional, Computed, Deprecated since v1.145.0) Operating system release, using `image_type` instead.
* `pre_user_data` - (Optional, Available since v1.232.0) Node pre custom data, base64-encoded, the script executed before the node is initialized. 
* `private_pool_options` - (Optional, List) Private node pool configuration. See [`private_pool_options`](#private_pool_options) below.
* `ram_role_name` - (Optional, ForceNew, Computed, Available since v1.242.0) The name of the Worker RAM role.
* If it is empty, the default Worker RAM role created in the cluster will be used.
* If the specified RAM role is not empty, the specified RAM role must be a **Common Service role**, and its **trusted service** configuration must be **cloud server**. For more information, see [Create a common service role](https://help.aliyun.com/document_detail/116800.html). If the specified RAM role is not the default Worker RAM role created in the cluster, the role name cannot start with 'KubernetesMasterRole-'or 'KubernetesWorkerRole.

-> **NOTE:**  This parameter is only supported for ACK-managed clusters of 1.22 or later versions.

* `rds_instances` - (Optional, List) The list of RDS instances.
* `resource_group_id` - (Optional, Computed) The ID of the resource group
* `rolling_policy` - (Optional, List) Rotary configuration. See [`rolling_policy`](#rolling_policy) below.

-> **NOTE:** This parameter only applies during resource update. If modified in isolation without other property changes, Terraform will not trigger any action.

* `runtime_name` - (Optional, Computed) The runtime name of containers. If not set, the cluster runtime will be used as the node pool runtime. If you select another container runtime, see [Comparison of Docker, containerd, and Sandboxed-Container](https://www.alibabacloud.com/help/doc-detail/160313.htm).
* `runtime_version` - (Optional, Computed) The runtime version of containers. If not set, the cluster runtime will be used as the node pool runtime.
* `scaling_config` - (Optional, Computed, List) Automatic scaling configuration. See [`scaling_config`](#scaling_config) below.
* `scaling_policy` - (Optional, Computed) Scaling group mode, default value: `release`. Valid values:
  - `release`: in the standard mode, scaling is performed by creating and releasing ECS instances based on the usage of the application resource value.
  - `recycle`: in the speed mode, scaling is performed through creation, shutdown, and startup to increase the speed of scaling again (computing resources are not charged during shutdown, only storage fees are charged, except for local disk models).
* `security_group_id` - (Optional, ForceNew, Computed, Deprecated since v1.145.0) The security group ID of the node pool. This field has been replaced by `security_group_ids`, please use the `security_group_ids` field instead.
* `security_group_ids` - (Optional, ForceNew, Computed, Set) Multiple security groups can be configured for a node pool. If both `security_group_ids` and `security_group_id` are configured, `security_group_ids` takes effect. This field cannot be modified.
* `security_hardening_os` - (Optional, ForceNew) Alibaba Cloud OS security reinforcement. Default value: `false`. Value:
  -`true`: enable Alibaba Cloud OS security reinforcement.
  -`false`: does not enable Alibaba Cloud OS security reinforcement.
* `soc_enabled` - (Optional, ForceNew) Whether enable worker node to support soc security reinforcement, its valid value `true` or `false`. Default to `false` and apply to AliyunLinux series. See [SOC Reinforcement](https://help.aliyun.com/document_detail/196148.html).

-> **NOTE:**  It is forbidden to set both `security_hardening_os` and `soc_enabled` to `true` at the same time.

* `spot_instance_pools` - (Optional, Int) The number of instance types that are available. Auto Scaling creates preemptible instances of multiple instance types that are available at the lowest cost. Valid values: 1 to 10.
* `spot_instance_remedy` - (Optional) Specifies whether to supplement preemptible instances when the number of preemptible instances drops below the specified minimum number. If you set the value to true, Auto Scaling attempts to create a new preemptible instance when the system notifies that an existing preemptible instance is about to be reclaimed. Valid values: `true`: enables the supplementation of preemptible instances. `false`: disables the supplementation of preemptible instances.
* `spot_price_limit` - (Optional, List) The current single preemptible instance type market price range configuration. See [`spot_price_limit`](#spot_price_limit) below.
* `spot_strategy` - (Optional, Computed) The preemptible instance type. Value:
  - `NoSpot` : Non-preemptible instance.
  - `SpotWithPriceLimit` : Set the upper limit of the preemptible instance price.
  - `SpotAsPriceGo` : The system automatically bids, following the actual price of the current market.
* `system_disk_bursting_enabled` - (Optional) Specifies whether to enable the burst feature for system disks. Valid values:`true`: enables the burst feature. `false`: disables the burst feature. This parameter is supported only when `system_disk_category` is set to `cloud_auto`.
* `system_disk_categories` - (Optional, Computed, List) The multi-disk categories of the system disk. When a high-priority disk type cannot be used, Auto Scaling automatically tries to create a system disk with the next priority disk category. Valid values see `system_disk_category`.
* `system_disk_category` - (Optional, Computed) The category of the system disk for nodes. Default value: `cloud_efficiency`. Valid values:
  - `cloud`: basic disk.
  - `cloud_efficiency`: ultra disk.
  - `cloud_ssd`: standard SSD.
  - `cloud_essd`: ESSD.
  - `cloud_auto`: ESSD AutoPL disk.
  - `cloud_essd_entry`: ESSD Entry disk.
* `system_disk_encrypt_algorithm` - (Optional) The encryption algorithm used by the system disk. Value range: aes-256.
* `system_disk_encrypted` - (Optional) Whether to encrypt the system disk. Value range: `true`: encryption. `false`: Do not encrypt.
* `system_disk_kms_key` - (Optional) The ID of the KMS key used by the system disk.
* `system_disk_performance_level` - (Optional) The system disk performance of the node takes effect only for the ESSD disk.
  - `PL0`: maximum random read/write IOPS 10000 for a single disk.
  - `PL1`: maximum random read/write IOPS 50000 for a single disk.
  - `PL2`: highest random read/write IOPS 100000 for a single disk.
  - `PL3`: maximum random read/write IOPS 1 million for a single disk.
* `system_disk_provisioned_iops` - (Optional, Int) The predefined IOPS of a system disk. Valid values: 0 to min{50,000, 1,000 × Capacity - Baseline IOPS}. Baseline IOPS = min{1,800 + 50 × Capacity, 50,000}. This parameter is supported only when `system_disk_category` is set to `cloud_auto`.
* `system_disk_size` - (Optional, Int) The size of the system disk. Unit: GiB. The value of this parameter must be at least 1 and greater than or equal to the image size. Default value: 40 or the size of the image, whichever is larger.
  - Basic disk: 20 to 500.
  - ESSD (cloud_essd): The valid values vary based on the performance level of the ESSD. PL0 ESSD: 1 to 2048. PL1 ESSD: 20 to 2048. PL2 ESSD: 461 to 2048. PL3 ESSD: 1261 to 2048.
  - ESSD AutoPL disk (cloud_auto): 1 to 2048.
  - Other disk categories: 20 to 2048.
* `system_disk_snapshot_policy_id` - (Optional) The ID of the automatic snapshot policy used by the system disk.
* `tags` - (Optional, Map) Add tags only for ECS instances. The maximum length of the tag key is 128 characters. The tag key and value cannot start with aliyun or acs:, or contain https:// or http://.
* `taints` - (Optional, List) A List of Kubernetes taints to assign to the nodes. Detailed below. More information in [Taints and Toleration](https://kubernetes.io/docs/concepts/scheduling-eviction/taint-and-toleration/). See [`taints`](#taints) below.
* `tee_config` - (Optional, ForceNew, Computed, List) The configuration about confidential computing for the cluster. See [`tee_config`](#tee_config) below.
* `type` - (Optional, ForceNew, Computed, Available since v1.252.0) Node pool type, value range:
  -'ess': common node pool (including hosting function and auto scaling function).
  -'lingjun': Lingjun node pool.
* `unschedulable` - (Optional) Whether the node after expansion can be scheduled.
* `update_nodes` - (Optional) Synchronously update node labels and taints.

-> **NOTE:** This parameter only applies during resource update. If modified in isolation without other property changes, Terraform will not trigger any action.

* `user_data` - (Optional) Node custom data, base64-encoded.
* `vswitch_ids` - (Optional, List) The vswitches used by node pool workers.

* `kms_encrypted_password` - (Optional, Available since v1.177.0) An KMS encrypts password used to a cs kubernetes. You have to specify one of `password` `key_name` `kms_encrypted_password` fields.
* `kms_encryption_context` - (Optional, Available since v1.177.0) An KMS encryption context used to decrypt `kms_encrypted_password` before creating or updating a cs kubernetes with `kms_encrypted_password`. See [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm). It is valid when `kms_encrypted_password` is set.
* `format_disk` - (Optional, Available since v1.127.0) After you select this check box, if data disks have been attached to the specified ECS instances and the file system of the last data disk is uninitialized, the system automatically formats the last data disk to ext4 and mounts the data disk to /var/lib/docker and /var/lib/kubelet. The original data on the disk will be cleared. Make sure that you back up data in advance. If no data disk is mounted on the ECS instance, no new data disk will be purchased. Default is `false`.
* `instances` - (Optional, Available since v1.127.0) The instance list. Add existing nodes under the same cluster VPC to the node pool. 
* `node_count` (Optional, Deprecated) The worker node number of the node pool. From version 1.111.0, `node_count` is not required.
* `keep_instance_name` - (Optional, Available since v1.127.0) Add an existing instance to the node pool, whether to keep the original instance name. It is recommended to set to `true`.
* `rollout_policy` - (Optional, Deprecated since 1.185.0) Rollout policy is used to specify the strategy when the node pool is rolling update. This field works when node pool updating. Please use `rolling_policy` to instead it from provider version 1.185.0. See [`rollout_policy`](#rollout_policy) below.

The following arguments will be discarded. Please use new fields as soon as possible:
* `name` - (Deprecated since v1.219.0). Field 'name' has been deprecated from provider version 1.219.0. New field 'node_pool_name' instead.

### `auto_mode`

The auto_mode supports the following:
* `enabled` - (Optional, ForceNew, Available since v1.266.0) Whether to enable auto mode. Valid values: 
  - `true`: enables Smart Managed Configuration. 
  - `false`: disables Smart Managed Configuration.

-> **NOTE:** When `auto_mode.enable` is set to `true`, the system will automatically manage the node pool with optimized default configurations. **All parameters except the following can be specified or modified:**

**Parameters That Can Be Specified or Modified:**
  - `scaling_config.max_size`: default `50`, can be specified during creation and modified afterward
  - `scaling_config.min_size`: default `0`, can be specified during creation and modified afterward
  - `instance_types`: can be specified during creation and modified afterward. **Note:** `instance_types` and `instance_patterns` are mutually exclusive - you can only specify one of them.
  - `instance_patterns`: has default instance specification configuration (4-16 CPU cores, 8-32GB memory, etc.), can be specified during creation and modified afterward. **Note:** `instance_patterns` and `instance_types` are mutually exclusive - you can only specify one of them.
  - `data_disks`: can be specified during creation and modified afterward. If not specified during creation, default data disk configuration will be used (120GB size, supports cloud_auto, cloud_essd, cloud_ssd types). When specified, can configure `size`, `category`, `categories`, `performance_level`, `provisioned_iops`, `bursting_enabled`
  - `resource_group_id`: can be specified during creation and modified afterward
  - `vswitch_ids`: can be modified during creation and modified afterward
  - `tags`: can be modified during creation and modified afterward
  - `labels`: can be specified during creation and modified afterward
  - `taints`: can be specified during creation and modified afterward
  - `kubelet_configuration`:  can be specified during creation and modified afterward

**All Other Parameters:**

All other parameters (including but not limited to `management`, `scaling_config.enable`, `scaling_config.type`, `install_cloud_monitor`, `cpu_policy`, `node_name_mode`, `runtime_name`, `runtime_version`, `unschedulable`, `user_data`, `pre_user_data`, `auto_renew`, `auto_renew_period`, `cis_enabled`, `compensate_with_on_demand`, `deployment_set_id`, `image_id`, `image_type`, `instance_charge_type`, `instance_metadata_options`, `internet_charge_type`, `internet_max_bandwidth_out`, `key_name`, `login_as_non_root`, `password`, `multi_az_policy`, `on_demand_base_capacity`, `on_demand_percentage_above_base_capacity`, `period`, `period_unit`, `platform`, `private_pool_options`, `ram_role_name`, `rds_instances`, `scaling_policy`, `security_group_id`, `security_group_ids`, `security_hardening_os`, `soc_enabled`, `spot_instance_pools`, `spot_instance_remedy`, `spot_price_limit`, `spot_strategy`, all `system_disk_*` parameters, `tee_config`, etc.) will be automatically set to default values during node pool creation (user-specified values will be ignored), and cannot be modified after creation.

### `data_disks`

The data_disks supports the following:
* `auto_format` - (Optional, Available since v1.229.0) Whether to automatically mount the data disk. Valid values: true and false.
* `auto_snapshot_policy_id` - (Optional) The ID of the automatic snapshot policy that you want to apply to the system disk.
* `bursting_enabled` - (Optional) Whether the data disk is enabled with Burst (performance Burst). This is configured when the disk type is cloud_auto.
* `category` - (Optional) The type of data disk. Default value: `cloud_efficiency`. Valid values:
  - `cloud`: basic disk.
  - `cloud_efficiency`: ultra disk.
  - `cloud_ssd`: standard SSD.
  - `cloud_essd`: Enterprise SSD (ESSD).
  - `cloud_auto`: ESSD AutoPL disk.
  - `cloud_essd_entry`: ESSD Entry disk.
  - `elastic_ephemeral_disk_premium`: premium elastic ephemeral disk.
  - `elastic_ephemeral_disk_standard`: standard elastic ephemeral disk.
* `device` - (Optional) The mount target of data disk N. Valid values of N: 1 to 16. If you do not specify this parameter, the system automatically assigns a mount target when Auto Scaling creates ECS instances. The name of the mount target ranges from /dev/xvdb to /dev/xvdz.
* `encrypted` - (Optional) Specifies whether to encrypt data disks. Valid values: true and false. Default to `false`.
* `file_system` - (Optional, Available since v1.229.0) The type of the mounted file system. Works when auto_format is true. Optional value: `ext4`, `xfs`.
* `kms_key_id` - (Optional) The kms key id used to encrypt the data disk. It takes effect when `encrypted` is true.
* `mount_target` - (Optional, Available since v1.229.0) The Mount path. Works when auto_format is true.
* `name` - (Optional, Computed) The length is 2~128 English or Chinese characters. It must start with an uppercase or lowr letter or a Chinese character and cannot start with http:// or https. Can contain numbers, colons (:), underscores (_), or dashes (-). It will be overwritten if auto_format is set.
* `performance_level` - (Optional) Worker node data disk performance level, when `category` values `cloud_essd`, the optional values are `PL0`, `PL1`, `PL2` or `PL3`, but the specific performance level is related to the disk capacity. For more information, see [Enhanced SSDs](https://www.alibabacloud.com/help/doc-detail/122389.htm). Default is `PL1`.
* `provisioned_iops` - (Optional, Int) The read/write IOPS preconfigured for the data disk, which is configured when the disk type is cloud_auto.
* `size` - (Optional, Int) The size of a data disk, Its valid value range [40~32768] in GB. Default to `40`.
* `snapshot_id` - (Optional) The ID of the snapshot that you want to use to create data disk N. Valid values of N: 1 to 16. If you specify this parameter, DataDisk.N.Size is ignored. The size of the disk is the same as the size of the specified snapshot. If you specify a snapshot that is created on or before July 15, 2013, the operation fails and InvalidSnapshot.TooOld is returned.

### `eflo_node_group`

The eflo_node_group supports the following:
* `cluster_id` - (Optional, Available since v1.252.0) The ID of the associated Lingjun cluster is required when creating a Lingjun node pool.
* `group_id` - (Optional) When creating a Lingjun node pool, you need the Lingjun group ID of the associated Lingjun cluster.

### `instance_metadata_options`

The instance_metadata_options supports the following:
* `http_tokens` - (Optional, ForceNew, Available since v1.266.0) ECS instance metadata access mode configuration. Value range:

  - 'optional': Compatible with both normal mode and reinforced mode.
  - 'required': Enables only hardening mode (IMDSv2). When enabled, applications in the node cannot access the ECS instance metadata in normal mode. Ensure that the component and operating system versions in the cluster meet the minimum version requirements. For more information, see [accessing ECS instance metadata in hardened mode only](https://www.alibabacloud.com/help/ack/ack-managed-and-ack-dedicated/security-and-compliance/secure-access-to-ecs-instance-metadata).

Default value: 'optional '.

 This parameter is only supported for ACK-managed clusters of 1.28 or later versions. 

### `instance_patterns`

The instance_patterns supports the following:
* `cores` - (Optional, Int, Available since v1.266.0) The number of vCPU cores of the instance type. Example value: 8.
* `cpu_architectures` - (Optional, List, Available since v1.266.0) The CPU architecture of the instance. Value range:
  - X86
  - ARM
* `excluded_instance_types` - (Optional, List, Available since v1.266.0) Instance specifications to be excluded. You can exclude individual specifications or entire specification families by using the wildcard character (*). For example:
  - ecs.c6.large: indicates that the ecs.c6.large instance type is excluded.
  - ecs.c6. *: indicates that the instance specification of the entire c6 specification family is excluded.
* `instance_categories` - (Optional, List, Available since v1.266.0) Instance classification. Value range:
  - General-purpose: Universal.
  - Compute-optimized: Compute type.
  - Memory-optimized: Memory type.
  - Big data: Big data type.
  - Local SSDs: Local SSD type.
  - High Clock Speed: High frequency type.
  - Enhanced: Enhanced.
  - Shared: Shared.
  - ECS Bare Metal: elastic Bare Metal server.
  - High Performance Compute: High Performance Compute.
* `instance_family_level` - (Required, Available since v1.266.0) Instance specification family level, value range:
  - EntryLevel: entry-level, that is, shared instance specifications. The cost is lower, but the stability of instance computing performance cannot be guaranteed. Applicable to business scenarios with low CPU usage. For more information, see Shared.
  - EnterpriseLevel: Enterprise level. Stable performance and exclusive resources, suitable for business scenarios that require high stability. For more information, see Instance Specification Family.
* `instance_type_families` - (Optional, List, Available since v1.266.0) Specifies the instance type family. Example values:["ecs.g8i","ecs.c8i"]
* `max_cpu_cores` - (Optional, Int, Available since v1.266.0) The maximum number of vCPU cores of the instance type. Example value: 8. MaxCpuCores cannot exceed 4 times of MinCpuCores.
* `max_memory_size` - (Optional, Float, Available since v1.266.0) The maximum memory of the instance type. Unit: GiB, example value: 8,MaxMemoryCores does not support more than 4 times MinMemoryCores.
* `memory` - (Optional, Float, Available since v1.266.0) The memory size of the instance type, in GiB. Example value: 8.
* `min_cpu_cores` - (Optional, Int, Available since v1.266.0) The minimum number of vCPU cores of the instance type. Example value: 4. MaxCpuCores cannot exceed 4 times of MinCpuCores.
* `min_memory_size` - (Optional, Float, Available since v1.266.0) The minimum memory of the instance type. Unit: GiB, example value: 4,MaxMemoryCores does not support more than 4 times MinMemoryCores.

### `kubelet_configuration`

The kubelet_configuration supports the following:
* `allowed_unsafe_sysctls` - (Optional, List) Allowed sysctl mode whitelist.
* `cluster_dns` - (Optional, List, Available since v1.242.0) The list of IP addresses of the cluster DNS servers.
* `container_log_max_files` - (Optional) The maximum number of log files that can exist in each container.
* `container_log_max_size` - (Optional) The maximum size that can be reached before a log file is rotated.
* `container_log_max_workers` - (Optional, Available since v1.242.0) Specifies the maximum number of concurrent workers required to perform log rotation operations.
* `container_log_monitor_interval` - (Optional, Available since v1.242.0) Specifies the duration for which container logs are monitored for log rotation.
* `cpu_cfs_quota` - (Optional, Available since v1.242.0) CPU CFS quota constraint switch.
* `cpu_cfs_quota_period` - (Optional, Available since v1.242.0) CPU CFS quota period value.
* `cpu_manager_policy` - (Optional) Same as cpuManagerPolicy. The name of the policy to use. Requires the CPUManager feature gate to be enabled. Valid value is `none` or `static`.
* `event_burst` - (Optional) Same as eventBurst. The maximum size of a burst of event creations, temporarily allows event creations to burst to this number, while still not exceeding `event_record_qps`. It is only used when `event_record_qps` is greater than 0. Valid value is `[0-100]`.
* `event_record_qps` - (Optional) Same as eventRecordQPS. The maximum event creations per second. If 0, there is no limit enforced. Valid value is `[0-50]`.
* `eviction_hard` - (Optional, Map) Same as evictionHard. The map of signal names to quantities that defines hard eviction thresholds. For example: `{"memory.available" = "300Mi"}`.
* `eviction_soft` - (Optional, Map) Same as evictionSoft. The map of signal names to quantities that defines soft eviction thresholds. For example: `{"memory.available" = "300Mi"}`.
* `eviction_soft_grace_period` - (Optional, Map) Same as evictionSoftGracePeriod. The map of signal names to quantities that defines grace periods for each soft eviction signal. For example: `{"memory.available" = "30s"}`.
* `feature_gates` - (Optional, Map) Feature switch to enable configuration of experimental features.
* `image_gc_high_threshold_percent` - (Optional, Available since v1.242.0) If the image usage exceeds this threshold, image garbage collection will continue.
* `image_gc_low_threshold_percent` - (Optional, Available since v1.242.0) Image garbage collection is not performed when the image usage is below this threshold.
* `kube_api_burst` - (Optional) Same as kubeAPIBurst. The burst to allow while talking with kubernetes api-server. Valid value is `[0-100]`.
* `kube_api_qps` - (Optional) Same as kubeAPIQPS. The QPS to use while talking with kubernetes api-server. Valid value is `[0-50]`.
* `kube_reserved` - (Optional, Map) Same as kubeReserved. The set of ResourceName=ResourceQuantity (e.g. cpu=200m,memory=150G) pairs that describe resources reserved for kubernetes system components. Currently, cpu, memory and local storage for root file system are supported. See [compute resources](https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/) for more details.
* `max_pods` - (Optional) The maximum number of running pods.
* `memory_manager_policy` - (Optional, Available since v1.242.0) The policy to be used by the memory manager.
* `pod_pids_limit` - (Optional, Available since v1.242.0) The maximum number of PIDs that can be used in a Pod.
* `read_only_port` - (Optional) Read-only port number.
* `registry_burst` - (Optional) Same as registryBurst. The maximum size of burst pulls, temporarily allows pulls to burst to this number, while still not exceeding `registry_pull_qps`. Only used if `registry_pull_qps` is greater than 0. Valid value is `[0-100]`.
* `registry_pull_qps` - (Optional) Same as registryPullQPS. The limit of registry pulls per second. Setting it to `0` means no limit. Valid value is `[0-50]`.
* `reserved_memory` - (Optional, List, Available since v1.242.0) Reserve memory for NUMA nodes. See [`reserved_memory`](#kubelet_configuration-reserved_memory) below.
* `serialize_image_pulls` - (Optional) Same as serializeImagePulls. When enabled, it tells the Kubelet to pull images one at a time. We recommend not changing the default value on nodes that run docker daemon with version < 1.9 or an Aufs storage backend. Valid value is `true` or `false`.
* `server_tls_bootstrap` - (Optional, Available since v1.266.0) Used to enable the kubelet server certificate signing and rotation via CSR.
* `system_reserved` - (Optional, Map) Same as systemReserved. The set of ResourceName=ResourceQuantity (e.g. cpu=200m,memory=150G) pairs that describe resources reserved for non-kubernetes components. Currently, only cpu and memory are supported. See [compute resources](https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/) for more details.
* `topology_manager_policy` - (Optional, Available since v1.242.0) Name of the Topology Manager policy used.
* `tracing` - (Optional, List) OpenTelemetry tracks the configuration information for client settings versioning. See [`tracing`](#kubelet_configuration-tracing) below.

### `kubelet_configuration-reserved_memory`

The kubelet_configuration-reserved_memory supports the following:
* `limits` - (Optional, Map, Available since v1.242.0) Memory resource limit.
* `numa_node` - (Optional, Int) The NUMA node.

### `kubelet_configuration-tracing`

The kubelet_configuration-tracing supports the following:
* `endpoint` - (Optional, Available since v1.242.0) The endpoint of the collector.
* `sampling_rate_per_million` - (Optional) Number of samples to be collected per million span.

### `labels`

The labels supports the following:
* `key` - (Required) The label key.
* `value` - (Optional) The label value.

### `management`

The management supports the following:
* `auto_repair` - (Optional, Computed) Whether to enable automatic repair. Valid values: `true`: Automatic repair. `false`: not automatically repaired.
* `auto_repair_policy` - (Optional, Computed, List) Automatic repair node policy. See [`auto_repair_policy`](#management-auto_repair_policy) below.
* `auto_upgrade` - (Optional, Computed) Specifies whether to enable auto update. Valid values: `true`: enables auto update. `false`: disables auto update.
* `auto_upgrade_policy` - (Optional, Computed, List) The auto update policy. See [`auto_upgrade_policy`](#management-auto_upgrade_policy) below.
* `auto_vul_fix` - (Optional, Computed) Specifies whether to automatically patch CVE vulnerabilities. Valid values: `true`, `false`.
* `auto_vul_fix_policy` - (Optional, Computed, List) The auto CVE patching policy. See [`auto_vul_fix_policy`](#management-auto_vul_fix_policy) below.
* `enable` - (Optional, Computed) Specifies whether to enable the managed node pool feature. Valid values: `true`: enables the managed node pool feature. `false`: disables the managed node pool feature. Other parameters in this section take effect only when you specify enable=true.
* `max_unavailable` - (Optional, Int) Maximum number of unavailable nodes. Default value: 1. Value range:\[1,1000\].
* `surge` - (Optional, Int, Deprecated since v1.219.0) Number of additional nodes. You have to specify one of surge, surge_percentage.
* `surge_percentage` - (Optional, Int, Deprecated since v1.219.0) Proportion of additional nodes. You have to specify one of surge, surge_percentage.

### `management-auto_repair_policy`

The management-auto_repair_policy supports the following:
* `restart_node` - (Optional, Computed) Whether to allow node restart.

### `management-auto_upgrade_policy`

The management-auto_upgrade_policy supports the following:
* `auto_upgrade_kubelet` - (Optional, Computed) Specifies whether  to automatically update the kubelet. Valid values: `true`: yes; `false`: no.

### `management-auto_vul_fix_policy`

The management-auto_vul_fix_policy supports the following:
* `restart_node` - (Optional, Computed) Specifies whether to automatically restart nodes after patching CVE vulnerabilities. Valid values: `true`, `false`.
* `vul_level` - (Optional, Computed) The severity levels of vulnerabilities that is allowed to automatically patch. Multiple severity levels are separated by commas (,).

### `private_pool_options`

The private_pool_options supports the following:
* `private_pool_options_id` - (Optional) The ID of the private node pool.
* `private_pool_options_match_criteria` - (Optional) The type of private node pool. This parameter specifies the type of the private pool that you want to use to create instances. A private node pool is generated when an elasticity assurance or a capacity reservation service takes effect. The system selects a private node pool to launch instances. Valid values: `Open`: specifies an open private node pool. The system selects an open private node pool to launch instances. If no matching open private node pool is available, the resources in the public node pool are used. `Target`: specifies a private node pool. The system uses the resources of the specified private node pool to launch instances. If the specified private node pool is unavailable, instances cannot be started. `None`: no private node pool is used. The resources of private node pools are not used to launch the instances.

### `rolling_policy`

The rolling_policy supports the following:
* `max_parallelism` - (Optional, Int) The maximum number of unusable nodes.

### `scaling_config`

The scaling_config supports the following:
* `eip_bandwidth` - (Optional, Int) Peak EIP bandwidth. Its valid value range [1~500] in Mbps. It works if `is_bond_eip=true`. Default to `5`.
* `eip_internet_charge_type` - (Optional) EIP billing type. `PayByBandwidth`: Charged at fixed bandwidth. `PayByTraffic`: Billed as used traffic. Default: `PayByBandwidth`. It works if `is_bond_eip=true`, conflict with `internet_charge_type`. EIP and public network IP can only choose one.
* `enable` - (Optional) Whether to enable automatic scaling. Value:
  - `true`: enables the node pool auto-scaling function.
  - `false`: Auto scaling is not enabled. When the value is false, other `auto_scaling` configuration parameters do not take effect.
* `is_bond_eip` - (Optional) Whether to bind EIP for an instance. Default: `false`.
* `max_size` - (Optional, Int) Max number of instances in a auto scaling group, its valid value range [0~1000]. `max_size` has to be greater than `min_size`.
* `min_size` - (Optional, Int) Min number of instances in a auto scaling group, its valid value range [0~1000].
* `type` - (Optional) Instance classification, not required. Vaild value: `cpu`, `gpu`, `gpushare` and `spot`. Default: `cpu`. The actual instance type is determined by `instance_types`.

### `spot_price_limit`

The spot_price_limit supports the following:
* `instance_type` - (Optional) The type of the preemptible instance.
* `price_limit` - (Optional) The maximum price of a single instance.

### `taints`

The taints supports the following:
* `effect` - (Optional) The scheduling policy.
* `key` - (Required) The key of a taint.
* `value` - (Optional) The value of a taint.

### `tee_config`

The tee_config supports the following:
* `tee_enable` - (Optional, ForceNew) Specifies whether to enable confidential computing for the cluster.

### `rollout_policy`

The rollout_policy mapping supports the following:
* `max_unavailable` - (Optional, Deprecated since 1.185.0) Maximum number of unavailable nodes during rolling upgrade. The value of this field should be greater than `0`, and if it's set to a number less than or equal to `0`, the default setting will be used. Please use `max_parallelism` to instead it from provider version 1.185.0.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<cluster_id>:<node_pool_id>`.
* `node_pool_id` - The first ID of the resource.
* `scaling_group_id` - The ID of the scaling group.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 90 mins) Used when create the Nodepool.
* `delete` - (Defaults to 60 mins) Used when delete the Nodepool.
* `update` - (Defaults to 60 mins) Used when update the Nodepool.

## Import

Container Service for Kubernetes (ACK) Nodepool can be imported using the id, e.g.

```shell
$ terraform import alicloud_cs_kubernetes_node_pool.example <cluster_id>:<node_pool_id>
```