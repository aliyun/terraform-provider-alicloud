---
subcategory: "Cloud Monitor Service"
layout: "alicloud"
page_title: "Alicloud: alicloud_cloud_monitor_service_hybrid_double_write"
description: |-
  Provides a Alicloud Cloud Monitor Service Hybrid Double Write resource.
---

# alicloud_cloud_monitor_service_hybrid_double_write

Provides a Cloud Monitor Service Hybrid Double Write resource. 

For information about Cloud Monitor Service Hybrid Double Write and how to use it, see [What is Hybrid Double Write](https://next.api.alibabacloud.com/document/Cms/2018-03-08/CreateHybridDoubleWrite).

-> **NOTE:** Available since v1.210.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cloud_monitor_service_hybrid_double_write&exampleId=7744a17d-acb2-4f18-83cf-9d0f88acf1bbc3655837&activeTab=example&spm=docs.r.cloud_monitor_service_hybrid_double_write.0.7744a17dac&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}

data "alicloud_account" "default" {
}

resource "alicloud_cms_namespace" "source" {
  namespace = var.name
}

resource "alicloud_cms_namespace" "default" {
  namespace = "${var.name}-source"
}

resource "alicloud_cloud_monitor_service_hybrid_double_write" "default" {
  source_namespace = alicloud_cms_namespace.source.id
  source_user_id   = data.alicloud_account.default.id
  namespace        = alicloud_cms_namespace.default.id
  user_id          = data.alicloud_account.default.id
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cloud_monitor_service_hybrid_double_write&spm=docs.r.cloud_monitor_service_hybrid_double_write.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `source_namespace` - (Required, ForceNew) Source Namespace.
* `source_user_id` - (Required, ForceNew) Source UserId.
* `namespace` - (Required, ForceNew) Target Namespace.
* `user_id` - (Required, ForceNew) Target UserId.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Hybrid Double Write. It formats as `<source_namespace>:<source_user_id>`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Hybrid Double Write.
* `delete` - (Defaults to 5 mins) Used when delete the Hybrid Double Write.

## Import

Cloud Monitor Service Hybrid Double Write can be imported using the id, e.g.

```shell
$ terraform import alicloud_cloud_monitor_service_hybrid_double_write.example <source_namespace>:<source_user_id>
```
