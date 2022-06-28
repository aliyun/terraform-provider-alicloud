---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_db_instances"
sidebar_current: "docs-alicloud-datasource-db-instances"
description: |-
    Provides a collection of RDS instances according to the specified filters.
---

# alicloud\_db\_instances

The `alicloud_db_instances` data source provides a collection of RDS instances available in Alibaba Cloud account.
Filters support regular expression for the instance name, searches by tags, and other filters which are listed below.

## Example Usage

```
data "alicloud_db_instances" "db_instances_ds" {
  name_regex = "data-\\d+"
  status     = "Running"
  tags       = {
    "type" = "database",
    "size" = "tiny"
  }

}

output "first_db_instance_id" {
  value = "${data.alicloud_db_instances.db_instances_ds.instances.0.id}"
}
```

## Argument Reference

The following arguments are supported:
* `enable_details` - (Optional, Available in 1.135.0+) Default to `false`. Set it to `true` can output parameter template about resource attributes.
* `name_regex` - (Optional) A regex string to filter results by instance name.
* `ids` - (Optional, Available 1.52.0+) A list of RDS instance IDs. 
* `engine` - (Optional) Database type. Options are `MySQL`, `SQLServer`, `PostgreSQL` and `PPAS`. If no value is specified, all types are returned.
* `status` - (Optional) Status of the instance.
* `db_type` - (Optional) `Primary` for primary instance, `Readonly` for read-only instance, `Guard` for disaster recovery instance, and `Temp` for temporary instance.
* `vpc_id` - (Optional) Used to retrieve instances belong to specified VPC.
* `vswitch_id` - (Optional) Used to retrieve instances belong to specified `vswitch` resources.
* `connection_mode` - (Optional) `Standard` for standard access mode and `Safe` for high security access mode.
* `tags` - (Optional) A map of tags assigned to the DB instances. 
Note: Before 1.60.0, the value's format is a `json` string which including `TagKey` and `TagValue`. `TagKey` cannot be null, and `TagValue` can be empty. Format example `"{\"key1\":\"value1\"}"`
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of RDS instance IDs. 
* `names` - A list of RDS instance names. 
* `instances` - A list of RDS instances. Each element contains the following attributes:
  * `id` - The ID of the RDS instance.
  * `name` - The name of the RDS instance.
  * `charge_type` - Billing method. Value options: `Postpaid` for Pay-As-You-Go and `Prepaid` for subscription.
  * `db_type` - `Primary` for primary instance, `Readonly` for read-only instance, `Guard` for disaster recovery instance, and `Temp` for temporary instance.
  * `region_id` - Region ID the instance belongs to.
  * `create_time` - Creation time of the instance.
  * `expire_time` - Expiration time. Pay-As-You-Go instances never expire.
  * `status` - Status of the instance.
  * `engine` - Database type. Options are `MySQL`, `SQLServer`, `PostgreSQL` and `PPAS`. If no value is specified, all types are returned.
  * `engine_version` - Database version.
  * `net_type` - `Internet` for public network or `Intranet` for private network.
  * `connection_mode` - `Standard` for standard access mode and `Safe` for high security access mode.
  * `instance_type` - Sizing of the RDS instance.
  * `availability_zone` - Availability zone.
  * `master_instance_id` - ID of the primary instance. If this parameter is not returned, the current instance is a primary instance.
  * `guard_instance_id` - If a disaster recovery instance is attached to the current instance, the ID of the disaster recovery instance applies.
  * `temp_instance_id` - If a temporary instance is attached to the current instance, the ID of the temporary instance applies.
  * `readonly_instance_ids` - A list of IDs of read-only instances attached to the primary instance.
  * `vpc_id` - ID of the VPC the instance belongs to.
  * `vswitch_id` - ID of the VSwitch the instance belongs to.
  * `port` - (Available in 1.70.3+) RDS database connection port.
  * `connection_string` - (Available in 1.70.3+) RDS database connection string.
  * `instance_storage` - (Available in 1.70.3+) User-defined DB instance storage space.
  * `db_instance_storage_type` - (Available in 1.70.3+) The storage type of the instance.
  * `master_zone` - (Available in 1.101.0+) The master zone of the instance.
  * `zone_id_slave_a` - (Available in 1.101.0+) The region ID of the secondary instance if you create a secondary instance. If you set this parameter to the same value as the ZoneId parameter, the instance is deployed in a single zone. Otherwise, the instance is deployed in multiple zones. 
  * `zone_id_slave_b` - (Available in 1.101.0+) The region ID of the log instance if you create a log instance. If you set this parameter to the same value as the ZoneId parameter, the instance is deployed in a single zone. Otherwise, the instance is deployed in multiple zones.
  * `ssl_expire_time` - (Available in 1.124.1+) The time when the server certificate expires. The time follows the ISO 8601 standard in the yyyy-MM-ddTHH:mm:ssZ format. The time is displayed in UTC.
  * `require_update` - (Available in 1.124.1+) Indicates whether the server certificate needs to be updated.
    - Valid values for ApsaraDB RDS for MySQL and ApsaraDB RDS for SQL Server:
      - No
      - Yes
    - Valid values for ApsaraDB RDS for PostgreSQL:
      - 0: no
      - 1: yes
  * `acl` - (Available in 1.124.1+) The method that is used to verify the identities of clients. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs. Valid values:
      - cert
      - perfer
      - verify-ca
      - verify-full (supported only when the instance runs PostgreSQL 12 or later)
  * `ca_type` - (Available in 1.124.1+) The type of the server certificate. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs. Valid values:
      - aliyun: a cloud certificate
      - custom: a custom certificate
  * `client_ca_cert` - (Available in 1.124.1+) The public key of the CA that issues client certificates. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs.
  * `client_cert_revocation_list` - (Available in 1.124.1+) The certificate revocation list (CRL) that contains revoked client certificates. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs.
  * `last_modify_status` - (Available in 1.124.1+) The status of the SSL link. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs. Valid values:
      - success
      - setting
      - failed
  * `modify_status_reason` - (Available in 1.124.1+) The reason why the SSL link stays in the current state. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs.
  * `replication_acl` - (Available in 1.124.1+) The method that is used to verify the replication permission. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs. Valid values:
      - cert
      - perfer
      - verify-ca
      - verify-full (supported only when the instance runs PostgreSQL 12 or later)
  * `require_update_item` - (Available in 1.124.1+) The server certificate that needs to be updated. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs.
  * `require_update_reason` - (Available in 1.124.1+) The reason why the server certificate needs to be updated. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs.
  * `ssl_create_time` - (Available in 1.124.1+) The time when the server certificate was created. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs. In addition, this parameter is valid only when the CAType parameter is set to aliyun.
  * `ssl_enabled` - (Available in 1.124.1+) Indicates whether SSL encryption is enabled. Valid values:
      - on: enabled
      - off: disabled
  * `server_ca_url` - (Available in 1.124.1+) The URL of the CA that issues the server certificate. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs.
  * `server_cert` - (Available in 1.124.1+) The content of the server certificate. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs.
  * `server_key` - (Available in 1.124.1+) The private key of the server certificate. This parameter is supported only when the instance runs PostgreSQL with standard or enhanced SSDs.
  * `creator` - (Available in 1.124.3+) The creator of the encryption key.
  * `delete_date` - (Available in 1.124.3+) The estimated time when the encryption key will be deleted. The time follows the ISO 8601 standard in the yyyy-MM-ddTHH:mm:ssZ format. The time is displayed in UTC.
  * `description` - (Available in 1.124.3+) The description of the encryption key.
  * `encryption_key` - (Available in 1.124.3+) The ID of the encryption key.
  * `encryption_key_status` - (Available in 1.124.3+) The status of the encryption key. Valid values:
      - Enabled
      - Disabled
  * `key_usage` - (Available in 1.124.3+) The purpose of the encryption key.
  * `material_expire_time` - (Available in 1.124.3+) The time when the encryption key expires. The time follows the ISO 8601 standard in the yyyy-MM-ddTHH:mm:ssZ format. The time is displayed in UTC.
  * `origin` - (Available in 1.124.3+) The source of the encryption key.
  * `parameters` - (Available in 1.135.0+) Parameter list.
      * `parameter_name` - The name of the parameter.
      * `parameter_value` - The default value of the parameter.
      * `force_modify` - Indicates whether the parameter can be modified. Valid values: true | false
      * `force_restart` - Indicates whether the modified parameter takes effect only after a database restart. Valid values: true | false
      * `checking_code` - The value range of the parameter.
      * `parameter_description` - The description of the parameter.
  * `deletion_protection` - (Available in 1.167.0+) Indicates whether the release protection feature is enabled for the instance. Valid values:
      * **true**: The release protection feature is enabled.
      * **false**: The release protection feature is disabled.