---
layout: "alicloud"
page_title: "Alicloud: alicloud_ccn_instances"
sidebar_current: "docs-alicloud-datasource-ccn-instances"
description: |-
    Provides a list of CCN(Cloud Enterprise Network) instances owned by an Alibaba Cloud account.
---

# alicloud\_ccn\_instances

This data source provides CCN instances available to the user.

-> **NOTE:** Available in 1.59.0+

## Example Usage

```
data "alicloud_ccn_instances" "foo" {
  ids               = ["fake-de-id"]
  name_regex        = "^foo"
  output_file       = "/tmp/ccndomainextension"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of CCN instances IDs.
* `name_regex` - (Optional) A regex string to filter CCN instances by name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of CCN instances IDs.
* `names` - A list of CCN instances names. 
* `instances` - A list of CCN instances. Each element contains the following attributes:
  * `id` - ID of the CCN instance.
  * `name` - Name of the CCN instance.
  * `cidr_block` - CidrBlock of the CCN instance.
  * `is_default` - IsDefault of the CCN instance.