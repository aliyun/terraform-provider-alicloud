---
layout: "alicloud"
page_title: "Alicloud: alicloud_ccn_instance"
sidebar_current: "docs-alicloud-resource-ccn-instance"
description: |-
  Provides a Alicloud CCN instance resource.
---

# alicloud\_ccn_instance

Provides a CCN instance resource. Cloud Enterprise Network (CCN) is a service that allows you to create a global network for rapidly building a distributed business system with a hybrid cloud computing solution. CCN enables you to build a secure, private, and enterprise-class interconnected network between VPCs in different regions and your local data centers. CCN provides enterprise-class scalability that automatically responds to your dynamic computing requirements.

For information about CCN and how to use it, see [What is Cloud Enterprise Network](https://www.alibabacloud.com/help/doc-detail/93667.htm).

## Example Usage

Basic Usage

```
resource "alicloud_ccn_instance" "ccn" {
  name        = "tf_test_foo"
  description = "an example for ccn"
  cidr_block  = "192.168.0.0/24"
  is_default  = true
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the CCN instance. Defaults to null.
* `description` - (Optional) The description of the CCN instance. Defaults to null.
* `cidr_block` - (Optional) The CidrBlock of the CCN instance. Defaults to null.
* `is_default` - (Required) The IsDefault of the CCN instance.


### Timeouts

-> **NOTE:** Available in 1.48.0+.

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when creating the ccn instance (until it reaches the initial `Active` status). 
* `delete` - (Defaults to 3 mins) Used when terminating the ccn instance. 

## Attributes Reference

The following attributes are exported:

* `id` - The CcnId of the CCN instance.
* `name` - The name of the CCN instance.
* `description` - The description of the CCN instance.
* `cidr_block` - The CidrBlock of the CCN instance.
* `status` - The Status of the CCN instance.



## Import

CCN instance can be imported using the id, e.g.

```
$ terraform import alicloud_ccn_instance.example ccn-abc123456
```

