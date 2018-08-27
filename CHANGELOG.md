## 1.13.0 (Unreleased)

FEATURES:

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
