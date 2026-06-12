# Affected Products

| # | Product Code | Resources | Data Sources |
|---|---|---|---|
| 1 | `CS` (Container Service / ACK) | `alicloud_cs_serverless_kubernetes`, `alicloud_cs_edge_kubernetes`, `alicloud_cs_kubernetes_node_pool`, `alicloud_cs_kubernetes`, `alicloud_cs_managed_kubernetes` | `alicloud_cs_managed_kubernetes_clusters`, `alicloud_cs_cluster_credential`, `alicloud_cs_edge_kubernetes_clusters`, `alicloud_cs_kubernetes_clusters`, `alicloud_cs_serverless_kubernetes_clusters` |
| 2 | `CloudAPI` (API Gateway) | `alicloud_api_gateway_instance` | — |
| 3 | `cr` (Container Registry) | `alicloud_cr_repo` | `alicloud_cr_repos` |
| 4 | `clickhouse` (ClickHouse) | `alicloud_click_house_db_cluster` | `alicloud_click_house_db_clusters` |
| 5 | `amqp-open` (RabbitMQ) | `alicloud_amqp_queue` | — |
| 6 | `Cms` (Cloud Monitor) | `alicloud_cms_alarm` | — |
| 7 | `Eci` (Elastic Container Instance) | `alicloud_eci_container_group` | — |
| 8 | `Ecs` (Elastic Compute) | `alicloud_disk` / `alicloud_ecs_disk`, `alicloud_instance` | `alicloud_instance_types` |
| 9 | `Slb` (Server Load Balancer) | `alicloud_slb` / `alicloud_slb_load_balancer`, `alicloud_slb_listener`, `alicloud_slb_master_slave_server_group` | `alicloud_slbs` / `alicloud_slb_load_balancers`, `alicloud_slb_master_slave_server_groups` |
| 10 | `ResourceManager` | `alicloud_resource_manager_role`, `alicloud_resource_manager_resource_group`, `alicloud_resource_manager_policy_version` | — |
| 11 | `pvtz` (Private Zone) | `alicloud_pvtz_zone` | `alicloud_pvtz_zones` |
| 12 | `R-kvstore` (Redis) | `alicloud_kvstore_instance` | — |
| 13 | `Vpc` (Express Connect) | `alicloud_express_connect_virtual_border_router` | — |
| 14 | `hbr` (Hybrid Backup Recovery) | `alicloud_hbr_vault` | — |
| 15 | `ROS` (Resource Orchestration) | `alicloud_ros_stack_group` | — |
| 16 | `Vpc` (NAT Gateway) | `alicloud_nat_gateway`, `alicloud_vpc_nat_ip` | — |
| 17 | `Cloudfw` (Cloud Firewall) | `alicloud_cloud_firewall_instance` | `alicloud_cloud_firewall_control_policies` |
| 18 | `hitsdb` (Lindorm) | `alicloud_lindorm_instance` | — |
| 19 | `sae` (Serverless App Engine) | `alicloud_sae_application` | — |
| 20 | `Config` (Cloud Config) | — | `alicloud_config_rules` |
| 21 | `Vpc` (Elastic IP) | — | `alicloud_eips` / `alicloud_eip_addresses` |
| 22 | `HBase` | — | `alicloud_hbase_zones` |
| 23 | `Kms` (Key Management) | — | `alicloud_kms_key_versions` |
| 24 | `ots` (Table Store) | — | `alicloud_ots_instances` |
| 25 | `Ram` (Resource Access Mgmt) | — | `alicloud_ram_policies`, `alicloud_ram_users` |
| 26 | `Rds` (ApsaraDB RDS) | — | `alicloud_db_instance_classes` |

# Migration Examples

## 1. Field renamed (use a different field)

```hcl
# Before
resource "alicloud_cs_serverless_kubernetes" "example" {
  vswitch_id = "vsw-abc123"
}

# After
resource "alicloud_cs_serverless_kubernetes" "example" {
  vswitch_ids = ["vsw-abc123"]
}
```

## 2. Field moved to a different resource

