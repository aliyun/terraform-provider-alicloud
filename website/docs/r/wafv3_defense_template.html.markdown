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

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_wafv3_defense_template&exampleId=4885bb2d-ced3-0337-0228-cfd2963930d65749a9f7&activeTab=example&spm=docs.r.wafv3_defense_template.0.4885bb2dce&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

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
* `defense_scene` - (Required, ForceNew) The WAF protection scenario to be created. Valid values:
  - `waf_group`: indicates basic protection.
  - `antiscan`: indicates scan protection.
  - `ip_blacklist`: indicates the IP address blacklist.
  - `custom_acl`: indicates a custom rule.
  - `whitelist`: indicates the whitelist.
  - `region_block`: indicates that the region is blocked.
  - `custom_response`: indicates a custom response.
  - `cc`: indicates CC protection.
  - `tamperproof`: Indicates that the web page is tamper-proof.
  - `dlp`: Indicates information leakage protection.
  - `spike_throttle`: indicates peak traffic throttling.

* `defense_template_name` - (Required) The name of the protection rule template.
* `description` - (Optional) The description of the protection rule template.
* `instance_id` - (Required, ForceNew) The ID of the Web Application Firewall (WAF) instance.
* `resource_groups` - (Optional, Set, Available since v1.263.0) The name of the protected object group. After a protection template resource is created, you can modify the bound protection object group.
* `resource_manager_resource_group_id` - (Optional) The ID of the Alibaba Cloud resource group.

-> **NOTE:** This parameter only applies during resource update. If modified in isolation without other property changes, Terraform will not trigger any action.

* `resources` - (Optional, Set, Available since v1.257.0) The list of protected objects to be bound. After a protection template resource is created, you can modify the bound protected objects.
* `status` - (Required) The status of the protection rule template. Valid values:
  - `0`: disabled.
  - `1`: enabled.
* `template_origin` - (Required, ForceNew) The origin of the protection rule template that you want to create. Set the value to `custom`. The value specifies that the protection rule template is a custom template.
* `template_type` - (Required, ForceNew) The type of the protection rule template. Valid values:
  - **user_default:** default template.
  - **user_custom:** custom template.

## Attributes Reference

The following attributes are exported:
* `id` - The ID of the resource supplied above.The value is formulated as `<instance_id>:<defense_template_id>`.
* `defense_template_id` - Template ID

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:
* `create` - (Defaults to 5 mins) Used when create the Defense Template.
* `delete` - (Defaults to 5 mins) Used when delete the Defense Template.
* `update` - (Defaults to 5 mins) Used when update the Defense Template.

## Import

WAFV3 Defense Template can be imported using the id, e.g.

```shell
$ terraform import alicloud_wafv3_defense_template.example <instance_id>:<defense_template_id>
```