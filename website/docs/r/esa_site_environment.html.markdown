---
subcategory: "ESA"
layout: "alicloud"
page_title: "Alicloud: alicloud_esa_site_environment"
description: |-
  Provides a Alicloud ESA Site Environment resource.
---

# alicloud_esa_site_environment

Provides a ESA Site Environment resource.



For information about ESA Site Environment and how to use it, see [What is Site Environment](https://next.api.alibabacloud.com/document/ESA/2024-09-10/CreateEnvironment).

-> **NOTE:** Available since v2.2.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_esa_site_environment&exampleId=21b95e9d-d9f5-ff0e-d0d7-82ccffc9d519593379de&activeTab=example&spm=docs.r.esa_site_environment.0.21b95e9dd9&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_esa_sites" "default" {
  plan_subscribe_type = "advancedplan"
}

resource "alicloud_esa_site_environment" "default" {
  site_id               = data.alicloud_esa_sites.default.sites.0.id
  environment_name      = var.name
  rule                  = "(ip.src in {1.1.1.1})"
  next_environment_name = "Production"
}
```


📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_esa_site_environment&spm=docs.r.esa_site_environment.example&intl_lang=EN_US)


## Argument Reference

The following arguments are supported:

* `site_id` - (Required, ForceNew) The site ID.
* `environment_name` - (Required, ForceNew) The name of the environment.
* `rule` - (Required) The environment rule.
* `next_environment_name` - (Required, ForceNew) The name of the environment with the next priority.
* `read_only` - (Optional, Bool) Specifies whether the environment is read-only. Valid values: `true`, `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Site Environment. It formats as `<site_id>:<environment_name>`.
* `create_time` - The time when the environment was created.
* `update_time` - The time when the environment was updated.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Site Environment.
* `update` - (Defaults to 5 mins) Used when update the Site Environment.
* `delete` - (Defaults to 5 mins) Used when delete the Site Environment.

## Import

ESA Site Environment can be imported using the id, e.g.

```shell
$ terraform import alicloud_esa_site_environment.example <site_id>:<environment_name>
```
