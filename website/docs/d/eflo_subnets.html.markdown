---
subcategory: "Eflo"
layout: "alicloud"
page_title: "Alicloud: alicloud_eflo_subnets"
sidebar_current: "docs-alicloud-datasource-eflo-subnets"
description: |-
  Provides a list of Eflo Subnet owned by an Alibaba Cloud account.
---

# alicloud_eflo_subnets

This data source provides Eflo Subnet available to the user.[What is Subnet](https://help.aliyun.com/document_detail/604977.html)

-> **NOTE:** Available in 1.204.0+

## Example Usage

```terraform
data "alicloud_eflo_subnets" "default" {
  name_regex  = alicloud_eflo_subnet.default.name
  subnet_name = "SubnetTestForTerraform"
  vpd_id      = var.vpdId
  zone_id     = var.zoneId
}

output "alicloud_eflo_subnet_example_id" {
  value = data.alicloud_eflo_subnets.default.subnets.0.id
}
```

## Argument Reference

The following arguments are supported:
* `resource_group_id` - (ForceNew,Optional) Resource Group ID.
* `status` - (ForceNew,Optional) The status of the resource.
* `subnet_id` - (ForceNew,Optional) Primary key ID.
* `subnet_name` - (ForceNew,Optional) The Subnet name.
* `type` - (ForceNew,Optional) Eflo subnet usage type, optional value: 
  - General type is not filled in
  - OOB:OOB type 
  - LB: LB type
* `vpd_id` - (ForceNew,Optional) The Eflo VPD ID.
* `zone_id` - (ForceNew,Optional) The zone ID of the resource.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Group Metric Rule name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `names` - A list of name of Subnets.
* `subnets` - A list of Subnet Entries. Each element contains the following attributes:
  * `cidr` - Network segment
  * `create_time` - The creation time of the resource
  * `gmt_modified` - Modification time
  * `message` - Error message
  * `resource_group_id` - Resource Group ID.
  * `status` - The status of the resource.
  * `subnet_id` - The Eflo subnet ID.
  * `subnet_name` - The Subnet name.
  * `type` - Eflo subnet usage type.
  * `vpd_id` - Eflo VPD ID.
  * `id` - The ID of the resource.
  * `zone_id` - The zone ID of the resource.