```hcl
# Before
resource "alicloud_cs_managed_kubernetes" "example" {
  worker_instance_types = ["ecs.n4.large"]
  worker_number         = 3
  worker_disk_size      = 40
  worker_disk_category  = "cloud_efficiency"
  kube_config           = "~/.kube/config"
}

# After - remove worker fields, manage via node pool instead
resource "alicloud_cs_managed_kubernetes" "example" {
  # worker_* fields removed
}

resource "alicloud_cs_kubernetes_node_pool" "example" {
  cluster_id           = alicloud_cs_managed_kubernetes.example.id
  instance_types       = ["ecs.n4.large"]
  desired_size         = "3"
  system_disk_size     = 40
  system_disk_category = "cloud_efficiency"
}

data "alicloud_cs_cluster_credential" "example" {
  cluster_id  = alicloud_cs_managed_kubernetes.example.id
  output_file = "~/.kube/config"
}
```

## 3. Field removed with no replacement

```hcl
# Before
resource "alicloud_cs_kubernetes" "example" {
  force_update      = true
  availability_zone = "cn-hangzhou-a"
}

# After - just remove the fields
resource "alicloud_cs_kubernetes" "example" {
}
```

## 4. TypeMap changed to TypeList

```hcl
# Before - accessing as a map
output "cluster_api_server" {
  value = alicloud_cs_managed_kubernetes.example.connections["api_server_internet"]
}

output "ca_cert" {
  value = alicloud_cs_managed_kubernetes.example.certificate_authority["cluster_cert"]
}

# After - accessing as a list (first element)
output "cluster_api_server" {
  value = alicloud_cs_managed_kubernetes.example.connections[0].api_server_internet
}

output "ca_cert" {
  value = alicloud_cs_managed_kubernetes.example.certificate_authority[0].cluster_cert
}
```

# Resources

---

`alicloud_cs_serverless_kubernetes` (Deprecated)
- Removing the removed field `vswitch_id`, use `vswitch_ids` instead.
- Removing the removed field `force_update`.
- Removing the removed field `create_v2_cluster`.

---

`alicloud_cs_edge_kubernetes` (Deprecated)
- Removing the removed field `force_update`.
- Change the type of `runtime`, `certificate_authority`, `connections` from `TypeMap` to `TypeList`.

---

`alicloud_cs_kubernetes_node_pool`
- Removing the removed field `vpc_id`.
- Removing the removed field `rollout_policy`, use `rolling_policy` instead.

---

`alicloud_cs_kubernetes` (Deprecated)
- Removing the removed field `worker_vswitch_ids`, use resource `alicloud_cs_kubernetes_node_pool.vswitch_ids` instead.
- Removing the removed field `worker_instance_types`, use resource `alicloud_cs_kubernetes_node_pool.instance_types` instead.
- Removing the removed field `worker_number`, use resource `alicloud_cs_kubernetes_node_pool.desired_size` instead.
- Removing the removed field `worker_disk_size`, use resource `alicloud_cs_kubernetes_node_pool.system_disk_size` instead.
- Removing the removed field `worker_disk_category`, use resource `alicloud_cs_kubernetes_node_pool.system_disk_category` instead.
- Removing the removed field `worker_disk_performance_level`, use resource `alicloud_cs_kubernetes_node_pool.system_disk_performance_level` instead.
- Removing the removed field `worker_disk_snapshot_policy_id`, use resource `alicloud_cs_kubernetes_node_pool.system_disk_snapshot_policy_id` instead.
- Removing the removed field `worker_data_disk_size`, use resource `alicloud_cs_kubernetes_node_pool.data_disks.size` instead.
- Removing the removed field `worker_data_disk_category`, use resource `alicloud_cs_kubernetes_node_pool.data_disks.category` instead.
- Removing the removed field `worker_data_disks`, use resource `alicloud_cs_kubernetes_node_pool.data_disks` instead.
- Removing the removed field `worker_instance_charge_type`, use resource `alicloud_cs_kubernetes_node_pool.instance_charge_type` instead.
- Removing the removed field `worker_period_unit`, use resource `alicloud_cs_kubernetes_node_pool.period_unit` instead.
- Removing the removed field `worker_period`, use resource `alicloud_cs_kubernetes_node_pool.period` instead.
- Removing the removed field `worker_auto_renew`, use resource `alicloud_cs_kubernetes_node_pool.auto_renew` instead.
- Removing the removed field `worker_auto_renew_period`, use resource `alicloud_cs_kubernetes_node_pool.auto_renew_period` instead.
- Removing the removed field `exclude_autoscaler_nodes`, use resource `alicloud_cs_kubernetes_node_pool` instead.
- Removing the removed field `cpu_policy`, use resource `alicloud_cs_kubernetes_node_pool.cpu_policy` instead.
- Removing the removed field `node_port_range`, use resource `alicloud_cs_kubernetes_node_pool` instead.
- Removing the removed field `taints`, use resource `alicloud_cs_kubernetes_node_pool.taints` instead.
- Removing the removed field `kube_config`, use data source `alicloud_cs_cluster_credential.output_file` instead.
- Removing the removed field `worker_nodes`, use resource `alicloud_cs_kubernetes_node_pool` instead.
- Removing the removed field `worker_instance_type`, use resource `alicloud_cs_kubernetes_node_pool.instance_types` instead.
- Removing the removed field `vswitch_ids`, use field `master_vswitch_ids` instead.
- Removing the removed field `master_instance_type`, use field `master_instance_types` instead.
- Removing the removed field `force_update`.
- Removing the removed field `availability_zone`.
- Removing the removed field `vswitch_id`, use field `master_vswitch_ids` instead.
- Removing the removed field `worker_numbers`, use resource `alicloud_cs_kubernetes_node_pool.desired_size` instead.
- Removing the removed field `nodes`, use field `master_nodes` instead.
- Removing the removed field `log_config`, use field `addons` instead.
- Removing the removed field `cluster_network_type`, use field `addons` instead.
- Removing the removed field `user_data`, use resource `alicloud_cs_kubernetes_node_pool.user_data` instead.
- Change the type of `runtime`, `certificate_authority`, `connections` from `TypeMap` to `TypeList`.

