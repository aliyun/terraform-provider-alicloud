## 2.3.0 (Unreleased)

## 2.2.0 (November 5, 2026)

NOTES:

- resource/alicloud_dts_subscription_job and resource/alicloud_dts_synchronization_job: require the canonical `DtsJobId` for imports and job identity checks. If existing state contains a DTS instance ID, back up the state, verify the intended job, and re-import its `DtsJobId` as described in the resource documentation. ([#10412](https://github.com/aliyun/terraform-provider-alicloud/issues/10412))
- resource/alicloud_esa_image_transform: `enable` is now required. Set it explicitly to `on` or `off` when upgrading. ([#11176](https://github.com/aliyun/terraform-provider-alicloud/issues/11176))
- resource/alicloud_cloud_firewall_vpc_cen_tr_firewall: changes to `tr_attachment_master_zone` or `tr_attachment_slave_zone` now force resource replacement. ([#11295](https://github.com/aliyun/terraform-provider-alicloud/issues/11295))
- resource/alicloud_threat_detection_service_linked_role: default `service_linked_role` to `AliyunServiceRoleForSas` and refuse to adopt an existing role automatically. State created without an explicit role name on v1.283.0-v2.1.0 may require re-import; follow the migration instructions in the resource documentation before applying. ([#11455](https://github.com/aliyun/terraform-provider-alicloud/issues/11455))

FEATURES:

- **New Resource:** `alicloud_dts_check_job` ([#10244](https://github.com/aliyun/terraform-provider-alicloud/issues/10244))
- **New Resource:** `alicloud_direct_mail_dedicated_ip_pool` ([#10601](https://github.com/aliyun/terraform-provider-alicloud/issues/10601))
- **New Resource:** `alicloud_direct_mail_config_set` ([#10601](https://github.com/aliyun/terraform-provider-alicloud/issues/10601))
- **New Resource:** `alicloud_arms_notification_policy` ([#10694](https://github.com/aliyun/terraform-provider-alicloud/issues/10694))
- **New Resource:** `alicloud_arms_webhook_contact` ([#10694](https://github.com/aliyun/terraform-provider-alicloud/issues/10694))
- **New Resource:** `alicloud_tag_created_by` ([#11143](https://github.com/aliyun/terraform-provider-alicloud/issues/11143))
- **New Resource:** `alicloud_governance_service_activation` ([#11157](https://github.com/aliyun/terraform-provider-alicloud/issues/11157))
- **New Resource:** `alicloud_esa_site_environment` ([#11179](https://github.com/aliyun/terraform-provider-alicloud/issues/11179))
- **New Resource:** `alicloud_dts_gad_instance` ([#11378](https://github.com/aliyun/terraform-provider-alicloud/issues/11378))
- **New Data Source:** `alicloud_cms_metric_list` ([#10150](https://github.com/aliyun/terraform-provider-alicloud/issues/10150))
- **New Data Source:** `alicloud_direct_mail_dedicated_ip_pools` ([#10601](https://github.com/aliyun/terraform-provider-alicloud/issues/10601))
- **New Data Source:** `alicloud_direct_mail_config_sets` ([#10601](https://github.com/aliyun/terraform-provider-alicloud/issues/10601))
- **New Data Source:** `alicloud_rds_db_instance_performance` ([#10802](https://github.com/aliyun/terraform-provider-alicloud/issues/10802))
- **New Data Source:** `alicloud_oos_inventories` ([#11429](https://github.com/aliyun/terraform-provider-alicloud/issues/11429))
- **New Data Source:** `alicloud_rds_databases` ([#11453](https://github.com/aliyun/terraform-provider-alicloud/issues/11453))

ENHANCEMENTS:

- resource/alicloud_dts_migration_job: support `data_check_configure`. ([#10244](https://github.com/aliyun/terraform-provider-alicloud/issues/10244))
- resource/alicloud_polardb_cluster: support TDE automatic key rotation for PostgreSQL-compatible and Oracle syntax-compatible clusters. ([#10385](https://github.com/aliyun/terraform-provider-alicloud/issues/10385))
- resource/alicloud_kvstore_instance and resource/alicloud_redis_tair_instance: support Elastic Burst and per-shard additional bandwidth. ([#10584](https://github.com/aliyun/terraform-provider-alicloud/issues/10584), [#11372](https://github.com/aliyun/terraform-provider-alicloud/issues/11372))
- resource/alicloud_cms_addon_release: support `entity_rules` and deprecate `dry_run`. ([#10835](https://github.com/aliyun/terraform-provider-alicloud/issues/10835))
- resource/alicloud_vod_domain: support `ssl_protocol`, `cert_id`, `cert_name`, `cert_type`, `cert_region`, `ssl_pub` and `ssl_pri`. ([#11054](https://github.com/aliyun/terraform-provider-alicloud/issues/11054))
- resource/alicloud_nas_file_system: support `auto_upgrade_config` for CPFS SE file systems; data-source/alicloud_nas_file_systems: export `auto_upgrade_config`. ([#11085](https://github.com/aliyun/terraform-provider-alicloud/issues/11085))
- resource/alicloud_ecs_disk and resource/alicloud_disk: support `force_delete` to convert and release PrePaid data disks on destroy. ([#11108](https://github.com/aliyun/terraform-provider-alicloud/issues/11108))
- resource/alicloud_rds_backup: export the create-time `backup_job_id`. ([#11109](https://github.com/aliyun/terraform-provider-alicloud/issues/11109))
- resource/alicloud_pvtz_zone_record: support AAAA (IPv6) records. ([#11138](https://github.com/aliyun/terraform-provider-alicloud/issues/11138))
- resource/alicloud_polardb_backup_policy: support advanced backup policies and custom advanced data backup rules. Switching from `Normal` to `Advanced` is irreversible. ([#11149](https://github.com/aliyun/terraform-provider-alicloud/issues/11149))
- resource/alicloud_esa_image_transform: support `auto_avif` and `auto_webp`, and export `config_type`. ([#11176](https://github.com/aliyun/terraform-provider-alicloud/issues/11176))
- resource/alicloud_db_instance: support `force_encryption` for PostgreSQL. ([#11223](https://github.com/aliyun/terraform-provider-alicloud/issues/11223))
- resource/alicloud_dts_synchronization_instance: support updating `auto_pay` for existing subscription instances. ([#11224](https://github.com/aliyun/terraform-provider-alicloud/issues/11224))
- resource/alicloud_privatelink_vpc_endpoint_connection: support `traffic_control_mode` and export `endpoint_owner_id`; data-source/alicloud_privatelink_vpc_endpoint_connections: export both attributes. ([#11277](https://github.com/aliyun/terraform-provider-alicloud/issues/11277))
- resource/alicloud_cloud_firewall_vpc_cen_tr_firewall: support `firewall_vpc_id`, `firewall_vswitch_id`, and the `manual` and `semi-managed` route modes; make `firewall_vpc_cidr` optional and deprecate `firewall_subnet_cidr`, `tr_attachment_master_cidr` and `tr_attachment_slave_cidr`. ([#11295](https://github.com/aliyun/terraform-provider-alicloud/issues/11295))
- resource/alicloud_ecd_desktop: support `desktop_timers`; resource/alicloud_ecd_policy_group: support `internet_communication_protocol`, `recording_duration`, `client_types`, `end_user_apply_admin_coordinate` and `end_user_group_coordinate`. ([#11345](https://github.com/aliyun/terraform-provider-alicloud/issues/11345))
- resource/alicloud_realtime_compute_vvp_instance: support `default_namespace_resource_spec`. ([#11392](https://github.com/aliyun/terraform-provider-alicloud/issues/11392))
- resource/alicloud_kvstore_instance and resource/alicloud_redis_tair_instance: support write-only `password_wo` and `password_wo_version`. ([#11409](https://github.com/aliyun/terraform-provider-alicloud/issues/11409))
- resource/alicloud_mongodb_instance and resource/alicloud_mongodb_sharding_instance: support write-only `account_password_wo` and `account_password_wo_version`. ([#11410](https://github.com/aliyun/terraform-provider-alicloud/issues/11410))
- resource/alicloud_db_database: support `account_name`, `account_privilege` and `collation_name`. ([#11411](https://github.com/aliyun/terraform-provider-alicloud/issues/11411))
- resource/alicloud_elasticsearch_instance: support write-only `password_wo` and `password_wo_version`. ([#11412](https://github.com/aliyun/terraform-provider-alicloud/issues/11412))
- resource/alicloud_elasticsearch_instance: support `search_node_configuration` and `index_node_configuration`. ([#11428](https://github.com/aliyun/terraform-provider-alicloud/issues/11428))
- resource/alicloud_instance: support the `cloud_essd_xc0` and `cloud_essd_xc1` disk categories. ([#11436](https://github.com/aliyun/terraform-provider-alicloud/issues/11436))
- docs: use `maxmemory-policy` in the alicloud_kvstore_instance example configurations. ([#11438](https://github.com/aliyun/terraform-provider-alicloud/issues/11438))
- docs: correct the ID format and import example for alicloud_cen_transit_router_route_table_propagation. ([#11464](https://github.com/aliyun/terraform-provider-alicloud/issues/11464))
- resource/alicloud_wafv3_domain: support `redirect.backend_ports`. ([#11466](https://github.com/aliyun/terraform-provider-alicloud/issues/11466))

BUG FIXES:

- resource/alicloud_alb_acl, resource/alicloud_alb_acl_entry_attachment and resource/alicloud_alb_listener_acl_attachment: retry `AclStatusNotStable` errors. ([#10157](https://github.com/aliyun/terraform-provider-alicloud/issues/10157))
- resource/alicloud_ess_scaling_configuration: preserve out-of-band `user_data` and tags, reconcile tag removal, and reject negative IPv6 address counts during planning. ([#10343](https://github.com/aliyun/terraform-provider-alicloud/issues/10343))
- resource/alicloud_dts_subscription_job and resource/alicloud_dts_synchronization_job: use the canonical job ID when updating job names and reject mismatched job identities. ([#10412](https://github.com/aliyun/terraform-provider-alicloud/issues/10412))
- resource/alicloud_cs_kubernetes and data-source/alicloud_cs_cluster_credential: return kubeconfig parsing errors instead of panicking, and propagate state-setting errors. ([#10596](https://github.com/aliyun/terraform-provider-alicloud/issues/10596))
- resource/alicloud_apig_plugin: fix read errors when using the AI gateway. ([#11081](https://github.com/aliyun/terraform-provider-alicloud/issues/11081))
- resource/alicloud_nas_file_system: wait for asynchronously persisted `description` and `resource_group_id` values after creation. ([#11085](https://github.com/aliyun/terraform-provider-alicloud/issues/11085))
- resource/alicloud_privatelink_vpc_endpoint_service: include `zone_id` in the resource set hash. ([#11104](https://github.com/aliyun/terraform-provider-alicloud/issues/11104))
- resource/alicloud_redis_tair_instance: retry `SocketTimeout.Exception` errors. ([#11105](https://github.com/aliyun/terraform-provider-alicloud/issues/11105))
- resource/alicloud_ess_eci_scaling_configuration: fix empty-string defaults in environment variable attributes. ([#11107](https://github.com/aliyun/terraform-provider-alicloud/issues/11107))
- resource/alicloud_rds_backup: correctly parse the composite resource ID when deleting a backup. ([#11109](https://github.com/aliyun/terraform-provider-alicloud/issues/11109))
- resource/alicloud_rocketmq_instance: send `trace_on` only when explicitly configured. ([#11110](https://github.com/aliyun/terraform-provider-alicloud/issues/11110))
- resource/alicloud_tag_created_by: propagate the lifecycle context through retries. ([#11157](https://github.com/aliyun/terraform-provider-alicloud/issues/11157))
- resource/alicloud_alikafka_topic: avoid perpetual `configs` diffs caused by server-managed keys and JSON value representation. ([#11158](https://github.com/aliyun/terraform-provider-alicloud/issues/11158))
- provider: restore SLS endpoint cache reuse across API calls. ([#11208](https://github.com/aliyun/terraform-provider-alicloud/issues/11208))
- resource/alicloud_oos_service_setting: correctly detect OSS delivery changes, send the complete delivery configuration on update, and avoid a read panic. ([#11210](https://github.com/aliyun/terraform-provider-alicloud/issues/11210))
- resource/alicloud_esa_certificate: fix invalid certificate type errors during creation. ([#11215](https://github.com/aliyun/terraform-provider-alicloud/issues/11215))
- resource/alicloud_nlb_server_group: replace the invalid `TCPSSL` protocol value with `TCP_UDP`. ([#11222](https://github.com/aliyun/terraform-provider-alicloud/issues/11222))
- resource/alicloud_ecd_policy_group: prevent read drift for policy attributes and configured `client_types`. ([#11345](https://github.com/aliyun/terraform-provider-alicloud/issues/11345))
- resource/alicloud_polardb_cluster: make `encrypt_new_tables` Optional and Computed to prevent plan drift. ([#11426](https://github.com/aliyun/terraform-provider-alicloud/issues/11426))
- provider: honor `NO_PROXY` for tea-roa clients and preserve endpoint ports when bypassing a proxy. ([#11437](https://github.com/aliyun/terraform-provider-alicloud/issues/11437))
- resource/alicloud_threat_detection_service_linked_role: reconcile role status with Security Center and prevent unintended adoption of existing roles. ([#11455](https://github.com/aliyun/terraform-provider-alicloud/issues/11455))
- provider: restore skip-aware validation for DTS, IMS and PAI Workspace resources. ([#11459](https://github.com/aliyun/terraform-provider-alicloud/issues/11459))
- provider: preserve configured client timeouts when rebuilding the SLS SDK configuration. ([#11461](https://github.com/aliyun/terraform-provider-alicloud/issues/11461))
- provider: honor configured client timeouts in Function Compute and OSS v2 clients. ([#11463](https://github.com/aliyun/terraform-provider-alicloud/issues/11463))
- resource/alicloud_vpc: make `system_route_table_description` and `system_route_table_name` Optional and Computed to prevent persistent diffs when omitted. ([#11468](https://github.com/aliyun/terraform-provider-alicloud/issues/11468))
- resource/alicloud_dts_gad_instance: use the configured provider client through ResourceBase metadata. ([#11476](https://github.com/aliyun/terraform-provider-alicloud/issues/11476))
- provider: derive legacy SDK timeouts from provider settings and preserve defaults when timeout settings are omitted. ([#11477](https://github.com/aliyun/terraform-provider-alicloud/issues/11477))
- provider: refresh `credentials_uri` credentials before expiration. ([#11479](https://github.com/aliyun/terraform-provider-alicloud/issues/11479))
- data-source/alicloud_cs_kubernetes_addons: preserve error context; resource/alicloud_polardb_endpoint: prevent endpoint-specific retry rules from affecting other operations. ([#11483](https://github.com/aliyun/terraform-provider-alicloud/issues/11483))

## 2.1.0 (October 22, 2026)

FEATURES:

- **New Resource:** `alicloud_polardb_gateway_model_service` ([#10582](https://github.com/aliyun/terraform-provider-alicloud/issues/10582))
- **New Resource:** `alicloud_polardb_gateway_model_api` ([#10582](https://github.com/aliyun/terraform-provider-alicloud/issues/10582))
- **New Resource:** `alicloud_polardb_gateway_consumer_group` ([#10585](https://github.com/aliyun/terraform-provider-alicloud/issues/10585))
- **New Resource:** `alicloud_polardb_gateway_consumer` ([#10585](https://github.com/aliyun/terraform-provider-alicloud/issues/10585))
- **New Resource:** `alicloud_polardb_gateway_budget_policy` ([#10586](https://github.com/aliyun/terraform-provider-alicloud/issues/10586))
- **New Resource:** `alicloud_polardb_gateway_rate_limit_policy` ([#10586](https://github.com/aliyun/terraform-provider-alicloud/issues/10586))
- **New Resource:** `alicloud_polardb_gateway_cost_rule` ([#10592](https://github.com/aliyun/terraform-provider-alicloud/issues/10592))
- **New Resource:** `alicloud_polardb_polar_fs` ([#10623](https://github.com/aliyun/terraform-provider-alicloud/issues/10623))
- **New Resource:** `alicloud_polardb_polar_fs_endpoint` ([#10623](https://github.com/aliyun/terraform-provider-alicloud/issues/10623))
- **New Resource:** `alicloud_ecs_snapshot_lock` ([#10773](https://github.com/aliyun/terraform-provider-alicloud/issues/10773))
- **New Resource:** `alicloud_oss_access_point_policy` ([#11076](https://github.com/aliyun/terraform-provider-alicloud/issues/11076))
- **New Data Source:** `alicloud_polardb_gateway_model_services` ([#10582](https://github.com/aliyun/terraform-provider-alicloud/issues/10582))
- **New Data Source:** `alicloud_polardb_gateway_model_apis` ([#10582](https://github.com/aliyun/terraform-provider-alicloud/issues/10582))
- **New Data Source:** `alicloud_polardb_gateway_consumer_groups` ([#10585](https://github.com/aliyun/terraform-provider-alicloud/issues/10585))
- **New Data Source:** `alicloud_polardb_gateway_consumers` ([#10585](https://github.com/aliyun/terraform-provider-alicloud/issues/10585))
- **New Data Source:** `alicloud_polardb_gateway_budget_policies` ([#10586](https://github.com/aliyun/terraform-provider-alicloud/issues/10586))
- **New Data Source:** `alicloud_polardb_gateway_rate_limit_policies` ([#10586](https://github.com/aliyun/terraform-provider-alicloud/issues/10586))
- **New Data Source:** `alicloud_polardb_gateway_cost_rules` ([#10592](https://github.com/aliyun/terraform-provider-alicloud/issues/10592))
- **New Data Source:** `alicloud_polardb_polar_fs_instances` ([#10623](https://github.com/aliyun/terraform-provider-alicloud/issues/10623))
- **New Data Source:** `alicloud_polardb_lakebase_tenant_token` ([#10623](https://github.com/aliyun/terraform-provider-alicloud/issues/10623))
- **New Data Source:** `alicloud_sae_application_groups` ([#10929](https://github.com/aliyun/terraform-provider-alicloud/issues/10929))
- **New Data Source:** `alicloud_sae_application_instances` ([#11039](https://github.com/aliyun/terraform-provider-alicloud/issues/11039))
- **New Ephemeral Resource:** `alicloud_sts_token` — fetches temporary access credentials through AssumeRole without persisting them to state or plan. Requires Terraform 1.10 or later. Served by the framework provider. ([#10788](https://github.com/aliyun/terraform-provider-alicloud/issues/10788))
- **New Ephemeral Resource:** `alicloud_ram_access_key` — creates a temporary RAM access key for use during Terraform operations and deletes it when the ephemeral resource closes. Requires Terraform 1.10 or later. Served by the framework provider. ([#11047](https://github.com/aliyun/terraform-provider-alicloud/issues/11047))

ENHANCEMENTS:

- resource/alicloud_mse_gateway: support resource_group_id and tags. ([#10358](https://github.com/aliyun/terraform-provider-alicloud/issues/10358))
- resource/alicloud_dts_subscription_job: support source_endpoint_ssl and fail fast when prechecks fail. ([#10490](https://github.com/aliyun/terraform-provider-alicloud/issues/10490))
- resource/alicloud_pai_workspace_experiment: support labels. ([#10698](https://github.com/aliyun/terraform-provider-alicloud/issues/10698))
- resource/alicloud_message_service_topic: support enable_sse, sse_type, sse_algorithm and kms_key_id. ([#10706](https://github.com/aliyun/terraform-provider-alicloud/issues/10706))
- resource/alicloud_kvstore_instance and resource/alicloud_redis_tair_instance: export the CA certificate when SSL is enabled. ([#10757](https://github.com/aliyun/terraform-provider-alicloud/issues/10757))
- resource/alicloud_gpdb_instance: support effective_time and sql_collector_status. ([#10763](https://github.com/aliyun/terraform-provider-alicloud/issues/10763))
- docs: update the README for the 2.x provider line. ([#10796](https://github.com/aliyun/terraform-provider-alicloud/issues/10796))
- resource/alicloud_kvstore_instance: support global_security_group_list. ([#10798](https://github.com/aliyun/terraform-provider-alicloud/issues/10798))
- resource/alicloud_elasticsearch_instance: support in-place resource_group_id updates. ([#10805](https://github.com/aliyun/terraform-provider-alicloud/issues/10805))
- docs: align alicloud_ecs_capacity_reservation end_time semantics with the API and document platform as ForceNew. ([#10807](https://github.com/aliyun/terraform-provider-alicloud/issues/10807))
- resource/alicloud_instance: support managed_host_id and configurable additional_attributes. ([#10816](https://github.com/aliyun/terraform-provider-alicloud/issues/10816))
- resource/alicloud_cms_integration_policy: support entity_group.entity_user_id, entity_group.cluster_namespace and entity_group.disable_policy_share. ([#10823](https://github.com/aliyun/terraform-provider-alicloud/issues/10823))
- resource/alicloud_dts_migration_instance and resource/alicloud_dts_migration_job: support MariaDB source and destination engines with case-insensitive validation. ([#10825](https://github.com/aliyun/terraform-provider-alicloud/issues/10825))
- resource/alicloud_privatelink_vpc_endpoint: allow endpoint_type Reverse. ([#10826](https://github.com/aliyun/terraform-provider-alicloud/issues/10826))
- resource/alicloud_pai_service: support service_uid. ([#10827](https://github.com/aliyun/terraform-provider-alicloud/issues/10827))
- resource/alicloud_rds_clone_db_instance: apply tags atomically during CloneDBInstance creation. ([#10828](https://github.com/aliyun/terraform-provider-alicloud/issues/10828))
- docs: document the HTTPS coordinating node prerequisite for alicloud_elasticsearch_instance. ([#10832](https://github.com/aliyun/terraform-provider-alicloud/issues/10832))
- data-source/alicloud_oss_bucket_objects: export x_oss_hash_crc64ecma. ([#10859](https://github.com/aliyun/terraform-provider-alicloud/issues/10859))
- docs: fix NOTE callout placement so attribute guidance is attached to the relevant field. ([#11030](https://github.com/aliyun/terraform-provider-alicloud/issues/11030))
- resource/alicloud_slb_acl_entry_attachment: support batch ACL entries. ([#11055](https://github.com/aliyun/terraform-provider-alicloud/issues/11055))
- docs: add an alicloud_privatelink_vpc_endpoint example using vpc_endpoint_zone. ([#11063](https://github.com/aliyun/terraform-provider-alicloud/issues/11063))
- docs: document the requirement to configure at least one delivery target for alicloud_actiontrail_trail. ([#11075](https://github.com/aliyun/terraform-provider-alicloud/issues/11075))
- docs: pin the region in the alicloud_privatelink_vpc_endpoint endpoint zone example to match the endpoint service. ([#11093](https://github.com/aliyun/terraform-provider-alicloud/issues/11093))

BUG FIXES:

- resource/alicloud_dts_synchronization_job: reconcile instance_class against the actual job class; resource/alicloud_dts_synchronization_instance: accept the full class set and document instance_class as create-only. ([#10173](https://github.com/aliyun/terraform-provider-alicloud/issues/10173))
- data-source/alicloud_polardb_node_classes: normalize db_type and db_version engine matching. ([#10280](https://github.com/aliyun/terraform-provider-alicloud/issues/10280))
- resource/alicloud_ecs_disk: confirm category conversion before enabling bursting. ([#10303](https://github.com/aliyun/terraform-provider-alicloud/issues/10303))
- resource/alicloud_oss_access_point: clear stale vpc_configuration for non-VPC access points to avoid unnecessary replacement. ([#10484](https://github.com/aliyun/terraform-provider-alicloud/issues/10484))
- resource/alicloud_elasticsearch_instance: prevent a provider crash when parsing Kibana private network responses. ([#10600](https://github.com/aliyun/terraform-provider-alicloud/issues/10600))
- data-source/alicloud_express_connect_router_express_connect_routers: add the missing tags element schema. ([#10711](https://github.com/aliyun/terraform-provider-alicloud/issues/10711))
- data-source/alicloud_arms_environments: add the missing tags element schema and refine read handling. ([#10725](https://github.com/aliyun/terraform-provider-alicloud/issues/10725))
- resource/alicloud_click_house_db_cluster: avoid a panic when cold storage State is missing. ([#10729](https://github.com/aliyun/terraform-provider-alicloud/issues/10729))
- resource/alicloud_ots_instance: add the missing tags element schema. ([#10753](https://github.com/aliyun/terraform-provider-alicloud/issues/10753))
- resource/alicloud_message_service_queue: add the missing tags element schema. ([#10764](https://github.com/aliyun/terraform-provider-alicloud/issues/10764))
- resource/alicloud_nlb_server_group: suppress order-only health_check_http_code diffs. ([#10785](https://github.com/aliyun/terraform-provider-alicloud/issues/10785))
- provider: resolve credential profiles separately for each provider instance to prevent profile leakage between aliases. ([#10800](https://github.com/aliyun/terraform-provider-alicloud/issues/10800))
- provider: report credential profile read and parse errors during configuration instead of silently ignoring them; validate profile field types to prevent panics. ([#10801](https://github.com/aliyun/terraform-provider-alicloud/issues/10801))
- resource/alicloud_dts_migration_job: wait for start and suspend operations to reach their target states during updates, and report precheck or migration failures. ([#10825](https://github.com/aliyun/terraform-provider-alicloud/issues/10825))
- resource/alicloud_db_instance and resource/alicloud_ddoscoo_instance: return an error when destroying PrePaid instances instead of silently removing them from Terraform state. ([#10833](https://github.com/aliyun/terraform-provider-alicloud/issues/10833))
- resource/alicloud_adb_db_cluster_lake_version, resource/alicloud_cloud_storage_gateway_gateway, resource/alicloud_ecp_instance, resource/alicloud_ecs_dedicated_host, resource/alicloud_ga_bandwidth_package, resource/alicloud_ga_basic_accelerator, resource/alicloud_hologram_instance, resource/alicloud_mongodb_serverless_instance, resource/alicloud_mongodb_sharding_instance, resource/alicloud_nat_gateway, resource/alicloud_ocean_base_instance, resource/alicloud_realtime_compute_vvp_instance, resource/alicloud_rds_ddr_instance, resource/alicloud_reserved_instance, resource/alicloud_simple_application_server_instance, resource/alicloud_vpn_gateway and resource/alicloud_yundun_dbaudit_instance: return an error when destroying PrePaid instances instead of silently removing them from Terraform state. ([#10834](https://github.com/aliyun/terraform-provider-alicloud/issues/10834))
- resource/alicloud_security_group_rule: return readable errors for malformed import IDs instead of panicking. ([#10837](https://github.com/aliyun/terraform-provider-alicloud/issues/10837))
- resource/alicloud_vpn_gateway_vpn_attachment: preserve explicit tunnel updates and unknown values; resource/alicloud_kvstore_instance: handle missing EngineVersion values. ([#10857](https://github.com/aliyun/terraform-provider-alicloud/issues/10857))
- provider: return errors for invalid API value conversions instead of panicking. ([#10884](https://github.com/aliyun/terraform-provider-alicloud/issues/10884))
- resource/alicloud_redis_global_security_ip_group: fix read errors. ([#10907](https://github.com/aliyun/terraform-provider-alicloud/issues/10907))
- resource/alicloud_router_interface: preserve deletion errors; resource/alicloud_express_connect_router_interface: wait for deletion; resource/alicloud_vpc: retry router interface dependency errors. ([#11037](https://github.com/aliyun/terraform-provider-alicloud/issues/11037))
- resource/alicloud_cloud_firewall_instance_v2: release subscription instances on destroy and reject unsupported billing changes during planning. ([#11060](https://github.com/aliyun/terraform-provider-alicloud/issues/11060))
- resource/alicloud_cdn_domain_config: retry FlowControlError responses. ([#11070](https://github.com/aliyun/terraform-provider-alicloud/issues/11070))
- resource/alicloud_cs_kubernetes and Kubernetes cluster data sources: tolerate nodes with missing IP addresses instead of panicking. ([#11087](https://github.com/aliyun/terraform-provider-alicloud/issues/11087))
- provider: add the regional eflo_controller endpoint fallback. ([#11094](https://github.com/aliyun/terraform-provider-alicloud/issues/11094))

## 2.0.0 (October 12, 2026)

This is the first stable release of the 2.x line. It rolls up every change merged from the 1.x line since v2.0.0-beta5 — no additional 1.x section has been cut since `1.294.0` (`1.295.0` remains unreleased) — plus the v2-only changes below. Changes that merged to the 1.x line after the v2.0.0-beta5 cut are also listed below; with the 1.x line no longer maintained after this release, they will never appear in a released 1.x section.

FEATURES:

- **New Ephemeral Resource:** `alicloud_cr_ee_authorization_token` — fetches a temporary Container Registry Enterprise Edition login credential through the `GetAuthorizationToken` API without persisting it to the state or plan, for use in provider configuration, locals, provisioners, and other ephemeral resources. Requires Terraform 1.10 or later. Served by the framework provider. ([#10776](https://github.com/aliyun/terraform-provider-alicloud/issues/10776))
- **New Resource:** `alicloud_cloud_firewall_fw_switch`. ([#10748](https://github.com/aliyun/terraform-provider-alicloud/issues/10748))

ENHANCEMENTS:

- resource/alicloud_nas_file_system: support Agentic storage type for standard file systems; skip general-purpose NAS recycle bin/ACL calls on Agentic file systems. ([#10774](https://github.com/aliyun/terraform-provider-alicloud/issues/10774))
- resource/alicloud_ebs_disk_replica_group: support enable_rtc and updatable rpo; data-source/alicloud_ebs_disk_replica_groups: export enable_rtc. ([#10775](https://github.com/aliyun/terraform-provider-alicloud/issues/10775))
- docs: add MySQL 8.4 version to RDS documentation. ([#10722](https://github.com/aliyun/terraform-provider-alicloud/issues/10722))

BUG FIXES:

- resource/alicloud_nlb_server_group: preserve health_check_http_code order in read. ([#10762](https://github.com/aliyun/terraform-provider-alicloud/issues/10762))

## 2.0.0-beta5 (October 10, 2026)

This beta rolls up every change merged from the 1.x line since v2.0.0-beta4 — see the `1.291.0`, `1.292.0`, `1.293.0` and `1.294.0` sections of [CHANGELOG.md](CHANGELOG.md) — plus the v2-only changes below. Changes that merged to the 1.x line after the `1.294.0` changelog cut are also listed below, ahead of their appearance in any released 1.x section.

FEATURES:

- **New Ephemeral Resource:** `alicloud_kms_secret` — reads a KMS secret's value through the `GetSecretValue` API without persisting it to the state or plan, for use in provider configuration, write-only attributes, and other ephemeral resources. Requires Terraform 1.10 or later. Served by the framework provider. ([#10650](https://github.com/aliyun/terraform-provider-alicloud/issues/10650))
- **New Ephemeral Resource:** `alicloud_oos_secret_parameter` — reads an OOS secret parameter's decrypted value through the `GetSecretParameter` API without persisting it to the state or plan, for use in provider configuration, write-only attributes, and other ephemeral resources. Requires Terraform 1.10 or later. Served by the framework provider. ([#10662](https://github.com/aliyun/terraform-provider-alicloud/issues/10662))
- **New Ephemeral Resource:** `alicloud_kms_plaintext` — decrypts a KMS ciphertext through the `Decrypt` API without persisting the plaintext to the state or plan, for use in provider configuration, write-only attributes, and other ephemeral resources. Requires Terraform 1.10 or later. Served by the framework provider. ([#10677](https://github.com/aliyun/terraform-provider-alicloud/issues/10677))
- **New Ephemeral Resource:** `alicloud_cs_cluster_credential` — fetches an ACK cluster kubeconfig through the `DescribeClusterUserKubeconfig` API without persisting the credential to the state or plan, for use in provider configuration, locals, provisioners, and other ephemeral resources. Requires Terraform 1.10 or later. Served by the framework provider. ([#10747](https://github.com/aliyun/terraform-provider-alicloud/issues/10747))
- **New Action:** `alicloud_ssl_certificates_service_pca_cert_sync` — synchronizes a batch of existing PCA (private CA) client certificates to the SSL Certificates service in one shot, the same operation as "Batch Synchronize to SSL Certificates" in the PCA console. Requires Terraform 1.14 or later. Served by the framework provider. ([#10567](https://github.com/aliyun/terraform-provider-alicloud/issues/10567))
- **New Resource:** `alicloud_ims_user` — the IMS user resource, served by the framework provider. ([#10515](https://github.com/aliyun/terraform-provider-alicloud/issues/10515))

ENHANCEMENTS:

- resource/alicloud_oos_parameter: support the write-only attribute `value_wo` with `value_wo_version` as its change trigger, so the parameter value can be managed without being persisted in the state or plan. ([#10505](https://github.com/aliyun/terraform-provider-alicloud/issues/10505))
- resource/alicloud_kms_ciphertext: support the write-only attribute `plaintext_wo` with `plaintext_wo_version` as its change trigger, so the plaintext can be managed without being persisted in the state or plan. ([#10598](https://github.com/aliyun/terraform-provider-alicloud/issues/10598))
- resource/alicloud_ssl_certificates_service_certificate: support the write-only attribute `key_wo` with `key_wo_version` as its change trigger, so the certificate private key can be managed without being persisted in the state or plan. ([#10602](https://github.com/aliyun/terraform-provider-alicloud/issues/10602))
- resource/alicloud_rds_account: support the write-only attribute `account_password_wo` with `account_password_wo_version` as its change trigger, so the account password can be managed without being persisted in the state or plan. ([#10615](https://github.com/aliyun/terraform-provider-alicloud/issues/10615))
- resource/alicloud_gpdb_account: support the write-only attribute `account_password_wo` with `account_password_wo_version` as its change trigger, so the account password can be managed without being persisted in the state or plan. ([#10619](https://github.com/aliyun/terraform-provider-alicloud/issues/10619))
- resource/alicloud_polardb_account: support the write-only attribute `account_password_wo` with `account_password_wo_version` as its change trigger, so the account password can be managed without being persisted in the state or plan. ([#10636](https://github.com/aliyun/terraform-provider-alicloud/issues/10636))
- resource/alicloud_kvstore_account: support the write-only attribute `account_password_wo` with `account_password_wo_version` as its change trigger, so the account password can be managed without being persisted in the state or plan. ([#10637](https://github.com/aliyun/terraform-provider-alicloud/issues/10637))
- resource/alicloud_nlb_server_group: support ip_version_affinity_mode; data-source/alicloud_nlb_server_groups: support ip_version_affinity_mode. ([#10708](https://github.com/aliyun/terraform-provider-alicloud/issues/10708))
- data-source/alicloud_cloud_control_prices: add desire_attributes_json to support nested pricing attributes. ([#10742](https://github.com/aliyun/terraform-provider-alicloud/issues/10742))
- resource/alicloud_apig_policy: retry on transient 409 Conflict.DomainIsProcessing. ([#10736](https://github.com/aliyun/terraform-provider-alicloud/issues/10736))
- provider: support HTTP proxy for TableStore clients.
- resource/alicloud_image_import: migrate nvme_support from the disk_device_mapping block to the features block, updating it in place via ModifyImageAttribute instead of forcing a re-import. ([#10607](https://github.com/aliyun/terraform-provider-alicloud/issues/10607))
- docs: clarify auto_repair requires node-lifecycle-controller and ack-node-problem-detector addons for alicloud_cs_kubernetes_node_pool. ([#10741](https://github.com/aliyun/terraform-provider-alicloud/issues/10741))

BUG FIXES:

- resource/alicloud_oss_bucket_website: fix create failing with ResourceNotfound due to delayed config visibility. ([#10734](https://github.com/aliyun/terraform-provider-alicloud/issues/10734))
- resource/alicloud_nas_file_system: suppress vpc_id/vswitch_id replace diff; resource/alicloud_nas_mount_target: mark vswitch_id Computed. ([#10461](https://github.com/aliyun/terraform-provider-alicloud/issues/10461))
- resource/alicloud_vpn_gateway_vpn_attachment: detect nested tunnel updates while preserving set semantics. ([#10523](https://github.com/aliyun/terraform-provider-alicloud/issues/10523))
- resource/alicloud_cloud_firewall_address_book: prevent address list ordering drift. ([#10750](https://github.com/aliyun/terraform-provider-alicloud/issues/10750))
- resource/alicloud_threat_detection_check_config: prevent selected checks ordering drift. ([#10745](https://github.com/aliyun/terraform-provider-alicloud/issues/10745))
- resource/alicloud_gpdb_instance, resource/alicloud_click_house_db_cluster and resource/alicloud_mse_cluster: return an error when destroying PrePaid instances instead of silently removing them from Terraform state. ([#10740](https://github.com/aliyun/terraform-provider-alicloud/issues/10740))
- resource/alicloud_slb_load_balancer, resource/alicloud_eip_address and resource/alicloud_alikafka_instance: return an error when destroying PrePaid instances instead of silently removing them from Terraform state. ([#10746](https://github.com/aliyun/terraform-provider-alicloud/issues/10746))
- provider: the `1.294.0` kubeconfig error propagation and nil guards for resource/alicloud_cs_kubernetes, resource/alicloud_cs_managed_kubernetes and resource/alicloud_cs_serverless_kubernetes were reverted on the 1.x line and are not part of this beta. ([#10755](https://github.com/aliyun/terraform-provider-alicloud/issues/10755))

## 2.0.0-beta4 (August 27, 2026)

This beta rolls up every change merged from the 1.x line since v2.0.0-beta3 — see the `1.290.0` section of [CHANGELOG.md](CHANGELOG.md) — plus the v2-only changes below.

BREAKING CHANGES:

- resource/alicloud_amqp_exchange: `create_time` changed from `TypeInt` to `TypeString`, fixing a provider panic on 32-bit builds. No configuration or state migration is needed; expressions that consumed it as a number now need `tonumber`. ([#10254](https://github.com/aliyun/terraform-provider-alicloud/issues/10254))
- resource/alicloud_amqp_instance: `create_time` changed from `TypeInt` to `TypeString`, fixing a provider panic on 32-bit builds. No configuration or state migration is needed; expressions that consumed it as a number now need `tonumber`. ([#10254](https://github.com/aliyun/terraform-provider-alicloud/issues/10254))
- resource/alicloud_amqp_static_account: `create_time` changed from `TypeInt` to `TypeString`, fixing a provider panic on 32-bit builds. No configuration or state migration is needed; expressions that consumed it as a number now need `tonumber`. ([#10254](https://github.com/aliyun/terraform-provider-alicloud/issues/10254))

ENHANCEMENTS:

- provider: move ProtoV5ProviderServerFactory from package framework into the provider subpackage, resolving a cyclic import between the `alicloud` root and `framework` packages. ([#10269](https://github.com/aliyun/terraform-provider-alicloud/issues/10269))
- provider: migrate `helper/resource` retry helpers to the `helper/retry` subpackage, removing the last SDK v1 retry dependency from the provider. ([#10255](https://github.com/aliyun/terraform-provider-alicloud/issues/10255), [#10266](https://github.com/aliyun/terraform-provider-alicloud/issues/10266))

## 2.0.0-beta3 (August 19, 2026)

This beta rolls up every change merged from the 1.x line since v2.0.0-beta2 — see the `1.288.0` and `1.289.0` sections of [CHANGELOG.md](CHANGELOG.md) — plus the v2-only changes below.

FEATURES:

- **Terraform Plugin Framework:** the provider binary now serves a `terraform-plugin-framework` provider alongside the existing `terraform-plugin-sdk/v2` one, muxed behind `tf5muxserver`. Protocol version, provider address and provider schema are unchanged; new resources and data sources can be authored framework-native, and the framework-only concept types (ephemeral resources, provider functions, list resources, actions) become available for the first time.
- **New Data Source:** `alicloud_ims_default_domain` — the default domain of the Alibaba Cloud account, backed by the Ims `GetDefaultDomain` API. Served by the framework provider.
- **New Function:** `arn_build` — builds an Alibaba Cloud Resource Name (ARN) of the form `acs:<service>:<region>:<account_id>:<resource>` from its constituent parts, invoked as `provider::alicloud::arn_build(service, region, account_id, resource)`. Mirrors the AWS provider's `arn_build` function. Requires Terraform 1.8+ with the provider declared in `required_providers`.
- **New Function:** `arn_parse` — parses an Alibaba Cloud Resource Name (ARN) of the form `acs:<service>:<region>:<account_id>:<resource>` into an object with `service`, `region`, `account_id` and `resource` attributes, invoked as `provider::alicloud::arn_parse(arn)`. The attributes are the arguments of `arn_build`, so the two compose. Any colons beyond the fourth belong to the resource and are kept as part of it. Unlike `arn_build`, which formats whatever it is given, `arn_parse` reports an error for a string that is not an ARN. Mirrors the AWS provider's `arn_parse` function. Requires Terraform 1.8+ with the provider declared in `required_providers`.

ENHANCEMENTS:

- deps: bump `golang.org/x/crypto` to 0.52.0 and `google.golang.org/grpc` to 1.82.1.

BUG FIXES:

- provider: close the partial-state window before the success return in resource/alicloud_reserved_instance and resource/alicloud_cs_kubernetes_node_pool. SDK v2 deleted the `d.SetPartial` whitelist that `d.Partial(true)` filtered against, inverting the flag's meaning from "persist only these keys" to "persist nothing", so a window still open at a successful return wrote prior state instead of what `d.Set` staged: an `alicloud_reserved_instance` attribute update was discarded, and an `alicloud_cs_kubernetes_node_pool` scale-down reported success while leaving the old `node_count` in state and re-planning the same diff. Eleven `Partial` calls that cannot affect any outcome have also been removed. ([#10210](https://github.com/aliyun/terraform-provider-alicloud/issues/10210))

## 2.0.0-beta2 (August 4, 2026)

This beta rolls up every change merged from the 1.x line since v2.0.0-beta1 — see the `1.287.0` section of [CHANGELOG.md](CHANGELOG.md) — plus the v2-only fix below.

BUG FIXES:

- provider: report the actual provider version in the API `User-Agent`. `.goreleaser.yml` injected `-X main.version`, but `main.go` declares no `version` variable, so the linker discarded the value silently and v2.0.0-beta1 identified itself to Alibaba Cloud as `Terraform-Provider/1.286.0`. The ldflag now targets `alicloud/connectivity.providerVersion`, the variable that actually feeds the User-Agent, and the baked-in constant tracks the v2 line so unstamped local builds no longer claim v1. The stale `-X main.commit` flag has been dropped. ([#10109](https://github.com/aliyun/terraform-provider-alicloud/issues/10109))

## 2.0.0-beta1 (July 27, 2026)

FEATURES:

- **Terraform Plugin SDK v2:** the whole provider (all resources, data sources and the acceptance-test harness) has been migrated from `terraform-plugin-sdk` v1 to `terraform-plugin-sdk/v2`.

BREAKING CHANGES:

- provider: migrate from `terraform-plugin-sdk` v1 to `terraform-plugin-sdk/v2`. Terraform CLI 0.11 and earlier are no longer supported. Before upgrading, move to v1.282.0 or later and confirm `terraform plan` is clean. See the [Version 2 Upgrade Guide](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/guides/version-2-upgrade).
- provider: `d.SetPartial` is not supported by `terraform-plugin-sdk/v2`, so all 1437 call sites across 210 resources have been removed. Partial-state tracking on a failed apply is therefore no longer in effect.
- resource/alicloud_cr_repo, data-source/alicloud_cr_repos: the backing Container Registry API (version `2016-06-07`) is no longer available, so both are non-functional and must be removed from your configuration. `domain_list` also changed from `TypeMap` to `TypeList`.
- resource/alicloud_api_gateway_instance: `to_connect_vpc_ip_block` changed from `TypeMap` to `TypeList`. Existing state is upgraded automatically by a v0 -> v1 state upgrader; configurations must be updated from map syntax to block syntax.
- resource/alicloud_cs_edge_kubernetes, resource/alicloud_cs_kubernetes: `runtime`, `certificate_authority` and `connections` changed from `TypeMap` to `TypeList`. Existing state is upgraded automatically; update attribute references from map-key syntax to list-index syntax.
- resource/alicloud_cs_managed_kubernetes: `certificate_authority` and `connections` changed from `TypeMap` to `TypeList`. Existing state is upgraded automatically.
- resource/alicloud_cs_edge_kubernetes: cluster creation now goes through the ROA Container Service `CreateCluster` API instead of the legacy `CreateManagedKubernetesCluster`, and exactly one default Edge node pool must be configured.
- data-source/alicloud_instance_types: `instance_types.gpu`, `instance_types.burstable_instance` and `instance_types.local_storage` changed from `TypeMap` to `TypeList`.
- data-source/alicloud_db_instance_classes: `instance_classes.storage_range` changed from `TypeMap` to `TypeList`.
- data-source/alicloud_cs_cluster_credential: `certificate_authority` changed from `TypeMap` to `TypeList`.
- data-source/alicloud_cs_edge_kubernetes_clusters, data-source/alicloud_cs_kubernetes_clusters, data-source/alicloud_cs_managed_kubernetes_clusters, data-source/alicloud_cs_serverless_kubernetes_clusters: `clusters.connections` changed from `TypeMap` to `TypeList`.
- provider: `terraform-plugin-sdk/v2` dropped support for the `Removed:` schema property, so the 157 placeholder schema entries across 39 resources and data sources that existed only to emit "this field has been removed in vX" errors have been deleted. These arguments already had no effect in 1.x; configurations that still reference them now fail with Terraform's generic "Unsupported argument" error instead of the provider's explanatory message. Affected resources include `alicloud_instance`, `alicloud_nat_gateway`, `alicloud_cms_alarm`, `alicloud_slb_listener`, `alicloud_lindorm_instance`, `alicloud_sae_application`, `alicloud_ros_stack_group`, `alicloud_eci_container_group` and the `alicloud_cs_*` family.

DEPRECATIONS:

- provider: sync deprecation notices from the documentation into the schema. 72 resources and data sources that were already documented as deprecated (including `alicloud_mns_*`, `alicloud_network_acl_entries`, `alicloud_network_acl_attachment`, `alicloud_slb_attachment`, `alicloud_router_interface*`, `alicloud_scdn_*`, `alicloud_waf_domain`, `alicloud_waf_instance`, `alicloud_tsdb_instance`, `alicloud_kvstore_backup_policy`, `alicloud_rdc_organization`, `alicloud_config_delivery_channel`, `alicloud_emr_cluster`) now emit a deprecation warning at plan time.
- resource/alicloud_dbfs_*, data-source/alicloud_dbfs_*: mark all DBFS resources and data sources as deprecated since v1.279.0.
- data-source/alicloud_bp_studio_applications: deprecate `applications.topo_url`; it is no longer returned by the server.

ENHANCEMENTS:

- docs: add the Version 2 upgrade guide (`guides/version-2-upgrade`) documenting every type change and state migration in this release.
- provider: add v0 -> v1 state upgraders for `alicloud_api_gateway_instance`, `alicloud_cr_repo`, `alicloud_cs_edge_kubernetes`, `alicloud_cs_kubernetes` and `alicloud_cs_managed_kubernetes` so existing state is migrated in place.
- provider: migrate the acceptance-test plumbing from `Providers` to `ProviderFactories`.
- resource/alicloud_cs_edge_kubernetes: add the computed `slb_id` attribute and a dedicated delete path.
- resource/alicloud_cen_transit_router: mark `support_multicast` as Computed.
- resource/alicloud_polardb_database: mark `account_name` as Computed.
- resource/alicloud_emrv2_cluster: apply the auto scaling policy via `PutAutoScalingPolicy` after `CreateNodeGroup`.
- resource/alicloud_ebs_solution_instance: switch tag creation to the `TagResources` API.
- resource/alicloud_gpdb_instance: wait for `seg_node_num` to converge after an update.
- resource/alicloud_dms_enterprise_workspace: send `WorkspaceRegion` on create.
- resource/alicloud_ram_login_profile: only send `Password` when it actually changed.
- resource/alicloud_selectdb_db_instance: simplify how the upgrade target version is resolved.
- resource/alicloud_ecd_command: retry on error code `CloudAssistant.NotReady`.
- resource/alicloud_ecd_desktop: retry on error code `INTERNAL_ERROR`.
- resource/alicloud_ecd_network_package: add `Updating` to the pending states when waiting.
- resource/alicloud_polardb_cluster: retry `DescribeParameters` on throttling, and retry on `OperationDenied.DBClusterStatus`.
- resource/alicloud_dms_airflow: retry on `DMS.AIRFLOW.InternalError`.
- resource/alicloud_resource_manager_policy_attachment: add create and delete timeouts, and retry on `EntityNotExist.Role`.
- resource/alicloud_ots_tunnel: treat `OTSPermissionDenied` as temporarily unavailable and retry.
- provider: `NotFoundError` now recognizes OTS tunnel errors — `ParamInvalid` with "tunnel not exist", and `PermissionDenied` with "is not found" or "instance is not running" — as not-found.
- resource/alicloud_dms_enterprise_user: treat the `DELETE` state as not-found.
- resource/alicloud_ens_instance: raise the create, update and delete timeouts to 60 minutes.
- resource/alicloud_ens_eip_instance_attachment: raise the create and delete timeouts to 10 minutes.
- resource/alicloud_edas_k8s_cluster: raise the delete timeout to 60 minutes and treat `ClusterImportStatus` 4 as already deleted.
- resource/alicloud_data_works_project: wait for `PaiTaskEnabled` to converge, align throttling retries and degrade gracefully when the waiter cannot resolve.
- resource/alicloud_simple_application_server_custom_image: wait for the asynchronous create and delete to complete.
- resource/alicloud_vpc_ipam_service: treat `ExceedPurchaseLimit` as already-opened so opening the service is idempotent.
- resource/alicloud_compute_nest_service_instance: suppress the spurious `operation_metadata.resources` ForceNew diff.
- resource/alicloud_oss_access_point: sync status before the next operation.
- build: add `terraform-registry-manifest.json` declaring Terraform plugin protocol version 5.0.
- deps: bump Go module dependencies and vendor `hc-install` v0.9.5.

BUG FIXES:

- provider: guard the SDK v2 unknown-value sentinel in the credential configuration. ([#10038](https://github.com/aliyun/terraform-provider-alicloud/issues/10038))
- provider: remove `d.Set` calls targeting attributes absent from the schema, which panic under SDK v2. Affects resource/alicloud_cassandra_cluster, resource/alicloud_cassandra_data_center, resource/alicloud_mse_nacos_config, resource/alicloud_polardb_cluster, resource/alicloud_slb_domain_extension, resource/alicloud_cloud_firewall_vpc_firewall, resource/alicloud_cloud_firewall_vpc_firewall_cen (`id`); resource/alicloud_cs_edge_kubernetes and resource/alicloud_cs_kubernetes (various); resource/alicloud_vpc_route_entry (`name`); resource/alicloud_vpc_public_ip_address_pool (`zones`); resource/alicloud_security_center_service_linked_role (`product_name`); resource/alicloud_db_readonly_instance (`ssl_status`); resource/alicloud_yundun_dbaudit_instance (`region_id`); resource/alicloud_bastionhost_instance (`ad_auth_server.has_password`); data-source/alicloud_cs_kubernetes_node_pools (`names`). These attributes were never declared, so they were never present in state.
- resource/alicloud_nas_smb_acl_attachment: fix a panic caused by a boolean-to-string conversion.
- resource/alicloud_governance_account: format `account_id` as an int in Read to fix a panic.
- resource/alicloud_rds_db_proxy: fix an `ssl_status` panic.
- resource/alicloud_image: cast `disk_device_mapping.size` to an int in Read.
- resource/alicloud_image_import: read back through its own Read function instead of `alicloud_image`'s.
- resource/alicloud_ros_stack: marshal `stack_policy_body` to a JSON string in Read.
- resource/alicloud_express_connect_router_interface: cast `has_reservation_data` to a string in Read.
- resource/alicloud_instance: cast `deployment_set_group_no` to a string in Read.
- resource/alicloud_adb_backup_policy: cast `backup_retention_period` to a string in Read.
- resource/alicloud_gpdb_external_data_service: parse `service_id` as an integer in Read.
- resource/alicloud_gpdb_hadoop_data_source, resource/alicloud_gpdb_remote_adb_data_source: stop setting the integer id attributes from string ID parts.
- resource/alicloud_amqp_instance: only send `security_group_id` and `vpc_id` when they are set, and parse `create_time` from RFC3339 into the integer attribute.
- resource/alicloud_ecs_instance_set: read `tags` from the nested `Tags.Tag` structure.
- resource/alicloud_mse_engine_namespace: declare five missing Computed schema fields.
- resource/alicloud_config_rule: convert `resource_types_scope` to a comma-separated string for the deprecated `scope_compliance_resource_types`, and drop a redundant type assertion on a `TypeList` `d.Get`.
- resource/alicloud_click_house_db_cluster, resource/alicloud_click_house_account: fix the `resource_group_id` update wait and the account-type perpetual diff.
- resource/alicloud_polardb_account_privilege: sort `db_names` to avoid plan drift.
- resource/alicloud_sls_logtail_config: fix the persistent diff on `input_detail` and `output_detail`.
- resource/alicloud_sls_logtail_pipeline_config: fix the non-empty plan caused by API auto-populated fields.
- resource/alicloud_kms_instance: remove the invalid `instance_id` set.
- resource/alicloud_max_compute_quota, data-source/alicloud_maxcompute_service: repair the service data source and pay-as-you-go quota re-creation.
- resource/alicloud_rds_upgrade_db_instance, resource/alicloud_rds_clone_db_instance: fix the `ssl_enabled` typo.
- resource/alicloud_dts_migration_job: stop translating `payment_type` to and from `PayType`.
- resource/alicloud_ram_role_attachment: propagate create errors instead of swallowing them, and refresh state after create.
- resource/alicloud_yundun_dbaudit_instance: use `ALIYUN::DBAUDIT::INSTANCE` as the tag resource type, and drop the obsolete hardcoded `SeriesCode`/`NetworkType`/`RegionId` spec parameters.
- resource/alicloud_ens_instance_security_group_attachment: read `security_group_id` from the resource ID instead of the describe response.
- resource/alicloud_ens_security_group: treat `InvalidSecurityGroupId.NotFound` as not-found in Read.
- resource/alicloud_esa_routine: fix `NotFound` handling.
- resource/alicloud_fcv3_function: fix clearing `gpu_config` on update.
- resource/alicloud_pai_workspace_datasetversion: treat error code `201300002` as `NotFound`.
- resource/alicloud_pai_service, resource/alicloud_pai_workspace_dataset: fix SDK v2 read regressions.
- resource/alicloud_auto_provisioning_group: fix the resource and the `EcsSccSupportedRegions` region list.
- resource/alicloud_simple_application_server_firewall_rule: handle the asynchronous delete.
- data-source/alicloud_bastionhost_instances: report `storage` and `bandwidth` per instance instead of at the top level.
- data-source/alicloud_adb_db_cluster_lake_versions: cast `expired` to a string.
- data-source/alicloud_emrv2_clusters: return an empty result when `ListClusters` reports `ResourceNotFound`.
- data-source/alicloud_simple_application_server_snapshots: make the `status` filter case-insensitive.
- data-source/alicloud_simple_application_server_disks: make the `disk_type` filter case-insensitive.
- provider: resolve `go vet` warnings and lint issues across the package.

This file tracks the v2.x of the provider. Changes to the v1.x stay in [CHANGELOG.md](CHANGELOG.md).

For migration instructions, see the [Version 2 Upgrade Guide](https://registry.terraform.io/providers/aliyun/alicloud/latest/docs/guides/version-2-upgrade).
