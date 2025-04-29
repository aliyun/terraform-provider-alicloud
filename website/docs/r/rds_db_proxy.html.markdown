---
subcategory: "RDS"
layout: "alicloud"
page_title: "Alicloud: alicloud_rds_db_proxy"
sidebar_current: "docs-alicloud-resource-rds-db-proxy"
description: |-
  Provides an RDS instance read write splitting connection resource.
---

# alicloud_rds_db_proxy

Information about RDS database exclusive agent and its usage, see [What is RDS DB Proxy](https://www.alibabacloud.com/help/en/apsaradb-for-rds/latest/api-rds-2014-08-15-modifydbproxy).

-> **NOTE:** Available since v1.193.0.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_rds_db_proxy&exampleId=17286700-d627-0c0f-5c6c-18d91eb8572e7b07c365&activeTab=example&spm=docs.r.rds_db_proxy.0.17286700d6&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}
data "alicloud_db_zones" "default" {
  engine         = "MySQL"
  engine_version = "5.6"
}

resource "alicloud_vpc" "default" {
  vpc_name   = var.name
  cidr_block = "172.16.0.0/16"
}
resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = "172.16.0.0/24"
  zone_id      = data.alicloud_db_zones.default.zones.0.id
  vswitch_name = var.name
}

resource "alicloud_security_group" "default" {
  name   = var.name
  vpc_id = alicloud_vpc.default.id
}

resource "alicloud_db_instance" "default" {
  engine                   = "MySQL"
  engine_version           = "5.7"
  instance_type            = "rds.mysql.c1.large"
  instance_storage         = "20"
  instance_charge_type     = "Postpaid"
  instance_name            = var.name
  vswitch_id               = alicloud_vswitch.default.id
  db_instance_storage_type = "local_ssd"
}

resource "alicloud_db_readonly_instance" "default" {
  zone_id               = alicloud_db_instance.default.zone_id
  master_db_instance_id = alicloud_db_instance.default.id
  engine_version        = alicloud_db_instance.default.engine_version
  instance_storage      = alicloud_db_instance.default.instance_storage
  instance_type         = alicloud_db_instance.default.instance_type
  instance_name         = "${var.name}readonly"
  vswitch_id            = alicloud_vswitch.default.id
}

