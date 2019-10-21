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

-> **NOTE:** Available in 1.59.0+

## Example Usage

Basic Usage

```
resource "alicloud_cen_instance" "default" {
 name = "tf-testAccCenConfigName"
 description = "tf-testAccCenConfigDescription"
}

resource "alicloud_ccn_instance" "ccn" {
  name        = "tf_test_foo"
  description = "an example for ccn"
  cidr_block  = "192.168.0.0/24"
  is_default  = true
  cen_id = "${alicloud_cen_instance.default.id}"
  total_count =  "1"
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the CCN instance. The name can contain 2 to 128 characters including a-z, A-Z, 0-9, periods, underlines, and hyphens. The name must start with an English letter, but cannot start with http:// or https://.
* `description` - (Optional) The description of the CCN instance. The description can contain 2 to 256 characters. The description must start with English letters, but cannot start with http:// or https://.
* `cidr_block` - (Optional) The CidrBlock of the CCN instance. Defaults to null.
* `is_default` - (Required) Created by default. If the client does not have ccn in the binding, it will create a ccn for the user to replace.
* `total_count` - (Optional) The total count of Grant or Revoke instance to cen. If total_count is 1, run Grant. If total_count is 0, run Revoke.
* `cen_id` - (Optional) The CenId of the CEN instance, About Grant instance.


## Attributes Reference

The following attributes are exported:

* `id` - The CcnId of the CCN instance. For example "ccn-xxx".
* `name` - The name of the CCN instance. 
* `description` - The description of the CCN instance.
* `cidr_block` - The CidrBlock of the CCN instance.
* `status` - The Status of the CCN instance. For example "Active"
* `total_count` - The total count of Grant or Revoke instance to cen.
* `cen_id` - The CenId of the CEN instance, About Grant instance.


## Import

CCN instance can be imported using the id, e.g.

```
$ terraform import alicloud_ccn_instance.example ccn-abc123456
```

