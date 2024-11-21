---
subcategory: "Security Center"
layout: "alicloud"
page_title: "Alicloud: alicloud_security_center_service_linked_role"
sidebar_current: "docs-alicloud-resource-security-center-service-linked-role"
description: |-
  Provides a Alicloud Security to create the Security Center service-linked roles(SLR).
---

# alicloud\_security\_center\_service\_linked\_role


Using this resource can create SecurityCenter service-linked role : `AliyunServiceRolePolicyForSas`.  This Role is a Resource Access Management (RAM) role, which to obtain permissions to access another Alibaba Cloud service.


For information about Security Center Service Role and how to use it, see [What is Security Center](https://www.alibabacloud.com/help/en/doc-detail/42302.htm).

-> **NOTE:** Available in v1.142.0+.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_security_center_service_linked_role&exampleId=4b8f0b70-e60f-b1bd-ef5e-9782b94bcebb07945dbd&activeTab=example&spm=docs.r.security_center_service_linked_role.0.4b8f0b70e6&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_security_center_service_linked_role" "service_linked_role" {
}

```



## Attributes Reference

The following attributes are exported:
* `id` - The ID of the DataSource. The value is same as `product_name`. Valid Value: `AliyunServiceRolePolicyForSas`.
* `status` - The status of the service Associated role. Valid Values: `true`: Created. `false`: not created.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Slr.
* `delete` - (Defaults to 1 mins, Available in 1.163.0+.) Used when delete the Slr.

## Import

SecurityCenter service-linked roles(SLR) can be imported using the id, e.g.

```shell
$ terraform import alicloud_security_center_service_linked_role.example <product_name>
```
