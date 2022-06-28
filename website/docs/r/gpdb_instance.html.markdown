---
subcategory: "AnalyticDB for PostgreSQL (GPDB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_gpdb_instance"
sidebar_current: "docs-alicloud-resource-gpdb-instance"
description: |-
  Provides a AnalyticDB for PostgreSQL instance resource.
---

# alicloud\_gpdb\_instance

Provides a AnalyticDB for PostgreSQL instance resource supports replica set instances only. the AnalyticDB for PostgreSQL provides stable, reliable, and automatic scalable database services. 
You can see detail product introduction [here](https://www.alibabacloud.com/help/doc-detail/35387.htm)

-> **NOTE:**  Available in 1.47.0+

-> **NOTE:**  The following regions don't support create Classic network Gpdb instance.
[`ap-southeast-2`,`ap-southeast-3`,`ap-southeast-5`,`ap-south-1`,`me-east-1`,`ap-northeast-1`,`eu-west-1`,`us-east-1`,`eu-central-1`,`cn-shanghai-finance-1`,`cn-shenzhen-finance-1`,`cn-hangzhou-finance`]

-> **NOTE:**  Create instance or change instance would cost 10~15 minutes. Please make full preparation.

-> **NOTE:**  This resource is used to manage a Reserved Storage Mode instance, and creating a new reserved storage mode instance is no longer supported since v1.127.0. 
You can still use this resource to manage the instance which has been already created, but can not create a new one. 

## Example Usage

### Create a Gpdb instance

```
data "alicloud_zones" "default" {
  available_resource_creation = "Gpdb"
}

resource "alicloud_vpc" "default" {
  name       = "vpc-123456"
  cidr_block = "172.16.0.0/16"
}

resource "alicloud_vswitch" "default" {
  zone_id           = data.alicloud_zones.default.zones[0].id
  vpc_id            = alicloud_vpc.default.id
  cidr_block        = "172.16.0.0/24"
  vswitch_name      = "vpc-123456"
}

resource "alicloud_gpdb_instance" "example" {
  description          = "tf-gpdb-test"
  engine               = "gpdb"
  engine_version       = "4.3"
  instance_class       = "gpdb.group.segsdx2"
  instance_group_count = "2"
  vswitch_id           = alicloud_vswitch.default.id
  security_ip_list     = ["10.168.1.12", "100.69.7.112"]
}
```

## Argument Reference

The following arguments are supported:

* `engine` (Required, ForceNew) Database engine: gpdb. System Default value: gpdb.
* `engine_version` - (Required, ForceNew) Database version. Value options can refer to the latest docs [CreateDBInstance](https://www.alibabacloud.com/help/doc-detail/86908.htm) `EngineVersion`.
* `instance_class` - (Required) Instance specification. see [Instance specifications](https://www.alibabacloud.com/help/doc-detail/86942.htm).
* `instance_group_count` - (Required) The number of groups. Valid values: [2,4,8,16,32]
* `description` - (Optional) The name of DB instance. It a string of 2 to 256 characters.
* `instance_charge_type` - (Optional, ForceNew) Valid values are `PrePaid`, `PostPaid`,System default to `PostPaid`.
* `zone_id` - (Optional, ForceNew) The Zone to launch the DB instance. it supports multiple zone.
If it is a multi-zone and `vswitch_id` is specified, the vswitch must in one of them.
The multiple zone ID can be retrieved by setting `multi` to "true" in the data source `alicloud_zones`.
* `vswitch_id` - (Optional, ForceNew) The virtual switch ID to launch DB instances in one VPC.
* `security_ip_list` - (Optional) List of IP addresses allowed to access all databases of an instance. The list contains up to 1,000 IP addresses, separated by commas. Supported formats include 0.0.0.0/0, 10.23.12.24 (IP), and 10.23.12.24/24 (Classless Inter-Domain Routing (CIDR) mode. /24 represents the length of the prefix in an IP address. The range of the prefix length is [1,32]).
* `tags` - (Optional, Available in v1.55.3+) A mapping of tags to assign to the resource.

### Timeouts

-> **NOTE:** Available in 1.53.0+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 mins) Used when creating the DB instance (until it reaches the initial `Running` status). 
* `delete` - (Defaults to 10 mins) Used when terminating the ADB PG instance.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Instance.

## Import

AnalyticDB for PostgreSQL can be imported using the id, e.g.

```
$ terraform import alicloud_gpdb_instance.example gp-bp1291daeda44194
```
