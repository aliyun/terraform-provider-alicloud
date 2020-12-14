---
subcategory: "Cloud Bastionhost"
layout: "alicloud"
page_title: "Alicloud: alicloud_yundun_bastionhost_instance"
sidebar_current: "docs-alicloud-resource-yundun-bastionhost-instance"
description: |-
  Provides a Alicloud Cloud Bastionhost Instance Resource.
---

# alicloud_yundun_bastionhost_instance

Cloud Bastionhost instance resource ("Yundun_bastionhost" is the short term of this product). 
For information about Resource Manager Resource Directory and how to use it, see [What is Bastionhost](https://www.alibabacloud.com/help/en/doc-detail/52922.htm).

-> **NOTE:** The endpoint of bssopenapi used only support "business.aliyuncs.com" at present.

-> **NOTE:** Available in 1.63.0+ .

-> **NOTE:** In order to destroy Cloud Bastionhost instance , users are required to apply for white list first

## Example Usage

Basic Usage

```terraform
provider "alicloud" {
  endpoints {
    bssopenapi = "business.aliyuncs.com"
  }
}

resource "alicloud_yundun_bastionhost_instance" "default" {
  description        = "Terraform-test"
  license_code       = "bhah_ent_50_asset"
  period             = "1"
  vswitch_id         = "v-testVswitch"
  security_group_ids = "sg-test"
}
```
## Argument Reference

The following arguments are supported:

* `license_code` - (Required)  The package type of Cloud Bastionhost instance. You can query more supported types through the [DescribePricingModule](https://help.aliyun.com/document_detail/96469.html).
* `description` - (Required) Description of the instance. This name can have a string of 1 to 63 characters.
* `period` - (ForceNew) Duration for initially producing the instance. Valid values: [1~9], 12, 24, 36. Default to 1. At present, the provider does not support modify "period".
* `vswitch_id` - (Required, ForceNew) VSwitch ID configured to Bastionhost.
* `security_group_ids` - (Required) security group IDs configured to Bastionhost.
* `tags` - (Optional, Available in v1.67.0+) A mapping of tags to assign to the resource.
* `resource_group_id` - (Optional, Available in v1.87.0+) The Id of resource group which the Bastionhost Instance belongs. If not set, the resource is created in the default resource group.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the instance resource of Yundun_bastionhost.

## Import

Yundun_bastionhost instance can be imported using the id, e.g.

```
$ terraform import alicloud_yundun_bastionhost.example bastionhost-exampe123456
```