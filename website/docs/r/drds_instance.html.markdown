---
subcategory: "Distributed Relational Database Service (DRDS)"
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

For information about DRDS and how to use it, see [What is DRDS](https://www.alibabacloud.com/help/product/29657.htm).

-> **NOTE:** At present, DRDS instance only can be supported in the regions: cn-shenzhen, cn-beijing, cn-hangzhou, cn-hongkong, cn-qingdao, ap-southeast-1.

-> **NOTE:** Currently, this resource only support `Domestic Site Account`.

## Example Usage

```
resource "alicloud_drds_instance" "default" {
  description          = "drds instance"
  instance_charge_type = "PostPaid"
  zone_id              = "cn-hangzhou-e"
  vswitch_id           = "vsw-bp1jlu3swk8rq2yoi40ey"
  instance_series      = "drds.sn1.4c8g"
  specification        = "drds.sn1.4c8g.8C16G"
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Required) Description of the DRDS instance, This description can have a string of 2 to 256 characters.
* `zone_id` - (Required from v1.91.0, ForceNew) The Zone to launch the DRDS instance.
* `instance_charge_type` - (Optional, ForceNew) Valid values are `PrePaid`, `PostPaid`, Default to `PostPaid`.
* `vswitch_id` - (Required from v1.91.0, ForceNew) The VSwitch ID to launch in.
* `instance_series` - (Required, ForceNew) The parameter of the instance series. **NOTE:**  `drds.sn1.4c8g`,`drds.sn1.8c16g`,`drds.sn1.16c32g`,`drds.sn1.32c64g` are no longer supported. Valid values:
    - `drds.sn2.4c16g` Starter Edition.
    - `drds.sn2.8c32g` Standard Edition.
    - `drds.sn2.16c64g` Enterprise Edition.
* `specification` - (Required, ForceNew) User-defined DRDS instance specification. Value range:
    - `drds.sn1.4c8g` for DRDS instance Starter version; 
        - value range : `drds.sn1.4c8g.8c16g`, `drds.sn1.4c8g.16c32g`, `drds.sn1.4c8g.32c64g`, `drds.sn1.4c8g.64c128g`
    - `drds.sn1.8c16g` for DRDS instance Standard edition;
        - value range : `drds.sn1.8c16g.16c32g`, `drds.sn1.8c16g.32c64g`, `drds.sn1.8c16g.64c128g`
    - `drds.sn1.16c32g` for DRDS instance Enterprise Edition;
        - value range : `drds.sn1.16c32g.32c64g`, `drds.sn1.16c32g.64c128g`
    - `drds.sn1.32c64g` for DRDS instance Extreme Edition;
        - value range : `drds.sn1.32c64g.128c256g`
       
### Timeouts

-> **NOTE:** Available in 1.49.0+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when creating the drds instance (until it reaches running status). 
* `delete` - (Defaults to 10 mins) Used when terminating the drds instance. 
       
       
## Attributes Reference

The following attributes are exported:

* `id` - The DRDS instance ID.

## Import

Distributed Relational Database Service (DRDS) can be imported using the id, e.g.

```
$ terraform import alicloud_drds_instance.example drds-abc123456
```
