---
subcategory: "Smart Access Gateway (Smartag)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_connect_networks"
sidebar_current: "docs-alicloud-datasource-cloud-connect-networks"
description: |-
    Provides a list of CCN(Cloud Connect Network) instances owned by an Alibaba Cloud account.
---

# alicloud_cloud_connect_networks

This data source provides Cloud Connect Networks available to the user.

-> **NOTE:** Available since v1.59.0+

-> **NOTE:** Only the following regions support create Cloud Connect Network. [`cn-shanghai`, `cn-shanghai-finance-1`, `cn-hongkong`, `ap-southeast-1`, `ap-southeast-3`, `ap-southeast-5`, `ap-northeast-1`, `eu-central-1`]

## Example Usage

```terraform
data "alicloud_cloud_connect_networks" "default" {
  ids        = ["${alicloud_cloud_connect_network.default.id}"]
  name_regex = "^terraform-example"
}
resource "alicloud_cloud_connect_network" "default" {
  name        = "terraform-example"
  description = "terraform-example"
  cidr_block  = "192.168.0.0/24"
  is_default  = true
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
* `networks` - A list of CCN instances. Each element contains the following attributes:
  * `id` - ID of the CCN instance.
  * `name` - Name of the CCN instance.
  * `description` - The description of the CCN instance.
  * `cidr_block` - CidrBlock of the CCN instance.
  * `is_default` - IsDefault of the CCN instance.