resource "alicloud_rds_db_proxy" "default" {
  instance_id                          = alicloud_db_instance.default.id
  instance_network_type                = "VPC"
  vpc_id                               = alicloud_db_instance.default.vpc_id
  vswitch_id                           = alicloud_db_instance.default.vswitch_id
  db_proxy_instance_num                = 2
  db_proxy_connection_prefix           = "example"
  db_proxy_connect_string_port         = 3306
  db_proxy_endpoint_read_write_mode    = "ReadWrite"
  read_only_instance_max_delay_time    = 90
  db_proxy_features                    = "TransactionReadSqlRouteOptimizeStatus:1;ConnectionPersist:1;ReadWriteSpliting:1"
  read_only_instance_distribution_type = "Custom"
  read_only_instance_weight {
    instance_id = alicloud_db_instance.default.id
    weight      = "100"
  }
  read_only_instance_weight {
    instance_id = alicloud_db_readonly_instance.default.id
    weight      = "500"
  }
}
```

-> **NOTE:** Resource `alicloud_rds_db_proxy` should be created after `alicloud_db_readonly_instance`, so the `depends_on` statement is necessary.

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The Id of instance that can run database.
* `db_proxy_instance_num` - (Required)The number of proxy instances that are enabled. Valid values: 1 to 60.
* `instance_network_type` - (Required, ForceNew)The network type of the instance. Set the value to VPC.
* `vpc_id` - (Required, ForceNew)The ID of the virtual private cloud (VPC) to which the instance belongs.
* `vswitch_id` - (Required, ForceNew)The ID of the vSwitch that is associated with the specified VPC.
* `db_proxy_connection_prefix` - (Optional)The new prefix of the proxy endpoint. Enter a prefix.
* `db_proxy_connect_string_port` - (Optional)The port number that is associated with the proxy endpoint.
* `db_proxy_instance_type` - (Optional, Available since v1.230.0) The database proxy type. Valid values:
  - common: universal proxy.
  - exclusive: Exclusive proxy (default).
* `effective_time` - (Optional)When modifying the number of proxy instances,The time when you want to apply the new database proxy settings.Valid values:
  - Immediate: ApsaraDB RDS immediately applies the new settings.
  - MaintainTime: ApsaraDB RDS applies the new settings during the maintenance window that you specified. For more information, see Modify the maintenance window.
  - SpecificTime: ApsaraDB RDS applies the new settings at a specified point in time.

-> **NOTE:** Note If you set the EffectiveTime parameter to SpecificTime, you must specify the EffectiveSpecificTime parameter.

* `effective_specific_time` - (Optional) The point in time at which you want to apply the new database proxy settings. Specify the time in the ISO 8601 standard in the yyyy-MM-ddTHH:mm:ssZ format. The time must be in UTC.
* `read_only_instance_max_delay_time` - (Optional) The maximum latency threshold that is allowed for read/write splitting. If the latency on a read-only instance exceeds the threshold that you specified, ApsaraDB RDS no longer forwards read requests to the read-only instance. If you do not specify this parameter, the default value of this parameter is retained. Unit: seconds. Valid values: 0 to 3600.

-> **NOTE:** Note If the instance runs PostgreSQL, you can enable only the read/write splitting feature, which is specified by ReadWriteSpliting.

* `db_proxy_features` - (Optional) The features that you want to enable for the proxy endpoint. If you specify more than one feature, separate the features with semicolons (;). Format: Feature 1:Status;Feature 2:Status;.... Do not add a semicolon (;) at the end of the last value. Valid feature values:
  - ReadWriteSpliting: read/write splitting.
  - ConnectionPersist: connection pooling.
  - TransactionReadSqlRouteOptimizeStatus: transaction splitting.
    Valid status values:
    - 1: enabled.
    - 0: disabled.

-> **NOTE:** Note You must specify this parameter only when the read/write splitting feature is enabled.

* `read_only_instance_distribution_type` - (Optional) The policy that is used to allocate read weights. Valid values:
  - Standard: ApsaraDB RDS automatically allocates read weights to the instance and its read-only instances based on the specifications of the instances.
  - Custom: You must manually allocate read weights to the instance and its read-only instances.

-> **NOTE:** Note If you set the ReadOnlyInstanceDistributionType parameter to Custom, you must specify the ReadOnlyInstanceWeight parameter.

* `read_only_instance_weight` - (Optional) A list of the read weights of the instance and its read-only instances.  It contains two sub-fields(instance_id and weight). Read weights increase in increments of 100, and the maximum read weight is 10000. See [`read_only_instance_weight`](#read_only_instance_weight) below.
* `db_proxy_endpoint_read_write_mode` - (Optional) The read and write attributes of the proxy terminal. Valid values:
  - ReadWrite: The proxy terminal connects to the primary instance and can receive both read and write requests.
  - ReadOnly: The proxy terminal does not connect to the primary instance and can receive only read requests. This is the default value.

-> **NOTE:** Note This setting causes your instance to restart. Proceed with caution.

* `db_proxy_ssl_enabled` - (Optional) The SSL configuration setting that you want to apply on the instance. Valid values:
  - Close: disables SSL encryption.
  - Open: enables SSL encryption or modifies the endpoint that requires SSL encryption.
  - Update: updates the validity period of the SSL certificate.
* `upgrade_time` - (Optional) The time when you want to upgrade the database proxy version of the instance. Valid values:
  - MaintainTime: ApsaraDB RDS performs the upgrade during the maintenance window that you specified. This is the default value. For more information, see Modify the maintenance window.
  - Immediate: ApsaraDB RDS immediately performs the upgrade.
  - SpecificTime: ApsaraDB RDS performs the upgrade at a specified point in time.
* `switch_time` - (Optional) The point in time at which you want to upgrade the database proxy version of the instance. Specify the time in the ISO 8601 standard in the yyyy-MM-ddTHH:mm:ssZ format. The time must be in UTC.
* `resource_group_id` - (Optional) The ID of the resource group.

### `read_only_instance_weight`

The read_only_instance_weight mapping supports the following:

* `instance_id` - (Required) The Id of the instance and its read-only instances that can run database.
* `weight` - (Required) Weight of instances that can run the database and their read-only instances. Read weights increase in increments of 100, and the maximum read weight is 10000.

## Attributes Reference

The following attributes are exported:

* `id` - The Id of DB instance.
* `net_type` - Network type of proxy connection address.
* `db_proxy_endpoint_aliases` - Remarks of agent terminal.
* `db_proxy_endpoint_id` - Proxy connection address ID.
* `db_proxy_connection_string` - Connection instance string.
* `ssl_expired_time` - The time when the certificate expires.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 60 mins) Use when opening exclusive agent (until it reaches the initial `Running` status).
* `update` - (Defaults to 30 mins) Used when updating exclusive agent (until it reaches the initial `Running` status).
* `delete` - (Defaults to 20 mins) Use when closing exclusive agent.

## Import

RDS database proxy feature can be imported using the id, e.g.

```shell
$ terraform import alicloud_rds_db_proxy.example abc12345678
```