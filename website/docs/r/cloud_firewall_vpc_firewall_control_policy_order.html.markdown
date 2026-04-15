---
subcategory: "Cloud Firewall"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_firewall_vpc_firewall_control_policy_order"
description: |-
  Provides a Alicloud Cloud Firewall Vpc Firewall Control Policy Order resource.
---

# alicloud_cloud_firewall_vpc_firewall_control_policy_order

Provides a Cloud Firewall Vpc Firewall Control Policy Order resource.

ACL priority of the VPC border firewall.

For information about Cloud Firewall Vpc Firewall Control Policy Order and how to use it, see [What is Vpc Firewall Control Policy Order](https://next.api.alibabacloud.com/document/Cloudfw/2017-12-07/ModifyVpcFirewallControlPolicyPosition).

-> **NOTE:** Available since v1.276.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = ""
}


resource "alicloud_cloud_firewall_vpc_firewall_control_policy_order" "default" {
  order           = "1"
  vpc_firewall_id = "cen-38mhpjiqwbkfullqdj"
  lang            = "zh"
  acl_uuid        = "b71137c7-23f0-411d-b6a0-8a2f1977fe6f"
}
```

### Deleting `alicloud_cloud_firewall_vpc_firewall_control_policy_order` or removing it from your configuration

Terraform cannot destroy resource `alicloud_cloud_firewall_vpc_firewall_control_policy_order`. Terraform will remove this resource from the state file, however resources may remain.

## Argument Reference

The following arguments are supported:
* `acl_uuid` - (Optional, ForceNew, Computed) The unique identifier ID of the access control policy.  
  When modifying an access control policy, you must provide its unique identifier ID. You can obtain this ID by calling the [DescribeVpcFirewallControlPolicy](https://help.aliyun.com/document_detail/159758.html) API.
* `lang` - (Optional) The language type used for requests and responses.  

  Valid values:  
  - `zh`: Chinese  
  - `en`: English

  -> **NOTE:** This parameter is only evaluated during resource creation and update. Modifying it in isolation will not trigger any action.

* `order` - (Required) The new priority of the access control policy after modification.  

  -> **NOTE:**  For the valid range of the new priority, see the [API for querying the effective priority range](https://help.aliyun.com/document_detail/474145.html).

* `vpc_firewall_id` - (Required, ForceNew) The ID of the access control policy group for the VPC border firewall. You can obtain this ID by calling the [DescribeVpcFirewallAclGroupList](https://help.aliyun.com/document_detail/159760.html) API.  

    Valid values:  
  - When the VPC border firewall protects Cloud Enterprise Network (CEN), the policy group ID is the CEN instance ID.  

  Example: cen-ervw0g12b5jbw*\*\*\*  
  - When the VPC border firewall protects Express Connect, the policy group ID is the VPC border firewall instance ID.  

  Example: vfw-a42bbb7b887148c9*\*\*\*.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<vpc_firewall_id>:<acl_uuid>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Vpc Firewall Control Policy Order.
* `update` - (Defaults to 5 mins) Used when update the Vpc Firewall Control Policy Order.

## Import

Cloud Firewall Vpc Firewall Control Policy Order can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_firewall_vpc_firewall_control_policy_order.example <vpc_firewall_id>:<acl_uuid>
```