---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_version"
description: |-
  Provides a Alicloud ESA Version resource.
---

# alicloud_esa_version

Provides a ESA Version resource.



For information about ESA Version and how to use it, see [What is Version](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CloneVersion).

-> **NOTE:** Available since v1.251.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_esa_version&exampleId=4ffb01a7-5a0d-65a8-0b71-88c6361414a8c8d770d5&activeTab=example&spm=docs.r.esa_version.0.4ffb01a75a&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "bcd72239.com"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "enterpriseplan"
}

resource "alicloud_esa_site" "default" {
  site_name          = var.name
  instance_id        = data.alicloud_esa_sites.default.sites.0.instance_id
  coverage           = "overseas"
  access_type        = "NS"
  version_management = true
}

resource "alicloud_esa_version" "default" {
  site_id        = alicloud_esa_site.default.id
  description    = "example"
  origin_version = "0"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_esa_version&spm=docs.r.esa_version.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `description` - (Optional) The Site version's description.
* `site_id` - (Required, ForceNew) The site ID, which can be obtained by calling the ListSites API.
* `origin_version` - (Required, ForceNew, Int) The version number of the site configuration. For sites that have enabled configuration version management, this parameter can be used to specify the effective version of the configuration site, which defaults to version 0.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<site_id>:<site_version>`.
* `create_time` - The creation time. The date format follows ISO8601 notation and uses UTC time. The format is yyyy-MM-ddTHH:mm:ssZ.
* `site_version` - The version number of the site configuration. For sites that have enabled configuration version management, this parameter can be used to specify the effective version of the configuration site, which defaults to version 0.
* `status` - Site version status:：`online`.：`configuring`._faild`：`configure_faild`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Site Version.
* `delete` - (Defaults to 5 mins) Used when delete the Site Version.
* `update` - (Defaults to 5 mins) Used when update the Site Version.

## Import

ESA Site Version can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_version.example <site_id>:<site_version>
```