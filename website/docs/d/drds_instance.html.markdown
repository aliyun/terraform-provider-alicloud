---
layout: "alicloud"
page_title: "Alicloud: alicloud_drds_instances"
sidebar_current: "docs-alicloud-datasource-drds-instances"
description: |-
    Provides a collection of DRDS instances according to the specified filters.
---

# alicloud\_drds\_instances

The `alicloud_db_instances` data source provides a collection of DRDS instances available in Alibaba Cloud account.
Filters support regular expression for the instance name, searches by tags, and other filters which are listed below.

## Example Usage

```
data "alicloud_drds_instances" "drds_instances_ds" {
  name_regex = "drds-\\d+"
  status     = "Running"
  tags       = <<EOF
{
  "type": "private",
}
EOF
}

output "first_drds_instance_id" {
  value = "${data.alicloud_drds_instances.0.d}"
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
    - `drds.sn1.32c64g` for DRDS instance Extreme Edition;
    
~> **NOTE:** Because of replace DRDS instance nodes, change DRDS instance type and specification would cost 1~5 minutes. Please make full preparation before changing them.

## Attributes Reference

The following attributes are exported:

* `id` - The DRDS instance ID.
