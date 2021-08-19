---
subcategory: "Application Load Balancer (ALB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_alb_acl"
sidebar_current: "docs-alicloud-resource-alb-acl"
description: |-
  Provides a Alicloud Application Load Balancer (ALB) Acl resource.
---

# alicloud\_alb\_acl

Provides a Application Load Balancer (ALB) Acl resource.

For information about ALB Acl and how to use it, see [What is Acl](https://www.alibabacloud.com/help/doc-detail/213617.htm).

-> **NOTE:** Available in v1.133.0+.

## Example Usage

Basic Usage

```terraform
resource "alicloud_alb_acl" "example" {
  acl_name = "example_value"
}

```

## Argument Reference

The following arguments are supported:

* `acl_entries` - (Optional) ACL Entries.
* `acl_name` - (Optional, Computed) The name of the ACL. The name must be 2 to 128 characters in length, and can contain letters, digits, hyphens (-) and underscores (_). It must start with a letter.
* `dry_run` - (Optional) Specifies whether to precheck the API request. Valid values: `true`: only prechecks the API request. If you select this option, the specified endpoint service is not created after the request passes the precheck. The system prechecks the required parameters, request format, and service limits. If the request fails the precheck, the corresponding error message is returned. If the request passes the precheck, the DryRunOperation error code is returned. `false` (default): checks the request. After the request passes the check, an HTTP 2xx status code is returned and the operation is performed.
* `resource_group_id` - (Optional, Computed, ForceNew) Resource Group to Which the Number.

#### Block acl_entries

The acl_entries supports the following: 

* `description` - (Optional) The description of the ACL entry. The description must be 1 to 256 characters in length, and can contain letters, digits, hyphens (-), forward slashes (/), periods (.),and underscores (_). It can also contain Chinese characters.
* `entry` - (Optional) The IP address for the ACL entry.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Acl.
* `status` - The state of the ACL. Valid values:`Provisioning` , `Available` and `Configuring`.  `Provisioning`: The ACL is being created. `Available`: The ACL is available. `Configuring`: The ACL is being configured.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 16 mins) Used when create the Acl.
* `delete` - (Defaults to 16 mins) Used when delete the Acl.
* `update` - (Defaults to 16 mins) Used when update the Acl.

## Import

ALB Acl can be imported using the id, e.g.

```
$ terraform import alicloud_alb_acl.example <id>
```