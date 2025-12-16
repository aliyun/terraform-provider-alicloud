---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_namespace"
sidebar_current: "docs-alicloud-resource-container-registry"
description: |-
  Provides a Alicloud resource to manage Container Registry namespaces.
---

# alicloud_cr_namespace

This resource will help you to manager Container Registry namespaces, see [What is Namespace](https://www.alibabacloud.com/help/en/acr/developer-reference/api-cr-2018-12-01-createnamespace).

-> **NOTE:** Available since v1.34.0.

-> **NOTE:** You need to set your registry password in Container Registry console before use this resource.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cr_namespace&exampleId=2f2d0fca-8564-48f6-df56-08cb39ad22b9389ee75a&activeTab=example&spm=docs.r.cr_namespace.0.2f2d0fca85&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "terraform-example"
}

resource "random_integer" "default" {
  min = 10000000
  max = 99999999
}

resource "alicloud_cr_namespace" "example" {
  name               = "${var.name}-${random_integer.default.result}"
  auto_create        = false
  default_visibility = "PUBLIC"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cr_namespace&spm=docs.r.cr_namespace.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) Name of Container Registry namespace.
* `auto_create` - (Required) Boolean, when it set to true, repositories are automatically created when pushing new images. If it set to false, you create repository for images before pushing.
* `default_visibility` - (Required) `PUBLIC` or `PRIVATE`, default repository visibility in this namespace.

## Attributes Reference

The following attributes are exported:

* `id` - The id of Container Registry namespace. The value is same as its name.

## Import

Container Registry namespace can be imported using the namespace, e.g.

```shell
$ terraform import alicloud_cr_namespace.default my-namespace
```
