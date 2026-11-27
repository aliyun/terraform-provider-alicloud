---
subcategory: "Eflo"
layout: "alicloud"
page_title: "Alicloud: alicloud_eflo_vpds"
sidebar_current: "docs-alicloud-datasource-eflo-vpds"
description: |-
  Provides a list of Eflo Vpd owned by an Alibaba Cloud account.
---

# alicloud_eflo_vpds

This data source provides Eflo Vpd available to the user.[What is Vpd](https://next.api.alibabacloud.com/document/eflo/2022-05-30/CreateVpd)

-> **NOTE:** Available since v1.201.0.

## Example Usage

```terraform
data "alicloud_eflo_vpds" "default" {
  ids        = ["${alicloud_eflo_vpd.default.id}"]
  name_regex = alicloud_eflo_vpd.default.name
  vpd_name   = "RMC-Terraform-Test"
}

output "alicloud_eflo_vpd_example_id" {
  value = data.alicloud_eflo_vpds.default.vpds.0.id
}
```

## Argument Reference

The following arguments are supported:
* `resource_group_id` - (Optional) The Resource group id
* `status` - (Optional) The Vpd status. Valid values: `Available`, `Not Available`, `Executing`, `Deleting`,
* `vpd_id` - (Optional) The id of the vpd.
* `vpd_name` - (Optional) The Name of the VPD.
* `ids` - (Optional, Computed) A list of Vpd IDs.
* `name_regex` - (Optional) A regex string to filter results by Group Metric Rule name.
* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Vpd IDs.
* `names` - A list of name of Vpds.
* `vpds` - A list of Vpd Entries.

---

### `vpds`

The vpds supports the following attributes:

* `cidr` - CIDR network segment
* `create_time` - The creation time of the resource
* `gmt_modified` - Modification time
* `id` - The id of the vpd.
* `resource_group_id` - Resource group id
* `status` - The Vpd status.
* `vpd_id` - The id of the vpd.
* `vpd_name` - The Name of the VPD.
