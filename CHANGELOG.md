## 1.95.0 (Unreleased)
## 1.94.0 (August 24, 2020)

- **New Resource:** `alicloud_dcdn_domain` ([#2744](https://github.com/aliyun/terraform-provider-alicloud/issues/2744))
- **New Resource:** `alicloud_mse_cluster` ([#2733](https://github.com/aliyun/terraform-provider-alicloud/issues/2733))
- **New Resource:** `alicloud_resource_manager_policy_attachment` ([#2696](https://github.com/aliyun/terraform-provider-alicloud/issues/2696))
- **Data Source:** `alicloud_dcdn_domains` ([#2744](https://github.com/aliyun/terraform-provider-alicloud/issues/2744))
- **Data Source:** `alicloud_mse_clusters` ([#2733](https://github.com/aliyun/terraform-provider-alicloud/issues/2733))
- **Data Source:** `alicloud_resource_manager_policy_attachments` ([#2696](https://github.com/aliyun/terraform-provider-alicloud/issues/2696))

IMPROVEMENTS:

- Support allocate and release public connection for redis ([#2748](https://github.com/aliyun/terraform-provider-alicloud/issues/2748))
- Support to set warn and info level alarm ([#2743](https://github.com/aliyun/terraform-provider-alicloud/issues/2743))
- waf domain support setting resource_group_id and more attributes ([#2740](https://github.com/aliyun/terraform-provider-alicloud/issues/2740))
- resource dnat supports "import" feature ([#2735](https://github.com/aliyun/terraform-provider-alicloud/issues/2735))
- Add func sweep and Change testcase frequency ([#2726](https://github.com/aliyun/terraform-provider-alicloud/issues/2726))
- Correct provider docs order ([#2723](https://github.com/aliyun/terraform-provider-alicloud/issues/2723))
- Remove github.com/hashicorp/terraform import and use terraform-plugin-sdk instead ([#2722](https://github.com/aliyun/terraform-provider-alicloud/issues/2722))
- Add test sweep for eci_image_cache ([#2720](https://github.com/aliyun/terraform-provider-alicloud/issues/2720))
- modify alicloud_cen_instance_attachment ([#2714](https://github.com/aliyun/terraform-provider-alicloud/issues/2714))

BUG FIXES:

- fix the bug of create emr kafka cluster error ([#2754](https://github.com/aliyun/terraform-provider-alicloud/issues/2754))
- fix common bandwidth package idempotent issue when Adding and Removeing instance ([#2750](https://github.com/aliyun/terraform-provider-alicloud/issues/2750))
- fix website document error using `terraform` tag ([#2749](https://github.com/aliyun/terraform-provider-alicloud/issues/2749))
- Fix registry rendering of page ([#2747](https://github.com/aliyun/terraform-provider-alicloud/issues/2747))
- fix ci test website-test error ([#2742](https://github.com/aliyun/terraform-provider-alicloud/issues/2742))
- fix datasource for ResourceManager for Policy Attachment ([#2730](https://github.com/aliyun/terraform-provider-alicloud/issues/2730))
- fix_ecs_snapshot ([#2709](https://github.com/aliyun/terraform-provider-alicloud/issues/2709))

## 1.93.0 (August 12, 2020)

- **New Resource:** `alicloud_oos_execution` ([#2679](https://github.com/aliyun/terraform-provider-alicloud/issues/2679))
- **New Resource:** `alicloud_edas_k8s_cluster` ([#2678](https://github.com/aliyun/terraform-provider-alicloud/issues/2678))
- **Data Source:** `alicloud_oos_execution` ([#2679](https://github.com/aliyun/terraform-provider-alicloud/issues/2679))

IMPROVEMENTS:

- Add sweep func for adb cluster test ([#2716](https://github.com/aliyun/terraform-provider-alicloud/issues/2716))
- Add default vpc for drds ([#2713](https://github.com/aliyun/terraform-provider-alicloud/issues/2713))
- ADB MySQL output output connection string after creation ([#2699](https://github.com/aliyun/terraform-provider-alicloud/issues/2699))
- add .goreleaser.yml ([#2698](https://github.com/aliyun/terraform-provider-alicloud/issues/2698))
- transfer terraform-provider-alicloud to aliyun from terraform-providers ([#2697](https://github.com/aliyun/terraform-provider-alicloud/issues/2697))
- Add purge cluster api for cassandra sweeper ([#2693](https://github.com/aliyun/terraform-provider-alicloud/issues/2693))
- Add default vpc for mongodb ([#2689](https://github.com/aliyun/terraform-provider-alicloud/issues/2689))
- Add default vpc for kvstore ([#2688](https://github.com/aliyun/terraform-provider-alicloud/issues/2688))
- Add sweeper for cassandra cluster ([#2687](https://github.com/aliyun/terraform-provider-alicloud/issues/2687))
- Support 'resoruce_group_id' attribute for ImportKeyPair ([#2683](https://github.com/aliyun/terraform-provider-alicloud/issues/2683))
- Support to get NotFound error in read method ([#2682](https://github.com/aliyun/terraform-provider-alicloud/issues/2682))
- UPDATE CHANGELOG ([#2681](https://github.com/aliyun/terraform-provider-alicloud/issues/2681))
- Support specify security group when create instance ([#2680](https://github.com/aliyun/terraform-provider-alicloud/issues/2680))
- improve(slb) update slb_backend_server add parameter server_ip ([#2651](https://github.com/aliyun/terraform-provider-alicloud/issues/2651))

BUG FIXES:

- update: fix dnat query errror by only use forwardTableId ([#2712](https://github.com/aliyun/terraform-provider-alicloud/issues/2712))
- fix: create rds sql_collector_status bug ([#2690](https://github.com/aliyun/terraform-provider-alicloud/issues/2690))
- fix(edas): improve sweeper test ([#2686](https://github.com/aliyun/terraform-provider-alicloud/issues/2686))
- fix cassandra doc and add describe not found error ([#2685](https://github.com/aliyun/terraform-provider-alicloud/issues/2685))
- fix doc: attach AliyunMNSNotificationRolePolicy to role ([#2572](https://github.com/aliyun/terraform-provider-alicloud/issues/2572))
- docs: fix typos and grammar in Alicloud Provider ([#2559](https://github.com/aliyun/terraform-provider-alicloud/issues/2559))
- fix_markdown_auto_provisioning_group ([#2543](https://github.com/aliyun/terraform-provider-alicloud/issues/2543))
- fix_markdown_snapshot_policy ([#2540](https://github.com/aliyun/terraform-provider-alicloud/issues/2540))

## 1.92.0 (July 31, 2020)

- **New Resource:** `alicloud_oos_template` ([#2670](https://github.com/aliyun/terraform-provider-alicloud/issues/2670))
- **Data Source:** `alicloud_oos_template` ([#2670](https://github.com/aliyun/terraform-provider-alicloud/issues/2670))

IMPROVEMENTS:

- modify alicloud_cen_bandwidth_package_attachment ([#2675](https://github.com/aliyun/terraform-provider-alicloud/issues/2675))
- UPDATE CHANGELOG ([#2671](https://github.com/aliyun/terraform-provider-alicloud/issues/2671))
- upgrade resource of Nas AccessGroup ([#2667](https://github.com/aliyun/terraform-provider-alicloud/issues/2667))
- Supports setting the kms id for oss bucket ([#2662](https://github.com/aliyun/terraform-provider-alicloud/issues/2662))
- Support service_account_issuer and api_audiences in alicloud_cs_kubernetes and alicloud_cs_managed_kubernetes ([#2573](https://github.com/aliyun/terraform-provider-alicloud/issues/2573))

BUG FIXES:

- Fix ess kms disk ([#2668](https://github.com/aliyun/terraform-provider-alicloud/issues/2668))

## 1.91.0 (July 24, 2020)

- **New Resource:** `alicloud_ecs_dedicated_host` ([#2652](https://github.com/aliyun/terraform-provider-alicloud/issues/2652))
- **Data Source:** `alicloud_ecs_dedicated_hosts` ([#2652](https://github.com/aliyun/terraform-provider-alicloud/issues/2652))

IMPROVEMENTS:

- improve test case name ([#2672](https://github.com/aliyun/terraform-provider-alicloud/issues/2672))
- update: add nat binded eip ipaddress ([#2669](https://github.com/aliyun/terraform-provider-alicloud/issues/2669))
- correct log alert testcase name ([#2663](https://github.com/aliyun/terraform-provider-alicloud/issues/2663))
- add example moudle for OSS bucket ([#2661](https://github.com/aliyun/terraform-provider-alicloud/issues/2661))
- cs cluster support data disks ([#2657](https://github.com/aliyun/terraform-provider-alicloud/issues/2657))
- drds support internation ([#2654](https://github.com/aliyun/terraform-provider-alicloud/issues/2654))
- UPDATE CHANGELOG ([#2649](https://github.com/aliyun/terraform-provider-alicloud/issues/2649))
- modify cen_instance ([#2644](https://github.com/aliyun/terraform-provider-alicloud/issues/2644))
- add ability to enable ZRS on bucket creation ([#2605](https://github.com/aliyun/terraform-provider-alicloud/issues/2605))

BUG FIXES:

- fix ddh testcase ([#2665](https://github.com/aliyun/terraform-provider-alicloud/issues/2665))
- fix dms instance ([#2656](https://github.com/aliyun/terraform-provider-alicloud/issues/2656))
- fix_markdown_ess_scheduled_task ([#2655](https://github.com/aliyun/terraform-provider-alicloud/issues/2655))
- fix slb_listener creates NewCommonRequest error handling ([#2653](https://github.com/aliyun/terraform-provider-alicloud/issues/2653))

## 1.90.1 (July 15, 2020)

IMPROVEMENTS:

- perf: rds ssl and tde limitation ([#2645](https://github.com/aliyun/terraform-provider-alicloud/issues/2645))
- add isp support to cbwp ([#2642](https://github.com/aliyun/terraform-provider-alicloud/issues/2642))
- Remove the resource_group_id parameter when querying the system disk ([#2641](https://github.com/aliyun/terraform-provider-alicloud/issues/2641))
- Add 'testAcc' prefix for test case name ([#2636](https://github.com/aliyun/terraform-provider-alicloud/issues/2636))
- Support DescribeInstanceSystemDisk method return the error message ([#2635](https://github.com/aliyun/terraform-provider-alicloud/issues/2635))
- UPDATE CHANGELOG ([#2634](https://github.com/aliyun/terraform-provider-alicloud/issues/2634))

BUG FIXES:

- fix cassandra doc  ([#2648](https://github.com/aliyun/terraform-provider-alicloud/issues/2648))
- fix_markdown_ess_scheduled_task ([#2647](https://github.com/aliyun/terraform-provider-alicloud/issues/2647))
- fix WAF instance testcase ([#2640](https://github.com/aliyun/terraform-provider-alicloud/issues/2640))
- fix testcase name ([#2638](https://github.com/aliyun/terraform-provider-alicloud/issues/2638))
- fix_instance ([#2632](https://github.com/aliyun/terraform-provider-alicloud/issues/2632))

## 1.90.0 (July 10, 2020)

- **New Resource:** `alicloud_container_registry_enterprise_sync_rule` ([#2607](https://github.com/aliyun/terraform-provider-alicloud/issues/2607))
- **New Resource:** `alicloud_dms_user` ([#2604](https://github.com/aliyun/terraform-provider-alicloud/issues/2604))
- **Data Source:** `alicloud_cr_ee_sync_rules` ([#2630](https://github.com/aliyun/terraform-provider-alicloud/issues/2630))
- **Data Source:** `alicloud_eci_image_cache` ([#2627](https://github.com/aliyun/terraform-provider-alicloud/issues/2627))
- **Data Source:** `alicloud_waf_instance` ([#2617](https://github.com/aliyun/terraform-provider-alicloud/issues/2617))
- **Data Source:** `alicloud_dms_user` ([#2604](https://github.com/aliyun/terraform-provider-alicloud/issues/2604))

IMPROVEMENTS:

- support the CNAME of CDN domain new ([#2622](https://github.com/aliyun/terraform-provider-alicloud/issues/2622))
- UPDATE CHANGELOG ([#2594](https://github.com/aliyun/terraform-provider-alicloud/issues/2594))
- Feature/disable addon ([#2590](https://github.com/aliyun/terraform-provider-alicloud/issues/2590))
- set system default and make fmt ([#2480](https://github.com/aliyun/terraform-provider-alicloud/issues/2480))

BUG FIXES:

- fix_ess_scheduled_task ([#2628](https://github.com/aliyun/terraform-provider-alicloud/issues/2628))
- fix testcase of WAF instance datasource ([#2625](https://github.com/aliyun/terraform-provider-alicloud/issues/2625))
- fix oss lifecycle rule match the whole bucket by default ([#2621](https://github.com/aliyun/terraform-provider-alicloud/issues/2621))
- fix ack uat ([#2618](https://github.com/aliyun/terraform-provider-alicloud/issues/2618))

## 1.89.0 (July 03, 2020)

- **New Resource:** `alicloud_eci_image_cache` ([#2615](https://github.com/aliyun/terraform-provider-alicloud/issues/2615))

IMPROVEMENTS:

- improve(alikafka): using default vswitch to run alikafka testcases ([#2591](https://github.com/aliyun/terraform-provider-alicloud/issues/2591))
- run 'go mod vendor' to sync ([#2587](https://github.com/aliyun/terraform-provider-alicloud/issues/2587))
- update waf SDK ([#2616](https://github.com/aliyun/terraform-provider-alicloud/issues/2616))
- umodify cen_route_map ([#2606](https://github.com/aliyun/terraform-provider-alicloud/issues/2606))
- modify cen_bandwidth_package ([#2603](https://github.com/aliyun/terraform-provider-alicloud/issues/2603))
- support region cn-wulanchabu ([#2599](https://github.com/aliyun/terraform-provider-alicloud/issues/2599))
- Add version_stage filter ([#2597](https://github.com/aliyun/terraform-provider-alicloud/issues/2597))
- support releasing ddoscoo instance ([#2595](https://github.com/aliyun/terraform-provider-alicloud/issues/2595))
- Support modify system_disk_size online ([#2593](https://github.com/aliyun/terraform-provider-alicloud/issues/2593))
- Changelog ([#2584](https://github.com/aliyun/terraform-provider-alicloud/issues/2584))

BUG FIXES:

- fix cms site monitor document ([#2614](https://github.com/aliyun/terraform-provider-alicloud/issues/2614))
- fix_kms_ecs_disk ([#2613](https://github.com/aliyun/terraform-provider-alicloud/issues/2613))
- fix_markdown_ess_notification ([#2598](https://github.com/aliyun/terraform-provider-alicloud/issues/2598))
- fix kms secret, secret version doc ([#2596](https://github.com/aliyun/terraform-provider-alicloud/issues/2596))
- fix testcase for pvtz_zone ([#2592](https://github.com/aliyun/terraform-provider-alicloud/issues/2592))
- fix the dns_record test case bug ([#2588](https://github.com/aliyun/terraform-provider-alicloud/issues/2588))
- fix_markdown_launch_template ([#2541](https://github.com/aliyun/terraform-provider-alicloud/issues/2541))

## 1.88.0 (June 22, 2020)

- **New Resource:** `alicloud_cen_vbr_health_check` ([#2575](https://github.com/aliyun/terraform-provider-alicloud/issues/2575))
- **Data Source:** `alicloud_cen_private_zones` ([#2564](https://github.com/aliyun/terraform-provider-alicloud/issues/2564))
- **Data Source:** `alicloud_dms_enterprise_instances` ([#2557](https://github.com/aliyun/terraform-provider-alicloud/issues/2557))
- **Data Source:** `alicloud_cassandra` ([#2574](https://github.com/aliyun/terraform-provider-alicloud/issues/2574))
- **Data Source:** `alicloud_kms_secret_versions` ([#2583](https://github.com/aliyun/terraform-provider-alicloud/issues/2583))

IMPROVEMENTS:

- skip instance prepaid testcase ([#2585](https://github.com/aliyun/terraform-provider-alicloud/issues/2585))
- Support setting NO_PROXY and upgrade go sdk ([#2581](https://github.com/aliyun/terraform-provider-alicloud/issues/2581))
- Features/atoscaler_use_worker_token ([#2578](https://github.com/aliyun/terraform-provider-alicloud/issues/2578))
- Features/knock autoscaler off nodes ([#2571](https://github.com/aliyun/terraform-provider-alicloud/issues/2571))
- modify cen_instance_attachment ([#2566](https://github.com/aliyun/terraform-provider-alicloud/issues/2566))
- gpdb doc change "tf-gpdb-test"" to "tf-gpdb-test" ([#2561](https://github.com/aliyun/terraform-provider-alicloud/issues/2561))
- UPDATE CHANGELOG ([#2555](https://github.com/aliyun/terraform-provider-alicloud/issues/2555))
- cassandra cluster ([#2522](https://github.com/aliyun/terraform-provider-alicloud/issues/2522))

BUG FIXES:

- Fix the fc-function testcase and markdown ([#2569](https://github.com/aliyun/terraform-provider-alicloud/issues/2569))
- fix name spelling mistake ([#2558](https://github.com/aliyun/terraform-provider-alicloud/issues/2558))

## 1.87.0 (June 12, 2020)

- **Data Source:** `alicloud_container_registry_enterprise_repos` ([#2538](https://github.com/aliyun/terraform-provider-alicloud/issues/2538))
- **Data Source:** `alicloud_container_registry_enterprise_namespaces` ([#2530](https://github.com/aliyun/terraform-provider-alicloud/issues/2530))
- **Data Source:** `alicloud_container_registry_enterprise_instances` ([#2526](https://github.com/aliyun/terraform-provider-alicloud/issues/2526))
- **Data Source:** `alicloud_cen_route_maps` ([#2554](https://github.com/aliyun/terraform-provider-alicloud/issues/2554))

IMPROVEMENTS:

- adapter schedulerrule ([#2537](https://github.com/aliyun/terraform-provider-alicloud/issues/2537))
- UPDATE CHANHELOG ([#2535](https://github.com/aliyun/terraform-provider-alicloud/issues/2535))
- improve_user_experience ([#2491](https://github.com/aliyun/terraform-provider-alicloud/issues/2491))
- add testcase ([#2556](https://github.com/aliyun/terraform-provider-alicloud/issues/2556))
- improve(elasticsearch): resource support to open or close network, and modify the kibana whitelist in private network ([#2548](https://github.com/aliyun/terraform-provider-alicloud/issues/2548))
- support "resource_group_id" for Bastionhost Instance ([#2544](https://github.com/aliyun/terraform-provider-alicloud/issues/2544))
- support "resource_group_id" for DBaudit Instance ([#2539](https://github.com/aliyun/terraform-provider-alicloud/issues/2539))
- Automatically generate dns_domain datasource ([#2549](https://github.com/aliyun/terraform-provider-alicloud/issues/2549))

BUG FIXES:

- Fix image export ([#2542](https://github.com/aliyun/terraform-provider-alicloud/issues/2542))
- fix: perf create rds pg ([#2533](https://github.com/aliyun/terraform-provider-alicloud/issues/2533))
- fix_markdown_ess_scalinggroup ([#2529](https://github.com/aliyun/terraform-provider-alicloud/issues/2529))
- fix_markdown_image_import ([#2520](https://github.com/aliyun/terraform-provider-alicloud/issues/2520))
- fix_markdown_disk ([#2504](https://github.com/aliyun/terraform-provider-alicloud/issues/2504))
- fix_markdown_image_s ([#2546](https://github.com/aliyun/terraform-provider-alicloud/issues/2546))

## 1.86.0 (June 05, 2020)

- **New Resource:** `alicloud_container_registry_enterprise_repo` ([#2525](https://github.com/aliyun/terraform-provider-alicloud/issues/2525))
- **New Resource:** `alicloud_Container_registry_enterprise_namespace` ([#2519](https://github.com/aliyun/terraform-provider-alicloud/issues/2519))
- **New Resource:** `alicloud_ddoscoo_scheduler_rule` ([#2476](https://github.com/aliyun/terraform-provider-alicloud/issues/2476))
- **New Resource:** `alicloud_resource_manager_policies` ([#2474](https://github.com/aliyun/terraform-provider-alicloud/issues/2474))
- **Data Source:** `alicloud_waf_domains` ([#2498](https://github.com/aliyun/terraform-provider-alicloud/issues/2498))
- **Data Source:** `alicloud_kms_secrets` ([#2515](https://github.com/aliyun/terraform-provider-alicloud/issues/2515))
- **Data Source:** `alicloud_alidns_domain_records` ([#2503](https://github.com/aliyun/terraform-provider-alicloud/issues/2503))
- **Data Source:** `alicloud_resource_manager_resource_directories` ([#2499](https://github.com/aliyun/terraform-provider-alicloud/issues/2499))
- **Data Source:** `alicloud_resource_manager_handshakes` ([#2489](https://github.com/aliyun/terraform-provider-alicloud/issues/2489))
- **Data Source:** `alicloud_resource_manager_accounts` ([#2488](https://github.com/aliyun/terraform-provider-alicloud/issues/2488))
- **Data Source:** `alicloud_resource_manager_roles` ([#2483](https://github.com/aliyun/terraform-provider-alicloud/issues/2483))

IMPROVEMENTS:
- support "resource_group_id" for Elasticsearch instance ([#2528](https://github.com/aliyun/terraform-provider-alicloud/issues/2528))
- Added new feature of encrypting data node disk ([#2521](https://github.com/aliyun/terraform-provider-alicloud/issues/2521))
- support "resource_group_id" for Private Zone ([#2518](https://github.com/aliyun/terraform-provider-alicloud/issues/2518))
- support "resource_group_id" for DB instance ([#2514](https://github.com/aliyun/terraform-provider-alicloud/issues/2514))
- 更新sdk到v1.61.230 ([#2510](https://github.com/aliyun/terraform-provider-alicloud/issues/2510))
- support "resource_group_id" for kvstore instance ([#2509](https://github.com/aliyun/terraform-provider-alicloud/issues/2509))
- Add Log Dashboard ([#2502](https://github.com/aliyun/terraform-provider-alicloud/issues/2502))
- UPDATE CHANGELOG ([#2497](https://github.com/aliyun/terraform-provider-alicloud/issues/2497))
- Control the instance start and stop through the status attribute ([#2464](https://github.com/aliyun/terraform-provider-alicloud/issues/2464))

BUG FIXES:
- fix_markdown_image_import ([#2516](https://github.com/aliyun/terraform-provider-alicloud/issues/2516))
- fix ecs 'status' attribute bug ([#2512](https://github.com/aliyun/terraform-provider-alicloud/issues/2512))
- fix_markdown_ess_scalinggroup_vserver_groups ([#2508](https://github.com/aliyun/terraform-provider-alicloud/issues/2508))
- fix_markdown_disk_attachment ([#2501](https://github.com/aliyun/terraform-provider-alicloud/issues/2501))
- fix_markdown_disk_attachment ([#2481](https://github.com/aliyun/terraform-provider-alicloud/issues/2481))

## 1.85.0 (May 29, 2020)

- **New Resource:** `alicloud_alidns_record` ([#2495](https://github.com/aliyun/terraform-provider-alicloud/issues/2495))
- **New Resource:** `alicloud_kms_key` ([#2444](https://github.com/aliyun/terraform-provider-alicloud/issues/2444))
- **New Resource:** `alicloud_kms_keyversion` ([#2471](https://github.com/aliyun/terraform-provider-alicloud/issues/2471))
- **Data Source:** `alicloud_resource_manager_policy_versions` ([#2496](https://github.com/aliyun/terraform-provider-alicloud/issues/2496))
- **Data Source:** `alicloud_kms_key_versions` ([#2494](https://github.com/aliyun/terraform-provider-alicloud/issues/2494))
- **Data Source:** `alicloud_alidns_domain_group` ([#2482](https://github.com/aliyun/terraform-provider-alicloud/issues/2482))

IMPROVEMENTS:

- 增加cdn_config删除错误码 ([#2490](https://github.com/aliyun/terraform-provider-alicloud/issues/2490))
- UPDATE CHANGELOG.md ([#2477](https://github.com/aliyun/terraform-provider-alicloud/issues/2477))
- Alicloud edas docs modify ([#2473](https://github.com/aliyun/terraform-provider-alicloud/issues/2473))

BUG FIXES:

- fix_markdown_reserved_instance ([#2478](https://github.com/aliyun/terraform-provider-alicloud/issues/2478))
- fix_markdown_network_interfaces ([#2475](https://github.com/aliyun/terraform-provider-alicloud/issues/2475))
- fix_apg ([#2472](https://github.com/aliyun/terraform-provider-alicloud/issues/2472))

## 1.84.0 (May 22, 2020)

- **New Resource:** `alicloud_alidns_domain_group.` ([#2454](https://github.com/aliyun/terraform-provider-alicloud/issues/2454))
- **New Resource:** `alicloud_resource_manager_resource_directory` ([#2459](https://github.com/aliyun/terraform-provider-alicloud/issues/2459))
- **New Resource:** `alicloud_resource_manager_policy_version` ([#2457](https://github.com/aliyun/terraform-provider-alicloud/issues/2457))
- **Data Source:** `alicloud_resource_manager_folders` ([#2467](https://github.com/aliyun/terraform-provider-alicloud/issues/2467))
- **Data Source:** `alicloud_alidns_instance.` ([#2468](https://github.com/aliyun/terraform-provider-alicloud/issues/2468))
- **Data Source:** `alicloud_resource_manager_resource_groups` ([#2462](https://github.com/aliyun/terraform-provider-alicloud/issues/2462))

IMPROVEMENTS:

- Update CHANGELOG.md ([#2455](https://github.com/aliyun/terraform-provider-alicloud/issues/2455))

BUG FIXES:

- fix autoscaler configmap update ([#2377](https://github.com/aliyun/terraform-provider-alicloud/issues/2377))
- fix eip association failed cause by snat entry's snat_ip update bug ([#2440](https://github.com/aliyun/terraform-provider-alicloud/issues/2440))
- fix_tag_validation ([#2445](https://github.com/aliyun/terraform-provider-alicloud/issues/2445))
- fix polardb connection string output bug ([#2453](https://github.com/aliyun/terraform-provider-alicloud/issues/2453))
- fix the bug of TestAccAlicloudEmrCluster_local_storage nodeCount less… ([#2458](https://github.com/aliyun/terraform-provider-alicloud/issues/2458))
- fix_markdown_key_pair_attachment ([#2460](https://github.com/aliyun/terraform-provider-alicloud/issues/2460))
- fix_finger_print ([#2463](https://github.com/aliyun/terraform-provider-alicloud/issues/2463))
- fix_markdown_instance_type_families ([#2465](https://github.com/aliyun/terraform-provider-alicloud/issues/2465))
- fix_markdown_alicloud_network_interface_attachment ([#2469](https://github.com/aliyun/terraform-provider-alicloud/issues/2469))

## 1.83.0 (May 15, 2020)

- **New Resource:** `alicloud_waf_instance` ([#2456](https://github.com/aliyun/terraform-provider-alicloud/issues/2456))
- **New Resource:** `alicloud_resource_manager_account` ([#2441](https://github.com/aliyun/terraform-provider-alicloud/issues/2441))
- **New Resource:** `alicloud_resource_manager_policy` ([#2439](https://github.com/aliyun/terraform-provider-alicloud/issues/2439))
- **New Resource:** `alicloud_resource_manager_handshake` ([#2432](https://github.com/aliyun/terraform-provider-alicloud/issues/2432))
- **New Resource:** `alicloud_cen_private_zone` ([#2421](https://github.com/aliyun/terraform-provider-alicloud/issues/2421))

BUG FIXES:

- fix_markdown_instance ([#2436](https://github.com/aliyun/terraform-provider-alicloud/issues/2436))

## 1.82.0 (May 08, 2020)

- **New Resource:** `alicloud_resource_manager_handshake` ([#2425](https://github.com/aliyun/terraform-provider-alicloud/issues/2425))
- **New Resource:** `alicloud_resource_manager_folder` ([#2425](https://github.com/aliyun/terraform-provider-alicloud/issues/2425))
- **New Resource:** `alicloud_resource_manager_resource_group` ([#2422](https://github.com/aliyun/terraform-provider-alicloud/issues/2422))
- **New Resource:** `alicloud_waf_domain` ([#2414](https://github.com/aliyun/terraform-provider-alicloud/issues/2414))
- **New Resource:** `alicloud_resource_manager_role` ([#2405](https://github.com/aliyun/terraform-provider-alicloud/issues/2405))
- **New Resource:** `alicloud_edas_application` ([#2384](https://github.com/aliyun/terraform-provider-alicloud/issues/2384))
- **New Resource:** `alicloud_edas_deploy_group` ([#2384](https://github.com/aliyun/terraform-provider-alicloud/issues/2384))
- **New Resource:** `alicloud_edas_application_scale` ([#2384](https://github.com/aliyun/terraform-provider-alicloud/issues/2384))
- **New Resource:** `alicloud_edas_slb_attachment` ([#2384](https://github.com/aliyun/terraform-provider-alicloud/issues/2384))
- **New Resource:** `alicloud_edas_cluster` ([#2384](https://github.com/aliyun/terraform-provider-alicloud/issues/2384))
- **New Resource:** `alicloud_edas_instance_cluster_attachment` ([#2384](https://github.com/aliyun/terraform-provider-alicloud/issues/2384))
- **New Resource:** `alicloud_edas_application_deployment` ([#2384](https://github.com/aliyun/terraform-provider-alicloud/issues/2384))
- **New Resource:** `alicloud_cen_route_map` ([#2371](https://github.com/aliyun/terraform-provider-alicloud/issues/2371))
- **Data Source:** `alicloud_edas_applications` ([#2384](https://github.com/aliyun/terraform-provider-alicloud/issues/2384))
- **Data Source:** `alicloud_edas_deploy_groups` ([#2384](https://github.com/aliyun/terraform-provider-alicloud/issues/2384))
- **Data Source:** `alicloud_edas_clusters` ([#2384](https://github.com/aliyun/terraform-provider-alicloud/issues/2384))

IMPROVEMENTS:

- ci supports edas and resourceManager dependencies ([#2424](https://github.com/aliyun/terraform-provider-alicloud/issues/2424))
- add missing "security_group_id" attribute declaration to schema ([#2417](https://github.com/aliyun/terraform-provider-alicloud/issues/2417))
- Update go sdk to 1.61.155 ([#2413](https://github.com/aliyun/terraform-provider-alicloud/issues/2413))
- optimized create emr cluster test case ([#2397](https://github.com/aliyun/terraform-provider-alicloud/issues/2397))

BUG FIXES:

- fix_markdown_instance  documentation ([#2430](https://github.com/aliyun/terraform-provider-alicloud/issues/2430))
- fix_markdown_slb_vpc  documentation ([#2429](https://github.com/aliyun/terraform-provider-alicloud/issues/2429))
- fix log audit document  documentation ([#2415](https://github.com/aliyun/terraform-provider-alicloud/issues/2415))
- Fix regression in writeToFile ([#2412](https://github.com/aliyun/terraform-provider-alicloud/issues/2412))

## 1.81.0 (May 01, 2020)

- **New Resource:** `alicloud_hbase_instance` ([#2395](https://github.com/aliyun/terraform-provider-alicloud/issues/2395))
- **New Resource:** `alicloud_adb_connection` ([#2392](https://github.com/aliyun/terraform-provider-alicloud/issues/2392))
- **New Resource:** `alicloud_cs_kubernetes` ([#2391](https://github.com/aliyun/terraform-provider-alicloud/issues/2391))
- **New Resource:** `alicloud_dms_enterprise_instance` ([#2390](https://github.com/aliyun/terraform-provider-alicloud/issues/2390))
- **Data Source:** `alicloud_polardb_node_classes` ([#2369](https://github.com/aliyun/terraform-provider-alicloud/issues/2369))

IMPROVEMENTS:

- Update go sdk to 1.61.155 ([#2413](https://github.com/aliyun/terraform-provider-alicloud/issues/2413))
- add test Parameter dependence ([#2402](https://github.com/aliyun/terraform-provider-alicloud/issues/2402))
- improve dms docs ([#2401](https://github.com/aliyun/terraform-provider-alicloud/issues/2401))
- Add sls log audit ([#2389](https://github.com/aliyun/terraform-provider-alicloud/issues/2389))
- Update CHANGELOG.md ([#2386](https://github.com/aliyun/terraform-provider-alicloud/issues/2386))
- return connection_string after polardb cluster created ([#2379](https://github.com/aliyun/terraform-provider-alicloud/issues/2379))

BUG FIXES:

- Fix regression in writeToFile ([#2412](https://github.com/aliyun/terraform-provider-alicloud/issues/2412))
- fix(log_audit): resolve cannot get region ([#2411](https://github.com/aliyun/terraform-provider-alicloud/issues/2411))
- Let WriteFile return write/delete errors ([#2408](https://github.com/aliyun/terraform-provider-alicloud/issues/2408))
- fix adb documents bug ([#2388](https://github.com/aliyun/terraform-provider-alicloud/issues/2388))
- fix del rds ins task ([#2387](https://github.com/aliyun/terraform-provider-alicloud/issues/2387))

## 1.80.1 (April 24, 2020)

IMPROVEMENTS:

- update emr tag resourceType from instance to cluste ([#2383](https://github.com/aliyun/terraform-provider-alicloud/issues/2383))
- improve(cen_instances): support tags ([#2376](https://github.com/aliyun/terraform-provider-alicloud/issues/2376))
- improve(sdk): upgraded the sdk and made compatibility ([#2373](https://github.com/aliyun/terraform-provider-alicloud/issues/2373))
- Update oss_bucket.html.markdown ([#2359](https://github.com/aliyun/terraform-provider-alicloud/issues/2359))

BUG FIXES:

- fix adb account & documents bug ([#2382](https://github.com/aliyun/terraform-provider-alicloud/issues/2382))
- fix(oss): fix the bug of setting a object acl with the wrong option ([#2366](https://github.com/aliyun/terraform-provider-alicloud/issues/2366))

## 1.80.0 (April 17, 2020)

- **New Resource:** `alicloud_dns_domain_attachmen` ([#2365](https://github.com/aliyun/terraform-provider-alicloud/issues/2365))
- **New Resource:** `alicloud_dns_instance` ([#2361](https://github.com/aliyun/terraform-provider-alicloud/issues/2361))
- **New Resource:** `alicloud_polardb_endpoint` ([#2321](https://github.com/aliyun/terraform-provider-alicloud/issues/2321))
- **Data Source:** `alicloud_dns_domain_txt_guid` ([#2357](https://github.com/aliyun/terraform-provider-alicloud/issues/2357))
- **Data Source:** `alicloud_kms_aliases` ([#2353](https://github.com/aliyun/terraform-provider-alicloud/issues/2353))

IMPROVEMENTS:

- improve(cen_instance): support tags ([#2374](https://github.com/aliyun/terraform-provider-alicloud/issues/2374))
- improve(rds_instance): remove checking zone id ([#2372](https://github.com/aliyun/terraform-provider-alicloud/issues/2372))
- ADB support scale in/out ([#2367](https://github.com/aliyun/terraform-provider-alicloud/issues/2367))
- improve(skd): upgraded the sdk and made compatibility ([#2363](https://github.com/aliyun/terraform-provider-alicloud/issues/2363))
- improve(cen_flowlogs): append output_file ([#2362](https://github.com/aliyun/terraform-provider-alicloud/issues/2362))
- remove checking instance type before creating ecs instance ang ess configuration ([#2358](https://github.com/aliyun/terraform-provider-alicloud/issues/2358))

BUG FIXES:

- Fix alikafka topic tag crate bug ([#2370](https://github.com/aliyun/terraform-provider-alicloud/issues/2370))
- Fix assign variable bug ([#2368](https://github.com/aliyun/terraform-provider-alicloud/issues/2368))
- fix(oss): fix creation_date info displays incorrectly bug ([#2364](https://github.com/aliyun/terraform-provider-alicloud/issues/2364))

## 1.79.0 (April 10, 2020)

- **New Resource:** `alicloud_auto_provisioning_group` ([#2303](https://github.com/aliyun/terraform-provider-alicloud/issues/2303))

IMPROVEMENTS:

- optimize retryable request for alikafka ([#2350](https://github.com/aliyun/terraform-provider-alicloud/issues/2350))
- Update data_source_alicloud_fc_triggers.go ([#2348](https://github.com/aliyun/terraform-provider-alicloud/issues/2348))
- Add error retry in delete method ([#2346](https://github.com/aliyun/terraform-provider-alicloud/issues/2346))
- improve(vpc): vpc and vswitch supported timeouts settings ([#2345](https://github.com/aliyun/terraform-provider-alicloud/issues/2345))
- Update fc_function.html.markdown ([#2342](https://github.com/aliyun/terraform-provider-alicloud/issues/2342))
- Update fc_trigger.html.markdown ([#2341](https://github.com/aliyun/terraform-provider-alicloud/issues/2341))
- improve(polardb): modify the vsw specified ([#2261](https://github.com/aliyun/terraform-provider-alicloud/issues/2261))
- eip associate add clientToken support ([#2247](https://github.com/aliyun/terraform-provider-alicloud/issues/2247))

BUG FIXES:

- fix(polardb): fix the bug of parameters modification ([#2354](https://github.com/aliyun/terraform-provider-alicloud/issues/2354))
- fix(ram): fix datasource ram users convince bug ([#2352](https://github.com/aliyun/terraform-provider-alicloud/issues/2352))
- fix(rds): fix the bug of parameters modification ([#2351](https://github.com/aliyun/terraform-provider-alicloud/issues/2351))
- fix(managed_kubernetes): resolve field version diff issue ([#2349](https://github.com/aliyun/terraform-provider-alicloud/issues/2349))
- private_zone has the wrong description ([#2344](https://github.com/aliyun/terraform-provider-alicloud/issues/2344))
- Fix create rds bug ([#2343](https://github.com/aliyun/terraform-provider-alicloud/issues/2343))
- fix(db_instance): resolve deleting db instance bug ([#2317](https://github.com/aliyun/terraform-provider-alicloud/issues/2317))
- user want api server's public ip when they set endpoint_public_access_enabled to true. the second parameter in DescribeClusterUserConfig is "privateIpAddress", so it should be endpoint_public_access_enabled's negative value ([#2290](https://github.com/aliyun/terraform-provider-alicloud/issues/2290))



## 1.78.0 (April 03, 2020)

- **New Resource:** `alicloud_log_alert` ([#2325](https://github.com/aliyun/terraform-provider-alicloud/issues/2325))
- **Data Source:** `alicloud_cen_flowlogs` ([#2336](https://github.com/aliyun/terraform-provider-alicloud/issues/2336))

IMPROVEMENTS:

- improve(cen_flowlogs): add more parameters ([#2338](https://github.com/aliyun/terraform-provider-alicloud/issues/2338))
- improve(mongodb): supported ssl setting ([#2335](https://github.com/aliyun/terraform-provider-alicloud/issues/2335))
- Add statistics attribute support ErrorCodeMaximum for cms ([#2328](https://github.com/aliyun/terraform-provider-alicloud/issues/2328))
- alicloud_kms_secret: mark secret_data as sensitive ([#2322](https://github.com/aliyun/terraform-provider-alicloud/issues/2322))

BUG FIXES:

- fix create rds instance bug ([#2334](https://github.com/aliyun/terraform-provider-alicloud/issues/2334))
- fix(nas_rule): resolve pagesize bug ([#2330](https://github.com/aliyun/terraform-provider-alicloud/issues/2330))

## 1.77.0 (April 01, 2020)

- **New Resource:** `alicloud_kms_alias` ([#2307](https://github.com/aliyun/terraform-provider-alicloud/issues/2307))
- **New Resource:** `alicloud_maxcompute_project` ([#1681](https://github.com/aliyun/terraform-provider-alicloud/issues/1681))

IMPROVEMENTS:

- improve(kms_secret): improve tags ([#2313](https://github.com/aliyun/terraform-provider-alicloud/issues/2313))

BUG FIXES:

- fix(ram_user): resolve importing ram user notfound bug ([#2319](https://github.com/aliyun/terraform-provider-alicloud/issues/2319))
- fix adb ci bug ([#2312](https://github.com/aliyun/terraform-provider-alicloud/issues/2312))
- fix(db_instance): resolve postgres testcase ([#2311](https://github.com/aliyun/terraform-provider-alicloud/issues/2311))

## 1.76.0 (March 27, 2020)

- **New Resource:** `alicloud_kms_secret` ([#2310](https://github.com/aliyun/terraform-provider-alicloud/issues/2310))

IMPROVEMENTS:

- change default autoscaler tag and refactor docs of managed kubernetes ([#2309](https://github.com/aliyun/terraform-provider-alicloud/issues/2309))
- improve(sdk): update provider sdk and make it compatible ([#2306](https://github.com/aliyun/terraform-provider-alicloud/issues/2306))
- add security group id and TDE for mongodb sharding ([#2298](https://github.com/aliyun/terraform-provider-alicloud/issues/2298))
- improve cen_instance and cen_flowlog ([#2297](https://github.com/aliyun/terraform-provider-alicloud/issues/2297))
- added support for isp_cities in site_monito ([#2296](https://github.com/aliyun/terraform-provider-alicloud/issues/2296))
- add security group id for kvstore ([#2292](https://github.com/aliyun/terraform-provider-alicloud/issues/2292)](https://github.com/aliyun/terraform-provider-alicloud/pull/2292))
- support parameter DesiredCapacity ([#2277](https://github.com/aliyun/terraform-provider-alicloud/issues/2277))

BUG FIXES:

-fix(cen):modify resource_alicloud_cen_instance_grant ([#2293](https://github.com/aliyun/terraform-provider-alicloud/issues/2293))


## 1.75.0 (March 20, 2020)

- **Data Source:** `alicloud_adb_zones` ([#2248](https://github.com/aliyun/terraform-provider-alicloud/issues/2248))
- **Data Source:** `alicloud_slb_zones` ([#2244](https://github.com/aliyun/terraform-provider-alicloud/issues/2244))
- **Data Source:** `alicloud_elasticsearch_zones` ([#2243](https://github.com/aliyun/terraform-provider-alicloud/issues/2243))

IMPROVEMENTS:

- imporve(db_instance): support force_restart ([#2287](https://github.com/aliyun/terraform-provider-alicloud/issues/2287))
- imporve the zones markdown ([#2285](https://github.com/aliyun/terraform-provider-alicloud/issues/2285))
- Add Terway and other kubernetes params to resource ([#2284](https://github.com/aliyun/terraform-provider-alicloud/issues/2284))

BUG FIXES:

- fix ADB data source zones ([#2283](https://github.com/aliyun/terraform-provider-alicloud/issues/2283))
- fix polardb data source zones ([#2274](https://github.com/aliyun/terraform-provider-alicloud/issues/2274))
- fix adb ci bug ([#2272](https://github.com/aliyun/terraform-provider-alicloud/issues/2272))
- fix(cen):modify resource_alicloud_cen_instance_attachment ([#2269](https://github.com/aliyun/terraform-provider-alicloud/issues/2269))

## 1.74.1 (March 17, 2020)

IMPROVEMENTS:

- improve(alikafka_instance): suspend kafka prepaid test ([#2264](https://github.com/aliyun/terraform-provider-alicloud/issues/2264))
- improve(gpdb): modify the vsw specified ([#2260](https://github.com/aliyun/terraform-provider-alicloud/issues/2260))
- improve(elasticsearch): modify the vsw specified ([#2239](https://github.com/aliyun/terraform-provider-alicloud/issues/2239))

BUG FIXES:

- tagResource bug fix ([#2266](https://github.com/aliyun/terraform-provider-alicloud/issues/2266))
- fix(kvstore_instance): resolve auto_renew incorrect value ([#2265](https://github.com/aliyun/terraform-provider-alicloud/issues/2265))

## 1.74.0 (March 16, 2020)

- **Data Source:** `alicloud_fc_zones` ([#2256](https://github.com/aliyun/terraform-provider-alicloud/issues/2256))
- **Data Source:** `alicloud_polardb_zones` ([#2250](https://github.com/aliyun/terraform-provider-alicloud/issues/2250))

IMPROVEMENTS:

- improve(hbase): modify the vsw specified ([#2259](https://github.com/aliyun/terraform-provider-alicloud/issues/2259))
- improve(elasticsearch): data_source support tags ([#2257](https://github.com/aliyun/terraform-provider-alicloud/issues/2257))
- rename polardb test name ([#2255](https://github.com/aliyun/terraform-provider-alicloud/issues/2255))
- corrct cen_flowlog docs ([#2254](https://github.com/aliyun/terraform-provider-alicloud/issues/2254))
- Adjust the return error mode ([#2252](https://github.com/aliyun/terraform-provider-alicloud/issues/2252))
- improve(elasticsearch): resource support tags ([#2251](https://github.com/aliyun/terraform-provider-alicloud/issues/2251))

BUG FIXES:

- fix(cms_alarm): resolve the effective_interval default value ([#2253](https://github.com/aliyun/terraform-provider-alicloud/issues/2253))

## 1.73.0 (March 13, 2020)

- **New Resource:** `alicloud_cen_flowlog` ([#2229](https://github.com/aliyun/terraform-provider-alicloud/issues/2229))
- **Data Source:** `alicloud_gpdb_zones` ([#2241](https://github.com/aliyun/terraform-provider-alicloud/issues/2241))
- **Data Source:** `alicloud_hbase_zones` ([#2240](https://github.com/aliyun/terraform-provider-alicloud/issues/2240))
- **Data Source:** `alicloud_mongodb_zones` ([#2238](https://github.com/aliyun/terraform-provider-alicloud/issues/2238))
- **Data Source:** `alicloud_kvstore_zones` ([#2236](https://github.com/aliyun/terraform-provider-alicloud/issues/2236))
- **Data Source:** `alicloud_db_zones` ([#2235](https://github.com/aliyun/terraform-provider-alicloud/issues/2235))

IMPROVEMENTS:

- improve(ecs): supported auto snapshop policy ([#2245](https://github.com/aliyun/terraform-provider-alicloud/issues/2245))
- add flowlog docs in the alicloud.erb ([#2237](https://github.com/aliyun/terraform-provider-alicloud/issues/2237))
- fix(elasticsearch): update the sdk ([#2234](https://github.com/aliyun/terraform-provider-alicloud/issues/2234))
- add new version aliyungo ([#2232](https://github.com/aliyun/terraform-provider-alicloud/issues/2232))
- terraform format examples [2231]
- Hbase tags ([#2228](https://github.com/aliyun/terraform-provider-alicloud/issues/2228))
- mongodb support TDE [GH2207]

BUG FIXES:

- fix(cms_alarm): resolve the effective_interval format bug ([#2242](https://github.com/aliyun/terraform-provider-alicloud/issues/2242))
- fix SQLServer testcase ([#2233](https://github.com/aliyun/terraform-provider-alicloud/issues/2233))
- fix(es): fix ci bug ([#2230](https://github.com/aliyun/terraform-provider-alicloud/issues/2230))

## 1.72.0 (March 06, 2020)

- **New Resource:** `alicloud_cms_site_monitor` ([#2191](https://github.com/aliyun/terraform-provider-alicloud/issues/2191))
- **Data Source:** `alicloud_ess_alarms` ([#2215](https://github.com/aliyun/terraform-provider-alicloud/issues/2215))
- **Data Source:** `alicloud_ess_notifications` ([#2161](https://github.com/aliyun/terraform-provider-alicloud/issues/2161))
- **Data Source:** `alicloud_ess_scheduled_tasks` ([#2160](https://github.com/aliyun/terraform-provider-alicloud/issues/2160))

IMPROVEMENTS:

- improve(mns_topic_subscription): remove the validate ([#2225](https://github.com/aliyun/terraform-provider-alicloud/issues/2225))
- Support the parameter of 'protocol' ([#2214](https://github.com/aliyun/terraform-provider-alicloud/issues/2214))
- improve sweeper test ([#2212](https://github.com/aliyun/terraform-provider-alicloud/issues/2212))
- supported bootstrap action when create a new emr cluster instance ([#2210](https://github.com/aliyun/terraform-provider-alicloud/issues/2210))

BUG FIXES:

- fix sweep test bug ([#2223](https://github.com/aliyun/terraform-provider-alicloud/issues/2223))
- fix the bug of RAM user cannot be destroyed ([#2219](https://github.com/aliyun/terraform-provider-alicloud/issues/2219))
- fix(elasticsearch_instance): resolve the ci bug ([#2216](https://github.com/aliyun/terraform-provider-alicloud/issues/2216))
- fix(slb): fix slb listener fields and rules creation bug ([#2209](https://github.com/aliyun/terraform-provider-alicloud/issues/2209))

## 1.71.2 (February 28, 2020)

IMPROVEMENTS:

- improve alikafka sweeper test ([#2206](https://github.com/aliyun/terraform-provider-alicloud/issues/2206))
- added filter parameter instance type about data source emr_instance_t…  ([#2205](https://github.com/aliyun/terraform-provider-alicloud/issues/2205))
- improve(polardb): fix update polardb cluster db_node_class will delete instance ([#2203](https://github.com/aliyun/terraform-provider-alicloud/issues/2203))
- improve(cen): add more sweeper test for cen and update go sdk ([#2201](https://github.com/aliyun/terraform-provider-alicloud/issues/2201))
- improve(mns_topic_subscription): supports json ([#2200](https://github.com/aliyun/terraform-provider-alicloud/issues/2200))
- update go sdk to 1.61.1 ([#2197](https://github.com/aliyun/terraform-provider-alicloud/issues/2197))
- improve(snat): add snat_entry_name for this resource ([#2196](https://github.com/aliyun/terraform-provider-alicloud/issues/2196))
- add sweeper for polardb and hbase ([#2195](https://github.com/aliyun/terraform-provider-alicloud/issues/2195))
- improve(nat_gateways): add output vpc_id ([#2194](https://github.com/aliyun/terraform-provider-alicloud/issues/2194))
- add retry for throttling when setting tags ([#2193](https://github.com/aliyun/terraform-provider-alicloud/issues/2193))
- improve(client): remove useless goSdkMutex ([#2192](https://github.com/aliyun/terraform-provider-alicloud/issues/2192))

BUG FIXES:

- fix(cms_alarm): resolve the creating rule dunplicated ([#2211](https://github.com/aliyun/terraform-provider-alicloud/issues/2211))
- fix(ess): fix create ess scaling group error ([#2208](https://github.com/aliyun/terraform-provider-alicloud/issues/2208))
- fix(ess): fix the bug of creating ess scaling group ([#2204](https://github.com/aliyun/terraform-provider-alicloud/issues/2204))
fix(common_bandwidth): resolve BandwidthPackageOperation.conflict ([#2199](https://github.com/aliyun/terraform-provider-alicloud/issues/2199))

## 1.71.1 (February 21, 2020)

IMPROVEMENTS:

- update SnatEntry test case ([#2187](https://github.com/aliyun/terraform-provider-alicloud/issues/2187))
- improve(vpcs): support outputting tags ([#2184](https://github.com/aliyun/terraform-provider-alicloud/issues/2184))
- improve(instance): remove sdk mutex and improve instance creating speed ([#2181](https://github.com/aliyun/terraform-provider-alicloud/issues/2181))
- (improve market_products): supports more filter parameter ([#2177](https://github.com/aliyun/terraform-provider-alicloud/issues/2177))
- add heyuan region and datasource market supports available_region ([#2176](https://github.com/aliyun/terraform-provider-alicloud/issues/2176))
- improve(ecs_instance): add tags into runInstances request ([#2175](https://github.com/aliyun/terraform-provider-alicloud/issues/2175))
- improve(ecs_instance): improve security groups ([#2174](https://github.com/aliyun/terraform-provider-alicloud/issues/2174))
- improve(fc_function): remove useless code ([#2173](https://github.com/aliyun/terraform-provider-alicloud/issues/2173))
- add support for create snat entry with source_dir ([#2170](https://github.com/aliyun/terraform-provider-alicloud/issues/2170))

BUG FIXES:

- fix(instance): resolve LastTokenProcessing error when modifying nework spec ([#2186](https://github.com/aliyun/terraform-provider-alicloud/issues/2186))
- fix(instance): resolve modifying network spec LastOrderProcessing error ([#2185](https://github.com/aliyun/terraform-provider-alicloud/issues/2185))
- fix(instance): resolve volume_tags diff bug when new resource ([#2182](https://github.com/aliyun/terraform-provider-alicloud/issues/2182))
- fix(image): fix the bug of created image by disk ([#2180](https://github.com/aliyun/terraform-provider-alicloud/issues/2180))
- Fix creating instance with multiple security groups ([#2168](https://github.com/aliyun/terraform-provider-alicloud/issues/2168))


## 1.71.0 (February 14, 2020)

- **New Resource:** `alicloud_adb_account` ([#2169](https://github.com/aliyun/terraform-provider-alicloud/issues/2169))
- **New Resource:** `alicloud_adb_backup_policy` ([#2169](https://github.com/aliyun/terraform-provider-alicloud/issues/2169))
- **Data Source:** `alicloud_adb_clusters` ([#2153](https://github.com/aliyun/terraform-provider-alicloud/issues/2153))

IMPROVEMENTS:

- add market product image id ([#2171](https://github.com/aliyun/terraform-provider-alicloud/issues/2171))
- fixed regional sts endpoint ([#2167](https://github.com/aliyun/terraform-provider-alicloud/issues/2167))
- improve(cms): add computed for effective_interval ([#2163](https://github.com/aliyun/terraform-provider-alicloud/issues/2163))

BUG FIXES:

- fix(ram_login_profile): resolve not found when deleting and deprecate alicloud_slb_attachment ([#2164](https://github.com/aliyun/terraform-provider-alicloud/issues/2164))
- fix(db_account_privilege): resolve privilege timeout bug on PosrgreSql ([#2159](https://github.com/aliyun/terraform-provider-alicloud/issues/2159))

## 1.70.3 (February 06, 2020)

IMPROVEMENTS:

- improve(db_instances): add more parameters ([#2158](https://github.com/aliyun/terraform-provider-alicloud/issues/2158))
- improve(kvstore_account): correct test case ([#2154](https://github.com/aliyun/terraform-provider-alicloud/issues/2154))

BUG FIXES:

- Update go SDK to fix redis bug ([#2149](https://github.com/aliyun/terraform-provider-alicloud/issues/2149))

## 1.70.2 (January 31, 2020)

IMPROVEMENTS:

- improve(slb): improve set slb tags ([#2147](https://github.com/aliyun/terraform-provider-alicloud/issues/2147))
- improve(ram_login_profile): resolve EntityNotExist.User ([#2146](https://github.com/aliyun/terraform-provider-alicloud/issues/2146))
- improve client endpoint ([#2144](https://github.com/aliyun/terraform-provider-alicloud/issues/2144))
- improve(client): add a method when load endpoint from local file  ([#2141](https://github.com/aliyun/terraform-provider-alicloud/issues/2141))
- improve(error): remove useless error codes ([#2140](https://github.com/aliyun/terraform-provider-alicloud/issues/2140))
- improve(provider): change IsExceptedError to IsExpectedErrors ([#2139](https://github.com/aliyun/terraform-provider-alicloud/issues/2139))
- improve(instance): remove the useless method ([#2138](https://github.com/aliyun/terraform-provider-alicloud/issues/2138))

BUG FIXES:

- fix(ram): resolve ram resources not found error ([#2143](https://github.com/aliyun/terraform-provider-alicloud/issues/2143))
- fix(slb): resolve listTags throttling ([#2142](https://github.com/aliyun/terraform-provider-alicloud/issues/2142))
- fix(instance): resolve the untag bug ([#2137](https://github.com/aliyun/terraform-provider-alicloud/issues/2137))
- fix vpn bug ([#2065](https://github.com/aliyun/terraform-provider-alicloud/issues/2065))

## 1.70.1 (January 23, 2020)

IMPROVEMENTS:

- added data source emr main versions parameter filter: cluster_type ([#2130](https://github.com/aliyun/terraform-provider-alicloud/issues/2130))
- Features/upgrade cluster ([#2129](https://github.com/aliyun/terraform-provider-alicloud/issues/2129))
- improve(mongodb_sharding): removing the limitation of node_storage ([#2128](https://github.com/aliyun/terraform-provider-alicloud/issues/2128))
- improve(hbase): add precheck for the test cases ([#2127](https://github.com/aliyun/terraform-provider-alicloud/issues/2127))
- Support update alikafka topic partition num and remark ([#2096](https://github.com/aliyun/terraform-provider-alicloud/issues/2096))

BUG FIXES:

- fix the bug of create emr gateway failed and optimized status delay time ([#2124](https://github.com/aliyun/terraform-provider-alicloud/issues/2124))

## 1.70.0 (January 17, 2020)

- **Data Source:** `alicloud_polardb_accounts` ([#2091](https://github.com/aliyun/terraform-provider-alicloud/issues/2091))
- **Data Source:** `alicloud_polardb_databases` ([#2091](https://github.com/aliyun/terraform-provider-alicloud/issues/2091))

IMPROVEMENTS:

- improve(slb_listener): add document for health_check_method ([#2121](https://github.com/aliyun/terraform-provider-alicloud/issues/2121))
- modify:cen_instance_grant.html.markdown("alicloud_cen_instance_grant.foo") ([#2120](https://github.com/aliyun/terraform-provider-alicloud/issues/2120))
- improve drds and rds sweeper test ([#2119](https://github.com/aliyun/terraform-provider-alicloud/issues/2119))
- improve(kvstore_instance_classes): add validateFunc for engine ([#2117](https://github.com/aliyun/terraform-provider-alicloud/issues/2117))
- improve(instance): close partial ([#2115](https://github.com/aliyun/terraform-provider-alicloud/issues/2115))
- improve(instance): supports setting auto_release_time#2095 ([#2105](https://github.com/aliyun/terraform-provider-alicloud/issues/2105))
- improve(slb) update slb_listener add health_check_method ([#2102](https://github.com/aliyun/terraform-provider-alicloud/issues/2102))
- improve(elasticsearch): resource support to renew a PrePaid instance ([#2099](https://github.com/aliyun/terraform-provider-alicloud/issues/2099))
- improve(rds): feature rds support sql audit records ([#2082](https://github.com/aliyun/terraform-provider-alicloud/issues/2082))

BUG FIXES:

- fix(drds_instance): resolve parsing response error ([#2118](https://github.com/aliyun/terraform-provider-alicloud/issues/2118))
- fix:cen_instance.html.markdown(modify docs of name & description) ([#2116](https://github.com/aliyun/terraform-provider-alicloud/issues/2116))
- fix(rds): fix rds modify sql collector policy bug ([#2110](https://github.com/aliyun/terraform-provider-alicloud/issues/2110))
- fix(rds): fix rds modify db instance spec bug ([#2108](https://github.com/aliyun/terraform-provider-alicloud/issues/2108))
- fix(drds_instance): resolve parsing failed when creating ([#2106](https://github.com/aliyun/terraform-provider-alicloud/issues/2106))
- fix(snat_entry): resolve the error OperationUnsupported.EipNatBWPCheck ([#2104](https://github.com/aliyun/terraform-provider-alicloud/issues/2104))
- fix(pvtz_zone): correct the docs error ([#2097](https://github.com/aliyun/terraform-provider-alicloud/issues/2097))

## 1.69.1 (January 13, 2020)

IMPROVEMENTS:

- improve(market): supported new field 'search_term' ([#2090](https://github.com/aliyun/terraform-provider-alicloud/issues/2090))

BUG FIXES:

- fix(instance_types): resolve a bug results from the filed spelling error ([#2093](https://github.com/aliyun/terraform-provider-alicloud/issues/2093))

## 1.69.0 (January 13, 2020)

- **New Resource:** `alicloud_market_order` ([#2084](https://github.com/aliyun/terraform-provider-alicloud/issues/2084))
- **New Resource:** `alicloud_image_import` ([#2051](https://github.com/aliyun/terraform-provider-alicloud/issues/2051))
- **Data Source:** `alicloud_market_product` ([#2070](https://github.com/aliyun/terraform-provider-alicloud/issues/2070))

IMPROVEMENTS:

- improve(api_gateway_group): add outputs sub_domain and vpc_domain ([#2088](https://github.com/aliyun/terraform-provider-alicloud/issues/2088))
- improve(hbase): expose the hbase docs ([#2087](https://github.com/aliyun/terraform-provider-alicloud/issues/2087))
- improve(pvtz): supports proxy_pattern, user_client_ip and lang ([#2086](https://github.com/aliyun/terraform-provider-alicloud/issues/2086))
- improve(test): support force sleep while running testcase ([#2081](https://github.com/aliyun/terraform-provider-alicloud/issues/2081))
- improve(emr): improve sweeper test for emr ([#2078](https://github.com/aliyun/terraform-provider-alicloud/issues/2078))
- improve(slb): update slb_server_certificate ([#2077](https://github.com/aliyun/terraform-provider-alicloud/issues/2077))
- improve(elasticsearch): resource elasticsearch_instance support update for instance_charge_type ([#2073](https://github.com/aliyun/terraform-provider-alicloud/issues/2073))
- improve(instance_types): support gpu amount and gpu spec ([#2069](https://github.com/aliyun/terraform-provider-alicloud/issues/2069))
- modify(market): modify the attributes of market products datasource ([#2068](https://github.com/aliyun/terraform-provider-alicloud/issues/2068))
- improve(listener): supports description ([#2067](https://github.com/aliyun/terraform-provider-alicloud/issues/2067))
- improve(testcase): change test image name_regex ([#2066](https://github.com/aliyun/terraform-provider-alicloud/issues/2066))
- improve(instances): improve its efficiency when fetching its disk mappings ([#2062](https://github.com/aliyun/terraform-provider-alicloud/issues/2062))
- improve(image): correct docs ([#2061](https://github.com/aliyun/terraform-provider-alicloud/issues/2061))
- improve(instances): supports ram_role_name ([#2060](https://github.com/aliyun/terraform-provider-alicloud/issues/2060))
- improve(db_instance): support security_group_ids ([#2056](https://github.com/aliyun/terraform-provider-alicloud/issues/2056))
- improve(api): added app_code in attribute apps ([#2055](https://github.com/aliyun/terraform-provider-alicloud/issues/2055))
- improve(rds): feature rds backup policy improve the functions ([#2042](https://github.com/aliyun/terraform-provider-alicloud/issues/2042))

BUG FIXES:

- fix(ram_roles): getRole not found error ([#2089](https://github.com/aliyun/terraform-provider-alicloud/issues/2089))
- fix(polardb): fix polardb add parameters bug ([#2083](https://github.com/aliyun/terraform-provider-alicloud/issues/2083))
- fix(ecs): fix the bug of ecs instance not supported ([#2063](https://github.com/aliyun/terraform-provider-alicloud/issues/2063))
- fix(image): fix image disk mapping size bug ([#2052](https://github.com/aliyun/terraform-provider-alicloud/issues/2052))


## 1.68.0 (January 06, 2020)

- **New Resource:** `alicloud_export_image` ([#2036](https://github.com/aliyun/terraform-provider-alicloud/issues/2036))
- **New Resource:** `alicloud_image_share_permission` ([#2026](https://github.com/aliyun/terraform-provider-alicloud/issues/2026))
- **New Resource:** `alicloud_polardb_endpoint_address` ([#2020](https://github.com/aliyun/terraform-provider-alicloud/issues/2020))
- **Data Source:** `alicloud_polardb_endpoints` ([#2020](https://github.com/aliyun/terraform-provider-alicloud/issues/2020))

IMPROVEMENTS:

- improve(db_readonly_instance): supports tags ([#2050](https://github.com/aliyun/terraform-provider-alicloud/issues/2050))
- improve(db_instance): support setting db_instance_storage_type ([#2048](https://github.com/aliyun/terraform-provider-alicloud/issues/2048))
- switch r-kvstore sdk to r_kvstor ([#2047](https://github.com/aliyun/terraform-provider-alicloud/issues/2047))
- ess-rules/groups/configrution three markdown files different with codes ([#2046](https://github.com/aliyun/terraform-provider-alicloud/issues/2046))
- improve(polardb): feature polardb support tags #2045
- improve(slb): update slb_server_certificate ([#2044](https://github.com/aliyun/terraform-provider-alicloud/issues/2044))
- Modify naming issues in alicloud_image_copy ([#2040](https://github.com/aliyun/terraform-provider-alicloud/issues/2040))
- rollback hbase resource and datasource because of them need to be improved ([#2034](https://github.com/aliyun/terraform-provider-alicloud/issues/2034))
- improve(sasl): Implement update method for sasl user ([#2027](https://github.com/aliyun/terraform-provider-alicloud/issues/2027))


BUG FIXES:

- fix(security_group): fix enterprise sg does not support inner access policy from issue #1961 @yongzhang ([#2049](https://github.com/aliyun/terraform-provider-alicloud/issues/2049))

## 1.67.0 (December 27, 2019)

- **New Resource:** `alicloud_hbase_instance` ([#2012](https://github.com/aliyun/terraform-provider-alicloud/issues/2012))
- **New Resource:** `alicloud_polardb_account_privilege` ([#2005](https://github.com/aliyun/terraform-provider-alicloud/issues/2005))
- **New Resource:** `alicloud_polardb_account` ([#1998](https://github.com/aliyun/terraform-provider-alicloud/issues/1998))
- **Data Source:** `alicloud_hbase_instances` ([#2012](https://github.com/aliyun/terraform-provider-alicloud/issues/2012))


IMPROVEMENTS:

- rollback hbase resource and datasource because of them need to be improved ([#2033](https://github.com/aliyun/terraform-provider-alicloud/issues/2033))
- improve(ci): add hbase ci ([#2031](https://github.com/aliyun/terraform-provider-alicloud/issues/2031))
- improve(mongodb): hidden some security ips ([#2025](https://github.com/aliyun/terraform-provider-alicloud/issues/2025))
- improve(cdn_config): remove the args private_oss_tbl ([#2024](https://github.com/aliyun/terraform-provider-alicloud/issues/2024))
- improve(ons): update sdk and remove PreventCache in ons ([#2014](https://github.com/aliyun/terraform-provider-alicloud/issues/2014))
- add resource group id ([#2010](https://github.com/aliyun/terraform-provider-alicloud/issues/2010))
- improve(kvstore): remove type 'Super' for kvstore ([#2009](https://github.com/aliyun/terraform-provider-alicloud/issues/2009))
- improve(emr): support emr cluster tag ([#2008](https://github.com/aliyun/terraform-provider-alicloud/issues/2008))
- feat(alicloud/yundun_bastionhost): Add support for Cloud Bastionhost ([#2006](https://github.com/aliyun/terraform-provider-alicloud/issues/2006))

BUG FIXES:

- fix(rds): fix policy sqlserver test bug ([#2030](https://github.com/aliyun/terraform-provider-alicloud/issues/2030))
- fix(rds): fix rds resource alicloud_db_backup_policy`s log_backup bug ([#2017](https://github.com/aliyun/terraform-provider-alicloud/issues/2017))
- fix(db_backup_policy): fix postgresql backup policy bug ([#2003](https://github.com/aliyun/terraform-provider-alicloud/issues/2003))


## 1.66.0 (December 20, 2019)

- **New Resource:** `alicloud_kvstore_account` ([#1993](https://github.com/aliyun/terraform-provider-alicloud/issues/1993))
- **New Resource:** `alicloud_copy_image` ([#1978](https://github.com/aliyun/terraform-provider-alicloud/issues/1978))
- **New Resource:** `alicloud_polardb_database` ([#1996](https://github.com/aliyun/terraform-provider-alicloud/issues/1996))
- **New Resource:** `alicloud_polardb_backup_policy` ([#1991](https://github.com/aliyun/terraform-provider-alicloud/issues/1991))
- **New Resource:** `alicloud_poloardb_cluster` ([#1978](https://github.com/aliyun/terraform-provider-alicloud/issues/1978))
- **New Resource:** `alicloud_alikafka_sasl_acl` （[#2000](https://github.com/aliyun/terraform-provider-alicloud/issues/2000)
- **New Resource:** `alicloud_alikafka_sasl_user` （[#2000](https://github.com/aliyun/terraform-provider-alicloud/issues/2000)
- **Data Source:** `alicloud_poloardb_clusters` ([#1978](https://github.com/aliyun/terraform-provider-alicloud/issues/1978))
- **Data Source:** `alicloud_alikafka_sasl_acls` （[#2000](https://github.com/aliyun/terraform-provider-alicloud/issues/2000)）
- **Data Source:** `alicloud_alikafka_sasl_users`（[#2000](https://github.com/aliyun/terraform-provider-alicloud/issues/2000)）

IMPROVEMENTS:


- improve(SLS): Support SLS logstore index json keys ([#1999](https://github.com/aliyun/terraform-provider-alicloud/issues/1999))
- improve(acl): add missing tags for acl, keypair and son ([#1997](https://github.com/aliyun/terraform-provider-alicloud/issues/1997))
- improve(market): product datasource supported name regex and ids filter ([#1992](https://github.com/aliyun/terraform-provider-alicloud/issues/1992))
- improve(instance): improve auto_renew_period setting ([#1990](https://github.com/aliyun/terraform-provider-alicloud/issues/1990))
- documenting the replica_set_name attribute in mongodb @chanind ([#1989](https://github.com/aliyun/terraform-provider-alicloud/issues/1989))
- improve(period): improve computing period method ([#1988](https://github.com/aliyun/terraform-provider-alicloud/issues/1988))
- improve(prepaid): support computing period by week ([#1985](https://github.com/aliyun/terraform-provider-alicloud/issues/1985))
- improve(prepaid): add a method to fix period importer diff ([#1984](https://github.com/aliyun/terraform-provider-alicloud/issues/1984))
- improve(mongoDB): supported field 'tags' ([#1980](https://github.com/aliyun/terraform-provider-alicloud/issues/1980))
- add output to ssl vpn client cert @ionTea ([#1979](https://github.com/aliyun/terraform-provider-alicloud/issues/1979))

BUG FIXES:

- fix(db_backup_policy): fix postgresql backup policy bug ([#2002](https://github.com/aliyun/terraform-provider-alicloud/issues/2002))
- fix(fc): fixed bug from issue #1961 @yokzy88 ([#1987](https://github.com/aliyun/terraform-provider-alicloud/issues/1987))
- fix(vpn): fix bug from issue #1965 @chanind ([#1981](https://github.com/aliyun/terraform-provider-alicloud/issues/1981))
- fix(vpn): added Computed to field `vswitch_id` @chanind ([#1977](https://github.com/aliyun/terraform-provider-alicloud/issues/1977))

## 1.65.0 (December 13, 2019)

- **New Resource:** `alicloud_reserved_instance` ([#1967](https://github.com/aliyun/terraform-provider-alicloud/issues/1967))
- **New Resource:** `alicloud_cs_kubernetes_autoscaler` ([#1956](https://github.com/aliyun/terraform-provider-alicloud/issues/1956))
- **New Data Source:** `alicloud_caller_identity` ([#1944](https://github.com/aliyun/terraform-provider-alicloud/issues/1944))
- **New Resource:** `alicloud_sag_client_user` ([#1807](https://github.com/aliyun/terraform-provider-alicloud/issues/1807))

IMPROVEMENTS:

- improve(ess_vserver_groups): improve docs ([#1976](https://github.com/aliyun/terraform-provider-alicloud/issues/1976))
- improve(kvstore_instance): set period using createtime and endtime ([#1971](https://github.com/aliyun/terraform-provider-alicloud/issues/1971))
- improve(slb_server_group): set servers to computed and avoid diff when using ess_scaling_vserver_group ([#1970](https://github.com/aliyun/terraform-provider-alicloud/issues/1970))
- improve(k8s):add AccessKey and AccessKeySecret instead of RamRole ([#1966](https://github.com/aliyun/terraform-provider-alicloud/issues/1966))

BUG FIXES:

- fix(market): remove the suggested_price check for avoid the error of testcase ([#1964](https://github.com/aliyun/terraform-provider-alicloud/issues/1964))
- fix(provider): Resolve issues from aone. ([#1963](https://github.com/aliyun/terraform-provider-alicloud/issues/1963))
- fix(market): remove the tags check for avoid the error of testcase ([#1962](https://github.com/aliyun/terraform-provider-alicloud/issues/1962))
- fix(alicloud/yundun_bastionhost): Bastionhost RAM policy authorization bug fix([#1960](https://github.com/aliyun/terraform-provider-alicloud/issues/1960))
- fix(datahub): fix updating datahub topic comment bug ([#1959](https://github.com/aliyun/terraform-provider-alicloud/issues/1959))
- fix(kms): correct test case errors  ([#1958](https://github.com/aliyun/terraform-provider-alicloud/issues/1958))
- fix(validator): update package github.com/denverdino/aliyungo/cdn ([#1946](https://github.com/aliyun/terraform-provider-alicloud/issues/1946))

## 1.64.0 (December 06, 2019)

- **New Data Source:** `alicloud_market_products` ([#1941](https://github.com/aliyun/terraform-provider-alicloud/issues/1941))
- **New Resource:** `alicloud_cloud_connect_network_attachment` ([#1933](https://github.com/aliyun/terraform-provider-alicloud/issues/1933))
- **New Resource:** `alicloud_image` ([#1913](https://github.com/aliyun/terraform-provider-alicloud/issues/1913))

IMPROVEMENTS:

- improve(docs): improve module guide ([#1957](https://github.com/aliyun/terraform-provider-alicloud/issues/1957))
- improve(db_account_privilege): supports more privileges ([#1945](https://github.com/aliyun/terraform-provider-alicloud/issues/1945))
- improve(datasources): remove sorted_by testcase results from some internal limitation ([#1943](https://github.com/aliyun/terraform-provider-alicloud/issues/1943))
- improve(sdk): Updated sdk to v1.60.280 and modified drds fields ([#1938](https://github.com/aliyun/terraform-provider-alicloud/issues/1938))
- improve(snat): update example to support for snat's creation with multi eips ([#1931](https://github.com/aliyun/terraform-provider-alicloud/issues/1931))
- improve(ess): resource alicloud_ess_scalinggroup_vserver_groups support parameter ([#1919](https://github.com/aliyun/terraform-provider-alicloud/issues/1919))
- improve(db_instance): make 'instance_types' 'db_instance_class' 'kvstore_instance_class' support price ([#1749](https://github.com/aliyun/terraform-provider-alicloud/issues/1749))

BUG FIXES:

- fix(alikafka): fix bug in when doing alikafka instance multi acc test ([#1947](https://github.com/aliyun/terraform-provider-alicloud/issues/1947))
- fix(CSKubernetes): fix 3az test case ([#1942](https://github.com/aliyun/terraform-provider-alicloud/issues/1942))
- fix(cdn_domain_new): constant timeout waiting for server cert ([#1937](https://github.com/aliyun/terraform-provider-alicloud/issues/1937))
- fix(pvtz_zone_record): allow SRV records ([#1936](https://github.com/aliyun/terraform-provider-alicloud/issues/1936))
- fix(Serverless Kubernetes): fix #1867 add serverless kube_config ([#1923](https://github.com/aliyun/terraform-provider-alicloud/issues/1923))

## 1.63.0 (December 02, 2019)

- **New Resource:** `alicloud_cloud_connect_network_grant` ([#1921](https://github.com/aliyun/terraform-provider-alicloud/issues/1921))
- **New Data Source:** `alicloud_yundun_bastionhost_instances` ([#1894](https://github.com/aliyun/terraform-provider-alicloud/issues/1894))
- **New Resource:** `alicloud_yundun_bastionhost_instance` ([#1894](https://github.com/aliyun/terraform-provider-alicloud/issues/1894))
- **New Data Source:** `alicloud_kms_ciphertext` ([#1858](https://github.com/aliyun/terraform-provider-alicloud/issues/1858))
- **New Data Source:** `alicloud_kms_plaintext` ([#1858](https://github.com/aliyun/terraform-provider-alicloud/issues/1858))
- **New Resource:** `alicloud_kms_ciphertext` ([#1858](https://github.com/aliyun/terraform-provider-alicloud/issues/1858))
- **New Resource:** `alicloud_sag_dnat_entry` ([#1823](https://github.com/aliyun/terraform-provider-alicloud/issues/1823))

IMPROVEMENTS:

- improve(vpc): add module support for vpc, vswitch and route entry ([#1934](https://github.com/aliyun/terraform-provider-alicloud/issues/1934))
- improve(db_instance): tags supports case sensitive ([#1930](https://github.com/aliyun/terraform-provider-alicloud/issues/1930))
- improve(mongodb_instance): adding replica_set_name to output from alicloud_mongodb_instance ([#1929](https://github.com/aliyun/terraform-provider-alicloud/issues/1929))
- improve(slb): add a new field delete_protection_validation ([#1927](https://github.com/aliyun/terraform-provider-alicloud/issues/1927))
- improve(kms): improve kms testcases use new method ([#1926](https://github.com/aliyun/terraform-provider-alicloud/issues/1926))
- improve(provider): added 'Computed : true' to all 'ids' fields. ([#1924](https://github.com/aliyun/terraform-provider-alicloud/issues/1924))
- improve(validator): Delete TagNum Count ([#1920](https://github.com/aliyun/terraform-provider-alicloud/issues/1920))
- improve(sag_dnat_entry): modify docs "add subcategory" ([#1918](https://github.com/aliyun/terraform-provider-alicloud/issues/1918))
- improve(sdk): upgrade alibaba go sdk ([#1917](https://github.com/aliyun/terraform-provider-alicloud/issues/1917))
- improve(db_database):update db_database doc ([#1916](https://github.com/aliyun/terraform-provider-alicloud/issues/1916))
- improve(validator): shift validator to offical ones ([#1912](https://github.com/aliyun/terraform-provider-alicloud/issues/1912))
- improve(alikafka): Support pre paid instance & Support tag resource ([#1873](https://github.com/aliyun/terraform-provider-alicloud/issues/1873))

BUG FIXES:

- fix(dns_record): fix dns_record testcase bug ([#1892](https://github.com/aliyun/terraform-provider-alicloud/issues/1892))
- fix(ecs): FIX: query system disk does not exist because no resource_group_id is specified ([#1884](https://github.com/aliyun/terraform-provider-alicloud/issues/1884))

## 1.62.2 (November 26, 2019)

IMPROVEMENTS:

- improve(mongodb): feature mongodb support postpaid to prepaid ([#1908](https://github.com/aliyun/terraform-provider-alicloud/issues/1908))
- improve(kvstore_instance_classes): skip unsupported zones ([#1901](https://github.com/aliyun/terraform-provider-alicloud/issues/1901))

BUG FIXES:

- fix(pvtz_attachment): fix vpc_ids diff error ([#1911](https://github.com/aliyun/terraform-provider-alicloud/issues/1911))
- fix(kafka): remove the const endpoint ([#1910](https://github.com/aliyun/terraform-provider-alicloud/issues/1910))
- fix(ess): modify the type of from Set to List. ([#1905](https://github.com/aliyun/terraform-provider-alicloud/issues/1905))
- fix managedkubernetes demo  documentation ([#1903](https://github.com/aliyun/terraform-provider-alicloud/issues/1903))
- fix the bug of TestAccAlicloudEmrCluster_local_storage failed ([#1902](https://github.com/aliyun/terraform-provider-alicloud/issues/1902))
- fix(db_instance): fix postgre testcase ([#1899](https://github.com/aliyun/terraform-provider-alicloud/issues/1899))
- fix(db_instance): test case ([#1898](https://github.com/aliyun/terraform-provider-alicloud/issues/1898))

## 1.62.1 (November 22, 2019)

IMPROVEMENTS:

- improve(db_instance): add new field auto_upgrade_minor_version to set minor version ([#1897](https://github.com/aliyun/terraform-provider-alicloud/issues/1897))
- imporve(docs): add AODC warning ([#1893](https://github.com/aliyun/terraform-provider-alicloud/issues/1893))
- improve(kvstore_instance): correct its docs ([#1891](https://github.com/aliyun/terraform-provider-alicloud/issues/1891))
- improve(pvtz_zone_attachment):make pvtz_zone_attachment support different region vpc ([#1890](https://github.com/aliyun/terraform-provider-alicloud/issues/1890))
- improve(db, kvstore): add auto pay when changing instance charge type ([#1889](https://github.com/aliyun/terraform-provider-alicloud/issues/1889))
- improve(cs): Do not assume `private_zone` is returned from API ([#1885](https://github.com/aliyun/terraform-provider-alicloud/issues/1885))
- improve(cs): modify the value of 'new_nat_gateway' to avoid errors. ([#1882](https://github.com/aliyun/terraform-provider-alicloud/issues/1882))
- improve(docs): Terraform registry docs ([#1881](https://github.com/aliyun/terraform-provider-alicloud/issues/1881))
- improve(rds): feature support high security access mode not submitted ([#1880](https://github.com/aliyun/terraform-provider-alicloud/issues/1880))
- improve(oss): add transitions to life-cycle ([#1879](https://github.com/aliyun/terraform-provider-alicloud/issues/1879))
- improve(db_instance): feature support high security access mode not submitted ([#1878](https://github.com/aliyun/terraform-provider-alicloud/issues/1878))
- improve(scalingconfiguration): supports changing password_inherit ([#1877](https://github.com/aliyun/terraform-provider-alicloud/issues/1877))
- improve(zones): use describeAvailableResource API to get rds available zones ([#1876](https://github.com/aliyun/terraform-provider-alicloud/issues/1876))
- improve(kvstore_instance_classess): improve test case error caused by no stock ([#1875](https://github.com/aliyun/terraform-provider-alicloud/issues/1875))
- improve(mns): mns support sts access ([#1871](https://github.com/aliyun/terraform-provider-alicloud/issues/1871))
- improve(elasticsearch): Added retry to avoid CreateInstance TokenPreviousRequestProcessError error ([#1870](https://github.com/aliyun/terraform-provider-alicloud/issues/1870))
- improve(kvstore_instance_engines): improve its code ([#1864](https://github.com/aliyun/terraform-provider-alicloud/issues/1864))
- improve(kvstore): remove memcache filter from datasource test [[#1863](https://github.com/aliyun/terraform-provider-alicloud/issues/1863)] 
- improve(oss_bucket_object):make oss_bucket_object support KMS encryption ([#1860](https://github.com/aliyun/terraform-provider-alicloud/issues/1860))
- improve(provider): added endpoint for resources. ([#1855](https://github.com/aliyun/terraform-provider-alicloud/issues/1855))

BUG FIXES:

- fix(db_backup_policy): add limitation when modify sqlservr policy ([#1896](https://github.com/aliyun/terraform-provider-alicloud/issues/1896))
- fix(kvstore): remove the kvstore instance password limitation ([#1886](https://github.com/aliyun/terraform-provider-alicloud/issues/1886))
- fix(mongodb_instances): fix filetering bug ([#1874](https://github.com/aliyun/terraform-provider-alicloud/issues/1874))
- fix(mongodb_instances): fix name_regex bug ([#1865](https://github.com/aliyun/terraform-provider-alicloud/issues/1865))
- fix(key_pair):fix key_pair testcase bug ([#1862](https://github.com/aliyun/terraform-provider-alicloud/issues/1862))
- fix(autoscaling): fix autoscaling bugs. ([#1832](https://github.com/aliyun/terraform-provider-alicloud/issues/1832))

## 1.62.0 (November 13, 2019)

- **New Resource:** `alicloud_yundun_dbaudit_instance` ([#1819](https://github.com/aliyun/terraform-provider-alicloud/issues/1819))
- **New Data Source:** `alicloud_yundun_dbaudit_instances` ([#1819](https://github.com/aliyun/terraform-provider-alicloud/issues/1819))

IMPROVEMENTS:

- improve(ess_scalingconfiguration): support password_inherit ([#1856](https://github.com/aliyun/terraform-provider-alicloud/issues/1856))
- improve docs and add ci for yundun dbaudit ([#1853](https://github.com/aliyun/terraform-provider-alicloud/issues/1853))

BUG FIXES:

- fix(provider): fix the bug: slice bounds out of range ([#1854](https://github.com/aliyun/terraform-provider-alicloud/issues/1854))

## 1.61.0 (November 12, 2019)

- **New Resource:** `alicloud_sag_snat_entry` ([#1799](https://github.com/aliyun/terraform-provider-alicloud/issues/1799))

IMPROVEMENTS:

- improve(provider): add default value for configuration_source ([#1852](https://github.com/aliyun/terraform-provider-alicloud/issues/1852))
- improve(ess): add module guide for the ess resources ([#1850](https://github.com/aliyun/terraform-provider-alicloud/issues/1850))
- improve(instance): postpaid instance supported 'dry_run' ([#1845](https://github.com/aliyun/terraform-provider-alicloud/issues/1845))
- improve(rds): fix for hidden dts ip list ([#1844](https://github.com/aliyun/terraform-provider-alicloud/issues/1844))
- improve(resource_alicloud_db_database): support Mohawk_100_BIN ([#1838](https://github.com/aliyun/terraform-provider-alicloud/issues/1838))
- perf(alicloud_db_backup_policy,db_instances):perf rds document desc ([#1836](https://github.com/aliyun/terraform-provider-alicloud/issues/1836))
- change sideBar ([#1830](https://github.com/aliyun/terraform-provider-alicloud/issues/1830))
- modify CloudConnectNetwork_multi ([#1828](https://github.com/aliyun/terraform-provider-alicloud/issues/1828))
- improve(alikafka): Added name for vpcs and vswitches ([#1827](https://github.com/aliyun/terraform-provider-alicloud/issues/1827))
- support to create emr gateway cluster instance ([#1821](https://github.com/aliyun/terraform-provider-alicloud/issues/1821))

BUG FIXES:

- fix(cs_kubenrnetes): fix terraform docs documentation ([#1851](https://github.com/aliyun/terraform-provider-alicloud/issues/1851))
- fix(nat_gateway):fix nat_gateway period bug ([#1841](https://github.com/aliyun/terraform-provider-alicloud/issues/1841))
- fix waitfor method nil bug ([#1840](https://github.com/aliyun/terraform-provider-alicloud/issues/1840))
- fix(ess): use GetOkExists to avoid some potential bugs ([#1835](https://github.com/aliyun/terraform-provider-alicloud/issues/1835))

## 1.60.0 (November 01, 2019)

- **New Data Source:** `alicloud_emr_disk_types` ([#1805](https://github.com/aliyun/terraform-provider-alicloud/issues/1805))
- **New Data Source:** `alicloud_dns_resolution_lines` ([#1800](https://github.com/aliyun/terraform-provider-alicloud/issues/1800))
- **New Resource:** `alicloud_sag_qos` ([#1790](https://github.com/aliyun/terraform-provider-alicloud/issues/1790))
- **New Resource:** `alicloud_sag_qos_policy` ([#1790](https://github.com/aliyun/terraform-provider-alicloud/issues/1790))
- **New Resource:** `alicloud_sag_qos_car` ([#1790](https://github.com/aliyun/terraform-provider-alicloud/issues/1790))
- **New Resource:** `alicloud_sag_acl` ([#1788](https://github.com/aliyun/terraform-provider-alicloud/issues/1788))
- **New Resource:** `alicloud_sag_acl_rule` ([#1788](https://github.com/aliyun/terraform-provider-alicloud/issues/1788))
- **New Data Source:** `alicloud_sag_acls` ([#1788](https://github.com/aliyun/terraform-provider-alicloud/issues/1788))
- **New Resource:** `alicloud_slb_domain_extension` ([#1756](https://github.com/aliyun/terraform-provider-alicloud/issues/1756))
- **New Data Source:** `alicloud_slb_domain_extensions` ([#1756](https://github.com/aliyun/terraform-provider-alicloud/issues/1756))

IMPROVEMENTS:

- alicloud_ess_scheduled_task supports Cron type ([#1824](https://github.com/aliyun/terraform-provider-alicloud/issues/1824))
- vpc product datasource support resource_group_id ([#1822](https://github.com/aliyun/terraform-provider-alicloud/issues/1822))
- imporve(instance): modified the argument reference in doc. ([#1815](https://github.com/aliyun/terraform-provider-alicloud/issues/1815))
- Add resource_group_id to data_source_alicloud_route_tables ([#1814](https://github.com/aliyun/terraform-provider-alicloud/issues/1814))
- use homedir to expand shared_credentials_file value and add environment variable for it ([#1811](https://github.com/aliyun/terraform-provider-alicloud/issues/1811))
- Add password to resource_alicloud_ess_scalingconfiguration ([#1810](https://github.com/aliyun/terraform-provider-alicloud/issues/1810))
- add ids for db_instance_classess and remove limitation for db_database resource ([#1803](https://github.com/aliyun/terraform-provider-alicloud/issues/1803))
- improve(db_instances):update tags type from string to map ([#1802](https://github.com/aliyun/terraform-provider-alicloud/issues/1802))
- improve(instance): field 'user_data' supported update ([#1798](https://github.com/aliyun/terraform-provider-alicloud/issues/1798))
- add doc of cloud_connect_network ([#1791](https://github.com/aliyun/terraform-provider-alicloud/issues/1791))
- improve(slb): updated slb attachment testcase. ([#1758](https://github.com/aliyun/terraform-provider-alicloud/issues/1758))

BUG FIXES:

- fix(tag): fix api gw, gpdb, kvstore datasource bug ([#1817](https://github.com/aliyun/terraform-provider-alicloud/issues/1817))
- fix(rds): fix creating db account empty pointer bug ([#1812](https://github.com/aliyun/terraform-provider-alicloud/issues/1812))
- fix(slb_listener): fix server_certificate_id diff bug and add sag ci([#1808](https://github.com/aliyun/terraform-provider-alicloud/issues/1808))
- fix(vpc): fix DescribeTag bug for vpc's datasource ([#1801](https://github.com/aliyun/terraform-provider-alicloud/issues/1801))

## 1.59.0 (October 25, 2019)

- **New Resource:** `alicloud_cloud_connect_network` ([#1784](https://github.com/aliyun/terraform-provider-alicloud/issues/1784))
- **New Resource:** `alicloud_alikafka_instance` ([#1764](https://github.com/aliyun/terraform-provider-alicloud/issues/1764))
- **New Data Source:** `alicloud_cloud_connect_networks` ([#1784](https://github.com/aliyun/terraform-provider-alicloud/issues/1784))
- **New Data Source:** `alicloud_emr_instance_types` ([#1773](https://github.com/aliyun/terraform-provider-alicloud/issues/1773))
- **New Data Source:** `alicloud_emr_main_versions` [[#1773](https://github.com/aliyun/terraform-provider-alicloud/issues/1773)] 
- **New Data Source:** `alicloud_alikafka_instances` ([#1764](https://github.com/aliyun/terraform-provider-alicloud/issues/1764))
- **New Data Source:** `alicloud_file_crc64_checksum` ([#1722](https://github.com/aliyun/terraform-provider-alicloud/issues/1722))

IMPROVEMENTS:

- improve(slb_listener): deprecate ssl_certificate_id and use server_certificate_id instead ([#1797](https://github.com/aliyun/terraform-provider-alicloud/issues/1797))
- improve(slb): improve slb docs ([#1796](https://github.com/aliyun/terraform-provider-alicloud/issues/1796))
- improve(slb_listener): add retry for StartLoadBalancerListener ([#1794](https://github.com/aliyun/terraform-provider-alicloud/issues/1794))
- improve(fc_trigger):change testcase dependence resource cdn_domain to new ([#1793](https://github.com/aliyun/terraform-provider-alicloud/issues/1793))
- improve(zones): using describeAvailableResource instead of DescribeZones for RKvstore ([#1789](https://github.com/aliyun/terraform-provider-alicloud/issues/1789))
- Update ssl_vpn_server.html.markdown ([#1786](https://github.com/aliyun/terraform-provider-alicloud/issues/1786))
- add resource_group_id to dns ([#1781](https://github.com/aliyun/terraform-provider-alicloud/issues/1781))
- improve(provider): modified the kms field conflict to diffsuppress ([#1780](https://github.com/aliyun/terraform-provider-alicloud/issues/1780))
- Always set PolicyDocument for RAM policy update ([#1777](https://github.com/aliyun/terraform-provider-alicloud/issues/1777))
- rename cs_serveless_kubernetes to cs_serverless_kubernetes ([#1776](https://github.com/aliyun/terraform-provider-alicloud/issues/1776))
- improve(slb): updated slb server_group testcase ([#1753](https://github.com/aliyun/terraform-provider-alicloud/issues/1753))
- improve(fc_function):support code_checksum ([#1722](https://github.com/aliyun/terraform-provider-alicloud/issues/1722))

BUG FIXES:

- fix(slb): address_type diff bug ([#1795](https://github.com/aliyun/terraform-provider-alicloud/issues/1795))
- fix(ddosbgp): the docs error ([#1782](https://github.com/aliyun/terraform-provider-alicloud/issues/1782))
- fix(instance):fix credit_specification bug ([#1778](https://github.com/aliyun/terraform-provider-alicloud/issues/1778))

## 1.58.1 (October 22, 2019)

IMPROVEMENTS:

- add missing resource ddosbgp_instance docs index ([#1775](https://github.com/aliyun/terraform-provider-alicloud/issues/1775))

BUG FIXES:

- fix(common_bandwidth_package): fix common bandwidth package resource_group_id forcenew bug ([#1772](https://github.com/aliyun/terraform-provider-alicloud/issues/1772))

## 1.58.0 (October 18, 2019)

- **New Data Source:** `alicloud_cs_serverless_kubernetes_clusters` ([#1746](https://github.com/aliyun/terraform-provider-alicloud/issues/1746))
- **New Resource:** `alicloud_cs_serverless_kubernetes` ([#1746](https://github.com/aliyun/terraform-provider-alicloud/issues/1746))

IMPROVEMENTS:

- Make `resource_group_id` to computed ([#1771](https://github.com/aliyun/terraform-provider-alicloud/issues/1771))
- Add tag for `resource_group_id` in the docs ([#1770](https://github.com/aliyun/terraform-provider-alicloud/issues/1770))
- add resource_group_id to vpc, slb resources and data sources and revise corresponding docs ([#1769](https://github.com/aliyun/terraform-provider-alicloud/issues/1769))
- improve(security_group):make security_group support resource_group_id ([#1762](https://github.com/aliyun/terraform-provider-alicloud/issues/1762))
- Add resource_group_id to common_bandwidth_package(resource&data_source) ([#1761](https://github.com/aliyun/terraform-provider-alicloud/issues/1761))
- improve(cen): added precheck for testcases ([#1759](https://github.com/aliyun/terraform-provider-alicloud/issues/1759))
- improve(security_group):support security_group_type ([#1755](https://github.com/aliyun/terraform-provider-alicloud/issues/1755))
- Add missing routing rules for alicloud_dns_record ([#1754](https://github.com/aliyun/terraform-provider-alicloud/issues/1754))
- improve(slb): updated slb serverCertificate testcase ([#1751](https://github.com/aliyun/terraform-provider-alicloud/issues/1751))
- improve(slb): updated slb rule testcase ([#1748](https://github.com/aliyun/terraform-provider-alicloud/issues/1748))
- Improve(alicloud_ess_scaling_rule): support TargetTrackingScalingRule and StepScalingRule ([#1744](https://github.com/aliyun/terraform-provider-alicloud/issues/1744))
- improve(cdn): added adddebug for tags APIs ([#1741](https://github.com/aliyun/terraform-provider-alicloud/issues/1741))
- improve(slb): updated slb ca_certificate testcase ([#1740](https://github.com/aliyun/terraform-provider-alicloud/issues/1740))
- improve(slb): updated slb acl testcase ([#1739](https://github.com/aliyun/terraform-provider-alicloud/issues/1739))
- improve(slb): updated slb slb_attachment testcase ([#1738](https://github.com/aliyun/terraform-provider-alicloud/issues/1738))
- use a new ram role instead of hardcode name about emr unit test case and example ([#1732](https://github.com/aliyun/terraform-provider-alicloud/issues/1732))
- Revision of goReportCard.com suggestions ([#1729](https://github.com/aliyun/terraform-provider-alicloud/issues/1729))
- improve(cs): resources supports timeouts setting ([#1679](https://github.com/aliyun/terraform-provider-alicloud/issues/1679))

BUG FIXES:

- fix(instance):fix instance test bug ([#1768](https://github.com/aliyun/terraform-provider-alicloud/issues/1768))
- fix slb sweep bug and add region for role test ([#1752](https://github.com/aliyun/terraform-provider-alicloud/issues/1752))
- fix (cs) : log_config support create new project ([#1745](https://github.com/aliyun/terraform-provider-alicloud/issues/1745))
- fix(cs): modify the new_nat_gateway field in testcase to avoid InstanceRouterEntryNotExist error ([#1733](https://github.com/aliyun/terraform-provider-alicloud/issues/1733))
- fix(mongodb):fix password encrypt bug ([#1730](https://github.com/aliyun/terraform-provider-alicloud/issues/1730))
- fix typo in worker_instance_types description ([#1726](https://github.com/aliyun/terraform-provider-alicloud/issues/1726))

## 1.57.1 (October 11, 2019)

IMPROVEMENTS:

- improve:improve some resource support encrypt password ([#1727](https://github.com/aliyun/terraform-provider-alicloud/issues/1727))
- improve(sdk): updated sdk to v1.60.191 ([#1725](https://github.com/aliyun/terraform-provider-alicloud/issues/1725))
- update tablestore package ([#1719](https://github.com/aliyun/terraform-provider-alicloud/issues/1719))
- managekubernetes support sls ([#1718](https://github.com/aliyun/terraform-provider-alicloud/issues/1718))
- improve(instance): support encrypt password when creating or updating ecs instance ([#1711](https://github.com/aliyun/terraform-provider-alicloud/issues/1711))
- update golang image version ([#1709](https://github.com/aliyun/terraform-provider-alicloud/issues/1709))
- Added credit_specification to ECS instance resource ([#1705](https://github.com/aliyun/terraform-provider-alicloud/issues/1705))
- improve(slb_rule): remove `name` forcenew and make it can be updated ([#1703](https://github.com/aliyun/terraform-provider-alicloud/issues/1703))
- upgrade terraform package ([#1702](https://github.com/aliyun/terraform-provider-alicloud/issues/1702))
- improve emr test case, update document ([#1698](https://github.com/aliyun/terraform-provider-alicloud/issues/1698))
- improve emr test case ([#1697](https://github.com/aliyun/terraform-provider-alicloud/issues/1697))
- improve(provider): update go version to 1.12 ([#1686](https://github.com/aliyun/terraform-provider-alicloud/issues/1686))
- impove(slb_listener) slb listener support same port ([#1655](https://github.com/aliyun/terraform-provider-alicloud/issues/1655))

BUG FIXES:

- fix go clean error in the ci ([#1710](https://github.com/aliyun/terraform-provider-alicloud/issues/1710))

## 1.57.0 (September 27, 2019)

- **New Resource:** `alicloud_ddosbgp_instance` ([#1650](https://github.com/aliyun/terraform-provider-alicloud/issues/1650))
- **New Data Source:** `alicloud_ddosbgp_instances` ([#1650](https://github.com/aliyun/terraform-provider-alicloud/issues/1650))
- **New Resource:** `alicloud_emr_cluster` ([#1644](https://github.com/aliyun/terraform-provider-alicloud/issues/1644))
- **New Resource:** `alicloud_vpn_route_entry` ([#1613](https://github.com/aliyun/terraform-provider-alicloud/issues/1613))

IMPROVEMENTS:

- improve(ci): add new job emr ([#1695](https://github.com/aliyun/terraform-provider-alicloud/issues/1695))
- improve(elasticsearch): added retry setting to avoid InstanceStatusNotSupportCurrentAction and InstanceActivating error ([#1693](https://github.com/aliyun/terraform-provider-alicloud/issues/1693))
- improve useragent setting ([#1692](https://github.com/aliyun/terraform-provider-alicloud/issues/1692))
- improve(ecs):add resource_group_id to ecs ([#1690](https://github.com/aliyun/terraform-provider-alicloud/issues/1690))
- improve(sls): improve sls notfounderror ([#1689](https://github.com/aliyun/terraform-provider-alicloud/issues/1689))
- improve(kafka): added retry to aviod GetTopicList Throttling.User error ([#1688](https://github.com/aliyun/terraform-provider-alicloud/issues/1688))
- improve(ci): add ddosbgp job ([#1687](https://github.com/aliyun/terraform-provider-alicloud/issues/1687))
- improve: rds,redis,mongodb remove the enumeration ([#1684](https://github.com/aliyun/terraform-provider-alicloud/issues/1684))
- Update the default period of the ddosbgp instance to 12, add the bandwidth value 201, and update the test case ([#1683](https://github.com/aliyun/terraform-provider-alicloud/issues/1683))
- improve(elasticsearch): added wait setting for retry ([#1678](https://github.com/aliyun/terraform-provider-alicloud/issues/1678))
- improve(provider): change the ubuntu version to 18 ([#1677](https://github.com/aliyun/terraform-provider-alicloud/issues/1677))
- improve(provider): support provider test ([#1675](https://github.com/aliyun/terraform-provider-alicloud/issues/1675))
- ddoscoo instance only support upgrade currently ([#1673](https://github.com/aliyun/terraform-provider-alicloud/issues/1673))

BUG FIXES:

- fix unsupport account site for test ([#1696](https://github.com/aliyun/terraform-provider-alicloud/issues/1696))
- fix(ram user): supported backward compatible ([#1685](https://github.com/aliyun/terraform-provider-alicloud/issues/1685))

## 1.56.0 (September 20, 2019)

- **New Resource:** `alicloud_alikafka_consumer_group` ([#1658](https://github.com/aliyun/terraform-provider-alicloud/issues/1658))
- **New Data Source:** `alicloud_alikafka_consumer_groups` ([#1658](https://github.com/aliyun/terraform-provider-alicloud/issues/1658))
- **New Resource:** `alicloud_alikafka_topic` ([#1642](https://github.com/aliyun/terraform-provider-alicloud/issues/1642))
- **New Data Source:** `alicloud_alikafka_topics` ([#1642](https://github.com/aliyun/terraform-provider-alicloud/issues/1642))

IMPROVEMENTS:

- improve(elasticsearch): Added retry to avoid UpdateInstance ConcurrencyUpdateInstanceConflict error. ([#1669](https://github.com/aliyun/terraform-provider-alicloud/issues/1669))
- fix(security_group_rule):fix description bug  ([#1668](https://github.com/aliyun/terraform-provider-alicloud/issues/1668))
- improve: rds,redis,mongodb support modify maintain time ([#1665](https://github.com/aliyun/terraform-provider-alicloud/issues/1665))
- add missing field ALICLOUD_INSTANCE_ID ([#1664](https://github.com/aliyun/terraform-provider-alicloud/issues/1664))
- improve(sdk): update sdk to v1.60.164 ([#1663](https://github.com/aliyun/terraform-provider-alicloud/issues/1663))
- improve(ci): add ci test for alikafka ([#1662](https://github.com/aliyun/terraform-provider-alicloud/issues/1662))
- improve(provider): rename source_name to configuration_source ([#1661](https://github.com/aliyun/terraform-provider-alicloud/issues/1661))
- improve(cen): Added wait time to avoid CreateCen Operation.Blocking error ([#1660](https://github.com/aliyun/terraform-provider-alicloud/issues/1660))
- improve(provider): add a new field source_name to mark template ([#1657](https://github.com/aliyun/terraform-provider-alicloud/issues/1657))
- improve(vpc): Added retry to avoid ListTagResources Throttling error ([#1652](https://github.com/aliyun/terraform-provider-alicloud/issues/1652))
- update VPNgateway resource vswitchId field ([#1643](https://github.com/aliyun/terraform-provider-alicloud/issues/1643))

BUG FIXES:

- fix(ess_alarm):The 'ForceNew' attribute of input parameter 'scaling_group_id' is set 'True'. ([#1671](https://github.com/aliyun/terraform-provider-alicloud/issues/1671))
- fix(testCommon):fix test common bug ([#1666](https://github.com/aliyun/terraform-provider-alicloud/issues/1666))

## 1.55.4 (September 17, 2019)

IMPROVEMENTS:

- improve(table store): set primary key to forcenew ([#1654](https://github.com/aliyun/terraform-provider-alicloud/issues/1654))
- improve(docs): Added sensitive tag for the doc which has password ([#1653](https://github.com/aliyun/terraform-provider-alicloud/issues/1653))
- improve(provider): add the provider verison in the useragent ([#1651](https://github.com/aliyun/terraform-provider-alicloud/issues/1651))
- improve(images): modified the testcase of images datasource ([#1648](https://github.com/aliyun/terraform-provider-alicloud/issues/1648))
- improve(security_group_id):update description to support for modify ([#1647](https://github.com/aliyun/terraform-provider-alicloud/issues/1647))
- impove(slb):add new allowed spec for slb ([#1646](https://github.com/aliyun/terraform-provider-alicloud/issues/1646))
- improve(provider):support ecs_role_name + assume_role ([#1639](https://github.com/aliyun/terraform-provider-alicloud/issues/1639))
- improve(example): update the examples to the format of terraform version 0.12 ([#1633](https://github.com/aliyun/terraform-provider-alicloud/issues/1633))
- improve(instance):remove bandwidth limit ([#1630](https://github.com/aliyun/terraform-provider-alicloud/issues/1630))
- improve(gpdb): gpdb instance supported tags ([#1615](https://github.com/aliyun/terraform-provider-alicloud/issues/1615))

BUG FIXES:

- fix(security_group):fix security_group bug ([#1640](https://github.com/aliyun/terraform-provider-alicloud/issues/1640))
- fix(rds): add diffsuppressfunc to rds tags ([#1602](https://github.com/aliyun/terraform-provider-alicloud/issues/1602))

## 1.55.3 (September 09, 2019)

IMPROVEMENTS:

- improve(slb): midified the sweep rules of slb ([#1631](https://github.com/aliyun/terraform-provider-alicloud/issues/1631))
- improve(slb): add new field resource_group_id ([#1629](https://github.com/aliyun/terraform-provider-alicloud/issues/1629))
- improve(example): update the examples to the format of the new version ([#1625](https://github.com/aliyun/terraform-provider-alicloud/issues/1625))
- improve(api gateway): api gateway app supported tags ([#1622](https://github.com/aliyun/terraform-provider-alicloud/issues/1622))
- improve(vpc): vpc resources and datasources supported tags ([#1621](https://github.com/aliyun/terraform-provider-alicloud/issues/1621))
- improve(kvstore): kvstore instance supported tags ([#1619](https://github.com/aliyun/terraform-provider-alicloud/issues/1619))
- update example to support for snat's creation with multi eips ([#1554](https://github.com/aliyun/terraform-provider-alicloud/issues/1554))

BUG FIXES:

- fix(common_bandwidth_package):make ratio ForceNew ([#1626](https://github.com/aliyun/terraform-provider-alicloud/issues/1626))
- fix(disk):fix disk detach bug ([#1610](https://github.com/aliyun/terraform-provider-alicloud/issues/1610))
- fix:resource security_group 'inner_access_policy' replaces 'inner_access',resource slb 'address_type' replaces 'internet' ([#1594](https://github.com/aliyun/terraform-provider-alicloud/issues/1594))

## 1.55.2 (August 30, 2019)

IMPROVEMENTS:

- imporve(elasticsearch): modified availability zone of elasticsearch instance. ([#1617](https://github.com/aliyun/terraform-provider-alicloud/issues/1617))
- improve(ram & actiontrail): added precheck for resources testcases. ([#1616](https://github.com/aliyun/terraform-provider-alicloud/issues/1616))
- improve(cdn): cdn domain supported tags. ([#1609](https://github.com/aliyun/terraform-provider-alicloud/issues/1609))
- improve(db_readonly_instance):improve db_readonly_instance testcase ([#1607](https://github.com/aliyun/terraform-provider-alicloud/issues/1607))
- imporve(cdn) modified wait time of cdn domain creation. ([#1606](https://github.com/aliyun/terraform-provider-alicloud/issues/1606))
- improve(drds): modified drds supported regions ([#1605](https://github.com/aliyun/terraform-provider-alicloud/issues/1605))
- improve(CI): change sweeper time ([#1600](https://github.com/aliyun/terraform-provider-alicloud/issues/1600))
- improve(rds): fix db_instance apply error after import ([#1599](https://github.com/aliyun/terraform-provider-alicloud/issues/1599))
- improve(ons_topic):retry when Throttling.User ([#1598](https://github.com/aliyun/terraform-provider-alicloud/issues/1598))
- Improve(ddoscoo): Improve its resource and datasource use common method ([#1591](https://github.com/aliyun/terraform-provider-alicloud/issues/1591))
- Improve(slb):slb support set AddressIpVersion ([#1587](https://github.com/aliyun/terraform-provider-alicloud/issues/1587))
- Improve(cs_kubernetes): Improve its resource and datasource use common method ([#1584](https://github.com/aliyun/terraform-provider-alicloud/issues/1584))
- Improve(cs_managed_kubernetes): Improve its resource and datasource use common method ([#1581](https://github.com/aliyun/terraform-provider-alicloud/issues/1581))

BUG FIXES:

- fix(ons):fix ons error Throttling.User ([#1608](https://github.com/aliyun/terraform-provider-alicloud/issues/1608))
- fix(ons): fix the create group error in testcase ([#1604](https://github.com/aliyun/terraform-provider-alicloud/issues/1604))

## 1.55.1 (August 23, 2019)

IMPROVEMENTS:

- improve(ons_instance): set instance name using random ([#1597](https://github.com/aliyun/terraform-provider-alicloud/issues/1597))
- add support to Ipsec_pfs field be set with "disabled" and add example files ([#1589](https://github.com/aliyun/terraform-provider-alicloud/issues/1589))
- improve(slb): sweep the protected slb ([#1588](https://github.com/aliyun/terraform-provider-alicloud/issues/1588))
- Improve(ram): ram resources supports import ([#1586](https://github.com/aliyun/terraform-provider-alicloud/issues/1586))
- improve(tags): modified test case to check the upper case letters in tags ([#1585](https://github.com/aliyun/terraform-provider-alicloud/issues/1585))
- improve(Document):improve document demo about set ([#1580](https://github.com/aliyun/terraform-provider-alicloud/issues/1580))
- Update RouteEntry Resource RouteEntryName Field ([#1578](https://github.com/aliyun/terraform-provider-alicloud/issues/1578))
- improve(ci):supplement log ([#1577](https://github.com/aliyun/terraform-provider-alicloud/issues/1577))
- improve(sdk):update alibaba-cloud-sdk-go(1.60.107) ([#1575](https://github.com/aliyun/terraform-provider-alicloud/issues/1575))
- Rename resource name that is not start with a letter ([#1573](https://github.com/aliyun/terraform-provider-alicloud/issues/1573))
- Improve(datahub_topic): Improve resource use common method ([#1565](https://github.com/aliyun/terraform-provider-alicloud/issues/1565))
- Improve(datahub_subscription): Improve resource use common method ([#1556](https://github.com/aliyun/terraform-provider-alicloud/issues/1556))
- Improve(datahub_project): Improve resource use common method ([#1555](https://github.com/aliyun/terraform-provider-alicloud/issues/1555))

BUG FIXES:

- fix(vsw): fix bug from github issue ([#1593](https://github.com/aliyun/terraform-provider-alicloud/issues/1593))
- fix(instance):update instance testcase ([#1590](https://github.com/aliyun/terraform-provider-alicloud/issues/1590))
- fix(ci):fix CI statistics bug ([#1576](https://github.com/aliyun/terraform-provider-alicloud/issues/1576))
- Fix typo ([#1574](https://github.com/aliyun/terraform-provider-alicloud/issues/1574))
- fix(disks):fix dataSource test case bug ([#1566](https://github.com/aliyun/terraform-provider-alicloud/issues/1566))

## 1.55.0 (August 16, 2019)

- **New Resource:** `alicloud_ess_notification` ([#1549](https://github.com/aliyun/terraform-provider-alicloud/issues/1549))

IMPROVEMENTS:

- improve(key_pair):update key_pair document ([#1563](https://github.com/aliyun/terraform-provider-alicloud/issues/1563))
- improve(CI): add default bucket and region for CI ([#1561](https://github.com/aliyun/terraform-provider-alicloud/issues/1561))
- improve(CI): terraform CI log ([#1557](https://github.com/aliyun/terraform-provider-alicloud/issues/1557))
- Improve(ots_instance_attachment): Improve its resource and datasource use common method ([#1552](https://github.com/aliyun/terraform-provider-alicloud/issues/1552))
- Improve(ots_instance): Improve its resource and datasource use common method ([#1551](https://github.com/aliyun/terraform-provider-alicloud/issues/1551))
- Improve(ram): ram policy attachment resources supports import ([#1550](https://github.com/aliyun/terraform-provider-alicloud/issues/1550))
- Improve(ots_table): Improve its resource and datasource use common method ([#1546](https://github.com/aliyun/terraform-provider-alicloud/issues/1546))
- Improve(router_interface): modified testcase multi count ([#1545](https://github.com/aliyun/terraform-provider-alicloud/issues/1545))
- Improve(images): removed image alinux check in datasource ([#1543](https://github.com/aliyun/terraform-provider-alicloud/issues/1543))
- Improve(logtail_config): Improve resource use common method ([#1500](https://github.com/aliyun/terraform-provider-alicloud/issues/1500))

BUG FIXES:

- bugfix：throw notFoundError when scalingGroup is not found ([#1572](https://github.com/aliyun/terraform-provider-alicloud/issues/1572))
- fix(sweep): modified the error return to run sweep completely ([#1569](https://github.com/aliyun/terraform-provider-alicloud/issues/1569))
- fix(CI): remove the useless code ([#1564](https://github.com/aliyun/terraform-provider-alicloud/issues/1564))
- fix(CI): fix pipeline grammar error ([#1562](https://github.com/aliyun/terraform-provider-alicloud/issues/1562))
- Fix log document ([#1559](https://github.com/aliyun/terraform-provider-alicloud/issues/1559))
- modify(cs): skip the testcases of cs_application and cs_swarm ([#1553](https://github.com/aliyun/terraform-provider-alicloud/issues/1553))
- fix kvstore unexpected state 'Changing' ([#1539](https://github.com/aliyun/terraform-provider-alicloud/issues/1539))

## 1.54.0 (August 12, 2019)

- **New Data Source:** `alicloud_slb_master_slave_server_groups` ([#1531](https://github.com/aliyun/terraform-provider-alicloud/issues/1531))
- **New Resource:** `alicloud_slb_master_slave_server_group` ([#1531](https://github.com/aliyun/terraform-provider-alicloud/issues/1531))
- **New Data Source:** `alicloud_instance_type_families` ([#1519](https://github.com/aliyun/terraform-provider-alicloud/issues/1519))

IMPROVEMENTS:

- improve(provider):profile,role_arn,session_name,session_expiration support ENV ([#1537](https://github.com/aliyun/terraform-provider-alicloud/issues/1537))
- support sg description ([#1536](https://github.com/aliyun/terraform-provider-alicloud/issues/1536))
- support mac address ([#1535](https://github.com/aliyun/terraform-provider-alicloud/issues/1535))
- improve(sdk): update sdk and modify api_gateway strconvs ([#1533](https://github.com/aliyun/terraform-provider-alicloud/issues/1533))
- Improve(pvtz_zone_record): Improve resource use common method ([#1528](https://github.com/aliyun/terraform-provider-alicloud/issues/1528))
- improve(alicloud_ess_scaling_group): support 'COST_OPTIMIZED' mode of autoscaling group ([#1527](https://github.com/aliyun/terraform-provider-alicloud/issues/1527))
- Improve(pvtz_zone): Improve its and attachment resources use common method ([#1525](https://github.com/aliyun/terraform-provider-alicloud/issues/1525))
- remove useless trigger in vpn ci ([#1522](https://github.com/aliyun/terraform-provider-alicloud/issues/1522))
- Improve(cr_repo): Improve resource use common method ([#1515](https://github.com/aliyun/terraform-provider-alicloud/issues/1515))
- Improve(cr_namespace): Improve resource use common method ([#1509](https://github.com/aliyun/terraform-provider-alicloud/issues/1509))
- improve(kvstore): kvstore_instance resource supports timeouts setting ([#1445](https://github.com/aliyun/terraform-provider-alicloud/issues/1445))

BUG FIXES:

- Fix(alicloud_logstore_index) Repair parameter description document ([#1532](https://github.com/aliyun/terraform-provider-alicloud/issues/1532))
- fix(sweep): modified the region of prefixes ([#1526](https://github.com/aliyun/terraform-provider-alicloud/issues/1526))
- fix(mongodb_instance): fix notfound error when describing it ([#1521](https://github.com/aliyun/terraform-provider-alicloud/issues/1521))

## 1.53.0 (August 02, 2019)

- **New Resource:** `alicloud_ons_group` ([#1506](https://github.com/aliyun/terraform-provider-alicloud/issues/1506))
- **New Resource:** `alicloud_ess_scalinggroup_vserver_groups` ([#1503](https://github.com/aliyun/terraform-provider-alicloud/issues/1503))
- **New Resource:** `alicloud_slb_backend_server` ([#1498](https://github.com/aliyun/terraform-provider-alicloud/issues/1498))
- **New Resource:** `alicloud_ons_topic` ([#1483](https://github.com/aliyun/terraform-provider-alicloud/issues/1483))
- **New Data Source:** `alicloud_ons_groups` ([#1506](https://github.com/aliyun/terraform-provider-alicloud/issues/1506))
- **New Data source:** `alicloud_slb_backend_servers` ([#1498](https://github.com/aliyun/terraform-provider-alicloud/issues/1498))
- **New Data Source:** `alicloud_ons_topics` ([#1483](https://github.com/aliyun/terraform-provider-alicloud/issues/1483))


IMPROVEMENTS:

- improve(dns_record): add diffsuppressfunc to avoid DomainRecordDuplicate error. ([#1518](https://github.com/aliyun/terraform-provider-alicloud/issues/1518))
- remove useless import ([#1517](https://github.com/aliyun/terraform-provider-alicloud/issues/1517))
- remove empty fields in managed k8s, add force_update, add multiple az support ([#1516](https://github.com/aliyun/terraform-provider-alicloud/issues/1516))
- improve(fc_function):fc_function support sweeper ([#1513](https://github.com/aliyun/terraform-provider-alicloud/issues/1513))
- improve(fc_trigger):fc_trigger support sweeper ([#1512](https://github.com/aliyun/terraform-provider-alicloud/issues/1512))
- Improve(logtail_attachment): Improve resource use common method [[#1508](https://github.com/aliyun/terraform-provider-alicloud/issues/1508)] 
- improve(slb):update testcase ([#1507](https://github.com/aliyun/terraform-provider-alicloud/issues/1507))
- improve(disk):update disk_attachment ([#1501](https://github.com/aliyun/terraform-provider-alicloud/issues/1501))
- add(slb_backend_server): slb backend server resource & data source ([#1498](https://github.com/aliyun/terraform-provider-alicloud/issues/1498))
- Improve(log_machine_group): Improve resources use common method ([#1497](https://github.com/aliyun/terraform-provider-alicloud/issues/1497))
- Improve(log_project): Improve resource use common method ([#1496](https://github.com/aliyun/terraform-provider-alicloud/issues/1496))
- improve(network_interface): enhance sweeper test ([#1495](https://github.com/aliyun/terraform-provider-alicloud/issues/1495))
- Improve(log_store): Improve resources use common method ([#1494](https://github.com/aliyun/terraform-provider-alicloud/issues/1494))
- improve(instance_type):update testcase config ([#1493](https://github.com/aliyun/terraform-provider-alicloud/issues/1493))
- Improve(mns_topic_subscription): Improve its resource use common method ([#1492](https://github.com/aliyun/terraform-provider-alicloud/issues/1492))
- improve(disk):suppurt delete_auto_snapshot delete_with_instance enable_auto_snapshot ([#1491](https://github.com/aliyun/terraform-provider-alicloud/issues/1491))
- Improve(mns_topic): Improve its resource use common method ([#1488](https://github.com/aliyun/terraform-provider-alicloud/issues/1488))
- Improve(api_gateway): api_gateway_api added testcases ([#1487](https://github.com/aliyun/terraform-provider-alicloud/issues/1487))
- Improve(mns_queue): Improve its resource use common method ([#1485](https://github.com/aliyun/terraform-provider-alicloud/issues/1485))
- improve(customer_gateway):create add retry ([#1477](https://github.com/aliyun/terraform-provider-alicloud/issues/1477))
- improve(gpdb): resources supports timeouts setting ([#1476](https://github.com/aliyun/terraform-provider-alicloud/issues/1476))
- improve(fc_triggers): Added ids filter to datasource ([#1475](https://github.com/aliyun/terraform-provider-alicloud/issues/1475))
- improve(fc_services): Added ids filter to datasource ([#1474](https://github.com/aliyun/terraform-provider-alicloud/issues/1474))
- improve(fc_functions): Added ids filter to datasource ([#1473](https://github.com/aliyun/terraform-provider-alicloud/issues/1473))
- improve(instance_types):update instance_types filter condition ([#1472](https://github.com/aliyun/terraform-provider-alicloud/issues/1472))
- improve(pvtz_zone__domain): Added ids filter to datasource ([#1471](https://github.com/aliyun/terraform-provider-alicloud/issues/1471))
- improve(cr_repos): Added names to datasource attributes ([#1470](https://github.com/aliyun/terraform-provider-alicloud/issues/1470))
- improve(cr_namespaces): Added names to datasource attributes ([#1469](https://github.com/aliyun/terraform-provider-alicloud/issues/1469))
- improve(cdn): Added region to domain name and modified sweep rules ([#1466](https://github.com/aliyun/terraform-provider-alicloud/issues/1466))
- improve(ram_roles): Added ids filter to datasource ([#1461](https://github.com/aliyun/terraform-provider-alicloud/issues/1461))
- improve(ram_users): Added ids filter to datasource ([#1459](https://github.com/aliyun/terraform-provider-alicloud/issues/1459))
- improve(pvtz_zones): Added ids filter and added names to datasource attributes ([#1458](https://github.com/aliyun/terraform-provider-alicloud/issues/1458))
- improve(nas_mount_targets): Added ids filter to datasource ([#1453](https://github.com/aliyun/terraform-provider-alicloud/issues/1453))
- improve(nas_file_systems): Added descriptions to datasource attributes ([#1450](https://github.com/aliyun/terraform-provider-alicloud/issues/1450))
- improve(nas_access_rules): Added ids filter to datasource ([#1448](https://github.com/aliyun/terraform-provider-alicloud/issues/1448))
- improve(mongodb_instance): supports timeouts setting ([#1446](https://github.com/aliyun/terraform-provider-alicloud/issues/1446))
- improve(nas_access_groups): Added names to its attributes ([#1444](https://github.com/aliyun/terraform-provider-alicloud/issues/1444))
- improve(mns_topics): Added names to datasource attributes ([#1442](https://github.com/aliyun/terraform-provider-alicloud/issues/1442))
- improve(mns_topic_subscriptions): Added names to datasource attributes ([#1441](https://github.com/aliyun/terraform-provider-alicloud/issues/1441))
- improve(mns_queues): Added names to datasource attributes ([#1439](https://github.com/aliyun/terraform-provider-alicloud/issues/1439))

BUG FIXES:

- Fix(logstore_index): Invalid update parameter change ([#1505](https://github.com/aliyun/terraform-provider-alicloud/issues/1505))
- fix(api_gateway): fix can't get resource id when stage_names set ([#1486](https://github.com/aliyun/terraform-provider-alicloud/issues/1486))
- fix(kvstore_instance): resource kvstore_instance add Retry while ModifyInstanceSpec err ([#1484](https://github.com/aliyun/terraform-provider-alicloud/issues/1484))
- fix(cen): modified the timeouts of cen instance to avoid errors ([#1451](https://github.com/aliyun/terraform-provider-alicloud/issues/1451))

## 1.52.2 (July 20, 2019)

IMPROVEMENTS:

- improve(eip_association): supporting to set PrivateIPAddress  documentation ([#1480](https://github.com/aliyun/terraform-provider-alicloud/issues/1480))
- improve(mongodb_instances): Added ids filter to datasource ([#1478](https://github.com/aliyun/terraform-provider-alicloud/issues/1478))
- improve(dns_domain): Added ids filter to datasource ([#1468](https://github.com/aliyun/terraform-provider-alicloud/issues/1468))
- improve(cdn): Added retry to avoid ServiceBusy error ([#1467](https://github.com/aliyun/terraform-provider-alicloud/issues/1467))
- improve(dns_records): Added ids filter to datasource ([#1464](https://github.com/aliyun/terraform-provider-alicloud/issues/1464))
- improve(dns_groups): Added ids filter and added names to datasource attributes ([#1463](https://github.com/aliyun/terraform-provider-alicloud/issues/1463))
- improve(stateConfig):update stateConfig error ([#1462](https://github.com/aliyun/terraform-provider-alicloud/issues/1462))
- improve(kvstore): Added ids filter to datasource ([#1457](https://github.com/aliyun/terraform-provider-alicloud/issues/1457))
- improve(cas): Added precheck to testcases ([#1456](https://github.com/aliyun/terraform-provider-alicloud/issues/1456))
- improve(rds): db_instance and db_readonly_instance resource modify timeouts 20mins to 30mins ([#1455](https://github.com/aliyun/terraform-provider-alicloud/issues/1455))
- add CI for the alicloud provider ([#1449](https://github.com/aliyun/terraform-provider-alicloud/issues/1449))
- improve(api_gateway_apps): Deprecated api_id ([#1426](https://github.com/aliyun/terraform-provider-alicloud/issues/1426))
- improve(api_gateway_apis): Added ids filter to datasource ([#1425](https://github.com/aliyun/terraform-provider-alicloud/issues/1425))
- improve(slb_server_group): remove the maximum limitation of adding backend servers ([#1416](https://github.com/aliyun/terraform-provider-alicloud/issues/1416))
- improve(cdn): cdn_domain_config added testcases ([#1405](https://github.com/aliyun/terraform-provider-alicloud/issues/1405))

BUG FIXES:

- fix(kvstore_instance): resource kvstore_instance add Retry while ModifyInstanceSpec err ([#1465](https://github.com/aliyun/terraform-provider-alicloud/issues/1465))
- fix(slb): fix slb testcase can not find instance types' bug ([#1454](https://github.com/aliyun/terraform-provider-alicloud/issues/1454))

## 1.52.1 (July 16, 2019)

IMPROVEMENTS:

- improve(disk): support online resize ([#1447](https://github.com/aliyun/terraform-provider-alicloud/issues/1447))
- improve(rds): db_readonly_instance resource supports timeouts setting ([#1438](https://github.com/aliyun/terraform-provider-alicloud/issues/1438))
- improve(rds):improve db_readonly_instance TestAccAlicloudDBReadonlyInstance_multi testcase ([#1432](https://github.com/aliyun/terraform-provider-alicloud/issues/1432))
- improve(key_pairs): Added ids filter to datasource ([#1431](https://github.com/aliyun/terraform-provider-alicloud/issues/1431))
- improve(elasticsearch): Added ids filter and added descriptions to datasource attributes ([#1430](https://github.com/aliyun/terraform-provider-alicloud/issues/1430))
- improve(drds): Added descriptions to attributes of datasource ([#1429](https://github.com/aliyun/terraform-provider-alicloud/issues/1429))
- improve(rds):update ppas not support regions ([#1428](https://github.com/aliyun/terraform-provider-alicloud/issues/1428))
- improve(api_gateway_groups): Added ids filter to datasource ([#1427](https://github.com/aliyun/terraform-provider-alicloud/issues/1427))
- improve(docs): Reformat abnormal inline HCL code in docs ([#1423](https://github.com/aliyun/terraform-provider-alicloud/issues/1423))
- improve(mns):modified mns_queues.html ([#1422](https://github.com/aliyun/terraform-provider-alicloud/issues/1422))
- improve(rds): db_instance resource supports timeouts setting ([#1409](https://github.com/aliyun/terraform-provider-alicloud/issues/1409))
- improve(kms): modified the args of kms_keys datasource ([#1407](https://github.com/aliyun/terraform-provider-alicloud/issues/1407))
- improve(kms_key): modify the param `description` to forcenew ([#1406](https://github.com/aliyun/terraform-provider-alicloud/issues/1406))

BUG FIXES:

- fix(db_instance): modified the target state of state config ([#1437](https://github.com/aliyun/terraform-provider-alicloud/issues/1437))
- fix(db_readonly_instance): fix invalid status error when updating and deleting ([#1435](https://github.com/aliyun/terraform-provider-alicloud/issues/1435))
- fix(ots_table): fix setting deviation_cell_version_in_sec error ([#1434](https://github.com/aliyun/terraform-provider-alicloud/issues/1434))
- fix(db_backup_policy): resource db_backup_policy testcase use datasource db_instance_classes ([#1424](https://github.com/aliyun/terraform-provider-alicloud/issues/1424))

## 1.52.0 (July 12, 2019)

- **New Datasource:** `alicloud_ons_instances` ([#1411](https://github.com/aliyun/terraform-provider-alicloud/issues/1411))

IMPROVEMENTS:

- improve(vpc):add ids filter ([#1420](https://github.com/aliyun/terraform-provider-alicloud/issues/1420))
- improve(db_instances): Added ids filter and added names to datasource attributes ([#1419](https://github.com/aliyun/terraform-provider-alicloud/issues/1419))
- improve(cas): Added ids filter and added names to datasource attributes ([#1417](https://github.com/aliyun/terraform-provider-alicloud/issues/1417))
- docs(format): Convert inline HCL configs to canonical format ([#1415](https://github.com/aliyun/terraform-provider-alicloud/issues/1415))
- improve(gpdb_instance):add vpc name ([#1413](https://github.com/aliyun/terraform-provider-alicloud/issues/1413))
- improve(provider): add a new parameter `skip_region_validation` in the provider config ([#1404](https://github.com/aliyun/terraform-provider-alicloud/issues/1404))
- improve(cdn): cdn_domain support certificate config ([#1393](https://github.com/aliyun/terraform-provider-alicloud/issues/1393))
- improve(rds): resource db_instance support update for instance_charge_type ([#1389](https://github.com/aliyun/terraform-provider-alicloud/issues/1389))

BUG FIXES:

- fix(db_instance):fix db_instance testcase vsw availability_zone ([#1418](https://github.com/aliyun/terraform-provider-alicloud/issues/1418))
- fix(api_gateway): modified the testcase to avoid errors ([#1410](https://github.com/aliyun/terraform-provider-alicloud/issues/1410))
- fix(db_readonly_instance): extend the waiting time for spec modification ([#1408](https://github.com/aliyun/terraform-provider-alicloud/issues/1408))
- fix(db_readonly_instance): add retryable error content in instance spec modification and deletion ([#1403](https://github.com/aliyun/terraform-provider-alicloud/issues/1403))

## 1.51.0 (July 08, 2019)

- **New Date Source:** `alicloud_kvstore_instance_engines` ([#1371](https://github.com/aliyun/terraform-provider-alicloud/issues/1371))
- **New Resource:** `alicloud_ons_instance` ([#1333](https://github.com/aliyun/terraform-provider-alicloud/issues/1333))

IMPROVEMENTS:

- improve(db_instance): improve db_instance MAZ testcase ([#1391](https://github.com/aliyun/terraform-provider-alicloud/issues/1391))
- improve(cs_kubernetes): add importIgnore parameters in the importer testcase ([#1387](https://github.com/aliyun/terraform-provider-alicloud/issues/1387))
- Remove govendor commands in CI ([#1386](https://github.com/aliyun/terraform-provider-alicloud/issues/1386))
- improve(slb_vserver_group): support attaching eni ([#1384](https://github.com/aliyun/terraform-provider-alicloud/issues/1384))
- improve(db_instance_classes): add new parameter db_instance_class ([#1383](https://github.com/aliyun/terraform-provider-alicloud/issues/1383))
- improve(images): Add os_name_en to the attributes of images datasource ([#1380](https://github.com/aliyun/terraform-provider-alicloud/issues/1380))
- improve(disk): the snapshot_id conflicts with encrypted ([#1378](https://github.com/aliyun/terraform-provider-alicloud/issues/1378))
- Improve(cs_kubernetes): add some importState ignore fields in the importer testcase ([#1377](https://github.com/aliyun/terraform-provider-alicloud/issues/1377))
- Improve(oss_bucket): Add names for its attributes of datasource ([#1374](https://github.com/aliyun/terraform-provider-alicloud/issues/1374))
- improve(common_test):update common_test for terraform 0.12 ([#1372](https://github.com/aliyun/terraform-provider-alicloud/issues/1372))
- Improve(cs_kubernetes): add import ignore parameter `log_config` ([#1370](https://github.com/aliyun/terraform-provider-alicloud/issues/1370))
- improve(slb):support slb instance delete protection ([#1369](https://github.com/aliyun/terraform-provider-alicloud/issues/1369))
- improve(slb_rule): support health check config ([#1367](https://github.com/aliyun/terraform-provider-alicloud/issues/1367))
- Improve(oss_bucket_object): Improve its use common method ([#1366](https://github.com/aliyun/terraform-provider-alicloud/issues/1366))
- improve(drds_instance): Added precheck to its testcases ([#1364](https://github.com/aliyun/terraform-provider-alicloud/issues/1364))
- Improve(oss_bucket): Improve its resource use common method ([#1353](https://github.com/aliyun/terraform-provider-alicloud/issues/1353))
- improve(launch_template): support update method ([#1327](https://github.com/aliyun/terraform-provider-alicloud/issues/1327))
- improve(snapshot): support setting timeouts ([#1304](https://github.com/aliyun/terraform-provider-alicloud/issues/1304))
- improve(instance):update testcase ([#1199](https://github.com/aliyun/terraform-provider-alicloud/issues/1199))

BUG FIXES:

- fix(instnace): fix missing dry_run when creating instance ([#1401](https://github.com/aliyun/terraform-provider-alicloud/issues/1401))
- fix(oss_bucket): fix oss bucket deleting timeout error ([#1400](https://github.com/aliyun/terraform-provider-alicloud/issues/1400))
- fix(route_entry):fix route_entry create bug ([#1398](https://github.com/aliyun/terraform-provider-alicloud/issues/1398))
- fix(instance):fix testcase name too length bug ([#1396](https://github.com/aliyun/terraform-provider-alicloud/issues/1396))
- fix(vswitch):fix vswitch describe method wrapErrorf bug ([#1392](https://github.com/aliyun/terraform-provider-alicloud/issues/1392))
- fix(slb_rule): fix testcase bug ([#1390](https://github.com/aliyun/terraform-provider-alicloud/issues/1390))
- fix(db_backup_policy): pg10 of category 'basic' modify log_backup error ([#1388](https://github.com/aliyun/terraform-provider-alicloud/issues/1388))
- fix(cen):Add deadline to cen datasources and modify timeout for DescribeCenBandwidthPackages ([#1381](https://github.com/aliyun/terraform-provider-alicloud/issues/1381))
- fix(kvstore): kvstore_instance PostPaid to PrePaid error ([#1375](https://github.com/aliyun/terraform-provider-alicloud/issues/1375))
- fix(cen): fixed its not display error message, added CenThrottlingUser retry ([#1373](https://github.com/aliyun/terraform-provider-alicloud/issues/1373))

## 1.50.0 (July 01, 2019)

IMPROVEMENTS:

- Remove cs kubernetes autovpc testcases ([#1368](https://github.com/aliyun/terraform-provider-alicloud/issues/1368))
- disable nav-visible in the alicloud.erb file ([#1365](https://github.com/aliyun/terraform-provider-alicloud/issues/1365))
- Improve sweeper test and remove some needless waiting ([#1361](https://github.com/aliyun/terraform-provider-alicloud/issues/1361))
- This is a Terraform 0.12 compatible release of this provider ([#1356](https://github.com/aliyun/terraform-provider-alicloud/issues/1356))
- Deprecated resource `alicloud_cms_alarm` parameter start_time, end_time and removed notify_type based on the latest go sdk ([#1356](https://github.com/aliyun/terraform-provider-alicloud/issues/1356))
- Adapt to new parameters of dedicated kubernetes cluster ([#1354](https://github.com/aliyun/terraform-provider-alicloud/issues/1354))

BUG FIXES:

- Fix alicloud_cas_certificate setId bug ([#1368](https://github.com/aliyun/terraform-provider-alicloud/issues/1368))
- Fix oss bucket datasource testcase based on the 0.12 syntax ([#1362](https://github.com/aliyun/terraform-provider-alicloud/issues/1362))
- Fix deleting mongodb instance "NotFound" bug ([#1359](https://github.com/aliyun/terraform-provider-alicloud/issues/1359))

## 1.49.0 (June 28, 2019)

- **New Date Source:** `alicloud_kvstore_instance_classes` ([#1315](https://github.com/aliyun/terraform-provider-alicloud/issues/1315))

IMPROVEMENTS:

- remove the skipped testcase ([#1349](https://github.com/aliyun/terraform-provider-alicloud/issues/1349))
- Move some import testcase into resource testcase ([#1348](https://github.com/aliyun/terraform-provider-alicloud/issues/1348))
- Support attach & detach operation for loadbalancers and dbinstances ([#1346](https://github.com/aliyun/terraform-provider-alicloud/issues/1346))
- update security_group_rule md ([#1345](https://github.com/aliyun/terraform-provider-alicloud/issues/1345))
- Improve mongodb,rds testcase ([#1339](https://github.com/aliyun/terraform-provider-alicloud/issues/1339))
- Deprecate field statement and use field document to replace ([#1338](https://github.com/aliyun/terraform-provider-alicloud/issues/1338))
- Add function BuildStateConf for common timeouts setting ([#1330](https://github.com/aliyun/terraform-provider-alicloud/issues/1330))
- drds_instance resource supports timeouts setting ([#1329](https://github.com/aliyun/terraform-provider-alicloud/issues/1329))
- add support get Ak from config file ([#1328](https://github.com/aliyun/terraform-provider-alicloud/issues/1328))
- Improve api_gateway_vpc use common method. ([#1323](https://github.com/aliyun/terraform-provider-alicloud/issues/1323))
- Organize official documents in alphabetical order ([#1322](https://github.com/aliyun/terraform-provider-alicloud/issues/1322))
- improve snapshot_policy testcase ([#1313](https://github.com/aliyun/terraform-provider-alicloud/issues/1313))
- Improve api_gateway_group use common method ([#1311](https://github.com/aliyun/terraform-provider-alicloud/issues/1311))
- Improve api_gateway_app use common method ([#1306](https://github.com/aliyun/terraform-provider-alicloud/issues/1306))

BUG FIXES:

- bugfix: modify ess loadbalancers batch size ([#1352](https://github.com/aliyun/terraform-provider-alicloud/issues/1352))
- fix instance OperationConflict bug ([#1351](https://github.com/aliyun/terraform-provider-alicloud/issues/1351))
- fix(nas): convert some retrable error to nonretryable ([#1344](https://github.com/aliyun/terraform-provider-alicloud/issues/1344))
- fix mongodb testcase ([#1341](https://github.com/aliyun/terraform-provider-alicloud/issues/1341))
- fix log_store fields cannot be changed ([#1337](https://github.com/aliyun/terraform-provider-alicloud/issues/1337))
- fix(nas): fix error handling ([#1336](https://github.com/aliyun/terraform-provider-alicloud/issues/1336))
- fix db_instance_classes,db_instance_engines ([#1331](https://github.com/aliyun/terraform-provider-alicloud/issues/1331))
- fix sls-logconfig config_name field to name ([#1326](https://github.com/aliyun/terraform-provider-alicloud/issues/1326))
- fix db_instance_engines testcase ([#1325](https://github.com/aliyun/terraform-provider-alicloud/issues/1325))
- fix forward_entries testcase bug ([#1324](https://github.com/aliyun/terraform-provider-alicloud/issues/1324))

## 1.48.0 (June 21, 2019)

- **New Resource:** `alicloud_gpdb_connection` ([#1290](https://github.com/aliyun/terraform-provider-alicloud/issues/1290))

IMPROVEMENTS:

- Improve rds testcase zone_id ([#1321](https://github.com/aliyun/terraform-provider-alicloud/issues/1321))
- feature: support enable/disable action for resource alicloud_ess_alarm ([#1320](https://github.com/aliyun/terraform-provider-alicloud/issues/1320))
- cen_instance resource supports timeouts setting ([#1318](https://github.com/aliyun/terraform-provider-alicloud/issues/1318))
- added importer support for security_group_rule ([#1317](https://github.com/aliyun/terraform-provider-alicloud/issues/1317))
- add multi_zone for db_instance_classes and db_instance_engines ([#1310](https://github.com/aliyun/terraform-provider-alicloud/issues/1310))
- Update Eip Resource Isp Field ([#1303](https://github.com/aliyun/terraform-provider-alicloud/issues/1303))
- Improve db_instance,db_read_write_splitting_connection,db_readonly_instance testcase ([#1300](https://github.com/aliyun/terraform-provider-alicloud/issues/1300))
- Improve api_gateway_api use common method ([#1299](https://github.com/aliyun/terraform-provider-alicloud/issues/1299))
- Add name for cen bandwidth package testcase ([#1298](https://github.com/aliyun/terraform-provider-alicloud/issues/1298))
- Improve db testcase ([#1294](https://github.com/aliyun/terraform-provider-alicloud/issues/1294))
- elasticsearch_instance resource supports timeouts setting ([#1268](https://github.com/aliyun/terraform-provider-alicloud/issues/1268))

BUG FIXES:

- bugfix: remove the 'ForceNew' attribute of 'vswitch_ids' from resource alicloud_ess_scaling_group ([#1316](https://github.com/aliyun/terraform-provider-alicloud/issues/1316))
- managed k8s no longer returns vswitchids and instancetypes, fix crash ([#1314](https://github.com/aliyun/terraform-provider-alicloud/issues/1314))
- fix db_instance_classes ([#1309](https://github.com/aliyun/terraform-provider-alicloud/issues/1309))
- fix oss lifecycle nil pointer bug ([#1307](https://github.com/aliyun/terraform-provider-alicloud/issues/1307))
- Fix cen_bandwidth_limit Throttling.User bug ([#1305](https://github.com/aliyun/terraform-provider-alicloud/issues/1305))
- fix disk_attachment test bug ([#1302](https://github.com/aliyun/terraform-provider-alicloud/issues/1302))

## 1.47.0 (June 17, 2019)

- **New Date Source:** `alicloud_gpdb_instances` ([#1279](https://github.com/aliyun/terraform-provider-alicloud/issues/1279))
- **New Resource:** `alicloud_gpdb_instance` ([#1260](https://github.com/aliyun/terraform-provider-alicloud/issues/1260))

IMPROVEMENTS:

- fc_trigger datasource support outputting ids and names ([#1286](https://github.com/aliyun/terraform-provider-alicloud/issues/1286))
- add fc_trigger support cdn_events ([#1285](https://github.com/aliyun/terraform-provider-alicloud/issues/1285))
- modify apigateway-fc example ([#1284](https://github.com/aliyun/terraform-provider-alicloud/issues/1284))
- Added PGP encrypt Support for ram access key ([#1280](https://github.com/aliyun/terraform-provider-alicloud/issues/1280))
- Update Eip Resource Isp Field ([#1275](https://github.com/aliyun/terraform-provider-alicloud/issues/1275))
- Improve fc_service use common method ([#1269](https://github.com/aliyun/terraform-provider-alicloud/issues/1269))
- Improve fc_function use common method ([#1266](https://github.com/aliyun/terraform-provider-alicloud/issues/1266))
- update dns_group testcase name ([#1265](https://github.com/aliyun/terraform-provider-alicloud/issues/1265))
- update slb sdk ([#1263](https://github.com/aliyun/terraform-provider-alicloud/issues/1263))
- improve vpn_connection testcase ([#1257](https://github.com/aliyun/terraform-provider-alicloud/issues/1257))
- Improve cen_route_entries use common method ([#1249](https://github.com/aliyun/terraform-provider-alicloud/issues/1249))
- Improve cen_bandwidth_package_attachment resource use common method ([#1240](https://github.com/aliyun/terraform-provider-alicloud/issues/1240))
- Improve cen_bandwidth_package resource use common method ([#1237](https://github.com/aliyun/terraform-provider-alicloud/issues/1237))

BUG FIXES:

- feat(nas): fix error report ([#1293](https://github.com/aliyun/terraform-provider-alicloud/issues/1293))
- temp fix no value returned by cs openapi ([#1289](https://github.com/aliyun/terraform-provider-alicloud/issues/1289))
- fix disk device_name bug ([#1288](https://github.com/aliyun/terraform-provider-alicloud/issues/1288))
- fix sql server instance storage set bug ([#1283](https://github.com/aliyun/terraform-provider-alicloud/issues/1283))
- fix db_instance_classes storage_range bug ([#1282](https://github.com/aliyun/terraform-provider-alicloud/issues/1282))
- fc_service datasource support outputting ids and names ([#1278](https://github.com/aliyun/terraform-provider-alicloud/issues/1278))
- fix log_store ListShards InternalServerError bug ([#1277](https://github.com/aliyun/terraform-provider-alicloud/issues/1277))
- fix slb_listener docs bug ([#1276](https://github.com/aliyun/terraform-provider-alicloud/issues/1276))
- fix clientToken bug ([#1272](https://github.com/aliyun/terraform-provider-alicloud/issues/1272))
- fix(nas): fix document and nas_access_rules ([#1271](https://github.com/aliyun/terraform-provider-alicloud/issues/1271))
- docs(version) Added 6.7 supported and fixed bug of version difference ([#1270](https://github.com/aliyun/terraform-provider-alicloud/issues/1270))
- fix(nas): fix documents ([#1267](https://github.com/aliyun/terraform-provider-alicloud/issues/1267))
- fix(nas): describe mount target & access rule ([#1264](https://github.com/aliyun/terraform-provider-alicloud/issues/1264))

## 1.46.0 (June 10, 2019)

- **New Resource:** `alicloud_ram_account_password_policy` ([#1212](https://github.com/aliyun/terraform-provider-alicloud/issues/1212))
- **New Date Source:** `alicloud_db_instance_engines` ([#1201](https://github.com/aliyun/terraform-provider-alicloud/issues/1201))
- **New Date Source:** `alicloud_db_instance_classes` ([#1201](https://github.com/aliyun/terraform-provider-alicloud/issues/1201))

IMPROVEMENTS:

- refactor(nas): move import to resource ([#1254](https://github.com/aliyun/terraform-provider-alicloud/issues/1254))
- Improve ess_scalingconfiguration use common method ([#1250](https://github.com/aliyun/terraform-provider-alicloud/issues/1250))
- improve ssl_vpn_client_cert testcase ([#1248](https://github.com/aliyun/terraform-provider-alicloud/issues/1248))
- Improve ram_account_password_policy resource use common method ([#1247](https://github.com/aliyun/terraform-provider-alicloud/issues/1247))
- add pending status for resource instance when creating ([#1245](https://github.com/aliyun/terraform-provider-alicloud/issues/1245))
- resource instance supports timeouts configure ([#1244](https://github.com/aliyun/terraform-provider-alicloud/issues/1244))
- added webhook support for alarms ([#1243](https://github.com/aliyun/terraform-provider-alicloud/issues/1243))
- improve common test method ([#1242](https://github.com/aliyun/terraform-provider-alicloud/issues/1242))
- Update Eip Association Resource ([#1238](https://github.com/aliyun/terraform-provider-alicloud/issues/1238))
- improve ssl_vpn_server testcase ([#1235](https://github.com/aliyun/terraform-provider-alicloud/issues/1235))
- Improve ess_scalingconfigurations datasource use common method ([#1234](https://github.com/aliyun/terraform-provider-alicloud/issues/1234))
- improve vpn_customer_gateway testcase ([#1232](https://github.com/aliyun/terraform-provider-alicloud/issues/1232))
- Improve cen_instance_grant use common method ([#1230](https://github.com/aliyun/terraform-provider-alicloud/issues/1230))
- improve vpn_gateway testcase ([#1229](https://github.com/aliyun/terraform-provider-alicloud/issues/1229))
- Improve cen_bandwidth_limit use common method ([#1227](https://github.com/aliyun/terraform-provider-alicloud/issues/1227))
- Feature/support multi instance types ([#1226](https://github.com/aliyun/terraform-provider-alicloud/issues/1226))
- Improve ess_attachment use common method ([#1225](https://github.com/aliyun/terraform-provider-alicloud/issues/1225))
- Improve ess_alarm use common method ([#1218](https://github.com/aliyun/terraform-provider-alicloud/issues/1218))
- Add support for assume_role in provider block ([#1217](https://github.com/aliyun/terraform-provider-alicloud/issues/1217))
- Improve cen_instance_attachment resource use common method. ([#1216](https://github.com/aliyun/terraform-provider-alicloud/issues/1216))
- add db instance engines and db instance classes data source support ([#1201](https://github.com/aliyun/terraform-provider-alicloud/issues/1201))
- Handle alicloud_cs_*_kubernetes resource NotFound error properly ([#1191](https://github.com/aliyun/terraform-provider-alicloud/issues/1191))

BUG FIXES:

- fix slb_attachment classic testcase ([#1259](https://github.com/aliyun/terraform-provider-alicloud/issues/1259))
- fix oss bucket update bug ([#1258](https://github.com/aliyun/terraform-provider-alicloud/issues/1258))
- fix scalingConfiguration is inconsistent with the information that is returned by describe, when the input parameter user_data is base64 ([#1256](https://github.com/aliyun/terraform-provider-alicloud/issues/1256))
- fix slb_attachment err ObtainIpFail ([#1253](https://github.com/aliyun/terraform-provider-alicloud/issues/1253))
- Fix password to comliant with the default password policy ([#1241](https://github.com/aliyun/terraform-provider-alicloud/issues/1241))
- fix cr repo details, improve cs and cr docs ([#1239](https://github.com/aliyun/terraform-provider-alicloud/issues/1239))
- fix(nas): fix unittest bugs ([#1236](https://github.com/aliyun/terraform-provider-alicloud/issues/1236))
- fix slb_ca_certificate err ServiceIsConfiguring ([#1233](https://github.com/aliyun/terraform-provider-alicloud/issues/1233))
- fix reset account_password don't work ([#1231](https://github.com/aliyun/terraform-provider-alicloud/issues/1231))
- fix(nas): fix testcase errors ([#1184](https://github.com/aliyun/terraform-provider-alicloud/issues/1184))

## 1.45.0 (May 29, 2019)

FEATURES:

- **New Resource:** `alicloud_network_acl_entries` ([#1208](https://github.com/aliyun/terraform-provider-alicloud/issues/1208))

IMPROVEMENTS:

- update changeLog ([#1224](https://github.com/aliyun/terraform-provider-alicloud/issues/1224))
- support oss object versioning ([#1121](https://github.com/aliyun/terraform-provider-alicloud/issues/1121))
- update instance dataSource doc ([#1215](https://github.com/aliyun/terraform-provider-alicloud/issues/1215))
- update oss buket encryption configuration ([#1214](https://github.com/aliyun/terraform-provider-alicloud/issues/1214))
- support oss bucket tags ([#1213](https://github.com/aliyun/terraform-provider-alicloud/issues/1213))
- support oss bucket encryption configuration ([#1210](https://github.com/aliyun/terraform-provider-alicloud/issues/1210))
- Improve cen_instances use common method ([#1206](https://github.com/aliyun/terraform-provider-alicloud/issues/1206))
- support set oss bucket stroage class ([#1204](https://github.com/aliyun/terraform-provider-alicloud/issues/1204))
- Improve ess_lifecyclehook resource use common method ([#1196](https://github.com/aliyun/terraform-provider-alicloud/issues/1196))
- Improve ess_scalinggroup use common method ([#1192](https://github.com/aliyun/terraform-provider-alicloud/issues/1192))
- Improve ess_scheduled_task resource use common method ([#1175](https://github.com/aliyun/terraform-provider-alicloud/issues/1175))
- improve route_table testcase ([#1109](https://github.com/aliyun/terraform-provider-alicloud/issues/1109))

BUG FIXES:

- fix nat_gateway and network_interface testcase bug ([#1211](https://github.com/aliyun/terraform-provider-alicloud/issues/1211))
- Fix ram testcases name length bug ([#1205](https://github.com/aliyun/terraform-provider-alicloud/issues/1205))
- fix actiontrail bug ([#1203](https://github.com/aliyun/terraform-provider-alicloud/issues/1203))

## 1.44.0 (May 24, 2019)

FEATURES:

- **New Resource:** `alicloud_network_acl_attachment` ([#1187](https://github.com/aliyun/terraform-provider-alicloud/issues/1187))

IMPROVEMENTS:

- update CHANGELOG.md ([#1209](https://github.com/aliyun/terraform-provider-alicloud/issues/1209))
- Skip instance some testcases to avoid qouta limit ([#1195](https://github.com/aliyun/terraform-provider-alicloud/issues/1195))
- Added the multi zone's instance supported ([#1194](https://github.com/aliyun/terraform-provider-alicloud/issues/1194))
- remove multi test of ram_account_alias ([#1186](https://github.com/aliyun/terraform-provider-alicloud/issues/1186))
- Improve ram_role_attachment resource use common method ([#1185](https://github.com/aliyun/terraform-provider-alicloud/issues/1185))
- Improve ess_scalingrule use common method ([#1183](https://github.com/aliyun/terraform-provider-alicloud/issues/1183))
- update mongodb instance resource document ([#1182](https://github.com/aliyun/terraform-provider-alicloud/issues/1182))
- Improve ram_role resource use common method ([#1181](https://github.com/aliyun/terraform-provider-alicloud/issues/1181))
- Correct the oss bucket docs ([#1178](https://github.com/aliyun/terraform-provider-alicloud/issues/1178))
- add slb classic not support regions ([#1176](https://github.com/aliyun/terraform-provider-alicloud/issues/1176))
- Dev versioning ([#1174](https://github.com/aliyun/terraform-provider-alicloud/issues/1174))
- Improve ram_user_policy_attachment resource use common method ([#1172](https://github.com/aliyun/terraform-provider-alicloud/issues/1172))
- Improve ram_role_policy_attachment resource use common method ([#1171](https://github.com/aliyun/terraform-provider-alicloud/issues/1171))
- improve router_interface testcase ([#1170](https://github.com/aliyun/terraform-provider-alicloud/issues/1170))
- Improve ram_policy resource use common method ([#1166](https://github.com/aliyun/terraform-provider-alicloud/issues/1166))
- Improve slb_listeners datasource use common method ([#1165](https://github.com/aliyun/terraform-provider-alicloud/issues/1165))
- add name attribute for forward_entry ([#1164](https://github.com/aliyun/terraform-provider-alicloud/issues/1164))
- Improve ram_group_policy_attachment resource use common method ([#1163](https://github.com/aliyun/terraform-provider-alicloud/issues/1163))
- Improve ram_group_membership resource use common method ([#1159](https://github.com/aliyun/terraform-provider-alicloud/issues/1159))
- Improve ram_login_profile resource use common method ([#1158](https://github.com/aliyun/terraform-provider-alicloud/issues/1158))
- Improve ram_group resource use common method ([#1150](https://github.com/aliyun/terraform-provider-alicloud/issues/1150))

BUG FIXES:

- Fix ram_user sweeper ([#1200](https://github.com/aliyun/terraform-provider-alicloud/issues/1200))
- Fix ram group import bug ([#1198](https://github.com/aliyun/terraform-provider-alicloud/issues/1198))
- fix router_interface dataSource testcase bug ([#1197](https://github.com/aliyun/terraform-provider-alicloud/issues/1197))
- fix forward_entry multi testcase bug ([#1189](https://github.com/aliyun/terraform-provider-alicloud/issues/1189))
- fix api gw and network acl sweeper test error ([#1180](https://github.com/aliyun/terraform-provider-alicloud/issues/1180))
- fix ram user diff bug ([#1179](https://github.com/aliyun/terraform-provider-alicloud/issues/1179))
- Fix ram account alias multi testcase bug ([#1169](https://github.com/aliyun/terraform-provider-alicloud/issues/1169))

## 1.43.0 (May 17, 2019)

FEATURES:

- **New Resource:** `alicloud_network_acl` (([#1151](https://github.com/aliyun/terraform-provider-alicloud/issues/1151))

IMPROVEMENTS:

- change ecs instance instance_charge_type modifying position ([#1168](https://github.com/aliyun/terraform-provider-alicloud/issues/1168))
- AutoScaling support multiple security groups ([#1167](https://github.com/aliyun/terraform-provider-alicloud/issues/1167))
- Update ots and vpc document ([#1162](https://github.com/aliyun/terraform-provider-alicloud/issues/1162))
- Improve some slb datasource ([#1155](https://github.com/aliyun/terraform-provider-alicloud/issues/1155))
- improve forward_entry testcase ([#1152](https://github.com/aliyun/terraform-provider-alicloud/issues/1152))
- improve slb_attachment resource use common method ([#1148](https://github.com/aliyun/terraform-provider-alicloud/issues/1148))
- Improve ram_account_alias resource use common method  ([#1147](https://github.com/aliyun/terraform-provider-alicloud/issues/1147))
- slb instance support updating specification ([#1145](https://github.com/aliyun/terraform-provider-alicloud/issues/1145))
- improve slb_server_group resource use common method ([#1144](https://github.com/aliyun/terraform-provider-alicloud/issues/1144))
- add note for SLB that intl account does not support creating PrePaid instance ([#1143](https://github.com/aliyun/terraform-provider-alicloud/issues/1143))
- Update ots document ([#1142](https://github.com/aliyun/terraform-provider-alicloud/issues/1142))
- improve slb_server_certificate resource use common method ([#1139](https://github.com/aliyun/terraform-provider-alicloud/issues/1139))

BUG FIXES:

- Fix ram account alias notfound bug ([#1161](https://github.com/aliyun/terraform-provider-alicloud/issues/1161))
- fix(nas): refactor testcases ([#1157](https://github.com/aliyun/terraform-provider-alicloud/issues/1157))

## 1.42.0 (May 10, 2019)

FEATURES:

- **New Resource:** `alicloud_snapshot_policy` ([#989](https://github.com/aliyun/terraform-provider-alicloud/issues/989))

IMPROVEMENTS:

- improve mongodb and db sweeper test ([#1138](https://github.com/aliyun/terraform-provider-alicloud/issues/1138))
- Alicloud_ots_table: add max version offset ([#1137](https://github.com/aliyun/terraform-provider-alicloud/issues/1137))
- update disk category ([#1135](https://github.com/aliyun/terraform-provider-alicloud/issues/1135))
- Update Route Entry Resource ([#1134](https://github.com/aliyun/terraform-provider-alicloud/issues/1134))
- update images testcase check condition ([#1133](https://github.com/aliyun/terraform-provider-alicloud/issues/1133))
- bugfix: ess alarm apply recreate ([#1131](https://github.com/aliyun/terraform-provider-alicloud/issues/1131))
- improve slb_listener resource use common method ([#1130](https://github.com/aliyun/terraform-provider-alicloud/issues/1130))
- mongodb sharding instance add backup policy support ([#1127](https://github.com/aliyun/terraform-provider-alicloud/issues/1127))
- Improve ram_users datasource use common method ([#1126](https://github.com/aliyun/terraform-provider-alicloud/issues/1126))
- Improve ram_policies datasource use common method ([#1125](https://github.com/aliyun/terraform-provider-alicloud/issues/1125))
- rds datasource test case remove connection mode check ([#1124](https://github.com/aliyun/terraform-provider-alicloud/issues/1124))
- Add missing bracket ([#1123](https://github.com/aliyun/terraform-provider-alicloud/issues/1123))
- add support sha256 ([#1122](https://github.com/aliyun/terraform-provider-alicloud/issues/1122))
- Improve ram_groups datasource use common method ([#1121](https://github.com/aliyun/terraform-provider-alicloud/issues/1121))
- Modified the sweep rules in ram_roles testcases ([#1116](https://github.com/aliyun/terraform-provider-alicloud/issues/1116))
- improve instance testcase ([#1114](https://github.com/aliyun/terraform-provider-alicloud/issues/1114))
- Improve slb_ca_certificate resource use common method ([#1113](https://github.com/aliyun/terraform-provider-alicloud/issues/1113))
- Improve ram_roles datasource use common method ([#1112](https://github.com/aliyun/terraform-provider-alicloud/issues/1112))
- Improve slb datasource use common method ([#1111](https://github.com/aliyun/terraform-provider-alicloud/issues/1111))
- Improve ram_account_alias use common method ([#1108](https://github.com/aliyun/terraform-provider-alicloud/issues/1108))
- update data_source_alicoud_mongo_instances and add test case ([#1107](https://github.com/aliyun/terraform-provider-alicloud/issues/1107))
- add mongodb backup policy support, test case, document ([#1106](https://github.com/aliyun/terraform-provider-alicloud/issues/1106))
- update route_entry and forward_entry document ([#1096](https://github.com/aliyun/terraform-provider-alicloud/issues/1096))
- Improve slb_acl resource use common method ([#1092](https://github.com/aliyun/terraform-provider-alicloud/issues/1092))
- improve snat_entry testcase ([#1091](https://github.com/aliyun/terraform-provider-alicloud/issues/1091))
- Improve slb resource use common method ([#1090](https://github.com/aliyun/terraform-provider-alicloud/issues/1090))
- improve nat_gateway testcase ([#1089](https://github.com/aliyun/terraform-provider-alicloud/issues/1089))
- Modify table to entry ([#1088](https://github.com/aliyun/terraform-provider-alicloud/issues/1088))
- Modified the error code returned when timeout of upgrading instance ([#1085](https://github.com/aliyun/terraform-provider-alicloud/issues/1085))
- improve db backup policy test case ([#1083](https://github.com/aliyun/terraform-provider-alicloud/issues/1083))

BUG FIXES:

- Fix scalinggroup id is not found before creating scaling configuration  ([#1119](https://github.com/aliyun/terraform-provider-alicloud/issues/1119))
- fix slb instance sets tags bug ([#1105](https://github.com/aliyun/terraform-provider-alicloud/issues/1105))
- fix not support outputfile ([#1095](https://github.com/aliyun/terraform-provider-alicloud/issues/1095))
- Bugfix/slb import server group ([#1093](https://github.com/aliyun/terraform-provider-alicloud/issues/1093))
- Fix fc_triggers datasource when type is mns_topic ([#1086](https://github.com/aliyun/terraform-provider-alicloud/issues/1086))

## 1.41.0 (April 29, 2019)

IMPROVEMENTS:

- Improve fc_trigger support mns_topic modify config ([#1082](https://github.com/aliyun/terraform-provider-alicloud/issues/1082))
- Rds sdk-update ([#1078](https://github.com/aliyun/terraform-provider-alicloud/issues/1078))
- update some eip method name ([#1077](https://github.com/aliyun/terraform-provider-alicloud/issues/1077))
- improve vswitch testcase  ([#1076](https://github.com/aliyun/terraform-provider-alicloud/issues/1076))
- add rand for db_instances testcase ([#1074](https://github.com/aliyun/terraform-provider-alicloud/issues/1074))
- Improve fc_trigger support mns_topic ([#1073](https://github.com/aliyun/terraform-provider-alicloud/issues/1073))
- remove zone_id setting in the db instance testcase ([#1069](https://github.com/aliyun/terraform-provider-alicloud/issues/1069))
- change database default zone id to avoid some unsupported cases ([#1067](https://github.com/aliyun/terraform-provider-alicloud/issues/1067))
- add oss bucket policy implementation ([#1066](https://github.com/aliyun/terraform-provider-alicloud/issues/1066))
- improve vpc testcase ([#1065](https://github.com/aliyun/terraform-provider-alicloud/issues/1065))
- Change password to Yourpassword ([#1063](https://github.com/aliyun/terraform-provider-alicloud/issues/1063))
- Improve kvstore_instance datasource use common method ([#1062](https://github.com/aliyun/terraform-provider-alicloud/issues/1062))
- improve eip testcase ([#1058](https://github.com/aliyun/terraform-provider-alicloud/issues/1058))
- Improve kvstore_instance testcase use common method ([#1052](https://github.com/aliyun/terraform-provider-alicloud/issues/1052))
- improve mongodb testcase ([#1050](https://github.com/aliyun/terraform-provider-alicloud/issues/1050))
- update network_interface dataSource basic testcase config ([#1049](https://github.com/aliyun/terraform-provider-alicloud/issues/1049))
- Improve kvstore_backup_policy testcase use common method ([#1044](https://github.com/aliyun/terraform-provider-alicloud/issues/1044))

BUG FIXES:

- Fix fc_triggers datasource when type is mns_topic ([#1086](https://github.com/aliyun/terraform-provider-alicloud/issues/1086))
- Fix kvstore_instance multi ([#1080](https://github.com/aliyun/terraform-provider-alicloud/issues/1080))
- fix eip_association bug when snat or forward be released ([#1075](https://github.com/aliyun/terraform-provider-alicloud/issues/1075))
- Fix db_readonly_instance instance_name ([#1071](https://github.com/aliyun/terraform-provider-alicloud/issues/1071))
- fixed DB log backup policy bug when the log_retention_period does not input ([#1056](https://github.com/aliyun/terraform-provider-alicloud/issues/1056))
- fix cms diff bug and improve its testcases ([#1057](https://github.com/aliyun/terraform-provider-alicloud/issues/1057))


## 1.40.0 (April 20, 2019)

FEATURES:

- **New Resource:** `alicloud_mongodb_sharding_instance` ([#1017](https://github.com/aliyun/terraform-provider-alicloud/issues/1017))
- **New Data Source:** `alicloud_snapshots` ([#988](https://github.com/aliyun/terraform-provider-alicloud/issues/988))
- **New Resource:** `alicloud_snapshot` ([#954](https://github.com/aliyun/terraform-provider-alicloud/issues/954))

IMPROVEMENTS:

- Fix db_instance can't find method DescribeDbInstance ([#1046](https://github.com/aliyun/terraform-provider-alicloud/issues/1046))
- update network_interface testcase config ([#1045](https://github.com/aliyun/terraform-provider-alicloud/issues/1045))
- Update Nat Gateway Resource ([#1043](https://github.com/aliyun/terraform-provider-alicloud/issues/1043))
- improve network_interface dataSource testcase ([#1042](https://github.com/aliyun/terraform-provider-alicloud/issues/1042))
- improve network_interface resource testcase ([#1041](https://github.com/aliyun/terraform-provider-alicloud/issues/1041))
- Improve db_database db_instance db_readonly_instance db_readwrite_splitting_connection ([#1040](https://github.com/aliyun/terraform-provider-alicloud/issues/1040))
- improve key_pair resource testcase ([#1039](https://github.com/aliyun/terraform-provider-alicloud/issues/1039))
- improve key_pair dataSource testcase ([#1038](https://github.com/aliyun/terraform-provider-alicloud/issues/1038))
- make fmt ess_scalinggroups ([#1036](https://github.com/aliyun/terraform-provider-alicloud/issues/1036))
- improve test common method ([#1030](https://github.com/aliyun/terraform-provider-alicloud/issues/1030))
- Update cen data source document ([#1029](https://github.com/aliyun/terraform-provider-alicloud/issues/1029))
- fix Error method [[#1024](https://github.com/aliyun/terraform-provider-alicloud/issues/1024)] 
- Update Nat Gateway Token ([#1020](https://github.com/aliyun/terraform-provider-alicloud/issues/1020))
- update RAM website document ([#1019](https://github.com/aliyun/terraform-provider-alicloud/issues/1019))
- add computed for resource_group_id ([#1018](https://github.com/aliyun/terraform-provider-alicloud/issues/1018))
- remove ram validators and update website docs ([#1016](https://github.com/aliyun/terraform-provider-alicloud/issues/1016))
- improve test common method, support 'TestMatchResourceAttr' check ([#1012](https://github.com/aliyun/terraform-provider-alicloud/issues/1012))
- resource group support for creating new VPC ([#1010](https://github.com/aliyun/terraform-provider-alicloud/issues/1010))
- Improve cs_cluster sweeper test removing retained resources ([#1002](https://github.com/aliyun/terraform-provider-alicloud/issues/1002))
- improve security_group testcase use common method ([#995](https://github.com/aliyun/terraform-provider-alicloud/issues/995))
- fix vpn change local_subnet and remote_subnet bug ([#994](https://github.com/aliyun/terraform-provider-alicloud/issues/994))
- improve disk dataSource testcase use common method ([#990](https://github.com/aliyun/terraform-provider-alicloud/issues/990))
- fix(nas): use new sdk ([#984](https://github.com/aliyun/terraform-provider-alicloud/issues/984))
- Feature/slb listener redirect http to https ([#981](https://github.com/aliyun/terraform-provider-alicloud/issues/981))
- improve disk and diskAttachment resource testcase use testCommon method ([#978](https://github.com/aliyun/terraform-provider-alicloud/issues/978))
- improve dns dataSource testcase use testCommon method ([#971](https://github.com/aliyun/terraform-provider-alicloud/issues/971))

BUG FIXES:

- Fix ess go sdk compatibility ([#1032](https://github.com/aliyun/terraform-provider-alicloud/issues/1032))
- Update sdk to fix timeout bug ([#1015](https://github.com/aliyun/terraform-provider-alicloud/issues/1015))
- Fix Eip And VSwitch ClientToken bug ([#1000](https://github.com/aliyun/terraform-provider-alicloud/issues/1000))
- fix db_account diff bug and add some notes for it ([#999](https://github.com/aliyun/terraform-provider-alicloud/issues/999))
- fix vpn gateway Period bug ([#993](https://github.com/aliyun/terraform-provider-alicloud/issues/993))


## 1.39.0 (April 09, 2019)

FEATURES:

- **New Data Source:** `alicloud_ots_instance_attachments` ([#986](https://github.com/aliyun/terraform-provider-alicloud/issues/986))
- **New Data Source:** `alicloud_ssl_vpc_servers` ([#985](https://github.com/aliyun/terraform-provider-alicloud/issues/985))
- **New Data Source:** `alicloud_ssl_vpn_client_certs` ([#986](https://github.com/aliyun/terraform-provider-alicloud/issues/986))
- **New Data Source:** `alicloud_ess_scaling_rules` ([#976](https://github.com/aliyun/terraform-provider-alicloud/issues/976))
- **New Data Source:** `alicloud_ess_scaling_configurations` ([#974](https://github.com/aliyun/terraform-provider-alicloud/issues/974))
- **New Data Source:** `alicloud_ess_scaling_groups` ([#973](https://github.com/aliyun/terraform-provider-alicloud/issues/973))
- **New Data Source:** `alicloud_ddoscoo_instances` ([#967](https://github.com/aliyun/terraform-provider-alicloud/issues/967))
- **New Data Source:** `alicloud_ots_instances` ([#946](https://github.com/aliyun/terraform-provider-alicloud/issues/946))

IMPROVEMENTS:

- Improve instance type updating testcase ([#979](https://github.com/aliyun/terraform-provider-alicloud/issues/979))
- support changing prepaid instance type ([#977](https://github.com/aliyun/terraform-provider-alicloud/issues/977))
- Improve db_account db_account_privilege db_backup_policy db_connection ([#963](https://github.com/aliyun/terraform-provider-alicloud/issues/963))

BUG FIXES:

- Fix Nat GW ClientToken bug ([#983](https://github.com/aliyun/terraform-provider-alicloud/issues/983))
- Fix print error bug after DescribeDBInstanceById ([#980](https://github.com/aliyun/terraform-provider-alicloud/issues/980))

## 1.38.0 (April 03, 2019)

FEATURES:

- **New Resource:** `alicloud_ddoscoo_instance` ([#952](https://github.com/aliyun/terraform-provider-alicloud/issues/952))

IMPROVEMENTS:

- update dns_group describe method ([#966](https://github.com/aliyun/terraform-provider-alicloud/issues/966))
- update ram_policy resource testcase ([#964](https://github.com/aliyun/terraform-provider-alicloud/issues/964))
- improve ram_policy resource update method ([#960](https://github.com/aliyun/terraform-provider-alicloud/issues/960))
- ecs prepaid instance supports changing instance type ([#949](https://github.com/aliyun/terraform-provider-alicloud/issues/949))
- update mongodb instance test case for multiAZ ([#947](https://github.com/aliyun/terraform-provider-alicloud/issues/947))
- add test common method ,improve dns resource testcase ([#927](https://github.com/aliyun/terraform-provider-alicloud/issues/927))


BUG FIXES:

- Fix drds instance sweeper test bug ([#955](https://github.com/aliyun/terraform-provider-alicloud/issues/955))

## 1.37.0 (March 29, 2019)

FEATURES:

- **New Resource:** `alicloud_mongodb_instance` ([#881](https://github.com/aliyun/terraform-provider-alicloud/issues/881))
- **New Resource:** `alicloud_cen_instance_grant` ([#857](https://github.com/aliyun/terraform-provider-alicloud/issues/857))
- **New Data Source:** `alicloud_forward_entries` ([#922](https://github.com/aliyun/terraform-provider-alicloud/issues/922))
- **New Data Source:** `alicloud_snat_entries` ([#920](https://github.com/aliyun/terraform-provider-alicloud/issues/920))
- **New Data Source:** `alicloud_nat_gateways` ([#918](https://github.com/aliyun/terraform-provider-alicloud/issues/918))
- **New Data Source:** `alicloud_route_entries` ([#915](https://github.com/aliyun/terraform-provider-alicloud/issues/915))

IMPROVEMENTS:

- Add missing outputs for datasource dns_records, security groups, vpcs and vswitches ([#943](https://github.com/aliyun/terraform-provider-alicloud/issues/943))
- datasource dns_records add a output urls ([#942](https://github.com/aliyun/terraform-provider-alicloud/issues/942))
- modify stop instance timeout to 5min to avoid the exception timeout ([#941](https://github.com/aliyun/terraform-provider-alicloud/issues/941))
- datasource security_groups, vpcs and vswitches support outputs ids and names ([#939](https://github.com/aliyun/terraform-provider-alicloud/issues/939))
- Improve all of parameter's tag, like 'Required', 'ForceNew' ([#938](https://github.com/aliyun/terraform-provider-alicloud/issues/938))
- Improve pvtz_zone_record WrapError ([#934](https://github.com/aliyun/terraform-provider-alicloud/issues/934))
- Improve pvtz_zone_record create record ([#933](https://github.com/aliyun/terraform-provider-alicloud/issues/933))
- testSweepCRNamespace skip not supported region  ([#932](https://github.com/aliyun/terraform-provider-alicloud/issues/932))
- refine retry logic of resource tablestore to avoid the exception timeout ([#931](https://github.com/aliyun/terraform-provider-alicloud/issues/931))
- Improve pvtz resource datasource testcases ([#928](https://github.com/aliyun/terraform-provider-alicloud/issues/928))
- cr_repos fix docs link error ([#926](https://github.com/aliyun/terraform-provider-alicloud/issues/926))
- resource DB instance supports setting security group ([#925](https://github.com/aliyun/terraform-provider-alicloud/issues/925))
- resource DB instance supports setting monitor period ([#924](https://github.com/aliyun/terraform-provider-alicloud/issues/924))
- Skipping bandwidth package related test for international site account ([#917](https://github.com/aliyun/terraform-provider-alicloud/issues/917))
- Resource snat entry update id and support import ([#916](https://github.com/aliyun/terraform-provider-alicloud/issues/916))
- add docs about prerequisites for cs and cr  ([#914](https://github.com/aliyun/terraform-provider-alicloud/issues/914))
- add new schema environment_variables to fc_function.html.markdown ([#913](https://github.com/aliyun/terraform-provider-alicloud/issues/913))
- add skipping check for datasource route tables' testcases ([#911](https://github.com/aliyun/terraform-provider-alicloud/issues/911))
- modify ram_user id by userId ([#900](https://github.com/aliyun/terraform-provider-alicloud/issues/900))

BUG FIXES:

- Deprecate bucket `logging_isenable` and fix referer_config diff bug ([#937](https://github.com/aliyun/terraform-provider-alicloud/issues/937))
- fix ram user and group sweeper test bug ([#929](https://github.com/aliyun/terraform-provider-alicloud/issues/929))
- Fix the parameter bug when actiontrail is created ([#921](https://github.com/aliyun/terraform-provider-alicloud/issues/921))
- fix default pod_cidr in k8s docs ([#919](https://github.com/aliyun/terraform-provider-alicloud/issues/919))

## 1.36.0 (March 24, 2019)

FEATURES:

- **New Resource:** `alicloud_cas_certificate` ([#875](https://github.com/aliyun/terraform-provider-alicloud/issues/875))
- **New Data Source:** `alicloud_route_tables` ([#905](https://github.com/aliyun/terraform-provider-alicloud/issues/905))
- **New Data Source:** `alicloud_common_bandwidth_packages` ([#897](https://github.com/aliyun/terraform-provider-alicloud/issues/897))
- **New Data Source:** `alicloud_actiontrails` ([#891](https://github.com/aliyun/terraform-provider-alicloud/issues/891))
- **New Data Source:** `alicloud_cas_certificates` ([#875](https://github.com/aliyun/terraform-provider-alicloud/issues/875))

IMPROVEMENTS:

- Add wait method for disk and disk attachment ([#910](https://github.com/aliyun/terraform-provider-alicloud/issues/910))
- Add wait method for cen instance ([#909](https://github.com/aliyun/terraform-provider-alicloud/issues/909))
- add dns and dns_group test sweeper ([#906](https://github.com/aliyun/terraform-provider-alicloud/issues/906))
- fc_function add new schema environment_variables ([#904](https://github.com/aliyun/terraform-provider-alicloud/issues/904))
- support kv-store auto renewal option  documentation ([#902](https://github.com/aliyun/terraform-provider-alicloud/issues/902))
- Sort slb slave zone ids to avoid needless error ([#898](https://github.com/aliyun/terraform-provider-alicloud/issues/898))
- add region skip for container registry testcase ([#896](https://github.com/aliyun/terraform-provider-alicloud/issues/896))
- Add `enable_details` for alicloud_zones and support retrieving slb slave zones ([#893](https://github.com/aliyun/terraform-provider-alicloud/issues/893))
- Slb support setting master and slave zone id ([#887](https://github.com/aliyun/terraform-provider-alicloud/issues/887))
- improve disk and attachment resource testcase ([#886](https://github.com/aliyun/terraform-provider-alicloud/issues/886))
- Remove ModifySecurityGroupPolicy waiting and backend has fixed it ([#883](https://github.com/aliyun/terraform-provider-alicloud/issues/883))
- Improve cas resource and datasource testcases ([#882](https://github.com/aliyun/terraform-provider-alicloud/issues/882))
- Make db_connection resource code more standard ([#879](https://github.com/aliyun/terraform-provider-alicloud/issues/879))

BUG FIXES:

- Fix cen instance deleting bug ([#908](https://github.com/aliyun/terraform-provider-alicloud/issues/908))
- Fix cen create bug when one resion is China ([#903](https://github.com/aliyun/terraform-provider-alicloud/issues/903))
- fix cas_certificate sweeper test bug ([#899](https://github.com/aliyun/terraform-provider-alicloud/issues/899))
- Modify ram group's name's ForceNew to true ([#895](https://github.com/aliyun/terraform-provider-alicloud/issues/895))
- fix mount target deletion bugs ([#892](https://github.com/aliyun/terraform-provider-alicloud/issues/892))
- Fix link to BatchSetCdnDomainConfig document  documentation ([#885](https://github.com/aliyun/terraform-provider-alicloud/issues/885))
- fix rds instance parameter test case issue ([#880](https://github.com/aliyun/terraform-provider-alicloud/issues/880))

## 1.35.0 (March 18, 2019)

FEATURES:

- **New Resource:** `alicloud_cr_repo` ([#862](https://github.com/aliyun/terraform-provider-alicloud/issues/862))
- **New Resource:** `alicloud_actiontrail` ([#858](https://github.com/aliyun/terraform-provider-alicloud/issues/858))
- **New Data Source:** `alicloud_cr_repos` ([#868](https://github.com/aliyun/terraform-provider-alicloud/issues/868))
- **New Data Source:** `alicloud_cr_namespaces` ([#867](https://github.com/aliyun/terraform-provider-alicloud/issues/867))
- **New Data Source:** `alicloud_nas_file_systems` ([#864](https://github.com/aliyun/terraform-provider-alicloud/issues/864))
- **New Data Source:** `alicloud_nas_mount_targets` ([#864](https://github.com/aliyun/terraform-provider-alicloud/issues/864))
- **New Data Source:** `alicloud_drds_instances` ([#861](https://github.com/aliyun/terraform-provider-alicloud/issues/861))
- **New Data Source:** `alicloud_nas_access_rules` ([#860](https://github.com/aliyun/terraform-provider-alicloud/issues/860))
- **New Data Source:** `alicloud_nas_access_groups` ([#856](https://github.com/aliyun/terraform-provider-alicloud/issues/856))

IMPROVEMENTS:

- Improve actiontrail docs ([#878](https://github.com/aliyun/terraform-provider-alicloud/issues/878))
- Add account pre-check for common bandwidth package to avoid known error ([#877](https://github.com/aliyun/terraform-provider-alicloud/issues/877))
- Make dns resource code more standard ([#876](https://github.com/aliyun/terraform-provider-alicloud/issues/876))
- Improve dns resources' testcases ([#859](https://github.com/aliyun/terraform-provider-alicloud/issues/859))
- Add client token for vpn services ([#855](https://github.com/aliyun/terraform-provider-alicloud/issues/855))
- reback the lossing datasource ([#866](https://github.com/aliyun/terraform-provider-alicloud/issues/866))
- Improve drds instances testcases  documentation ([#863](https://github.com/aliyun/terraform-provider-alicloud/issues/863))
- Update sdk for vpc package ([#854](https://github.com/aliyun/terraform-provider-alicloud/issues/854))

BUG FIXES:

- Add waiting method to ensure the security group status is ok ([#873](https://github.com/aliyun/terraform-provider-alicloud/issues/873))
- Fix nas mount target notfound bug and improve nas datasource's testcases ([#872](https://github.com/aliyun/terraform-provider-alicloud/issues/872))
- Fix dns notfound bug ([#871](https://github.com/aliyun/terraform-provider-alicloud/issues/871))
- fix creating slb bug ([#870](https://github.com/aliyun/terraform-provider-alicloud/issues/870))
- fix elastic search sweeper test bug ([#865](https://github.com/aliyun/terraform-provider-alicloud/issues/865))


## 1.34.0 (March 13, 2019)

FEATURES:

- **New Resource:** `alicloud_nas_mount_target` ([#835](https://github.com/aliyun/terraform-provider-alicloud/issues/835))
- **New Resource:** `alicloud_cdn_domain_config` ([#829](https://github.com/aliyun/terraform-provider-alicloud/issues/829))
- **New Resource:** `alicloud_cr_namespace` ([#827](https://github.com/aliyun/terraform-provider-alicloud/issues/827))
- **New Resource:** `alicloud_nas_access_rule` ([#827](https://github.com/aliyun/terraform-provider-alicloud/issues/827))
- **New Resource:** `alicloud_cdn_domain_new` ([#787](https://github.com/aliyun/terraform-provider-alicloud/issues/787))
- **New Data Source:** `alicloud_cs_kubernetes_clusters` ([#818](https://github.com/aliyun/terraform-provider-alicloud/issues/818))

IMPROVEMENTS:

- Add drds instance docs ([#853](https://github.com/aliyun/terraform-provider-alicloud/issues/853))
- Improve resource mount target testcases ([#852](https://github.com/aliyun/terraform-provider-alicloud/issues/852))
- Add using note for spot instance ([#851](https://github.com/aliyun/terraform-provider-alicloud/issues/851))
- Resource alicloud_slb supports PrePaid ([#850](https://github.com/aliyun/terraform-provider-alicloud/issues/850))
- Add ssl_vpn_server and ssl_vpn_client_cert sweeper test ([#843](https://github.com/aliyun/terraform-provider-alicloud/issues/843))
- Improve vpn_gateway testcases and some sweeper test ([#842](https://github.com/aliyun/terraform-provider-alicloud/issues/842))
- Improve dns datasource testcases ([#841](https://github.com/aliyun/terraform-provider-alicloud/issues/841))
- Improve Eip and mns testcase ([#840](https://github.com/aliyun/terraform-provider-alicloud/issues/840))
- Add version notes in some docs ([#838](https://github.com/aliyun/terraform-provider-alicloud/issues/838))
- RDS resource supports auto-renewal ([#836](https://github.com/aliyun/terraform-provider-alicloud/issues/836))
- Deprecate the resource alicloud_cdn_domain ([#830](https://github.com/aliyun/terraform-provider-alicloud/issues/830))

BUG FIXES:

- Fix deleting dns record InternalError bug ([#848](https://github.com/aliyun/terraform-provider-alicloud/issues/848))
- fix log store and config sweeper test deleting bug ([#847](https://github.com/aliyun/terraform-provider-alicloud/issues/847))
- Fix drds resource no supporting client token ([#846](https://github.com/aliyun/terraform-provider-alicloud/issues/846))
- fix kms sweeper test deleting bug ([#844](https://github.com/aliyun/terraform-provider-alicloud/issues/844))
- fix kubernetes data resource ut and import error ([#839](https://github.com/aliyun/terraform-provider-alicloud/issues/839))
- Bugfix: destroying alicloud_ess_attachment timeout ([#834](https://github.com/aliyun/terraform-provider-alicloud/issues/834))
- fix cdn service func WaitForCdnDomain ([#833](https://github.com/aliyun/terraform-provider-alicloud/issues/833))
- deal with the error message in cen route entry ([#831](https://github.com/aliyun/terraform-provider-alicloud/issues/831))
- change bool to *bool in parameters of k8s clusters ([#828](https://github.com/aliyun/terraform-provider-alicloud/issues/828))
- Fix nas docs bug ([#825](https://github.com/aliyun/terraform-provider-alicloud/issues/825))
- create vpn gateway got "UnnecessarySslConnection" error when enable_ssl is false ([#822](https://github.com/aliyun/terraform-provider-alicloud/issues/822))

## 1.33.0 (March 05, 2019)

FEATURES:

- **New Resource:** `alicloud_nas_access_group` ([#817](https://github.com/aliyun/terraform-provider-alicloud/issues/817))
- **New Resource:** `alicloud_nas_file_system` ([#807](https://github.com/aliyun/terraform-provider-alicloud/issues/807))

IMPROVEMENTS:

- Improve nas resource docs ([#824](https://github.com/aliyun/terraform-provider-alicloud/issues/824))

BUG FIXES:

- bugfix: create vpn gateway got "UnnecessarySslConnection" error when enable_ssl is false ([#822](https://github.com/aliyun/terraform-provider-alicloud/issues/822))
- fix volume_tags diff bug when running testcases ([#816](https://github.com/aliyun/terraform-provider-alicloud/issues/816))

## 1.32.1 (March 03, 2019)

BUG FIXES:

- fix volume_tags diff bug when setting tags by alicloud_disk ([#815](https://github.com/aliyun/terraform-provider-alicloud/issues/815))

## 1.32.0 (March 01, 2019)

FEATURES:

- **New Resource:** `alicloud_db_readwrite_splitting_connection` ([#753](https://github.com/aliyun/terraform-provider-alicloud/issues/753))

IMPROVEMENTS:

- add slb_internet_enabled to managed kubernetes ([#806](https://github.com/aliyun/terraform-provider-alicloud/issues/806))
- update alicloud_slb_attachment usage example ([#805](https://github.com/aliyun/terraform-provider-alicloud/issues/805))
- rds support op tags  documentation ([#797](https://github.com/aliyun/terraform-provider-alicloud/issues/797))
- ForceNew for resource record and zone id updates for pvtz record ([#794](https://github.com/aliyun/terraform-provider-alicloud/issues/794))
- support volume tags for ecs instance disks ([#793](https://github.com/aliyun/terraform-provider-alicloud/issues/793))
- Improve instance and security group testcase for different account site ([#792](https://github.com/aliyun/terraform-provider-alicloud/issues/792))
- Add account site type setting to skip unsupported test cases automatically ([#790](https://github.com/aliyun/terraform-provider-alicloud/issues/790))
- update alibaba-cloud-sdk-go to use lastest useragent and modify errMessage when signature does not match  dependencies ([#788](https://github.com/aliyun/terraform-provider-alicloud/issues/788))
- make the timeout longer when cen attach/detach vpc ([#786](https://github.com/aliyun/terraform-provider-alicloud/issues/786))
- cen child instance attach after vsw created ([#785](https://github.com/aliyun/terraform-provider-alicloud/issues/785))
- kvstore support parameter configuration ([#784](https://github.com/aliyun/terraform-provider-alicloud/issues/784))
- Modify useragent to meet the standard of sdk ([#778](https://github.com/aliyun/terraform-provider-alicloud/issues/778))
- Modify kms client to dock with the alicloud official GO SDK ([#763](https://github.com/aliyun/terraform-provider-alicloud/issues/763))

BUG FIXES:

- fix rds readonly instance name update issue ([#812](https://github.com/aliyun/terraform-provider-alicloud/issues/812))
- fix import managed kubernetes test ([#809](https://github.com/aliyun/terraform-provider-alicloud/issues/809))
- fix rds parameter update issue ([#804](https://github.com/aliyun/terraform-provider-alicloud/issues/804))
- fix first create db with tags ([#803](https://github.com/aliyun/terraform-provider-alicloud/issues/803))
- Fix dns record ttl setting error and update bug ([#800](https://github.com/aliyun/terraform-provider-alicloud/issues/800))
- Fix vpc return custom route table bug ([#799](https://github.com/aliyun/terraform-provider-alicloud/issues/799))
- fix ssl vpn subnet can not pass comma separated string problem ([#780](https://github.com/aliyun/terraform-provider-alicloud/issues/780))
- fix(whitelist) Modified whitelist returned and filter the default values ([#779](https://github.com/aliyun/terraform-provider-alicloud/issues/779))

## 1.31.0 (February 19, 2019)

FEATURES:

- **New Resource:** `alicloud_db_readonly_instance` ([#755](https://github.com/aliyun/terraform-provider-alicloud/issues/755))

IMPROVEMENTS:

- support update deletion_protection option documentation ([#771](https://github.com/aliyun/terraform-provider-alicloud/issues/771))
- add three az k8s cluster docs  documentation ([#767](https://github.com/aliyun/terraform-provider-alicloud/issues/767))
- kvstore support vpc_auth_mode  dependencies ([#765](https://github.com/aliyun/terraform-provider-alicloud/issues/765))
- Fix sls logtail config collection error ([#762](https://github.com/aliyun/terraform-provider-alicloud/issues/762))
- Add attribute parameters to resource alicloud_db_instance  documentation ([#761](https://github.com/aliyun/terraform-provider-alicloud/issues/761))
- Add attribute parameters to resource alicloud_db_instance ([#761](https://github.com/aliyun/terraform-provider-alicloud/issues/761))
- Modify dns client to dock with the alicloud official GO SDK ([#750](https://github.com/aliyun/terraform-provider-alicloud/issues/750))

BUG FIXES:

- Fix cms_alarm updating notify_type bug ([#773](https://github.com/aliyun/terraform-provider-alicloud/issues/773))
- fix(error) Fixed bug of error code when timeout for upgrade instance ([#770](https://github.com/aliyun/terraform-provider-alicloud/issues/770))
- delete success if not found cen route when delete ([#753](https://github.com/aliyun/terraform-provider-alicloud/issues/753))

## 1.30.0 (February 04, 2019)

FEATURES:

- **New Resource:** `alicloud_elasticsearch_instance` ([#722](https://github.com/aliyun/terraform-provider-alicloud/issues/722))
- **New Resource:** `alicloud_logtail_attachment` ([#705](https://github.com/aliyun/terraform-provider-alicloud/issues/705))
- **New Data Source:** `alicloud_elasticsearch_instances` ([#739](https://github.com/aliyun/terraform-provider-alicloud/issues/739))

IMPROVEMENTS:

- Improve snat and forward testcases ([#749](https://github.com/aliyun/terraform-provider-alicloud/issues/749))
- delete data source roles limit of policy_type and policy_name ([#748](https://github.com/aliyun/terraform-provider-alicloud/issues/748))
- make k8s cluster deleting timeout longer ([#746](https://github.com/aliyun/terraform-provider-alicloud/issues/746))
- Improve nat_gateway testcases ([#743](https://github.com/aliyun/terraform-provider-alicloud/issues/743))
- Improve eip_association testcases ([#742](https://github.com/aliyun/terraform-provider-alicloud/issues/742))
- Improve elasticinstnace testcases for IPV6 supported ([#741](https://github.com/aliyun/terraform-provider-alicloud/issues/741))
- Add debug for db instance and ess group ([#740](https://github.com/aliyun/terraform-provider-alicloud/issues/740))
- Improve api_gateway_vpc_access testcases ([#738](https://github.com/aliyun/terraform-provider-alicloud/issues/738))
- Modify errors and  ram client to dock with the GO SDK ([#735](https://github.com/aliyun/terraform-provider-alicloud/issues/735))
- provider supports getting credential via ecs role name ([#731](https://github.com/aliyun/terraform-provider-alicloud/issues/731))
- Update testcases for cen region domain route entries ([#729](https://github.com/aliyun/terraform-provider-alicloud/issues/729))
- cs_kubernetes supports user_ca ([#726](https://github.com/aliyun/terraform-provider-alicloud/issues/726))
- Wrap resource elasticserarch_instance's error ([#725](https://github.com/aliyun/terraform-provider-alicloud/issues/725))
- Add note for kubernetes resource and improve its testcase ([#724](https://github.com/aliyun/terraform-provider-alicloud/issues/724))
- Datasource instance_types supports filter results and used to create kuberneters ([#723](https://github.com/aliyun/terraform-provider-alicloud/issues/723))
- Add ids parameter extraction in data source regions,zones,dns_domain,images and instance_types([#718](https://github.com/aliyun/terraform-provider-alicloud/issues/718))
- Improve dns group testcase ([#717](https://github.com/aliyun/terraform-provider-alicloud/issues/717))
- Improve security group rule testcase for classic ([#716](https://github.com/aliyun/terraform-provider-alicloud/issues/716))
- Improve security group creating request ([#715](https://github.com/aliyun/terraform-provider-alicloud/issues/715))
- Route entry supports Nat Gateway ([#713](https://github.com/aliyun/terraform-provider-alicloud/issues/713))
- Modify db account returning update to read after creating ([#711](https://github.com/aliyun/terraform-provider-alicloud/issues/711))
- Improve cdn testcase ([#708](https://github.com/aliyun/terraform-provider-alicloud/issues/708))
- Apply wraperror to security_group, security_group_rule, vswitch, disk ([#707](https://github.com/aliyun/terraform-provider-alicloud/issues/707))
- Improve cdn testcase ([#705](https://github.com/aliyun/terraform-provider-alicloud/issues/705))
- Add notes for datahub and improve its testcase ([#704](https://github.com/aliyun/terraform-provider-alicloud/issues/704))
- Improve security_group_rule resource and data source testcases ([#703](https://github.com/aliyun/terraform-provider-alicloud/issues/703))
- Improve kvstore backup policy ([#701](https://github.com/aliyun/terraform-provider-alicloud/issues/701))
- Improve pvtz attachment testcase ([#700](https://github.com/aliyun/terraform-provider-alicloud/issues/700))
- Modify pagesize on API DescribeVSWitches tp avoid ServiceUnavailable ([#698](https://github.com/aliyun/terraform-provider-alicloud/issues/698))
- Improve eip resource and data source testcases ([#697](https://github.com/aliyun/terraform-provider-alicloud/issues/697))

BUG FIXES:

- FIx cen route NotFoundRoute error when deleting ([#753](https://github.com/aliyun/terraform-provider-alicloud/issues/753))
- Fix log_store InternalServerError error ([#737](https://github.com/aliyun/terraform-provider-alicloud/issues/737))
- Fix cen region route entries testcase bug ([#734](https://github.com/aliyun/terraform-provider-alicloud/issues/734))
- Fix ots_table StorageServerBusy bug ([#733](https://github.com/aliyun/terraform-provider-alicloud/issues/733))
- Fix db_account setting description bug ([#732](https://github.com/aliyun/terraform-provider-alicloud/issues/732))
- Fix Router Entry Token Bug ([#730](https://github.com/aliyun/terraform-provider-alicloud/issues/730))
- Fix instance diff bug when updating its VPC attributes ([#728](https://github.com/aliyun/terraform-provider-alicloud/issues/728))
- Fix snat entry IncorretSnatEntryStatus error when deleting ([#714](https://github.com/aliyun/terraform-provider-alicloud/issues/714))
- Fix forward entry UnknownError error ([#712](https://github.com/aliyun/terraform-provider-alicloud/issues/712))
- Fix pvtz record Zone.NotExists error when deleting record ([#710](https://github.com/aliyun/terraform-provider-alicloud/issues/710))
- Fix modify kvstore policy not working bug ([#709](https://github.com/aliyun/terraform-provider-alicloud/issues/709))
- reattach the key pair after update OS image ([#699](https://github.com/aliyun/terraform-provider-alicloud/issues/699))
- Fix ServiceUnavailable error on VPC and VSW ([#695](https://github.com/aliyun/terraform-provider-alicloud/issues/695))

## 1.29.0 (January 21, 2019)

FEATURES:

- **New Resource:** `alicloud_logtail_config` ([#685](https://github.com/aliyun/terraform-provider-alicloud/issues/685))

IMPROVEMENTS:

- Apply wraperror to ess group ([#689](https://github.com/aliyun/terraform-provider-alicloud/issues/689))
- Add wraperror and apply it to vpc and eip ([#688](https://github.com/aliyun/terraform-provider-alicloud/issues/688))
- Improve vswitch resource and data source testcases ([#687](https://github.com/aliyun/terraform-provider-alicloud/issues/687))
- Improve security_group resource and data source testcases ([#686](https://github.com/aliyun/terraform-provider-alicloud/issues/686))
- Improve vpc resource and data source testcases ([#684](https://github.com/aliyun/terraform-provider-alicloud/issues/684))
- Modify the slb sever group testcase name ([#681](https://github.com/aliyun/terraform-provider-alicloud/issues/681))
- Improve sweeper testcases ([#680](https://github.com/aliyun/terraform-provider-alicloud/issues/680))
- Improve db instance's testcases ([#679](https://github.com/aliyun/terraform-provider-alicloud/issues/679))
- Improve ecs disk's testcases ([#678](https://github.com/aliyun/terraform-provider-alicloud/issues/678))
- Add multi_zone_ids for datasource alicloud_zones ([#677](https://github.com/aliyun/terraform-provider-alicloud/issues/677))
- Improve redis and memcache instance testcases ([#676](https://github.com/aliyun/terraform-provider-alicloud/issues/676))
- Improve ecs instance testcases ([#675](https://github.com/aliyun/terraform-provider-alicloud/issues/675))

BUG FIXES:

- Fix oss bucket docs error ([#692](https://github.com/aliyun/terraform-provider-alicloud/issues/692))
- Fix pvtz 'Zone.VpcExists' error ([#691](https://github.com/aliyun/terraform-provider-alicloud/issues/691))
- Fix multi-k8s testcase failed error ([#683](https://github.com/aliyun/terraform-provider-alicloud/issues/683))
- Fix pvtz attchment Zone.NotExists error ([#682](https://github.com/aliyun/terraform-provider-alicloud/issues/682))
- Fix deleting ram role error ([#674](https://github.com/aliyun/terraform-provider-alicloud/issues/674))
- Fix k8s cluster worker_period_unit type error ([#672](https://github.com/aliyun/terraform-provider-alicloud/issues/672))

## 1.28.0 (January 16, 2019)

IMPROVEMENTS:

- Ots service support https ([#669](https://github.com/aliyun/terraform-provider-alicloud/issues/669))
- check vswitch id when creating instance  documentation ([#668](https://github.com/aliyun/terraform-provider-alicloud/issues/668))
- Improve pvtz attachment test updating case ([#663](https://github.com/aliyun/terraform-provider-alicloud/issues/663))
- add vswitch id checker when creating k8s clusters ([#656](https://github.com/aliyun/terraform-provider-alicloud/issues/656))
- Improve cen instance testcase to avoid mistake query ([#655](https://github.com/aliyun/terraform-provider-alicloud/issues/655))
- Improve route entry retry strategy to avoid concurrence issue ([#654](https://github.com/aliyun/terraform-provider-alicloud/issues/654))
- Offline drds resource from website results from drds does not support idempotent ([#653](https://github.com/aliyun/terraform-provider-alicloud/issues/653))
- Support customer endpoints in the provider ([#652](https://github.com/aliyun/terraform-provider-alicloud/issues/652))
- Reback image filter to meet many non-ecs testcase ([#649](https://github.com/aliyun/terraform-provider-alicloud/issues/649))
- Improve ecs instance testcase by update instance type ([#646](https://github.com/aliyun/terraform-provider-alicloud/issues/646))
- Support cs client setting customer endpoint ([#643](https://github.com/aliyun/terraform-provider-alicloud/issues/643))
- do not poll nodes when k8s cluster is stable ([#641](https://github.com/aliyun/terraform-provider-alicloud/issues/641))
- Improve pvtz_zone testcase by using rand ([#639](https://github.com/aliyun/terraform-provider-alicloud/issues/639))
- support for zero node clusters in swarm container service ([#638](https://github.com/aliyun/terraform-provider-alicloud/issues/638))
- Slb listener can not be updated when load balancer instance is shared-performance ([#637](https://github.com/aliyun/terraform-provider-alicloud/issues/637))
- Improve db_account testcase and its docs ([#635](https://github.com/aliyun/terraform-provider-alicloud/issues/635))
- Adding https_config options to the alicloud_cdn_domain resource ([#605](https://github.com/aliyun/terraform-provider-alicloud/issues/605))

BUG FIXES:

- Fix slb OperationFailed.TokenIsProcessing error ([#667](https://github.com/aliyun/terraform-provider-alicloud/issues/667))
- Fix deleting log project requestTimeout error ([#666](https://github.com/aliyun/terraform-provider-alicloud/issues/666))
- Fix cs_kubernetes setting int value error ([#665](https://github.com/aliyun/terraform-provider-alicloud/issues/665))
- Fix pvtz zone attaching vpc system busy error ([#660](https://github.com/aliyun/terraform-provider-alicloud/issues/660))
- Fix ecs and ess tags read bug with ignore system tag ([#659](https://github.com/aliyun/terraform-provider-alicloud/issues/659))
- Fix cs cluster not found error and improve its testcase ([#658](https://github.com/aliyun/terraform-provider-alicloud/issues/658))
- Fix deleting pvtz zone not exist and internal error ([#657](https://github.com/aliyun/terraform-provider-alicloud/issues/657))
- Fix pvtz throttling user bug and improve WrapError ([#650](https://github.com/aliyun/terraform-provider-alicloud/issues/650))
- Fix ess group describing error ([#644](https://github.com/aliyun/terraform-provider-alicloud/issues/644))
- Fix pvtz throttling user bug and add WrapError ([#642](https://github.com/aliyun/terraform-provider-alicloud/issues/642))
- Fix kvstore instance docs ([#636](https://github.com/aliyun/terraform-provider-alicloud/issues/636))

## 1.27.0 (January 08, 2019)

IMPROVEMENTS:

- Improve slb instance docs ([#632](https://github.com/aliyun/terraform-provider-alicloud/issues/632))
- Upgrade to Go 1.11 ([#629](https://github.com/aliyun/terraform-provider-alicloud/issues/629))
- Remove ots https schema because of in some region only supports http ([#630](https://github.com/aliyun/terraform-provider-alicloud/issues/630))
- Support https for log client ([#623](https://github.com/aliyun/terraform-provider-alicloud/issues/623))
- Support https for ram, cdn, kms and fc client ([#622](https://github.com/aliyun/terraform-provider-alicloud/issues/622))
- Support https for dns client ([#621](https://github.com/aliyun/terraform-provider-alicloud/issues/621))
- Support https for services client using official sdk ([#619](https://github.com/aliyun/terraform-provider-alicloud/issues/619))
- Support mns client https and improve mns testcase ([#618](https://github.com/aliyun/terraform-provider-alicloud/issues/618))
- Support oss client https ([#617](https://github.com/aliyun/terraform-provider-alicloud/issues/617))
- Support change kvstore instance charge type ([#602](https://github.com/aliyun/terraform-provider-alicloud/issues/602))
- add region checks to kubernetes, multiaz kubernetes, swarm clusters ([#607](https://github.com/aliyun/terraform-provider-alicloud/issues/607))
- Add forcenew for ess lifecycle hook name and improve ess testcase by random name ([#603](https://github.com/aliyun/terraform-provider-alicloud/issues/603))
- Improve ess configuration testcase ([#600](https://github.com/aliyun/terraform-provider-alicloud/issues/600))
- Improve kvstore and ess schedule testcase ([#599](https://github.com/aliyun/terraform-provider-alicloud/issues/599))
- Improve apigateway testcase ([#593](https://github.com/aliyun/terraform-provider-alicloud/issues/593))
- Improve ram, ess schedule and cdn testcase ([#592](https://github.com/aliyun/terraform-provider-alicloud/issues/592))
- Improve kvstore client token ([#586](https://github.com/aliyun/terraform-provider-alicloud/issues/586))

BUG FIXES:

- Fix api gateway deleteing app bug ([#633](https://github.com/aliyun/terraform-provider-alicloud/issues/633))
- Fix cs_kubernetes missing name error ([#625](https://github.com/aliyun/terraform-provider-alicloud/issues/625))
- Fix api gateway groups filter bug ([#624](https://github.com/aliyun/terraform-provider-alicloud/issues/624))
- Fix ots instance description force new bug ([#616](https://github.com/aliyun/terraform-provider-alicloud/issues/616))
- Fix oss bucket object testcase destroy bug ([#605](https://github.com/aliyun/terraform-provider-alicloud/issues/605))
- Fix deleting ess group timeout bug ([#604](https://github.com/aliyun/terraform-provider-alicloud/issues/604))
- Fix deleting mns subscription bug ([#601](https://github.com/aliyun/terraform-provider-alicloud/issues/601))
- bug fix for the input of cen bandwidth limit ([#598](https://github.com/aliyun/terraform-provider-alicloud/issues/598))
- Fix log service timeout error ([#594](https://github.com/aliyun/terraform-provider-alicloud/issues/594))
- Fix record not found issue if pvtz records are more than 50 ([#590](https://github.com/aliyun/terraform-provider-alicloud/issues/590))
- Fix cen instance and bandwidth multi regions test case bug ([#588](https://github.com/aliyun/terraform-provider-alicloud/issues/588))

## 1.26.0 (December 20, 2018)

FEATURES:

- **New Resource:** `alicloud_cs_managed_kubernetes` ([#563](https://github.com/aliyun/terraform-provider-alicloud/issues/563))

IMPROVEMENTS:

- Improve ram client endpoint ([#584](https://github.com/aliyun/terraform-provider-alicloud/issues/584))
- Remove useless sweeper depencences for alicloud_instance sweeper testcase ([#582](https://github.com/aliyun/terraform-provider-alicloud/issues/582))
- Improve kvstore backup policy testcase ([#580](https://github.com/aliyun/terraform-provider-alicloud/issues/580))
- Improve the describing endpoint ([#579](https://github.com/aliyun/terraform-provider-alicloud/issues/579))
- VPN gateway supports 200/500/1000M bandwidth ([#577](https://github.com/aliyun/terraform-provider-alicloud/issues/577))
- skip private ip test in some regions ([#575](https://github.com/aliyun/terraform-provider-alicloud/issues/575))
- Add timeout and retry for tablestore client and Improve its testcases ([#569](https://github.com/aliyun/terraform-provider-alicloud/issues/569))
- Modify kvstore_instance password to Optional and improve its testcases ([#567](https://github.com/aliyun/terraform-provider-alicloud/issues/567))
- Improve datasource alicloud_vpcs testcase ([#566](https://github.com/aliyun/terraform-provider-alicloud/issues/566))
- Improve dns_domains testcase ([#561](https://github.com/aliyun/terraform-provider-alicloud/issues/561))
- Improve ram_role_attachment testcase ([#560](https://github.com/aliyun/terraform-provider-alicloud/issues/560))
- support PrePaid instances, image_id to be set when creating k8s clusters ([#559](https://github.com/aliyun/terraform-provider-alicloud/issues/559))
- Add retry and timemout for fc client ([#557](https://github.com/aliyun/terraform-provider-alicloud/issues/557))
- Datasource alicloud_zones supports filter FunctionCompute ([#555](https://github.com/aliyun/terraform-provider-alicloud/issues/555))
- Fix a bug that caused the alicloud_dns_record.routing attribute ([#554](https://github.com/aliyun/terraform-provider-alicloud/issues/554))
- Modify router interface prepaid test case  documentation ([#552](https://github.com/aliyun/terraform-provider-alicloud/issues/552))
- Resource alicloud_ess_scalingconfiguration supports system_disk_size ([#551](https://github.com/aliyun/terraform-provider-alicloud/issues/551))
- Improve datahub project testcase ([#548](https://github.com/aliyun/terraform-provider-alicloud/issues/548))
- resource alicloud_slb_listener support server group ([#545](https://github.com/aliyun/terraform-provider-alicloud/issues/545))
- Improve ecs instance and disk testcase with common case ([#544](https://github.com/aliyun/terraform-provider-alicloud/issues/544))

BUG FIXES:

- Fix provider compile error on 32bit ([#585](https://github.com/aliyun/terraform-provider-alicloud/issues/585))
- Fix table store no such host error with deleting and updating ([#583](https://github.com/aliyun/terraform-provider-alicloud/issues/583))
- Fix pvtz_record RecordInvalidConflict bug ([#581](https://github.com/aliyun/terraform-provider-alicloud/issues/581))
- fixed bug in backup policy update ([#521](https://github.com/aliyun/terraform-provider-alicloud/issues/521))
- Fix docs eip_association ([#578](https://github.com/aliyun/terraform-provider-alicloud/issues/578))
- Fix a bug about instance charge type change ([#576](https://github.com/aliyun/terraform-provider-alicloud/issues/576))
- Fix describing endpoint failed error ([#574](https://github.com/aliyun/terraform-provider-alicloud/issues/574))
- Fix table store describing no such host error ([#572](https://github.com/aliyun/terraform-provider-alicloud/issues/572))
- Fix table store creating timeout error ([#571](https://github.com/aliyun/terraform-provider-alicloud/issues/571))
- Fix kvstore instance class update error ([#570](https://github.com/aliyun/terraform-provider-alicloud/issues/570))
- Fix ess_scaling_group import bugs and improve ess schedule testcase ([#565](https://github.com/aliyun/terraform-provider-alicloud/issues/565))
- Fix alicloud rds related IncorrectStatus bug ([#558](https://github.com/aliyun/terraform-provider-alicloud/issues/558))
- Fix alicloud_fc_trigger's config diff bug ([#556](https://github.com/aliyun/terraform-provider-alicloud/issues/556))
- Fix oss bucket deleting failed error ([#550](https://github.com/aliyun/terraform-provider-alicloud/issues/550))
- Fix potential bugs of datahub and ram when the resource has been deleted ([#546](https://github.com/aliyun/terraform-provider-alicloud/issues/546))
- Fix pvtz_record describing bug ([#543](https://github.com/aliyun/terraform-provider-alicloud/issues/543))

## 1.25.0 (November 30, 2018)

IMPROVEMENTS:

- return a empty list when there is no any data source ([#540](https://github.com/aliyun/terraform-provider-alicloud/issues/540))
- Skip automatically the testcases which does not support API gateway ([#538](https://github.com/aliyun/terraform-provider-alicloud/issues/538))
- Improve common bandwidth package test case and remove PayBy95 ([#530](https://github.com/aliyun/terraform-provider-alicloud/issues/530))
- Update resource drds supported regions ([#534](https://github.com/aliyun/terraform-provider-alicloud/issues/534))
- Remove DB instance engine_version limitation ([#528](https://github.com/aliyun/terraform-provider-alicloud/issues/528))
- Skip automatically the testcases which does not support route table and classic drds ([#526](https://github.com/aliyun/terraform-provider-alicloud/issues/526))
- Skip automatically the testcases which does not support classic regions ([#524](https://github.com/aliyun/terraform-provider-alicloud/issues/524))
- datasource alicloud_slbs support tags ([#523](https://github.com/aliyun/terraform-provider-alicloud/issues/523))
- resouce alicloud_slb support tags ([#522](https://github.com/aliyun/terraform-provider-alicloud/issues/522))
- Skip automatically the testcases which does not support multi az regions ([#518](https://github.com/aliyun/terraform-provider-alicloud/issues/518))
- Add some region limitation guide for sone resources ([#517](https://github.com/aliyun/terraform-provider-alicloud/issues/517))
- Skip automatically the testcases which does not support some known regions ([#516](https://github.com/aliyun/terraform-provider-alicloud/issues/516))
- create instance with runinstances ([#514](https://github.com/aliyun/terraform-provider-alicloud/issues/514))
- support eni amount in data source instance types ([#512](https://github.com/aliyun/terraform-provider-alicloud/issues/512))
- Add a docs guides/getting-account to help user learn alibaba cloud account ([#510](https://github.com/aliyun/terraform-provider-alicloud/issues/510))

BUG FIXES:

- Fix route_entry concurrence bug and improve it testcases ([#537](https://github.com/aliyun/terraform-provider-alicloud/issues/537))
- Fix router interface prepaid purchase ([#529](https://github.com/aliyun/terraform-provider-alicloud/issues/529))
- Fix fc_service sweeper test bug ([#536](https://github.com/aliyun/terraform-provider-alicloud/issues/536))
- Fix drds creating VPC instance bug by adding vpc_id ([#531](https://github.com/aliyun/terraform-provider-alicloud/issues/531))
- fix a snat_entry bug without set id to empty ([#525](https://github.com/aliyun/terraform-provider-alicloud/issues/525))
- fix a bug of ram_use display name ([#519](https://github.com/aliyun/terraform-provider-alicloud/issues/519))
- fix a bug of instance testcase ([#513](https://github.com/aliyun/terraform-provider-alicloud/issues/513))
- Fix pvtz resource priority bug ([#511](https://github.com/aliyun/terraform-provider-alicloud/issues/511))

## 1.24.0 (November 21, 2018)

FEATURES:

- **New Resource:** `alicloud_drds_instance` ([#446](https://github.com/aliyun/terraform-provider-alicloud/issues/446))

IMPROVEMENTS:

- Improve drds_instance docs ([#509](https://github.com/aliyun/terraform-provider-alicloud/issues/509))
- Add a new test case for drds_instance ([#508](https://github.com/aliyun/terraform-provider-alicloud/issues/508))
- Improve provider config with Trim method ([#504](https://github.com/aliyun/terraform-provider-alicloud/issues/504))
- api gateway skip app relevant tests ([#500](https://github.com/aliyun/terraform-provider-alicloud/issues/500))
- update api resource that support to deploy api ([#498](https://github.com/aliyun/terraform-provider-alicloud/issues/498))
- Skip ram_groups a test case ([#496](https://github.com/aliyun/terraform-provider-alicloud/issues/496))
- support disk resize ([#490](https://github.com/aliyun/terraform-provider-alicloud/issues/490))
- cancel the limit of system disk size ([#489](https://github.com/aliyun/terraform-provider-alicloud/issues/489))
- Improve docs alicloud_db_database and alicloud_cs_kubernetes ([#488](https://github.com/aliyun/terraform-provider-alicloud/issues/488))
- Support creating data disk with instance ([#484](https://github.com/aliyun/terraform-provider-alicloud/issues/484))

BUG FIXES:

- Fix the sweeper test for CEN and CEN bandwidth package ([#505](https://github.com/aliyun/terraform-provider-alicloud/issues/505))
- Fix pvtz_zone_record update bug ([#503](https://github.com/aliyun/terraform-provider-alicloud/issues/503))
- Fix network_interface_attachment docs error ([#502](https://github.com/aliyun/terraform-provider-alicloud/issues/502))
- fix fix datahub bug when visit region of ap-southeast-1 ([#499](https://github.com/aliyun/terraform-provider-alicloud/issues/499))
- Fix examples/mns-topic parameter error ([#497](https://github.com/aliyun/terraform-provider-alicloud/issues/497))
- Fix db_connection not found error when deleting ([#495](https://github.com/aliyun/terraform-provider-alicloud/issues/495))
- fix error about the docs format  ([#492](https://github.com/aliyun/terraform-provider-alicloud/issues/492))

## 1.23.0 (November 13, 2018)

FEATURES:

- **New Resource:** `alicloud_api_gateway_app_attachment` ([#478](https://github.com/aliyun/terraform-provider-alicloud/issues/478))
- **New Resource:** `alicloud_network_interface_attachment` ([#474](https://github.com/aliyun/terraform-provider-alicloud/issues/474))
- **New Resource:** `alicloud_api_gateway_vpc_access` ([#472](https://github.com/aliyun/terraform-provider-alicloud/issues/472))
- **New Resource:** `alicloud_network_interface` ([#469](https://github.com/aliyun/terraform-provider-alicloud/issues/469))
- **New Resource:** `alicloud_common_bandwidth_package` ([#468](https://github.com/aliyun/terraform-provider-alicloud/issues/468))
- **New Data Source:** `alicloud_network_interfaces` ([#475](https://github.com/aliyun/terraform-provider-alicloud/issues/475))
- **New Data Source:** `alicloud_api_gateway_apps` ([#467](https://github.com/aliyun/terraform-provider-alicloud/issues/467))

IMPROVEMENTS:

- Add a new region eu-west-1 ([#486](https://github.com/aliyun/terraform-provider-alicloud/issues/486))
- remove unreachable codes ([#479](https://github.com/aliyun/terraform-provider-alicloud/issues/479))
- support enable/disable security enhancement strategy of alicloud_instance ([#471](https://github.com/aliyun/terraform-provider-alicloud/issues/471))
- alicloud_slb_listener support idle_timeout/request_timeout ([#463](https://github.com/aliyun/terraform-provider-alicloud/issues/463))

BUG FIXES:

- Fix cs_application cluster not found ([#480](https://github.com/aliyun/terraform-provider-alicloud/issues/480))
- fix the bug of security_group inner_access bug ([#477](https://github.com/aliyun/terraform-provider-alicloud/issues/477))
- Fix pagenumber built error ([#470](https://github.com/aliyun/terraform-provider-alicloud/issues/470))
- Fix cs_application cluster not found ([#480](https://github.com/aliyun/terraform-provider-alicloud/issues/480))

## 1.22.0 (November 02, 2018)

FEATURES:

- **New Resource:** `alicloud_api_gateway_api` ([#457](https://github.com/aliyun/terraform-provider-alicloud/issues/457))
- **New Resource:** `alicloud_api_gateway_app` ([#462](https://github.com/aliyun/terraform-provider-alicloud/issues/462))
- **New Reource:** `alicloud_common_bandwidth_package` ([#454](https://github.com/aliyun/terraform-provider-alicloud/issues/454))
- **New Data Source:** `alicloud_api_gateway_apis` ([#458](https://github.com/aliyun/terraform-provider-alicloud/issues/458))
- **New Data Source:** `cen_region_route_entries` ([#442](https://github.com/aliyun/terraform-provider-alicloud/issues/442))
- **New Data Source:** `alicloud_slb_ca_certificates` ([#452](https://github.com/aliyun/terraform-provider-alicloud/issues/452))

IMPROVEMENTS:

- Use product code to get common request domain ([#466](https://github.com/aliyun/terraform-provider-alicloud/issues/466))
- KVstore instance password supports at sign ([#465](https://github.com/aliyun/terraform-provider-alicloud/issues/465))
- Correct docs spelling error ([#464](https://github.com/aliyun/terraform-provider-alicloud/issues/464))
- alicloud_log_service : support update project and shard auto spit ([#461](https://github.com/aliyun/terraform-provider-alicloud/issues/461))
- Correct datasource alicloud_cen_route_entries docs error ([#460](https://github.com/aliyun/terraform-provider-alicloud/issues/460))
- Remove CDN default configuration ([#450](https://github.com/aliyun/terraform-provider-alicloud/issues/450))

BUG FIXES:

- set number of cen instances five for normal alicloud account testcases ([#459](https://github.com/aliyun/terraform-provider-alicloud/issues/459))

## 1.21.0 (October 30, 2018)

FEATURES:

- **New Data Source:** `alicloud_slb_server_certificates` ([#444](https://github.com/aliyun/terraform-provider-alicloud/issues/444))
- **New Data Source:** `alicloud_slb_acls` ([#443](https://github.com/aliyun/terraform-provider-alicloud/issues/443))
- **New Resource:** `alicloud_slb_ca_certificate` ([#438](https://github.com/aliyun/terraform-provider-alicloud/issues/438))
- **New Resource:** `alicloud_slb_server_certificate` ([#436](https://github.com/aliyun/terraform-provider-alicloud/issues/436))

IMPROVEMENTS:

- resource alicloud_slb_listener tcp protocol support established_timeout parameter ([#440](https://github.com/aliyun/terraform-provider-alicloud/issues/440))

BUG FIXES:

- Fix mns resource docs bug ([#441](https://github.com/aliyun/terraform-provider-alicloud/issues/441))

## 1.20.0 (October 22, 2018)

FEATURES:

- **New Resource:** `alicloud_slb_acl` ([#413](https://github.com/aliyun/terraform-provider-alicloud/issues/413))
- **New Resource:** `alicloud_cen_route_entry` ([#415](https://github.com/aliyun/terraform-provider-alicloud/issues/415))
- **New Data Source:** `alicloud_cen_route_entries` ([#424](https://github.com/aliyun/terraform-provider-alicloud/issues/424))

IMPROVEMENTS:

- Improve datahub_project sweeper test ([#435](https://github.com/aliyun/terraform-provider-alicloud/issues/435))
- Modify mns test case name ([#434](https://github.com/aliyun/terraform-provider-alicloud/issues/434))
- Improve fc_service sweeper test ([#433](https://github.com/aliyun/terraform-provider-alicloud/issues/433))
- Support provider thread safety ([#432](https://github.com/aliyun/terraform-provider-alicloud/issues/432))
- add tags to security group ([#423](https://github.com/aliyun/terraform-provider-alicloud/issues/423))
- Resource router_interface support PrePaid ([#425](https://github.com/aliyun/terraform-provider-alicloud/issues/425))
- resource alicloud_slb_listener support acl ([#426](https://github.com/aliyun/terraform-provider-alicloud/issues/426))
- change child instance type Vbr to VBR and replace some const variables ([#422](https://github.com/aliyun/terraform-provider-alicloud/issues/422))
- add slb_internet_enabled to Kubernetes Cluster ([#421](https://github.com/aliyun/terraform-provider-alicloud/issues/421))
- Hide AliCloud HaVip Attachment resource docs because of it is not public totally ([#420](https://github.com/aliyun/terraform-provider-alicloud/issues/420))
- Improve examples/ots-table ([#417](https://github.com/aliyun/terraform-provider-alicloud/issues/417))
- Improve examples ecs-vpc, ecs-new-vpc and api-gateway ([#416](https://github.com/aliyun/terraform-provider-alicloud/issues/416))

BUG FIXES:

- Fix reources' id description bugs ([#428](https://github.com/aliyun/terraform-provider-alicloud/issues/428))
- Fix alicloud_ess_scaling_configuration setting data_disk failed ([#427](https://github.com/aliyun/terraform-provider-alicloud/issues/427))

## 1.19.0 (October 13, 2018)

FEATURES:

- **New Resource:** `alicloud_api_gateway_group` ([#409](https://github.com/aliyun/terraform-provider-alicloud/issues/409))
- **New Resource:** `alicloud_datahub_subscription` ([#405](https://github.com/aliyun/terraform-provider-alicloud/issues/405))
- **New Resource:** `alicloud_datahub_topic` ([#404](https://github.com/aliyun/terraform-provider-alicloud/issues/404))
- **New Resource:** `alicloud_datahub_project` ([#403](https://github.com/aliyun/terraform-provider-alicloud/issues/403))
- **New Data Source:** `alicloud_api_gateway_groups` ([#412](https://github.com/aliyun/terraform-provider-alicloud/issues/412))
- **New Data Source:** `alicloud_cen_bandwidth_limits` ([#402](https://github.com/aliyun/terraform-provider-alicloud/issues/402))

IMPROVEMENTS:

- added need_slb attribute to cs swarm ([#414](https://github.com/aliyun/terraform-provider-alicloud/issues/414))
- Add new example/datahub ([#407](https://github.com/aliyun/terraform-provider-alicloud/issues/407))
- Add new example/datahub ([#406](https://github.com/aliyun/terraform-provider-alicloud/issues/406))
- Format examples ([#397](https://github.com/aliyun/terraform-provider-alicloud/issues/397))
- Add new example/kvstore ([#396](https://github.com/aliyun/terraform-provider-alicloud/issues/396))
- Remove useless datasource cache file ([#395](https://github.com/aliyun/terraform-provider-alicloud/issues/395))
- Add new example/pvtz ([#394](https://github.com/aliyun/terraform-provider-alicloud/issues/394))
- Improve example/ecs-key-pair ([#393](https://github.com/aliyun/terraform-provider-alicloud/issues/393))
- Change key pair file mode to 400 ([#392](https://github.com/aliyun/terraform-provider-alicloud/issues/392))

BUG FIXES:

- fix kubernetes's new_nat_gateway issue ([#410](https://github.com/aliyun/terraform-provider-alicloud/issues/410))
- modify the mns err info ([#400](https://github.com/aliyun/terraform-provider-alicloud/issues/400))
- Skip havip test case ([#399](https://github.com/aliyun/terraform-provider-alicloud/issues/399))
- modify the sweeptest nameprefix ([#398](https://github.com/aliyun/terraform-provider-alicloud/issues/398))

## 1.18.0 (October 09, 2018)

FEATURES:

- **New Resource:** `alicloud_havip` ([#378](https://github.com/aliyun/terraform-provider-alicloud/issues/378))
- **New Resource:** `alicloud_havip_attachment` ([#388](https://github.com/aliyun/terraform-provider-alicloud/issues/388))
- **New Resource:** `alicloud_mns_topic_subscription` ([#376](https://github.com/aliyun/terraform-provider-alicloud/issues/376))
- **New Resource:** `alicloud_route_table_attachment` ([#362](https://github.com/aliyun/terraform-provider-alicloud/issues/362))
- **New Resource:** `alicloud_cen_bandwidth_limit` ([#361](https://github.com/aliyun/terraform-provider-alicloud/issues/361))
- **New Resource:** `alicloud_mns_topic` ([#374](https://github.com/aliyun/terraform-provider-alicloud/issues/374))
- **New Resource:** `alicloud_mns_queue` ([#365](https://github.com/aliyun/terraform-provider-alicloud/issues/365))
- **New Resource:** `alicloud_cen_bandwidth_package_attachment` ([#354](https://github.com/aliyun/terraform-provider-alicloud/issues/354))
- **New Resource:** `alicloud_route_table` ([#356](https://github.com/aliyun/terraform-provider-alicloud/issues/356))
- **New Data Source:** `alicloud_mns_queues` ([#382](https://github.com/aliyun/terraform-provider-alicloud/issues/382))
- **New Data Source:** `alicloud_mns_topics` ([#384](https://github.com/aliyun/terraform-provider-alicloud/issues/384))
- **New Data Source:** `alicloud_mns_topic_subscriptions` ([#386](https://github.com/aliyun/terraform-provider-alicloud/issues/386))
- **New Data Source:** `alicloud_cen_bandwidth_packages` ([#367](https://github.com/aliyun/terraform-provider-alicloud/issues/367))
- **New Data Source:** `alicloud_vpn_connections` ([#366](https://github.com/aliyun/terraform-provider-alicloud/issues/366))
- **New Data Source:** `alicloud_vpn_gateways` ([#363](https://github.com/aliyun/terraform-provider-alicloud/issues/363))
- **New Data Source:** `alicloud_vpn_customer_gateways` ([#364](https://github.com/aliyun/terraform-provider-alicloud/issues/364))
- **New Data Source:** `alicloud_cen_instances` ([#342](https://github.com/aliyun/terraform-provider-alicloud/issues/342))

IMPROVEMENTS:

- Improve resource ram_policy's document validatefunc ([#385](https://github.com/aliyun/terraform-provider-alicloud/issues/385))
- RAM support useragent ([#383](https://github.com/aliyun/terraform-provider-alicloud/issues/383))
- add node_cidr_mas and log_config, fix worker_data_disk issue ([#368](https://github.com/aliyun/terraform-provider-alicloud/issues/368))
- Improve WaitForRouteTable and WaitForRouteTableAttachment method ([#375](https://github.com/aliyun/terraform-provider-alicloud/issues/375))
- Correct Function Compute conn ([#371](https://github.com/aliyun/terraform-provider-alicloud/issues/371))
- Improve datasource `images`'s docs ([#370](https://github.com/aliyun/terraform-provider-alicloud/issues/370))
- add worker_data_disk_category and worker_data_disk_size to kubernetes creation ([#355](https://github.com/aliyun/terraform-provider-alicloud/issues/355))

BUG FIXES:

- Fix alicloud_ram_user_policy_attachment EntityNotExist.User error ([#381](https://github.com/aliyun/terraform-provider-alicloud/issues/381))
- Add parameter 'force_delete' to support deleting 'PrePaid' instance ([#377](https://github.com/aliyun/terraform-provider-alicloud/issues/377))
- Add wait time to fix random detaching disk error ([#373](https://github.com/aliyun/terraform-provider-alicloud/issues/373))
- Fix cen_instances markdown ([#372](https://github.com/aliyun/terraform-provider-alicloud/issues/372))

## 1.17.0 (September 22, 2018)

FEATURES:

- **New Data Source:** `alicloud_fc_triggers` ([#351](https://github.com/aliyun/terraform-provider-alicloud/pull/351))
- **New Data Source:** `alicloud_oss_bucket_objects` ([#350](https://github.com/aliyun/terraform-provider-alicloud/pull/350))
- **New Data Source:** `alicloud_fc_functions` ([#349](https://github.com/aliyun/terraform-provider-alicloud/pull/349))
- **New Data Source:** `alicloud_fc_services` ([#348](https://github.com/aliyun/terraform-provider-alicloud/pull/348))
- **New Data Source:** `alicloud_oss_buckets` ([#345](https://github.com/aliyun/terraform-provider-alicloud/pull/345))
- **New Data Source:** `alicloud_disks` ([#343](https://github.com/aliyun/terraform-provider-alicloud/pull/343))
- **New Resource:** `alicloud_cen_bandwidth_package` ([#333](https://github.com/aliyun/terraform-provider-alicloud/pull/333))

IMPROVEMENTS:

- Update OSS Resources' link to English ([#352](https://github.com/aliyun/terraform-provider-alicloud/pull/352))
- Improve example/kubernetes to support multi-az ([#344](https://github.com/aliyun/terraform-provider-alicloud/pull/344))

## 1.16.0 (September 16, 2018)

FEATURES:

- **New Resource:** `alicloud_cen_instance_attachment` ([#327](https://github.com/aliyun/terraform-provider-alicloud/pull/327))

IMPROVEMENTS:

- Allow setting the scaling group balancing policy ([#339](https://github.com/aliyun/terraform-provider-alicloud/pull/339))
- cs_kubernetes supports multi-az ([#222](https://github.com/aliyun/terraform-provider-alicloud/pull/222))
- Improve client token using timestemp ([#326](https://github.com/aliyun/terraform-provider-alicloud/pull/326))

BUG FIXES:

- Fix alicloud db connection ([#341](https://github.com/aliyun/terraform-provider-alicloud/pull/341))
- Fix knstore productId ([#338](https://github.com/aliyun/terraform-provider-alicloud/pull/338))
- Fix retriving kvstore multi zones bug ([#337](https://github.com/aliyun/terraform-provider-alicloud/pull/337))
- Fix kvstore instance period bug ([#335](https://github.com/aliyun/terraform-provider-alicloud/pull/335))
- Fix kvstore docs bug ([#334](https://github.com/aliyun/terraform-provider-alicloud/pull/334))

## 1.15.0 (September 07, 2018)

FEATURES:

- **New Resource:** `alicloud_kvstore_backup_policy` ([#331](https://github.com/aliyun/terraform-provider-alicloud/pull/331))
- **New Resource:** `alicloud_kvstore_instance` ([#330](https://github.com/aliyun/terraform-provider-alicloud/pull/330))
- **New Data Source:** `alicloud_kvstore_instances` ([#329](https://github.com/aliyun/terraform-provider-alicloud/pull/329))
- **New Resource:** `alicloud_ess_alarm` ([#328](https://github.com/aliyun/terraform-provider-alicloud/pull/328))
- **New Resource:** `alicloud_ssl_vpn_client_cert` ([#317](https://github.com/aliyun/terraform-provider-alicloud/pull/317))
- **New Resource:** `alicloud_cen_instance` ([#312](https://github.com/aliyun/terraform-provider-alicloud/pull/312))
- **New Data Source:** `alicloud_slb_server_groups`  ([#324](https://github.com/aliyun/terraform-provider-alicloud/pull/324))
- **New Data Source:** `alicloud_slb_rules`  ([#323](https://github.com/aliyun/terraform-provider-alicloud/pull/323))
- **New Data Source:** `alicloud_slb_listeners`  ([#323](https://github.com/aliyun/terraform-provider-alicloud/pull/323))
- **New Data Source:** `alicloud_slb_attachments`  ([#322](https://github.com/aliyun/terraform-provider-alicloud/pull/322))
- **New Data Source:** `alicloud_slbs`  ([#321](https://github.com/aliyun/terraform-provider-alicloud/pull/321))
- **New Data Source:** `alicloud_account`  ([#319](https://github.com/aliyun/terraform-provider-alicloud/pull/319))
- **New Resource:** `alicloud_ssl_vpn_server` ([#313](https://github.com/aliyun/terraform-provider-alicloud/pull/313))

IMPROVEMENTS:

- Support sweeper to clean some resources coming from failed testcases ([#326](https://github.com/aliyun/terraform-provider-alicloud/pull/326))
- Improve function compute tst cases ([#325](https://github.com/aliyun/terraform-provider-alicloud/pull/325))
- Improve fc test case using new datasource `alicloud_account` ([#320](https://github.com/aliyun/terraform-provider-alicloud/pull/320))
- Base64 encode ESS scaling config user_data ([#315](https://github.com/aliyun/terraform-provider-alicloud/pull/315))
- Retrieve the account_id automatically if needed ([#314](https://github.com/aliyun/terraform-provider-alicloud/pull/314))

BUG FIXES:

- Fix DNS tests falied error ([#318](https://github.com/aliyun/terraform-provider-alicloud/pull/318))
- Fix DB database not found error ([#316](https://github.com/aliyun/terraform-provider-alicloud/pull/316))

## 1.14.0 (August 31, 2018)

FEATURES:

- **New Resource:** `alicloud_vpn_connection` ([#304](https://github.com/aliyun/terraform-provider-alicloud/pull/304))
- **New Resource:** `alicloud_vpn_customer_gateway` ([#299](https://github.com/aliyun/terraform-provider-alicloud/pull/299))

IMPROVEMENTS:

- Add 'force' to make key pair affect immediately ([#310](https://github.com/aliyun/terraform-provider-alicloud/pull/310))
- Improve http proxy support ([#307](https://github.com/aliyun/terraform-provider-alicloud/pull/307))
- Add flags to skip tests that use features not supported in all regions ([#306](https://github.com/aliyun/terraform-provider-alicloud/pull/306))
- Improve data source dns_domains test case ([#305](https://github.com/aliyun/terraform-provider-alicloud/pull/305))
- Change SDK config timeout ([#302](https://github.com/aliyun/terraform-provider-alicloud/pull/302))
- Support ClientToken for some request ([#301](https://github.com/aliyun/terraform-provider-alicloud/pull/301))
- Enlarge sdk default timeout to fix some timeout scenario ([#300](https://github.com/aliyun/terraform-provider-alicloud/pull/300))

BUG FIXES:

- Fix container cluster SDK timezone error ([#308](https://github.com/aliyun/terraform-provider-alicloud/pull/308))
- Fix network products throttling error ([#303](https://github.com/aliyun/terraform-provider-alicloud/pull/303))

## 1.13.0 (August 28, 2018)

FEATURES:

- **New Resource:** `alicloud_vpn_gateway` ([#298](https://github.com/aliyun/terraform-provider-alicloud/pull/298))
- **New Data Source:** `alicloud_mongo_instances` ([#221](https://github.com/aliyun/terraform-provider-alicloud/pull/221))
- **New Data Source:** `alicloud_pvtz_zone_records` ([#288](https://github.com/aliyun/terraform-provider-alicloud/pull/288))
- **New Data Source:** `alicloud_pvtz_zones` ([#287](https://github.com/aliyun/terraform-provider-alicloud/pull/287))
- **New Resource:** `alicloud_pvtz_zone_record` ([#286](https://github.com/aliyun/terraform-provider-alicloud/pull/286))
- **New Resource:** `alicloud_pvtz_zone_attachment` ([#285](https://github.com/aliyun/terraform-provider-alicloud/pull/285))
- **New Resource:** `alicloud_pvtz_zone` ([#284](https://github.com/aliyun/terraform-provider-alicloud/pull/284))
- **New Resource:** `alicloud_ess_lifecycle_hook` ([#283](https://github.com/aliyun/terraform-provider-alicloud/pull/283))
- **New Data Source:** `alicloud_router_interfaces` ([#269](https://github.com/aliyun/terraform-provider-alicloud/pull/269))

IMPROVEMENTS:

- Check pvtzconn error ([#295](https://github.com/aliyun/terraform-provider-alicloud/pull/295))
- For internationalize tests ([#294](https://github.com/aliyun/terraform-provider-alicloud/pull/294))
- Improve data source docs ([#293](https://github.com/aliyun/terraform-provider-alicloud/pull/293))
- Add SLB PayByBandwidth test case ([#292](https://github.com/aliyun/terraform-provider-alicloud/pull/292))
- Update vpc sdk to support new resource VPN gateway ([#291](https://github.com/aliyun/terraform-provider-alicloud/pull/291))
- Improve snat entry test case ([#290](https://github.com/aliyun/terraform-provider-alicloud/pull/290))
- Allow empty list of SLBs as arg to ESG ([#289](https://github.com/aliyun/terraform-provider-alicloud/pull/289))
- Improve docs vroute_entry ([#281](https://github.com/aliyun/terraform-provider-alicloud/pull/281))
- Improve examples/router_interface ([#278](https://github.com/aliyun/terraform-provider-alicloud/pull/278))
- Improve SLB instance test case ([#274](https://github.com/aliyun/terraform-provider-alicloud/pull/274))
- Improve alicloud_router_interface's test case ([#272](https://github.com/aliyun/terraform-provider-alicloud/pull/272))
- Improve data source alicloud_regions's test case ([#271](https://github.com/aliyun/terraform-provider-alicloud/pull/271))
- Add notes about ordering between two alicloud_router_interface_connections ([#270](https://github.com/aliyun/terraform-provider-alicloud/pull/270))
- Improve docs spelling error ([#268](https://github.com/aliyun/terraform-provider-alicloud/pull/268))
- ECS instance support more tags and update instance test cases ([#267](https://github.com/aliyun/terraform-provider-alicloud/pull/267))
- Improve OSS bucket test case ([#266](https://github.com/aliyun/terraform-provider-alicloud/pull/266))
- Fixing a broken link ([#265](https://github.com/aliyun/terraform-provider-alicloud/pull/265))
- Allow creation of slb vserver group with 0 servers ([#264](https://github.com/aliyun/terraform-provider-alicloud/pull/264))
- Improve SLB test cases results from international regions does support PayByBandwidth and ' Guaranteed-performance' instance ([#263](https://github.com/aliyun/terraform-provider-alicloud/pull/263))
- Improve EIP test cases results from international regions does support PayByBandwidth ([#262](https://github.com/aliyun/terraform-provider-alicloud/pull/262))
- Improve ESS test cases results from some region does support Classic Network ([#261](https://github.com/aliyun/terraform-provider-alicloud/pull/261))
- Recover nat gateway bandwidth pacakges to meet stock user requirements ([#260](https://github.com/aliyun/terraform-provider-alicloud/pull/260))
- Resource alicloud_slb_listener supports new field 'x-forwarded-for' ([#259](https://github.com/aliyun/terraform-provider-alicloud/pull/259))
- Resource alicloud_slb_listener supports new field 'gzip' ([#258](https://github.com/aliyun/terraform-provider-alicloud/pull/258))

BUG FIXES:

- Fix getting oss endpoint timeout error ([#282](https://github.com/aliyun/terraform-provider-alicloud/pull/282))
- Fix router interface connection error when 'opposite_interface_owner_id' is empty ([#277](https://github.com/aliyun/terraform-provider-alicloud/pull/277))
- Fix router interface connection error and deleting error ([#275](https://github.com/aliyun/terraform-provider-alicloud/pull/275))
- Fix disk detach error and improve test using dynamic zone and region ([#273](https://github.com/aliyun/terraform-provider-alicloud/pull/273))

## 1.12.0 (August 10, 2018)

IMPROVEMENTS:

- Improve `make build` ([#256](https://github.com/aliyun/terraform-provider-alicloud/pull/256))
- Improve examples slb and slb-vpc by modifying 'paybytraffic' to 'PayByTraffic' ([#256](https://github.com/aliyun/terraform-provider-alicloud/pull/256))
- Improve example/router-interface by adding resource alicloud_router_interface_connection ([#255](https://github.com/aliyun/terraform-provider-alicloud/pull/255))
- Support more specification of router interface ([#253](https://github.com/aliyun/terraform-provider-alicloud/pull/253))
- Improve resource alicloud_fc_service docs ([#252](https://github.com/aliyun/terraform-provider-alicloud/pull/252))
- Modify resource alicloud_fc_function 'handler' is required ([#251](https://github.com/aliyun/terraform-provider-alicloud/pull/251))
- Resource alicloud_router_interface support "import" function ([#249](https://github.com/aliyun/terraform-provider-alicloud/pull/249))
- Deprecate some field of alicloud_router_interface fields and use new resource instead ([#248](https://github.com/aliyun/terraform-provider-alicloud/pull/248))
- *New Resource*: _alicloud_router_interface_connection_ ([#247](https://github.com/aliyun/terraform-provider-alicloud/pull/247))

BUG FIXES:

- Fix network resource throttling error ([#257](https://github.com/aliyun/terraform-provider-alicloud/pull/257))
- Fix resource alicloud_fc_trigger "source_arn" inputting empty error ([#253](https://github.com/aliyun/terraform-provider-alicloud/pull/253))
- Fix describing vpcs with name_regex no results error ([#250](https://github.com/aliyun/terraform-provider-alicloud/pull/250))
- Fix creating slb listener in international region failed error ([#246](https://github.com/aliyun/terraform-provider-alicloud/pull/246))

## 1.11.0 (August 08, 2018)

IMPROVEMENTS:

- Resource alicloud_eip support name and description ([#244](https://github.com/aliyun/terraform-provider-alicloud/pull/244))
- Resource alicloud_eip support PrePaid ([#243](https://github.com/aliyun/terraform-provider-alicloud/pull/243))
- Correct version writting error ([#241](https://github.com/aliyun/terraform-provider-alicloud/pull/241))
- Change slb go sdk to official repo ([#240](https://github.com/aliyun/terraform-provider-alicloud/pull/240))
- Remove useless file website/fc_service.html.markdown ([#239](https://github.com/aliyun/terraform-provider-alicloud/pull/239))
- Update Go version to 1.10.1 to match new sdk ([#237](https://github.com/aliyun/terraform-provider-alicloud/pull/237))
- Support http(s) proxy ([#236](https://github.com/aliyun/terraform-provider-alicloud/pull/236))
- Add guide for installing goimports ([#233](https://github.com/aliyun/terraform-provider-alicloud/pull/233))
- Improve the makefile and README ([#232](https://github.com/aliyun/terraform-provider-alicloud/pull/232))

BUG FIXES:

- Fix losing key pair error after updating ecs instance ([#245](https://github.com/aliyun/terraform-provider-alicloud/pull/245))
- Fix BackendServer.configuring error when creating slb rule ([#242](https://github.com/aliyun/terraform-provider-alicloud/pull/242))
- Fix bug "...zoneinfo.zip: no such file or directory" happened in windows. ([#238](https://github.com/aliyun/terraform-provider-alicloud/pull/238))
- Fix ess_scalingrule InvalidScalingRuleId.NotFound error ([#234](https://github.com/aliyun/terraform-provider-alicloud/pull/234))

## 1.10.0 (July 27, 2018)

IMPROVEMENTS:

- Rds supports to create 10.0 PostgreSQL instance. ([#230](https://github.com/aliyun/terraform-provider-alicloud/pull/230))
- *New Resource*: _alicloud_fc_trigger_ ([#228](https://github.com/aliyun/terraform-provider-alicloud/pull/228))
- *New Resource*: _alicloud_fc_function_ ([#227](https://github.com/aliyun/terraform-provider-alicloud/pull/227))
- *New Resource*: _alicloud_fc_service_ 30([#226](https://github.com/aliyun/terraform-provider-alicloud/pull/226))
- Support new field 'instance_name' for _alicloud_ots_table_ ([#225](https://github.com/aliyun/terraform-provider-alicloud/pull/225))
- *New Resource*: _alicloud_ots_instance_attachment_ ([#224](https://github.com/aliyun/terraform-provider-alicloud/pull/224))
- *New Resource*: _alicloud_ots_instance_ ([#223](https://github.com/aliyun/terraform-provider-alicloud/pull/223))

BUG FIXES:

- Fix Snat entry not found error ([#229](https://github.com/aliyun/terraform-provider-alicloud/pull/229))

## 1.9.6 (July 24, 2018)

IMPROVEMENTS:

- Remove the number limitation of vswitch_ids, slb_ids and db_instance_ids ([#219](https://github.com/aliyun/terraform-provider-alicloud/pull/219))
- Reduce test nat gateway cost ([#218](https://github.com/aliyun/terraform-provider-alicloud/pull/218))
- Support creating zero-node swarm cluster ([#217](https://github.com/aliyun/terraform-provider-alicloud/pull/217))
- Improve security group and rule data source test case ([#216](https://github.com/aliyun/terraform-provider-alicloud/pull/216))
- Improve dns record resource test case ([#215](https://github.com/aliyun/terraform-provider-alicloud/pull/215))
- Improve test case destroy method ([#214](https://github.com/aliyun/terraform-provider-alicloud/pull/214))
- Improve ecs instance resource test case ([#213](https://github.com/aliyun/terraform-provider-alicloud/pull/213))
- Improve cdn resource test case ([#212](https://github.com/aliyun/terraform-provider-alicloud/pull/212))
- Improve kms resource test case ([#211](https://github.com/aliyun/terraform-provider-alicloud/pull/211))
- Improve key pair resource test case ([#210](https://github.com/aliyun/terraform-provider-alicloud/pull/210))
- Improve rds resource test case ([#209](https://github.com/aliyun/terraform-provider-alicloud/pull/209))
- Improve disk resource test case ([#208](https://github.com/aliyun/terraform-provider-alicloud/pull/208))
- Improve eip resource test case ([#207](https://github.com/aliyun/terraform-provider-alicloud/pull/207))
- Improve scaling service resource test case ([#206](https://github.com/aliyun/terraform-provider-alicloud/pull/206))
- Improve vpc and vswitch resource test case ([#205](https://github.com/aliyun/terraform-provider-alicloud/pull/205))
- Improve slb resource test case ([#204](https://github.com/aliyun/terraform-provider-alicloud/pull/204))
- Improve security group resource test case ([#203](https://github.com/aliyun/terraform-provider-alicloud/pull/203))
- Improve ram resource test case ([#202](https://github.com/aliyun/terraform-provider-alicloud/pull/202))
- Improve container cluster resource test case ([#201](https://github.com/aliyun/terraform-provider-alicloud/pull/201))
- Improve cloud monitor resource test case ([#200](https://github.com/aliyun/terraform-provider-alicloud/pull/200))
- Improve route and router interface resource test case ([#199](https://github.com/aliyun/terraform-provider-alicloud/pull/199))
- Improve dns resource test case ([#198](https://github.com/aliyun/terraform-provider-alicloud/pull/198))
- Improve oss resource test case ([#197](https://github.com/aliyun/terraform-provider-alicloud/pull/197))
- Improve ots table resource test case ([#196](https://github.com/aliyun/terraform-provider-alicloud/pull/196))
- Improve nat gateway resource test case ([#195](https://github.com/aliyun/terraform-provider-alicloud/pull/195))
- Improve log resource test case ([#194](https://github.com/aliyun/terraform-provider-alicloud/pull/194))
- Support changing ecs charge type from Prepaid to PostPaid ([#192](https://github.com/aliyun/terraform-provider-alicloud/pull/192))
- Add method to compare json template is equal ([#187](https://github.com/aliyun/terraform-provider-alicloud/pull/187))
- Remove useless file ([#191](https://github.com/aliyun/terraform-provider-alicloud/pull/191))

BUG FIXES:

- Fix CS kubernetes read error and CS app timeout ([#217](https://github.com/aliyun/terraform-provider-alicloud/pull/217))
- Fix getting location connection error ([#193](https://github.com/aliyun/terraform-provider-alicloud/pull/193))
- Fix CS kubernetes connection error ([#190](https://github.com/aliyun/terraform-provider-alicloud/pull/190))
- Fix Oss bucket diff error ([#189](https://github.com/aliyun/terraform-provider-alicloud/pull/189))

NOTES:

- From version 1.9.6, the deprecated resource alicloud_ram_alias file has been removed and the resource has been
replaced by alicloud_ram_account_alias. Details refer to [pull 191](https://github.com/aliyun/terraform-provider-alicloud/pull/191/commits/e3fd74591230ccb545bb4309b674d6df33b716b9)

## 1.9.5 (June 20, 2018)

IMPROVEMENTS:

- Improve log machine group docs ([#186](https://github.com/aliyun/terraform-provider-alicloud/pull/186))
- Support sts token for some resources ([#185](https://github.com/aliyun/terraform-provider-alicloud/pull/185))
- Support user agent for log service ([#184](https://github.com/aliyun/terraform-provider-alicloud/pull/184))
- *New Resource*: _alicloud_log_machine_group_ ([#183](https://github.com/aliyun/terraform-provider-alicloud/pull/183))
- *New Resource*: _alicloud_log_store_index_ ([#182](https://github.com/aliyun/terraform-provider-alicloud/pull/182))
- *New Resource*: _alicloud_log_store_ ([#181](https://github.com/aliyun/terraform-provider-alicloud/pull/181))
- *New Resource*: _alicloud_log_project_ ([#180](https://github.com/aliyun/terraform-provider-alicloud/pull/180))
- Improve example about cs_kubernetes ([#179](https://github.com/aliyun/terraform-provider-alicloud/pull/179))
- Add losing docs about cs_kubernetes ([#178](https://github.com/aliyun/terraform-provider-alicloud/pull/178))

## 1.9.4 (June 08, 2018)

IMPROVEMENTS:

- cs_kubernetes supports output worker nodes and master nodes ([#177](https://github.com/aliyun/terraform-provider-alicloud/pull/177))
- cs_kubernetes supports to output kube config and certificate ([#176](https://github.com/aliyun/terraform-provider-alicloud/pull/176))
- Add a example to deploy mysql and wordpress on kubernetes ([#175](https://github.com/aliyun/terraform-provider-alicloud/pull/175))
- Add a example to create swarm and deploy wordpress on it ([#174](https://github.com/aliyun/terraform-provider-alicloud/pull/174))
- Change ECS and ESS sdk to official go sdk ([#173](https://github.com/aliyun/terraform-provider-alicloud/pull/173))


## 1.9.3 (May 27, 2018)

IMPROVEMENTS:

- *New Data Source*: _alicloud_db_instances_ ([#161](https://github.com/aliyun/terraform-provider-alicloud/pull/161))
- Support to set auto renew for ECS instance ([#172](https://github.com/aliyun/terraform-provider-alicloud/pull/172))
- Improve cs_kubernetes, slb_listener and db_database docs ([#171](https://github.com/aliyun/terraform-provider-alicloud/pull/171))
- Add missing code for describing RDS zones ([#170](https://github.com/aliyun/terraform-provider-alicloud/pull/170))
- Add docs notes for windows os([#169](https://github.com/aliyun/terraform-provider-alicloud/pull/169))
- Add filter parameters and export parameters for instance types data source. ([#168](https://github.com/aliyun/terraform-provider-alicloud/pull/168))
- Add filter parameters for zones data source. ([#167](https://github.com/aliyun/terraform-provider-alicloud/pull/167))
- Remove kubernetes work_number limitation ([#165](https://github.com/aliyun/terraform-provider-alicloud/pull/165))
- Improve kubernetes examples ([#163](https://github.com/aliyun/terraform-provider-alicloud/pull/163))

BUG FIXES:

- Fix getting some instance types failed bug ([#166](https://github.com/aliyun/terraform-provider-alicloud/pull/166))
- Fix kubernetes out range index error ([#164](https://github.com/aliyun/terraform-provider-alicloud/pull/164))

## 1.9.2 (May 09, 2018)

IMPROVEMENTS:

- *New Resource*: _alicloud_ots_table_ ([#162](https://github.com/aliyun/terraform-provider-alicloud/pull/162))
- Fix SLB listener "OperationBusy" error ([#159](https://github.com/aliyun/terraform-provider-alicloud/pull/159))
- Prolong waiting time for creating kubernetes cluster to avoid timeout ([#158](https://github.com/aliyun/terraform-provider-alicloud/pull/158))
- Support load endpoint from environment variable or specified file ([#157](https://github.com/aliyun/terraform-provider-alicloud/pull/157))
- Update example ([#155](https://github.com/aliyun/terraform-provider-alicloud/pull/155))

BUG FIXES:

- Fix modifying instance host name failed bug ([#160](https://github.com/aliyun/terraform-provider-alicloud/pull/160))
- Fix SLB listener "OperationBusy" error ([#159](https://github.com/aliyun/terraform-provider-alicloud/pull/159))
- Fix deleting forward table not found error ([#154](https://github.com/aliyun/terraform-provider-alicloud/pull/154))
- Fix deleting slb listener error ([#150](https://github.com/aliyun/terraform-provider-alicloud/pull/150))
- Fix creating vswitch error ([#149](https://github.com/aliyun/terraform-provider-alicloud/pull/149))

## 1.9.1 (April 13, 2018)

IMPROVEMENTS:

- *New Resource*: _alicloud_cms_alarm_ ([#146](https://github.com/aliyun/terraform-provider-alicloud/pull/146))
- *New Resource*: _alicloud_cs_application_ ([#136](https://github.com/aliyun/terraform-provider-alicloud/pull/136))
- *New Datasource*: _alicloud_security_group_rules_ ([#135](https://github.com/aliyun/terraform-provider-alicloud/pull/135))
- Output application attribution service block ([#141](https://github.com/aliyun/terraform-provider-alicloud/pull/141))
- Output swarm attribution 'vpc_id' ([#140](https://github.com/aliyun/terraform-provider-alicloud/pull/140))
- Support to release eip after deploying swarm cluster. ([#139](https://github.com/aliyun/terraform-provider-alicloud/pull/139))
- Output swarm and kubernetes's nodes information and other attribution ([#138](https://github.com/aliyun/terraform-provider-alicloud/pull/138))
- Modify `size` to `node_number` ([#137](https://github.com/aliyun/terraform-provider-alicloud/pull/137))
- Set swarm ID before waiting its status ([#134](https://github.com/aliyun/terraform-provider-alicloud/pull/134))
- Add 'is_outdated' for cs_swarm and cs_kubernetes ([#133](https://github.com/aliyun/terraform-provider-alicloud/pull/133))
- Add warning when creating postgresql and ppas database ([#132](https://github.com/aliyun/terraform-provider-alicloud/pull/132))
- Add kubernetes example ([#142](https://github.com/aliyun/terraform-provider-alicloud/pull/142))
- Update sdk to support user-agent ([#143](https://github.com/aliyun/terraform-provider-alicloud/pull/143))
- Add eip unassociation retry times to avoid needless error ([#144](https://github.com/aliyun/terraform-provider-alicloud/pull/144))
- Add connections output for kubernetes cluster ([#145](https://github.com/aliyun/terraform-provider-alicloud/pull/145))

BUG FIXES:

- Fix vpc not found when vpc has been deleted ([#131](https://github.com/aliyun/terraform-provider-alicloud/pull/131))


## 1.9.0 (March 19, 2018)

IMPROVEMENTS:

- *New Resource*: _alicloud_cs_kubernetes_ ([#129](https://github.com/aliyun/terraform-provider-alicloud/pull/129))
- *New DataSource*: _alicloud_eips_ ([#123](https://github.com/aliyun/terraform-provider-alicloud/pull/123))
- Add server_group_id to slb listener resource ([#122](https://github.com/aliyun/terraform-provider-alicloud/pull/122))
- Rename _alicloud_container_cluster_ to _alicloud_cs_swarm_ ([#128](https://github.com/aliyun/terraform-provider-alicloud/pull/128))

BUG FIXES:

- Fix vpc description validate ([#125](https://github.com/aliyun/terraform-provider-alicloud/pull/125))
- Update SDK version to fix unresolving endpoint issue ([#126](https://github.com/aliyun/terraform-provider-alicloud/pull/126))
- Add waiting time after ECS bind ECS to ensure network is ok ([#127](https://github.com/aliyun/terraform-provider-alicloud/pull/127))

## 1.8.1 (March 09, 2018)

IMPROVEMENTS:

- DB instance supports multiple zone ([#120](https://github.com/aliyun/terraform-provider-alicloud/pull/120))
- Data source zones support to retrieve multiple zone ([#119](https://github.com/aliyun/terraform-provider-alicloud/pull/119))
- VPC supports alibaba cloud official go sdk ([#118](https://github.com/aliyun/terraform-provider-alicloud/pull/118))

BUG FIXES:

- Fix not found db instance bug when allocating connection ([#121](https://github.com/aliyun/terraform-provider-alicloud/pull/121))


## 1.8.0 (March 02, 2018)

IMPROVEMENTS:

- Support golang version 1.9 ([#114](https://github.com/aliyun/terraform-provider-alicloud/pull/114))
- RDS supports alibaba cloud official go sdk ([#113](https://github.com/aliyun/terraform-provider-alicloud/pull/113))
- Deprecated 'in_use' in eips datasource to fix conflict ([#115](https://github.com/aliyun/terraform-provider-alicloud/pull/115))
- Add encrypted argument to alicloud_disk resource（[#116](https://github.com/aliyun/terraform-provider-alicloud/pull/116))

BUG FIXES:

- Fix reading router interface failed bug ([#117](https://github.com/aliyun/terraform-provider-alicloud/pull/117))

## 1.7.2 (February 09, 2018)

IMPROVEMENTS:

- *New DataSource*: _alicloud_eips_ ([#110](https://github.com/aliyun/terraform-provider-alicloud/pull/110))
- *New DataSource*: _alicloud_vswitches_ ([#109](https://github.com/aliyun/terraform-provider-alicloud/pull/109))
- Support inner network segregation in one security group ([#112](https://github.com/aliyun/terraform-provider-alicloud/pull/112))

BUG FIXES:

- Fix creating Classic instance failed result in role_name ([#111](https://github.com/aliyun/terraform-provider-alicloud/pull/111))
- Fix eip is not exist in nat gateway when creating snat ([#108](https://github.com/aliyun/terraform-provider-alicloud/pull/108))

## 1.7.1 (February 02, 2018)

IMPROVEMENTS:

- Support setting instance_name for ESS scaling configuration ([#107](https://github.com/aliyun/terraform-provider-alicloud/pull/107))
- Support multiple vswitches for ESS scaling group and output slbIds and dbIds ([#105](https://github.com/aliyun/terraform-provider-alicloud/pull/105))
- Support to set internet_max_bandwidth_out is 0 for ESS configuration ([#103](https://github.com/aliyun/terraform-provider-alicloud/pull/103))
- Modify EIP default to PayByTraffic for international account ([#101](https://github.com/aliyun/terraform-provider-alicloud/pull/101))
- Deprecate nat gateway fileds 'spec' and 'bandwidth_packages' ([#100](https://github.com/aliyun/terraform-provider-alicloud/pull/100))
- Support to associate EIP with SLB and Nat Gateway ([#99](https://github.com/aliyun/terraform-provider-alicloud/pull/99))

BUG FIXES:

- fix a bug that can't create multiple VPC, vswitch and nat gateway at one time ([#102](https://github.com/aliyun/terraform-provider-alicloud/pull/102))
- fix a bug that can't import instance 'role_name' ([#104](https://github.com/aliyun/terraform-provider-alicloud/pull/104))
- fix a bug that creating ESS scaling group and configuration results from 'Throttling' ([#106](https://github.com/aliyun/terraform-provider-alicloud/pull/106))

## 1.7.0 (January 25, 2018)

IMPROVEMENTS:

- *New Resource*: _alicloud_kms_key_ ([#91](https://github.com/aliyun/terraform-provider-alicloud/pull/91))
- *New DataSource*: _alicloud_kms_keys_ ([#93](https://github.com/aliyun/terraform-provider-alicloud/pull/93))
- *New DataSource*: _alicloud_instances_ ([#94](https://github.com/aliyun/terraform-provider-alicloud/pull/94))
- Add a new output field "arn" for _alicloud_kms_key_ ([#92](https://github.com/aliyun/terraform-provider-alicloud/pull/92))
- Add a new field "specification" for _alicloud_slb_ ([#95](https://github.com/aliyun/terraform-provider-alicloud/pull/95))
- Improve security group rule's port range for "-1/-1" ([#96](https://github.com/aliyun/terraform-provider-alicloud/pull/96))

BUG FIXES:

- fix slb invalid status error when launching ESS scaling group ([#97](https://github.com/aliyun/terraform-provider-alicloud/pull/97))

## 1.6.2 (January 22, 2018)

IMPROVEMENTS:

- Modify db_connection prefix default value to "instance_id + 'tf'"([#90](https://github.com/aliyun/terraform-provider-alicloud/pull/90))
- Modify db_connection ID to make it more simple while importing it([#90](https://github.com/aliyun/terraform-provider-alicloud/pull/90))
- Add wait method to avoid useless status error while creating/modifying account or privilege or connection or database([#90](https://github.com/aliyun/terraform-provider-alicloud/pull/90))
- Support to set instnace name for RDS ([#88](https://github.com/aliyun/terraform-provider-alicloud/pull/88))
- Avoid container cluster cidr block conflicts with vswitch's ([#88](https://github.com/aliyun/terraform-provider-alicloud/pull/88))
- Output resource import information ([#87](https://github.com/aliyun/terraform-provider-alicloud/pull/87))

BUG FIXES:

- fix instance id not found and instane status not supported bug([#90](https://github.com/aliyun/terraform-provider-alicloud/pull/90))
- fix deleting slb_attachment resource failed bug ([#86](https://github.com/aliyun/terraform-provider-alicloud/pull/86))


## 1.6.1 (January 18, 2018)

IMPROVEMENTS:

- Support to modify instance type and network spec ([#84](https://github.com/aliyun/terraform-provider-alicloud/pull/84))
- Avoid needless error when creating security group rule ([#83](https://github.com/aliyun/terraform-provider-alicloud/pull/83))

BUG FIXES:

- fix creating cluster container failed bug ([#85](https://github.com/aliyun/terraform-provider-alicloud/pull/85))


## 1.6.0 (January 15, 2018)

IMPROVEMENTS:

- *New Resource*: _alicloud_ess_attachment_ ([#80](https://github.com/aliyun/terraform-provider-alicloud/pull/80))
- *New Resource*: _alicloud_slb_rule_ ([#79](https://github.com/aliyun/terraform-provider-alicloud/pull/79))
- *New Resource*: _alicloud_slb_server_group_ ([#78](https://github.com/aliyun/terraform-provider-alicloud/pull/78))
- Support Spot Instance ([#77](https://github.com/aliyun/terraform-provider-alicloud/pull/77))
- Output tip message when international account create SLB failed ([#75](https://github.com/aliyun/terraform-provider-alicloud/pull/75))
- Standardize the order of imports packages ([#74](https://github.com/aliyun/terraform-provider-alicloud/pull/74))
- Add "weight" for slb_attachment to improve the resource ([#81](https://github.com/aliyun/terraform-provider-alicloud/pull/81))

BUG FIXES:

- fix allocating RDS public connection conflict error ([#76](https://github.com/aliyun/terraform-provider-alicloud/pull/76))

## 1.5.3 (January 9, 2018)

BUG FIXES:
  * fix getting OSS endpoint failed error  ([#73](https://github.com/aliyun/terraform-provider-alicloud/pull/73))
  * fix describing dns record not found when deleting record ([#73](https://github.com/aliyun/terraform-provider-alicloud/pull/73))

## 1.5.2 (January 8, 2018)

BUG FIXES:
  * fix creating rds 'Prepaid' instance failed error  ([#70](https://github.com/aliyun/terraform-provider-alicloud/pull/70))

## 1.5.1 (January 5, 2018)

BUG FIXES:
  * modify security_token to Optional ([#69](https://github.com/aliyun/terraform-provider-alicloud/pull/69))

## 1.5.0 (January 4, 2018)

IMPROVEMENTS:

- *New Resource*: _alicloud_db_database_ ([#68](https://github.com/aliyun/terraform-provider-alicloud/pull/68))
- *New Resource*: _alicloud_db_backup_policy_ ([#68](https://github.com/aliyun/terraform-provider-alicloud/pull/68))
- *New Resource*: _alicloud_db_connection_ ([#67](https://github.com/aliyun/terraform-provider-alicloud/pull/67))
- *New Resource*: _alicloud_db_account_ ([#66](https://github.com/aliyun/terraform-provider-alicloud/pull/66))
- *New Resource*: _alicloud_db_account_privilege_ ([#66](https://github.com/aliyun/terraform-provider-alicloud/pull/66))
- resource/db_instance: remove some field to new resource ([#65](https://github.com/aliyun/terraform-provider-alicloud/pull/65))
- resource/instance: support to modify private ip, vswitch_id and instance charge type ([#65](https://github.com/aliyun/terraform-provider-alicloud/pull/65))

BUG FIXES:

- resource/dns-record: Fix dns record still exist after deleting it ([#65](https://github.com/aliyun/terraform-provider-alicloud/pull/65))
- resource/instance: fix deleting route entry error ([#69](https://github.com/aliyun/terraform-provider-alicloud/pull/69))


## 1.2.0 (December 15, 2017)

IMPROVEMENTS:
- resource/slb: wait for SLB active before return back ([#61](https://github.com/aliyun/terraform-provider-alicloud/pull/61))

BUG FIXES:

- resource/dns-record: Fix setting dns priority failed ([#58](https://github.com/aliyun/terraform-provider-alicloud/pull/58))
- resource/dns-record: Fix ESS attachs SLB failed ([#59](https://github.com/aliyun/terraform-provider-alicloud/pull/59))
- resource/dns-record: Fix security group not found error ([#59](https://github.com/aliyun/terraform-provider-alicloud/pull/59))


## 1.0.0 (December 11, 2017)

IMPROVEMENTS:

- *New Resource*: _alicloud_slb_listener_ ([#53](https://github.com/aliyun/terraform-provider-alicloud/pull/53))
- *New Resource*: _alicloud_cdn_domain_ ([#52](https://github.com/aliyun/terraform-provider-alicloud/pull/52))
- *New Resource*: _alicloud_dns_ ([#51](https://github.com/aliyun/terraform-provider-alicloud/pull/51))
- *New Resource*: _alicloud_dns_group_ ([#51](https://github.com/aliyun/terraform-provider-alicloud/pull/51))
- *New Resource*: _alicloud_dns_record_ ([#51](https://github.com/aliyun/terraform-provider-alicloud/pull/51))
- *New Resource*: _alicloud_ram_account_alias_ ([#50](https://github.com/aliyun/terraform-provider-alicloud/pull/50))
- *New Resource*: _alicloud_ram_login_profile_ ([#50](https://github.com/aliyun/terraform-provider-alicloud/pull/50))
- *New Resource*: _alicloud_ram_access_key_ ([#50](https://github.com/aliyun/terraform-provider-alicloud/pull/50))
- *New Resource*: _alicloud_ram_group_ ([#49](https://github.com/aliyun/terraform-provider-alicloud/pull/49))
- *New Resource*: _alicloud_ram_group_membership_ ([#49](https://github.com/aliyun/terraform-provider-alicloud/pull/49))
- *New Resource*: _alicloud_ram_group_policy_attachment_ ([#49](https://github.com/aliyun/terraform-provider-alicloud/pull/49))
- *New Resource*: _alicloud_ram_role_ ([#48](https://github.com/aliyun/terraform-provider-alicloud/pull/48))
- *New Resource*: _alicloud_ram_role_attachment_ ([#48](https://github.com/aliyun/terraform-provider-alicloud/pull/48))
- *New Resource*: _alicloud_ram_role_polocy_attachment_ ([#48](https://github.com/aliyun/terraform-provider-alicloud/pull/48))
- *New Resource*: _alicloud_container_cluster_ ([#47](https://github.com/aliyun/terraform-provider-alicloud/pull/47))
- *New Resource:* _alicloud_ram_policy_ ([#46](https://github.com/aliyun/terraform-provider-alicloud/pull/46))
- *New Resource*: _alicloud_ram_user_policy_attachment_ ([#46](https://github.com/aliyun/terraform-provider-alicloud/pull/46))
- *New Resource* _alicloud_ram_user_ ([#44](https://github.com/aliyun/terraform-provider-alicloud/pull/44))
- *New Datasource* _alicloud_ram_policies_ ([#46](https://github.com/aliyun/terraform-provider-alicloud/pull/46))
- *New Datasource* _alicloud_ram_users_ ([#44](https://github.com/aliyun/terraform-provider-alicloud/pull/44))
- *New Datasource*: _alicloud_ram_roles_ ([#48](https://github.com/aliyun/terraform-provider-alicloud/pull/48))
- *New Datasource*: _alicloud_ram_account_aliases_ ([#50](https://github.com/aliyun/terraform-provider-alicloud/pull/50))
- *New Datasource*: _alicloud_dns_domains_ ([#51](https://github.com/aliyun/terraform-provider-alicloud/pull/51))
- *New Datasource*: _alicloud_dns_groups_ ([#51](https://github.com/aliyun/terraform-provider-alicloud/pull/51))
- *New Datasource*: _alicloud_dns_records_ ([#51](https://github.com/aliyun/terraform-provider-alicloud/pull/51))
- resource/instance: add new parameter `role_name` ([#48](https://github.com/aliyun/terraform-provider-alicloud/pull/48))
- resource/slb: remove slb schema field `listeners` and using new listener resource to replace ([#55](https://github.com/aliyun/terraform-provider-alicloud/pull/55))
- resource/ess_scaling_configuration: add new parameters `key_name`, `role_name`, `user_data`, `force_delete` and `tags` ([#54](https://github.com/aliyun/terraform-provider-alicloud/pull/54))
- resource/ess_scaling_configuration: remove it importing ([#54](https://github.com/aliyun/terraform-provider-alicloud/pull/54))
- resource: format not found error ([#55](https://github.com/aliyun/terraform-provider-alicloud/pull/55))
- website: improve resource docs ([#56](https://github.com/aliyun/terraform-provider-alicloud/pull/56))
- examples: add new examples, like oss, key_pair, router_interface and so on ([#56](https://github.com/aliyun/terraform-provider-alicloud/pull/56))

- Added support for importing:
  - `alicloud_container_cluster` ([#47](https://github.com/aliyun/terraform-provider-alicloud/pull/47))
  - `alicloud_ram_policy` ([#46](https://github.com/aliyun/terraform-provider-alicloud/pull/46))
  - `alicloud_ram_user` ([#44](https://github.com/aliyun/terraform-provider-alicloud/pull/44))
  - `alicloud_ram_role` ([#48](https://github.com/aliyun/terraform-provider-alicloud/pull/48))
  - `alicloud_ram_groups` ([#49](https://github.com/aliyun/terraform-provider-alicloud/pull/49))
  - `alicloud_ram_login_profile` ([#50](https://github.com/aliyun/terraform-provider-alicloud/pull/50))
  - `alicloud_dns` ([#51](https://github.com/aliyun/terraform-provider-alicloud/pull/51))
  - `alicloud_dns_record` ([#51](https://github.com/aliyun/terraform-provider-alicloud/pull/51))
  - `alicloud_slb_listener` ([#53](https://github.com/aliyun/terraform-provider-alicloud/pull/53))
  - `alicloud_security_group` ([#55](https://github.com/aliyun/terraform-provider-alicloud/pull/55))
  - `alicloud_slb` ([#55](https://github.com/aliyun/terraform-provider-alicloud/pull/55))
  - `alicloud_vswitch` ([#55](https://github.com/aliyun/terraform-provider-alicloud/pull/55))
  - `alicloud_vroute_entry` ([#55](https://github.com/aliyun/terraform-provider-alicloud/pull/55))

BUG FIXES:

- resource/vroute_entry: Fix building route_entry concurrency issue ([#55](https://github.com/aliyun/terraform-provider-alicloud/pull/55))
- resource/vswitch: Fix building vswitch concurrency issue ([#55](https://github.com/aliyun/terraform-provider-alicloud/pull/55))
- resource/router_interface: Fix building router interface concurrency issue ([#55](https://github.com/aliyun/terraform-provider-alicloud/pull/55))
- resource/vpc: Fix building vpc concurrency issue ([#55](https://github.com/aliyun/terraform-provider-alicloud/pull/55))
- resource/slb_attachment: Fix attaching slb failed ([#55](https://github.com/aliyun/terraform-provider-alicloud/pull/55))

## 0.1.1 (December 11, 2017)

IMPROVEMENTS:

- *New Resource:* _alicloud_key_pair_ ([#27](https://github.com/aliyun/terraform-provider-alicloud/pull/27))
- *New Resource*: _alicloud_key_pair_attachment_ ([#28](https://github.com/aliyun/terraform-provider-alicloud/pull/28))
- *New Resource*: _alicloud_router_interface_ ([#40](https://github.com/aliyun/terraform-provider-alicloud/pull/40))
- *New Resource:* _alicloud_oss_bucket_ ([#10](https://github.com/aliyun/terraform-provider-alicloud/pull/10))
- *New Resource*: _alicloud_oss_bucket_object_ ([#14](https://github.com/aliyun/terraform-provider-alicloud/pull/14))
- *New Datasource* _alicloud_key_pairs_ ([#30](https://github.com/aliyun/terraform-provider-alicloud/pull/30))
- *New Datasource* _alicloud_vpcs_ ([#34](https://github.com/aliyun/terraform-provider-alicloud/pull/34))
- *New output_file* option for data sources: export data to a specified file ([#29](https://github.com/aliyun/terraform-provider-alicloud/pull/29))
- resource/instance:add new parameter `key_name` ([#31](https://github.com/aliyun/terraform-provider-alicloud/pull/31))
- resource/route_entry: new nexthop type 'RouterInterface' for route entry ([#41](https://github.com/aliyun/terraform-provider-alicloud/pull/41))
- resource/security_group_rule: Remove `cidr_ip` contribute "ConflictsWith" ([#39](https://github.com/aliyun/terraform-provider-alicloud/pull/39))
- resource/rds: add ability to change instance password ([#17](https://github.com/aliyun/terraform-provider-alicloud/pull/17))
- resource/rds: Add ability to import existing RDS resources ([#16](https://github.com/aliyun/terraform-provider-alicloud/pull/16))
- datasource/alicloud_zones: Add more options for filtering ([#19](https://github.com/aliyun/terraform-provider-alicloud/pull/19))
- Added support for importing:
  - `alicloud_vpc` ([#32](https://github.com/aliyun/terraform-provider-alicloud/pull/32))
  - `alicloud_route_entry` ([#33](https://github.com/aliyun/terraform-provider-alicloud/pull/33))
  - `alicloud_nat_gateway` ([#26](https://github.com/aliyun/terraform-provider-alicloud/pull/26))
  - `alicloud_ess_schedule` ([#25](https://github.com/aliyun/terraform-provider-alicloud/pull/25))
  - `alicloud_ess_scaling_group` ([#24](https://github.com/aliyun/terraform-provider-alicloud/pull/24))
  - `alicloud_instance` ([#23](https://github.com/aliyun/terraform-provider-alicloud/pull/23))
  - `alicloud_eip` ([#22](https://github.com/aliyun/terraform-provider-alicloud/pull/22))
  - `alicloud_disk` ([#21](https://github.com/aliyun/terraform-provider-alicloud/pull/21))

BUG FIXES:

- resource/disk_attachment: Fix issue attaching multiple disks and set disk_attachment's parameter 'device_name' as deprecated ([#9](https://github.com/aliyun/terraform-provider-alicloud/pull/9))
- resource/rds: Fix diff error about rds security_ips ([#13](https://github.com/aliyun/terraform-provider-alicloud/pull/13))
- resource/security_group_rule: Fix diff error when authorizing security group rules ([#15](https://github.com/aliyun/terraform-provider-alicloud/pull/15))
- resource/security_group_rule: Fix diff bug by modifying 'DestCidrIp' to 'DestGroupId' when running read ([#35](https://github.com/aliyun/terraform-provider-alicloud/pull/35))


## 0.1.0 (June 20, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
