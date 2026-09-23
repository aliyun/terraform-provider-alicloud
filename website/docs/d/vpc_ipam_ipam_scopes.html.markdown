---
subcategory: "Vpc Ipam"
layout: "alicloud"
page_title: "Alicloud: alicloud_vpc_ipam_ipam_scopes"
sidebar_current: "docs-alicloud-datasource-vpc-ipam-ipam-scopes"
description: |-
  Provides a list of Vpc Ipam Ipam Scope owned by an Alibaba Cloud account.
---

# alicloud_vpc_ipam_ipam_scopes

This data source provides Vpc Ipam Ipam Scope available to the user.[What is Ipam Scope](https://next.api.alibabacloud.com/document/VpcIpam/2023-02-28/CreateIpamScope)

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

resource "alicloud_vpc_ipam_ipam" "defaultIpam" {
  operating_region_list = ["cn-hangzhou"]
  ipam_name             = var.name
}

resource "alicloud_vpc_ipam_ipam_scope" "default" {
  ipam_scope_name        = var.name
  ipam_id                = alicloud_vpc_ipam_ipam.defaultIpam.id
  ipam_scope_description = "This is a ipam scope."
  ipam_scope_type        = "private"
  tags = {
    "k1" : "v1"
  }
}

data "alicloud_vpc_ipam_ipam_scopes" "default" {
  ipam_scope_name = alicloud_vpc_ipam_ipam_scope.default.ipam_scope_name
}

output "alicloud_vpc_ipam_ipam_scope_example_id" {
  value = data.alicloud_vpc_ipam_ipam_scopes.default.scopes.0.id
}
```

## Argument Reference

The following arguments are supported:
* `ipam_id` - (ForceNew, Optional) The id of the Ipam instance.
* `ipam_scope_id` - (ForceNew, Optional) The first ID of the resource.
* `ipam_scope_name` - (ForceNew, Optional) The name of the resource.
* `ipam_scope_type` - (ForceNew, Optional) IPAM scope of action type:**private**.> Currently, only the role scope of the private network is supported.
* `resource_group_id` - (ForceNew, Optional) The ID of the resource group.
* `tags` - (ForceNew, Optional) The tag of the resource.
* `ids` - (Optional, ForceNew, Computed) A list of Ipam Scope IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Group Metric Rule name.
* `output_file` - (Optional, ForceNew) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Ipam Scope IDs.
* `names` - A list of name of Ipam Scopes.
* `scopes` - A list of Ipam Scope Entries. Each element contains the following attributes:
  * `create_time` - The creation time of the resource.
  * `ipam_id` - The id of the Ipam instance.
  * `ipam_scope_description` - The description of the IPAM's scope of action.It must be 2 to 256 characters in length and must start with a lowercase letter, but cannot start with 'http:// 'or 'https. If it is not filled in, it is empty. The default value is empty.
  * `ipam_scope_id` - The first ID of the resource.
  * `ipam_scope_name` - The name of the resource.
  * `ipam_scope_type` - IPAM scope of action type:**private**.> Currently, only the role scope of the private network is supported.
  * `resource_group_id` - The ID of the resource group.
  * `status` - The status of the resource.
  * `tags` - The tag of the resource.
  * `id` - The ID of the resource supplied above.
  * `region_id` - The region ID of the resource.
