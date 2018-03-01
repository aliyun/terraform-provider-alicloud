---
layout: "alicloud"
page_title: "Alicloud: alicloud_eips"
sidebar_current: "docs-alicloud-datasource-eips"
description: |-
    Provides a list of EIP which owned by an Alicloud account.
---

# alicloud\_eips

The elastic ip address data source lists a list of eips resource information owned by an Alicloud account,
and each EIP including its basic attribution and association instance.

## Example Usage

```
data "alicloud_eips" "eips"{
    cidr_block="172.16.0.0/12"
    name_regex="^foo"
}

resource "alicloud_instance" "foo" {
    ...
    instance_name =  "in-the-eip"
    vswitch_id = "vsw-abc123456"
    ...
}

resource "alicloud_eip_association" "asso" {
    instance_id = "${alicloud_instance.foo.id}"
    allocation_id = "${data.alicloud_eips.eips.eips.0.id}"
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of EIP allocation ID.
* `ip_addresses` - (Optional) A list of EIP ip address ID.
* `in_use` - (Deprecated) It has been deprecated from provider version 1.8.0.
* `output_file` - (Optional) The name of file that can save eips data source after running `terraform plan`.

## Attributes Reference

The following attributes are exported:

* `eips` A list of eips. It contains several attributes to `Block EIPs`.

### Block EIPs

Attributes for eips:

* `id` - ID of the EIP.
* `status` - EIP status.
* `ip_address` - Address of the the EIP.
* `bandwidth` - EIP internat max bandwidth.
* `internet_charge_type` - EIP internet charge type.
* `instance_id` - ID of the instance with which EIP association.
* `instance_id` - Type of the instance with which EIP association.
* `creation_time` - Time of creation.
