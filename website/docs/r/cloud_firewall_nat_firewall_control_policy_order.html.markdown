---
subcategory: "Cloud Firewall"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_firewall_nat_firewall_control_policy_order"
description: |-
  Provides a Alicloud Cloud Firewall Nat Firewall Control Policy Order resource.
---

# alicloud_cloud_firewall_nat_firewall_control_policy_order

Provides a Cloud Firewall Nat Firewall Control Policy Order resource.

NAT border firewall ACL priority.

For information about Cloud Firewall Nat Firewall Control Policy Order and how to use it, see [What is Nat Firewall Control Policy Order](https://next.api.alibabacloud.com/document/Cloudfw/2017-12-07/ModifyNatFirewallControlPolicyPosition).

-> **NOTE:** Available since v1.276.0.

~> **NOTE:** The resource can be used to manage the ordering of resource `alicloud_cloud_firewall_nat_firewall_control_policy.new_order`.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cloud_firewall_nat_firewall_control_policy_order&exampleId=60ab672c-ffbe-70c6-4887-67810ea3266c4ce49904&activeTab=example&spm=docs.r.cloud_firewall_nat_firewall_control_policy_order.0.60ab672cff&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = ""
}


resource "alicloud_cloud_firewall_nat_firewall_control_policy_order" "default" {
  acl_uuid       = "a3b5e8f3-6d2c-4e26-b078-87cee0790106"
  nat_gateway_id = "ngw-2ze8hqi59w9wrm30zwgnq"
  direction      = "out"
  order          = "1"
}
```

### Deleting `alicloud_cloud_firewall_nat_firewall_control_policy_order` or removing it from your configuration

Terraform cannot destroy resource `alicloud_cloud_firewall_nat_firewall_control_policy_order`. Terraform will remove this resource from the state file, however resources may remain.


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cloud_firewall_nat_firewall_control_policy_order&spm=docs.r.cloud_firewall_nat_firewall_control_policy_order.example&intl_lang=EN_US)


## Argument Reference

The following arguments are supported:
* `acl_uuid` - (Required, ForceNew) The unique identifier ID of the security access control policy.
* `current_page` - (Optional) The page number of the current page for paginated queries.

  -> **NOTE:** This parameter is only evaluated during resource creation and update. Modifying it in isolation will not trigger any action.

* `direction` - (Optional, ForceNew, Computed) The traffic direction controlled by the access control policy. Valid values:
  - `out`: Access control for outbound traffic (from internal to external).
* `nat_gateway_id` - (Required, ForceNew) The ID of the NAT gateway to query.
* `order` - (Required) The priority at which the access control policy takes effect.
Priority numbers start from 1 and increment sequentially. A smaller priority number indicates a higher priority.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above. The value is formulated as `<acl_uuid>:<nat_gateway_id>:<direction>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Nat Firewall Control Policy Order.
* `update` - (Defaults to 5 mins) Used when update the Nat Firewall Control Policy Order.

## Import

Cloud Firewall Nat Firewall Control Policy Order can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_firewall_nat_firewall_control_policy_order.example <acl_uuid>:<nat_gateway_id>:<direction>
```