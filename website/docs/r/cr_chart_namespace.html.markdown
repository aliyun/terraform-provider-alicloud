---
subcategory: "Container Registry (CR)"
layout: "alicloud"
page_title: "Alicloud: alicloud_cr_chart_namespace"
sidebar_current: "docs-alicloud-resource-cr-chart-namespace"
description: |-
  Provides a Alicloud CR Chart Namespace resource.
---

# alicloud_cr_chart_namespace

Provides a CR Chart Namespace resource.

For information about CR Chart Namespace and how to use it, see [What is Chart Namespace](https://www.alibabacloud.com/help/en/acr/developer-reference/api-cr-2018-12-01-createchartnamespace).

-> **NOTE:** Available since v1.149.0.

## Example Usage

Basic Usage

<div style="display: block;margin-bottom: 40px;"><div class="oics-button" style="float: right;position: absolute;margin-bottom: 10px;">
  <a href="https://api.aliyun.com/terraform?resource=alicloud_cr_chart_namespace&exampleId=63bbe816-f1f3-80cd-6ac0-a0a3ab8c44076d2b361b&activeTab=example&spm=docs.r.cr_chart_namespace.0.63bbe816f1&intl_lang=EN_US" target="_blank">
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
```

## Argument Reference

The following arguments are supported:

* `auto_create_repo` - (Optional) Specifies whether to automatically create repositories in the namespace. Valid values:
  * `true` - automatically creates repositories in the namespace.
  * `false` - does not automatically create repositories in the namespace.
* `default_repo_type` - (Optional) DefaultRepoType. Valid values: `PRIVATE`, `PUBLIC`.
* `instance_id` - (Required, ForceNew) The ID of the Container Registry instance.
* `namespace_name` - (Required, ForceNew) The name of the namespace that you want to create.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Chart Namespace. The value formats as `<instance_id>:<namespace_name>`.

## Import

CR Chart Namespace can be imported using the id, e.g.

```shell
$ terraform import alicloud_cr_chart_namespace.example <instance_id>:<namespace_name>
```