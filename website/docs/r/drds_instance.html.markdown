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
  provider = "alicloud"
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
* `type` - (Required) DRDS instance type. Value options: 
    - private or 1
* `pay_type` - (Deprecated) It has been deprecated from version 1.5.0 and use 'instance_type' to replace.
* `vswitch_id` - (Required for a VPC SLB, Forces New Resource) The VSwitch ID to launch in.
* `instance_series` - (Required) User-defined DB instance storage space. Value range:
    - `drds.sn1.4c8g` for DRDS instance Starter version;
    - `drds.sn1.8c16g` for DRDS instance Standard edition;
    - `drds.sn1.16c32g` for DRDS instance Enterprise Edition;
    - `drds.sn1.32c64g1 for DRDS instance Extreme Edition;
    
~> **NOTE:** Because of replace DRDS instance nodes, change DRDS instance type and specification would cost 1~5 minutes. Please make full preparation before changing them.

## Attributes Reference

The following attributes are exported:

* `id` - The DRDS instance ID.

