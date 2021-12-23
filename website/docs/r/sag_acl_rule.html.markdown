---
subcategory: "Smart Access Gateway"
layout: "alicloud"
page_title: "Alicloud: alicloud_sag_acl_rule"
sidebar_current: "docs-alicloud-resource-sag-acl-rule"
description: |-
  Provides a Sag Acl Rule resource.
---

# alicloud\_sag\_acl\_rule

Provides a Sag Acl Rule resource. This topic describes how to configure an access control list (ACL) rule for a target Smart Access Gateway instance to permit or deny access to or from specified IP addresses in the ACL rule.

For information about Sag Acl Rule and how to use it, see [What is access control list (ACL) rule](https://www.alibabacloud.com/help/doc-detail/111483.htm).

-> **NOTE:** Available in 1.60.0+

-> **NOTE:** Only the following regions support create Cloud Connect Network. [`cn-shanghai`, `cn-shanghai-finance-1`, `cn-hongkong`, `ap-southeast-1`, `ap-southeast-2`, `ap-southeast-3`, `ap-southeast-5`, `ap-northeast-1`, `eu-central-1`]

## Example Usage

Basic Usage

```
resource "alicloud_sag_acl" "default" {
  name        = "tf-testAccSagAclName"
  sag_count   = "0"
}
resource "alicloud_sag_acl_rule" "default" {
  acl_id            = alicloud_sag_acl.default.id
  description       = "tf-testSagAclRule"
  policy            = "accept"
  ip_protocol       = "ALL"
  direction         = "in"
  source_cidr       = "10.10.1.0/24"
  source_port_range = "-1/-1"
  dest_cidr         = "192.168.1.0/24"
  dest_port_range   = "-1/-1"
  priority          = "1"
}
```
## Argument Reference

The following arguments are supported:

* `acl_id` - (Required) The ID of the ACL.
* `description` - (Optional) The description of the ACL rule. It must be 1 to 512 characters in length.
* `policy` - (Required) The policy used by the ACL rule. Valid values: accept|drop.
* `ip_protocol` - (Required) The protocol used by the ACL rule. The value is not case sensitive.
* `direction` - (Required) The direction of the ACL rule. Valid values: in|out.
* `source_cidr` - (Required) The source address. It is an IPv4 address range in the CIDR format. Default value: 0.0.0.0/0.
* `source_port_range` - (Required) The range of the source port. Valid value: 80/80.
* `dest_cidr` - (Required) The destination address. It is an IPv4 address range in CIDR format. Default value: 0.0.0.0/0.
* `dest_port_range` - (Required) The range of the destination port. Valid value: 80/80. 
* `priority` - (Optional) The priority of the ACL rule. Value range: 1 to 100. 


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the ACL rule. For example "acr-xxx".

## Import

The Sag Acl Rule can be imported using the id, e.g.

```
$ terraform import alicloud_sag_acl_rule.example acr-abc123456
```

