---
subcategory: "Cloud Control"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_control_resource_types"
sidebar_current: "docs-alicloud-datasource-cloud-control-resource-types"
description: |-
  Provides a list of Cloud Control Resource Type owned by an Alibaba Cloud account.
---

# alicloud_cloud_control_resource_types

This data source provides Cloud Control Resource Type available to the user.[What is Resource Type](https://next.api.aliyun.com/document/cloudcontrol/2022-08-30/GetResourceType)

-> **NOTE:** Available since v1.241.0.

## Example Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}

data "alicloud_cloud_control_resource_types" "default" {
  product = "VPC"
  ids     = ["VSwitch"]
}
```

## Argument Reference

The following arguments are supported:
* `product` - (Required, ForceNew) Product Code.
* `ids` - (Optional, ForceNew, Computed) A list of Resource Type IDs.
* `output_file` - (Optional, ForceNew) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Resource Type IDs.
* `types` - A list of Resource Type Entries. Each element contains the following attributes:
  * `create_only_properties` - Create an operation private parameter collection. The attributes are not returned in the resource query operation, but the parameters are required in the creation operation.
  * `delete_only_properties` - Delete operation private parameter collection. The attribute is not returned in the resource query operation, but the parameter is required in the delete operation.
  * `filter_properties` - A collection of attributes that can be used as the filter parameter during the list operation.
  * `get_only_properties` - Query operation private parameter collection. The attribute is not returned in the resource query operation, but the input parameter is required in the query operation.
  * `get_response_properties` - The collection of properties returned by the query.
  * `handlers` - Supported resource operation information (including RAM permissions).
    * `create` - Create operation association information.
      * `permissions` - The collection of required RAM permission information.
    * `delete` - Delete operation association information.
      * `permissions` - The collection of required RAM permission information.
    * `get` - Query operation association information.
      * `permissions` - The collection of required RAM permission information.
    * `list` - List operation association information.
      * `permissions` - The collection of required RAM permission information.
    * `update` - Update operation association information.
      * `permissions` - The collection of required RAM permission information.
  * `info` - Basic information about the resource type.
    * `charge_type` - Payment formpaid (paid)(free).
    * `delivery_scope` - Delivery Levelcenter (centralized deployment level)region (regional deployment level)zone (Availability zone deployment level).
    * `description` - Resource type description.
    * `title` - The resource type name.
  * `list_only_properties` - Enumerate the operation private parameter collection. The attributes are not returned in the resource query operation, but the parameters that need to be passed in the enumeration operation.
  * `list_response_properties` - Enumerates the returned property collection.
  * `primary_identifier` - Resource ID
  * `product` - Product Code.
  * `properties` - Resource attribute definition, where key is the attribute name and value is the attribute details.
  * `public_properties` - A collection of public attributes, which are the basic attributes of the resource. Non-Operation private parameters.
  * `read_only_properties` - A set of read-only parameters. It is returned only in the list or get Operation. It is not used as an input parameter during creation and change.
  * `required` - Resource creation required parameter collection.
  * `resource_type` - The resource type.
  * `sensitive_info_properties` - A collection of sensitive attributes, such as passwords.
  * `update_only_properties` - Update operation private parameter collection. The attributes are not returned in the resource query operation, but the parameters are required in the update operation.
  * `update_type_properties` - A collection of properties that can be modified.
  * `id` - The ID of the resource supplied above.
