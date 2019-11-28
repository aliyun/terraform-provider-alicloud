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

-> **NOTE:** The endpoint of bssopenapi used only support "business.aliyuncs.com" at present.

-> **NOTE:** Available in 1.63.0+ .

-> **NOTE:** In order to destroy Cloud Bastionhost instance , users are required to apply for white list first

## Example Usage

Basic Usage

```
  provider "alicloud" {
    endpoints {
        bssopenapi = "business.aliyuncs.com"
        }
  }

  resource "alicloud_yundun_bastionhost_instance" "default" {
        description        = "Terraform-test"
        plan_code          = "alpha.professional"
        period             = "1"
        vswitch_id         = "v-testVswitch"
        security_group_ids = "sg-test"
  }
```
## Argument Reference

The following arguments are supported:

* `plan_code` - (Required) Plan code of the Cloud DBaudit to produce. (alpha.professional, alpha.basic, alpha.premium) 
* `description` - (Required) Description of the instance. This name can have a string of 1 to 63 characters.
* `period` - (ForceNew) Duration for initially producing the instance. Valid values: [1~9], 12, 24, 36. Default to 1. At present, the provider does not support modify "period".
* `vswitch_id` - (Required, ForceNew) vSwtich ID configured to audit
* `security_group_ids` - (Required) vSwtich ID configured to audit

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the instance resource of Yundun_bastionhost.

## Import

Yundun_bastionhost instance can be imported using the id, e.g.

```
$ terraform import alicloud_yundun_bastionhost.example bastionhost-exampe123456
```