---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_custom_scene_policy"
description: |-
  Provides a Alicloud ESA Custom Scene Policy resource.
---

# alicloud_esa_custom_scene_policy

Provides a ESA Custom Scene Policy resource.



For information about ESA Custom Scene Policy and how to use it, see [What is Custom Scene Policy](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateCustomScenePolicy).

-> **NOTE:** Available since v1.253.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_esa_custom_scene_policy&exampleId=bc06bf91-0806-9ef8-3906-7dc392cd42ad365d85dc&activeTab=example&spm=docs.r.esa_custom_scene_policy.0.bc06bf9108&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "bcd58610.com"
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

resource "alicloud_esa_custom_scene_policy" "default" {
  end_time                 = "2025-08-07T17:00:00Z"
  create_time              = "2025-07-07T17:00:00Z"
  site_ids                 = alicloud_esa_site.default.id
  template                 = "promotion"
  custom_scene_policy_name = "example-policy"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_esa_custom_scene_policy&spm=docs.r.esa_custom_scene_policy.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:
* `create_time` - (Required) The time when the policy takes effect.
The time follows the ISO 8601 standard in the yyyy-MM-ddTHH:mm:ssZ format. The time is displayed in UTC.
* `custom_scene_policy_name` - (Required) The policy name.
* `end_time` - (Required) The time when the policy expires.
The time follows the ISO 8601 standard in the yyyy-MM-ddTHH:mm:ssZ format. The time is displayed in UTC.
* `site_ids` - (Required) The IDs of websites associated.
* `template` - (Required) The name of the policy template. Valid value:
  - `promotion`: major events.
* `status` - (Optional) Policy effective status. Valid values: `Disabled`, `Running`.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Custom Scene Policy.
* `delete` - (Defaults to 5 mins) Used when delete the Custom Scene Policy.
* `update` - (Defaults to 5 mins) Used when update the Custom Scene Policy.

## Import

ESA Custom Scene Policy can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_custom_scene_policy.example <id>
```