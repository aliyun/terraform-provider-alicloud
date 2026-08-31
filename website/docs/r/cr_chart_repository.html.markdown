---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_chart_repository"
sidebar_current: "docs-alicloud-resource-cr-chart-repository"
description: |-
  Provides a Alicloud CR Chart Repository resource.
---

# alicloud_cr_chart_repository

Provides a CR Chart Repository resource.

For information about CR Chart Repository and how to use it, see [What is Chart Repository](https://www.alibabacloud.com/help/en/acr/developer-reference/api-cr-2018-12-01-createchartrepository).

-> **NOTE:** Available since v1.149.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cr_chart_repository&exampleId=c0f976aa-4477-cd56-6da3-09ba814a248fe75a219d&activeTab=example&spm=docs.r.cr_chart_repository.0.c0f976aa44&intl_lang=EN_US" target="_blank">
    <img alt="Open in AliCloud" src="https://img.alicdn.com/imgextra/i1/O1CN01hjjqXv1uYUlY56FyX_!!6000000006049-55-tps-254-36.svg" style="max-height: 44px; max-width: 100%;">
  </a>
</div></div>

```terraform
variable "name" {
  default = "tf-example"
}

resource "random_integer" "default" {
  min = 100000
  max = 999999
}

resource "alicloud_cr_ee_instance" "example" {
  payment_type   = "Subscription"
  period         = 1
  renew_period   = 0
  renewal_status = "ManualRenewal"
  instance_type  = "Advanced"
  instance_name  = "${var.name}-${random_integer.default.result}"
}

resource "alicloud_cr_chart_namespace" "example" {
  instance_id    = alicloud_cr_ee_instance.example.id
  namespace_name = "${var.name}-${random_integer.default.result}"
}

resource "alicloud_cr_chart_repository" "example" {
  repo_namespace_name = alicloud_cr_chart_namespace.example.namespace_name
  instance_id         = alicloud_cr_chart_namespace.example.instance_id
  repo_name           = "${var.name}-${random_integer.default.result}"
}
```

📚 Need more examples? [VIEW MORE EXAMPLES](https://api.aliyun.com/terraform?activeTab=sample&source=Sample&sourcePath=OfficialSample:alicloud_cr_chart_repository&spm=docs.r.cr_chart_repository.example&intl_lang=EN_US)

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the Container Registry instance.
* `repo_name` - (Required, ForceNew) The name of the repository that you want to create.
* `repo_namespace_name` - (Required, ForceNew) The namespace to which the repository belongs.
* `repo_type` - (Optional) The default repository type. Valid values: `PUBLIC`,`PRIVATE`.
* `summary` - (Optional) The summary about the repository.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Chart Repository. The value formats as `<instance_id>:<repo_namespace_name>:<repo_name>`.

## Import

CR Chart Repository can be imported using the id, e.g.

```shell
$ terraform import alicloud_cr_chart_repository.example <instance_id>:<repo_namespace_name>:<repo_name>
```