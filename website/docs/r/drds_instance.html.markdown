---
layout: "alicloud"
page_title: "Alicloud: alicloud_drds_instance"
sidebar_current: "docs-alicloud-resource-drds-instance"
description: |-
  Provides an DRDS instance resource.
---

# alicloud\_drds\_instance

Distributed Relational Database Service (DRDS) is a lightweight (stateless), flexible, stable, and efficient middleware product independently developed by Alibaba Group to resolve scalability issues with single-host relational databases.
With its compatibility with MySQL protocols and syntaxes, DRDS enables database/table sharding, smooth scaling, configuration upgrade/downgrade,
transparent read/write splitting, and distributed transactions, providing O&M capabilities for distributed databases throughout their entire lifecycle.

## Example Usage

```
 resource "alicloud_drds_instance" "default" {
  description = "drds"
  type = "1"
  pay_type = "Postpaid"
  zone_id = "cn-hangzhou-e"
  vswitch_id = "vsw-bp1jlu3swk8rq2yoi40ey"
  instance_series = "drds.sn1.4c8g"
  specification = "drds.sn1.4c8g.8C16G"
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) Description of the DRDS instance, This description can have a string of 2 to 256 characters.
* `zone_id` - (Required) The Zone to launch the DRDS instance.
* `instance_charge_type` -  (Optional) Valid values are `Prepaid`, `Postpaid`, Default to `Postpaid`.
* `vswitch_id` - (Required for a VPC SLB, Forces New Resource) The VSwitch ID to launch in.
* `instance_series` - (Required) User-defined DB instance storage space. Value range:
    - `drds.sn1.4c8g` for DRDS instance Starter version;
    - `drds.sn1.8c16g` for DRDS instance Standard edition;
    - `drds.sn1.16c32g` for DRDS instance Enterprise Edition;
    - `drds.sn1.32c64g` for DRDS instance Extreme Edition;
* `specification` - "drds.sn1.4c8g.8C16G" 

-> **NOTE:** Because of replace DRDS instance nodes, change DRDS instance type and specification would cost 1~5 minutes. Please make full preparation before changing them.

-> **NOTE:** You can use [DRDS DOC](https://www.alibabacloud.com/help/product/29657.htm) to do it.

## Attributes Reference

The following attributes are exported:

* `id` - The DRDS instance ID.

## Import

Distributed Relational Database Service (DRDS) can be imported using the id, e.g.

```
$ terraform import alicloud_drds_instance.example drds-abc123456
```
