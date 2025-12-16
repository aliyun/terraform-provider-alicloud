---
subcategory: "Global Accelerator (GA)"
layout: "alicloud"
page_title: "Alicloud: alicloud_ga_accelerator_spare_ip_attachment"
sidebar_current: "docs-alicloud-resource-ga-accelerator-spare-ip-attachment"
description: |-
  Provides a Alicloud Global Accelerator (GA) Accelerator Spare Ip Attachment resource.
---

# alicloud_ga_accelerator_spare_ip_attachment

Provides a Global Accelerator (GA) Accelerator Spare Ip Attachment resource.

For information about Global Accelerator (GA) Accelerator Spare Ip Attachment and how to use it, see [What is Accelerator Spare Ip Attachment](https://www.alibabacloud.com/help/en/global-accelerator/latest/api-ga-2019-11-20-createspareips).

-> **NOTE:** Available since v1.167.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_ga_accelerator_spare_ip_attachment&exampleId=0c81abd6-707c-c4eb-cb1b-95bf9d7e7cf1303b11cc&activeTab=example&spm=docs.r.ga_accelerator_spare_ip_attachment.0.0c81abd670&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
resource "alicloud_ga_accelerator" "default" {
  duration         = 1
  spec             = "1"
  accelerator_name = "terraform-example"
  auto_use_coupon  = true
  description      = "terraform-example"
}

resource "alicloud_ga_bandwidth_package" "default" {
  bandwidth              = 100
  type                   = "Basic"
  bandwidth_type         = "Basic"
  payment_type           = "PayAsYouGo"
  billing_type           = "PayBy95"
  ratio                  = 30
  bandwidth_package_name = "terraform-example"
  auto_pay               = true
  auto_use_coupon        = true
}

resource "alicloud_ga_bandwidth_package_attachment" "default" {
  accelerator_id       = alicloud_ga_accelerator.default.id
  bandwidth_package_id = alicloud_ga_bandwidth_package.default.id
}

resource "alicloud_ga_accelerator_spare_ip_attachment" "default" {
  accelerator_id = alicloud_ga_bandwidth_package_attachment.default.accelerator_id
  spare_ip       = "127.0.0.1"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_ga_accelerator_spare_ip_attachment&spm=docs.r.ga_accelerator_spare_ip_attachment.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `accelerator_id` - (Required, ForceNew) The ID of the global acceleration instance.
* `dry_run` - (Optional) The dry run.
* `spare_ip` - (Required, ForceNew) The standby IP address of CNAME. When the acceleration area is abnormal, the traffic is switched to the standby IP address.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Accelerator Spare Ip Attachment. The value formats as `<accelerator_id>:<spare_ip>`.
* `status` - The status of the standby CNAME IP address.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Accelerator Spare Ip Attachment.
* `delete` - (Defaults to 1 mins) Used when delete the Accelerator Spare Ip Attachment.

## Import

Global Accelerator (GA) Accelerator Spare Ip Attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_ga_accelerator_spare_ip_attachment.example <accelerator_id>:<spare_ip>
```