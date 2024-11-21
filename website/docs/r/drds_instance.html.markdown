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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_drds_instance&exampleId=b7bb3e59-7f37-011d-7f04-8aec84bb572e24948d9f&activeTab=example&spm=docs.r.drds_instance.0.b7bb3e597f&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  region = "cn-beijing"
}
data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

variable "instance_series" {
  default = "drds.sn1.4c8g"
}

data "alicloud_vpcs" "default" {
  name_regex = "default-NODELETING"
}
data "alicloud_vswitches" "default" {
  vpc_id = data.alicloud_vpcs.default.ids.0
}

resource "alicloud_drds_instance" "default" {
  description          = "drds instance"
  instance_charge_type = "PostPaid"
  zone_id              = data.alicloud_vswitches.default.vswitches.0.zone_id
  vswitch_id           = data.alicloud_vswitches.default.vswitches.0.id
  instance_series      = var.instance_series
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
* `vpc_id` - (Optional, ForceNew, Available in v1.185.0+) The id of the VPC.
* `mysql_version` - (Optional, ForceNew, Available in v1.201.0+) The MySQL version supported by the instance, with the following range of values. `5`: Fully compatible with MySQL 5.x (default) `8`: Fully compatible with MySQL 8.0. This parameter takes effect when the primary instance is created, and the read-only instance has the same MySQL version as the primary instance by default.
       
### Timeouts

-> **NOTE:** Available in 1.49.0+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 10 mins) Used when creating the drds instance (until it reaches running status). 
* `delete` - (Defaults to 10 mins) Used when terminating the drds instance. 
       
       
## Attributes Reference

The following attributes are exported:

* `id` - The DRDS instance ID.
* `connection_string` - (Available in 1.196.0+) The connection string of the DRDS instance.
* `port` - (Available in 1.196.0+) The connection port of the DRDS instance.


## Import

Distributed Relational Database Service (DRDS) can be imported using the id, e.g.

```shell
$ terraform import alicloud_drds_instance.example drds-abc123456
```
