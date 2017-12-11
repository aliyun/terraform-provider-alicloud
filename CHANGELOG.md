## 1.0.0 (Unreleased)

IMPROVEMENTS:

- *New Resource*: _alicloud_slb_listener_ [GH-53]
- *New Resource*: _alicloud_cdn_domain_ [GH-52]
- *New Resource*: _alicloud_dns_ [GH-51]
- *New Resource*: _alicloud_dns_group_ [GH-51]
- *New Resource*: _alicloud_dns_record_ [GH-51]
- *New Resource*: _alicloud_ram_account_alias_ [GH-50]
- *New Resource*: _alicloud_ram_login_profile_ [GH-50]
- *New Resource*: _alicloud_ram_access_key_ [GH-50]
- *New Resource*: _alicloud_ram_group_ [GH-49]
- *New Resource*: _alicloud_ram_group_membership_ [GH-49]
- *New Resource*: _alicloud_ram_group_policy_attachment_ [GH-49]
- *New Resource*: _alicloud_ram_role_ [GH-48]
- *New Resource*: _alicloud_ram_role_attachment_ [GH-48]
- *New Resource*: _alicloud_ram_role_polocy_attachment_ [GH-48]
- *New Resource*: _alicloud_container_cluster_ [GH-47]
- *New Resource:* _alicloud_ram_policy_ [GH-46]
- *New Resource*: _alicloud_ram_user_policy_attachment_ [GH-46]
- *New Resource* _alicloud_ram_user_ [GH-44]
- *New Datasource* _alicloud_ram_policies_ [GH-46]
- *New Datasource* _alicloud_ram_users_ [GH-44]
- *New Datasource*: _alicloud_ram_roles_ [GH-48]
- *New Datasource*: _alicloud_ram_account_aliases_ [GH-50]
- *New Datasource*: _alicloud_dns_domains_ [GH-51]
- *New Datasource*: _alicloud_dns_groups_ [GH-51]
- *New Datasource*: _alicloud_dns_records_ [GH-51]
- resource/instance: add new parameter `role_name` [GH-48]
- resource/slb: remove slb schema field `listeners` and using new listener resource to replace [GH-55]
- resource/ess_scaling_configuration: add new parameters `key_name`, `role_name`, `user_data`, `force_delete` and `tags` [GH-54]
- resource/ess_scaling_configuration: remove it importing [GH-54]
- resource: format not found error [GH-55]
- website: improve resource docs [GH-56]
- examples: add new examples, like oss, key_pair, router_interface and so on [GH-56]

- Added support for importing:
  - `alicloud_container_cluster` [GH-47]
  - `alicloud_ram_policy` [GH-46]
  - `alicloud_ram_user` [GH-44]
  - `alicloud_ram_role` [GH-48]
  - `alicloud_ram_groups` [GH-49]
  - `alicloud_ram_login_profile` [GH-50]
  - `alicloud_dns` [GH-51]
  - `alicloud_dns_record` [GH-51]
  - `alicloud_slb_listener` [GH-53]
  - `alicloud_security_group` [GH-55]
  - `alicloud_slb` [GH-55]
  - `alicloud_vswitch` [GH-55]
  - `alicloud_vroute_entry` [GH-55]

BUG FIXES:

- resource/vroute_entry: Fix building route_entry concurrency issue [GH-55]
- resource/vswitch: Fix building vswitch concurrency issue [GH-55]
- resource/router_interface: Fix building router interface concurrency issue [GH-55]
- resource/vpc: Fix building vpc concurrency issue [GH-55]
- resource/slb_attachment: Fix attaching slb failed [GH-55]

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