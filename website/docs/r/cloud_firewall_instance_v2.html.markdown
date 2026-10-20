---
subcategory: "Cloud Firewall"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_firewall_instance_v2"
description: |-
  Provides a Alicloud Cloud Firewall Instance V2 resource.
---

# alicloud_cloud_firewall_instance_v2

Provides a Cloud Firewall Instance V2 resource.

Cloud Firewall instance.

For information about Cloud Firewall Instance V2 and how to use it, see [What is Instance V2](https://www.alibabacloud.com/help/en/product/90174.htm).

-> **NOTE:** Available since v1.269.0.

-> **NOTE:** The billing method of an existing Cloud Firewall instance cannot be converted through the API or Terraform. To convert a `PayAsYouGo` instance to `Subscription`, perform the smooth conversion in the [Cloud Firewall console](https://www.alibabacloud.com/help/en/cloud-firewall/cloudfirewall/product-overview/smoothly-change-the-billing-method-between-pay-as-you-go-and-subscription) (**Overview** -> **More** -> **Convert Pay-As-You-Go to Subscription**). Changing `payment_type` on an existing instance that has not been converted returns an error, and the instance is never destroyed or recreated. When `payment_type` changes together with `product_type` or `spec` (both ForceNew), the change is rejected during planning instead of planning a replacement. After the conversion, update `payment_type`, `product_type` and `spec` in the configuration to match the converted instance (the instance ID is retained), or import the converted instance with `terraform import`, and then run `terraform plan` to confirm that no changes are pending. Changing only `spec` or only `product_type` still replaces the instance.

-> **NOTE:** Destroying an instance releases it immediately and cannot be undone. A `PayAsYouGo` instance is released through the Cloud Firewall `ReleasePostInstance` API. A `Subscription` instance is unsubscribed through the BSS `RefundInstance` API with immediate release, so the unused part of the subscription is refunded according to the [unsubscription rules](https://www.alibabacloud.com/help/en/user-center/cancel-subscription/). If the refund is rejected, for example because the instance is not eligible for a refund, the refund quota is exhausted or the instance has an unsettled order, `terraform destroy` fails and the instance is kept so that the state still matches the cloud.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cloud_firewall_instance_v2&exampleId=7ae5ef1f-3b76-5210-699e-2873c8ad64ba251d14e0&activeTab=example&spm=docs.r.cloud_firewall_instance_v2.0.7ae5ef1f3b&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_cloud_firewall_instance_v2" "default" {
  payment_type = "PayAsYouGo"
  product_code = "cfw"
  product_type = "cfw_elasticity_public_cn"
  spec         = "payg_version"
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cloud_firewall_instance_v2&spm=docs.r.cloud_firewall_instance_v2.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `cfw_log` - (Optional, Bool) Whether to use log audit. Valid values:
  - `true`: Enabled.
  - `false`: Disabled.
* `modify_type` - (Optional) The type of modification. Valid values: `Upgrade`, `Downgrade`. **NOTE:** The `modify_type` is required when you execute an update operation.
* `payment_type` - (Required) The payment type of the resource. Valid values: `PayAsYouGo`, `Subscription`. **NOTE:** The billing method of an existing instance cannot be converted through the API. Changing this value on an existing instance is rejected unless the instance has already been converted to the target billing method in the Cloud Firewall console; in that case the state is aligned with the converted instance without destroying or recreating it. A change of this value combined with a change of `product_type` or `spec` is rejected during planning, so no replacement is planned.
* `period` - (Optional) The prepaid period. **NOTE:** If `payment_type` is set to `Subscription`, `period` is required.
* `product_code` - (Required, ForceNew) The product code. Valid values: `cfw`.
* `product_type` - (Required, ForceNew) The product type. Valid values: `cfw_elasticity_public_cn`, `cfw_elasticity_public_intl`, `cfw_sub_public_cn`, `cfw_sub_public_intl`.
* `renewal_duration` - (Optional, Int) The auto-renewal duration. **NOTE:** `renewal_duration` takes effect only if `payment_type` is set to `Subscription`, and `renewal_status` is set to `AutoRenewal`.
* `renewal_duration_unit` - (Optional) The unit of the auto-renewal period. Valid values:
  - `M`: Month.
  - `Y`: Year.
* `renewal_status` - (Optional) Whether to renew an instance automatically or not.
  - `AutoRenewal`: Auto renewal.
  - `ManualRenewal`: Manual renewal.
  - `NotRenewal`: No renewal any longer. After you specify this value, Alibaba Cloud stop sending notification of instance expiry, and only gives a brief reminder on the third day before the instance expiry.
* `sdl` - (Optional, Bool) Data leakage protection status. Valid values: `true`, `false`.
* `spec` - (Required, ForceNew) The edition of the Cloud Firewall instance. Valid values: `payg_version`, `premium_version`, `enterprise_version`,`ultimate_version`. **NOTE:** Changing `spec` replaces the instance; an in-place edition change of an existing instance is not supported through this resource.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.
* `create_time` - The creation time.
* `end_time` - The end time.
* `release_time` - The release time.
* `status` - The status of Cloud Firewall Instance.
* `user_status` - The user status of Cloud Firewall Instance.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 10 mins) Used when create the Instance V2.
* `delete` - (Defaults to 5 mins) Used when delete the Instance V2.
* `update` - (Defaults to 5 mins) Used when update the Instance V2.

## Import

Cloud Firewall Instance V2 can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_firewall_instance_v2.example <id>
```