---

`alicloud_cs_managed_kubernetes`
- Removing the removed field `worker_instance_types`, use resource `alicloud_cs_kubernetes_node_pool.instance_types` instead.
- Removing the removed field `worker_number`, use resource `alicloud_cs_kubernetes_node_pool.desired_size` instead.
- Removing the removed field `worker_disk_size`, use resource `alicloud_cs_kubernetes_node_pool.system_disk_size` instead.
- Removing the removed field `worker_disk_category`, use resource `alicloud_cs_kubernetes_node_pool.system_disk_category` instead.
- Removing the removed field `worker_disk_performance_level`, use resource `alicloud_cs_kubernetes_node_pool.system_disk_performance_level` instead.
- Removing the removed field `worker_disk_snapshot_policy_id`, use resource `alicloud_cs_kubernetes_node_pool.system_disk_snapshot_policy_id` instead.
- Removing the removed field `worker_data_disk_size`, use resource `alicloud_cs_kubernetes_node_pool.data_disks.size` instead.
- Removing the removed field `worker_data_disk_category`, use resource `alicloud_cs_kubernetes_node_pool.data_disks.category` instead.
- Removing the removed field `worker_instance_charge_type`, use resource `alicloud_cs_kubernetes_node_pool.instance_charge_type` instead.
- Removing the removed field `worker_data_disks`, use resource `alicloud_cs_kubernetes_node_pool.data_disks` instead.
- Removing the removed field `worker_period_unit`, use resource `alicloud_cs_kubernetes_node_pool.period_unit` instead.
- Removing the removed field `worker_period`, use resource `alicloud_cs_kubernetes_node_pool.period` instead.
- Removing the removed field `worker_auto_renew`, use resource `alicloud_cs_kubernetes_node_pool.auto_renew` instead.
- Removing the removed field `worker_auto_renew_period`, use resource `alicloud_cs_kubernetes_node_pool.auto_renew_period` instead.
- Removing the removed field `exclude_autoscaler_nodes`, use resource `alicloud_cs_kubernetes_node_pool` instead.
- Removing the removed field `enable_ssh`.
- Removing the removed field `password`, use resource `alicloud_cs_kubernetes_node_pool.password` instead.
- Removing the removed field `key_name`, use resource `alicloud_cs_kubernetes_node_pool.key_name` instead.
- Removing the removed field `kms_encrypted_password`, use resource `alicloud_cs_kubernetes_node_pool.kms_encrypted_password` instead.
- Removing the removed field `kms_encryption_context`, use resource `alicloud_cs_kubernetes_node_pool.kms_encryption_context` instead.
- Removing the removed field `image_id`, use resource `alicloud_cs_kubernetes_node_pool.image_id` instead.
- Removing the removed field `install_cloud_monitor`, use resource `alicloud_cs_kubernetes_node_pool.install_cloud_monitor` instead.
- Removing the removed field `cpu_policy`, use resource `alicloud_cs_kubernetes_node_pool.cpu_policy` instead.
- Removing the removed field `os_type`, use resource `alicloud_cs_kubernetes_node_pool` instead.
- Removing the removed field `platform`, use resource `alicloud_cs_kubernetes_node_pool.platform` instead.
- Removing the removed field `node_port_range`, use resource `alicloud_cs_kubernetes_node_pool` instead.
- Removing the removed field `runtime`, use resource `alicloud_cs_kubernetes_node_pool.runtime_name` and `alicloud_cs_kubernetes_node_pool.runtime_version` instead.
- Removing the removed field `taints`, use resource `alicloud_cs_kubernetes_node_pool.taints` instead.
- Removing the removed field `rds_instances`, use resource `alicloud_cs_kubernetes_node_pool.rds_instances` instead.
- Removing the removed field `user_data`, use resource `alicloud_cs_kubernetes_node_pool.user_data` instead.
- Removing the removed field `node_name_mode`, use resource `alicloud_cs_kubernetes_node_pool.node_name_mode` instead.
- Removing the removed field `worker_nodes`, use resource `alicloud_cs_kubernetes_node_pool` instead.
- Removing the removed field `kube_config`, use data source `alicloud_cs_cluster_credential.output_file` instead.
- Removing the removed field `availability_zone`.
- Removing the removed field `force_update`.
- Removing the removed field `worker_numbers`, use resource `alicloud_cs_kubernetes_node_pool.desired_size` instead.
- Removing the removed field `cluster_network_type`, use field `addons` instead.
- Removing the removed field `log_config`, use field `addons` instead.
- Removing the removed field `worker_instance_type`, use resource `alicloud_cs_kubernetes_node_pool.instance_types` instead.
- Change the type of `certificate_authority`, `connections` from `TypeMap` to `TypeList`.

