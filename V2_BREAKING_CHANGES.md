
# Resources

---
`alicloud_cs_serverless_kubernetes` （Deprecated）
- Removing the removed field `vswitch_id`， use `vswitch_ids` instead.
- Removing the removed field `force_update`.
- Removing the removed field `create_v2_cluster`.

---
`alicloud_cs_edge_kubernetes` （Deprecated）
- Removing the removed field `force_update`.
- Change the type of `runtime`, `certificate_authority`, `connections` from `TypeMap` to `TypeList`

---
`alicloud_cs_kubernetes_node_pool`
- Removing the removed field `vpc_id`.
- Removing the removed field `rollout_policy`, use `rolling_policy` instead.

---
`alicloud_cs_kubernetes` （Deprecated）
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
- Removing the removed field `kube_config`, use resource `alicloud_cs_cluster_credential.output_file` instead.
- Removing the removed field `worker_nodes`, use resource `alicloud_cs_kubernetes_node_pool` instead.
- Removing the removed field `worker_instance_type`, use field `worker_instance_types` instead.
- Removing the removed field `vswitch_ids`, use field `master_vswitch_ids` instead.
- Removing the removed field `master_instance_type`, use field `master_instance_types` instead.
- Removing the removed field `force_update`.
- Removing the removed field `availability_zone`.
- Removing the removed field `vswitch_id`, use field `master_vswitch_ids`  instead.
- Removing the removed field `worker_numbers`, use field `worker_number` instead.
- Removing the removed field `nodes`, use field `master_nodes` instead.
- Removing the removed field `log_config`, use field `addons` instead.
- Removing the removed field `cluster_network_type`, use field `addons` instead.
- Removing the removed field `user_data`, use resource `alicloud_cs_kubernetes_node_pool.user_data` instead.
- Change the type of `runtime`, `certificate_authority`, `connections` from `TypeMap` to `TypeList`

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
- Removing the removed field `worker_numbers`, use field `worker_number` instead.
- Removing the removed field `cluster_network_type`, use field `addons` instead.
- Removing the removed field `log_config`, use field `addons` instead.
- Removing the removed field `worker_instance_type`, use field `worker_instance_types` instead.
- Change the type of `certificate_authority`, `connections` from `TypeMap` to `TypeList`

---
`alicloud_api_gateway_instance`
- Change the type of `to_connect_vpc_ip_block` from `TypeMap` to `TypeList`

---
`alicloud_cr_repos`: Change the type of `domain_list` from `TypeMap` to `TypeList`


# Data sources:

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
- Change the type of `connections` from `TypeMap` to `TypeList`

---
`alicloud_cr_repo`: Change the type of `domain_list` from `TypeMap` to `TypeList`

---
`alicloud_cs_cluster_credential`
- Change the type of `certificate_authority` from `TypeMap` to `TypeList`

---
`alicloud_cs_edge_kubernetes_clusters`
- Change the type of `connections` from `TypeMap` to `TypeList`

---
`alicloud_cs_kubernetes_clusters`
- Change the type of `connections` from `TypeMap` to `TypeList`

---
`alicloud_cs_serverless_kubernetes_clusters`
- Change the type of `connections` from `TypeMap` to `TypeList`

---
`alicloud_db_instance_classes`
- Change the type of `storage_range` from `TypeMap` to `TypeList`

---
`alicloud_instance_types`
- Change the type of `gpu`, `burstable_instance`, `local_storage` from `TypeMap` to `TypeList`