## 1.38.0 (Unreleased)

FEATURES:

- **New Resource:** `alicloud_ddoscoo_instance` [GH-952]

IMPROVEMENTS:

- improve ram resource testcases [GH-961]
- ecs prepaid instance supports changing instance type [GH-949]
- update mongodb instance test case for multiAZ [GH-947]
- improve ram_policy resource update method [GH-960]

BUG FIXES:

- Fix drds instance sweeper test bug [GH-955]

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
