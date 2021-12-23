---
subcategory: "Bastion Host"
layout: "alicloud"
page_title: "Alicloud: alicloud_bastionhost_instance"
sidebar_current: "docs-alicloud-resource-bastionhost-instance"
description: |-
  Provides a Alicloud Bastion Host Instance Resource.
---

# alicloud_bastionhost_instance

-> **NOTE:** From the version 1.132.0, the resource has been renamed to `alicloud_bastionhost_instance`.

Cloud Bastion Host instance resource ("Yundun_bastionhost" is the short term of this product). 
For information about Resource Manager Resource Directory and how to use it, see [What is Bastionhost](https://www.alibabacloud.com/help/en/doc-detail/52922.htm).

-> **NOTE:** The endpoint of bssopenapi used only support "business.aliyuncs.com" at present.

-> **NOTE:** Available in 1.63.0+ .

-> **NOTE:** In order to destroy Cloud Bastionhost instance , users are required to apply for white list first

## Example Usage

Basic Usage

```terraform
resource "alicloud_bastionhost_instance" "default" {
  description        = "Terraform-test"
  license_code       = "bhah_ent_50_asset"
  period             = "1"
  vswitch_id         = "v-testVswitch"
  security_group_ids = ["sg-test", "sg-12345"]
}
```
## Argument Reference

The following arguments are supported:

* `license_code` - (Required)  The package type of Cloud Bastionhost instance. You can query more supported types through the [DescribePricingModule](https://help.aliyun.com/document_detail/96469.html).
* `description` - (Required) Description of the instance. This name can have a string of 1 to 63 characters.
* `period` - (Optional) Duration for initially producing the instance. Valid values: [1~9], 12, 24, 36. At present, the provider does not support modify "period".
-> **NOTE:** The attribute `period` is only used to create Subscription instance or modify the PayAsYouGo instance to Subscription. Once effect, it will not be modified that means running `terraform apply` will not effect the resource.
* `vswitch_id` - (Required, ForceNew) VSwitch ID configured to Bastionhost.
* `security_group_ids` - (Required) security group IDs configured to Bastionhost.
* `tags` - (Optional, Available in v1.67.0+) A mapping of tags to assign to the resource.
* `resource_group_id` - (Optional, Available in v1.87.0+) The Id of resource group which the Bastionhost Instance belongs. If not set, the resource is created in the default resource group.
* `enable_public_access` - (Optional, Available in v1.143.0+)  Whether to Enable the public internet access to a specified Bastionhost instance. The valid values: `true`, `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the instance resource of Bastionhost.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 min) Used when create the Instance.
* `update` - (Defaults to 20 min) Used when create the Instance.

## Import

Yundun_bastionhost instance can be imported using the id, e.g.

```
$ terraform import alicloud_bastionhost_instance.example bastionhost-exampe123456
```
