## 1.55.2 (Unreleased)

IMPROVEMENTS:

- improve(db_readonly_instance):improve db_readonly_instance testcase [GH-1607]
- improve(drds): modified drds supported regions [GH-1605]
- improve(CI): change sweeper time [GH-1600]
- improve(rds): fix db_instance apply error after import [GH-1599]
- improve(ons_topic):retry when Throttling.User [GH-1598]
- Improve(ddoscoo): Improve its resource and datasource use common method [GH-1591]
- Improve(slb):slb support set AddressIpVersion [GH-1587]
- Improve(cs_kubernetes): Improve its resource and datasource use common method [GH-1584]
- Improve(cs_managed_kubernetes): Improve its resource and datasource use common method [GH-1581]

BUG FIXES:

- fix(ons):fix ons error Throttling.User [GH-1608]
- fix(ons): fix the create group error in testcase [GH-1604]

## 1.55.1 (August 23, 2019)

IMPROVEMENTS:

- improve(ons_instance): set instance name using random ([#1597](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1597))
- add support to Ipsec_pfs field be set with "disabled" and add example files ([#1589](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1589))
- improve(slb): sweep the protected slb ([#1588](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1588))
- Improve(ram): ram resources supports import ([#1586](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1586))
- improve(tags): modified test case to check the upper case letters in tags ([#1585](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1585))
- improve(Document):improve document demo about set ([#1580](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1580))
- Update RouteEntry Resource RouteEntryName Field ([#1578](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1578))
- improve(ci):supplement log ([#1577](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1577))
- improve(sdk):update alibaba-cloud-sdk-go(1.60.107) ([#1575](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1575))
- Rename resource name that is not start with a letter ([#1573](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1573))
- Improve(datahub_topic): Improve resource use common method ([#1565](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1565))
- Improve(datahub_subscription): Improve resource use common method ([#1556](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1556))
- Improve(datahub_project): Improve resource use common method ([#1555](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1555))

BUG FIXES:

- fix(vsw): fix bug from github issue ([#1593](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1593))
- fix(instance):update instance testcase ([#1590](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1590))
- fix(ci):fix CI statistics bug ([#1576](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1576))
- Fix typo ([#1574](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1574))
- fix(disks):fix dataSource test case bug ([#1566](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1566))

## 1.55.0 (August 16, 2019)

- **New Resource:** `alicloud_ess_notification` ([#1549](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1549))

IMPROVEMENTS:

- improve(key_pair):update key_pair document ([#1563](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1563))
- improve(CI): add default bucket and region for CI ([#1561](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1561))
- improve(CI): terraform CI log ([#1557](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1557))
- Improve(ots_instance_attachment): Improve its resource and datasource use common method ([#1552](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1552))
- Improve(ots_instance): Improve its resource and datasource use common method ([#1551](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1551))
- Improve(ram): ram policy attachment resources supports import ([#1550](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1550))
- Improve(ots_table): Improve its resource and datasource use common method ([#1546](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1546))
- Improve(router_interface): modified testcase multi count ([#1545](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1545))
- Improve(images): removed image alinux check in datasource ([#1543](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1543))
- Improve(logtail_config): Improve resource use common method ([#1500](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1500))

BUG FIXES:

- bugfix：throw notFoundError when scalingGroup is not found ([#1572](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1572))
- fix(sweep): modified the error return to run sweep completely ([#1569](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1569))
- fix(CI): remove the useless code ([#1564](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1564))
- fix(CI): fix pipeline grammar error ([#1562](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1562))
- Fix log document ([#1559](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1559))
- modify(cs): skip the testcases of cs_application and cs_swarm ([#1553](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1553))
- fix kvstore unexpected state 'Changing' ([#1539](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1539))

## 1.54.0 (August 12, 2019)

- **New Data Source:** `alicloud_slb_master_slave_server_groups` ([#1531](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1531))
- **New Resource:** `alicloud_slb_master_slave_server_group` ([#1531](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1531))
- **New Data Source:** `alicloud_instance_type_families` ([#1519](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1519))

IMPROVEMENTS:

- improve(provider):profile,role_arn,session_name,session_expiration support ENV ([#1537](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1537))
- support sg description ([#1536](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1536))
- support mac address ([#1535](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1535))
- improve(sdk): update sdk and modify api_gateway strconvs ([#1533](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1533))
- Improve(pvtz_zone_record): Improve resource use common method ([#1528](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1528))
- improve(alicloud_ess_scaling_group): support 'COST_OPTIMIZED' mode of autoscaling group ([#1527](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1527))
- Improve(pvtz_zone): Improve its and attachment resources use common method ([#1525](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1525))
- remove useless trigger in vpn ci ([#1522](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1522))
- Improve(cr_repo): Improve resource use common method ([#1515](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1515))
- Improve(cr_namespace): Improve resource use common method ([#1509](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1509))
- improve(kvstore): kvstore_instance resource supports timeouts setting ([#1445](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1445))

BUG FIXES:

- Fix(alicloud_logstore_index) Repair parameter description document ([#1532](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1532))
- fix(sweep): modified the region of prefixes ([#1526](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1526))
- fix(mongodb_instance): fix notfound error when describing it ([#1521](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1521))

## 1.53.0 (August 02, 2019)

- **New Resource:** `alicloud_ons_group` ([#1506](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1506))
- **New Resource:** `alicloud_ess_scalinggroup_vserver_groups` ([#1503](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1503))
- **New Resource:** `alicloud_slb_backend_server` ([#1498](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1498))
- **New Resource:** `alicloud_ons_topic` ([#1483](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1483))
- **New Data Source:** `alicloud_ons_groups` ([#1506](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1506))
- **New Data source:** `alicloud_slb_backend_servers` ([#1498](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1498))
- **New Data Source:** `alicloud_ons_topics` ([#1483](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1483))


IMPROVEMENTS:

- improve(dns_record): add diffsuppressfunc to avoid DomainRecordDuplicate error. ([#1518](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1518))
- remove useless import ([#1517](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1517))
- remove empty fields in managed k8s, add force_update, add multiple az support ([#1516](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1516))
- improve(fc_function):fc_function support sweeper ([#1513](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1513))
- improve(fc_trigger):fc_trigger support sweeper ([#1512](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1512))
- Improve(logtail_attachment): Improve resource use common method [[#1508](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1508)] 
- improve(slb):update testcase ([#1507](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1507))
- improve(disk):update disk_attachment ([#1501](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1501))
- add(slb_backend_server): slb backend server resource & data source ([#1498](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1498))
- Improve(log_machine_group): Improve resources use common method ([#1497](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1497))
- Improve(log_project): Improve resource use common method ([#1496](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1496))
- improve(network_interface): enhance sweeper test ([#1495](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1495))
- Improve(log_store): Improve resources use common method ([#1494](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1494))
- improve(instance_type):update testcase config ([#1493](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1493))
- Improve(mns_topic_subscription): Improve its resource use common method ([#1492](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1492))
- improve(disk):suppurt delete_auto_snapshot delete_with_instance enable_auto_snapshot ([#1491](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1491))
- Improve(mns_topic): Improve its resource use common method ([#1488](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1488))
- Improve(api_gateway): api_gateway_api added testcases ([#1487](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1487))
- Improve(mns_queue): Improve its resource use common method ([#1485](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1485))
- improve(customer_gateway):create add retry ([#1477](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1477))
- improve(gpdb): resources supports timeouts setting ([#1476](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1476))
- improve(fc_triggers): Added ids filter to datasource ([#1475](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1475))
- improve(fc_services): Added ids filter to datasource ([#1474](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1474))
- improve(fc_functions): Added ids filter to datasource ([#1473](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1473))
- improve(instance_types):update instance_types filter condition ([#1472](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1472))
- improve(pvtz_zone__domain): Added ids filter to datasource ([#1471](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1471))
- improve(cr_repos): Added names to datasource attributes ([#1470](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1470))
- improve(cr_namespaces): Added names to datasource attributes ([#1469](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1469))
- improve(cdn): Added region to domain name and modified sweep rules ([#1466](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1466))
- improve(ram_roles): Added ids filter to datasource ([#1461](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1461))
- improve(ram_users): Added ids filter to datasource ([#1459](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1459))
- improve(pvtz_zones): Added ids filter and added names to datasource attributes ([#1458](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1458))
- improve(nas_mount_targets): Added ids filter to datasource ([#1453](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1453))
- improve(nas_file_systems): Added descriptions to datasource attributes ([#1450](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1450))
- improve(nas_access_rules): Added ids filter to datasource ([#1448](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1448))
- improve(mongodb_instance): supports timeouts setting ([#1446](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1446))
- improve(nas_access_groups): Added names to its attributes ([#1444](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1444))
- improve(mns_topics): Added names to datasource attributes ([#1442](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1442))
- improve(mns_topic_subscriptions): Added names to datasource attributes ([#1441](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1441))
- improve(mns_queues): Added names to datasource attributes ([#1439](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1439))

BUG FIXES:

- Fix(logstore_index): Invalid update parameter change ([#1505](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1505))
- fix(api_gateway): fix can't get resource id when stage_names set ([#1486](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1486))
- fix(kvstore_instance): resource kvstore_instance add Retry while ModifyInstanceSpec err ([#1484](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1484))
- fix(cen): modified the timeouts of cen instance to avoid errors ([#1451](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1451))

## 1.52.2 (July 20, 2019)

IMPROVEMENTS:

- improve(eip_association): supporting to set PrivateIPAddress  documentation ([#1480](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1480))
- improve(mongodb_instances): Added ids filter to datasource ([#1478](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1478))
- improve(dns_domain): Added ids filter to datasource ([#1468](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1468))
- improve(cdn): Added retry to avoid ServiceBusy error ([#1467](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1467))
- improve(dns_records): Added ids filter to datasource ([#1464](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1464))
- improve(dns_groups): Added ids filter and added names to datasource attributes ([#1463](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1463))
- improve(stateConfig):update stateConfig error ([#1462](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1462))
- improve(kvstore): Added ids filter to datasource ([#1457](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1457))
- improve(cas): Added precheck to testcases ([#1456](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1456))
- improve(rds): db_instance and db_readonly_instance resource modify timeouts 20mins to 30mins ([#1455](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1455))
- add CI for the alicloud provider ([#1449](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1449))
- improve(api_gateway_apps): Deprecated api_id ([#1426](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1426))
- improve(api_gateway_apis): Added ids filter to datasource ([#1425](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1425))
- improve(slb_server_group): remove the maximum limitation of adding backend servers ([#1416](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1416))
- improve(cdn): cdn_domain_config added testcases ([#1405](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1405))

BUG FIXES:

- fix(kvstore_instance): resource kvstore_instance add Retry while ModifyInstanceSpec err ([#1465](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1465))
- fix(slb): fix slb testcase can not find instance types' bug ([#1454](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1454))

## 1.52.1 (July 16, 2019)

IMPROVEMENTS:

- improve(disk): support online resize ([#1447](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1447))
- improve(rds): db_readonly_instance resource supports timeouts setting ([#1438](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1438))
- improve(rds):improve db_readonly_instance TestAccAlicloudDBReadonlyInstance_multi testcase ([#1432](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1432))
- improve(key_pairs): Added ids filter to datasource ([#1431](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1431))
- improve(elasticsearch): Added ids filter and added descriptions to datasource attributes ([#1430](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1430))
- improve(drds): Added descriptions to attributes of datasource ([#1429](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1429))
- improve(rds):update ppas not support regions ([#1428](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1428))
- improve(api_gateway_groups): Added ids filter to datasource ([#1427](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1427))
- improve(docs): Reformat abnormal inline HCL code in docs ([#1423](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1423))
- improve(mns):modified mns_queues.html ([#1422](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1422))
- improve(rds): db_instance resource supports timeouts setting ([#1409](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1409))
- improve(kms): modified the args of kms_keys datasource ([#1407](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1407))
- improve(kms_key): modify the param `description` to forcenew ([#1406](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1406))

BUG FIXES:

- fix(db_instance): modified the target state of state config ([#1437](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1437))
- fix(db_readonly_instance): fix invalid status error when updating and deleting ([#1435](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1435))
- fix(ots_table): fix setting deviation_cell_version_in_sec error ([#1434](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1434))
- fix(db_backup_policy): resource db_backup_policy testcase use datasource db_instance_classes ([#1424](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1424))

## 1.52.0 (July 12, 2019)

- **New Datasource:** `alicloud_ons_instances` ([#1411](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1411))

IMPROVEMENTS:

- improve(vpc):add ids filter ([#1420](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1420))
- improve(db_instances): Added ids filter and added names to datasource attributes ([#1419](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1419))
- improve(cas): Added ids filter and added names to datasource attributes ([#1417](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1417))
- docs(format): Convert inline HCL configs to canonical format ([#1415](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1415))
- improve(gpdb_instance):add vpc name ([#1413](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1413))
- improve(provider): add a new parameter `skip_region_validation` in the provider config ([#1404](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1404))
- improve(cdn): cdn_domain support certificate config ([#1393](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1393))
- improve(rds): resource db_instance support update for instance_charge_type ([#1389](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1389))

BUG FIXES:

- fix(db_instance):fix db_instance testcase vsw availability_zone ([#1418](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1418))
- fix(api_gateway): modified the testcase to avoid errors ([#1410](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1410))
- fix(db_readonly_instance): extend the waiting time for spec modification ([#1408](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1408))
- fix(db_readonly_instance): add retryable error content in instance spec modification and deletion ([#1403](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1403))

## 1.51.0 (July 08, 2019)

- **New Date Source:** `alicloud_kvstore_instance_engines` ([#1371](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1371))
- **New Resource:** `alicloud_ons_instance` ([#1333](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1333))

IMPROVEMENTS:

- improve(db_instance): improve db_instance MAZ testcase ([#1391](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1391))
- improve(cs_kubernetes): add importIgnore parameters in the importer testcase ([#1387](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1387))
- Remove govendor commands in CI ([#1386](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1386))
- improve(slb_vserver_group): support attaching eni ([#1384](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1384))
- improve(db_instance_classes): add new parameter db_instance_class ([#1383](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1383))
- improve(images): Add os_name_en to the attributes of images datasource ([#1380](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1380))
- improve(disk): the snapshot_id conflicts with encrypted ([#1378](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1378))
- Improve(cs_kubernetes): add some importState ignore fields in the importer testcase ([#1377](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1377))
- Improve(oss_bucket): Add names for its attributes of datasource ([#1374](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1374))
- improve(common_test):update common_test for terraform 0.12 ([#1372](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1372))
- Improve(cs_kubernetes): add import ignore parameter `log_config` ([#1370](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1370))
- improve(slb):support slb instance delete protection ([#1369](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1369))
- improve(slb_rule): support health check config ([#1367](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1367))
- Improve(oss_bucket_object): Improve its use common method ([#1366](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1366))
- improve(drds_instance): Added precheck to its testcases ([#1364](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1364))
- Improve(oss_bucket): Improve its resource use common method ([#1353](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1353))
- improve(launch_template): support update method ([#1327](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1327))
- improve(snapshot): support setting timeouts ([#1304](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1304))
- improve(instance):update testcase ([#1199](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1199))

BUG FIXES:

- fix(instnace): fix missing dry_run when creating instance ([#1401](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1401))
- fix(oss_bucket): fix oss bucket deleting timeout error ([#1400](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1400))
- fix(route_entry):fix route_entry create bug ([#1398](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1398))
- fix(instance):fix testcase name too length bug ([#1396](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1396))
- fix(vswitch):fix vswitch describe method wrapErrorf bug ([#1392](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1392))
- fix(slb_rule): fix testcase bug ([#1390](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1390))
- fix(db_backup_policy): pg10 of category 'basic' modify log_backup error ([#1388](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1388))
- fix(cen):Add deadline to cen datasources and modify timeout for DescribeCenBandwidthPackages ([#1381](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1381))
- fix(kvstore): kvstore_instance PostPaid to PrePaid error ([#1375](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1375))
- fix(cen): fixed its not display error message, added CenThrottlingUser retry ([#1373](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1373))

## 1.50.0 (July 01, 2019)

IMPROVEMENTS:

- Remove cs kubernetes autovpc testcases ([#1368](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1368))
- disable nav-visible in the alicloud.erb file ([#1365](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1365))
- Improve sweeper test and remove some needless waiting ([#1361](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1361))
- This is a Terraform 0.12 compatible release of this provider ([#1356](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1356))
- Deprecated resource `alicloud_cms_alarm` parameter start_time, end_time and removed notify_type based on the latest go sdk ([#1356](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1356))
- Adapt to new parameters of dedicated kubernetes cluster ([#1354](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1354))

BUG FIXES:

- Fix alicloud_cas_certificate setId bug ([#1368](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1368))
- Fix oss bucket datasource testcase based on the 0.12 syntax ([#1362](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1362))
- Fix deleting mongodb instance "NotFound" bug ([#1359](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1359))

## 1.49.0 (June 28, 2019)

- **New Date Source:** `alicloud_kvstore_instance_classes` ([#1315](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1315))

IMPROVEMENTS:

- remove the skipped testcase ([#1349](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1349))
- Move some import testcase into resource testcase ([#1348](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1348))
- Support attach & detach operation for loadbalancers and dbinstances ([#1346](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1346))
- update security_group_rule md ([#1345](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1345))
- Improve mongodb,rds testcase ([#1339](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1339))
- Deprecate field statement and use field document to replace ([#1338](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1338))
- Add function BuildStateConf for common timeouts setting ([#1330](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1330))
- drds_instance resource supports timeouts setting ([#1329](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1329))
- add support get Ak from config file ([#1328](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1328))
- Improve api_gateway_vpc use common method. ([#1323](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1323))
- Organize official documents in alphabetical order ([#1322](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1322))
- improve snapshot_policy testcase ([#1313](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1313))
- Improve api_gateway_group use common method ([#1311](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1311))
- Improve api_gateway_app use common method ([#1306](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1306))

BUG FIXES:

- bugfix: modify ess loadbalancers batch size ([#1352](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1352))
- fix instance OperationConflict bug ([#1351](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1351))
- fix(nas): convert some retrable error to nonretryable ([#1344](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1344))
- fix mongodb testcase ([#1341](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1341))
- fix log_store fields cannot be changed ([#1337](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1337))
- fix(nas): fix error handling ([#1336](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1336))
- fix db_instance_classes,db_instance_engines ([#1331](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1331))
- fix sls-logconfig config_name field to name ([#1326](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1326))
- fix db_instance_engines testcase ([#1325](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1325))
- fix forward_entries testcase bug ([#1324](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1324))

## 1.48.0 (June 21, 2019)

- **New Resource:** `alicloud_gpdb_connection` ([#1290](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1290))

IMPROVEMENTS:

- Improve rds testcase zone_id ([#1321](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1321))
- feature: support enable/disable action for resource alicloud_ess_alarm ([#1320](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1320))
- cen_instance resource supports timeouts setting ([#1318](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1318))
- added importer support for security_group_rule ([#1317](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1317))
- add multi_zone for db_instance_classes and db_instance_engines ([#1310](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1310))
- Update Eip Resource Isp Field ([#1303](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1303))
- Improve db_instance,db_read_write_splitting_connection,db_readonly_instance testcase ([#1300](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1300))
- Improve api_gateway_api use common method ([#1299](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1299))
- Add name for cen bandwidth package testcase ([#1298](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1298))
- Improve db testcase ([#1294](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1294))
- elasticsearch_instance resource supports timeouts setting ([#1268](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1268))

BUG FIXES:

- bugfix: remove the 'ForceNew' attribute of 'vswitch_ids' from resource alicloud_ess_scaling_group ([#1316](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1316))
- managed k8s no longer returns vswitchids and instancetypes, fix crash ([#1314](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1314))
- fix db_instance_classes ([#1309](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1309))
- fix oss lifecycle nil pointer bug ([#1307](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1307))
- Fix cen_bandwidth_limit Throttling.User bug ([#1305](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1305))
- fix disk_attachment test bug ([#1302](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1302))

## 1.47.0 (June 17, 2019)

- **New Date Source:** `alicloud_gpdb_instances` ([#1279](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1279))
- **New Resource:** `alicloud_gpdb_instance` ([#1260](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1260))

IMPROVEMENTS:

- fc_trigger datasource support outputting ids and names ([#1286](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1286))
- add fc_trigger support cdn_events ([#1285](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1285))
- modify apigateway-fc example ([#1284](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1284))
- Added PGP encrypt Support for ram access key ([#1280](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1280))
- Update Eip Resource Isp Field ([#1275](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1275))
- Improve fc_service use common method ([#1269](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1269))
- Improve fc_function use common method ([#1266](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1266))
- update dns_group testcase name ([#1265](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1265))
- update slb sdk ([#1263](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1263))
- improve vpn_connection testcase ([#1257](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1257))
- Improve cen_route_entries use common method ([#1249](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1249))
- Improve cen_bandwidth_package_attachment resource use common method ([#1240](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1240))
- Improve cen_bandwidth_package resource use common method ([#1237](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1237))

BUG FIXES:

- feat(nas): fix error report ([#1293](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1293))
- temp fix no value returned by cs openapi ([#1289](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1289))
- fix disk device_name bug ([#1288](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1288))
- fix sql server instance storage set bug ([#1283](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1283))
- fix db_instance_classes storage_range bug ([#1282](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1282))
- fc_service datasource support outputting ids and names ([#1278](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1278))
- fix log_store ListShards InternalServerError bug ([#1277](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1277))
- fix slb_listener docs bug ([#1276](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1276))
- fix clientToken bug ([#1272](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1272))
- fix(nas): fix document and nas_access_rules ([#1271](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1271))
- docs(version) Added 6.7 supported and fixed bug of version difference ([#1270](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1270))
- fix(nas): fix documents ([#1267](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1267))
- fix(nas): describe mount target & access rule [1264]

## 1.46.0 (June 10, 2019)

- **New Resource:** `alicloud_ram_account_password_policy` ([#1212](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1212))
- **New Date Source:** `alicloud_db_instance_engines` ([#1201](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1201))
- **New Date Source:** `alicloud_db_instance_classes` ([#1201](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1201))

IMPROVEMENTS:

- refactor(nas): move import to resource ([#1254](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1254))
- Improve ess_scalingconfiguration use common method ([#1250](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1250))
- improve ssl_vpn_client_cert testcase ([#1248](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1248))
- Improve ram_account_password_policy resource use common method ([#1247](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1247))
- add pending status for resource instance when creating ([#1245](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1245))
- resource instance supports timeouts configure ([#1244](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1244))
- added webhook support for alarms ([#1243](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1243))
- improve common test method ([#1242](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1242))
- Update Eip Association Resource ([#1238](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1238))
- improve ssl_vpn_server testcase ([#1235](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1235))
- Improve ess_scalingconfigurations datasource use common method ([#1234](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1234))
- improve vpn_customer_gateway testcase ([#1232](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1232))
- Improve cen_instance_grant use common method ([#1230](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1230))
- improve vpn_gateway testcase ([#1229](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1229))
- Improve cen_bandwidth_limit use common method ([#1227](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1227))
- Feature/support multi instance types ([#1226](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1226))
- Improve ess_attachment use common method ([#1225](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1225))
- Improve ess_alarm use common method ([#1218](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1218))
- Add support for assume_role in provider block ([#1217](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1217))
- Improve cen_instance_attachment resource use common method. ([#1216](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1216))
- add db instance engines and db instance classes data source support ([#1201](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1201))
- Handle alicloud_cs_*_kubernetes resource NotFound error properly ([#1191](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1191))

BUG FIXES:

- fix slb_attachment classic testcase ([#1259](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1259))
- fix oss bucket update bug ([#1258](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1258))
- fix scalingConfiguration is inconsistent with the information that is returned by describe, when the input parameter user_data is base64 ([#1256](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1256))
- fix slb_attachment err ObtainIpFail ([#1253](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1253))
- Fix password to comliant with the default password policy ([#1241](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1241))
- fix cr repo details, improve cs and cr docs ([#1239](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1239))
- fix(nas): fix unittest bugs ([#1236](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1236))
- fix slb_ca_certificate err ServiceIsConfiguring ([#1233](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1233))
- fix reset account_password don't work ([#1231](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1231))
- fix(nas): fix testcase errors ([#1184](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1184))

## 1.45.0 (May 29, 2019)

FEATURES:

- **New Resource:** `alicloud_network_acl_entries` ([#1208](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1208))

IMPROVEMENTS:

- update changeLog ([#1224](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1224))
- support oss object versioning ([#1121](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1121))
- update instance dataSource doc ([#1215](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1215))
- update oss buket encryption configuration ([#1214](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1214))
- support oss bucket tags ([#1213](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1213))
- support oss bucket encryption configuration ([#1210](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1210))
- Improve cen_instances use common method ([#1206](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1206))
- support set oss bucket stroage class ([#1204](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1204))
- Improve ess_lifecyclehook resource use common method ([#1196](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1196))
- Improve ess_scalinggroup use common method ([#1192](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1192))
- Improve ess_scheduled_task resource use common method ([#1175](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1175))
- improve route_table testcase ([#1109](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1109))

BUG FIXES:

- fix nat_gateway and network_interface testcase bug ([#1211](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1211))
- Fix ram testcases name length bug ([#1205](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1205))
- fix actiontrail bug ([#1203](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1203))

## 1.44.0 (May 24, 2019)

FEATURES:

- **New Resource:** `alicloud_network_acl_attachment` ([#1187](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1187))

IMPROVEMENTS:

- update CHANGELOG.md ([#1209](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1209))
- Skip instance some testcases to avoid qouta limit ([#1195](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1195))
- Added the multi zone's instance supported ([#1194](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1194))
- remove multi test of ram_account_alias ([#1186](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1186))
- Improve ram_role_attachment resource use common method ([#1185](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1185))
- Improve ess_scalingrule use common method ([#1183](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1183))
- update mongodb instance resource document ([#1182](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1182))
- Improve ram_role resource use common method ([#1181](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1181))
- Correct the oss bucket docs ([#1178](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1178))
- add slb classic not support regions ([#1176](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1176))
- Dev versioning ([#1174](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1174))
- Improve ram_user_policy_attachment resource use common method ([#1172](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1172))
- Improve ram_role_policy_attachment resource use common method ([#1171](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1171))
- improve router_interface testcase ([#1170](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1170))
- Improve ram_policy resource use common method ([#1166](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1166))
- Improve slb_listeners datasource use common method ([#1165](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1165))
- add name attribute for forward_entry ([#1164](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1164))
- Improve ram_group_policy_attachment resource use common method ([#1163](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1163))
- Improve ram_group_membership resource use common method ([#1159](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1159))
- Improve ram_login_profile resource use common method ([#1158](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1158))
- Improve ram_group resource use common method ([#1150](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1150))

BUG FIXES:

- Fix ram_user sweeper ([#1200](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1200))
- Fix ram group import bug ([#1198](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1198))
- fix router_interface dataSource testcase bug ([#1197](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1197))
- fix forward_entry multi testcase bug ([#1189](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1189))
- fix api gw and network acl sweeper test error ([#1180](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1180))
- fix ram user diff bug ([#1179](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1179))
- Fix ram account alias multi testcase bug ([#1169](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1169))

## 1.43.0 (May 17, 2019)

FEATURES:

- **New Resource:** `alicloud_network_acl` (([#1151](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1151))

IMPROVEMENTS:

- change ecs instance instance_charge_type modifying position ([#1168](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1168))
- AutoScaling support multiple security groups ([#1167](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1167))
- Update ots and vpc document ([#1162](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1162))
- Improve some slb datasource ([#1155](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1155))
- improve forward_entry testcase ([#1152](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1152))
- improve slb_attachment resource use common method ([#1148](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1148))
- Improve ram_account_alias resource use common method  ([#1147](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1147))
- slb instance support updating specification ([#1145](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1145))
- improve slb_server_group resource use common method ([#1144](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1144))
- add note for SLB that intl account does not support creating PrePaid instance ([#1143](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1143))
- Update ots document ([#1142](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1142))
- improve slb_server_certificate resource use common method ([#1139](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1139))

BUG FIXES:

- Fix ram account alias notfound bug ([#1161](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1161))
- fix(nas): refactor testcases ([#1157](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1157))

## 1.42.0 (May 10, 2019)

FEATURES:

- **New Resource:** `alicloud_snapshot_policy` ([#989](https://github.com/terraform-providers/terraform-provider-alicloud/issues/989))

IMPROVEMENTS:

- improve mongodb and db sweeper test ([#1138](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1138))
- Alicloud_ots_table: add max version offset ([#1137](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1137))
- update disk category ([#1135](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1135))
- Update Route Entry Resource ([#1134](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1134))
- update images testcase check condition ([#1133](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1133))
- bugfix: ess alarm apply recreate ([#1131](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1131))
- improve slb_listener resource use common method ([#1130](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1130))
- mongodb sharding instance add backup policy support ([#1127](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1127))
- Improve ram_users datasource use common method ([#1126](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1126))
- Improve ram_policies datasource use common method ([#1125](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1125))
- rds datasource test case remove connection mode check ([#1124](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1124))
- Add missing bracket ([#1123](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1123))
- add support sha256 ([#1122](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1122))
- Improve ram_groups datasource use common method ([#1121](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1121))
- Modified the sweep rules in ram_roles testcases ([#1116](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1116))
- improve instance testcase ([#1114](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1114))
- Improve slb_ca_certificate resource use common method ([#1113](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1113))
- Improve ram_roles datasource use common method ([#1112](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1112))
- Improve slb datasource use common method ([#1111](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1111))
- Improve ram_account_alias use common method ([#1108](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1108))
- update data_source_alicoud_mongo_instances and add test case ([#1107](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1107))
- add mongodb backup policy support, test case, document ([#1106](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1106))
- update route_entry and forward_entry document ([#1096](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1096))
- Improve slb_acl resource use common method ([#1092](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1092))
- improve snat_entry testcase ([#1091](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1091))
- Improve slb resource use common method ([#1090](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1090))
- improve nat_gateway testcase ([#1089](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1089))
- Modify table to entry ([#1088](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1088))
- Modified the error code returned when timeout of upgrading instance ([#1085](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1085))
- improve db backup policy test case ([#1083](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1083))

BUG FIXES:

- Fix scalinggroup id is not found before creating scaling configuration  ([#1119](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1119))
- fix slb instance sets tags bug ([#1105](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1105))
- fix not support outputfile ([#1095](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1095))
- Bugfix/slb import server group ([#1093](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1093))
- Fix fc_triggers datasource when type is mns_topic ([#1086](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1086))

## 1.41.0 (April 29, 2019)

IMPROVEMENTS:

- Improve fc_trigger support mns_topic modify config ([#1082](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1082))
- Rds sdk-update ([#1078](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1078))
- update some eip method name ([#1077](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1077))
- improve vswitch testcase  ([#1076](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1076))
- add rand for db_instances testcase ([#1074](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1074))
- Improve fc_trigger support mns_topic ([#1073](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1073))
- remove zone_id setting in the db instance testcase ([#1069](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1069))
- change database default zone id to avoid some unsupported cases ([#1067](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1067))
- add oss bucket policy implementation ([#1066](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1066))
- improve vpc testcase ([#1065](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1065))
- Change password to Yourpassword ([#1063](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1063))
- Improve kvstore_instance datasource use common method ([#1062](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1062))
- improve eip testcase ([#1058](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1058))
- Improve kvstore_instance testcase use common method ([#1052](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1052))
- improve mongodb testcase ([#1050](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1050))
- update network_interface dataSource basic testcase config ([#1049](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1049))
- Improve kvstore_backup_policy testcase use common method ([#1044](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1044))

BUG FIXES:

- Fix fc_triggers datasource when type is mns_topic ([#1086](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1086))
- Fix kvstore_instance multi ([#1080](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1080))
- fix eip_association bug when snat or forward be released ([#1075](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1075))
- Fix db_readonly_instance instance_name ([#1071](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1071))
- fixed DB log backup policy bug when the log_retention_period does not input ([#1056](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1056))
- fix cms diff bug and improve its testcases ([#1057](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1057))


## 1.40.0 (April 20, 2019)

FEATURES:

- **New Resource:** `alicloud_mongodb_sharding_instance` ([#1017](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1017))
- **New Data Source:** `alicloud_snapshots` ([#988](https://github.com/terraform-providers/terraform-provider-alicloud/issues/988))
- **New Resource:** `alicloud_snapshot` ([#954](https://github.com/terraform-providers/terraform-provider-alicloud/issues/954))

IMPROVEMENTS:

- Fix db_instance can't find method DescribeDbInstance ([#1046](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1046))
- update network_interface testcase config ([#1045](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1045))
- Update Nat Gateway Resource ([#1043](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1043))
- improve network_interface dataSource testcase ([#1042](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1042))
- improve network_interface resource testcase ([#1041](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1041))
- Improve db_database db_instance db_readonly_instance db_readwrite_splitting_connection ([#1040](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1040))
- improve key_pair resource testcase ([#1039](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1039))
- improve key_pair dataSource testcase ([#1038](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1038))
- make fmt ess_scalinggroups ([#1036](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1036))
- improve test common method ([#1030](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1030))
- Update cen data source document ([#1029](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1029))
- fix Error method [[#1024](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1024)] 
- Update Nat Gateway Token ([#1020](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1020))
- update RAM website document ([#1019](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1019))
- add computed for resource_group_id ([#1018](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1018))
- remove ram validators and update website docs ([#1016](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1016))
- improve test common method, support 'TestMatchResourceAttr' check ([#1012](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1012))
- resource group support for creating new VPC ([#1010](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1010))
- Improve cs_cluster sweeper test removing retained resources ([#1002](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1002))
- improve security_group testcase use common method ([#995](https://github.com/terraform-providers/terraform-provider-alicloud/issues/995))
- fix vpn change local_subnet and remote_subnet bug ([#994](https://github.com/terraform-providers/terraform-provider-alicloud/issues/994))
- improve disk dataSource testcase use common method ([#990](https://github.com/terraform-providers/terraform-provider-alicloud/issues/990))
- fix(nas): use new sdk ([#984](https://github.com/terraform-providers/terraform-provider-alicloud/issues/984))
- Feature/slb listener redirect http to https ([#981](https://github.com/terraform-providers/terraform-provider-alicloud/issues/981))
- improve disk and diskAttachment resource testcase use testCommon method ([#978](https://github.com/terraform-providers/terraform-provider-alicloud/issues/978))
- improve dns dataSource testcase use testCommon method ([#971](https://github.com/terraform-providers/terraform-provider-alicloud/issues/971))

BUG FIXES:

- Fix ess go sdk compatibility ([#1032](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1032))
- Update sdk to fix timeout bug ([#1015](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1015))
- Fix Eip And VSwitch ClientToken bug ([#1000](https://github.com/terraform-providers/terraform-provider-alicloud/issues/1000))
- fix db_account diff bug and add some notes for it ([#999](https://github.com/terraform-providers/terraform-provider-alicloud/issues/999))
- fix vpn gateway Period bug ([#993](https://github.com/terraform-providers/terraform-provider-alicloud/issues/993))


## 1.39.0 (April 09, 2019)

FEATURES:

- **New Data Source:** `alicloud_ots_instance_attachments` ([#986](https://github.com/terraform-providers/terraform-provider-alicloud/issues/986))
- **New Data Source:** `alicloud_ssl_vpc_servers` ([#985](https://github.com/terraform-providers/terraform-provider-alicloud/issues/985))
- **New Data Source:** `alicloud_ssl_vpn_client_certs` ([#986](https://github.com/terraform-providers/terraform-provider-alicloud/issues/986))
- **New Data Source:** `alicloud_ess_scaling_rules` ([#976](https://github.com/terraform-providers/terraform-provider-alicloud/issues/976))
- **New Data Source:** `alicloud_ess_scaling_configurations` ([#974](https://github.com/terraform-providers/terraform-provider-alicloud/issues/974))
- **New Data Source:** `alicloud_ess_scaling_groups` ([#973](https://github.com/terraform-providers/terraform-provider-alicloud/issues/973))
- **New Data Source:** `alicloud_ddoscoo_instances` ([#967](https://github.com/terraform-providers/terraform-provider-alicloud/issues/967))
- **New Data Source:** `alicloud_ots_instances` ([#946](https://github.com/terraform-providers/terraform-provider-alicloud/issues/946))

IMPROVEMENTS:

- Improve instance type updating testcase ([#979](https://github.com/terraform-providers/terraform-provider-alicloud/issues/979))
- support changing prepaid instance type ([#977](https://github.com/terraform-providers/terraform-provider-alicloud/issues/977))
- Improve db_account db_account_privilege db_backup_policy db_connection ([#963](https://github.com/terraform-providers/terraform-provider-alicloud/issues/963))

BUG FIXES:

- Fix Nat GW ClientToken bug ([#983](https://github.com/terraform-providers/terraform-provider-alicloud/issues/983))
- Fix print error bug after DescribeDBInstanceById ([#980](https://github.com/terraform-providers/terraform-provider-alicloud/issues/980))

## 1.38.0 (April 03, 2019)

FEATURES:

- **New Resource:** `alicloud_ddoscoo_instance` ([#952](https://github.com/terraform-providers/terraform-provider-alicloud/issues/952))

IMPROVEMENTS:

- update dns_group describe method ([#966](https://github.com/terraform-providers/terraform-provider-alicloud/issues/966))
- update ram_policy resource testcase ([#964](https://github.com/terraform-providers/terraform-provider-alicloud/issues/964))
- improve ram_policy resource update method ([#960](https://github.com/terraform-providers/terraform-provider-alicloud/issues/960))
- ecs prepaid instance supports changing instance type ([#949](https://github.com/terraform-providers/terraform-provider-alicloud/issues/949))
- update mongodb instance test case for multiAZ ([#947](https://github.com/terraform-providers/terraform-provider-alicloud/issues/947))
- add test common method ,improve dns resource testcase ([#927](https://github.com/terraform-providers/terraform-provider-alicloud/issues/927))


BUG FIXES:

- Fix drds instance sweeper test bug ([#955](https://github.com/terraform-providers/terraform-provider-alicloud/issues/955))

## 1.37.0 (March 29, 2019)

FEATURES:

- **New Resource:** `alicloud_mongodb_instance` ([#881](https://github.com/terraform-providers/terraform-provider-alicloud/issues/881))
- **New Resource:** `alicloud_cen_instance_grant` ([#857](https://github.com/terraform-providers/terraform-provider-alicloud/issues/857))
- **New Data Source:** `alicloud_forward_entries` ([#922](https://github.com/terraform-providers/terraform-provider-alicloud/issues/922))
- **New Data Source:** `alicloud_snat_entries` ([#920](https://github.com/terraform-providers/terraform-provider-alicloud/issues/920))
- **New Data Source:** `alicloud_nat_gateways` ([#918](https://github.com/terraform-providers/terraform-provider-alicloud/issues/918))
- **New Data Source:** `alicloud_route_entries` ([#915](https://github.com/terraform-providers/terraform-provider-alicloud/issues/915))

IMPROVEMENTS:

- Add missing outputs for datasource dns_records, security groups, vpcs and vswitches ([#943](https://github.com/terraform-providers/terraform-provider-alicloud/issues/943))
- datasource dns_records add a output urls ([#942](https://github.com/terraform-providers/terraform-provider-alicloud/issues/942))
- modify stop instance timeout to 5min to avoid the exception timeout ([#941](https://github.com/terraform-providers/terraform-provider-alicloud/issues/941))
- datasource security_groups, vpcs and vswitches support outputs ids and names ([#939](https://github.com/terraform-providers/terraform-provider-alicloud/issues/939))
- Improve all of parameter's tag, like 'Required', 'ForceNew' ([#938](https://github.com/terraform-providers/terraform-provider-alicloud/issues/938))
- Improve pvtz_zone_record WrapError ([#934](https://github.com/terraform-providers/terraform-provider-alicloud/issues/934))
- Improve pvtz_zone_record create record ([#933](https://github.com/terraform-providers/terraform-provider-alicloud/issues/933))
- testSweepCRNamespace skip not supported region  ([#932](https://github.com/terraform-providers/terraform-provider-alicloud/issues/932))
- refine retry logic of resource tablestore to avoid the exception timeout ([#931](https://github.com/terraform-providers/terraform-provider-alicloud/issues/931))
- Improve pvtz resource datasource testcases ([#928](https://github.com/terraform-providers/terraform-provider-alicloud/issues/928))
- cr_repos fix docs link error ([#926](https://github.com/terraform-providers/terraform-provider-alicloud/issues/926))
- resource DB instance supports setting security group ([#925](https://github.com/terraform-providers/terraform-provider-alicloud/issues/925))
- resource DB instance supports setting monitor period ([#924](https://github.com/terraform-providers/terraform-provider-alicloud/issues/924))
- Skipping bandwidth package related test for international site account ([#917](https://github.com/terraform-providers/terraform-provider-alicloud/issues/917))
- Resource snat entry update id and support import ([#916](https://github.com/terraform-providers/terraform-provider-alicloud/issues/916))
- add docs about prerequisites for cs and cr  ([#914](https://github.com/terraform-providers/terraform-provider-alicloud/issues/914))
- add new schema environment_variables to fc_function.html.markdown ([#913](https://github.com/terraform-providers/terraform-provider-alicloud/issues/913))
- add skipping check for datasource route tables' testcases ([#911](https://github.com/terraform-providers/terraform-provider-alicloud/issues/911))
- modify ram_user id by userId ([#900](https://github.com/terraform-providers/terraform-provider-alicloud/issues/900))

BUG FIXES:

- Deprecate bucket `logging_isenable` and fix referer_config diff bug ([#937](https://github.com/terraform-providers/terraform-provider-alicloud/issues/937))
- fix ram user and group sweeper test bug ([#929](https://github.com/terraform-providers/terraform-provider-alicloud/issues/929))
- Fix the parameter bug when actiontrail is created ([#921](https://github.com/terraform-providers/terraform-provider-alicloud/issues/921))
- fix default pod_cidr in k8s docs ([#919](https://github.com/terraform-providers/terraform-provider-alicloud/issues/919))

## 1.36.0 (March 24, 2019)

FEATURES:

- **New Resource:** `alicloud_cas_certificate` ([#875](https://github.com/terraform-providers/terraform-provider-alicloud/issues/875))
- **New Data Source:** `alicloud_route_tables` ([#905](https://github.com/terraform-providers/terraform-provider-alicloud/issues/905))
- **New Data Source:** `alicloud_common_bandwidth_packages` ([#897](https://github.com/terraform-providers/terraform-provider-alicloud/issues/897))
- **New Data Source:** `alicloud_actiontrails` ([#891](https://github.com/terraform-providers/terraform-provider-alicloud/issues/891))
- **New Data Source:** `alicloud_cas_certificates` ([#875](https://github.com/terraform-providers/terraform-provider-alicloud/issues/875))

IMPROVEMENTS:

- Add wait method for disk and disk attachment ([#910](https://github.com/terraform-providers/terraform-provider-alicloud/issues/910))
- Add wait method for cen instance ([#909](https://github.com/terraform-providers/terraform-provider-alicloud/issues/909))
- add dns and dns_group test sweeper ([#906](https://github.com/terraform-providers/terraform-provider-alicloud/issues/906))
- fc_function add new schema environment_variables ([#904](https://github.com/terraform-providers/terraform-provider-alicloud/issues/904))
- support kv-store auto renewal option  documentation ([#902](https://github.com/terraform-providers/terraform-provider-alicloud/issues/902))
- Sort slb slave zone ids to avoid needless error ([#898](https://github.com/terraform-providers/terraform-provider-alicloud/issues/898))
- add region skip for container registry testcase ([#896](https://github.com/terraform-providers/terraform-provider-alicloud/issues/896))
- Add `enable_details` for alicloud_zones and support retrieving slb slave zones ([#893](https://github.com/terraform-providers/terraform-provider-alicloud/issues/893))
- Slb support setting master and slave zone id ([#887](https://github.com/terraform-providers/terraform-provider-alicloud/issues/887))
- improve disk and attachment resource testcase ([#886](https://github.com/terraform-providers/terraform-provider-alicloud/issues/886))
- Remove ModifySecurityGroupPolicy waiting and backend has fixed it ([#883](https://github.com/terraform-providers/terraform-provider-alicloud/issues/883))
- Improve cas resource and datasource testcases ([#882](https://github.com/terraform-providers/terraform-provider-alicloud/issues/882))
- Make db_connection resource code more standard ([#879](https://github.com/terraform-providers/terraform-provider-alicloud/issues/879))

BUG FIXES:

- Fix cen instance deleting bug ([#908](https://github.com/terraform-providers/terraform-provider-alicloud/issues/908))
- Fix cen create bug when one resion is China ([#903](https://github.com/terraform-providers/terraform-provider-alicloud/issues/903))
- fix cas_certificate sweeper test bug ([#899](https://github.com/terraform-providers/terraform-provider-alicloud/issues/899))
- Modify ram group's name's ForceNew to true ([#895](https://github.com/terraform-providers/terraform-provider-alicloud/issues/895))
- fix mount target deletion bugs ([#892](https://github.com/terraform-providers/terraform-provider-alicloud/issues/892))
- Fix link to BatchSetCdnDomainConfig document  documentation ([#885](https://github.com/terraform-providers/terraform-provider-alicloud/issues/885))
- fix rds instance parameter test case issue ([#880](https://github.com/terraform-providers/terraform-provider-alicloud/issues/880))

## 1.35.0 (March 18, 2019)

FEATURES:

- **New Resource:** `alicloud_cr_repo` ([#862](https://github.com/terraform-providers/terraform-provider-alicloud/issues/862))
- **New Resource:** `alicloud_actiontrail` ([#858](https://github.com/terraform-providers/terraform-provider-alicloud/issues/858))
- **New Data Source:** `alicloud_cr_repos` ([#868](https://github.com/terraform-providers/terraform-provider-alicloud/issues/868))
- **New Data Source:** `alicloud_cr_namespaces` ([#867](https://github.com/terraform-providers/terraform-provider-alicloud/issues/867))
- **New Data Source:** `alicloud_nas_file_systems` ([#864](https://github.com/terraform-providers/terraform-provider-alicloud/issues/864))
- **New Data Source:** `alicloud_nas_mount_targets` ([#864](https://github.com/terraform-providers/terraform-provider-alicloud/issues/864))
- **New Data Source:** `alicloud_drds_instances` ([#861](https://github.com/terraform-providers/terraform-provider-alicloud/issues/861))
- **New Data Source:** `alicloud_nas_access_rules` ([#860](https://github.com/terraform-providers/terraform-provider-alicloud/issues/860))
- **New Data Source:** `alicloud_nas_access_groups` ([#856](https://github.com/terraform-providers/terraform-provider-alicloud/issues/856))

IMPROVEMENTS:

- Improve actiontrail docs ([#878](https://github.com/terraform-providers/terraform-provider-alicloud/issues/878))
- Add account pre-check for common bandwidth package to avoid known error ([#877](https://github.com/terraform-providers/terraform-provider-alicloud/issues/877))
- Make dns resource code more standard ([#876](https://github.com/terraform-providers/terraform-provider-alicloud/issues/876))
- Improve dns resources' testcases ([#859](https://github.com/terraform-providers/terraform-provider-alicloud/issues/859))
- Add client token for vpn services ([#855](https://github.com/terraform-providers/terraform-provider-alicloud/issues/855))
- reback the lossing datasource ([#866](https://github.com/terraform-providers/terraform-provider-alicloud/issues/866))
- Improve drds instances testcases  documentation ([#863](https://github.com/terraform-providers/terraform-provider-alicloud/issues/863))
- Update sdk for vpc package ([#854](https://github.com/terraform-providers/terraform-provider-alicloud/issues/854))

BUG FIXES:

- Add waiting method to ensure the security group status is ok ([#873](https://github.com/terraform-providers/terraform-provider-alicloud/issues/873))
- Fix nas mount target notfound bug and improve nas datasource's testcases ([#872](https://github.com/terraform-providers/terraform-provider-alicloud/issues/872))
- Fix dns notfound bug ([#871](https://github.com/terraform-providers/terraform-provider-alicloud/issues/871))
- fix creating slb bug ([#870](https://github.com/terraform-providers/terraform-provider-alicloud/issues/870))
- fix elastic search sweeper test bug ([#865](https://github.com/terraform-providers/terraform-provider-alicloud/issues/865))


## 1.34.0 (March 13, 2019)

FEATURES:

- **New Resource:** `alicloud_nas_mount_target` ([#835](https://github.com/terraform-providers/terraform-provider-alicloud/issues/835))
- **New Resource:** `alicloud_cdn_domain_config` ([#829](https://github.com/terraform-providers/terraform-provider-alicloud/issues/829))
- **New Resource:** `alicloud_cr_namespace` ([#827](https://github.com/terraform-providers/terraform-provider-alicloud/issues/827))
- **New Resource:** `alicloud_nas_access_rule` ([#827](https://github.com/terraform-providers/terraform-provider-alicloud/issues/827))
- **New Resource:** `alicloud_cdn_domain_new` ([#787](https://github.com/terraform-providers/terraform-provider-alicloud/issues/787))
- **New Data Source:** `alicloud_cs_kubernetes_clusters` ([#818](https://github.com/terraform-providers/terraform-provider-alicloud/issues/818))

IMPROVEMENTS:

- Add drds instance docs ([#853](https://github.com/terraform-providers/terraform-provider-alicloud/issues/853))
- Improve resource mount target testcases ([#852](https://github.com/terraform-providers/terraform-provider-alicloud/issues/852))
- Add using note for spot instance ([#851](https://github.com/terraform-providers/terraform-provider-alicloud/issues/851))
- Resource alicloud_slb supports PrePaid ([#850](https://github.com/terraform-providers/terraform-provider-alicloud/issues/850))
- Add ssl_vpn_server and ssl_vpn_client_cert sweeper test ([#843](https://github.com/terraform-providers/terraform-provider-alicloud/issues/843))
- Improve vpn_gateway testcases and some sweeper test ([#842](https://github.com/terraform-providers/terraform-provider-alicloud/issues/842))
- Improve dns datasource testcases ([#841](https://github.com/terraform-providers/terraform-provider-alicloud/issues/841))
- Improve Eip and mns testcase ([#840](https://github.com/terraform-providers/terraform-provider-alicloud/issues/840))
- Add version notes in some docs ([#838](https://github.com/terraform-providers/terraform-provider-alicloud/issues/838))
- RDS resource supports auto-renewal ([#836](https://github.com/terraform-providers/terraform-provider-alicloud/issues/836))
- Deprecate the resource alicloud_cdn_domain ([#830](https://github.com/terraform-providers/terraform-provider-alicloud/issues/830))

BUG FIXES:

- Fix deleting dns record InternalError bug ([#848](https://github.com/terraform-providers/terraform-provider-alicloud/issues/848))
- fix log store and config sweeper test deleting bug ([#847](https://github.com/terraform-providers/terraform-provider-alicloud/issues/847))
- Fix drds resource no supporting client token ([#846](https://github.com/terraform-providers/terraform-provider-alicloud/issues/846))
- fix kms sweeper test deleting bug ([#844](https://github.com/terraform-providers/terraform-provider-alicloud/issues/844))
- fix kubernetes data resource ut and import error ([#839](https://github.com/terraform-providers/terraform-provider-alicloud/issues/839))
- Bugfix: destroying alicloud_ess_attachment timeout ([#834](https://github.com/terraform-providers/terraform-provider-alicloud/issues/834))
- fix cdn service func WaitForCdnDomain ([#833](https://github.com/terraform-providers/terraform-provider-alicloud/issues/833))
- deal with the error message in cen route entry ([#831](https://github.com/terraform-providers/terraform-provider-alicloud/issues/831))
- change bool to *bool in parameters of k8s clusters ([#828](https://github.com/terraform-providers/terraform-provider-alicloud/issues/828))
- Fix nas docs bug ([#825](https://github.com/terraform-providers/terraform-provider-alicloud/issues/825))
- create vpn gateway got "UnnecessarySslConnection" error when enable_ssl is false ([#822](https://github.com/terraform-providers/terraform-provider-alicloud/issues/822))

## 1.33.0 (March 05, 2019)

FEATURES:

- **New Resource:** `alicloud_nas_access_group` ([#817](https://github.com/terraform-providers/terraform-provider-alicloud/issues/817))
- **New Resource:** `alicloud_nas_file_system` ([#807](https://github.com/terraform-providers/terraform-provider-alicloud/issues/807))

IMPROVEMENTS:

- Improve nas resource docs ([#824](https://github.com/terraform-providers/terraform-provider-alicloud/issues/824))

BUG FIXES:

- bugfix: create vpn gateway got "UnnecessarySslConnection" error when enable_ssl is false ([#822](https://github.com/terraform-providers/terraform-provider-alicloud/issues/822))
- fix volume_tags diff bug when running testcases ([#816](https://github.com/terraform-providers/terraform-provider-alicloud/issues/816))

## 1.32.1 (March 03, 2019)

BUG FIXES:

- fix volume_tags diff bug when setting tags by alicloud_disk ([#815](https://github.com/terraform-providers/terraform-provider-alicloud/issues/815))

## 1.32.0 (March 01, 2019)

FEATURES:

- **New Resource:** `alicloud_db_readwrite_splitting_connection` ([#753](https://github.com/terraform-providers/terraform-provider-alicloud/issues/753))

IMPROVEMENTS:

- add slb_internet_enabled to managed kubernetes ([#806](https://github.com/terraform-providers/terraform-provider-alicloud/issues/806))
- update alicloud_slb_attachment usage example ([#805](https://github.com/terraform-providers/terraform-provider-alicloud/issues/805))
- rds support op tags  documentation ([#797](https://github.com/terraform-providers/terraform-provider-alicloud/issues/797))
- ForceNew for resource record and zone id updates for pvtz record ([#794](https://github.com/terraform-providers/terraform-provider-alicloud/issues/794))
- support volume tags for ecs instance disks ([#793](https://github.com/terraform-providers/terraform-provider-alicloud/issues/793))
- Improve instance and security group testcase for different account site ([#792](https://github.com/terraform-providers/terraform-provider-alicloud/issues/792))
- Add account site type setting to skip unsupported test cases automatically ([#790](https://github.com/terraform-providers/terraform-provider-alicloud/issues/790))
- update alibaba-cloud-sdk-go to use lastest useragent and modify errMessage when signature does not match  dependencies ([#788](https://github.com/terraform-providers/terraform-provider-alicloud/issues/788))
- make the timeout longer when cen attach/detach vpc ([#786](https://github.com/terraform-providers/terraform-provider-alicloud/issues/786))
- cen child instance attach after vsw created ([#785](https://github.com/terraform-providers/terraform-provider-alicloud/issues/785))
- kvstore support parameter configuration ([#784](https://github.com/terraform-providers/terraform-provider-alicloud/issues/784))
- Modify useragent to meet the standard of sdk ([#778](https://github.com/terraform-providers/terraform-provider-alicloud/issues/778))
- Modify kms client to dock with the alicloud official GO SDK ([#763](https://github.com/terraform-providers/terraform-provider-alicloud/issues/763))

BUG FIXES:

- fix rds readonly instance name update issue ([#812](https://github.com/terraform-providers/terraform-provider-alicloud/issues/812))
- fix import managed kubernetes test ([#809](https://github.com/terraform-providers/terraform-provider-alicloud/issues/809))
- fix rds parameter update issue ([#804](https://github.com/terraform-providers/terraform-provider-alicloud/issues/804))
- fix first create db with tags ([#803](https://github.com/terraform-providers/terraform-provider-alicloud/issues/803))
- Fix dns record ttl setting error and update bug ([#800](https://github.com/terraform-providers/terraform-provider-alicloud/issues/800))
- Fix vpc return custom route table bug ([#799](https://github.com/terraform-providers/terraform-provider-alicloud/issues/799))
- fix ssl vpn subnet can not pass comma separated string problem ([#780](https://github.com/terraform-providers/terraform-provider-alicloud/issues/780))
- fix(whitelist) Modified whitelist returned and filter the default values ([#779](https://github.com/terraform-providers/terraform-provider-alicloud/issues/779))

## 1.31.0 (February 19, 2019)

FEATURES:

- **New Resource:** `alicloud_db_readonly_instance` ([#755](https://github.com/terraform-providers/terraform-provider-alicloud/issues/755))

IMPROVEMENTS:

- support update deletion_protection option documentation ([#771](https://github.com/terraform-providers/terraform-provider-alicloud/issues/771))
- add three az k8s cluster docs  documentation ([#767](https://github.com/terraform-providers/terraform-provider-alicloud/issues/767))
- kvstore support vpc_auth_mode  dependencies ([#765](https://github.com/terraform-providers/terraform-provider-alicloud/issues/765))
- Fix sls logtail config collection error ([#762](https://github.com/terraform-providers/terraform-provider-alicloud/issues/762))
- Add attribute parameters to resource alicloud_db_instance  documentation ([#761](https://github.com/terraform-providers/terraform-provider-alicloud/issues/761))
- Add attribute parameters to resource alicloud_db_instance ([#761](https://github.com/terraform-providers/terraform-provider-alicloud/issues/761))
- Modify dns client to dock with the alicloud official GO SDK ([#750](https://github.com/terraform-providers/terraform-provider-alicloud/issues/750))

BUG FIXES:

- Fix cms_alarm updating notify_type bug ([#773](https://github.com/terraform-providers/terraform-provider-alicloud/issues/773))
- fix(error) Fixed bug of error code when timeout for upgrade instance ([#770](https://github.com/terraform-providers/terraform-provider-alicloud/issues/770))
- delete success if not found cen route when delete ([#753](https://github.com/terraform-providers/terraform-provider-alicloud/issues/753))

## 1.30.0 (February 04, 2019)

FEATURES:

- **New Resource:** `alicloud_elasticsearch_instance` ([#722](https://github.com/terraform-providers/terraform-provider-alicloud/issues/722))
- **New Resource:** `alicloud_logtail_attachment` ([#705](https://github.com/terraform-providers/terraform-provider-alicloud/issues/705))
- **New Data Source:** `alicloud_elasticsearch_instances` ([#739](https://github.com/terraform-providers/terraform-provider-alicloud/issues/739))

IMPROVEMENTS:

- Improve snat and forward testcases ([#749](https://github.com/terraform-providers/terraform-provider-alicloud/issues/749))
- delete data source roles limit of policy_type and policy_name ([#748](https://github.com/terraform-providers/terraform-provider-alicloud/issues/748))
- make k8s cluster deleting timeout longer ([#746](https://github.com/terraform-providers/terraform-provider-alicloud/issues/746))
- Improve nat_gateway testcases ([#743](https://github.com/terraform-providers/terraform-provider-alicloud/issues/743))
- Improve eip_association testcases ([#742](https://github.com/terraform-providers/terraform-provider-alicloud/issues/742))
- Improve elasticinstnace testcases for IPV6 supported ([#741](https://github.com/terraform-providers/terraform-provider-alicloud/issues/741))
- Add debug for db instance and ess group ([#740](https://github.com/terraform-providers/terraform-provider-alicloud/issues/740))
- Improve api_gateway_vpc_access testcases ([#738](https://github.com/terraform-providers/terraform-provider-alicloud/issues/738))
- Modify errors and  ram client to dock with the GO SDK ([#735](https://github.com/terraform-providers/terraform-provider-alicloud/issues/735))
- provider supports getting credential via ecs role name ([#731](https://github.com/terraform-providers/terraform-provider-alicloud/issues/731))
- Update testcases for cen region domain route entries ([#729](https://github.com/terraform-providers/terraform-provider-alicloud/issues/729))
- cs_kubernetes supports user_ca ([#726](https://github.com/terraform-providers/terraform-provider-alicloud/issues/726))
- Wrap resource elasticserarch_instance's error ([#725](https://github.com/terraform-providers/terraform-provider-alicloud/issues/725))
- Add note for kubernetes resource and improve its testcase ([#724](https://github.com/terraform-providers/terraform-provider-alicloud/issues/724))
- Datasource instance_types supports filter results and used to create kuberneters ([#723](https://github.com/terraform-providers/terraform-provider-alicloud/issues/723))
- Add ids parameter extraction in data source regions,zones,dns_domain,images and instance_types([#718](https://github.com/terraform-providers/terraform-provider-alicloud/issues/718))
- Improve dns group testcase ([#717](https://github.com/terraform-providers/terraform-provider-alicloud/issues/717))
- Improve security group rule testcase for classic ([#716](https://github.com/terraform-providers/terraform-provider-alicloud/issues/716))
- Improve security group creating request ([#715](https://github.com/terraform-providers/terraform-provider-alicloud/issues/715))
- Route entry supports Nat Gateway ([#713](https://github.com/terraform-providers/terraform-provider-alicloud/issues/713))
- Modify db account returning update to read after creating ([#711](https://github.com/terraform-providers/terraform-provider-alicloud/issues/711))
- Improve cdn testcase ([#708](https://github.com/terraform-providers/terraform-provider-alicloud/issues/708))
- Apply wraperror to security_group, security_group_rule, vswitch, disk ([#707](https://github.com/terraform-providers/terraform-provider-alicloud/issues/707))
- Improve cdn testcase ([#705](https://github.com/terraform-providers/terraform-provider-alicloud/issues/705))
- Add notes for datahub and improve its testcase ([#704](https://github.com/terraform-providers/terraform-provider-alicloud/issues/704))
- Improve security_group_rule resource and data source testcases ([#703](https://github.com/terraform-providers/terraform-provider-alicloud/issues/703))
- Improve kvstore backup policy ([#701](https://github.com/terraform-providers/terraform-provider-alicloud/issues/701))
- Improve pvtz attachment testcase ([#700](https://github.com/terraform-providers/terraform-provider-alicloud/issues/700))
- Modify pagesize on API DescribeVSWitches tp avoid ServiceUnavailable ([#698](https://github.com/terraform-providers/terraform-provider-alicloud/issues/698))
- Improve eip resource and data source testcases ([#697](https://github.com/terraform-providers/terraform-provider-alicloud/issues/697))

BUG FIXES:

- FIx cen route NotFoundRoute error when deleting ([#753](https://github.com/terraform-providers/terraform-provider-alicloud/issues/753))
- Fix log_store InternalServerError error ([#737](https://github.com/terraform-providers/terraform-provider-alicloud/issues/737))
- Fix cen region route entries testcase bug ([#734](https://github.com/terraform-providers/terraform-provider-alicloud/issues/734))
- Fix ots_table StorageServerBusy bug ([#733](https://github.com/terraform-providers/terraform-provider-alicloud/issues/733))
- Fix db_account setting description bug ([#732](https://github.com/terraform-providers/terraform-provider-alicloud/issues/732))
- Fix Router Entry Token Bug ([#730](https://github.com/terraform-providers/terraform-provider-alicloud/issues/730))
- Fix instance diff bug when updating its VPC attributes ([#728](https://github.com/terraform-providers/terraform-provider-alicloud/issues/728))
- Fix snat entry IncorretSnatEntryStatus error when deleting ([#714](https://github.com/terraform-providers/terraform-provider-alicloud/issues/714))
- Fix forward entry UnknownError error ([#712](https://github.com/terraform-providers/terraform-provider-alicloud/issues/712))
- Fix pvtz record Zone.NotExists error when deleting record ([#710](https://github.com/terraform-providers/terraform-provider-alicloud/issues/710))
- Fix modify kvstore policy not working bug ([#709](https://github.com/terraform-providers/terraform-provider-alicloud/issues/709))
- reattach the key pair after update OS image ([#699](https://github.com/terraform-providers/terraform-provider-alicloud/issues/699))
- Fix ServiceUnavailable error on VPC and VSW ([#695](https://github.com/terraform-providers/terraform-provider-alicloud/issues/695))

## 1.29.0 (January 21, 2019)

FEATURES:

- **New Resource:** `alicloud_logtail_config` ([#685](https://github.com/terraform-providers/terraform-provider-alicloud/issues/685))

IMPROVEMENTS:

- Apply wraperror to ess group ([#689](https://github.com/terraform-providers/terraform-provider-alicloud/issues/689))
- Add wraperror and apply it to vpc and eip ([#688](https://github.com/terraform-providers/terraform-provider-alicloud/issues/688))
- Improve vswitch resource and data source testcases ([#687](https://github.com/terraform-providers/terraform-provider-alicloud/issues/687))
- Improve security_group resource and data source testcases ([#686](https://github.com/terraform-providers/terraform-provider-alicloud/issues/686))
- Improve vpc resource and data source testcases ([#684](https://github.com/terraform-providers/terraform-provider-alicloud/issues/684))
- Modify the slb sever group testcase name ([#681](https://github.com/terraform-providers/terraform-provider-alicloud/issues/681))
- Improve sweeper testcases ([#680](https://github.com/terraform-providers/terraform-provider-alicloud/issues/680))
- Improve db instance's testcases ([#679](https://github.com/terraform-providers/terraform-provider-alicloud/issues/679))
- Improve ecs disk's testcases ([#678](https://github.com/terraform-providers/terraform-provider-alicloud/issues/678))
- Add multi_zone_ids for datasource alicloud_zones ([#677](https://github.com/terraform-providers/terraform-provider-alicloud/issues/677))
- Improve redis and memcache instance testcases ([#676](https://github.com/terraform-providers/terraform-provider-alicloud/issues/676))
- Improve ecs instance testcases ([#675](https://github.com/terraform-providers/terraform-provider-alicloud/issues/675))

BUG FIXES:

- Fix oss bucket docs error ([#692](https://github.com/terraform-providers/terraform-provider-alicloud/issues/692))
- Fix pvtz 'Zone.VpcExists' error ([#691](https://github.com/terraform-providers/terraform-provider-alicloud/issues/691))
- Fix multi-k8s testcase failed error ([#683](https://github.com/terraform-providers/terraform-provider-alicloud/issues/683))
- Fix pvtz attchment Zone.NotExists error ([#682](https://github.com/terraform-providers/terraform-provider-alicloud/issues/682))
- Fix deleting ram role error ([#674](https://github.com/terraform-providers/terraform-provider-alicloud/issues/674))
- Fix k8s cluster worker_period_unit type error ([#672](https://github.com/terraform-providers/terraform-provider-alicloud/issues/672))

## 1.28.0 (January 16, 2019)

IMPROVEMENTS:

- Ots service support https ([#669](https://github.com/terraform-providers/terraform-provider-alicloud/issues/669))
- check vswitch id when creating instance  documentation ([#668](https://github.com/terraform-providers/terraform-provider-alicloud/issues/668))
- Improve pvtz attachment test updating case ([#663](https://github.com/terraform-providers/terraform-provider-alicloud/issues/663))
- add vswitch id checker when creating k8s clusters ([#656](https://github.com/terraform-providers/terraform-provider-alicloud/issues/656))
- Improve cen instance testcase to avoid mistake query ([#655](https://github.com/terraform-providers/terraform-provider-alicloud/issues/655))
- Improve route entry retry strategy to avoid concurrence issue ([#654](https://github.com/terraform-providers/terraform-provider-alicloud/issues/654))
- Offline drds resource from website results from drds does not support idempotent ([#653](https://github.com/terraform-providers/terraform-provider-alicloud/issues/653))
- Support customer endpoints in the provider ([#652](https://github.com/terraform-providers/terraform-provider-alicloud/issues/652))
- Reback image filter to meet many non-ecs testcase ([#649](https://github.com/terraform-providers/terraform-provider-alicloud/issues/649))
- Improve ecs instance testcase by update instance type ([#646](https://github.com/terraform-providers/terraform-provider-alicloud/issues/646))
- Support cs client setting customer endpoint ([#643](https://github.com/terraform-providers/terraform-provider-alicloud/issues/643))
- do not poll nodes when k8s cluster is stable ([#641](https://github.com/terraform-providers/terraform-provider-alicloud/issues/641))
- Improve pvtz_zone testcase by using rand ([#639](https://github.com/terraform-providers/terraform-provider-alicloud/issues/639))
- support for zero node clusters in swarm container service ([#638](https://github.com/terraform-providers/terraform-provider-alicloud/issues/638))
- Slb listener can not be updated when load balancer instance is shared-performance ([#637](https://github.com/terraform-providers/terraform-provider-alicloud/issues/637))
- Improve db_account testcase and its docs ([#635](https://github.com/terraform-providers/terraform-provider-alicloud/issues/635))
- Adding https_config options to the alicloud_cdn_domain resource ([#605](https://github.com/terraform-providers/terraform-provider-alicloud/issues/605))

BUG FIXES:

- Fix slb OperationFailed.TokenIsProcessing error ([#667](https://github.com/terraform-providers/terraform-provider-alicloud/issues/667))
- Fix deleting log project requestTimeout error ([#666](https://github.com/terraform-providers/terraform-provider-alicloud/issues/666))
- Fix cs_kubernetes setting int value error ([#665](https://github.com/terraform-providers/terraform-provider-alicloud/issues/665))
- Fix pvtz zone attaching vpc system busy error ([#660](https://github.com/terraform-providers/terraform-provider-alicloud/issues/660))
- Fix ecs and ess tags read bug with ignore system tag ([#659](https://github.com/terraform-providers/terraform-provider-alicloud/issues/659))
- Fix cs cluster not found error and improve its testcase ([#658](https://github.com/terraform-providers/terraform-provider-alicloud/issues/658))
- Fix deleting pvtz zone not exist and internal error ([#657](https://github.com/terraform-providers/terraform-provider-alicloud/issues/657))
- Fix pvtz throttling user bug and improve WrapError ([#650](https://github.com/terraform-providers/terraform-provider-alicloud/issues/650))
- Fix ess group describing error ([#644](https://github.com/terraform-providers/terraform-provider-alicloud/issues/644))
- Fix pvtz throttling user bug and add WrapError ([#642](https://github.com/terraform-providers/terraform-provider-alicloud/issues/642))
- Fix kvstore instance docs ([#636](https://github.com/terraform-providers/terraform-provider-alicloud/issues/636))

## 1.27.0 (January 08, 2019)

IMPROVEMENTS:

- Improve slb instance docs ([#632](https://github.com/terraform-providers/terraform-provider-alicloud/issues/632))
- Upgrade to Go 1.11 ([#629](https://github.com/terraform-providers/terraform-provider-alicloud/issues/629))
- Remove ots https schema because of in some region only supports http ([#630](https://github.com/terraform-providers/terraform-provider-alicloud/issues/630))
- Support https for log client ([#623](https://github.com/terraform-providers/terraform-provider-alicloud/issues/623))
- Support https for ram, cdn, kms and fc client ([#622](https://github.com/terraform-providers/terraform-provider-alicloud/issues/622))
- Support https for dns client ([#621](https://github.com/terraform-providers/terraform-provider-alicloud/issues/621))
- Support https for services client using official sdk ([#619](https://github.com/terraform-providers/terraform-provider-alicloud/issues/619))
- Support mns client https and improve mns testcase ([#618](https://github.com/terraform-providers/terraform-provider-alicloud/issues/618))
- Support oss client https ([#617](https://github.com/terraform-providers/terraform-provider-alicloud/issues/617))
- Support change kvstore instance charge type ([#602](https://github.com/terraform-providers/terraform-provider-alicloud/issues/602))
- add region checks to kubernetes, multiaz kubernetes, swarm clusters ([#607](https://github.com/terraform-providers/terraform-provider-alicloud/issues/607))
- Add forcenew for ess lifecycle hook name and improve ess testcase by random name ([#603](https://github.com/terraform-providers/terraform-provider-alicloud/issues/603))
- Improve ess configuration testcase ([#600](https://github.com/terraform-providers/terraform-provider-alicloud/issues/600))
- Improve kvstore and ess schedule testcase ([#599](https://github.com/terraform-providers/terraform-provider-alicloud/issues/599))
- Improve apigateway testcase ([#593](https://github.com/terraform-providers/terraform-provider-alicloud/issues/593))
- Improve ram, ess schedule and cdn testcase ([#592](https://github.com/terraform-providers/terraform-provider-alicloud/issues/592))
- Improve kvstore client token ([#586](https://github.com/terraform-providers/terraform-provider-alicloud/issues/586))

BUG FIXES:

- Fix api gateway deleteing app bug ([#633](https://github.com/terraform-providers/terraform-provider-alicloud/issues/633))
- Fix cs_kubernetes missing name error ([#625](https://github.com/terraform-providers/terraform-provider-alicloud/issues/625))
- Fix api gateway groups filter bug ([#624](https://github.com/terraform-providers/terraform-provider-alicloud/issues/624))
- Fix ots instance description force new bug ([#616](https://github.com/terraform-providers/terraform-provider-alicloud/issues/616))
- Fix oss bucket object testcase destroy bug ([#605](https://github.com/terraform-providers/terraform-provider-alicloud/issues/605))
- Fix deleting ess group timeout bug ([#604](https://github.com/terraform-providers/terraform-provider-alicloud/issues/604))
- Fix deleting mns subscription bug ([#601](https://github.com/terraform-providers/terraform-provider-alicloud/issues/601))
- bug fix for the input of cen bandwidth limit ([#598](https://github.com/terraform-providers/terraform-provider-alicloud/issues/598))
- Fix log service timeout error ([#594](https://github.com/terraform-providers/terraform-provider-alicloud/issues/594))
- Fix record not found issue if pvtz records are more than 50 ([#590](https://github.com/terraform-providers/terraform-provider-alicloud/issues/590))
- Fix cen instance and bandwidth multi regions test case bug ([#588](https://github.com/terraform-providers/terraform-provider-alicloud/issues/588))

## 1.26.0 (December 20, 2018)

FEATURES:

- **New Resource:** `alicloud_cs_managed_kubernetes` ([#563](https://github.com/terraform-providers/terraform-provider-alicloud/issues/563))

IMPROVEMENTS:

- Improve ram client endpoint ([#584](https://github.com/terraform-providers/terraform-provider-alicloud/issues/584))
- Remove useless sweeper depencences for alicloud_instance sweeper testcase ([#582](https://github.com/terraform-providers/terraform-provider-alicloud/issues/582))
- Improve kvstore backup policy testcase ([#580](https://github.com/terraform-providers/terraform-provider-alicloud/issues/580))
- Improve the describing endpoint ([#579](https://github.com/terraform-providers/terraform-provider-alicloud/issues/579))
- VPN gateway supports 200/500/1000M bandwidth ([#577](https://github.com/terraform-providers/terraform-provider-alicloud/issues/577))
- skip private ip test in some regions ([#575](https://github.com/terraform-providers/terraform-provider-alicloud/issues/575))
- Add timeout and retry for tablestore client and Improve its testcases ([#569](https://github.com/terraform-providers/terraform-provider-alicloud/issues/569))
- Modify kvstore_instance password to Optional and improve its testcases ([#567](https://github.com/terraform-providers/terraform-provider-alicloud/issues/567))
- Improve datasource alicloud_vpcs testcase ([#566](https://github.com/terraform-providers/terraform-provider-alicloud/issues/566))
- Improve dns_domains testcase ([#561](https://github.com/terraform-providers/terraform-provider-alicloud/issues/561))
- Improve ram_role_attachment testcase ([#560](https://github.com/terraform-providers/terraform-provider-alicloud/issues/560))
- support PrePaid instances, image_id to be set when creating k8s clusters ([#559](https://github.com/terraform-providers/terraform-provider-alicloud/issues/559))
- Add retry and timemout for fc client ([#557](https://github.com/terraform-providers/terraform-provider-alicloud/issues/557))
- Datasource alicloud_zones supports filter FunctionCompute ([#555](https://github.com/terraform-providers/terraform-provider-alicloud/issues/555))
- Fix a bug that caused the alicloud_dns_record.routing attribute ([#554](https://github.com/terraform-providers/terraform-provider-alicloud/issues/554))
- Modify router interface prepaid test case  documentation ([#552](https://github.com/terraform-providers/terraform-provider-alicloud/issues/552))
- Resource alicloud_ess_scalingconfiguration supports system_disk_size ([#551](https://github.com/terraform-providers/terraform-provider-alicloud/issues/551))
- Improve datahub project testcase ([#548](https://github.com/terraform-providers/terraform-provider-alicloud/issues/548))
- resource alicloud_slb_listener support server group ([#545](https://github.com/terraform-providers/terraform-provider-alicloud/issues/545))
- Improve ecs instance and disk testcase with common case ([#544](https://github.com/terraform-providers/terraform-provider-alicloud/issues/544))

BUG FIXES:

- Fix provider compile error on 32bit ([#585](https://github.com/terraform-providers/terraform-provider-alicloud/issues/585))
- Fix table store no such host error with deleting and updating ([#583](https://github.com/terraform-providers/terraform-provider-alicloud/issues/583))
- Fix pvtz_record RecordInvalidConflict bug ([#581](https://github.com/terraform-providers/terraform-provider-alicloud/issues/581))
- fixed bug in backup policy update ([#521](https://github.com/terraform-providers/terraform-provider-alicloud/issues/521))
- Fix docs eip_association ([#578](https://github.com/terraform-providers/terraform-provider-alicloud/issues/578))
- Fix a bug about instance charge type change ([#576](https://github.com/terraform-providers/terraform-provider-alicloud/issues/576))
- Fix describing endpoint failed error ([#574](https://github.com/terraform-providers/terraform-provider-alicloud/issues/574))
- Fix table store describing no such host error ([#572](https://github.com/terraform-providers/terraform-provider-alicloud/issues/572))
- Fix table store creating timeout error ([#571](https://github.com/terraform-providers/terraform-provider-alicloud/issues/571))
- Fix kvstore instance class update error ([#570](https://github.com/terraform-providers/terraform-provider-alicloud/issues/570))
- Fix ess_scaling_group import bugs and improve ess schedule testcase ([#565](https://github.com/terraform-providers/terraform-provider-alicloud/issues/565))
- Fix alicloud rds related IncorrectStatus bug ([#558](https://github.com/terraform-providers/terraform-provider-alicloud/issues/558))
- Fix alicloud_fc_trigger's config diff bug ([#556](https://github.com/terraform-providers/terraform-provider-alicloud/issues/556))
- Fix oss bucket deleting failed error ([#550](https://github.com/terraform-providers/terraform-provider-alicloud/issues/550))
- Fix potential bugs of datahub and ram when the resource has been deleted ([#546](https://github.com/terraform-providers/terraform-provider-alicloud/issues/546))
- Fix pvtz_record describing bug ([#543](https://github.com/terraform-providers/terraform-provider-alicloud/issues/543))

## 1.25.0 (November 30, 2018)

IMPROVEMENTS:

- return a empty list when there is no any data source ([#540](https://github.com/terraform-providers/terraform-provider-alicloud/issues/540))
- Skip automatically the testcases which does not support API gateway ([#538](https://github.com/terraform-providers/terraform-provider-alicloud/issues/538))
- Improve common bandwidth package test case and remove PayBy95 ([#530](https://github.com/terraform-providers/terraform-provider-alicloud/issues/530))
- Update resource drds supported regions ([#534](https://github.com/terraform-providers/terraform-provider-alicloud/issues/534))
- Remove DB instance engine_version limitation ([#528](https://github.com/terraform-providers/terraform-provider-alicloud/issues/528))
- Skip automatically the testcases which does not support route table and classic drds ([#526](https://github.com/terraform-providers/terraform-provider-alicloud/issues/526))
- Skip automatically the testcases which does not support classic regions ([#524](https://github.com/terraform-providers/terraform-provider-alicloud/issues/524))
- datasource alicloud_slbs support tags ([#523](https://github.com/terraform-providers/terraform-provider-alicloud/issues/523))
- resouce alicloud_slb support tags ([#522](https://github.com/terraform-providers/terraform-provider-alicloud/issues/522))
- Skip automatically the testcases which does not support multi az regions ([#518](https://github.com/terraform-providers/terraform-provider-alicloud/issues/518))
- Add some region limitation guide for sone resources ([#517](https://github.com/terraform-providers/terraform-provider-alicloud/issues/517))
- Skip automatically the testcases which does not support some known regions ([#516](https://github.com/terraform-providers/terraform-provider-alicloud/issues/516))
- create instance with runinstances ([#514](https://github.com/terraform-providers/terraform-provider-alicloud/issues/514))
- support eni amount in data source instance types ([#512](https://github.com/terraform-providers/terraform-provider-alicloud/issues/512))
- Add a docs guides/getting-account to help user learn alibaba cloud account ([#510](https://github.com/terraform-providers/terraform-provider-alicloud/issues/510))

BUG FIXES:

- Fix route_entry concurrence bug and improve it testcases ([#537](https://github.com/terraform-providers/terraform-provider-alicloud/issues/537))
- Fix router interface prepaid purchase ([#529](https://github.com/terraform-providers/terraform-provider-alicloud/issues/529))
- Fix fc_service sweeper test bug ([#536](https://github.com/terraform-providers/terraform-provider-alicloud/issues/536))
- Fix drds creating VPC instance bug by adding vpc_id ([#531](https://github.com/terraform-providers/terraform-provider-alicloud/issues/531))
- fix a snat_entry bug without set id to empty ([#525](https://github.com/terraform-providers/terraform-provider-alicloud/issues/525))
- fix a bug of ram_use display name ([#519](https://github.com/terraform-providers/terraform-provider-alicloud/issues/519))
- fix a bug of instance testcase ([#513](https://github.com/terraform-providers/terraform-provider-alicloud/issues/513))
- Fix pvtz resource priority bug ([#511](https://github.com/terraform-providers/terraform-provider-alicloud/issues/511))

## 1.24.0 (November 21, 2018)

FEATURES:

- **New Resource:** `alicloud_drds_instance` ([#446](https://github.com/terraform-providers/terraform-provider-alicloud/issues/446))

IMPROVEMENTS:

- Improve drds_instance docs ([#509](https://github.com/terraform-providers/terraform-provider-alicloud/issues/509))
- Add a new test case for drds_instance ([#508](https://github.com/terraform-providers/terraform-provider-alicloud/issues/508))
- Improve provider config with Trim method ([#504](https://github.com/terraform-providers/terraform-provider-alicloud/issues/504))
- api gateway skip app relevant tests ([#500](https://github.com/terraform-providers/terraform-provider-alicloud/issues/500))
- update api resource that support to deploy api ([#498](https://github.com/terraform-providers/terraform-provider-alicloud/issues/498))
- Skip ram_groups a test case ([#496](https://github.com/terraform-providers/terraform-provider-alicloud/issues/496))
- support disk resize ([#490](https://github.com/terraform-providers/terraform-provider-alicloud/issues/490))
- cancel the limit of system disk size ([#489](https://github.com/terraform-providers/terraform-provider-alicloud/issues/489))
- Improve docs alicloud_db_database and alicloud_cs_kubernetes ([#488](https://github.com/terraform-providers/terraform-provider-alicloud/issues/488))
- Support creating data disk with instance ([#484](https://github.com/terraform-providers/terraform-provider-alicloud/issues/484))

BUG FIXES:

- Fix the sweeper test for CEN and CEN bandwidth package ([#505](https://github.com/terraform-providers/terraform-provider-alicloud/issues/505))
- Fix pvtz_zone_record update bug ([#503](https://github.com/terraform-providers/terraform-provider-alicloud/issues/503))
- Fix network_interface_attachment docs error ([#502](https://github.com/terraform-providers/terraform-provider-alicloud/issues/502))
- fix fix datahub bug when visit region of ap-southeast-1 ([#499](https://github.com/terraform-providers/terraform-provider-alicloud/issues/499))
- Fix examples/mns-topic parameter error ([#497](https://github.com/terraform-providers/terraform-provider-alicloud/issues/497))
- Fix db_connection not found error when deleting ([#495](https://github.com/terraform-providers/terraform-provider-alicloud/issues/495))
- fix error about the docs format  ([#492](https://github.com/terraform-providers/terraform-provider-alicloud/issues/492))

## 1.23.0 (November 13, 2018)

FEATURES:

- **New Resource:** `alicloud_api_gateway_app_attachment` ([#478](https://github.com/terraform-providers/terraform-provider-alicloud/issues/478))
- **New Resource:** `alicloud_network_interface_attachment` ([#474](https://github.com/terraform-providers/terraform-provider-alicloud/issues/474))
- **New Resource:** `alicloud_api_gateway_vpc_access` ([#472](https://github.com/terraform-providers/terraform-provider-alicloud/issues/472))
- **New Resource:** `alicloud_network_interface` ([#469](https://github.com/terraform-providers/terraform-provider-alicloud/issues/469))
- **New Resource:** `alicloud_common_bandwidth_package` ([#468](https://github.com/terraform-providers/terraform-provider-alicloud/issues/468))
- **New Data Source:** `alicloud_network_interfaces` ([#475](https://github.com/terraform-providers/terraform-provider-alicloud/issues/475))
- **New Data Source:** `alicloud_api_gateway_apps` ([#467](https://github.com/terraform-providers/terraform-provider-alicloud/issues/467))

IMPROVEMENTS:

- Add a new region eu-west-1 ([#486](https://github.com/terraform-providers/terraform-provider-alicloud/issues/486))
- remove unreachable codes ([#479](https://github.com/terraform-providers/terraform-provider-alicloud/issues/479))
- support enable/disable security enhancement strategy of alicloud_instance ([#471](https://github.com/terraform-providers/terraform-provider-alicloud/issues/471))
- alicloud_slb_listener support idle_timeout/request_timeout ([#463](https://github.com/terraform-providers/terraform-provider-alicloud/issues/463))

BUG FIXES:

- Fix cs_application cluster not found ([#480](https://github.com/terraform-providers/terraform-provider-alicloud/issues/480))
- fix the bug of security_group inner_access bug ([#477](https://github.com/terraform-providers/terraform-provider-alicloud/issues/477))
- Fix pagenumber built error ([#470](https://github.com/terraform-providers/terraform-provider-alicloud/issues/470))
- Fix cs_application cluster not found ([#480](https://github.com/terraform-providers/terraform-provider-alicloud/issues/480))

## 1.22.0 (November 02, 2018)

FEATURES:

- **New Resource:** `alicloud_api_gateway_api` ([#457](https://github.com/terraform-providers/terraform-provider-alicloud/issues/457))
- **New Resource:** `alicloud_api_gateway_app` ([#462](https://github.com/terraform-providers/terraform-provider-alicloud/issues/462))
- **New Reource:** `alicloud_common_bandwidth_package` ([#454](https://github.com/terraform-providers/terraform-provider-alicloud/issues/454))
- **New Data Source:** `alicloud_api_gateway_apis` ([#458](https://github.com/terraform-providers/terraform-provider-alicloud/issues/458))
- **New Data Source:** `cen_region_route_entries` ([#442](https://github.com/terraform-providers/terraform-provider-alicloud/issues/442))
- **New Data Source:** `alicloud_slb_ca_certificates` ([#452](https://github.com/terraform-providers/terraform-provider-alicloud/issues/452))

IMPROVEMENTS:

- Use product code to get common request domain ([#466](https://github.com/terraform-providers/terraform-provider-alicloud/issues/466))
- KVstore instance password supports at sign ([#465](https://github.com/terraform-providers/terraform-provider-alicloud/issues/465))
- Correct docs spelling error ([#464](https://github.com/terraform-providers/terraform-provider-alicloud/issues/464))
- alicloud_log_service : support update project and shard auto spit ([#461](https://github.com/terraform-providers/terraform-provider-alicloud/issues/461))
- Correct datasource alicloud_cen_route_entries docs error ([#460](https://github.com/terraform-providers/terraform-provider-alicloud/issues/460))
- Remove CDN default configuration ([#450](https://github.com/terraform-providers/terraform-provider-alicloud/issues/450))

BUG FIXES:

- set number of cen instances five for normal alicloud account testcases ([#459](https://github.com/terraform-providers/terraform-provider-alicloud/issues/459))

## 1.21.0 (October 30, 2018)

FEATURES:

- **New Data Source:** `alicloud_slb_server_certificates` ([#444](https://github.com/terraform-providers/terraform-provider-alicloud/issues/444))
- **New Data Source:** `alicloud_slb_acls` ([#443](https://github.com/terraform-providers/terraform-provider-alicloud/issues/443))
- **New Resource:** `alicloud_slb_ca_certificate` ([#438](https://github.com/terraform-providers/terraform-provider-alicloud/issues/438))
- **New Resource:** `alicloud_slb_server_certificate` ([#436](https://github.com/terraform-providers/terraform-provider-alicloud/issues/436))

IMPROVEMENTS:

- resource alicloud_slb_listener tcp protocol support established_timeout parameter ([#440](https://github.com/terraform-providers/terraform-provider-alicloud/issues/440))

BUG FIXES:

- Fix mns resource docs bug ([#441](https://github.com/terraform-providers/terraform-provider-alicloud/issues/441))

## 1.20.0 (October 22, 2018)

FEATURES:

- **New Resource:** `alicloud_slb_acl` ([#413](https://github.com/terraform-providers/terraform-provider-alicloud/issues/413))
- **New Resource:** `alicloud_cen_route_entry` ([#415](https://github.com/terraform-providers/terraform-provider-alicloud/issues/415))
- **New Data Source:** `alicloud_cen_route_entries` ([#424](https://github.com/terraform-providers/terraform-provider-alicloud/issues/424))

IMPROVEMENTS:

- Improve datahub_project sweeper test ([#435](https://github.com/terraform-providers/terraform-provider-alicloud/issues/435))
- Modify mns test case name ([#434](https://github.com/terraform-providers/terraform-provider-alicloud/issues/434))
- Improve fc_service sweeper test ([#433](https://github.com/terraform-providers/terraform-provider-alicloud/issues/433))
- Support provider thread safety ([#432](https://github.com/terraform-providers/terraform-provider-alicloud/issues/432))
- add tags to security group ([#423](https://github.com/terraform-providers/terraform-provider-alicloud/issues/423))
- Resource router_interface support PrePaid ([#425](https://github.com/terraform-providers/terraform-provider-alicloud/issues/425))
- resource alicloud_slb_listener support acl ([#426](https://github.com/terraform-providers/terraform-provider-alicloud/issues/426))
- change child instance type Vbr to VBR and replace some const variables ([#422](https://github.com/terraform-providers/terraform-provider-alicloud/issues/422))
- add slb_internet_enabled to Kubernetes Cluster ([#421](https://github.com/terraform-providers/terraform-provider-alicloud/issues/421))
- Hide AliCloud HaVip Attachment resource docs because of it is not public totally ([#420](https://github.com/terraform-providers/terraform-provider-alicloud/issues/420))
- Improve examples/ots-table ([#417](https://github.com/terraform-providers/terraform-provider-alicloud/issues/417))
- Improve examples ecs-vpc, ecs-new-vpc and api-gateway ([#416](https://github.com/terraform-providers/terraform-provider-alicloud/issues/416))

BUG FIXES:

- Fix reources' id description bugs ([#428](https://github.com/terraform-providers/terraform-provider-alicloud/issues/428))
- Fix alicloud_ess_scaling_configuration setting data_disk failed ([#427](https://github.com/terraform-providers/terraform-provider-alicloud/issues/427))

## 1.19.0 (October 13, 2018)

FEATURES:

- **New Resource:** `alicloud_api_gateway_group` ([#409](https://github.com/terraform-providers/terraform-provider-alicloud/issues/409))
- **New Resource:** `alicloud_datahub_subscription` ([#405](https://github.com/terraform-providers/terraform-provider-alicloud/issues/405))
- **New Resource:** `alicloud_datahub_topic` ([#404](https://github.com/terraform-providers/terraform-provider-alicloud/issues/404))
- **New Resource:** `alicloud_datahub_project` ([#403](https://github.com/terraform-providers/terraform-provider-alicloud/issues/403))
- **New Data Source:** `alicloud_api_gateway_groups` ([#412](https://github.com/terraform-providers/terraform-provider-alicloud/issues/412))
- **New Data Source:** `alicloud_cen_bandwidth_limits` ([#402](https://github.com/terraform-providers/terraform-provider-alicloud/issues/402))

IMPROVEMENTS:

- added need_slb attribute to cs swarm ([#414](https://github.com/terraform-providers/terraform-provider-alicloud/issues/414))
- Add new example/datahub ([#407](https://github.com/terraform-providers/terraform-provider-alicloud/issues/407))
- Add new example/datahub ([#406](https://github.com/terraform-providers/terraform-provider-alicloud/issues/406))
- Format examples ([#397](https://github.com/terraform-providers/terraform-provider-alicloud/issues/397))
- Add new example/kvstore ([#396](https://github.com/terraform-providers/terraform-provider-alicloud/issues/396))
- Remove useless datasource cache file ([#395](https://github.com/terraform-providers/terraform-provider-alicloud/issues/395))
- Add new example/pvtz ([#394](https://github.com/terraform-providers/terraform-provider-alicloud/issues/394))
- Improve example/ecs-key-pair ([#393](https://github.com/terraform-providers/terraform-provider-alicloud/issues/393))
- Change key pair file mode to 400 ([#392](https://github.com/terraform-providers/terraform-provider-alicloud/issues/392))

BUG FIXES:

- fix kubernetes's new_nat_gateway issue ([#410](https://github.com/terraform-providers/terraform-provider-alicloud/issues/410))
- modify the mns err info ([#400](https://github.com/terraform-providers/terraform-provider-alicloud/issues/400))
- Skip havip test case ([#399](https://github.com/terraform-providers/terraform-provider-alicloud/issues/399))
- modify the sweeptest nameprefix ([#398](https://github.com/terraform-providers/terraform-provider-alicloud/issues/398))

## 1.18.0 (October 09, 2018)

FEATURES:

- **New Resource:** `alicloud_havip` ([#378](https://github.com/terraform-providers/terraform-provider-alicloud/issues/378))
- **New Resource:** `alicloud_havip_attachment` ([#388](https://github.com/terraform-providers/terraform-provider-alicloud/issues/388))
- **New Resource:** `alicloud_mns_topic_subscription` ([#376](https://github.com/terraform-providers/terraform-provider-alicloud/issues/376))
- **New Resource:** `alicloud_route_table_attachment` ([#362](https://github.com/terraform-providers/terraform-provider-alicloud/issues/362))
- **New Resource:** `alicloud_cen_bandwidth_limit` ([#361](https://github.com/terraform-providers/terraform-provider-alicloud/issues/361))
- **New Resource:** `alicloud_mns_topic` ([#374](https://github.com/terraform-providers/terraform-provider-alicloud/issues/374))
- **New Resource:** `alicloud_mns_queue` ([#365](https://github.com/terraform-providers/terraform-provider-alicloud/issues/365))
- **New Resource:** `alicloud_cen_bandwidth_package_attachment` ([#354](https://github.com/terraform-providers/terraform-provider-alicloud/issues/354))
- **New Resource:** `alicloud_route_table` ([#356](https://github.com/terraform-providers/terraform-provider-alicloud/issues/356))
- **New Data Source:** `alicloud_mns_queues` ([#382](https://github.com/terraform-providers/terraform-provider-alicloud/issues/382))
- **New Data Source:** `alicloud_mns_topics` ([#384](https://github.com/terraform-providers/terraform-provider-alicloud/issues/384))
- **New Data Source:** `alicloud_mns_topic_subscriptions` ([#386](https://github.com/terraform-providers/terraform-provider-alicloud/issues/386))
- **New Data Source:** `alicloud_cen_bandwidth_packages` ([#367](https://github.com/terraform-providers/terraform-provider-alicloud/issues/367))
- **New Data Source:** `alicloud_vpn_connections` ([#366](https://github.com/terraform-providers/terraform-provider-alicloud/issues/366))
- **New Data Source:** `alicloud_vpn_gateways` ([#363](https://github.com/terraform-providers/terraform-provider-alicloud/issues/363))
- **New Data Source:** `alicloud_vpn_customer_gateways` ([#364](https://github.com/terraform-providers/terraform-provider-alicloud/issues/364))
- **New Data Source:** `alicloud_cen_instances` ([#342](https://github.com/terraform-providers/terraform-provider-alicloud/issues/342))

IMPROVEMENTS:

- Improve resource ram_policy's document validatefunc ([#385](https://github.com/terraform-providers/terraform-provider-alicloud/issues/385))
- RAM support useragent ([#383](https://github.com/terraform-providers/terraform-provider-alicloud/issues/383))
- add node_cidr_mas and log_config, fix worker_data_disk issue ([#368](https://github.com/terraform-providers/terraform-provider-alicloud/issues/368))
- Improve WaitForRouteTable and WaitForRouteTableAttachment method ([#375](https://github.com/terraform-providers/terraform-provider-alicloud/issues/375))
- Correct Function Compute conn ([#371](https://github.com/terraform-providers/terraform-provider-alicloud/issues/371))
- Improve datasource `images`'s docs ([#370](https://github.com/terraform-providers/terraform-provider-alicloud/issues/370))
- add worker_data_disk_category and worker_data_disk_size to kubernetes creation ([#355](https://github.com/terraform-providers/terraform-provider-alicloud/issues/355))

BUG FIXES:

- Fix alicloud_ram_user_policy_attachment EntityNotExist.User error ([#381](https://github.com/terraform-providers/terraform-provider-alicloud/issues/381))
- Add parameter 'force_delete' to support deleting 'PrePaid' instance ([#377](https://github.com/terraform-providers/terraform-provider-alicloud/issues/377))
- Add wait time to fix random detaching disk error ([#373](https://github.com/terraform-providers/terraform-provider-alicloud/issues/373))
- Fix cen_instances markdown ([#372](https://github.com/terraform-providers/terraform-provider-alicloud/issues/372))

## 1.17.0 (September 22, 2018)

FEATURES:

- **New Data Source:** `alicloud_fc_triggers` ([#351](https://github.com/terraform-providers/terraform-provider-alicloud/pull/351))
- **New Data Source:** `alicloud_oss_bucket_objects` ([#350](https://github.com/terraform-providers/terraform-provider-alicloud/pull/350))
- **New Data Source:** `alicloud_fc_functions` ([#349](https://github.com/terraform-providers/terraform-provider-alicloud/pull/349))
- **New Data Source:** `alicloud_fc_services` ([#348](https://github.com/terraform-providers/terraform-provider-alicloud/pull/348))
- **New Data Source:** `alicloud_oss_buckets` ([#345](https://github.com/terraform-providers/terraform-provider-alicloud/pull/345))
- **New Data Source:** `alicloud_disks` ([#343](https://github.com/terraform-providers/terraform-provider-alicloud/pull/343))
- **New Resource:** `alicloud_cen_bandwidth_package` ([#333](https://github.com/terraform-providers/terraform-provider-alicloud/pull/333))

IMPROVEMENTS:

- Update OSS Resources' link to English ([#352](https://github.com/terraform-providers/terraform-provider-alicloud/pull/352))
- Improve example/kubernetes to support multi-az ([#344](https://github.com/terraform-providers/terraform-provider-alicloud/pull/344))

## 1.16.0 (September 16, 2018)

FEATURES:

- **New Resource:** `alicloud_cen_instance_attachment` ([#327](https://github.com/terraform-providers/terraform-provider-alicloud/pull/327))

IMPROVEMENTS:

- Allow setting the scaling group balancing policy ([#339](https://github.com/terraform-providers/terraform-provider-alicloud/pull/339))
- cs_kubernetes supports multi-az ([#222](https://github.com/terraform-providers/terraform-provider-alicloud/pull/222))
- Improve client token using timestemp ([#326](https://github.com/terraform-providers/terraform-provider-alicloud/pull/326))

BUG FIXES:

- Fix alicloud db connection ([#341](https://github.com/terraform-providers/terraform-provider-alicloud/pull/341))
- Fix knstore productId ([#338](https://github.com/terraform-providers/terraform-provider-alicloud/pull/338))
- Fix retriving kvstore multi zones bug ([#337](https://github.com/terraform-providers/terraform-provider-alicloud/pull/337))
- Fix kvstore instance period bug ([#335](https://github.com/terraform-providers/terraform-provider-alicloud/pull/335))
- Fix kvstore docs bug ([#334](https://github.com/terraform-providers/terraform-provider-alicloud/pull/334))

## 1.15.0 (September 07, 2018)

FEATURES:

- **New Resource:** `alicloud_kvstore_backup_policy` ([#331](https://github.com/terraform-providers/terraform-provider-alicloud/pull/331))
- **New Resource:** `alicloud_kvstore_instance` ([#330](https://github.com/terraform-providers/terraform-provider-alicloud/pull/330))
- **New Data Source:** `alicloud_kvstore_instances` ([#329](https://github.com/terraform-providers/terraform-provider-alicloud/pull/329))
- **New Resource:** `alicloud_ess_alarm` ([#328](https://github.com/terraform-providers/terraform-provider-alicloud/pull/328))
- **New Resource:** `alicloud_ssl_vpn_client_cert` ([#317](https://github.com/terraform-providers/terraform-provider-alicloud/pull/317))
- **New Resource:** `alicloud_cen_instance` ([#312](https://github.com/terraform-providers/terraform-provider-alicloud/pull/312))
- **New Data Source:** `alicloud_slb_server_groups`  ([#324](https://github.com/terraform-providers/terraform-provider-alicloud/pull/324))
- **New Data Source:** `alicloud_slb_rules`  ([#323](https://github.com/terraform-providers/terraform-provider-alicloud/pull/323))
- **New Data Source:** `alicloud_slb_listeners`  ([#323](https://github.com/terraform-providers/terraform-provider-alicloud/pull/323))
- **New Data Source:** `alicloud_slb_attachments`  ([#322](https://github.com/terraform-providers/terraform-provider-alicloud/pull/322))
- **New Data Source:** `alicloud_slbs`  ([#321](https://github.com/terraform-providers/terraform-provider-alicloud/pull/321))
- **New Data Source:** `alicloud_account`  ([#319](https://github.com/terraform-providers/terraform-provider-alicloud/pull/319))
- **New Resource:** `alicloud_ssl_vpn_server` ([#313](https://github.com/terraform-providers/terraform-provider-alicloud/pull/313))

IMPROVEMENTS:

- Support sweeper to clean some resources coming from failed testcases ([#326](https://github.com/terraform-providers/terraform-provider-alicloud/pull/326))
- Improve function compute tst cases ([#325](https://github.com/terraform-providers/terraform-provider-alicloud/pull/325))
- Improve fc test case using new datasource `alicloud_account` ([#320](https://github.com/terraform-providers/terraform-provider-alicloud/pull/320))
- Base64 encode ESS scaling config user_data ([#315](https://github.com/terraform-providers/terraform-provider-alicloud/pull/315))
- Retrieve the account_id automatically if needed ([#314](https://github.com/terraform-providers/terraform-provider-alicloud/pull/314))

BUG FIXES:

- Fix DNS tests falied error ([#318](https://github.com/terraform-providers/terraform-provider-alicloud/pull/318))
- Fix DB database not found error ([#316](https://github.com/terraform-providers/terraform-provider-alicloud/pull/316))

## 1.14.0 (August 31, 2018)

FEATURES:

- **New Resource:** `alicloud_vpn_connection` ([#304](https://github.com/terraform-providers/terraform-provider-alicloud/pull/304))
- **New Resource:** `alicloud_vpn_customer_gateway` ([#299](https://github.com/terraform-providers/terraform-provider-alicloud/pull/299))

IMPROVEMENTS:

- Add 'force' to make key pair affect immediately ([#310](https://github.com/terraform-providers/terraform-provider-alicloud/pull/310))
- Improve http proxy support ([#307](https://github.com/terraform-providers/terraform-provider-alicloud/pull/307))
- Add flags to skip tests that use features not supported in all regions ([#306](https://github.com/terraform-providers/terraform-provider-alicloud/pull/306))
- Improve data source dns_domains test case ([#305](https://github.com/terraform-providers/terraform-provider-alicloud/pull/305))
- Change SDK config timeout ([#302](https://github.com/terraform-providers/terraform-provider-alicloud/pull/302))
- Support ClientToken for some request ([#301](https://github.com/terraform-providers/terraform-provider-alicloud/pull/301))
- Enlarge sdk default timeout to fix some timeout scenario ([#300](https://github.com/terraform-providers/terraform-provider-alicloud/pull/300))

BUG FIXES:

- Fix container cluster SDK timezone error ([#308](https://github.com/terraform-providers/terraform-provider-alicloud/pull/308))
- Fix network products throttling error ([#303](https://github.com/terraform-providers/terraform-provider-alicloud/pull/303))

## 1.13.0 (August 28, 2018)

FEATURES:

- **New Resource:** `alicloud_vpn_gateway` ([#298](https://github.com/terraform-providers/terraform-provider-alicloud/pull/298))
- **New Data Source:** `alicloud_mongo_instances` ([#221](https://github.com/terraform-providers/terraform-provider-alicloud/pull/221))
- **New Data Source:** `alicloud_pvtz_zone_records` ([#288](https://github.com/terraform-providers/terraform-provider-alicloud/pull/288))
- **New Data Source:** `alicloud_pvtz_zones` ([#287](https://github.com/terraform-providers/terraform-provider-alicloud/pull/287))
- **New Resource:** `alicloud_pvtz_zone_record` ([#286](https://github.com/terraform-providers/terraform-provider-alicloud/pull/286))
- **New Resource:** `alicloud_pvtz_zone_attachment` ([#285](https://github.com/terraform-providers/terraform-provider-alicloud/pull/285))
- **New Resource:** `alicloud_pvtz_zone` ([#284](https://github.com/terraform-providers/terraform-provider-alicloud/pull/284))
- **New Resource:** `alicloud_ess_lifecycle_hook` ([#283](https://github.com/terraform-providers/terraform-provider-alicloud/pull/283))
- **New Data Source:** `alicloud_router_interfaces` ([#269](https://github.com/terraform-providers/terraform-provider-alicloud/pull/269))

IMPROVEMENTS:

- Check pvtzconn error ([#295](https://github.com/terraform-providers/terraform-provider-alicloud/pull/295))
- For internationalize tests ([#294](https://github.com/terraform-providers/terraform-provider-alicloud/pull/294))
- Improve data source docs ([#293](https://github.com/terraform-providers/terraform-provider-alicloud/pull/293))
- Add SLB PayByBandwidth test case ([#292](https://github.com/terraform-providers/terraform-provider-alicloud/pull/292))
- Update vpc sdk to support new resource VPN gateway ([#291](https://github.com/terraform-providers/terraform-provider-alicloud/pull/291))
- Improve snat entry test case ([#290](https://github.com/terraform-providers/terraform-provider-alicloud/pull/290))
- Allow empty list of SLBs as arg to ESG ([#289](https://github.com/terraform-providers/terraform-provider-alicloud/pull/289))
- Improve docs vroute_entry ([#281](https://github.com/terraform-providers/terraform-provider-alicloud/pull/281))
- Improve examples/router_interface ([#278](https://github.com/terraform-providers/terraform-provider-alicloud/pull/278))
- Improve SLB instance test case ([#274](https://github.com/terraform-providers/terraform-provider-alicloud/pull/274))
- Improve alicloud_router_interface's test case ([#272](https://github.com/terraform-providers/terraform-provider-alicloud/pull/272))
- Improve data source alicloud_regions's test case ([#271](https://github.com/terraform-providers/terraform-provider-alicloud/pull/271))
- Add notes about ordering between two alicloud_router_interface_connections ([#270](https://github.com/terraform-providers/terraform-provider-alicloud/pull/270))
- Improve docs spelling error ([#268](https://github.com/terraform-providers/terraform-provider-alicloud/pull/268))
- ECS instance support more tags and update instance test cases ([#267](https://github.com/terraform-providers/terraform-provider-alicloud/pull/267))
- Improve OSS bucket test case ([#266](https://github.com/terraform-providers/terraform-provider-alicloud/pull/266))
- Fixing a broken link ([#265](https://github.com/terraform-providers/terraform-provider-alicloud/pull/265))
- Allow creation of slb vserver group with 0 servers ([#264](https://github.com/terraform-providers/terraform-provider-alicloud/pull/264))
- Improve SLB test cases results from international regions does support PayByBandwidth and ' Guaranteed-performance' instance ([#263](https://github.com/terraform-providers/terraform-provider-alicloud/pull/263))
- Improve EIP test cases results from international regions does support PayByBandwidth ([#262](https://github.com/terraform-providers/terraform-provider-alicloud/pull/262))
- Improve ESS test cases results from some region does support Classic Network ([#261](https://github.com/terraform-providers/terraform-provider-alicloud/pull/261))
- Recover nat gateway bandwidth pacakges to meet stock user requirements ([#260](https://github.com/terraform-providers/terraform-provider-alicloud/pull/260))
- Resource alicloud_slb_listener supports new field 'x-forwarded-for' ([#259](https://github.com/terraform-providers/terraform-provider-alicloud/pull/259))
- Resource alicloud_slb_listener supports new field 'gzip' ([#258](https://github.com/terraform-providers/terraform-provider-alicloud/pull/258))

BUG FIXES:

- Fix getting oss endpoint timeout error ([#282](https://github.com/terraform-providers/terraform-provider-alicloud/pull/282))
- Fix router interface connection error when 'opposite_interface_owner_id' is empty ([#277](https://github.com/terraform-providers/terraform-provider-alicloud/pull/277))
- Fix router interface connection error and deleting error ([#275](https://github.com/terraform-providers/terraform-provider-alicloud/pull/275))
- Fix disk detach error and improve test using dynamic zone and region ([#273](https://github.com/terraform-providers/terraform-provider-alicloud/pull/273))

## 1.12.0 (August 10, 2018)

IMPROVEMENTS:

- Improve `make build` ([#256](https://github.com/terraform-providers/terraform-provider-alicloud/pull/256))
- Improve examples slb and slb-vpc by modifying 'paybytraffic' to 'PayByTraffic' ([#256](https://github.com/terraform-providers/terraform-provider-alicloud/pull/256))
- Improve example/router-interface by adding resource alicloud_router_interface_connection ([#255](https://github.com/terraform-providers/terraform-provider-alicloud/pull/255))
- Support more specification of router interface ([#253](https://github.com/terraform-providers/terraform-provider-alicloud/pull/253))
- Improve resource alicloud_fc_service docs ([#252](https://github.com/terraform-providers/terraform-provider-alicloud/pull/252))
- Modify resource alicloud_fc_function 'handler' is required ([#251](https://github.com/terraform-providers/terraform-provider-alicloud/pull/251))
- Resource alicloud_router_interface support "import" function ([#249](https://github.com/terraform-providers/terraform-provider-alicloud/pull/249))
- Deprecate some field of alicloud_router_interface fields and use new resource instead ([#248](https://github.com/terraform-providers/terraform-provider-alicloud/pull/248))
- *New Resource*: _alicloud_router_interface_connection_ ([#247](https://github.com/terraform-providers/terraform-provider-alicloud/pull/247))

BUG FIXES:

- Fix network resource throttling error ([#257](https://github.com/terraform-providers/terraform-provider-alicloud/pull/257))
- Fix resource alicloud_fc_trigger "source_arn" inputting empty error ([#253](https://github.com/terraform-providers/terraform-provider-alicloud/pull/253))
- Fix describing vpcs with name_regex no results error ([#250](https://github.com/terraform-providers/terraform-provider-alicloud/pull/250))
- Fix creating slb listener in international region failed error ([#246](https://github.com/terraform-providers/terraform-provider-alicloud/pull/246))

## 1.11.0 (August 08, 2018)

IMPROVEMENTS:

- Resource alicloud_eip support name and description ([#244](https://github.com/terraform-providers/terraform-provider-alicloud/pull/244))
- Resource alicloud_eip support PrePaid ([#243](https://github.com/terraform-providers/terraform-provider-alicloud/pull/243))
- Correct version writting error ([#241](https://github.com/terraform-providers/terraform-provider-alicloud/pull/241))
- Change slb go sdk to official repo ([#240](https://github.com/terraform-providers/terraform-provider-alicloud/pull/240))
- Remove useless file website/fc_service.html.markdown ([#239](https://github.com/terraform-providers/terraform-provider-alicloud/pull/239))
- Update Go version to 1.10.1 to match new sdk ([#237](https://github.com/terraform-providers/terraform-provider-alicloud/pull/237))
- Support http(s) proxy ([#236](https://github.com/terraform-providers/terraform-provider-alicloud/pull/236))
- Add guide for installing goimports ([#233](https://github.com/terraform-providers/terraform-provider-alicloud/pull/233))
- Improve the makefile and README ([#232](https://github.com/terraform-providers/terraform-provider-alicloud/pull/232))

BUG FIXES:

- Fix losing key pair error after updating ecs instance ([#245](https://github.com/terraform-providers/terraform-provider-alicloud/pull/245))
- Fix BackendServer.configuring error when creating slb rule ([#242](https://github.com/terraform-providers/terraform-provider-alicloud/pull/242))
- Fix bug "...zoneinfo.zip: no such file or directory" happened in windows. ([#238](https://github.com/terraform-providers/terraform-provider-alicloud/pull/238))
- Fix ess_scalingrule InvalidScalingRuleId.NotFound error ([#234](https://github.com/terraform-providers/terraform-provider-alicloud/pull/234))

## 1.10.0 (July 27, 2018)

IMPROVEMENTS:

- Rds supports to create 10.0 PostgreSQL instance. ([#230](https://github.com/terraform-providers/terraform-provider-alicloud/pull/230))
- *New Resource*: _alicloud_fc_trigger_ ([#228](https://github.com/terraform-providers/terraform-provider-alicloud/pull/228))
- *New Resource*: _alicloud_fc_function_ ([#227](https://github.com/terraform-providers/terraform-provider-alicloud/pull/227))
- *New Resource*: _alicloud_fc_service_ 30([#226](https://github.com/terraform-providers/terraform-provider-alicloud/pull/226))
- Support new field 'instance_name' for _alicloud_ots_table_ ([#225](https://github.com/terraform-providers/terraform-provider-alicloud/pull/225))
- *New Resource*: _alicloud_ots_instance_attachment_ ([#224](https://github.com/terraform-providers/terraform-provider-alicloud/pull/224))
- *New Resource*: _alicloud_ots_instance_ ([#223](https://github.com/terraform-providers/terraform-provider-alicloud/pull/223))

BUG FIXES:

- Fix Snat entry not found error ([#229](https://github.com/terraform-providers/terraform-provider-alicloud/pull/229))

## 1.9.6 (July 24, 2018)

IMPROVEMENTS:

- Remove the number limitation of vswitch_ids, slb_ids and db_instance_ids ([#219](https://github.com/terraform-providers/terraform-provider-alicloud/pull/219))
- Reduce test nat gateway cost ([#218](https://github.com/terraform-providers/terraform-provider-alicloud/pull/218))
- Support creating zero-node swarm cluster ([#217](https://github.com/terraform-providers/terraform-provider-alicloud/pull/217))
- Improve security group and rule data source test case ([#216](https://github.com/terraform-providers/terraform-provider-alicloud/pull/216))
- Improve dns record resource test case ([#215](https://github.com/terraform-providers/terraform-provider-alicloud/pull/215))
- Improve test case destroy method ([#214](https://github.com/terraform-providers/terraform-provider-alicloud/pull/214))
- Improve ecs instance resource test case ([#213](https://github.com/terraform-providers/terraform-provider-alicloud/pull/213))
- Improve cdn resource test case ([#212](https://github.com/terraform-providers/terraform-provider-alicloud/pull/212))
- Improve kms resource test case ([#211](https://github.com/terraform-providers/terraform-provider-alicloud/pull/211))
- Improve key pair resource test case ([#210](https://github.com/terraform-providers/terraform-provider-alicloud/pull/210))
- Improve rds resource test case ([#209](https://github.com/terraform-providers/terraform-provider-alicloud/pull/209))
- Improve disk resource test case ([#208](https://github.com/terraform-providers/terraform-provider-alicloud/pull/208))
- Improve eip resource test case ([#207](https://github.com/terraform-providers/terraform-provider-alicloud/pull/207))
- Improve scaling service resource test case ([#206](https://github.com/terraform-providers/terraform-provider-alicloud/pull/206))
- Improve vpc and vswitch resource test case ([#205](https://github.com/terraform-providers/terraform-provider-alicloud/pull/205))
- Improve slb resource test case ([#204](https://github.com/terraform-providers/terraform-provider-alicloud/pull/204))
- Improve security group resource test case ([#203](https://github.com/terraform-providers/terraform-provider-alicloud/pull/203))
- Improve ram resource test case ([#202](https://github.com/terraform-providers/terraform-provider-alicloud/pull/202))
- Improve container cluster resource test case ([#201](https://github.com/terraform-providers/terraform-provider-alicloud/pull/201))
- Improve cloud monitor resource test case ([#200](https://github.com/terraform-providers/terraform-provider-alicloud/pull/200))
- Improve route and router interface resource test case ([#199](https://github.com/terraform-providers/terraform-provider-alicloud/pull/199))
- Improve dns resource test case ([#198](https://github.com/terraform-providers/terraform-provider-alicloud/pull/198))
- Improve oss resource test case ([#197](https://github.com/terraform-providers/terraform-provider-alicloud/pull/197))
- Improve ots table resource test case ([#196](https://github.com/terraform-providers/terraform-provider-alicloud/pull/196))
- Improve nat gateway resource test case ([#195](https://github.com/terraform-providers/terraform-provider-alicloud/pull/195))
- Improve log resource test case ([#194](https://github.com/terraform-providers/terraform-provider-alicloud/pull/194))
- Support changing ecs charge type from Prepaid to PostPaid ([#192](https://github.com/terraform-providers/terraform-provider-alicloud/pull/192))
- Add method to compare json template is equal ([#187](https://github.com/terraform-providers/terraform-provider-alicloud/pull/187))
- Remove useless file ([#191](https://github.com/terraform-providers/terraform-provider-alicloud/pull/191))

BUG FIXES:

- Fix CS kubernetes read error and CS app timeout ([#217](https://github.com/terraform-providers/terraform-provider-alicloud/pull/217))
- Fix getting location connection error ([#193](https://github.com/terraform-providers/terraform-provider-alicloud/pull/193))
- Fix CS kubernetes connection error ([#190](https://github.com/terraform-providers/terraform-provider-alicloud/pull/190))
- Fix Oss bucket diff error ([#189](https://github.com/terraform-providers/terraform-provider-alicloud/pull/189))

NOTES:

- From version 1.9.6, the deprecated resource alicloud_ram_alias file has been removed and the resource has been
replaced by alicloud_ram_account_alias. Details refer to [pull 191](https://github.com/terraform-providers/terraform-provider-alicloud/pull/191/commits/e3fd74591230ccb545bb4309b674d6df33b716b9)

## 1.9.5 (June 20, 2018)

IMPROVEMENTS:

- Improve log machine group docs ([#186](https://github.com/terraform-providers/terraform-provider-alicloud/pull/186))
- Support sts token for some resources ([#185](https://github.com/terraform-providers/terraform-provider-alicloud/pull/185))
- Support user agent for log service ([#184](https://github.com/terraform-providers/terraform-provider-alicloud/pull/184))
- *New Resource*: _alicloud_log_machine_group_ ([#183](https://github.com/terraform-providers/terraform-provider-alicloud/pull/183))
- *New Resource*: _alicloud_log_store_index_ ([#182](https://github.com/terraform-providers/terraform-provider-alicloud/pull/182))
- *New Resource*: _alicloud_log_store_ ([#181](https://github.com/terraform-providers/terraform-provider-alicloud/pull/181))
- *New Resource*: _alicloud_log_project_ ([#180](https://github.com/terraform-providers/terraform-provider-alicloud/pull/180))
- Improve example about cs_kubernetes ([#179](https://github.com/terraform-providers/terraform-provider-alicloud/pull/179))
- Add losing docs about cs_kubernetes ([#178](https://github.com/terraform-providers/terraform-provider-alicloud/pull/178))

## 1.9.4 (June 08, 2018)

IMPROVEMENTS:

- cs_kubernetes supports output worker nodes and master nodes ([#177](https://github.com/terraform-providers/terraform-provider-alicloud/pull/177))
- cs_kubernetes supports to output kube config and certificate ([#176](https://github.com/terraform-providers/terraform-provider-alicloud/pull/176))
- Add a example to deploy mysql and wordpress on kubernetes ([#175](https://github.com/terraform-providers/terraform-provider-alicloud/pull/175))
- Add a example to create swarm and deploy wordpress on it ([#174](https://github.com/terraform-providers/terraform-provider-alicloud/pull/174))
- Change ECS and ESS sdk to official go sdk ([#173](https://github.com/terraform-providers/terraform-provider-alicloud/pull/173))


## 1.9.3 (May 27, 2018)

IMPROVEMENTS:

- *New Data Source*: _alicloud_db_instances_ ([#161](https://github.com/terraform-providers/terraform-provider-alicloud/pull/161))
- Support to set auto renew for ECS instance ([#172](https://github.com/terraform-providers/terraform-provider-alicloud/pull/172))
- Improve cs_kubernetes, slb_listener and db_database docs ([#171](https://github.com/terraform-providers/terraform-provider-alicloud/pull/171))
- Add missing code for describing RDS zones ([#170](https://github.com/terraform-providers/terraform-provider-alicloud/pull/170))
- Add docs notes for windows os([#169](https://github.com/terraform-providers/terraform-provider-alicloud/pull/169))
- Add filter parameters and export parameters for instance types data source. ([#168](https://github.com/terraform-providers/terraform-provider-alicloud/pull/168))
- Add filter parameters for zones data source. ([#167](https://github.com/terraform-providers/terraform-provider-alicloud/pull/167))
- Remove kubernetes work_number limitation ([#165](https://github.com/terraform-providers/terraform-provider-alicloud/pull/165))
- Improve kubernetes examples ([#163](https://github.com/terraform-providers/terraform-provider-alicloud/pull/163))

BUG FIXES:

- Fix getting some instance types failed bug ([#166](https://github.com/terraform-providers/terraform-provider-alicloud/pull/166))
- Fix kubernetes out range index error ([#164](https://github.com/terraform-providers/terraform-provider-alicloud/pull/164))

## 1.9.2 (May 09, 2018)

IMPROVEMENTS:

- *New Resource*: _alicloud_ots_table_ ([#162](https://github.com/terraform-providers/terraform-provider-alicloud/pull/162))
- Fix SLB listener "OperationBusy" error ([#159](https://github.com/terraform-providers/terraform-provider-alicloud/pull/159))
- Prolong waiting time for creating kubernetes cluster to avoid timeout ([#158](https://github.com/terraform-providers/terraform-provider-alicloud/pull/158))
- Support load endpoint from environment variable or specified file ([#157](https://github.com/terraform-providers/terraform-provider-alicloud/pull/157))
- Update example ([#155](https://github.com/terraform-providers/terraform-provider-alicloud/pull/155))

BUG FIXES:

- Fix modifying instance host name failed bug ([#160](https://github.com/terraform-providers/terraform-provider-alicloud/pull/160))
- Fix SLB listener "OperationBusy" error ([#159](https://github.com/terraform-providers/terraform-provider-alicloud/pull/159))
- Fix deleting forward table not found error ([#154](https://github.com/terraform-providers/terraform-provider-alicloud/pull/154))
- Fix deleting slb listener error ([#150](https://github.com/terraform-providers/terraform-provider-alicloud/pull/150))
- Fix creating vswitch error ([#149](https://github.com/terraform-providers/terraform-provider-alicloud/pull/149))

## 1.9.1 (April 13, 2018)

IMPROVEMENTS:

- *New Resource*: _alicloud_cms_alarm_ ([#146](https://github.com/terraform-providers/terraform-provider-alicloud/pull/146))
- *New Resource*: _alicloud_cs_application_ ([#136](https://github.com/terraform-providers/terraform-provider-alicloud/pull/136))
- *New Datasource*: _alicloud_security_group_rules_ ([#135](https://github.com/terraform-providers/terraform-provider-alicloud/pull/135))
- Output application attribution service block ([#141](https://github.com/terraform-providers/terraform-provider-alicloud/pull/141))
- Output swarm attribution 'vpc_id' ([#140](https://github.com/terraform-providers/terraform-provider-alicloud/pull/140))
- Support to release eip after deploying swarm cluster. ([#139](https://github.com/terraform-providers/terraform-provider-alicloud/pull/139))
- Output swarm and kubernetes's nodes information and other attribution ([#138](https://github.com/terraform-providers/terraform-provider-alicloud/pull/138))
- Modify `size` to `node_number` ([#137](https://github.com/terraform-providers/terraform-provider-alicloud/pull/137))
- Set swarm ID before waiting its status ([#134](https://github.com/terraform-providers/terraform-provider-alicloud/pull/134))
- Add 'is_outdated' for cs_swarm and cs_kubernetes ([#133](https://github.com/terraform-providers/terraform-provider-alicloud/pull/133))
- Add warning when creating postgresql and ppas database ([#132](https://github.com/terraform-providers/terraform-provider-alicloud/pull/132))
- Add kubernetes example ([#142](https://github.com/terraform-providers/terraform-provider-alicloud/pull/142))
- Update sdk to support user-agent ([#143](https://github.com/terraform-providers/terraform-provider-alicloud/pull/143))
- Add eip unassociation retry times to avoid needless error ([#144](https://github.com/terraform-providers/terraform-provider-alicloud/pull/144))
- Add connections output for kubernetes cluster ([#145](https://github.com/terraform-providers/terraform-provider-alicloud/pull/145))

BUG FIXES:

- Fix vpc not found when vpc has been deleted ([#131](https://github.com/terraform-providers/terraform-provider-alicloud/pull/131))


## 1.9.0 (March 19, 2018)

IMPROVEMENTS:

- *New Resource*: _alicloud_cs_kubernetes_ ([#129](https://github.com/terraform-providers/terraform-provider-alicloud/pull/129))
- *New DataSource*: _alicloud_eips_ ([#123](https://github.com/terraform-providers/terraform-provider-alicloud/pull/123))
- Add server_group_id to slb listener resource ([#122](https://github.com/terraform-providers/terraform-provider-alicloud/pull/122))
- Rename _alicloud_container_cluster_ to _alicloud_cs_swarm_ ([#128](https://github.com/terraform-providers/terraform-provider-alicloud/pull/128))

BUG FIXES:

- Fix vpc description validate ([#125](https://github.com/terraform-providers/terraform-provider-alicloud/pull/125))
- Update SDK version to fix unresolving endpoint issue ([#126](https://github.com/terraform-providers/terraform-provider-alicloud/pull/126))
- Add waiting time after ECS bind ECS to ensure network is ok ([#127](https://github.com/terraform-providers/terraform-provider-alicloud/pull/127))

## 1.8.1 (March 09, 2018)

IMPROVEMENTS:

- DB instance supports multiple zone ([#120](https://github.com/terraform-providers/terraform-provider-alicloud/pull/120))
- Data source zones support to retrieve multiple zone ([#119](https://github.com/terraform-providers/terraform-provider-alicloud/pull/119))
- VPC supports alibaba cloud official go sdk ([#118](https://github.com/terraform-providers/terraform-provider-alicloud/pull/118))

BUG FIXES:

- Fix not found db instance bug when allocating connection ([#121](https://github.com/terraform-providers/terraform-provider-alicloud/pull/121))


## 1.8.0 (March 02, 2018)

IMPROVEMENTS:

- Support golang version 1.9 ([#114](https://github.com/terraform-providers/terraform-provider-alicloud/pull/114))
- RDS supports alibaba cloud official go sdk ([#113](https://github.com/terraform-providers/terraform-provider-alicloud/pull/113))
- Deprecated 'in_use' in eips datasource to fix conflict ([#115](https://github.com/terraform-providers/terraform-provider-alicloud/pull/115))
- Add encrypted argument to alicloud_disk resource（[#116](https://github.com/terraform-providers/terraform-provider-alicloud/pull/116))

BUG FIXES:

- Fix reading router interface failed bug ([#117](https://github.com/terraform-providers/terraform-provider-alicloud/pull/117))

## 1.7.2 (February 09, 2018)

IMPROVEMENTS:

- *New DataSource*: _alicloud_eips_ ([#110](https://github.com/terraform-providers/terraform-provider-alicloud/pull/110))
- *New DataSource*: _alicloud_vswitches_ ([#109](https://github.com/terraform-providers/terraform-provider-alicloud/pull/109))
- Support inner network segregation in one security group ([#112](https://github.com/terraform-providers/terraform-provider-alicloud/pull/112))

BUG FIXES:

- Fix creating Classic instance failed result in role_name ([#111](https://github.com/terraform-providers/terraform-provider-alicloud/pull/111))
- Fix eip is not exist in nat gateway when creating snat ([#108](https://github.com/terraform-providers/terraform-provider-alicloud/pull/108))

## 1.7.1 (February 02, 2018)

IMPROVEMENTS:

- Support setting instance_name for ESS scaling configuration ([#107](https://github.com/terraform-providers/terraform-provider-alicloud/pull/107))
- Support multiple vswitches for ESS scaling group and output slbIds and dbIds ([#105](https://github.com/terraform-providers/terraform-provider-alicloud/pull/105))
- Support to set internet_max_bandwidth_out is 0 for ESS configuration ([#103](https://github.com/terraform-providers/terraform-provider-alicloud/pull/103))
- Modify EIP default to PayByTraffic for international account ([#101](https://github.com/terraform-providers/terraform-provider-alicloud/pull/101))
- Deprecate nat gateway fileds 'spec' and 'bandwidth_packages' ([#100](https://github.com/terraform-providers/terraform-provider-alicloud/pull/100))
- Support to associate EIP with SLB and Nat Gateway ([#99](https://github.com/terraform-providers/terraform-provider-alicloud/pull/99))

BUG FIXES:

- fix a bug that can't create multiple VPC, vswitch and nat gateway at one time ([#102](https://github.com/terraform-providers/terraform-provider-alicloud/pull/102))
- fix a bug that can't import instance 'role_name' ([#104](https://github.com/terraform-providers/terraform-provider-alicloud/pull/104))
- fix a bug that creating ESS scaling group and configuration results from 'Throttling' ([#106](https://github.com/terraform-providers/terraform-provider-alicloud/pull/106))

## 1.7.0 (January 25, 2018)

IMPROVEMENTS:

- *New Resource*: _alicloud_kms_key_ ([#91](https://github.com/terraform-providers/terraform-provider-alicloud/pull/91))
- *New DataSource*: _alicloud_kms_keys_ ([#93](https://github.com/terraform-providers/terraform-provider-alicloud/pull/93))
- *New DataSource*: _alicloud_instances_ ([#94](https://github.com/terraform-providers/terraform-provider-alicloud/pull/94))
- Add a new output field "arn" for _alicloud_kms_key_ ([#92](https://github.com/terraform-providers/terraform-provider-alicloud/pull/92))
- Add a new field "specification" for _alicloud_slb_ ([#95](https://github.com/terraform-providers/terraform-provider-alicloud/pull/95))
- Improve security group rule's port range for "-1/-1" ([#96](https://github.com/terraform-providers/terraform-provider-alicloud/pull/96))

BUG FIXES:

- fix slb invalid status error when launching ESS scaling group ([#97](https://github.com/terraform-providers/terraform-provider-alicloud/pull/97))

## 1.6.2 (January 22, 2018)

IMPROVEMENTS:

- Modify db_connection prefix default value to "instance_id + 'tf'"([#90](https://github.com/terraform-providers/terraform-provider-alicloud/pull/90))
- Modify db_connection ID to make it more simple while importing it([#90](https://github.com/terraform-providers/terraform-provider-alicloud/pull/90))
- Add wait method to avoid useless status error while creating/modifying account or privilege or connection or database([#90](https://github.com/terraform-providers/terraform-provider-alicloud/pull/90))
- Support to set instnace name for RDS ([#88](https://github.com/terraform-providers/terraform-provider-alicloud/pull/88))
- Avoid container cluster cidr block conflicts with vswitch's ([#88](https://github.com/terraform-providers/terraform-provider-alicloud/pull/88))
- Output resource import information ([#87](https://github.com/terraform-providers/terraform-provider-alicloud/pull/87))

BUG FIXES:

- fix instance id not found and instane status not supported bug([#90](https://github.com/terraform-providers/terraform-provider-alicloud/pull/90))
- fix deleting slb_attachment resource failed bug ([#86](https://github.com/terraform-providers/terraform-provider-alicloud/pull/86))


## 1.6.1 (January 18, 2018)

IMPROVEMENTS:

- Support to modify instance type and network spec ([#84](https://github.com/terraform-providers/terraform-provider-alicloud/pull/84))
- Avoid needless error when creating security group rule ([#83](https://github.com/terraform-providers/terraform-provider-alicloud/pull/83))

BUG FIXES:

- fix creating cluster container failed bug ([#85](https://github.com/terraform-providers/terraform-provider-alicloud/pull/85))


## 1.6.0 (January 15, 2018)

IMPROVEMENTS:

- *New Resource*: _alicloud_ess_attachment_ ([#80](https://github.com/terraform-providers/terraform-provider-alicloud/pull/80))
- *New Resource*: _alicloud_slb_rule_ ([#79](https://github.com/terraform-providers/terraform-provider-alicloud/pull/79))
- *New Resource*: _alicloud_slb_server_group_ ([#78](https://github.com/terraform-providers/terraform-provider-alicloud/pull/78))
- Support Spot Instance ([#77](https://github.com/terraform-providers/terraform-provider-alicloud/pull/77))
- Output tip message when international account create SLB failed ([#75](https://github.com/terraform-providers/terraform-provider-alicloud/pull/75))
- Standardize the order of imports packages ([#74](https://github.com/terraform-providers/terraform-provider-alicloud/pull/74))
- Add "weight" for slb_attachment to improve the resource ([#81](https://github.com/terraform-providers/terraform-provider-alicloud/pull/81))

BUG FIXES:

- fix allocating RDS public connection conflict error ([#76](https://github.com/terraform-providers/terraform-provider-alicloud/pull/76))

## 1.5.3 (January 9, 2018)

BUG FIXES:
  * fix getting OSS endpoint failed error  ([#73](https://github.com/terraform-providers/terraform-provider-alicloud/pull/73))
  * fix describing dns record not found when deleting record ([#73](https://github.com/terraform-providers/terraform-provider-alicloud/pull/73))

## 1.5.2 (January 8, 2018)

BUG FIXES:
  * fix creating rds 'Prepaid' instance failed error  ([#70](https://github.com/terraform-providers/terraform-provider-alicloud/pull/70))

## 1.5.1 (January 5, 2018)

BUG FIXES:
  * modify security_token to Optional ([#69](https://github.com/terraform-providers/terraform-provider-alicloud/pull/69))

## 1.5.0 (January 4, 2018)

IMPROVEMENTS:

- *New Resource*: _alicloud_db_database_ ([#68](https://github.com/terraform-providers/terraform-provider-alicloud/pull/68))
- *New Resource*: _alicloud_db_backup_policy_ ([#68](https://github.com/terraform-providers/terraform-provider-alicloud/pull/68))
- *New Resource*: _alicloud_db_connection_ ([#67](https://github.com/terraform-providers/terraform-provider-alicloud/pull/67))
- *New Resource*: _alicloud_db_account_ ([#66](https://github.com/terraform-providers/terraform-provider-alicloud/pull/66))
- *New Resource*: _alicloud_db_account_privilege_ ([#66](https://github.com/terraform-providers/terraform-provider-alicloud/pull/66))
- resource/db_instance: remove some field to new resource ([#65](https://github.com/terraform-providers/terraform-provider-alicloud/pull/65))
- resource/instance: support to modify private ip, vswitch_id and instance charge type ([#65](https://github.com/terraform-providers/terraform-provider-alicloud/pull/65))

BUG FIXES:

- resource/dns-record: Fix dns record still exist after deleting it ([#65](https://github.com/terraform-providers/terraform-provider-alicloud/pull/65))
- resource/instance: fix deleting route entry error ([#69](https://github.com/terraform-providers/terraform-provider-alicloud/pull/69))


## 1.2.0 (December 15, 2017)

IMPROVEMENTS:
- resource/slb: wait for SLB active before return back ([#61](https://github.com/terraform-providers/terraform-provider-alicloud/pull/61))

BUG FIXES:

- resource/dns-record: Fix setting dns priority failed ([#58](https://github.com/terraform-providers/terraform-provider-alicloud/pull/58))
- resource/dns-record: Fix ESS attachs SLB failed ([#59](https://github.com/terraform-providers/terraform-provider-alicloud/pull/59))
- resource/dns-record: Fix security group not found error ([#59](https://github.com/terraform-providers/terraform-provider-alicloud/pull/59))


## 1.0.0 (December 11, 2017)

IMPROVEMENTS:

- *New Resource*: _alicloud_slb_listener_ ([#53](https://github.com/terraform-providers/terraform-provider-alicloud/pull/53))
- *New Resource*: _alicloud_cdn_domain_ ([#52](https://github.com/terraform-providers/terraform-provider-alicloud/pull/52))
- *New Resource*: _alicloud_dns_ ([#51](https://github.com/terraform-providers/terraform-provider-alicloud/pull/51))
- *New Resource*: _alicloud_dns_group_ ([#51](https://github.com/terraform-providers/terraform-provider-alicloud/pull/51))
- *New Resource*: _alicloud_dns_record_ ([#51](https://github.com/terraform-providers/terraform-provider-alicloud/pull/51))
- *New Resource*: _alicloud_ram_account_alias_ ([#50](https://github.com/terraform-providers/terraform-provider-alicloud/pull/50))
- *New Resource*: _alicloud_ram_login_profile_ ([#50](https://github.com/terraform-providers/terraform-provider-alicloud/pull/50))
- *New Resource*: _alicloud_ram_access_key_ ([#50](https://github.com/terraform-providers/terraform-provider-alicloud/pull/50))
- *New Resource*: _alicloud_ram_group_ ([#49](https://github.com/terraform-providers/terraform-provider-alicloud/pull/49))
- *New Resource*: _alicloud_ram_group_membership_ ([#49](https://github.com/terraform-providers/terraform-provider-alicloud/pull/49))
- *New Resource*: _alicloud_ram_group_policy_attachment_ ([#49](https://github.com/terraform-providers/terraform-provider-alicloud/pull/49))
- *New Resource*: _alicloud_ram_role_ ([#48](https://github.com/terraform-providers/terraform-provider-alicloud/pull/48))
- *New Resource*: _alicloud_ram_role_attachment_ ([#48](https://github.com/terraform-providers/terraform-provider-alicloud/pull/48))
- *New Resource*: _alicloud_ram_role_polocy_attachment_ ([#48](https://github.com/terraform-providers/terraform-provider-alicloud/pull/48))
- *New Resource*: _alicloud_container_cluster_ ([#47](https://github.com/terraform-providers/terraform-provider-alicloud/pull/47))
- *New Resource:* _alicloud_ram_policy_ ([#46](https://github.com/terraform-providers/terraform-provider-alicloud/pull/46))
- *New Resource*: _alicloud_ram_user_policy_attachment_ ([#46](https://github.com/terraform-providers/terraform-provider-alicloud/pull/46))
- *New Resource* _alicloud_ram_user_ ([#44](https://github.com/terraform-providers/terraform-provider-alicloud/pull/44))
- *New Datasource* _alicloud_ram_policies_ ([#46](https://github.com/terraform-providers/terraform-provider-alicloud/pull/46))
- *New Datasource* _alicloud_ram_users_ ([#44](https://github.com/terraform-providers/terraform-provider-alicloud/pull/44))
- *New Datasource*: _alicloud_ram_roles_ ([#48](https://github.com/terraform-providers/terraform-provider-alicloud/pull/48))
- *New Datasource*: _alicloud_ram_account_aliases_ ([#50](https://github.com/terraform-providers/terraform-provider-alicloud/pull/50))
- *New Datasource*: _alicloud_dns_domains_ ([#51](https://github.com/terraform-providers/terraform-provider-alicloud/pull/51))
- *New Datasource*: _alicloud_dns_groups_ ([#51](https://github.com/terraform-providers/terraform-provider-alicloud/pull/51))
- *New Datasource*: _alicloud_dns_records_ ([#51](https://github.com/terraform-providers/terraform-provider-alicloud/pull/51))
- resource/instance: add new parameter `role_name` ([#48](https://github.com/terraform-providers/terraform-provider-alicloud/pull/48))
- resource/slb: remove slb schema field `listeners` and using new listener resource to replace ([#55](https://github.com/terraform-providers/terraform-provider-alicloud/pull/55))
- resource/ess_scaling_configuration: add new parameters `key_name`, `role_name`, `user_data`, `force_delete` and `tags` ([#54](https://github.com/terraform-providers/terraform-provider-alicloud/pull/54))
- resource/ess_scaling_configuration: remove it importing ([#54](https://github.com/terraform-providers/terraform-provider-alicloud/pull/54))
- resource: format not found error ([#55](https://github.com/terraform-providers/terraform-provider-alicloud/pull/55))
- website: improve resource docs ([#56](https://github.com/terraform-providers/terraform-provider-alicloud/pull/56))
- examples: add new examples, like oss, key_pair, router_interface and so on ([#56](https://github.com/terraform-providers/terraform-provider-alicloud/pull/56))

- Added support for importing:
  - `alicloud_container_cluster` ([#47](https://github.com/terraform-providers/terraform-provider-alicloud/pull/47))
  - `alicloud_ram_policy` ([#46](https://github.com/terraform-providers/terraform-provider-alicloud/pull/46))
  - `alicloud_ram_user` ([#44](https://github.com/terraform-providers/terraform-provider-alicloud/pull/44))
  - `alicloud_ram_role` ([#48](https://github.com/terraform-providers/terraform-provider-alicloud/pull/48))
  - `alicloud_ram_groups` ([#49](https://github.com/terraform-providers/terraform-provider-alicloud/pull/49))
  - `alicloud_ram_login_profile` ([#50](https://github.com/terraform-providers/terraform-provider-alicloud/pull/50))
  - `alicloud_dns` ([#51](https://github.com/terraform-providers/terraform-provider-alicloud/pull/51))
  - `alicloud_dns_record` ([#51](https://github.com/terraform-providers/terraform-provider-alicloud/pull/51))
  - `alicloud_slb_listener` ([#53](https://github.com/terraform-providers/terraform-provider-alicloud/pull/53))
  - `alicloud_security_group` ([#55](https://github.com/terraform-providers/terraform-provider-alicloud/pull/55))
  - `alicloud_slb` ([#55](https://github.com/terraform-providers/terraform-provider-alicloud/pull/55))
  - `alicloud_vswitch` ([#55](https://github.com/terraform-providers/terraform-provider-alicloud/pull/55))
  - `alicloud_vroute_entry` ([#55](https://github.com/terraform-providers/terraform-provider-alicloud/pull/55))

BUG FIXES:

- resource/vroute_entry: Fix building route_entry concurrency issue ([#55](https://github.com/terraform-providers/terraform-provider-alicloud/pull/55))
- resource/vswitch: Fix building vswitch concurrency issue ([#55](https://github.com/terraform-providers/terraform-provider-alicloud/pull/55))
- resource/router_interface: Fix building router interface concurrency issue ([#55](https://github.com/terraform-providers/terraform-provider-alicloud/pull/55))
- resource/vpc: Fix building vpc concurrency issue ([#55](https://github.com/terraform-providers/terraform-provider-alicloud/pull/55))
- resource/slb_attachment: Fix attaching slb failed ([#55](https://github.com/terraform-providers/terraform-provider-alicloud/pull/55))

## 0.1.1 (December 11, 2017)

IMPROVEMENTS:

- *New Resource:* _alicloud_key_pair_ ([#27](https://github.com/terraform-providers/terraform-provider-alicloud/pull/27))
- *New Resource*: _alicloud_key_pair_attachment_ ([#28](https://github.com/terraform-providers/terraform-provider-alicloud/pull/28))
- *New Resource*: _alicloud_router_interface_ ([#40](https://github.com/terraform-providers/terraform-provider-alicloud/pull/40))
- *New Resource:* _alicloud_oss_bucket_ ([#10](https://github.com/terraform-providers/terraform-provider-alicloud/pull/10))
- *New Resource*: _alicloud_oss_bucket_object_ ([#14](https://github.com/terraform-providers/terraform-provider-alicloud/pull/14))
- *New Datasource* _alicloud_key_pairs_ ([#30](https://github.com/terraform-providers/terraform-provider-alicloud/pull/30))
- *New Datasource* _alicloud_vpcs_ ([#34](https://github.com/terraform-providers/terraform-provider-alicloud/pull/34))
- *New output_file* option for data sources: export data to a specified file ([#29](https://github.com/terraform-providers/terraform-provider-alicloud/pull/29))
- resource/instance:add new parameter `key_name` ([#31](https://github.com/terraform-providers/terraform-provider-alicloud/pull/31))
- resource/route_entry: new nexthop type 'RouterInterface' for route entry ([#41](https://github.com/terraform-providers/terraform-provider-alicloud/pull/41))
- resource/security_group_rule: Remove `cidr_ip` contribute "ConflictsWith" ([#39](https://github.com/terraform-providers/terraform-provider-alicloud/pull/39))
- resource/rds: add ability to change instance password ([#17](https://github.com/terraform-providers/terraform-provider-alicloud/pull/17))
- resource/rds: Add ability to import existing RDS resources ([#16](https://github.com/terraform-providers/terraform-provider-alicloud/pull/16))
- datasource/alicloud_zones: Add more options for filtering ([#19](https://github.com/terraform-providers/terraform-provider-alicloud/pull/19))
- Added support for importing:
  - `alicloud_vpc` ([#32](https://github.com/terraform-providers/terraform-provider-alicloud/pull/32))
  - `alicloud_route_entry` ([#33](https://github.com/terraform-providers/terraform-provider-alicloud/pull/33))
  - `alicloud_nat_gateway` ([#26](https://github.com/terraform-providers/terraform-provider-alicloud/pull/26))
  - `alicloud_ess_schedule` ([#25](https://github.com/terraform-providers/terraform-provider-alicloud/pull/25))
  - `alicloud_ess_scaling_group` ([#24](https://github.com/terraform-providers/terraform-provider-alicloud/pull/24))
  - `alicloud_instance` ([#23](https://github.com/terraform-providers/terraform-provider-alicloud/pull/23))
  - `alicloud_eip` ([#22](https://github.com/terraform-providers/terraform-provider-alicloud/pull/22))
  - `alicloud_disk` ([#21](https://github.com/terraform-providers/terraform-provider-alicloud/pull/21))

BUG FIXES:

- resource/disk_attachment: Fix issue attaching multiple disks and set disk_attachment's parameter 'device_name' as deprecated ([#9](https://github.com/terraform-providers/terraform-provider-alicloud/pull/9))
- resource/rds: Fix diff error about rds security_ips ([#13](https://github.com/terraform-providers/terraform-provider-alicloud/pull/13))
- resource/security_group_rule: Fix diff error when authorizing security group rules ([#15](https://github.com/terraform-providers/terraform-provider-alicloud/pull/15))
- resource/security_group_rule: Fix diff bug by modifying 'DestCidrIp' to 'DestGroupId' when running read ([#35](https://github.com/terraform-providers/terraform-provider-alicloud/pull/35))


## 0.1.0 (June 20, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
