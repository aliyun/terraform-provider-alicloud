---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_cache_reserve_instance"
description: |-
  Provides a Alicloud ESA Cache Reserve Instance resource.
---

# alicloud_esa_cache_reserve_instance

Provides a ESA Cache Reserve Instance resource.



For information about ESA Cache Reserve Instance and how to use it, see [What is Cache Reserve Instance](https://next.api.alibabacloud.com/document/ESA/2024-09-10/PurchaseCacheReserve).

-> **NOTE:** Available since v1.247.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_esa_cache_reserve_instance&exampleId=47a5e858-308c-c1b4-9ebd-c694f5b5915c37c9017c&activeTab=example&spm=docs.r.esa_cache_reserve_instance.0.47a5e85830&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

provider "alicloud" {
  region = "cn-hangzhou"
}


resource "alicloud_esa_cache_reserve_instance" "default" {
  quota_gb     = "10240"
  cr_region    = "CN-beijing"
  auto_renew   = true
  period       = "1"
  payment_type = "Subscription"
  auto_pay     = true
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_esa_cache_reserve_instance&spm=docs.r.esa_cache_reserve_instance.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `auto_pay` - (Optional) Automatic payment.
* `auto_renew` - (Optional) Whether to auto-renew:
  - `true`: Auto-renew.
  - `false`: Do not auto-renew.
* `cr_region` - (Optional, ForceNew) Cache holding area
  - `HK`: Hong Kong, China
  - `CN`: Mainland China
* `payment_type` - (Required, ForceNew) Specifies whether to enable auto payment.
* `period` - (Optional, ForceNew, Int) Purchase period (unit: month).
* `quota_gb` - (Optional, Int) Cache retention specification (unit: GB).

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - Instance purchase time.
* `status` - The status of the cache reserve instance. , it is unavailable.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Cache Reserve Instance.
* `delete` - (Defaults to 5 mins) Used when delete the Cache Reserve Instance.
* `update` - (Defaults to 5 mins) Used when update the Cache Reserve Instance.

## Import

ESA Cache Reserve Instance can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_cache_reserve_instance.example <id>
```