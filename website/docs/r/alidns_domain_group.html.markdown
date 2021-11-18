---
subcategory: "DNS"
layout: "alicloud"
page_title: "Alicloud: alicloud_alidns_domain_group"
sidebar_current: "docs-alicloud-resource-alidns-domain-group"
description: |-
  Provides a Alidns Domain Group resource.
---

# alicloud\_alidns\_domain\_group

Provides a Alidns Domain Group resource. For information about Alidns Domain Group and how to use it, see [What is Resource Alidns Domain Group](https://www.alibabacloud.com/help/en/doc-detail/29762.htm).

-> **NOTE:** Available in v1.84.0+.

## Example Usage

```terraform
# Add a new Alinds Domain Group.
resource "alicloud_alidns_domain_group" "example" {
  domain_group_name = "tf-testDG"
}
```
## Argument Reference

The following arguments are supported:

* `group_name` - (Optional, Deprecated in v1.97.0+) The Name of the domain group. The `group_name` has been deprecated from provider version 1.97.0. Please use `domain_group_name` instead.
* `domain_group_name` - (Optional, Available in v1.97.0+) The Name of the domain group. The `domain_group_name` is required when the value of the `group_name`  is Empty.
* `lang` - (Optional) User language. 

## Attributes Reference

The following attributes are exported:

* `id` - This ID of this domain group resource.

## Import

Alidns domain group can be imported using the id, e.g.

```
$ terraform import alicloud_alidns_domain_group.example 0932eb3ddee7499085c4d13d45*****
```
