---
subcategory: "DMS Enterprise"
layout: "alicloud"
page_title: "Alicloud: alicloud_dms_enterprise_instances"
sidebar_current: "docs-alicloud-datasource-dms-enterprise-instances"
description: |-
    Provides a list of available DMS Enterprise Instances.
---

# alicloud\_dms\_enterprise\_instances

This data source provides a list of DMS Enterprise Instances in an Alibaba Cloud account according to the specified filters.

-> **NOTE:** Available in 1.88.0+

## Example Usage

```terraform
# Declare the data source
data "alicloud_dms_enterprise_instances" "dms_enterprise_instances_ds" {
  net_type      = "CLASSIC"
  instance_type = "mysql"
  env_type      = "test"
  name_regex    = "tf_testAcc"
  output_file   = "dms_enterprise_instances.json"
}

output "first_database_instance_id" {
  value = "${data.alicloud_dms_enterprise_instances.dms_enterprise_instances_ds.instances.0.instance_id}"
}
```

## Argument Reference

The following arguments are supported:

* `status` - (Optional) Filter the results by status of the DMS Enterprise Instances. Valid values: `NORMAL`, `UNAVAILABLE`, `UNKNOWN`, `DELETED`, `DISABLE`.
* `env_type` - (Optional) The type of the environment to which the database instance belongs.
* `instance_source` - (Optional) The source of the database instance.
* `instance_state` - (Optional) The status of the database instance.
* `net_type` - (Optional) The network type of the database instance. Valid values: CLASSIC and VPC. For more information about the valid values, see the description of the RegisterInstance operation.
* `search_key` - (Optional) The keyword used to query database instances.
* `tid` - (Optional) The ID of the tenant in Data Management (DMS) Enterprise.
* `name_regex` - (Optional, Available in v1.100.0+) A regex string to filter the results by the DMS Enterprise Instance instance_alias.
* `instance_alias_regex` - (Optional) A regex string to filter the results by the DMS Enterprise Instance instance_alias.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of DMS Enterprise IDs (Each of them consists of host:port).
* `names` - A list of DMS Enterprise names.
* `instances` - A list of KMS keys. Each element contains the following attributes:
  * `data_link_name` - The name of the data link for the database instance.
  * `database_password` - The logon password of the database instance.
  * `database_user` - The logon username of the database instance.
  * `dba_id` - The ID of the database administrator (DBA) of the database instance.
  * `dba_nick_name` - The nickname of the DBA.
  * `ddl_online` - Indicates whether the online data description language (DDL) service was enabled for the database instance.
  * `ecs_instance_id` - The ID of the Elastic Compute Service (ECS) instance to which the database instance belongs.
  * `ecs_region` - The region where the database instance resides.
  * `env_type` - The type of the environment to which the database instance belongs..
  * `export_timeout` - The timeout period for exporting the database instance.
  * `host` - The endpoint of the database instance.
  * `instance_alias` - The alias of the database instance.
  * `instance_id` - The ID of the database instance.
  * `instance_source` - The ID of the database instance.
  * `instance_type` - The ID of the database instance.
  * `port` - The connection port of the database instance.
  * `query_timeout` - The timeout period for querying the database instance.
  * `safe_rule_id` - The ID of the security rule for the database instance.
  * `sid` - The system ID (SID) of the database instance.
  * `status` - The status of the database instance.
  * `use_dsql` - Indicates whether cross-database query was enabled for the database instance.
  * `vpc_id` - The ID of the Virtual Private Cloud (VPC) to which the database instance belongs.
