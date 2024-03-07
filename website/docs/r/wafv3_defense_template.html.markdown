---
subcategory: "Web Application Firewall(WAF)"
layout: "alicloud"
page_title: "Alicloud: alicloud_wafv3_defense_template"
description: |-
  Provides a Alicloud WAFV3 Defense Template resource.
---

# alicloud_wafv3_defense_template

Provides a WAFV3 Defense Template resource. 

For information about WAFV3 Defense Template and how to use it, see [What is Defense Template](https://www.alibabacloud.com/help/en/web-application-firewall/latest/api-waf-openapi-2021-10-01-createdefensetemplate).

-> **NOTE:** Available since v1.218.0.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "terraform-example"
}

data "alicloud_wafv3_instances" "default" {
}

resource "alicloud_wafv3_defense_template" "default" {
  status                = "1"
  instance_id           = data.alicloud_wafv3_instances.default.ids.0
  defense_template_name = var.name

  template_type                      = "user_custom"
  template_origin                    = "custom"
  defense_scene                      = "antiscan"
  resource_manager_resource_group_id = "example"
  description                        = var.name
}
```

## Argument Reference

The following arguments are supported:
* `defense_scene` - (Required, ForceNew) The module to which the protection rule that you want to create belongs. Value:
  - **waf_group**: the basic protection rule module.
  - **antiscan**: the scan protection module.
  - **ip_blacklist**: the IP address blacklist module.
  - **custom_acl**: the custom rule module.
  - **whitelist**: the whitelist module.
  - **region_block**: the region blacklist module.
  - **custom_response**: the custom response module.
  - **cc**: the HTTP flood protection module.
  - **tamperproof**: the website tamper-proofing module.
  - **dlp**: the data leakage prevention module.
* `defense_template_name` - (Required) The name of the protection rule template.
* `description` - (Optional) The description of the protection rule template. .
* `instance_id` - (Required, ForceNew) The ID of the Web Application Firewall (WAF) instance. .
* `resource_manager_resource_group_id` - (Optional) The ID of the Alibaba Cloud resource group. .
* `status` - (Required) The status of the protection rule template. Valid values:
  - **0**: disabled.
  - **1**: enabled.
* `template_origin` - (Required, ForceNew) The origin of the protection rule template that you want to create. Set the value to **custom**. The value specifies that the protection rule template is a custom template. .
* `template_type` - (Required, ForceNew) The type of the protection rule template. Valid values:
  - **user_default:** default template.
  - **user_custom:** custom template.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<instance_id>:<defense_template_id>`.
* `defense_template_id` - Template ID.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Defense Template.
* `delete` - (Defaults to 5 mins) Used when delete the Defense Template.
* `update` - (Defaults to 5 mins) Used when update the Defense Template.

## Import

WAFV3 Defense Template can be imported using the id, e.g.

```shell
$ terraform import alicloud_wafv3_defense_template.example <instance_id>:<defense_template_id>
```