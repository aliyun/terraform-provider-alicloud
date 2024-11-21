---
subcategory: "Cloud DBAudit (DBAudit)"
layout: "alicloud"
page_title: "Alicloud: alicloud_yundun_dbaudit_instance"
sidebar_current: "docs-alicloud-resource-yundun-dbaudit-instance"
description: |-
  Provides a Alicloud Cloud DBaudit Instance Resource.
---

# alicloud_yundun_dbaudit_instance

Cloud DBaudit instance resource ("Yundun_dbaudit" is the short term of this product).

-> **NOTE:** The endpoint of bssopenapi used only support "business.aliyuncs.com" at present.

-> **NOTE:** Available since v1.62.0+.

-> **NOTE:** In order to destroy Cloud DBaudit instance , users are required to apply for white list first

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_yundun_dbaudit_instance&exampleId=e879359e-94c4-313b-85a0-25da5642294787674a37&activeTab=example&spm=docs.r.yundun_dbaudit_instance.0.e879359e94&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
provider "alicloud" {
  endpoints {
    bssopenapi = "business.aliyuncs.com"
  }
  region = "cn-hangzhou"
}

data "alicloud_zones" "default" {
  available_resource_creation = "VSwitch"
}

data "alicloud_vpcs" "default" {
  name_regex = "^default-NODELETING$"
}
data "alicloud_vswitches" "default" {
  vpc_id  = data.alicloud_vpcs.default.ids.0
  zone_id = data.alicloud_zones.default.zones.0.id
}

resource "alicloud_yundun_dbaudit_instance" "default" {
  description = "tf-example"
  plan_code   = "alpha.professional"
  period      = "1"
  vswitch_id  = data.alicloud_vswitches.default.ids.0
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

```shell
$ terraform import alicloud_yundun_dbaudit_instance.example dbaudit-exampe123456
```
