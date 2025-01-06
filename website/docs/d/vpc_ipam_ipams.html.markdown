---
subcategory: "Vpc Ipam"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_ipam_ipams"
sidebar_current: "docs-alicloud-datasource-vpc-ipam-ipams"
description: |-
  Provides a list of Vpc Ipam Ipam owned by an Alibaba Cloud account.
---

# alicloud_vpc_ipam_ipams

This data source provides Vpc Ipam Ipam available to the user.[What is Ipam](https://www.alibabacloud.com/help/en/)

-> **NOTE:** Available since v1.241.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_resource_manager_resource_groups" "default" {}


resource "alicloud_vpc_ipam_ipam" "default" {
  ipam_description      = "This is my first Ipam."
  ipam_name             = var.name
  operating_region_list = ["cn-hangzhou"]
}

data "alicloud_vpc_ipam_ipams" "default" {
  ids        = ["${alicloud_vpc_ipam_ipam.default.id}"]
  name_regex = alicloud_vpc_ipam_ipam.default.ipam_name
  ipam_name  = var.name
}

output "alicloud_vpc_ipam_ipam_example_id" {
  value = data.alicloud_vpc_ipam_ipams.default.ipams.0.id
}
```

## Argument Reference

The following arguments are supported:
* `ipam_id` - (ForceNew, Optional) The first ID of the resource.
* `ipam_name` - (ForceNew, Optional) The name of the resource.
* `resource_group_id` - (ForceNew, Optional) The ID of the resource group.
* `tags` - (ForceNew, Optional) The tag of the resource.
* `ids` - (Optional, ForceNew, Computed) A list of Ipam IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Group Metric Rule name.
* `output_file` - (Optional, ForceNew) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Ipam IDs.
* `names` - A list of name of Ipams.
* `ipams` - A list of Ipam Entries. Each element contains the following attributes:
  * `create_time` - The creation time of the resource.
  * `default_resource_discovery_association_id` - After an IPAM is created, the association between the resource discovery created by the system by default and the IPAM.
  * `default_resource_discovery_id` - After IPAM is created, the system creates resource discovery by default.
  * `ipam_description` - The description of IPAM.It must be 2 to 256 characters in length and must start with an uppercase letter or a Chinese character, but cannot start with 'http: // 'or 'https. If the description is not filled in, it is blank. The default value is blank.
  * `ipam_id` - The first ID of the resource.
  * `ipam_name` - The name of the resource.
  * `private_default_scope_id` - After an IPAM is created, the scope of the private network IPAM created by the system by default.
  * `public_default_scope_id` - After an IPAM is created, the public network IPAM is created by default.
  * `region_id` - The region ID of the resource.
  * `resource_discovery_association_count` - The number of resource discovery objects associated with IPAM.
  * `resource_group_id` - The ID of the resource group.
  * `status` - The status of the resource.
  * `tags` - The tag of the resource.
  * `id` - The ID of the resource supplied above.
