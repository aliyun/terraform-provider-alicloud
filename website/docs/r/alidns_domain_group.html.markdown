---
subcategory: "Alidns"
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

```
# Add a new Alinds Domain Group.
resource "alicloud_alidns_domain_group" "example" {
  group_name = "tf-testDG"
}
```
## Argument Reference

The following arguments are supported:

* `group_name` - (Required) Name of the domain group. 
* `lang` - (Optional) User language. 

## Attributes Reference

The following attributes are exported:

* `id` - This ID of this domain group resource.

## Import

Alidns domain group can be imported using the id, e.g.

```
$ terraform import alicloud_alidns_domain_group.example 0932eb3ddee7499085c4d13d45*****
```