---

`alicloud_api_gateway_instance`
- Change the type of `to_connect_vpc_ip_block` from `TypeMap` to `TypeList`.

---

`alicloud_cr_repo`
- Change the type of `domain_list` from `TypeMap` to `TypeList`.

---

`alicloud_click_house_db_cluster`
- Removing the removed field `db_cluster_access_white_list.db_cluster_ip_array_attribute`.

---

`alicloud_amqp_queue`
- Removing the removed field `exclusive_state`.

---

`alicloud_cms_alarm`
- Removing the removed field `operator`, use `escalations_critical.comparison_operator` instead.
- Removing the removed field `statistics`, use `escalations_critical.statistics` instead.
- Removing the removed field `threshold`, use `escalations_critical.threshold` instead.
- Removing the removed field `triggered_count`, use `escalations_critical.times` instead.
- Removing the removed field `notify_type`.

---

`alicloud_eci_container_group`
- Removing the removed field `eci_security_context`.

---

`alicloud_disk` / `alicloud_ecs_disk`
- Removing the removed field `dedicated_block_storage_cluster_id`.

---

`alicloud_slb_master_slave_server_group`
- Removing the removed field `servers.is_backup`.

---

`alicloud_slb` / `alicloud_slb_load_balancer`
- Removing the removed field `internet`, use `address_type` instead.

---

`alicloud_slb_listener`
- Removing the removed fields `lb_port`, `instance_port`, `lb_protocol`.

---

`alicloud_resource_manager_role`
- Removing the removed field `create_date`.

---

`alicloud_resource_manager_resource_group`
- Removing the removed field `create_date`.

---

`alicloud_resource_manager_policy_version`
- Removing the removed fields `create_date`, `version_id`.

---

`alicloud_pvtz_zone`
- Removing the removed fields `creation_time`, `update_time`.

---

`alicloud_kvstore_instance`
- Removing the removed field `modify_mode`.

---

`alicloud_instance`
- Removing the removed field `io_optimized`.
- Removing the removed field `subnet_id`, use `vswitch_id` instead.

---

`alicloud_express_connect_virtual_border_router`
- Removing the removed field `include_cross_account_vbr`.

---

`alicloud_hbr_vault`
- Removing the removed field `redundancy_type`.

---

`alicloud_ros_stack_group`
- Removing the removed fields `account_ids`, `operation_description`, `operation_preferences`, `region_ids`.

---

`alicloud_nat_gateway`
- Removing the removed fields `bandwidth_packages`, `bandwidth_packages.ip_count`, `bandwidth_packages.bandwidth`, `bandwidth_packages.zone`, `bandwidth_packages.public_ip_addresses`, `bandwidth_package_ids`.
- Removing the removed field `spec`, use `specification` instead.

---

`alicloud_vpc_nat_ip`
- Removing the removed field `nat_ip_cidr_id`.

---

