---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_bandwidth_package_attachment"
sidebar_current: "docs-alicloud-resource-ga-bandwidth-package-attachment"
description: |-
  Provides a Alicloud Global Accelerator (GA) Bandwidth Package Attachment resource.
---

# alicloud_ga_bandwidth_package_attachment

Provides a Global Accelerator (GA) Bandwidth Package Attachment resource.

For information about Global Accelerator (GA) Bandwidth Package Attachment and how to use it, see [What is Bandwidth Package Attachment](https://www.alibabacloud.com/help/en/global-accelerator/latest/api-ga-2019-11-20-bandwidthpackageaddaccelerator).

-> **NOTE:** Available since v1.113.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ga_bandwidth_package_attachment&exampleId=fdd25e31-816e-745a-9fdd-10fe6fdaea88470b76d6&activeTab=example&spm=docs.r.ga_bandwidth_package_attachment.0.fdd25e3181&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_ga_accelerator" "default" {
  duration        = 1
  auto_use_coupon = true
  spec            = "1"
}

resource "alicloud_ga_bandwidth_package" "default" {
  bandwidth      = 100
  type           = "Basic"
  bandwidth_type = "Basic"
  payment_type   = "PayAsYouGo"
  billing_type   = "PayBy95"
  ratio          = 30
}

resource "alicloud_ga_bandwidth_package_attachment" "default" {
  accelerator_id       = alicloud_ga_accelerator.default.id
  bandwidth_package_id = alicloud_ga_bandwidth_package.default.id
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ga_bandwidth_package_attachment&spm=docs.r.ga_bandwidth_package_attachment.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `accelerator_id` - (Required, ForceNew) The ID of the Global Accelerator instance.
* `bandwidth_package_id` - (Required) The ID of the Bandwidth Package. **NOTE:** From version 1.192.0, `bandwidth_package_id` can be modified.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Bandwidth Package Attachment. It formats as `<accelerator_id>:<bandwidth_package_id>`.
-> **NOTE:** Before provider version 1.120.0, it formats as `<bandwidth_package_id>`.
* `accelerators` - Accelerators bound with current Bandwidth Package.
* `status` - State of Bandwidth Package.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Bandwidth Package Attachment.
* `update` - (Defaults to 5 mins) Used when update the Bandwidth Package Attachment.
* `delete` - (Defaults to 5 mins) Used when delete the Bandwidth Package Attachment.

## Import

Ga Bandwidth Package Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_ga_bandwidth_package_attachment.example <accelerator_id>:<bandwidth_package_id>
```
