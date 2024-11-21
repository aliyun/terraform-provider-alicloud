---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_repo"
sidebar_current: "docs-alicloud-resource-container-registry"
description: |-
  Provides a Alicloud resource to manage Container Registry repositories.
---

# alicloud_cr_repo

This resource will help you to manager Container Registry repositories, see [What is Repository](https://www.alibabacloud.com/help/en/acr/developer-reference/api-cr-2018-12-01-createrepository).

-> **NOTE:** Available since v1.35.0.

-> **NOTE:** You need to set your registry password in Container Registry console before use this resource.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cr_repo&exampleId=d66267ae-7391-500c-c95b-3407bb9bb12afbd2568c&activeTab=example&spm=docs.r.cr_repo.0.d66267ae73&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}
resource "alicloud_cr_namespace" "example" {
  name               = var.name
  auto_create        = false
  default_visibility = "PUBLIC"
}

resource "alicloud_cr_repo" "example" {
  namespace = alicloud_cr_namespace.example.name
  name      = var.name
  summary   = "this is summary of my new repo"
  repo_type = "PRIVATE"
  detail    = "this is a private repo"
}
```

## Argument Reference

The following arguments are supported:

* `namespace` - (Required, ForceNew) Name of container registry namespace where repository is located.
* `name` - (Required, ForceNew) Name of container registry repository.
* `summary` - (Required) The repository general information. It can contain 1 to 80 characters.
* `repo_type` - (Required) `PUBLIC` or `PRIVATE`, repo's visibility.
* `detail` - (Optional) The repository specific information. MarkDown format is supported, and the length limit is 2000.

## Attributes Reference

The following attributes are exported:

* `id` - The id of Container Registry repository. The value is in format `namespace/repository`.
* `domain_list` - (Optional) The repository domain list.
  * `public` - Domain of public endpoint.
  * `internal` - Domain of internal endpoint, only in some regions.
  * `vpc` - Domain of vpc endpoint.

## Import

Container Registry repository can be imported using the `namespace/repository`, e.g.

```shell
$ terraform import alicloud_cr_repo.default `my-namespace/my-repo`
```