`alicloud_cloud_firewall_instance`
- Removing the removed field `cfw_service`.

---

`alicloud_lindorm_instance`
- Removing the removed fields `core_num`, `group_name`, `phoenix_node_count`, `phoenix_node_specification`, `upgrade_type`.

---

`alicloud_sae_application`
- Removing the removed fields `version_id`, `mount_desc`, `mount_host`, `nas_id`.

# Data Sources

---

`alicloud_cs_managed_kubernetes_clusters` (Deprecated)
- Removing the removed field `slb_internet_enabled` from block `clusters`.
- Removing the removed field `vswitch_ids` from block `clusters`.
- Removing the removed field `worker_instance_types` from block `clusters`.
- Removing the removed field `worker_numbers` from block `clusters`.
- Removing the removed field `key_name` from block `clusters`.
- Removing the removed field `pod_cidr` from block `clusters`.
- Removing the removed field `service_cidr` from block `clusters`.
- Removing the removed field `cluster_network_type` from block `clusters`.
- Removing the removed field `log_config` from block `clusters`.
- Removing the removed field `image_id` from block `clusters`.
- Removing the removed field `worker_disk_size` from block `clusters`.
- Removing the removed field `worker_disk_category` from block `clusters`.
- Removing the removed field `worker_data_disk_size` from block `clusters`.
- Removing the removed field `worker_data_disk_category` from block `clusters`.
- Removing the removed field `worker_instance_charge_type` from block `clusters`.
- Removing the removed field `worker_period_unit` from block `clusters`.
- Removing the removed field `worker_period` from block `clusters`.
- Removing the removed field `worker_auto_renew` from block `clusters`.
- Removing the removed field `worker_auto_renew_period` from block `clusters`.
- Change the type of `clusters.connections` from `TypeMap` to `TypeList`.

---

`alicloud_cr_repos`
- Change the type of `repos.domain_list` from `TypeMap` to `TypeList`.

---

`alicloud_cs_cluster_credential`
- Change the type of `certificate_authority` from `TypeMap` to `TypeList`.

---

`alicloud_cs_edge_kubernetes_clusters`
- Change the type of `clusters.connections` from `TypeMap` to `TypeList`.

---

`alicloud_cs_kubernetes_clusters`
- Change the type of `clusters.connections` from `TypeMap` to `TypeList`.

---

`alicloud_cs_serverless_kubernetes_clusters`
- Change the type of `clusters.connections` from `TypeMap` to `TypeList`.

---

`alicloud_db_instance_classes`
- Change the type of `instance_classes.storage_range` from `TypeMap` to `TypeList`.

---

`alicloud_instance_types`
- Change the type of `instance_types.gpu`, `instance_types.burstable_instance`, `instance_types.local_storage` from `TypeMap` to `TypeList`.

---

`alicloud_click_house_db_clusters`
- Removing the removed field `clusters.db_cluster_access_white_list.db_cluster_ip_array_attribute`.

---

`alicloud_cloud_firewall_control_policies`
- Removing the removed field `source_ip`.

---

`alicloud_config_rules`
- Removing the removed field `multi_account`, use resource `alicloud_config_aggregate_config_rule` instead.
- Removing the removed field `member_id`, use resource `alicloud_config_aggregate_config_rule` instead.
- Removing the removed field `message_type`.

---

`alicloud_eips` / `alicloud_eip_addresses`
- Removing the removed field `in_use`.

---

`alicloud_hbase_zones`
- Removing the removed fields `multi`, `zones.multi_zone_ids`.

---

`alicloud_kms_key_versions`
- Removing the removed field `versions.creation_date`, use `versions.create_time` instead.

---

`alicloud_ots_instances`
- Removing the removed field `instances.network`, use `instances.network_type_acl` and `instances.network_source_acl` instead.
- Removing the removed field `instances.entity_quota`, use `instances.table_quota` instead.

---

`alicloud_pvtz_zones`
- Removing the removed fields `zones.creation_time`, `zones.update_time`.

---

`alicloud_ram_policies`
- Removing the removed field `policies.user_name`.

---

`alicloud_ram_users`
- Removing the removed field `users.last_login_date`.

---

`alicloud_slbs` / `alicloud_slb_load_balancers`
- Removing the removed field `master_availability_zone`, use `master_zone_id` instead.
- Removing the removed field `slave_availability_zone`, use `slave_zone_id` instead.

---

`alicloud_slb_master_slave_server_groups`
- Removing the removed field `groups.servers.is_backup`.
