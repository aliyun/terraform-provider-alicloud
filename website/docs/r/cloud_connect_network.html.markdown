---
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_connect_network"
sidebar_current: "docs-alicloud-resource-cloud-connect-network"
description: |-
  Provides a Alicloud Cloud Connect Network resource.
---

# alicloud\_cloud_connect_network

Provides a cloud connect network resource. Cloud Enterprise Network (cloud connect network) is a service that allows you to create a global network for rapidly building a distributed business system with a hybrid cloud computing solution. CCN enables you to build a secure, private, and enterprise-class interconnected network between VPCs in different regions and your local data centers. CCN provides enterprise-class scalability that automatically responds to your dynamic computing requirements.

For information about cloud connect network and how to use it, see [What is Cloud Enterprise Network](https://www.alibabacloud.com/help/doc-detail/93667.htm).

-> **NOTE:** Available in 1.59.0+

## Example Usage

Basic Usage

```
resource "alicloud_cloud_connect_network" "foo" {
  name        = "tf-testAccCloudConnectNetworkName"
  description = "tf-testAccCloudConnectNetworkDescription"
  cidr_block  = "192.168.0.0/24"
  is_default  = true
}
```
## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the CCN instance. The name can contain 2 to 128 characters including a-z, A-Z, 0-9, periods, underlines, and hyphens. The name must start with an English letter, but cannot start with http:// or https://.
* `description` - (Optional) The description of the CCN instance. The description can contain 2 to 256 characters. The description must start with English letters, but cannot start with http:// or https://.
* `cidr_block` - (Optional) The CidrBlock of the CCN instance. Defaults to null.
* `is_default` - (Required) Created by default. If the client does not have ccn in the binding, it will create a ccn for the user to replace.


## Attributes Reference

The following attributes are exported:

* `id` - The CcnId of the CCN instance. For example "ccn-xxx".
* `name` - The name of the CCN instance. 
* `description` - The description of the CCN instance.
* `cidr_block` - The CidrBlock of the CCN instance.
* `status` - The Status of the CCN instance. For example "Active"


## Import

The cloud connect network instance can be imported using the id, e.g.

```
$ terraform import alicloud_cloud_connect_network.example ccn-abc123456
```

