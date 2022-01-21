---
subcategory: "Cloud DBAudit"
layout: "alicloud"
page_title: "Alicloud: alicloud_yundun_dbaudit_instance"
sidebar_current: "docs-alicloud-resource-yundun-dbaudit-instance"
description: |-
  Provides a Alicloud Cloud DBaudit Instance Resource.
---

# alicloud\_yundun_dbaudit_instance

Cloud DBaudit instance resource ("Yundun_dbaudit" is the short term of this product).

-> **NOTE:** The endpoint of bssopenapi used only support "business.aliyuncs.com" at present.

-> **NOTE:** Available in 1.62.0+ .

-> **NOTE:** In order to destroy Cloud DBaudit instance , users are required to apply for white list first

## Example Usage

Basic Usage

```
provider "alicloud" {
  endpoints {
    bssopenapi = "business.aliyuncs.com"
  }
}

resource "alicloud_yundun_dbaudit_instance" "default" {
  description       = "Terraform-test"
  plan_code         = "alpha.professional"
  period            = "1"
  vswitch_id        = "v-testVswitch"
}
```

### Deleting `alicloud_yundun_dbaudit_instance` or removing it from your configuration

The `alicloud_yundun_dbaudit_instance` resource allows you to manage yundun dbaudit instance, but Terraform cannot destroy it.
Deleting the subscription resource or removing it from your configuration
will remove it from your state file and management, but will not destroy the `yundun_dbaudit_instance`.
You can resume managing the subscription yundun_dbaudit_instance via the AlibabaCloud Console.


## Argument Reference

The following arguments are supported:

* `plan_code` - (Required) Plan code of the Cloud DBAudit to produce. (alpha.professional, alpha.basic, alpha.premium) 
* `description` - (Required) Description of the instance. This name can have a string of 1 to 63 characters.
* `period` - (Required) Duration for initially producing the instance. Valid values: [1~9], 12, 24, 36. At present, the provider does not support modify "period".
-> **NOTE:** The attribute `period` is only used to create Subscription instance or modify the PayAsYouGo instance to Subscription. Once effect, it will not be modified that means running `terraform apply` will not effect the resource.
* `vswitch_id` - (Required, ForceNew) vSwtich ID configured to audit
* `tags` - (Optional, Available in v1.67.0+) A mapping of tags to assign to the resource.
* `resource_group_id` - (Optional, Available in v1.87.0+) The Id of resource group which the DBaudit Instance belongs. If not set, the resource is created in the default resource group.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the instance resource of Yundun_dbaudit.

## Import

Yundun_dbaudit instance can be imported using the id, e.g.

```
$ terraform import alicloud_yundun_dbaudit_instance.example dbaudit-exampe123456
```
