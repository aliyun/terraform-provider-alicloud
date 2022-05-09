---
subcategory: "Resource Manager"
layout: "alicloud"
page_title: "Alicloud: alicloud_resource_manager_resource_directory"
sidebar_current: "docs-alicloud-resource-resource-manager-resource-directory"
description: |-
  Provides a Alicloud Resource Manager Resource Directory resource.
---

# alicloud\_resource\_manager\_resource\_directory

Provides a Resource Manager Resource Directory resource. Resource Directory enables you to establish an organizational structure for the resources used by applications of your enterprise. You can plan, build, and manage the resources in a centralized manner by using only one resource directory.

For information about Resource Manager Resource Directory and how to use it, see [What is Resource Manager Resource Directory](https://www.alibabacloud.com/help/en/doc-detail/94475.htm).

-> **NOTE:** Available in v1.84.0+.

-> **NOTE:** An account can only be used to enable a resource directory after it passes enterprise real-name verification. An account that only passed individual real-name verification cannot be used to enable a resource directory.

-> **NOTE:** Before you destroy the resource, make sure that the following requirements are met:
  - All member accounts must be removed from the resource directory. 
  - All folders except the root folder must be deleted from the resource directory.
  
## Example Usage

Basic Usage

```terraform
resource "alicloud_resource_manager_resource_directory" "example" {
  // Enable the control policy
  status = "Enabled"
}
```
## Argument Reference

The following arguments are supported:

* `status` - (Optional, Available in v1.120.0+) The status of control policy. Valid values:`Enabled` and `Disabled`. 

## Attributes Reference

* `id` - The ID of the resource directory.
* `master_account_id` - The ID of the master account.
* `master_account_name` - The name of the master account.
* `root_folder_id` - The ID of the root folder.

-> **NOTE:** Available in 1.120.0+.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `update` - (Defaults to 6 mins) Used when update the control policy status.

## Import

Resource Manager Resource Directory can be imported using the id, e.g.

```
$ terraform import alicloud_resource_manager_resource_directory.example rd-s3****
```
