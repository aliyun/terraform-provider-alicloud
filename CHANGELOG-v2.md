## 2.0.0-beta3 (August 17, 2026)

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
