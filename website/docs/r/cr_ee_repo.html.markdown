---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_ee_repo"
sidebar_current: "docs-alicloud-resource-cr-ee-repo"
description: |-
  Provides a Alicloud Container Registry Enterprise Edition Repository resource.
---

# alicloud_cr_ee_repo

Provides a Container Registry Enterprise Edition Repository resource.

For information about Container Registry Enterprise Edition Repository and how to use it, see [What is Repository](https://www.alibabacloud.com/help/en/acr/developer-reference/api-cr-2018-12-01-createrepository)

-> **NOTE:** Available since v1.86.0.

-> **NOTE:** You need to set your registry password in Container Registry Enterprise Edition console before use this resource.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cr_ee_repo&exampleId=c886163f-ce34-3536-e9ab-da6647a4420bc9c8213d&activeTab=example&spm=docs.r.cr_ee_repo.0.c886163fce&intl_lang=EN_US" target="_blank">
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

resource "alicloud_cr_ee_instance" "default" {
  payment_type   = "Subscription"
  period         = 1
  renew_period   = 0
  renewal_status = "ManualRenewal"
  instance_type  = "Advanced"
  instance_name  = "${var.name}-${random_integer.default.result}"
}

resource "alicloud_cr_ee_namespace" "default" {
  instance_id        = alicloud_cr_ee_instance.default.id
  name               = "${var.name}-${random_integer.default.result}"
  auto_create        = false
  default_visibility = "PUBLIC"
}

resource "alicloud_cr_ee_repo" "example" {
  instance_id = alicloud_cr_ee_instance.default.id
  namespace   = alicloud_cr_ee_namespace.default.name
  name        = "${var.name}-${random_integer.default.result}"
  repo_type   = "PUBLIC"
  summary     = "this is summary of my new repo"
  detail      = "this is a public repo"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cr_ee_repo&spm=docs.r.cr_ee_repo.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the Container Registry Enterprise Edition instance.
* `namespace` - (Required, ForceNew) The name of the namespace to which the image repository belongs.
* `name` - (Required, ForceNew) The name of the image repository.
* `repo_type` - (Required) The type of the repository. Valid values:
  - `PUBLIC`: The repository is a public repository.
  - `PRIVATE`: The repository is a private repository.
* `summary` - (Required) The summary about the repository.
* `detail` - (Optional) The description of the repository.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Repository. It formats as `<instance_id>:<namespace>:<name>`.
* `repo_id` - The ID of the repository.

## Import

Container Registry Enterprise Edition Repository can be imported using the id, e.g.

```shell
$ terraform import alicloud_cr_ee_repo.example <instance_id>:<namespace>:<name>
```
