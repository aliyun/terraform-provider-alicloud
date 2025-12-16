---
subcategory: "Application Load Balancer (ALB)"
layout: "alicloud"
page_title: "Alicloud: alicloud_alb_acl_entry_attachment"
sidebar_current: "docs-alicloud-resource-alb-acl-entry-attachment"
description: |-
  Provides a Acl entry attachment resource.
---

# alicloud_alb_acl_entry_attachment

For information about acl entry attachment and how to use it, see [Configure an acl entry](https://www.alibabacloud.com/help/en/slb/application-load-balancer/developer-reference/api-alb-2020-06-16-addentriestoacl).

-> **NOTE:** Available since v1.166.0.

## Example Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_alb_acl_entry_attachment&exampleId=47eecb52-ea5d-fc63-d81b-4856fc336b42987e0db2&activeTab=example&spm=docs.r.alb_acl_entry_attachment.0.47eecb52ea&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf_example"
}
data "alicloud_resource_manager_resource_groups" "default" {}

resource "alicloud_alb_acl" "default" {
  acl_name          = var.name
  resource_group_id = data.alicloud_resource_manager_resource_groups.default.groups.0.id
}

resource "alicloud_alb_acl_entry_attachment" "default" {
  acl_id      = alicloud_alb_acl.default.id
  entry       = "168.10.10.0/24"
  description = var.name
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_alb_acl_entry_attachment&spm=docs.r.alb_acl_entry_attachment.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `acl_id` - (Required, ForceNew) The ID of the Acl.
* `entry` - (Required, ForceNew) The CIDR blocks.
* `description` - (Optional, ForceNew) The description of the entry.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource. The value formats as `<acl_id>:<entry>`.
* `status` - The Status of the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the resource.
* `delete` - (Defaults to 5 mins) Used when delete the resource.

## Import

Acl entry attachment can be imported using the id, e.g.

```shell
$ terraform import alicloud_alb_acl_entry_attachment.example <acl_id>:<entry>
```
