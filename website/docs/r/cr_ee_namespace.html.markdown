---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_ee_namespace"
sidebar_current: "docs-alicloud-resource-cr-ee-namespace"
description: |-
  Provides a Alicloud resource to manage Container Registry Enterprise Edition namespaces.
---

# alicloud_cr_ee_namespace

This resource will help you to manager Container Registry Enterprise Edition namespaces.

For information about Container Registry Enterprise Edition namespaces and how to use it, see [Create a Namespace](https://www.alibabacloud.com/help/en/acr/developer-reference/api-cr-2018-12-01-createnamespace)

-> **NOTE:** Available since v1.86.0.

-> **NOTE:** You need to set your registry password in Container Registry Enterprise Edition console before use this resource.

## Example Usage
<div class="oics-button" style="float: right;margin: 0 0 -40px 0;">
  <a href="https://api.aliyun.com/api-tools/terraform?resource=alicloud_cr_ee_namespace&exampleId=84a53ef2-51c5-ee91-41dd-fc139db1c54ac4bf77fd&activeTab=example&spm=docs.r.cr_ee_namespace.0.84a53ef251" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; margin: 32px auto; max-width: 100%;">
  </a>
</div>

Basic Usage

```terraform
variable "name" {
  default = "terraform-example-name"
}
resource "alicloud_cr_ee_instance" "example" {
  payment_type   = "Subscription"
  period         = 1
  renew_period   = 0
  renewal_status = "ManualRenewal"
  instance_type  = "Advanced"
  instance_name  = var.name
}

resource "alicloud_cr_ee_namespace" "example" {
  instance_id        = alicloud_cr_ee_instance.example.id
  name               = var.name
  auto_create        = false
  default_visibility = "PUBLIC"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) ID of Container Registry Enterprise Edition instance.
* `name` - (Required, ForceNew) Name of Container Registry Enterprise Edition namespace. It can contain 2 to 30 characters.
* `auto_create` - (Required) Boolean, when it set to true, repositories are automatically created when pushing new images. If it set to false, you create repository for images before pushing.
* `default_visibility` - (Required) `PUBLIC` or `PRIVATE`, default repository visibility in this namespace.

## Attributes Reference

The following attributes are exported:

* `id` - ID of Container Registry Enterprise Edition namespace. The value is in format `{instance_id}:{namespace}` .

## Import

Container Registry Enterprise Edition namespace can be imported using the `{instance_id}:{namespace}`, e.g.

```shell
$ terraform import alicloud_cr_ee_namespace.default cri-xxx:my-namespace
```
